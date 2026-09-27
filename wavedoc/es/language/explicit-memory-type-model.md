---
translation_set_id: memory-model
path: language/explicit-memory-type-model
locale: es
group: language
group_order: 2
order: 9
title: 9. Punteros, modificación de valores y vidas útiles.
summary: Aprenda la dirección de, la desreferenciación, la modificación de un valor mediante un puntero y los punteros colgantes.
---

## Distinguir entre valor y ubicación de almacenamiento

El número entero 42 y la dirección donde se almacena el número entero son valores diferentes. El puntero apunta a la ubicación de almacenamiento. Pasar una dirección permite que la función lea o cambie el almacenamiento de la persona que llama.

Este capítulo cubre la toma de direcciones, la desreferenciación, la modificación del valor original, la aritmética de punteros y la vida útil. La asignación dinámica se trata en el siguiente capítulo. Comience con direcciones de variables locales y elementos de matriz.

## Tomar direcciones y desreferenciar

<!-- wave-example: book-pointer-first -->
```wave
fun main() {
    var value: i32 = 42;
    var address: ptr<i32> = &value;

    println("value={}", value);
    println("through pointer={}", deref address);

    deref address = 99;
    println("changed={}", value);
}
```

Resultado de la ejecución:

```text
value=42
through pointer=42
changed=99
```

`&value` obtiene la dirección, `deref address` lee o escribe el valor de esa dirección, y así sucesivamente. En lugar de almacenar 99 en address, se escribió 99 en el número entero señalado por address. La propia variable address continúa apuntando a value.

`ptr<i32>` es un tipo de puntero para acceder al almacenamiento i32. El tipo no registra una longitud ni proporciona una desasignación automática.

## Cambiar el puntero en sí

<!-- wave-example: book-pointer-reassign -->
```wave
fun main() {
    var first: i32 = 10;
    var second: i32 = 20;
    var selected: ptr<i32> = &first;

    selected = &second;
    deref selected = 25;

    println("{} {}", first, second);
}
```

Resultado de la ejecución:

```text
10 25
```

`selected = &second` almacena otra dirección en una variable de puntero. El valor de first no cambia. Luego, si escribe el valor como deref, second cambia. Separar “cambio de dirección” y “cambio de valor a través de dirección” en oraciones separadas reducirá la confusión.

## Hacer que una función cambie el original

<!-- wave-example: book-pointer-increment -->
```wave
fun increment(value: ptr<i32>) {
    deref value = deref value + 1;
}

fun main() {
    var count: i32 = 4;

    increment(&count);
    increment(&count);

    println("{}", count);
}
```

Resultado de la ejecución:

```text
6
```

La función recibe la dirección del recuento, en lugar de su valor, 4. Cambiar ese almacenamiento también cambia el recuento en la persona que llama. La función requiere una dirección i32 válida y grabable. Pasar null viola este requisito.

Cada función define si acepta null. Si no es así, la persona que llama debe proporcionar una dirección válida. Si es así, la función debe incluir una ruta que maneje null.

## Función a procesar null

<!-- wave-example: book-pointer-null -->
```wave
fun try_increment(value: ptr<i32>) -> bool {
    if (value == null) {
        return false;
    }

    deref value = deref value + 1;
    return true;
}

fun main() {
    var count: i32 = 7;

    if (!try_increment(null)) {
        println("no value");
    }

    if (try_increment(&count)) {
        println("count={}", count);
    }
}
```

Resultado de la ejecución:

```text
no value
count=8
```

Una verificación null maneja solo la ausencia de una dirección. Convertir un número arbitrario que no sea null en un puntero no crea memoria válida. La lectura y la escritura también requieren una vida útil válida, un tamaño suficiente, una alineación correcta y los permisos de acceso adecuados.

## Direcciones de matriz y aritmética de punteros por elementos

<!-- wave-example: book-pointer-array -->
```wave
fun main() {
    var values: array<i32, 3> = [10, 20, 30];
    var first: ptr<i32> = &values[0];
    var second: ptr<i32> = first + 1;

    println("{}", deref first);
    println("{}", deref second);
    println("{}", deref first[2]);
}
```

Resultado de la ejecución:

```text
10
20
30
```

Agregar 1 al puntero lo mueve un elemento del tipo de destino. El siguiente elemento en i32 y el siguiente elemento en u8 tienen diferentes números de bytes de desplazamiento. Si multiplica `first + 1` nuevamente por el tamaño de letra y lo suma, se moverá a una posición no deseada.

La indexación del puntero también debe realizarse dentro del rango válido. first no recuerda la longitud 3 de la matriz en sí, por lo que al pasar el rango a la función, utiliza un formulario que recibe tanto un puntero como una longitud.

## Pasar el rango de lectura a una función

