---
translation_set_id: language-async-and-never
path: language/async-and-never
locale: en
group: language
group_order: 2
order: 13
title: 13. Asynchronous functions and Future
summary: Learn the role of delayed execution Future, await, and block_on.
---

## Expressing waiting tasks

For waiting operations such as files, sockets, and timers, it is necessary to distinguish between continuing calculations and waiting for completion. The asynchronous function expresses the result to be completed as Future. Adding async does not automatically create a new thread or change all synchronous calls to asynchronous.

Read this chapter after Functions, Pointers and Error Handling. The example is a native program that uses the `std::task` launcher and runs each file as `wavec run main.wave`.

## Create Future and receive results

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

Execution result:

```text
42
```

i64 written in the declaration of calculate is the value obtained after completion. The result of the call itself is Future<i64>. In normal main, run Future with block_on and receive the completed result.

Call shutdown after cleaning up all tasks to release the executor’s resources. Do not free a buffer while a task still uses it, and do not ignore unfinished tasks.

## Invocation and body execution are different

Asynchronous functions execute lazily. It must be distinguished from a regular function that executes its body immediately after calling it.

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

Execution result:

```text
before=0
after=1 result=7
```

When Future was created, entered was still 0. After driving execution, the body is executed and becomes 1. This is why you should not treat the task as completed just because Future is stored in a variable.

## Waiting inside an asynchronous function

Inside the async function, await waits for the completion of another Future. The result of an expression await is a completion value.

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

Execution result:

```text
42
```

process waits for Future of twice and then adds 2 to the completion value. You can use the value similar to the result of a call to a regular function, but while waiting, you can pass the opportunity for execution to another task.

yield_now yields cooperative execution opportunities. Failure to give even once in a long computation loop can slow down other tasks. Asynchronous CPU is not a device for automatically parallelizing and distributing computations.

## Schedule multiple tasks

You can schedule tasks with spawn and wait for their respective results. This is an example of checking the final result without relying on the intermediate output order of the two operations.

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

Execution result:

```text
60
```

Each task you schedule has a place where it waits for results. Instead of creating a task and forgetting its handle, you need to decide who will check its completion. await Do not consider the sequence and any order in which internal operations are executed to be the same.

## Consume Future once

Future is treated as a single consumption handle. Copying the same Future doesn't make it wait like two different tasks. Do not block_on or await again for Future that has already been completed.

If you need the same result in multiple places, instead of consuming Future multiple times, store the completed value and pass it along according to the copy/share rules for that value. You should also check whether there are pointers or owned resources within the value.

## Difference between timer and synchronous wait

async When waiting within a function, you can use `await task::sleep_ms(...)`. Calling synchronous sleep blocks the current flow of execution, which may also affect the progress of other tasks in the executor.

I don't expect the time waited to be exactly the number of milliseconds requested. Depending on your scheduling and other tasks, you may wake up late. When implementing timeouts, we use clock and deadline to measure elapsed time, as opposed to waiting for the original full time again each time.

## Buffer Lifetime and Cancellation

A buffer passed to asynchronous I/O must remain valid even while the calling function is suspended. Freeing or reallocating it before completion or cancellation cleanup can leave the operation with an invalid address.

It cannot be assumed that the cancellation request and the task completion are at the same time. Check the results of cancel series API and wait for necessary completion before releasing resources. Please read [See task](/docs/en/stdlib/task) for detailed calling rules.

## common misunderstanding

|think|actually check|
| --- | --- |
|async I called and it’s over.|Did you actually run and complete Future?|
|async All calls within the function are asynchronous|Is the called API synchronous or asynchronous?|
|Future Copy duplicates a task|Are you consuming the same handle repeatedly?|
|Since you canceled it, you can release the buffer immediately.|Have you finished organizing the work after cancellation?|
|Intermediate output order is always fixed|Does it explicitly wait only for the order required for the result?|

## Exercise and complete solution

Create a pipeline that waits for three asynchronous functions in a row. Returns the result of doubling and adding 5.

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

Execution result:

```text
25
```

This example is intentionally a sequential dependency. transform requires the result of read_value, so simply doing all spawn does not make the relationship go away. The starting point of asynchronous design is the distinction between independent operations and operations that require results.

Once you've completed the basics, connect to real external resources with [File reading practice](/docs/en/practice/file-reader) and [TCP Practice](/docs/en/practice/tcp-client).

## void and never

A regular function that omits a return type may return to the point of call without a value. The never type is written as `!`, meaning it does not return to the call point normally. A representative example is the process termination function.

Example illustrating the declaration:

```text
fun log(message: str)              // 작업 후 돌아옴, 결과값 없음
fun stop(code: i32) -> !           // 정상 반환하지 않음
async fun work() -> i64            // 호출 결과는 Future<i64>
```

There is no attempt to make never a common stored value. Do not write functions that are declared as non-returning to have a normal return path. If resource cleanup is required before exiting, the caller must do so first.

## Full example of a function that doesn't return

If you save it as main.wave and run it, it ends with exit code 0 without output. stop does not return to the caller, so it is declared as `-> !`.

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
