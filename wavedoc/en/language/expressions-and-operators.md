---
translation_set_id: expressions
path: language/expressions-and-operators
locale: en
group: language
group_order: 2
order: 3
title: 3. Arithmetic, comparisons, and conversions
summary: Learn calculation order, integer division, bitwise operations, and cast.
---

## View calculation results and calculation types together

An expression is code that calculates a value. Variable names, literals, function calls, and expressions that connect multiple values ​​with operators are all expressions. In mathematics, even if it looks like the same equation, the result will be different depending on whether it is an integer or a real number, and how many bits it contains.

Starting with simple calculations, this chapter introduces parentheses, division, logical operations, bitwise operations, and casts. Each example is a complete main.wave file that you can run with `wavec run main.wave`.

## Range enclosed in parentheses

<!-- wave-example: book-precedence -->
```wave
fun main() {
    var first: i32 = 2 + 3 * 4;
    var second: i32 = (2 + 3) * 4;

    println("{} {}", first, second);
}
```

Execution result:

```text
14 20
```

Multiplication is evaluated before addition, so the first expression is 2+12. In the second equation, we take the sum in parentheses, which is 5, and then multiply by 4. The goal is not to use fewer parentheses. It is recommended to use it so that the reader can easily understand the scope of the calculation.

Even when using the same operator multiple times, the direction of chaining is important. `20 - 5 - 3` is `(20 - 5) - 3`, which is 12. `20 - (5 - 3)` is 18. The exact full sequence is in [operator reference](/docs/en/language/expressions-and-operators).

## Integer division and remainder

Dividing two integers does not produce a fractional floating-point result. Use division for the quotient and the remainder operator for the remainder.

<!-- wave-example: book-division -->
```wave
fun main() {
    var items: i32 = 17;
    var box_size: i32 = 5;
    var full_boxes: i32 = items / box_size;
    var remaining: i32 = items % box_size;

    println("boxes={} remaining={}", full_boxes, remaining);
    println("negative quotient={}", -17 / 5);
}
```

Execution result:

```text
boxes=3 remaining=2
negative quotient=-3
```

Dividing 17 items into groups of 5 produces 3 complete groups and 2 remaining items. Signed division truncates toward zero, so -17/5 is -3. This differs from rounding down toward negative infinity.

You cannot divide by 0. signed The value obtained by dividing the minimum value by -1 does not fall into the same type. Functions that take these inputs should check before dividing or use checked math API.

## The point of conversion changes the result.

The following two expressions are both stored in the variable f64, but the calculation process is different.

<!-- wave-example: book-division-conversion -->
```wave
fun main() {
    var already_divided: f64 = (7 / 2) as f64;
    var floating_division: f64 = (7 as f64) / (2 as f64);

    if (already_divided == 3.0) {
        println("integer division lost the fraction");
    }

    if (floating_division == 3.5) {
        println("floating division kept the fraction");
    }
}
```

Execution result:

```text
integer division lost the fraction
floating division kept the fraction
```

The first expression performs integer division to obtain 3, then converts it to f64. The second converts the operands to f64 before performing floating-point division. Choosing a wider type for the final variable cannot recover information lost earlier.

Floating-point values are approximations. Two results that look like the same decimal value may not be suitable for an exact equality comparison. Choose a tolerance that fits the units and scale of the problem; one fixed epsilon is not suitable for every calculation.

## Comparison made bool

`<`, `<=`, `>`, `>=`, `==`, `!=` check the relationship. One equal sign `=` is an assignment, and two equal signs `==` are equal comparisons.

<!-- wave-example: book-comparisons -->
```wave
fun main() {
    var age: i32 = 20;
    var minimum: i32 = 18;
    var eligible: bool = age >= minimum;

    if (eligible) {
        println("eligible");
    }

    if (age != minimum) {
        println("not exactly the boundary");
    }
}
```

Execution result:

```text
eligible
not exactly the boundary
```

“At least 18” includes 18; “greater than 18” excludes it. Test 17, 18, and 19 to check this boundary. Mixed-type comparisons depend on signedness and width, so converting both operands to the intended type can make the comparison clearer.

