---
translation_set_id: learn-allocation
path: language/allocation
locale: en
group: language
group_order: 2
order: 10
title: 10. Memory allocation and resource management
summary: Learn about allocation failures, initialization, scope, and freeing.
---

## When do you need dynamic storage space?

A fixed-size array includes its length in its type. Use dynamic memory when the amount of data is known only at runtime, such as a file size or an input length. Free each allocation when you no longer need it.

In this chapter, you will manage a small allocation, resize it, and then use a Buffer. Passing a pointer is different from transferring ownership. Identify which resources each function owns as you follow the examples.

## Allocate, check, use, and free

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

Execution result:

```text
42 0
```

The program has four steps: request 4 bytes, check for null, access only the valid range, and free the allocation. On success, mem_alloc_zeroed initializes the memory to zero, so the second byte is zero even though the program has not written to it.

Do not assume any initial contents for memory returned by mem_alloc. Initialize each region before reading it. A zero or negative allocation size returns null. A positive-size allocation can also fail to obtain memory.

## unit of size

The size argument of a memory allocation API is measured in bytes. To allocate ten integers, multiply the element size by the number of elements. Check that this multiplication does not overflow.

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

Execution result:

```text
1 2 3
```

count is an element count; bytes is a byte count. Pointer arithmetic moves in units of i32, but freeing the allocation requires its original size in bytes. size_of uses the target type layout, making the relationship to the element type explicit.

When changing the result of size_of to i64 for a general very large type, the conversion range must also be considered. Here we use i32, which has a known size.

## Clean up even on failure paths

If another operation fails after allocation, free the memory before returning early. An ownership table helps identify cleanup paths you might otherwise miss.

|Step|Resources owned by|If you fail|
| --- | --- | --- |
|Before allocation|None|return it right away|
|After successful allocation|data and original size|data Return after release|
|After successful reallocation|New address and new size|New address release|
|After release|None|Do not use old address|

Overwriting a pointer variable and losing the original address also loses the information needed to free the allocation. This causes a memory leak. Conversely, freeing the same allocation through two owners causes a double free.

## Reallocation to increase size

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

Execution result:

```text
7 9
```

First store the result in next. If allocation of the larger block fails, the original data remains valid and can still be freed. On success, the old allocation is freed and the new address must be used. Initialize the newly added region before reading it.

The new size in this example is positive. A request with new_size=0 instead attempts to free the existing allocation and returns null. Therefore, a null result does not always mean that the old allocation is still valid. Call mem_free directly when you need to check whether freeing succeeded.

## When to Recheck Borrowed Pointers

data It is wrong to save an internal pointer and use it after reallocation. This is because the address of the new data may be different. If you need an internal location, you can store offset instead of the address and recalculate based on the new data after success.

Freeing or reallocating memory also affects code that has borrowed it. Check whether another operation is still using that memory. A buffer passed to an asynchronous operation must remain valid until the operation has finished.

## The byte list contains Buffer

Managing a list of bytes with frequently changing lengths while manually reallocating requires handling both len and cap, expansion failures, and size calculations. Buffer of std bundles these operations.

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

Execution result:

```text
bytes=11
```

len is the number of bytes in use; cap is the allocated capacity. Appending data grows the allocation when necessary. Using Buffer does not remove the caller’s responsibility to free it.

buffer_append_str does not add NUL to the end of the string. Therefore, message.data should not be output directly as str. Bytes are output with the I/O function, which takes the length, or explicitly constructs a string representation.

## Exercise and approach

Add the bytes 0 to 9 one by one to Buffer and get the sum. It must be freed when each append fails, and reads are only performed in the range len. The full solution and boundary failure can be checked by following the example in [Buffer How to use](/docs/en/stdlib/buffer).

Try marking allocation, reallocation, and deallocation calls in your code. For each successful allocation, you must be able to describe who owns it and which path frees it.
