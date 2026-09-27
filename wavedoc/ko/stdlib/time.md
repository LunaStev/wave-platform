---
translation_set_id: stdlib-time
path: stdlib/time
locale: ko
group: stdlib
group_order: 1
order: 11
title: time: 시간 값·측정·대기
summary: Duration의 단위와 realtime·monotonic clock의 차이를 설명합니다.
---

## Duration

`std::time::duration`의 `Duration`은 seconds와 nanoseconds를 가집니다. 정규화된 nanoseconds 범위는 0부터 999999999까지입니다. 밀리초 1000은 1초입니다.

```text
time_duration_from_ms(milliseconds: i64) -> Duration
time_duration_from_seconds(seconds: i64) -> Duration
time_duration_checked_add(left: Duration, right: Duration) -> DurationResult
time_duration_to_ns(value: Duration) -> DurationValueResult
```

checked 연산과 정수 변환 결과에는 ok와 value가 있습니다. 넓은 Duration을 하나의 i64 나노초 값으로 바꾸면 범위를 넘을 수 있으므로 ok를 먼저 확인합니다.

## 측정과 시각 구분

`std::time::clock`의 `time_now_realtime(tp: ptr<TimeSpec>) -> i64`는 달력 시각에 대응합니다. 시스템 시각 보정으로 바뀔 수 있으므로 경과 시간 측정에는 `time_now_monotonic`을 사용합니다. 출력 저장소는 호출자가 제공하고 상태가 성공일 때만 sec/nsec를 읽습니다.

```text
std::time::sleep
time_sleep(duration: Duration) -> i64
time_sleep_ms(ms: i64) -> i64
time_sleep_us(us: i64) -> i64
time_sleep_ns(ns: i64) -> i64
time_sleep_seconds(seconds: i64) -> i64
```

음수 대기는 오류이며 0은 즉시 성공합니다. interruption 이후에는 monotonic deadline까지 남은 시간만 대기합니다. 스케줄링 때문에 실제 깨어나는 시점이 늦을 수 있으므로 정밀한 실행 시각을 보장하는 함수로 사용하지 않습니다. async 작업 안에서 동기 sleep을 호출하면 실행기 진행을 막을 수 있어 `task::sleep_ms`를 선택합니다.

## 단위 변환 예제

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

실행 결과:

```text
1 500000000
1500000000
```

## 시간을 더하고 단위를 바꾸기

750ms와 800ms를 더하면 1초 550000000나노초입니다. 초와 나노초를 각각 더하는 대신 checked_add를 사용하면 자리올림과 표현 범위를 함께 검사할 수 있습니다.

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

실행 결과:

```text
1s 550000000ns
1550ms
```

## 음수 시간 간격

Duration은 음수도 표현합니다. -1ms는 seconds=-1, nanoseconds=999000000으로 정규화합니다. 두 필드를 합친 값이 음수 1밀리초입니다. nanoseconds 필드만 보고 양수라고 판단하면 안 됩니다.

시간 간격 계산에서는 음수가 유효하지만, sleep에 음수를 넘기는 것은 오류입니다. 남은 대기 시간을 계산할 때 이미 마감 시각을 지났다면 대기하지 않고 다음 처리를 진행합니다.

## 시간 API 선택

| 목적 | 선택 | 결과의 의미 |
| --- | --- | --- |
| 두 시점 사이의 경과 시간 | monotonic clock | 시스템 시각 보정과 독립된 간격 |
| 실제 달력 시각 | realtime clock | 시스템이 설정한 시각 |
| 동기 프로그램의 대기 | `time_sleep_ms` | 호출 흐름이 대기 |
| async 작업의 대기 | `task::sleep_ms` | 다른 작업에 실행 기회를 넘김 |

나노초를 반환하는 시계라도 실제 측정 정밀도가 1나노초라는 뜻은 아닙니다. 성능을 비교할 때는 짧은 작업을 여러 번 반복한 전체 시간을 측정하고, 입출력처럼 측정 대상과 무관한 작업은 구간 밖으로 옮깁니다.
