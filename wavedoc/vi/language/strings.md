---
translation_set_id: learn-strings
path: language/strings
locale: vi
group: language
group_order: 2
order: 7
title: 7. Chuỗi, ký tự và byte
summary: Phân biệt giữa chuỗi và char, UTF-8 độ dài byte, NUL, dữ liệu tìm kiếm và nhị phân.
---

## các chữ cái trên màn hình và byte trong bộ nhớ

Bạn nhìn thấy các chữ cái trên màn hình nhưng byte được lưu trữ trong bộ nhớ. Đặc biệt trong trường hợp một ký tự có nhiều UTF-8 byte, chẳng hạn như trong tiếng Hàn, rất dễ mắc lỗi nếu “độ dài” và “số ký tự” được sử dụng thay thế cho nhau.

Chương này phân biệt giữa str và char, kết thúc NUL, escape, vị trí tìm kiếm và dữ liệu nhị phân. Mỗi ví dụ là một chương trình hoàn chỉnh.

## Chuỗi ký tự và đầu ra

<!-- wave-example: book-string-literals -->
```wave
fun main() {
    var greeting: str = "안녕하세요";

    println("{}", greeting);
    println("line one\nline two");
    println("quote: \"Wave\"");
}
```

Kết quả thực hiện:

```text
안녕하세요
line one
line two
quote: "Wave"
```

Các ký tự thông thường trong dấu ngoặc kép được biểu thị dưới dạng UTF-8. escape biểu thị các byte khó ghi trực tiếp từ nguồn. `\n` là byte ngắt dòng và không in hai ký tự, dấu gạch chéo ngược và n. Để in dấu gạch chéo ngược, hãy sử dụng `\\`.

## Độ dài là số byte

<!-- wave-example: book-utf8-length -->
```wave
import("std::string::len")::{
    len
};

fun main() {
    println("ASCII={}", len("Wave"));
    println("Korean={}", len("한"));
    println("mixed={}", len("Wave한"));
}
```

Kết quả thực hiện:

```text
ASCII=4
Korean=3
mixed=7
```

`len` đếm byte trước khi kết thúc NUL, không hiển thị ký tự. Ký tự `한` có ba byte trong UTF-8. Số lượng ký tự, số điểm mã Unicode và số byte thường không thể thay thế cho nhau. Độ rộng hiển thị còn phụ thuộc vào các yếu tố như phông chữ và ký tự kết hợp.

Do đó, các hàm cắt một chuỗi ở vị trí byte tùy ý và hiển thị nó trên màn hình phải xem xét riêng ranh giới Unicode. Xác định rõ ràng các điều kiện đầu vào, cho dù đó là chương trình chỉ xử lý văn bản ASCII hay văn bản chung Unicode.

## NUL Phần cuối và độ dài

str sử dụng 0 byte để biểu thị sự kết thúc. len không bao gồm byte cuối cùng đó trong độ dài. Sẽ có lỗi khi đặt NUL bên trong một chuỗi ký tự. Đây là một ví dụ về lỗi cố ý:

```wave
fun main() {
    var text: str = "left\x00right";
}
```

`\xNN` trong nguồn chỉ định một byte có chính xác hai chữ số thập lục phân. `\x41` đại diện cho byte A là 65. Bởi vì việc viết một ký tự thông thường là UTF-8 và việc chèn một byte tùy ý là khác nhau, nên không phải tất cả str có thể được viết là `\xNN` đều hợp lệ UTF-8.

## char không chứa toàn bộ ký tự Unicode

char là giá trị ký tự 8 bit không dấu. Bạn có thể sử dụng các hằng đại diện cho các giá trị trong một phạm vi byte đơn, chẳng hạn như `'A'`. `'한'` là lỗi vì nó không nằm trong phạm vi này. `"한"` là một str riêng biệt có nhiều UTF-8 byte.

Đừng cố gắng luôn đặt một chữ cái trong một char. Trước tiên, bạn phải quyết định xem đơn vị cần để xử lý văn bản là byte hay điểm mã Unicode.

## so sánh chuỗi

Để so sánh nội dung chuỗi, hãy sử dụng hàm std. Dưới đây là chương trình kiểm tra sự khác biệt về nội dung và kiểu chữ giống nhau.

<!-- wave-example: book-string-compare -->
```wave
import("std::string::cmp")::{
    eq, starts_with, ends_with
};

fun main() {
    if (eq("Wave", "Wave")) {
        println("same bytes");
    }

    if (!eq("Wave", "wave")) {
        println("case differs");
    }

    if (starts_with("report.txt", "report") && ends_with("report.txt", ".txt")) {
        println("name matches");
    }
}
```

Kết quả thực hiện:

```text
same bytes
case differs
name matches
```

