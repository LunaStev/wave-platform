---
translation_set_id: lexical
path: language/lexical-structure
locale: zh
group: language
group_order: 2
order: 14
title: 词汇结构
summary: 描述标识符、文字、分隔符、关键字和类型名称。
---

## 标识符

标识符命名变量、函数、类型和字段。该名称区分大小写，可以包含字母、数字和 `_` 的任意组合。第一个字母中不能使用数字。字符Unicode也可用于标识符。

```wave
var request_count: i64 = 0;
var 이름: str = "Wave";
```

在实际项目中，为了工具兼容性和可搜索性，建议使用一致的命名约定。

## 句子和分隔符

大多数声明和表达式语句以`;`结尾。带主体的语句（例如函数、条件语句、循环语句和结构）使用 `{ ... }` 块。

```wave
var answer: i32 = 42;

fun double(value: i32) -> i32 {
    return value * 2;
}
```

## 字面意思

```wave
var integer: i32 = 42;
var decimal: f64 = 3.14;
var text: str = "Wave";
var letter: char = 'W';
var enabled: bool = true;
var address: ptr<u8> = null;
```

您可以使用整数、浮点数、字符串、字符、布尔值和 `null` 文字。使用 `null` 作为指针值。

## 关键字和类型名称

Wave语法中使用的主要关键字如下。

`pub`, `fun`, `extern`, `export`, `type`, `enum`, `static`, `var`, `deref`, `const`, `if`, `else`, `proto`, `struct`, `while`, `for`, `in`, `out`, `clobber`, `as`, `asm`, `import`, `return`, `continue`, `print`, `input`, `println`, `match`, `break`, `true`, `false`, `null`.

内置类型名称包括 `bool`、`char`、`byte`、`str`、整数和浮点类型、`ptr` 和 `array`。指针的书写形式为`ptr<T>`，定长数组的书写形式为`array<T, N>`。

## 字符串和字符 escape

|符号|意义|
| --- | --- |
| `\n` |LF 断线|
| `\r` | CR |
| `\t` |选项卡|
| `\\` |反斜杠|
| `\"` |双引号|
| `\xNN` |一个字节指定为恰好两个十六进制数字|

一般字符串字符保存为UTF-8。由于 `\xNN` 保留一个字节，因此无法保证整个字符串是有效的 UTF-8。字符串文字中的NUL（包括`\x00`）是编译错误。对于包含零的数据，请使用字节数组和长度。

`char` 文字必须适合 8 位值。超出该范围的字符（例如`'한'`）是错误的。它与字符串`"한"`不同。

源代码中的LF、CRLF 和CR 均被视为一个逻辑换行符。这些是有关源位置和注释终止的规则，并不意味着它们会更改文件数据的实际字节。

其他语法名称包括`variant`、`async`和`await`，异步值表示为`Future<T>`。上面的独立`var`声明块是函数内部的代码片段。

[字符串类](/docs/zh/language/arrays) · [评论](/docs/zh/language/comments)

## 故意失败的例子

如果运行check下面的程序，则会出现内部NUL错误。如果需要 0 字节，请使用 `[97, 0, 98]` 字节数组。

<!-- wave-example: reject-string-nul -->
```wave
fun main() {
    var text: str = "a\x00b";
}
```

下面的字符文字也超出了 8 位范围，因此是编译错误。要表示字符串 UTF-8，请使用 `str` 和双引号。

<!-- wave-example: reject-wide-char -->
```wave
fun main() {
    var letter: char = '한';
}
```
