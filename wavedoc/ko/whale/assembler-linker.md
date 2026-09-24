---
translation_set_id: whale-assembler-linker
path: whale/assembler-linker
locale: ko
group: whale
group_order: 5
order: 7
title: 어셈블러와 정적 링크
summary: 피연산자 인코딩, 섹션 배치, 심볼 바인딩과 실행 진입점을 설명합니다.
---

## 어셈블리와 오브젝트

Whale 어셈블러는 AMD64 명령을 기계어 바이트와 재배치 정보로 변환합니다. ELF64 오브젝트는 이 정보와 섹션·심볼을 담습니다. 어셈블은 Whale 자체 구현으로 수행하며 외부 어셈블러가 필요하지 않습니다.

재배치 가능한 오브젝트에는 최종 주소를 아직 모르는 참조가 있을 수 있습니다. 이 주소를 해결하는 과정이 링크입니다. 어셈블이 성공했다고 모든 외부 심볼을 해결했거나 실행 파일을 만들었다고 판단해서는 안 됩니다.

## 함수 어셈블하기

다음 코드를 `answer.asm`으로 저장합니다.

```asm
section .text
global answer

answer:
    mov eax, 42
    ret
```

```sh
whale asm --amd64 answer.asm -o answer.o
```

`.text` 내용은 `b8 2a 00 00 00 c3`이며, `mov eax, 42` 다음에 `ret`이 옵니다. ELF64 오브젝트는 `answer`를 공개합니다. 프로세스 시작 코드가 없는 호출 가능한 함수이며 실행 파일은 아닙니다. 이 예제의 명령은 현재 어셈블러에서 처리할 수 있습니다.

## Rust API로 오브젝트 구성하기

다음은 `object` 크레이트로 같은 함수 바이트를 기록하는 완전한 예제입니다. 출력 타깃과 전역 심볼을 명시합니다.

```rust
use object::{ObjectFile, ObjectSymbol, SectionKind, SymbolBinding, SymbolVisibility, Target};

fn main() {
    let mut object = ObjectFile::with_target(Target::X86_64WhaleLinux.object_target());
    let text = object.add_section(".text", SectionKind::Text, 16);
    object.sections[text].data = vec![0xb8, 0x2a, 0x00, 0x00, 0x00, 0xc3];
    object.symbols.push(ObjectSymbol {
        name: "answer".into(),
        section_index: Some(text),
        value: 0,
        size: 6,
        binding: SymbolBinding::Global,
        visibility: SymbolVisibility::Default,
    });
    let elf = object.write().unwrap();
    assert_eq!(&elf[..7], b"\x7fELF\x02\x01\x01");
    assert_eq!(u16::from_le_bytes([elf[18], elf[19]]), 62);
    std::fs::write("answer.o", elf).unwrap();
}
```

Assertion은 ELF 클래스, 바이트 순서, machine 식별자를 확인합니다. `value: 0`은 `.text` 내부 offset이며 `size: 6`은 심볼의 바이트 크기입니다. 잘못된 섹션 참조나 범위는 직렬화 오류입니다. 다른 machine이나 바이트 순서도 AMD64로 표시하지 않고 거부합니다.

## 두 오브젝트의 심볼 해석하기

현재 `linker` 크레이트는 심볼 해석을 제공합니다. 다음 실행 가능한 예제는 두 오브젝트에 각각 로컬 `helper`를 정의한 뒤, 같은 이름을 두 번 공개하면 오류가 나는지 확인합니다.

```rust
use linker::core::symbol_table::{SymbolKey, SymbolTable};
use object::{
    ObjectFile, ObjectFormat, ObjectSymbol, SectionKind, SymbolBinding, SymbolVisibility,
};

fn input(binding: SymbolBinding) -> ObjectFile {
    let mut object = ObjectFile::new(ObjectFormat::ELF64);
    let text = object.add_section(".text", SectionKind::Text, 1);
    object.sections[text].data = vec![0xc3];
    object.symbols.push(ObjectSymbol {
        name: "helper".into(),
        section_index: Some(text),
        value: 0,
        size: 1,
        binding,
        visibility: SymbolVisibility::Default,
    });
    object
}

fn main() {
    let mut symbols = SymbolTable::new();
    symbols
        .resolve(&[input(SymbolBinding::Local), input(SymbolBinding::Local)])
        .unwrap();
    for object_index in 0..2 {
        let key = SymbolKey::Local {
            object_index,
            name: "helper".into(),
        };
        assert_eq!(symbols.symbols[&key].object_index, Some(object_index));
    }
    let error = symbols
        .resolve(&[input(SymbolBinding::Global), input(SymbolBinding::Global)])
        .unwrap_err();
    assert_eq!(error, "Duplicate global symbol: helper");
    println!("{error}");
}
```

