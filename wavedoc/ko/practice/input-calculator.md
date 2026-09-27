---
translation_set_id: practice-input-calculator
path: practice/input-calculator
locale: ko
group: practice
group_order: 4
order: 1
title: 실습: 입력을 검사하는 계산기
summary: 입력, 범위 검사, 함수와 종료 코드를 연결합니다.
---

## 목표와 실행

수량과 단가를 입력받아 합계를 계산합니다. 정수 두 개를 공백이나 줄바꿈으로 구분해 입력합니다. 이 예제는 수량 1~1000, 단가 0~100000만 받으므로 i32 곱셈 범위 안에서 계산합니다.

`main.wave`에 저장하고 `wavec run main.wave`를 실행한 뒤 `3 1200`을 입력하십시오. 프로그램이 출력하는 결과는 아래와 같습니다. 터미널에서 입력한 문자가 보이는 것은 프로그램 출력과 별개입니다.

<!-- wave-example: input-calculator -->
```wave
fun total(quantity: i32, price: i32) -> i32 {
    return quantity * price;
}

fun main() -> i32 {
    var quantity: i32 = 0;
    var price: i32 = 0;
    input("{} {}", quantity, price);
    if (quantity < 1 || quantity > 1000 || price < 0 || price > 100000) {
        println("out of range");
        return 1;
    }

    println("total={}", total(quantity, price));
    return 0;
}
```

실행 결과:

```text
total=3600
```

## 실패도 확인하기

`0 1200`을 넣으면 `out of range`와 종료 코드 1을 기대합니다. `input`의 숫자 파싱 실패와 프로그램의 업무 범위 검사는 다른 단계입니다. 숫자가 아닌 토큰·타입 범위 초과·필요한 입력 전 EOF는 입력 실패입니다. 이 내장 입력은 오류를 반환해 재입력시키는 인터페이스가 아닙니다. 복구 가능한 파서가 필요하면 [io](/docs/ko/stdlib/files-io)로 바이트를 읽어 검증 과정을 직접 구성합니다.

## 확장 연습과 해설

할인율을 세 번째 입력으로 받아 0~100인지 검사하십시오. 큰 중간 곱셈을 피하려면 i64로 계산 범위를 넓히고 결과를 좁힐 때 범위를 확인해야 합니다. 단순히 마지막 결과 타입만 넓히면 중간 계산을 이미 좁은 타입에서 수행했을 수 있습니다.

[콘솔 I/O](/docs/ko/language/console-io-and-formatting) · [다음: 파일 읽기](/docs/ko/practice/file-reader)

## 계산 범위를 먼저 정하는 이유

가장 큰 입력은 수량 1000과 단가 100000입니다. 두 값을 곱하면 100000000이므로 i32 범위 안에 들어갑니다. 이 범위 검사가 있어야 total 함수의 곱셈 결과를 그대로 사용할 수 있습니다.

입력을 받는 main은 입력과 오류 메시지를 담당하고, total은 계산만 담당합니다. 나중에 파일에서 주문을 읽도록 바꿔도 계산 함수는 그대로 사용할 수 있습니다.

## 할인 계산 완성하기

할인율을 추가하면 합계에 100까지 곱할 수 있습니다. 중간 계산부터 i64로 수행하도록 변환한 뒤 할인율을 적용합니다. 정수 나눗셈의 소수 부분은 버리므로, 이 예제는 할인 후 금액을 정수 단위로 절삭합니다.

<!-- wave-example: book-calculator-discount -->
```wave
fun discounted_total(quantity: i32, price: i32, discount: i32) -> i64 {
    var subtotal: i64 = (quantity as i64) * (price as i64);
    var remaining: i64 = 100 - (discount as i64);
    return subtotal * remaining / 100;
}

fun main() -> i32 {
    var quantity: i32 = 0;
    var price: i32 = 0;
    var discount: i32 = 0;

    input("{} {} {}", quantity, price, discount);

    if (quantity < 1 || quantity > 1000) {
        println("invalid quantity");
        return 1;
    }

    if (price < 0 || price > 100000) {
        println("invalid price");
        return 2;
    }

    if (discount < 0 || discount > 100) {
        println("invalid discount");
        return 3;
    }

    println("total={}", discounted_total(quantity, price, discount));
    return 0;
}
```

실행 결과:

```text
total=3240
```

위 결과는 `3 1200 10`을 입력했을 때입니다. 원래 합계 3600에서 10%를 뺀 3240이 출력됩니다.

| 입력 | 예상 결과 | 확인하는 경로 |
| --- | --- | --- |
| `3 1200 0` | `total=3600` | 할인 없음 |
| `3 1200 100` | `total=0` | 전액 할인 |
| `3 1200 101` | `invalid discount` | 할인율 범위 초과 |
| `1000 100000 0` | `total=100000000` | 최대 입력 |
| `0 1200 10` | `invalid quantity` | 수량 범위 미달 |

## 다음 연습

함수를 바꿔 소수 금액을 반올림해 보십시오. 양수 금액만 받는 이 프로그램에서는 100으로 나누기 전에 50을 더하면 정수 단위 반올림이 됩니다. 1개에 99, 할인율 50을 입력했을 때 절삭 결과 49와 반올림 결과 50을 비교할 수 있습니다.
