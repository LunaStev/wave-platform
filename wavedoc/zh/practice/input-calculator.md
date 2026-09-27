---
translation_set_id: practice-input-calculator
path: practice/input-calculator
locale: zh
group: practice
group_order: 4
order: 1
title: 项目：验证输入的计算器
summary: 关联输入、边界检查、函数和退出代码。
---

## 目标和行动

输入数量和单价并计算总计。输入两个整数，以空格或换行符分隔。本例仅接受数量 1 到 1000，单价 0 到 100000，因此是在乘法范围i32 内计算的。

保存到`main.wave`，运行`wavec run main.wave`，然后输入`3 1200`。程序输出的结果如下。在终端中输入的字符的可见性与程序输出是分开的。

<!-- wave-example: input-calculator -->
```wave
fun total(quantity: i32, price: i32) -> i32 {
    return quantity * price;
}

fun main() -> i32 {
    var quantity: i32 = 0;
    var price: i32 = 0;
    input("{} {}", quantity, price);
    if (quantity < 1 || quantity > 1000 || price < 0 || price > 100000) {
        println("out of range");
        return 1;
    }

    println("total={}", total(quantity, price));
    return 0;
}
```

执行结果：

```text
total=3600
```

## 还要检查是否有故障

如果您输入 `0 1200`，它需要 `out of range` 并退出代码 1。 `input` 中的数字解析失败和检查程序的工作范围是两个不同的步骤。非数字标记/类型超出范围/在所需输入EOF之前是输入失败。这个内置输入不是一个返回错误并需要重新输入的接口。如果您需要可恢复的解析器，请通过使用 [io](/docs/zh/stdlib/files-io) 读取字节来自行配置验证过程。

## 扩展练习和评论

将折扣率作为第三个输入，检查是否为 0 到 100。为避免中间较大的乘法，需要使用 i64 扩大计算范围，并在缩小结果时检查范围。如果您只是加宽最后一个结果类型，则可能已经在窄类型上执行了中间计算。

[控制台I/O](/docs/zh/language/console-io-and-formatting) · [下一步：读取文件](/docs/zh/practice/file-reader)

## 为什么要先确定计算范围呢？

最大的输入是数量1000和单价100000。两个值相乘是100000000，所以落在范围i32内。需要进行范围检查，以便可以按原样使用total函数的乘法结果。

main，接收输入，负责输入和错误消息，total只负责计算。即使您后来更改为从文件中读取命令，您仍然可以使用计算功能。

## 完整的折扣计算

您可以通过添加折扣百分比将总数乘以 100。转换为执行中间计算i64，然后应用折扣率。由于整数除法的小数部分被丢弃，因此本示例将折扣后的金额截断为整数。

<!-- wave-example: book-calculator-discount -->
```wave
fun discounted_total(quantity: i32, price: i32, discount: i32) -> i64 {
    var subtotal: i64 = (quantity as i64) * (price as i64);
    var remaining: i64 = 100 - (discount as i64);
    return subtotal * remaining / 100;
}

fun main() -> i32 {
    var quantity: i32 = 0;
    var price: i32 = 0;
    var discount: i32 = 0;

    input("{} {} {}", quantity, price, discount);

    if (quantity < 1 || quantity > 1000) {
        println("invalid quantity");
        return 1;
    }

    if (price < 0 || price > 100000) {
        println("invalid price");
        return 2;
    }

    if (discount < 0 || discount > 100) {
        println("invalid discount");
        return 3;
    }

    println("total={}", discounted_total(quantity, price, discount));
    return 0;
}
```

执行结果：

```text
total=3240
```

以上结果是输入`3 1200 10`时的结果。输出为 3240，从原始总数 3600 中减去 10%。

|输入|预期结果|检查路径|
| --- | --- | --- |
| `3 1200 0` | `total=3600` |没有折扣|
| `3 1200 100` | `total=0` |全额折扣|
| `3 1200 101` | `invalid discount` |折扣率超出范围|
| `1000 100000 0` | `total=100000000` |最大输入|
| `0 1200 10` | `invalid quantity` |数量低于范围|

## 下一个练习

尝试更改函数以舍入小数位数。在此程序中，仅接受正数，在除以 100 之前加上 50 将四舍五入为整数。当您为 1 输入 99 且折扣率为 50 时，您可以比较 49 的切割结果和 50 的舍入结果。
