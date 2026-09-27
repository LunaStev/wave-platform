---
translation_set_id: whale-numeric-operations
path: whale/numeric-operations
locale: es
group: whale
group_order: 1
order: 7
title: Operaciones numéricas
summary: Describe operaciones de números enteros wrap, checked, cambios, errores de conversión de tipos y resultados de punto flotante.
---

## representación entera

Un entero de bits N tiene bits de valor N. Los enteros sin signo varían de 0 a 2^N − 1, y los enteros con signo varían de −2^(N−1) a 2^(N−1) − 1. El signedness de la operación determina la interpretación de la cadena de bits.

La siguiente tabla explica los resultados del cálculo. IR El ejemplo de código utiliza la representación de la impresora actual.

## Suma, resta, multiplicación.

Los enteros básicos add·sub·mul contienen los bits N inferiores del resultado. overflow no causa trap. La operación checked devuelve el mismo resultado wrap, así como Bool, que indica si el resultado matemático está fuera del rango del signed o unsigned correspondiente.

|operación|Wrap Resultado| Checked overflow |
| --- | --- | --- |
| u8: 255 + 1 | 0 | true |
| i8: 127 + 1 | −128 | true |
| u8: 0 − 1 | 255 | true |
| i8: 12 × 3 | 36 | false |

Las interfaces de los lenguajes que interrumpen la ejecución en overflow deben usar un `trap_if` explícito para el overflow resultado de la operación checked. La operación predeterminada no aplica implícitamente la política overflow del idioma de origen.

### Wrap y IR expresando un cheque explícito

Los siguientes módulos están configurados como Rust builder y son la salida de impresora actual que pasó el verificador. El analizador de texto y el backend de ejecución aún no están disponibles.

```text
module {
  format_version 2
  semantics_version 1
  target "x86_64-whale-linux"
  datalayout { ptr=64, endian=little }

  declare @f0 "add_u8": whale () -> u8, linkage internal
  declare @f1 "require_no_overflow": whale () -> u8, linkage internal

  fn @add_u8() -> u8, id @f0 {
  entry:
    %v0: u8 = const u8 255
    %v1: u8 = const u8 1
    %v2: u8 = add u8 %v0, %v1
    ret u8 %v2
  }

  fn @require_no_overflow() -> u8, id @f1 {
  entry:
    %v3: u8 = const u8 255
    %v4: u8 = const u8 1
    %v5: tuple<u8, bool> = uadd_chk u8 %v3, %v4
    %v6: u8 = extract %v5, 0
    %v7: bool = extract %v5, 1
    trap_if bool %v7, reason="integer overflow"
    ret u8 %v6
  }

}
```

Según el contrato aritmético, `add_u8` ajusta 256 al resultado 0 de 8 bits. `require_no_overflow` extrae el resultado ajustado con `extract ..., 0` y el indicador de desbordamiento Bool con `extract ..., 1`. Si el indicador es true, `trap_if` detiene la ejecución antes de la devolución. Estos son los resultados especificados, no los resultados de un intérprete implementado.

## División y resto

La división de números enteros o el resto por cero provoca una trampa. Dividiendo el valor mínimo con signo por −1 se ajusta al valor mínimo. El resto en este caso es cero.

|operación|resultado|
| --- | --- |
| i8: −128 / −1 | −128 |
| i8: −128 % −1 | 0 |
|Dividir por número entero 0| trap |
|resto para el entero 0| trap |

## turno

Al cambiar el valor de bit N, la cadena de bits de count se interpreta como unsigned y se utiliza el resto dividido por N. No genera trap solo porque count está fuera del rango de 0 a N−1.

En valores de 8 bits, count 0·8·16 son todos desplazamientos de 0 bits. La cadena de bits `11111111` de 8 bits count se desplaza 7 bits. Esto es lo mismo incluso si esta cadena de bits representa signed −1. unsigned El análisis se realiza antes que el resto de los cálculos.

Los idiomas de origen que rechazan count negativo o excesivo deben realizar una verificación explícita antes del cambio.

