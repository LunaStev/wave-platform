---
translation_set_id: stdlib-time
path: stdlib/time
locale: ru
group: stdlib
group_order: 1
order: 11
title: time: продолжительность, измерение и ожидание
summary: Объясните разницу между единицами измерения Duration и realtime·monotonic clock.
---

## Duration

`Duration` в `std::time::duration` имеет seconds и nanoseconds. Нормализованный диапазон nanoseconds составляет от 0 до 999999999. 1000 миллисекунд равны 1 секунде.

```text
time_duration_from_ms(milliseconds: i64) -> Duration
time_duration_from_seconds(seconds: i64) -> Duration
time_duration_checked_add(left: Duration, right: Duration) -> DurationResult
time_duration_to_ns(value: Duration) -> DurationValueResult
```

Результаты операции checked и целочисленного преобразования включают ok и value. Замена широкого Duration на одно наносекундное значение i64 может выйти за пределы допустимого диапазона, поэтому сначала проверьте ok.

## Различие между измерением и видением

От `std::time::clock` до `time_now_realtime(tp: ptr<TimeSpec>) -> i64` соответствуют календарному времени. Используйте `time_now_monotonic` для измерения прошедшего времени, так как оно может измениться при поправках системных часов. Выходное хранилище предоставляется вызывающей стороной и читается как sec/nsec только в том случае, если статус успешен.

```text
std::time::sleep
time_sleep(duration: Duration) -> i64
time_sleep_ms(ms: i64) -> i64
time_sleep_us(us: i64) -> i64
time_sleep_ns(ns: i64) -> i64
time_sleep_seconds(seconds: i64) -> i64
```

Отрицательное значение ожидания означает ошибку, а 0 — немедленный успех. После interruption будет ждаться только время, оставшееся до monotonic deadline. Поскольку фактическое время пробуждения может быть задержано из-за планирования, оно не используется как функция, гарантирующая точное время выполнения. Вызов синхронного sleep в задаче async может помешать работе исполнителя, поэтому выберите `task::sleep_ms`.

## Пример преобразования единиц измерения

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

Результат выполнения:

```text
1 500000000
1500000000
```

## Добавьте время и измените единицы измерения

750 мс и 800 мс в сумме дают 1 секунду и 550000000 наносекунд. Вместо добавления секунд и наносекунд по отдельности вы можете использовать checked_add, чтобы проверить перенос и дальность вместе.

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

Результат выполнения:

```text
1s 550000000ns
1550ms
```

## отрицательный интервал времени

Duration также представляет отрицательные числа. -1 мс нормализуется до seconds=-1, nanoseconds=999000000. Оба поля вместе имеют отрицательную величину в 1 миллисекунду. Не следует судить о том, что это положительное число, просто взглянув на поле nanoseconds.

Отрицательные числа допустимы при расчете временных интервалов, но передача отрицательного числа в sleep является ошибкой. При расчете оставшегося времени ожидания, если срок уже прошел, следующая обработка продолжится без ожидания.

## Выберите время API

|цель|выбрать|Что означают результаты|
| --- | --- | --- |
|Время, прошедшее между двумя моментами времени| monotonic clock |Интервал, независимый от системы визуальной коррекции|
|фактическое календарное время| realtime clock |Время, установленное системой|
|Ожидание синхронной программы| `time_sleep_ms` |Поток вызовов поставлен в очередь|
|async Ожидание задания| `task::sleep_ms` |Передача возможности выполнения другой задаче|

Даже часы, которые показывают наносекунды, не означают, что фактическая точность измерения составляет 1 наносекунду. При сравнении производительности измерьте общее время многократного повторения короткой задачи и переместите из раздела задачи, не связанные с целью измерения, например ввод/вывод.
