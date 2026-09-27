---
translation_set_id: types
path: language/declarations-and-types
locale: ko
group: language
group_order: 2
order: 2
title: 2. 변수, 타입과 스코프
summary: 지역 변수, 정수 범위, bool과 스코프를 배웁니다.
---

## 값을 이름으로 다루기

가격을 여러 곳에 직접 적으면 가격을 바꿀 때 모든 위치를 찾아야 합니다. 변수는 값에 이름을 붙이고 그 이름을 통해 읽고 변경하는 저장 공간입니다. 이번 장에서는 선언, 대입, 타입의 범위와 블록 안에서 이름이 보이는 범위를 배웁니다.

아래 프로그램들은 각각 별도의 main.wave 전체입니다. 예제 하나를 저장한 뒤 `wavec run main.wave`로 실행하고 다음 예제로 교체하십시오.

## 선언과 초기화

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

실행 결과:

```text
price=1200
quantity=3
total=3600
```

선언을 네 부분으로 나눠 읽습니다.

| 부분 | 이 예제 | 역할 |
| --- | --- | --- |
| 선언 키워드 | var | 지역 변수 생성 |
| 이름 | price | 이후 사용할 식별자 |
| 타입 | i32 | 저장할 값의 종류와 범위 |
| 초기값 | 1200 | 처음 저장할 값 |

타입 앞의 콜론과 초기값 앞의 등호는 역할이 다릅니다. 이름에는 의미를 담으십시오. 이 예제에서 price는 단가이고 quantity는 개수입니다. 같은 i32라도 서로 바꿔 사용하면 문법 오류 없이 잘못된 계산을 할 수 있습니다.

## 대입은 관계를 유지하는 수식이 아님

변수에 계산 결과를 저장하면 그 시점의 값이 들어갑니다. 계산식을 기억해서 나중에 자동으로 다시 평가하지 않습니다.

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

실행 결과:

```text
before recalculation=3600
after recalculation=6000
```

`quantity = 5`는 이미 존재하는 변수에 새 값을 씁니다. `var quantity`처럼 다시 선언하는 것과 구분해야 합니다. `total`도 다시 대입하기 전에는 3600입니다. 여러 변수 사이의 관계를 프로그램이 유지해야 한다면 관계가 바뀔 때 계산을 수행하도록 작성해야 합니다.

## 이전 값으로 다음 값 계산하기

대입문의 오른쪽을 먼저 계산하고 그 결과를 왼쪽 저장 공간에 씁니다. 수학의 등식과 달리 `count = count + 1`은 유효한 갱신입니다.

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

실행 결과:

```text
9
```

값의 진행은 0 → 1 → 3 → 9입니다. `+=`, `*=`는 계산과 저장을 함께 표현합니다. 계산 순서를 따라 적어 보면 결과가 예상과 다를 때 어느 단계에서 생각이 달랐는지 찾을 수 있습니다.

## 정수 타입의 폭과 부호

`i`로 시작하면 signed, `u`로 시작하면 unsigned 정수입니다. 뒤 숫자는 비트 수입니다. 비트 수가 커지면 표현할 수 있는 범위가 늘어나고 저장 공간도 늘어납니다.

| 타입 | 최솟값 | 최댓값 | 예시 용도 |
| --- | --- | --- | --- |
| i8 | -128 | 127 | 작은 부호 있는 값 |
| u8 | 0 | 255 | 한 바이트 |
| i16 | -32768 | 32767 | 작은 정수 데이터 |
| u16 | 0 | 65535 | 포트·16비트 필드 |
| i32 | -2147483648 | 2147483647 | 일반적인 작은 정수 계산 |
| u32 | 0 | 4294967295 | 32비트 비트 필드 |

64·128·256·512·1024비트 signed/unsigned 정수도 있습니다. “넓은 정수가 있으니 아무 계산이나 안전하다”는 뜻은 아닙니다. 선택한 폭에서 계산 범위를 넘을 수 있으므로 필요한 범위를 먼저 생각합니다. `isz`와 `usz`는 컴파일 대상의 주소 크기를 따릅니다.

큰 리터럴을 작은 타입에 저장하는 것과 의도적으로 cast해서 비트를 버리는 것은 구분해야 합니다. 단순히 오류를 없애려고 작은 타입으로 변환하면 값 자체가 바뀔 수 있습니다. 변환은 다음 장에서 다룹니다.

## 실수와 bool

`f32`·`f64`는 부동소수점 수입니다. 정수와 달리 소수부를 표현할 수 있지만 모든 십진수를 정확히 저장하지는 못합니다. 금액을 작은 정수 단위로 관리하는 이유 중 하나입니다.

bool은 참과 거짓을 나타냅니다. 다음처럼 비교 결과를 저장하면 조건에 이름을 붙일 수 있습니다.

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

실행 결과:

```text
purchase allowed
```

`can_buy`가 balance와 cost의 변화를 자동으로 따라가는 것도 아닙니다. 두 값이 바뀐 뒤 현재 상태가 필요하면 비교를 다시 수행합니다.

