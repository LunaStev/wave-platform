---
translation_set_id: stdlib-math-debug
path: stdlib/math-debug
locale: ko
group: stdlib
group_order: 1
order: 14
title: math와 debug: 계산 보조와 진단
summary: 범위 검사가 있는 수학 함수와 진단 출력을 사용합니다.
---

## 범위를 검사하는 정수 함수

```text
std::math::int
min_i32(a: i32, b: i32) -> i32
max_i32(a: i32, b: i32) -> i32
abs_i32_checked(x: i32) -> MathResult<i32>
clamp_i32_checked(x: i32, lo: i32, hi: i32) -> MathResult<i32>
div_floor_i32_checked(a: i32, b: i32) -> MathResult<i32>
div_ceil_i32_checked(a: i32, b: i32) -> MathResult<i32>
```

`MathResult<T>`는 `value`와 `error`를 가집니다. `error == MATH_ERROR_NONE`을 확인한 뒤 value를 사용합니다. 오류 상수는 `std::math::result`에서 가져옵니다. 최소 signed 정수의 절댓값은 같은 signed 타입으로 표현할 수 없습니다. checked 절댓값은 이를 실패로 알립니다. clamp는 lo가 hi보다 큰 잘못된 범위를 거부합니다. 나눗셈은 0인 분모와 표현 범위를 검사합니다. floor와 ceil은 언어의 0 방향 절삭 나눗셈과 반올림 방향이 다릅니다.

## 실수 분류

`std::math::float`의 `is_nan_f64`, `is_infinite_f64`, `is_finite_f64`는 특수 값을 구분합니다. f32 함수도 제공합니다. `float_to_bits_f64(value: f64) -> u64`는 저장 비트를 얻는 함수이며 `value as u64`의 수치 변환과 다릅니다. NaN은 자기 자신과도 같지 않으므로 `value == nan`으로 검사하지 않습니다.

## 진단

`std::debug::core`의 `debug_assert(condition: bool, message: str)`는 거짓 조건에서 진단 후 종료합니다. 사용자 입력처럼 정상적으로 실패할 수 있는 상황은 결과값으로 처리하고, 반드시 성립해야 할 프로그램 내부 조건을 확인할 때 assert를 사용합니다.

아래 프로그램을 `main.wave`로 저장하여 실행합니다.

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

실행 결과:

```text
4 9
```

## 음수 나눗셈의 반올림 방향

-7을 3으로 나누는 방법을 비교합니다. `/`는 0 쪽으로 잘라 -2가 됩니다. floor는 더 작은 정수 -3, ceil은 더 큰 정수 -2를 선택합니다.

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

실행 결과:

```text
truncate=-2 floor=-3 ceil=-2
absolute value is outside i32
```

결과값만 확인하면 실패 시 들어 있는 대체 값과 실제 계산 결과를 구분할 수 없습니다. error를 먼저 검사하는 순서를 지킵니다. 음수 좌표를 일정한 크기의 구간에 넣을 때는 floor, 필요한 묶음 수를 올림할 때는 ceil이 유용합니다.

## 어떤 오류를 처리해야 하나

| 상황 | 오류 | 처리 예 |
| --- | --- | --- |
| 0으로 나누기 | `MATH_ERROR_DIVIDE_BY_ZERO` | 분모 입력을 다시 받음 |
| 결과를 타입에 담을 수 없음 | `MATH_ERROR_OVERFLOW` | 더 넓은 타입으로 계산하거나 입력을 거부 |
| clamp에서 최솟값이 최댓값보다 큼 | `MATH_ERROR_INVALID_ARGUMENT` | 설정 범위를 수정 |

assert는 오류를 복구하는 도구가 아닙니다. 사용자 입력의 오류는 조건문과 반환값으로 처리하고, 계산을 마친 뒤 반드시 성립해야 할 내부 조건을 debug_assert로 확인합니다.