## Logical operations and short-circuit evaluation

`&&` checks if both are true, `||` checks if more than one is true, and `!` flips true/false. This operation does not always execute the expression on the right-hand side.

<!-- wave-example: book-short-circuit -->
```wave
fun report() -> bool {
    println("right side evaluated");
    return true;
}

fun main() {
    var enabled: bool = false;

    if (enabled && report()) {
        println("both true");
    }

    if (!enabled || report()) {
        println("at least one true");
    }
}
```

Execution result:

```text
at least one true
```

In the first condition, enabled is false, so there is no need to look to the right. In the second, !enabled is true, so the right-hand side is also not needed. Therefore, the output of report never appears.

You can use this to check the denominator before division. The function body fragment `if (divisor != 0 && value / divisor > 2) { ... }` does not perform division when the denominator is 0. However, this test alone does not resolve other boundaries such as signed minimum value/-1.

`&&` takes precedence over `||`. In complex policies, indicate your intent in parentheses, such as `(member && active) || admin`.

## Narrowing or Widening Integers

`as` is an explicit cast. The high-order bits that are discarded when narrowing an integer cannot be recovered by widening it again.

<!-- wave-example: book-cast-chain -->
```wave
fun main() {
    var original: i32 = 300;
    var small: u8 = original as u8;
    var widened: i32 = small as i32;

    println("{} -> {} -> {}", original, small, widened);

    var negative: i8 = -1;
    var signed_value: i32 = negative as i32;
    var same_bits: u8 = negative as u8;

    println("{} {}", signed_value, same_bits);
}
```

Execution result:

```text
300 -> 44 -> 44
-1 255
```

The lower 8 bits of 300 are 44. If you widen -1 to the signed type, the sign will be expanded to maintain -1. If you interpret the same 8 bits as unsigned, it is 255.

Transformations that seek to preserve in-range values and transformations that attempt to manipulate storage bits have different purposes. If you have user input, first check if it is the destination range and convert it. The presence of cast does not guarantee that the value was in a safe range.

## Replace with bool

Converting an integer to bool yields false for zero and true otherwise. This does not truncate to the lowest bit: 2 also converts to true. For floating-point values, only +0.0 and -0.0 convert to false; all other values, including NaN and infinities, convert to true.

<!-- wave-example: book-bool-conversion -->
```wave
fun main() {
    var zero: i32 = 0;
    var two: i32 = 2;

    if (!(zero as bool)) {
        println("zero is false");
    }

    if (two as bool) {
        println("two is true");
    }
}
```

Execution result:

```text
zero is false
two is true
```

Pointer-to-bool conversions are not supported. Compare the pointer with null explicitly, for example with `pointer != null`. Whether a non-null address is safe to read is a separate question.

## bit operations

`&`, `|`, `^`, `~` cover each bit of the integer. It can be used to express permissions or functions in bits. Below, 1 is read permission and 2 is write permission.

<!-- wave-example: book-bit-flags -->
```wave
const READ: u32 = 1;
const WRITE: u32 = 2;

fun main() {
    var permissions: u32 = READ | WRITE;

    if ((permissions & WRITE) != 0) {
        println("write enabled");
    }

    permissions = permissions & ~WRITE;
    println("remaining={}", permissions);
}
```

Execution result:

```text
write enabled
remaining=1
```

Add a bit with OR and check for the presence of a specific bit with AND. Create a mask with only the relevant bits being 0 with `~WRITE` and erase it. The bitwise operations `&`·`|` are different operators from the short-circuit evaluation `&&`·`||` of bool.

## Width and number of shifts

A left shift moves bits left and discards high bits outside the operand’s width. A right shift uses sign extension for signed values and zero extension for unsigned values. The result always has the left operand’s type.

<!-- wave-example: book-shifts -->
```wave
fun main() {
    var bits: u8 = 129;
    var shifted: u8 = bits << 1;
    var negative: i32 = -8;
    var positive: u32 = 8;

    println("left={}", shifted);
    println("right={} {}", negative >> 1, positive >> 1);
}
```

