---
translation_set_id: expressions
path: language/expressions-and-operators
locale: vi
group: language
group_order: 2
order: 3
title: 3. Số học, so sánh và chuyển đổi
summary: Tìm hiểu thứ tự tính toán, chia số nguyên, các phép toán theo bit và cast.
---

## Xem kết quả tính toán và các loại phép tính cùng nhau

Biểu thức là mã tính toán một giá trị. Tên biến, chữ, lệnh gọi hàm và biểu thức kết nối nhiều giá trị với toán tử đều là biểu thức. Trong toán học, ngay cả khi nó trông giống như một phương trình, kết quả sẽ khác nhau tùy thuộc vào việc đó là số nguyên hay số thực và số bit chứa trong đó.

Bắt đầu với các phép tính đơn giản, chương này giới thiệu các dấu ngoặc đơn, phép chia, các phép toán logic, các phép toán theo bit và ép kiểu. Mỗi ví dụ là một tệp main.wave hoàn chỉnh mà bạn có thể chạy với `wavec run main.wave`.

## Phạm vi được đặt trong dấu ngoặc đơn

<!-- wave-example: book-precedence -->
```wave
fun main() {
    var first: i32 = 2 + 3 * 4;
    var second: i32 = (2 + 3) * 4;

    println("{} {}", first, second);
}
```

Kết quả thực hiện:

```text
14 20
```

Phép nhân được thực hiện trước phép cộng, vì vậy biểu thức đầu tiên là 2+12. Trong phương trình thứ hai, chúng ta lấy tổng trong ngoặc đơn là 5 rồi nhân với 4. Mục đích không phải là sử dụng ít dấu ngoặc đơn hơn. Nên sử dụng nó để người đọc có thể dễ dàng hiểu được phạm vi tính toán.

Ngay cả khi sử dụng cùng một toán tử nhiều lần thì hướng của chuỗi vẫn quan trọng. `20 - 5 - 3` là `(20 - 5) - 3`, tức là 12. `20 - (5 - 3)` là 18. Chuỗi đầy đủ chính xác nằm trong [tham chiếu toán tử](/docs/vi/language/expressions-and-operators).

## Phép chia số nguyên và số dư

Chia hai số nguyên không tạo ra kết quả có dấu phẩy động phân số. Sử dụng phép chia cho thương và toán tử số dư cho số dư.

<!-- wave-example: book-division -->
```wave
fun main() {
    var items: i32 = 17;
    var box_size: i32 = 5;
    var full_boxes: i32 = items / box_size;
    var remaining: i32 = items % box_size;

    println("boxes={} remaining={}", full_boxes, remaining);
    println("negative quotient={}", -17 / 5);
}
```

Kết quả thực hiện:

```text
boxes=3 remaining=2
negative quotient=-3
```

Chia 17 món thành các nhóm 5 món được 3 nhóm hoàn chỉnh và 2 món còn lại. Phép chia có dấu rút ngắn về 0, nên -17/5 là -3. Điều này khác với việc làm tròn xuống tới âm vô cực.

Bạn không thể chia cho 0. signed Giá trị thu được khi chia giá trị nhỏ nhất cho -1 không cùng loại. Các hàm nhận những đầu vào này phải kiểm tra trước khi chia hoặc sử dụng toán checked toán API.

## Điểm chuyển đổi làm thay đổi kết quả.

Hai biểu thức sau đây đều được lưu trữ trong biến f64, nhưng quá trình tính toán thì khác.

<!-- wave-example: book-division-conversion -->
```wave
fun main() {
    var already_divided: f64 = (7 / 2) as f64;
    var floating_division: f64 = (7 as f64) / (2 as f64);

    if (already_divided == 3.0) {
        println("integer division lost the fraction");
    }

    if (floating_division == 3.5) {
        println("floating division kept the fraction");
    }
}
```

Kết quả thực hiện:

```text
integer division lost the fraction
floating division kept the fraction
```

Biểu thức đầu tiên thực hiện phép chia số nguyên để thu được 3, sau đó chuyển đổi nó thành f64. Lệnh thứ hai chuyển đổi toán hạng thành f64 trước khi thực hiện phép chia dấu phẩy động. Việc chọn loại rộng hơn cho biến cuối cùng không thể khôi phục thông tin bị mất trước đó.

