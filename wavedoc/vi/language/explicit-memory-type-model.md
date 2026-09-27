---
translation_set_id: memory-model
path: language/explicit-memory-type-model
locale: vi
group: language
group_order: 2
order: 9
title: 9. Con trỏ, sửa đổi giá trị và thời gian tồn tại
summary: Tìm hiểu địa chỉ, hội thảo, sửa đổi giá trị thông qua con trỏ và con trỏ lơ lửng.
---

## Phân biệt giá trị và vị trí lưu trữ

Số nguyên 42 và địa chỉ nơi số nguyên được lưu trữ là các giá trị khác nhau. Con trỏ trỏ đến vị trí lưu trữ. Việc truyền địa chỉ cho phép hàm này đọc hoặc thay đổi bộ nhớ của người gọi.

Chương này đề cập đến việc lấy địa chỉ, hội thảo, sửa đổi giá trị ban đầu, số học con trỏ và thời gian tồn tại. Phân bổ động sẽ được đề cập trong chương tiếp theo. Bắt đầu với địa chỉ của các biến cục bộ và các phần tử mảng.

## Lấy địa chỉ và hội thảo

<!-- wave-example: book-pointer-first -->
```wave
fun main() {
    var value: i32 = 42;
    var address: ptr<i32> = &value;

    println("value={}", value);
    println("through pointer={}", deref address);

    deref address = 99;
    println("changed={}", value);
}
```

Kết quả thực hiện:

```text
value=42
through pointer=42
changed=99
```

`&value` lấy địa chỉ, `deref address` đọc hoặc ghi giá trị của địa chỉ đó, v.v. Thay vì lưu trữ 99 trong address, 99 được ghi vào số nguyên được trỏ bởi address. Bản thân biến address tiếp tục trỏ đến value.

`ptr<i32>` là loại con trỏ để truy cập bộ lưu trữ i32. Loại này không ghi lại độ dài hoặc cung cấp khả năng phân bổ tự động.

## Tự thay đổi con trỏ

<!-- wave-example: book-pointer-reassign -->
```wave
fun main() {
    var first: i32 = 10;
    var second: i32 = 20;
    var selected: ptr<i32> = &first;

    selected = &second;
    deref selected = 25;

    println("{} {}", first, second);
}
```

Kết quả thực hiện:

```text
10 25
```

`selected = &second` lưu trữ một địa chỉ khác trong một biến con trỏ. Giá trị của first không thay đổi. Sau đó, nếu bạn viết giá trị là deref, thì second sẽ thay đổi. Việc tách “thay đổi địa chỉ” và “thay đổi giá trị thông qua địa chỉ” thành các câu riêng biệt sẽ giảm nhầm lẫn.

## Tạo một hàm thay đổi bản gốc

<!-- wave-example: book-pointer-increment -->
```wave
fun increment(value: ptr<i32>) {
    deref value = deref value + 1;
}

fun main() {
    var count: i32 = 4;

    increment(&count);
    increment(&count);

    println("{}", count);
}
```

Kết quả thực hiện:

```text
6
```

Hàm nhận địa chỉ của số đếm, thay vì giá trị của nó, 4. Việc thay đổi bộ nhớ đó cũng thay đổi số lượng trong lệnh gọi. Hàm này yêu cầu địa chỉ i32 hợp lệ, có thể ghi. Việc vượt qua null vi phạm yêu cầu này.

Mỗi hàm xác định xem nó có chấp nhận null hay không. Nếu không, người gọi phải cung cấp một địa chỉ hợp lệ. Nếu đúng như vậy, hàm phải bao gồm một đường dẫn xử lý null.

## Chức năng xử lý null

<!-- wave-example: book-pointer-null -->
```wave
fun try_increment(value: ptr<i32>) -> bool {
    if (value == null) {
        return false;
    }

    deref value = deref value + 1;
    return true;
}

fun main() {
    var count: i32 = 7;

    if (!try_increment(null)) {
        println("no value");
    }

    if (try_increment(&count)) {
        println("count={}", count);
    }
}
```

Kết quả thực hiện:

```text
no value
count=8
```

Kiểm tra null chỉ xử lý việc thiếu địa chỉ. Việc chuyển đổi một số không phải null tùy ý thành một con trỏ sẽ không tạo ra bộ nhớ hợp lệ. Việc đọc và ghi cũng yêu cầu thời gian tồn tại hợp lệ, kích thước vừa đủ, căn chỉnh chính xác và quyền truy cập phù hợp.

## Địa chỉ mảng và số học con trỏ theo phần tử

<!-- wave-example: book-pointer-array -->
```wave
fun main() {
    var values: array<i32, 3> = [10, 20, 30];
    var first: ptr<i32> = &values[0];
    var second: ptr<i32> = first + 1;

    println("{}", deref first);
    println("{}", deref second);
    println("{}", deref first[2]);
}
```

Kết quả thực hiện:

```text
10
20
30
```

Thêm 1 vào con trỏ sẽ di chuyển nó một phần tử của loại mục tiêu. Phần tử tiếp theo trong i32 và phần tử tiếp theo trong u8 có số byte dịch chuyển khác nhau. Nếu bạn nhân `first + 1` lần nữa với kích thước loại và thêm nó, nó sẽ di chuyển đến vị trí không mong muốn.

Việc lập chỉ mục con trỏ cũng phải được thực hiện trong phạm vi hợp lệ. first không nhớ chính độ dài mảng 3 nên khi truyền phạm vi cho hàm, nó sử dụng dạng nhận cả con trỏ và độ dài.

## Chuyển phạm vi đọc cho một hàm

