---
translation_set_id: learn-arrays-strings
path: language/arrays
locale: zh
group: language
group_order: 2
order: 6
title: 6. 数组和迭代
summary: 学习固定大小数组、索引、遍历、复制和搜索。
---

## 同一类型的多个值

如果将三个分数分别创建为score1、score2和score3，则当数字发生变化时，声明和计算都必须更改。数组将一定数量的相同类型的元素组合在一起。使用循环，您可以将相同的规则应用于每个元素。

本章涵盖创建、索引、修改、迭代、搜索和聚合数组。字符串也支持索引，但其含义有所不同，因此下一章将介绍字符串。

## 在类型中输入长度

<!-- wave-example: book-array-create -->
```wave
fun main() {
    var scores: array<i32, 3> = [70, 80, 90];

    println("first={}", scores[0]);
    println("second={}", scores[1]);
    println("last={}", scores[2]);
}
```

执行结果：

```text
first=70
second=80
last=90
```

`array<i32, 3>`中的i32是元素类型，3是元素数量。我们存储的是 3 个整数。这并不意味着字节数为 3。数组文字中的元素数必须与声明的长度匹配。

索引从0开始。第一个元素是0，最后一个元素是length-1。 scores[3] 是超出范围的访问，而不是第三个元素。

## 改变元素

<!-- wave-example: book-array-update -->
```wave
fun main() {
    var scores: array<i32, 3> = [70, 80, 90];

    scores[1] = 85;
    scores[0] += 5;

    println("{} {} {}", scores[0], scores[1], scores[2]);
}
```

执行结果：

```text
75 85 90
```

更改特定元素的存储，而无需重新创建整个数组。索引表达式也可以是计算的结果，但必须确保该值在某个范围内。当使用外部输入作为索引时，将检查负数和上限。

## 迭代数组

<!-- wave-example: book-array-sum -->
```wave
fun main() {
    var scores: array<i32, 4> = [60, 70, 80, 90];
    var total: i32 = 0;

    for (var index: i32 = 0; index < 4; index += 1) {
        total += scores[index];
    }

    println("total={} average={}", total, total / 4);
}
```

执行结果：

```text
total=300 average=75
```

每次迭代都会读取不同索引处的元素。 sum 变量必须在迭代之外初始化。每次在循环体内将其初始化为 0 会产生不正确的结果，例如只保留最后一个元素。

用于计算平均值的整数除法会丢弃小数部分。对于浮点平均值，请在除法之前先转换总和。对于较大的数组或较大的值，还要确保累加器类型可以表示总和。

## 仅聚合部分元素

可以通过条件语句和遍历相结合来实现过滤。这里我们统计得分为 80 或更高的元素数量。

<!-- wave-example: book-array-filter -->
```wave
fun main() {
    var scores: array<i32, 5> = [60, 80, 90, 75, 100];
    var passed: i32 = 0;

    for (var index: i32 = 0; index < 5; index += 1) {
        if (scores[index] >= 80) {
            passed += 1;
        }
    }

    println("passed={}", passed);
}
```

执行结果：

```text
passed=3
```

索引和元素值必须分开。检查`index >= 80`比较位置，而不是分数。两者都可以是i32，所以单独根据类型很难发现这个语义错误。

## 查找第一个匹配位置

首先，决定如何显示未找到的结果。在此示例中，有效索引为 0 到 4，因此我们使用 -1 作为失败标记。

<!-- wave-example: book-array-find -->
```wave
fun main() {
    var values: array<i32, 5> = [8, 3, 8, 1, 5];
    var target: i32 = 8;
    var found: i32 = -1;

    for (var index: i32 = 0; index < 5; index += 1) {
        if (values[index] == target) {
            found = index;
            break;
        }
    }

    if (found >= 0) {
        println("found at {}", found);
    } else {
        println("not found");
    }
}
```

执行结果：

```text
found at 0
```

第一个位置0也是正常结果。如果你用`found > 0`检查是否成功，你会被误认为没有找到第一个元素。同样的道理，用bool作为cast来判断成功也是错误的。

如果删除 break，后续匹配将覆盖 found，为您提供最后一场匹配的位置。由于单个语句可以改变函数的契约，因此“搜索”的描述也必须具体写明它是在第一个位置还是最后一个位置。

## 复制数组元素

要将数组的值复制到另一个存储空间，可以逐个元素读取并分配它们。即使复制后更改一个整数元素，另一整数元素也不会更改。

<!-- wave-example: book-array-copy -->
```wave
fun main() {
    var original: array<i32, 3> = [1, 2, 3];
    var copied: array<i32, 3>;

    for (var index: i32 = 0; index < 3; index += 1) {
        copied[index] = original[index];
    }

    copied[0] = 99;

    println("original={}", original[0]);
    println("copied={}", copied[0]);
}
```

执行结果：

```text
original=1
copied=99
```

如果元素是指针，则复制它们会复制它们的地址。它不会复制它们指向的单独内存。在管理所有权时，这种区别很重要。

## 初始化及有效范围

它并不假设声明为无初始值的数组的所有元素都可以读取。如果只记录了部分号码，则必须单独管理实际初始化的号码。这与库读取函数返回的长度可能小于整个缓冲区容量的原因相同。

由于数组长度包含在类型中，因此在执行过程中不会任意增长。大小增长的字节列表使用动态存储，例如`Buffer`。更改数组的长度需要考虑类型、初始值、遍历上限以及依赖于该长度的计算。

## 练习：最大值和位置

找出最大值及其在数组`[4, 9, 2, 9, 1]`中首次出现的位置。如果它是所有元素都可能为负的常规函数，则最大值不应初始化为 0。

### 完整解答

<!-- wave-example: book-array-solution -->
```wave
fun main() {
    var values: array<i32, 5> = [4, 9, 2, 9, 1];
    var maximum: i32 = values[0];
    var position: i32 = 0;

    for (var index: i32 = 1; index < 5; index += 1) {
        if (values[index] > maximum) {
            maximum = values[index];
            position = index;
        }
    }

    println("max={} first={}", maximum, position);
}
```

执行结果：

```text
max=9 first=1
```

将第一个元素作为初始参考并与第二个元素进行比较。从`>`开始，即使再次出现相同的最大值，位置也不会改变。将其更改为`>=`成为最后一个位置。任何长度为零的接口都必须在读取第一个元素之前处理空输入。
