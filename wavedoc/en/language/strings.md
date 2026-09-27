---
translation_set_id: learn-strings
path: language/strings
locale: en
group: language
group_order: 2
order: 7
title: 7. Strings, characters, and bytes
summary: Distinguish between string and char, UTF-8 byte length, NUL, search and binary data.
---

## letters on screen and bytes in memory

You see letters on the screen, but bytes are stored in memory. Especially in cases where one character has multiple UTF-8 bytes, such as in Korean, it is easy to make a mistake if “length” and “number of characters” are used interchangeably.

This chapter distinguishes between str and char, ending NUL, escape, search location and binary data. The examples are each a complete program.

## String literals and output

<!-- wave-example: book-string-literals -->
```wave
fun main() {
    var greeting: str = "안녕하세요";

    println("{}", greeting);
    println("line one\nline two");
    println("quote: \"Wave\"");
}
```

Execution result:

```text
안녕하세요
line one
line two
quote: "Wave"
```

Regular characters within double quotation marks are expressed as UTF-8. escape indicates bytes that are difficult to write directly from the source. `\n` is a line break byte and does not print two characters, a backslash and n. To print the backslash itself, use `\\`.

## Length is number of bytes

<!-- wave-example: book-utf8-length -->
```wave
import("std::string::len")::{
    len
};

fun main() {
    println("ASCII={}", len("Wave"));
    println("Korean={}", len("한"));
    println("mixed={}", len("Wave한"));
}
```

Execution result:

```text
ASCII=4
Korean=3
mixed=7
```

`len` counts bytes before the terminating NUL, not visible characters. The character `한` takes three bytes in UTF-8. Character counts, Unicode code point counts, and byte counts are not generally interchangeable. Display width also depends on factors such as fonts and combining characters.

Therefore, functions that cut a string at an arbitrary byte position and display it on the screen must separately consider the Unicode boundary. Clearly define the input conditions, whether it is a program that only processes ASCII or general Unicode text.

## NUL End and Length

str uses a 0 byte to indicate the end. len does not include that last byte in the length. It is an error to put NUL inside a string literal. Here is an example of an intentional error:

```wave
fun main() {
    var text: str = "left\x00right";
}
```

`\xNN` in the source specifies one byte with exactly two hexadecimal digits. `\x41` represents byte A which is 65. Because writing a regular character as UTF-8 and inserting an arbitrary byte are different, not all str that can be written as `\xNN` are valid UTF-8.

## char does not contain the entire character Unicode

char is an unsigned 8-bit character value. You can use literals that represent values ​​in a single byte range, such as `'A'`. `'한'` is an error because it does not fall within this range. `"한"` is a separate str with multiple UTF-8 bytes.

Don't try to always put one letter in one char. You must first decide whether the units needed to process text are bytes or Unicode code points.

## string comparison

To compare string contents, use the function std. Below is a program that checks for identical content and case differences.

<!-- wave-example: book-string-compare -->
```wave
import("std::string::cmp")::{
    eq, starts_with, ends_with
};

fun main() {
    if (eq("Wave", "Wave")) {
        println("same bytes");
    }

    if (!eq("Wave", "wave")) {
        println("case differs");
    }

    if (starts_with("report.txt", "report") && ends_with("report.txt", ".txt")) {
        println("name matches");
    }
}
```

Execution result:

```text
same bytes
case differs
name matches
```

This comparison compares byte strings. It does not automatically perform language-specific case conversion or Unicode normalization. Even when comparing file names, the file name equality rules in OS are not the same as simple string comparisons.

## Units and failures of search results

<!-- wave-example: book-string-search -->
```wave
import("std::string::find")::{
    find, count
};

fun main() {
    var first: i32 = find("banana", "na");
    var missing: i32 = find("banana", "xy");
    var matches: i32 = count("aaaa", "aa");

    println("first={} missing={}", first, missing);
    println("matches={}", matches);
}
```

Execution result:

```text
first=2 missing=-1
matches=2
```

find returns the first position or -1. Index 0 is also a success, so it is checked with `result >= 0`. count is not a position, but the number of non-overlapping matches. Above, aa is number 2 as it matches 0~1 and 2~3.

Bin needle is also part of the contract. find returns 0, contains returns true, and count returns 0. Don't assume that just because the function name is in the same module, even the return method is the same.

## Removing spaces is different from creating a new string

trim_range returns the range excluding spaces without modifying or copying the original text. Since we receive an output pointer, we first prepare an integer to store the result.

<!-- wave-example: book-string-trim-range -->
```wave
import("std::string::trim")::{
    trim_range
};

fun main() {
    var start: i32 = 0;
    var end: i32 = 0;

    trim_range("  Wave  ", &start, &end);

    println("start={} end={} bytes={}", start, end, end - start);
}
```

Execution result:

```text
start=2 end=6 bytes=4
```

The range is `[start, end)`. It includes the start but not the end, so its length is end-start. Adding start to the starting address of the original does not automatically create NUL at end location. You will need to carry the range separately or prepare a new string space.

## Binary data has a separate length

Data containing zeros is not covered by the end-of-string rules. It uses byte arrays and lengths.

<!-- wave-example: book-binary-not-string -->
```wave
fun main() {
    var bytes: array<u8, 3> = [65, 0, 66];

    for (var index: i32 = 0; index < 3; index += 1) {
        println("{}", bytes[index]);
    }
}
```

Execution result:

```text
65
0
66
```

The second 0 is the actual data. If you interpret this as str, it is treated as ending at the first 0 and you cannot see the following 66. Conversely, if you change an array without NUL to str with cast, there is a risk of reading beyond the array. cast is not an operation to add a termination byte.

## Exercise: Examining file names

Check if the file name ends with `.wave`, and if the string contains `test`, output it to a test file. This exercise only checks for byte patterns in the name and does not address actual file existence.

### Complete solution

<!-- wave-example: book-string-solution -->
```wave
import("std::string::cmp")::{
    ends_with
};
import("std::string::find")::{
    contains
};

fun classify(name: str) {
    if (!ends_with(name, ".wave")) {
        println("other file");
        return;
    }

    if (contains(name, "test")) {
        println("Wave test file");
    } else {
        println("Wave source file");
    }
}

fun main() {
    classify("main.wave");
    classify("parser_test.wave");
    classify("notes.txt");
}
```

Execution result:

```text
Wave source file
Wave test file
other file
```

There are separate policies regarding how to handle the capital letter `.WAVE` and whether to consider it as a test even if test is included in the entire path. Even if it is a small function, its operation can be accurately described only if it is determined what input it is targeting.
