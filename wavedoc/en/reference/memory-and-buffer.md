---
translation_set_id: memory-buffer
path: reference/memory-and-buffer
locale: en
group: stdlib
group_order: 1
order: 4
title: mem: Allocation, reallocation, and layout
summary: Describes size in bytes, allocation failures, reallocation boundaries, and freeing responsibilities.
---

## Allocation and deallocation

```text
std::mem::alloc
mem_alloc(size: i64) -> ptr<u8>
mem_alloc_zeroed(size: i64) -> ptr<u8>
mem_free(p: ptr<u8>, size: i64) -> i64
mem_realloc(old_ptr: ptr<u8>, old_size: i64, new_size: i64) -> ptr<u8>
```

Size is in bytes. An allocation of size 0 or less returns null. If allocation also fails for positive sizes, it may be null. Do not assume the initial contents of `mem_alloc`, but use `mem_alloc_zeroed` if zero initialization is required.

The caller owns each successful allocation and must pass its original size when freeing it. `mem_free(null, size)` returns 0. A non-null pointer paired with a nonpositive size is an error. Never access or free an allocation after it has already been freed.

## Classification in case of reallocation

|request|action|
| --- | --- |
|new size is positive and success|`min(old_size, new_size)` Copy byte and deallocate previous one|
|New allocation of positive size fails|null Return, maintain existing allocation|
|`old_ptr == null`, positive new size|behaves like a new assignment|
| `new_size == 0` |Attempts to free a valid previous allocation and returns null|
|old_size=0 for negative size or existing pointer|null Return|

A null result from reallocating to size zero does not establish whether freeing succeeded. Call `mem_free` directly if you need its status. After growing an allocation, initialize the newly added region yourself.

## Example of preserving an existing pointer

Below is the case where the new size is positive. Save it as `main.wave` and run it.

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

Execution result:

```text
7
```

If you overwrite it with `data = mem_realloc(...)` before confirming the failure, you may lose the existing address. If reallocation is successful, the old address and pointers pointing to it are unused.

## The size and alignment of a target type

```text
std::mem::layout
size_of<T>() -> u64
align_of<T>() -> u64
```

Both values are the layout of the compilation target, not the computer it is running on. `size_of` includes tail padding and does not generate or evaluate a value. When multiplying the number of elements by their size, we check for overflow. You can use `mem_size_mul_checked` and `mem_size_add_checked` of `std::mem::ops`.

`mem_copy` is used to copy a range that does not overlap, and `mem_move` is used to copy a range that may overlap. Neither can determine the actual allocation length from pointers alone, so the caller must guarantee bounds. A list of bytes of varying size can be managed with [Buffer](/docs/en/stdlib/buffer).

## Layout query example

Save it as main.wave and run it. The size and alignment of the target covered in the document, i32, are each 4 bytes, so `4 4` is output. Check the values ​​of other types, especially structures and pointers, by target.

<!-- wave-example: layout-api -->
```wave
import("std::mem::layout")::{
    size_of, align_of
};

fun main() {
    println("{} {}", size_of<i32>(), align_of<i32>());
}
```

## Check for overflow in size calculation

You need to check if `count * element_size` is valid before passing it to the assignment function. If you allocate a small space with the overflow value and write as much as the original number, it will go out of bounds.

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

Execution result:

```text
bytes=12
overflow rejected
```

Failure results are not written to the allocation size. You can't detect all overflows just by calculating them with regular arithmetic and seeing if the result is negative. To calculate the size that needs to be inspected, use the checked function from the beginning.

## For overlapping copies, mem_move

When moving part of the same array backwards, the input and output areas overlap. Do not pass overlapping ranges to mem_copy, use mem_move.

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

Execution result:

```text
1
1
2
3
4
```

The original first four bytes move one position to the right. Copying them forward manually can read values that have already been overwritten, accidentally producing all ones. An overlap-aware API handles the copy direction for you.

## Write a function that passes ownership

For functions that return memory, it is best to provide a return address on success and the size required to free it. If the caller has to guess the size, it is prone to incorrect release. If a function returns a borrowed address, it should not be freed by the caller and describes the lifetime of the original.

Across function boundaries, you should be able to trace `allocator → owner → deallocation`. A pointer variable’s name or type does not automatically determine ownership.
