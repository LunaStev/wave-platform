---
translation_set_id: modules-ffi
path: language/modules-imports-and-ffi
locale: zh
group: language
group_order: 2
order: 11
title: 11. 模块和通用代码
summary: 了解如何获取公共名称和显式类型参数。
---

## 分割文件的原因

随着程序的增长，通过将相关函数分组在一起而不是将所有函数放在main.wave中更容易找到它。模块边界决定哪些名称暴露给其他代码。泛型是一种以不同类型重用同一任务的工具，与文件分离无关。

在本章中，我们创建一个两个文件的程序并学习模块别名、选择import以及通用函数和结构。

## 两个文件程序

在同一目录下创建helpers.wave和main.wave。

helpers.wave 全部：

```wave
pub fun double(value: i32) -> i32 {
    return value * 2;
}
```

main.wave 全部：

<!-- wave-example: book-module-two-files -->
```wave
import("./helpers")::{
    double
};

fun main() {
    var result: i32 = double(21);

    println("{}", result);
}
```

执行结果：

```text
42
```

在终端中运行`wavec run main.wave`。 helpers.wave 也不单独运行。所需源通过import 连接。

helpers中函数前面的pub表示可以被其他模块导入。不需要对外暴露的辅助功能不需要公开。即使您更改了模块的内部实现，也可以通过保留公共函数的契约来减少对所使用代码的更改。

## 相对路径的基线

`./helpers` 相对于创建import 句子的源文件的目录。运行程序时，将其与写入文件I/O的工作目录分开。运行时查找文件import的步骤和查找input.txt的步骤是不同的。

对于本地import，可以省略`.wave`扩展名。如果划分了目录，则根据位置写出路径`./module`。不要将本地相对路径与获取包依赖项名称的路径混淆。

## 选择import和别名

选择 import 只会导致在当前文件中直接使用所需的公共名称。如果存在名称冲突或者您想显示该函数属于哪个模块，请使用别名。

<!-- wave-example: book-module-alias -->
```wave
import("std::string::len" as strings);

fun main() {
    var length: i32 = strings::len("Wave");

    println("{}", length);
}
```

执行结果：

```text
4
```

strings 是该文件中定义的模块别名。 `strings::len` 使用模块的名称。用于字段访问的点和用于模块分隔的`::`是不同的符号。

请勿在一句话中同时使用选项 import 和别名 import。无论您采用哪种风格，请在整个文件中一致地使用它，以便轻松读取名称的来源。

## 标准库和包

路径`std::`指向标准库。用户发布API和import所需的模块。并非所有标准库函数都会自动放入当前命名空间中。

外部包路径从包名称开始。包的位置由编译器选项或包管理器提供。首先了解本地模块的边界，然后在[Vex 使用方法](/docs/zh/whale/vex-package-manager)中学习如何管理依赖关系。

## 每种类型使用相同的函数

以下函数逐字返回其输入： 对于 i32 和 str，请使用类型参数 T 以避免编写相同的代码两次。

<!-- wave-example: book-generic-identity -->
```wave
fun identity<T>(value: T) -> T {
    return value;
}

fun main() {
    var number: i32 = identity<i32>(42);
    var text: str = identity<str>("Wave");

    println("{} {}", number, text);
}
```

执行结果：

```text
42 Wave
```

T是输入类型的地方，不是执行过程中传递的整数值。通过指定类型参数（如 `<i32>`）来调用它。普通用户泛型函数不会省略类型参数。

identity<str> 不会通过再次分配字符串字节来重复它们。按原样返回值。泛型的语法不会改变数据的复制和所有权规则。

## 通用机构要求的操作

仅仅因为有类型参数并不意味着所有操作都可以用于所有类型。下面的minimum应该用作可以比较的实际类型。

