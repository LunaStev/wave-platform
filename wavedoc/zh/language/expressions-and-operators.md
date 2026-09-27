---
translation_set_id: expressions
path: language/expressions-and-operators
locale: zh
group: language
group_order: 2
order: 3
title: 3. 算术、比较和转换
summary: 学习计算顺序、整数除法、按位运算和cast。
---

## 一起查看计算结果和计算类型

表达式是计算值的代码。变量名、文字、函数调用以及用运算符连接多个值的表达式都是表达式。在数学中，即使它看起来像同一个方程，但根据它是整数还是实数以及包含多少位，结果也会有所不同。

本章从简单的计算开始，介绍括号、除法、逻辑运算、按位运算和强制转换。每个示例都是一个完整的 main.wave 文件，您可以使用 `wavec run main.wave` 运行。

## 括号内的范围

<!-- wave-example: book-precedence -->
```wave
fun main() {
    var first: i32 = 2 + 3 * 4;
    var second: i32 = (2 + 3) * 4;

    println("{} {}", first, second);
}
```

执行结果：

```text
14 20
```

乘法先于加法计算，因此第一个表达式是 2+12。在第二个方程中，我们将括号中的总和设为 5，然后乘以 4。目标不是使用更少的括号。建议使用它，以便读者可以轻松了解计算的范围。

即使多次使用相同的运算符，链接的方向也很重要。 `20 - 5 - 3` 是 `(20 - 5) - 3`，即 12。`20 - (5 - 3)` 是 18。确切的完整序列位于 [操作员参考](/docs/zh/language/expressions-and-operators)。

## 整数除法和余数

两个整数相除不会产生小数浮点结果。对商使用除法，对余数使用余数运算符。

<!-- wave-example: book-division -->
```wave
fun main() {
    var items: i32 = 17;
    var box_size: i32 = 5;
    var full_boxes: i32 = items / box_size;
    var remaining: i32 = items % box_size;

    println("boxes={} remaining={}", full_boxes, remaining);
    println("negative quotient={}", -17 / 5);
}
```

执行结果：

```text
boxes=3 remaining=2
negative quotient=-3
```

将 17 个项目分成 5 个组，生成 3 个完整组和 2 个剩余项目。有符号除法向零截断，因此 -17/5 为 -3。这与向下舍入到负无穷大不同。

不能除以 0。 signed 最小值除以 -1 得到的值不属于同一类型。接受这些输入的函数应在除法之前进行检查或使用checked数学API。

## 转换点会改变结果。

下面两个表达式都存储在变量f64中，但计算过程不同。

<!-- wave-example: book-division-conversion -->
```wave
fun main() {
    var already_divided: f64 = (7 / 2) as f64;
    var floating_division: f64 = (7 as f64) / (2 as f64);

    if (already_divided == 3.0) {
        println("integer division lost the fraction");
    }

    if (floating_division == 3.5) {
        println("floating division kept the fraction");
    }
}
```

执行结果：

```text
integer division lost the fraction
floating division kept the fraction
```

第一个表达式进行整数除法得到3，然后将其转换为f64。第二个在执行浮点除法之前将操作数转换为f64。为最终变量选择更宽的类型无法恢复之前丢失的信息。

浮点值是近似值。看起来像相同十进制值的两个结果可能不适合进行精确相等的比较。选择适合问题的单位和规模的容差；一个固定的 epsilon 并不适合每一种计算。

## 比较bool

`<`、`<=`、`>`、`>=`、`==`、`!=` 检查关系。一个等号`=`是赋值，两个等号`==`是相等比较。

<!-- wave-example: book-comparisons -->
```wave
fun main() {
    var age: i32 = 20;
    var minimum: i32 = 18;
    var eligible: bool = age >= minimum;

    if (eligible) {
        println("eligible");
    }

    if (age != minimum) {
        println("not exactly the boundary");
    }
}
```

执行结果：

```text
eligible
not exactly the boundary
```

“至少18”包括18； “大于18”排除它。测试 17、18 和 19 以检查此边界。混合类型比较取决于符号性和宽度，因此将两个操作数转换为预期类型可以使比较更清晰。

