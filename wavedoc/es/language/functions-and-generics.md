---
translation_set_id: functions
path: language/functions-and-generics
locale: es
group: language
group_order: 2
order: 5
title: 5. Diseño y composición de funciones.
summary: Aprenda parámetros, valores de retorno, valores predeterminados y paso por valor.
---

## A partir de código repetitivo

Las funciones son una herramienta para reducir la sintaxis, pero también son una herramienta para demarcar tareas. Separar lo que toma como entrada, lo que calcula y los resultados que devuelve le permite comprender su programa en partes más pequeñas.

En este capítulo, comenzamos con un programa que escribe cálculos de descuento varias veces. Cada ejemplo se completa en main.wave y se ejecuta como `wavec run main.wave`. Aún no estamos dividiendo los archivos.

<!-- wave-example: book-function-before -->
```wave
fun main() {
    var first_price: i32 = 2000;
    var first_discount: i32 = first_price * 10 / 100;
    var first_total: i32 = first_price - first_discount;
    var second_price: i32 = 5000;
    var second_discount: i32 = second_price * 10 / 100;
    var second_total: i32 = second_price - second_discount;
    println("{} {}", first_total, second_total);
}
```

Resultado de la ejecución:

```text
1800 4500
```

Los dos cálculos difieren sólo en el precio y tienen la misma estructura. Al cambiar las reglas de descuento, deberás editar ambos lugares. Si cambia solo un lado, obtendrá resultados diferentes para productos que requieren la misma política.

## Determinar la entrada y la salida.

Mueva los cálculos redundantes a funciones. El valor modificado se recibe como entrada price y se devuelve el precio calculado.

<!-- wave-example: book-function-extract -->
```wave
fun discounted(price: i32) -> i32 {
    var discount: i32 = price * 10 / 100;
    return price - discount;
}

fun main() {
    println("{} {}", discounted(2000), discounted(5000));
}
```

Resultado de la ejecución:

```text
1800 4500
```

El nombre de la función es discounted. `price: i32` entre paréntesis está la declaración del parámetro y `-> i32` es el tipo de resultado. La variable local discount en el cuerpo se usa solo dentro de esta función.

`discounted(2000)` es una expresión que llama a una función. El 2000 entre paréntesis es el argumento que realmente pasas. Dado que el valor devuelto por la función se convierte en el resultado de esta expresión de llamada, se puede utilizar directamente como argumento para println.

|terminología|código|significado|
| --- | --- | --- |
|parámetro| price |Nombre de entrada especificado al declarar la función|
|factor| 2000 |Valor pasado al llamar|
|tipo de retorno| i32 |Tipo de valor producido por la expresión de llamada|
|declaración de devolución| return price - discount |Pasa el resultado y finaliza esta llamada.|

## Seguir la orden de convocatoria y ejecución.

Simplemente escribir una declaración de función en el código fuente no hace que su cuerpo se ejecute inmediatamente. Se ejecuta cuando se alcanza el punto llamado en main.

<!-- wave-example: book-function-trace -->
```wave
fun calculate(value: i32) -> i32 {
    println("inside: {}", value);
    return value * 2;
}

fun main() {
    println("before");
    var result: i32 = calculate(7);
    println("after: {}", result);
}
```

Resultado de la ejecución:

```text
before
inside: 7
after: 14
```

El orden de progreso es el primer resultado de main, el texto principal de calculate y el último resultado de main. Cuando se ejecuta `return`, la llamada a calculate finaliza con el resultado 14 y se completa la inicialización de result de main.

Si se mezclan varias llamadas a funciones en una expresión y el orden de los efectos secundarios es importante, divida las llamadas en declaraciones separadas. Los ejemplos de este capítulo también almacenan los resultados de las llamadas que deben rastrearse en variables locales.

## múltiples parámetros

Si también recibe una tasa de descuento como entrada, puede calcular varias pólizas con la misma función.

<!-- wave-example: book-function-parameters -->
```wave
fun discounted(price: i32, percent: i32) -> i32 {
    return price - price * percent / 100;
}

fun main() {
    var standard: i32 = discounted(2000, 10);
    var special: i32 = discounted(2000, 25);
    println("standard={}", standard);
    println("special={}", special);
}
```

Resultado de la ejecución:

```text
standard=1800
special=1500
```

El orden de los argumentos debe coincidir con la declaración. Si ambos parámetros son i32, es difícil distinguir su significado mediante la verificación de tipos únicamente, incluso si se cambia el orden. Defina claramente el nombre de la función y el nombre del parámetro, y escriba la ubicación de la llamada de una manera fácil de leer.

