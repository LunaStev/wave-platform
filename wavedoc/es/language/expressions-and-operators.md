---
translation_set_id: expressions
path: language/expressions-and-operators
locale: es
group: language
group_order: 2
order: 3
title: 3. Aritmética, comparaciones y conversiones
summary: Aprenda el orden de cálculo, la división de enteros, las operaciones bit a bit y cast.
---

## Ver los resultados de los cálculos y los tipos de cálculos juntos

Una expresión es un código que calcula un valor. Los nombres de variables, los literales, las llamadas a funciones y las expresiones que conectan múltiples valores con operadores son todas expresiones. En matemáticas, aunque parezca la misma ecuación, el resultado será diferente dependiendo de si es un número entero o real y de cuántos bits contiene.

Comenzando con cálculos simples, este capítulo presenta paréntesis, división, operaciones lógicas, operaciones bit a bit y conversiones. Cada ejemplo es un archivo main.wave completo que puede ejecutar con `wavec run main.wave`.

## Rango encerrado entre paréntesis

<!-- wave-example: book-precedence -->
```wave
fun main() {
    var first: i32 = 2 + 3 * 4;
    var second: i32 = (2 + 3) * 4;

    println("{} {}", first, second);
}
```

Resultado de la ejecución:

```text
14 20
```

La multiplicación se evalúa antes que la suma, por lo que la primera expresión es 2+12. En la segunda ecuación, tomamos la suma entre paréntesis, que es 5, y luego la multiplicamos por 4. El objetivo es no utilizar menos paréntesis. Se recomienda utilizarlo para que el lector pueda comprender fácilmente el alcance del cálculo.

Incluso cuando se utiliza el mismo operador varias veces, la dirección del encadenamiento es importante. `20 - 5 - 3` es `(20 - 5) - 3`, que es 12. `20 - (5 - 3)` es 18. La secuencia completa exacta está en [referencia del operador](/docs/es/language/expressions-and-operators).

## División de números enteros y resto

La división de dos números enteros no produce un resultado de punto flotante fraccionario. Utilice la división para el cociente y el operador resto para el resto.

<!-- wave-example: book-division -->
```wave
fun main() {
    var items: i32 = 17;
    var box_size: i32 = 5;
    var full_boxes: i32 = items / box_size;
    var remaining: i32 = items % box_size;

    println("boxes={} remaining={}", full_boxes, remaining);
    println("negative quotient={}", -17 / 5);
}
```

Resultado de la ejecución:

```text
boxes=3 remaining=2
negative quotient=-3
```

Dividir 17 elementos en grupos de 5 produce 3 grupos completos y 2 elementos restantes. La división con signo se trunca hacia cero, por lo que -17/5 es -3. Esto difiere de redondear hacia abajo hacia el infinito negativo.

No se puede dividir entre 0. signed El valor obtenido al dividir el valor mínimo entre -1 no entra en el mismo tipo. Las funciones que toman estas entradas deben verificarse antes de dividir o usar checked matemáticas API.

## El punto de conversión cambia el resultado.

Las dos expresiones siguientes se almacenan en la variable f64, pero el proceso de cálculo es diferente.

<!-- wave-example: book-division-conversion -->
```wave
fun main() {
    var already_divided: f64 = (7 / 2) as f64;
    var floating_division: f64 = (7 as f64) / (2 as f64);

    if (already_divided == 3.0) {
        println("integer division lost the fraction");
    }

    if (floating_division == 3.5) {
        println("floating division kept the fraction");
    }
}
```

Resultado de la ejecución:

```text
integer division lost the fraction
floating division kept the fraction
```

La primera expresión realiza una división de enteros para obtener 3 y luego lo convierte a f64. El segundo convierte los operandos a f64 antes de realizar la división en punto flotante. Elegir un tipo más amplio para la variable final no puede recuperar la información perdida anteriormente.

Los valores de coma flotante son aproximaciones. Es posible que dos resultados que parezcan el mismo valor decimal no sean adecuados para una comparación de igualdad exacta. Elija una tolerancia que se ajuste a las unidades y escala del problema; un épsilon fijo no es adecuado para todos los cálculos.

## Comparación realizada bool

`<`, `<=`, `>`, `>=`, `==`, `!=` verifique la relación. Un signo igual `=` es una asignación y dos signos iguales `==` son comparaciones iguales.

<!-- wave-example: book-comparisons -->
```wave
fun main() {
    var age: i32 = 20;
    var minimum: i32 = 18;
    var eligible: bool = age >= minimum;

    if (eligible) {
        println("eligible");
    }

    if (age != minimum) {
        println("not exactly the boundary");
    }
}
```

Resultado de la ejecución:

```text
eligible
not exactly the boundary
```

“Al menos 18” incluye 18; “mayor de 18” lo excluye. Pruebe 17, 18 y 19 para verificar este límite. Las comparaciones de tipos mixtos dependen del signo y el ancho, por lo que convertir ambos operandos al tipo deseado puede hacer que la comparación sea más clara.

