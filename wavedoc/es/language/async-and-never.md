---
translation_set_id: language-async-and-never
path: language/async-and-never
locale: es
group: language
group_order: 2
order: 13
title: 13. Funciones asíncronas y Future
summary: Conozca el papel de la ejecución retrasada Future, await y block_on.
---

## Expresar tareas de espera.

Para operaciones en espera como archivos, sockets y temporizadores, es necesario distinguir entre cálculos continuos y espera a que se completen. La función asincrónica expresa el resultado a completar como Future. Agregar async no crea automáticamente un nuevo hilo ni cambia todas las llamadas sincrónicas a asincrónicas.

Lea este capítulo después de Funciones, punteros y manejo de errores. El ejemplo es un programa nativo que utiliza el iniciador `std::task` y ejecuta cada archivo como `wavec run main.wave`.

## Crea Future y recibe resultados

<!-- wave-example: book-async-basic -->
```wave
import("std::task" as task);

async fun calculate(value: i64) -> i64 {
    return value * 2;
}

fun main() {
    var pending: Future<i64> = calculate(21);
    var result: i64 = task::block_on(pending);

    task::shutdown();
    println("{}", result);
}
```

Resultado de la ejecución:

```text
42
```

i64 escrito en la declaración de calculate es el valor obtenido después de la finalización. El resultado de la llamada en sí es Future<i64>. En condiciones normales main, ejecute Future con block_on y reciba el resultado completo.

Llame al cierre después de limpiar todas las tareas para liberar los recursos del ejecutor. No libere un búfer mientras una tarea todavía lo utiliza y no ignore las tareas sin terminar.

## La invocación y la ejecución del cuerpo son diferentes.

Las funciones asincrónicas se ejecutan de forma perezosa. Debe distinguirse de una función normal que ejecuta su cuerpo inmediatamente después de llamarla.

<!-- wave-example: book-async-lazy -->
```wave
import("std::task" as task);

static entered: i32 = 0;

async fun work() -> i32 {
    entered += 1;
    return 7;
}

fun main() {
    var pending: Future<i32> = work();

    println("before={}", entered);

    var result: i32 = task::block_on(pending);

    task::shutdown();
    println("after={} result={}", entered, result);
}
```

Resultado de la ejecución:

```text
before=0
after=1 result=7
```

Cuando se creó Future, entered todavía era 0. Después de impulsar la ejecución, el cuerpo se ejecuta y se convierte en 1. Es por eso que no debe considerar la tarea como completada solo porque Future está almacenado en una variable.

## Esperando dentro de una función asincrónica

Dentro de la función async, await espera la finalización de otro Future. El resultado de una expresión await es un valor de finalización.

<!-- wave-example: book-async-nested -->
```wave
import("std::task" as task);

async fun twice(value: i32) -> i32 {
    await task::yield_now();
    return value * 2;
}

async fun process() -> i32 {
    var value: i32 = await twice(20);
    return value + 2;
}

fun main() {
    var answer: i32 = task::block_on(process());

    task::shutdown();
    println("{}", answer);
}
```

Resultado de la ejecución:

```text
42
```

process espera Future de twice y luego suma 2 al valor de finalización. Puede utilizar un valor similar al resultado de una llamada a una función normal, pero mientras espera, puede pasar la oportunidad de ejecución a otra tarea.

yield_now genera oportunidades de ejecución cooperativa. No dar ni una sola vez en un ciclo de cálculo largo puede ralentizar otras tareas. El asincrónico CPU no es un dispositivo para paralelizar y distribuir cálculos automáticamente.

## Programe múltiples tareas

Puedes programar tareas con spawn y esperar sus respectivos resultados. Este es un ejemplo de cómo verificar el resultado final sin depender del orden de salida intermedio de las dos operaciones.

<!-- wave-example: book-async-spawn -->
```wave
import("std::task" as task);

async fun compute(value: i32) -> i32 {
    await task::yield_now();
    return value * 2;
}

async fun combine() -> i32 {
    var first: Future<i32> = task::spawn(compute(10));
    var second: Future<i32> = task::spawn(compute(20));

    var left: i32 = await first;
    var right: i32 = await second;

    return left + right;
}

fun main() {
    var total: i32 = task::block_on(combine());

    task::shutdown();
    println("{}", total);
}
```

Resultado de la ejecución:

```text
60
```

Cada tarea que programes tiene un lugar donde espera los resultados. En lugar de crear una tarea y olvidar su identificador, debe decidir quién comprobará su finalización. await No considere que la secuencia y el orden en el que se ejecutan las operaciones internas sean iguales.

