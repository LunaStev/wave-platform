---
translation_set_id: stdlib-resolution
path: stdlib/resolution
locale: vi
group: stdlib
group_order: 1
order: 8
title: net.resolve: Phân giải địa chỉ
summary: Mô tả việc cắt bớt địa chỉ số, danh sách tên do người gọi sở hữu và kết quả.
---

## Chọn phương pháp tra cứu

Giá trị mặc định cho Linux, resolver, lấy địa chỉ số IPv4/IPv6 và một cổng số. Tên máy chủ hoặc tên dịch vụ không được truyền ngầm tới hệ thống DNS. Ví dụ: `127.0.0.1` và `8080` là đầu vào số và `example.com` và `http` là tên.

Để sử dụng danh sách tên được cung cấp rõ ràng, hãy sử dụng `std::net::resolve_table`. Bạn có thể kết nối trực tiếp tới địa chỉ số hoặc tìm kiếm bằng cách đăng ký địa chỉ vào danh sách tên.

## Tra cứu cơ bản

```text
std::net::resolve
net_resolve_tcp(host: str, service: str, output: ptr<SocketAddr>, capacity: i64) -> NetResolveResult
net_resolve_udp(host: str, service: str, output: ptr<SocketAddr>, capacity: i64) -> NetResolveResult
net_resolve_host(host: str, service: str, output: ptr<SocketAddr>, capacity: i64) -> NetResolveResult
```

Mảng đầu ra được cung cấp bởi người gọi. Đơn vị của capacity không phải là byte mà là `SocketAddr` số phần tử. Bạn chỉ có thể tìm kiếm số đếm với capacity=0 và null output. Cửa hàng tra cứu nội bộ được dọn sạch trước khi trả về và không sở hữu mảng đầu ra.

|trường kết quả|ý nghĩa|
| --- | --- |
| `ok` |Truy vấn có thành công hay không|
| `count` |Tổng số kết quả địa chỉ được hỗ trợ|
| `written` |Số phần tử thực sự được ghi vào mảng đầu ra|
| `truncated` |Công suất đầu ra có nhỏ hơn tổng kết quả không?|
| `error.kind` |Các phân loại được chuẩn hóa như INVALID, NOT_FOUND, UNSUPPORTED|
| `error.native_code` |Mã gốc để xác định nguyên nhân|

Ngay cả khi thành công, written có thể bằng 0. Hãy kiểm tra `ok && written > 0` trước khi sử dụng output[0].

## Danh sách tên bạn tự cung cấp

```text
net_resolve_from_table(host: str, service: str, protocol: i32,
    entries: ptr<NetResolveEntry>, entry_count: i64,
    output: ptr<SocketAddr>, capacity: i64) -> NetResolveResult
```

`NetResolveEntry` trong `std::net::resolve_table` có host, service, protocol, address. Tên máy chủ ASCII không phân biệt chữ hoa chữ thường và tên dịch vụ được so sánh chính xác. TCP/UDP/Đối với tất cả các truy vấn giao thức, các hằng số của mô-đun được sử dụng. Các mục và chuỗi thuộc sở hữu của người gọi và không được giữ lại sau cuộc gọi.

Các kết quả trùng khớp được sao chép theo thứ tự đầu vào và các bản sao được duy trì. Không gian đầu ra và không gian lưu trữ đầu vào không được chồng chéo. Nếu tên không có trong danh sách thì đó là NOT_FOUND và sẽ không thử lại với tên DNS bên ngoài.

[TCP Cách sử dụng](/docs/vi/stdlib/tcp) · [TCP Thực hành khách hàng](/docs/vi/practice/tcp-client)

## Tra cứu địa chỉ dạng số

Chuẩn bị một mảng để lưu kết quả truy vấn và kiểm tra xem địa chỉ có được ghi lại hay không. Ví dụ bên dưới không yêu cầu máy chủ vì việc tra cứu địa chỉ không tự tạo ra kết nối.

<!-- wave-example: book-resolve-numeric -->
```wave
import("std::net::address")::{
    SocketAddr
};
import("std::net::resolve")::{
    NetResolveResult,
    net_resolve_tcp
};

fun main() -> i32 {
    var addresses: array<SocketAddr, 4>;
    var result: NetResolveResult = net_resolve_tcp(
        "127.0.0.1",
        "8080",
        &addresses[0],
        4
    );

    if (!result.ok || result.written == 0) {
        println("address unavailable");
        return 1;
    }

    println("address ready");
    return 0;
}
```

Kết quả thực hiện:

```text
address ready
```

Để thực sự kết nối, hãy chuyển addresses[0] cho chuỗi hàm tcp_connect_addr. Ngay cả khi bạn nhận được địa chỉ, ở giai đoạn kết nối, bạn có thể biết liệu máy chủ có đang chạy trên cổng đó hay không.

## Tìm kiếm theo danh sách tên

Danh sách địa chỉ tùy chỉnh rất hữu ích cho việc liên kết tên dựa trên tệp cấu hình hoặc danh sách dịch vụ của chương trình. Ví dụ sau ánh xạ tên api tới cổng 8080 cục bộ.

<!-- wave-example: book-resolve-table -->
```wave
import("std::net::address")::{
    SocketAddr,
    socket_addr_from_v4,
    socket_addr_v4_loopback
};
import("std::net::resolve")::{
    NetResolveResult
};
import("std::net::resolve_table")::{
    NetResolveEntry,
    RESOLVE_TCP,
    net_resolve_from_table
};

fun main() -> i32 {
    var entries: array<NetResolveEntry, 1>;
    entries[0] = NetResolveEntry {
        host: "api",
        service: "http",
        protocol: RESOLVE_TCP,
        address: socket_addr_from_v4(socket_addr_v4_loopback(8080))
    };

    var addresses: array<SocketAddr, 2>;
    var result: NetResolveResult = net_resolve_from_table(
        "API",
        "http",
        RESOLVE_TCP,
        &entries[0],
        1,
        &addresses[0],
        2
    );

    if (!result.ok || result.written != 1) {
        return 1;
    }

    println("matched={}", result.written);
    return 0;
}
```

Kết quả thực hiện:

```text
matched=1
```

API và api khớp với nhau khi so sánh tên máy chủ. http và HTTP không khớp nhau trong ví dụ này vì chúng khác nhau khi so sánh tên dịch vụ. Nếu bạn đăng ký nhiều địa chỉ có cùng tên thì kết quả sẽ xuất hiện theo thứ tự đã nhập. Nếu bạn chỉ nhận được một phần của mảng nhỏ, chỉ đọc tối đa written và nếu bạn cần toàn bộ danh sách, hãy chuẩn bị khoảng trống cho count và tìm kiếm lại.
