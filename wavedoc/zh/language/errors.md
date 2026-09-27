---
translation_set_id: learn-errors
path: language/errors
locale: zh
group: language
group_order: 2
order: 12
title: 12. 错误的表示和恢复
summary: 将错误与正常值分开，并清除故障路径中的资源。
---

## 失效度函数结果

丢失文件或越界输入在程序中很自然地发生。处理错误不仅仅是打印一条消息。它是区分失败、检查已完成工作的状态、组织已获得的资源，然后选择是继续还是终止的过程。

在本章中，我们从一个小函数的失败表示开始，然后进展到结果结构、variant、提前返回和资源清理。

## 失败标记不得与成功值重叠

在数组搜索中使用 -1 作为未找到的原因是因为有效索引为 0 或更大。另一方面，在任何整数都可以是正常结果的计算中，如果将-1指定为错误，则无法将其与正常值-1区分开。

0 也是一个经常被误解的值。空串长度0、首位置0、传输字节数0对于不同的函数有不同的含义。不要仅仅因为返回值非零就判断成功。

## 共同回报成功与价值

<!-- wave-example: book-error-result -->
```wave
struct Division {
    ok: bool;
    value: i32;
}

fun divide_nonnegative(left: i32, right: i32) -> Division {
    if (left < 0 || right <= 0) {
        return Division {
            ok: false,
            value: 0
        };
    }

    return Division {
        ok: true,
        value: left / right
    };
}

fun main() {
    var result: Division = divide_nonnegative(0, 3);

    if (!result.ok) {
        println("invalid input");
        return;
    }

    println("value={}", result.value);
}
```

执行结果：

```text
value=0
```

正常结果也可能是 0。不要查看value并猜测是否成功，而是先检查ok。即使失败结果中存在字段value，也不意味着它就是要使用的值。

本例的输入规则是左侧为0或更大，右侧为正。作用域显示在函数名称和说明中，以区别于典型的 signed 整数除法。

## 区分错误原因

根据失败原因添加错误信息以提供额外指导或恢复。还有一种方法可以将函数分为检查输入范围的函数和计算输入范围的函数。

<!-- wave-example: book-error-codes -->
```wave
struct Check {
    ok: bool;
    error: i32;
}

fun check_quantity(value: i32) -> Check {
    if (value < 1) {
        return Check {
            ok: false,
            error: 1
        };
    }

    if (value > 100) {
        return Check {
            ok: false,
            error: 2
        };
    }

    return Check {
        ok: true,
        error: 0
    };
}

fun main() {
    var result: Check = check_quantity(101);

    if (!result.ok) {
        match (result.error) {
            1 => {
                println("quantity is too small");
            }
            2 => {
                println("quantity is too large");
            }
            _ => {
                println("invalid quantity");
            }
        }
    }
}
```

执行结果：

```text
quantity is too large
```

数字的含义由该函数定义。不能将其视为与其他库中的错误号 1 或 2 相同。在 public API 中，命名错误常量或类型可以使调用者不必记住随机数。

## 用 variant 分隔结果

成功值和错误不能同时存在的关系可以表示为variant。

<!-- wave-example: book-error-variant -->
```wave
variant ByteResult {
    Value(u8),
    OutOfRange(i32)
}

fun to_byte(value: i32) -> ByteResult {
    if (value < 0 || value > 255) {
        return ByteResult::OutOfRange(value);
    }

    return ByteResult::Value(value as u8);
}

fun main() {
    var result: ByteResult = to_byte(300);

    match (result) {
        ByteResult::Value(value) => {
            println("byte={}", value);
        }
        ByteResult::OutOfRange(original) => {
            println("out of range: {}", original);
        }
    }
}
```

执行结果：

```text
out of range: 300
```

错误包含原始输入值。对于调用者来说，解释问题比简单地返回false更容易。通过不在日志中留下密码或令牌等敏感输入，还考虑了数据的性质。

## 提前返程，让正常路线更容易阅读

