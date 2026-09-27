---
translation_set_id: data-types
path: language/structures-enums-and-aliases
locale: zh
group: language
group_order: 2
order: 8
title: 8. 结构体、枚举和变体
summary: 了解字段、结构体初始化以及enum和variant的角色。
---

## 将数据之间的关系表达为类型

如果将产品的价格和数量分别作为变量传递，仅通过查看代码很难判断这两个值是否属于同一个产品。结构对相关字段进行分组。 enum 表示命名状态，variant 将每种情况的不同数据存储在一起。

这三个函数不是相互替代的语法。根据你想表达的内容来选择。

|要表达的东西|选择|是的|
| --- | --- | --- |
|多个字段同时存在| struct |产品单价及数量|
|命名整数状态| enum |等待/继续/完成|
|每种情况都有不同的数据| variant |成功值或错误|

## 结构声明和价值创造

<!-- wave-example: book-struct-first -->
```wave
struct Product {
    price: i32;
    quantity: i32;
}

fun main() {
    var item: Product = Product {
        price: 1500,
        quantity: 2
    };

    println("{} {}", item.price, item.quantity);
}
```

执行结果：

```text
1500 2
```

声明中的字段以分号结尾，创建值时，字段和值用冒号连接并用逗号分隔。定义类型和创建实际值是两个不同的步骤。声明类型 Product 不会自动为一个产品创建存储空间。

该字段的访问方式为 `item.price`。如果创建多个相同类型的item，则可以在每个中存储不同的值。

## 将结构体传递给函数

<!-- wave-example: book-struct-function -->
```wave
struct Product {
    price: i32;
    quantity: i32;
}

fun subtotal(item: Product) -> i32 {
    return item.price * item.quantity;
}

fun main() {
    var item: Product = Product {
        price: 1500,
        quantity: 2
    };

    println("{}", subtotal(item));

    item.quantity = 3;
    println("{}", subtotal(item));
}
```

执行结果：

```text
3000
4500
```

该函数接收属于同一产品的单价和数量的关系作为类型。这是一个读取作为值传递的整数字段结构的函数。任何想要修改调用者存储的函数都可以设计为接受指针。

如果结构体具有指针字段，则复制值也会复制地址。这不是深度复制到单独分配的函数。包含文件句柄或Buffer等资源的类型必须一起指定复制和释放规则。

## 方法及proto

相关函数可以以方法形式分组。 proto是将结构体的方法写成单独的块的方法。

<!-- wave-example: book-struct-method -->
```wave
struct Counter {
    value: i32;
}

proto Counter {
    fun current(self: Counter) -> i32 {
        return self.value;
    }
}

fun main() {
    var counter: Counter = Counter {
        value: 3
    };

    println("{}", counter.current());
}
```

执行结果：

```text
3
```

`self: Counter` 是接收值的参数。使用方法调用表示法不会自动使其成为修改原始方法的方法。请一起阅读文中self的类型及其作用。

您不能仅仅因为向某个字段附加了一个方法就期望该字段仅具有有效状态。如果用户可以使用公共字段创建无效组合，您的函数应该检查它们或提供创建规则。

## 将状态命名为 enum

<!-- wave-example: book-enum-state -->
```wave
enum State -> i32 {
    Ready = 0,
    Running,
    Finished,
}

fun main() {
    var state: State = State::Ready;

    if (state == State::Ready) {
        println("ready");
    }

    state = State::Running;

    if (state == State::Running) {
        println("running");
    }
}
```

执行结果：

```text
ready
running
```

`-> i32` 是表达式中使用的整数类型。第一个值设置为0，后续省略的值比前一个值大1。使用 State::Ready 和 State::Running 可以揭示其含义，而不是简单地比较代码中的 0 和 1。

enum 拥有名称并不会自动限制状态转换。是否可以从Finished返回到Running等规则必须以函数的形式实现。

## 使用 variant 连接案例和数据

如果只有成功时才有值，失败时需要错误信息，可以表示为variant。

