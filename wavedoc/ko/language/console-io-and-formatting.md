---
translation_set_id: console-io-formatting
path: language/console-io-and-formatting
locale: ko
group: language
group_order: 2
order: 16
title: 콘솔 입출력과 포매팅
summary: print, println, input 문장과 자리표시자 규칙을 설명합니다.
---

## 입출력 문장

Wave는 `print`, `println`, `input`을 콘솔 입출력 문장으로 제공합니다.

```wave
fun main() {
    var count: i32 = 0;
    input("{}", count);
    print("count = ");
    println("{}", count);
}
```

각 문장은 `;`로 끝납니다. 첫 번째 인자는 반드시 문자열 리터럴이어야 합니다. 변수나 계산한 문자열을 format 인자로 사용할 수 없습니다.

## 자리표시자

정확히 두 글자인 `{}`만 자리표시자입니다.

```wave
println("name = {}, score = {}", name, score);
```

자리표시자 수와 뒤따르는 식의 수는 정확히 같아야 합니다. 수가 다르면 문법 오류입니다.

```wave
println("{} {}", one);
// 오류: 자리표시자 2개, 값 1개
println("plain text", one);
// 오류: 자리표시자 없음, 값이 남음
```

다른 형태의 중괄호는 일반 텍스트로 남습니다. 이 문법에는 이름이나 번호가 붙은 자리표시자가 없습니다.

## print와 println

`print`는 포매팅한 텍스트를 그대로 출력하고 `println`은 줄바꿈을 덧붙입니다.

```wave
print("loading...");
println("done");
```

포매팅 인자에는 정수, 부동소수점, 문자열과 포인터 같은 스칼라 값을 사용합니다. 배열과 구조체는 포매팅 인자로 사용할 수 없습니다.

## input 대상

`input`은 읽은 값을 대상에 저장하므로 format 뒤의 모든 식은 쓰기 가능한 위치여야 합니다.

```wave
var number: i32 = 0;
input("{}", number);
```

변수, 필드와 역참조한 저장 위치를 대상으로 사용할 수 있습니다. 리터럴과 계산 결과는 입력 대상으로 사용할 수 없습니다.

입력값을 요청한 타입으로 모두 변환하지 못하면 프로그램은 실패 상태로 종료합니다.

## 런타임 경계

이 문장들은 hosted 환경의 콘솔 입출력을 사용합니다. 프리스탠딩 환경에서는 커널이나 장치가 제공하는 입출력을 함수 또는 FFI 경계로 정의해야 합니다.

## 입력 값과 범위

bool 입력은 `0`과 `1`만 허용합니다. 2를 true로 해석하거나 `true`라는 문자열을 같은 입력으로 받지 않습니다. 정수 입력은 대상 정수 폭의 범위 안에 있어야 합니다. 128·256·512·1024비트 정수도 타입의 전체 폭을 기준으로 처리합니다.

형식 오류, 범위 초과, 필요한 입력 전 EOF는 실패입니다. 내장 input은 실패를 반환해 재입력하는 함수가 아니라 프로세스를 실패 종료시키는 입력 기능입니다. 복구 가능한 입력 처리가 필요하면 io로 바이트를 읽고 별도 파서를 구성합니다.

[입력 계산기 실습](/docs/ko/practice/input-calculator) · [파일과 io](/docs/ko/stdlib/files-io)
