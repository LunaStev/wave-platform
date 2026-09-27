---
translation_set_id: string-bytes
path: reference/string-and-bytes
locale: en
group: stdlib
group_order: 1
order: 3
title: string: Length, searching, and ranges
summary: NUL Describes the byte unit of the termination string API and the return value.
---

## String storage and argument conditions

The `str` argument to this module must be an accessible NUL termination byte. Length and search index are in bytes. Regular characters are stored as UTF-8, but byte search is Unicode without normalization or character-by-character splitting. Do not assume that the returned index is a character boundary.

## Compare with length

```text
std::string::len
len(s: str) -> i32
is_empty(s: str) -> bool

std::string::cmp
eq(a: str, b: str) -> bool
cmp(a: str, b: str) -> i32
starts_with(s: str, prefix: str) -> bool
ends_with(s: str, suffix: str) -> bool
```

`len` excludes the last NUL. `cmp` The order is judged by the sign of the result. The return value is not interpreted as Unicode character order or language-specific dictionary sorting. These functions do not allocate memory and do not change their input.

## search

Get the name you need, like `import("std::string::find")::{find, contains, count};`.

|function declaration|result|
| --- | --- |
| `find(s: str, needle: str) -> i32` |First match location. -1 if not present, 0 for empty needle|
| `contains(s: str, needle: str) -> bool` |Included or not. Empty needle is true|
| `count(s: str, needle: str) -> i32` |Number of non-overlapping matches. Bin needle is 0|
| `find_char(s: str, c: u8) -> i32` |first position of byte or -1|
| `rfind_char(s: str, c: u8) -> i32` |Last position of byte or -1|
| `contains_char(s: str, c: u8) -> bool` |Existence of that byte|
| `count_char(s: str, c: u8) -> i32` |the number of bytes in question|

`c` in the name `*_char` is a byte, not a Unicode code point. The NUL itself at the end of the string is not included in the search target.

## Range excluding spaces

```text
std::string::trim
trim_left_index(s: str) -> i32
trim_right_index(s: str) -> i32
trim_range(s: str, out_start: ptr<i32>, out_end: ptr<i32>)
```

`trim_range` writes the semi-open range `[start, end)` excluding spaces ASCII to the output argument. Both output pointers must point to writable integers. It does not modify the original text or create new strings. If everything is blank, it becomes an empty range.

## Running example

Save it to `main.wave` and run `wavec run main.wave`.

<!-- wave-example: string-api -->
```wave
import("std::string::find")::{
    find, count
};
import("std::string::trim")::{
    trim_range
};

fun main() {
    var start: i32 = 0;
    var end: i32 = 0;
    trim_range("  Wave  ", &start, &end);
    println("{} {}", start, end);
    println("{} {}", find("banana", "na"), count("aaaa", "aa"));
}
```

Execution result:

```text
2 6
2 2
```

## Related features

The classification/case conversion of `std::string::ascii` is for the range ASCII. `djb2_32` and `fnv1a_64` of `std::string::hash` are not used for cryptographic hashes or password storage. For data containing NUL, use [bytes](/docs/en/stdlib/bytes).

## Patterns that overlap with empty search terms

Seeing the edge behavior of the search function with actual values makes it easier to determine the calling conditions.

<!-- wave-example: book-string-empty-search -->
```wave
import("std::string::find")::{find, contains, count};

fun main() {
    println("empty find={}", find("Wave", ""));
    println("empty count={}", count("Wave", ""));
    println("nonoverlapping={}", count("aaaa", "aa"));

    if (contains("Wave", "")) {
        println("empty needle is contained");
    }
}
```

Execution result:

```text
empty find=0
empty count=0
nonoverlapping=2
empty needle is contained
```

The success value of find, 0, is the first position. A success value of 0 for count is the result of a no match or an empty search term rule. No two values ​​are treated equally. If case ignoring or Unicode normalization is required, separate policies must be implemented before and after this byte search.

## trim Copying range to new string

There is no new ending NUL in the range returned by trim_range. When copying to a separate destination, reserve length + 1 space and write the last byte directly as 0.

<!-- wave-example: book-trim-copy -->
```wave
import("std::string::trim")::{trim_range};

fun main() -> i32 {
    var original: str = "  Wave  ";
    var start: i32 = 0;
    var end: i32 = 0;
    var destination: array<u8, 16>;

    trim_range(original, &start, &end);

    var length: i32 = end - start;

    if (length >= 16) {
        return 1;
    }

    for (var index: i32 = 0; index < length; index += 1) {
        destination[index] = original[start + index];
    }

    destination[length] = 0;
    println("{}", &destination[0] as str);
    return 0;
}
```

Execution result:

```text
Wave
```

A string of length 16 will not fit into this destination. This is because you need the last NUL. Even if the length is 0, writing destination[0]=0 results in a valid empty string. The destination local array lives until the end of main, so we print within it.

## String API Order of use

When designing a string API, specify whether its input is NUL-terminated, whether indexes count bytes, and whether the result borrows a source range or owns a new allocation. A borrowed range depends on the source’s lifetime. An allocated result must specify who frees it.

Read [String Learning Chapter](/docs/en/language/strings) for basic concepts and [bytes](/docs/en/stdlib/bytes) for data including NUL.