Sự so sánh này so sánh các chuỗi byte. Nó không tự động thực hiện chuyển đổi trường hợp cụ thể theo ngôn ngữ hoặc chuẩn hóa Unicode. Ngay cả khi so sánh tên tệp, quy tắc bình đẳng tên tệp trong OS không giống như so sánh chuỗi đơn giản.

## Đơn vị và lỗi của kết quả tìm kiếm

<!-- wave-example: book-string-search -->
```wave
import("std::string::find")::{
    find, count
};

fun main() {
    var first: i32 = find("banana", "na");
    var missing: i32 = find("banana", "xy");
    var matches: i32 = count("aaaa", "aa");

    println("first={} missing={}", first, missing);
    println("matches={}", matches);
}
```

Kết quả thực hiện:

```text
first=2 missing=-1
matches=2
```

find trả về vị trí đầu tiên hoặc -1. Chỉ số 0 cũng thành công nên nó được kiểm tra bằng `result >= 0`. count không phải là vị trí mà là số lượng trận đấu không trùng nhau. Ở trên, aa là số 2 vì nó khớp với 0~1 và 2~3.

Thùng needle cũng là một phần của hợp đồng. find trả về 0, contains trả về true và count trả về 0. Đừng cho rằng chỉ vì tên hàm nằm trong cùng một mô-đun nên ngay cả phương thức trả về cũng giống nhau.

## Xóa dấu cách khác với việc tạo chuỗi mới

trim_range trả về phạm vi không bao gồm khoảng trắng mà không sửa đổi hoặc sao chép văn bản gốc. Vì chúng ta nhận được một con trỏ đầu ra nên trước tiên chúng ta chuẩn bị một số nguyên để lưu kết quả.

<!-- wave-example: book-string-trim-range -->
```wave
import("std::string::trim")::{
    trim_range
};

fun main() {
    var start: i32 = 0;
    var end: i32 = 0;

    trim_range("  Wave  ", &start, &end);

    println("start={} end={} bytes={}", start, end, end - start);
}
```

Kết quả thực hiện:

```text
start=2 end=6 bytes=4
```

Phạm vi là `[start, end)`. Nó bao gồm phần đầu nhưng không có phần cuối, vì vậy độ dài của nó là end-start. Việc thêm start vào địa chỉ bắt đầu của bản gốc không tự động tạo ra NUL tại vị trí end. Bạn sẽ cần phải mang theo phạm vi riêng hoặc chuẩn bị một không gian chuỗi mới.

## Dữ liệu nhị phân có độ dài riêng

Dữ liệu chứa số 0 không nằm trong quy tắc cuối chuỗi. Nó sử dụng mảng byte và độ dài.

<!-- wave-example: book-binary-not-string -->
```wave
fun main() {
    var bytes: array<u8, 3> = [65, 0, 66];

    for (var index: i32 = 0; index < 3; index += 1) {
        println("{}", bytes[index]);
    }
}
```

Kết quả thực hiện:

```text
65
0
66
```

Số 0 thứ hai là dữ liệu thực tế. Nếu bạn hiểu giá trị này là str thì nó được coi là kết thúc ở số 0 đầu tiên và bạn không thể nhìn thấy 66 tiếp theo. Ngược lại, nếu bạn thay đổi một mảng không có NUL thành str bằng cast, thì sẽ có nguy cơ đọc vượt quá mảng. cast không phải là thao tác thêm byte kết thúc.

## Bài tập: Kiểm tra tên file

Kiểm tra xem tên tệp có kết thúc bằng `.wave` hay không và nếu chuỗi chứa `test`, hãy xuất nó thành tệp kiểm tra. Bài tập này chỉ kiểm tra các mẫu byte trong tên và không đề cập đến sự tồn tại của tệp thực tế.

### Lời giải đầy đủ

<!-- wave-example: book-string-solution -->
```wave
import("std::string::cmp")::{
    ends_with
};
import("std::string::find")::{
    contains
};

fun classify(name: str) {
    if (!ends_with(name, ".wave")) {
        println("other file");
        return;
    }

    if (contains(name, "test")) {
        println("Wave test file");
    } else {
        println("Wave source file");
    }
}

fun main() {
    classify("main.wave");
    classify("parser_test.wave");
    classify("notes.txt");
}
```

Kết quả thực hiện:

```text
Wave source file
Wave test file
other file
```

Có các chính sách riêng về cách xử lý chữ in hoa `.WAVE` và liệu có nên coi đó là một bài kiểm tra ngay cả khi test được bao gồm trong toàn bộ đường dẫn hay không. Ngay cả khi nó là một chức năng nhỏ, hoạt động của nó chỉ có thể được mô tả chính xác nếu nó được xác định đầu vào mà nó nhắm mục tiêu.
