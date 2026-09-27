---
translation_set_id: quick-reference
path: reference/syntax-quick-reference
locale: en
group: reference
group_order: 5
order: 3
title: Syntax quick reference
summary: Frequently used declarations, control flow, types, pointers, and FFI grammar are organized on one page.
---

## declaration

```wave
var value: i32 = 1;
var next: i32 = value + 1;
const LIMIT: i32 = 64;
static total: i64 = 0;
type Identifier = u64;
```

`var` is the region, `const`/`static` are top-level declarations. Local variables explicitly declare their type.

## Functions

```wave
fun max(left: i32, right: i32) -> i32 {
    if (left > right) {
        return left;
    }

    return right;
}
```

## generic

```wave
fun identity<T>(value: T) -> T {
    return value;
}

var value: i32 = identity<i32>(10);
```

When calling a generic, specify a type argument.

## Structure and enum

```wave
struct Pair {
    left: i32;
    right: i32;
}

enum Result -> i32 {
    Ok = 0,
    Error,
}
```

## Conditions and loops

```wave
if (ready) {
    println("ready");
}

while (count < 10) {
    count += 1;
}

for (var i: i32 = 0; i < 10; i += 1) {
    println("{}", i);
}

match (status) {
    Ready => {
        println("ready");
    }
    0 => {
        println("zero");
    }
    _ => {
        println("other");
    }
}
```

The headers of `if`, `while`, `for`, and `match` use parentheses.

## Arrays and pointers

```wave
var values: array<i32, 4> = [1, 2, 3, 4];
var p: ptr<i32> = &values[0];
var first: i32 = deref p;
```

## console input/output

```wave
print("value = ");
println("{}", value);
input("{}", value);
```

The first argument is a string literal. Each exact `{}` placeholder requires an expression following it, and the `input` target must be assignable.

## import and public items

```wave
import("std::string::len");
import("./helpers" as helpers);
import("math")::{
    Vector
};

pub fun add(left: i32, right: i32) -> i32 {
    return left + right;
}

pub import("./extra"):: {
    increment
};
```

The local path starts with `./`. The alias import specifies the module name, and the select import imports the required public entries into this file's namespace. `pub import` re-exports the selected items.

## FFI

```wave
extern(c) fun native_call(value: i32) -> i32;

export(c) fun wave_call(value: i32) -> i32 {
    return value + 1;
}
```

## Target conditional item

```wave
#[target(os="linux", arch="riscv64")]
extern(c) fun platform_call(value: i32) -> i32;
```

Support condition keys are `arch`, `os`, `env`, `abi`. Properties control the next top-level item.

## Inline assembly

```wave
var result: i64 = 0;
asm {
    "mv a0, a1"
    in("a1") 7
    out("a0") result
    clobber("memory")
}
```

Instruction text and register names are target dependent. Declare all inputs, outputs and hidden clobber required for the block.

## source inspection

```shell
wavec build main.wave --emit=check
wavec print supported-targets
wavec print supported-emit-kinds
```

## Learning and Example Range

Examples of local variables and statements separately displayed outside the function are code fragments inserted into the body of the function. Complete running examples and exercises follow in [Wave Learning Process](/docs/en/getting-started/overview). Please check [Standard library](/docs/en/stdlib) for detailed rules of memory and external functions.
