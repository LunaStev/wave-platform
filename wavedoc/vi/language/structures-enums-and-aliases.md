---
translation_set_id: data-types
path: language/structures-enums-and-aliases
locale: vi
group: language
group_order: 2
order: 8
title: 8. Cấu trúc, enum và biến thể
summary: Tìm hiểu về các trường, khởi tạo cấu trúc và vai trò của enum và variant.
---

## Thể hiện mối quan hệ giữa các dữ liệu dưới dạng các loại

Nếu giá và số lượng sản phẩm được truyền dưới dạng biến thì rất khó để biết chỉ bằng cách nhìn vào mã liệu hai giá trị có thuộc về cùng một sản phẩm hay không. Cấu trúc nhóm các trường liên quan. enum đại diện cho một trạng thái được đặt tên và variant lưu trữ các dữ liệu khác nhau cùng nhau cho từng trường hợp.

Ba hàm này không phải là cú pháp thay thế lẫn nhau. Chọn tùy thuộc vào những gì bạn muốn thể hiện.

|một cái gì đó để thể hiện|chọn|vâng|
| --- | --- | --- |
|Nhiều trường tồn tại đồng thời| struct |Đơn giá và số lượng sản phẩm|
|Trạng thái số nguyên được đặt tên| enum |Chờ/Tiếp tục/Hoàn thành|
|Dữ liệu khác nhau trong từng trường hợp| variant |giá trị thành công hoặc lỗi|

## Khai báo cấu trúc và tạo giá trị

<!-- wave-example: book-struct-first -->
```wave
struct Product {
    price: i32;
    quantity: i32;
}

fun main() {
    var item: Product = Product {
        price: 1500,
        quantity: 2
    };

    println("{} {}", item.price, item.quantity);
}
```

Kết quả thực hiện:

```text
1500 2
```

Các trường trong phần khai báo kết thúc bằng dấu chấm phẩy và khi tạo giá trị, các trường và giá trị được kết nối bằng dấu hai chấm và phân tách bằng dấu phẩy. Xác định loại và tạo giá trị thực tế là hai bước khác nhau. Khai báo loại Product không tự động tạo dung lượng lưu trữ cho một sản phẩm.

Trường được truy cập dưới dạng `item.price`. Nếu bạn tạo nhiều item cùng loại, bạn có thể lưu trữ các giá trị khác nhau trong mỗi loại.

## Truyền cấu trúc cho hàm

<!-- wave-example: book-struct-function -->
```wave
struct Product {
    price: i32;
    quantity: i32;
}

fun subtotal(item: Product) -> i32 {
    return item.price * item.quantity;
}

fun main() {
    var item: Product = Product {
        price: 1500,
        quantity: 2
    };

    println("{}", subtotal(item));

    item.quantity = 3;
    println("{}", subtotal(item));
}
```

Kết quả thực hiện:

```text
3000
4500
```

Hàm nhận được dưới dạng mối quan hệ giữa đơn giá và số lượng thuộc về cùng một sản phẩm. Đây là hàm đọc cấu trúc trường số nguyên được truyền dưới dạng giá trị. Bất kỳ chức năng nào muốn sửa đổi bộ nhớ của người gọi đều có thể được thiết kế để chấp nhận một con trỏ.

Nếu cấu trúc có trường con trỏ, giá trị sao chép cũng sao chép địa chỉ. Đây không phải là chức năng sao chép sâu để phân bổ riêng biệt. Các loại chứa tài nguyên như phần xử lý tệp hoặc Buffer phải chỉ định các quy tắc sao chép và phát hành cùng nhau.

## Phương pháp và proto

Các hàm liên quan có thể được nhóm lại dưới dạng phương thức. proto là phương pháp viết phương thức của cấu trúc dưới dạng một khối riêng biệt.

<!-- wave-example: book-struct-method -->
```wave
struct Counter {
    value: i32;
}

proto Counter {
    fun current(self: Counter) -> i32 {
        return self.value;
    }
}

fun main() {
    var counter: Counter = Counter {
        value: 3
    };

    println("{}", counter.current());
}
```

Kết quả thực hiện:

```text
3
```

