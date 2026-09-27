---
translation_set_id: modules-ffi
path: language/modules-imports-and-ffi
locale: es
group: language
group_order: 2
order: 11
title: 11. Módulos y código genérico.
summary: Obtenga información sobre cómo obtener nombres públicos y argumentos de tipo explícitos.
---

## Razones para dividir archivos

A medida que el programa crece, es más fácil encontrarlo agrupando funciones relacionadas en lugar de colocar todas las funciones en main.wave. Los límites de los módulos determinan qué nombres están expuestos a otro código. Los genéricos son una herramienta que reutiliza la misma tarea con diferentes tipos, independientemente de la separación de archivos.

En este capítulo, creamos un programa de dos archivos y aprendemos los alias de los módulos, la selección import y funciones y estructuras genéricas.

## programa de dos archivos

Cree helpers.wave y main.wave en el mismo directorio.

helpers.wave Todos:

```wave
pub fun double(value: i32) -> i32 {
    return value * 2;
}
```

main.wave Todos:

<!-- wave-example: book-module-two-files -->
```wave
import("./helpers")::{
    double
};

fun main() {
    var result: i32 = double(21);

    println("{}", result);
}
```

Resultado de la ejecución:

```text
42
```

Ejecute `wavec run main.wave` en la terminal. helpers.wave tampoco se ejecuta por separado. La fuente requerida se conecta a través de import.

El pub delante de la función en helpers indica que otros módulos pueden importarlo. Las funciones auxiliares que no necesitan exponerse al mundo exterior no necesitan hacerse públicas. Incluso si cambia la implementación interna de un módulo, puede reducir los cambios en el código que utiliza manteniendo el contrato de la función pública.

## Línea de base para rutas relativas

`./helpers` es relativo al directorio del archivo fuente que creó la oración import. Al ejecutar el programa, sepárelo del directorio de trabajo donde está escrito el archivo I/O. Los pasos para encontrar el archivo import y los pasos para encontrar input.txt mientras se ejecuta son diferentes.

Para import local, se puede omitir la extensión `.wave`. Si dividiste el directorio, escribe la ruta `./module` según la ubicación. No confunda las rutas relativas locales con rutas que obtienen los nombres de las dependencias del paquete.

## Seleccione import y alias

La selección import hace que solo los nombres públicos deseados se utilicen directamente en el archivo actual. Si hay un conflicto de nombres o desea revelar a qué módulo pertenece la función, utilice un alias.

<!-- wave-example: book-module-alias -->
```wave
import("std::string::len" as strings);

fun main() {
    var length: i32 = strings::len("Wave");

    println("{}", length);
}
```

Resultado de la ejecución:

```text
4
```

strings es el alias del módulo definido en este archivo. `strings::len` usa el nombre del módulo. El punto para el acceso al campo y `::` para la separación de módulos son notaciones diferentes.

No utilice la opción import y el alias import juntos en una oración. Sea cual sea su estilo, utilícelo de forma coherente en todo el archivo para que el origen del nombre sea fácilmente legible.

## Bibliotecas y paquetes estándar

La ruta `std::` apunta a la biblioteca estándar. El usuario publica API y import los módulos requeridos. No todas las funciones de la biblioteca estándar se colocan automáticamente en el espacio de nombres actual.

Las rutas de los paquetes externos comienzan desde el nombre del paquete. La ubicación del paquete la proporcionan las opciones del compilador o el administrador de paquetes. Primero aprenda los límites con los módulos locales y luego aprenda cómo administrar las dependencias en [Vex Cómo utilizar](/docs/es/whale/vex-package-manager).

## Utilice la misma función para cada tipo

Las siguientes funciones devuelven su entrada palabra por palabra: Para i32 y str, use el parámetro de tipo T para evitar escribir el mismo código dos veces.

<!-- wave-example: book-generic-identity -->
```wave
fun identity<T>(value: T) -> T {
    return value;
}

fun main() {
    var number: i32 = identity<i32>(42);
    var text: str = identity<str>("Wave");

    println("{} {}", number, text);
}
```

Resultado de la ejecución:

```text
42 Wave
```

T es un lugar para ingresar el tipo, no el valor entero pasado durante la ejecución. Llámelo especificando el argumento de tipo como `<i32>`. Las funciones genéricas de usuario normales no omiten argumentos de tipo.

identity<str> no duplica bytes de cadena al asignarlos nuevamente. Devuelve el valor tal cual. La gramática de los genéricos no cambia las reglas de copia y propiedad de los datos.

## Operaciones requeridas por el cuerpo genérico.

El hecho de que exista un parámetro de tipo no significa que todas las operaciones se puedan utilizar en todos los tipos. minimum a continuación debe usarse como un tipo real que se puede comparar.

