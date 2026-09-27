---
translation_set_id: standard-library
path: reference/standard-library
locale: vi
group: stdlib
group_order: 1
order: 1
title: Hướng dẫn thư viện chuẩn
summary: Cách tìm module phù hợp với mục đích của bạn và đọc các lỗi cũng như quy tắc sở hữu của hàm.
---

## Tìm các tính năng bạn cần

Thư viện tiêu chuẩn là import với đường dẫn `std::module::file`. Ngay cả khi tên giống nhau, các hàm có thể trả về lỗi theo những cách khác nhau. Trước tiên hãy đọc [API Cách đọc](/docs/vi/stdlib/contracts), sau đó chuyển đến mô-đun bạn cần trong bảng sau.

|Những gì tôi muốn làm|tài liệu|Chính import|
| --- | --- | --- |
|Độ dài chuỗi/so sánh/tìm kiếm| [string](/docs/vi/reference/string-and-bytes) | `std::string::len`, `cmp`, `find`, `trim` |
|Phân bổ/sao chép/kích thước bộ nhớ| [mem](/docs/vi/reference/memory-and-buffer) | `std::mem::alloc`, `ops`, `layout` |
|Danh sách byte có kích thước khác nhau| [buffer](/docs/vi/stdlib/buffer) | `std::buffer::alloc`, `read`, `write` |
|Đọc/ghi nhị phân| [bytes](/docs/vi/stdlib/bytes) | `std::bytes::types`, `cursor`, `leb128` |
|Tệp·Mô tả I/O| [fs và io](/docs/vi/stdlib/files-io) | `std::fs::file`, `std::io::fd` |
|Cài đặt môi trường/kết hợp đường dẫn| [path và env](/docs/vi/stdlib/path-env) | `std::path::copy`, `std::env::environ` |
|Đo thời gian/chờ đợi| [time](/docs/vi/stdlib/time) | `std::time::duration`, `clock`, `sleep` |
|Tìm kiếm danh sách tên/địa chỉ bằng số| [net.resolve](/docs/vi/stdlib/resolution) | `std::net::resolve`, `resolve_table` |
|TCP Kết nối/Truyền tải| [net.tcp](/docs/vi/stdlib/tcp) | `std::net::tcp`, `address`, `error` |
|OS Số ngẫu nhiên| [random](/docs/vi/stdlib/random) | `std::random::fill` |
|Quá trình·OS Ranh giới| [chức năng hệ thống](/docs/vi/reference/system-io-network-process) | `std::process::core`, `spawn`, `std::sys` |
|Thực thi tác vụ không đồng bộ| [task](/docs/vi/stdlib/task) | `std::task` |
|Trợ lý toán học/chẩn đoán| [math và debug](/docs/vi/stdlib/math-debug) | `std::math::int`, `float`, `std::debug::core` |

## Ví dụ về lần sử dụng đầu tiên

Chương trình bên dưới sử dụng một chức năng từ std mà không cần tải xuống gói riêng. Lưu nó dưới dạng `main.wave` và chạy dưới dạng `wavec run main.wave`.

<!-- wave-example: library-import -->
```wave
import("std::string::len")::{
    len
};

fun main() {
    println("{}", len("Wave"));
}
```

Kết quả là `4`. Để hiểu từng bước của ví dụ tương tự, hãy đọc [Dây](/docs/vi/language/strings).

## Tương thích với trình biên dịch std

Xác nhận đường dẫn đã chọn bằng `wavec print std-path`. Khi sử dụng std từ một lần thanh toán khác, hãy chỉ định đường dẫn là `wavec --std-root /absolute/path/to/std check main.wave`. Nếu đường dẫn được chỉ định không hợp lệ hoặc không tương thích, lỗi sẽ được hiển thị.

## biên giới nền tảng

Phân biệt giữa các hàm tính toán như chuỗi/byte và các hàm OS như tệp/socket. Việc nhận ra mục tiêu không đảm bảo rằng tất cả các máy chủ API sẽ được cung cấp. Đọc các mục nhập nền tảng cho [Mục tiêu hỗ trợ](/docs/vi/whale/build-link-targets) và mỗi API cùng nhau. `std::sys` là giao diện cấp thấp hơn và các chương trình di động sẽ sử dụng mô-đun cấp cao hơn trước tiên.