출력:

```text
Duplicate global symbol: helper
```

두 로컬 정의는 `object_index`가 다른 별도 키를 갖습니다. 두 전역 정의는 충돌합니다. 이 예제는 메모리에 구성한 오브젝트의 심볼을 직접 해석합니다. `.o` 파일 읽기, 재배치 적용, 실행 파일 출력은 수행하지 않습니다. 아래에서 설명하는 전체 정의 선택 정책 중 weak 우선순위 등은 아직 구현되지 않았습니다.

## 리터럴과 메모리 피연산자

리터럴은 실제 명령의 인코딩 범위를 검사할 때까지 폭과 부호를 보존합니다. 범위에 들어가지 않는 값은 조용히 잘리지 않고 오류가 되어야 합니다.

명령의 다른 정보로 메모리 폭을 결정할 수 없을 때 명시적인 크기 표기가 필요합니다. 예를 들어 레지스터 피연산자가 폭을 정할 수 있지만, 메모리와 즉시값만 있는 경우에는 모호할 수 있습니다. 어셈블러가 모호한 폭을 임의로 추측해서는 안 됩니다.

AMD64에서 심볼 하나로 이루어진 메모리 피연산자는 기본적으로 RIP-relative입니다. 명시적 rel/abs로 주소 지정 방식을 선택합니다. 리터럴 안의 알 수 없는 escape는 오류입니다.

## 섹션과 정렬

| 섹션 내용 | 정렬 동작 |
| --- | --- |
| 코드 | NOP 명령 삽입 |
| 초기화된 데이터 | 0 바이트 삽입 |
| BSS | 파일 payload를 추가하지 않고 논리 메모리 크기를 증가 |

파일 크기와 메모리 크기는 구별됩니다. BSS는 메모리를 예약하지만 같은 크기의 0 바이트를 오브젝트 파일에 저장할 필요는 없습니다. 사용자 정의 섹션은 속성을 가지며 심볼은 바인딩과 타입 정보를 유지합니다.

## 심볼 식별

함수와 변수는 IR 내부에서 서로 다른 식별자를 사용합니다. 외부 연결은 프런트엔드가 명시한 `link_name`을 사용합니다. Whale은 충돌하는 공개 심볼 중 하나의 이름을 자동으로 바꾸지 않고 명시된 이름을 유지합니다.

따라서 내부 함수와 변수의 이름이 모두 `item`인 것은 가능하지만, 둘을 같은 외부 이름으로 공개하면 오류가 될 수 있습니다. 내부 이름 공간이 분리되어 있다고 외부 이름 공간도 자동으로 분리되지는 않습니다.

오브젝트 로컬 심볼의 범위는 해당 입력 오브젝트입니다. 전역 심볼은 오브젝트 사이의 해석에 참여합니다. 확인된 함수·데이터 충돌은 오류입니다. NOTYPE 심볼은 더 구체적인 타입을 제공하지 않는 입력과의 호환을 유지합니다. 타입이 없다는 사실만으로 함수나 데이터라고 단정하지 않습니다.

## 정의 선택

| 정의 또는 참조 | 결과 |
| --- | --- |
| Strong과 strong | 중복 정의 오류 |
| Strong과 weak | Strong 정의 선택 |
| Weak와 weak | 입력 순서상 첫 정의 선택 |
| 미해결 strong 참조 | 링크 오류 |
| 미해결 weak 참조 | 초기 정적 프로파일에서는 미지원 오류 |

여러 weak 정의가 있으면 입력 순서가 결과에 영향을 줍니다. 결정적인 링크를 위해 전달된 입력 순서를 일관되게 사용해야 합니다.

## 정적 실행 파일 출력

정적 native 프로파일은 진입점을 명시한 ELF ET_EXEC를 생성합니다. `main`이라는 함수 이름에서 진입점을 추론하지 않습니다. 해당 함수를 호출하는 시작 코드도 자동 삽입하지 않습니다.

섹션 제거, 동일 코드 합치기, 심볼 제거를 자동으로 수행하지 않습니다. 파일 배치는 실제로 저장하는 바이트와 실행 시 예약하는 메모리를 별도로 계산해야 합니다.

CLI의 완전한 정적 실행 파일 생성 경로는 아직 제공되지 않습니다. `whale asm`은 재배치 가능한 오브젝트를 만들고, `whale object`는 원시 바이트를 오브젝트로 감쌉니다. 제공 여부는 [툴체인 개요](overview), ABI 요구 사항은 [AMD64 타깃](amd64-target)을 참고하세요.
