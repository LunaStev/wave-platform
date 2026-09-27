---
translation_set_id: stdlib-bytes
path: stdlib/bytes
locale: vi
group: stdlib
group_order: 1
order: 6
title: bytes: Phạm vi, con trỏ và ULEB128
summary: Mô tả các byte đọc/ghi có độ dài view và duy trì trạng thái trong trường hợp bị lỗi.
---

## Nó khác với một chuỗi như thế nào?

Dữ liệu byte có thể chứa số 0, vì vậy hãy chuyển một con trỏ cùng với độ dài. `Bytes` và `BytesMut` là các chế độ xem không sở hữu và chỉ hợp lệ khi bộ nhớ cơ bản vẫn hợp lệ. `BytesMut` yêu cầu bộ nhớ có thể ghi.

## Tạo cursor

```text
std::bytes::cursor
bytes_reader(data: ptr<u8>, len: i64) -> ByteReader
bytes_writer(data: ptr<u8>, len: i64) -> ByteWriter
bytes_reader_read_be_u16(reader: ptr<ByteReader>, out_value: ptr<u16>) -> i32
bytes_writer_write_be_u16(writer: ptr<ByteWriter>, value: u16) -> i32
```

Nhập `ByteReader` và `ByteWriter` từ `std::bytes::types`. Trường vị trí của họ xác định vị trí của hoạt động tiếp theo. Độ dài của chúng là tổng số byte có thể truy cập được. Việc xây dựng con trỏ không sao chép hoặc phân bổ bộ nhớ cơ bản.

`be` là big-endian, `le` là little-endian. Nếu loại tệp là big-endian, hãy sử dụng hàm `be` bất kể thứ tự byte của máy chủ CPU. Có các hàm 16, 32 và 64 bit signed/unsigned đọc/ghi và một byte.

## Lỗi và bảo quản trạng thái

`BYTES_OK` từ `std::bytes::errors` là 0. INVALID biểu thị phạm vi không hợp lệ, EOF đầu vào không đủ, NO_SPACE không đủ công suất đầu ra và OVERFLOW một giá trị nằm ngoài phạm vi có thể biểu thị. Đã kiểm tra vị trí tiến của các thao tác con trỏ chỉ sau khi toàn bộ thao tác thành công. Việc đọc không thành công cũng giữ nguyên giá trị đầu ra.

## ULEB128

```text
std::bytes::leb128
bytes_reader_read_uleb128_u64(reader: ptr<ByteReader>, output_value: ptr<u64>) -> i32
bytes_writer_write_uleb128_u64(writer: ptr<ByteWriter>, value: u64) -> i32
```

ULEB128 lưu trữ số nguyên 64 bit không dấu trong một số byte thay đổi, sử dụng tối đa 10 byte. Nếu không đủ dung lượng, bộ ghi sẽ giữ nguyên cả vị trí của nó và byte đích. Người đọc phân biệt đầu vào không đầy đủ với giá trị vượt quá u64. Chấp nhận mã hóa không tối thiểu, chấm dứt.

Để tạo một tin nhắn thực tế và thấy lỗi đầu vào ngắn, hãy tiếp tục tới [Thực hành tin nhắn nhị phân](/docs/vi/practice/binary-message). Đừng cố xuất chuỗi byte chứa số 0 dưới dạng `str`.

## Đọc các byte giống nhau theo các thứ tự khác nhau

Thứ tự byte là quy tắc lưu trữ số. Nếu bạn đọc hai byte 1 và 2 là big-endian thì nó là 1×256+2 và nếu bạn đọc chúng là little-endian thì nó là 2×256+1. Chọn dựa trên quy tắc mạng hoặc loại tệp.

<!-- wave-example: book-bytes-endian -->
```wave
import("std::bytes::types")::{Bytes, bytes_view};
import("std::bytes::read")::{bytes_read_be_u16, bytes_read_le_u16};

fun main() -> i32 {
    var data: array<u8, 2> = [1, 2];
    var view: Bytes = bytes_view(&data[0], 2);
    var big: u16 = 0;
    var little: u16 = 0;

    if (bytes_read_be_u16(view, 0, &big) < 0) {
        return 1;
    }

    if (bytes_read_le_u16(view, 0, &little) < 0) {
        return 2;
    }

    println("be={} le={}", big, little);
    return 0;
}
```

