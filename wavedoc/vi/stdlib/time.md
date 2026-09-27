---
translation_set_id: stdlib-time
path: stdlib/time
locale: vi
group: stdlib
group_order: 1
order: 11
title: time: Thời lượng, đo lường và chờ đợi
summary: Giải thích sự khác biệt giữa đơn vị của Duration và realtime·monotonic clock.
---

## Duration

`Duration` trong `std::time::duration` có seconds và nanoseconds. Phạm vi nanoseconds được chuẩn hóa là từ 0 đến 999999999. 1000 mili giây bằng 1 giây.

```text
time_duration_from_ms(milliseconds: i64) -> Duration
time_duration_from_seconds(seconds: i64) -> Duration
time_duration_checked_add(left: Duration, right: Duration) -> DurationResult
time_duration_to_ns(value: Duration) -> DurationValueResult
```

checked kết quả hoạt động và chuyển đổi số nguyên bao gồm ok và value. Việc thay thế Duration rộng bằng một giá trị nano giây i64 có ​​thể nằm ngoài phạm vi, vì vậy trước tiên hãy kiểm tra ok.

## Phân biệt giữa đo lường và tầm nhìn

`std::time::clock` đến `time_now_realtime(tp: ptr<TimeSpec>) -> i64` tương ứng với thời gian theo lịch. Sử dụng `time_now_monotonic` để đo thời gian đã trôi qua vì điều này có thể thay đổi khi điều chỉnh đồng hồ hệ thống. Cửa hàng đầu ra do người gọi cung cấp và chỉ đọc sec/nsec nếu trạng thái thành công.

```text
std::time::sleep
time_sleep(duration: Duration) -> i64
time_sleep_ms(ms: i64) -> i64
time_sleep_us(us: i64) -> i64
time_sleep_ns(ns: i64) -> i64
time_sleep_seconds(seconds: i64) -> i64
```

Chờ âm là lỗi và 0 là thành công ngay lập tức. Sau interruption, chỉ còn lại thời gian cho đến monotonic deadline sẽ được chờ. Vì thời gian thức dậy thực tế có thể bị trì hoãn do lập kế hoạch nên nó không được sử dụng như một chức năng đảm bảo thời gian thực hiện chính xác. Việc gọi đồng bộ sleep trong tác vụ async có thể cản trở tiến trình của người thực thi, vì vậy hãy chọn `task::sleep_ms`.

## Ví dụ chuyển đổi đơn vị

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

Kết quả thực hiện:

```text
1 500000000
1500000000
```

## Thêm thời gian và thay đổi đơn vị

750 mili giây và 800 mili giây cộng lại thành 1 giây và 550000000 nano giây. Thay vì thêm giây và nano giây riêng biệt, bạn có thể sử dụng checked_add để kiểm tra mức độ mang và phạm vi cùng nhau.

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

Kết quả thực hiện:

```text
1s 550000000ns
1550ms
```

## khoảng thời gian âm

Duration cũng đại diện cho số âm. -1ms chuẩn hóa thành seconds=-1, nanoseconds=999000000. Hai trường kết hợp là âm 1 mili giây. Bạn không nên đánh giá rằng đó là số dương chỉ bằng cách nhìn vào trường nanoseconds.

Số âm hợp lệ trong tính toán khoảng thời gian, nhưng việc chuyển số âm tới sleep là một lỗi. Khi tính thời gian chờ còn lại, nếu đã qua thời hạn đó thì quá trình xử lý tiếp theo sẽ được tiến hành mà không phải chờ đợi.

## Chọn thời gian API

|mục đích|chọn|Kết quả có ý nghĩa gì|
| --- | --- | --- |
|Thời gian trôi qua giữa hai thời điểm| monotonic clock |Khoảng thời gian độc lập với hiệu chỉnh hình ảnh của hệ thống|
|thời gian thực tế theo lịch| realtime clock |Thời gian do hệ thống đặt|
|Chờ chương trình đồng bộ| `time_sleep_ms` |Luồng cuộc gọi được xếp hàng đợi|
|async Đang chờ nhiệm vụ| `task::sleep_ms` |Chuyển cơ hội thực hiện sang nhiệm vụ khác|

Ngay cả đồng hồ trả về nano giây cũng không có nghĩa là độ chính xác đo thực tế là 1 nano giây. Khi so sánh hiệu suất, hãy đo tổng thời gian lặp lại một nhiệm vụ ngắn nhiều lần và di chuyển các nhiệm vụ không liên quan đến mục tiêu đo lường, chẳng hạn như đầu vào/đầu ra, ra khỏi phần.
