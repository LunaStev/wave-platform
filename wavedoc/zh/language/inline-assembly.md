---
translation_set_id: assembly
path: language/inline-assembly
locale: zh
group: language
group_order: 2
order: 19
title: 内联装配
summary: 描述asm块的命令字符串、in/out操作数和clobber的约定。
---

## asm块

`asm`是直接插入目标架构指令的低级语法。

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

块内的字符串文字作为汇编指令列表传递。

## 输入和输出

```wave
var result: i64 = 0;
asm {
    "mov rax, rdi"
    in("rdi") 123
    out("rax") result
}
```

- `in("reg") expression` 将值 Wave 连接到输入操作数。
- `out("reg") target` 将输出值写入可分配的Wave 目标。
- 寄存器名称可以写为字符串或标识符。

输入操作数可以包括变量、整数/字符串文字、`&identifier`、`deref identifier` 和负数。

## clobber

如果块更改寄存器或内存状态而不是显式输出，则会将其记录在`clobber(...)`中。

```wave
asm {
    "nop"
    clobber("rax", "rcx", "memory")
}
```

## 使用时检查

- 指令语法必须与目标体系结构和LLVM内联汇编合约匹配。
- 不要随意破坏根据调用约定必须保留的寄存器。
- 对于读取或写入内存的块，声明clobber，包括`memory`。
- 如果可能，将特定于架构的asm隔离在一个小函数后面。

内联汇编的行为和可移植性并不能仅由语言类型来保证。

## 学习和示例范围

[练习完整的程序](/docs/zh/getting-started/overview) · [标准库](/docs/zh/stdlib)