<!-- wave-example: book-variant-result -->
```wave
variant Result {
    Value(i32),
    Error(i32)
}

fun divide(left: i32, right: i32) -> Result {
    if (left < 0 || right <= 0) {
        return Result::Error(1);
    }

    return Result::Value(left / right);
}

fun show(result: Result) {
    match (result) {
        Result::Value(value) => {
            println("value={}", value);
        }
        Result::Error(code) => {
            println("error={}", code);
        }
    }
}

fun main() {
    show(divide(12, 3));
    show(divide(12, 0));
}
```

执行结果：

```text
value=4
error=1
```

Result::Value 和 Result::Error 各自包含 payload。即使它们包含相同的整数类型，也会区分某些情况。调用者使用 match 检查大小写，并在 arm 中使用 payload。

上面的divide是一个不支持负操作数的小例子。由于指定了输入范围，因此不应将其与覆盖常规 signed 除法的所有边界的函数混淆。

## 如果 payload 不存在

并非所有情况都需要数据。无值状态可以表示为单独的情况。

<!-- wave-example: book-variant-empty -->
```wave
variant Lookup {
    Found(i32),
    Missing
}

fun main() {
    var result: Lookup = Lookup::Missing;

    match (result) {
        Lookup::Found(value) => {
            println("found={}", value);
        }
        Lookup::Missing => {
            println("missing");
        }
    }
}
```

执行结果：

```text
missing
```

我们没有将一个任意整数保留为“none”，而是使用了案例Missing。无论成功值是多少，其意义并不重叠。

match至`_`处理其余情况。如果您希望在添加新案例时再次审查每个调用者，最好明确分隔所有案例。无论您选择哪种方法，请确保没有未处理的输入。

## 结构与variant一起使用

选择不同的数据可以表示为variant，将属于一个案例的多个字段分组可以表示为一个结构体。例如，如果订单处理结果是成功，则可以设计为包含收据结构，如果为失败，则可以设计为包含错误号。

即使该值包含其他值，内存生命周期规则也不会消失。如果variant中存储了一个指针，则该指针是否有效以及由谁来释放是单独确定的。即使以外部文件格式保存，您也必须为每个字段指定编码，而不是按原样转储结构内存。

## 练习与完整解答

创建仅接受 0 到 100 之间分数的判定。有效分数表示为 Grade(score)，其余分数表示为 Invalid。

<!-- wave-example: book-model-solution -->
```wave
variant CheckedScore {
    Grade(i32),
    Invalid
}

fun check_score(score: i32) -> CheckedScore {
    if (score < 0 || score > 100) {
        return CheckedScore::Invalid;
    }

    return CheckedScore::Grade(score);
}

fun main() {
    var result: CheckedScore = check_score(87);

    match (result) {
        CheckedScore::Grade(score) => {
            println("accepted={}", score);
        }
        CheckedScore::Invalid => {
            println("invalid score");
        }
    }
}
```

执行结果：

```text
accepted=87
```

将输入更改为 -1、0、100、101 以检查边界。由于成功和失败不共享相同的整数空间，因此调用者降低了意外将错误值添加到平均计算中的风险。


## 我应该选择哪种表达方式？

|数据的形状|合适的表达|是的|
| --- | --- | --- |
|同一类型的多个值|数组|10分|
|多个字段相互关联|结构体|姓名和分数|
|命名状态值| enum | Ready, Running, Stopped |
|因州而异的附加数据| variant | Value(i32), Error(str) |
|现有类型的上下文名称|类型别名| UserId = u64 |

选择数据结构时，不仅要考虑要存储的值，还要考虑可以表示哪些错误状态。同时具有成功/失败和值/错误字段的结构可能会创建不正确的组合，但对于每种情况，variant 可以区分为 payload。

## 内存布局和外部数据

该结构的内存可能包含空白空间以确保字段之间的对齐。简单地添加字段大小并不总是等于结构的总大小。当您需要知道大小和对齐方式时，请使用[mem 布局功能](/docs/zh/reference/memory-and-buffer)。

对于文件或网络消息，可以通过[bytes](/docs/zh/stdlib/bytes)函数按顺序编码字段来明确确定字节顺序和长度。将结构传递给其他语言时，请将[FFI](/docs/zh/language/modules-imports-and-ffi)的外部声明与目标ABI对齐。

[结构学习与实践](/docs/zh/language/structures-enums-and-aliases) · [variant](/docs/zh/language/variants)
