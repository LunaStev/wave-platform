---
translation_set_id: data-types
path: language/structures-enums-and-aliases
locale: ko
group: language
group_order: 2
order: 8
title: 8. 구조체, enum과 variant
summary: 필드, 구조체 초기화, enum과 variant의 역할을 배웁니다.
---

## 데이터 사이의 관계를 타입으로 표현하기

상품의 가격과 개수를 각각 변수로 전달하면 두 값이 같은 상품에 속하는지 코드만 보고 알기 어렵습니다. 구조체는 관련된 필드를 묶습니다. enum은 이름 붙은 상태를 나타내고, variant는 경우마다 다른 데이터를 함께 저장합니다.

세 기능은 서로 대체하는 문법이 아닙니다. 무엇을 표현하려는지에 따라 선택합니다.

| 표현할 것 | 선택 | 예 |
| --- | --- | --- |
| 동시에 존재하는 여러 필드 | struct | 상품의 단가와 수량 |
| 이름이 있는 정수 상태 | enum | 대기·진행·완료 |
| 경우마다 다른 데이터 | variant | 성공값 또는 오류 |

## 구조체 선언과 값 생성

<!-- wave-example: book-struct-first -->
```wave
struct Product {
    price: i32;
    quantity: i32;
}

fun main() {
    var item: Product = Product {
        price: 1500,
        quantity: 2
    };

    println("{} {}", item.price, item.quantity);
}
```

실행 결과:

```text
1500 2
```

선언의 필드는 세미콜론으로 끝나고, 값을 만들 때 필드와 값은 콜론으로 연결하고 쉼표로 구분합니다. 타입을 정의한 것과 실제 값을 만든 것은 다른 단계입니다. Product라는 타입을 선언했다고 상품 하나의 저장 공간이 자동으로 생기지는 않습니다.

필드는 `item.price`처럼 접근합니다. 같은 타입의 item을 여러 개 만들면 각각 다른 값을 저장할 수 있습니다.

## 구조체를 함수에 전달하기

<!-- wave-example: book-struct-function -->
```wave
struct Product {
    price: i32;
    quantity: i32;
}

fun subtotal(item: Product) -> i32 {
    return item.price * item.quantity;
}

fun main() {
    var item: Product = Product {
        price: 1500,
        quantity: 2
    };

    println("{}", subtotal(item));

    item.quantity = 3;
    println("{}", subtotal(item));
}
```

실행 결과:

```text
3000
4500
```

함수는 단가와 수량이 같은 상품에 속한다는 관계를 타입으로 받습니다. 값으로 전달한 정수 필드 구조체를 읽는 함수입니다. 호출자의 저장 공간을 수정하려는 함수라면 포인터를 받도록 설계할 수 있습니다.

구조체에 포인터 필드가 있으면 값 복사는 주소도 복사합니다. 별도의 할당까지 깊게 복사하는 기능은 아닙니다. 파일 핸들이나 Buffer 같은 자원을 담은 타입은 복사와 해제 규칙을 함께 정해야 합니다.

## 메서드와 proto

관련 함수를 메서드 형태로 묶을 수 있습니다. proto는 구조체의 메서드를 별도 블록으로 작성하는 방법입니다.

<!-- wave-example: book-struct-method -->
```wave
struct Counter {
    value: i32;
}

proto Counter {
    fun current(self: Counter) -> i32 {
        return self.value;
    }
}

fun main() {
    var counter: Counter = Counter {
        value: 3
    };

    println("{}", counter.current());
}
```

실행 결과:

```text
3
```

`self: Counter`는 값을 받는 매개변수입니다. 메서드 호출 표기를 쓴다고 자동으로 원본을 수정하는 메서드가 되는 것은 아닙니다. self의 타입과 본문에서 수행하는 작업을 함께 읽으십시오.

메서드를 붙였다는 이유로 필드가 유효한 상태만 가질 것이라고 기대할 수도 없습니다. 공개 필드로 사용자가 만들 수 있는 잘못된 조합이 있다면 함수에서 검사하거나 생성 규칙을 제공해야 합니다.

## enum으로 상태에 이름 붙이기

<!-- wave-example: book-enum-state -->
```wave
enum State -> i32 {
    Ready = 0,
    Running,
    Finished,
}

fun main() {
    var state: State = State::Ready;

    if (state == State::Ready) {
        println("ready");
    }

    state = State::Running;

    if (state == State::Running) {
        println("running");
    }
}
```

실행 결과:

```text
ready
running
```

`-> i32`는 표현에 사용하는 정수 타입입니다. 첫 값을 0으로 정했고 뒤의 생략된 값은 앞 값보다 1 큽니다. 코드에서 단순히 0과 1을 비교하는 것보다 State::Ready와 State::Running을 쓰면 뜻이 드러납니다.

enum 이름이 있다고 상태 전이가 자동으로 제한되는 것은 아닙니다. Finished에서 Running으로 돌아갈 수 있는지 같은 규칙은 함수로 구현해야 합니다.

## variant로 경우와 데이터 연결하기

성공했을 때만 값이 있고 실패했을 때는 오류 정보가 필요하다면 variant로 나타낼 수 있습니다.

