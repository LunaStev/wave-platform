---
translation_set_id: modules-ffi
path: language/modules-imports-and-ffi
locale: ko
group: language
group_order: 2
order: 11
title: 11. 파일을 나누고 제네릭으로 재사용하기
summary: 공개 이름 가져오기와 명시적 타입 인자를 배웁니다.
---

## 파일을 나누는 이유

프로그램이 커지면 모든 함수를 main.wave에 두기보다 관련 기능끼리 묶는 편이 찾기 쉽습니다. 모듈 경계에서는 어떤 이름을 다른 코드에 공개할지 정합니다. 제네릭은 파일 분리와 별개로 타입이 다른 같은 작업을 재사용하는 도구입니다.

이번 장에서는 두 파일짜리 프로그램을 만든 뒤 모듈 별칭, 선택 import, 제네릭 함수와 구조체를 배웁니다.

## 두 파일짜리 프로그램

같은 디렉터리에 helpers.wave와 main.wave를 만듭니다.

helpers.wave 전체:

```wave
pub fun double(value: i32) -> i32 {
    return value * 2;
}
```

main.wave 전체:

<!-- wave-example: book-module-two-files -->
```wave
import("./helpers")::{
    double
};

fun main() {
    var result: i32 = double(21);

    println("{}", result);
}
```

실행 결과:

```text
42
```

터미널에서 `wavec run main.wave`로 실행합니다. helpers.wave도 따로 실행하는 것이 아닙니다. import를 통해 필요한 소스가 연결됩니다.

helpers의 함수 앞에 붙은 pub은 다른 모듈이 가져올 수 있음을 나타냅니다. 외부에 공개할 필요가 없는 보조 함수는 공개하지 않아도 됩니다. 모듈 내부 구현을 바꾸더라도 공개 함수의 계약을 지키면 사용하는 코드의 변경을 줄일 수 있습니다.

## 상대 경로의 기준

`./helpers`는 import 문장을 작성한 소스 파일의 디렉터리를 기준으로 합니다. 프로그램 실행 시 파일 I/O가 쓰는 작업 디렉터리와 구분하십시오. import 파일을 찾는 단계와 실행 중 input.txt를 찾는 단계는 다릅니다.

로컬 import에서는 `.wave` 확장자를 생략할 수 있습니다. 디렉터리를 나누었다면 위치에 맞게 `./이름` 경로를 작성합니다. 패키지 의존성의 이름을 가져오는 경로와 로컬 상대 경로를 혼동하지 마십시오.

## 선택 import와 별칭

선택 import는 원하는 공개 이름만 현재 파일에서 직접 사용하게 합니다. 이름 충돌이 있거나 어느 모듈의 함수인지 드러내고 싶으면 별칭을 사용합니다.

<!-- wave-example: book-module-alias -->
```wave
import("std::string::len" as strings);

fun main() {
    var length: i32 = strings::len("Wave");

    println("{}", length);
}
```

실행 결과:

```text
4
```

strings는 이 파일에서 정한 모듈 별칭입니다. `strings::len`은 그 모듈의 이름을 사용합니다. 필드 접근의 점과 모듈 구분의 `::`는 다른 표기입니다.

선택 import와 별칭 import를 하나의 문장에 함께 쓰지 않습니다. 어떤 스타일이든 파일 전체에서 이름의 출처를 쉽게 읽을 수 있도록 일관되게 사용합니다.

## 표준 라이브러리와 패키지

`std::` 경로는 표준 라이브러리를 가리킵니다. 사용자는 필요한 모듈의 공개 API를 import합니다. 표준 라이브러리 함수가 모두 자동으로 현재 이름 공간에 들어오는 것은 아닙니다.

외부 패키지 경로는 패키지 이름에서 시작합니다. 패키지의 위치는 컴파일러 옵션이나 패키지 관리자가 제공합니다. 먼저 로컬 모듈로 경계를 익힌 뒤 [Vex 사용법](/docs/ko/whale/vex-package-manager)에서 의존성을 관리하는 방법을 배우면 됩니다.

## 타입별로 같은 함수를 쓰기

다음 함수는 입력을 그대로 반환합니다. i32와 str을 위해 같은 코드를 두 번 쓰지 않도록 타입 매개변수 T를 사용합니다.

<!-- wave-example: book-generic-identity -->
```wave
fun identity<T>(value: T) -> T {
    return value;
}

fun main() {
    var number: i32 = identity<i32>(42);
    var text: str = identity<str>("Wave");

    println("{} {}", number, text);
}
```

실행 결과:

```text
42 Wave
```

T는 실행 중 전달하는 정수 값이 아니라 타입을 넣는 자리입니다. `<i32>`처럼 타입 인자를 명시해서 호출합니다. 일반 사용자 제네릭 함수에서는 타입 인자를 생략하지 않습니다.

identity<str>은 문자열 바이트를 새로 할당해서 복제하지 않습니다. 값을 그대로 반환합니다. 제네릭이라는 문법이 데이터의 복사·소유권 규칙을 바꾸지는 않습니다.

## 제네릭 본문이 요구하는 연산

