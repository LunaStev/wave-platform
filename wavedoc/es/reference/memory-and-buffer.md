---
translation_set_id: memory-buffer
path: reference/memory-and-buffer
locale: es
group: stdlib
group_order: 1
order: 4
title: mem: Asignación, reasignación y diseño
summary: Describe el tamaño en bytes, errores de asignación, límites de reasignación y responsabilidades de liberación.
---

## Asignación y desasignación

```text
std::mem::alloc
mem_alloc(size: i64) -> ptr<u8>
mem_alloc_zeroed(size: i64) -> ptr<u8>
mem_free(p: ptr<u8>, size: i64) -> i64
mem_realloc(old_ptr: ptr<u8>, old_size: i64, new_size: i64) -> ptr<u8>
```

El tamaño está en bytes. Una asignación de tamaño 0 o menos devuelve null. Si la asignación también falla para tamaños positivos, puede ser null. No asuma el contenido inicial de `mem_alloc`, pero use `mem_alloc_zeroed` si se requiere inicialización cero.

La persona que llama es propietaria de cada asignación exitosa y debe pasar su tamaño original al liberarla. `mem_free(null, size)` devuelve 0. Un puntero que no sea null emparejado con un tamaño no positivo es un error. Nunca acceda ni libere una asignación después de que ya haya sido liberada.

## Clasificación en caso de reasignación

|solicitud|acción|
| --- | --- |
|El nuevo tamaño es positivo y exitoso.|`min(old_size, new_size)` Copiar byte y desasignar el anterior|
|Fracasa nueva asignación de tamaño positivo|null Devolver, mantener la asignación existente|
|`old_ptr == null`, nuevo tamaño positivo|se comporta como una nueva tarea|
| `new_size == 0` |Intenta liberar una asignación anterior válida y devuelve null|
|old_size=0 para tamaño negativo o puntero existente|null Volver|

Un resultado null de la reasignación al tamaño cero no establece si la liberación se realizó correctamente. Llame a `mem_free` directamente si necesita su estado. Después de aumentar una asignación, inicialice usted mismo la región recién agregada.

## Ejemplo de conservación de un puntero existente

A continuación se muestra el caso en el que el nuevo tamaño es positivo. Guárdelo como `main.wave` y ejecútelo.

<!-- wave-example: reallocation -->
```wave
import("std::mem::alloc")::{
    mem_alloc, mem_realloc, mem_free
};

fun main() -> i32 {
    var data: ptr<u8> = mem_alloc(4);
    if (data == null) {
        return 1;
    }
    deref data[0] = 7;
    var grown: ptr<u8> = mem_realloc(data, 4, 8);
    if (grown == null) {
        mem_free(data, 4);
        return 2;
    }
    data = grown;
    println("{}", deref data[0]);
    if (mem_free(data, 8) < 0) {
        return 3;
    }

    return 0;
}
```

Resultado de la ejecución:

```text
7
```

Si lo sobrescribe con `data = mem_realloc(...)` antes de confirmar el error, puede perder la dirección existente. Si la reasignación se realiza correctamente, la dirección anterior y los punteros que apuntan a ella no se utilizan.

## El tamaño y la alineación de un tipo de objetivo.

```text
std::mem::layout
size_of<T>() -> u64
align_of<T>() -> u64
```

Ambos valores son el diseño del destino de compilación, no la computadora en la que se ejecuta. `size_of` incluye relleno de cola y no genera ni evalúa un valor. Al multiplicar el número de elementos por su tamaño, comprobamos si hay desbordamiento. Puedes usar `mem_size_mul_checked` y `mem_size_add_checked` de `std::mem::ops`.

`mem_copy` se usa para copiar un rango que no se superpone y `mem_move` se usa para copiar un rango que puede superponerse. Ninguno de los dos puede determinar la longitud real de la asignación únicamente a partir de punteros, por lo que la persona que llama debe garantizar los límites. Se puede gestionar una lista de bytes de diferentes tamaños con [Buffer](/docs/es/stdlib/buffer).

## Ejemplo de consulta de diseño

Guárdelo como main.wave y ejecútelo. El tamaño y la alineación del objetivo cubierto en el documento, i32, son cada uno de 4 bytes, por lo que se genera `4 4`. Verifique los valores de otros tipos, especialmente estructuras y punteros, por destino.

<!-- wave-example: layout-api -->
```wave
import("std::mem::layout")::{
    size_of, align_of
};

fun main() {
    println("{} {}", size_of<i32>(), align_of<i32>());
}
```

## Compruebe si hay desbordamiento en el cálculo del tamaño

Debe verificar si `count * element_size` es válido antes de pasarlo a la función de asignación. Si asigna un espacio pequeño con el valor de desbordamiento y escribe tanto como el número original, se saldrá de los límites.

<!-- wave-example: book-memory-size-check -->
```wave
import("std::mem::ops")::{mem_size_add_checked, mem_size_mul_checked};

fun main() {
    var result: i64 = 99;

    if (mem_size_mul_checked(3, 4, &result) == 0) {
        println("bytes={}", result);
    }

    if (mem_size_add_checked(9223372036854775807, 1, &result) < 0) {
        println("overflow rejected");
    }
}
```

Resultado de la ejecución:

```text
bytes=12
overflow rejected
```

Los resultados de las fallas no se escriben en el tamaño de la asignación. No se pueden detectar todos los desbordamientos simplemente calculándolos con aritmética regular y viendo si el resultado es negativo. Para calcular el tamaño que debe inspeccionarse, utilice la función checked desde el principio.

## Para copias superpuestas, mem_move

Al mover parte de la misma matriz hacia atrás, las áreas de entrada y salida se superponen. No pase rangos superpuestos a mem_copy, use mem_move.

<!-- wave-example: book-memory-overlap -->
```wave
import("std::mem::ops")::{mem_move};

fun main() {
    var data: array<u8, 5> = [1, 2, 3, 4, 5];

    mem_move(&data[1], &data[0], 4);

    for (var index: i32 = 0; index < 5; index += 1) {
        println("{}", data[index]);
    }
}
```

Resultado de la ejecución:

```text
1
1
2
3
4
```

Los primeros cuatro bytes originales se mueven una posición hacia la derecha. Copiarlos manualmente puede leer valores que ya se han sobrescrito, generando accidentalmente todos unos. Una API que reconoce la superposición maneja la dirección de la copia por usted.

## Escribe una función que pase la propiedad.

Para funciones que devuelven memoria, es mejor proporcionar una dirección de retorno en caso de éxito y el tamaño necesario para liberarla. Si la persona que llama tiene que adivinar el tamaño, es probable que se publique incorrectamente. Si una función devuelve una dirección prestada, la persona que llama no debe liberarla y describe la vida útil de la original.

A través de los límites de las funciones, debería poder rastrear `allocator → owner → deallocation`. El nombre o tipo de una variable de puntero no determina automáticamente la propiedad.
