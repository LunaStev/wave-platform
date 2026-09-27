---
translation_set_id: learn-strings
path: language/strings
locale: ko
group: language
group_order: 2
order: 7
title: 7. 문자열, 문자와 바이트
summary: 문자열과 char, UTF-8 바이트 길이, NUL, 검색과 바이너리 데이터를 구분합니다.
---

## 화면의 글자와 메모리의 바이트

화면에는 글자가 보이지만 메모리에는 바이트가 저장됩니다. 특히 한국어처럼 한 글자가 여러 UTF-8 바이트인 경우 “길이”와 “몇 번째 문자”를 같은 뜻으로 사용하면 실수하기 쉽습니다.

이번 장에서는 str과 char의 차이, NUL 종료, escape, 검색 위치와 바이너리 데이터를 구분합니다. 예제는 각각 전체 프로그램입니다.

## 문자열 리터럴과 출력

<!-- wave-example: book-string-literals -->
```wave
fun main() {
    var greeting: str = "안녕하세요";

    println("{}", greeting);
    println("line one\nline two");
    println("quote: \"Wave\"");
}
```

실행 결과:

```text
안녕하세요
line one
line two
quote: "Wave"
```

큰따옴표 안의 일반 문자는 UTF-8로 표현합니다. escape는 소스에서 직접 적기 어려운 바이트를 나타냅니다. `\n`은 줄바꿈 한 바이트이며 역슬래시와 n 두 글자를 출력하는 것이 아닙니다. 역슬래시 자체를 출력하려면 `\\`를 사용합니다.

## 길이는 바이트 수

<!-- wave-example: book-utf8-length -->
```wave
import("std::string::len")::{
    len
};

fun main() {
    println("ASCII={}", len("Wave"));
    println("Korean={}", len("한"));
    println("mixed={}", len("Wave한"));
}
```

실행 결과:

```text
ASCII=4
Korean=3
mixed=7
```

len은 화면에 보이는 글자 수가 아니라 마지막 NUL 이전의 바이트 수를 반환합니다. `한`은 UTF-8 세 바이트입니다. 글자 수, Unicode 코드 포인트 수, 바이트 수는 일반적으로 서로 다른 값입니다. 표시 폭까지 구하려면 글꼴과 결합 문자 등 추가 규칙도 필요합니다.

따라서 임의의 바이트 위치에서 문자열을 잘라 화면에 출력하는 기능은 Unicode 경계를 별도로 고려해야 합니다. ASCII만 처리하는 프로그램인지 일반 Unicode 텍스트인지 입력 조건을 분명히 정하십시오.

## NUL 종료와 길이

str은 끝을 나타내는 0 바이트를 사용합니다. len은 그 마지막 바이트를 길이에 포함하지 않습니다. 문자열 리터럴 내부에 NUL을 넣는 것은 오류입니다. 다음은 의도적인 오류 예제입니다.

```wave
fun main() {
    var text: str = "left\x00right";
}
```

소스의 `\xNN`은 정확히 두 자리 16진수로 한 바이트를 지정합니다. `\x41`은 65인 바이트 A를 나타냅니다. 일반 문자를 UTF-8로 쓰는 것과 임의의 한 바이트를 넣는 것은 다르므로 `\xNN`을 쓸 수 있는 모든 str이 유효한 UTF-8인 것은 아닙니다.

## char는 Unicode 문자 전체를 담지 않음

char는 부호 없는 8비트 문자 값입니다. `'A'`처럼 한 바이트 범위의 값을 나타내는 리터럴을 사용할 수 있습니다. `'한'`은 이 범위에 들어가지 않으므로 오류입니다. `"한"`은 여러 UTF-8 바이트를 가진 str이라 별개입니다.

한 글자를 항상 char 하나에 넣으려 하지 마십시오. 텍스트를 처리하는 데 필요한 단위가 바이트인지 Unicode 코드 포인트인지 먼저 정해야 합니다.

## 문자열 비교

문자열 내용 비교에는 std의 함수를 사용합니다. 다음은 동일한 내용과 대소문자 차이를 검사하는 프로그램입니다.

<!-- wave-example: book-string-compare -->
```wave
import("std::string::cmp")::{
    eq, starts_with, ends_with
};

fun main() {
    if (eq("Wave", "Wave")) {
        println("same bytes");
    }

    if (!eq("Wave", "wave")) {
        println("case differs");
    }

    if (starts_with("report.txt", "report") && ends_with("report.txt", ".txt")) {
        println("name matches");
    }
}
```

