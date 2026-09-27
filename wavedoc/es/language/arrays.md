---
translation_set_id: learn-arrays-strings
path: language/arrays
locale: es
group: language
group_order: 2
order: 6
title: 6. Matrices e iteración
summary: Aprenda matrices de tamaño fijo, indexación, recorrido, copia y búsqueda.
---

## Múltiples valores del mismo tipo

Si crea las tres puntuaciones por separado como score1, score2 y score3, tanto las declaraciones como los cálculos deben cambiarse cuando cambia el número. Los arrays agrupan un número determinado de elementos del mismo tipo. Usando bucles, puedes aplicar las mismas reglas a cada elemento.

Este capítulo cubre la creación, indexación, modificación, iteración, búsqueda y agregación de matrices. Las cadenas también admiten la indexación, pero su significado difiere, por lo que se tratan en el siguiente capítulo.

## Introduzca la longitud en tipo

<!-- wave-example: book-array-create -->
```wave
fun main() {
    var scores: array<i32, 3> = [70, 80, 90];

    println("first={}", scores[0]);
    println("second={}", scores[1]);
    println("last={}", scores[2]);
}
```

Resultado de la ejecución:

```text
first=70
second=80
last=90
```

i32 en `array<i32, 3>` es el tipo de elemento y 3 es el número de elementos. Lo que almacenamos son 3 números enteros. No significa que la cantidad de bytes sea 3. La cantidad de elementos en un literal de matriz debe coincidir con la longitud declarada.

Los índices comienzan desde 0. El primer elemento es 0, el último elemento tiene longitud-1. scores[3] es un acceso fuera de alcance, no un tercer elemento.

## cambiar elemento

<!-- wave-example: book-array-update -->
```wave
fun main() {
    var scores: array<i32, 3> = [70, 80, 90];

    scores[1] = 85;
    scores[0] += 5;

    println("{} {} {}", scores[0], scores[1], scores[2]);
}
```

Resultado de la ejecución:

```text
75 85 90
```

Cambie el almacenamiento de elementos específicos sin recrear toda la matriz. Una expresión de índice también puede ser el resultado de un cálculo, pero debe asegurarse de que el valor esté dentro de un rango. Cuando se utiliza una entrada externa como índice, se verifican tanto el número negativo como el límite superior.

## Iterando sobre una matriz

<!-- wave-example: book-array-sum -->
```wave
fun main() {
    var scores: array<i32, 4> = [60, 70, 80, 90];
    var total: i32 = 0;

    for (var index: i32 = 0; index < 4; index += 1) {
        total += scores[index];
    }

    println("total={} average={}", total, total / 4);
}
```

Resultado de la ejecución:

```text
total=300 average=75
```

Cada iteración lee un elemento en un índice diferente. La variable de suma debe inicializarse fuera de la iteración. Inicializarlo a 0 cada vez dentro del cuerpo del bucle producirá resultados incorrectos, como dejar solo el último elemento.

La división de números enteros utilizada para calcular el promedio descarta la parte fraccionaria. Para obtener un promedio de punto flotante, convierta la suma antes de dividir. Para matrices o valores más grandes, asegúrese también de que el tipo de acumulador pueda representar la suma.

## Agregando solo algunos elementos

El filtrado se puede lograr combinando declaraciones condicionales y transversales. Aquí contamos el número de elementos con una puntuación de 80 o superior.

<!-- wave-example: book-array-filter -->
```wave
fun main() {
    var scores: array<i32, 5> = [60, 80, 90, 75, 100];
    var passed: i32 = 0;

    for (var index: i32 = 0; index < 5; index += 1) {
        if (scores[index] >= 80) {
            passed += 1;
        }
    }

    println("passed={}", passed);
}
```

Resultado de la ejecución:

```text
passed=3
```

Los valores de índice y elemento deben estar separados. El examen `index >= 80` compara posiciones, no puntuaciones. Ambos pueden ser i32, por lo que es difícil encontrar este error semántico basándose únicamente en el tipo.

