---
translation_set_id: memory-model
path: language/explicit-memory-type-model
locale: ko
group: language
group_order: 2
order: 9
title: 9. 포인터, 원본 변경과 수명
summary: 주소 얻기, 역참조, 원본 변경과 댕글링 포인터를 배웁니다.
---

## 값과 저장 위치 구분하기

정수 42와 그 정수가 저장된 주소는 다른 값입니다. 포인터는 저장 위치를 가리킵니다. 주소를 전달하면 함수가 호출자의 저장 공간을 읽거나 변경할 수 있습니다.

이번 장에서는 주소 얻기, 역참조, 원본 변경, 포인터 산술과 수명을 배웁니다. 동적 할당은 다음 장에서 다룹니다. 먼저 지역 변수와 배열의 주소로 시작합니다.

## 주소 얻기와 역참조

<!-- wave-example: book-pointer-first -->
```wave
fun main() {
    var value: i32 = 42;
    var address: ptr<i32> = &value;

    println("value={}", value);
    println("through pointer={}", deref address);

    deref address = 99;
    println("changed={}", value);
}
```

실행 결과:

```text
value=42
through pointer=42
changed=99
```

`&value`는 주소를 얻고 `deref address`는 그 주소의 값을 읽거나 쓰는 식입니다. address에 99를 저장한 것이 아니라 address가 가리키는 정수에 99를 썼습니다. 변수 address 자체는 계속 value를 가리킵니다.

`ptr<i32>`는 i32 저장 공간에 접근할 포인터 타입입니다. 타입 이름에는 길이나 자동 해제 규칙이 들어 있지 않습니다.

## 포인터 자체를 바꾸기

<!-- wave-example: book-pointer-reassign -->
```wave
fun main() {
    var first: i32 = 10;
    var second: i32 = 20;
    var selected: ptr<i32> = &first;

    selected = &second;
    deref selected = 25;

    println("{} {}", first, second);
}
```

실행 결과:

```text
10 25
```

`selected = &second`는 포인터 변수에 다른 주소를 저장합니다. first의 값은 변경하지 않습니다. 그다음 deref로 값을 쓰면 second가 바뀝니다. “주소의 변경”과 “주소를 통한 값 변경”을 문장별로 구분하면 혼란이 줄어듭니다.

## 함수가 원본을 바꾸게 하기

<!-- wave-example: book-pointer-increment -->
```wave
fun increment(value: ptr<i32>) {
    deref value = deref value + 1;
}

fun main() {
    var count: i32 = 4;

    increment(&count);
    increment(&count);

    println("{}", count);
}
```

실행 결과:

```text
6
```

함수에 count의 값 4가 아니라 주소를 넘겼습니다. 함수는 그 저장 공간을 변경하므로 호출자의 count도 바뀝니다. 이 함수의 입력 조건은 유효하고 쓰기 가능한 i32 주소라는 것입니다. null을 전달해도 괜찮은 함수가 아닙니다.

null을 허용할지 여부는 함수마다 정할 수 있습니다. 허용하지 않는 함수를 호출할 때는 호출자가 주소를 보장해야 합니다. 허용한다면 함수가 null을 처리하는 경로를 제공해야 합니다.

## null을 처리하는 함수

<!-- wave-example: book-pointer-null -->
```wave
fun try_increment(value: ptr<i32>) -> bool {
    if (value == null) {
        return false;
    }

    deref value = deref value + 1;
    return true;
}

fun main() {
    var count: i32 = 7;

    if (!try_increment(null)) {
        println("no value");
    }

    if (try_increment(&count)) {
        println("count={}", count);
    }
}
```

실행 결과:

```text
no value
count=8
```

null 검사는 주소가 없다는 경우만 처리합니다. null이 아닌 임의의 숫자를 포인터로 바꿨다고 유효한 메모리가 생기지는 않습니다. 올바른 수명, 크기, 정렬과 접근 권한까지 있어야 읽고 쓸 수 있습니다.

## 배열의 주소와 요소 단위 이동

<!-- wave-example: book-pointer-array -->
```wave
fun main() {
    var values: array<i32, 3> = [10, 20, 30];
    var first: ptr<i32> = &values[0];
    var second: ptr<i32> = first + 1;

    println("{}", deref first);
    println("{}", deref second);
    println("{}", deref first[2]);
}
```

실행 결과:

```text
10
20
30
```

포인터에 1을 더하면 대상 타입 한 요소만큼 이동합니다. i32의 다음 요소와 u8의 다음 요소는 이동 바이트 수가 다릅니다. `first + 1`에 다시 타입 크기를 곱해서 더하면 원하지 않는 위치로 이동합니다.

포인터 인덱싱 역시 유효 범위 안에서 해야 합니다. first는 배열 길이 3을 스스로 기억하지 않으므로 함수에 범위를 전달할 때는 포인터와 길이를 함께 받는 형태를 사용합니다.

## 읽을 범위를 함수에 전달하기

