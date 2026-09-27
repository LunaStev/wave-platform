---
translation_set_id: learn-errors
path: language/errors
locale: en
group: language
group_order: 2
order: 12
title: 12. Representing and recovering from errors
summary: Separate errors from normal values and clean up resources from failure paths.
---

## Result of failure degree function

Missing files or out-of-bounds input occur naturally in programs. Handling an error is more than just printing a message. It is the process of distinguishing failures, checking the status of work already done, organizing acquired resources, and then choosing whether to continue or terminate.

In this chapter, we start with the failure representation of a small function and progress to the result structure, variant, early return, and resource cleanup.

## Failure markers must not overlap success values

The reason -1 is used as not found in an array search is because the effective index is 0 or more. On the other hand, in calculations where any integer can be a normal result, if -1 is designated as an error, it cannot be distinguished from the normal value -1.

0 is also a commonly misunderstood value. Empty string length 0, first position 0, and number of bytes transferred 0 have different meanings for different functions. Don't judge success just because the return value is non-zero.

## Return success and value together

<!-- wave-example: book-error-result -->
```wave
struct Division {
    ok: bool;
    value: i32;
}

fun divide_nonnegative(left: i32, right: i32) -> Division {
    if (left < 0 || right <= 0) {
        return Division {
            ok: false,
            value: 0
        };
    }

    return Division {
        ok: true,
        value: left / right
    };
}

fun main() {
    var result: Division = divide_nonnegative(0, 3);

    if (!result.ok) {
        println("invalid input");
        return;
    }

    println("value={}", result.value);
}
```

Execution result:

```text
value=0
```

A normal result may also be 0. Instead of looking at value and guessing whether it was successful, check ok first. Even if the field value exists in the failed result, it does not mean that it is the value to use.

The input rules for this example are that the left side is 0 or greater and the right side is positive. The scope is shown in the function name and description to distinguish it from the typical signed integer division.

## Distinguish between causes of errors

Add error information for additional guidance or recovery depending on the reason for the failure. There is also a way to divide the function into a function that checks the input range and a function that calculates it.

<!-- wave-example: book-error-codes -->
```wave
struct Check {
    ok: bool;
    error: i32;
}

fun check_quantity(value: i32) -> Check {
    if (value < 1) {
        return Check {
            ok: false,
            error: 1
        };
    }

    if (value > 100) {
        return Check {
            ok: false,
            error: 2
        };
    }

    return Check {
        ok: true,
        error: 0
    };
}

fun main() {
    var result: Check = check_quantity(101);

    if (!result.ok) {
        match (result.error) {
            1 => {
                println("quantity is too small");
            }
            2 => {
                println("quantity is too large");
            }
            _ => {
                println("invalid quantity");
            }
        }
    }
}
```

Execution result:

```text
quantity is too large
```

The meaning of the numbers is defined by this function. It cannot be considered the same as error number 1 or 2 in other libraries. In public API, naming error constants or types saves the caller from having to memorize random numbers.

## Separate results with variant

The relationship that a success value and an error cannot exist at the same time can be expressed as variant.

<!-- wave-example: book-error-variant -->
```wave
variant ByteResult {
    Value(u8),
    OutOfRange(i32)
}

fun to_byte(value: i32) -> ByteResult {
    if (value < 0 || value > 255) {
        return ByteResult::OutOfRange(value);
    }

    return ByteResult::Value(value as u8);
}

fun main() {
    var result: ByteResult = to_byte(300);

    match (result) {
        ByteResult::Value(value) => {
            println("byte={}", value);
        }
        ByteResult::OutOfRange(original) => {
            println("out of range: {}", original);
        }
    }
}
```

Execution result:

```text
out of range: 300
```

The error contained the original input value. It's easier for the caller to explain the problem than simply returning false. The nature of the data is also taken into consideration by not leaving sensitive inputs such as passwords or tokens in the log.