执行多步检查时，无需将所有好的代码都放入if。如果失败，您可以立即返回并继续下面的成功路径。

<!-- wave-example: book-error-early-return -->
```wave
fun show_total(quantity: i32, price: i32) {
    if (quantity < 1 || quantity > 100) {
        println("invalid quantity");
        return;
    }

    if (price < 0 || price > 100000) {
        println("invalid price");
        return;
    }

    var total: i32 = quantity * price;
    println("total={}", total);
}

fun main() {
    show_total(3, 1200);
    show_total(0, 1200);
    show_total(3, -1);
}
```

执行结果：

```text
total=3600
invalid quantity
invalid price
```

通过测试后，您可以利用quantity和price在指定范围内的事实。设置范围以使中间乘法也落入i32之内。添加提前回报时，您还应该检查此时您是否已经拥有任何资源。

## 故障路径的内存清理

<!-- wave-example: book-error-cleanup -->
```wave
import("std::buffer::alloc")::{
    Buffer, buffer_init, buffer_free
};
import("std::buffer::write")::{
    buffer_append_str
};

fun make_message(allowed: bool) -> i32 {
    var data: Buffer;

    if (buffer_init(&data, 16) < 0) {
        return 1;
    }

    if (!allowed) {
        buffer_free(&data);
        return 2;
    }

    if (buffer_append_str(&data, "ready") < 0) {
        buffer_free(&data);
        return 3;
    }

    println("bytes={}", data.len);

    if (buffer_free(&data) < 0) {
        return 4;
    }

    return 0;
}

fun main() {
    println("status={}", make_message(false));
}
```

执行结果：

```text
status=2
```

即使在不输出数据的故障路径上也可以释放Buffer。重要的是要清楚地管理每个分支机构所拥有的内容，而不是将每个故障都变成单个return。

还有 API 清理本身失败。决定如何保留原始作品中的错误并清除错误。这个小例子首先返回任务错误号。较大的程序可能会分别记录两者。

## 部分成功不会自动取消

如果将一些字节写入文件然后写入失败，则已写入的字节不会丢失。网络的另一端可能也收到了一些数据。从头开始重复相同的任务可能会导致重复记录。

相反，读取 checked byte cursor 会在失败时保留位置和输出值。这些函数可以接收更多输入并在同一位置重试。不要将“如果失败，则不会发生任何变化”规则应用于每个API，而是检查该函数的文档。

## 可恢复的错误和陷阱

无效的文件路径或无效的用户输入可以报告为错误值，以便调用者可以恢复。由无效的运行时移位计数或浮点到整数转换引起的陷阱是一种不同的机制。

必须继续执行的程序必须在危险操作之前检查其输入。 assert 也无意被滥用作为处理用户输入的正常失败的方法。如果您需要让用户有机会重新输入其输入，请返回结果以继续控制流程。

## 练习与完整解答

编写一个函数，按索引读取数组元素，当索引为负数或大于或等于长度时失败。使用结果结构，因为零可以是有效的元素值。

<!-- wave-example: book-error-solution -->
```wave
struct Lookup {
    ok: bool;
    value: i32;
}

fun at(values: ptr<i32>, length: i32, index: i32) -> Lookup {
    if (index < 0 || index >= length) {
        return Lookup {
            ok: false,
            value: 0
        };
    }

    return Lookup {
        ok: true,
        value: deref values[index]
    };
}

fun main() {
    var values: array<i32, 3> = [0, 10, 20];
    var first: Lookup = at(&values[0], 3, 0);
    var outside: Lookup = at(&values[0], 3, 3);

    if (first.ok) {
        println("value={}", first.value);
    }

    if (!outside.ok) {
        println("out of bounds");
    }
}
```

执行结果：

```text
value=0
out of bounds
```

调用者的条件是指针和length代表实际可读的数组。这不是一个仅通过检查索引就可以使随机地址安全的函数。请单独阅读该函数负责哪些检查以及调用者必须保证哪些条件。
