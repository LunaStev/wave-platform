---
translation_set_id: whale-memory-model
path: whale/memory-model
locale: ko
group: whale
group_order: 5
order: 4
title: 메모리 모델
summary: 할당 추적, 초기화된 값의 읽기, 포인터 산술, 레이아웃과 문자열 저장 규칙입니다.
---

## 할당 추적

추적 포인터는 주소에 할당 ID·세대·범위·offset·접근 권한을 연결합니다. 메모리 모델은 할당의 수명과 초기화 상태도 관리합니다. 접근은 이 조건들을 만족해야 하며 위반하면 trap입니다.

같은 물리 주소가 재사용되더라도 세대는 각 수명을 구분합니다. 주소가 있다는 사실만으로 포인터가 유효하거나 호출자에게 해당 저장 공간의 접근 권한이 있다고 판단하지 않습니다.

native 주소는 64비트를 유지합니다. 별도의 shadow metadata가 복사·저장·호출·반환을 통해 포인터와 함께 전달됩니다. 초기 native 메모리 범위는 추적 가능한 스택과 전역 할당입니다. C 경계에는 명시적 어댑터가 필요하며 임의 외부 메모리의 소유권 이전은 이 범위에 포함하지 않습니다.

## 초기화와 읽기

저장 공간을 선언했다고 값이 초기화되는 것은 아닙니다. 실제로 읽는 바이트 범위의 초기화 여부를 검사합니다. 값의 일부라도 미초기화 상태로 읽으면 trap입니다. 읽기 결과를 0이나 불특정 값으로 대체하지 않습니다.

값 읽기의 초기화 검사는 padding 바이트를 제외합니다. 예를 들어 구조체의 모든 필드가 초기화되었다면, 필드 정렬 때문에 생긴 빈 바이트가 미초기화 상태라는 이유만으로 값 읽기가 잘못되지는 않습니다.

메모리 복사는 바이트와 함께 초기화 상태를 전달합니다. 미초기화 저장 공간을 복사해도 초기화된 저장 공간으로 바뀌지 않습니다. 나중에 목적지에서 값을 읽을 때도 원본을 읽을 때와 같은 검사를 적용합니다.

BSS의 물리적인 바이트가 0이라는 사실만으로 IR 변수의 초기화가 성립하지 않습니다.

### 초기화된 스칼라 저장 공간의 IR

다음 모듈은 builder로 구성해 검증기를 통과했습니다. 값을 읽기 전에 `store`하며, 세 메모리 명령 모두 0이 아닌 2의 거듭제곱 정렬을 명시합니다.

```text
module {
  target "x86_64-whale-linux"
  datalayout { ptr=64, endian=little }

  fn @initialized_local() -> i32 {
  entry:
    %v0: ptr<i32> = alloca i32, align 4
    %v1: i32 = const i32 42
    store i32 %v1, ptr<i32> %v0, align 4
    %v2: i32 = load i32, ptr<i32> %v0, align 4
    ret i32 %v2
  }

}
```

`alloca`는 i32 저장 공간을 만들고, `store`는 42를 쓰며, `load`는 반환할 값을 정의합니다. 이 코드는 typed IR 프린터 출력이며 native 실행 기록이 아닙니다. 정렬을 3으로 바꾸면 검증 오류입니다. Store를 제거하면 이 메모리 모델에서는 미초기화 읽기 trap이 필요하지만, 초기화 추적과 해당 런타임 trap 검사는 아직 제공되지 않습니다. 현재 검증기를 통과했다는 사실만으로 그러한 읽기가 안전하다고 판단해서는 안 됩니다.

## 포인터 산술과 비교

주소 계산 overflow는 trap입니다. 할당 범위의 바로 다음을 가리키는 one-past 포인터는 만들 수 있지만, 이를 통해 메모리에 접근할 수는 없습니다.

포인터 동등성에는 할당 정체성을 사용합니다. 숫자 주소만으로 정체성이 성립하지는 않습니다. 서로 다른 할당의 포인터 사이에서 순서를 비교하거나 차이를 구하면 trap입니다. 정수에서 주소를 재구성해도 접근 권한은 복구되지 않습니다.

### GEP

GEP는 요소와 필드 단위로 주소를 계산합니다. 값을 읽는 명령이 아닙니다.

첫 인덱스는 기준 포인터가 가리키는 타입의 요소 단위 offset입니다. 이후 인덱스는 배열 요소 또는 구조체·튜플 필드를 선택합니다. 구조체·튜플의 필드 인덱스는 컴파일 시점의 필드 순번이며 바이트 offset이 아닙니다. 선택된 타입이 결과 포인터 타입을 결정합니다.

기준 타입이 `ptr<array<i32, 4>>`일 때 인덱스 `[0, 2]`는 배열의 세 번째 i32 요소를 선택하며 결과는 `ptr<i32>`입니다. 첫 인덱스를 1로 지정하면 i32 요소 하나가 아니라 4개 요소를 가진 배열 하나만큼 이동합니다. 이 예제는 인덱스 의미를 설명하며 텍스트 명령 문법은 아닙니다.