Esta función asume una pequeña cantidad y una tasa de descuento válida. No maneja precios negativos, ratios superiores a 100 ni desbordamientos a mitad de multiplicación. Al crear una función, debes describir no solo el cuerpo sino también las condiciones de entrada. El programa de finalización posterior separa los pasos de inspección.

## argumento predeterminado

Puede proporcionar valores utilizados con frecuencia como predeterminados.

<!-- wave-example: book-function-default -->
```wave
fun discounted(price: i32, percent: i32 = 10) -> i32 {
    return price - price * percent / 100;
}

fun main() {
    println("{}", discounted(2000));
    println("{}", discounted(2000, 25));
}
```

Resultado de la ejecución:

```text
1800
1500
```

La primera llamada omite el segundo argumento y usa 10. La segunda llamada usa el 25 especificado. Los valores predeterminados se dejan para los parámetros opcionales al final. No se utiliza como gramática escribir un espacio en blanco para omitir sólo el primer argumento.

Cambiar el valor predeterminado cambia el comportamiento de omitir llamadas. Los valores predeterminados de las funciones públicas también son parte del comportamiento en el que confía. Esta es la razón por la que las llamadas con argumentos especificados y las llamadas con argumentos omitidos se prueban por separado.

## significa pasar por valor

Pasar un valor entero separa el valor que recibe la función del almacenamiento de variables de la persona que llama. Calcular un resultado no cambia automáticamente las variables de la persona que llama.

<!-- wave-example: book-function-value -->
```wave
fun next(value: i32) -> i32 {
    return value + 1;
}

fun main() {
    var count: i32 = 4;
    var later: i32 = next(count);
    println("count={} later={}", count, later);
    count = next(count);
    println("count={}", count);
}
```

Resultado de la ejecución:

```text
count=4 later=5
count=5
```

La primera llamada dice count e inicializa later. count sigue siendo 4. Después de la segunda llamada, el resultado se asigna a count, por lo que pasa a ser 5. El diseño de retorno por valor hace evidente para la persona que llama dónde han cambiado los datos.

Si desea cambiar el espacio de almacenamiento original dentro de una función, puede pasar un puntero. Esto se trata en [capítulo puntero](/docs/es/language/explicit-memory-type-model). Incluso si pasa un puntero, debe distinguir entre el valor del puntero en sí y el espacio de almacenamiento para su dirección.

## No dejes caminos que no regresan

Una función que devuelve un valor debe proporcionar resultados en todas las rutas requeridas. No deja el camino final así:

```wave
// 의도적으로 잘못된 함수 조각

fun positive(value: i32) -> i32 {
    if (value > 0) {
        return value;
    }
}
```

Si value es menor o igual a 0, no hay ningún valor establecido para devolver. Una vez que haya establecido las reglas previstas, debe escribir todas sus rutas.

<!-- wave-example: book-function-paths -->
```wave
fun positive_or_zero(value: i32) -> i32 {
    if (value > 0) {
        return value;
    }

    return 0;
}

fun main() {
    println("{} {} {}", positive_or_zero(-2), positive_or_zero(0), positive_or_zero(8));
}
```

Resultado de la ejecución:

```text
0 0 8
```

La llamada que se ejecuta por primera vez return no continúa con return a continuación. El último return se alcanza solo cuando la condición es falsa. Verificamos la regla para tres casos: positivo, 0 y negativo.

## Función sin resultado

Si solo realiza operaciones como imprimir, puede omitir el tipo de devolución. Una función sin resultado también se puede finalizar antes de tiempo con `return;`.

<!-- wave-example: book-function-void -->
```wave
fun show_positive(value: i32) {
    if (value <= 0) {
        return;
    }

    println("positive={}", value);
}

fun main() {
    show_positive(-3);
    show_positive(6);
}
```

Resultado de la ejecución:

```text
positive=6
```

La primera llamada regresa sin imprimir nada. Las impresiones de la segunda llamada: “Sin resultado” y “No regresa al punto de llamada” son diferentes. Las funciones que no regresan, como las funciones que finalizan un proceso, se clasifican por su tipo de retorno `!`.

## Componer un programa completo con múltiples funciones.

Ahora dividimos la validación de entrada, el cálculo y la salida en diferentes funciones. Como hemos restringido el rango de precios, la multiplicación intermedia en este ejemplo está dentro del rango i32.