Giá trị dấu phẩy động là giá trị gần đúng. Hai kết quả trông giống nhau có thể không phù hợp để so sánh bằng nhau một cách chính xác. Chọn dung sai phù hợp với đơn vị và quy mô của bài toán; một epsilon cố định không phù hợp cho mọi phép tính.

## Đã thực hiện so sánh bool

`<`, `<=`, `>`, `>=`, `==`, `!=` kiểm tra mối quan hệ. Một dấu bằng `=` là phép gán và hai dấu bằng `==` là phép so sánh bằng nhau.

<!-- wave-example: book-comparisons -->
```wave
fun main() {
    var age: i32 = 20;
    var minimum: i32 = 18;
    var eligible: bool = age >= minimum;

    if (eligible) {
        println("eligible");
    }

    if (age != minimum) {
        println("not exactly the boundary");
    }
}
```

Kết quả thực hiện:

```text
eligible
not exactly the boundary
```

“Ít nhất 18” bao gồm 18; "lớn hơn 18" loại trừ nó. Kiểm tra 17, 18 và 19 để kiểm tra ranh giới này. So sánh kiểu hỗn hợp phụ thuộc vào độ ký và chiều rộng, do đó việc chuyển đổi cả hai toán hạng sang loại dự định có thể làm cho việc so sánh rõ ràng hơn.

## Hoạt động logic và đánh giá ngắn mạch

`&&` kiểm tra xem cả hai đều đúng hay không, `||` kiểm tra xem có nhiều hơn một điều đúng hay không và `!` lật đúng/sai. Thao tác này không phải lúc nào cũng thực hiện biểu thức ở phía bên phải.

<!-- wave-example: book-short-circuit -->
```wave
fun report() -> bool {
    println("right side evaluated");
    return true;
}

fun main() {
    var enabled: bool = false;

    if (enabled && report()) {
        println("both true");
    }

    if (!enabled || report()) {
        println("at least one true");
    }
}
```

Kết quả thực hiện:

```text
at least one true
```

Ở điều kiện đầu tiên, enabled là sai nên không cần phải nhìn sang bên phải. Trong trường hợp thứ hai, !enabled là đúng nên vế phải cũng không cần thiết. Do đó, đầu ra của report không bao giờ xuất hiện.

Bạn có thể sử dụng điều này để kiểm tra mẫu số trước khi chia. Đoạn nội dung hàm `if (divisor != 0 && value / divisor > 2) { ... }` không thực hiện phép chia khi mẫu số bằng 0. Tuy nhiên, riêng thử nghiệm này không giải quyết được các ranh giới khác như giá trị tối thiểu signed/-1.

`&&` được ưu tiên hơn `||`. Trong các chính sách phức tạp, hãy cho biết ý định của bạn trong dấu ngoặc đơn, chẳng hạn như `(member && active) || admin`.

## Thu hẹp hoặc mở rộng số nguyên

`as` là một diễn viên rõ ràng. Các bit bậc cao bị loại bỏ khi thu hẹp một số nguyên không thể được phục hồi bằng cách mở rộng lại.

<!-- wave-example: book-cast-chain -->
```wave
fun main() {
    var original: i32 = 300;
    var small: u8 = original as u8;
    var widened: i32 = small as i32;

    println("{} -> {} -> {}", original, small, widened);

    var negative: i8 = -1;
    var signed_value: i32 = negative as i32;
    var same_bits: u8 = negative as u8;

    println("{} {}", signed_value, same_bits);
}
```

Kết quả thực hiện:

```text
300 -> 44 -> 44
-1 255
```

8 bit thấp hơn của 300 là 44. Nếu bạn mở rộng -1 thành loại signed, dấu sẽ được mở rộng để duy trì -1. Nếu bạn diễn giải 8 bit giống như unsigned thì nó là 255.

Các phép biến đổi tìm cách bảo toàn các giá trị trong phạm vi và các phép biến đổi cố gắng thao tác các bit lưu trữ có các mục đích khác nhau. Nếu bạn có đầu vào của người dùng, trước tiên hãy kiểm tra xem đó có phải là phạm vi đích hay không và chuyển đổi nó. Sự hiện diện của cast không đảm bảo rằng giá trị đó nằm trong phạm vi an toàn.

## Thay thế bằng bool

