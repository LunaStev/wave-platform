---
translation_set_id: quick-reference
path: reference/syntax-quick-reference
locale: vi
group: reference
group_order: 5
order: 3
title: Cú pháp tham khảo nhanh
summary: Các khai báo, luồng điều khiển, loại, con trỏ và ngữ pháp FFI thường được sử dụng được sắp xếp trên một trang.
---

## khai báo

```wave
var value: i32 = 1;
var next: i32 = value + 1;
const LIMIT: i32 = 64;
static total: i64 = 0;
type Identifier = u64;
```

`var` là vùng, `const`/`static` là các khai báo cấp cao nhất. Các biến cục bộ khai báo rõ ràng kiểu của chúng.

## Hàm

```wave
fun max(left: i32, right: i32) -> i32 {
    if (left > right) {
        return left;
    }

    return right;
}
```

## chung chung

```wave
fun identity<T>(value: T) -> T {
    return value;
}

var value: i32 = identity<i32>(10);
```

Khi gọi một cái chung, hãy chỉ định một đối số kiểu.

## Cấu trúc và enum

```wave
struct Pair {
    left: i32;
    right: i32;
}

enum Result -> i32 {
    Ok = 0,
    Error,
}
```

## Điều kiện và vòng lặp

```wave
if (ready) {
    println("ready");
}

while (count < 10) {
    count += 1;
}

for (var i: i32 = 0; i < 10; i += 1) {
    println("{}", i);
}

match (status) {
    Ready => {
        println("ready");
    }
    0 => {
        println("zero");
    }
    _ => {
        println("other");
    }
}
```

Tiêu đề của `if`, `while`, `for` và `match` sử dụng dấu ngoặc đơn.

## Mảng và con trỏ

```wave
var values: array<i32, 4> = [1, 2, 3, 4];
var p: ptr<i32> = &values[0];
var first: i32 = deref p;
```

## đầu vào/đầu ra của bàn điều khiển

```wave
print("value = ");
println("{}", value);
input("{}", value);
```

Đối số đầu tiên là một chuỗi ký tự. Mỗi phần giữ chỗ `{}` chính xác yêu cầu một biểu thức theo sau nó và mục tiêu `input` phải có thể gán được.

## import và các vật phẩm công cộng

```wave
import("std::string::len");
import("./helpers" as helpers);
import("math")::{
    Vector
};

pub fun add(left: i32, right: i32) -> i32 {
    return left + right;
}

pub import("./extra"):: {
    increment
};
```

Đường dẫn cục bộ bắt đầu bằng `./`. Bí danh import chỉ định tên mô-đun và phần chọn import nhập các mục công khai bắt buộc vào không gian tên của tệp này. `pub import` tái xuất các mục đã chọn.

## FFI

```wave
extern(c) fun native_call(value: i32) -> i32;

export(c) fun wave_call(value: i32) -> i32 {
    return value + 1;
}
```

## Mục tiêu có điều kiện

```wave
#[target(os="linux", arch="riscv64")]
extern(c) fun platform_call(value: i32) -> i32;
```

Các phím điều kiện hỗ trợ là `arch`, `os`, `env`, `abi`. Thuộc tính kiểm soát mục cấp cao nhất tiếp theo.

## Lắp ráp nội tuyến

```wave
var result: i64 = 0;
asm {
    "mv a0, a1"
    in("a1") 7
    out("a0") result
    clobber("memory")
}
```

Văn bản hướng dẫn và tên đăng ký phụ thuộc vào mục tiêu. Khai báo tất cả các đầu vào, đầu ra và ẩn clobber cần thiết cho khối.

## kiểm tra nguồn

```shell
wavec build main.wave --emit=check
wavec print supported-targets
wavec print supported-emit-kinds
```

## Phạm vi học tập và ví dụ

Ví dụ về các biến cục bộ và các câu lệnh được hiển thị riêng biệt bên ngoài hàm là các đoạn mã được chèn vào phần thân của hàm. Các ví dụ và bài tập chạy hoàn chỉnh có trong [Wave Quá trình học tập](/docs/vi/getting-started/overview). Vui lòng kiểm tra [Thư viện chuẩn](/docs/vi/stdlib) để biết quy tắc chi tiết về bộ nhớ và các chức năng bên ngoài.
