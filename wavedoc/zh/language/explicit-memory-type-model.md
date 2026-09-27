---
translation_set_id: memory-model
path: language/explicit-memory-type-model
locale: zh
group: language
group_order: 2
order: 9
title: 9. 指针、修改值和生命周期
summary: 学习寻址、取消引用、通过指针修改值以及悬空指针。
---

## 区分价值和存储位置

整数42和该整数存储的地址是不同的值。指针指向存储位置。传递地址允许函数读取或更改调用者的存储。

本章介绍了获取地址、取消引用、修改原始值、指针算术和生存期。动态分配将在下一章介绍。从局部变量和数组元素的地址开始。

## 获取地址并取消引用

<!-- wave-example: book-pointer-first -->
```wave
fun main() {
    var value: i32 = 42;
    var address: ptr<i32> = &value;

    println("value={}", value);
    println("through pointer={}", deref address);

    deref address = 99;
    println("changed={}", value);
}
```

执行结果：

```text
value=42
through pointer=42
changed=99
```

`&value` 获取地址，`deref address` 读取或写入该地址的值，依此类推。 99 不是存储在address 中，而是写入address 指向的整数。变量address本身继续指向value。

`ptr<i32>`是用于访问i32存储的指针类型。该类型不记录长度或提供自动释放。

## 改变指针本身

<!-- wave-example: book-pointer-reassign -->
```wave
fun main() {
    var first: i32 = 10;
    var second: i32 = 20;
    var selected: ptr<i32> = &first;

    selected = &second;
    deref selected = 25;

    println("{} {}", first, second);
}
```

执行结果：

```text
10 25
```

`selected = &second` 将另一个地址存储在指针变量中。 first 的值不会改变。然后，如果将值写入 deref，second 会发生变化。将“地址更改”和“通过地址更改值”分成单独的句子将减少混乱。

## 让一个函数改变原来的函数

<!-- wave-example: book-pointer-increment -->
```wave
fun increment(value: ptr<i32>) {
    deref value = deref value + 1;
}

fun main() {
    var count: i32 = 4;

    increment(&count);
    increment(&count);

    println("{}", count);
}
```

执行结果：

```text
6
```

该函数接收 count 的地址，而不是其值 4。更改该存储也会更改调用者中的 count。该函数需要一个有效的、可写的i32地址。通过null违反了此要求。

每个函数定义是否接受null。如果没有，呼叫者必须提供有效的地址。如果是，该函数必须包含处理 null 的路径。

## 处理函数null

<!-- wave-example: book-pointer-null -->
```wave
fun try_increment(value: ptr<i32>) -> bool {
    if (value == null) {
        return false;
    }

    deref value = deref value + 1;
    return true;
}

fun main() {
    var count: i32 = 7;

    if (!try_increment(null)) {
        println("no value");
    }

    if (try_increment(&count)) {
        println("count={}", count);
    }
}
```

执行结果：

```text
no value
count=8
```

null 检查仅处理地址缺失的情况。将任意非null数字转换为指针不会创建有效的内存。读取和写入还需要有效的生存期、足够的大小、正确的对齐方式以及适当的访问权限。

## 数组地址和逐元素指针算术

<!-- wave-example: book-pointer-array -->
```wave
fun main() {
    var values: array<i32, 3> = [10, 20, 30];
    var first: ptr<i32> = &values[0];
    var second: ptr<i32> = first + 1;

    println("{}", deref first);
    println("{}", deref second);
    println("{}", deref first[2]);
}
```

执行结果：

```text
10
20
30
```

向指针加 1 将其移动到目标类型的一个元素。 i32 中的下一个元素和 u8 中的下一个元素具有不同的移位字节数。如果您再次将`first + 1`乘以字体大小并相加，它将移动到不需要的位置。

指针索引也必须在有效范围内完成。 first 不记得数组长度 3 本身，因此在将范围传递给函数时，它使用同时接收指针和长度的形式。