Kết quả thực hiện:

```text
be=258 le=513
```

view mượn một mảng. Vì không có sự phân bổ hoặc sao chép riêng biệt nên mảng chỉ có thể được sử dụng khi nó hợp lệ. offset tính bằng byte và việc đọc 16 bit cần 2 byte từ vị trí đó.

## Chọn offset API và cursor API

Hàm read/write, lấy đối số offset, thuận tiện cho các định dạng đọc trực tiếp một vị trí trường được chỉ định. Đối với các luồng trong đó vị trí tiếp theo phụ thuộc vào độ dài của trường trước đó, cursor với position là thuận tiện.

Khi trộn cả hai phải phân rõ đâu là chuẩn: cursor.position hay tách riêng offset. Tránh sai lầm khi thêm cùng một vị trí hai lần hoặc di chuyển đến vị trí tiếp theo mà không đọc thành công.

## Kiểm tra trạng thái từ đầu vào ngắn

<!-- wave-example: book-bytes-short-read -->
```wave
import("std::bytes::types")::{ByteReader};
import("std::bytes::cursor")::{bytes_reader, bytes_reader_read_be_u16};
import("std::bytes::errors")::{BYTES_ERROR_EOF};

fun main() {
    var data: array<u8, 1> = [1];
    var reader: ByteReader = bytes_reader(&data[0], 1);
    var value: u16 = 99;
    var status: i32 = bytes_reader_read_be_u16(&reader, &value);

    if (status == BYTES_ERROR_EOF) {
        println("position={} value={}", reader.position, value);
    }
}
```

Kết quả thực hiện:

```text
position=0 value=99
```

Không cần hai byte nên giá trị đầu ra và vị trí được giữ nguyên. Thuộc tính này hữu ích trong các thiết kế thử đọc lại cùng một trường sau khi có thêm đầu vào. Tuy nhiên, nếu không gian lưu trữ được trỏ tới bởi view đã được phân bổ lại thì địa chỉ cũng phải được cập nhật.

## Ranh giới của ULEB128

0~127 sử dụng một byte, 128 trở đi sử dụng nhiều byte hơn. Bit thứ tự cao của mỗi byte cho biết liệu dữ liệu có tuân theo hay không. Giá trị nằm ngoài u64 hoặc giá trị đầu vào liên tục quá dài là OVERFLOW, khác với EOF, đơn giản là có ít đầu vào hơn.

Trực tiếp kiểm tra xem Dung lượng không đủ writer có duy trì trạng thái hay không.

<!-- wave-example: book-leb128-space -->
```wave
import("std::bytes::types")::{ByteWriter};
import("std::bytes::cursor")::{bytes_writer};
import("std::bytes::leb128")::{bytes_writer_write_uleb128_u64};
import("std::bytes::errors")::{BYTES_ERROR_NO_SPACE};

fun main() {
    var data: array<u8, 1> = [85];
    var writer: ByteWriter = bytes_writer(&data[0], 1);
    var status: i32 = bytes_writer_write_uleb128_u64(&writer, 128);

    if (status == BYTES_ERROR_NO_SPACE) {
        println("position={} byte={}", writer.position, data[0]);
    }
}
```

Kết quả thực hiện:

```text
position=0 byte=85
```

128 yêu cầu hai byte nhưng chỉ có một khoảng trắng. Sau khi thất bại, byte 85 đầu tiên vẫn còn. Việc kiểm tra các loại lỗi và bảo toàn trạng thái cùng nhau mô tả ranh giới tốt hơn so với kiểm tra thành công roundtrip đơn giản.

## Trình tự tạo trình phân tích cú pháp tin nhắn

1. Đọc tiêu đề cố định và kiểm tra loại và phiên bản.
2. Đọc độ dài và so sánh nó với phạm vi đầu vào còn lại.
3. Chỉ truyền dữ liệu cần thiết tới view hoặc một bộ đệm riêng.
4. Nếu định dạng yêu cầu toàn bộ tin nhắn thì các byte bổ sung cũng sẽ được kiểm tra.
5. Phân biệt giữa EOF và lỗi định dạng không hợp lệ và chuyển chúng cho người gọi.

Bạn có thể tạo một chương trình bằng cách kết nối các trường trong [Thực hành tin nhắn nhị phân](/docs/vi/practice/binary-message).
