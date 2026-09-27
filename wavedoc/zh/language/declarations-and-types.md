---
translation_set_id: types
path: language/declarations-and-types
locale: zh
group: language
group_order: 2
order: 2
title: 2. 变量、类型和范围
summary: 学习局部变量、整数作用域、bool 和作用域。
---

## 按名称处理值

如果您直接在多个位置写下价格，则更改价格时必须找到所有位置。变量是一个存储空间，它为值提供名称，并允许您使用这些名称读取和更改它们。在本章中，您将了解声明、赋值、类型的范围以及块内可见名称的范围。

下面的程序是单独的main.wave的每个部分。保存一个示例，将其运行为`wavec run main.wave`，并将其替换为下一个示例。

## 声明和初始化

<!-- wave-example: book-variable-first -->
```wave
fun main() {
    var price: i32 = 1200;
    var quantity: i32 = 3;
    println("price={}", price);
    println("quantity={}", quantity);
    println("total={}", price * quantity);
}
```

执行结果：

```text
price=1200
quantity=3
total=3600
```

阅读该声明的四个部分。

|部分|这个例子|角色|
| --- | --- | --- |
|声明关键字| var |创建局部变量|
|姓名| price |稍后使用的标识符|
|类型| i32 |要存储的值的类型和范围|
|初始值| 1200 |要存储的第一个值|

类型前的冒号和初始值前的等号具有不同的作用。让你的名字有意义。在此示例中，price 是单价，quantity 是数量。即使相同的i32互换使用，也可能在没有语法错误的情况下做出错误的计算。

## 分配不是维持关系的公式。

如果将计算结果保存在变量中，则将输入该点的值。它不记住计算并在以后自动重新评估它们。

<!-- wave-example: book-assignment-snapshot -->
```wave
fun main() {
    var price: i32 = 1200;
    var quantity: i32 = 3;
    var total: i32 = price * quantity;
    quantity = 5;
    println("before recalculation={}", total);
    total = price * quantity;
    println("after recalculation={}", total);
}
```

执行结果：

```text
before recalculation=3600
after recalculation=6000
```

`quantity = 5` 将新值写入已存在的变量。它必须与像`var quantity`这样的重新声明区分开来。再次替换之前`total`也是3600。如果程序必须维护多个变量之间的关系，则应编写该程序以在关系发生变化时执行计算。

## 根据前一个值计算下一个值

先计算赋值语句右侧，并将结果写入左侧存储空间。与数学方程不同，`count = count + 1` 是有效的更新。

<!-- wave-example: book-update-variable -->
```wave
fun main() {
    var count: i32 = 0;
    count = count + 1;
    count += 2;
    count *= 3;
    println("{}", count);
}
```

执行结果：

```text
9
```

数值的级数为0 → 1 → 3 → 9。 `+=`、`*=` 一起表示计算和存储。如果你写下计算顺序，你就可以发现在哪个阶段，当结果与你的预期不同时，你的想法有所不同。

## 整数类型的宽度和符号

如果以`i`开头，则变为signed，如果以`u`开头，则变为unsigned。最后一个数字是位数。随着位数的增加，可表达的范围增大，存储空间也增大。

|类型|最小值|最大值|使用示例|
| --- | --- | --- | --- |
| i8 | -128 | 127 |小符号值|
| u8 | 0 | 255 |一个字节|
| i16 | -32768 | 32767 |小整数数据|
| u16 | 0 | 65535 |端口/16 位字段|
| i32 | -2147483648 | 2147483647 |常见小整数计算|
| u32 | 0 | 4294967295 |32位位域|

Wave 还提供 64、128、256、512 和 1024 位的有符号和无符号整数。更宽的类型并不能使每次计算都安全：结果仍然可能超出所选范围。首先选择您需要的范围。 `isz` 和`usz` 遵循目标地址宽度。

必须区分在小类型中存储大文字和通过执行cast有意丢弃位。仅仅为了消除错误而转换为较小的类型可能会改变值本身。下一章将介绍转换。

## 浮点类型和布尔值

`f32`·`f64` 是浮点数。与整数不同，它们可以表示小数部分，但不能准确存储所有小数。这就是为什么金额以小整数单位进行管理的原因之一。

bool代表真与假。您可以通过保存比较结果来为您的条件命名，如下所示：

<!-- wave-example: book-named-condition -->
```wave
fun main() {
    var balance: i32 = 5000;
    var cost: i32 = 3600;
    var can_buy: bool = balance >= cost;
    if (can_buy) {
        println("purchase allowed");
    } else {
        println("not enough money");
    }
}
```

