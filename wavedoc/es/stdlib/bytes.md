---
translation_set_id: stdlib-bytes
path: stdlib/bytes
locale: es
group: stdlib
group_order: 1
order: 6
title: bytes: rangos, cursores y ULEB128
summary: Describe bytes de lectura/escritura con longitud view y preservación del estado en caso de falla.
---

## ¿En qué se diferencia de una cuerda?

Los datos de bytes pueden contener cero, así que pase un puntero junto con una longitud. `Bytes` y `BytesMut` son vistas que no son propietarias y son válidas solo mientras el almacenamiento subyacente siga siendo válido. `BytesMut` requiere almacenamiento grabable.

## Crear cursor

```text
std::bytes::cursor
bytes_reader(data: ptr<u8>, len: i64) -> ByteReader
bytes_writer(data: ptr<u8>, len: i64) -> ByteWriter
bytes_reader_read_be_u16(reader: ptr<ByteReader>, out_value: ptr<u16>) -> i32
bytes_writer_write_be_u16(writer: ptr<ByteWriter>, value: u16) -> i32
```

Importe `ByteReader` y `ByteWriter` desde `std::bytes::types`. Su campo de posición identifica la ubicación de la próxima operación. Su longitud es el recuento total de bytes accesibles. La construcción de cualquiera de los cursores no copia ni asigna la memoria subyacente.

`be` es big-endian, `le` es little-endian. Si el tipo de archivo es big-endian, utilice la función `be` independientemente del orden de bytes del host CPU. Hay funciones de lectura/escritura y de un byte signed/unsigned de 16, 32 y 64 bits.

## Errores y preservación del estado.

`BYTES_OK` de `std::bytes::errors` es 0. INVALID indica un rango no válido, EOF entrada insuficiente, NO_SPACE capacidad de salida insuficiente y OVERFLOW un valor fuera del rango representable. Las operaciones del cursor marcadas avanzan la posición solo después de que toda la operación se haya realizado correctamente. Una lectura fallida también conserva el valor de salida.

## ULEB128

```text
std::bytes::leb128
bytes_reader_read_uleb128_u64(reader: ptr<ByteReader>, output_value: ptr<u64>) -> i32
bytes_writer_write_uleb128_u64(writer: ptr<ByteWriter>, value: u64) -> i32
```

ULEB128 almacena un entero de 64 bits sin signo en un número variable de bytes, utilizando como máximo 10 bytes. Si no hay suficiente espacio, el escritor conserva tanto su posición como los bytes de destino. El lector distingue una entrada incompleta de un valor que excede u64. Se aceptan codificaciones terminadas y no mínimas.

Para crear un mensaje real y ver un error de entrada breve, continúe con [Práctica de mensajes binarios](/docs/es/practice/binary-message). No intente generar una cadena de bytes que contenga ceros como `str`.

## Leer los mismos bytes en diferentes órdenes.

El orden de los bytes es la regla de almacenamiento de los números. Si lees los dos bytes 1 y 2 como big-endian, es 1×256+2, y si los lees como little-endian, es 2×256+1. Seleccione según las reglas de red o tipo de archivo.

<!-- wave-example: book-bytes-endian -->
```wave
import("std::bytes::types")::{Bytes, bytes_view};
import("std::bytes::read")::{bytes_read_be_u16, bytes_read_le_u16};

fun main() -> i32 {
    var data: array<u8, 2> = [1, 2];
    var view: Bytes = bytes_view(&data[0], 2);
    var big: u16 = 0;
    var little: u16 = 0;

    if (bytes_read_be_u16(view, 0, &big) < 0) {
        return 1;
    }

    if (bytes_read_le_u16(view, 0, &little) < 0) {
        return 2;
    }

    println("be={} le={}", big, little);
    return 0;
}
```

Resultado de la ejecución:

```text
be=258 le=513
```

view toma prestada una matriz. Dado que no hay una asignación o copia separada, la matriz solo se puede utilizar mientras sea válida. offset está en bytes y la lectura de 16 bits requiere 2 bytes desde esa posición.

## Seleccione offset API y cursor API

La función read/write, que toma el argumento offset, es conveniente para formatos que leen directamente una posición de campo específica. Para transmisiones donde la siguiente posición depende de la longitud del campo anterior, es conveniente cursor con position.

Al mezclar los dos, dejar claro cuál es el estándar: cursor.position o separar offset. Evite el error de agregar la misma ubicación dos veces o pasar a la siguiente ubicación sin una lectura exitosa.

## Verificar el estado desde una entrada corta

<!-- wave-example: book-bytes-short-read -->
```wave
import("std::bytes::types")::{ByteReader};
import("std::bytes::cursor")::{bytes_reader, bytes_reader_read_be_u16};
import("std::bytes::errors")::{BYTES_ERROR_EOF};

fun main() {
    var data: array<u8, 1> = [1];
    var reader: ByteReader = bytes_reader(&data[0], 1);
    var value: u16 = 99;
    var status: i32 = bytes_reader_read_be_u16(&reader, &value);

    if (status == BYTES_ERROR_EOF) {
        println("position={} value={}", reader.position, value);
    }
}
```

Resultado de la ejecución:

```text
position=0 value=99
```

No se necesitan dos bytes, por lo que se conservan el valor de salida y la ubicación. Esta propiedad es útil en diseños que reintentan leer el mismo campo después de adquirir más entradas. Sin embargo, si se ha reasignado el espacio de almacenamiento señalado por view, la dirección también debe actualizarse.

## Límite de ULEB128

0~127 usa un byte, 128 en adelante usa más bytes. El bit de orden superior de cada byte indica si le siguen datos. Un valor fuera de u64 o una entrada continua que es demasiado larga es OVERFLOW, que es diferente de EOF, que simplemente tiene menos entradas.

Verifique directamente si Capacidad insuficiente writer conserva su estado.

<!-- wave-example: book-leb128-space -->
```wave
import("std::bytes::types")::{ByteWriter};
import("std::bytes::cursor")::{bytes_writer};
import("std::bytes::leb128")::{bytes_writer_write_uleb128_u64};
import("std::bytes::errors")::{BYTES_ERROR_NO_SPACE};

fun main() {
    var data: array<u8, 1> = [85];
    var writer: ByteWriter = bytes_writer(&data[0], 1);
    var status: i32 = bytes_writer_write_uleb128_u64(&writer, 128);

    if (status == BYTES_ERROR_NO_SPACE) {
        println("position={} byte={}", writer.position, data[0]);
    }
}
```

Resultado de la ejecución:

```text
position=0 byte=85
```

128 requiere dos bytes pero solo un espacio. Tras un fallo también permanece el primer byte 85. Verificar los tipos de error y la preservación del estado juntos describe el límite mejor que una simple verificación de éxito roundtrip.

## Secuencia de creación del analizador de mensajes

1. Lea el encabezado fijo y verifique el tipo y la versión.
2. Lea la longitud y compárela con el rango de entrada restante.
3. Pase solo los datos necesarios a view o a un búfer independiente.
4. Si el formato requiere el mensaje completo, también se verifican los bytes adicionales.
5. Distingue entre EOF y errores de formato no válido y los transmite a la persona que llama.

Puede crear un programa conectando campos en [Práctica de mensajes binarios](/docs/es/practice/binary-message).