## Operaciones lógicas y evaluación de cortocircuitos.

`&&` comprueba si ambos son verdaderos, `||` comprueba si más de uno es verdadero y `!` cambia verdadero/falso. Esta operación no siempre ejecuta la expresión del lado derecho.

<!-- wave-example: book-short-circuit -->
```wave
fun report() -> bool {
    println("right side evaluated");
    return true;
}

fun main() {
    var enabled: bool = false;

    if (enabled && report()) {
        println("both true");
    }

    if (!enabled || report()) {
        println("at least one true");
    }
}
```

Resultado de la ejecución:

```text
at least one true
```

En la primera condición, enabled es falso, por lo que no es necesario mirar hacia la derecha. En el segundo, !enabled es verdadero, por lo que el lado derecho tampoco es necesario. Por lo tanto, el resultado de report nunca aparece.

Puedes usar esto para verificar el denominador antes de la división. El fragmento del cuerpo de la función `if (divisor != 0 && value / divisor > 2) { ... }` no realiza división cuando el denominador es 0. Sin embargo, esta prueba por sí sola no resuelve otros límites como signed valor mínimo/-1.

`&&` tiene prioridad sobre `||`. En políticas complejas, indique su intención entre paréntesis, como `(member && active) || admin`.

## Estrechamiento o ampliación de números enteros

`as` es un elenco explícito. Los bits de orden superior que se descartan al reducir un número entero no se pueden recuperar ampliándolo nuevamente.

<!-- wave-example: book-cast-chain -->
```wave
fun main() {
    var original: i32 = 300;
    var small: u8 = original as u8;
    var widened: i32 = small as i32;

    println("{} -> {} -> {}", original, small, widened);

    var negative: i8 = -1;
    var signed_value: i32 = negative as i32;
    var same_bits: u8 = negative as u8;

    println("{} {}", signed_value, same_bits);
}
```

Resultado de la ejecución:

```text
300 -> 44 -> 44
-1 255
```

Los 8 bits inferiores de 300 son 44. Si amplía -1 al tipo signed, el signo se ampliará para mantener -1. Si interpreta los mismos 8 bits como unsigned, es 255.

Las transformaciones que buscan preservar valores dentro del rango y las transformaciones que intentan manipular bits de almacenamiento tienen propósitos diferentes. Si tiene información del usuario, primero verifique si es el rango de destino y conviértalo. La presencia de cast no garantiza que el valor esté en un rango seguro.

## Reemplazar con bool

Al convertir un número entero a bool se obtiene false para cero y true en caso contrario. Esto no se trunca al bit más bajo: 2 también se convierte a true. Para valores de punto flotante, solo +0,0 y -0,0 se convierten a false; todos los demás valores, incluidos NaN e infinitos, se convierten a true.

<!-- wave-example: book-bool-conversion -->
```wave
fun main() {
    var zero: i32 = 0;
    var two: i32 = 2;

    if (!(zero as bool)) {
        println("zero is false");
    }

    if (two as bool) {
        println("two is true");
    }
}
```

Resultado de la ejecución:

```text
zero is false
two is true
```

No se admiten conversiones de puntero a bool. Compare el puntero con null explícitamente, por ejemplo con `pointer != null`. Si una dirección que no es null es segura de leer es una cuestión aparte.

## operaciones de bits

`&`, `|`, `^`, `~` cubren cada bit del número entero. Se puede utilizar para expresar permisos o funciones en bits. A continuación, 1 es permiso de lectura y 2 es permiso de escritura.

<!-- wave-example: book-bit-flags -->
```wave
const READ: u32 = 1;
const WRITE: u32 = 2;

fun main() {
    var permissions: u32 = READ | WRITE;

    if ((permissions & WRITE) != 0) {
        println("write enabled");
    }

    permissions = permissions & ~WRITE;
    println("remaining={}", permissions);
}
```

Resultado de la ejecución:

```text
write enabled
remaining=1
```

Agregue un bit con OR y verifique la presencia de un bit específico con AND. Cree una máscara con solo los bits relevantes siendo 0 con `~WRITE` y bórrela. Las operaciones bit a bit `&`·`|` son operadores diferentes de la evaluación de cortocircuito `&&`·`||` de bool.

## Ancho y número de turnos

Un desplazamiento a la izquierda mueve los bits hacia la izquierda y descarta los bits altos fuera del ancho del operando. Un desplazamiento a la derecha utiliza extensión de signo para valores con signo y extensión cero para valores sin signo. El resultado siempre tiene el tipo de operando izquierdo.

<!-- wave-example: book-shifts -->
```wave
fun main() {
    var bits: u8 = 129;
    var shifted: u8 = bits << 1;
    var negative: i32 = -8;
    var positive: u32 = 8;

    println("left={}", shifted);
    println("right={} {}", negative >> 1, positive >> 1);
}
```

Resultado de la ejecución:

```text
left=2
right=-4 4
```