## 초기화하지 않은 저장 공간

`var value: i32;`는 저장 공간만 선언하는 형태입니다. 읽기 전에 유효한 값을 써야 합니다. 선언했다고 자동으로 0이 들어간다고 가정하지 마십시오. 입문 과정에서는 값을 바로 알 수 있으면 선언과 동시에 초기화하는 편이 이해하기 쉽습니다.

출력 인자로 값을 받는 라이브러리 호출에서는 먼저 공간을 선언한 뒤 성공했을 때 읽는 경우도 있습니다. 그때는 함수의 성공 결과를 확인해야 합니다. 실패한 호출 뒤 초기화되지 않은 출력값을 읽는 실수를 피하십시오.

## 블록과 이름의 유효 범위

블록은 중괄호로 묶은 코드 영역입니다. 안쪽에서 같은 이름의 변수를 새로 선언하면 해당 블록 안에서는 새 변수가 사용됩니다. 이를 shadowing이라고 합니다.

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

실행 결과:

```text
inner=15
inner changed=16
outer=10
```

안쪽 선언의 초기식 `value + 5`는 바깥 value를 읽습니다. 새 변수의 초기화가 끝난 뒤 안쪽 value는 15입니다. 안쪽 값을 16으로 바꾸어도 바깥 저장 공간은 바뀌지 않습니다. 블록이 끝나면 다시 바깥 value가 보입니다.

반대로 안쪽 블록에서 `var` 없이 `value += 1`만 실행하면 보이는 기존 변수를 변경합니다. 새 선언인지 기존 값의 변경인지 키워드를 보고 구분하십시오.

## 지역 변수와 최상위 저장소

함수 밖에서는 const와 static을 사용할 수 있습니다. const는 상수 값을 표현하고 static은 실행 동안 유지되는 저장 공간입니다. 지역 변수와 같은 방식으로 모든 곳에서 선언할 수 있는 것은 아닙니다.

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

실행 결과:

```text
visit=1
visit=2
limit=3
```

visit를 두 번 호출해도 static은 호출마다 새로 0이 되지 않습니다. 반면 함수 안에서 `var visits: i32 = 0;`으로 선언하면 호출할 때마다 지역 저장 공간을 초기화합니다. 공유되는 변경 가능한 상태는 동작을 추적하기 어렵게 만들 수 있으므로 함수의 입력·출력으로 해결할 수 있는지 먼저 생각합니다.

## 흔한 실수

- 선언과 대입을 혼동해 같은 이름을 불필요하게 다시 선언하는 경우.
- 계산 결과를 저장한 변수가 입력 변수의 변경을 자동으로 따라간다고 생각하는 경우.
- 타입이 같으니 개수와 바이트 수 같은 단위도 같다고 생각하는 경우.
- 초기화 없이 읽거나 함수가 실패했는데 출력 인자를 읽는 경우.
- 지역 변수의 이름이 블록 밖에서도 보인다고 생각하는 경우.

오류가 난 이름을 찾을 때는 그 이름의 선언 위치와 중괄호 범위를 함께 확인합니다.

## 연습: 재고 변경 계산

초기 재고는 20개이고 3개씩 두 번 판매합니다. 남은 재고와 판매한 전체 개수를 출력하십시오. 재고를 바꿀 때마다 같은 변수를 갱신하고 판매량도 별도로 누적합니다.

### 전체 풀이

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

실행 결과:

```text
stock=14 sold=6
```

재고와 판매량이 함께 바뀌어야 합니다. 둘 중 하나만 갱신하면 값 사이의 관계가 깨집니다. 반복되는 주문 처리는 함수와 반복문을 배우면서 묶어 보겠습니다.


## 정수와 부동소수점 타입

정수 타입은 다음과 같습니다.

- 부호 있음: `i8`, `i16`, `i32`, `i64`, `i128`, `i256`, `i512`, `i1024`
- 부호 없음: `u8`, `u16`, `u32`, `u64`, `u128`, `u256`, `u512`, `u1024`
- 주소 크기 정수: `isz`, `usz`
- 부동소수점: `f32`, `f64`

`isz`는 주소 크기에 맞는 부호 있는 정수 타입이고, `usz`는 주소 크기에 맞는 부호 없는 정수 타입입니다.

## 기타 내장 타입

| 타입 | 용도 |
| --- | --- |
| `bool` | `true` 또는 `false` |
| `char` | 부호 없는 8비트 문자 값. 임의의 Unicode 코드 포인트 타입이 아님 |
| `byte` | 8비트 바이트 값 |
| `str` | NUL로 끝나는 문자열 바이트열 |
| `ptr<T>` | `T`를 대상으로 하는 포인터 |
| `array<T, N>` | 요소 타입 `T`, 길이 `N`인 고정 길이 배열 |

사용자 정의 구조체와 열거형, 타입 별칭도 타입 위치에 사용할 수 있습니다.

`var`는 지역 변수를 선언하는 문법입니다. 타입 별칭은 같은 타입을 코드의 문맥에 맞는 이름으로 표현하는 문법입니다.
