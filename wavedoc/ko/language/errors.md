---
translation_set_id: learn-errors
path: language/errors
locale: ko
group: language
group_order: 2
order: 12
title: 12. 오류를 표현하고 복구하기
summary: 오류와 정상 값을 구분하고 실패 경로에서 자원을 정리합니다.
---

## 실패도 함수의 결과

파일이 없거나 입력이 범위를 벗어나는 일은 프로그램에서 자연스럽게 발생합니다. 오류를 처리한다는 것은 메시지 하나를 출력하는 데 그치지 않습니다. 실패를 구분하고, 이미 진행한 작업의 상태를 확인하고, 확보한 자원을 정리한 뒤 계속할지 종료할지 선택하는 과정입니다.

이번 장에서는 작은 함수의 실패 표현부터 시작해 결과 구조체, variant, 조기 반환과 자원 정리로 발전시킵니다.

## 실패 표지가 성공값과 겹치면 안 됨

배열 검색에서 -1을 못 찾음으로 사용한 것은 유효 인덱스가 0 이상이기 때문입니다. 반면 어떤 정수든 정상 결과일 수 있는 계산에서는 -1을 오류로 정하면 정상값 -1과 구별할 수 없습니다.

0도 흔히 오해하는 값입니다. 빈 문자열 길이 0, 첫 위치 0, 전송한 바이트 수 0은 함수마다 의미가 다릅니다. 반환값이 0이 아니라는 이유만으로 성공이라고 판단하지 마십시오.

## 성공 여부와 값을 함께 반환하기

<!-- wave-example: book-error-result -->
```wave
struct Division {
    ok: bool;
    value: i32;
}

fun divide_nonnegative(left: i32, right: i32) -> Division {
    if (left < 0 || right <= 0) {
        return Division {
            ok: false,
            value: 0
        };
    }

    return Division {
        ok: true,
        value: left / right
    };
}

fun main() {
    var result: Division = divide_nonnegative(0, 3);

    if (!result.ok) {
        println("invalid input");
        return;
    }

    println("value={}", result.value);
}
```

실행 결과:

```text
value=0
```

정상 결과도 0일 수 있습니다. value를 보고 성공 여부를 추측하지 않고 ok를 먼저 확인합니다. 실패한 결과에 value 필드가 존재하더라도 사용할 값이라는 뜻은 아닙니다.

이 예제의 입력 규칙은 왼쪽이 0 이상, 오른쪽이 양수라는 것입니다. 함수 이름과 설명에 범위를 드러내 일반적인 signed 정수 나눗셈과 구분했습니다.

## 오류 원인도 구분하기

실패 이유에 따라 다른 안내나 복구를 하려면 오류 정보를 추가합니다. 입력 범위를 검사하는 함수와 계산하는 함수를 나누는 방법도 있습니다.

<!-- wave-example: book-error-codes -->
```wave
struct Check {
    ok: bool;
    error: i32;
}

fun check_quantity(value: i32) -> Check {
    if (value < 1) {
        return Check {
            ok: false,
            error: 1
        };
    }

    if (value > 100) {
        return Check {
            ok: false,
            error: 2
        };
    }

    return Check {
        ok: true,
        error: 0
    };
}

fun main() {
    var result: Check = check_quantity(101);

    if (!result.ok) {
        match (result.error) {
            1 => {
                println("quantity is too small");
            }
            2 => {
                println("quantity is too large");
            }
            _ => {
                println("invalid quantity");
            }
        }
    }
}
```

실행 결과:

```text
quantity is too large
```

숫자의 의미는 이 함수가 정의합니다. 다른 라이브러리의 오류 번호 1이나 2와 같다고 볼 수 없습니다. 공개 API에서는 오류 상수나 타입에 이름을 붙이면 호출자가 임의의 숫자를 외우지 않아도 됩니다.

## variant로 결과 분리하기

성공값과 오류가 동시에 존재할 수 없다는 관계를 variant로 표현할 수 있습니다.

<!-- wave-example: book-error-variant -->
```wave
variant ByteResult {
    Value(u8),
    OutOfRange(i32)
}

fun to_byte(value: i32) -> ByteResult {
    if (value < 0 || value > 255) {
        return ByteResult::OutOfRange(value);
    }

    return ByteResult::Value(value as u8);
}

fun main() {
    var result: ByteResult = to_byte(300);

    match (result) {
        ByteResult::Value(value) => {
            println("byte={}", value);
        }
        ByteResult::OutOfRange(original) => {
            println("out of range: {}", original);
        }
    }
}
```

실행 결과:

```text
out of range: 300
```

오류에 원래 입력값을 담았습니다. 단순히 false만 돌려주는 것보다 호출자가 문제를 설명하기 쉽습니다. 비밀번호나 토큰 같은 민감한 입력은 그대로 로그에 남기지 않는 식으로 데이터의 성격도 고려합니다.

## 조기 반환으로 정상 경로 읽기 쉽게 만들기

