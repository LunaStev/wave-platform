---
translation_set_id: modules-ffi
path: language/modules-imports-and-ffi
locale: vi
group: language
group_order: 2
order: 11
title: 11. Mô-đun và mã chung
summary: Tìm hiểu về cách lấy tên công khai và đối số loại rõ ràng.
---

## Lý do chia file

Khi chương trình phát triển, bạn sẽ dễ dàng tìm thấy nó hơn bằng cách nhóm các hàm liên quan lại với nhau thay vì đặt tất cả các hàm vào main.wave. Ranh giới mô-đun xác định tên nào được hiển thị với mã khác. Generics là công cụ tái sử dụng cùng một tác vụ với các loại khác nhau, không phụ thuộc vào việc phân tách tệp.

Trong chương này, chúng ta tạo một chương trình gồm hai tệp và tìm hiểu các bí danh mô-đun, lựa chọn import cũng như các hàm và cấu trúc chung.

## chương trình hai tập tin

Tạo helpers.wave và main.wave trong cùng một thư mục.

helpers.wave Tất cả:

```wave
pub fun double(value: i32) -> i32 {
    return value * 2;
}
```

main.wave Tất cả:

<!-- wave-example: book-module-two-files -->
```wave
import("./helpers")::{
    double
};

fun main() {
    var result: i32 = double(21);

    println("{}", result);
}
```

Kết quả thực hiện:

```text
42
```

Chạy `wavec run main.wave` trong thiết bị đầu cuối. helpers.wave cũng không được chạy riêng lẻ. Nguồn yêu cầu được kết nối qua import.

pub ở phía trước hàm trong helpers cho biết rằng nó có thể được nhập bởi các mô-đun khác. Những chức năng phụ trợ không cần tiếp xúc với thế giới bên ngoài cũng không cần phải công khai. Ngay cả khi bạn thay đổi cách triển khai nội bộ của một mô-đun, bạn có thể giảm bớt các thay đổi đối với mã bạn sử dụng bằng cách giữ nguyên hợp đồng của chức năng công khai.

## Đường cơ sở cho đường dẫn tương đối

`./helpers` liên quan đến thư mục của tệp nguồn đã tạo ra câu import. Khi chạy chương trình, hãy tách chương trình ra khỏi thư mục làm việc nơi ghi tệp I/O. Các bước tìm tệp import và các bước tìm input.txt trong khi chạy là khác nhau.

Đối với import cục bộ, có thể bỏ qua phần mở rộng `.wave`. Nếu bạn chia thư mục, hãy viết đường dẫn `./module` theo vị trí. Đừng nhầm lẫn các đường dẫn tương đối cục bộ với các đường dẫn tìm nạp tên của các gói phụ thuộc.

## Chọn import và bí danh

Lựa chọn import chỉ khiến các tên công khai mong muốn được sử dụng trực tiếp trong tệp hiện tại. Nếu có xung đột về tên hoặc bạn muốn tiết lộ chức năng đó thuộc về mô-đun nào, hãy sử dụng bí danh.

<!-- wave-example: book-module-alias -->
```wave
import("std::string::len" as strings);

fun main() {
    var length: i32 = strings::len("Wave");

    println("{}", length);
}
```

Kết quả thực hiện:

```text
4
```

strings là bí danh mô-đun được xác định trong tệp này. `strings::len` sử dụng tên của mô-đun. Dấu chấm để truy cập trường và `::` để phân tách mô-đun là các ký hiệu khác nhau.

Không sử dụng tùy chọn import và bí danh import cùng nhau trong một câu. Dù phong cách của bạn là gì, hãy sử dụng nó một cách nhất quán trong toàn bộ tệp để có thể dễ dàng đọc được nguồn gốc của tên.

## Thư viện và gói tiêu chuẩn

Đường dẫn `std::` trỏ đến thư viện chuẩn. Người dùng xuất bản API và import các mô-đun được yêu cầu. Không phải tất cả các chức năng thư viện tiêu chuẩn đều được tự động đặt vào không gian tên hiện tại.

Đường dẫn gói bên ngoài bắt đầu từ tên gói. Vị trí của gói được cung cấp bởi các tùy chọn trình biên dịch hoặc trình quản lý gói. Trước tiên, hãy tìm hiểu các ranh giới với các mô-đun cục bộ, sau đó tìm hiểu cách quản lý các phần phụ thuộc trong [Vex Cách sử dụng](/docs/vi/whale/vex-package-manager).

## Sử dụng chức năng giống nhau cho từng loại

Các hàm sau trả về nguyên văn đầu vào của chúng: Đối với i32 và str, hãy sử dụng tham số loại T để tránh viết cùng một mã hai lần.

<!-- wave-example: book-generic-identity -->
```wave
fun identity<T>(value: T) -> T {
    return value;
}

fun main() {
    var number: i32 = identity<i32>(42);
    var text: str = identity<str>("Wave");

    println("{} {}", number, text);
}
```

Kết quả thực hiện:

```text
42 Wave
```

T là nơi nhập kiểu chứ không phải nhập giá trị số nguyên được truyền trong quá trình thực thi. Gọi nó bằng cách chỉ định đối số loại như `<i32>`. Các hàm chung của người dùng thông thường không bỏ qua các đối số kiểu.

identity<str> không trùng lặp byte chuỗi bằng cách phân bổ lại chúng. Trả về giá trị như cũ. Ngữ pháp của generics không thay đổi các quy tắc sao chép và quyền sở hữu dữ liệu.

## Các hoạt động được yêu cầu bởi cơ thể chung