<!-- wave-example: book-function-program -->
```wave
fun valid_order(price: i32, quantity: i32, percent: i32) -> bool {
    return price >= 0 && price <= 100000
        && quantity >= 1 && quantity <= 100
        && percent >= 0 && percent <= 100;
}

fun discounted_unit(price: i32, percent: i32) -> i32 {
    return price - price * percent / 100;
}

fun order_total(price: i32, quantity: i32, percent: i32) -> i32 {
    var unit: i32 = discounted_unit(price, percent);
    return unit * quantity;
}

fun show_order(price: i32, quantity: i32, percent: i32) {
    if (!valid_order(price, quantity, percent)) {
        println("invalid order");
        return;
    }

    var total: i32 = order_total(price, quantity, percent);
    println("total={}", total);
}

fun main() {
    show_order(2000, 3, 10);
    show_order(2000, 0, 10);
    show_order(2000, 3, 120);
}
```

Resultado de la ejecución:

```text
total=5400
invalid order
invalid order
```

Responda una pregunta para cada función. valid_order determina si se permite la entrada, discounted_unit determina cuánto es un descuento, order_total determina cuánto es el precio total y show_order determina qué mostrar.

Las funciones más pequeñas no son necesariamente mejores. Darle un nombre a cada expresión puede dificultar el seguimiento. Reglas separadas cuando sean significativas para ser reutilizadas en otro lugar o cuando haya reglas que deban explicarse y verificarse de forma independiente.

## problemas de practica

1. Escriba maximum, que devuelve el mayor de dos números enteros.
2. Escribe una función clamp que devuelva un valor entre dos límites. En la siguiente solución, low <= high se establece como condición de llamada.
3. Cree una función que agregue impuestos a un valor y combínela con la función de suma del pedido. Primero determine el rango y los puntos de corte de números enteros.

### Solución: función con límite

<!-- wave-example: book-function-solution -->
```wave
fun maximum(left: i32, right: i32) -> i32 {
    if (left > right) {
        return left;
    }

    return right;
}

fun clamp(value: i32, low: i32, high: i32) -> i32 {
    if (value < low) {
        return low;
    }
    if (value > high) {
        return high;
    }

    return value;
}

fun main() {
    println("max={}", maximum(7, 4));
    println("{} {} {}", clamp(-3, 0, 10), clamp(6, 0, 10), clamp(20, 0, 10));
}
```

Resultado de la ejecución:

```text
max=7
0 6 10
```

clamp tiene tres rutas: menor que el rango, dentro del rango y mayor que el rango. Verifique agregando manualmente también los valores límite 0 y 10. Si desea postular hasta low > high, deberá decidir cómo expresar el fracaso. [Capítulo de manejo de errores](/docs/es/language/errors) aborda este problema mediante estructuras de resultados.


## llamada recursiva

Una función puede llamarse a sí misma. La recursividad requiere una condición de salida en la que no se realizan más llamadas y un proceso en el que cada llamada se acerca a esa condición.

<!-- wave-example: book-ref-recursion -->
```wave
fun factorial(value: i32) -> i32 {
    if (value <= 1) {
        return 1;
    }

    return value * factorial(value - 1);
}

fun main() {
    println("{}", factorial(5));
}
```

Resultado de la ejecución:

```text
120
```

El cálculo de 5 conduce a `5 * factorial(4)`, que devuelve 1 cuando se alcanza. El valor de retorno se pasa secuencialmente a la llamada anterior, lo que da como resultado 120. Esta función es un ejemplo para explicar los números enteros positivos pequeños. Con entradas grandes, debe considerar el rango de resultados y la profundidad de la llamada. Puede evitar el problema del aumento de la profundidad de las llamadas escribiendo la misma operación en un bucle.

## Errores comunes

|fenómeno|comprobar|
| --- | --- |
|Error con argumentos insuficientes o demasiados|Número de parámetros y valores predeterminados opcionales|
|El tipo de devolución no coincide|return Tipo de expresión y declaración de función|
|No regresar de un camino específico|¿Regresa incluso si la condición es falsa?|
|Falta un argumento de tipo genérico|Después del nombre de la función `<Type>`|
|El original se cambia después de llamar a una función de puntero.|¿Es una función que solo lee el valor o una función que lo modifica?|

Las funciones exportadas externamente, como `export(c)`, utilizan firmas específicas. Las funciones genéricas en sí no se pueden exportar con una convención de llamada externa. `ptr<T>` y `array<T, N>` son los tipos de memoria integrados del lenguaje y los distinguen de las declaraciones de estructura genéricas del usuario.

[aprendizaje de funciones](/docs/es/language/functions-and-generics) · [Módulos de aprendizaje y genéricos.](/docs/es/language/modules-imports-and-ffi) · [FFI](/docs/es/language/modules-imports-and-ffi)
