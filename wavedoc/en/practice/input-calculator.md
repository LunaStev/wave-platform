---
translation_set_id: practice-input-calculator
path: practice/input-calculator
locale: en
group: practice
group_order: 4
order: 1
title: Project: A calculator that validates input
summary: Associate inputs, bounds checks, functions, and exit codes.
---

## Goals and Action

Enter the quantity and unit price and calculate the total. Enter two integers separated by a space or line break. This example only accepts quantities 1 to 1000 and unit prices 0 to 100000, so it is calculated within the multiplication range i32.

Save it to `main.wave`, run `wavec run main.wave`, then enter `3 1200`. The results output by the program are as follows. The visibility of characters entered in the terminal is separate from the program output.

<!-- wave-example: input-calculator -->
```wave
fun total(quantity: i32, price: i32) -> i32 {
    return quantity * price;
}

fun main() -> i32 {
    var quantity: i32 = 0;
    var price: i32 = 0;
    input("{} {}", quantity, price);
    if (quantity < 1 || quantity > 1000 || price < 0 || price > 100000) {
        println("out of range");
        return 1;
    }

    println("total={}", total(quantity, price));
    return 0;
}
```

Execution result:

```text
total=3600
```

## Also check for failures

If you put `0 1200`, it expects `out of range` and exit code 1. Numeric parsing failure in `input` and checking the program's scope of work are two different steps. Non-numeric token/type exceeds range/before required input EOF is an input failure. This built-in input is not an interface that returns an error and requires re-input. If you need a recoverable parser, configure the verification process yourself by reading the bytes with [io](/docs/en/stdlib/files-io).

## Extended exercises and commentary

Take the discount rate as the third input and check whether it is 0 to 100. To avoid large intermediate multiplications, you need to widen the calculation range with i64 and check the range when narrowing down the result. If you simply widen the last result type, intermediate calculations may have already been performed on narrow types.

[Console I/O](/docs/en/language/console-io-and-formatting) · [Next: Reading files](/docs/en/practice/file-reader)

## Why decide the calculation scope first?

The largest inputs are quantity 1000 and unit price 100000. Multiplying the two values ​​is 100000000, so it falls within the range i32. This range check is required so that the multiplication result of the total function can be used as is.

main, which receives input, is responsible for input and error messages, and total is only responsible for calculations. Even if you later change to reading orders from a file, you can still use the calculation function.

## Complete discount calculation

You can multiply the total by 100 by adding the discount percentage. Convert to perform intermediate calculations as i64 and then apply the discount rate. Since the fractional part of integer division is discarded, this example truncates the amount after the discount to an integer.

<!-- wave-example: book-calculator-discount -->
```wave
fun discounted_total(quantity: i32, price: i32, discount: i32) -> i64 {
    var subtotal: i64 = (quantity as i64) * (price as i64);
    var remaining: i64 = 100 - (discount as i64);
    return subtotal * remaining / 100;
}

fun main() -> i32 {
    var quantity: i32 = 0;
    var price: i32 = 0;
    var discount: i32 = 0;

    input("{} {} {}", quantity, price, discount);

    if (quantity < 1 || quantity > 1000) {
        println("invalid quantity");
        return 1;
    }

    if (price < 0 || price > 100000) {
        println("invalid price");
        return 2;
    }

    if (discount < 0 || discount > 100) {
        println("invalid discount");
        return 3;
    }

    println("total={}", discounted_total(quantity, price, discount));
    return 0;
}
```

Execution result:

```text
total=3240
```

The above result is when `3 1200 10` is entered. The output is 3240, which is 10% subtracted from the original total of 3600.

|input|expected result|path to check|
| --- | --- | --- |
| `3 1200 0` | `total=3600` |no discount|
| `3 1200 100` | `total=0` |Full discount|
| `3 1200 101` | `invalid discount` |Discount rate exceeds range|
| `1000 100000 0` | `total=100000000` |maximum input|
| `0 1200 10` | `invalid quantity` |Quantity below range|

## next exercise

Try changing the function to round decimal amounts. In this program, which accepts only positive amounts, adding 50 before dividing by 100 will round to whole numbers. When you enter 99 for 1 and a discount rate of 50, you can compare the cutting result of 49 and the rounding result of 50.
