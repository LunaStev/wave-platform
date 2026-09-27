---
translation_set_id: stdlib-time
path: stdlib/time
locale: zh
group: stdlib
group_order: 1
order: 11
title: time:持续时间、测量和等待
summary: 解释一下Duration和realtime·monotonicclock的单位区别。
---

## Duration

`std::time::duration` 中的`Duration` 有 seconds 和 nanoseconds。标准化的nanoseconds范围是0到999999999。1000毫秒等于1秒。

```text
time_duration_from_ms(milliseconds: i64) -> Duration
time_duration_from_seconds(seconds: i64) -> Duration
time_duration_checked_add(left: Duration, right: Duration) -> DurationResult
time_duration_to_ns(value: Duration) -> DurationValueResult
```

checked运算和整数转换结果包括ok和value。用单个 i64 纳秒值替换宽 Duration 可能会超出范围，因此请先检查 ok。

## 区分测量和视觉

`std::time::clock` 到 `time_now_realtime(tp: ptr<TimeSpec>) -> i64` 对应于日历时间。使用 `time_now_monotonic` 进行经过时间测量，因为这可能会随着系统时钟校正而改变。输出存储由调用者提供，并且仅在状态成功时才读取sec/nsec。

```text
std::time::sleep
time_sleep(duration: Duration) -> i64
time_sleep_ms(ms: i64) -> i64
time_sleep_us(us: i64) -> i64
time_sleep_ns(ns: i64) -> i64
time_sleep_seconds(seconds: i64) -> i64
```

负等待是一个错误，0 是立即成功。 interruption之后，仅等待到monotonicdeadline为止的剩余时间。由于实际唤醒时间可能因调度而延迟，因此不作为保证精确执行时间的函数。在async任务中调用同步sleep会阻止执行器进度，因此选择`task::sleep_ms`。

## 单位换算示例

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

执行结果：

```text
1 500000000
1500000000
```

## 添加时间并更改单位

750ms 和 800ms 加起来就是 1 秒和 550000000 纳秒。您可以使用 checked_add 一起检查进位和范围，而不是分别添加秒和纳秒。

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

执行结果：

```text
1s 550000000ns
1550ms
```

## 负时间间隔

Duration也代表负数。 -1ms 归一化为 seconds=-1，nanoseconds=999000000。两个字段加起来为负 1 毫秒。您不应该仅通过查看nanoseconds字段来判断它是正数。

负数在时间间隔计算中是有效的，但将负数传递给sleep是错误的。在计算剩余等待时间时，如果已经过了截止时间，则不等待而进行下一个处理。

## 选择时间API

|目的|选择|结果意味着什么|
| --- | --- | --- |
|两个时间点之间经过的时间| monotonic clock |间隔与系统视觉校正无关|
|实际日历时间| realtime clock |系统设定的时间|
|等待同步程序| `time_sleep_ms` |呼叫流程已排队|
|async 等待任务| `task::sleep_ms` |将执行机会传递给另一个任务|

即使是返回纳秒的时钟也不意味着实际测量精度为1纳秒。在比较性能时，测量重复多次短任务的总时间，并将与测量目标无关的任务（例如输入/输出）移出该部分。
