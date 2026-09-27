---
translation_set_id: learn-errors
path: language/errors
locale: vi
group: language
group_order: 2
order: 12
title: 12. Đại diện và khắc phục lỗi
summary: Tách lỗi khỏi các giá trị bình thường và dọn sạch tài nguyên khỏi đường dẫn lỗi.
---

## Kết quả của chức năng mức độ thất bại

Việc thiếu tệp hoặc đầu vào vượt quá giới hạn xảy ra một cách tự nhiên trong các chương trình. Xử lý lỗi không chỉ đơn thuần là in một tin nhắn. Đó là quá trình phân biệt các lỗi, kiểm tra trạng thái công việc đã hoàn thành, tổ chức các nguồn lực có được và sau đó chọn tiếp tục hay chấm dứt.

Trong chương này, chúng ta bắt đầu với việc biểu diễn lỗi của một hàm nhỏ và tiến tới cấu trúc kết quả, variant, trả về sớm và dọn dẹp tài nguyên.

## Các điểm đánh dấu lỗi không được chồng lên các giá trị thành công

Lý do -1 được sử dụng vì không tìm thấy trong tìm kiếm mảng là vì chỉ mục hiệu quả là 0 trở lên. Mặt khác, trong các phép tính mà bất kỳ số nguyên nào cũng có thể là kết quả bình thường, nếu -1 được chỉ định là lỗi thì không thể phân biệt được nó với giá trị bình thường -1.

0 cũng là một giá trị thường bị hiểu lầm. Độ dài chuỗi trống 0, vị trí đầu tiên 0 và số byte được truyền 0 có ý nghĩa khác nhau đối với các hàm khác nhau. Đừng đánh giá thành công chỉ vì giá trị trả về khác 0.

## Trả lại thành công và giá trị cùng nhau

<!-- wave-example: book-error-result -->
```wave
struct Division {
    ok: bool;
    value: i32;
}

fun divide_nonnegative(left: i32, right: i32) -> Division {
    if (left < 0 || right <= 0) {
        return Division {
            ok: false,
            value: 0
        };
    }

    return Division {
        ok: true,
        value: left / right
    };
}

fun main() {
    var result: Division = divide_nonnegative(0, 3);

    if (!result.ok) {
        println("invalid input");
        return;
    }

    println("value={}", result.value);
}
```

Kết quả thực hiện:

```text
value=0
```

Kết quả bình thường cũng có thể là 0. Thay vì nhìn vào value và đoán xem nó có thành công hay không, trước tiên hãy kiểm tra ok. Ngay cả khi trường value tồn tại trong kết quả không thành công, điều đó không có nghĩa đó là giá trị được sử dụng.

Quy tắc đầu vào cho ví dụ này là vế trái bằng 0 hoặc lớn hơn và vế phải là dương. Phạm vi được hiển thị trong tên và mô tả hàm để phân biệt với phép chia số nguyên signed điển hình.

## Phân biệt nguyên nhân gây ra lỗi

Thêm thông tin lỗi để được hướng dẫn bổ sung hoặc phục hồi tùy theo lý do lỗi. Ngoài ra còn có cách chia hàm thành hàm kiểm tra phạm vi đầu vào và hàm tính toán.

<!-- wave-example: book-error-codes -->
```wave
struct Check {
    ok: bool;
    error: i32;
}

fun check_quantity(value: i32) -> Check {
    if (value < 1) {
        return Check {
            ok: false,
            error: 1
        };
    }

    if (value > 100) {
        return Check {
            ok: false,
            error: 2
        };
    }

    return Check {
        ok: true,
        error: 0
    };
}

fun main() {
    var result: Check = check_quantity(101);

    if (!result.ok) {
        match (result.error) {
            1 => {
                println("quantity is too small");
            }
            2 => {
                println("quantity is too large");
            }
            _ => {
                println("invalid quantity");
            }
        }
    }
}
```

Kết quả thực hiện:

```text
quantity is too large
```

Ý nghĩa của các con số được xác định bởi chức năng này. Nó không thể được coi là giống như lỗi số 1 hoặc 2 trong các thư viện khác. Ở chế độ công khai API, việc đặt tên các hằng hoặc loại lỗi giúp người gọi không phải ghi nhớ các số ngẫu nhiên.

## Tách các kết quả bằng variant

Mối quan hệ mà giá trị thành công và lỗi không thể tồn tại cùng lúc có thể được biểu thị bằng variant.

<!-- wave-example: book-error-variant -->
```wave
variant ByteResult {
    Value(u8),
    OutOfRange(i32)
}

fun to_byte(value: i32) -> ByteResult {
    if (value < 0 || value > 255) {
        return ByteResult::OutOfRange(value);
    }

    return ByteResult::Value(value as u8);
}

fun main() {
    var result: ByteResult = to_byte(300);

    match (result) {
        ByteResult::Value(value) => {
            println("byte={}", value);
        }
        ByteResult::OutOfRange(original) => {
            println("out of range: {}", original);
        }
    }
}
```

Kết quả thực hiện:

```text
out of range: 300
```

Lỗi chứa giá trị đầu vào ban đầu. Người gọi sẽ dễ dàng giải thích vấn đề hơn là chỉ quay lại false. Bản chất của dữ liệu cũng được xem xét bằng cách không để lại các thông tin đầu vào nhạy cảm như mật khẩu hoặc mã thông báo trong nhật ký.

## Làm cho các tuyến đường thông thường dễ đọc hơn bằng cách quay lại sớm

Không cần phải đặt tất cả mã tốt vào sâu bên trong if khi thực hiện kiểm tra nhiều bước. Nếu thất bại, bạn có thể quay lại ngay lập tức và tiếp tục con đường thành công bên dưới.

