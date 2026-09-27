---
translation_set_id: stdlib-math-debug
path: stdlib/math-debug
locale: vi
group: stdlib
group_order: 1
order: 14
title: math và debug: Tiện ích toán học và chẩn đoán
summary: Sử dụng các hàm toán học với đầu ra chẩn đoán và kiểm tra phạm vi.
---

## Hàm số nguyên để kiểm tra phạm vi

```text
std::math::int
min_i32(a: i32, b: i32) -> i32
max_i32(a: i32, b: i32) -> i32
abs_i32_checked(x: i32) -> MathResult<i32>
clamp_i32_checked(x: i32, lo: i32, hi: i32) -> MathResult<i32>
div_floor_i32_checked(a: i32, b: i32) -> MathResult<i32>
div_ceil_i32_checked(a: i32, b: i32) -> MathResult<i32>
```

`MathResult<T>` chứa `value` và `error`. Nhập các hằng lỗi từ `std::math::result` và chỉ sử dụng `value` sau khi kiểm tra `error == MATH_ERROR_NONE`. Trị tuyệt đối của số nguyên có dấu nhỏ nhất không thể biểu diễn bằng cùng kiểu, nên hàm trị tuyệt đối có kiểm tra sẽ báo lỗi. `clamp` từ chối trường hợp `lo > hi`. Phép chia kiểm tra số chia bằng 0 và kết quả vượt phạm vi. `floor` và `ceil` làm tròn khác với phép chia số nguyên, vốn bỏ phần thập phân theo hướng về 0.

## Phân loại giá trị dấu phẩy động

`is_nan_f64`, `is_infinite_f64` và `is_finite_f64` trong `std::math::float` phân biệt các giá trị đặc biệt. Chức năng f32 cũng được cung cấp. `float_to_bits_f64(value: f64) -> u64` là chức năng lấy các bit lưu trữ và khác với chức năng chuyển đổi số của `value as u64`. NaN thậm chí còn không bằng chính nó nên không được kiểm tra với `value == nan`.

## chẩn đoán

`debug_assert(condition: bool, message: str)` trong `std::debug::core` chấm dứt sau khi chẩn đoán tình trạng sai. Các tình huống thường có thể không thành công, chẳng hạn như đầu vào của người dùng, được xử lý bằng giá trị kết quả và assert được sử dụng khi kiểm tra các điều kiện chương trình nội bộ phải được đáp ứng.

Lưu chương trình bên dưới dưới dạng `main.wave` và chạy nó.

<!-- wave-example: math-api -->
```wave
import("std::math::int")::{
    min_i32, max_i32
};
import("std::debug::core")::{
    debug_assert
};

fun main() {
    var low: i32 = min_i32(9, 4);
    var high: i32 = max_i32(9, 4);
    debug_assert(low <= high, "invalid range");
    println("{} {}", low, high);
}
```

Kết quả thực hiện:

```text
4 9
```

## Hướng làm tròn cho phép chia âm

So sánh cách chia -7 cho 3. `/` bị rút gọn về 0 thành -2. floor chọn số nguyên nhỏ hơn -3 và ceil chọn số nguyên lớn hơn -2.

<!-- wave-example: book-math-rounding -->
```wave
import("std::math::int")::{
    div_floor_i32_checked,
    div_ceil_i32_checked,
    abs_i32_checked
};
import("std::math::result")::{
    MathResult,
    MATH_ERROR_NONE,
    MATH_ERROR_OVERFLOW
};

fun main() -> i32 {
    var floor: MathResult<i32> = div_floor_i32_checked(-7, 3);
    var ceil: MathResult<i32> = div_ceil_i32_checked(-7, 3);

    if (floor.error != MATH_ERROR_NONE || ceil.error != MATH_ERROR_NONE) {
        return 1;
    }

    println("truncate={} floor={} ceil={}", -7 / 3, floor.value, ceil.value);

    var absolute: MathResult<i32> = abs_i32_checked(-2147483648);

    if (absolute.error == MATH_ERROR_OVERFLOW) {
        println("absolute value is outside i32");
    }

    return 0;
}
```

Kết quả thực hiện:

```text
truncate=-2 floor=-3 ceil=-2
absolute value is outside i32
```

Nếu chỉ kiểm tra giá trị kết quả, bạn không thể phân biệt giữa giá trị thay thế được đưa vào trong trường hợp thất bại và kết quả tính toán thực tế. Thực hiện theo thứ tự kiểm tra error trước. floor rất hữu ích khi đặt tọa độ âm vào một khoảng có kích thước nhất định và ceil rất hữu ích khi làm tròn số bó được yêu cầu.

## Tôi nên xử lý những lỗi nào?

|tình huống|lỗi|Ví dụ xử lý|
| --- | --- | --- |
|Chia cho số 0| `MATH_ERROR_DIVIDE_BY_ZERO` |Nhập lại mẫu số|
|Kết quả không thể được lưu trữ trong loại| `MATH_ERROR_OVERFLOW` |Tính toán với loại rộng hơn hoặc từ chối đầu vào|
|Mức tối thiểu lớn hơn mức tối đa trong clamp| `MATH_ERROR_INVALID_ARGUMENT` |Sửa đổi phạm vi cài đặt|

assert không phải là công cụ sửa lỗi. Lỗi trong dữ liệu đầu vào của người dùng được xử lý bằng cách sử dụng các câu lệnh có điều kiện và giá trị trả về, đồng thời sau khi hoàn thành phép tính, các điều kiện nội bộ phải đáp ứng sẽ được kiểm tra bằng debug_assert.