## Consumir Future una vez

Future se trata como una manija de consumo único. Copiar el mismo Future no lo hace esperar como dos tareas diferentes. No vuelva a hacer block_on o await para Future que ya se ha completado.

Si necesita el mismo resultado en varios lugares, en lugar de consumir Future varias veces, almacene el valor completo y páselo de acuerdo con las reglas de copiar/compartir para ese valor. También debe verificar si hay punteros o recursos propios dentro del valor.

## Diferencia entre temporizador y espera síncrona

async Cuando espera dentro de una función, puede usar `await task::sleep_ms(...)`. La llamada sincrónica sleep bloquea el flujo de ejecución actual, lo que también puede afectar el progreso de otras tareas en el ejecutor.

No espero que el tiempo de espera sea exactamente el número de milisegundos solicitados. Dependiendo de su programación y otras tareas, es posible que se despierte tarde. Al implementar tiempos de espera, utilizamos clock y deadline para medir el tiempo transcurrido, en lugar de esperar cada vez el tiempo completo original.

## Duración y cancelación del búfer

Un búfer pasado a E/S asincrónicas debe seguir siendo válido incluso cuando la función de llamada esté suspendida. Liberarlo o reasignarlo antes de completarlo o cancelar la limpieza puede dejar la operación con una dirección no válida.

No se puede suponer que la solicitud de cancelación y la finalización de la tarea sean al mismo tiempo. Verifique los resultados de la serie cancel API y espere a que se complete lo necesario antes de liberar recursos. Lea [Ver task](/docs/es/stdlib/task) para conocer las reglas de llamadas detalladas.

## malentendido común

|pensar|realmente comprobar|
| --- | --- |
|async Llamé y se acabó.|¿Realmente ejecutaste y completaste Future?|
|async Todas las llamadas dentro de la función son asincrónicas|¿El llamado API es síncrono o asíncrono?|
|Future Copiar duplica una tarea|¿Estás consumiendo el mismo mango repetidamente?|
|Como lo canceló, puede liberar el búfer inmediatamente.|¿Has terminado de organizar el trabajo después de la cancelación?|
|El orden de salida intermedio siempre es fijo|¿Espera explícitamente sólo el orden requerido para el resultado?|

## Ejercicio y solución completa

Cree una canalización que espere tres funciones asincrónicas seguidas. Devuelve el resultado de duplicar y sumar 5.

<!-- wave-example: book-async-solution -->
```wave
import("std::task" as task);

async fun read_value() -> i32 {
    return 10;
}

async fun transform(value: i32) -> i32 {
    await task::yield_now();
    return value * 2;
}

async fun pipeline() -> i32 {
    var initial: i32 = await read_value();
    var changed: i32 = await transform(initial);

    return changed + 5;
}

fun main() {
    var result: i32 = task::block_on(pipeline());

    task::shutdown();
    println("{}", result);
}
```

Resultado de la ejecución:

```text
25
```

Este ejemplo es intencionalmente una dependencia secuencial. transform requiere el resultado de read_value, por lo que simplemente hacer todo spawn no hace que la relación desaparezca. El punto de partida del diseño asincrónico es la distinción entre operaciones independientes y operaciones que requieren resultados.

Una vez que haya completado los conceptos básicos, conéctese a recursos externos reales con [Práctica de lectura de archivos.](/docs/es/practice/file-reader) y [TCP Práctica](/docs/es/practice/tcp-client).

## void y never

Una función normal que omite un tipo de retorno puede regresar al punto de llamada sin un valor. El tipo never se escribe como `!`, lo que significa que no regresa al pulsador normalmente. Un ejemplo representativo es la función de terminación del proceso.

Ejemplo que ilustra la declaración:

```text
fun log(message: str)              // 작업 후 돌아옴, 결과값 없음
fun stop(code: i32) -> !           // 정상 반환하지 않음
async fun work() -> i64            // 호출 결과는 Future<i64>
```

No se intenta hacer de never un valor almacenado común. No escriba funciones declaradas como sin retorno para que tengan una ruta de retorno normal. Si se requiere la limpieza de recursos antes de salir, la persona que llama debe hacerlo primero.

## Ejemplo completo de una función que no regresa

Si lo guarda como main.wave y lo ejecuta, termina con el código de salida 0 sin salida. stop no regresa a la persona que llama, por lo que se declara como `-> !`.

<!-- wave-example: never-function -->
```wave
import("std::process::core")::{
    proc_exit
};

fun stop() -> ! {
    proc_exit(0);
}

fun main() {
    stop();
}
```
