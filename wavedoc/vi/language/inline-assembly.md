---
translation_set_id: assembly
path: language/inline-assembly
locale: vi
group: language
group_order: 2
order: 19
title: Lắp ráp nội tuyến
summary: Mô tả hợp đồng của chuỗi lệnh của khối asm, toán hạng in/out và clobber.
---

## khối asm

`asm` là cú pháp cấp thấp để chèn trực tiếp các hướng dẫn của kiến trúc đích.

```wave
fun read_value() -> i64 {
    var result: i64 = 0;
    asm {
        "mov rax, 123"
        out("rax") result
    }

    return result;
}
```

Các chuỗi ký tự trong một khối được truyền dưới dạng danh sách lệnh tập hợp.

## đầu vào và đầu ra

```wave
var result: i64 = 0;
asm {
    "mov rax, rdi"
    in("rdi") 123
    out("rax") result
}
```

- `in("reg") expression` nối giá trị Wave vào toán hạng đầu vào.
- `out("reg") target` ghi giá trị đầu ra vào mục tiêu Wave có thể gán.
- Tên đăng ký có thể được viết dưới dạng chuỗi hoặc mã định danh.

Toán hạng đầu vào có thể bao gồm các biến, số nguyên/chuỗi ký tự, `&identifier`, `deref identifier` và số âm.

## clobber

Nếu một khối thay đổi trạng thái thanh ghi hoặc bộ nhớ không phải là đầu ra rõ ràng, thì khối đó sẽ được ghi vào `clobber(...)`.

```wave
asm {
    "nop"
    clobber("rax", "rcx", "memory")
}
```

## Kiểm tra khi sử dụng

- Cú pháp hướng dẫn phải phù hợp với kiến trúc đích và hợp đồng lắp ráp nội tuyến LLVM.
- Không được tự ý hủy các thanh ghi phải được bảo quản theo quy ước gọi.
- Đối với các khối đọc hoặc ghi bộ nhớ, hãy khai báo clobber, bao gồm `memory`.
- Nếu có thể, hãy tách biệt asm dành riêng cho kiến trúc đằng sau một hàm nhỏ.

Hành vi và tính di động của tập hợp nội tuyến không được đảm bảo chỉ bằng loại ngôn ngữ.

## Phạm vi học tập và ví dụ

[Thực hành với toàn bộ chương trình](/docs/vi/getting-started/overview) · [Thư viện chuẩn](/docs/vi/stdlib)