<!-- wave-example: book-error-early-return -->
```wave
fun show_total(quantity: i32, price: i32) {
    if (quantity < 1 || quantity > 100) {
        println("invalid quantity");
        return;
    }

    if (price < 0 || price > 100000) {
        println("invalid price");
        return;
    }

    var total: i32 = quantity * price;
    println("total={}", total);
}

fun main() {
    show_total(3, 1200);
    show_total(0, 1200);
    show_total(3, -1);
}
```

Kết quả thực hiện:

```text
total=3600
invalid quantity
invalid price
```

Sau khi vượt qua bài kiểm tra, bạn có thể tận dụng thực tế là quantity và price nằm trong phạm vi được chỉ định. Phạm vi được đặt sao cho phép nhân trung gian cũng nằm trong i32. Khi thêm lợi tức sớm, bạn cũng nên kiểm tra xem liệu bạn đã sở hữu bất kỳ tài nguyên nào vào thời điểm đó hay chưa.

## Dọn dẹp bộ nhớ của các đường dẫn lỗi

<!-- wave-example: book-error-cleanup -->
```wave
import("std::buffer::alloc")::{
    Buffer, buffer_init, buffer_free
};
import("std::buffer::write")::{
    buffer_append_str
};

fun make_message(allowed: bool) -> i32 {
    var data: Buffer;

    if (buffer_init(&data, 16) < 0) {
        return 1;
    }

    if (!allowed) {
        buffer_free(&data);
        return 2;
    }

    if (buffer_append_str(&data, "ready") < 0) {
        buffer_free(&data);
        return 3;
    }

    println("bytes={}", data.len);

    if (buffer_free(&data) < 0) {
        return 4;
    }

    return 0;
}

fun main() {
    println("status={}", make_message(false));
}
```

Kết quả thực hiện:

```text
status=2
```

Phát hành Buffer ngay cả trên các đường dẫn lỗi không xuất ra dữ liệu. Thay vì biến mọi thất bại thành một return duy nhất, điều quan trọng là phải quản lý rõ ràng những gì thuộc sở hữu của mỗi chi nhánh.

Ngoài ra còn có API trong đó quá trình dọn dẹp không thành công. Quyết định cách bạn sẽ bảo tồn các lỗi và dọn dẹp các lỗi trong tác phẩm gốc. Ví dụ nhỏ này trước tiên trả về số lỗi của tác vụ. Các chương trình lớn hơn có thể ghi cả hai một cách riêng biệt.

## Thành công một phần không tự động bị hủy

Nếu bạn ghi một số byte vào một tệp và sau đó ghi không thành công thì các byte đã ghi sẽ không bị mất. Đầu bên kia của mạng cũng có thể đã nhận được một số dữ liệu. Việc lặp lại cùng một tác vụ từ đầu có thể dẫn đến các bản ghi trùng lặp.

Ngược lại, việc đọc checked byte cursor bảo toàn vị trí và giá trị đầu ra khi lỗi. Các chức năng này có thể nhận thêm đầu vào và thử lại ở cùng một vị trí. Thay vì áp dụng quy tắc “nếu thất bại, không có gì thay đổi” cho mọi API, hãy kiểm tra tài liệu về chức năng đó.

## Lỗi và bẫy có thể phục hồi

Đường dẫn tệp không hợp lệ hoặc đầu vào của người dùng không hợp lệ có thể được báo cáo dưới dạng giá trị lỗi để người gọi có thể khôi phục. Một cái bẫy gây ra bởi số lần thay đổi thời gian chạy không hợp lệ hoặc chuyển đổi dấu phẩy động thành số nguyên là một cơ chế khác.

Các chương trình phải tiếp tục thực thi phải kiểm tra đầu vào của chúng trước khi thực hiện các hoạt động nguy hiểm. assert cũng không nhằm mục đích bị lạm dụng như một phương tiện xử lý các lỗi đầu vào của người dùng. Nếu bạn cần cho người dùng cơ hội nhập lại dữ liệu đầu vào của họ, hãy trả về kết quả để tiếp tục quy trình điều khiển.

## Bài tập và lời giải đầy đủ

Viết hàm đọc một phần tử mảng theo chỉ mục và không thành công khi chỉ mục đó âm hoặc lớn hơn hoặc bằng độ dài. Sử dụng cấu trúc kết quả vì số 0 có thể là giá trị phần tử hợp lệ.

<!-- wave-example: book-error-solution -->
```wave
struct Lookup {
    ok: bool;
    value: i32;
}

fun at(values: ptr<i32>, length: i32, index: i32) -> Lookup {
    if (index < 0 || index >= length) {
        return Lookup {
            ok: false,
            value: 0
        };
    }

    return Lookup {
        ok: true,
        value: deref values[index]
    };
}

fun main() {
    var values: array<i32, 3> = [0, 10, 20];
    var first: Lookup = at(&values[0], 3, 0);
    var outside: Lookup = at(&values[0], 3, 3);

    if (first.ok) {
        println("value={}", first.value);
    }

    if (!outside.ok) {
        println("out of bounds");
    }
}
```

Kết quả thực hiện:

```text
value=0
out of bounds
```

Điều kiện của người gọi là con trỏ và length đại diện cho mảng thực tế có thể đọc được. Đây không phải là chức năng giúp các địa chỉ ngẫu nhiên trở nên an toàn chỉ bằng cách kiểm tra chỉ mục. Vui lòng đọc riêng nội dung kiểm tra mà hàm chịu trách nhiệm và những điều kiện mà người gọi phải đảm bảo.
