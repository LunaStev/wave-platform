---
translation_set_id: stdlib-task
path: stdlib/task
locale: es
group: stdlib
group_order: 1
order: 13
title: task: Ejecución y limpieza de Future
summary: Describe el consumo único, la ejecución y la cancelación después de la limpieza de operaciones asincrónicas.
---

## Básico API

Importado como `import("std::task" as task);`.

```text
block_on<T>(future: Future<T>) -> T
spawn<T>(future: Future<T>) -> Future<T>
cancel<T>(future: Future<T>) -> bool
yield_now() -> Future<void>
sleep_ms(milliseconds: i64) -> Future<void>
cancel_and_join<T>(future: Future<T>) -> Future<void>
shutdown()
```

`task::block_on(pending)` se ejecuta hasta que Future se completa y devuelve el resultado. El tipo de resultado se determina a partir del Future pasado.

## reglas de la vida

Future es un identificador de consumo único. No pienso en copiar valores como dos operaciones independientes. También es responsable de esperar los resultados o cancelar/organizar las tareas programadas con `spawn`. Future ya consumido con `await` y `block_on` no se volverá a consumir.

cancelar solicitudes de cancelación; por sí solo no garantiza que la limpieza haya terminado. Espere a que se complete lo requerido antes de liberar recursos. No libere la memoria prestada por E/S asincrónicas mientras una tarea aún pueda acceder a ella. Llame al `shutdown` después de que las tareas hayan completado o terminado la limpieza de cancelación.

## ejecución colaborativa

Los cálculos prolongados y las llamadas blocking sincrónicas pueden ralentizar el progreso general del ejecutor. La espera asincrónica con yield genera una oportunidad de ejecución. El bloqueo I/O no se vuelve asíncrono solo porque esté dentro de una función async.

Consulte la secuencia de ejecución y limpieza en el programa completo en [Introducción al código asincrónico](/docs/es/language/async-and-never).

## Ceder y esperar a que se complete

El siguiente programa produce la ejecución en medio de una operación y devuelve un resultado. Dado que yield no es una terminación de función, el código posterior a await se ejecuta continuamente.

<!-- wave-example: book-task-yield -->
```wave
import("std::task" as task);

async fun calculate() -> i32 {
    println("started");
    await task::yield_now();
    println("resumed");
    return 42;
}

fun main() -> i32 {
    var pending: Future<i32> = calculate();
    var result: i32 = task::block_on(pending);

    println("result={}", result);
    task::shutdown();
    return 0;
}
```

Resultado de la ejecución:

```text
started
resumed
result=42
```

No hay otras operaciones en este ejemplo, por lo que el orden de salida es coherente. En un programa que tiene múltiples tareas spawn, otras tareas pueden continuar en el punto yield, por lo que no depende del orden de salida de las diferentes tareas.

## Orden de organización de las tareas.

1. Prepara el espacio de almacenamiento y los recursos que necesitas para tu trabajo.
2. Cree Future y ejecútelo como await, block_on o spawn.
3. Si necesita resultados, espere hasta que finalice.
4. Si canceló una tarea en ejecución, espere a que se limpie.
5. Limpia buffers, archivos y sockets tomados prestados por la tarea.
6. Si no queda trabajo, llame al shutdown.

Future Dejar el alcance de una variable es una cosa y limpiar el trabajo de forma segura es otra. En particular, si pasa la dirección de una matriz local de función a una tarea async, la tarea debe terminar de usar esa matriz antes de que la función regrese.

## async Funciones divisorias y funciones generales

Los cálculos puros se pueden separar en funciones regulares. Adjunte async a la función que necesita expresar espera y espere hasta que se complete con await dentro de esa función. Simplemente empaquetar una función síncrona que lleva mucho tiempo, como leer un archivo, con la función async no da la oportunidad de ejecutar otras tareas.

Puede comparar al crear y ejecutar Future en [aprendizaje asincrónico](/docs/es/language/async-and-never).
