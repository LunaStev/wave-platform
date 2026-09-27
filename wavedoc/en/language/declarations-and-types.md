---
translation_set_id: types
path: language/declarations-and-types
locale: en
group: language
group_order: 2
order: 2
title: 2. Variables, types, and scope
summary: Learn local variables, integer scope, bool, and scopes.
---

## Treating values by name

If you write your prices directly in multiple places, you'll have to find all of them when you change the price. A variable is a storage space that gives names to values ​​and allows you to read and change them using those names. In this chapter, you will learn about the scope of declarations, assignments, types, and the scope of visible names within blocks.

The programs below are each part of a separate main.wave. Save one example, run it as `wavec run main.wave`, and replace it with the next example.

## Declaration and initialization

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

Execution result:

```text
price=1200
quantity=3
total=3600
```

Read the declaration in four parts.

|part|this example|role|
| --- | --- | --- |
|declaration keyword| var |Create local variable|
|name| price |Identifier to be used later|
|type| i32 |Type and range of values to store|
|initial value| 1200 |First value to store|

The colon before the type and the equal sign before the initial value have different roles. Make your name meaningful. In this example, price is the unit price and quantity is the quantity. Even if the same i32 is used interchangeably, incorrect calculations can be made without grammatical errors.

## Assignment is not a formula that maintains relationships.

If you save the calculation result in a variable, the value at that point will be entered. It does not remember calculations and automatically re-evaluate them later.

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

Execution result:

```text
before recalculation=3600
after recalculation=6000
```

`quantity = 5` writes a new value to an already existing variable. It must be distinguished from redeclaration like `var quantity`. `total` is also 3600 before substituting again. If a program must maintain relationships between multiple variables, it should be written to perform calculations when the relationships change.

## Calculate next value from previous value

Calculate the right side of the assignment statement first and write the result to the storage space on the left. Unlike equations in mathematics, `count = count + 1` is a valid update.

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

Execution result:

```text
9
```

The progression of values is 0 → 1 → 3 → 9. `+=`, `*=` express calculation and storage together. If you write down the calculation order, you can find out at what stage your thoughts were different when the results are different from what you expected.

## Width and sign of integer types

If it starts with `i`, it becomes signed, and if it starts with `u`, it becomes unsigned. The last number is the number of bits. As the number of bits increases, the range that can be expressed increases and storage space also increases.

|type|minimum value|maximum value|Example use|
| --- | --- | --- | --- |
| i8 | -128 | 127 |small signed value|
| u8 | 0 | 255 |one byte|
| i16 | -32768 | 32767 |small integer data|
| u16 | 0 | 65535 |Port/16-bit field|
| i32 | -2147483648 | 2147483647 |Common small integer calculations|
| u32 | 0 | 4294967295 |32-bit bit field|

Wave also provides signed and unsigned integers of 64, 128, 256, 512, and 1024 bits. A wider type does not make every calculation safe: the result can still exceed the selected range. Choose the range you need first. `isz` and `usz` follow the target address width.

A distinction must be made between storing a large literal in a small type and intentionally discarding bits by doing cast. Converting to a smaller type simply to eliminate errors may change the value itself. Transformations are covered in the next chapter.

## Floating-point types and bool

`f32`·`f64` are floating point numbers. Unlike integers, they can represent decimal parts, but they cannot store all decimal numbers exactly. This is one reason why amounts are managed in small integer units.

bool represents true and false. You can give your condition a name by saving the comparison results like this:

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

Execution result:

```text
purchase allowed
```

`can_buy` does not automatically follow the changes in balance and cost. After the two values ​​are swapped, the comparison is performed again if the current state is needed.

## uninitialized storage space

`var value: i32;` is a form that declares only the storage space. A valid value must be written before reading. Don't assume that 0 is automatically inserted just because you declare it. In an introductory course, it is easier to understand if the value is immediately known and initialized at the same time as declaration.

In library calls that receive values as output arguments, there are cases where a space is declared first and then read when successful. At that time, you should check the success result of the function. Avoid the mistake of reading uninitialized output after a failed call.

## Valid range of blocks and names

A block is a region of code enclosed in curly braces. If you declare a new variable with the same name inside, the new variable will be used inside that block. This is called shadowing.

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

Execution result:

```text
inner=15
inner changed=16
outer=10
```

The initial expression `value + 5` in the inner declaration reads the outer value. After initialization of the new variable is completed, the inside value is 15. Even if you change the inner value to 16, the outer storage space does not change. After the block, you will see value outside again.

Conversely, if you only execute `value += 1` without `var` in the inner block, the existing visible variables will be changed. Look at the keyword to determine whether it is a new declaration or a change to an existing value.

## Local variables and top-level storage

Outside the function, you can use const and static. const represents a constant value and static is the storage space maintained during execution. They cannot be declared everywhere in the same way as local variables.

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

Execution result:

```text
visit=1
visit=2
limit=3
```

Even if visit is called twice, static will not be reset to 0 on each call. On the other hand, if you declare it as `var visits: i32 = 0;` inside a function, it will initialize the local storage every time you call it. Shared mutable state can make behavior difficult to track, so first consider whether it can be resolved with the function's inputs and outputs.

## common mistakes

- When declaration and assignment are confused and the same name is declared again unnecessarily.
- If you think that the variable storing the calculation result automatically follows changes in the input variable.
- If you think that since the type is the same, the units such as number and number of bytes are also the same.
- When reading without initialization or when the output argument is read when the function fails.
- If you think the name of a local variable is visible outside the block.

When looking for a name in error, check the name's declaration position and brace range.

## Exercise: Calculating Inventory Changes

Initial stock is 20 units and sold twice, 3 units each. Print out remaining inventory and total units sold. Whenever inventory is changed, the same variables are updated and sales volume is also accumulated separately.

### Complete solution

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

Execution result:

```text
stock=14 sold=6
```

Inventory and sales volume must change together. Updating either one breaks the relationship between the values. We will group repeated order processing by learning functions and loop statements.


## Integer and floating point types

Integer types are as follows:

- Signed: `i8`, `i16`, `i32`, `i64`, `i128`, `i256`, `i512`, `i1024`
- Unsigned: `u8`, `u16`, `u32`, `u64`, `u128`, `u256`, `u512`, `u1024`
- Address size integer: `isz`, `usz`
- Floating point: `f32`, `f64`

`isz` is a signed integer type that matches the address size, and `usz` is an unsigned integer type that matches the address size.

## Other built-in types

|type|Use|
| --- | --- |
| `bool` |`true` or `false`|
| `char` |An unsigned 8-bit character value. Not arbitrary Unicode code point type|
| `byte` |8-bit byte value|
| `str` |String bytes ending with NUL|
| `ptr<T>` |Pointer targeting `T`|
| `array<T, N>` |Fixed-length array with element type `T` and length `N`|

User-defined structures, enumerations, and type aliases can also be used in type locations.

`var` is the syntax for declaring local variables. A type alias is a grammar that expresses the same type with a name that fits the context of the code.
