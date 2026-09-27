---
translation_set_id: console-io-formatting
path: language/console-io-and-formatting
locale: zh
group: language
group_order: 2
order: 16
title: 控制台输入、输出和格式化
summary: print、println、input 描述句子和占位符规则。
---

## 输入/输出语句

Wave 提供`print`、`println`、`input` 作为控制台输入/输出语句。

```wave
fun main() {
    var count: i32 = 0;
    input("{}", count);
    print("count = ");
    println("{}", count);
}
```

每个句子都以`;`结尾。第一个参数必须是字符串文字。变量或计算字符串不能用作 format 参数。

## 占位符

只有两个字符 `{}` 是占位符。

```wave
println("name = {}, score = {}", name, score);
```

占位符的数量和后面的表达式的数量必须完全相同。如果数字不同，则属于语法错误。

```wave
println("{} {}", one);
// 오류: 자리표시자 2개, 값 1개
println("plain text", one);
// 오류: 자리표시자 없음, 값이 남음
```

其他形式的花括号保留为纯文本。此语法中没有命名或编号的占位符。

## print 和 println

`print` 按原样打印格式化文本，`println` 添加换行符。

```wave
print("loading...");
println("done");
```

格式化参数使用标量值，例如整数、浮点数、字符串和指针。数组和结构不能用作格式化参数。

## input 目标

`input` 将读取的值存储在目标中，因此 format 之后的所有表达式都必须是可写位置。

```wave
var number: i32 = 0;
input("{}", number);
```

变量、字段和取消引用的存储位置可以用作目标。文字和计算结果不能用作输入。

如果无法将所有输入值转换为请求的类型，则程序会以失败状态退出。

## 运行时边界

这些语句使用来自 hosted 环境的控制台输入和输出。在独立环境中，内核或设备提供的输入/输出必须定义为函数或FFI边界。

## 输入值及范围

bool 输入仅接受 `0` 和 `1`。它不会将 2 解释为 true 或接受字符串 `true` 作为相同的输入。整数输入必须在目标整数宽度范围内。 128、256、512 和 1024 位整数也会根据类型的整体宽度进行处理。

在所需输入EOF失败之前，格式错误，超出范围。内置的input不是一个返回失败并重新进入的函数，而是一个失败时终止进程的输入函数。如果需要可恢复的输入处理，请使用 io 读取字节并构造一个单独的解析器。

[输入计算器练习](/docs/zh/practice/input-calculator) · [文件和io](/docs/zh/stdlib/files-io)