Chuyển đổi một số nguyên thành bool mang lại false cho số 0 và true nếu ngược lại. Điều này không cắt bớt bit thấp nhất: 2 cũng chuyển đổi thành true. Đối với các giá trị dấu phẩy động, chỉ +0,0 và -0,0 chuyển đổi thành false; tất cả các giá trị khác, bao gồm NaN và vô cực, chuyển đổi thành true.

<!-- wave-example: book-bool-conversion -->
```wave
fun main() {
    var zero: i32 = 0;
    var two: i32 = 2;

    if (!(zero as bool)) {
        println("zero is false");
    }

    if (two as bool) {
        println("two is true");
    }
}
```

Kết quả thực hiện:

```text
zero is false
two is true
```

Chuyển đổi con trỏ tới bool không được hỗ trợ. So sánh con trỏ với null một cách rõ ràng, ví dụ: với `pointer != null`. Liệu địa chỉ không phải null có an toàn để đọc hay không là một câu hỏi riêng.

## hoạt động bit

`&`, `|`, `^`, `~` bao gồm từng bit của số nguyên. Nó có thể được sử dụng để thể hiện các quyền hoặc chức năng theo bit. Bên dưới, 1 là quyền đọc và 2 là quyền ghi.

<!-- wave-example: book-bit-flags -->
```wave
const READ: u32 = 1;
const WRITE: u32 = 2;

fun main() {
    var permissions: u32 = READ | WRITE;

    if ((permissions & WRITE) != 0) {
        println("write enabled");
    }

    permissions = permissions & ~WRITE;
    println("remaining={}", permissions);
}
```

Kết quả thực hiện:

```text
write enabled
remaining=1
```

Thêm một bit bằng OR và kiểm tra sự hiện diện của một bit cụ thể bằng AND. Tạo mặt nạ chỉ có các bit liên quan là 0 với `~WRITE` và xóa nó. Các phép toán theo bit `&`·`|` là các toán tử khác với đánh giá ngắn mạch `&&`·`||` của bool.

## Chiều rộng và số ca

Phép dịch trái di chuyển các bit sang trái và loại bỏ các bit cao nằm ngoài độ rộng của toán hạng. Dịch chuyển phải sử dụng phần mở rộng dấu cho các giá trị có dấu và phần mở rộng bằng 0 cho các giá trị không dấu. Kết quả luôn có kiểu toán hạng bên trái.

<!-- wave-example: book-shifts -->
```wave
fun main() {
    var bits: u8 = 129;
    var shifted: u8 = bits << 1;
    var negative: i32 = -8;
    var positive: u32 = 8;

    println("left={}", shifted);
    println("right={} {}", negative >> 1, positive >> 1);
}
```

Kết quả thực hiện:

```text
left=2
right=-4 4
```

Việc dịch chuyển giá trị u8 129 còn lại sẽ loại bỏ bit cao nhất của nó và để lại 2. Giá trị số lần thay đổi ban đầu phải không âm và nhỏ hơn độ rộng bit của toán hạng bên trái. Đối với u8, số lượng hợp lệ là từ 0 đến 7. Số lượng hằng số không hợp lệ là lỗi thời gian biên dịch; số lượng thời gian chạy không hợp lệ sẽ gây ra bẫy.

## Chuyển đổi giá trị dấu phẩy động thành số nguyên

Chuyển đổi dấu phẩy động thành số nguyên trước tiên sẽ bỏ phần thập phân theo hướng về 0, sau đó kiểm tra phạm vi số nguyên đích. Các kết quả NaN, vô cực và ngoài phạm vi đều không hợp lệ. Phép chuyển đổi hằng không hợp lệ tạo ra lỗi thời gian biên dịch; chuyển đổi thời gian chạy không hợp lệ gây ra bẫy.

Bẫy không trả về giá trị lỗi từ hàm. Đối với các lỗi chuyển đổi có thể phục hồi, hãy thiết kế giao diện kiểm tra phạm vi trước khi chuyển đổi. Để thu được các bit được lưu trữ của một giá trị dấu phẩy động, hãy sử dụng các hàm chuyển đổi bit trong `std::math::float` thay vì truyền số.

## Bài tập và lời giải

Viết hàm chia 137 won đổi thành 50 won đơn vị và số dư, đồng thời chỉ thay đổi các số nguyên trong phạm vi từ 0 đến 255 đến u8. Ví dụ này minh họa bước xác thực dưới dạng một hàm riêng biệt cho biết nằm ngoài phạm vi là -1.