<!-- wave-example: book-variant-result -->
```wave
variant Result {
    Value(i32),
    Error(i32)
}

fun divide(left: i32, right: i32) -> Result {
    if (left < 0 || right <= 0) {
        return Result::Error(1);
    }

    return Result::Value(left / right);
}

fun show(result: Result) {
    match (result) {
        Result::Value(value) => {
            println("value={}", value);
        }
        Result::Error(code) => {
            println("error={}", code);
        }
    }
}

fun main() {
    show(divide(12, 3));
    show(divide(12, 0));
}
```

실행 결과:

```text
value=4
error=1
```

Result::Value와 Result::Error는 각각 payload를 담는 경우입니다. 같은 정수 타입을 담더라도 어느 경우인지 구분됩니다. 호출자는 match로 경우를 확인하고 해당 arm 안에서 payload를 사용합니다.

위 divide는 음수 피연산자를 지원하지 않는 작은 예제입니다. 입력 범위를 명시했기 때문에 일반적인 signed 나눗셈의 모든 경계를 다루는 함수와 혼동하지 않아야 합니다.

## payload가 없는 경우

모든 경우에 데이터가 필요한 것은 아닙니다. 값이 없는 상태를 별도 경우로 표현할 수 있습니다.

<!-- wave-example: book-variant-empty -->
```wave
variant Lookup {
    Found(i32),
    Missing
}

fun main() {
    var result: Lookup = Lookup::Missing;

    match (result) {
        Lookup::Found(value) => {
            println("found={}", value);
        }
        Lookup::Missing => {
            println("missing");
        }
    }
}
```

실행 결과:

```text
missing
```

임의의 정수 하나를 “없음”으로 예약하는 대신 Missing이라는 경우를 사용했습니다. 성공값이 어떤 정수여도 의미가 겹치지 않습니다.

match에서 `_`는 나머지 경우를 처리합니다. 새 경우를 추가했을 때 각 호출부가 다시 검토되기를 원한다면 모든 경우를 명시적으로 나누는 편이 좋습니다. 어떤 방법을 선택하든 처리하지 않은 입력이 없도록 확인합니다.

## 구조체와 variant를 함께 쓰기

서로 다른 데이터를 선택하는 것은 variant, 한 경우에 속하는 여러 필드를 묶는 것은 구조체로 표현할 수 있습니다. 예를 들어 주문 처리 결과가 성공이면 영수증 구조체, 실패이면 오류 번호를 담도록 설계할 수 있습니다.

값 안에 다른 값이 포함되어도 메모리 수명 규칙은 사라지지 않습니다. variant에 포인터를 담았다면 그 포인터가 유효한지, 누가 해제할지는 별도로 정합니다. 외부 파일 형식으로 저장할 때도 구조체 메모리를 그대로 덤프하지 말고 필드별 인코딩을 정해야 합니다.

## 연습과 전체 풀이

0~100 사이의 점수만 허용하는 판정 결과를 만드십시오. 유효한 점수는 Grade(score), 나머지는 Invalid로 표현합니다.

<!-- wave-example: book-model-solution -->
```wave
variant CheckedScore {
    Grade(i32),
    Invalid
}

fun check_score(score: i32) -> CheckedScore {
    if (score < 0 || score > 100) {
        return CheckedScore::Invalid;
    }

    return CheckedScore::Grade(score);
}

fun main() {
    var result: CheckedScore = check_score(87);

    match (result) {
        CheckedScore::Grade(score) => {
            println("accepted={}", score);
        }
        CheckedScore::Invalid => {
            println("invalid score");
        }
    }
}
```

실행 결과:

```text
accepted=87
```

입력을 -1, 0, 100, 101로 바꿔 경계를 확인하십시오. 성공과 실패가 같은 정수 공간을 공유하지 않으므로 호출부에서 오류값을 실수로 평균 계산에 더하는 위험을 줄일 수 있습니다.


## 어떤 표현을 선택할까

| 데이터의 모양 | 적합한 표현 | 예 |
| --- | --- | --- |
| 같은 타입의 값 여러 개 | 배열 | 점수 10개 |
| 서로 관련된 여러 필드 | 구조체 | 이름과 점수 |
| 이름을 붙인 상태 값 | enum | Ready, Running, Stopped |
| 상태마다 다른 추가 데이터 | variant | Value(i32), Error(str) |
| 기존 타입에 문맥상의 이름 | 타입 별칭 | UserId = u64 |

데이터 구조를 선택할 때는 저장할 값뿐 아니라, 어떤 잘못된 상태를 표현할 수 있는지도 생각합니다. 성공 여부와 값·오류 필드를 모두 가진 구조체는 잘못된 조합을 만들 수 있지만 variant는 경우별 payload로 구분할 수 있습니다.

## 메모리 배치와 외부 데이터

구조체의 메모리에는 필드 사이 정렬을 맞추기 위한 빈 공간이 들어갈 수 있습니다. 필드 크기를 단순히 더한 값이 구조체 전체 크기와 항상 같지는 않습니다. 크기와 정렬을 알아야 할 때는 [mem 레이아웃 함수](/docs/ko/reference/memory-and-buffer)를 사용합니다.

파일이나 네트워크 메시지는 [bytes](/docs/ko/stdlib/bytes) 함수로 필드를 순서대로 인코딩하면 바이트 순서와 길이를 명확하게 정할 수 있습니다. 다른 언어와 구조체를 전달할 때는 [FFI](/docs/ko/language/modules-imports-and-ffi)의 외부 선언과 대상 ABI를 맞춥니다.

[구조체 학습과 연습](/docs/ko/language/structures-enums-and-aliases) · [variant](/docs/ko/language/variants)