## Encuentra la ubicación del primer partido

Primero, decida cómo mostrar los resultados que no se encontraron. En este ejemplo, los índices válidos son del 0 al 4, por lo que utilizamos -1 como marcador de error.

<!-- wave-example: book-array-find -->
```wave
fun main() {
    var values: array<i32, 5> = [8, 3, 8, 1, 5];
    var target: i32 = 8;
    var found: i32 = -1;

    for (var index: i32 = 0; index < 5; index += 1) {
        if (values[index] == target) {
            found = index;
            break;
        }
    }

    if (found >= 0) {
        println("found at {}", found);
    } else {
        println("not found");
    }
}
```

Resultado de la ejecución:

```text
found at 0
```

La primera posición 0 también es un resultado normal. Si verifica el éxito con `found > 0`, se confundirá con no encontrar el primer elemento. Por la misma razón, juzgar el éxito utilizando bool como cast es incorrecto.

Si elimina break, las coincidencias posteriores sobrescribirán found, dándole la posición de la última coincidencia. Dado que una sola declaración puede cambiar el contrato de una función, la descripción de "búsqueda" también debe escribirse específicamente si está en la primera o última posición.

## Copiar elementos de matriz

Para copiar los valores de una matriz a otro espacio de almacenamiento, puede leerlos y asignarlos elemento por elemento. Incluso si cambia un elemento entero después de copiar, el otro elemento entero no cambia.

<!-- wave-example: book-array-copy -->
```wave
fun main() {
    var original: array<i32, 3> = [1, 2, 3];
    var copied: array<i32, 3>;

    for (var index: i32 = 0; index < 3; index += 1) {
        copied[index] = original[index];
    }

    copied[0] = 99;

    println("original={}", original[0]);
    println("copied={}", copied[0]);
}
```

Resultado de la ejecución:

```text
original=1
copied=99
```

Si los elementos son punteros, copiarlos copia sus direcciones. No duplica la memoria separada a la que apuntan. Esta distinción es importante a la hora de gestionar la propiedad.

## Inicialización y alcance efectivo.

No se supone que se puedan leer todos los elementos de una matriz declarada sin un valor inicial. Si solo se registra un número parcial, el número inicializado real se debe gestionar por separado. Esta es la misma razón por la que la longitud devuelta por la función de lectura de la biblioteca puede ser menor que la capacidad total del búfer.

Debido a que la longitud de la matriz está contenida en el tipo, no crece arbitrariamente durante la ejecución. Las listas de bytes que aumentan de tamaño utilizan almacenamiento dinámico como `Buffer`. Cambiar la longitud de una matriz requiere considerar el tipo, el valor inicial, el límite superior transversal y los cálculos que dependen de esa longitud.

## Ejercicio: máximos y ubicaciones

Encuentre el valor máximo y la posición en la que aparece por primera vez en la matriz `[4, 9, 2, 9, 1]`. Si se trata de una función normal donde todos los elementos pueden ser negativos, el valor máximo no debe inicializarse en 0.

### Solución completa

<!-- wave-example: book-array-solution -->
```wave
fun main() {
    var values: array<i32, 5> = [4, 9, 2, 9, 1];
    var maximum: i32 = values[0];
    var position: i32 = 0;

    for (var index: i32 = 1; index < 5; index += 1) {
        if (values[index] > maximum) {
            maximum = values[index];
            position = index;
        }
    }

    println("max={} first={}", maximum, position);
}
```

Resultado de la ejecución:

```text
max=9 first=1
```

Tome el primer elemento como referencia inicial y compare con el segundo. Desde `>`, la posición no cambia incluso si vuelve a aparecer el mismo valor máximo. Cámbielo a `>=` para que se convierta en la última posición. Cualquier interfaz que pueda tener longitud cero debe procesar una entrada vacía antes de leer el primer elemento.
