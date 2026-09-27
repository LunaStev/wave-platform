---
translation_set_id: stdlib-math-debug
path: stdlib/math-debug
locale: ru
group: stdlib
group_order: 1
order: 14
title: math и debug: Математические функции и диагностика
summary: Используйте математические функции для проверки диапазона и вывода диагностики.
---

## Целочисленная функция для проверки диапазона

```text
std::math::int
min_i32(a: i32, b: i32) -> i32
max_i32(a: i32, b: i32) -> i32
abs_i32_checked(x: i32) -> MathResult<i32>
clamp_i32_checked(x: i32, lo: i32, hi: i32) -> MathResult<i32>
div_floor_i32_checked(a: i32, b: i32) -> MathResult<i32>
div_ceil_i32_checked(a: i32, b: i32) -> MathResult<i32>
```

`MathResult<T>` содержит `value` и `error`. Импортируйте константы ошибок из `std::math::result` и используйте `value` только после проверки `error == MATH_ERROR_NONE`. Модуль минимального знакового целого не представим в том же типе, поэтому функция вычисления модуля с проверкой сообщает об ошибке. `clamp` отклоняет `lo > hi`. Деление проверяет нулевой делитель и переполнение диапазона. `floor` и `ceil` округляют иначе, чем целочисленное деление, которое отбрасывает дробную часть в направлении нуля.

## Классификация значений с плавающей запятой

`is_nan_f64`, `is_infinite_f64` и `is_finite_f64` в `std::math::float` различают специальные значения. Также предусмотрена функция f32. `float_to_bits_f64(value: f64) -> u64` — это функция получения битов памяти, которая отличается от числового преобразования `value as u64`. NaN даже не равен самому себе, поэтому он не сверяется с `value == nan`.

## диагноз

`debug_assert(condition: bool, message: str)` в `std::debug::core` завершается после диагностики ложного состояния. Ситуации, которые обычно могут привести к сбою, например ввод данных пользователем, обрабатываются с помощью значения результата, а assert используется при проверке внутренних условий программы, которые должны быть выполнены.

Сохраните приведенную ниже программу как `main.wave` и запустите ее.

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

Результат выполнения:

```text
4 9
```

## Направление округления для отрицательного деления

Сравните, как разделить -7 на 3. `/` усекается в сторону 0 и становится -2. floor выбирает меньшее целое число -3, а ceil выбирает большее целое число -2.

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

Результат выполнения:

```text
truncate=-2 floor=-3 ceil=-2
absolute value is outside i32
```

Если вы проверяете только значение результата, вы не сможете отличить замещающее значение, включенное в случае сбоя, от фактического результата расчета. Сначала следуйте порядку проверки error. floor полезен при помещении отрицательных координат в интервал определенного размера, а ceil полезен при округлении необходимого количества связок.

## Какие ошибки мне следует обрабатывать?

|ситуация|ошибка|Пример обработки|
| --- | --- | --- |
|Разделить на ноль| `MATH_ERROR_DIVIDE_BY_ZERO` |Снова принимает знаменатель|
|Результат не может быть сохранен в типе| `MATH_ERROR_OVERFLOW` |Вычислить более широкий тип или отклонить ввод|
|Минимум больше максимума в clamp| `MATH_ERROR_INVALID_ARGUMENT` |Изменить диапазон настроек|

assert не является инструментом исправления ошибок. Ошибки пользовательского ввода обрабатываются с помощью условных операторов и возвращаемых значений, а после завершения вычислений внутренние условия, которые должны быть выполнены, проверяются с помощью debug_assert.
