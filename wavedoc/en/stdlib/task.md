---
translation_set_id: stdlib-task
path: stdlib/task
locale: en
group: stdlib
group_order: 1
order: 13
title: task: Running and cleaning up Futures
summary: Describes single consumption, execution, and cancellation-after-cleanup of asynchronous operations.
---

## Basic API

Imported as `import("std::task" as task);`.

```text
block_on<T>(future: Future<T>) -> T
spawn<T>(future: Future<T>) -> Future<T>
cancel<T>(future: Future<T>) -> bool
yield_now() -> Future<void>
sleep_ms(milliseconds: i64) -> Future<void>
cancel_and_join<T>(future: Future<T>) -> Future<void>
shutdown()
```

`task::block_on(pending)` runs until Future completes and returns the result. The result type is determined from the passed Future.

## life rules

Future is a one-time consumption handle. I don't think of copying values ​​as two independent operations. You are also responsible for waiting for results or canceling/organizing tasks scheduled with `spawn`. Future already consumed with `await` and `block_on` will not be consumed again.

cancel requests cancellation; it does not by itself guarantee that cleanup has finished. Wait for the required completion before releasing resources. Do not free memory borrowed by asynchronous I/O while a task can still access it. Call `shutdown` after tasks have completed or finished cancellation cleanup.

## collaborative execution

Long computations and synchronous blocking calls can slow down the overall progress of the executor. Asynchronous wait with yield yields execution opportunity. Blocking I/O does not become asynchronous just because it is inside a function async.

Check out the execution and cleanup sequence in the full program in [Introduction to asynchronous code](/docs/en/language/async-and-never).

## Yield and wait for completion

The following program yields execution in the middle of an operation and returns a result. Since yield is not a function termination, the code after await is executed continuously.

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

Execution result:

```text
started
resumed
result=42
```

There are no other operations in this example, so the output order is consistent. In a program that has multiple tasks spawn, other tasks can proceed at the yield point, so it does not depend on the output order of the different tasks.

## Order of organizing tasks

1. Prepare the storage space and resources you need for your work.
2. Create Future and run it as await, block_on, or spawn.
3. If you need results, wait until completion.
4. If you canceled a running task, wait for it to clean up.
5. Cleans up buffers, files, and sockets borrowed by the task.
6. With no work remaining, call shutdown.

Future Leaving the scope of a variable is one thing, and having the work safely cleaned up is another. In particular, if you pass the address of a function-local array to a task async, the task must finish using that array before the function returns.

## async Dividing functions and general functions

Pure calculations can be separated into regular functions. Attach async to the function that needs to express wait, and wait for completion with await within that function. Just wrapping a synchronous function that takes a long time, such as reading a file, with the async function does not give other tasks a chance to execute.

You can compare when creating and running Future in [asynchronous learning](/docs/en/language/async-and-never).
