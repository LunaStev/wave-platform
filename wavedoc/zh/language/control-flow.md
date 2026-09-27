---
translation_set_id: control-flow
path: language/control-flow
locale: zh
group: language
group_order: 2
order: 4
title: 4. 条件、循环和边界值
summary: 了解if、for、while以及重复语句的范围。
---

## 选择执行路径

上一章中的程序从上到下执行语句。真正的程序必须根据输入和状态执行不同的操作。条件语句选择执行路径，循环将相同的规则应用于多个值。

将每个示例保存在main.wave中并运行它。阅读代码时，在纸上写下当前变量值、接下来要测试的条件以及要执行的语句的顺序。跟随流程进行练习比记住结果更重要。

## if 和 else

条件写在括号中，文本括在大括号中。在下面的示例中，将balance替换为500或2000来确定执行哪个分支。

<!-- wave-example: book-if-else -->
```wave
fun main() {
    var balance: i32 = 2000;
    var price: i32 = 1200;

    if (balance >= price) {
        balance -= price;
        println("bought, balance={}", balance);
    } else {
        println("not enough money");
    }
}
```

执行结果：

```text
bought, balance=800
```

它不执行这两个块。如果条件为真，则执行第一个块。如果条件为假，则执行else块。即使balance等于price也允许购买，所以我写了`>=`。如果将其更改为`>`，则操作将更改相同的金额。

## 多个条件的顺序

您可以将条件与 else if 连接起来。由于我们首先从上面执行一个满足的分支，因此我们需要考虑是否要先检查较大的边界还是较小的边界。

<!-- wave-example: book-grade -->
```wave
fun grade(score: i32) {
    if (score < 0 || score > 100) {
        println("invalid");
    } else if (score >= 90) {
        println("A");
    } else if (score >= 80) {
        println("B");
    } else {
        println("C");
    }
}

fun main() {
    grade(95);
    grade(80);
    grade(79);
    grade(101);
}
```

执行结果：

```text
A
B
C
invalid
```

无效分数首先被拒绝，然后进行评分。如果您将 `score >= 80` 放在开头，您将永远不会到达 A 分支，因为 95 度进入该分支。检查条件的顺序以及每个条件的正确性。

## 不要更改条件表达式中的值

if、while 或 for 条件中不允许进行赋值、复合赋值以及递增或递减运算。使用`==`进行比较。要更新值然后测试它，请编写两个单独的语句。

函数内片段的正确形式：

```wave
value = read_value();

if (value == expected) {
    println("matched");
}
```

read_value 和 expected 未在此片段中定义，因此它不是按原样运行的完整程序。这里显示的规则是“状态更改后比较”。

## while: 当条件成立时

<!-- wave-example: book-while-countdown -->
```wave
fun main() {
    var remaining: i32 = 3;

    while (remaining > 0) {
        println("{}", remaining);
        remaining -= 1;
    }

    println("finished at {}", remaining);
}
```

执行结果：

```text
3
2
1
finished at 0
```

在进入身体之前检查条件。如果 remaining 从一开始就是 0，则主体永远不会被执行。如果省略主体末尾的减量，则条件保持为真并且循环不会结束。

编写循环后，检查“是什么使它更接近终止条件？”如果是等待输入的循环，则输入变化或EOF发挥作用，如果是数字循环，则索引更新发挥作用。

## for: 初始化/条件/更新

for 表示重复所必需的三个部分。

<!-- wave-example: book-for-sum -->
```wave
fun main() {
    var total: i32 = 0;

    for (var number: i32 = 1; number <= 5; number += 1) {
        total += number;
    }

    println("sum={}", total);
}
```

执行结果：

```text
sum=15
```

1. 将number初始化为1。这一步是一次。
2. 检查 number <= 5。如果为 false，则结束迭代。
3. 在文本中将number添加到total。
4. 将 number 加 1 并返回条件检查。

不假定在 for 中声明的重复变量可以在重复后使用。如果您的设计在迭代后需要一个值，请在外部声明它并明确初始化位置。

## 包含和排除边界

1 到 5 的自然和是 `<= 5`。另一方面，长度为 5 的数组的索引应该使用 `< 5`。这是因为数组索引从 0 开始，到 4 结束。

不要将“运行五次”与“最多且包括值 5”混淆。您可以通过一起查看开始值和结束值来确定重复次数。输入为空且只有一个元素的情况有利于检测边界错误。

## continue 和 break

continue 跳过本次迭代的其余部分，break 结束最近的迭代。

<!-- wave-example: book-loop-control -->
```wave
fun main() {
    var total: i32 = 0;

    for (var number: i32 = 1; number <= 10; number += 1) {
        if (number == 3) {
            continue;
        }

        if (number == 6) {
            break;
        }

        total += number;
    }

    println("{}", total);
}
```

执行结果：

```text
12
```

总数为 1、2、4、5。如果for 中遇到continue，则继续续订过程。由于while没有像for那样的单独更新表达式，因此您必须小心不要遗漏continue之前所需的任何状态更改。

如果循环嵌套，则一个 break 无法完成所有循环。如果您需要在多个阶段停止，请将工作封装在一个函数中，并通过使用 return 或检查外部迭代中的终止条件来表明您的意图。

## 将大小写除以 match

比较多个相同值的情况时，可以使用match。每个arm的主体是一个块。

<!-- wave-example: book-match-number -->
```wave
fun describe(status: i32) {
    match (status) {
        200 => {
            println("ok");
        }
        404 => {
            println("missing");
        }
        _ => {
            println("other");
        }
    }
}

fun main() {
    describe(200);
    describe(404);
    describe(500);
}
```

执行结果：

```text
ok
missing
other
```

`_` 是处理其余部分的模式。请勿在同一个 match 中放置重复项。 variant，根据值的不同，具有不同的数据类型，在[数据模型章节](/docs/zh/language/structures-enums-and-aliases)中进行了介绍。

## 完整示例：计算满足条件的数字

求 1 到 10 中偶数的个数和和。由于 count 和 sum 是不同的信息，因此它们作为变量进行累加。

<!-- wave-example: book-loop-statistics -->
```wave
fun main() {
    var count: i32 = 0;
    var total: i32 = 0;

    for (var number: i32 = 1; number <= 10; number += 1) {
        if (number % 2 == 0) {
            count += 1;
            total += number;
        }
    }

    println("count={} total={}", count, total);
}
```

执行结果：

```text
count=5 total=30
```

偶数为 2、4、6、8 和 10，因此数字为 5，总和为 30。即使表达式很短，如果您首先检查可以手动获得的小范围内的结果，也很容易验证迭代边界。

## 练习与完整解答

只添加1到20之间3的倍数，但不要添加任何加起来超过30的值。我们需要区分“添加然后检查是否结束”和“检查是否结束然后添加”。

<!-- wave-example: book-loop-solution -->
```wave
fun main() {
    var total: i32 = 0;

    for (var number: i32 = 1; number <= 20; number += 1) {
        if (number % 3 != 0) {
            continue;
        }

        if (total + number > 30) {
            break;
        }

        total += number;
    }

    println("{}", total);
}
```

执行结果：

```text
30
```

值 3+6+9+12 之和为 30，因此不添加下一个值 15。这些小输入是安全的，但对于大整数，检查`total + number`本身可能会溢出。编写支票不会自动处理每种边界情况。
