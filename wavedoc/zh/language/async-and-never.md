---
translation_set_id: language-async-and-never
path: language/async-and-never
locale: zh
group: language
group_order: 2
order: 13
title: 13. 异步函数与 Future
summary: 了解延迟执行Future、await和block_on的作用。
---

## 表达等待任务

对于文件、套接字、定时器等等待操作，需要区分继续计算和等待完成。异步函数将要完成的结果表示为Future。添加 async 不会自动创建新线程或将所有同步调用更改为异步。

在函数、指针和错误处理之后阅读本章。该示例是一个本机程序，它使用 `std::task` 启动器并将每个文件作为 `wavec run main.wave` 运行。

## 创建 Future 并接收结果

<!-- wave-example: book-async-basic -->
```wave
import("std::task" as task);

async fun calculate(value: i64) -> i64 {
    return value * 2;
}

fun main() {
    var pending: Future<i64> = calculate(21);
    var result: i64 = task::block_on(pending);

    task::shutdown();
    println("{}", result);
}
```

执行结果：

```text
42
```

calculate的声明中写的i64是完成后得到的值。调用本身的结果是Future<i64>。正常情况下main，用block_on运行Future，即可得到完成的结果。

清理完所有任务后调用 shutdown 来释放执行器的资源。当任务仍在使用缓冲区时，不要释放缓冲区，也不要忽略未完成的任务。

## 调用和主体执行不同

异步函数延迟执行。它必须与调用后立即执行其主体的常规函数​​区分开来。

<!-- wave-example: book-async-lazy -->
```wave
import("std::task" as task);

static entered: i32 = 0;

async fun work() -> i32 {
    entered += 1;
    return 7;
}

fun main() {
    var pending: Future<i32> = work();

    println("before={}", entered);

    var result: i32 = task::block_on(pending);

    task::shutdown();
    println("after={} result={}", entered, result);
}
```

执行结果：

```text
before=0
after=1 result=7
```

创建Future时，entered还是0，驱动执行后，body执行完就变成1了。这就是为什么你不应该因为Future存储在变量中就认为任务已完成。

## 在异步函数内等待

在async函数内，await等待另一个Future的完成。表达式await的结果是完成值。

<!-- wave-example: book-async-nested -->
```wave
import("std::task" as task);

async fun twice(value: i32) -> i32 {
    await task::yield_now();
    return value * 2;
}

async fun process() -> i32 {
    var value: i32 = await twice(20);
    return value + 2;
}

fun main() {
    var answer: i32 = task::block_on(process());

    task::shutdown();
    println("{}", answer);
}
```

执行结果：

```text
42
```

process等待twice的Future，然后将完成值加2。您可以使用与调用常规函数的结果类似的值，但在等待期间，您可以将执行机会传递给另一个任务。

yield_now 产生合作执行机会。在长计算循环中即使只给出一次失败也会减慢其他任务的速度。异步CPU不是自动并行和分布计算的设备。

## 安排多项任务

您可以使用spawn安排任务并等待各自的结果。这是一个不依赖两个操作的中间输出顺序来检查最终结果的示例。

<!-- wave-example: book-async-spawn -->
```wave
import("std::task" as task);

async fun compute(value: i32) -> i32 {
    await task::yield_now();
    return value * 2;
}

async fun combine() -> i32 {
    var first: Future<i32> = task::spawn(compute(10));
    var second: Future<i32> = task::spawn(compute(20));

    var left: i32 = await first;
    var right: i32 = await second;

    return left + right;
}

fun main() {
    var total: i32 = task::block_on(combine());

    task::shutdown();
    println("{}", total);
}
```

执行结果：

```text
60
```

您安排的每个任务都有一个等待结果的地方。您需要决定由谁来检查其完成情况，而不是创建任务并忘记其句柄。 await 不要认为内部操作执行的顺序和顺序是相同的。

## 食用Future一次

Future 被视为单个消费句柄。复制相同的 Future 不会使其像两个不同的任务一样等待。请勿对已完成的Future再次执行block_on或await。

如果您需要在多个位置获得相同的结果，请存储完整的值并根据该值的复制/共享规则传递它，而不是多次使用 Future。您还应该检查值中是否有指针或拥有的资源。

## 定时器和同步等待的区别

async 在函数内等待时，可以使用`await task::sleep_ms(...)`。调用同步sleep会阻塞当前的执行流程，这也可能会影响执行器中其他任务的进度。

我不认为等待的时间恰好是请求的毫秒数。根据您的日程安排和其他任务，您可能会起晚。在实现超时时，我们使用clock和deadline来测量经过的时间，而不是每次都再次等待原始的完整时间。

## 缓冲区寿命和取消

即使调用函数被挂起，传递给异步 I/O 的缓冲区也必须保持有效。在完成或取消清理之前释放或重新分配它可能会使操作留下无效地址。

不能假设取消请求和任务完成是同时发生的。检查cancel系列API的结果并等待必要的完成后再释放资源。详细调用规则请阅读[参见task](/docs/zh/stdlib/task)。

## 常见的误解

|认为|实际上检查|
| --- | --- |
|async 我打电话过去了。|你真的运行并完成了Future吗？|
|async 函数内的所有调用都是异步的|被调用的API是同步还是异步？|
|Future 复制重复任务|您是否重复使用同一个手柄？|
|既然你取消了它，你可以立即释放缓冲区。|取消后的工作你整理完了吗？|
|中间输出顺序始终是固定的|它是否明确地只等待结果所需的顺序？|

## 练习与完整解答

创建一个连续等待三个异步函数的管道。返回乘以 5 后的结果。

<!-- wave-example: book-async-solution -->
```wave
import("std::task" as task);

async fun read_value() -> i32 {
    return 10;
}

async fun transform(value: i32) -> i32 {
    await task::yield_now();
    return value * 2;
}

async fun pipeline() -> i32 {
    var initial: i32 = await read_value();
    var changed: i32 = await transform(initial);

    return changed + 5;
}

fun main() {
    var result: i32 = task::block_on(pipeline());

    task::shutdown();
    println("{}", result);
}
```

执行结果：

```text
25
```

这个例子是故意的顺序依赖。 transform 需要 read_value 的结果，因此简单地执行所有 spawn 并不会使关系消失。异步设计的出发点是区分独立操作和需要结果的操作。

完成基础知识后，即可使用 [文件阅读练习](/docs/zh/practice/file-reader) 和 [TCP 练习](/docs/zh/practice/tcp-client) 连接到真实的外部资源。

## void 和 never

省略返回类型的常规函数可能会返回没有值的调用点。 never类型写为`!`，表示不会正常返回调用点。一个代表性的例子是进程终止函数。

说明声明的示例：

```text
fun log(message: str)              // 작업 후 돌아옴, 결과값 없음
fun stop(code: i32) -> !           // 정상 반환하지 않음
async fun work() -> i64            // 호출 결과는 Future<i64>
```

没有尝试使 never 成为通用存储值。不要编写声明为非返回的函数以具有正常的返回路径。如果退出前需要清理资源，则调用者必须首先执行此操作。

## 不返回函数的完整示例

如果将其保存为 main.wave 并运行它，它将以退出代码 0 结束，没有输出。 stop不会返回给调用者，因此被声明为`-> !`。

<!-- wave-example: never-function -->
```wave
import("std::process::core")::{
    proc_exit
};

fun stop() -> ! {
    proc_exit(0);
}

fun main() {
    stop();
}
```
