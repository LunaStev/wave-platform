---
translation_set_id: learn-arrays-strings
path: language/arrays
locale: ko
group: language
group_order: 2
order: 6
title: 6. 배열과 데이터 순회
summary: 고정 길이 배열, 인덱스, UTF-8 바이트와 문자열 끝을 배웁니다.
---

## 같은 종류의 여러 값

점수 세 개를 score1, score2, score3으로 따로 만들면 개수가 바뀔 때 선언과 계산을 모두 바꿔야 합니다. 배열은 같은 타입의 요소를 정해진 수만큼 묶습니다. 반복문을 사용하면 요소마다 같은 규칙을 적용할 수 있습니다.

이번 장에서는 배열 생성, 인덱스, 변경, 순회, 검색과 집계를 배웁니다. 문자열도 인덱싱할 수 있지만 의미가 다르므로 다음 장에서 별도로 다룹니다.

## 타입에 길이를 적기

<!-- wave-example: book-array-create -->
```wave
fun main() {
    var scores: array<i32, 3> = [70, 80, 90];

    println("first={}", scores[0]);
    println("second={}", scores[1]);
    println("last={}", scores[2]);
}
```

실행 결과:

```text
first=70
second=80
last=90
```

`array<i32, 3>`의 i32는 요소 타입이고 3은 요소 수입니다. 저장하는 것은 정수 3개입니다. 바이트 수 3이라는 뜻이 아닙니다. 배열 리터럴의 요소 개수는 선언한 길이와 맞아야 합니다.

인덱스는 0부터 시작합니다. 첫 요소가 0, 마지막 요소가 길이-1입니다. scores[3]은 세 번째 요소가 아니라 유효 범위를 벗어난 접근입니다.

## 요소 변경

<!-- wave-example: book-array-update -->
```wave
fun main() {
    var scores: array<i32, 3> = [70, 80, 90];

    scores[1] = 85;
    scores[0] += 5;

    println("{} {} {}", scores[0], scores[1], scores[2]);
}
```

실행 결과:

```text
75 85 90
```

배열 전체를 다시 만들지 않고 특정 요소의 저장 공간을 변경합니다. 인덱스 식도 계산 결과일 수 있지만 그 값이 범위 안인지 보장해야 합니다. 외부 입력을 인덱스로 쓸 때는 음수 여부와 상한을 모두 검사합니다.

## 배열 전체를 순회하기

<!-- wave-example: book-array-sum -->
```wave
fun main() {
    var scores: array<i32, 4> = [60, 70, 80, 90];
    var total: i32 = 0;

    for (var index: i32 = 0; index < 4; index += 1) {
        total += scores[index];
    }

    println("total={} average={}", total, total / 4);
}
```

실행 결과:

```text
total=300 average=75
```

반복마다 다른 인덱스의 요소를 읽습니다. 합계 변수는 반복 밖에서 초기화해야 합니다. 반복 본문 안에서 매번 0으로 초기화하면 마지막 요소만 남는 식의 잘못된 결과가 생깁니다.

평균의 정수 나눗셈은 소수부를 버립니다. 소수 평균이 필요하면 합계를 실수로 바꾼 뒤 나누십시오. 더 큰 배열이나 값에서는 합계를 담을 타입이 충분히 넓은지도 확인해야 합니다.

## 일부 요소만 집계하기

조건문과 순회를 조합하면 필터링할 수 있습니다. 여기서는 80점 이상인 요소의 개수를 셉니다.

<!-- wave-example: book-array-filter -->
```wave
fun main() {
    var scores: array<i32, 5> = [60, 80, 90, 75, 100];
    var passed: i32 = 0;

    for (var index: i32 = 0; index < 5; index += 1) {
        if (scores[index] >= 80) {
            passed += 1;
        }
    }

    println("passed={}", passed);
}
```

실행 결과:

```text
passed=3
```

인덱스와 요소 값은 구분해야 합니다. `index >= 80`을 검사하면 점수가 아니라 위치를 비교하게 됩니다. 둘 다 i32일 수 있어 타입만으로 이런 의미 오류를 찾기 어렵습니다.

