---
translation_set_id: ecosystem
path: whale/ecosystem
locale: vi
group: whale
group_order: 1
order: 2
title: Thành phần chuỗi công cụ
summary: Mô tả vai trò của chuỗi công cụ cấp thấp riêng biệt, Whale và các ranh giới thành phần của hệ sinh thái Wave.
---

## WhaleIran

Whale là một chuỗi công cụ cấp thấp xử lý các biểu diễn lắp ráp và trung gian. Các thành phần xử lý các tập hợp, đối tượng, liên kết và biểu diễn trung gian được thiết kế để có thể tái sử dụng trên Wave và các công cụ tạo mã gốc khác.

Whale không phải là tên của toàn bộ môi trường phát triển Wave. Trách nhiệm của từng dự án được phân chia như sau:

|dự án|trách nhiệm|
| --- | --- |
| `wavec` |Wave Kiểm tra nguồn và tạo tệp thực thi.|
| Vex |Quản lý gói Wave, manifest, biểu đồ phụ thuộc, lockfile và các bản dựng gói.|
| Whale |Cung cấp các thành phần assembler, object, linker và IR độc lập.|
| Wave `std` |Thời gian chạy và hệ thống API được cung cấp dưới dạng mô-đun nguồn Wave.|

## thành phần

Whale workspace bao gồm 4 khu vực thư viện chính:

- `assembler`: Mã hóa, AMD64 Phân tích cú pháp/Mã hóa, section, symbol và relocation
- `object`: mô hình tệp đối tượng và ELF64 writer
- `linker`: Lớp liên kết
- `ir`: Whale IR loại, builder, đầu ra, xác minh và tùy chọn frontend socket

Tệp thực thi `whale` cung cấp cho vùng này các lệnh `asm`, `object`, `link` và `ir`.

## biên giới công cụ

Chương trình Wave được xây dựng dưới dạng `wavec`. Khi xử lý trực tiếp với các tập hợp, tệp đối tượng và IR, hãy sử dụng lệnh `whale`.

Việc cài đặt Whale không làm thay đổi phương thức xây dựng của `wavec`. Vex sử dụng `wavec` để xây dựng gói Wave và Whale chạy gói này trực tiếp trong một tác vụ xử lý các tạo phẩm cấp thấp.

## Xác minh có thể giao được

Khi liên kết cấu phần phần mềm Whale với quy trình xây dựng của bạn, hãy đảm bảo rằng object format và mục tiêu architecture khớp với nhau. symbol và relocation có thể được kiểm tra bằng các công cụ độc lập như `readelf` và `objdump`. Các bản dựng sử dụng IR socket phải sử dụng socket schema từ cùng một nhà sản xuất với Whale.
