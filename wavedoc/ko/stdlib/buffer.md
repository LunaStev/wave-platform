---
translation_set_id: stdlib-buffer
path: stdlib/buffer
locale: ko
group: stdlib
group_order: 1
order: 5
title: buffer: 크기가 변하는 바이트 저장소
summary: Buffer 초기화, 추가, 조회, 용량과 해제 규칙을 설명합니다.
---

## Buffer의 의미

`std::buffer::types`의 `Buffer`는 `data: ptr<u8>`, `len: i64`, `cap: i64`를 가집니다. len은 초기화되어 사용 중인 바이트 수, cap은 할당된 전체 바이트 수입니다. 항상 `0 <= len <= cap`을 유지합니다. 문자열 NUL 종료를 자동으로 보장하지 않습니다.

## 기본 API

| 모듈 | 선언 | 의미 |
| --- | --- | --- |
| `std::buffer::alloc` | `buffer_init(out_buf: ptr<Buffer>, capacity: i64) -> i64` | 새 저장소 초기화. 이미 소유한 버퍼에 재호출하지 않음 |
| 같은 모듈 | `buffer_free(buf: ptr<Buffer>) -> i64` | 할당 해제. 성공하면 빈 상태 |
| 같은 모듈 | `buffer_reserve(buf: ptr<Buffer>, required_cap: i64) -> i64` | 최소 전체 용량 확보. len 유지 |
| 같은 모듈 | `buffer_clear(buf: ptr<Buffer>) -> i64` | 용량을 유지하고 len=0 |
| 같은 모듈 | `buffer_resize(buf: ptr<Buffer>, new_len: i64, value: u8) -> i64` | 길이 변경, 새 바이트를 value로 채움 |
| `std::buffer::write` | `buffer_push(buf: ptr<Buffer>, value: u8) -> i64` | 한 바이트 추가 |
| 같은 모듈 | `buffer_append(buf: ptr<Buffer>, src: ptr<u8>, size: i64) -> i64` | size바이트 복사하여 추가 |
| 같은 모듈 | `buffer_append_str(buf: ptr<Buffer>, s: str) -> i64` | NUL을 제외한 문자열 바이트 추가 |
| `std::buffer::read` | `buffer_get(buf: ptr<Buffer>, index: i64, out_value: ptr<u8>) -> i64` | 범위 안의 한 바이트 읽기 |

상태 반환 API는 성공 시 `BUFFER_OK`(0)를 반환합니다. `std::buffer::error`의 INVALID, BOUNDS, OVERFLOW, ALLOC 오류를 구분하십시오. 숫자를 OS errno로 해석하지 않습니다. `buffer_new`는 할당 실패를 빈 Buffer로 표현하므로 실패를 구분해야 할 때는 `buffer_init`을 사용합니다.

## 실행 예제

<!-- wave-example: buffer-api -->
```wave
import("std::buffer::alloc")::{
    Buffer, buffer_init, buffer_free
};
import("std::buffer::write")::{
    buffer_append_str, buffer_push
};
import("std::buffer::read")::{
    buffer_get
};

fun main() -> i32 {
    var data: Buffer;
    if (buffer_init(&data, 0) < 0) {
        return 1;
    }
    if (buffer_append_str(&data, "Hi") < 0 || buffer_push(&data, 33) < 0) {
        buffer_free(&data);
        return 2;
    }

    var value: u8 = 0;
    if (buffer_get(&data, 2, &value) < 0) {
        buffer_free(&data);
        return 3;
    }

    println("{} {}", data.len, value);
    if (buffer_free(&data) < 0) {
        return 4;
    }

    return 0;
}
```

실행 결과:

```text
3 33
```

`main.wave`로 저장해 `wavec run main.wave`로 실행합니다. 초기 용량 0은 실패가 아니라 유효한 빈 버퍼입니다. 추가 과정에서 공간을 확보합니다.

## 수명과 실패

성장하는 연산은 data 주소를 바꿀 수 있습니다. 이전에 빌린 주소를 저장해 두고 성장 후 사용하지 않습니다. Buffer를 값으로 복사해도 메모리는 깊은 복사되지 않으므로 소유자는 한 곳으로 정합니다.

`buffer_get`은 실패하면 출력 인자를 바꾸지 않습니다. 반면 편의 함수 `buffer_at`은 오류도 0으로 나타내므로 실제 0 바이트와 실패를 구분하려면 `buffer_get`을 사용합니다. 공개 필드를 직접 바꿔 유효하지 않은 len/cap을 만들지 않습니다.

