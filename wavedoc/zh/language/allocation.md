---
translation_set_id: learn-allocation
path: language/allocation
locale: zh
group: language
group_order: 2
order: 10
title: 10.内存分配和资源管理
summary: 了解分配失败、初始化、范围和释放。
---

## 什么时候需要动态存储空间？

固定大小的数组在其类型中包含其长度。当数据量仅在运行时已知时（例如文件大小或输入长度），请使用动态内存。当您不再需要每个分配时，将其释放。

在本章中，您将管理一个小型分配，调整其大小，然后使用Buffer。传递指针与转移所有权不同。按照示例操作时，确定每个函数拥有哪些资源。

## 分配、检查、使用、释放

<!-- wave-example: book-alloc-lifecycle -->
```wave
import("std::mem::alloc")::{
    mem_alloc_zeroed, mem_free
};

fun main() -> i32 {
    var size: i64 = 4;
    var data: ptr<u8> = mem_alloc_zeroed(size);

    if (data == null) {
        println("allocation failed");
        return 1;
    }

    deref data[0] = 42;
    println("{} {}", deref data[0], deref data[1]);

    var status: i64 = mem_free(data, size);

    if (status < 0) {
        return 2;
    }

    return 0;
}
```

执行结果：

```text
42 0
```

该程序有四个步骤：请求 4 个字节，检查null，仅访问有效范围，并释放分配。成功后，mem_alloc_zeroed将内存初始化为零，因此即使程序尚未写入第二个字节也是零。

不要假设mem_alloc返回的内存的任何初始内容。在读取每个区域之前对其进行初始化。零或负分配大小返回null。正大小分配也可能无法获取内存。

## 尺寸单位

内存分配 API 的大小参数以字节为单位。要分配十个整数，请将元素大小乘以元素数量。检查这个乘法没有溢出。

<!-- wave-example: book-alloc-typed -->
```wave
import("std::mem::alloc")::{
    mem_alloc, mem_free
};
import("std::mem::layout")::{
    size_of
};
import("std::mem::ops")::{
    mem_size_mul_checked
};

fun main() -> i32 {
    var count: i64 = 3;
    var bytes: i64 = 0;
    var item_size: i64 = size_of<i32>() as i64;

    if (mem_size_mul_checked(count, item_size, &bytes) < 0) {
        return 1;
    }

    var raw: ptr<u8> = mem_alloc(bytes);

    if (raw == null) {
        return 2;
    }

    var values: ptr<i32> = raw as ptr<i32>;

    for (var index: i64 = 0; index < count; index += 1) {
        deref values[index] = (index + 1) as i32;
    }

    println("{} {} {}", deref values[0], deref values[1], deref values[2]);

    if (mem_free(raw, bytes) < 0) {
        return 3;
    }

    return 0;
}
```

执行结果：

```text
1 2 3
```

count 是元素计数； bytes 是字节数。指针算术以i32为单位移动，但释放分配需要其原始大小（以字节为单位）。 size_of 使用目标类型布局，使与元素类型的关系明确。

对于一般非常大的类型，将size_of的结果更改为i64时，还必须考虑转换范围。这里我们使用i32，它的大小已知。

## 即使在故障路径上也能进行清理

如果分配后另一个操作失败，请在提前返回之前释放内存。所有权表有助于识别您可能会错过的清理路径。

|步骤|拥有的资源|如果你失败了|
| --- | --- | --- |
|分配前|无|立即归还|
|分配成功后|data 和原始尺寸|data 发布后返回|
|重新分配成功后|新地址和新尺寸|新地址发布|
|释放后|无|不要使用旧地址|

覆盖指针变量并丢失原始地址也会丢失释放分配所需的信息。这会导致内存泄漏。相反，通过两个所有者释放相同的分配会导致双重释放。

## 重新分配以增加规模

<!-- wave-example: book-alloc-grow -->
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
    var next: ptr<u8> = mem_realloc(data, 4, 8);

    if (next == null) {
        mem_free(data, 4);
        return 2;
    }

    data = next;
    deref data[4] = 9;
    println("{} {}", deref data[0], deref data[4]);

    if (mem_free(data, 8) < 0) {
        return 3;
    }

    return 0;
}
```

执行结果：

```text
7 9
```

首先将结果存储在 next 中。如果较大块的分配失败，原始数据仍然有效并且仍然可以被释放。成功后，旧的分配将被释放，并且必须使用新的地址。在读取新添加的区域之前先对其进行初始化。

本例中的新尺寸为正值。 new_size=0 的请求会尝试释放现有分配并返回 null。因此，null结果并不总是意味着旧的分配仍然有效。当需要检查释放是否成功时，直接调用mem_free即可。

## 何时重新检查借用的指针

data 保存内部指针并在重新分配后使用它是错误的。这是因为新的data的地址可能不同。如果需要内部位置，可以存储offset代替地址，成功后根据新的data重新计算。

释放或重新分配内存也会影响借用它的代码。检查其他操作是否仍在使用该内存。传递给异步操作的缓冲区必须保持有效，直到操作完成。

## 字节列表包含Buffer

在手动重新分配时管理长度频繁变化的字节列表需要处理len和cap、扩展失败和大小计算。 std 的Buffer 捆绑了这些操作。

<!-- wave-example: book-buffer-progressive -->
```wave
import("std::buffer::alloc")::{
    Buffer, buffer_init, buffer_free
};
import("std::buffer::write")::{
    buffer_append_str
};

fun main() -> i32 {
    var message: Buffer;

    if (buffer_init(&message, 0) < 0) {
        return 1;
    }

    if (buffer_append_str(&message, "Hello") < 0) {
        buffer_free(&message);
        return 2;
    }

    if (buffer_append_str(&message, ", Wave") < 0) {
        buffer_free(&message);
        return 3;
    }

    println("bytes={}", message.len);

    if (buffer_free(&message) < 0) {
        return 4;
    }

    return 0;
}
```

执行结果：

```text
bytes=11
```

len 是正在使用的字节数； cap 是分配的容量。必要时附加数据会增加分配。使用Buffer并不能免除调用者释放它的责任。

buffer_append_str 不会将 NUL 添加到字符串末尾。因此，message.data不应该直接输出为str。使用 I/O 函数输出字节，该函数获取长度，或显式构造字符串表示形式。

## 练习与解题思路

将字节0到9一一加到Buffer并得到和。每次追加失败时必须释放它，并且仅在len范围内执行读取。完整解和边界失效可以按照[Buffer 使用方法](/docs/zh/stdlib/buffer)中的示例进行检查。

尝试在代码中标记分配、重新分配和取消分配调用。对于每个成功的分配，您必须能够描述谁拥有它以及哪条路径释放它。
