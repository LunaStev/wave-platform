---
translation_set_id: stdlib-random
path: stdlib/random
locale: ko
group: stdlib
group_order: 1
order: 10
title: random: OS 난수 채우기
summary: OS 엔트로피로 버퍼를 채우고 부분 실패를 처리합니다.
---

## API

```text
std::random::fill
random_available() -> bool
random_fill(buffer: ptr<u8>, size: i64) -> RandomFillResult
```

size는 바이트 수이고 저장 공간은 호출자가 제공합니다. `RandomFillResult`에는 ok, written, error가 있습니다. 성공하면 written은 요청 길이입니다. 실패하면 written은 유효하게 채워진 접두사의 길이이며 나머지를 난수로 사용하지 않습니다.

`random_available`은 OS 난수 기능의 지원 여부를 알려줍니다. 개별 요청의 성공 여부는 random_fill의 결과로 확인합니다. OS 엔트로피만 사용하고 실패 시 시간값이나 약한 PRNG로 대체하지 않습니다. size=0은 null과 함께 전달해도 성공합니다. 음수 길이 또는 양수 길이에 null은 오류입니다.

## 실행 예제

<!-- wave-example: random-api -->
```wave
import("std::random::fill")::{
    RandomFillResult, random_fill
};

fun main() -> i32 {
    var bytes: array<u8, 16>;
    var result: RandomFillResult = random_fill(&bytes[0], 16);
    if (!result.ok) {
        println("random error={}", result.error);
        return 1;
    }

    println("filled={}", result.written);
    return 0;
}
```

실행 결과:

```text
filled=16
```

`main.wave`로 저장해 실행합니다. 바이트 내용은 매번 다르므로 특정 값을 예상하지 않습니다. 실패하면 result.error로 원인을 확인합니다. 난수 바이트를 문자로 그대로 출력하지 말고 필요하면 별도 인코딩을 사용하십시오.

## 요청이 잘못된 경우

0바이트 요청은 쓸 내용이 없으므로 성공합니다. 반면 양수 길이에 null을 넘기면 데이터를 쓸 수 없어 실패합니다. 아래 코드는 메모리를 할당하지 않고 두 상황을 비교합니다.

<!-- wave-example: book-random-boundaries -->
```wave
import("std::random::fill")::{
    RandomFillResult,
    random_fill
};

fun main() -> i32 {
    var empty: RandomFillResult = random_fill(null, 0);

    if (!empty.ok || empty.written != 0) {
        return 1;
    }

    println("empty request succeeded");

    var invalid: RandomFillResult = random_fill(null, 16);

    if (invalid.ok) {
        return 2;
    }

    println("missing buffer rejected");
    return 0;
}
```

실행 결과:

```text
empty request succeeded
missing buffer rejected
```

## 일부만 채워진 버퍼 처리

16바이트를 요청했지만 실패하면서 written=8이 반환되었다면 앞의 8바이트만 채워진 상태입니다. 16바이트짜리 식별자를 만드는 작업이라면 성공한 식별자가 아니므로 전체 결과를 버리고 실패를 전달합니다. 남은 8바이트를 0으로 채운 뒤 성공으로 처리하면 안 됩니다.

난수를 정수 범위에 넣을 때도 분포를 고려해야 합니다. 균일한 u8 값에 `% 10`을 적용하면 0~5가 6~9보다 더 자주 나옵니다. 256이 10으로 나누어떨어지지 않기 때문입니다. 공정한 선택이 필요하다면 250~255를 버리고 다시 뽑은 뒤 나머지 연산을 적용하는 방식으로 편향을 제거할 수 있습니다.

난수 바이트의 저장 공간과 수명은 호출자가 관리합니다. 배열을 사용하면 배열 범위 안에서 처리하고, 동적 메모리를 사용하면 사용이 끝난 뒤 해제합니다. 반환 구조체가 버퍼를 대신 소유하지는 않습니다.
