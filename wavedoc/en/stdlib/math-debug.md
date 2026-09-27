---
translation_set_id: stdlib-math-debug
path: stdlib/math-debug
locale: en
group: stdlib
group_order: 1
order: 14
title: math and debug: Mathematical utilities and diagnostics
summary: Use mathematical functions with range checking and diagnostic output.
---

## Integer function to check range

```text
std::math::int
min_i32(a: i32, b: i32) -> i32
max_i32(a: i32, b: i32) -> i32
abs_i32_checked(x: i32) -> MathResult<i32>
clamp_i32_checked(x: i32, lo: i32, hi: i32) -> MathResult<i32>
div_floor_i32_checked(a: i32, b: i32) -> MathResult<i32>
div_ceil_i32_checked(a: i32, b: i32) -> MathResult<i32>
```

`MathResult<T>` contains `value` and `error`. Use value only after checking `error == MATH_ERROR_NONE`; import the error constants from `std::math::result`. The absolute value of the minimum signed integer cannot be represented in the same type, so the checked absolute-value function reports failure. clamp rejects lo greater than hi. Division checks for a zero divisor and range overflow. floor and ceil round differently from the language’s integer division, which truncates toward zero.

## Classifying floating-point values

`is_nan_f64`, `is_infinite_f64`, and `is_finite_f64` in `std::math::float` distinguish special values. The f32 function is also provided. `float_to_bits_f64(value: f64) -> u64` is a function to obtain storage bits and is different from the numeric conversion of `value as u64`. NaN is not even equal to itself, so it is not checked against `value == nan`.

## diagnosis

`debug_assert(condition: bool, message: str)` in `std::debug::core` terminates after diagnosing a false condition. Situations that may normally fail, such as user input, are handled with the result value, and assert is used when checking internal program conditions that must be met.

Save the program below as `main.wave` and run it.

<!-- wave-example: math-api -->
```wave
import("std::math::int")::{
    min_i32, max_i32
};
import("std::debug::core")::{
    debug_assert
};

fun main() {
    var low: i32 = min_i32(9, 4);
    var high: i32 = max_i32(9, 4);
    debug_assert(low <= high, "invalid range");
    println("{} {}", low, high);
}
```

Execution result:

```text
4 9
```

## Rounding direction for negative division

Compare how to divide -7 by 3. `/` is truncated toward 0 to become -2. floor selects the smaller integer -3, and ceil selects the larger integer -2.

<!-- wave-example: book-math-rounding -->
```wave
import("std::math::int")::{
    div_floor_i32_checked,
    div_ceil_i32_checked,
    abs_i32_checked
};
import("std::math::result")::{
    MathResult,
    MATH_ERROR_NONE,
    MATH_ERROR_OVERFLOW
};

fun main() -> i32 {
    var floor: MathResult<i32> = div_floor_i32_checked(-7, 3);
    var ceil: MathResult<i32> = div_ceil_i32_checked(-7, 3);

    if (floor.error != MATH_ERROR_NONE || ceil.error != MATH_ERROR_NONE) {
        return 1;
    }

    println("truncate={} floor={} ceil={}", -7 / 3, floor.value, ceil.value);

    var absolute: MathResult<i32> = abs_i32_checked(-2147483648);

    if (absolute.error == MATH_ERROR_OVERFLOW) {
        println("absolute value is outside i32");
    }

    return 0;
}
```

Execution result:

```text
truncate=-2 floor=-3 ceil=-2
absolute value is outside i32
```

If you only check the result value, you cannot distinguish between the substitute value included in case of failure and the actual calculation result. Follow the order of checking error first. floor is useful when putting negative coordinates into an interval of a certain size, and ceil is useful when rounding up the required number of bundles.

## What errors should I handle?

|situation|error|Processing example|
| --- | --- | --- |
|Divide by zero| `MATH_ERROR_DIVIDE_BY_ZERO` |Takes denominator input again|
|Result cannot be stored in type| `MATH_ERROR_OVERFLOW` |Calculate with wider type or reject input|
|The minimum is greater than the maximum in clamp| `MATH_ERROR_INVALID_ARGUMENT` |Modify setting range|

assert is not an error repair tool. Errors in user input are handled using conditional statements and return values, and after completing the calculation, the internal conditions that must be met are checked with debug_assert.
