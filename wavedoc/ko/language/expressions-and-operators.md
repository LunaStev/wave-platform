---
translation_set_id: expressions
path: language/expressions-and-operators
locale: ko
group: language
group_order: 2
order: 3
title: 3. 계산, 비교와 형 변환
summary: 계산 순서, 정수 나눗셈, 비트 연산과 cast를 배웁니다.
---

## 계산 결과와 계산 타입을 함께 보기

표현식은 값을 계산하는 코드입니다. 변수 이름, 리터럴, 함수 호출, 여러 값을 연산자로 연결한 식이 모두 표현식입니다. 수학에서 같은 식처럼 보여도 정수인지 실수인지, 몇 비트인지에 따라 결과가 달라집니다.

이번 장에서는 단순 계산에서 출발해 괄호, 나눗셈, 논리 연산, 비트 연산과 cast를 배웁니다. 예제는 각각 main.wave 전체이며 `wavec run main.wave`로 실행합니다.

## 괄호로 묶는 범위

<!-- wave-example: book-precedence -->
```wave
fun main() {
    var first: i32 = 2 + 3 * 4;
    var second: i32 = (2 + 3) * 4;

    println("{} {}", first, second);
}
```

실행 결과:

```text
14 20
```

곱셈이 덧셈보다 먼저 계산되어 첫 식은 2+12가 됩니다. 두 번째 식에서는 괄호 안의 합 5를 구한 다음 4를 곱합니다. 괄호를 적게 쓰는 것이 목표는 아닙니다. 읽는 사람이 계산 범위를 쉽게 알 수 있도록 사용하는 것이 좋습니다.

같은 연산자를 여러 번 쓸 때도 묶이는 방향이 중요합니다. `20 - 5 - 3`은 `(20 - 5) - 3`이므로 12입니다. `20 - (5 - 3)`은 18입니다. 정확한 전체 순서는 [연산자 참조](/docs/ko/language/expressions-and-operators)에 있습니다.

## 정수 나눗셈과 나머지

정수끼리 나누면 소수부를 가진 실수 결과가 생기지 않습니다. 몫과 나머지를 각각 구할 수 있습니다.

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

실행 결과:

```text
boxes=3 remaining=2
negative quotient=-3
```

17개를 5개씩 묶으면 완성된 상자 3개와 남은 2개가 있습니다. signed 나눗셈은 0 방향으로 절삭하므로 -17/5는 -3입니다. 바닥 함수처럼 항상 작은 정수 쪽으로 내리는 것과 다릅니다.

0으로 나눌 수는 없습니다. signed 최소값을 -1로 나눈 값도 같은 타입에 들어가지 않습니다. 이런 입력을 받는 함수는 나누기 전에 검사하거나 checked 수학 API를 사용해야 합니다.

## 변환하는 시점이 결과를 바꿈

다음 두 식은 모두 f64 변수에 저장되지만 계산 과정이 다릅니다.

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

실행 결과:

```text
integer division lost the fraction
floating division kept the fraction
```

첫 식은 정수 몫 3을 얻은 뒤 f64로 바꿉니다. 둘째 식은 피연산자부터 f64로 바꾸어 실수 나눗셈을 합니다. 마지막 변수 타입만 넓게 정해도 앞 단계에서 잃은 정보가 되살아나는 것은 아닙니다.

실수는 근삿값입니다. 계산 결과가 눈으로 같은 십진수처럼 보여도 정확한 동등 비교가 적절하지 않을 수 있습니다. 필요한 오차 범위는 문제의 단위와 크기에 맞춰 정합니다. 모든 계산에 같은 고정 epsilon을 붙이는 것도 해결책은 아닙니다.

## 비교는 bool을 만듦

`<`, `<=`, `>`, `>=`, `==`, `!=`는 관계를 검사합니다. 등호 하나 `=`는 대입이고 둘 `==`는 동등 비교입니다.

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

실행 결과:

```text
eligible
not exactly the boundary
```

