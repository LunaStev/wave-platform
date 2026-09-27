---
translation_set_id: practice-file-reader
path: practice/file-reader
locale: es
group: practice
group_order: 4
order: 2
title: Proyecto: leer un archivo y contar bytes
summary: Buffer, archivo I/O, vincula el manejo de fallas y la liberación de memoria.
---

## listo

Cree input.txt en el directorio de trabajo que contiene `Wave` seguido de una única nueva línea LF. El archivo contiene entonces 5 bytes. Con CRLF contiene 6 bytes; una lista de materiales UTF-8 agrega más bytes. Verifique la codificación de archivos y los finales de línea del editor.

Guarde el programa como `main.wave` y ejecútelo como `wavec run main.wave` desde el mismo directorio. Las rutas relativas son relativas al directorio de trabajo en ejecución.

<!-- wave-example: file-reader -->
```wave
import("std::buffer::alloc")::{
    Buffer, buffer_init, buffer_free
};
import("std::fs::file")::{
    read_to_end
};

fun main() -> i32 {
    var data: Buffer;
    if (buffer_init(&data, 0) < 0) {
        return 1;
    }

    var count: i64 = read_to_end("input.txt", &data);
    if (count < 0) {
        println("read failed={}", count);
        buffer_free(&data);
        return 2;
    }

    var lines: i64 = 0;
    for (var i: i64 = 0; i < data.len; i += 1) {
        if (deref data.data[i] == 10) {
            lines += 1;
        }
    }

    println("bytes={} LF={}", count, lines);
    if (buffer_free(&data) < 0) {
        return 3;
    }

    return 0;
}
```

Resultado de la ejecución:

```text
bytes=5 LF=1
```

## Acciones y Responsabilidades

`read_to_end` abre y cierra el archivo, pero la liberación Buffer es responsabilidad de la persona que llama. Se desactiva tanto en las rutas de éxito como en las de error de lectura. Dado que se trata de datos con una cantidad de bytes, no asumimos que sea una cadena que termina en NUL.

LF El número de líneas y el número de líneas que la gente piensa no siempre son los mismos. Si la última línea no contiene LF, no se incluye en el recuento LF de este programa. También verifique si hay fallas de lectura cambiando el nombre del archivo. Los números de error específicos pueden variar según su entorno.

## Ejercicios extendidos y comentarios.

Si desea manejar archivos muy grandes, reemplácelo con una matriz de tamaño fijo y una iteración `io_read`. Procesa solo rangos de retorno positivos y termina en 0. Puede acumular recuentos de bytes y recuentos LF sin tener que mantener todo el archivo en la memoria. Si lo abrió usted mismo, también cierra el descriptor en cualquier ruta de salida.

[Ver fs y io](/docs/es/stdlib/files-io) · [Ver Buffer](/docs/es/stdlib/buffer)

## Siga el flujo de procesamiento

1. Inicializar bin Buffer. Aún no hay contenido del archivo.
2. read_to_end lee el archivo y aumenta el espacio requerido.
3. Si la lectura es exitosa, se verifican los bytes en el rango data.len.
4. Siempre que encontramos el valor de byte 10 de LF, incrementamos lines.
5. Imprima los resultados y suelte Buffer.

data.cap es el espacio de almacenamiento reservado y data.len es la longitud de datos válida. Si cambia la condición de repetición a cap, se leerán los bytes que no estaban en el archivo, así que use len. count es el número de bytes agregados en esta llamada a read_to_end. Este ejemplo comienza con un Buffer vacío, por lo que count y data.len son iguales.

## Cambiar entrada para comprobar

|input.txt Contenido| bytes | LF |razón|
| --- | --- | --- | --- |
|archivo vacío| 0 | 0 |No hay bytes para leer|
| `Wave` | 4 | 0 |Sin salto de línea final|
| `Wave` + LF | 5 | 1 |Datos hasta el último LF|
| `A` + LF + `B` + LF | 4 | 2 |Cuenta dos LF|
| `Wave` + CRLF | 6 | 1 |CR también es de 1 byte, pero solo se cuenta LF.|

Para contar archivos sin LF en la última línea como una línea, agregue 1 al número de líneas si el archivo no está vacío y el último byte no es 10. Primero debe verificar si data.len es 0 antes de poder acceder al último elemento.

## Expandir a archivos grandes

El método actual de archivar todo el contenido es conveniente para una posterior relectura o recuperación de los datos. Si solo necesita la cantidad de bytes y el recuento LF, tiene sentido reutilizar un búfer de tamaño fijo.

En el bucle io_read del ejemplo [Leer un archivo en un búfer de tamaño fijo](/docs/es/stdlib/files-io), simplemente cuente LF por el número de bytes devueltos. Se puede procesar con una memoria igual al tamaño del búfer, no con la longitud total del archivo. Si la lectura devuelve 0, se acabó; si es negativo, es un error.
