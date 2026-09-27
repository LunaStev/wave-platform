---
translation_set_id: lexical
path: language/lexical-structure
locale: en
group: language
group_order: 2
order: 14
title: Lexical structure
summary: Describes identifiers, literals, delimiters, keywords, and type names.
---

## identifier

Identifiers name variables, functions, types, and fields. The name is case sensitive and can contain any combination of letters, numbers, and `_`. Numbers cannot be used in the first letter. The characters Unicode can also be used in identifiers.

```wave
var request_count: i64 = 0;
var 이름: str = "Wave";
```

In real projects, it is recommended to use a consistent naming convention for tool compatibility and searchability.

## Sentences and Separators

Most declaration and expression statements end with `;`. Statements with a body, such as functions, conditional statements, loop statements, and structures, use the `{ ... }` block.

```wave
var answer: i32 = 42;

fun double(value: i32) -> i32 {
    return value * 2;
}
```

## literal

```wave
var integer: i32 = 42;
var decimal: f64 = 3.14;
var text: str = "Wave";
var letter: char = 'W';
var enabled: bool = true;
var address: ptr<u8> = null;
```

You can use integers, floating point numbers, strings, characters, booleans, and `null` literals. Use `null` for pointer values.

## Keywords and Type Names

The main keywords used in the Wave grammar are as follows.

`pub`, `fun`, `extern`, `export`, `type`, `enum`, `static`, `var`, `deref`, `const`, `if`, `else`, `proto`, `struct`, `while`, `for`, `in`, `out`, `clobber`, `as`, `asm`, `import`, `return`, `continue`, `print`, `input`, `println`, `match`, `break`, `true`, `false`, `null`.

Built-in type names include `bool`, `char`, `byte`, `str`, integer and floating point types, `ptr` and `array`. Pointers are written in the form `ptr<T>`, and fixed-length arrays are written in the form `array<T, N>`.

## Strings and Characters escape

|notation|meaning|
| --- | --- |
| `\n` |LF Line break|
| `\r` | CR |
| `\t` |tab|
| `\\` |Backslash|
| `\"` |double quotation marks|
| `\xNN` |One byte specified as exactly two hexadecimal digits|

General string characters are saved as UTF-8. Since `\xNN` preserves one byte, there is no guarantee that the entire string is a valid UTF-8. NUL (including `\x00`) inside a string literal is a compilation error. For data containing zeros, use a byte array and length.

`char` The literal must fit into an 8-bit value. Characters that exceed that range, such as `'한'`, are errors. It is different from the string `"한"`.

LF, CRLF, and CR alone in the source are each treated as one logical line break. These are rules about source location and comment termination and do not mean that they change the actual bytes of the file data.

Additional grammar names include `variant`, `async`, and `await`, and the asynchronous value is expressed as `Future<T>`. The above independent `var` declaration block is a code fragment inside a function.

[string class](/docs/en/language/arrays) · [Comments](/docs/en/language/comments)

## Example of intentional failure

If you run the program below check, an internal NUL error should occur. If you need 0 bytes, use the `[97, 0, 98]` byte array.

<!-- wave-example: reject-string-nul -->
```wave
fun main() {
    var text: str = "a\x00b";
}
```

The character literal below also exceeds the 8-bit range, so it is a compilation error. To represent the string UTF-8, use `str` and double quotation marks.

<!-- wave-example: reject-wide-char -->
```wave
fun main() {
    var letter: char = '한';
}
```
