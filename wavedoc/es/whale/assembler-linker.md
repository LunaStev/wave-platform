---
translation_set_id: whale-assembler-linker
path: whale/assembler-linker
locale: es
group: whale
group_order: 1
order: 11
title: Ensamblado y enlace estático
summary: Describe la codificación de operandos, la ubicación de las secciones, la vinculación de símbolos y los puntos de entrada de ejecución.
---

## Ensamblajes y objetos

El ensamblador Whale convierte la instrucción AMD64 en bytes de código de máquina e información de reubicación. El objeto ELF64 contiene esta información y secciones/símbolos. El ensamblaje lo realiza la propia implementación de Whale y no requiere un ensamblador externo.

Los objetos reubicables pueden tener referencias cuya dirección final aún no se conoce. El proceso de resolución de esta dirección es un enlace. Un ensamblaje exitoso no debe inferir que se han resuelto todos los símbolos externos o que se ha creado un ejecutable.

## Montar la función

Guarde el siguiente código como `answer.asm`.

```asm
section .text
global answer

answer:
    mov eax, 42
    ret
```

```sh
whale asm --amd64 answer.asm -o answer.o
```

El contenido de `.text` es `b8 2a 00 00 00 c3` y `mov eax, 42` va seguido de `ret`. El objeto ELF64 expone `answer`. Es una función invocable sin código de inicio del proceso y no es un archivo ejecutable. Las instrucciones de este ejemplo pueden ser procesadas por el ensamblador actual.

## Construyendo objetos con Rust API

Aquí hay un ejemplo completo de cómo escribir los mismos bytes de función en la caja `object`. Especifica destinos de salida y símbolos globales.

```rust
use object::{ObjectFile, ObjectSymbol, SectionKind, SymbolBinding, SymbolVisibility, Target};

fn main() {
    let mut object = ObjectFile::with_target(Target::X86_64WhaleLinux.object_target());
    let text = object.add_section(".text", SectionKind::Text, 16);
    object.sections[text].data = vec![0xb8, 0x2a, 0x00, 0x00, 0x00, 0xc3];
    object.symbols.push(ObjectSymbol {
        name: "answer".into(),
        section_index: Some(text),
        value: 0,
        size: 6,
        binding: SymbolBinding::Global,
        visibility: SymbolVisibility::Default,
    });
    let elf = object.write().unwrap();
    assert_eq!(&elf[..7], b"\x7fELF\x02\x01\x01");
    assert_eq!(u16::from_le_bytes([elf[18], elf[19]]), 62);
    std::fs::write("answer.o", elf).unwrap();
}
```

Assertion comprueba la clase ELF, orden de bytes, identificador machine. `value: 0` es `.text` dentro de offset y `size: 6` es el tamaño de bytes del símbolo. Las referencias o rangos de secciones no válidos son errores de serialización. Rechaza otros machine u órdenes de bytes sin marcarlos como AMD64.

## Interpretar los símbolos de dos objetos.

Actualmente, la caja `linker` proporciona interpretación de símbolos. El siguiente ejemplo ejecutable define un `helper` local en dos objetos y luego verifica si se produce un error si el mismo nombre se publica dos veces.

```rust
use linker::core::symbol_table::{SymbolKey, SymbolTable};
use object::{
    ObjectFile, ObjectFormat, ObjectSymbol, SectionKind, SymbolBinding, SymbolVisibility,
};

fn input(binding: SymbolBinding) -> ObjectFile {
    let mut object = ObjectFile::new(ObjectFormat::ELF64);
    let text = object.add_section(".text", SectionKind::Text, 1);
    object.sections[text].data = vec![0xc3];
    object.symbols.push(ObjectSymbol {
        name: "helper".into(),
        section_index: Some(text),
        value: 0,
        size: 1,
        binding,
        visibility: SymbolVisibility::Default,
    });
    object
}

fn main() {
    let mut symbols = SymbolTable::new();
    symbols
        .resolve(&[input(SymbolBinding::Local), input(SymbolBinding::Local)])
        .unwrap();
    for object_index in 0..2 {
        let key = SymbolKey::Local {
            object_index,
            name: "helper".into(),
        };
        assert_eq!(symbols.symbols[&key].object_index, Some(object_index));
    }
    let error = symbols
        .resolve(&[input(SymbolBinding::Global), input(SymbolBinding::Global)])
        .unwrap_err();
    assert_eq!(error, "Duplicate global symbol: helper");
    println!("{error}");
}
```

Salida:

```text
Duplicate global symbol: helper
```

Las dos definiciones locales tienen claves separadas con `object_index` diferentes. Las dos definiciones globales entran en conflicto. Este ejemplo interpreta directamente los símbolos del objeto que ha construido en la memoria. `.o` No se realiza la lectura de archivos, la aplicación de reubicación ni la salida de archivos ejecutables. De las políticas de selección de definición completa que se describen a continuación, las prioridades weak, etc. aún no se han implementado.

## Literales y operandos de memoria.

Los literales conservan su ancho y señalización hasta que se examina el rango de codificación de la instrucción real. Los valores que están fuera de rango deberían generar un error en lugar de truncarse silenciosamente.