<!-- wave-example: book-expression-solution -->
```wave
fun checked_byte(value: i32) -> i32 {
    if (value < 0 || value > 255) {
        return -1;
    }

    var small: u8 = value as u8;
    return small as i32;
}

fun main() {
    var money: i32 = 137;

    println("coins={} remainder={}", money / 50, money % 50);
    println("{} {} {}", checked_byte(0), checked_byte(255), checked_byte(256));
}
```

Kết quả thực hiện:

```text
coins=2 remainder=37
0 255 -1
```

Lý do tại sao -1 có thể được sử dụng làm chỉ báo lỗi là vì phạm vi thành công là từ 0 đến 255. Nếu bất kỳ số nguyên nào có thể là giá trị thành công thì cần có một biểu diễn kết quả khác. Chúng ta tiếp tục thiết kế này trong chương xử lý lỗi sau.


## ưu tiên

Độ ưu tiên của toán tử như sau, bắt đầu từ cao nhất:

1. Các biểu thức cơ bản và truy cập hậu tố: gọi hàm, truy cập trường, lập chỉ mục, hậu tố `++`·`--`
2. Các phép toán một ngôi: `!`, `~`, `&`, `deref`, tiền tố `++`·`--`, một ngôi `+`·`-`
3. `as` chuyển đổi loại
4. `*`, `/`, `%`
5. `+`, `-`
6. `<<`, `>>`
7. `<`, `<=`, `>`, `>=`
8. `==`, `!=`
9. Bit `&`
10. Bit `^`
11. Bit `|`
12. `&&`
13. `||`
14. Bài tập và bài tập ghép

Phân công chuỗi kết hợp từ bên phải. Khi trộn các loại toán tử khác nhau, hãy sử dụng dấu ngoặc đơn để thể hiện rõ thứ tự đánh giá.

## Chủ đề có thể thay thế

Việc gán, `++` và `--` yêu cầu biểu thức biểu thị vị trí lưu trữ, chẳng hạn như biến, trường, phần tử mảng hoặc con trỏ hủy tham chiếu. Không được phép ghi vào `const`.

## sự thay đổi

Kết quả của phép dịch luôn có kiểu toán hạng bên trái; kiểu toán hạng bên phải không mở rộng tính toán. Dịch chuyển trái loại bỏ các bit cao vượt quá độ rộng đó. Dịch chuyển sang phải các giá trị có dấu mở rộng dấu và mở rộng các giá trị không dấu bằng 0.

Số ca phải là số nguyên có giá trị ban đầu thỏa mãn `0 <= n < LHS bit width`. Nó được kiểm tra trước khi cắt bớt thành loại nhỏ hơn. Số lượng hằng số không hợp lệ là lỗi thời gian biên dịch; số lượng thời gian chạy không hợp lệ sẽ gây ra bẫy.

## Chuyển đổi liên quan đến giá trị bool và dấu phẩy động

Chuyển đổi số nguyên thành bool tạo ra false cho số 0 và true cho mọi giá trị khác. Chuyển đổi dấu phẩy động sang bool chỉ tạo ra false cho +0,0 và -0,0; NaN và vô cực dương hoặc âm tạo ra true. Không hỗ trợ chuyển đổi con trỏ sang bool: so sánh rõ ràng với `pointer != null`.

Chuyển đổi dấu phẩy động thành số nguyên bỏ phần thập phân theo hướng về 0 và sau đó kiểm tra phạm vi đích. Các kết quả NaN, vô cực và ngoài phạm vi đều không hợp lệ. Phép chuyển đổi hằng không hợp lệ là lỗi thời gian biên dịch; chuyển đổi thời gian chạy không hợp lệ gây ra bẫy. Một cái bẫy không phải là một lỗi có thể phục hồi được.

`&&` và `||` là đánh giá ngắn mạch. Các tác dụng phụ của toán hạng bên phải chưa được thực hiện sẽ không xảy ra. Bạn có thể kiểm tra kết quả với các giá trị nhỏ trong [lớp số học](/docs/vi/language/expressions-and-operators).

## sự thay đổi không đúng mục đích

Số lượng dịch chuyển cho giá trị 8 bit phải từ 0 đến 7. Các chương trình bên dưới phải bị từ chối trước khi thực hiện.

<!-- wave-example: reject-shift-count -->
```wave
fun main() {
    var value: u8 = 1;
    var result: u8 = value << 8;
}
```
