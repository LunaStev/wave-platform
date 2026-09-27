---
translation_set_id: functions
path: language/functions-and-generics
locale: vi
group: language
group_order: 2
order: 5
title: 5. Thiết kế và soạn thảo chức năng
summary: Tìm hiểu các tham số, giá trị trả về, giá trị mặc định và truyền theo giá trị.
---

## Bắt đầu từ mã lặp đi lặp lại

Hàm là một công cụ để rút gọn cú pháp, nhưng chúng cũng là một công cụ để phân định các nhiệm vụ. Việc tách những gì nó lấy làm đầu vào, những gì nó tính toán và những kết quả nó trả về cho phép bạn hiểu chương trình của mình theo từng phần nhỏ hơn.

Trong chương này, chúng ta bắt đầu với một chương trình viết các phép tính chiết khấu nhiều lần. Mỗi ví dụ hoàn chỉnh trong main.wave và chạy dưới dạng `wavec run main.wave`. Chúng tôi chưa phân chia các tập tin.

<!-- wave-example: book-function-before -->
```wave
fun main() {
    var first_price: i32 = 2000;
    var first_discount: i32 = first_price * 10 / 100;
    var first_total: i32 = first_price - first_discount;
    var second_price: i32 = 5000;
    var second_discount: i32 = second_price * 10 / 100;
    var second_total: i32 = second_price - second_discount;
    println("{} {}", first_total, second_total);
}
```

Kết quả thực hiện:

```text
1800 4500
```

Hai phép tính chỉ khác nhau về giá và có cấu trúc giống nhau. Khi thay đổi quy tắc giảm giá, bạn cần chỉnh sửa cả hai địa điểm. Nếu chỉ thay đổi một bên, bạn sẽ nhận được kết quả khác nhau đối với những sản phẩm yêu cầu cùng chính sách.

## Xác định đầu vào và đầu ra

Di chuyển các phép tính dư thừa sang hàm. Giá trị đã thay đổi được nhận dưới dạng đầu vào price và giá đã tính sẽ được trả về.

<!-- wave-example: book-function-extract -->
```wave
fun discounted(price: i32) -> i32 {
    var discount: i32 = price * 10 / 100;
    return price - discount;
}

fun main() {
    println("{} {}", discounted(2000), discounted(5000));
}
```

Kết quả thực hiện:

```text
1800 4500
```

Tên hàm là discounted. `price: i32` trong ngoặc đơn là khai báo tham số và `-> i32` là loại kết quả. Biến cục bộ discount trong nội dung chỉ được sử dụng trong hàm này.

`discounted(2000)` là biểu thức gọi hàm. 2000 trong ngoặc đơn là đối số bạn thực sự vượt qua. Vì giá trị được hàm trả về trở thành kết quả của biểu thức gọi này nên nó có thể được sử dụng trực tiếp làm đối số cho println.

|thuật ngữ|mã|ý nghĩa|
| --- | --- | --- |
|tham số| price |Tên đầu vào được chỉ định khi khai báo hàm|
|yếu tố| 2000 |Giá trị được chuyển khi gọi|
|kiểu trả về| i32 |Loại giá trị được tạo bởi biểu thức cuộc gọi|
|tuyên bố trở lại| return price - discount |Chuyển kết quả và kết thúc cuộc gọi này|

## Thực hiện theo thứ tự gọi và thực hiện

Chỉ viết một khai báo hàm trong nguồn không khiến phần thân của nó được thực thi ngay lập tức. Thực thi khi đạt đến điểm được gọi trong main.

<!-- wave-example: book-function-trace -->
```wave
fun calculate(value: i32) -> i32 {
    println("inside: {}", value);
    return value * 2;
}

fun main() {
    println("before");
    var result: i32 = calculate(7);
    println("after: {}", result);
}
```

Kết quả thực hiện:

```text
before
inside: 7
after: 14
```

Thứ tự tiến trình là đầu ra đầu tiên của main, văn bản chính của calculate và đầu ra cuối cùng của main. Khi `return` được thực thi, lệnh gọi tới calculate kết thúc với kết quả 14 và quá trình khởi tạo result của main hoàn tất.

