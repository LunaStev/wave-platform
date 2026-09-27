---
translation_set_id: stdlib-buffer
path: stdlib/buffer
locale: en
group: stdlib
group_order: 1
order: 5
title: buffer: Growable byte storage
summary: Buffer Describes initialization, addition, inquiry, capacity and release rules.
---

## Meaning of Buffer

`Buffer` in `std::buffer::types` has `data: ptr<u8>`, `len: i64`, and `cap: i64`. len is the number of bytes initialized and in use, and cap is the total number of bytes allocated. Always maintain `0 <= len <= cap`. The string NUL does not automatically guarantee termination.

## Basic API

|module|declaration|meaning|
| --- | --- | --- |
| `std::buffer::alloc` | `buffer_init(out_buf: ptr<Buffer>, capacity: i64) -> i64` |Initialize new repository. Do not recall into buffers already owned|
|same module| `buffer_free(buf: ptr<Buffer>) -> i64` |Deallocate. Empty if successful|
|same module| `buffer_reserve(buf: ptr<Buffer>, required_cap: i64) -> i64` |Ensure minimum full capacity. len Maintained|
|same module| `buffer_clear(buf: ptr<Buffer>) -> i64` |Maintain capacity and len=0|
|same module| `buffer_resize(buf: ptr<Buffer>, new_len: i64, value: u8) -> i64` |Change length, fill new bytes with value|
| `std::buffer::write` | `buffer_push(buf: ptr<Buffer>, value: u8) -> i64` |add one byte|
|same module| `buffer_append(buf: ptr<Buffer>, src: ptr<u8>, size: i64) -> i64` |sizeCopy and add byte|
|same module| `buffer_append_str(buf: ptr<Buffer>, s: str) -> i64` |Add string bytes excluding NUL|
| `std::buffer::read` | `buffer_get(buf: ptr<Buffer>, index: i64, out_value: ptr<u8>) -> i64` |Read one byte in range|

Status return API returns `BUFFER_OK`(0) on success. Distinguish between errors INVALID, BOUNDS, OVERFLOW, and ALLOC from `std::buffer::error`. The number is not interpreted as OS errno. `buffer_new` represents an allocation failure as an empty Buffer, so when you need to distinguish between failures, use `buffer_init`.

## Running example

<!-- wave-example: buffer-api -->
```wave
import("std::buffer::alloc")::{
    Buffer, buffer_init, buffer_free
};
import("std::buffer::write")::{
    buffer_append_str, buffer_push
};
import("std::buffer::read")::{
    buffer_get
};

fun main() -> i32 {
    var data: Buffer;
    if (buffer_init(&data, 0) < 0) {
        return 1;
    }
    if (buffer_append_str(&data, "Hi") < 0 || buffer_push(&data, 33) < 0) {
        buffer_free(&data);
        return 2;
    }

    var value: u8 = 0;
    if (buffer_get(&data, 2, &value) < 0) {
        buffer_free(&data);
        return 3;
    }

    println("{} {}", data.len, value);
    if (buffer_free(&data) < 0) {
        return 4;
    }

    return 0;
}
```

Execution result:

```text
3 33
```

Save it as `main.wave` and run it as `wavec run main.wave`. An initial capacity of 0 is not a failure, but a valid empty buffer. Frees up space during further processing.

## Lifespan and Failure

Growing the buffer can change data. Do not use a previously borrowed address after an operation that may reallocate. Copying the Buffer structure does not duplicate its allocation, so give that allocation a single owner.

`buffer_get` does not change the output arguments if it fails. On the other hand, the convenience function `buffer_at` also represents errors as 0, so use `buffer_get` to distinguish between actual 0 bytes and failures. Avoid creating invalid len/cap by directly changing public fields.

[Memory API](/docs/en/reference/memory-and-buffer) · [Practice reading a file as Buffer](/docs/en/practice/file-reader)

## Observe length and capacity separately

