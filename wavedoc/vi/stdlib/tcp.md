---
translation_set_id: stdlib-tcp
path: stdlib/tcp
locale: vi
group: stdlib
group_order: 1
order: 9
title: net.tcp: Kết nối và chuyển giao
summary: TCP Mô tả cấu trúc kết quả, trách nhiệm chuyển giao và ngắt kết nối một phần.
---

## Kiểm tra kết quả kết nối

```text
std::net::tcp
tcp_connect_addr(addr: SocketAddr) -> NetResult<TcpStream>
tcp_bind_loopback(port: u16) -> NetResult<TcpListener>
tcp_accept(listener: TcpListener) -> NetResult<TcpStream>
tcp_read(stream: TcpStream, buf: ptr<u8>, size: i64) -> i64
tcp_write_all(stream: TcpStream, buf: ptr<u8>, size: i64) -> i64
tcp_read_exact(stream: TcpStream, buf: ptr<u8>, size: i64) -> i64
tcp_close(stream: TcpStream) -> NetError
tcp_close_listener(listener: TcpListener) -> NetError
```

Hàm liên kết `NetResult<T>` có `ok`, `value` và `error`. Chỉ sử dụng value nếu thành công. `NetError` chứa phân loại lỗi chuẩn hóa và native_code, đồng thời có thể kiểm tra thành công bằng `error.kind == NET_ERROR_NONE`. Hằng số này được lấy từ `std::net::error`.

Luồng được kết nối và luồng nhận được bằng accept phải được đóng tương ứng. Đóng trình nghe không có nghĩa là đóng tất cả các luồng mà nó đã chấp nhận. Không đóng từng bản sao của cấu trúc.

## TCP không bảo toàn ranh giới tin nhắn

Đừng cho rằng mọi thứ bạn viết một lần sẽ quay lại với bạn sau khi bạn đọc nó. Giao thức phải được phân cách bằng tiền tố độ dài, dấu phân cách hoặc độ dài cố định. Độ dài cố định được xử lý bởi `tcp_read_exact` và các luồng có độ dài không xác định được xử lý bằng cách lặp lại đọc và xử lý EOF.

Số âm là lỗi và trong chiều dài đọc dương, 0 là điểm cuối của đầu kia. write_all Một số dữ liệu có thể đã được truyền trước khi lỗi. Nếu bạn truyền lại cùng một nội dung từ đầu, nó có thể bị trùng lặp.

## Chờ đợi và giới hạn thời gian

Chức năng blocking mặc định có thể đợi phía bên kia rất lâu. Đối với các chương trình yêu cầu giới hạn thời gian, hãy chọn dòng `tcp_connect_addr_timeout`, `tcp_read_timeout`, `tcp_write_timeout`. Kiểm tra yếu tố mili giây và hậu quả lỗi, đồng thời thiết kế lộ trình mà phía bên kia không phản hồi. Sử dụng không đồng bộ [task](/docs/vi/stdlib/task) cùng với mạng không đồng bộ tương ứng API.

[tra cứu địa chỉ](/docs/vi/stdlib/resolution) · [Khách hàng địa phương đã hoàn thành](/docs/vi/practice/tcp-client)

## Tiêu chí lựa chọn chức năng đọc

|hành động cần thiết|chức năng|giá trị cần kiểm tra|
| --- | --- | --- |
|Đọc một số dữ liệu đến| `tcp_read` |Lỗi âm, chấm dứt bằng 0, số byte dương|
|Đọc một độ dài nhất định| `tcp_read_exact` |Độ dài yêu cầu đã được đọc đầy đủ chưa?|
|Gửi nội dung nhất định đến cuối| `tcp_write_all` |Độ dài yêu cầu và giá trị trả về|
|Giới hạn thời gian chờ|timeout Dòng sản phẩm|Kết quả trả về và lỗi timeout|

Đối với giao thức có trường độ dài bốn byte, theo sau là phần thân, trước tiên hãy đọc trường độ dài đầy đủ, kiểm tra xem độ dài có vượt quá mức tối đa cho phép hay không, sau đó phân bổ không gian cho phần thân. Không phân bổ lượng lớn bộ nhớ từ độ dài không được xác thực do thiết bị ngang hàng cung cấp.

## Khi nào cần đóng kết nối

Ngay cả khi hàm đọc trả về 0, tài nguyên luồng cục bộ vẫn còn. Sau khi đọc xong hãy gọi tcp_close. Nếu bạn sử dụng cùng một đường dẫn dọn dẹp ngay cả sau khi xảy ra lỗi truyền tải, bạn sẽ không quên đóng các đường dẫn bình thường và đường dẫn bị lỗi.

Trình nghe là tài nguyên nhận các kết nối mới và luồng là tài nguyên liên lạc đã được kết nối. Khi tạo máy chủ, bạn quản lý cả hai loại. Mỗi lần accept thành công, một luồng mới sẽ được tạo nên chúng tôi sẽ đóng luồng đó sau khi xử lý và đóng trình nghe khi quá trình lặp lại chấp nhận của máy chủ kết thúc.

## Hãy tự mình thử

[Địa phương TCP Thực hành khách hàng](/docs/vi/practice/tcp-client) chứa mã hoàn chỉnh cho máy chủ và máy khách. Bạn có thể phân biệt giữa tra cứu địa chỉ và lỗi kết nối bằng cách so sánh các trường hợp máy chủ được chạy trước, khi không có máy chủ và khi số cổng khác nhau.
