---
translation_set_id: stdlib-math-debug
path: stdlib/math-debug
locale: zh
group: stdlib
group_order: 1
order: 14
title: math 和 debug: 数学工具与诊断
summary: 使用数学函数进行范围检查和诊断输出。
---

## 检查范围的整数函数

```text
std::math::int
min_i32(a: i32, b: i32) -> i32
max_i32(a: i32, b: i32) -> i32
abs_i32_checked(x: i32) -> MathResult<i32>
clamp_i32_checked(x: i32, lo: i32, hi: i32) -> MathResult<i32>
div_floor_i32_checked(a: i32, b: i32) -> MathResult<i32>
div_ceil_i32_checked(a: i32, b: i32) -> MathResult<i32>
```

`MathResult<T>` 包含 `value` 和 `error`。从 `std::math::result` 导入错误常量，确认 `error == MATH_ERROR_NONE` 后才能使用 `value`。有符号整数最小值的绝对值无法用同一类型表示，因此带检查的绝对值函数会报告失败。`clamp` 拒绝 `lo > hi`。除法会检查除数为零及结果溢出的情况。`floor` 和 `ceil` 的取整方式不同于向零截断的整数除法。

## 对浮点值进行分类

`std::math::float`中的`is_nan_f64`、`is_infinite_f64`、`is_finite_f64`区分特殊值。还提供f32功能。 `float_to_bits_f64(value: f64) -> u64`是获取存储位的函数，与`value as u64`的数值转换不同。 NaN 甚至不等于其自身，因此不会针对 `value == nan` 进行检查。

## 诊断

`std::debug::core` 中的`debug_assert(condition: bool, message: str)` 在诊断出错误条件后终止。通常可能失败的情况，例如用户输入，通过结果值进行处理，并且在检查必须满足的内部程序条件时使用assert。

将下面的程序保存为`main.wave`并运行。

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

执行结果：

```text
4 9
```

## 负除法的舍入方向

比较如何将 -7 除以 3。`/` 向 0 截断，变成 -2。 floor 选择较小的整数 -3，ceil 选择较大的整数 -2。

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

执行结果：

```text
truncate=-2 floor=-3 ceil=-2
absolute value is outside i32
```

如果只检查结果值，则无法区分失败时包含的替代值和实际计算结果。按照先检查error的顺序。当将负坐标放入一定大小的区间时，floor很有用，而当向上舍入所需的束数时，ceil很有用。

## 我应该处理哪些错误？

|情况|错误|加工实例|
| --- | --- | --- |
|除以零| `MATH_ERROR_DIVIDE_BY_ZERO` |再次输入分母|
|结果无法存储在类型中| `MATH_ERROR_OVERFLOW` |使用更宽的类型进行计算或拒绝输入|
|clamp 中的最小值大于最大值| `MATH_ERROR_INVALID_ARGUMENT` |修改设定范围|

assert不是错误修复工具。使用条件语句和返回值处理用户输入的错误，完成计算后，使用debug_assert检查必须满足的内部条件。
