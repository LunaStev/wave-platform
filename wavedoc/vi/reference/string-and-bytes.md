---
translation_set_id: string-bytes
path: reference/string-and-bytes
locale: vi
group: stdlib
group_order: 1
order: 3
title: string: Độ dài, tìm kiếm và phạm vi
summary: NUL Mô tả đơn vị byte của chuỗi kết thúc API và giá trị trả về.
---

## Điều kiện lưu trữ và đối số chuỗi

Đối số `str` cho mô-đun này phải là byte kết thúc NUL có thể truy cập được. Độ dài và chỉ mục tìm kiếm được tính bằng byte. Các ký tự thông thường được lưu trữ dưới dạng UTF-8, nhưng tìm kiếm theo byte là Unicode mà không cần chuẩn hóa hoặc phân tách từng ký tự. Đừng cho rằng chỉ mục trả về là ranh giới ký tự.

## So sánh với chiều dài

```text
std::string::len
len(s: str) -> i32
is_empty(s: str) -> bool

std::string::cmp
eq(a: str, b: str) -> bool
cmp(a: str, b: str) -> i32
starts_with(s: str, prefix: str) -> bool
ends_with(s: str, suffix: str) -> bool
```

`len` không bao gồm NUL cuối cùng. `cmp` Thứ tự được đánh giá bằng dấu của kết quả. Giá trị trả về không được hiểu là thứ tự ký tự Unicode hoặc sắp xếp từ điển theo ngôn ngữ cụ thể. Các chức năng này không phân bổ bộ nhớ và không thay đổi đầu vào của chúng.

## tìm kiếm

Lấy tên bạn cần, chẳng hạn như `import("std::string::find")::{find, contains, count};`.

|khai báo hàm|kết quả|
| --- | --- |
| `find(s: str, needle: str) -> i32` |Địa điểm trận đấu đầu tiên. -1 nếu không có, 0 nếu trống needle|
| `contains(s: str, needle: str) -> bool` |Bao gồm hoặc không. Trống needle là true|
| `count(s: str, needle: str) -> i32` |Số lượng kết quả trùng khớp không trùng lặp. Thùng needle là 0|
| `find_char(s: str, c: u8) -> i32` |vị trí đầu tiên của byte hoặc -1|
| `rfind_char(s: str, c: u8) -> i32` |Vị trí cuối cùng của byte hoặc -1|
| `contains_char(s: str, c: u8) -> bool` |Sự tồn tại của byte đó|
| `count_char(s: str, c: u8) -> i32` |số byte được đề cập|

`c` trong tên `*_char` là một byte, không phải là điểm mã Unicode. Bản thân NUL ở cuối chuỗi không được đưa vào mục tiêu tìm kiếm.

## Phạm vi không bao gồm khoảng trắng

```text
std::string::trim
trim_left_index(s: str) -> i32
trim_right_index(s: str) -> i32
trim_range(s: str, out_start: ptr<i32>, out_end: ptr<i32>)
```

`trim_range` ghi phạm vi nửa mở `[start, end)` không bao gồm khoảng trắng ASCII vào đối số đầu ra. Cả hai con trỏ đầu ra phải trỏ đến số nguyên có thể ghi. Nó không sửa đổi văn bản gốc hoặc tạo chuỗi mới. Nếu mọi thứ đều trống, nó sẽ trở thành một phạm vi trống.

## Ví dụ đang chạy

Lưu nó vào `main.wave` và chạy `wavec run main.wave`.

<!-- wave-example: string-api -->
```wave
import("std::string::find")::{
    find, count
};
import("std::string::trim")::{
    trim_range
};

fun main() {
    var start: i32 = 0;
    var end: i32 = 0;
    trim_range("  Wave  ", &start, &end);
    println("{} {}", start, end);
    println("{} {}", find("banana", "na"), count("aaaa", "aa"));
}
```

Kết quả thực hiện:

```text
2 6
2 2
```

## Tính năng liên quan

Việc chuyển đổi phân loại/trường hợp của `std::string::ascii` dành cho phạm vi ASCII. `djb2_32` và `fnv1a_64` của `std::string::hash` không được sử dụng để băm mật mã hoặc lưu trữ mật khẩu. Đối với dữ liệu chứa NUL, hãy sử dụng [bytes](/docs/vi/stdlib/bytes).

## Các mẫu trùng lặp với cụm từ tìm kiếm trống

Xem hành vi biên của hàm tìm kiếm với các giá trị thực tế giúp xác định điều kiện gọi dễ dàng hơn.

<!-- wave-example: book-string-empty-search -->
```wave
import("std::string::find")::{find, contains, count};

fun main() {
    println("empty find={}", find("Wave", ""));
    println("empty count={}", count("Wave", ""));
    println("nonoverlapping={}", count("aaaa", "aa"));

    if (contains("Wave", "")) {
        println("empty needle is contained");
    }
}
```

Kết quả thực hiện:

```text
empty find=0
empty count=0
nonoverlapping=2
empty needle is contained
```

Giá trị thành công của find, 0, là vị trí đầu tiên. Giá trị thành công bằng 0 cho count là kết quả của quy tắc cụm từ tìm kiếm trống hoặc không khớp. Không có hai giá trị nào được đối xử như nhau. Nếu cần bỏ qua trường hợp hoặc cần chuẩn hóa Unicode thì các chính sách riêng biệt phải được triển khai trước và sau tìm kiếm byte này.

## trim Sao chép phạm vi sang chuỗi mới

Không có phần kết thúc mới NUL trong phạm vi được trả về bởi trim_range. Khi sao chép sang một đích riêng biệt, hãy đặt trước độ dài + 1 dấu cách và ghi trực tiếp byte cuối cùng là 0.

<!-- wave-example: book-trim-copy -->
```wave
import("std::string::trim")::{trim_range};

fun main() -> i32 {
    var original: str = "  Wave  ";
    var start: i32 = 0;
    var end: i32 = 0;
    var destination: array<u8, 16>;

    trim_range(original, &start, &end);

    var length: i32 = end - start;

    if (length >= 16) {
        return 1;
    }

    for (var index: i32 = 0; index < length; index += 1) {
        destination[index] = original[start + index];
    }

    destination[length] = 0;
    println("{}", &destination[0] as str);
    return 0;
}
```

Kết quả thực hiện:

```text
Wave
```

Một chuỗi có độ dài 16 sẽ không vừa với đích này. Điều này là do bạn cần NUL cuối cùng. Ngay cả khi độ dài bằng 0, việc viết destination[0]=0 sẽ dẫn đến một chuỗi trống hợp lệ. Mảng cục bộ đích tồn tại cho đến hết main, vì vậy chúng tôi in bên trong mảng đó.

## Chuỗi API Thứ tự sử dụng

Khi thiết kế API chuỗi, hãy chỉ định xem đầu vào của nó có được kết thúc bằng NUL hay không, liệu các chỉ mục có đếm byte hay không và liệu kết quả có mượn một phạm vi nguồn hay sở hữu một phân bổ mới hay không. Phạm vi mượn phụ thuộc vào thời gian tồn tại của nguồn. Một kết quả được phân bổ phải chỉ rõ ai giải phóng nó.

Đọc [Chương học về chuỗi](/docs/vi/language/strings) để biết các khái niệm cơ bản và [bytes](/docs/vi/stdlib/bytes) để biết dữ liệu bao gồm NUL.