Execution result:

```text
left=2
right=-4 4
```

Shifting the u8 value 129 left by one discards its highest bit and leaves 2. The original shift-count value must be nonnegative and less than the left operand’s bit width. For u8, valid counts are 0 through 7. An invalid constant count is a compile-time error; an invalid runtime count causes a trap.

## Converting floating-point values to integers

Floating-point-to-integer conversion first truncates toward zero, then checks the destination integer range. NaN, infinity, and out-of-range results are invalid. Invalid constant conversions produce a compile-time error; invalid runtime conversions cause a trap.

A trap does not return an error value from the function. For recoverable conversion failures, design an interface that checks the range before converting. To obtain the stored bits of a floating-point value, use the bit-conversion functions in `std::math::float` instead of a numeric cast.

## Exercise and solution

Write a function that divides 137 won in change into 50 won units and the remainder, and changes only integers in the range 0 to 255 to u8. This example demonstrates the validation step as a separate function that indicates out-of-range as -1.

<!-- wave-example: book-expression-solution -->
```wave
fun checked_byte(value: i32) -> i32 {
    if (value < 0 || value > 255) {
        return -1;
    }

    var small: u8 = value as u8;
    return small as i32;
}

fun main() {
    var money: i32 = 137;

    println("coins={} remainder={}", money / 50, money % 50);
    println("{} {} {}", checked_byte(0), checked_byte(255), checked_byte(256));
}
```

Execution result:

```text
coins=2 remainder=37
0 255 -1
```

The reason why -1 can be used as a failure indicator is because the success range is 0 to 255. If any integer can be a success value, a different result representation is needed. We continue this design in the error handling chapter later.


## priority

Operator precedence is as follows, starting from highest:

1. Basic expressions and postfix access: function calls, field access, indexing, postfix `++`·`--`
2. Unary operations: `!`, `~`, `&`, `deref`, prefix `++`·`--`, unary `+`·`-`
3. `as` type conversion
4. `*`, `/`, `%`
5. `+`, `-`
6. `<<`, `>>`
7. `<`, `<=`, `>`, `>=`
8. `==`, `!=`
9. Bit `&`
10. Bit `^`
11. Bit `|`
12. `&&`
13. `||`
14. Assignment and compound assignment

Chain assignment combines from the right. When mixing different types of operators, use parentheses to clearly express the order of evaluation.

## Subjects that can be substituted

Assignment, `++`, and `--` require an expression denoting a storage location, such as a variable, field, array element, or dereferenced pointer. Writing to a `const` is not allowed.

## shift

The result of a shift always has the left operand’s type; the right operand’s type does not widen the computation. Left shifts discard high bits beyond that width. Right shifts sign-extend signed values and zero-extend unsigned values.

The shift count must be an integer whose original value satisfies `0 <= n < LHS bit width`. It is checked before any truncation to a smaller type. An invalid constant count is a compile-time error; an invalid runtime count causes a trap.

## Conversions involving bool and floating-point values

Integer-to-bool conversion produces false for zero and true for every other value. Floating-point-to-bool conversion produces false only for +0.0 and -0.0; NaN and positive or negative infinity produce true. Pointer-to-bool casts are unsupported: compare explicitly with `pointer != null`.

Floating-point-to-integer conversion truncates toward zero and then checks the destination range. NaN, infinity, and out-of-range results are invalid. Invalid constant conversions are compile-time errors; invalid runtime conversions cause a trap. A trap is not a recoverable error return.

`&&` and `||` are short circuit evaluations. The side effects of the unexecuted right operand do not occur. You can check the results with small values ​​in [arithmetic class](/docs/en/language/expressions-and-operators).

## shift that fails on purpose

The number of shifts for an 8-bit value must be 0 to 7. The programs below should be rejected before execution.

<!-- wave-example: reject-shift-count -->
```wave
fun main() {
    var value: u8 = 1;
    var result: u8 = value << 8;
}
```