`self: Counter` là tham số nhận giá trị. Việc sử dụng ký hiệu gọi phương thức không tự động biến nó thành phương thức sửa đổi bản gốc. Mời các bạn cùng đọc loại self và chức năng của nó trong văn bản.

Bạn không thể mong đợi một trường chỉ có trạng thái hợp lệ chỉ vì bạn đã đính kèm một phương thức vào trường đó. Nếu có những kết hợp không hợp lệ mà người dùng có thể tạo bằng các trường công khai, thì hàm của bạn nên kiểm tra chúng hoặc cung cấp quy tắc tạo.

## Đặt tên tiểu bang bằng enum

<!-- wave-example: book-enum-state -->
```wave
enum State -> i32 {
    Ready = 0,
    Running,
    Finished,
}

fun main() {
    var state: State = State::Ready;

    if (state == State::Ready) {
        println("ready");
    }

    state = State::Running;

    if (state == State::Running) {
        println("running");
    }
}
```

Kết quả thực hiện:

```text
ready
running
```

`-> i32` là loại số nguyên được sử dụng trong biểu thức. Giá trị đầu tiên được đặt thành 0 và các giá trị bị bỏ qua tiếp theo lớn hơn giá trị trước đó 1. Thay vì chỉ so sánh 0 và 1 trong mã của bạn, việc sử dụng State::Ready và State::Running sẽ tiết lộ ý nghĩa.

enum Việc có tên không tự động hạn chế việc chuyển đổi trạng thái. Các quy tắc như có thể quay lại từ Finished đến Running phải được triển khai dưới dạng hàm hay không.

## Kết nối các trường hợp và dữ liệu với variant

Nếu chỉ có một giá trị trong trường hợp thành công và thông tin lỗi là cần thiết trong trường hợp thất bại, thì giá trị đó có thể được biểu thị dưới dạng variant.

<!-- wave-example: book-variant-result -->
```wave
variant Result {
    Value(i32),
    Error(i32)
}

fun divide(left: i32, right: i32) -> Result {
    if (left < 0 || right <= 0) {
        return Result::Error(1);
    }

    return Result::Value(left / right);
}

fun show(result: Result) {
    match (result) {
        Result::Value(value) => {
            println("value={}", value);
        }
        Result::Error(code) => {
            println("error={}", code);
        }
    }
}

fun main() {
    show(divide(12, 3));
    show(divide(12, 0));
}
```

Kết quả thực hiện:

```text
value=4
error=1
```

Result::Value và Result::Error đều chứa payload. Ngay cả khi chúng chứa cùng loại số nguyên, một số trường hợp nhất định vẫn được phân biệt. Người gọi sẽ kiểm tra trường hợp bằng match và sử dụng payload trong đó arm.

divide trên là một ví dụ nhỏ không hỗ trợ toán hạng âm. Vì phạm vi đầu vào được chỉ định nên không nên nhầm lẫn nó với một hàm bao trùm tất cả các ranh giới của phép chia signed thông thường.

## Nếu payload không có mặt

Dữ liệu không cần thiết trong mọi trường hợp. Trạng thái không có giá trị có thể được biểu diễn dưới dạng một trường hợp riêng biệt.

<!-- wave-example: book-variant-empty -->
```wave
variant Lookup {
    Found(i32),
    Missing
}

fun main() {
    var result: Lookup = Lookup::Missing;

    match (result) {
        Lookup::Found(value) => {
            println("found={}", value);
        }
        Lookup::Missing => {
            println("missing");
        }
    }
}
```

Kết quả thực hiện:

```text
missing
```

Thay vì đặt trước một số nguyên tùy ý là “none”, chúng tôi đã sử dụng trường hợp Missing. Dù giá trị thành công là bao nhiêu thì ý nghĩa cũng không trùng lặp.

match đến `_` xử lý các trường hợp còn lại. Nếu bạn muốn mỗi người gọi được xem xét lại khi một trường hợp mới được thêm vào, tốt hơn hết bạn nên tách biệt tất cả các trường hợp một cách rõ ràng. Cho dù bạn chọn phương pháp nào, hãy đảm bảo không có dữ liệu đầu vào chưa được xử lý.