<!-- wave-example: book-generic-minimum -->
```wave
fun minimum<T>(left: T, right: T) -> T {
    if (left < right) {
        return left;
    }

    return right;
}

fun main() {
    var small: i32 = minimum<i32>(7, 4);
    var wide: i64 = minimum<i64>(100, 20);

    println("{} {}", small, wide);
}
```

Resultado de la ejecución:

```text
4 20
```

Si cambia el argumento de tipo, `<` y el retorno utilizado en el texto debe ser del tipo correspondiente. Al leer errores genéricos, verifique tanto la combinación de tipos llamada como la operación requerida por el cuerpo de la función.

## estructura genérica

Creemos Pair, que vincula dos valores diferentes.

<!-- wave-example: book-generic-pair -->
```wave
struct Pair<A, B> {
    first: A;
    second: B;
}

fun main() {
    var item: Pair<i32, str> = Pair<i32, str> {
        first: 7,
        second: "seven"
    };

    println("{} {}", item.first, item.second);
}
```

Resultado de la ejecución:

```text
7 seven
```

Pair<i32, str> y Pair<i64, str> son tipos específicos diferentes. El orden de los argumentos de tipo también es significativo. Lea el código de declaración y generación para ver dónde se determinan los tipos de first y second.

## Nombre y contrato de API publicados

Cuando publica una función, especifica no solo el nombre sino también la unidad de entrada, el valor de retorno, el error y la propiedad. Por ejemplo, el bucle que escribirá la persona que llama depende de si read lee la longitud máxima o la longitud exacta.

pub es el alcance público entre los módulos Wave. Esto es diferente de export (c), que exporta símbolos externos para que otros idiomas los llamen. Puedes ver un ejemplo completo que vincula los dos idiomas en [Ver FFI](/docs/es/language/modules-imports-and-ffi).

## Ejercicio y solución completa

Cree una función pública square en math.wave y llámela como alias en main.wave para imprimir la potencia de 3 y 5.

math.wave:

```wave
pub fun square(value: i32) -> i32 {
    return value * value;
}
```

main.wave:

<!-- wave-example: book-module-solution -->
```wave
import("./math" as math);

fun main() {
    println("{} {}", math::square(3), math::square(5));
}
```

Resultado de la ejecución:

```text
9 25
```

Si el error es que la función ha desaparecido, verifique primero la ruta import y pub. Si hay un conflicto de nombres, verifique si la llamada tiene un alias. Si se trata de un error de tipo, verifique la entrada de la función y el tipo del argumento pasado. No intente resolver diferentes problemas con una corrección de ruta.


## C Función de importación

```wave
extern(c) fun puts(text: ptr<i8>) -> i32;
```

El nombre ABI puede ir seguido del nombre del símbolo real como una cadena.

```wave
extern(c, "native_symbol") fun local_name(value: i32) -> i32;
```

## Wave Función de exportación

```wave
export(c) fun wave_add(left: i32, right: i32) -> i32 {
    return left + right;
}
```

`extern` y `export` se pueden utilizar como funciones y bloques únicos. La función exportada debe tener la firma específica ABI y por lo tanto no puede ser genérica.

## Propiedad de condición objetivo

Las propiedades de condición de destino se pueden adjuntar a elementos de nivel superior.

```wave
#[target(os="linux", arch="x86_64")]
extern(c) fun platform_call(value: i32) -> i32;
```

Las claves de condición son `arch`, `os`, `env`, `abi` y las propiedades se aplican al siguiente elemento de nivel superior.

## Conéctese con la función C que usted mismo escribió

Esta práctica de laboratorio es para el entorno nativo con el compilador C. Concatena una función de número entero sin asignación de biblioteca ni procesamiento de cadenas.

`native.c`:

```c
#include <stdint.h>
int32_t native_double(int32_t value) { return value * 2; }
```

`main.wave`:

```wave
extern(c) fun native_double(value: i32) -> i32;

fun main() {
    println("{}", native_double(21));
}
```

Ejecute Linux/macOS desde una terminal en el mismo directorio de trabajo.

```shell
cc -c native.c -o native.o
wavec build main.wave native.o -o ffi-example
./ffi-example
```

El resultado esperado es `42`. En el shell de desarrollador de MSVC en Windows, cree object con `cl /c native.c /Fonative.obj` y conéctese a `wavec build main.wave native.obj -o ffi-example.exe`. Las arquitecturas de origen y destino de object deben ser las mismas. Los ejemplos utilizan sólo valores pequeños. Para pasar un valor grande a la función C, el rango de multiplicación de la página C también debe garantizarse por separado.

Las rutas de archivos locales comienzan con `./`.