Nếu nhiều lệnh gọi hàm được kết hợp trong một biểu thức và thứ tự của các tác dụng phụ là quan trọng, hãy chia các lệnh gọi thành các câu lệnh riêng biệt. Các ví dụ trong chương này cũng lưu trữ kết quả của các lệnh gọi cần được theo dõi trong các biến cục bộ.

## nhiều tham số

Nếu bạn cũng nhận được tỷ lệ chiết khấu làm đầu vào, bạn có thể tính toán nhiều chính sách có cùng chức năng.

<!-- wave-example: book-function-parameters -->
```wave
fun discounted(price: i32, percent: i32) -> i32 {
    return price - price * percent / 100;
}

fun main() {
    var standard: i32 = discounted(2000, 10);
    var special: i32 = discounted(2000, 25);
    println("standard={}", standard);
    println("special={}", special);
}
```

Kết quả thực hiện:

```text
standard=1800
special=1500
```

Thứ tự của các đối số phải khớp với khai báo. Nếu cả hai tham số đều là i32, rất khó để phân biệt ý nghĩa của chúng chỉ bằng cách kiểm tra loại, ngay cả khi thứ tự bị thay đổi. Xác định rõ ràng tên hàm và tên tham số, đồng thời ghi vị trí cuộc gọi theo cách dễ đọc.

Hàm này giả định một số tiền nhỏ và tỷ lệ chiết khấu hợp lệ. Không xử lý giá âm, tỷ lệ lớn hơn 100 hoặc tràn nhân giữa. Khi tạo một hàm, bạn phải mô tả không chỉ phần thân mà còn cả các điều kiện đầu vào. Chương trình hoàn thiện sau này sẽ tách biệt các bước kiểm tra.

## đối số mặc định

Bạn có thể cung cấp các giá trị được sử dụng thường xuyên làm mặc định.

<!-- wave-example: book-function-default -->
```wave
fun discounted(price: i32, percent: i32 = 10) -> i32 {
    return price - price * percent / 100;
}

fun main() {
    println("{}", discounted(2000));
    println("{}", discounted(2000, 25));
}
```

Kết quả thực hiện:

```text
1800
1500
```

Cuộc gọi đầu tiên bỏ qua đối số thứ hai và sử dụng 10. Cuộc gọi thứ hai sử dụng 25 được chỉ định. Các giá trị mặc định được để lại cho các tham số tùy chọn ở cuối. Nó không được sử dụng như một ngữ pháp để viết một khoảng trống nhằm chỉ bỏ qua đối số đầu tiên.

Việc thay đổi mặc định sẽ thay đổi hành vi bỏ qua cuộc gọi. Các giá trị mặc định của các hàm public cũng là một phần của hành vi mà bạn dựa vào. Đây là lý do tại sao các cuộc gọi có đối số được chỉ định và các cuộc gọi có đối số bị bỏ qua được kiểm tra riêng biệt.

## có nghĩa là truyền theo giá trị

Việc truyền một giá trị số nguyên sẽ tách biệt giá trị mà hàm nhận được từ bộ lưu trữ biến của hàm gọi. Việc tính toán kết quả không tự động thay đổi các biến của người gọi.

<!-- wave-example: book-function-value -->
```wave
fun next(value: i32) -> i32 {
    return value + 1;
}

fun main() {
    var count: i32 = 4;
    var later: i32 = next(count);
    println("count={} later={}", count, later);
    count = next(count);
    println("count={}", count);
}
```

Kết quả thực hiện:

```text
count=4 later=5
count=5
```

Cuộc gọi đầu tiên đọc count và khởi tạo later. count vẫn là 4. Sau lệnh gọi thứ hai, kết quả được gán cho count, do đó nó trở thành 5. Thiết kế trả về theo giá trị giúp người gọi biết rõ dữ liệu đã thay đổi ở đâu.

Nếu bạn muốn thay đổi không gian lưu trữ ban đầu trong một hàm, bạn có thể truyền con trỏ. Điều này được đề cập trong [chương con trỏ](/docs/vi/language/explicit-memory-type-model). Ngay cả khi bạn truyền một con trỏ, bạn phải phân biệt giữa chính giá trị con trỏ và không gian lưu trữ cho địa chỉ của nó.

