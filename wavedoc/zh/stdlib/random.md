---
translation_set_id: stdlib-random
path: stdlib/random
locale: zh
group: stdlib
group_order: 1
order: 10
title: random:用操作系统随机性填充缓冲区
summary: OS 用熵填充缓冲区并处理部分失败。
---

## API

```text
std::random::fill
random_available() -> bool
random_fill(buffer: ptr<u8>, size: i64) -> RandomFillResult
```

size 是字节数，调用者提供存储。 `RandomFillResult` 包含正常、已写入和错误。成功后，写入的长度等于请求的长度。失败时，写入标识有效的、填充的前缀；不要使用剩余字节作为随机数据。

`random_available` 告诉您是否支持OS 随机数功能。单个请求是否成功通过random_fill的结果进行检查。仅使用 OS 熵，并且在失败时不会依赖于时间值或弱 PRNG。即使与null一起传递，size=0也会成功。 null 是负或正长度的错误。

## 运行示例

<!-- wave-example: random-api -->
```wave
import("std::random::fill")::{
    RandomFillResult, random_fill
};

fun main() -> i32 {
    var bytes: array<u8, 16>;
    var result: RandomFillResult = random_fill(&bytes[0], 16);
    if (!result.ok) {
        println("random error={}", result.error);
        return 1;
    }

    println("filled={}", result.written);
    return 0;
}
```

执行结果：

```text
filled=16
```

保存为`main.wave`并运行。每次的字节内容都不同，因此不需要特定值。如果失败，请通过result.error检查原因。如果需要，不要按字面意思输出随机字节，而是使用单独的编码。

## 如果您的请求不正确

零字节请求会成功，因为不需要写入任何内容。传递具有正长度的null会失败，因为没有目标缓冲区。以下程序在不分配内存的情况下比较这些情况。

<!-- wave-example: book-random-boundaries -->
```wave
import("std::random::fill")::{
    RandomFillResult,
    random_fill
};

fun main() -> i32 {
    var empty: RandomFillResult = random_fill(null, 0);

    if (!empty.ok || empty.written != 0) {
        return 1;
    }

    println("empty request succeeded");

    var invalid: RandomFillResult = random_fill(null, 16);

    if (invalid.ok) {
        return 2;
    }

    println("missing buffer rejected");
    return 0;
}
```

执行结果：

```text
empty request succeeded
missing buffer rejected
```

## 处理部分填充的缓冲区

如果请求16个字节失败，返回written=8，则只填充前8个字节。如果任务是创建一个 16 字节的标识符，那么它不是一个成功的标识符，因此我们丢弃整个结果并报告失败。您不应该将剩余的8个字节填充为0然后将其视为成功。

将随机字节映射到整数范围需要小心。将 `% 10` 应用于均匀分布的 u8 值会使 0-5 比 6-9 更有可能出现，因为 256 不能被 10 整除。要消除这种偏差，请拒绝值 250-255，再次绘制，然后仅对接受的值应用余数运算。

随机字节的存储和生命周期由调用者管理。使用数组时，在数组范围内处理，使用动态内存时，使用后释放。相反，返回结构并不拥有缓冲区。
