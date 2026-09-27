---
translation_set_id: functions
path: language/functions-and-generics
locale: ko
group: language
group_order: 2
order: 5
title: 5. 함수를 설계하고 조합하기
summary: 매개변수, 반환값, 기본값과 값 전달을 배웁니다.
---

## 반복되는 코드에서 출발하기

함수는 문법을 줄이기 위한 도구이기도 하지만, 작업의 경계를 정하는 도구입니다. 입력으로 무엇을 받고, 무엇을 계산하며, 어떤 결과를 반환하는지 분리하면 프로그램을 작은 단위로 이해할 수 있습니다.

이번 장에서는 할인 계산을 여러 번 작성하는 프로그램에서 출발합니다. 각 예제는 main.wave 전체이며 `wavec run main.wave`로 실행합니다. 아직 파일을 나누지는 않습니다.

<!-- wave-example: book-function-before -->
```wave
fun main() {
    var first_price: i32 = 2000;
    var first_discount: i32 = first_price * 10 / 100;
    var first_total: i32 = first_price - first_discount;
    var second_price: i32 = 5000;
    var second_discount: i32 = second_price * 10 / 100;
    var second_total: i32 = second_price - second_discount;
    println("{} {}", first_total, second_total);
}
```

실행 결과:

```text
1800 4500
```

두 계산은 가격만 다르고 구조가 같습니다. 할인 규칙을 바꿀 때 두 곳을 모두 수정해야 합니다. 한쪽만 바꾸면 같은 정책을 적용해야 하는 상품에서 서로 다른 결과가 나옵니다.

## 입력과 출력 정하기

중복된 계산을 함수로 옮깁니다. 바뀌는 값은 price라는 입력으로 받고, 계산한 가격은 반환합니다.

<!-- wave-example: book-function-extract -->
```wave
fun discounted(price: i32) -> i32 {
    var discount: i32 = price * 10 / 100;
    return price - discount;
}

fun main() {
    println("{} {}", discounted(2000), discounted(5000));
}
```

실행 결과:

```text
1800 4500
```

함수 이름은 discounted입니다. 괄호 안 `price: i32`는 매개변수 선언이고 `-> i32`는 결과 타입입니다. 본문의 지역 변수 discount는 이 함수 안에서만 사용됩니다.

`discounted(2000)`은 함수를 호출하는 식입니다. 괄호 안 2000은 실제로 전달하는 인자입니다. 함수가 반환한 값이 이 호출식의 결과가 되므로 println의 인자로 바로 사용할 수 있습니다.

| 용어 | 코드 | 뜻 |
| --- | --- | --- |
| 매개변수 | price | 함수를 선언할 때 정한 입력 이름 |
| 인자 | 2000 | 호출할 때 전달한 값 |
| 반환 타입 | i32 | 호출식이 만들어 내는 값의 타입 |
| 반환문 | return price - discount | 결과를 전달하고 이번 호출 종료 |

## 호출과 실행 순서 따라가기

함수 선언을 소스에 적는 것만으로 본문이 즉시 실행되지는 않습니다. main에서 호출한 지점에 도달했을 때 실행합니다.

<!-- wave-example: book-function-trace -->
```wave
fun calculate(value: i32) -> i32 {
    println("inside: {}", value);
    return value * 2;
}

fun main() {
    println("before");
    var result: i32 = calculate(7);
    println("after: {}", result);
}
```

실행 결과:

```text
before
inside: 7
after: 14
```

진행 순서는 main의 첫 출력, calculate 본문, main의 마지막 출력입니다. `return`이 실행되면 calculate 호출이 결과 14로 끝나고 main의 result 초기화가 완료됩니다.

여러 함수 호출이 한 식에 섞여 있고 부작용 순서가 중요하다면 호출을 별도 문장으로 나누십시오. 이 장의 예제도 추적이 필요한 호출 결과를 지역 변수에 저장합니다.

## 여러 매개변수

할인율도 입력으로 받으면 같은 함수로 여러 정책을 계산할 수 있습니다.

<!-- wave-example: book-function-parameters -->
```wave
fun discounted(price: i32, percent: i32) -> i32 {
    return price - price * percent / 100;
}

fun main() {
    var standard: i32 = discounted(2000, 10);
    var special: i32 = discounted(2000, 25);
    println("standard={}", standard);
    println("special={}", special);
}
```

실행 결과:

```text
standard=1800
special=1500
```

인자의 순서가 선언과 맞아야 합니다. 두 매개변수가 모두 i32이면 순서를 바꾸어도 타입 검사만으로 의미를 구별하기 어렵습니다. 함수 이름과 매개변수 이름을 분명히 정하고 호출하는 곳도 읽기 쉽게 작성합니다.

