---
translation_set_id: learn-strings
path: language/strings
locale: zh
group: language
group_order: 2
order: 7
title: 7. 字符串、字符和字节
summary: 区分字符串和char、UTF-8字节长度、NUL、搜索和二进制数据。
---

## 屏幕上的字母和内存中的字节

您在屏幕上看到字母，但字节存储在内存中。特别是在一个字符有多个UTF-8字节的情况下，例如韩语，如果“长度”和“字符数”互换使用，很容易出错。

本章区分str和char，结尾NUL，escape，搜索位置和二进制数据。每个示例都是一个完整的程序。

## 字符串文字和输出

<!-- wave-example: book-string-literals -->
```wave
fun main() {
    var greeting: str = "안녕하세요";

    println("{}", greeting);
    println("line one\nline two");
    println("quote: \"Wave\"");
}
```

执行结果：

```text
안녕하세요
line one
line two
quote: "Wave"
```

双引号内的常规字符表示为UTF-8。 escape 表示难以直接从源写入的字节。 `\n` 是换行符字节，不会打印两个字符：反斜杠和n。要打印反斜杠本身，请使用 `\\`。

## 长度是字节数

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

执行结果：

```text
ASCII=4
Korean=3
mixed=7
```

`len` 计算终止 NUL 之前的字节，不可见字符。字符`한`在UTF-8中占用三个字节。字符计数、Unicode 代码点计数和字节计数通常不可互换。显示宽度还取决于字体和组合字符等因素。

因此，在任意字节位置剪切字符串并将其显示在屏幕上的函数必须单独考虑Unicode边界。明确定义输入条件，无论是只处理ASCII的程序还是一般的Unicode文本。

## NUL 末端和长度

str 使用0字节表示结束。 len 不包括长度中的最后一个字节。将 NUL 放入字符串文字中是错误的。以下是故意错误的示例：

```wave
fun main() {
    var text: str = "left\x00right";
}
```

源代码中的`\xNN`指定一个字节，恰好有两个十六进制数字。 `\x41` 代表字节A，即 65。由于将常规字符写入UTF-8 和插入任意字节是不同的，所以并非所有可写入`\xNN` 的str 都是有效的UTF-8。

## char 不包含整个字符 Unicode

char 是一个无符号8位字符值。您可以使用表示单字节范围内的值的文字，例如`'A'`。 `'한'` 是一个错误，因为它不属于这个范围。 `"한"` 是一个单独的 str，具有多个 UTF-8 字节。

不要试图总是将一个字母放在一个char。您必须首先决定处理文本所需的单位是字节还是Unicode代码点。

## 字符串比较

要比较字符串内容，请使用函数std。下面是一个检查相同内容和大小写差异的程序。

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

执行结果：

```text
same bytes
case differs
name matches
```

此比较比较字节字符串。它不会自动执行特定于语言的大小写转换或 Unicode 规范化。即使在比较文件名时，OS中的文件名相等规则也与简单的字符串比较不同。

## 搜索结果的单位和失败

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

执行结果：

```text
first=2 missing=-1
matches=2
```

find 返回第一个位置或-1。索引0也成功，因此用`result >= 0`进行检查。 count不是位置，而是非重叠匹配的数量。上面，aa 是数字 2，因为它匹配 0~1 和 2~3。

Bin needle 也是合同的一部分。 find 返回 0，contains 返回 true，count 返回 0。不要以为函数名在同一个模块中，就连返回方法都一样。

## 删除空格与创建新字符串不同

trim_range 返回排除空格的范围，无需修改或复制原始文本。由于我们收到一个输出指针，因此我们首先准备一个整数来存储结果。

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

执行结果：

```text
start=2 end=6 bytes=4
```

范围是`[start, end)`。它包括开头但不包括结尾，因此其长度为end-start。将 start 添加到原始地址的起始地址不会自动在 end 位置创建 NUL。您需要单独携带范围或准备一个新的字符串空间。

## 二进制数据有单独的长度

包含零的数据不属于字符串结尾规则的范围。它使用字节数组和长度。

<!-- wave-example: book-binary-not-string -->
```wave
fun main() {
    var bytes: array<u8, 3> = [65, 0, 66];

    for (var index: i32 = 0; index < 3; index += 1) {
        println("{}", bytes[index]);
    }
}
```

执行结果：

```text
65
0
66
```

第二个0是实际数据。如果将其解释为 str，则将其视为以第一个 0 结尾，并且看不到后面的 66。相反，如果将不带 NUL 的数组更改为带 cast 的 str，则存在读取超出数组范围的风险。 cast 不是添加终止字节的操作。

## 练习：检查文件名

检查文件名是否以`.wave`结尾，如果字符串中包含`test`，则输出到测试文件。此练习仅检查名称中的字节模式，并不解决实际文件存在的问题。

### 完整解答

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

执行结果：

```text
Wave source file
Wave test file
other file
```

关于如何处理大写字母`.WAVE`以及即使test包含在整个路径中是否将其视为测试，都有单独的策略。即使它是一个小函数，只有确定它的目标输入是什么，才能准确地描述它的操作。