Chỉ vì có một tham số kiểu không có nghĩa là tất cả các thao tác có thể được sử dụng trên tất cả các kiểu. minimum bên dưới nên được sử dụng như một loại thực tế có thể so sánh được.

<!-- wave-example: book-generic-minimum -->
```wave
fun minimum<T>(left: T, right: T) -> T {
    if (left < right) {
        return left;
    }

    return right;
}

fun main() {
    var small: i32 = minimum<i32>(7, 4);
    var wide: i64 = minimum<i64>(100, 20);

    println("{} {}", small, wide);
}
```

Kết quả thực hiện:

```text
4 20
```

Nếu bạn thay đổi đối số loại, `<` và kết quả trả về được sử dụng trong văn bản phải thuộc loại tương ứng. Khi đọc các lỗi chung, hãy kiểm tra cả tổ hợp kiểu được gọi và thao tác mà thân hàm yêu cầu.

## cấu trúc chung

Hãy tạo Pair, liên kết hai giá trị khác nhau.

<!-- wave-example: book-generic-pair -->
```wave
struct Pair<A, B> {
    first: A;
    second: B;
}

fun main() {
    var item: Pair<i32, str> = Pair<i32, str> {
        first: 7,
        second: "seven"
    };

    println("{} {}", item.first, item.second);
}
```

Kết quả thực hiện:

```text
7 seven
```

Pair<i32, str> và Pair<i64, str> là các loại cụ thể khác nhau. Thứ tự của các đối số kiểu cũng có ý nghĩa. Đọc phần khai báo và mã tạo để xem loại first và second được xác định ở đâu.

## Tên và hợp đồng của API đã được phát hành

Khi bạn xuất bản một hàm, bạn không chỉ chỉ định tên mà còn cả đơn vị đầu vào, giá trị trả về, lỗi và quyền sở hữu. Ví dụ: vòng lặp mà người gọi sẽ viết tùy thuộc vào việc read đọc độ dài tối đa hay độ dài chính xác.

pub là phạm vi công khai giữa các mô-đun Wave. Điều này khác với export (c), xuất các ký hiệu bên ngoài để các ngôn ngữ khác gọi. Bạn có thể xem ví dụ đầy đủ về liên kết hai ngôn ngữ tại [Xem FFI](/docs/vi/language/modules-imports-and-ffi).

## Bài tập và lời giải đầy đủ

Tạo một hàm công khai square trên math.wave và gọi nó là bí danh trên main.wave để in lũy thừa của 3 và 5.

math.wave:

```wave
pub fun square(value: i32) -> i32 {
    return value * value;
}
```

main.wave:

<!-- wave-example: book-module-solution -->
```wave
import("./math" as math);

fun main() {
    println("{} {}", math::square(3), math::square(5));
}
```

Kết quả thực hiện:

```text
9 25
```

Nếu lỗi là chức năng đã biến mất, trước tiên hãy kiểm tra đường dẫn import và pub. Nếu có xung đột tên, hãy kiểm tra xem cuộc gọi có bí danh hay không. Nếu đó là lỗi loại, hãy kiểm tra đầu vào của hàm và loại đối số được truyền. Đừng cố gắng giải quyết các vấn đề khác nhau bằng cách chỉnh sửa một đường dẫn.


## C Chức năng nhập

```wave
extern(c) fun puts(text: ptr<i8>) -> i32;
```

Tên ABI có thể được theo sau bởi tên biểu tượng thực tế dưới dạng chuỗi.

```wave
extern(c, "native_symbol") fun local_name(value: i32) -> i32;
```

## Wave Chức năng xuất

```wave
export(c) fun wave_add(left: i32, right: i32) -> i32 {
    return left + right;
}
```

`extern` và `export` có thể được sử dụng dưới dạng các hàm và khối đơn lẻ. Hàm được xuất phải có chữ ký cụ thể ABI và do đó không thể chung chung.

## Thuộc tính điều kiện mục tiêu

Thuộc tính điều kiện mục tiêu có thể được gắn vào các mục cấp cao nhất.

```wave
#[target(os="linux", arch="x86_64")]
extern(c) fun platform_call(value: i32) -> i32;
```

Các khóa điều kiện là `arch`, `os`, `env`, `abi` và các thuộc tính áp dụng cho mục cấp cao nhất tiếp theo.

## Kết nối với hàm C do chính bạn viết

Phòng thí nghiệm này dành cho môi trường gốc với trình biên dịch C. Nối một hàm số nguyên mà không cần cấp phát thư viện hoặc xử lý chuỗi.

`native.c`:

```c
#include <stdint.h>
int32_t native_double(int32_t value) { return value * 2; }
```

`main.wave`:

```wave
extern(c) fun native_double(value: i32) -> i32;

fun main() {
    println("{}", native_double(21));
}
```

Chạy Linux/macOS từ một thiết bị đầu cuối trong cùng thư mục làm việc.

```shell
cc -c native.c -o native.o
wavec build main.wave native.o -o ffi-example
./ffi-example
```

Đầu ra dự kiến là `42`. Trong trình bao nhà phát triển của MSVC trong Windows, tạo object với `cl /c native.c /Fonative.obj` và kết nối với `wavec build main.wave native.obj -o ffi-example.exe`. Kiến trúc nguồn và đích của object phải giống nhau. Các ví dụ chỉ sử dụng các giá trị nhỏ. Để chuyển một giá trị lớn cho hàm C, phạm vi nhân của trang C cũng phải được đảm bảo riêng.

Đường dẫn tệp cục bộ bắt đầu bằng `./`.