## 将读取范围传递给函数

<!-- wave-example: book-pointer-range -->
```wave
fun sum(values: ptr<i32>, count: i32) -> i32 {
    var total: i32 = 0;

    for (var index: i32 = 0; index < count; index += 1) {
        total += deref values[index];
    }

    return total;
}

fun main() {
    var values: array<i32, 4> = [2, 4, 6, 8];

    println("first two={}", sum(&values[0], 2));
    println("all={}", sum(&values[0], 4));
}
```

执行结果：

```text
first two=6
all=20
```

count的单位是元素个数。此函数要求调用者具有 count 可读 i32。传递大于实际数组的长度违反了约定。即使您使用相同的 ptr<u8>·i64 组合 API，您也应该在文档中检查长度是以字节为单位还是以元素为单位。

## 寿命：地址的有效期是多长？

局部变量在调用和块的生命周期内使用。如果您返回函数内局部变量的地址供调用者稍后读取，则该存储空间可能已达到其生命周期的终点。

这是一个不应该实现的糟糕设计：

```wave
fun invalid_address() -> ptr<i32> {
    var local: i32 = 42;

    return &local;
}
```

如果需要一个值，则返回i32。如果它需要写入调用者提供的存储空间，它需要一个指针作为输入。如果您需要单独的存储来保留调用之外的内容，请显式分配它并传递释放它的责任。

## 两个指针指向同一个存储空间

复制指针会创建指向同一地址的另一个名称。不复制内存。

<!-- wave-example: book-pointer-alias -->
```wave
fun main() {
    var value: i32 = 1;
    var first: ptr<i32> = &value;
    var second: ptr<i32> = first;

    deref second = 9;

    println("{} {}", value, deref first);
}
```

执行结果：

```text
9 9
```

通过second改变的结果也可以通过first看到。一旦源内存被释放，两个指针就变得不可用。将 null 分配给一个指针变量不会自动更改其他副本。

## 练习：交换两个整数

编写一个函数，获取两个 i32 地址并交换它们的值。第一个值在覆盖之前必须存储在临时变量中。检查即使两次传递同一个地址该值是否保留。

### 完整解答

<!-- wave-example: book-pointer-swap -->
```wave
fun swap(left: ptr<i32>, right: ptr<i32>) {
    var saved: i32 = deref left;

    deref left = deref right;
    deref right = saved;
}

fun main() {
    var first: i32 = 3;
    var second: i32 = 8;

    swap(&first, &second);
    println("{} {}", first, second);

    swap(&first, &first);
    println("{}", first);
}
```

执行结果：

```text
8 3
8
```

该函数还要求两个地址都指向有效的、可写的整数存储。要处理null，请添加一个指示成功或失败的结果，如 try_increment 中。


## **Wave Explicit Memory Type Model**

Wave的指针设计是基于**Wave Explicit Memory Type Model**。该模型将指针和数组定义为语言级别的显式内存类型，而不是语法技巧或库抽象。

`ptr<T>`是指向存储`T`值的内存地址的类型，`array<T, N>`是定长内存类型，连续存储`T`的`N`值。因此，指针和数组的结构按原样在函数参数、返回值、结构字段和其他类型中显示。

## null

```wave
var buffer: ptr<u8> = null;
if (buffer == null) {
    println("no buffer");
}
```

`null`是一个未指向有效内存地址的指针值。 `null` 只能赋值给`ptr<T>` 类型，不能用作整数、布尔值或数组值。

当分配或查找函数没有结果时，可以返回`null`。在取消引用此类结果之前检查`null`。取消引用 `null` 指针不会访问有效的存储。

## 指针转换

当您需要更改地址或其他指针表示形式时，请使用`as`。

```wave
var raw: i64 = 0;
var p: ptr<u8> = raw as ptr<u8>;
```

仅在低级边界上使用整数和指针之间的转换，并考虑目标平台的地址宽度和ABI。