<!-- wave-example: book-pointer-range -->
```wave
fun sum(values: ptr<i32>, count: i32) -> i32 {
    var total: i32 = 0;

    for (var index: i32 = 0; index < count; index += 1) {
        total += deref values[index];
    }

    return total;
}

fun main() {
    var values: array<i32, 4> = [2, 4, 6, 8];

    println("first two={}", sum(&values[0], 2));
    println("all={}", sum(&values[0], 4));
}
```

실행 결과:

```text
first two=6
all=20
```

count의 단위는 요소 수입니다. 이 함수는 count개의 읽을 수 있는 i32가 있다는 조건을 호출자에게 요구합니다. 실제 배열보다 큰 길이를 넘기면 계약을 위반합니다. 같은 ptr<u8>·i64 조합을 사용하는 API라도 길이가 바이트 수인지 요소 수인지 문서에서 확인해야 합니다.

## 수명: 주소가 언제까지 유효한가

지역 변수는 해당 호출과 블록의 유효 기간 안에서 사용합니다. 함수 안의 지역 변수 주소를 반환해 호출자가 나중에 읽게 하면 그 저장 공간의 수명이 끝나 있을 수 있습니다.

다음은 실행하면 안 되는 잘못된 설계의 조각입니다.

```wave
fun invalid_address() -> ptr<i32> {
    var local: i32 = 42;

    return &local;
}
```

값 하나가 필요하면 i32를 반환합니다. 호출자가 마련한 저장 공간에 써야 한다면 포인터를 입력으로 받습니다. 호출 이후에도 유지할 별도 저장 공간이 필요하면 명시적으로 할당하고 해제 책임을 전달합니다.

## 같은 저장 공간을 가리키는 두 포인터

포인터를 복사하면 같은 주소를 가리키는 이름이 하나 더 생깁니다. 메모리를 복제하지 않습니다.

<!-- wave-example: book-pointer-alias -->
```wave
fun main() {
    var value: i32 = 1;
    var first: ptr<i32> = &value;
    var second: ptr<i32> = first;

    deref second = 9;

    println("{} {}", value, deref first);
}
```

실행 결과:

```text
9 9
```

second를 통해 바꾼 결과가 first를 통해서도 보입니다. 원본 메모리가 해제되면 두 포인터 모두 사용할 수 없습니다. 한 포인터 변수에 null을 대입해도 다른 복사본까지 자동으로 바뀌지는 않습니다.

## 연습: 두 정수 교환

두 i32 주소를 받아 값을 교환하는 함수를 작성하십시오. 첫 값을 덮어쓰기 전에 임시 변수에 저장해야 합니다. 동일한 주소를 두 번 전달해도 값이 유지되는지 확인하십시오.

### 전체 풀이

<!-- wave-example: book-pointer-swap -->
```wave
fun swap(left: ptr<i32>, right: ptr<i32>) {
    var saved: i32 = deref left;

    deref left = deref right;
    deref right = saved;
}

fun main() {
    var first: i32 = 3;
    var second: i32 = 8;

    swap(&first, &second);
    println("{} {}", first, second);

    swap(&first, &first);
    println("{}", first);
}
```

실행 결과:

```text
8 3
8
```

이 함수도 두 주소가 유효한 쓰기 가능한 정수라는 조건을 가집니다. null 처리까지 필요하면 이전 try_increment와 같은 결과 반환을 추가할 수 있습니다.


## **Wave Explicit Memory Type Model**

Wave의 포인터 설계는 **Wave Explicit Memory Type Model**을 기반으로 합니다. 이 모델은 포인터와 배열을 문법적 트릭이나 라이브러리 추상화가 아닌, 언어 차원의 명시적인 메모리 타입으로 정의합니다.

`ptr<T>`는 `T` 값을 저장한 메모리 주소를 가리키는 타입이고, `array<T, N>`은 `T` 값 `N`개를 연속해서 저장하는 고정 길이 메모리 타입입니다. 따라서 함수 인자, 반환값, 구조체 필드와 다른 타입 안에서도 포인터와 배열의 구조가 그대로 드러납니다.

## null

```wave
var buffer: ptr<u8> = null;
if (buffer == null) {
    println("no buffer");
}
```

`null`은 유효한 메모리 주소를 가리키지 않는 포인터 값입니다. `null`은 `ptr<T>` 타입에만 대입할 수 있으며 정수, 불리언이나 배열 값으로 사용할 수 없습니다.

메모리 할당이나 검색처럼 결과가 없을 수 있는 함수는 `null`을 반환할 수 있습니다. 이런 결과는 `null`과 비교한 뒤에만 역참조합니다. `null` 포인터를 역참조하면 유효한 저장소에 접근할 수 없습니다.

## 포인터 변환

주소나 다른 포인터 표현을 바꿔야 할 때는 `as`를 사용합니다.

```wave
var raw: i64 = 0;
var p: ptr<u8> = raw as ptr<u8>;
```

정수와 포인터 사이의 변환은 저수준 경계에서만 사용하고, 대상 플랫폼의 주소 너비와 ABI를 고려하십시오.
