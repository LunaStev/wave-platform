---
translation_set_id: stdlib-buffer
path: stdlib/buffer
locale: vi
group: stdlib
group_order: 1
order: 5
title: buffer: Bộ nhớ byte có thể phát triển
summary: Buffer Mô tả các quy tắc khởi tạo, bổ sung, truy vấn, dung lượng và giải phóng.
---

## Ý nghĩa của Buffer

`Buffer` trong `std::buffer::types` có `data: ptr<u8>`, `len: i64` và `cap: i64`. len là số byte được khởi tạo và sử dụng, và cap là tổng số byte được phân bổ. Luôn duy trì `0 <= len <= cap`. Chuỗi NUL không tự động đảm bảo chấm dứt.

## Cơ bản API

|mô-đun|khai báo|ý nghĩa|
| --- | --- | --- |
| `std::buffer::alloc` | `buffer_init(out_buf: ptr<Buffer>, capacity: i64) -> i64` |Khởi tạo kho lưu trữ mới. Không thu hồi vào bộ đệm đã sở hữu|
|cùng một mô-đun| `buffer_free(buf: ptr<Buffer>) -> i64` |Phân bổ. Trống nếu thành công|
|cùng một mô-đun| `buffer_reserve(buf: ptr<Buffer>, required_cap: i64) -> i64` |Đảm bảo đầy đủ công suất tối thiểu. len Đã bảo trì|
|cùng một mô-đun| `buffer_clear(buf: ptr<Buffer>) -> i64` |Duy trì công suất và len=0|
|cùng một mô-đun| `buffer_resize(buf: ptr<Buffer>, new_len: i64, value: u8) -> i64` |Thay đổi độ dài, điền byte mới bằng value|
| `std::buffer::write` | `buffer_push(buf: ptr<Buffer>, value: u8) -> i64` |thêm một byte|
|cùng một mô-đun| `buffer_append(buf: ptr<Buffer>, src: ptr<u8>, size: i64) -> i64` |sizeSao chép và thêm byte|
|cùng một mô-đun| `buffer_append_str(buf: ptr<Buffer>, s: str) -> i64` |Thêm byte chuỗi ngoại trừ NUL|
| `std::buffer::read` | `buffer_get(buf: ptr<Buffer>, index: i64, out_value: ptr<u8>) -> i64` |Đọc một byte trong phạm vi|

Trả về trạng thái API trả về `BUFFER_OK`(0) khi thành công. Phân biệt các lỗi INVALID, BOUNDS, OVERFLOW, và ALLOC từ `std::buffer::error`. Con số này không được hiểu là OS errno. `buffer_new` thể hiện lỗi phân bổ dưới dạng trống Buffer, vì vậy khi bạn cần phân biệt giữa các lỗi, hãy sử dụng `buffer_init`.

## Ví dụ đang chạy

<!-- wave-example: buffer-api -->
```wave
import("std::buffer::alloc")::{
    Buffer, buffer_init, buffer_free
};
import("std::buffer::write")::{
    buffer_append_str, buffer_push
};
import("std::buffer::read")::{
    buffer_get
};

fun main() -> i32 {
    var data: Buffer;
    if (buffer_init(&data, 0) < 0) {
        return 1;
    }
    if (buffer_append_str(&data, "Hi") < 0 || buffer_push(&data, 33) < 0) {
        buffer_free(&data);
        return 2;
    }

    var value: u8 = 0;
    if (buffer_get(&data, 2, &value) < 0) {
        buffer_free(&data);
        return 3;
    }

    println("{} {}", data.len, value);
    if (buffer_free(&data) < 0) {
        return 4;
    }

    return 0;
}
```

Kết quả thực hiện:

```text
3 33
```

Lưu nó dưới dạng `main.wave` và chạy nó dưới dạng `wavec run main.wave`. Dung lượng ban đầu bằng 0 không phải là lỗi mà là bộ đệm trống hợp lệ. Giải phóng không gian trong quá trình xử lý tiếp theo.

## Tuổi thọ và thất bại

Việc tăng bộ đệm có thể thay đổi dữ liệu. Không sử dụng địa chỉ đã mượn trước đó sau một hoạt động có thể phân bổ lại. Việc sao chép cấu trúc Buffer sẽ không trùng lặp quá trình phân bổ của nó, vì vậy hãy cấp phân bổ đó cho một chủ sở hữu duy nhất.

`buffer_get` không thay đổi đối số đầu ra nếu thất bại. Mặt khác, hàm tiện lợi `buffer_at` cũng biểu thị lỗi là 0, vì vậy hãy sử dụng `buffer_get` để phân biệt giữa 0 byte thực tế và lỗi. Tránh tạo len/cap không hợp lệ bằng cách thay đổi trực tiếp các trường công khai.