<!-- wave-example: book-pointer-range -->
```wave
fun sum(values: ptr<i32>, count: i32) -> i32 {
    var total: i32 = 0;

    for (var index: i32 = 0; index < count; index += 1) {
        total += deref values[index];
    }

    return total;
}

fun main() {
    var values: array<i32, 4> = [2, 4, 6, 8];

    println("first two={}", sum(&values[0], 2));
    println("all={}", sum(&values[0], 4));
}
```

Kết quả thực hiện:

```text
first two=6
all=20
```

Đơn vị của count là số phần tử. Chức năng này yêu cầu người gọi phải có count có thể đọc được i32. Truyền một độ dài lớn hơn mảng thực tế sẽ vi phạm hợp đồng. Ngay cả khi bạn sử dụng cùng một kết hợp ptr<u8>·i64, API, bạn vẫn nên kiểm tra trong tài liệu xem độ dài được tính theo byte hay theo phần tử.

## Tuổi thọ: Địa chỉ có hiệu lực trong bao lâu?

Các biến cục bộ được sử dụng trong suốt thời gian tồn tại của cuộc gọi và khối. Nếu bạn trả về địa chỉ của một biến cục bộ bên trong một hàm để người gọi đọc sau, thì không gian lưu trữ đó có thể đã hết thời gian sử dụng.

Đây là một phần thiết kế tồi không nên triển khai:

```wave
fun invalid_address() -> ptr<i32> {
    var local: i32 = 42;

    return &local;
}
```

Nếu một giá trị được yêu cầu, nó sẽ trả về i32. Nếu nó cần ghi vào vùng lưu trữ do người gọi cung cấp, thì nó sẽ lấy một con trỏ làm đầu vào. Nếu bạn cần bộ nhớ riêng để giữ lại ngoài cuộc gọi, hãy phân bổ nó một cách rõ ràng và chuyển giao trách nhiệm giải phóng nó.

## Hai con trỏ trỏ đến cùng một không gian lưu trữ

Sao chép một con trỏ sẽ tạo ra một tên khác trỏ đến cùng một địa chỉ. Không trùng lặp bộ nhớ.

<!-- wave-example: book-pointer-alias -->
```wave
fun main() {
    var value: i32 = 1;
    var first: ptr<i32> = &value;
    var second: ptr<i32> = first;

    deref second = 9;

    println("{} {}", value, deref first);
}
```

Kết quả thực hiện:

```text
9 9
```

Bạn cũng có thể xem kết quả được thay đổi qua second qua first. Khi bộ nhớ nguồn được giải phóng, cả hai con trỏ đều không thể sử dụng được. Việc gán null cho một biến con trỏ không tự động thay đổi các bản sao khác.

## Bài tập: Hoán đổi hai số nguyên

Viết hàm nhận hai địa chỉ i32 và trao đổi giá trị của chúng. Giá trị đầu tiên phải được lưu trữ trong một biến tạm thời trước khi ghi đè lên nó. Kiểm tra xem giá trị có được giữ lại ngay cả khi bạn chuyển cùng một địa chỉ hai lần hay không.

### Lời giải đầy đủ

<!-- wave-example: book-pointer-swap -->
```wave
fun swap(left: ptr<i32>, right: ptr<i32>) {
    var saved: i32 = deref left;

    deref left = deref right;
    deref right = saved;
}

fun main() {
    var first: i32 = 3;
    var second: i32 = 8;

    swap(&first, &second);
    println("{} {}", first, second);

    swap(&first, &first);
    println("{}", first);
}
```

Kết quả thực hiện:

```text
8 3
8
```

Hàm này cũng yêu cầu cả hai địa chỉ để trỏ đến bộ lưu trữ số nguyên hợp lệ, có thể ghi được. Để xử lý null, hãy thêm kết quả biểu thị thành công hay thất bại, như trong try_increment.


## **Wave Explicit Memory Type Model**

Thiết kế con trỏ của Wave dựa trên **Wave Explicit Memory Type Model**. Mô hình này định nghĩa các con trỏ và mảng là các kiểu bộ nhớ rõ ràng ở cấp độ ngôn ngữ, thay vì các thủ thuật cú pháp hoặc trừu tượng hóa thư viện.

`ptr<T>` là loại trỏ đến địa chỉ bộ nhớ lưu trữ giá trị `T` và `array<T, N>` là loại bộ nhớ có độ dài cố định lưu trữ các giá trị `N` của `T` liên tiếp. Do đó, cấu trúc của con trỏ và mảng được bộc lộ như trong các đối số hàm, giá trị trả về, trường cấu trúc và các kiểu khác.

## null

```wave
var buffer: ptr<u8> = null;
if (buffer == null) {
    println("no buffer");
}
```

`null` là giá trị con trỏ không trỏ đến địa chỉ bộ nhớ hợp lệ. `null` chỉ có thể được gán cho loại `ptr<T>` và không thể được sử dụng làm giá trị số nguyên, Boolean hoặc mảng.

Hàm phân bổ hoặc tra cứu có thể trả về `null` khi không có kết quả. Kiểm tra `null` trước khi hủy tham chiếu kết quả như vậy. Việc hủy tham chiếu con trỏ `null` không truy cập vào bộ nhớ hợp lệ.

## chuyển đổi con trỏ

Khi bạn cần thay đổi địa chỉ hoặc cách biểu diễn con trỏ khác, hãy sử dụng `as`.

```wave
var raw: i64 = 0;
var p: ptr<u8> = raw as ptr<u8>;
```

Chỉ sử dụng chuyển đổi giữa số nguyên và con trỏ trên các ranh giới cấp thấp, đồng thời xem xét độ rộng địa chỉ của nền tảng đích và ABI.
