---
translation_set_id: stdlib-buffer
path: stdlib/buffer
locale: zh
group: stdlib
group_order: 1
order: 5
title: buffer:可增长的字节存储
summary: Buffer 描述了初始化、添加、查询、容量和释放规则。
---

## Buffer的含义

`std::buffer::types` 中的`Buffer` 有 `data: ptr<u8>`、`len: i64` 和 `cap: i64`。 len 是初始化和使用中的字节数，cap 是分配的总字节数。始终保持`0 <= len <= cap`。字符串 NUL 不会自动保证终止。

## 基本款API

|模块|声明|意义|
| --- | --- | --- |
| `std::buffer::alloc` | `buffer_init(out_buf: ptr<Buffer>, capacity: i64) -> i64` |初始化新存储库。不要召回已拥有的缓冲区|
|相同模块| `buffer_free(buf: ptr<Buffer>) -> i64` |解除分配。如果成功则为空|
|相同模块| `buffer_reserve(buf: ptr<Buffer>, required_cap: i64) -> i64` |确保最低满负荷。 len 已维护|
|相同模块| `buffer_clear(buf: ptr<Buffer>) -> i64` |维持容量且len=0|
|相同模块| `buffer_resize(buf: ptr<Buffer>, new_len: i64, value: u8) -> i64` |更改长度，用value填充新字节|
| `std::buffer::write` | `buffer_push(buf: ptr<Buffer>, value: u8) -> i64` |添加一个字节|
|相同模块| `buffer_append(buf: ptr<Buffer>, src: ptr<u8>, size: i64) -> i64` |size复制并添加字节|
|相同模块| `buffer_append_str(buf: ptr<Buffer>, s: str) -> i64` |添加不包括 NUL 的字符串字节|
| `std::buffer::read` | `buffer_get(buf: ptr<Buffer>, index: i64, out_value: ptr<u8>) -> i64` |读取范围内的一个字节|

状态返回 API 成功时返回 `BUFFER_OK`(0)。区分错误 INVALID、BOUNDS、OVERFLOW 和 ALLOC 与 `std::buffer::error`。该数字不解释为 OS errno。 `buffer_new` 将分配失败表示为空的Buffer，因此当需要区分失败时，请使用`buffer_init`。

## 运行示例

<!-- wave-example: buffer-api -->
```wave
import("std::buffer::alloc")::{
    Buffer, buffer_init, buffer_free
};
import("std::buffer::write")::{
    buffer_append_str, buffer_push
};
import("std::buffer::read")::{
    buffer_get
};

fun main() -> i32 {
    var data: Buffer;
    if (buffer_init(&data, 0) < 0) {
        return 1;
    }
    if (buffer_append_str(&data, "Hi") < 0 || buffer_push(&data, 33) < 0) {
        buffer_free(&data);
        return 2;
    }

    var value: u8 = 0;
    if (buffer_get(&data, 2, &value) < 0) {
        buffer_free(&data);
        return 3;
    }

    println("{} {}", data.len, value);
    if (buffer_free(&data) < 0) {
        return 4;
    }

    return 0;
}
```

执行结果：

```text
3 33
```

将其保存为`main.wave`并以`wavec run main.wave`运行。初始容量为 0 并不是失败，而是有效的空缓冲区。在进一步处理过程中释放空间。

## 寿命和故障

增加缓冲区可以更改数据。在可能重新分配的操作之后，请勿使用先前借用的地址。复制 Buffer 结构不会重复其分配，因此请为该分配指定一个所有者。

如果失败，`buffer_get`不会更改输出参数。另一方面，便利函数`buffer_at`也将错误表示为0，因此使用`buffer_get`来区分实际的0字节和失败。通过直接更改公共字段来避免创建无效的len/cap。

[内存API](/docs/zh/reference/memory-and-buffer) · [练习读取文件 Buffer](/docs/zh/practice/file-reader)

## 分别观察长度和容量