타입 매개변수가 있다고 모든 연산을 모든 타입에 사용할 수 있는 것은 아닙니다. 아래 minimum은 비교가 가능한 실제 타입으로 사용해야 합니다.

<!-- wave-example: book-generic-minimum -->
```wave
fun minimum<T>(left: T, right: T) -> T {
    if (left < right) {
        return left;
    }

    return right;
}

fun main() {
    var small: i32 = minimum<i32>(7, 4);
    var wide: i64 = minimum<i64>(100, 20);

    println("{} {}", small, wide);
}
```

실행 결과:

```text
4 20
```

타입 인자를 바꾸면 본문에서 사용하는 `<`와 반환이 해당 타입에서 성립해야 합니다. 제네릭 오류를 읽을 때는 호출한 타입 조합과 함수 본문이 요구하는 연산을 함께 확인합니다.

## 제네릭 구조체

서로 다른 두 값을 묶는 Pair를 만들어 봅니다.

<!-- wave-example: book-generic-pair -->
```wave
struct Pair<A, B> {
    first: A;
    second: B;
}

fun main() {
    var item: Pair<i32, str> = Pair<i32, str> {
        first: 7,
        second: "seven"
    };

    println("{} {}", item.first, item.second);
}
```

실행 결과:

```text
7 seven
```

Pair<i32, str>와 Pair<i64, str>는 서로 다른 구체적 타입입니다. 타입 인자의 순서도 의미가 있습니다. first와 second의 타입이 어디에서 결정되는지 선언과 생성 코드를 연결해서 읽으십시오.

## 공개 API의 이름과 계약

함수를 공개할 때는 이름뿐 아니라 입력 단위, 반환값, 실패와 소유권을 함께 정합니다. 예를 들어 read가 최대 길이를 읽는지 정확한 길이를 읽는지에 따라 호출자가 작성할 반복문이 달라집니다.

pub은 Wave 모듈 사이의 공개 범위입니다. 다른 언어가 호출할 외부 심볼을 내보내는 export(c)와는 다른 기능입니다. [FFI 참조](/docs/ko/language/modules-imports-and-ffi)에서 두 언어를 연결하는 완성 예제를 확인할 수 있습니다.

## 연습과 전체 풀이

math.wave에 공개 함수 square를 만들고 main.wave에서 별칭으로 불러 3과 5의 제곱을 출력하십시오.

math.wave:

```wave
pub fun square(value: i32) -> i32 {
    return value * value;
}
```

main.wave:

<!-- wave-example: book-module-solution -->
```wave
import("./math" as math);

fun main() {
    println("{} {}", math::square(3), math::square(5));
}
```

실행 결과:

```text
9 25
```

함수가 사라졌다는 오류라면 import 경로와 pub을 먼저 확인합니다. 이름 충돌이라면 별칭이 있는 호출인지 확인합니다. 타입 오류라면 함수의 입력과 전달한 인자 타입을 확인합니다. 서로 다른 문제를 경로 수정 하나로 해결하려 하지 마십시오.


## C 함수 가져오기

```wave
extern(c) fun puts(text: ptr<i8>) -> i32;
```

ABI 이름 뒤에는 실제 심볼 이름을 문자열로 지정할 수 있습니다.

```wave
extern(c, "native_symbol") fun local_name(value: i32) -> i32;
```

## Wave 함수 내보내기

```wave
export(c) fun wave_add(left: i32, right: i32) -> i32 {
    return left + right;
}
```

`extern`과 `export`는 단일 함수와 블록 형태로 사용할 수 있습니다. 내보내는 함수는 구체적인 ABI 시그니처를 가져야 하므로 제네릭일 수 없습니다.

## 대상 조건 속성

최상위 항목에는 대상 조건 속성을 붙일 수 있습니다.

```wave
#[target(os="linux", arch="x86_64")]
extern(c) fun platform_call(value: i32) -> i32;
```

조건 키는 `arch`, `os`, `env`, `abi`이며 속성은 바로 다음 최상위 항목에 적용됩니다.

## 직접 작성한 C 함수와 연결하기

이 실습은 C 컴파일러가 있는 네이티브 환경용입니다. 라이브러리의 할당이나 문자열 처리 없이 정수 함수 하나를 연결합니다.

`native.c`:

```c
#include <stdint.h>
int32_t native_double(int32_t value) { return value * 2; }
```

`main.wave`:

```wave
extern(c) fun native_double(value: i32) -> i32;

fun main() {
    println("{}", native_double(21));
}
```

Linux/macOS에서 같은 작업 디렉터리의 터미널로 실행합니다.

```shell
cc -c native.c -o native.o
wavec build main.wave native.o -o ffi-example
./ffi-example
```

기대 출력은 `42`입니다. Windows의 MSVC 개발자 셸에서는 `cl /c native.c /Fonative.obj`로 object를 만들고 `wavec build main.wave native.obj -o ffi-example.exe`로 연결합니다. 소스와 object의 대상 아키텍처가 같아야 합니다. 예제는 작은 값만 사용합니다. C 함수에 큰 값을 넘기려면 C 쪽 곱셈의 범위도 별도로 보장해야 합니다.

로컬 파일 경로는 `./`로 시작합니다.
