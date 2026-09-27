---
translation_set_id: memory-buffer
path: reference/memory-and-buffer
locale: zh
group: stdlib
group_order: 1
order: 4
title: mem:分配、重新分配和布局
summary: 描述大小（以字节为单位）、分配失败、重新分配边界和释放责任。
---

## 分配和解除分配

```text
std::mem::alloc
mem_alloc(size: i64) -> ptr<u8>
mem_alloc_zeroed(size: i64) -> ptr<u8>
mem_free(p: ptr<u8>, size: i64) -> i64
mem_realloc(old_ptr: ptr<u8>, old_size: i64, new_size: i64) -> ptr<u8>
```

大小以字节为单位。大小为 0 或更小的分配会返回 null。如果正大小的分配也失败，则可能是null。不要假设`mem_alloc`的初始内容，但如果需要零初始化，请使用`mem_alloc_zeroed`。

调用者拥有每个成功的分配，并且在释放它时必须传递其原始大小。 `mem_free(null, size)` 返回 0。非 null 指针与非正大小配对是错误的。切勿在释放后访问或释放分配。

## 重新分配时的分类

|请求|行动|
| --- | --- |
|新尺寸积极且成功|`min(old_size, new_size)` 复制字节并释放前一个字节|
|正大小的新分配失败|null 返回，维持现有分配|
|`old_ptr == null`，正新尺寸|行为就像一个新任务|
| `new_size == 0` |尝试释放有效的先前分配并返回 null|
|old_size=0 表示负大小或现有指针|null 返回|

重新分配到大小零的null结果并不能确定释放是否成功。如果您需要其状态，请直接致电`mem_free`。增加分配后，请自行初始化新添加的区域。

## 保留现有指针的示例

下面是新尺寸为正值的情况。将其保存为`main.wave`并运行。

<!-- wave-example: reallocation -->
```wave
import("std::mem::alloc")::{
    mem_alloc, mem_realloc, mem_free
};

fun main() -> i32 {
    var data: ptr<u8> = mem_alloc(4);
    if (data == null) {
        return 1;
    }
    deref data[0] = 7;
    var grown: ptr<u8> = mem_realloc(data, 4, 8);
    if (grown == null) {
        mem_free(data, 4);
        return 2;
    }
    data = grown;
    println("{}", deref data[0]);
    if (mem_free(data, 8) < 0) {
        return 3;
    }

    return 0;
}
```

执行结果：

```text
7
```

如果在确认失败之前用`data = mem_realloc(...)`覆盖的话，可能会丢失现有的地址。如果重新分配成功，则旧地址和指向它的指针将不被使用。

## 目标类型的大小和对齐方式

```text
std::mem::layout
size_of<T>() -> u64
align_of<T>() -> u64
```

这两个值都是编译目标的布局，而不是运行它的计算机的布局。 `size_of` 包括尾部填充，并且不生成或评估值。当将元素数量乘以其大小时，我们检查是否溢出。您可以使用`std::mem::ops`的`mem_size_mul_checked`和`mem_size_add_checked`。

`mem_copy` 用于复制不重叠的范围，`mem_move` 用于复制可能重叠的范围。两者都无法仅根据指针确定实际分配长度，因此调用者必须保证边界。可以使用 [Buffer](/docs/zh/stdlib/buffer) 管理不同大小的字节列表。

## 布局查询示例

将其保存为main.wave并运行。文档中涉及的目标i32的大小和对齐方式均为4字节，因此输出`4 4`。通过目标检查其他类型的值，尤其是结构体和指针。

<!-- wave-example: layout-api -->
```wave
import("std::mem::layout")::{
    size_of, align_of
};

fun main() {
    println("{} {}", size_of<i32>(), align_of<i32>());
}
```

## 检查尺寸计算是否溢出

在将 `count * element_size` 传递给赋值函数之前，您需要检查它是否有效。如果你用溢出值分配一个小空间并写入与原始数量一样多的内容，就会出界。

<!-- wave-example: book-memory-size-check -->
```wave
import("std::mem::ops")::{mem_size_add_checked, mem_size_mul_checked};

fun main() {
    var result: i64 = 99;

    if (mem_size_mul_checked(3, 4, &result) == 0) {
        println("bytes={}", result);
    }

    if (mem_size_add_checked(9223372036854775807, 1, &result) < 0) {
        println("overflow rejected");
    }
}
```

执行结果：

```text
bytes=12
overflow rejected
```

失败结果不会写入分配大小。仅通过常规算术计算并查看结果是否为负，无法检测所有溢出。要计算需要检查的尺寸，请从一开始就使用 checked 函数。

## 对于重叠副本，mem_move

当向后移动同一数组的一部分时，输入和输出区域会重叠。不要将重叠范围传递给mem_copy，请使用mem_move。

<!-- wave-example: book-memory-overlap -->
```wave
import("std::mem::ops")::{mem_move};

fun main() {
    var data: array<u8, 5> = [1, 2, 3, 4, 5];

    mem_move(&data[1], &data[0], 4);

    for (var index: i32 = 0; index < 5; index += 1) {
        println("{}", data[index]);
    }
}
```

执行结果：

```text
1
1
2
3
4
```

原来的前四个字节向右移动一个位置。手动向前复制它们可能会读取已被覆盖的值，从而意外地生成所有值。重叠感知 API 会为您处理复制方向。

## 编写一个传递所有权的函数

对于返回内存的函数，最好在成功时提供返回地址以及释放内存所需的大小。如果调用者必须猜测大小，则很容易出现错误释放。如果函数返回借用地址，则调用者不应释放该地址，并描述原始地址的生命周期。

跨越函数边界，您应该能够跟踪`allocator → owner → deallocation`。指针变量的名称或类型不会自动确定所有权。