reserve 释放存储空间，但不会增加 len。 resize 更改实际使用的长度并将扩展部分初始化为指定字节。 clear 仅将使用的长度设置为 0，允许重复使用分配。

将以下程序保存为main.wave并运行。它不依赖于容量的确切增长倍数；它只是确保您拥有所需的空间。

<!-- wave-example: book-buffer-length-capacity -->
```wave
import("std::buffer::alloc")::{
    Buffer,
    buffer_init,
    buffer_reserve,
    buffer_resize,
    buffer_clear,
    buffer_free
};

fun main() -> i32 {
    var data: Buffer;

    if (buffer_init(&data, 0) < 0) {
        return 1;
    }

    if (buffer_reserve(&data, 20) < 0) {
        buffer_free(&data);
        return 2;
    }

    println("reserved length={}", data.len);

    if (buffer_resize(&data, 3, 7) < 0) {
        buffer_free(&data);
        return 3;
    }

    println("resized length={} first={}", data.len, deref data.data[0]);

    if (buffer_clear(&data) < 0) {
        buffer_free(&data);
        return 4;
    }

    println("cleared length={}", data.len);

    if (buffer_free(&data) < 0) {
        return 5;
    }

    return 0;
}
```

执行结果：

```text
reserved length=0
resized length=3 first=7
cleared length=0
```

当增加到resize时，我们传递了value=7，因此我们看到的所有三个新字节都是7。仅使用reserve保护的空间不会被读取为初始化数据。 cap 将保留在 clear 之后，并且可以再次添加到相同的 Buffer。

## 区分零字节和查找失败

buffer_get 返回状态并将实际字节写入作为输出参数。即使数据为0，也是正常成功。

<!-- wave-example: book-buffer-zero-bounds -->
```wave
import("std::buffer::alloc")::{Buffer, buffer_init, buffer_free};
import("std::buffer::write")::{buffer_push};
import("std::buffer::read")::{buffer_get};
import("std::buffer::error")::{BUFFER_ERR_BOUNDS};

fun main() -> i32 {
    var data: Buffer;

    if (buffer_init(&data, 0) < 0) {
        return 1;
    }

    if (buffer_push(&data, 0) < 0) {
        buffer_free(&data);
        return 2;
    }

    var value: u8 = 99;

    if (buffer_get(&data, 0, &value) < 0) {
        buffer_free(&data);
        return 3;
    }

    println("stored={}", value);
    value = 99;

    if (buffer_get(&data, 1, &value) == BUFFER_ERR_BOUNDS) {
        println("outside, preserved={}", value);
    }

    if (buffer_free(&data) < 0) {
        return 4;
    }

    return 0;
}
```

执行结果：

```text
stored=0
outside, preserved=99
```

第一次命中为成功，读数为 0，第二次命中为出界失败。即使失败后仍保留value=99，也不意味着它是从缓冲区读取的值。请务必一起检查状态。

## 练习题：字节累加

要添加数字 0 到 9，请重复 buffer_push 并检查每个结果。将总和存储在i64中，并仅读取范围`0 <= index < data.len`。处理完缓冲区后，我们在成功和失败路径上调用buffer_free。

<!-- wave-example: book-buffer-sum -->
```wave
import("std::buffer::alloc")::{Buffer, buffer_init, buffer_free};
import("std::buffer::write")::{buffer_push};

fun main() -> i32 {
    var data: Buffer;

    if (buffer_init(&data, 0) < 0) {
        return 1;
    }

    for (var value: i32 = 0; value < 10; value += 1) {
        if (buffer_push(&data, value as u8) < 0) {
            buffer_free(&data);
            return 2;
        }
    }

    var total: i64 = 0;

    for (var index: i64 = 0; index < data.len; index += 1) {
        total += deref data.data[index] as i64;
    }

    println("sum={}", total);

    if (buffer_free(&data) < 0) {
        return 3;
    }

    return 0;
}
```

执行结果：

```text
sum=45
```
