---
translation_set_id: stdlib-time
path: stdlib/time
locale: ja
group: stdlib
group_order: 1
order: 11
title: time: 継続時間、測定、待機時間
summary: Durationの単位とrealtime・monotonic clockの違いを説明します。
---

## Duration

`std::time::duration`の`Duration`にはsecondsとnanosecondsがあります。正規化されたnanosecondsの範囲は0から999999999までです。ミリ秒1000は1秒です。

```text
time_duration_from_ms(milliseconds: i64) -> Duration
time_duration_from_seconds(seconds: i64) -> Duration
time_duration_checked_add(left: Duration, right: Duration) -> DurationResult
time_duration_to_ns(value: Duration) -> DurationValueResult
```

checked演算と整数変換結果にはokとvalueがあります。広いDurationを1つのi64ナノ秒の値に変更すると範囲を超える可能性があるため、okを最初に確認します。

## 測定と視覚の区別

`std::time::clock`の`time_now_realtime(tp: ptr<TimeSpec>) -> i64`はカレンダー時刻に対応します。システム時刻補正に変更できるため、経過時間測定には`time_now_monotonic`を使用します。出力リポジトリは、呼び出し元によって提供され、ステータスが成功した場合にのみsec/nsecを読み込みます。

```text
std::time::sleep
time_sleep(duration: Duration) -> i64
time_sleep_ms(ms: i64) -> i64
time_sleep_us(us: i64) -> i64
time_sleep_ns(ns: i64) -> i64
time_sleep_seconds(seconds: i64) -> i64
```

負の待機はエラーで、ゼロはすぐに成功します。 interruption以降はmonotonicdeadlineまでの残り時間だけ待機します。スケジューリングのため、実際に目が覚める時間が遅れる可能性があるため、正確な実行時間を保証する関数としては使用しません。 asyncタスク内で同期sleepを呼び出すとランチャーの進行を防ぐことができ、`task::sleep_ms`を選択します。

## 単位変換の例

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

実行結果：

```text
1 500000000
1500000000
```

## 時間を加えて単位を変える

750msと800msを加えると1秒550000000ナノ秒です。秒とナノ秒をそれぞれ追加する代わりに、checked_addを使用すると、桁上げと表現範囲を一緒に調べることができます。

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

実行結果：

```text
1s 550000000ns
1550ms
```

## 負の時間間隔

Durationは負の数も表現します。 -1msはseconds=-1、nanoseconds=999000000に正規化します。 2つのフィールドを合わせた値は負の1ミリ秒です。 nanoseconds フィールドのみ見て正数と判断してはいけません。

時間間隔計算では負数が有効ですが、sleepに負数を渡すのはエラーです。残りの待ち時間を計算する際に、すでに終了時刻を過ぎている場合は待機せずに次の処理に進む。

## 時間API選択

|目的|選択|結果の意味|
| --- | --- | --- |
|2つの視点間の経過時間| monotonic clock |システム時刻補正と独立した間隔|
|実際のカレンダー時刻| realtime clock |システムが設定した時刻|
|同期プログラムの待機| `time_sleep_ms` |コールフローが待機|
|asyncジョブの待機| `task::sleep_ms` |他のタスクに実行する機会を引き渡す|

ナノ秒を返す時計でも実際の測定精度が1ナノ秒という意味ではありません。性能を比較するときは、短い操作を何度も繰り返した全時間を測定し、入出力のように測定対象とは無関係な作業は区間外に移動します。
