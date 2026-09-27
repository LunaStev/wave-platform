---
translation_set_id: memory-model
path: language/explicit-memory-type-model
locale: en
group: language
group_order: 2
order: 9
title: 9. Pointers, modifying values, and lifetimes
summary: Learn address-of, dereferencing, modifying a value through a pointer, and dangling pointers.
---

## Distinguish between value and storage location

The integer 42 and the address where the integer is stored are different values. The pointer points to the storage location. Passing an address allows the function to read or change the caller's storage.

This chapter covers taking addresses, dereferencing, modifying the original value, pointer arithmetic, and lifetimes. Dynamic allocation is covered in the next chapter. Start with addresses of local variables and array elements.

## Taking addresses and dereferencing

<!-- wave-example: book-pointer-first -->
```wave
fun main() {
    var value: i32 = 42;
    var address: ptr<i32> = &value;

    println("value={}", value);
    println("through pointer={}", deref address);

    deref address = 99;
    println("changed={}", value);
}
```

Execution result:

```text
value=42
through pointer=42
changed=99
```

`&value` gets the address, `deref address` reads or writes the value of that address, and so on. Instead of storing 99 in address, 99 was written to the integer pointed to by address. The variable address itself continues to point to value.

`ptr<i32>` is a pointer type for accessing i32 storage. The type does not record a length or provide automatic deallocation.

## Changing the pointer itself

<!-- wave-example: book-pointer-reassign -->
```wave
fun main() {
    var first: i32 = 10;
    var second: i32 = 20;
    var selected: ptr<i32> = &first;

    selected = &second;
    deref selected = 25;

    println("{} {}", first, second);
}
```

Execution result:

```text
10 25
```

`selected = &second` stores another address in a pointer variable. The value of first does not change. Then, if you write the value as deref, second changes. Separating “change of address” and “change of value through address” into separate sentences will reduce confusion.

## Making a function change the original

<!-- wave-example: book-pointer-increment -->
```wave
fun increment(value: ptr<i32>) {
    deref value = deref value + 1;
}

fun main() {
    var count: i32 = 4;

    increment(&count);
    increment(&count);

    println("{}", count);
}
```

Execution result:

```text
6
```

The function receives the address of count, rather than its value, 4. Changing that storage also changes count in the caller. The function requires a valid, writable i32 address. Passing null violates this requirement.

Each function defines whether it accepts null. If it does not, the caller must provide a valid address. If it does, the function must include a path that handles null.

## Function to process null

<!-- wave-example: book-pointer-null -->
```wave
fun try_increment(value: ptr<i32>) -> bool {
    if (value == null) {
        return false;
    }

    deref value = deref value + 1;
    return true;
}

fun main() {
    var count: i32 = 7;

    if (!try_increment(null)) {
        println("no value");
    }

    if (try_increment(&count)) {
        println("count={}", count);
    }
}
```

Execution result:

```text
no value
count=8
```

A null check handles only the absence of an address. Converting an arbitrary non-null number to a pointer does not create valid memory. Reading and writing also require a valid lifetime, sufficient size, correct alignment, and the appropriate access permissions.

## Array addresses and element-wise pointer arithmetic

<!-- wave-example: book-pointer-array -->
```wave
fun main() {
    var values: array<i32, 3> = [10, 20, 30];
    var first: ptr<i32> = &values[0];
    var second: ptr<i32> = first + 1;

    println("{}", deref first);
    println("{}", deref second);
    println("{}", deref first[2]);
}
```

Execution result:

```text
10
20
30
```

Adding 1 to the pointer moves it one element of the target type. The next element in i32 and the next element in u8 have different numbers of shift bytes. If you multiply `first + 1` again by the type size and add it, it will move to an unwanted position.

Pointer indexing must also be done within the valid range. first does not remember the array length 3 itself, so when passing the range to the function, it uses a form that receives both a pointer and length.

## Passing the read range to a function

<!-- wave-example: book-pointer-range -->
```wave
fun sum(values: ptr<i32>, count: i32) -> i32 {
    var total: i32 = 0;

    for (var index: i32 = 0; index < count; index += 1) {
        total += deref values[index];
    }

    return total;
}

fun main() {
    var values: array<i32, 4> = [2, 4, 6, 8];

    println("first two={}", sum(&values[0], 2));
    println("all={}", sum(&values[0], 4));
}
```

Execution result:

```text
first two=6
all=20
```

The unit of count is the number of elements. This function requires the caller to have count readable i32. Passing a length greater than the actual array violates the contract. Even if you use the same ptr<u8>·i64 combination, API, you should check in the documentation whether the length is in bytes or in elements.

## Lifespan: How long is the address valid?

Local variables are used within the lifetime of the call and block. If you return the address of a local variable inside a function for the caller to read later, that storage space may have reached the end of its lifetime.

Here's a piece of bad design that shouldn't be implemented:

```wave
fun invalid_address() -> ptr<i32> {
    var local: i32 = 42;

    return &local;
}
```

If one value is required, it returns i32. If it needs to write to storage space provided by the caller, it takes a pointer as input. If you need separate storage to retain beyond the call, allocate it explicitly and pass responsibility for freeing it.

## Two pointers pointing to the same storage space

Copying a pointer creates another name pointing to the same address. Does not duplicate memory.

<!-- wave-example: book-pointer-alias -->
```wave
fun main() {
    var value: i32 = 1;
    var first: ptr<i32> = &value;
    var second: ptr<i32> = first;

    deref second = 9;

    println("{} {}", value, deref first);
}
```

Execution result:

```text
9 9
```

The results changed through second can also be seen through first. Once the source memory is freed, both pointers become unusable. Assigning null to one pointer variable does not automatically change the other copies.

## Exercise: Exchange two integers

Write a function that takes two i32 addresses and exchanges their values. The first value must be stored in a temporary variable before overwriting it. Check if the value is retained even if you pass the same address twice.

### Complete solution

<!-- wave-example: book-pointer-swap -->
```wave
fun swap(left: ptr<i32>, right: ptr<i32>) {
    var saved: i32 = deref left;

    deref left = deref right;
    deref right = saved;
}

fun main() {
    var first: i32 = 3;
    var second: i32 = 8;

    swap(&first, &second);
    println("{} {}", first, second);

    swap(&first, &first);
    println("{}", first);
}
```

Execution result:

```text
8 3
8
```

This function also requires both addresses to point to valid, writable integer storage. To handle null, add a result indicating success or failure, as in try_increment.


## **Wave Explicit Memory Type Model**

The pointer design of Wave is based on **Wave Explicit Memory Type Model**. This model defines pointers and arrays as explicit memory types at the language level, rather than syntactic tricks or library abstractions.

`ptr<T>` is a type that points to the memory address that stores the `T` value, and `array<T, N>` is a fixed-length memory type that stores `N` values of `T` in succession. Therefore, the structure of pointers and arrays is revealed as is in function arguments, return values, structure fields, and other types.

## null

```wave
var buffer: ptr<u8> = null;
if (buffer == null) {
    println("no buffer");
}
```

`null` is a pointer value that does not point to a valid memory address. `null` can only be assigned to the `ptr<T>` type and cannot be used as an integer, Boolean, or array value.

An allocation or lookup function can return `null` when it has no result. Check for `null` before dereferencing such a result. Dereferencing a `null` pointer does not access valid storage.

## pointer conversion

When you need to change the address or other pointer representation, use `as`.

```wave
var raw: i64 = 0;
var p: ptr<u8> = raw as ptr<u8>;
```

Use conversions between integers and pointers only on low-level boundaries, and consider the address width of the target platform and ABI.
