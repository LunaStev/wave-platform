---
translation_set_id: control-flow
path: language/control-flow
locale: ko
group: language
group_order: 2
order: 4
title: 4. 조건, 반복과 경계값
summary: if, for, while과 반복문의 범위를 배웁니다.
---

## 실행 경로를 선택하기

앞 장의 프로그램은 문장을 위에서 아래로 실행했습니다. 실제 프로그램은 입력과 상태에 따라 다른 작업을 해야 합니다. 조건문은 실행 경로를 선택하고 반복문은 같은 규칙을 여러 값에 적용합니다.

예제는 각각 main.wave에 저장해 실행합니다. 코드를 읽을 때는 현재 변수 값, 다음에 검사할 조건, 실행할 문장 순서를 종이에 적어 보십시오. 결과를 외우는 것보다 흐름을 따라가는 연습이 중요합니다.

## if와 else

조건은 괄호 안에 쓰고 본문은 중괄호로 묶습니다. 다음 예제에서 balance를 500 또는 2000으로 바꿔 어느 분기가 실행되는지 확인하십시오.

<!-- wave-example: book-if-else -->
```wave
fun main() {
    var balance: i32 = 2000;
    var price: i32 = 1200;

    if (balance >= price) {
        balance -= price;
        println("bought, balance={}", balance);
    } else {
        println("not enough money");
    }
}
```

실행 결과:

```text
bought, balance=800
```

두 블록을 모두 실행하는 것이 아닙니다. 조건이 참이면 첫 블록, 거짓이면 else 블록을 실행합니다. balance가 price와 같을 때도 구매를 허용하므로 `>=`를 썼습니다. `>`로 바꾸면 같은 금액에서 동작이 달라집니다.

## 여러 조건의 순서

else if로 조건을 이어 붙일 수 있습니다. 위에서 먼저 만족한 분기 하나를 실행하므로 큰 경계부터 검사할지 작은 경계부터 검사할지 생각해야 합니다.

<!-- wave-example: book-grade -->
```wave
fun grade(score: i32) {
    if (score < 0 || score > 100) {
        println("invalid");
    } else if (score >= 90) {
        println("A");
    } else if (score >= 80) {
        println("B");
    } else {
        println("C");
    }
}

fun main() {
    grade(95);
    grade(80);
    grade(79);
    grade(101);
}
```

실행 결과:

```text
A
B
C
invalid
```

유효하지 않은 점수를 먼저 거부한 뒤 등급을 나눕니다. `score >= 80`을 맨 앞에 두면 95도 그 분기에 들어가므로 A 분기에 도달하지 못합니다. 조건 각각의 정확성뿐 아니라 조건들의 순서도 검사하십시오.

## 조건식에서 값을 바꾸지 않기

Wave의 if·while·for 조건식에는 대입, 복합 대입, 증감 연산을 넣지 않습니다. 비교하려면 `==`를 사용합니다. 읽으면서 동시에 갱신하고 싶다면 두 문장으로 나눕니다.

함수 내부 조각의 올바른 형태:

```wave
value = read_value();

if (value == expected) {
    println("matched");
}
```

read_value와 expected는 이 조각에서 정의하지 않았으므로 그대로 실행하는 전체 프로그램은 아닙니다. 여기서 보여 주는 규칙은 “상태 변경 후 비교”입니다.

## while: 조건이 유지되는 동안

<!-- wave-example: book-while-countdown -->
```wave
fun main() {
    var remaining: i32 = 3;

    while (remaining > 0) {
        println("{}", remaining);
        remaining -= 1;
    }

    println("finished at {}", remaining);
}
```

실행 결과:

```text
3
2
1
finished at 0
```

조건은 본문에 들어가기 전에 검사합니다. 처음부터 remaining이 0이면 본문을 한 번도 실행하지 않습니다. 본문 마지막의 감소가 빠지면 조건이 계속 참이라 반복이 끝나지 않습니다.

반복문을 작성한 뒤 “무엇이 종료 조건에 가까워지게 하는가?”를 확인하십시오. 입력을 기다리는 루프라면 입력 변화나 EOF가, 숫자 루프라면 인덱스 갱신이 그 역할을 합니다.

## for: 초기화·조건·갱신

for는 반복에 필요한 세 부분을 모아 표현합니다.

<!-- wave-example: book-for-sum -->
```wave
fun main() {
    var total: i32 = 0;

    for (var number: i32 = 1; number <= 5; number += 1) {
        total += number;
    }

    println("sum={}", total);
}
```

실행 결과:

```text
sum=15
```