초기 native 타깃은 크기 0 요소에 대한 포인터 산술을 거부합니다. 주소를 계산해도 이후 접근에 필요한 수명·범위·초기화·권한 검사는 없어지지 않습니다.

## 데이터 레이아웃

출력 타깃이 크기·정렬·필드 offset·배열 stride를 결정합니다. 구조체와 튜플의 필드 순서는 선언 순서를 유지합니다. 컴파일러를 실행하는 호스트의 레이아웃을 출력 타깃의 레이아웃으로 가정해서는 안 됩니다.

| 값 | 저장 규칙 |
| --- | --- |
| Bool, signed i1, unsigned u1 | 최소 1바이트 |
| 빈 구조체·튜플 | 크기 0, 정렬 1 |
| 배열 | 타깃 레이아웃이 요소 stride를 결정 |
| 구조체·튜플 | 선언 순서와 타깃이 요구하는 정렬을 유지 |

완성된 IR의 정렬은 0이 아닌 2의 거듭제곱입니다. 자동 정렬은 이 IR을 생성하기 전에 결정해야 합니다. 이 프로파일에서 packed 레이아웃·union·bitfield는 지원하지 않으며 거부해야 합니다.

### 출력 레이아웃 조회

Rust API는 빌드 호스트와 무관하게 저장 레이아웃을 계산합니다. 다음 예제에는 u64 필드 앞의 padding 7바이트와 마지막 padding 6바이트가 있습니다.

```rust
use ir::{allocation_align, layout_of, Target, Type};

fn main() {
    let target = Target::X86_64WhaleLinux;
    let record = Type::Struct(vec![Type::U8, Type::U64, Type::U16]);
    let layout = layout_of(&record, target).unwrap();
    assert_eq!((layout.size, layout.align), (24, 8));
    assert_eq!(layout.field_offsets, [0, 8, 16]);

    let array = Type::Array(Box::new(record), 3);
    let layout = layout_of(&array, target).unwrap();
    assert_eq!((layout.size, layout.align), (72, 8));
    assert_eq!(layout.element_stride, Some(24));
    assert_eq!(allocation_align(&array, target).unwrap(), 16);
}
```

```text
struct{u8, u64, u16}: size 24, natural alignment 8
field 0: byte 0
field 1: byte 8
field 2: byte 16
array of 3: size 72, element stride 24
standalone array placement alignment: 16
```

`layout_of`의 크기와 배열 stride에는 마지막 padding이 포함됩니다. 구조체와 튜플은 같은 필드 순서 규칙을 사용합니다. Bool·i1·u1은 각각 1바이트를 차지합니다. 빈 구조체와 튜플은 크기 0·정렬 1이며, 길이 0인 배열은 요소의 자연 정렬을 유지합니다. `void`에는 저장 레이아웃이 없지만 `ptr<void>`는 8바이트입니다.

필드와 배열 요소에는 자연 정렬을 사용합니다. `allocation_align`은 16바이트 이상인 독립된 지역·전역 배열에 최소 16바이트 정렬을 요구하는 SysV AMD64 규칙을 적용합니다. 배열 필드의 정렬이나 요소 stride를 늘리지는 않습니다. AST lowering은 지역 저장 공간에 이 배치 정렬 조회를 사용합니다.

크기 곱셈, 필드 offset 덧셈, padding 계산의 overflow는 `LayoutError::Overflow`로 반환합니다. `layout_of`의 복합 타입 중첩 한도는 128단계이며 `layout_of_with_limit`으로 호출자가 한도를 지정할 수 있습니다. 배열 요소 수만큼 저장 공간을 할당하지 않습니다. `pointer_stride`는 native 포인터 산술에서 크기 0인 pointee를 거부하지만 해당 타입의 저장 레이아웃 자체는 유효합니다. Packed·union·bitfield는 지원되는 타입 표현이 없습니다. 이 저장 레이아웃 조회는 복합 타입 호출 규약이나 런타임 범위 검사를 구현하지 않습니다.

## 문자열과 C 경계

문자열은 길이가 명시된 불변 바이트열입니다. 기본 인코딩은 UTF-8입니다. 내부 NUL을 허용하며 자동으로 끝 NUL을 붙이지 않습니다. O0는 내용이 같은 문자열 객체를 자동으로 합치지 않습니다.

따라서 `A`, NUL, `B`로 구성된 바이트열의 길이는 3입니다. 내부 NUL을 거부하는 명시적 C 문자열 변환에는 사용할 수 없습니다. 프런트엔드가 이를 조용히 `A`로 잘라서는 안 됩니다.

외부 C·원시 주소·인라인 어셈블리는 별도 계약 경계입니다. 추적 메모리의 런타임 검사가 외부 코드의 모든 잘못된 동작까지 검출한다고 보장하지 않습니다.
