---
translation_set_id: functions
path: language/functions-and-generics
locale: zh
group: language
group_order: 2
order: 5
title: 5. 设计和编写函数
summary: 了解参数、返回值、默认值和按值传递。
---

## 从重复的代码开始

函数是减少语法的工具，但它们也是划分任务的工具。将其作为输入的内容、计算的内容以及返回的结果分开，可以让您以较小的块来理解您的程序。

在本章中，我们从一个多次编写折扣计算的程序开始。每个示例都是整个 main.wave 并使用“wavec run main.wave”运行。我们还没有分割文件。

<!-- wave-example: book-function-before -->
```wave
fun main() {
    var first_price: i32 = 2000;
    var first_discount: i32 = first_price * 10 / 100;
    var first_total: i32 = first_price - first_discount;
    var second_price: i32 = 5000;
    var second_discount: i32 = second_price * 10 / 100;
    var second_total: i32 = second_price - second_discount;
    println("{} {}", first_total, second_total);
}
```

执行结果：

```text
1800 4500
```

这两种计算仅在价格上有所不同，并且具有相同的结构。更改折扣规则时，您需要编辑这两个位置。如果仅更改一侧，对于需要相同策略的产品，您将获得不同的结果。

## 确定输入和输出

将多余的计算移至函数中。更改后的值作为输入price接收，并返回计算出的价格。

<!-- wave-example: book-function-extract -->
```wave
fun discounted(price: i32) -> i32 {
    var discount: i32 = price * 10 / 100;
    return price - discount;
}

fun main() {
    println("{} {}", discounted(2000), discounted(5000));
}
```

执行结果：

```text
1800 4500
```

函数名称为discounted。括号中的`price: i32`是参数声明，`-> i32`是结果类型。主体中的局部变量discount仅在该函数内使用。

`discounted(2000)` 是调用函数的表达式。括号中的 2000 是您实际传递的参数。由于函数返回的值成为该调用表达式的结果，因此可以直接用作println的参数。

|术语|代码|意义|
| --- | --- | --- |
|参数| price |声明函数时指定的输入名称|
|因素| 2000 |调用时传递的值|
|返回类型| i32 |调用表达式生成的值的类型|
|返回声明| return price - discount |传递结果并结束本次调用|

## 遵循调用和执行顺序

仅在源代码中编写函数声明不会导致其主体立即执行。当到达 main 中调用的点时执行。

<!-- wave-example: book-function-trace -->
```wave
fun calculate(value: i32) -> i32 {
    println("inside: {}", value);
    return value * 2;
}

fun main() {
    println("before");
    var result: i32 = calculate(7);
    println("after: {}", result);
}
```

执行结果：

```text
before
inside: 7
after: 14
```

进度顺序是第一个输出main、calculate正文、最后一个输出main。当执行`return`时，对calculate的调用结束，结果为14，main的result初始化完成。

如果多个函数调用混合在一个表达式中，并且副作用的顺序很重要，请将调用分解为单独的语句。本章中的示例还将需要跟踪的调用结果存储在局部变量中。

## 多个参数

如果您还收到折扣率作为输入，则可以使用相同的函数计算多个保单。

<!-- wave-example: book-function-parameters -->
```wave
fun discounted(price: i32, percent: i32) -> i32 {
    return price - price * percent / 100;
}

fun main() {
    var standard: i32 = discounted(2000, 10);
    var special: i32 = discounted(2000, 25);
    println("standard={}", standard);
    println("special={}", special);
}
```

执行结果：

```text
standard=1800
special=1500
```

参数的顺序必须与声明匹配。如果两个参数都是i32，即使改变顺序，仅通过类型检查也很难区分它们的含义。明确定义函数名和参数名，并以易于阅读的方式写出调用位置。

此函数假设金额较小且具有有效的折扣率。不处理负价格、大于 100 的比率或中间乘法溢出。创建函数时，不仅需要描述函数体，还需要描述输入条件。后期完成程序将检查步骤分开。

## 默认参数

您可以提供常用值作为默认值。

<!-- wave-example: book-function-default -->
```wave
fun discounted(price: i32, percent: i32 = 10) -> i32 {
    return price - price * percent / 100;
}

fun main() {
    println("{}", discounted(2000));
    println("{}", discounted(2000, 25));
}
```

执行结果：

```text
1800
1500
```

第一次调用省略第二个参数并使用 10。第二次调用使用指定的 25。最后为可选参数保留默认值。它不被用作编写空格以仅省略第一个参数的语法。

更改默认值会更改跳过调用的行为。公共函数的默认值也是你所依赖的行为的一部分。这就是为什么带有指定参数的调用和带有省略参数的调用是分开测试的。

## 意味着按值传递

传递整数值将函数接收的值与调用者的变量存储区分开。计算结果不会自动更改调用者的变量。

<!-- wave-example: book-function-value -->
```wave
fun next(value: i32) -> i32 {
    return value + 1;
}

fun main() {
    var count: i32 = 4;
    var later: i32 = next(count);
    println("count={} later={}", count, later);
    count = next(count);
    println("count={}", count);
}
```

执行结果：

```text
count=4 later=5
count=5
```

第一次调用读取count并初始化later。 count仍然是4。第二次调用后，结果被赋值给count，所以变成了5。按值返回的设计让调用者很明显数据发生了变化。

如果想改变函数内原有的存储空间，可以传递指针。 [指针章节](/docs/zh/language/explicit-memory-type-model) 对此进行了介绍。即使传递指针，也必须区分指针值本身和其地址的存储空间。

## 不走不归路

返回值的函数必须提供所有所需路径的结果。它不会留下这样的最终路径：