<!-- wave-example: book-generic-minimum -->
```wave
fun minimum<T>(left: T, right: T) -> T {
    if (left < right) {
        return left;
    }

    return right;
}

fun main() {
    var small: i32 = minimum<i32>(7, 4);
    var wide: i64 = minimum<i64>(100, 20);

    println("{} {}", small, wide);
}
```

执行结果：

```text
4 20
```

如果更改类型参数，则文本中使用的`<`和返回必须是相应的类型。读取通用错误时，请检查调用的类型组合和函数体所需的操作。

## 通用结构

让我们创建 Pair，它绑定两个不同的值。

<!-- wave-example: book-generic-pair -->
```wave
struct Pair<A, B> {
    first: A;
    second: B;
}

fun main() {
    var item: Pair<i32, str> = Pair<i32, str> {
        first: 7,
        second: "seven"
    };

    println("{} {}", item.first, item.second);
}
```

执行结果：

```text
7 seven
```

Pair<i32, str> 和 Pair<i64, str> 是不同的具体类型。类型参数的顺序也有意义。阅读声明和生成代码，看看first和second的类型是在哪里确定的。

## API姓名及合约公布

当您发布函数时，您不仅指定名称，还指定输入单位、返回值、失败和所有权。例如，调用者将编写的循环取决于read是否读取最大长度或确切长度。

pub是模块之间的公共作用域Wave。这与export（c）不同，它导出外部符号供其他语言调用。您可以在[参见FFI](/docs/zh/language/modules-imports-and-ffi)查看链接两种语言的完整示例。

## 练习与完整解答

在 math.wave 上创建一个公共函数 square 并在 main.wave 上将其作为别名调用，以打印 3 和 5 的幂。

math.wave:

```wave
pub fun square(value: i32) -> i32 {
    return value * value;
}
```

main.wave:

<!-- wave-example: book-module-solution -->
```wave
import("./math" as math);

fun main() {
    println("{} {}", math::square(3), math::square(5));
}
```

执行结果：

```text
9 25
```

如果错误是功能消失了，请先检查import路径和pub。如果存在名称冲突，请检查调用是否有别名。如果是类型错误，请检查函数的输入和传递的参数的类型。不要试图通过一种路径修正来解决不同的问题。


## C 导入功能

```wave
extern(c) fun puts(text: ptr<i8>) -> i32;
```

名称 ABI 后面可以跟有字符串形式的实际交易品种名称。

```wave
extern(c, "native_symbol") fun local_name(value: i32) -> i32;
```

## Wave 导出功能

```wave
export(c) fun wave_add(left: i32, right: i32) -> i32 {
    return left + right;
}
```

`extern` 和 `export` 可用作单个函数和块。导出的函数必须具有特定的签名ABI，因此不能是通用的。

## 目标条件属性

目标条件属性可以附加到顶级项目。

```wave
#[target(os="linux", arch="x86_64")]
extern(c) fun platform_call(value: i32) -> i32;
```

条件键为 `arch`、`os`、`env`、`abi`，并且属性适用于下一个顶级项目。

## 连接你自己写的C函数

本实验适用于使用 C 编译器的本机环境。连接整数函数，无需库分配或字符串处理。

`native.c`:

```c
#include <stdint.h>
int32_t native_double(int32_t value) { return value * 2; }
```

`main.wave`:

```wave
extern(c) fun native_double(value: i32) -> i32;

fun main() {
    println("{}", native_double(21));
}
```

从同一工作目录中的终端运行 Linux/macOS。

```shell
cc -c native.c -o native.o
wavec build main.wave native.o -o ffi-example
./ffi-example
```

预期输出为`42`。在Windows中的MSVC的开发者shell中，使用`cl /c native.c /Fonative.obj`创建object并连接到`wavec build main.wave native.obj -o ffi-example.exe`。 object的源架构和目标架构必须相同。这些示例仅使用较小的值。要将大值传递给C函数，还必须单独保证C页的乘法范围。

本地文件路径以`./`开头。
