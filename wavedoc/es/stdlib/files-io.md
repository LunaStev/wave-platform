---
translation_set_id: stdlib-files-io
path: stdlib/files-io
locale: es
group: stdlib
group_order: 1
order: 7
title: fs e io: archivos y transferencias de bytes
summary: Describe la vida útil del archivo, la lectura completa, la transferencia parcial y los estados posteriores al fallo.
---

## Seleccionar archivo API

La función de conveniencia de `std::fs::file` recibe el camino y realiza la apertura y cierre necesarios. La persona que llama debe cerrar las funciones que devuelven un descriptor.

|declaración|Resultados exitosos y precauciones.|
| --- | --- |
| `open_read(path: str) -> i64` |Abrir descriptor. Los números negativos son errores.|
| `create(path: str) -> i64` |Cree o elimine el contenido del archivo existente. Devuelve el descriptor propietario.|
| `open_append(path: str) -> i64` |Abrir o crear para agregar|
| `size(path: str) -> i64` |Número de bytes. Los números negativos son errores.|
| `read_into(path: str, dst: ptr<u8>, dst_cap: i64) -> i64` |Número de bytes en todo el archivo. La falta de capacidad es un error.|
| `read_to_end(path: str, dst_buffer: ptr<Buffer>) -> i64` |Agregue el archivo después del Buffer existente y devuelva el monto adicional|
| `write(path: str, src: ptr<u8>, len: i64) -> i64` |Número de bytes escritos que reemplazan el contenido existente|
| `append(path: str, src: ptr<u8>, len: i64) -> i64` |Número de bytes agregados al final.|
| `remove(path: str) -> i64` |estado de eliminación. el fracaso es negativo|

false de `exists(path)` por sí solo no puede distinguir entre archivos faltantes y errores de permisos. Asegúrese de verificar el resultado de apertura real, ya que el estado puede cambiar entre la verificación de existencia y la apertura.

## Nivel bajo I/O

```text
std::io::fd
io_read(fd: i64, buf: ptr<u8>, len: i64) -> i64
io_write(fd: i64, buf: ptr<u8>, len: i64) -> i64
io_read_exact(fd: i64, buf: ptr<u8>, len: i64) -> i64
io_write_all(fd: i64, buf: ptr<u8>, len: i64) -> i64
io_close(fd: i64) -> i64
```

El resultado positivo de `io_read` es el número de bytes leídos y 0 en una solicitud de longitud positiva es EOF. `io_write` puede escribirse menos de lo solicitado. Si se requiere una transferencia completa, use la función exact/all. Aún así, no asumimos que la falla revierta el estado externo, ya que algunas transferencias pueden haber ocurrido antes del error.

`io_read_exact` es un error si EOF se encuentra antes de la longitud requerida. `read_into` devuelve `IO_ERR_NO_SPACE` si el búfer está lleno y es posible que ya se hayan escrito algunos bytes. La función de lectura no agrega automáticamente NUL al final de la cadena.

## Buffer y manejo de errores

En caso de falla, `read_to_end` restaura la lente original, pero es posible que su capacidad y dirección de datos hayan cambiado. La persona que llama debe liberar el Buffer después del éxito o del fracaso. Las API de escritura de archivos no garantizan el reemplazo atómico de archivos.

Desde [Práctica de lectura de archivos.](/docs/es/practice/file-reader), puede ejecutar el programa desde import para lanzarlo. Considere las diferencias de ruta/permisos en Linux/macOS/Windows/FreeBSD y las restricciones de directorio accesible en WASI. No interpreta directamente el valor del descriptor como un identificador sin formato de otro OS.

## Leer archivos grandes en buffers pequeños

Las operaciones que no requieren que todo el archivo se coloque en la memoria se pueden manejar con búferes fijos e iteraciones de lectura. El siguiente programa imprime el contenido de input.txt y cuenta el número total de bytes leídos. Guarde uno `Wave` y LF en el archivo de entrada.