## Sử dụng cấu trúc và variant cùng nhau

Việc chọn dữ liệu khác nhau có thể được biểu thị dưới dạng variant và việc nhóm nhiều trường thuộc một trường hợp có thể được biểu thị dưới dạng cấu trúc. Ví dụ: nếu kết quả xử lý đơn hàng thành công, nó có thể được thiết kế để chứa cấu trúc biên nhận và nếu thất bại, nó có thể được thiết kế để chứa một số lỗi.

Quy tắc tuổi thọ bộ nhớ không biến mất ngay cả khi giá trị chứa các giá trị khác. Nếu một con trỏ được lưu trữ trong variant, liệu con trỏ có hợp lệ hay không và ai sẽ giải phóng nó sẽ được xác định riêng. Ngay cả khi lưu ở định dạng tệp bên ngoài, bạn phải chỉ định mã hóa cho từng trường thay vì bỏ bộ nhớ cấu trúc như hiện tại.

## Bài tập và lời giải đầy đủ

Tạo phán quyết chỉ chấp nhận điểm từ 0 đến 100. Điểm hợp lệ được biểu thị dưới dạng Grade(score), phần còn lại được biểu thị dưới dạng Invalid.

<!-- wave-example: book-model-solution -->
```wave
variant CheckedScore {
    Grade(i32),
    Invalid
}

fun check_score(score: i32) -> CheckedScore {
    if (score < 0 || score > 100) {
        return CheckedScore::Invalid;
    }

    return CheckedScore::Grade(score);
}

fun main() {
    var result: CheckedScore = check_score(87);

    match (result) {
        CheckedScore::Grade(score) => {
            println("accepted={}", score);
        }
        CheckedScore::Invalid => {
            println("invalid score");
        }
    }
}
```

Kết quả thực hiện:

```text
accepted=87
```

Thay đổi đầu vào thành -1, 0, 100, 101 để kiểm tra ranh giới. Vì thành công và thất bại không chia sẻ cùng một không gian số nguyên nên người gọi sẽ giảm nguy cơ vô tình thêm các giá trị lỗi vào phép tính trung bình.


## Tôi nên chọn cách diễn đạt nào?

|hình dạng dữ liệu|cách diễn đạt phù hợp|vâng|
| --- | --- | --- |
|Nhiều giá trị cùng loại|Mảng|10 điểm|
|Nhiều lĩnh vực liên quan đến nhau|cấu trúc|tên và điểm|
|Giá trị trạng thái được đặt tên| enum | Ready, Running, Stopped |
|Dữ liệu bổ sung thay đổi theo tiểu bang| variant | Value(i32), Error(str) |
|Tên theo ngữ cảnh cho một loại hiện có|gõ bí danh| UserId = u64 |

Khi chọn cấu trúc dữ liệu, hãy xem xét không chỉ các giá trị bạn muốn lưu trữ mà còn cả những trạng thái sai mà bạn có thể biểu thị. Cấu trúc có cả trường thành công/thất bại và trường giá trị/lỗi có thể tạo ra sự kết hợp không chính xác, nhưng variant có thể được phân biệt thành payload cho từng trường hợp.

## Vị trí bộ nhớ và dữ liệu bên ngoài

Bộ nhớ của cấu trúc có thể chứa không gian trống để đảm bảo sự liên kết giữa các trường. Việc thêm kích thước trường đơn giản không phải lúc nào cũng bằng tổng kích thước của cấu trúc. Khi bạn cần biết kích thước và căn chỉnh, hãy sử dụng [mem Chức năng bố cục](/docs/vi/reference/memory-and-buffer).

Đối với các tệp hoặc tin nhắn mạng, thứ tự và độ dài byte có thể được xác định rõ ràng bằng cách mã hóa các trường theo thứ tự bằng hàm [bytes](/docs/vi/stdlib/bytes). Khi chuyển cấu trúc sang các ngôn ngữ khác, hãy căn chỉnh khai báo bên ngoài của [FFI](/docs/vi/language/modules-imports-and-ffi) với mục tiêu ABI.

[Cấu trúc học tập và thực hành](/docs/vi/language/structures-enums-and-aliases) · [variant](/docs/vi/language/variants)