“18 이상”과 “18 초과”는 경계에서 다릅니다. 코드가 맞는지 확인할 때 17, 18, 19를 각각 넣는 이유입니다. 비교에 서로 다른 타입이 섞이면 부호와 폭에 영향을 받으므로 의도한 타입으로 맞춘 뒤 비교하는 편이 명확합니다.

## 논리 연산과 단락 평가

`&&`는 둘 다 참인지, `||`는 하나 이상 참인지 검사하고 `!`는 참·거짓을 뒤집습니다. 이 연산은 오른쪽 식을 항상 실행하지 않습니다.

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

실행 결과:

```text
at least one true
```

첫 조건에서는 enabled가 거짓이므로 오른쪽을 볼 필요가 없습니다. 두 번째에서는 !enabled가 참이므로 역시 오른쪽이 필요 없습니다. 따라서 report의 출력이 한 번도 나타나지 않습니다.

이를 이용해 나눗셈 전에 분모를 확인할 수 있습니다. 함수 본문 조각 `if (divisor != 0 && value / divisor > 2) { ... }`에서는 분모가 0일 때 나눗셈을 수행하지 않습니다. 단, 이 검사만으로 signed 최소값/-1 같은 다른 경계까지 해결되는 것은 아닙니다.

`&&`가 `||`보다 우선합니다. 복합 정책에서는 `(member && active) || admin`처럼 괄호로 의도를 드러내십시오.

## 정수를 좁히거나 넓히기

`as`는 명시적 형 변환입니다. 정수를 좁힐 때 버리는 상위 비트는 다시 넓힌다고 복구되지 않습니다.

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

실행 결과:

```text
300 -> 44 -> 44
-1 255
```

300의 낮은 8비트는 44입니다. -1을 signed 타입으로 넓히면 부호를 확장하여 -1을 유지합니다. 같은 8비트를 unsigned로 해석하면 255입니다.

범위가 맞는 값을 보존하려는 변환과 저장 비트를 다루려는 변환은 목적이 다릅니다. 사용자 입력을 받았다면 먼저 목적지 범위인지 확인하고 변환하십시오. cast가 있다는 이유로 값이 안전한 범위였다고 보장되지 않습니다.

## bool로 바꾸기

정수는 0이면 false, 그 외는 true입니다. 낮은 1비트만 남기는 변환이 아니므로 2도 true입니다. 실수는 +0.0과 -0.0만 false이고 NaN·무한대를 포함한 나머지는 true입니다.

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

실행 결과:

```text
zero is false
two is true
```

포인터를 bool로 cast하지 말고 `pointer != null`처럼 비교합니다. null이 아닌 주소라고 실제로 읽어도 안전한지는 별개의 문제입니다.

## 비트 연산

`&`, `|`, `^`, `~`는 정수의 각 비트를 다룹니다. 권한이나 기능을 비트로 표현할 때 사용할 수 있습니다. 아래에서는 1이 읽기 권한, 2가 쓰기 권한입니다.

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

실행 결과:

```text
write enabled
remaining=1
```

OR로 비트를 추가하고 AND로 특정 비트가 있는지 검사합니다. `~WRITE`로 해당 비트만 0인 마스크를 만들어 지웁니다. 비트 연산의 `&`·`|`는 bool의 단락 평가 `&&`·`||`와 다른 연산자입니다.

## 시프트의 폭과 횟수

왼쪽 시프트는 비트를 왼쪽으로 이동하고 폭 밖의 상위 비트를 버립니다. 오른쪽 시프트는 signed이면 부호를, unsigned이면 0을 채웁니다. 결과 타입은 항상 왼쪽 피연산자의 타입입니다.

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

실행 결과:

```text
left=2
right=-4 4
```

u8의 129를 왼쪽으로 한 칸 옮기면 가장 높은 비트가 버려지고 2가 남습니다. 이동 횟수는 원래 값이 0 이상이고 왼쪽 타입의 비트 폭보다 작아야 합니다. u8에서는 0~7입니다. 잘못된 상수 횟수는 컴파일 오류, 실행 중 잘못된 횟수는 trap입니다.

## 실수에서 정수로