이 함수는 작은 금액과 유효한 할인율을 전제로 합니다. 음수 가격, 100보다 큰 비율, 중간 곱셈 오버플로를 처리하지 않습니다. 함수를 만들 때는 본문뿐 아니라 입력 조건도 설명해야 합니다. 뒤의 완성 프로그램에서는 검사 단계를 분리합니다.

## 기본 인자

자주 사용하는 값을 기본값으로 제공할 수 있습니다.

<!-- wave-example: book-function-default -->
```wave
fun discounted(price: i32, percent: i32 = 10) -> i32 {
    return price - price * percent / 100;
}

fun main() {
    println("{}", discounted(2000));
    println("{}", discounted(2000, 25));
}
```

실행 결과:

```text
1800
1500
```

첫 호출은 두 번째 인자를 생략해 10을 사용합니다. 두 번째 호출은 명시한 25를 사용합니다. 기본값은 뒤쪽의 생략 가능한 매개변수에 둡니다. 첫 인자만 빠뜨리기 위해 빈칸을 적는 문법으로 사용하지 않습니다.

기본값을 바꾸면 생략 호출의 동작이 바뀝니다. 공개 함수의 기본값도 사용자가 의존하는 동작의 일부입니다. 인자를 명시한 호출과 생략한 호출을 각각 시험하는 이유입니다.

## 값으로 전달한다는 뜻

정수 값을 전달하면 함수가 받는 값과 호출자의 변수 저장 공간이 구분됩니다. 결과를 계산했다고 호출자의 변수가 자동으로 변경되지는 않습니다.

<!-- wave-example: book-function-value -->
```wave
fun next(value: i32) -> i32 {
    return value + 1;
}

fun main() {
    var count: i32 = 4;
    var later: i32 = next(count);
    println("count={} later={}", count, later);
    count = next(count);
    println("count={}", count);
}
```

실행 결과:

```text
count=4 later=5
count=5
```

첫 호출은 count를 읽어서 later를 초기화합니다. count는 여전히 4입니다. 두 번째 호출 뒤에는 결과를 count에 대입했으므로 5가 됩니다. 값을 반환하는 설계는 데이터가 어디서 변경되는지 호출부에서 드러나게 합니다.

원본 저장 공간을 함수 안에서 바꾸고 싶다면 포인터를 전달할 수 있습니다. 이는 [포인터 장](/docs/ko/language/explicit-memory-type-model)에서 다룹니다. 포인터를 전달해도 포인터 값 자체와 그 주소의 저장 공간은 구분해야 합니다.

## 반환하지 않는 경로를 남기지 않기

값을 반환하는 함수는 필요한 모든 경로에서 결과를 제공해야 합니다. 다음처럼 마지막 경로를 남겨두지 않습니다.

```wave
// 의도적으로 잘못된 함수 조각

fun positive(value: i32) -> i32 {
    if (value > 0) {
        return value;
    }
}
```

value가 0 이하이면 어떤 값을 반환할지 정해져 있지 않습니다. 의도한 규칙을 정한 다음 모든 경로를 작성해야 합니다.

<!-- wave-example: book-function-paths -->
```wave
fun positive_or_zero(value: i32) -> i32 {
    if (value > 0) {
        return value;
    }

    return 0;
}

fun main() {
    println("{} {} {}", positive_or_zero(-2), positive_or_zero(0), positive_or_zero(8));
}
```

실행 결과:

```text
0 0 8
```

첫 return을 실행한 호출에서는 아래 return으로 진행하지 않습니다. 조건이 거짓일 때만 마지막 return에 도달합니다. 양수, 0, 음수의 세 경우로 규칙을 확인했습니다.

## 결과가 없는 함수

출력 같은 작업만 수행한다면 반환 타입을 생략할 수 있습니다. 결과가 없는 함수도 `return;`으로 일찍 끝낼 수 있습니다.

<!-- wave-example: book-function-void -->
```wave
fun show_positive(value: i32) {
    if (value <= 0) {
        return;
    }

    println("positive={}", value);
}

fun main() {
    show_positive(-3);
    show_positive(6);
}
```

실행 결과:

```text
positive=6
```

첫 호출은 아무것도 출력하지 않고 돌아옵니다. 두 번째 호출은 출력합니다. “결과값이 없음”과 “호출 지점으로 돌아오지 않음”은 다릅니다. 프로세스를 종료하는 함수처럼 돌아오지 않는 함수는 반환 타입 `!`로 구분합니다.

## 여러 함수로 완성 프로그램 구성

이제 입력 검증, 계산, 출력을 서로 다른 함수로 나눕니다. 가격 범위를 제한했으므로 이 예제의 중간 곱셈은 i32 범위 안입니다.

