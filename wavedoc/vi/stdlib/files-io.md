---
translation_set_id: stdlib-files-io
path: stdlib/files-io
locale: vi
group: stdlib
group_order: 1
order: 7
title: fs và io: Truyền tệp và byte
summary: Mô tả thời gian tồn tại của tệp, trạng thái đọc toàn bộ, chuyển một phần và sau lỗi.
---

## Chọn tệp API

Chức năng tiện lợi của `std::fs::file` nhận đường dẫn và thực hiện việc mở và đóng cần thiết. Các hàm trả về một bộ mô tả phải được người gọi đóng lại.

|khai báo|Kết quả thành công và biện pháp phòng ngừa|
| --- | --- |
| `open_read(path: str) -> i64` |Mở mô tả. Số âm là lỗi|
| `create(path: str) -> i64` |Tạo hoặc xóa nội dung tập tin hiện có. Trả về bộ mô tả sở hữu|
| `open_append(path: str) -> i64` |Mở hoặc tạo để bổ sung|
| `size(path: str) -> i64` |Số byte. Số âm là lỗi|
| `read_into(path: str, dst: ptr<u8>, dst_cap: i64) -> i64` |Số byte trong toàn bộ tập tin. Thiếu năng lực là một lỗi|
| `read_to_end(path: str, dst_buffer: ptr<Buffer>) -> i64` |Thêm tệp sau Buffer hiện có và trả lại số tiền bổ sung|
| `write(path: str, src: ptr<u8>, len: i64) -> i64` |Số byte được ghi thay thế nội dung hiện có|
| `append(path: str, src: ptr<u8>, len: i64) -> i64` |Số byte được thêm vào cuối|
| `remove(path: str) -> i64` |tình trạng loại bỏ. thất bại là tiêu cực|

Riêng false của `exists(path)` không thể phân biệt được giữa thiếu file và lỗi cấp phép. Hãy chắc chắn kiểm tra kết quả mở thực tế vì trạng thái có thể thay đổi giữa việc kiểm tra sự tồn tại và mở.

## Mức độ thấp I/O

```text
std::io::fd
io_read(fd: i64, buf: ptr<u8>, len: i64) -> i64
io_write(fd: i64, buf: ptr<u8>, len: i64) -> i64
io_read_exact(fd: i64, buf: ptr<u8>, len: i64) -> i64
io_write_all(fd: i64, buf: ptr<u8>, len: i64) -> i64
io_close(fd: i64) -> i64
```

Kết quả dương tính của `io_read` là số byte được đọc và 0 trong yêu cầu có độ dài dương là EOF. `io_write` có thể được viết ít hơn yêu cầu. Nếu cần chuyển toàn bộ, hãy sử dụng chức năng exact/all. Tuy nhiên, chúng tôi không cho rằng lỗi sẽ hoàn nguyên trạng thái bên ngoài vì một số lần chuyển có thể đã xảy ra trước khi xảy ra lỗi.

`io_read_exact` là lỗi nếu gặp EOF trước độ dài yêu cầu. `read_into` trả về `IO_ERR_NO_SPACE` nếu bộ đệm đầy và một số byte có thể đã được ghi. Chức năng đọc không tự động thêm NUL vào cuối chuỗi.

## Buffer và xử lý lỗi

Nếu không thành công, `read_to_end` khôi phục len gốc nhưng dung lượng và địa chỉ dữ liệu của nó có thể đã thay đổi. Người gọi phải giải phóng Buffer sau khi thành công hoặc thất bại. API ghi tệp không đảm bảo thay thế tệp nguyên tử.

Từ [Luyện đọc file](/docs/vi/practice/file-reader), bạn có thể chạy chương trình từ import để phát hành. Hãy xem xét sự khác biệt về đường dẫn/quyền trong Linux/macOS/Windows/FreeBSD và các hạn chế về thư mục có thể truy cập trong WASI. Nó không diễn giải trực tiếp giá trị mô tả dưới dạng mã điều khiển thô của một OS khác.

## Đọc các tệp lớn vào bộ đệm nhỏ

Các thao tác không yêu cầu đặt toàn bộ tệp vào bộ nhớ có thể được xử lý bằng bộ đệm cố định và đọc lặp lại. Chương trình sau in nội dung của input.txt và đếm tổng số byte đã đọc. Lưu một `Wave` và LF trong tệp đầu vào.

