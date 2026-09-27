---
translation_set_id: learn-allocation
path: language/allocation
locale: ko
group: language
group_order: 2
order: 10
title: 10. 메모리 할당과 자원 관리
summary: 할당 실패, 초기화, 유효 범위와 해제를 배웁니다.
---

## 언제 동적 저장 공간이 필요한가

고정 배열은 크기를 타입에 적습니다. 파일 크기나 입력 길이처럼 실행 중 알게 되는 양의 데이터를 저장하려면 동적 메모리를 사용할 수 있습니다. 할당한 공간은 더 이상 필요하지 않을 때 해제해야 합니다.

이번 장에서는 작은 할당을 직접 관리한 뒤 재할당과 Buffer로 발전시킵니다. 포인터만 전달하는 것과 소유권을 넘기는 것은 다릅니다. 각 함수가 어떤 자원을 소유하는지 따라가며 읽으십시오.

## 할당·검사·사용·해제

<!-- wave-example: book-alloc-lifecycle -->
```wave
import("std::mem::alloc")::{
    mem_alloc_zeroed, mem_free
};

fun main() -> i32 {
    var size: i64 = 4;
    var data: ptr<u8> = mem_alloc_zeroed(size);

    if (data == null) {
        println("allocation failed");
        return 1;
    }

    deref data[0] = 42;
    println("{} {}", deref data[0], deref data[1]);

    var status: i64 = mem_free(data, size);

    if (status < 0) {
        return 2;
    }

    return 0;
}
```

실행 결과:

```text
42 0
```

네 단계가 보입니다. 4바이트를 요청하고, null인지 확인하고, 유효 범위에서 사용하고, 마지막에 해제합니다. mem_alloc_zeroed는 성공한 메모리를 0으로 초기화하므로 아직 쓰지 않은 두 번째 바이트가 0입니다.

일반 mem_alloc의 초기 내용은 가정하지 않습니다. 읽을 영역은 먼저 초기화해야 합니다. 크기 0이나 음수의 할당은 null을 반환하며 양수여도 메모리 확보에 실패할 수 있습니다.

## 크기의 단위

메모리 할당 API의 size는 바이트 수입니다. 정수 열 개의 공간을 얻는다면 요소 크기에 개수를 곱해야 합니다. 곱셈이 범위를 넘지 않는지도 검사합니다.

<!-- wave-example: book-alloc-typed -->
```wave
import("std::mem::alloc")::{
    mem_alloc, mem_free
};
import("std::mem::layout")::{
    size_of
};
import("std::mem::ops")::{
    mem_size_mul_checked
};

fun main() -> i32 {
    var count: i64 = 3;
    var bytes: i64 = 0;
    var item_size: i64 = size_of<i32>() as i64;

    if (mem_size_mul_checked(count, item_size, &bytes) < 0) {
        return 1;
    }

    var raw: ptr<u8> = mem_alloc(bytes);

    if (raw == null) {
        return 2;
    }

    var values: ptr<i32> = raw as ptr<i32>;

    for (var index: i64 = 0; index < count; index += 1) {
        deref values[index] = (index + 1) as i32;
    }

    println("{} {} {}", deref values[0], deref values[1], deref values[2]);

    if (mem_free(raw, bytes) < 0) {
        return 3;
    }

    return 0;
}
```

실행 결과:

```text
1 2 3
```

count는 요소 수, bytes는 바이트 수입니다. 포인터로 요소에 접근할 때는 i32 단위로 이동하지만 해제할 때는 원래 할당한 바이트 수를 전달합니다. size_of는 컴파일 대상의 크기를 돌려주므로 임의의 숫자를 넣는 것보다 타입과 연결이 명확합니다.

일반적인 아주 큰 타입에 size_of 결과를 i64로 바꾸는 경우에는 그 변환 범위도 고려해야 합니다. 여기서는 크기가 알려진 i32를 사용합니다.

## 실패 경로에서도 정리하기

할당 뒤 다른 작업이 실패하면 조기 return 전에 이미 얻은 메모리를 정리합니다. 소유권을 표로 적으면 놓치는 경로를 찾기 쉽습니다.

| 단계 | 소유한 자원 | 실패하면 |
| --- | --- | --- |
| 할당 전 | 없음 | 바로 반환 |
| 할당 성공 후 | data와 원래 크기 | data 해제 후 반환 |
| 재할당 성공 후 | 새 주소와 새 크기 | 새 주소 해제 |
| 해제 후 | 없음 | 이전 주소 사용 금지 |