## 逻辑运算和短路评估

`&&` 检查两者是否都为真，`||` 检查是否有多个为真，`!` 翻转真/假。此操作并不总是执行右侧的表达式。

<!-- wave-example: book-short-circuit -->
```wave
fun report() -> bool {
    println("right side evaluated");
    return true;
}

fun main() {
    var enabled: bool = false;

    if (enabled && report()) {
        println("both true");
    }

    if (!enabled || report()) {
        println("at least one true");
    }
}
```

执行结果：

```text
at least one true
```

在第一个条件中，enabled为假，因此无需向右查找。在第二个中，!enabled 为 true，因此右侧也不需要。因此，report的输出永远不会出现。

您可以使用它来检查除法之前的分母。函数体片段`if (divisor != 0 && value / divisor > 2) { ... }`在分母为0时不执行除法。但是，仅此测试无法解析其他边界，例如signed最小值/-1。

`&&` 优先于`||`。在复杂的策略中，请在括号中表明您的意图，例如`(member && active) || admin`。

## 缩小或扩大整数

`as` 是显式强制转换。缩小整数时被丢弃的高位无法通过再次扩大整数来恢复。

<!-- wave-example: book-cast-chain -->
```wave
fun main() {
    var original: i32 = 300;
    var small: u8 = original as u8;
    var widened: i32 = small as i32;

    println("{} -> {} -> {}", original, small, widened);

    var negative: i8 = -1;
    var signed_value: i32 = negative as i32;
    var same_bits: u8 = negative as u8;

    println("{} {}", signed_value, same_bits);
}
```

执行结果：

```text
300 -> 44 -> 44
-1 255
```

300 的低 8 位是 44。如果将 -1 加宽为signed 类型，则符号将扩展为维持 -1。如果将相同的 8 位解释为unsigned，则为 255。

寻求保留范围内值的转换和尝试操纵存储位的转换具有不同的目的。如果您有用户输入，请首先检查它是否是目标范围并进行转换。 cast的存在并不能保证该值处于安全范围内。

## 替换为bool

将整数转换为 bool 会产生 false（零），否则产生 true。这不会截断到最低位：2 也会转换为 true。对于浮点值，仅+0.0和-0.0转换为false；所有其他值，包括 NaN 和无穷大，都转换为 true。

<!-- wave-example: book-bool-conversion -->
```wave
fun main() {
    var zero: i32 = 0;
    var two: i32 = 2;

    if (!(zero as bool)) {
        println("zero is false");
    }

    if (two as bool) {
        println("two is true");
    }
}
```

执行结果：

```text
zero is false
two is true
```

不支持将指针转换为 bool。显式地将指针与 null 进行比较，例如与 `pointer != null` 进行比较。非null地址是否可以安全读取是一个单独的问题。

## 位运算

`&`、`|`、`^`、`~` 覆盖整数的每一位。它可用于以位表示权限或功能。下面，1是读权限，2是写权限。

<!-- wave-example: book-bit-flags -->
```wave
const READ: u32 = 1;
const WRITE: u32 = 2;

fun main() {
    var permissions: u32 = READ | WRITE;

    if ((permissions & WRITE) != 0) {
        println("write enabled");
    }

    permissions = permissions & ~WRITE;
    println("remaining={}", permissions);
}
```

执行结果：

```text
write enabled
remaining=1
```

使用 OR 添加一位，并使用 AND 检查是否存在特定位。使用 `~WRITE` 创建一个仅相关位为 0 的掩码并将其擦除。按位运算`&`·`|` 与bool 的短路运算`&&`·`||` 是不同的运算符。

## 宽度和班次数量

左移将位向左移动并丢弃操作数宽度之外的高位。右移对有符号值使用符号扩展，对无符号值使用零扩展。结果始终具有左操作数的类型。

<!-- wave-example: book-shifts -->
```wave
fun main() {
    var bits: u8 = 129;
    var shifted: u8 = bits << 1;
    var negative: i32 = -8;
    var positive: u32 = 8;

    println("left={}", shifted);
    println("right={} {}", negative >> 1, positive >> 1);
}
```

