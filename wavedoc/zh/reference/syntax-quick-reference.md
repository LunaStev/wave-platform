---
translation_set_id: quick-reference
path: reference/syntax-quick-reference
locale: zh
group: reference
group_order: 5
order: 3
title: 语法快速参考
summary: 常用的声明、控制流、类型、指针和FFI语法都组织在一页上。
---

## 声明

```wave
var value: i32 = 1;
var next: i32 = value + 1;
const LIMIT: i32 = 64;
static total: i64 = 0;
type Identifier = u64;
```

`var` 是区域，`const`/`static` 是顶级声明。局部变量显式声明其类型。

## 函数

```wave
fun max(left: i32, right: i32) -> i32 {
    if (left > right) {
        return left;
    }

    return right;
}
```

## 通用的

```wave
fun identity<T>(value: T) -> T {
    return value;
}

var value: i32 = identity<i32>(10);
```

调用泛型时，指定类型参数。

## 结构及enum

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

## 条件和循环

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

`if`、`while`、`for`和`match`的标头使用括号。

## 数组和指针

```wave
var values: array<i32, 4> = [1, 2, 3, 4];
var p: ptr<i32> = &values[0];
var first: i32 = deref p;
```

## 控制台输入/输出

```wave
print("value = ");
println("{}", value);
input("{}", value);
```

第一个参数是字符串文字。每个精确的 `{}` 占位符后面都需要一个表达式，并且 `input` 目标必须是可分配的。

## import及公共物品

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

本地路径以`./`开头。别名 import 指定模块名称，选择 import 将所需的公共条目导入到该文件的命名空间中。 `pub import` 重新导出所选项目。

## FFI

```wave
extern(c) fun native_call(value: i32) -> i32;

export(c) fun wave_call(value: i32) -> i32 {
    return value + 1;
}
```

## 目标条件项

```wave
#[target(os="linux", arch="riscv64")]
extern(c) fun platform_call(value: i32) -> i32;
```

支持的条件键为`arch`、`os`、`env`、`abi`。属性控制下一个顶级项目。

## 内联装配

```wave
var result: i64 = 0;
asm {
    "mv a0, a1"
    in("a1") 7
    out("a0") result
    clobber("memory")
}
```

指令文本和寄存器名称取决于目标。声明该块所需的所有输入、输出和隐藏clobber。

## 源头检验

```shell
wavec build main.wave --emit=check
wavec print supported-targets
wavec print supported-emit-kinds
```

## 学习和示例范围

单独显示在函数外部的局部变量和语句的示例是插入到函数体中的代码片段。完整的运行示例和练习如下[Wave 学习过程](/docs/zh/getting-started/overview)。内存和外部函数的详细规则请查看[标准库](/docs/zh/stdlib)。
