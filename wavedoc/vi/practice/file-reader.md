---
translation_set_id: practice-file-reader
path: practice/file-reader
locale: vi
group: practice
group_order: 4
order: 2
title: Đề tài: Đọc file và đếm byte
summary: Buffer, tệp I/O, xử lý lỗi liên kết và giải phóng bộ nhớ.
---

## sẵn sàng

Tạo input.txt trong thư mục làm việc chứa `Wave` theo sau là một dòng mới LF. Tệp sau đó chứa 5 byte. Với CRLF nó chứa 6 byte; BOM UTF-8 bổ sung thêm byte. Kiểm tra mã hóa tập tin và kết thúc dòng của trình soạn thảo.

Lưu chương trình dưới dạng `main.wave` và chạy nó dưới dạng `wavec run main.wave` từ cùng thư mục. Đường dẫn tương đối có liên quan đến thư mục làm việc đang chạy.

<!-- wave-example: file-reader -->
```wave
import("std::buffer::alloc")::{
    Buffer, buffer_init, buffer_free
};
import("std::fs::file")::{
    read_to_end
};

fun main() -> i32 {
    var data: Buffer;
    if (buffer_init(&data, 0) < 0) {
        return 1;
    }

    var count: i64 = read_to_end("input.txt", &data);
    if (count < 0) {
        println("read failed={}", count);
        buffer_free(&data);
        return 2;
    }

    var lines: i64 = 0;
    for (var i: i64 = 0; i < data.len; i += 1) {
        if (deref data.data[i] == 10) {
            lines += 1;
        }
    }

    println("bytes={} LF={}", count, lines);
    if (buffer_free(&data) < 0) {
        return 3;
    }

    return 0;
}
```

Kết quả thực hiện:

```text
bytes=5 LF=1
```

## Hành động và trách nhiệm

`read_to_end` mở và đóng tệp, nhưng việc giải phóng Buffer là trách nhiệm của người gọi. Tắt cả đường dẫn thành công và đường dẫn đọc-thất bại. Vì đây là dữ liệu có số byte nên chúng tôi không cho rằng đó là một chuỗi kết thúc bằng NUL.

LF Số đường và số đường mà người ta nghĩ không phải lúc nào cũng giống nhau. Nếu dòng cuối cùng không chứa LF thì nó không được tính vào số lượng LF cho chương trình này. Đồng thời kiểm tra lỗi đọc bằng cách đổi tên tệp. Số lỗi cụ thể có thể khác nhau tùy thuộc vào môi trường của bạn.

## Bài tập và bình luận mở rộng

Nếu bạn muốn xử lý các tệp rất lớn, hãy thay thế nó bằng một mảng có kích thước cố định và lặp lại `io_read`. Chỉ xử lý phạm vi trả về dương và kết thúc ở 0. Bạn có thể tích lũy số byte và số lượng LF mà không cần phải giữ toàn bộ tệp trong bộ nhớ. Nếu bạn tự mở nó, nó cũng sẽ đóng bộ mô tả trên bất kỳ đường dẫn thoát nào.

[Xem fs và io](/docs/vi/stdlib/files-io) · [Xem Buffer](/docs/vi/stdlib/buffer)

## Thực hiện theo quy trình xử lý

1. Khởi tạo thùng Buffer. Chưa có nội dung tập tin.
2. read_to_end đọc tệp và tăng dung lượng cần thiết.
3. Nếu đọc thành công, các byte trong phạm vi data.len sẽ được kiểm tra.
4. Bất cứ khi nào chúng tôi gặp giá trị byte 10 của LF, chúng tôi sẽ tăng lines.
5. In kết quả và phát hành Buffer.

data.cap là không gian lưu trữ dành riêng và data.len là độ dài dữ liệu hợp lệ. Nếu bạn thay đổi điều kiện lặp lại thành cap, các byte không có trong tệp sẽ được đọc, vì vậy hãy sử dụng len. count là số byte được thêm vào trong lệnh gọi tới read_to_end này. Ví dụ này bắt đầu bằng Buffer trống, vì vậy count và data.len bằng nhau.

## Thay đổi đầu vào để kiểm tra

|input.txt Nội dung| bytes | LF |lý do|
| --- | --- | --- | --- |
|tập tin trống| 0 | 0 |Không có byte để đọc|
| `Wave` | 4 | 0 |Không ngắt dòng cuối cùng|
| `Wave` + LF | 5 | 1 |Dữ liệu cập nhật đến cuối cùng LF|
| `A` + LF + `B` + LF | 4 | 2 |Đếm hai LF|
| `Wave` + CRLF | 6 | 1 |CR cũng là 1 byte nhưng chỉ tính LF.|

Để đếm các tệp không có LF ở dòng cuối cùng dưới dạng một dòng, hãy thêm 1 vào số dòng nếu tệp không trống và byte cuối cùng không phải là 10. Trước tiên, bạn phải kiểm tra xem data.len có phải là 0 hay không trước khi bạn có thể truy cập phần tử cuối cùng.

## Mở rộng sang các tệp lớn

Phương pháp lưu trữ toàn bộ nội dung hiện nay thuận tiện cho việc đọc lại hoặc truy xuất dữ liệu sau này. Nếu bạn chỉ cần số byte và số lượng LF, việc sử dụng lại bộ đệm có kích thước cố định sẽ hợp lý.

Trong vòng lặp io_read trong ví dụ [Đọc tệp vào bộ đệm có kích thước cố định](/docs/vi/stdlib/files-io), chỉ cần đếm LF theo số byte được trả về. Nó có thể được xử lý với bộ nhớ bằng kích thước bộ đệm chứ không phải toàn bộ chiều dài của tệp. Nếu quá trình đọc trả về 0 thì kết thúc; nếu nó âm, đó là một lỗi.
