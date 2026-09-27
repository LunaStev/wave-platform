---
translation_set_id: string-bytes
path: reference/string-and-bytes
locale: zh
group: stdlib
group_order: 1
order: 3
title: string:长度、搜索和范围
summary: NUL 描述终止字符串API 的字节单位以及返回值。
---

## 字符串存储和参数条件

该模块的`str`参数必须是可访问的NUL终止字节。长度和搜索索引以字节为单位。常规字符存储为 UTF-8，但字节搜索为 Unicode，没有标准化或逐字符拆分。不要假设返回的索引是字符边界。

## 与长度比较

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

`len` 不包括最后一个NUL。 `cmp` 通过结果的符号判断顺序。返回值不被解释为 Unicode 字符顺序或特定于语言的字典排序。这些函数不分配内存，也不更改其输入。

## 搜索

获取您需要的名称，例如`import("std::string::find")::{find, contains, count};`。

|函数声明|结果|
| --- | --- |
| `find(s: str, needle: str) -> i32` |第一场比赛地点。 -1 如果不存在，0 为空 needle|
| `contains(s: str, needle: str) -> bool` |包括或不包括。空needle是true|
| `count(s: str, needle: str) -> i32` |非重叠匹配的数量。 Bin needle 为 0|
| `find_char(s: str, c: u8) -> i32` |字节的第一个位置或-1|
| `rfind_char(s: str, c: u8) -> i32` |字节的最后位置或-1|
| `contains_char(s: str, c: u8) -> bool` |该字节的存在|
| `count_char(s: str, c: u8) -> i32` |有问题的字节数|

名称中的 `c` `*_char` 是一个字节，而不是 Unicode 代码点。字符串末尾的NUL本身不包含在搜索目标中。

## 范围不包括空格

```text
std::string::trim
trim_left_index(s: str) -> i32
trim_right_index(s: str) -> i32
trim_range(s: str, out_start: ptr<i32>, out_end: ptr<i32>)
```

`trim_range` 将半开范围 `[start, end)`（不包括空格 ASCII）写入输出参数。两个输出指针都必须指向可写整数。它不会修改原始文本或创建新字符串。如果一切都是空白，它就变成一个空范围。

## 运行示例

将其保存到`main.wave`并运行`wavec run main.wave`。

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

执行结果：

```text
2 6
2 2
```

## 相关功能

`std::string::ascii` 的分类/大小写转换适用于范围 ASCII。 `djb2_32` 和 `std::string::hash` 的 `fnv1a_64` 不用于加密哈希或密码存储。对于包含NUL的数据，请使用[bytes](/docs/zh/stdlib/bytes)。

## 与空搜索词重叠的模式

通过实际值查看搜索函数的边缘行为可以更轻松地确定调用条件。

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

执行结果：

```text
empty find=0
empty count=0
nonoverlapping=2
empty needle is contained
```

find 的成功值为 0，为第一个位置。 count 的成功值为 0 是不匹配或空搜索词规则的结果。没有两个价值观会受到同等对待。如果需要忽略大小写或Unicode标准化，则必须在此字节搜索之前和之后实施单独的策略。

## trim 将范围复制到新字符串

trim_range返回的范围内没有新的结尾NUL。复制到单独目的地时，保留长度+1空间，最后一个字节直接写为0。

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

执行结果：

```text
Wave
```

长度为 16 的字符串不适合此目的地。这是因为您需要最后一个NUL。即使长度为 0，写入 destination[0]=0 也会产生有效的空字符串。目标本地数组一直存在到 main 末尾，因此我们在其中进行打印。

## 字符串 API 使用顺序

设计字符串 API 时，请指定其输入是否以 NUL 结尾、索引是否计数字节以及结果是否借用源范围或拥有新分配。借用的范围取决于源的生命周期。分配的结果必须指定谁释放它。

阅读[弦乐学习篇](/docs/zh/language/strings)了解基本概念，阅读[bytes](/docs/zh/stdlib/bytes)了解数据，包括NUL。
