---
translation_set_id: memory-buffer
path: reference/memory-and-buffer
locale: ko
group: stdlib
group_order: 1
order: 4
title: mem: 할당·재할당·레이아웃
summary: 바이트 단위 크기, 할당 실패, 재할당 경계와 해제 책임을 설명합니다.
---

## 할당과 해제

```text
std::mem::alloc
mem_alloc(size: i64) -> ptr<u8>
mem_alloc_zeroed(size: i64) -> ptr<u8>
mem_free(p: ptr<u8>, size: i64) -> i64
mem_realloc(old_ptr: ptr<u8>, old_size: i64, new_size: i64) -> ptr<u8>
```

크기는 바이트 단위입니다. 0 이하 크기의 할당은 null을 반환합니다. 양수 크기도 할당에 실패하면 null일 수 있습니다. `mem_alloc`의 초기 내용을 가정하지 말고, 0 초기화가 필요하면 `mem_alloc_zeroed`를 사용합니다.

성공한 할당은 호출자가 소유합니다. 해제에는 원래 할당 크기를 전달합니다. `mem_free(null, size)`는 0을 반환하고, null이 아닌 포인터와 0 이하 크기의 조합은 오류입니다. 이미 해제한 포인터를 다시 해제하거나 접근하지 않습니다.

## 재할당의 경우 구분

| 요청 | 동작 |
| --- | --- |
| 새 크기가 양수이고 성공 | `min(old_size, new_size)` 바이트를 복사하고 이전 할당 해제 |
| 양수 크기의 새 할당 실패 | null 반환, 기존 할당 유지 |
| `old_ptr == null`, 양수 새 크기 | 새 할당처럼 동작 |
| `new_size == 0` | 유효한 이전 할당의 해제를 시도하고 null 반환 |
| 음수 크기 또는 기존 포인터에 old_size=0 | null 반환 |

0 크기로 재할당한 결과만으로 해제 성공 여부를 확인할 수 없습니다. 해제 결과를 확인해야 한다면 `mem_free`를 직접 호출하십시오. 재할당 후 새로 늘어난 영역은 직접 초기화합니다.

## 기존 포인터를 보존하는 예제

아래는 새 크기가 양수인 경우입니다. `main.wave`로 저장해 실행합니다.

<!-- wave-example: reallocation -->
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
    var grown: ptr<u8> = mem_realloc(data, 4, 8);
    if (grown == null) {
        mem_free(data, 4);
        return 2;
    }
    data = grown;
    println("{}", deref data[0]);
    if (mem_free(data, 8) < 0) {
        return 3;
    }

    return 0;
}
```

실행 결과:

```text
7
```

실패 확인 전에 `data = mem_realloc(...)`로 덮어쓰면 기존 주소를 잃을 수 있습니다. 재할당이 성공하면 이전 주소와 그 내부를 가리키던 포인터는 사용하지 않습니다.

## 대상 타입의 크기와 정렬

```text
std::mem::layout
size_of<T>() -> u64
align_of<T>() -> u64
```

두 값은 실행 중인 컴퓨터가 아니라 컴파일 대상의 레이아웃입니다. `size_of`에는 꼬리 패딩이 포함되며, 값을 생성하거나 평가하지 않습니다. 요소 개수에 크기를 곱할 때는 오버플로를 검사합니다. `std::mem::ops`의 `mem_size_mul_checked`와 `mem_size_add_checked`를 사용할 수 있습니다.

`mem_copy`는 겹치지 않는 범위, `mem_move`는 겹칠 수 있는 범위 복사에 사용합니다. 어느 쪽도 포인터만으로 실제 할당 길이를 알아내지 못하므로 호출자가 범위를 보장해야 합니다. 크기가 변하는 바이트 목록은 [Buffer](/docs/ko/stdlib/buffer)로 관리할 수 있습니다.

## 레이아웃 조회 예제

main.wave로 저장해 실행합니다. 문서에서 다루는 타깃의 i32 크기와 정렬은 각각 4바이트이므로 `4 4`를 출력합니다. 다른 타입, 특히 구조체와 포인터의 값은 대상별로 확인하십시오.

<!-- wave-example: layout-api -->
```wave
import("std::mem::layout")::{
    size_of, align_of
};

fun main() {
    println("{} {}", size_of<i32>(), align_of<i32>());
}
```

## 크기 계산에서 오버플로 확인하기

할당 함수에 넘기기 전에 `count * element_size`가 유효한지 확인해야 합니다. 넘친 값으로 작은 공간을 할당하고 원래 개수만큼 쓰면 경계를 벗어나게 됩니다.

<!-- wave-example: book-memory-size-check -->
```wave
import("std::mem::ops")::{mem_size_add_checked, mem_size_mul_checked};

fun main() {
    var result: i64 = 99;

    if (mem_size_mul_checked(3, 4, &result) == 0) {
        println("bytes={}", result);
    }

    if (mem_size_add_checked(9223372036854775807, 1, &result) < 0) {
        println("overflow rejected");
    }
}
```

실행 결과:

```text
bytes=12
overflow rejected
```

실패한 결과를 할당 크기로 쓰지 않습니다. 일반 산술로 계산한 뒤 결과가 음수인지 보는 것만으로 모든 오버플로를 검출할 수도 없습니다. 검사가 필요한 크기 계산에는 처음부터 checked 함수를 사용합니다.

## 겹치는 복사에는 mem_move

같은 배열의 일부를 뒤로 옮길 때는 입력과 출력 영역이 겹칩니다. mem_copy에 겹치는 범위를 전달하지 말고 mem_move를 사용합니다.

<!-- wave-example: book-memory-overlap -->
```wave
import("std::mem::ops")::{mem_move};

fun main() {
    var data: array<u8, 5> = [1, 2, 3, 4, 5];

    mem_move(&data[1], &data[0], 4);

    for (var index: i32 = 0; index < 5; index += 1) {
        println("{}", data[index]);
    }
}
```

실행 결과:

```text
1
1
2
3
4
```

원래 앞의 네 바이트가 한 칸 뒤로 이동합니다. 직접 앞에서부터 대입하면 이미 덮어쓴 값을 다시 읽어 전부 1이 되는 실수를 할 수 있습니다. 겹침을 처리하는 API를 선택하면 방향 선택을 직접 구현할 필요가 없습니다.

## 소유권을 전달하는 함수 작성

메모리를 반환하는 함수는 성공 시 반환 주소와 해제에 필요한 크기를 함께 제공하는 편이 좋습니다. 호출자가 그 크기를 추측해야 하면 잘못된 해제가 생기기 쉽습니다. 빌린 주소를 반환하는 함수라면 호출자가 해제하면 안 된다는 점과 원본 수명을 설명합니다.

함수 경계를 넘어서도 `할당자 → 소유자 → 해제 지점`을 연결해 읽을 수 있어야 합니다. 포인터 변수의 이름이나 타입만으로 소유권이 자동으로 결정되는 것은 아닙니다.