<!-- wave-example: book-file-stream -->
```wave
import("std::fs::file")::{open_read};
import("std::io::fd")::{io_read, io_write_all, io_close};
import("std::io::consts")::{IO_STDOUT_FD};

fun main() -> i32 {
    var descriptor: i64 = open_read("input.txt");

    if (descriptor < 0) {
        return 1;
    }

    var buffer: array<u8, 4>;
    var total: i64 = 0;

    while (true) {
        var count: i64 = io_read(descriptor, &buffer[0], 4);

        if (count < 0) {
            io_close(descriptor);
            return 2;
        }

        if (count == 0) {
            break;
        }

        if (io_write_all(IO_STDOUT_FD, &buffer[0], count) < 0) {
            io_close(descriptor);
            return 3;
        }

        total += count;
    }

    if (io_close(descriptor) < 0) {
        return 4;
    }

    println("bytes={}", total);
    return 0;
}
```

Resultado de la ejecución:

```text
Wave
bytes=5
```

La capacidad del buffer es 4, pero la última lectura puede ser de 1 byte. Siempre pasamos el count real a la salida. Si escribe la matriz completa, se pueden generar incluso bytes antiguos y no leídos.

El programa abrió descriptor de input.txt, así que ciérrelo. La salida estándar no es un recurso recién adquirido en esta función, por lo que no se cierra arbitrariamente al final del ejemplo.

## Lecturas completas y capacidad insuficiente

read_into recibe un espacio de almacenamiento fijo que contiene todo el archivo. Si el espacio es insuficiente, se trunca silenciosamente y devuelve NO_SPACE sin éxito.

<!-- wave-example: book-file-capacity -->
```wave
import("std::fs::file")::{read_into};
import("std::io::consts")::{IO_ERR_NO_SPACE};

fun main() {
    var data: array<u8, 2>;
    var status: i64 = read_into("input.txt", &data[0], 2);

    if (status == IO_ERR_NO_SPACE) {
        println("destination too small");
    }
}
```

Resultado de la ejecución:

```text
destination too small
```

No se supone que la lectura fallida no haya cambiado en absoluto el byte de destino. No lo utilice como contenido de archivo terminado, prepare un repositorio más grande o elija el método streaming. Incluso si consulta el tamaño primero, el resultado de la lectura real es el juicio final, ya que el archivo puede cambiar entre la consulta y la lectura.

## Diferencia entre escribir y adjuntar archivos

write reemplaza el contenido existente y append se agrega al final. La cantidad de bytes a almacenar se puede obtener directamente de la longitud de la cadena y pasarla. El NUL al final de la cadena generalmente no se incluye en el contenido del archivo de texto.

<!-- wave-example: book-file-write-append -->
```wave
import("std::fs::file")::{write, append, size, remove};
import("std::string::len")::{len};

fun main() -> i32 {
    var first: str = "Wave";
    var second: str = " study";

    if (write("output.txt", first as ptr<u8>, len(first) as i64) < 0) {
        return 1;
    }

    if (append("output.txt", second as ptr<u8>, len(second) as i64) < 0) {
        return 2;
    }

    println("bytes={}", size("output.txt"));

    if (remove("output.txt") < 0) {
        return 3;
    }

    return 0;
}
```

Resultado de la ejecución:

```text
bytes=10
```

Este ejemplo crea, reemplaza y finalmente elimina output.txt en el directorio de trabajo. Ejecute desde el directorio de práctica sin archivos existentes. El editor o programa de almacenamiento real puede requerir políticas de almacenamiento independientes, como archivos temporales y reemplazo.

## API Tabla de selección

|situación|seleccionar|
| --- | --- |
|Leer un archivo pequeño completo en un búfer fijo| read_into |
|Conserva todo el contenido sin saber el tamaño|read_to_end y Buffer|
|Procesar el contenido en orden en lugar de almacenarlo en su totalidad|open_read + io_read repetir|
|Leer registros de longitud fija| io_read_exact |
|Transmitir cadena de bytes completa| io_write_all |
|Manejo de archivos ya abiertos|fd función en lugar de función de ruta|

Después de seleccionar una función, verifique cómo cambian el búfer, la ubicación del archivo y los datos externos en caso de falla.
