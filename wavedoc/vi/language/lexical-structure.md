---
translation_set_id: lexical
path: language/lexical-structure
locale: vi
group: language
group_order: 2
order: 14
title: Cấu trúc từ vựng
summary: Mô tả các mã định danh, chữ, dấu phân cách, từ khóa và tên loại.
---

## định danh

Mã định danh đặt tên cho các biến, hàm, loại và trường. Tên có phân biệt chữ hoa chữ thường và có thể chứa bất kỳ tổ hợp chữ cái, số và `_` nào. Số không thể được sử dụng trong chữ cái đầu tiên. Các ký tự Unicode cũng có thể được sử dụng trong mã định danh.

```wave
var request_count: i64 = 0;
var 이름: str = "Wave";
```

Trong các dự án thực tế, nên sử dụng quy ước đặt tên nhất quán để có khả năng tương thích và tìm kiếm của công cụ.

## Câu và dấu phân cách

Hầu hết các câu lệnh khai báo và biểu thức đều kết thúc bằng `;`. Các câu lệnh có phần thân, chẳng hạn như hàm, câu lệnh điều kiện, câu lệnh vòng lặp và cấu trúc, sử dụng khối `{ ... }`.

```wave
var answer: i32 = 42;

fun double(value: i32) -> i32 {
    return value * 2;
}
```

## theo nghĩa đen

```wave
var integer: i32 = 42;
var decimal: f64 = 3.14;
var text: str = "Wave";
var letter: char = 'W';
var enabled: bool = true;
var address: ptr<u8> = null;
```

Bạn có thể sử dụng số nguyên, số dấu phẩy động, chuỗi, ký tự, boolean và chữ `null`. Sử dụng `null` cho các giá trị con trỏ.

## Từ khóa và tên loại

Các từ khóa chính được sử dụng trong ngữ pháp Wave như sau.

`pub`, `fun`, `extern`, `export`, `type`, `enum`, `static`, `var`, `deref`, `const`, `if`, `else`, `proto`, `struct`, `while`, `for`, `in`, `out`, `clobber`, `as`, `asm`, `import`, `return`, `continue`, `print`, `input`, `println`, `match`, `break`, `true`, `false`, `null`.

Tên loại tích hợp bao gồm `bool`, `char`, `byte`, `str`, loại số nguyên và dấu phẩy động, `ptr` và `array`. Con trỏ được viết dưới dạng `ptr<T>` và mảng có độ dài cố định được viết dưới dạng `array<T, N>`.

## Chuỗi và ký tự escape

|ký hiệu|ý nghĩa|
| --- | --- |
| `\n` |LF Ngắt dòng|
| `\r` | CR |
| `\t` |tab|
| `\\` |Dấu gạch chéo ngược|
| `\"` |dấu ngoặc kép|
| `\xNN` |Một byte được chỉ định chính xác bằng hai chữ số thập lục phân|

Các ký tự chuỗi chung được lưu dưới dạng UTF-8. Vì `\xNN` bảo toàn một byte nên không có gì đảm bảo rằng toàn bộ chuỗi là UTF-8 hợp lệ. NUL (bao gồm `\x00`) bên trong một chuỗi ký tự là lỗi biên dịch. Đối với dữ liệu chứa số 0, hãy sử dụng mảng byte và độ dài.

`char` Chữ phải khớp với giá trị 8 bit. Các ký tự vượt quá phạm vi đó, chẳng hạn như `'한'`, là lỗi. Nó khác với chuỗi `"한"`.

Riêng LF, CRLF và CR trong nguồn đều được coi là một ngắt dòng logic. Đây là các quy tắc về vị trí nguồn và chấm dứt nhận xét và không có nghĩa là chúng thay đổi byte thực tế của dữ liệu tệp.

Tên ngữ pháp bổ sung bao gồm `variant`, `async` và `await` và giá trị không đồng bộ được biểu thị dưới dạng `Future<T>`. Khối khai báo `var` độc lập ở trên là một đoạn mã bên trong một hàm.

[lớp chuỗi](/docs/vi/language/arrays) · [Bình luận](/docs/vi/language/comments)

## Ví dụ về lỗi cố ý

Nếu bạn chạy chương trình bên dưới check thì sẽ xảy ra lỗi NUL nội bộ. Nếu bạn cần 0 byte, hãy sử dụng mảng byte `[97, 0, 98]`.

<!-- wave-example: reject-string-nul -->
```wave
fun main() {
    var text: str = "a\x00b";
}
```

Ký tự bên dưới cũng vượt quá phạm vi 8 bit nên đây là lỗi biên dịch. Để biểu thị chuỗi UTF-8, hãy sử dụng `str` và dấu ngoặc kép.

<!-- wave-example: reject-wide-char -->
```wave
fun main() {
    var letter: char = '한';
}
```