执行结果：

```text
purchase allowed
```

`can_buy` 不会自动跟随 balance 和 cost 中的更改。两个值交换后，如果需要当前状态，则再次进行比较。

## 未初始化的存储空间

`var value: i32;`是只声明存储空间的形式。读取之前必须写入有效值。不要仅仅因为声明了 0 就认为它会自动插入。在入门课程中，如果在声明的同时立即知道并初始化该值，则更容易理解。

在接收值作为输出参数的库调用中，有时会先声明一个空格，然后在成功时读取该空格。此时，您应该检查该函数的成功结果。避免调用失败后读取未初始化输出的错误。

## 块和名称的有效范围

块是用大括号括起来的代码区域。如果您在内部声明一个同名的新变量，则该新变量将在该块内使用。这称为shadowing。

<!-- wave-example: book-shadowing -->
```wave
fun main() {
    var value: i32 = 10;

    if (true) {
        var value: i32 = value + 5;
        println("inner={}", value);
        value += 1;
        println("inner changed={}", value);
    }

    println("outer={}", value);
}
```

执行结果：

```text
inner=15
inner changed=16
outer=10
```

内部声明中的初始表达式`value + 5`读取外部value。新变量初始化完成后，内部value为15。即使将内部值改为16，外部存储空间也不会改变。过了街区，你会再次在外面看到value。

相反，如果在内部块中只执行`value += 1`而不执行`var`，则现有的可见变量将被更改。查看关键字以确定它是新声明还是对现有值的更改。

## 局部变量和顶级存储

在函数之外，您可以使用const和static。 const代表常量值，static是执行期间维护的存储空间。它们不能像局部变量一样在任何地方声明。

<!-- wave-example: book-storage -->
```wave
const LIMIT: i32 = 3;
static visits: i32 = 0;

fun visit() {
    visits += 1;
    println("visit={}", visits);
}

fun main() {
    visit();
    visit();
    println("limit={}", LIMIT);
}
```

执行结果：

```text
visit=1
visit=2
limit=3
```

即使visit被调用两次，static也不会在每次调用时重置为0。另一方面，如果您在函数内将其声明为 `var visits: i32 = 0;`，则每次调用它时都会初始化本地存储。共享可变状态可能会使行为难以跟踪，因此首先考虑是否可以通过函数的输入和输出来解决它。

## 常见错误

- 当声明和赋值混淆并且不必要地再次声明相同的名称时。
- 如果您认为存储计算结果的变量会自动跟随输入变量的变化。
- 如果你认为既然类型相同，那么个数、字节数等单位也相同。
- 在未初始化的情况下读取或在函数失败时读取输出参数时。
- 如果您认为局部变量的名称在块外可见。

当查找错误的名称时，请检查名称的声明位置和大括号范围。

## 练习：计算库存变化

初始库存为 20 件，分两次出售，每次 3 件。打印剩余库存和销售总量。每当库存发生变化时，相同的变量都会更新，并且销量也会单独累积。

### 完整解答

<!-- wave-example: book-stock-solution -->
```wave
fun main() {
    var stock: i32 = 20;
    var sold: i32 = 0;
    var order: i32 = 3;
    stock -= order;
    sold += order;
    stock -= order;
    sold += order;
    println("stock={} sold={}", stock, sold);
}
```

执行结果：

```text
stock=14 sold=6
```

库存和销量必须一起变化。更新其中任何一个都会破坏值之间的关系。我们将通过学习函数和循环语句来对重复的订单处理进行分组。


## 整数和浮点类型

整数类型如下：

- 签名：`i8`、`i16`、`i32`、`i64`、`i128`、`i256`、`i512`、`i1024`
- 无符号：`u8`、`u16`、`u32`、`u64`、`u128`、`u256`、`u512`、`u1024`
- 地址大小整数：`isz`、`usz`
- 浮点：`f32`、`f64`

`isz`是与地址大小匹配的有符号整数类型，`usz`是与地址大小匹配的无符号整数类型。

## 其他内置类型

|类型|使用|
| --- | --- |
| `bool` |`true` 或 `false`|
| `char` |无符号 8 位字符值。不是任意 Unicode 代码点类型|
| `byte` |8 位字节值|
| `str` |以 NUL 结尾的字符串字节|
| `ptr<T>` |指针定位 `T`|
| `array<T, N>` |元素类型为 `T` 且长度为 `N` 的固定长度数组|

用户定义的结构、枚举和类型别名也可以在类型位置中使用。

`var`是声明局部变量的语法。类型别名是一种语法，它用适合代码上下文的名称来表达相同的类型。