검사를 여러 단계 수행할 때 모든 정상 코드를 깊은 if 안에 넣을 필요는 없습니다. 실패하면 바로 반환하고, 아래에는 성공 경로를 이어 쓸 수 있습니다.

<!-- wave-example: book-error-early-return -->
```wave
fun show_total(quantity: i32, price: i32) {
    if (quantity < 1 || quantity > 100) {
        println("invalid quantity");
        return;
    }

    if (price < 0 || price > 100000) {
        println("invalid price");
        return;
    }

    var total: i32 = quantity * price;
    println("total={}", total);
}

fun main() {
    show_total(3, 1200);
    show_total(0, 1200);
    show_total(3, -1);
}
```

실행 결과:

```text
total=3600
invalid quantity
invalid price
```

검사를 통과한 뒤에는 quantity와 price가 정해진 범위 안이라는 사실을 이용할 수 있습니다. 범위는 중간 곱셈도 i32에 들어가도록 정했습니다. 조기 반환을 추가할 때는 그 지점에서 이미 소유한 자원이 있는지 함께 확인해야 합니다.

## 실패 경로의 메모리 정리

<!-- wave-example: book-error-cleanup -->
```wave
import("std::buffer::alloc")::{
    Buffer, buffer_init, buffer_free
};
import("std::buffer::write")::{
    buffer_append_str
};

fun make_message(allowed: bool) -> i32 {
    var data: Buffer;

    if (buffer_init(&data, 16) < 0) {
        return 1;
    }

    if (!allowed) {
        buffer_free(&data);
        return 2;
    }

    if (buffer_append_str(&data, "ready") < 0) {
        buffer_free(&data);
        return 3;
    }

    println("bytes={}", data.len);

    if (buffer_free(&data) < 0) {
        return 4;
    }

    return 0;
}

fun main() {
    println("status={}", make_message(false));
}
```

실행 결과:

```text
status=2
```

데이터를 출력하지 않는 실패 경로에서도 Buffer를 해제합니다. 모든 실패를 하나의 return으로 바꾸는 것보다 각 지점에서 무엇을 소유하는지 명확히 관리하는 것이 중요합니다.

정리 자체가 실패하는 API도 있습니다. 원래 작업의 오류와 정리 오류를 어떻게 보존할지 정하십시오. 이 작은 예제는 작업 오류 번호를 우선 반환합니다. 더 큰 프로그램은 둘을 각각 기록할 수 있습니다.

## 부분 성공은 자동으로 취소되지 않음

파일에 일부 바이트를 쓴 뒤 쓰기가 실패하면 이미 쓴 바이트가 사라지지 않습니다. 네트워크에서도 상대편이 일부 데이터를 받았을 수 있습니다. 같은 작업을 처음부터 반복하면 중복 기록이 생길 수 있습니다.

반대로 checked byte cursor 읽기는 실패했을 때 위치와 출력값을 보존합니다. 이런 함수는 입력을 더 받은 뒤 같은 위치에서 다시 시도할 수 있습니다. “실패하면 아무것도 안 바뀐다”는 규칙을 모든 API에 적용하지 말고 해당 함수의 문서를 확인합니다.

## 복구 가능한 오류와 trap

잘못된 파일 경로나 사용자 입력은 오류값으로 전달해 복구할 수 있게 설계할 수 있습니다. 잘못된 런타임 시프트 횟수나 유효하지 않은 실수→정수 변환의 trap은 같은 인터페이스가 아닙니다.

계속 실행해야 하는 프로그램은 위험한 연산 전에 입력을 검사해야 합니다. assert도 사용자 입력의 정상적인 실패를 처리하는 수단으로 남용하지 않습니다. 사용자에게 다시 입력할 기회를 주어야 한다면 반환 결과를 통해 제어 흐름을 이어갑니다.

## 연습과 전체 풀이

배열에서 인덱스로 값을 조회하되 음수나 상한 이상이면 실패하도록 함수를 작성하십시오. 성공값이 0일 수도 있으므로 결과 구조체를 사용합니다.

<!-- wave-example: book-error-solution -->
```wave
struct Lookup {
    ok: bool;
    value: i32;
}

fun at(values: ptr<i32>, length: i32, index: i32) -> Lookup {
    if (index < 0 || index >= length) {
        return Lookup {
            ok: false,
            value: 0
        };
    }

    return Lookup {
        ok: true,
        value: deref values[index]
    };
}

fun main() {
    var values: array<i32, 3> = [0, 10, 20];
    var first: Lookup = at(&values[0], 3, 0);
    var outside: Lookup = at(&values[0], 3, 3);

    if (first.ok) {
        println("value={}", first.value);
    }

    if (!outside.ok) {
        println("out of bounds");
    }
}
```

실행 결과:

```text
value=0
out of bounds
```

포인터와 length가 실제 읽을 수 있는 배열을 나타내는 것은 호출자의 조건입니다. 인덱스 검사만으로 임의의 주소까지 안전하게 만드는 함수는 아닙니다. 함수가 책임지는 검사와 호출자가 보장할 조건을 나누어 읽으십시오.
