---
translation_set_id: lexical
path: language/lexical-structure
locale: ko
group: language
group_order: 2
order: 14
title: 어휘 구조
summary: 식별자, 리터럴, 구분자, 키워드와 타입 이름을 설명합니다.
---

## 식별자

식별자는 변수, 함수, 타입과 필드 등에 이름을 붙입니다. 이름은 대소문자를 구분하며 문자, 숫자와 `_`를 조합할 수 있습니다. 첫 글자에는 숫자를 사용할 수 없습니다. Unicode 문자도 식별자에 사용할 수 있습니다.

```wave
var request_count: i64 = 0;
var 이름: str = "Wave";
```

실제 프로젝트에서는 도구 호환성과 검색 편의성을 위해 일관된 이름 규칙을 정해 사용하는 것이 좋습니다.

## 문장과 구분자

대부분의 선언과 표현식 문장은 `;`로 끝납니다. 함수·조건문·반복문·구조체처럼 본문을 가지는 구문은 `{ ... }` 블록을 사용합니다.

```wave
var answer: i32 = 42;

fun double(value: i32) -> i32 {
    return value * 2;
}
```

## 리터럴

```wave
var integer: i32 = 42;
var decimal: f64 = 3.14;
var text: str = "Wave";
var letter: char = 'W';
var enabled: bool = true;
var address: ptr<u8> = null;
```

정수·부동소수점·문자열·문자·불리언과 `null` 리터럴을 사용할 수 있습니다. `null`은 포인터 값에 사용하십시오.

## 키워드와 타입 이름

Wave 문법에서 사용하는 주요 키워드는 다음과 같습니다.

`pub`, `fun`, `extern`, `export`, `type`, `enum`, `static`, `var`, `deref`, `const`, `if`, `else`, `proto`, `struct`, `while`, `for`, `in`, `out`, `clobber`, `as`, `asm`, `import`, `return`, `continue`, `print`, `input`, `println`, `match`, `break`, `true`, `false`, `null`.

내장 타입 이름에는 `bool`, `char`, `byte`, `str`, 정수·부동소수점 타입, `ptr`과 `array`가 있습니다. 포인터는 `ptr<T>`, 고정 길이 배열은 `array<T, N>` 형태로 적습니다.

## 문자열과 문자 escape

| 표기 | 의미 |
| --- | --- |
| `\n` | LF 줄바꿈 |
| `\r` | CR |
| `\t` | 탭 |
| `\\` | 역슬래시 |
| `\"` | 큰따옴표 |
| `\xNN` | 정확히 두 자리 16진수로 지정한 한 바이트 |

일반 문자열 문자는 UTF-8로 저장합니다. `\xNN`은 한 바이트를 보존하므로 문자열 전체가 유효한 UTF-8이라는 보장은 없습니다. 문자열 리터럴 내부 NUL(`\x00` 포함)은 컴파일 오류입니다. 0을 포함하는 데이터에는 바이트 배열과 길이를 사용합니다.

`char` 리터럴은 8비트 값에 들어가야 합니다. `'한'`처럼 그 범위를 넘는 문자는 오류입니다. 문자열 `"한"`과는 다릅니다.

소스의 LF, CRLF, 단독 CR은 각각 하나의 논리적 줄바꿈으로 처리합니다. 이는 소스 위치와 주석 종료에 관한 규칙이며 파일 데이터의 실제 바이트를 바꾼다는 뜻은 아닙니다.

추가 문법 이름으로 `variant`, `async`, `await`가 있고 비동기 값은 `Future<T>`로 표현합니다. 위의 독립된 `var` 선언 블록은 함수 내부 코드 조각입니다.

[문자열 수업](/docs/ko/language/arrays) · [주석](/docs/ko/language/comments)

## 의도적으로 실패하는 예제

아래 프로그램을 check하면 내부 NUL 오류가 나야 합니다. 0 바이트가 필요하면 `[97, 0, 98]` 바이트 배열을 사용하십시오.

<!-- wave-example: reject-string-nul -->
```wave
fun main() {
    var text: str = "a\x00b";
}
```

아래 문자 리터럴도 8비트 범위를 넘으므로 컴파일 오류입니다. UTF-8 문자열을 표현하려면 `str`과 큰따옴표를 사용합니다.

<!-- wave-example: reject-wide-char -->
```wave
fun main() {
    var letter: char = '한';
}
```