<!-- wave-example: book-function-program -->
```wave
fun valid_order(price: i32, quantity: i32, percent: i32) -> bool {
    return price >= 0 && price <= 100000
        && quantity >= 1 && quantity <= 100
        && percent >= 0 && percent <= 100;
}

fun discounted_unit(price: i32, percent: i32) -> i32 {
    return price - price * percent / 100;
}

fun order_total(price: i32, quantity: i32, percent: i32) -> i32 {
    var unit: i32 = discounted_unit(price, percent);
    return unit * quantity;
}

fun show_order(price: i32, quantity: i32, percent: i32) {
    if (!valid_order(price, quantity, percent)) {
        println("invalid order");
        return;
    }

    var total: i32 = order_total(price, quantity, percent);
    println("total={}", total);
}

fun main() {
    show_order(2000, 3, 10);
    show_order(2000, 0, 10);
    show_order(2000, 3, 120);
}
```

실행 결과:

```text
total=5400
invalid order
invalid order
```

함수마다 질문 하나에 답합니다. valid_order는 입력이 허용되는가, discounted_unit은 한 개의 할인 가격이 얼마인가, order_total은 전체 가격이 얼마인가, show_order는 무엇을 보여줄 것인가를 담당합니다.

작은 함수가 무조건 좋은 것은 아닙니다. 한 식마다 이름을 붙이면 오히려 추적이 어려울 수 있습니다. 다른 곳에서 재사용할 의미가 있거나, 독립적으로 설명·검증할 규칙이 있을 때 분리합니다.

## 연습 문제

1. 정수 두 개 중 큰 값을 반환하는 maximum을 작성하십시오.
2. 두 경계 사이의 값을 돌려주는 clamp 함수를 작성하십시오. 아래 풀이에서는 low <= high를 호출 조건으로 정합니다.
3. 값에 세금을 더하는 함수를 만들어 주문 합계 함수와 조합하십시오. 범위와 정수 절삭 시점을 먼저 정하십시오.

### 풀이: 경계를 갖는 함수

<!-- wave-example: book-function-solution -->
```wave
fun maximum(left: i32, right: i32) -> i32 {
    if (left > right) {
        return left;
    }

    return right;
}

fun clamp(value: i32, low: i32, high: i32) -> i32 {
    if (value < low) {
        return low;
    }
    if (value > high) {
        return high;
    }

    return value;
}

fun main() {
    println("max={}", maximum(7, 4));
    println("{} {} {}", clamp(-3, 0, 10), clamp(6, 0, 10), clamp(20, 0, 10));
}
```

실행 결과:

```text
max=7
0 6 10
```

clamp는 범위보다 작음, 범위 안, 범위보다 큼의 세 경로가 있습니다. 경계값 0과 10도 직접 추가해 확인하십시오. low > high인 경우까지 지원하려면 실패를 어떻게 표현할지 정해야 합니다. [오류 처리 장](/docs/ko/language/errors)에서 결과 구조체를 사용해 이 문제를 다룹니다.


## 재귀 호출

함수는 자기 자신을 호출할 수 있습니다. 재귀에서는 더 이상 호출하지 않는 종료 조건과, 매 호출이 그 조건에 가까워지는 과정이 필요합니다.

<!-- wave-example: book-ref-recursion -->
```wave
fun factorial(value: i32) -> i32 {
    if (value <= 1) {
        return 1;
    }

    return value * factorial(value - 1);
}

fun main() {
    println("{}", factorial(5));
}
```

실행 결과:

```text
120
```

5의 계산은 `5 * factorial(4)`로 이어지고, 1에 도달하면 1을 반환합니다. 반환값이 이전 호출로 차례로 전달되어 120이 됩니다. 이 함수는 작은 양의 정수를 설명하기 위한 예제입니다. 큰 입력에서는 결과 범위와 호출 깊이를 고려해야 합니다. 반복문으로 같은 작업을 작성하면 호출 깊이가 늘어나는 문제를 피할 수 있습니다.

## 자주 생기는 오류

| 현상 | 확인할 것 |
| --- | --- |
| 인자가 부족하거나 많다는 오류 | 매개변수 수와 생략 가능한 기본값 |
| 반환 타입이 맞지 않음 | return 식의 타입과 함수 선언 |
| 특정 경로에서 반환하지 않음 | 조건이 거짓인 경우까지 반환하는지 |
| 제네릭 타입 인자 누락 | 함수 이름 뒤 `<실제 타입>` |
| 포인터 함수 호출 뒤 원본이 바뀜 | 값을 읽기만 하는 함수인지 수정하는 함수인지 |

`export(c)`처럼 외부에 내보내는 함수는 구체적인 서명을 사용합니다. 제네릭 함수 자체를 외부 호출 규약으로 내보낼 수는 없습니다. `ptr<T>`와 `array<T, N>`은 언어의 내장 메모리 타입이며 사용자 제네릭 구조체 선언과 구분합니다.

[함수 학습](/docs/ko/language/functions-and-generics) · [모듈과 제네릭 학습](/docs/ko/language/modules-imports-and-ffi) · [FFI](/docs/ko/language/modules-imports-and-ffi)