[메모리 API](/docs/ko/reference/memory-and-buffer) · [파일을 Buffer로 읽는 실습](/docs/ko/practice/file-reader)

## 길이와 용량을 따로 관찰하기

reserve는 저장 공간을 확보하지만 len을 늘리지 않습니다. resize는 실제 사용 길이를 바꾸고 늘어난 부분을 지정한 바이트로 초기화합니다. clear는 사용 길이만 0으로 만들어 할당을 재사용하게 합니다.

다음 프로그램을 main.wave로 저장해 실행합니다. 용량의 정확한 성장 배수에는 의존하지 않고 필요한 공간을 확보했는지만 확인합니다.

<!-- wave-example: book-buffer-length-capacity -->
```wave
import("std::buffer::alloc")::{
    Buffer,
    buffer_init,
    buffer_reserve,
    buffer_resize,
    buffer_clear,
    buffer_free
};

fun main() -> i32 {
    var data: Buffer;

    if (buffer_init(&data, 0) < 0) {
        return 1;
    }

    if (buffer_reserve(&data, 20) < 0) {
        buffer_free(&data);
        return 2;
    }

    println("reserved length={}", data.len);

    if (buffer_resize(&data, 3, 7) < 0) {
        buffer_free(&data);
        return 3;
    }

    println("resized length={} first={}", data.len, deref data.data[0]);

    if (buffer_clear(&data) < 0) {
        buffer_free(&data);
        return 4;
    }

    println("cleared length={}", data.len);

    if (buffer_free(&data) < 0) {
        return 5;
    }

    return 0;
}
```

실행 결과:

```text
reserved length=0
resized length=3 first=7
cleared length=0
```

resize로 늘릴 때 value=7을 전달했으므로 새로 보이는 세 바이트는 모두 7입니다. reserve로 확보하기만 한 공간을 초기화된 데이터처럼 읽지 않습니다. clear 뒤에도 cap은 유지되며 같은 Buffer에 다시 추가할 수 있습니다.

## 0 바이트와 조회 실패 구분하기

buffer_get은 상태를 반환하고 실제 바이트는 출력 인자로 씁니다. 데이터가 0인 경우도 정상 성공입니다.

<!-- wave-example: book-buffer-zero-bounds -->
```wave
import("std::buffer::alloc")::{Buffer, buffer_init, buffer_free};
import("std::buffer::write")::{buffer_push};
import("std::buffer::read")::{buffer_get};
import("std::buffer::error")::{BUFFER_ERR_BOUNDS};

fun main() -> i32 {
    var data: Buffer;

    if (buffer_init(&data, 0) < 0) {
        return 1;
    }

    if (buffer_push(&data, 0) < 0) {
        buffer_free(&data);
        return 2;
    }

    var value: u8 = 99;

    if (buffer_get(&data, 0, &value) < 0) {
        buffer_free(&data);
        return 3;
    }

    println("stored={}", value);
    value = 99;

    if (buffer_get(&data, 1, &value) == BUFFER_ERR_BOUNDS) {
        println("outside, preserved={}", value);
    }

    if (buffer_free(&data) < 0) {
        return 4;
    }

    return 0;
}
```

실행 결과:

```text
stored=0
outside, preserved=99
```

첫 조회는 0을 읽는 성공이고 두 번째 조회는 범위 밖 실패입니다. 실패 뒤 value=99가 유지되어도 그것이 버퍼에서 읽은 값이라는 뜻은 아닙니다. 반드시 상태를 함께 확인합니다.

## 연습 풀이: 바이트 누적

0부터 9까지 추가하려면 buffer_push를 반복하고 각 결과를 확인합니다. 합계는 i64에 저장하고 `0 <= index < data.len` 범위만 읽습니다. 버퍼를 다룬 뒤에는 성공과 실패 경로 모두에서 buffer_free를 호출합니다.

<!-- wave-example: book-buffer-sum -->
```wave
import("std::buffer::alloc")::{Buffer, buffer_init, buffer_free};
import("std::buffer::write")::{buffer_push};

fun main() -> i32 {
    var data: Buffer;

    if (buffer_init(&data, 0) < 0) {
        return 1;
    }

    for (var value: i32 = 0; value < 10; value += 1) {
        if (buffer_push(&data, value as u8) < 0) {
            buffer_free(&data);
            return 2;
        }
    }

    var total: i64 = 0;

    for (var index: i64 = 0; index < data.len; index += 1) {
        total += deref data.data[index] as i64;
    }

    println("sum={}", total);

    if (buffer_free(&data) < 0) {
        return 3;
    }

    return 0;
}
```

실행 결과:

```text
sum=45
```