## Đừng rời bỏ những con đường không trở lại

Hàm trả về một giá trị phải cung cấp kết quả trên tất cả các đường dẫn được yêu cầu. Nó không để lại con đường cuối cùng như thế này:

```wave
// 의도적으로 잘못된 함수 조각

fun positive(value: i32) -> i32 {
    if (value > 0) {
        return value;
    }
}
```

Nếu value nhỏ hơn hoặc bằng 0 thì không có giá trị nào được đặt để trả về. Khi bạn đã thiết lập các quy tắc dự định của mình, bạn cần viết ra tất cả các tuyến đường của mình.

<!-- wave-example: book-function-paths -->
```wave
fun positive_or_zero(value: i32) -> i32 {
    if (value > 0) {
        return value;
    }

    return 0;
}

fun main() {
    println("{} {} {}", positive_or_zero(-2), positive_or_zero(0), positive_or_zero(8));
}
```

Kết quả thực hiện:

```text
0 0 8
```

Cuộc gọi thực hiện đầu tiên return không tiếp tục tới return bên dưới. return cuối cùng chỉ đạt được khi điều kiện sai. Chúng tôi đã kiểm tra quy tắc cho ba trường hợp: dương, 0 và âm.

## Chức năng không có kết quả

Nếu bạn chỉ thực hiện các thao tác như in, bạn có thể bỏ qua kiểu trả về. Một chức năng không có kết quả cũng có thể bị chấm dứt sớm bằng `return;`.

<!-- wave-example: book-function-void -->
```wave
fun show_positive(value: i32) {
    if (value <= 0) {
        return;
    }

    println("positive={}", value);
}

fun main() {
    show_positive(-3);
    show_positive(6);
}
```

Kết quả thực hiện:

```text
positive=6
```

Cuộc gọi đầu tiên trở lại mà không in bất cứ điều gì. Cuộc gọi thứ hai in ra: “Không có kết quả” và “Không quay lại điểm gọi” là khác nhau. Các hàm không trả về, chẳng hạn như các hàm kết thúc một quá trình, được phân loại theo kiểu trả về của chúng `!`.

## Soạn thảo một chương trình hoàn chỉnh với nhiều chức năng

Bây giờ chúng tôi chia xác thực đầu vào, tính toán và đầu ra thành các hàm khác nhau. Vì chúng tôi đã giới hạn phạm vi giá nên phép nhân trung gian trong ví dụ này nằm trong phạm vi i32.

<!-- wave-example: book-function-program -->
```wave
fun valid_order(price: i32, quantity: i32, percent: i32) -> bool {
    return price >= 0 && price <= 100000
        && quantity >= 1 && quantity <= 100
        && percent >= 0 && percent <= 100;
}

fun discounted_unit(price: i32, percent: i32) -> i32 {
    return price - price * percent / 100;
}

fun order_total(price: i32, quantity: i32, percent: i32) -> i32 {
    var unit: i32 = discounted_unit(price, percent);
    return unit * quantity;
}

fun show_order(price: i32, quantity: i32, percent: i32) {
    if (!valid_order(price, quantity, percent)) {
        println("invalid order");
        return;
    }

    var total: i32 = order_total(price, quantity, percent);
    println("total={}", total);
}

fun main() {
    show_order(2000, 3, 10);
    show_order(2000, 0, 10);
    show_order(2000, 3, 120);
}
```

Kết quả thực hiện:

```text
total=5400
invalid order
invalid order
```

Trả lời một câu hỏi cho mỗi chức năng. valid_order xác định xem có cho phép đầu vào hay không, discounted_unit xác định mức giảm giá của một lần, order_total xác định tổng giá là bao nhiêu và show_order xác định nội dung sẽ hiển thị.

Các chức năng nhỏ hơn không nhất thiết phải tốt hơn. Việc đặt tên cho mỗi biểu thức có thể gây khó khăn cho việc theo dõi. Tách riêng các quy tắc khi chúng có ý nghĩa để sử dụng lại ở nơi khác hoặc khi có các quy tắc cần được giải thích và xác minh một cách độc lập.