## conversión de tipo

|conversión|significado|
| --- | --- |
| Zero extension |Aumente el ancho llenando bits de orden superior con 0|
| Sign extension |Aumente el ancho duplicando el bit del signo.|
|corte de bits|Mantenga solo los bits de orden inferior correspondientes al ancho de destino|
|Beat reinterpretación|Interpretar la misma cadena de bits como diferentes tipos|
|Conversión numérica sin pérdidas|Si no se puede expresar conservando el valor numérico, trap|

Por ejemplo, convertir `11111111` de 8 bits en zero extension de 16 bits se convierte en `0000000011111111` y sign extension se convierte en `1111111111111111`. Incluso si los bits de entrada son los mismos, son operaciones diferentes.

La conversión de flotante a int se trunca hacia cero y luego verifica el rango de enteros. NaN y el infinito provocan una trampa. Al convertir a i8, 127,9 se convierte en 127, mientras que 128,0 queda atrapado. Bool se convierte a un número entero 0 o 1. i1 con signo no puede representar 1, por lo que no puede ser el destino de esta conversión.

Convertir una dirección a un número entero no restaura el acceso del puntero a ese número entero. Consulte [validez del puntero](memory-model).

## aritmética de coma flotante

Los valores de punto flotante tienen la cadena de bits exacta f16·f32·f64. La operación redondea al valor más cercano en el ancho declarado y, si está exactamente en el medio, utiliza ties-to-even, que selecciona un valor incluso con los bits menos significativos de los dígitos significativos.

La operación predeterminada es fast-math, implícita FMA, que no permite el manejo de cero forzado de valores pequeños. La multiplicación seguida de la suma conserva sus respectivos pasos de redondeo y el backend no debe combinarlos implícitamente en una sola operación.

El resultado de una operación numérica puede ser NaN o infinito. Los NaN de operaciones matemáticas y cambios de ancho se normalizan a una cantidad fija de NaN silencioso por ancho. Guarde y copie para conservar los bits NaN originales. Por lo tanto, el comportamiento es diferente al pasar la carga útil de NaN a la memoria sin operaciones aritméticas y al calcularla.

No expone indicadores de estado de punto flotante. Aunque la aritmética básica de punto flotante permite NaN·resultados infinitos, la conversión float→int aplica las reglas trap anteriores.


### Almacenar constantes exactas

Utilice `FloatBits` variant o cualquier cadena de bits hexadecimales del ancho exacto. La igualdad de los valores almacenados se compara con una cadena de bits que contiene 0 negativo y NaN payload. El verificador rechaza si los anchos de tipo IR y payload son diferentes.

```rust
use ir::{FloatBits, ModuleBuilder, Target, Type};
fn main() {
    let bits = FloatBits::parse(32, "0xffc01234").unwrap();
    assert_eq!(bits, FloatBits::F32(0xffc01234));
    let target = Target::X86_64WhaleLinux;
    let mut module = ModuleBuilder::new(target.name(), target.data_layout());
    let mut function = module.begin_function("payload", vec![], Type::F32);
    let value = function.const_float_bits(Type::F32, bits);
    function.ret(Some(value));
    function.finish();
    let module = module.finish();
    ir::verify_module(&module).unwrap();
    assert!(ir::print_module(&module).contains("const f32 0xffc01234"));
    println!("{}", bits);
}
```

```text
0xffc01234
```

Una cadena de bits f16, f32 o f64 comienza con `0x`, seguida exactamente de 4, 8 o 16 dígitos hexadecimales, respectivamente. `0x80000000` representa f32 cero negativo; `0x7f800000` representa infinito positivo. `const_float` convierte un valor de host f64 numéricamente; utilice `const_float_bits` para conservar los bits originales. La representación exacta del almacenamiento no implica un backend de ejecución de punto flotante completo. La aritmética en tiempo de compilación todavía utiliza intermediarios del host f64, por lo que aún no se ha implementado el contrato de redondeo completo para cada ancho declarado.