Se requiere una notación de tamaño explícita cuando el ancho de la memoria no se puede determinar a partir de otra información en el comando. Por ejemplo, un operando de registro puede determinar el ancho, pero esto puede resultar ambiguo si sólo hay memoria y un valor inmediato. El ensamblador no debe adivinar anchos ambiguos al azar.

El operando de memoria que consta de un símbolo en AMD64 es básicamente RIP-relative. Seleccione el método de direccionamiento como explícito rel/abs. La incógnita escape en un literal es un error.

## Secciones y alineación

|Contenido de la sección|Comportamiento de alineación|
| --- | --- |
|código|NOP Insertar comando|
|datos inicializados|Insertar 0 bytes|
| BSS |Aumente el tamaño de la memoria lógica sin agregar el archivo payload|

Existe una distinción entre tamaño de archivo y tamaño de memoria. BSS reserva memoria, pero no requiere que se almacene el mismo tamaño de 0 bytes en el archivo objeto. Las secciones personalizadas tienen propiedades y los símbolos mantienen información vinculante y de tipo.

### Payload Reserva BSS sin asignación

Guarde lo siguiente como `buffer.asm`. La lógica BSS reserva 1 TiB, pero no asigna ni escribe 1 TiB durante el ensamblaje.

```asm
section .bss
global buffer
buffer:
    resb 1099511627776
buffer_end:

section .text
    ret
```

```sh
whale asm --amd64 buffer.asm -o buffer.o
```

El valor dentro de la sección de `buffer` es 0 y el valor de `buffer_end` es 1099511627776. El encabezado `.bss` es `SHT_NOBITS` con ese tamaño y no hay ningún archivo payload. Luego, cuando regrese a `.bss`, continúe usando la lógica offset. Las directivas de datos con ceros también aumentan el tamaño lógico. Se rechazan el valor inicial distinto de cero, la reubicación que se aplicará dentro de BSS y los comandos dentro de BSS.

Se puede utilizar la misma distinción para objetos y enlazadores API.

```rust
use linker::core::layout::Layout;
use object::{ObjectFile, ObjectFormat, SectionKind};

fn main() {
    let mut object = ObjectFile::new(ObjectFormat::ELF64);
    let text = object.add_section(".text", SectionKind::Text, 16);
    object.sections[text].data = vec![0xc3];
    let bss = object.add_section(".bss", SectionKind::Bss, 16);
    object.sections[bss].zero_fill = 1 << 40;
    assert!(object.sections[bss].data.is_empty());

    let elf = object.write_with_limit(4096).unwrap();
    assert!(elf.len() < 4096);
    let layout = Layout::compute(&[object], 0x1000).unwrap();
    assert_eq!(layout.file_size, 1);
    assert_eq!(layout.memory_size, (1 << 40) + 16);
    assert_eq!(layout.sections[bss].memory_address, 0x1010);
    assert_eq!(layout.sections[bss].file_size, 0);
}
```

```text
section  memory address  file bytes       memory bytes
.text    0x1000          1                1
.bss     0x1010          0                1099511627776
```

`Section::zero_fill` es el número de BSS bytes adicionales que no están almacenados en el archivo. El tamaño de memoria probado es `data.len() + zero_fill`. También se acepta el BSS `data` con relleno de ceros existente, pero su asignación se puede evitar con un vector `data` vacío. Las secciones distintas de BSS deben ser `zero_fill == 0`. Un espacio de almacenamiento físico 0 no implica el estado de inicialización de IR.

`Layout::compute` devuelve `Result`, que registra la correspondencia entre el objeto/sección de entrada, la alineación, el archivo offset, la dirección de memoria y ambos tamaños de todas las secciones. Se verifica la aritmética de dirección y alineación, se conserva el orden de entrada y la alineación 0 del objeto se trata como sin restricciones (alineación 1). BSS no mueve el cursor del archivo. Esta es la implementación payload, con el encabezado ejecutable, load segment con derechos de acceso, aplicar la reubicación es una tarea separada.

ELF writer rechaza overflow y truncamiento al reducir el ancho del campo. No se admiten números de sección extendidos y el número total de encabezados, incluida la tabla creada y las secciones reorganizadas, debe ser inferior a `0xff00`. El límite predeterminado para la salida serializada es 256 MiB. Puede especificar un límite de bytes que incluya la tabla padding· con `ObjectFile::write_with_limit` o `write_elf_with_limit`, y el límite no incluye el tamaño de memoria BSS que no está almacenado en un archivo. El tamaño de salida se verifica antes de asignar el vector de bytes final.

## Identificación de símbolos

Las funciones y variables utilizan diferentes identificadores dentro de IR. Las conexiones externas utilizan `link_name` especificado por la interfaz. Whale no cambia automáticamente el nombre de uno de los símbolos públicos en conflicto, manteniendo su nombre indicado.

La tabla de declaración IR registra referencias `FunctionId` escritas y valores de función `link_name` explícita. Las API de ensamblador/objeto/enlazador que aparecen a continuación siguen siendo interfaces separadas: la emisión nativa IR aún no transmite esas identidades al enlace final. IR la verificación de llamadas por sí sola no establece esta propiedad de un extremo a otro.

