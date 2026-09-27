---
translation_set_id: practice-tcp-client
path: practice/tcp-client
locale: vi
group: practice
group_order: 4
order: 4
title: Dự án: Máy khách TCP cục bộ
summary: Kết nối tra cứu địa chỉ số, lỗi kết nối, chuyển và đóng.
---

## Chuẩn bị và phạm vi

Chương trình này kết nối với `127.0.0.1:8080` trên bản gốc OS, gửi `ping` và LF rồi thoát. Bạn sẽ cần một máy chủ thử nghiệm trong một thiết bị đầu cuối riêng biệt. Nếu Python 3 xuất hiện, máy chủ tiếp theo sẽ nhận và hiển thị tối đa 5 byte từ một kết nối cục bộ.

```python
import socket
with socket.socket() as server:
    server.bind(("127.0.0.1", 8080))
    server.listen(1)
    connection, address = server.accept()
    with connection:
        message = b""
        while len(message) < 5:
            chunk = connection.recv(5 - len(message))
            if not chunk:
                break
            message += chunk
        print(repr(message))
```

Trước tiên hãy chạy máy chủ, sau đó lưu ứng dụng khách dưới dạng `main.wave` và chạy dưới dạng `wavec run main.wave`. Máy chủ xuất ra `b'ping\n'`. Nếu cổng đã được sử dụng, hãy thay đổi cổng trên cả máy chủ và máy khách cùng nhau.

<!-- wave-example: tcp-client -->
```wave
import("std::net::address")::{
    SocketAddr
};
import("std::net::resolve")::{
    NetResolveResult, net_resolve_tcp
};
import("std::net::error")::{
    NetResult, NetError, NET_ERROR_NONE
};
import("std::net::tcp")::{
    TcpStream, tcp_connect_addr_timeout, tcp_write_all, tcp_close
};

fun main() -> i32 {
    var addresses: array<SocketAddr, 4>;
    var resolved: NetResolveResult = net_resolve_tcp("127.0.0.1", "8080", &addresses[0], 4);
    if (!resolved.ok || resolved.written == 0) {
        println("resolve failed");
        return 1;
    }

    var connected: NetResult<TcpStream> = tcp_connect_addr_timeout(addresses[0], 1000);
    if (!connected.ok) {
        println("connect failed");
        return 2;
    }

    var sent: i64 = tcp_write_all(connected.value, "ping\n" as ptr<u8>, 5);
    var closed: NetError = tcp_close(connected.value);
    if (sent != 5 || closed.kind != NET_ERROR_NONE) {
        return 3;
    }

    println("sent=5");
    return 0;
}
```

Kết quả thành công từ khách hàng là `sent=5`. Nếu máy chủ không chạy, cách xử lý lỗi thông thường là in `connect failed` và thoát bằng 2.

## Hiểu mã

Phạm vi hợp lệ của mảng tra cứu lên tới written. Đây là một ví dụ nhỏ chỉ sử dụng địa chỉ đầu tiên. Trong các dịch vụ có nhiều ứng viên, bạn sẽ cần xác định chính sách kết nối cho từng ứng viên. Luồng được kết nối sẽ bị đóng bất kể truyền thành công hay thất bại. Thời gian chờ kết nối là 1000ms và đây không phải là ví dụ đảm bảo thời gian chờ cho toàn bộ quá trình truyền.

TCP không bảo toàn ranh giới truyền tải. Đây là lý do tại sao nó được viết để máy chủ có thể chạy recv nhiều lần. Cách tìm địa chỉ theo tên được đề cập trong [tra cứu địa chỉ](/docs/vi/stdlib/resolution).

## thực hành mở rộng

Yêu cầu máy chủ gửi phản hồi và yêu cầu khách hàng đọc nó. Sau khi xác định độ dài của phản hồi, bạn cần xử lý các lần đọc một phần, EOF, hết thời gian chờ và đóng lại với nhau.

## Vai trò của hai thiết bị đầu cuối

Thiết bị đầu cuối máy chủ chờ kết nối ở trạng thái listen. Khi khách hàng kết nối, accept trả về ổ cắm đã kết nối và vòng lặp recv thu thập dữ liệu. Thiết bị đầu cuối của khách hàng thực hiện tra cứu địa chỉ, kết nối, truyền và đóng theo thứ tự đó.

Mã Python là máy chủ đối tác của phòng thí nghiệm. Wave Nó được sử dụng để dễ dàng kiểm tra các byte được chương trình truyền vào và không được đặt trong cùng một tệp với mã máy khách. Máy chủ chấm dứt sau khi xử lý một kết nối, do đó, nó cũng khởi động lại trước khi chạy lại máy khách.

## Tại sao gửi 5 byte

Vì `ping` là 4 byte và LF là 1 byte nên độ dài truyền là 5. NUL ở cuối chuỗi không có trong tin nhắn. Máy chủ thu thập 5 byte và sau đó xuất chúng ra, vì vậy ngay cả khi các gói đến thành nhiều khối thì vẫn cho ra kết quả tương tự.

Nếu tcp_write_all trả về 5 thì hàm truyền cục bộ đã xử lý các byte được yêu cầu. Điều này không xác nhận rằng chương trình kia đã giải thích và lưu tin nhắn. Nếu giao thức yêu cầu xác nhận hoàn thành quá trình xử lý, máy chủ sẽ gửi phản hồi và máy khách sẽ đọc phản hồi.

## Con đường dẫn tới thất bại

|thay đổi|kết quả|Học gì|
| --- | --- | --- |
|Chạy mà không cần máy chủ|Kết nối không thành công, mã thoát 2|Ngay cả khi địa chỉ hợp lệ, máy chủ có thể không tồn tại|
|Đặt cổng máy chủ và máy khách khác nhau|Kết nối không thành công|Cả IP và cổng trong địa chỉ phải khớp nhau.|
|Thay đổi kích thước của máy chủ recv thành 1|5 byte giống nhau|TCP Đơn vị đọc khác với đơn vị tin nhắn|
|Chia khách hàng chuyển thành nhiều lần|nhận được theo thứ tự tương tự|Ranh giới tin nhắn được xác định bởi giao thức.|

Thời gian chờ kết nối 1000ms là cài đặt trong giai đoạn kết nối. Để hạn chế đọc và ghi, hãy sử dụng hàm timeout ở bước liên quan và nếu có giới hạn thời gian áp dụng cho toàn bộ chương trình, hãy tính thời gian còn lại và chuyển tiếp.
