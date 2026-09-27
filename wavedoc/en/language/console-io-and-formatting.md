---
translation_set_id: console-io-formatting
path: language/console-io-and-formatting
locale: en
group: language
group_order: 2
order: 16
title: Console input, output, and formatting
summary: print, println, input Describes sentences and placeholder rules.
---

## input/output statement

Wave provides `print`, `println`, and `input` as console input/output statements.

```wave
fun main() {
    var count: i32 = 0;
    input("{}", count);
    print("count = ");
    println("{}", count);
}
```

Each sentence ends with `;`. The first argument must be a string literal. Variables or calculated strings cannot be used as format arguments.

## placeholder

Only exactly two characters, `{}`, are placeholders.

```wave
println("name = {}, score = {}", name, score);
```

The number of placeholders and the number of expressions that follow must be exactly the same. If the numbers are different, it is a grammar error.

```wave
println("{} {}", one);
// 오류: 자리표시자 2개, 값 1개
println("plain text", one);
// 오류: 자리표시자 없음, 값이 남음
```

Other forms of curly braces are left as plain text. There are no named or numbered placeholders in this grammar.

## print and println

`print` prints the formatted text as is, and `println` adds a line break.

```wave
print("loading...");
println("done");
```

Formatting arguments use scalar values such as integers, floating point numbers, strings, and pointers. Arrays and structures cannot be used as formatting arguments.

## input Target

`input` stores the read value in the destination, so all expressions after format must be writable locations.

```wave
var number: i32 = 0;
input("{}", number);
```

Variables, fields and dereferenced storage locations can be used as targets. Literals and calculation results cannot be used as input.

If all input values cannot be converted to the requested type, the program exits with a failure status.

## runtime boundary

These statements use console input and output from the hosted environment. In a freestanding environment, input/output provided by the kernel or device must be defined as a function or FFI boundary.

## Input value and range

bool input only accepts `0` and `1`. It does not interpret 2 as true or accept the string `true` as the same input. Integer input must be within the range of the target integer width. 128, 256, 512, and 1024-bit integers are also processed based on the overall width of the type.

Format error, out of range, before required input EOF fails. The built-in input is not a function that returns failure and re-enters, but is an input function that terminates the process on failure. If you need recoverable input processing, read the bytes with io and construct a separate parser.

[Input Calculator Practice](/docs/en/practice/input-calculator) · [File and io](/docs/en/stdlib/files-io)