## vấn đề thực hành

1. Viết maximum, trả về số lớn hơn trong hai số nguyên.
2. Viết hàm clamp trả về giá trị giữa hai ranh giới. Trong giải pháp bên dưới, low <= high được đặt làm điều kiện gọi.
3. Tạo hàm cộng thuế vào một giá trị và kết hợp nó với hàm tổng thứ tự. Trước tiên, hãy xác định phạm vi và điểm cắt số nguyên.

### Giải: Hàm có ranh giới

<!-- wave-example: book-function-solution -->
```wave
fun maximum(left: i32, right: i32) -> i32 {
    if (left > right) {
        return left;
    }

    return right;
}

fun clamp(value: i32, low: i32, high: i32) -> i32 {
    if (value < low) {
        return low;
    }
    if (value > high) {
        return high;
    }

    return value;
}

fun main() {
    println("max={}", maximum(7, 4));
    println("{} {} {}", clamp(-3, 0, 10), clamp(6, 0, 10), clamp(20, 0, 10));
}
```

Kết quả thực hiện:

```text
max=7
0 6 10
```

clamp có ba đường dẫn: nhỏ hơn phạm vi, trong phạm vi và lớn hơn phạm vi. Kiểm tra bằng cách thêm thủ công các giá trị biên 0 và 10. Nếu bạn muốn đăng ký lên tới low > high, bạn phải quyết định cách thể hiện sự thất bại. [Chương xử lý lỗi](/docs/vi/language/errors) giải quyết vấn đề này bằng cách sử dụng cấu trúc kết quả.


## cuộc gọi đệ quy

Một hàm có thể gọi chính nó. Đệ quy yêu cầu một điều kiện thoát trong đó không còn cuộc gọi nào được thực hiện nữa và một quy trình trong đó mỗi cuộc gọi tiến gần hơn đến điều kiện đó.

<!-- wave-example: book-ref-recursion -->
```wave
fun factorial(value: i32) -> i32 {
    if (value <= 1) {
        return 1;
    }

    return value * factorial(value - 1);
}

fun main() {
    println("{}", factorial(5));
}
```

Kết quả thực hiện:

```text
120
```

Phép tính 5 dẫn đến `5 * factorial(4)`, trả về 1 khi đạt đến. Giá trị trả về được chuyển tuần tự cho lệnh gọi trước đó, dẫn đến 120. Hàm này là một ví dụ để giải thích các số nguyên dương nhỏ. Với đầu vào lớn, bạn cần xem xét phạm vi kết quả và độ sâu cuộc gọi. Bạn có thể tránh vấn đề tăng độ sâu cuộc gọi bằng cách viết thao tác tương tự vào một vòng lặp.

## Các lỗi thường gặp

|hiện tượng|Kiểm tra|
| --- | --- |
|Lỗi không đủ hoặc có quá nhiều đối số|Số lượng tham số và giá trị mặc định tùy chọn|
|Kiểu trả về không khớp|return Khai báo kiểu biểu thức và hàm|
|Không trở về từ một con đường cụ thể|Nó có trả về ngay cả khi điều kiện sai không?|
|Thiếu đối số loại chung|Sau tên hàm `<Type>`|
|Bản gốc được thay đổi sau khi hàm con trỏ được gọi.|Đây là hàm chỉ đọc giá trị hay hàm sửa đổi nó?|

Các hàm được xuất bên ngoài, chẳng hạn như `export(c)`, sử dụng chữ ký cụ thể. Bản thân các hàm chung không thể được xuất bằng quy ước gọi bên ngoài. `ptr<T>` và `array<T, N>` là các loại bộ nhớ tích hợp của ngôn ngữ và phân biệt chúng với các khai báo cấu trúc chung của người dùng.

[học chức năng](/docs/vi/language/functions-and-generics) · [Mô-đun học tập và khái quát](/docs/vi/language/modules-imports-and-ffi) · [FFI](/docs/vi/language/modules-imports-and-ffi)