```wave
// 의도적으로 잘못된 함수 조각

fun positive(value: i32) -> i32 {
    if (value > 0) {
        return value;
    }
}
```

如果value小于或等于0，则没有设置值可返回。一旦你建立了你想要的规则，你需要写出你的所有路线。

<!-- wave-example: book-function-paths -->
```wave
fun positive_or_zero(value: i32) -> i32 {
    if (value > 0) {
        return value;
    }

    return 0;
}

fun main() {
    println("{} {} {}", positive_or_zero(-2), positive_or_zero(0), positive_or_zero(8));
}
```

执行结果：

```text
0 0 8
```

首先执行return的调用不会继续执行下面的return。仅当条件为假时才达到最后一个return。我们检查了三种情况的规则：正、0 和负。

## 函数没有结果

如果只执行打印等操作，则可以省略返回类型。没有结果的函数也可以通过`return;`提前终止。

<!-- wave-example: book-function-void -->
```wave
fun show_positive(value: i32) {
    if (value <= 0) {
        return;
    }

    println("positive={}", value);
}

fun main() {
    show_positive(-3);
    show_positive(6);
}
```

执行结果：

```text
positive=6
```

第一个调用返回，但不打印任何内容。第二次调用打印：“无结果”和“不返回调用点”是不同的。不返回的函数，例如终止进程的函数，按其返回类型`!`进行分类。

## 编写一个具有多种功能的完整程序

现在我们将输入验证、计算和输出拆分为不同的函数。由于我们限制了价格范围，因此本例中的中间乘法在i32范围内。

<!-- wave-example: book-function-program -->
```wave
fun valid_order(price: i32, quantity: i32, percent: i32) -> bool {
    return price >= 0 && price <= 100000
        && quantity >= 1 && quantity <= 100
        && percent >= 0 && percent <= 100;
}

fun discounted_unit(price: i32, percent: i32) -> i32 {
    return price - price * percent / 100;
}

fun order_total(price: i32, quantity: i32, percent: i32) -> i32 {
    var unit: i32 = discounted_unit(price, percent);
    return unit * quantity;
}

fun show_order(price: i32, quantity: i32, percent: i32) {
    if (!valid_order(price, quantity, percent)) {
        println("invalid order");
        return;
    }

    var total: i32 = order_total(price, quantity, percent);
    println("total={}", total);
}

fun main() {
    show_order(2000, 3, 10);
    show_order(2000, 0, 10);
    show_order(2000, 3, 120);
}
```

执行结果：

```text
total=5400
invalid order
invalid order
```

为每个功能回答一个问题。 valid_order决定是否允许输入，discounted_unit决定一次折扣是多少，order_total决定总价是多少，show_order决定显示什么。

功能越小并不一定越好。为每个表达式命名可能会导致难以跟踪。当规则对于在其他地方重用有意义或有需要独立解释和验证的规则时，将规则分开。

## 练习题

1. 写入 maximum，它返回两个整数中较大的一个。
2. 编写一个函数 clamp 返回两个边界之间的值。在下面的解决方案中，low <= high 被设置为调用条件。
3. 创建一个对值加税的函数，并将其与订单总和函数结合起来。首先确定范围和整数切割点。

### 解：有边界的函数

<!-- wave-example: book-function-solution -->
```wave
fun maximum(left: i32, right: i32) -> i32 {
    if (left > right) {
        return left;
    }

    return right;
}

fun clamp(value: i32, low: i32, high: i32) -> i32 {
    if (value < low) {
        return low;
    }
    if (value > high) {
        return high;
    }

    return value;
}

fun main() {
    println("max={}", maximum(7, 4));
    println("{} {} {}", clamp(-3, 0, 10), clamp(6, 0, 10), clamp(20, 0, 10));
}
```

执行结果：

```text
max=7
0 6 10
```

clamp 有三种路径：小于范围、在范围内、大于范围。还可以通过手动添加边界值 0 和 10 进行检查。如果你想申请到low > high，你必须决定如何表达失败。 [错误处理章节](/docs/zh/language/errors) 使用结果结构解决此问题。


## 递归调用

函数可以调用自身。递归需要一个不再进行调用的退出条件，以及一个每个调用都更接近该条件的过程。

<!-- wave-example: book-ref-recursion -->
```wave
fun factorial(value: i32) -> i32 {
    if (value <= 1) {
        return 1;
    }

    return value * factorial(value - 1);
}

fun main() {
    println("{}", factorial(5));
}
```

执行结果：

```text
120
```

5 的计算导致`5 * factorial(4)`，达到时返回 1。返回值按顺序传递给前一个调用，结果为 120。此函数是解释小正整数的示例。对于大量输入，您需要考虑结果范围和调用深度。您可以通过在循环中编写相同的操作来避免增加调用深度的问题。

## 常见错误

|现象|检查|
| --- | --- |
|参数不足或太多的错误|参数数量和可选默认值|
|返回类型不匹配|return 表达式类型和函数声明|
|不从特定路径返回|即使条件为假它也会返回吗？|
|缺少泛型类型参数|函数名后`<Type>`|
|调用指针函数后，原来的内容发生了变化。|它是一个只读取值的函数还是一个修改它的函数？|

外部导出的函数，例如`export(c)`，使用特定的签名。通用函数本身无法通过外部调用约定导出。 `ptr<T>` 和 `array<T, N>` 是语言的内置内存类型，将它们与用户通用结构声明区分开来。

[功能学习](/docs/zh/language/functions-and-generics) · [学习模块和泛型](/docs/zh/language/modules-imports-and-ffi) · [FFI](/docs/zh/language/modules-imports-and-ffi)
