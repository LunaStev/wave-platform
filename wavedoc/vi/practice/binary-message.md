---
translation_set_id: practice-binary-message
path: practice/binary-message
locale: vi
group: practice
group_order: 4
order: 3
title: Dự án: Tạo một thông báo nhị phân
summary: Sử dụng thứ tự byte rõ ràng và ULEB128 và từ chối đầu vào ngắn.
---

## định dạng tin nhắn

2 byte đầu tiên lưu trữ các số loại big-endian u16, sau đó là các giá trị ULEB128 u64. Nếu bạn ghi bộ nhớ cấu trúc vào một tệp, nó sẽ bị ảnh hưởng bởi phần đệm và thứ tự byte, vì vậy hãy mã hóa nó theo trường.

Lưu nó dưới dạng `main.wave` và chạy nó.

<!-- wave-example: binary-message -->
```wave
import("std::bytes::types")::{
    ByteReader, ByteWriter
};
import("std::bytes::cursor")::{
    bytes_reader, bytes_writer, bytes_writer_write_be_u16, bytes_reader_read_be_u16
};
import("std::bytes::leb128")::{
    bytes_writer_write_uleb128_u64, bytes_reader_read_uleb128_u64
};

fun main() -> i32 {
    var data: array<u8, 12>;
    var writer: ByteWriter = bytes_writer(&data[0], 12);
    if (bytes_writer_write_be_u16(&writer, 7) < 0) {
        return 1;
    }
    if (bytes_writer_write_uleb128_u64(&writer, 300) < 0) {
        return 2;
    }

    var reader: ByteReader = bytes_reader(&data[0], writer.position);
    var kind: u16 = 0;
    var value: u64 = 0;
    if (bytes_reader_read_be_u16(&reader, &kind) < 0) {
        return 3;
    }
    if (bytes_reader_read_uleb128_u64(&reader, &value) < 0) {
        return 4;
    }

    println("kind={} value={} bytes={}", kind, value, reader.position);
    var short: ByteReader = bytes_reader(&data[0], 1);
    kind = 99;
    if (bytes_reader_read_be_u16(&short, &kind) >= 0) {
        return 5;
    }

    println("preserved={} {}", short.position, kind);
    return 0;
}
```

Kết quả thực hiện:

```text
kind=7 value=300 bytes=4
preserved=0 99
```

## Những gì cần kiểm tra

Trong reader, chúng tôi chuyển độ dài ghi thực tế chứ không phải tổng dung lượng mảng là 12. Điều này nhằm tránh việc đọc các byte ở cuối chưa được khởi tạo làm đầu vào. Ngay cả khi quá trình đọc đầu vào ngắn không thành công, position=0 và kind=99 vẫn được duy trì.

## Bài tập và bình luận mở rộng

Nếu định dạng không cho phép thêm byte ở cuối, hãy kiểm tra `reader.position == reader.len` sau khi phân tích cú pháp hoàn tất. Khi thêm trường độ dài, hãy đảm bảo trường đó không lớn hơn các byte còn lại của đầu vào và phép tính độ dài+offset không vượt quá phạm vi.

[Xem bytes](/docs/vi/stdlib/bytes)

## Nhìn vào byte thực tế

Loại số 7 là big-endian u16 và do đó `00 07`. Giá trị 300 trở thành ULEB128 thành `AC 02`. Toàn bộ tin nhắn là bốn byte sau:

```text
00 07 AC 02
───── ─────
종류  값
```

ULEB128 sử dụng 7 bit thấp cho mỗi byte làm giá trị và nếu bit cao là 1 thì điều đó cho biết byte tiếp theo sẽ theo sau. 7 bit thấp của 300 là 44 và phần còn lại là 2. Byte đầu tiên là 44 cộng với 128 dấu liên tiếp hoặc 172 hoặc 0xAC. Không có dấu tiếp tục trong byte cuối cùng 0x02.

## Công suất và thời gian sử dụng

u16 sử dụng 2 byte và ULEB128 của u64 sử dụng tối đa 10 byte, do đó, mảng 12 byte có thể lưu trữ cả hai trường. Giá trị 300 chỉ sử dụng 2 byte, tạo ra thông báo thực tế là 4 byte. Khi gửi tới một tệp hoặc ổ cắm, bạn gửi byte writer.position thay vì toàn bộ mảng.

position trong reader là vị trí đọc hiện tại. Sau khi đọc loại, nó trở thành 2 và sau khi đọc giá trị, nó trở thành 4. Để chuyển sang thông báo tiếp theo khi xảy ra lỗi, phải biết ranh giới của thông báo lỗi. Chính sách khôi phục tin nhắn không được thiết lập tự động chỉ vì chức năng đọc giữ nguyên vị trí.

## Bài tập về giá trị biên

Thay đổi các giá trị thành 0, 127, 128, 16383, 16384 để xác định độ dài mã hóa. Khi thay đổi từ 127 thành 128, độ dài của ULEB128 tăng từ 1 lên 2 và khi thay đổi từ 16383 thành 16384, độ dài tăng từ 2 lên 3. Tính tổng độ dài, bao gồm cả 2 byte của trường loại.

Các tin nhắn có byte cuối cùng bị cắt ngắn cũng được kiểm tra. Nếu độ dài tin nhắn ban đầu là 4 thì chúng tôi chuyển độ dài 3 đến reader. Trường loại đã được đọc nhưng việc đọc giá trị phải không thành công vì thiếu byte cuối cùng của ULEB128. Tại thời điểm này, hãy kiểm tra xem vị trí 2 và giá trị đầu ra ngay trước khi đọc ULEB128 có được duy trì hay không.
