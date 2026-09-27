---
translation_set_id: stdlib-buffer
path: stdlib/buffer
locale: es
group: stdlib
group_order: 1
order: 5
title: buffer: almacenamiento de bytes ampliable
summary: Buffer Describe las reglas de inicialización, adición, consulta, capacidad y liberación.
---

## Significado de Buffer

`Buffer` en `std::buffer::types` tiene `data: ptr<u8>`, `len: i64` y `cap: i64`. len es el número de bytes inicializados y en uso, y cap es el número total de bytes asignados. Mantenga siempre `0 <= len <= cap`. La cadena NUL no garantiza automáticamente la terminación.

## Básico API

|módulo|declaración|significado|
| --- | --- | --- |
| `std::buffer::alloc` | `buffer_init(out_buf: ptr<Buffer>, capacity: i64) -> i64` |Inicialice el nuevo repositorio. No recuperar en buffers que ya se poseen|
|mismo módulo| `buffer_free(buf: ptr<Buffer>) -> i64` |Desasignar. Vacío si tiene éxito|
|mismo módulo| `buffer_reserve(buf: ptr<Buffer>, required_cap: i64) -> i64` |Asegúrese de tener una capacidad total mínima. len Mantenido|
|mismo módulo| `buffer_clear(buf: ptr<Buffer>) -> i64` |Mantener capacidad y len=0|
|mismo módulo| `buffer_resize(buf: ptr<Buffer>, new_len: i64, value: u8) -> i64` |Cambie la longitud, complete nuevos bytes con value|
| `std::buffer::write` | `buffer_push(buf: ptr<Buffer>, value: u8) -> i64` |agregar un byte|
|mismo módulo| `buffer_append(buf: ptr<Buffer>, src: ptr<u8>, size: i64) -> i64` |sizeCopiar y agregar byte|
|mismo módulo| `buffer_append_str(buf: ptr<Buffer>, s: str) -> i64` |Agregue bytes de cadena excluyendo NUL|
| `std::buffer::read` | `buffer_get(buf: ptr<Buffer>, index: i64, out_value: ptr<u8>) -> i64` |Leer un byte en el rango|

El retorno de estado API devuelve `BUFFER_OK`(0) en caso de éxito. Distinga entre los errores INVALID, BOUNDS, OVERFLOW y ALLOC de `std::buffer::error`. El número no se interpreta como OS errno. `buffer_new` representa un error de asignación como un Buffer vacío, por lo que cuando necesite distinguir entre errores, utilice `buffer_init`.

## Ejemplo de ejecución

<!-- wave-example: buffer-api -->
```wave
import("std::buffer::alloc")::{
    Buffer, buffer_init, buffer_free
};
import("std::buffer::write")::{
    buffer_append_str, buffer_push
};
import("std::buffer::read")::{
    buffer_get
};

fun main() -> i32 {
    var data: Buffer;
    if (buffer_init(&data, 0) < 0) {
        return 1;
    }
    if (buffer_append_str(&data, "Hi") < 0 || buffer_push(&data, 33) < 0) {
        buffer_free(&data);
        return 2;
    }

    var value: u8 = 0;
    if (buffer_get(&data, 2, &value) < 0) {
        buffer_free(&data);
        return 3;
    }

    println("{} {}", data.len, value);
    if (buffer_free(&data) < 0) {
        return 4;
    }

    return 0;
}
```

Resultado de la ejecución:

```text
3 33
```

Guárdelo como `main.wave` y ejecútelo como `wavec run main.wave`. Una capacidad inicial de 0 no es un error, sino un búfer vacío válido. Libera espacio durante el procesamiento posterior.

## Esperanza de vida y fracaso

Hacer crecer el búfer puede cambiar los datos. No utilice una dirección prestada anteriormente después de una operación que pueda reasignarla. Copiar la estructura Buffer no duplica su asignación, así que asigne a esa asignación un único propietario.

`buffer_get` no cambia los argumentos de salida si falla. Por otro lado, la función de conveniencia `buffer_at` también representa errores como 0, así que use `buffer_get` para distinguir entre 0 bytes reales y fallas. Evite crear len/cap no válido cambiando directamente los campos públicos.

[Memoria API](/docs/es/reference/memory-and-buffer) · [Practica leer un archivo como Buffer](/docs/es/practice/file-reader)

