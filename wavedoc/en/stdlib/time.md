---
translation_set_id: stdlib-time
path: stdlib/time
locale: en
group: stdlib
group_order: 1
order: 11
title: time: Durations, measurement, and waiting
summary: Explain the difference between the units of Duration and realtime·monotonic clock.
---

## Duration

`Duration` in `std::time::duration` has seconds and nanoseconds. The normalized nanoseconds range is 0 to 999999999. 1000 milliseconds is equal to 1 second.

```text
time_duration_from_ms(milliseconds: i64) -> Duration
time_duration_from_seconds(seconds: i64) -> Duration
time_duration_checked_add(left: Duration, right: Duration) -> DurationResult
time_duration_to_ns(value: Duration) -> DurationValueResult
```

checked operation and integer conversion results include ok and value. Replacing the wide Duration with a single i64 nanosecond value may be out of range, so check ok first.

## Distinguish between measurement and vision

`std::time::clock` through `time_now_realtime(tp: ptr<TimeSpec>) -> i64` correspond to calendar times. Use `time_now_monotonic` for elapsed time measurements, as this may change with system clock corrections. The output store is provided by the caller and only reads sec/nsec if the status is success.

```text
std::time::sleep
time_sleep(duration: Duration) -> i64
time_sleep_ms(ms: i64) -> i64
time_sleep_us(us: i64) -> i64
time_sleep_ns(ns: i64) -> i64
time_sleep_seconds(seconds: i64) -> i64
```

A negative wait is an error, and 0 is an immediate success. After interruption, only the time remaining until monotonic deadline will be waited. Because the actual wake-up time may be delayed due to scheduling, it is not used as a function that guarantees precise execution time. Calling synchronous sleep within the async task can prevent executor progress, so choose `task::sleep_ms`.

## Unit conversion example

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

Execution result:

```text
1 500000000
1500000000
```

## Add time and change units

750ms and 800ms add up to 1 second and 550000000 nanoseconds. Instead of adding seconds and nanoseconds separately, you can use checked_add to check carry and range together.

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

Execution result:

```text
1s 550000000ns
1550ms
```

## negative time interval

Duration also represents negative numbers. -1ms normalizes to seconds=-1, nanoseconds=999000000. The two fields combined are negative 1 millisecond. You should not judge that it is a positive number just by looking at the nanoseconds field.

Negative numbers are valid in time interval calculations, but passing a negative number to sleep is an error. When calculating the remaining waiting time, if the deadline has already passed, the next processing will proceed without waiting.

## Select time API

|purpose|select|What the results mean|
| --- | --- | --- |
|Time elapsed between two points in time| monotonic clock |Interval independent of system visual correction|
|actual calendar time| realtime clock |Time set by the system|
|Waiting for synchronous program| `time_sleep_ms` |Call flow is queued|
|async Waiting for task| `task::sleep_ms` |Passing execution opportunity to another task|

Even a clock that returns nanoseconds does not mean that the actual measurement precision is 1 nanosecond. When comparing performance, measure the total time for repeating a short task multiple times, and move tasks unrelated to the measurement target, such as input/output, out of the section.
