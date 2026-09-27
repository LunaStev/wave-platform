---
translation_set_id: learn-allocation
path: language/allocation
locale: vi
group: language
group_order: 2
order: 10
title: 10. Phân bổ bộ nhớ và quản lý tài nguyên
summary: Tìm hiểu về lỗi phân bổ, khởi tạo, phạm vi và giải phóng.
---

## Khi nào bạn cần không gian lưu trữ động?

Mảng có kích thước cố định bao gồm độ dài của nó trong kiểu của nó. Sử dụng bộ nhớ động khi lượng dữ liệu chỉ được biết trong thời gian chạy, chẳng hạn như kích thước tệp hoặc độ dài đầu vào. Giải phóng mỗi lần phân bổ khi bạn không còn cần đến nó nữa.

Trong chương này, bạn sẽ quản lý một phân bổ nhỏ, thay đổi kích thước của nó và sau đó sử dụng Buffer. Truyền con trỏ khác với chuyển quyền sở hữu. Xác định những tài nguyên mà mỗi chức năng sở hữu khi bạn làm theo các ví dụ.

## Phân bổ, kiểm tra, sử dụng và miễn phí

<!-- wave-example: book-alloc-lifecycle -->
```wave
import("std::mem::alloc")::{
    mem_alloc_zeroed, mem_free
};

fun main() -> i32 {
    var size: i64 = 4;
    var data: ptr<u8> = mem_alloc_zeroed(size);

    if (data == null) {
        println("allocation failed");
        return 1;
    }

    deref data[0] = 42;
    println("{} {}", deref data[0], deref data[1]);

    var status: i64 = mem_free(data, size);

    if (status < 0) {
        return 2;
    }

    return 0;
}
```

Kết quả thực hiện:

```text
42 0
```

Chương trình có bốn bước: yêu cầu 4 byte, kiểm tra null, chỉ truy cập phạm vi hợp lệ và giải phóng phân bổ. Khi thành công, mem_alloc_zeroed khởi tạo bộ nhớ về 0, do đó byte thứ hai bằng 0 ngay cả khi chương trình chưa ghi vào nó.

Không giả sử bất kỳ nội dung ban đầu nào cho bộ nhớ được trả về bởi mem_alloc. Khởi tạo từng vùng trước khi đọc nó. Kích thước phân bổ bằng 0 hoặc âm trả về null. Việc phân bổ kích thước dương cũng có thể không lấy được bộ nhớ.

## đơn vị kích thước

Đối số kích thước của API cấp phát bộ nhớ được đo bằng byte. Để phân bổ mười số nguyên, nhân kích thước phần tử với số phần tử. Kiểm tra phép nhân này không bị tràn.

<!-- wave-example: book-alloc-typed -->
```wave
import("std::mem::alloc")::{
    mem_alloc, mem_free
};
import("std::mem::layout")::{
    size_of
};
import("std::mem::ops")::{
    mem_size_mul_checked
};

fun main() -> i32 {
    var count: i64 = 3;
    var bytes: i64 = 0;
    var item_size: i64 = size_of<i32>() as i64;

    if (mem_size_mul_checked(count, item_size, &bytes) < 0) {
        return 1;
    }

    var raw: ptr<u8> = mem_alloc(bytes);

    if (raw == null) {
        return 2;
    }

    var values: ptr<i32> = raw as ptr<i32>;

    for (var index: i64 = 0; index < count; index += 1) {
        deref values[index] = (index + 1) as i32;
    }

    println("{} {} {}", deref values[0], deref values[1], deref values[2]);

    if (mem_free(raw, bytes) < 0) {
        return 3;
    }

    return 0;
}
```

Kết quả thực hiện:

```text
1 2 3
```

count là số phần tử; byte là số byte. Số học con trỏ di chuyển theo đơn vị i32, nhưng việc giải phóng phân bổ yêu cầu kích thước ban đầu của nó tính bằng byte. size_of sử dụng bố cục loại mục tiêu, làm cho mối quan hệ với loại phần tử trở nên rõ ràng.

Khi thay đổi kết quả của size_of thành i64 đối với loại chung rất lớn, phạm vi chuyển đổi cũng phải được xem xét. Ở đây chúng tôi sử dụng i32, có kích thước đã biết.

## Dọn dẹp ngay cả trên đường dẫn lỗi

Nếu một thao tác khác không thành công sau khi cấp phát, hãy giải phóng bộ nhớ trước khi quay lại sớm. Bảng quyền sở hữu giúp xác định các đường dẫn dọn dẹp mà bạn có thể bỏ lỡ.

|Bước|Tài nguyên thuộc sở hữu của|Nếu bạn thất bại|
| --- | --- | --- |
|Trước khi phân bổ|không có|trả lại ngay|
|Sau khi phân bổ thành công|data và kích thước gốc|data Trở lại sau khi ra mắt|
|Sau khi phân bổ lại thành công|Địa chỉ mới và kích thước mới|Phát hành địa chỉ mới|
|Sau khi phát hành|không có|Không sử dụng địa chỉ cũ|