reserve frees up storage space, but does not increase len. resize changes the actual length used and initializes the extended portion to the specified byte. clear sets only the length used to 0, allowing allocation to be reused.

Save the following program as main.wave and run it. It doesn't rely on the exact growth multiple of capacity; it just ensures that you have the space you need.

<!-- wave-example: book-buffer-length-capacity -->
```wave
import("std::buffer::alloc")::{
    Buffer,
    buffer_init,
    buffer_reserve,
    buffer_resize,
    buffer_clear,
    buffer_free
};

fun main() -> i32 {
    var data: Buffer;

    if (buffer_init(&data, 0) < 0) {
        return 1;
    }

    if (buffer_reserve(&data, 20) < 0) {
        buffer_free(&data);
        return 2;
    }

    println("reserved length={}", data.len);

    if (buffer_resize(&data, 3, 7) < 0) {
        buffer_free(&data);
        return 3;
    }

    println("resized length={} first={}", data.len, deref data.data[0]);

    if (buffer_clear(&data) < 0) {
        buffer_free(&data);
        return 4;
    }

    println("cleared length={}", data.len);

    if (buffer_free(&data) < 0) {
        return 5;
    }

    return 0;
}
```

Execution result:

```text
reserved length=0
resized length=3 first=7
cleared length=0
```

When increasing to resize, we passed value=7, so all three new bytes we see are 7. Space only secured with reserve is not read as initialized data. cap will remain after clear and can be added again to the same Buffer.

## Distinguish between zero bytes and lookup failures

buffer_get returns the status and writes the actual bytes as output arguments. Even if the data is 0, it is a normal success.

<!-- wave-example: book-buffer-zero-bounds -->
```wave
import("std::buffer::alloc")::{Buffer, buffer_init, buffer_free};
import("std::buffer::write")::{buffer_push};
import("std::buffer::read")::{buffer_get};
import("std::buffer::error")::{BUFFER_ERR_BOUNDS};

fun main() -> i32 {
    var data: Buffer;

    if (buffer_init(&data, 0) < 0) {
        return 1;
    }

    if (buffer_push(&data, 0) < 0) {
        buffer_free(&data);
        return 2;
    }

    var value: u8 = 99;

    if (buffer_get(&data, 0, &value) < 0) {
        buffer_free(&data);
        return 3;
    }

    println("stored={}", value);
    value = 99;

    if (buffer_get(&data, 1, &value) == BUFFER_ERR_BOUNDS) {
        println("outside, preserved={}", value);
    }

    if (buffer_free(&data) < 0) {
        return 4;
    }

    return 0;
}
```

Execution result:

```text
stored=0
outside, preserved=99
```

The first hit is a success reading 0, the second hit is an out-of-bounds failure. Even if value=99 remains after a failure, it does not mean that it is the value read from the buffer. Be sure to check the status together.

## Practice solution: Byte accumulation

To add numbers 0 through 9, repeat buffer_push and check each result. Store the sum in i64 and read only the range `0 <= index < data.len`. After handling the buffer, we call buffer_free on both the success and failure paths.

<!-- wave-example: book-buffer-sum -->
```wave
import("std::buffer::alloc")::{Buffer, buffer_init, buffer_free};
import("std::buffer::write")::{buffer_push};

fun main() -> i32 {
    var data: Buffer;

    if (buffer_init(&data, 0) < 0) {
        return 1;
    }

    for (var value: i32 = 0; value < 10; value += 1) {
        if (buffer_push(&data, value as u8) < 0) {
            buffer_free(&data);
            return 2;
        }
    }

    var total: i64 = 0;

    for (var index: i64 = 0; index < data.len; index += 1) {
        total += deref data.data[index] as i64;
    }

    println("sum={}", total);

    if (buffer_free(&data) < 0) {
        return 3;
    }

    return 0;
}
```

Execution result:

```text
sum=45
```
