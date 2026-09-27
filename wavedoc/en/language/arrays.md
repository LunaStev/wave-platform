---
translation_set_id: learn-arrays-strings
path: language/arrays
locale: en
group: language
group_order: 2
order: 6
title: 6. Arrays and iteration
summary: Learn fixed-size arrays, indexing, traversal, copying, and searching.
---

## Multiple values of the same type

If you create the three scores separately as score1, score2, and score3, both declarations and calculations must be changed when the number changes. Arrays group together a set number of elements of the same type. Using loops, you can apply the same rules to each element.

This chapter covers creating, indexing, modifying, iterating over, searching, and aggregating arrays. Strings also support indexing, but its meaning differs, so strings are covered in the next chapter.

## Enter length in type

<!-- wave-example: book-array-create -->
```wave
fun main() {
    var scores: array<i32, 3> = [70, 80, 90];

    println("first={}", scores[0]);
    println("second={}", scores[1]);
    println("last={}", scores[2]);
}
```

Execution result:

```text
first=70
second=80
last=90
```

i32 in `array<i32, 3>` is the element type and 3 is the number of elements. What we store is 3 integers. It doesn't mean the number of bytes is 3. The number of elements in an array literal must match the declared length.

Indexes start from 0. The first element is 0, the last element is length-1. scores[3] is an out-of-scope access, not a third element.

## change element

<!-- wave-example: book-array-update -->
```wave
fun main() {
    var scores: array<i32, 3> = [70, 80, 90];

    scores[1] = 85;
    scores[0] += 5;

    println("{} {} {}", scores[0], scores[1], scores[2]);
}
```

Execution result:

```text
75 85 90
```

Change the storage of specific elements without recreating the entire array. An index expression can also be the result of a calculation, but you must ensure that the value is within a range. When using an external input as an index, both the negative number and the upper limit are checked.

## Iterating over an array

<!-- wave-example: book-array-sum -->
```wave
fun main() {
    var scores: array<i32, 4> = [60, 70, 80, 90];
    var total: i32 = 0;

    for (var index: i32 = 0; index < 4; index += 1) {
        total += scores[index];
    }

    println("total={} average={}", total, total / 4);
}
```

Execution result:

```text
total=300 average=75
```

Each iteration reads an element at a different index. The sum variable must be initialized outside of the iteration. Initializing it to 0 each time within the loop body will produce incorrect results, such as leaving only the last element.

The integer division used to calculate the average discards the fractional part. For a floating-point average, convert the sum before dividing. For larger arrays or larger values, also ensure that the accumulator type can represent the sum.

## Aggregating only some elements

Filtering can be achieved by combining conditional statements and traversal. Here we count the number of elements with a score of 80 or higher.

<!-- wave-example: book-array-filter -->
```wave
fun main() {
    var scores: array<i32, 5> = [60, 80, 90, 75, 100];
    var passed: i32 = 0;

    for (var index: i32 = 0; index < 5; index += 1) {
        if (scores[index] >= 80) {
            passed += 1;
        }
    }

    println("passed={}", passed);
}
```

Execution result:

```text
passed=3
```

Index and element values must be separated. Examining `index >= 80` compares positions, not scores. Both can be i32, so it is difficult to find this semantic error based on type alone.

## Find first match location

First, decide how to display the results that were not found. In this example, valid indices are 0 to 4, so we use -1 as the failure marker.

<!-- wave-example: book-array-find -->
```wave
fun main() {
    var values: array<i32, 5> = [8, 3, 8, 1, 5];
    var target: i32 = 8;
    var found: i32 = -1;

    for (var index: i32 = 0; index < 5; index += 1) {
        if (values[index] == target) {
            found = index;
            break;
        }
    }

    if (found >= 0) {
        println("found at {}", found);
    } else {
        println("not found");
    }
}
```

Execution result:

```text
found at 0
```

The first position 0 is also a normal result. If you check for success with `found > 0`, you will be mistaken for not finding the first element. For the same reason, judging success by using bool as cast is wrong.

If you remove break, subsequent matches will overwrite found, giving you the position of the last match. Since a single statement can change the contract of a function, the description of “search” must also be specifically written as to whether it is in the first or last position.

## Copying array elements

To copy the values of an array to another storage space, you can read and assign them element by element. Even if you change one integer element after copying, the other integer element does not change.

<!-- wave-example: book-array-copy -->
```wave
fun main() {
    var original: array<i32, 3> = [1, 2, 3];
    var copied: array<i32, 3>;

    for (var index: i32 = 0; index < 3; index += 1) {
        copied[index] = original[index];
    }

    copied[0] = 99;

    println("original={}", original[0]);
    println("copied={}", copied[0]);
}
```

Execution result:

```text
original=1
copied=99
```

If the elements are pointers, copying them copies their addresses. It does not duplicate the separate memory they point to. This distinction matters when managing ownership.

## Initialization and effective range

It does not assume that all elements of an array declared without an initial value can be read. If only a partial number is recorded, the actual initialized number must be managed separately. This is the same reason why the length returned by the library read function may be smaller than the entire buffer capacity.

Because the array length is contained in the type, it does not grow arbitrarily during execution. Lists of bytes that grow in size use dynamic storage such as `Buffer`. Changing the length of an array requires considering the type, initial value, traversal upper limit, and calculations that depend on that length.

## Exercise: Maximums and Locations

Find the maximum value and the position at which it first appears in the array `[4, 9, 2, 9, 1]`. If it is a regular function where all elements may be negative, the maximum value should not be initialized to 0.

### Complete solution

<!-- wave-example: book-array-solution -->
```wave
fun main() {
    var values: array<i32, 5> = [4, 9, 2, 9, 1];
    var maximum: i32 = values[0];
    var position: i32 = 0;

    for (var index: i32 = 1; index < 5; index += 1) {
        if (values[index] > maximum) {
            maximum = values[index];
            position = index;
        }
    }

    println("max={} first={}", maximum, position);
}
```

Execution result:

```text
max=9 first=1
```

Take the first element as the initial reference and compare from the second. Since `>`, the position does not change even if the same maximum value appears again. Change it to `>=` to become the last position. Any interface that can be zero-length must process an empty input before reading the first element.