## 첫 일치 위치 찾기

못 찾은 결과를 어떻게 나타낼지 먼저 정합니다. 이 예제에서 유효한 인덱스는 0~4이므로 -1을 실패 표지로 사용합니다.

<!-- wave-example: book-array-find -->
```wave
fun main() {
    var values: array<i32, 5> = [8, 3, 8, 1, 5];
    var target: i32 = 8;
    var found: i32 = -1;

    for (var index: i32 = 0; index < 5; index += 1) {
        if (values[index] == target) {
            found = index;
            break;
        }
    }

    if (found >= 0) {
        println("found at {}", found);
    } else {
        println("not found");
    }
}
```

실행 결과:

```text
found at 0
```

첫 위치 0도 정상 결과입니다. `found > 0`으로 성공을 검사하면 첫 요소를 못 찾았다고 오해합니다. bool로 cast해 성공 여부를 판단하는 것도 같은 이유로 잘못입니다.

break를 없애면 뒤의 일치가 found를 덮어쓰므로 마지막 일치 위치를 얻게 됩니다. 문장 하나가 함수의 계약을 바꿀 수 있으므로 “찾는다”는 설명도 첫 위치인지 마지막 위치인지 구체적으로 적어야 합니다.

## 값으로 복사한 배열

배열의 값을 다른 저장 공간에 복사하려면 요소별로 읽어서 대입할 수 있습니다. 복사를 마친 뒤 한쪽 정수 요소를 바꾸어도 다른 쪽 정수 요소는 바뀌지 않습니다.

<!-- wave-example: book-array-copy -->
```wave
fun main() {
    var original: array<i32, 3> = [1, 2, 3];
    var copied: array<i32, 3>;

    for (var index: i32 = 0; index < 3; index += 1) {
        copied[index] = original[index];
    }

    copied[0] = 99;

    println("original={}", original[0]);
    println("copied={}", copied[0]);
}
```

실행 결과:

```text
original=1
copied=99
```

요소가 포인터이면 배열 값에 들어 있는 주소가 복사됩니다. 주소가 가리키는 별도 메모리까지 깊게 복사하는 것은 아닙니다. 이 차이는 뒤에서 메모리 소유권을 다룰 때 중요해집니다.

## 초기화와 유효 범위

초기값 없이 선언한 배열의 모든 요소를 읽을 수 있다고 가정하지 않습니다. 일부만 기록했다면 실제 초기화한 개수를 별도로 관리해야 합니다. 라이브러리 읽기 함수가 반환한 길이가 버퍼 전체 용량보다 작을 수 있는 이유도 같습니다.

배열 길이가 타입에 들어 있으므로 실행 중 임의로 늘어나지 않습니다. 크기가 늘어나는 바이트 목록은 `Buffer` 같은 동적 저장소를 사용합니다. 배열의 길이를 바꾸려면 타입과 초기값, 순회 상한과 그 길이에 의존하는 계산을 함께 검토해야 합니다.

## 연습: 최댓값과 위치

배열 `[4, 9, 2, 9, 1]`에서 최댓값과 그 값이 처음 나타나는 위치를 구하십시오. 모든 요소가 음수일 수도 있는 일반 함수라면 최댓값을 0으로 초기화하면 안 됩니다.

### 전체 풀이

<!-- wave-example: book-array-solution -->
```wave
fun main() {
    var values: array<i32, 5> = [4, 9, 2, 9, 1];
    var maximum: i32 = values[0];
    var position: i32 = 0;

    for (var index: i32 = 1; index < 5; index += 1) {
        if (values[index] > maximum) {
            maximum = values[index];
            position = index;
        }
    }

    println("max={} first={}", maximum, position);
}
```

실행 결과:

```text
max=9 first=1
```

첫 요소를 초기 기준으로 삼고 두 번째부터 비교합니다. `>`이므로 같은 최댓값이 다시 나와도 위치를 바꾸지 않습니다. `>=`로 바꾸면 마지막 위치가 됩니다. 길이가 0일 수 있는 인터페이스라면 첫 요소를 읽기 전에 빈 입력을 처리해야 합니다.