[Ký ức API](/docs/vi/reference/memory-and-buffer) · [Luyện đọc file dưới dạng Buffer](/docs/vi/practice/file-reader)

## Quan sát chiều dài và công suất riêng biệt

reserve giải phóng không gian lưu trữ nhưng không tăng len. resize thay đổi độ dài thực tế được sử dụng và khởi tạo phần mở rộng thành byte được chỉ định. clear chỉ đặt độ dài được sử dụng thành 0, cho phép sử dụng lại phân bổ.

Lưu chương trình sau dưới dạng main.wave và chạy nó. Nó không dựa vào bội số tăng trưởng chính xác của công suất; nó chỉ đảm bảo rằng bạn có không gian bạn cần.

<!-- wave-example: book-buffer-length-capacity -->
```wave
import("std::buffer::alloc")::{
    Buffer,
    buffer_init,
    buffer_reserve,
    buffer_resize,
    buffer_clear,
    buffer_free
};

fun main() -> i32 {
    var data: Buffer;

    if (buffer_init(&data, 0) < 0) {
        return 1;
    }

    if (buffer_reserve(&data, 20) < 0) {
        buffer_free(&data);
        return 2;
    }

    println("reserved length={}", data.len);

    if (buffer_resize(&data, 3, 7) < 0) {
        buffer_free(&data);
        return 3;
    }

    println("resized length={} first={}", data.len, deref data.data[0]);

    if (buffer_clear(&data) < 0) {
        buffer_free(&data);
        return 4;
    }

    println("cleared length={}", data.len);

    if (buffer_free(&data) < 0) {
        return 5;
    }

    return 0;
}
```

Kết quả thực hiện:

```text
reserved length=0
resized length=3 first=7
cleared length=0
```

Khi tăng lên resize, chúng tôi đã chuyển value=7, vì vậy tất cả ba byte mới mà chúng tôi thấy là 7. Không gian chỉ được bảo mật bằng reserve không được đọc dưới dạng dữ liệu khởi tạo. cap sẽ vẫn còn sau clear và có thể được thêm lại vào cùng Buffer.

## Phân biệt giữa byte 0 và lỗi tra cứu

buffer_get trả về trạng thái và ghi byte thực tế làm đối số đầu ra. Ngay cả khi dữ liệu là 0 thì đó là thành công bình thường.

<!-- wave-example: book-buffer-zero-bounds -->
```wave
import("std::buffer::alloc")::{Buffer, buffer_init, buffer_free};
import("std::buffer::write")::{buffer_push};
import("std::buffer::read")::{buffer_get};
import("std::buffer::error")::{BUFFER_ERR_BOUNDS};

fun main() -> i32 {
    var data: Buffer;

    if (buffer_init(&data, 0) < 0) {
        return 1;
    }

    if (buffer_push(&data, 0) < 0) {
        buffer_free(&data);
        return 2;
    }

    var value: u8 = 99;

    if (buffer_get(&data, 0, &value) < 0) {
        buffer_free(&data);
        return 3;
    }

    println("stored={}", value);
    value = 99;

    if (buffer_get(&data, 1, &value) == BUFFER_ERR_BOUNDS) {
        println("outside, preserved={}", value);
    }

    if (buffer_free(&data) < 0) {
        return 4;
    }

    return 0;
}
```

Kết quả thực hiện:

```text
stored=0
outside, preserved=99
```

Lần truy cập đầu tiên là thành công khi đọc 0, lần truy cập thứ hai là thất bại ngoài giới hạn. Ngay cả khi value=99 vẫn còn sau một lỗi, điều đó không có nghĩa đó là giá trị được đọc từ bộ đệm. Hãy chắc chắn kiểm tra trạng thái cùng nhau.

## Giải pháp thực hành: Tích lũy byte

Để cộng các số từ 0 đến 9, hãy lặp lại buffer_push và kiểm tra từng kết quả. Lưu trữ tổng trong i64 và chỉ đọc phạm vi `0 <= index < data.len`. Sau khi xử lý bộ đệm, chúng tôi gọi buffer_free trên cả đường dẫn thành công và thất bại.

<!-- wave-example: book-buffer-sum -->
```wave
import("std::buffer::alloc")::{Buffer, buffer_init, buffer_free};
import("std::buffer::write")::{buffer_push};

fun main() -> i32 {
    var data: Buffer;

    if (buffer_init(&data, 0) < 0) {
        return 1;
    }

    for (var value: i32 = 0; value < 10; value += 1) {
        if (buffer_push(&data, value as u8) < 0) {
            buffer_free(&data);
            return 2;
        }
    }

    var total: i64 = 0;

    for (var index: i64 = 0; index < data.len; index += 1) {
        total += deref data.data[index] as i64;
    }

    println("sum={}", total);

    if (buffer_free(&data) < 0) {
        return 3;
    }

    return 0;
}
```

Kết quả thực hiện:

```text
sum=45
```