소수부를 0 방향으로 절삭한 뒤 목적지 정수 범위를 검사합니다. NaN·무한대와 범위를 벗어나는 결과는 유효하지 않습니다. 상수에서 확인되는 잘못된 변환은 컴파일 오류이고 실행 중에는 trap입니다.

trap은 함수가 오류값을 돌려주는 방식이 아닙니다. 복구 가능한 프로그램을 작성하려면 변환 전에 유효 범위를 검사하는 인터페이스를 설계해야 합니다. 실수의 저장 비트를 얻으려는 목적이라면 수치 cast 대신 `std::math::float`의 bits 함수를 사용합니다.

## 연습과 풀이

잔돈 137원을 50원 단위와 나머지로 나누고, 0~255 범위의 정수만 u8로 바꾸는 함수를 작성하십시오. 이 예제에서는 범위 밖을 -1로 표시하는 별도 함수로 검증 단계를 보여 줍니다.

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

실행 결과:

```text
coins=2 remainder=37
0 255 -1
```

-1을 실패 표지로 쓸 수 있는 것은 성공 범위가 0~255이기 때문입니다. 모든 정수가 성공값일 수 있다면 다른 결과 표현이 필요합니다. 뒤의 오류 처리 장에서 이 설계를 이어갑니다.


## 우선순위

연산자 우선순위는 높은 순서부터 다음과 같습니다.

1. 기본식과 후위 접근: 함수 호출, 필드 접근, 인덱싱, 후위 `++`·`--`
2. 단항 연산: `!`, `~`, `&`, `deref`, 전위 `++`·`--`, 단항 `+`·`-`
3. `as` 형 변환
4. `*`, `/`, `%`
5. `+`, `-`
6. `<<`, `>>`
7. `<`, `<=`, `>`, `>=`
8. `==`, `!=`
9. 비트 `&`
10. 비트 `^`
11. 비트 `|`
12. `&&`
13. `||`
14. 대입과 복합 대입

연쇄 대입은 오른쪽부터 결합합니다. 여러 종류의 연산자를 섞을 때는 괄호로 계산 순서를 명확히 표현하십시오.

## 대입 가능한 대상

대입과 `++`·`--`는 변수, 필드, 배열 인덱스, 역참조처럼 저장 위치를 나타내는 식에 사용합니다. `const`에 대한 쓰기는 허용되지 않습니다.

## 시프트

시프트 결과 타입은 항상 왼쪽 피연산자 타입입니다. 오른쪽의 타입이 계산 폭을 넓히지 않습니다. 왼쪽 시프트는 왼쪽 타입의 폭에서 넘친 상위 비트를 버립니다. 오른쪽 시프트는 signed이면 부호를, unsigned이면 0을 채웁니다.

이동 횟수는 정수이며 원래 값이 `0 <= n < 왼쪽 타입의 비트 폭`이어야 합니다. 작은 타입으로 잘라서 검사하지 않습니다. 잘못된 상수 횟수는 컴파일 오류이고 실행 중 잘못된 횟수는 trap입니다.

## bool과 실수 변환

정수→bool은 0일 때 false, 나머지는 true입니다. 실수→bool은 +0.0과 -0.0만 false이며 NaN과 무한대도 true입니다. 포인터→bool cast는 지원하지 않으므로 `pointer != null`처럼 비교합니다.

실수→정수는 0 방향 절삭 후 목적지 범위를 검사합니다. NaN·무한대와 범위 밖 결과는 유효하지 않습니다. 잘못된 상수 변환은 컴파일 오류, 실행 중 변환은 trap입니다. trap을 복구 가능한 오류 반환처럼 사용하지 마십시오.

`&&`와 `||`는 단락 평가합니다. 실행되지 않는 오른쪽 피연산자의 부작용은 발생하지 않습니다. [연산 수업](/docs/ko/language/expressions-and-operators)에서 작은 값으로 결과를 확인할 수 있습니다.

## 의도적으로 실패하는 시프트

8비트 값의 이동 횟수는 0~7이어야 합니다. 아래 프로그램은 실행 전에 거부되어야 합니다.

<!-- wave-example: reject-shift-count -->
```wave
fun main() {
    var value: u8 = 1;
    var result: u8 = value << 8;
}
```