Việc ghi đè một biến con trỏ và làm mất địa chỉ ban đầu cũng làm mất thông tin cần thiết để giải phóng việc phân bổ. Điều này gây ra rò rỉ bộ nhớ. Ngược lại, việc giải phóng cùng một khoản phân bổ thông qua hai chủ sở hữu sẽ gây ra sự miễn phí gấp đôi.

## Tái phân bổ để tăng kích thước

<!-- wave-example: book-alloc-grow -->
```wave
import("std::mem::alloc")::{
    mem_alloc, mem_realloc, mem_free
};

fun main() -> i32 {
    var data: ptr<u8> = mem_alloc(4);

    if (data == null) {
        return 1;
    }

    deref data[0] = 7;
    var next: ptr<u8> = mem_realloc(data, 4, 8);

    if (next == null) {
        mem_free(data, 4);
        return 2;
    }

    data = next;
    deref data[4] = 9;
    println("{} {}", deref data[0], deref data[4]);

    if (mem_free(data, 8) < 0) {
        return 3;
    }

    return 0;
}
```

Kết quả thực hiện:

```text
7 9
```

Đầu tiên lưu trữ kết quả trong phần tiếp theo. Nếu việc phân bổ khối lớn hơn không thành công, dữ liệu gốc vẫn hợp lệ và vẫn có thể được giải phóng. Nếu thành công, việc phân bổ cũ sẽ được giải phóng và địa chỉ mới phải được sử dụng. Khởi tạo vùng mới được thêm vào trước khi đọc nó.

Kích thước mới trong ví dụ này là dương. Thay vào đó, một yêu cầu có new_size=0 sẽ cố gắng giải phóng phân bổ hiện có và trả về null. Do đó, kết quả null không phải lúc nào cũng có nghĩa là phân bổ cũ vẫn hợp lệ. Gọi trực tiếp mem_free khi bạn cần kiểm tra xem việc giải phóng có thành công hay không.

## Khi nào cần kiểm tra lại con trỏ đã mượn

data Việc lưu con trỏ nội bộ và sử dụng nó sau khi phân bổ lại là sai. Điều này là do địa chỉ của data mới có thể khác. Nếu cần vị trí nội bộ, bạn có thể lưu trữ offset thay vì địa chỉ và tính toán lại dựa trên data mới sau khi thành công.

Việc giải phóng hoặc phân bổ lại bộ nhớ cũng ảnh hưởng đến mã đã mượn nó. Kiểm tra xem một thao tác khác có còn sử dụng bộ nhớ đó hay không. Bộ đệm được chuyển đến thao tác không đồng bộ phải duy trì hiệu lực cho đến khi thao tác kết thúc.

## Danh sách byte chứa Buffer

Việc quản lý danh sách byte có độ dài thay đổi thường xuyên trong khi phân bổ lại theo cách thủ công yêu cầu xử lý cả len và cap, lỗi mở rộng và tính toán kích thước. Buffer của std gộp các thao tác này lại.

<!-- wave-example: book-buffer-progressive -->
```wave
import("std::buffer::alloc")::{
    Buffer, buffer_init, buffer_free
};
import("std::buffer::write")::{
    buffer_append_str
};

fun main() -> i32 {
    var message: Buffer;

    if (buffer_init(&message, 0) < 0) {
        return 1;
    }

    if (buffer_append_str(&message, "Hello") < 0) {
        buffer_free(&message);
        return 2;
    }

    if (buffer_append_str(&message, ", Wave") < 0) {
        buffer_free(&message);
        return 3;
    }

    println("bytes={}", message.len);

    if (buffer_free(&message) < 0) {
        return 4;
    }

    return 0;
}
```

Kết quả thực hiện:

```text
bytes=11
```

len là số byte đang sử dụng; cap là công suất được phân bổ. Việc thêm dữ liệu sẽ tăng phân bổ khi cần thiết. Việc sử dụng Buffer không loại bỏ trách nhiệm giải phóng nó của người gọi.

buffer_append_str không thêm NUL vào cuối chuỗi. Do đó, message.data không được xuất trực tiếp dưới dạng str. Byte được xuất ra bằng hàm I/O, hàm này lấy độ dài hoặc tạo biểu diễn chuỗi một cách rõ ràng.

## Bài tập và hướng giải

Thêm từng byte từ 0 đến 9 vào Buffer và nhận tổng. Nó phải được giải phóng khi mỗi lần nối thêm không thành công và việc đọc chỉ được thực hiện trong phạm vi len. Bạn có thể kiểm tra giải pháp đầy đủ và lỗi biên bằng cách làm theo ví dụ trong [Buffer Cách sử dụng](/docs/vi/stdlib/buffer).

Hãy thử đánh dấu các lệnh gọi phân bổ, tái phân bổ và phân bổ trong mã của bạn. Đối với mỗi lần phân bổ thành công, bạn phải có khả năng mô tả ai sở hữu nó và đường dẫn nào sẽ giải phóng nó.