<!-- wave-example: book-pointer-range -->
```wave
fun sum(values: ptr<i32>, count: i32) -> i32 {
    var total: i32 = 0;

    for (var index: i32 = 0; index < count; index += 1) {
        total += deref values[index];
    }

    return total;
}

fun main() {
    var values: array<i32, 4> = [2, 4, 6, 8];

    println("first two={}", sum(&values[0], 2));
    println("all={}", sum(&values[0], 4));
}
```

Resultado de la ejecución:

```text
first two=6
all=20
```

La unidad de count es el número de elementos. Esta función requiere que la persona que llama tenga count legible i32. Pasar una longitud mayor que la matriz real viola el contrato. Incluso si usa la misma combinación ptr<u8>·i64, API, debe verificar en la documentación si la longitud está en bytes o en elementos.

## Vida útil: ¿Cuánto tiempo es válida la dirección?

Las variables locales se utilizan durante la vida útil de la llamada y el bloque. Si devuelve la dirección de una variable local dentro de una función para que la persona que llama la lea más tarde, es posible que ese espacio de almacenamiento haya llegado al final de su vida útil.

Aquí hay una pieza de mal diseño que no debería implementarse:

```wave
fun invalid_address() -> ptr<i32> {
    var local: i32 = 42;

    return &local;
}
```

Si se requiere un valor, devuelve i32. Si necesita escribir en el espacio de almacenamiento proporcionado por la persona que llama, toma un puntero como entrada. Si necesita almacenamiento separado para retenerlo más allá de la llamada, asígnelo explícitamente y transfiera la responsabilidad de liberarlo.

## Dos punteros que apuntan al mismo espacio de almacenamiento.

Copiar un puntero crea otro nombre que apunta a la misma dirección. No duplica la memoria.

<!-- wave-example: book-pointer-alias -->
```wave
fun main() {
    var value: i32 = 1;
    var first: ptr<i32> = &value;
    var second: ptr<i32> = first;

    deref second = 9;

    println("{} {}", value, deref first);
}
```

Resultado de la ejecución:

```text
9 9
```

Los resultados modificados mediante second también se pueden ver mediante first. Una vez que se libera la memoria de origen, ambos punteros quedan inutilizables. Asignar null a una variable de puntero no cambia automáticamente las otras copias.

## Ejercicio: intercambiar dos números enteros

Escriba una función que tome dos direcciones i32 e intercambie sus valores. El primer valor debe almacenarse en una variable temporal antes de sobrescribirlo. Compruebe si el valor se conserva incluso si pasa la misma dirección dos veces.

### Solución completa

<!-- wave-example: book-pointer-swap -->
```wave
fun swap(left: ptr<i32>, right: ptr<i32>) {
    var saved: i32 = deref left;

    deref left = deref right;
    deref right = saved;
}

fun main() {
    var first: i32 = 3;
    var second: i32 = 8;

    swap(&first, &second);
    println("{} {}", first, second);

    swap(&first, &first);
    println("{}", first);
}
```

Resultado de la ejecución:

```text
8 3
8
```

Esta función también requiere que ambas direcciones apunten a un almacenamiento de enteros válido y grabable. Para manejar null, agregue un resultado que indique éxito o fracaso, como en try_increment.


## **Wave Explicit Memory Type Model**

El diseño del puntero de Wave se basa en **Wave Explicit Memory Type Model**. Este modelo define punteros y matrices como tipos de memoria explícita a nivel de lenguaje, en lugar de trucos sintácticos o abstracciones de biblioteca.

`ptr<T>` es un tipo que apunta a la dirección de memoria que almacena el valor `T`, y `array<T, N>` es un tipo de memoria de longitud fija que almacena valores `N` de `T` en sucesión. Por lo tanto, la estructura de punteros y matrices se revela tal cual en argumentos de función, valores de retorno, campos de estructura y otros tipos.

## null

```wave
var buffer: ptr<u8> = null;
if (buffer == null) {
    println("no buffer");
}
```

`null` es un valor de puntero que no apunta a una dirección de memoria válida. `null` solo se puede asignar al tipo `ptr<T>` y no se puede utilizar como valor entero, booleano o de matriz.

Una función de asignación o búsqueda puede devolver `null` cuando no tiene ningún resultado. Verifique `null` antes de eliminar la referencia a dicho resultado. Al eliminar la referencia a un puntero `null` no se accede a un almacenamiento válido.

## conversión de puntero

Cuando necesite cambiar la dirección u otra representación del puntero, utilice `as`.

```wave
var raw: i64 = 0;
var p: ptr<u8> = raw as ptr<u8>;
```

Utilice conversiones entre números enteros y punteros solo en límites de bajo nivel y considere el ancho de la dirección de la plataforma de destino y ABI.
