---
translation_set_id: practice-input-calculator
path: practice/input-calculator
locale: es
group: practice
group_order: 4
order: 1
title: Proyecto: Una calculadora que valida la entrada
summary: Asocia entradas, comprobaciones de límites, funciones y códigos de salida.
---

## Metas y acción

Ingrese la cantidad y el precio unitario y calcule el total. Introduzca dos números enteros separados por un espacio o salto de línea. Este ejemplo solo acepta cantidades del 1 al 1000 y precios unitarios del 0 al 100000, por lo que se calcula dentro del rango de multiplicación i32.

Guárdelo en `main.wave`, ejecute `wavec run main.wave` y luego ingrese `3 1200`. Los resultados arrojados por el programa son los siguientes. La visibilidad de los caracteres ingresados ​​en la terminal es independiente de la salida del programa.

<!-- wave-example: input-calculator -->
```wave
fun total(quantity: i32, price: i32) -> i32 {
    return quantity * price;
}

fun main() -> i32 {
    var quantity: i32 = 0;
    var price: i32 = 0;
    input("{} {}", quantity, price);
    if (quantity < 1 || quantity > 1000 || price < 0 || price > 100000) {
        println("out of range");
        return 1;
    }

    println("total={}", total(quantity, price));
    return 0;
}
```

Resultado de la ejecución:

```text
total=3600
```

## También verifique si hay fallas

Si ingresa `0 1200`, espera `out of range` y el código de salida 1. La falla del análisis numérico en `input` y la verificación del alcance de trabajo del programa son dos pasos diferentes. El token/tipo no numérico excede el rango/antes de la entrada requerida EOF es una falla de entrada. Esta entrada integrada no es una interfaz que devuelve un error y requiere una nueva entrada. Si necesita un analizador recuperable, configure usted mismo el proceso de verificación leyendo los bytes con [io](/docs/es/stdlib/files-io).

## Ejercicios extendidos y comentarios.

Tome la tasa de descuento como tercera entrada y verifique si es de 0 a 100. Para evitar grandes multiplicaciones intermedias, debe ampliar el rango de cálculo con i64 y verificar el rango al reducir el resultado. Si simplemente amplía el último tipo de resultado, es posible que ya se hayan realizado cálculos intermedios en tipos estrechos.

[Consola I/O](/docs/es/language/console-io-and-formatting) · [Siguiente: Leer archivos](/docs/es/practice/file-reader)

## ¿Por qué decidir primero el alcance del cálculo?

Los insumos más grandes son la cantidad 1000 y el precio unitario 100000. Multiplicar los dos valores da 100000000, por lo que cae dentro del rango i32. Esta verificación de rango es necesaria para que el resultado de la multiplicación de la función total pueda usarse tal cual.

main, que recibe entradas, es responsable de las entradas y mensajes de error, y total solo es responsable de los cálculos. Incluso si luego cambia a leer órdenes de un archivo, aún puede usar la función de cálculo.

## Cálculo completo del descuento

Puedes multiplicar el total por 100 sumando el porcentaje de descuento. Convierta para realizar cálculos intermedios como i64 y luego aplique la tasa de descuento. Dado que se descarta la parte fraccionaria de la división de números enteros, este ejemplo trunca la cantidad después del descuento a un número entero.

<!-- wave-example: book-calculator-discount -->
```wave
fun discounted_total(quantity: i32, price: i32, discount: i32) -> i64 {
    var subtotal: i64 = (quantity as i64) * (price as i64);
    var remaining: i64 = 100 - (discount as i64);
    return subtotal * remaining / 100;
}

fun main() -> i32 {
    var quantity: i32 = 0;
    var price: i32 = 0;
    var discount: i32 = 0;

    input("{} {} {}", quantity, price, discount);

    if (quantity < 1 || quantity > 1000) {
        println("invalid quantity");
        return 1;
    }

    if (price < 0 || price > 100000) {
        println("invalid price");
        return 2;
    }

    if (discount < 0 || discount > 100) {
        println("invalid discount");
        return 3;
    }

    println("total={}", discounted_total(quantity, price, discount));
    return 0;
}
```

Resultado de la ejecución:

```text
total=3240
```

El resultado anterior es cuando se ingresa `3 1200 10`. La salida es 3240, que se resta el 10% del total original de 3600.

|entrada|resultado esperado|camino para comprobar|
| --- | --- | --- |
| `3 1200 0` | `total=3600` |sin descuento|
| `3 1200 100` | `total=0` |Descuento total|
| `3 1200 101` | `invalid discount` |La tasa de descuento excede el rango|
| `1000 100000 0` | `total=100000000` |entrada máxima|
| `0 1200 10` | `invalid quantity` |Cantidad por debajo del rango|

## próximo ejercicio

Intente cambiar la función para redondear cantidades decimales. En este programa, que acepta sólo cantidades positivas, sumar 50 antes de dividir por 100 redondeará a números enteros. Cuando ingresa 99 por 1 y una tasa de descuento de 50, puede comparar el resultado de corte de 49 y el resultado de redondeo de 50.