## Make normal routes easier to read with early returns

There is no need to put all the good code deep inside if when performing a multi-step check. If it fails, you can return immediately and continue the success path below.

<!-- wave-example: book-error-early-return -->
```wave
fun show_total(quantity: i32, price: i32) {
    if (quantity < 1 || quantity > 100) {
        println("invalid quantity");
        return;
    }

    if (price < 0 || price > 100000) {
        println("invalid price");
        return;
    }

    var total: i32 = quantity * price;
    println("total={}", total);
}

fun main() {
    show_total(3, 1200);
    show_total(0, 1200);
    show_total(3, -1);
}
```

Execution result:

```text
total=3600
invalid quantity
invalid price
```

After passing the test, you can take advantage of the fact that quantity and price are within the specified range. The range was set so that the intermediate multiplication also falls within i32. When adding an early return, you should also check to see if you already own any resources at that point.

## Memory cleanup of failure paths

<!-- wave-example: book-error-cleanup -->
```wave
import("std::buffer::alloc")::{
    Buffer, buffer_init, buffer_free
};
import("std::buffer::write")::{
    buffer_append_str
};

fun make_message(allowed: bool) -> i32 {
    var data: Buffer;

    if (buffer_init(&data, 16) < 0) {
        return 1;
    }

    if (!allowed) {
        buffer_free(&data);
        return 2;
    }

    if (buffer_append_str(&data, "ready") < 0) {
        buffer_free(&data);
        return 3;
    }

    println("bytes={}", data.len);

    if (buffer_free(&data) < 0) {
        return 4;
    }

    return 0;
}

fun main() {
    println("status={}", make_message(false));
}
```

Execution result:

```text
status=2
```

Release Buffer even on failure paths that do not output data. Rather than turning every failure into a single return, it is important to clearly manage what is owned at each branch.

There is also API where the cleanup itself fails. Decide how you will preserve errors and cleanup errors in the original work. This small example first returns the task error number. Larger programs may record both separately.

## Partial successes are not automatically canceled

If you write some bytes to a file and then the write fails, the bytes already written are not lost. The other end of the network may also have received some data. Repeating the same task from the beginning may result in duplicate records.

Conversely, reading checked byte cursor preserves the position and output value when it fails. These functions can receive more input and try again at the same location. Instead of applying the “if it fails, nothing changes” rule to every API, check the documentation for that function.

## Recoverable errors and traps

An invalid file path or invalid user input can be reported as an error value so that the caller can recover. A trap caused by an invalid runtime shift count or floating-point-to-integer conversion is a different mechanism.

Programs that must continue executing must inspect their input before hazardous operations. assert is also not intended to be abused as a means of handling graceful failures of user input. If you need to give the user a chance to re-enter their input, return results to continue the control flow.

## Exercise and complete solution

Write a function that reads an array element by index and fails when the index is negative or greater than or equal to the length. Use a result structure because zero can be a valid element value.

<!-- wave-example: book-error-solution -->
```wave
struct Lookup {
    ok: bool;
    value: i32;
}

fun at(values: ptr<i32>, length: i32, index: i32) -> Lookup {
    if (index < 0 || index >= length) {
        return Lookup {
            ok: false,
            value: 0
        };
    }

    return Lookup {
        ok: true,
        value: deref values[index]
    };
}

fun main() {
    var values: array<i32, 3> = [0, 10, 20];
    var first: Lookup = at(&values[0], 3, 0);
    var outside: Lookup = at(&values[0], 3, 3);

    if (first.ok) {
        println("value={}", first.value);
    }

    if (!outside.ok) {
        println("out of bounds");
    }
}
```

Execution result:

```text
value=0
out of bounds
```

It is the caller's condition that the pointer and length represent the actual readable array. This is not a function that makes random addresses safe just by checking the index. Please read separately what checks the function is responsible for and what conditions the caller must guarantee.