실행 결과:

```text
same bytes
case differs
name matches
```

이 비교는 바이트열을 비교합니다. 언어별 대소문자 변환이나 Unicode 정규화를 자동으로 수행하지 않습니다. 파일 이름을 비교할 때도 OS의 파일 이름 동등성 규칙과 단순 문자열 비교가 같은 것은 아닙니다.

## 검색 결과의 단위와 실패

<!-- wave-example: book-string-search -->
```wave
import("std::string::find")::{
    find, count
};

fun main() {
    var first: i32 = find("banana", "na");
    var missing: i32 = find("banana", "xy");
    var matches: i32 = count("aaaa", "aa");

    println("first={} missing={}", first, missing);
    println("matches={}", matches);
}
```

실행 결과:

```text
first=2 missing=-1
matches=2
```

find는 첫 위치나 -1을 반환합니다. 인덱스 0도 성공이므로 `result >= 0`으로 검사합니다. count는 위치가 아니라 겹치지 않는 일치 개수입니다. 위에서 aa는 0~1과 2~3에 일치해 2번입니다.

빈 needle도 계약의 일부입니다. find는 0, contains는 true, count는 0을 반환합니다. 함수 이름이 같은 모듈에 있다고 반환 방식까지 같다고 생각하지 마십시오.

## 공백 제거는 새 문자열 생성과 다름

trim_range는 원문을 수정하거나 복사하지 않고 공백을 제외한 범위를 돌려줍니다. 출력 포인터를 받으므로 결과를 저장할 정수를 먼저 준비합니다.

<!-- wave-example: book-string-trim-range -->
```wave
import("std::string::trim")::{
    trim_range
};

fun main() {
    var start: i32 = 0;
    var end: i32 = 0;

    trim_range("  Wave  ", &start, &end);

    println("start={} end={} bytes={}", start, end, end - start);
}
```

실행 결과:

```text
start=2 end=6 bytes=4
```

범위는 `[start, end)`입니다. 시작은 포함하고 끝은 포함하지 않으므로 길이가 end-start입니다. 원문의 시작 주소에 start를 더한다고 자동으로 end 위치에 NUL이 생기는 것은 아닙니다. 범위를 별도로 들고 다니거나 새 문자열 공간을 준비해야 합니다.

## 바이너리 데이터에는 별도 길이

0을 포함하는 데이터는 문자열 끝 규칙으로 다루지 않습니다. 바이트 배열과 길이를 사용합니다.

<!-- wave-example: book-binary-not-string -->
```wave
fun main() {
    var bytes: array<u8, 3> = [65, 0, 66];

    for (var index: i32 = 0; index < 3; index += 1) {
        println("{}", bytes[index]);
    }
}
```

실행 결과:

```text
65
0
66
```

두 번째 0은 실제 데이터입니다. 이를 str로 해석하면 첫 0에서 끝난 것으로 취급해 뒤의 66을 볼 수 없습니다. 반대로 NUL이 없는 배열을 str로 cast하면 배열 밖까지 읽을 위험이 있습니다. cast는 종료 바이트를 추가하는 작업이 아닙니다.

## 연습: 파일 이름 검사

파일 이름이 `.wave`로 끝나는지 검사하고, 문자열에 `test`가 들어 있으면 테스트 파일로 출력하십시오. 이 연습은 이름의 바이트 패턴만 검사하며 실제 파일 존재 여부는 다루지 않습니다.

### 전체 풀이

<!-- wave-example: book-string-solution -->
```wave
import("std::string::cmp")::{
    ends_with
};
import("std::string::find")::{
    contains
};

fun classify(name: str) {
    if (!ends_with(name, ".wave")) {
        println("other file");
        return;
    }

    if (contains(name, "test")) {
        println("Wave test file");
    } else {
        println("Wave source file");
    }
}

fun main() {
    classify("main.wave");
    classify("parser_test.wave");
    classify("notes.txt");
}
```

실행 결과:

```text
Wave source file
Wave test file
other file
```

대문자 `.WAVE`는 어떻게 처리할지, 경로 전체에 test가 들어 있어도 테스트로 볼지 등은 별도 정책입니다. 작은 함수라도 어떤 입력을 대상으로 하는지 정해야 동작을 정확히 설명할 수 있습니다.