<!-- wave-example: book-file-stream -->
```wave
import("std::fs::file")::{open_read};
import("std::io::fd")::{io_read, io_write_all, io_close};
import("std::io::consts")::{IO_STDOUT_FD};

fun main() -> i32 {
    var descriptor: i64 = open_read("input.txt");

    if (descriptor < 0) {
        return 1;
    }

    var buffer: array<u8, 4>;
    var total: i64 = 0;

    while (true) {
        var count: i64 = io_read(descriptor, &buffer[0], 4);

        if (count < 0) {
            io_close(descriptor);
            return 2;
        }

        if (count == 0) {
            break;
        }

        if (io_write_all(IO_STDOUT_FD, &buffer[0], count) < 0) {
            io_close(descriptor);
            return 3;
        }

        total += count;
    }

    if (io_close(descriptor) < 0) {
        return 4;
    }

    println("bytes={}", total);
    return 0;
}
```

Kết quả thực hiện:

```text
Wave
bytes=5
```

Dung lượng bộ đệm là 4, nhưng lần đọc cuối cùng có thể là 1 byte. Chúng tôi luôn chuyển count thực tế vào đầu ra. Nếu bạn ghi toàn bộ mảng, thậm chí cả các byte cũ, chưa đọc cũng có thể được xuất ra.

Chương trình đã mở descriptor trong số input.txt, vì vậy hãy đóng nó lại. Đầu ra tiêu chuẩn không phải là tài nguyên mới có được trong hàm này, vì vậy nó không bị đóng tùy ý ở cuối ví dụ.

## Đọc đầy đủ và không đủ năng lực

read_into nhận được một không gian lưu trữ cố định chứa toàn bộ tệp. Nếu không đủ dung lượng, nó sẽ âm thầm cắt bớt và trả về NO_SPACE nhưng không thành công.

<!-- wave-example: book-file-capacity -->
```wave
import("std::fs::file")::{read_into};
import("std::io::consts")::{IO_ERR_NO_SPACE};

fun main() {
    var data: array<u8, 2>;
    var status: i64 = read_into("input.txt", &data[0], 2);

    if (status == IO_ERR_NO_SPACE) {
        println("destination too small");
    }
}
```

Kết quả thực hiện:

```text
destination too small
```

Nó không cho rằng việc đọc thất bại hoàn toàn không thay đổi byte đích. Không sử dụng nó làm nội dung tệp hoàn chỉnh, hãy chuẩn bị một kho lưu trữ lớn hơn hoặc chọn phương thức streaming. Ngay cả khi bạn truy vấn kích thước trước, kết quả đọc thực tế vẫn là phán quyết cuối cùng vì tệp có thể thay đổi giữa truy vấn và đọc.

## Sự khác biệt giữa ghi và nối thêm tập tin

write thay thế nội dung hiện có và append được thêm vào cuối. Số byte cần lưu trữ có thể được lấy trực tiếp từ độ dài chuỗi và được truyền. NUL ở cuối chuỗi thường không có trong nội dung tệp văn bản.

<!-- wave-example: book-file-write-append -->
```wave
import("std::fs::file")::{write, append, size, remove};
import("std::string::len")::{len};

fun main() -> i32 {
    var first: str = "Wave";
    var second: str = " study";

    if (write("output.txt", first as ptr<u8>, len(first) as i64) < 0) {
        return 1;
    }

    if (append("output.txt", second as ptr<u8>, len(second) as i64) < 0) {
        return 2;
    }

    println("bytes={}", size("output.txt"));

    if (remove("output.txt") < 0) {
        return 3;
    }

    return 0;
}
```

Kết quả thực hiện:

```text
bytes=10
```

Ví dụ này tạo, thay thế và cuối cùng xóa output.txt trong thư mục làm việc. Chạy từ thư mục thực hành không có tập tin hiện có. Trình chỉnh sửa hoặc chương trình lưu trữ thực tế có thể yêu cầu các chính sách lưu trữ riêng biệt, chẳng hạn như các tệp tạm thời và thay thế.

## API Bảng tuyển chọn

|tình huống|chọn|
| --- | --- |
|Đọc toàn bộ tệp nhỏ vào bộ đệm cố định| read_into |
|Giữ toàn bộ nội dung mà không cần biết kích thước|read_to_end và Buffer|
|Xử lý nội dung theo thứ tự thay vì lưu trữ toàn bộ nội dung|open_read + io_read lặp lại|
|Đọc bản ghi có độ dài cố định| io_read_exact |
|Truyền toàn bộ chuỗi byte| io_write_all |
|Xử lý các tập tin đã mở|Hàm fd thay vì hàm đường dẫn|

Sau khi chọn một chức năng, hãy kiểm tra xem bộ đệm, vị trí tệp và dữ liệu bên ngoài thay đổi như thế nào trong trường hợp xảy ra lỗi.
