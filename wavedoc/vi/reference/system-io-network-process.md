---
translation_set_id: system-io
path: reference/system-io-network-process
locale: vi
group: stdlib
group_order: 1
order: 15
title: Chức năng và quy trình của hệ thống
summary: Mô tả ranh giới và thời gian tồn tại của quy trình của giao diện gốc API và OS.
---

## Tài liệu theo chức năng

Đọc [fs và io](/docs/vi/stdlib/files-io) để xử lý tệp, [TCP](/docs/vi/stdlib/tcp) để liên kết và [resolver](/docs/vi/stdlib/resolution) để tra cứu địa chỉ. Dưới đây là quy trình và quy tắc truy cập cấp thấp hơn OS.

## Quy trình cơ bản API

```text
std::process::core
proc_exit(code: i32) -> !
proc_getpid() -> i64
proc_getppid() -> i64
proc_execve(path: str, argv: ptr<ptr<i8>>, envp: ptr<ptr<i8>>) -> i64
proc_waitpid_raw(pid: i64, status: ptr<i32>, options: i32) -> i64
```

`proc_exit` không quay lại điểm gọi. Vui lòng thực hiện mọi thao tác dọn dẹp tập tin/bộ nhớ cần thiết trước khi tắt máy. `proc_execve` khác với chức năng tạo con thông thường vì nếu thành công, nó sẽ thay thế hình ảnh quy trình hiện có. raw argv/envp phải được chuẩn bị cho việc chấm dứt NUL của mỗi chuỗi, với con trỏ null biểu thị sự kết thúc.

Hàm spawn trong `std::process::spawn` xử lý kết quả tạo và hàm chờ xử lý trạng thái thoát của thành phần con. Tạo thành công và chấm dứt thành công chương trình là hai việc khác nhau. Khi bạn tạo một đường ống, cha mẹ và con phải đóng đầu không sử dụng để EOF được thông qua. Nếu bạn đợi trẻ thoát ra mà không đọc ống chụp, bộ đệm có thể đầy và đợi lẫn nhau.

## Tính di động và cách tiếp cận cấp thấp

fork/exec, bộ mô tả tệp và phần xử lý Windows không giống nhau về chức năng OS. Xác minh hỗ trợ cho mục tiêu đã chọn và coi unsupported như một đường dẫn lỗi thông thường. `std::sys` là giao diện dành riêng cho OS và không sử dụng lại cờ số cũng như bố cục của nó từ OS khác.

Khi liên kết trực tiếp với thư viện C bên ngoài, vui lòng đọc [FFI](/docs/vi/language/modules-imports-and-ffi). Không cần phải tự ý khai báo hàm libc để sử dụng hàm cha std API. Trước tiên hãy kiểm tra [Môi trường mục tiêu và liên kết](/docs/vi/whale/build-link-targets).
