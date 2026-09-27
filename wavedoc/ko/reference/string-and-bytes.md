---
translation_set_id: string-bytes
path: reference/string-and-bytes
locale: ko
group: stdlib
group_order: 1
order: 3
title: string: 문자열 길이·검색·범위
summary: NUL 종료 문자열의 바이트 단위 API와 반환값을 설명합니다.
---

## 문자열의 저장과 인자 조건

이 모듈의 `str` 인자는 접근 가능한 NUL 종료 바이트열이어야 합니다. 길이와 검색 인덱스는 바이트 단위입니다. 일반 문자는 UTF-8로 저장되지만 바이트 검색은 Unicode 정규화나 글자 단위 분할을 하지 않습니다. 반환된 인덱스가 글자 경계라고 가정하지 마십시오.

## 길이와 비교

```text
std::string::len
len(s: str) -> i32
is_empty(s: str) -> bool

std::string::cmp
eq(a: str, b: str) -> bool
cmp(a: str, b: str) -> i32
starts_with(s: str, prefix: str) -> bool
ends_with(s: str, suffix: str) -> bool
```

`len`은 마지막 NUL을 제외합니다. `cmp` 결과의 부호로 순서를 판단합니다. 반환값을 Unicode 문자 순서나 언어별 사전 정렬로 해석하지 않습니다. 이 함수들은 메모리를 할당하지 않고 입력을 변경하지 않습니다.

## 검색

`import("std::string::find")::{find, contains, count};`처럼 필요한 이름을 가져옵니다.

| 함수 선언 | 결과 |
| --- | --- |
| `find(s: str, needle: str) -> i32` | 첫 일치 위치. 없으면 -1, 빈 needle은 0 |
| `contains(s: str, needle: str) -> bool` | 포함 여부. 빈 needle은 true |
| `count(s: str, needle: str) -> i32` | 겹치지 않는 일치 개수. 빈 needle은 0 |
| `find_char(s: str, c: u8) -> i32` | 바이트의 첫 위치 또는 -1 |
| `rfind_char(s: str, c: u8) -> i32` | 바이트의 마지막 위치 또는 -1 |
| `contains_char(s: str, c: u8) -> bool` | 해당 바이트의 존재 여부 |
| `count_char(s: str, c: u8) -> i32` | 해당 바이트의 개수 |

`*_char`라는 이름의 `c`는 Unicode 코드 포인트가 아니라 한 바이트입니다. 문자열 끝의 NUL 자체는 검색 대상에 포함하지 않습니다.

## 공백을 제외한 범위

```text
std::string::trim
trim_left_index(s: str) -> i32
trim_right_index(s: str) -> i32
trim_range(s: str, out_start: ptr<i32>, out_end: ptr<i32>)
```

`trim_range`는 ASCII 공백을 제외한 반열린 범위 `[start, end)`를 출력 인자에 기록합니다. 두 출력 포인터는 쓰기 가능한 정수를 가리켜야 합니다. 원문을 수정하거나 새 문자열을 만들지 않습니다. 모두 공백이면 빈 범위가 됩니다.

## 실행 예제

`main.wave`에 저장하고 `wavec run main.wave`를 실행합니다.

<!-- wave-example: string-api -->
```wave
import("std::string::find")::{
    find, count
};
import("std::string::trim")::{
    trim_range
};

fun main() {
    var start: i32 = 0;
    var end: i32 = 0;
    trim_range("  Wave  ", &start, &end);
    println("{} {}", start, end);
    println("{} {}", find("banana", "na"), count("aaaa", "aa"));
}
```

실행 결과:

```text
2 6
2 2
```

## 관련 기능

`std::string::ascii`의 분류·대소문자 변환은 ASCII 범위용입니다. `std::string::hash`의 `djb2_32`와 `fnv1a_64`를 암호학적 해시나 비밀번호 저장에 사용하지 않습니다. NUL을 포함하는 데이터에는 [bytes](/docs/ko/stdlib/bytes)를 사용하십시오.

## 빈 검색어와 겹치는 패턴

검색 함수의 가장자리 동작을 실제 값으로 확인하면 호출 조건을 정하기 쉽습니다.

<!-- wave-example: book-string-empty-search -->
```wave
import("std::string::find")::{find, contains, count};

fun main() {
    println("empty find={}", find("Wave", ""));
    println("empty count={}", count("Wave", ""));
    println("nonoverlapping={}", count("aaaa", "aa"));

    if (contains("Wave", "")) {
        println("empty needle is contained");
    }
}
```

실행 결과:

```text
empty find=0
empty count=0
nonoverlapping=2
empty needle is contained
```

find의 성공값 0은 첫 위치입니다. count의 성공값 0은 일치가 없거나 빈 검색어 규칙의 결과입니다. 두 값을 같은 조건으로 처리하지 않습니다. 대소문자 무시나 Unicode 정규화가 필요하다면 이 바이트 검색 앞뒤에 별도 정책을 구현해야 합니다.

## trim 범위를 새 문자열로 복사하기

trim_range가 돌려주는 범위에는 종료 NUL이 새로 생기지 않습니다. 별도의 목적지에 복사할 때는 길이+1 공간을 확보하고 마지막 바이트를 직접 0으로 씁니다.

<!-- wave-example: book-trim-copy -->
```wave
import("std::string::trim")::{trim_range};

fun main() -> i32 {
    var original: str = "  Wave  ";
    var start: i32 = 0;
    var end: i32 = 0;
    var destination: array<u8, 16>;

    trim_range(original, &start, &end);

    var length: i32 = end - start;

    if (length >= 16) {
        return 1;
    }

    for (var index: i32 = 0; index < length; index += 1) {
        destination[index] = original[start + index];
    }

    destination[length] = 0;
    println("{}", &destination[0] as str);
    return 0;
}
```

실행 결과:

```text
Wave
```

길이 16인 문자열은 이 목적지에 들어가지 않습니다. 마지막 NUL까지 필요하기 때문입니다. 길이가 0인 경우에도 destination[0]=0을 기록하면 유효한 빈 문자열이 됩니다. 목적지 지역 배열은 main이 끝날 때까지 살아 있으므로 그 안에서 출력합니다.

## 문자열 API 사용 순서

문자열을 받는 함수를 설계할 때는 입력이 NUL 종료인지 확인하고, 인덱스가 바이트 단위임을 명시하고, 결과가 원문 범위인지 새 할당인지 구분합니다. 원문 범위를 반환하는 함수는 원문의 수명에 묶입니다. 새 할당을 반환하면 해제 책임도 문서화해야 합니다.

기초 개념은 [문자열 학습 장](/docs/ko/language/strings), NUL을 포함하는 데이터는 [bytes](/docs/ko/stdlib/bytes)를 읽으십시오.
