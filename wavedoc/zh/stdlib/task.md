---
translation_set_id: stdlib-task
path: stdlib/task
locale: zh
group: stdlib
group_order: 1
order: 13
title: task: 运行与清理 Future
summary: 描述异步操作的单次消费、执行和清理后取消。
---

## 基本款API

导入为`import("std::task" as task);`。

```text
block_on<T>(future: Future<T>) -> T
spawn<T>(future: Future<T>) -> Future<T>
cancel<T>(future: Future<T>) -> bool
yield_now() -> Future<void>
sleep_ms(milliseconds: i64) -> Future<void>
cancel_and_join<T>(future: Future<T>) -> Future<void>
shutdown()
```

`task::block_on(pending)` 运行直到 Future 完成并返回结果。结果类型由传递的Future确定。

## 生活规则

Future是一次性消费手柄。我不认为复制值是两个独立的操作。您还负责等待结果或取消/组织通过`spawn`安排的任务。 Future已与`await`一起消耗，`block_on`将不会再次消耗。

取消请求取消；它本身并不保证清理已经完成。在释放资源之前等待所需的完成。当任务仍可以访问异步 I/O 借用的内存时，请勿释放该内存。任务完成或完成取消清理后调用`shutdown`。

## 协同执行

长计算和同步 blocking 调用可能会减慢执行器的整体进度。使用 yield 进行异步等待会产生执行机会。阻塞I/O不会仅仅因为它位于函数async内部而变得异步。

查看[异步代码简介](/docs/zh/language/async-and-never)中完整程序的执行和清理顺序。

## 屈服并等待完成

以下程序在操作中间产生执行并返回结果。由于yield不是函数终止，因此await之后的代码会连续执行。

<!-- wave-example: book-task-yield -->
```wave
import("std::task" as task);

async fun calculate() -> i32 {
    println("started");
    await task::yield_now();
    println("resumed");
    return 42;
}

fun main() -> i32 {
    var pending: Future<i32> = calculate();
    var result: i32 = task::block_on(pending);

    println("result={}", result);
    task::shutdown();
    return 0;
}
```

执行结果：

```text
started
resumed
result=42
```

本例中没有其他操作，因此输出顺序是一致的。在具有多个任务spawn的程序中，其他任务可以在yield点继续进行，因此不依赖于不同任务的输出顺序。

## 组织任务的顺序

1. 准备好工作所需的存储空间和资源。
2. 创建 Future 并将其作为 await、block_on 或 spawn 运行。
3. 如果您需要结果，请等到完成。
4. 如果您取消了正在运行的任务，请等待它清理。
5. 清理任务借用的缓冲区、文件和套接字。
6. 如果没有剩余工作，请致电shutdown。

Future 离开变量的作用域是一回事，安全地清理工作则是另一回事。特别是，如果将函数本地数组的地址传递给任务async，则该任务必须在函数返回之前完成对该数组的使用。

## async 划分功能和通用功能

纯计算可以分解为常规函数。将async附加到需要表达等待的函数上，并在该函数内使用await等待完成。仅使用 async 函数包装需要很长时间的同步函数（例如读取文件）并不会给其他任务执行的机会。

您可以在[异步学习](/docs/zh/language/async-and-never)中创建并运行Future时进行比较。
