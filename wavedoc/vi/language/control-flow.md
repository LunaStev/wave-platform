---
translation_set_id: control-flow
path: language/control-flow
locale: vi
group: language
group_order: 2
order: 4
title: 4. Điều kiện, vòng lặp và giá trị biên
summary: Tìm hiểu if, for, while và phạm vi các câu lặp đi lặp lại.
---

## Chọn đường dẫn thực hiện

Chương trình ở chương trước thực hiện các câu lệnh từ trên xuống dưới. Một chương trình thực sự phải làm những việc khác nhau tùy thuộc vào đầu vào và trạng thái. Câu lệnh có điều kiện chọn đường dẫn thực thi và vòng lặp áp dụng cùng một quy tắc cho nhiều giá trị.

Lưu từng ví dụ vào main.wave và chạy nó. Khi đọc mã, hãy viết ra một tờ giấy các giá trị biến hiện tại, các điều kiện cần kiểm tra tiếp theo và thứ tự các câu lệnh sẽ được thực hiện. Điều quan trọng hơn là thực hành theo dòng chảy hơn là ghi nhớ kết quả.

## if và else

Các điều kiện được viết trong dấu ngoặc đơn và văn bản được đặt trong dấu ngoặc nhọn. Trong ví dụ sau, thay thế balance bằng 500 hoặc 2000 để xác định nhánh nào được thực thi.

<!-- wave-example: book-if-else -->
```wave
fun main() {
    var balance: i32 = 2000;
    var price: i32 = 1200;

    if (balance >= price) {
        balance -= price;
        println("bought, balance={}", balance);
    } else {
        println("not enough money");
    }
}
```

Kết quả thực hiện:

```text
bought, balance=800
```

Nó không thực thi cả hai khối. Nếu điều kiện đúng thì khối đầu tiên sẽ được thực thi. Nếu điều kiện sai, khối else sẽ được thực thi. Việc mua hàng được cho phép ngay cả khi balance bằng price, vì vậy tôi đã viết `>=`. Nếu bạn thay đổi nó thành `>`, thao tác sẽ thay đổi với số lượng tương tự.

## Thứ tự của nhiều điều kiện

Bạn có thể ghép các điều kiện với else if. Vì trước tiên chúng ta thực hiện một nhánh thỏa mãn từ phía trên, nên chúng ta cần cân nhắc xem liệu chúng ta muốn kiểm tra ranh giới lớn hơn hay ranh giới nhỏ hơn trước.

<!-- wave-example: book-grade -->
```wave
fun grade(score: i32) {
    if (score < 0 || score > 100) {
        println("invalid");
    } else if (score >= 90) {
        println("A");
    } else if (score >= 80) {
        println("B");
    } else {
        println("C");
    }
}

fun main() {
    grade(95);
    grade(80);
    grade(79);
    grade(101);
}
```

Kết quả thực hiện:

```text
A
B
C
invalid
```

Điểm không hợp lệ đầu tiên sẽ bị từ chối và sau đó được chấm điểm. Nếu bạn đặt `score >= 80` ở đầu thì bạn sẽ không bao giờ đến được nhánh A vì 95 độ đi vào nhánh đó. Kiểm tra thứ tự các điều kiện cũng như tính đúng đắn của từng điều kiện.

## Không thay đổi giá trị trong biểu thức điều kiện

Các phép gán, phép gán phức hợp và các phép toán tăng hoặc giảm không được phép trong các điều kiện if, while hoặc for. Sử dụng `==` để so sánh. Để cập nhật một giá trị và sau đó kiểm tra nó, hãy viết hai câu lệnh riêng biệt.

Dạng đúng của một đoạn bên trong hàm:

```wave
value = read_value();

if (value == expected) {
    println("matched");
}
```

read_value và expected không được xác định trong đoạn này, vì vậy đây không phải là chương trình đầy đủ chạy như bình thường. Quy tắc được hiển thị ở đây là “So sánh sau khi thay đổi trạng thái”.

## while: Trong khi điều kiện được giữ nguyên

<!-- wave-example: book-while-countdown -->
```wave
fun main() {
    var remaining: i32 = 3;

    while (remaining > 0) {
        println("{}", remaining);
        remaining -= 1;
    }

    println("finished at {}", remaining);
}
```

Kết quả thực hiện:

```text
3
2
1
finished at 0
```

Các điều kiện được kiểm tra trước khi vào cơ thể. Nếu remaining là 0 ngay từ đầu thì phần thân sẽ không bao giờ được thực thi. Nếu phần giảm ở cuối phần thân bị bỏ qua, điều kiện vẫn đúng và vòng lặp không kết thúc.

Sau khi viết vòng lặp, hãy kiểm tra “Điều gì khiến vòng lặp tiến gần hơn đến điều kiện kết thúc?” Nếu đó là vòng lặp chờ đầu vào, thay đổi đầu vào hoặc EOF đóng vai trò của nó và nếu đó là vòng lặp số, cập nhật chỉ mục sẽ đóng vai trò của nó.

## for: Khởi tạo/Điều kiện/Cập nhật

for thể hiện ba phần cần thiết để lặp lại.

<!-- wave-example: book-for-sum -->
```wave
fun main() {
    var total: i32 = 0;

    for (var number: i32 = 1; number <= 5; number += 1) {
        total += number;
    }

    println("sum={}", total);
}
```