Al desplazar el valor u8 129 a la izquierda en uno, se descarta su bit más alto y se deja 2. El valor de recuento de desplazamiento original debe ser no negativo y menor que el ancho de bits del operando izquierdo. Para u8, los recuentos válidos son del 0 al 7. Un recuento constante no válido es un error en tiempo de compilación; un recuento de tiempo de ejecución no válido provoca una trampa.

## Convertir valores de punto flotante a números enteros

La conversión de punto flotante a entero primero se trunca hacia cero y luego verifica el rango de enteros de destino. Los resultados NaN, infinito y fuera de rango no son válidos. Las conversiones constantes no válidas producen un error en tiempo de compilación; Las conversiones de tiempo de ejecución no válidas provocan una trampa.

Una trampa no devuelve un valor de error de la función. Para errores de conversión recuperables, diseñe una interfaz que verifique el rango antes de realizar la conversión. Para obtener los bits almacenados de un valor de punto flotante, utilice las funciones de conversión de bits en `std::math::float` en lugar de una conversión numérica.

## Ejercicio y solución

Escriba una función que divida 137 wones a cambio en 50 unidades de wones y el resto, y cambie solo números enteros en el rango de 0 a 255 a u8. Este ejemplo demuestra el paso de validación como una función separada que indica fuera de rango como -1.

<!-- wave-example: book-expression-solution -->
```wave
fun checked_byte(value: i32) -> i32 {
    if (value < 0 || value > 255) {
        return -1;
    }

    var small: u8 = value as u8;
    return small as i32;
}

fun main() {
    var money: i32 = 137;

    println("coins={} remainder={}", money / 50, money % 50);
    println("{} {} {}", checked_byte(0), checked_byte(255), checked_byte(256));
}
```

Resultado de la ejecución:

```text
coins=2 remainder=37
0 255 -1
```

La razón por la que -1 se puede utilizar como indicador de error es porque el rango de éxito es de 0 a 255. Si cualquier número entero puede ser un valor de éxito, se necesita una representación de resultado diferente. Continuamos con este diseño en el capítulo sobre manejo de errores más adelante.


## prioridad

La precedencia de los operadores es la siguiente, comenzando desde el más alto:

1. Expresiones básicas y acceso a postfijos: llamadas a funciones, acceso a campos, indexación, postfijos `++`·`--`
2. Operaciones unarias: `!`, `~`, `&`, `deref`, prefijo `++`·`--`, unario `+`·`-`
3. `as` conversión de tipo
4. `*`, `/`, `%`
5. `+`, `-`
6. `<<`, `>>`
7. `<`, `<=`, `>`, `>=`
8. `==`, `!=`
9. Bit `&`
10. Bit `^`
11. Bit `|`
12. `&&`
13. `||`
14. Asignación y asignación compuesta

La asignación de cadenas se combina desde la derecha. Al mezclar diferentes tipos de operadores, utilice paréntesis para expresar claramente el orden de evaluación.

## Materias que pueden ser sustituidas

La asignación, `++` y `--` requieren una expresión que indique una ubicación de almacenamiento, como una variable, un campo, un elemento de matriz o un puntero sin referencia. No se permite escribir en un `const`.

## turno

El resultado de un desplazamiento siempre tiene el tipo de operando izquierdo; el tipo de operando correcto no amplía el cálculo. Los desplazamientos a la izquierda descartan los bits altos más allá de ese ancho. Desplaza hacia la derecha los valores con signo extendidos al signo y los valores sin signo extendidos al cero.

El recuento de turnos debe ser un número entero cuyo valor original satisfaga `0 <= n < LHS bit width`. Se comprueba antes de cualquier truncamiento a un tipo más pequeño. Un recuento constante no válido es un error en tiempo de compilación; un recuento de tiempo de ejecución no válido provoca una trampa.

## Conversiones que involucran valores booleanos y de punto flotante

La conversión de entero a bool produce false para cero y true para cualquier otro valor. La conversión de punto flotante a bool produce false solo para +0,0 y -0,0; NaN e infinito positivo o negativo producen true. Las conversiones de puntero a bool no son compatibles: compárese explícitamente con `pointer != null`.

La conversión de punto flotante a entero se trunca hacia cero y luego verifica el rango de destino. Los resultados NaN, infinito y fuera de rango no son válidos. Las conversiones constantes no válidas son errores en tiempo de compilación; Las conversiones de tiempo de ejecución no válidas provocan una trampa. Una trampa no es una devolución de error recuperable.

`&&` y `||` son evaluaciones de cortocircuito. Los efectos secundarios del operando derecho no ejecutado no ocurren. Puedes comprobar los resultados con valores pequeños en [clase de aritmética](/docs/es/language/expressions-and-operators).

## cambio que falla a propósito

El número de desplazamientos para un valor de 8 bits debe ser de 0 a 7. Los programas siguientes deben rechazarse antes de su ejecución.

<!-- wave-example: reject-shift-count -->
```wave
fun main() {
    var value: u8 = 1;
    var result: u8 = value << 8;
}
```