执行结果：

```text
left=2
right=-4 4
```

将 u8 值 129 左移一位会丢弃其最高位并留下 2。原始移位计数值必须为非负且小于左操作数的位宽。对于u8，有效计数为 0 到 7。无效的常量计数是编译时错误；无效的运行时间计数会导致陷阱。

## 将浮点值转换为整数

浮点到整数的转换首先朝零截断，然后检查目标整数范围。 NaN、无穷大和超出范围的结果无效。无效的常量转换会产生编译时错误；无效的运行时转换会导致陷阱。

陷阱不会从函数返回错误值。对于可恢复的转换失败，设计一个在转换之前检查范围的接口。要获取浮点值的存储位，请使用`std::math::float`中的位转换函数而不是数字转换。

## 练习与解答

编写一个函数，将 137 韩元的找零分成 50 韩元和余数，并仅更改 0 到 255 到 u8 范围内的整数。此示例将验证步骤演示为一个单独的函数，该函数将超出范围指示为 -1。

<!-- wave-example: book-expression-solution -->
```wave
fun checked_byte(value: i32) -> i32 {
    if (value < 0 || value > 255) {
        return -1;
    }

    var small: u8 = value as u8;
    return small as i32;
}

fun main() {
    var money: i32 = 137;

    println("coins={} remainder={}", money / 50, money % 50);
    println("{} {} {}", checked_byte(0), checked_byte(255), checked_byte(256));
}
```

执行结果：

```text
coins=2 remainder=37
0 255 -1
```

之所以可以使用-1作为失败指示符，是因为成功的范围是0到255。如果任何整数都可以作为成功值，则需要不同的结果表示。我们在后面的错误处理章节中继续这个设计。


## 优先级

运算符优先级如下，从最高优先级开始：

1. 基本表达式和后缀访问：函数调用、字段访问、索引、后缀`++`·`--`
2. 一元运算：`!`、`~`、`&`、`deref`，前缀 `++`·`--`，一元 `+`·`-`
3. `as` 类型转换
4. `*`, `/`, `%`
5. `+`, `-`
6. `<<`, `>>`
7. `<`, `<=`, `>`, `>=`
8. `==`, `!=`
9. 位`&`
10. 位`^`
11. 位`|`
12. `&&`
13. `||`
14. 赋值和复合赋值

链分配从右侧组合。当混合不同类型的运算符时，使用括号来清楚地表达求值的顺序。

## 可替代的科目

赋值、`++`和`--`需要一个表示存储位置的表达式，例如变量、字段、数组元素或取消引用的指针。不允许写入`const`。

## 转变

移位的结果始终具有左操作数的类型；正确的操作数类型不会扩大计算范围。左移丢弃超出该宽度的高位。右移符号扩展有符号值和零扩展无符号值。

移位计数必须是一个整数，其原始值满足`0 <= n < LHS bit width`。在截断为较小类型之前会对其进行检查。无效的常量计数是一个编译时错误；无效的运行时间计数会导致陷阱。

## 涉及布尔值和浮点值的转换

整数到bool 的转换会生成 false（对于零）和 true（对于每个其他值）。浮点到bool的转换仅在+0.0和-0.0时产生false； NaN 和正无穷大或负无穷大产生 true。不支持将指针强制转换为 bool：与 `pointer != null` 显式比较。

浮点到整数的转换朝零方向截断，然后检查目标范围。 NaN、无穷大和超出范围的结果无效。无效的常量转换是编译时错误；无效的运行时转换会导致陷阱。陷阱不是可恢复的错误返回。

`&&` 和 `||` 是短路评估。未执行的右操作数的副作用不会发生。您可以在[算术课](/docs/zh/language/expressions-and-operators)中检查较小值的结果。

## 故意失败的转变

8 位值的移位数必须为 0 到 7。下面的程序在执行前应被拒绝。

<!-- wave-example: reject-shift-count -->
```wave
fun main() {
    var value: u8 = 1;
    var result: u8 = value << 8;
}
```
