---
translation_set_id: stdlib-time
path: stdlib/time
locale: es
group: stdlib
group_order: 1
order: 11
title: time: Duraciones, medición y espera.
summary: Explique la diferencia entre las unidades de Duration y realtime·monotonic clock.
---

## Duration

`Duration` en `std::time::duration` tiene seconds y nanoseconds. El rango normalizado de nanoseconds es de 0 a 999999999. 1000 milisegundos equivalen a 1 segundo.

```text
time_duration_from_ms(milliseconds: i64) -> Duration
time_duration_from_seconds(seconds: i64) -> Duration
time_duration_checked_add(left: Duration, right: Duration) -> DurationResult
time_duration_to_ns(value: Duration) -> DurationValueResult
```

Los resultados de la operación checked y de la conversión de enteros incluyen ok y value. Reemplazar el ancho Duration con un único valor de nanosegundos i64 puede estar fuera de rango, así que verifique ok primero.

## Distinguir entre medición y visión.

`std::time::clock` a `time_now_realtime(tp: ptr<TimeSpec>) -> i64` corresponden a tiempos del calendario. Utilice `time_now_monotonic` para mediciones de tiempo transcurrido, ya que esto puede cambiar con las correcciones del reloj del sistema. El almacén de resultados lo proporciona la persona que llama y solo lee sec/nsec si el estado es exitoso.

```text
std::time::sleep
time_sleep(duration: Duration) -> i64
time_sleep_ms(ms: i64) -> i64
time_sleep_us(us: i64) -> i64
time_sleep_ns(ns: i64) -> i64
time_sleep_seconds(seconds: i64) -> i64
```

Una espera negativa es un error y 0 es un éxito inmediato. Después de interruption, sólo se esperará el tiempo restante hasta monotonic deadline. Debido a que la hora de despertar real puede retrasarse debido a la programación, no se utiliza como una función que garantice un tiempo de ejecución preciso. Llamar sincrónico a sleep dentro de la tarea async puede impedir el progreso del ejecutor, así que elija `task::sleep_ms`.

## Ejemplo de conversión de unidades

<!-- wave-example: duration-api -->
```wave
import("std::time::duration")::{
    Duration, DurationValueResult, time_duration_from_ms, time_duration_to_ns
};

fun main() -> i32 {
    var duration: Duration = time_duration_from_ms(1500);
    var value: DurationValueResult = time_duration_to_ns(duration);
    if (!value.ok) {
        return 1;
    }

    println("{} {}", duration.seconds, duration.nanoseconds);
    println("{}", value.value);
    return 0;
}
```

Resultado de la ejecución:

```text
1 500000000
1500000000
```

## Agregar tiempo y cambiar unidades

750 ms y 800 ms suman 1 segundo y 550000000 nanosegundos. En lugar de sumar segundos y nanosegundos por separado, puede usar checked_add para verificar el transporte y el alcance juntos.

<!-- wave-example: book-duration-add -->
```wave
import("std::time::duration")::{
    Duration,
    DurationResult,
    DurationValueResult,
    time_duration_from_ms,
    time_duration_checked_add,
    time_duration_to_ms
};

fun main() -> i32 {
    var first: Duration = time_duration_from_ms(750);
    var second: Duration = time_duration_from_ms(800);
    var sum: DurationResult = time_duration_checked_add(first, second);

    if (!sum.ok) {
        return 1;
    }

    var milliseconds: DurationValueResult = time_duration_to_ms(sum.value);

    if (!milliseconds.ok) {
        return 2;
    }

    println("{}s {}ns", sum.value.seconds, sum.value.nanoseconds);
    println("{}ms", milliseconds.value);
    return 0;
}
```

Resultado de la ejecución:

```text
1s 550000000ns
1550ms
```

## intervalo de tiempo negativo

Duration también representa números negativos. -1ms se normaliza a seconds=-1, nanoseconds=999000000. Los dos campos combinados son 1 milisegundo negativo. No debes juzgar que es un número positivo con solo mirar el campo nanoseconds.

Los números negativos son válidos en los cálculos de intervalos de tiempo, pero pasar un número negativo a sleep es un error. A la hora de calcular el tiempo de espera restante, si el plazo ya ha transcurrido, se procederá al siguiente trámite sin esperas.

## Seleccione hora API

|propósito|seleccionar|¿Qué significan los resultados?|
| --- | --- | --- |
|Tiempo transcurrido entre dos puntos en el tiempo| monotonic clock |Intervalo independiente de la corrección visual del sistema.|
|tiempo real del calendario| realtime clock |Hora establecida por el sistema.|
|Esperando programa sincrónico| `time_sleep_ms` |El flujo de llamadas está en cola|
|async Esperando tarea| `task::sleep_ms` |Pasar la oportunidad de ejecución a otra tarea|

Incluso un reloj que devuelve nanosegundos no significa que la precisión real de la medición sea de 1 nanosegundo. Al comparar el rendimiento, mida el tiempo total para repetir una tarea corta varias veces y mueva las tareas no relacionadas con el objetivo de medición, como entrada/salida, fuera de la sección.
