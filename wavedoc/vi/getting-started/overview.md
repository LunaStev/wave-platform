---
translation_set_id: overview
path: getting-started/overview
locale: vi
group: getting-started
group_order: 1
order: 1
title: Tài liệu và hướng dẫn học Wave
summary: Tìm hiểu Wave từng bước, từ cài đặt đến chương trình thực tế, tra cứu các quy tắc ngôn ngữ và API thư viện chuẩn.
---

## Tìm hiểu Wave với hướng dẫn này

Học cách viết mã nguồn Wave, biên dịch và chạy mã đó cũng như kiểm tra kết quả. Nếu bạn chưa quen với lập trình, hãy làm theo trình tự dưới đây. Nếu bạn biết một ngôn ngữ khác, hãy chạy các ví dụ của từng chương và so sánh các quy tắc cũng như trường hợp ranh giới của nó với những gì bạn đã biết.

## Lộ trình học tập

|Bước|chương|Bạn sẽ học được gì|
| --- | --- | --- |
|thiết lập| [Cài đặt](/docs/vi/getting-started/install) |Chuẩn bị trình biên dịch và thư viện chuẩn và xác minh rằng chúng chạy|
| 1 | [Chương trình đầu tiên của bạn](/docs/vi/language/program-structure) |Tạo, kiểm tra và chạy tệp nguồn cũng như hiểu mã thoát|
| 2 | [Biến và loại](/docs/vi/language/declarations-and-types) |Lưu trữ giá trị và chọn loại có phạm vi được yêu cầu|
| 3 | [Toán tử và chuyển đổi](/docs/vi/language/expressions-and-operators) |Giải thích thứ tự đánh giá và kết quả chuyển đổi loại|
| 4 | [Điều kiện và vòng lặp](/docs/vi/language/control-flow) |Phân nhánh theo điều kiện và xử lý dữ liệu với các vòng lặp|
| 5 | [Hàm](/docs/vi/language/functions-and-generics) |Trích xuất các thao tác lặp lại thành các hàm|
| 6 | [Mảng](/docs/vi/language/arrays) |Truy cập các phần tử theo chỉ mục và lặp qua một mảng|
| 7 | [Dây](/docs/vi/language/strings) |Phân biệt các ký tự với byte và hiểu các ký tự thoát và độ dài chuỗi|
| 8 | [Cấu trúc và biến thể](/docs/vi/language/structures-enums-and-aliases) |Nhóm dữ liệu liên quan và thể hiện sự thành công và thất bại|
| 9 | [Con trỏ và vòng đời](/docs/vi/language/explicit-memory-type-model) |Sửa đổi giá trị ban đầu thông qua địa chỉ của nó và quản lý vòng đời của nó|
| 10 | [Bộ nhớ động](/docs/vi/language/allocation) |Xử lý lỗi phân bổ và bộ nhớ trống|
| 11 | [Mô-đun và thuốc generic](/docs/vi/language/modules-imports-and-ffi) |Chia mã trên các tệp và sử dụng lại các chức năng với các loại khác nhau|
| 12 | [Xử lý lỗi](/docs/vi/language/errors) |Kiểm tra kết quả và dọn sạch tài nguyên khi thất bại|
| 13 | [Giới thiệu về mã không đồng bộ](/docs/vi/language/async-and-never) |Tạo ra một Tương lai và chờ đợi nó kết thúc|

## Áp dụng kiến thức của bạn vào thực tế

Sau các chương cốt lõi, hãy xây dựng [máy tính đầu vào](/docs/vi/practice/input-calculator), [trình đọc tệp](/docs/vi/practice/file-reader), [thông báo nhị phân](/docs/vi/practice/binary-message) và [máy khách TCP](/docs/vi/practice/tcp-client). Kiểm tra cả trường hợp đầu vào thành công và thất bại trong từng dự án.

## Ba tab tài liệu

- **Wave**: Một khóa học ngôn ngữ có hướng dẫn và các dự án thực tế theo thứ tự.
- **[Thư viện tiêu chuẩn](/docs/vi/stdlib)**: API của mỗi mô-đun, giá trị trả về, lỗi, quy tắc quyền sở hữu và yêu cầu nền tảng.
- **[Whale](/docs/vi/whale)**: Xây dựng và liên kết, quản lý gói, sử dụng lệnh và chuỗi công cụ cấp thấp.

Các ví dụ phân biệt các chương trình hoàn chỉnh với các đoạn mã bên trong một hàm. Chạy các lệnh `wavec` trong một thiết bị đầu cuối và lưu các khối mã `wave` trong các tệp `.wave`. Đầu vào và đầu ra được hiển thị riêng biệt; ví dụ đọc đầu vào tiêu chuẩn chỉ định những gì cần nhập.

## Khi bạn gặp khó khăn

Sử dụng [Khắc phục sự cố](/docs/vi/reference/diagnostics) để phân biệt các vấn đề về cài đặt, kiểm tra nguồn, liên kết và thực thi. Tra cứu các quy tắc ngôn ngữ trong [tham chiếu nhanh cú pháp](/docs/vi/reference/syntax-quick-reference), các lệnh trong [tham chiếu trình biên dịch](/docs/vi/getting-started/compiler) và API trong [hướng dẫn thư viện tiêu chuẩn](/docs/vi/reference/standard-library).