## Observar la longitud y la capacidad por separado

reserve libera espacio de almacenamiento, pero no aumenta len. resize cambia la longitud real utilizada e inicializa la parte extendida al byte especificado. clear establece solo la longitud utilizada en 0, lo que permite reutilizar la asignación.

Guarde el siguiente programa como main.wave y ejecútelo. No depende del múltiplo de crecimiento exacto de la capacidad; simplemente garantiza que tenga el espacio que necesita.

<!-- wave-example: book-buffer-length-capacity -->
```wave
import("std::buffer::alloc")::{
    Buffer,
    buffer_init,
    buffer_reserve,
    buffer_resize,
    buffer_clear,
    buffer_free
};

fun main() -> i32 {
    var data: Buffer;

    if (buffer_init(&data, 0) < 0) {
        return 1;
    }

    if (buffer_reserve(&data, 20) < 0) {
        buffer_free(&data);
        return 2;
    }

    println("reserved length={}", data.len);

    if (buffer_resize(&data, 3, 7) < 0) {
        buffer_free(&data);
        return 3;
    }

    println("resized length={} first={}", data.len, deref data.data[0]);

    if (buffer_clear(&data) < 0) {
        buffer_free(&data);
        return 4;
    }

    println("cleared length={}", data.len);

    if (buffer_free(&data) < 0) {
        return 5;
    }

    return 0;
}
```

Resultado de la ejecución:

```text
reserved length=0
resized length=3 first=7
cleared length=0
```

Al aumentar a resize, pasamos value=7, por lo que los tres bytes nuevos que vemos son 7. El espacio solo asegurado con reserve no se lee como datos inicializados. cap permanecerá después de clear y se puede agregar nuevamente al mismo Buffer.

## Distinguir entre cero bytes y errores de búsqueda

buffer_get devuelve el estado y escribe los bytes reales como argumentos de salida. Incluso si los datos son 0, es un éxito normal.

<!-- wave-example: book-buffer-zero-bounds -->
```wave
import("std::buffer::alloc")::{Buffer, buffer_init, buffer_free};
import("std::buffer::write")::{buffer_push};
import("std::buffer::read")::{buffer_get};
import("std::buffer::error")::{BUFFER_ERR_BOUNDS};

fun main() -> i32 {
    var data: Buffer;

    if (buffer_init(&data, 0) < 0) {
        return 1;
    }

    if (buffer_push(&data, 0) < 0) {
        buffer_free(&data);
        return 2;
    }

    var value: u8 = 99;

    if (buffer_get(&data, 0, &value) < 0) {
        buffer_free(&data);
        return 3;
    }

    println("stored={}", value);
    value = 99;

    if (buffer_get(&data, 1, &value) == BUFFER_ERR_BOUNDS) {
        println("outside, preserved={}", value);
    }

    if (buffer_free(&data) < 0) {
        return 4;
    }

    return 0;
}
```

Resultado de la ejecución:

```text
stored=0
outside, preserved=99
```

El primer acierto es un éxito con lectura 0, el segundo acierto es un fallo fuera de límites. Incluso si value=99 permanece después de una falla, no significa que sea el valor leído del búfer. Asegúrese de verificar el estado juntos.

## Solución práctica: acumulación de bytes

Para sumar números del 0 al 9, repita buffer_push y verifique cada resultado. Guarde la suma en i64 y lea solo el rango `0 <= index < data.len`. Después de manejar el búfer, llamamos a buffer_free tanto en la ruta de éxito como en la de error.

<!-- wave-example: book-buffer-sum -->
```wave
import("std::buffer::alloc")::{Buffer, buffer_init, buffer_free};
import("std::buffer::write")::{buffer_push};

fun main() -> i32 {
    var data: Buffer;

    if (buffer_init(&data, 0) < 0) {
        return 1;
    }

    for (var value: i32 = 0; value < 10; value += 1) {
        if (buffer_push(&data, value as u8) < 0) {
            buffer_free(&data);
            return 2;
        }
    }

    var total: i64 = 0;

    for (var index: i64 = 0; index < data.len; index += 1) {
        total += deref data.data[index] as i64;
    }

    println("sum={}", total);

    if (buffer_free(&data) < 0) {
        return 3;
    }

    return 0;
}
```

Resultado de la ejecución:

```text
sum=45
```