Kết quả thực hiện:

```text
sum=15
```

1. Khởi tạo number thành 1. Bước này chỉ thực hiện một lần.
2. Kiểm tra number <= 5. Nếu sai, kết thúc vòng lặp.
3. Thêm number vào total trong văn bản.
4. Tăng number lên 1 và quay lại kiểm tra điều kiện.

Người ta không cho rằng các biến lặp lại được khai báo trong for có thể được sử dụng sau khi lặp lại. Nếu thiết kế của bạn yêu cầu một giá trị sau khi lặp, hãy khai báo giá trị đó bên ngoài và làm rõ vị trí khởi tạo.

## Ranh giới bao gồm và loại trừ

Tổng tự nhiên của 1 đến 5 là `<= 5`. Mặt khác, chỉ mục của mảng có độ dài 5 nên sử dụng `< 5`. Điều này là do các chỉ số mảng bắt đầu từ 0 và kết thúc ở 4.

Đừng nhầm lẫn giữa “chạy năm lần” với “tối đa và bao gồm giá trị 5”. Bạn có thể xác định số lần lặp lại bằng cách xem xét các giá trị bắt đầu và kết thúc cùng nhau. Các trường hợp đầu vào trống và chỉ có một phần tử sẽ tốt cho việc phát hiện các lỗi biên.

## continue và break

continue bỏ qua phần còn lại của lần lặp này và break kết thúc lần lặp gần nhất.

<!-- wave-example: book-loop-control -->
```wave
fun main() {
    var total: i32 = 0;

    for (var number: i32 = 1; number <= 10; number += 1) {
        if (number == 3) {
            continue;
        }

        if (number == 6) {
            break;
        }

        total += number;
    }

    println("{}", total);
}
```

Kết quả thực hiện:

```text
12
```

Các số trong tổng số là 1, 2, 4 và 5. Nếu gặp continue trong for, quá trình gia hạn sẽ được tiến hành. Vì while không có biểu thức cập nhật riêng như for, bạn phải cẩn thận để không bỏ sót bất kỳ thay đổi trạng thái nào được yêu cầu trước continue.

Nếu các vòng lặp được lồng nhau thì một break sẽ không hoàn thành tất cả các vòng lặp. Nếu bạn cần dừng ở nhiều giai đoạn, hãy đưa công việc vào một hàm và cho biết ý định của bạn bằng cách sử dụng return hoặc bằng cách kiểm tra cả điều kiện kết thúc trong lần lặp bên ngoài.

## Chia trường hợp cho match

Bạn có thể sử dụng match khi so sánh nhiều trường hợp có cùng giá trị. Phần thân của mỗi arm là một khối.

<!-- wave-example: book-match-number -->
```wave
fun describe(status: i32) {
    match (status) {
        200 => {
            println("ok");
        }
        404 => {
            println("missing");
        }
        _ => {
            println("other");
        }
    }
}

fun main() {
    describe(200);
    describe(404);
    describe(500);
}
```

Kết quả thực hiện:

```text
ok
missing
other
```

`_` là mẫu xử lý phần còn lại. Không đặt các bản sao trong cùng một match. variant, có các loại dữ liệu khác nhau tùy thuộc vào giá trị, được đề cập trong [Chương mô hình dữ liệu](/docs/vi/language/structures-enums-and-aliases).

## Ví dụ đầy đủ: Đếm số thỏa mãn điều kiện

Tìm số và tổng của các số chẵn từ 1 đến 10. Vì số và tổng là những thông tin khác nhau nên chúng được tích lũy dưới dạng biến.

<!-- wave-example: book-loop-statistics -->
```wave
fun main() {
    var count: i32 = 0;
    var total: i32 = 0;

    for (var number: i32 = 1; number <= 10; number += 1) {
        if (number % 2 == 0) {
            count += 1;
            total += number;
        }
    }

    println("count={} total={}", count, total);
}
```

Kết quả thực hiện:

```text
count=5 total=30
```

Các số chẵn là 2, 4, 6, 8 và 10, do đó số đó là 5 và tổng là 30. Ngay cả khi biểu thức ngắn, bạn vẫn có thể dễ dàng xác minh ranh giới lặp nếu trước tiên bạn kiểm tra kết quả trong một phạm vi nhỏ có thể lấy được bằng tay.

## Bài tập và lời giải đầy đủ

Chỉ cộng bội số của 3 từ 1 đến 20, nhưng không cộng bất kỳ giá trị nào có tổng lớn hơn 30. Chúng ta cần phân biệt giữa “cộng rồi kiểm tra xem đã hết chưa” và “kiểm tra xem đã hết chưa rồi cộng”.

<!-- wave-example: book-loop-solution -->
```wave
fun main() {
    var total: i32 = 0;

    for (var number: i32 = 1; number <= 20; number += 1) {
        if (number % 3 != 0) {
            continue;
        }

        if (total + number > 30) {
            break;
        }

        total += number;
    }

    println("{}", total);
}
```

Kết quả thực hiện:

```text
30
```

Các giá trị 3+6+9+12 có tổng bằng 30 nên giá trị tiếp theo, 15, không được thêm vào. Những đầu vào nhỏ này an toàn, nhưng đối với số nguyên lớn, kiểm tra `total + number` có thể tự tràn. Viết séc không tự động xử lý mọi trường hợp ranh giới.