1. number를 1로 초기화합니다. 이 단계는 한 번입니다.
2. number <= 5를 검사합니다. 거짓이면 반복을 끝냅니다.
3. 본문에서 total에 number를 더합니다.
4. number를 1 증가시키고 조건 검사로 돌아갑니다.

for 안에서 선언한 반복 변수를 반복 뒤에도 사용할 수 있다고 가정하지 않습니다. 반복 뒤 값이 필요한 설계라면 바깥에서 선언하고 초기화 위치를 분명히 하십시오.

## 포함 경계와 제외 경계

1부터 5까지의 합에는 `<= 5`가 자연스럽습니다. 반면 길이가 5인 배열의 인덱스에는 `< 5`를 사용해야 합니다. 배열 인덱스는 0부터 시작해 4에서 끝나기 때문입니다.

“다섯 번 실행”과 “5라는 값까지 포함”을 혼동하지 마십시오. 시작값과 종료 비교를 함께 봐야 반복 횟수를 알 수 있습니다. 입력이 비어 있는 경우와 요소가 하나뿐인 경우는 경계 실수를 발견하기 좋습니다.

## continue와 break

continue는 이번 반복의 나머지를 건너뛰고, break는 가장 가까운 반복문을 끝냅니다.

<!-- wave-example: book-loop-control -->
```wave
fun main() {
    var total: i32 = 0;

    for (var number: i32 = 1; number <= 10; number += 1) {
        if (number == 3) {
            continue;
        }

        if (number == 6) {
            break;
        }

        total += number;
    }

    println("{}", total);
}
```

실행 결과:

```text
12
```

합계에 들어간 수는 1, 2, 4, 5입니다. for에서 continue를 만나면 갱신식으로 진행합니다. while에는 for와 같은 별도 갱신식이 없으므로 continue 이전에 필요한 상태 변경을 빠뜨리지 않도록 주의해야 합니다.

반복문이 중첩되어 있으면 break 하나가 모든 반복을 끝내는 것이 아닙니다. 여러 단계에서 중단해야 한다면 함수로 작업을 묶고 return을 사용하거나 바깥 반복에서도 종료 조건을 확인하는 식으로 의도를 드러냅니다.

## match로 경우 나누기

같은 값의 여러 경우를 비교할 때 match를 사용할 수 있습니다. 각 arm의 본문은 블록입니다.

<!-- wave-example: book-match-number -->
```wave
fun describe(status: i32) {
    match (status) {
        200 => {
            println("ok");
        }
        404 => {
            println("missing");
        }
        _ => {
            println("other");
        }
    }
}

fun main() {
    describe(200);
    describe(404);
    describe(500);
}
```

실행 결과:

```text
ok
missing
other
```

`_`는 나머지를 처리하는 패턴입니다. 같은 match 안에 중복으로 두지 않습니다. 값에 따라 데이터 종류도 달라지는 variant는 [데이터 모델 장](/docs/ko/language/structures-enums-and-aliases)에서 다룹니다.

## 완성 예제: 조건에 맞는 수 집계

1부터 10까지 중 짝수의 개수와 합을 구합니다. 개수와 합은 다른 정보이므로 각각 변수로 누적합니다.

<!-- wave-example: book-loop-statistics -->
```wave
fun main() {
    var count: i32 = 0;
    var total: i32 = 0;

    for (var number: i32 = 1; number <= 10; number += 1) {
        if (number % 2 == 0) {
            count += 1;
            total += number;
        }
    }

    println("count={} total={}", count, total);
}
```

실행 결과:

```text
count=5 total=30
```

짝수는 2, 4, 6, 8, 10이므로 개수는 5, 합은 30입니다. 식이 짧아도 결과를 손으로 구할 수 있는 작은 범위에서 먼저 확인하면 반복 경계를 검증하기 쉽습니다.

## 연습과 전체 풀이

1부터 20까지 중 3의 배수만 더하되, 합계가 30을 넘게 되는 값은 더하지 않고 종료하십시오. “더한 뒤 넘었는지 확인”과 “넘을지 확인한 뒤 더하기”를 구분해야 합니다.

<!-- wave-example: book-loop-solution -->
```wave
fun main() {
    var total: i32 = 0;

    for (var number: i32 = 1; number <= 20; number += 1) {
        if (number % 3 != 0) {
            continue;
        }

        if (total + number > 30) {
            break;
        }

        total += number;
    }

    println("{}", total);
}
```

실행 결과:

```text
30
```

3+6+9+12는 30이고 다음 15는 더하지 않습니다. 예제 범위에서는 덧셈이 안전하지만 일반적인 큰 정수 입력에서는 검사식 `total + number` 자체의 오버플로도 고려해야 합니다. 검사 코드를 썼다는 사실만으로 모든 경계가 해결되지는 않습니다.