Por lo tanto, es posible que tanto una función interna como una variable tengan el nombre `item`, pero sería un error exponer ambas con el mismo nombre externo. La separación del espacio de nombres interno no separa automáticamente el espacio de nombres externo.

El alcance de un símbolo local de objeto es su objeto de entrada. Los símbolos globales participan en la interpretación entre objetos. Los conflictos de función/datos confirmados son errores. El símbolo NOTYPE mantiene la compatibilidad con entradas que no proporcionan un tipo más específico. El mero hecho de que no tenga tipo no significa que sea una función o un dato.

## seleccionar definición

|definición o referencia|resultado|
| --- | --- |
|Strong y strong|error de definición duplicada|
|Strong y weak|Strong Seleccionar definición|
|Weak y weak|Seleccione la primera definición en orden de entrada|
|Ver strong sin resolver|error de enlace|
|Ver sin resolver weak|Error no admitido en el perfil estático inicial|

Si hay varias definiciones de weak, el orden en que se ingresen afectará los resultados. Debe utilizar el orden de entrada pasado de forma coherente para enlaces deterministas.

## Salida ejecutable estática

El perfil estático native genera ELF ET_EXEC especificando el punto de entrada. El punto de entrada no se deduce del nombre de la función `main`. Tampoco inserta automáticamente el código de inicio que llama a esa función.

No elimina automáticamente secciones, no fusiona códigos idénticos ni elimina símbolos. La colocación de archivos requiere una contabilidad separada de los bytes que realmente guarda y la memoria que reserva en tiempo de ejecución.

La ruta completa de creación del ejecutable estático para CLI aún no está disponible. `whale asm` crea un objeto reubicable y `whale object` envuelve bytes sin formato en un objeto. Consulte [Descripción general de la cadena de herramientas](overview) para conocer la disponibilidad, ABI y [AMD64 Objetivo](amd64-target) para conocer los requisitos.


## Serialización de registros seleccionable Wave

Linux Whale compilaciones para x86_64 ahora pueden usar la implementación Wave para registros fijos ELF64 encabezado·sección·símbolo·RELA. También se proporciona una implementación predeterminada Rust. La selección del objetivo de salida, la verificación de objetos, la ubicación, la interpretación de símbolos y la asignación de búfer están a cargo de Rust. Al seleccionar Wave no se agregan arquitecturas compatibles ni un vinculador completo.

Instale Rust, LLVM 21 bibliotecas de desarrollo, el vinculador C y `ar`, luego cree la ruta seleccionada desde el repositorio Whale.

```sh
git clone https://github.com/wavefnd/Wave.git /tmp/whale-wave-bootstrap
git -C /tmp/whale-wave-bootstrap checkout --detach 8a465e30aeea4b817d925cdd0e8d08c1bb029c9a
python3 tools/build_wave_elf.py --wave-source /tmp/whale-wave-bootstrap --out-dir /tmp/whale-wave-elf
WHALE_WAVE_ELF_DIR=/tmp/whale-wave-elf cargo build --locked --all-features
WHALE_WAVE_ELF_DIR=/tmp/whale-wave-elf cargo test --locked --workspace --all-features
```

El script busca la corrección revision, rechaza el cambio a tracked, construye el compilador Wave, crea un objeto Wave con LLVM y lo agrupa con archive para enlaces estáticos. `WHALE_WAVE_ELF_DIR` debería tener esto archive. Solicitar explícitamente un archive no válido o un host no compatible es un error de compilación. Las compilaciones normales sin variables especificadas no requieren los compiladores `--all-features` ni Wave. El compilador Wave tampoco es necesario en tiempo de ejecución para el ejecutable Whale vinculado. La ruta para realizar la compilación cruzada de Whale con bootstrap aún no es compatible.

Por ejemplo, guarde lo siguiente como `return.asm` y ensamblelo en el binario generado:

```asm
section .text
global entry
entry:
    ret
```

```sh
target/debug/whale asm --amd64 return.asm -o return.o
```

```text
ELF class: ELF64
machine: AMD64 (62)
.text bytes: c3
```

El registro ABI lleva el tipo, el puntero de campo u64, el número de campos, el puntero de salida y la capacidad. El búfer es propiedad de la persona que llama Rust. No transfiere la propiedad de la asignación ni la expresión Rust enum/String/Vec a través de límites. La rutina Wave comprueba el recuento, la capacidad y el ancho del campo antes de escribir. El éxito devuelve el estado 0, el tipo/puntero/capacidad no válido devuelve 1 y el campo overflow devuelve 2. El puntero debe apuntar a buffers que estén activos y del tamaño correcto y que no se superpongan entre sí. El puntero sin formato C por sí solo no puede probar esta condición; wrapper lo satisface. La prueba compara toda la reubicación ELF, incluida BSS y signed addend, con la ruta Rust. Esta es una implementación parcial de Wave que utiliza Wave/LLVM bootstrap y no es completamente autohospedada.