포인터 변수에 다른 값을 대입해 버려 원래 주소를 잃으면 해제할 방법도 잃습니다. 이런 문제가 누수로 이어집니다. 반대로 같은 할당을 두 소유자가 각각 해제하면 중복 해제가 됩니다.

## 크기를 늘리는 재할당

<!-- wave-example: book-alloc-grow -->
```wave
import("std::mem::alloc")::{
    mem_alloc, mem_realloc, mem_free
};

fun main() -> i32 {
    var data: ptr<u8> = mem_alloc(4);

    if (data == null) {
        return 1;
    }

    deref data[0] = 7;
    var next: ptr<u8> = mem_realloc(data, 4, 8);

    if (next == null) {
        mem_free(data, 4);
        return 2;
    }

    data = next;
    deref data[4] = 9;
    println("{} {}", deref data[0], deref data[4]);

    if (mem_free(data, 8) < 0) {
        return 3;
    }

    return 0;
}
```

실행 결과:

```text
7 9
```

새 주소를 next에 먼저 받았습니다. 새 할당이 실패하면 기존 data가 유지되므로 그 주소로 정리할 수 있습니다. 성공하면 이전 공간은 해제되고 새 주소를 사용합니다. 늘어난 영역의 값은 직접 초기화합니다.

이 예제의 새 크기는 양수입니다. new_size=0은 기존 공간 해제를 시도하고 null을 반환하는 별도 동작이므로 “null이면 언제나 옛 공간이 살아 있다”고 생각하면 안 됩니다. 해제 상태를 확인해야 할 때는 mem_free를 직접 호출하십시오.

## 빌린 포인터를 다시 확인해야 하는 시점

data 내부를 가리키는 포인터를 저장해 두었다가 재할당 후 사용하는 것은 잘못입니다. 새 data의 주소가 달라질 수 있기 때문입니다. 내부 위치가 필요하면 주소 대신 offset을 저장하고 성공 후 새 data를 기준으로 다시 계산할 수 있습니다.

해제와 재할당은 소유자만의 문제가 아닙니다. 그 메모리를 빌린 다른 코드가 아직 사용하는지도 확인해야 합니다. 비동기 작업에 버퍼를 전달했다면 작업이 완료할 때까지 수명을 유지해야 합니다.

## 바이트 목록에는 Buffer

길이가 자주 바뀌는 바이트 목록을 직접 재할당하면서 관리하려면 len과 cap, 확장 실패, 크기 계산을 모두 처리해야 합니다. std의 Buffer는 이런 작업을 묶습니다.

<!-- wave-example: book-buffer-progressive -->
```wave
import("std::buffer::alloc")::{
    Buffer, buffer_init, buffer_free
};
import("std::buffer::write")::{
    buffer_append_str
};

fun main() -> i32 {
    var message: Buffer;

    if (buffer_init(&message, 0) < 0) {
        return 1;
    }

    if (buffer_append_str(&message, "Hello") < 0) {
        buffer_free(&message);
        return 2;
    }

    if (buffer_append_str(&message, ", Wave") < 0) {
        buffer_free(&message);
        return 3;
    }

    println("bytes={}", message.len);

    if (buffer_free(&message) < 0) {
        return 4;
    }

    return 0;
}
```

실행 결과:

```text
bytes=11
```

len은 사용 중인 바이트 수이고 cap은 할당 용량입니다. 용량이 부족하면 추가 과정에서 확장합니다. Buffer를 사용해도 마지막 해제 책임은 호출자에게 있습니다.

buffer_append_str은 문자열 끝의 NUL을 추가하지 않습니다. 따라서 message.data를 곧바로 str로 출력하면 안 됩니다. 바이트열은 길이를 받는 I/O 함수로 출력하거나 명시적으로 문자열 표현을 구성합니다.

## 연습과 풀이 방향

0부터 9까지의 바이트를 Buffer에 차례로 추가하고 합계를 구하십시오. 각 추가가 실패할 때 해제해야 하며 읽기는 len 범위에서만 수행합니다. 전체 풀이와 경계 실패는 [Buffer 사용법](/docs/ko/stdlib/buffer)의 예제를 따라 확인할 수 있습니다.

자기 코드에서 할당·재할당·해제 호출에 표시를 해 보십시오. 성공한 할당 하나마다 누가 소유하며 어느 경로에서 해제하는지 설명할 수 있어야 합니다.
