---
translation_set_id: practice-binary-message
path: practice/binary-message
locale: es
group: practice
group_order: 4
order: 3
title: Proyecto: Crear un mensaje binario
summary: Utilice ordenamiento explícito de bytes y ULEB128 y rechace la entrada corta.
---

## formato de mensaje

Los primeros 2 bytes almacenan los números de tipo big-endian u16 y luego los valores ULEB128 u64. Si escribe memoria de estructura en un archivo tal como está, se verá afectado por el relleno y el orden de bytes, así que codifíquelo por campo.

Guárdelo como `main.wave` y ejecútelo.

<!-- wave-example: binary-message -->
```wave
import("std::bytes::types")::{
    ByteReader, ByteWriter
};
import("std::bytes::cursor")::{
    bytes_reader, bytes_writer, bytes_writer_write_be_u16, bytes_reader_read_be_u16
};
import("std::bytes::leb128")::{
    bytes_writer_write_uleb128_u64, bytes_reader_read_uleb128_u64
};

fun main() -> i32 {
    var data: array<u8, 12>;
    var writer: ByteWriter = bytes_writer(&data[0], 12);
    if (bytes_writer_write_be_u16(&writer, 7) < 0) {
        return 1;
    }
    if (bytes_writer_write_uleb128_u64(&writer, 300) < 0) {
        return 2;
    }

    var reader: ByteReader = bytes_reader(&data[0], writer.position);
    var kind: u16 = 0;
    var value: u64 = 0;
    if (bytes_reader_read_be_u16(&reader, &kind) < 0) {
        return 3;
    }
    if (bytes_reader_read_uleb128_u64(&reader, &value) < 0) {
        return 4;
    }

    println("kind={} value={} bytes={}", kind, value, reader.position);
    var short: ByteReader = bytes_reader(&data[0], 1);
    kind = 99;
    if (bytes_reader_read_be_u16(&short, &kind) >= 0) {
        return 5;
    }

    println("preserved={} {}", short.position, kind);
    return 0;
}
```

Resultado de la ejecución:

```text
kind=7 value=300 bytes=4
preserved=0 99
```

## Que comprobar

En reader, pasamos la longitud escrita real, no la capacidad total de la matriz de 12. Esto es para evitar leer bytes finales no inicializados como entrada. Incluso si falla una lectura de entrada breve, position=0 y kind=99 se mantienen.

## Ejercicios extendidos y comentarios.

Si el formato no permite bytes adicionales al final, verifique `reader.position == reader.len` después de completar el análisis. Al agregar un campo de longitud, asegúrese de que no sea mayor que los bytes restantes de la entrada y que el cálculo de longitud+offset no exceda el rango.

[Ver bytes](/docs/es/stdlib/bytes)

## Mira los bytes reales

El tipo número 7 es big-endian u16 y por tanto `00 07`. El valor 300 pasa a ser de ULEB128 a `AC 02`. El mensaje completo consta de los siguientes cuatro bytes:

```text
00 07 AC 02
───── ─────
종류  값
```

ULEB128 utiliza los 7 bits bajos para cada byte para el valor, y si el bit alto es 1, indica que le sigue el siguiente byte. Los 7 bits inferiores de 300 son 44 y el resto es 2. El primer byte es 44 más 128 marcas consecutivas, o 172, o 0xAC. No hay ninguna marca de continuación en el último byte 0x02.

## Capacidad y duración de uso.

u16 usa 2 bytes y ULEB128 de u64 usa hasta 10 bytes, por lo que una matriz de 12 bytes puede almacenar ambos campos. Un valor de 300 utiliza sólo 2 bytes, lo que hace que el mensaje real tenga 4 bytes. Al enviar a un archivo o socket, envía el byte writer.position en lugar de toda la matriz.

position en reader es la posición de lectura actual. Después de leer el tipo, se convierte en 2 y después de leer el valor, se convierte en 4. Para pasar al siguiente mensaje cuando se produce un error, se debe conocer el límite del mensaje fallido. Una política de recuperación de mensajes no se establece automáticamente simplemente porque la función de lectura conserva la ubicación.

## Ejercicio de valor límite

Cambie los valores a 0, 127, 128, 16383, 16384 para determinar la longitud de codificación. Al cambiar de 127 a 128, la longitud de ULEB128 aumenta de 1 a 2, y al cambiar de 16383 a 16384, la longitud aumenta de 2 a 3. Calcule la longitud total, incluidos los 2 bytes del campo tipo.

También se comprueban los mensajes cuyo último byte está truncado. Si la longitud del mensaje original era 4, pasamos la longitud 3 a reader. El campo de tipo se lee, pero la lectura del valor debe fallar porque falta el último byte de ULEB128. En este momento, verifique si se mantienen la posición 2 y el valor de salida justo antes de la lectura ULEB128.
