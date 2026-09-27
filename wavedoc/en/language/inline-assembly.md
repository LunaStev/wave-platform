---
translation_set_id: assembly
path: language/inline-assembly
locale: en
group: language
group_order: 2
order: 19
title: Inline assembly
summary: Describes the contract of the asm block's command string, in/out operands, and clobber.
---

## asm block

`asm` is a low-level syntax for directly inserting instructions of the target architecture.

```wave
fun read_value() -> i64 {
    var result: i64 = 0;
    asm {
        "mov rax, 123"
        out("rax") result
    }

    return result;
}
```

String literals within a block are passed as an assembly instruction list.

## input and output

```wave
var result: i64 = 0;
asm {
    "mov rax, rdi"
    in("rdi") 123
    out("rax") result
}
```

- `in("reg") expression` wires the value Wave to the input operand.
- `out("reg") target` writes the output value to the assignable Wave target.
- Register names can be written as strings or identifiers.

Input operands can include variables, integer/string literals, `&identifier`, `deref identifier`, and negative numbers.

## clobber

If a block changes a register or memory state other than an explicit output, it is recorded in `clobber(...)`.

```wave
asm {
    "nop"
    clobber("rax", "rcx", "memory")
}
```

## Check when using

- The instruction syntax must match the target architecture and the LLVM inline assembly contract.
- Do not arbitrarily destroy registers that must be preserved according to the calling convention.
- For blocks that read or write memory, declare clobber, including `memory`.
- If possible, isolate architecture-specific asm behind a small function.

The behavior and portability of inline assembly are not guaranteed by language type alone.

## Learning and Example Range

[Practice with the full program](/docs/en/getting-started/overview) · [Standard library](/docs/en/stdlib)
