---
translation_set_id: whale-ir-reference
path: whale/ir-reference
locale: ko
group: whale
group_order: 5
order: 2
title: Whale IR 참조
summary: 타입, 식별자, 함수의 유효성, 평가 순서와 교환 형식을 설명합니다.
---

## 모듈과 식별자

모듈은 타깃 정보, 전역 정의, 함수로 구성됩니다. 값에는 명시적인 타입이 있습니다. 프런트엔드는 소스 언어의 이름·타입·오버로드·제네릭을 해결한 뒤 typed IR을 생성합니다.

함수와 전역 변수는 서로 다른 내부 이름 공간을 사용합니다. 따라서 함수와 변수가 같은 이름을 가질 수 있습니다. 내부 식별자는 외부 연결 이름인 `link_name`과 구별되며, 외부 이름은 프런트엔드가 명시합니다. Whale은 새로운 이름을 자동 생성해 외부 충돌을 해결하지 않습니다. [심볼과 링크](assembler-linker)를 참고하세요.

각 값 정의는 식별자를 갖습니다. 정의는 중복될 수 없고, 타입 메타데이터는 정의에 명시된 타입과 일치해야 합니다. 이름이 같은 선언이 서로를 가리는 경우에도 이름만으로 정의를 식별하지 않습니다.

## IR 구성과 읽기

다음은 `ir` 크레이트로 함수를 구성하고 검증한 뒤 typed IR을 출력하는 완전한 Rust 예제입니다.

```rust
use ir::{ModuleBuilder, Target, Type};

fn main() {
    let target = Target::lookup("x86_64-whale-linux").unwrap();
    let mut module = ModuleBuilder::new(target.name(), target.data_layout());
    let mut function = module.begin_function("answer", vec![], Type::I32);
    let left = function.const_i32(40);
    let right = function.const_i32(2);
    let answer = function.add(Type::I32, left, right);
    function.ret(Some(answer));
    function.finish();
    let module = module.finish();
    ir::verify_module(&module).unwrap();
    print!("{}", ir::print_module(&module));
}
```

프린터는 다음 IR을 출력합니다.

```text
module {
  format_version 1
  semantics_version 1
  target "x86_64-whale-linux"
  datalayout { ptr=64, endian=little }

  fn @answer() -> i32 {
  entry:
    %v0: i32 = const i32 40
    %v1: i32 = const i32 2
    %v2: i32 = add i32 %v0, %v1
    ret i32 %v2
  }

}
```

`%v0`과 `%v1`은 i32 상수의 정의입니다. `add`가 정의한 `%v2`를 함수의 i32 반환값으로 사용합니다. 반환값을 Bool로 바꾸면 함수 서명에 맞지 않아 검증 오류가 됩니다. 두 피연산자가 상수여도 O0에서는 덧셈 명령을 유지합니다.

이 코드는 실제 프린터 출력이며 텍스트 파서에 전달할 입력 파일이 아닙니다. 텍스트 파싱과 IR 실행은 아직 지원하지 않으며, 현재는 Rust builder로 이 모듈을 구성할 수 있습니다.

## 타입

| 타입 | 의미 |
| --- | --- |
| `bool` | 논리값 false 또는 true |
| `i1`, `i8`, `i16`, `i32`, `i64`, `i128` | 지정된 비트 폭의 부호 있는 정수 |
| `u1`, `u8`, `u16`, `u32`, `u64`, `u128` | 지정된 비트 폭의 부호 없는 정수 |
| `f16`, `f32`, `f64` | 지정된 비트 폭의 부동소수점 값 |
| `ptr<T>` | T 타입 값을 가리키는 포인터 |
| `array<T, N>` | 같은 타입의 원소 N개 |
| `struct{T, ...}` | 순서가 있는 구조체 필드 |
| `tuple<T, ...>` | 순서가 있는 튜플 원소 |
| `void` | 결과 없음 |

`bool`, `i1`, `u1`은 서로 다른 타입입니다. signed `i1`은 −1과 0, unsigned `u1`은 0과 1을 표현합니다. 정수 1은 암묵적으로 논리 조건이 되지 않습니다. 조건 분기, Select의 조건, `trap_if`는 Bool 피연산자를 요구합니다.

저장 크기는 값의 비트 수만으로 정해지지 않고 [타깃 레이아웃](memory-model)을 따릅니다. 예를 들어 `i1`의 값은 1비트지만 메모리에서는 최소 1바이트를 차지합니다.

## 함수와 호출

함수에는 전체 매개변수·결과 타입, 호출 규약, linkage를 명시합니다. 직접 호출과 간접 호출은 호출 대상의 서명과 일치해야 합니다. void 호출은 결과 ID가 없습니다. nonvoid 호출은 O0에서 결과를 사용하지 않더라도 결과 정의를 유지합니다.

반환은 함수의 결과 타입과 일치해야 합니다. void 반환은 값을 전달하지 않으며, nonvoid 반환은 선언된 결과 타입의 값을 전달합니다.

## 블록과 값의 사용 가능성

각 블록에는 고유한 식별자와 정확히 하나의 terminator가 있습니다. 분기 대상은 같은 함수에 속해야 합니다. 진입 블록은 반드시 존재하며 선행 간선과 phi를 가질 수 없습니다. 루프를 만들 때는 진입 블록에서 별도의 루프 헤더로 분기합니다.

실행 가능한 경로에서 값의 정의는 일반적인 사용 지점을 지배해야 합니다. 즉, 진입점에서 사용 지점으로 가는 모든 경로가 그 정의를 지나야 합니다. 같은 블록에서는 정의가 사용보다 앞에 있어야 합니다. 블록의 저장 순서는 지배 관계를 결정하지 않습니다.

진입점에서 left 또는 right로 분기한 뒤 join에서 합류한다고 가정합니다. left에서만 정의한 값은 join에서 일반 값으로 사용할 수 없습니다. right를 지나는 경로에는 정의가 없기 때문입니다. 각 선행 블록의 값을 입력으로 받는 phi로 합쳐야 합니다.

도달 불가능한 블록도 모듈에 보존됩니다. 검증기는 해당 블록의 식별자·타입·피연산자·분기 구조를 계속 검사합니다. 도달 불가능한 블록의 정의가 도달 가능한 경로의 일반 사용에 값을 공급할 수는 없습니다.

## Phi 명령

phi는 블록 안의 모든 일반 명령보다 앞에 놓입니다. 서로 다른 선행 블록마다 정확히 하나의 입력이 필요합니다. 입력 값은 phi와 같은 타입이어야 하며, 해당 선행 블록의 끝에서 사용할 수 있어야 합니다.

하나의 선행 블록에서 간선이 여러 개 들어와도 입력은 하나입니다. 루프 phi는 모듈 저장 순서상 나중에 나오는 블록의 값이라도 반복 간선에서 계산되는 값이면 참조할 수 있습니다. 누락·중복·무관한 선행 블록과 잘못된 타입의 입력은 검증 오류입니다.

## 평가와 선택

Whale AST는 호출 대상과 하위 표현식을 명시된 왼쪽부터 평가합니다. 프런트엔드는 단락 평가를 제어 흐름 분기로 표현합니다.

Select는 이미 계산된 값 중 하나를 선택합니다. 어느 쪽 입력의 계산도 생략하지 않습니다. 예를 들어 안전한 값을 선택하더라도 다른 입력을 계산하면서 발생한 trap을 피할 수 없습니다. 특정 경로에서만 실행해야 하는 계산은 조건부 블록 안에 배치해야 합니다.

## 검증과 trap

잘못된 IR은 검증 단계에서 거부합니다. 실행 조건 위반은 정의된 trap으로 처리하며, `undef`와 `poison`은 허용되는 값이 아닙니다. builder 오용, 중복 정의, 두 번째 terminator 추가는 구조적 오류로 반환해야 합니다.

trap은 이유·소스 위치·IR ID를 포함하며 이후 실행을 중단합니다. native 실행에서는 프로그램을 종료하고 인터프리터 API에서는 Trap 오류를 반환합니다. 이전 부작용은 보존하지만 버퍼 flush·소멸자 호출·스택 unwinding은 보장하지 않습니다.

이 보장은 검증된 IR과 추적 메모리에 적용됩니다. 외부 C·원시 주소·인라인 어셈블리는 별도 계약을 가지며, 그 경계 밖의 위반까지 항상 검출하지는 않습니다. [메모리 모델](memory-model)을 참고하세요.

## 교환 형식과 텍스트 표현

AST와 typed IR은 각각의 format version과 공통 semantics version을 사용합니다. 읽는 쪽은 무버전·알 수 없는 버전·필드·기능·중복 JSON 키를 거부해야 합니다. 생성자는 지원하지 않는 속성이 조용히 무시될 것이라고 가정해서는 안 됩니다.

정수는 비트 폭·signedness·문자열 숫자로 전달합니다. 부동소수점 상수는 폭과 정확한 비트열로 전달합니다. 텍스트 IR의 round-trip은 이름·ID·타입·상수·순서·속성·메타데이터를 보존해야 합니다. 공백과 주석 배치는 보존 대상이 아닙니다.

아래 스칼라 AST JSON 계약을 사용할 수 있습니다. 출력되는 typed IR에는 버전 정보가 포함되지만 텍스트 파서와 완전한 왕복 교환은 아직 미지원입니다.


### 버전이 명시된 AST JSON

다음을 `program.json`으로 저장합니다. Envelope의 네 필드는 모두 필수입니다. `program`의 `globals`와 `functions` 배열도 필수이며 빈 배열을 허용합니다. 함수의 이름·매개변수·반환 타입·본문은 필수입니다. 각 enum은 unit 이름 또는 variant 키 하나를 가진 객체로 표현합니다. Unit variant는 `{"Void":null}`처럼 null 값을 가진 객체도 허용하며, encoder는 unit 이름 `"Void"`로 출력합니다. `VarDecl.init`은 생략하거나 null로 지정할 수 있으며 다른 필수 필드는 생략할 수 없습니다.

```json
{
  "format_version": 1,
  "semantics_version": 1,
  "features": [],
  "program": {
    "globals": [],
    "functions": [
      {
        "name": "answer",
        "parameters": [],
        "return_type": {
          "Int": {
            "bits": 128,
            "signed": false
          }
        },
        "body": [
          {
            "Return": {
              "Lit": {
                "Int": {
                  "bits": 128,
                  "signed": false,
                  "value": "340282366920938463463374607431768211455"
                }
              }
            }
          }
        ]
      }
    ]
  }
}
```

```sh
cargo run --locked --features socket-cli -- ir lower program.json
```

```text
module {
  format_version 1
  semantics_version 1
  target "x86_64-whale-linux"
  datalayout { ptr=64, endian=little }

  fn @answer() -> u128 {
  entry:
    %v0: u128 = const u128 340282366920938463463374607431768211455
    ret u128 %v0
  }

}
```

정수 `value`는 10진 문자열입니다. Signed 정수의 선택적인 마이너스 뒤에 숫자를 쓰며 공백·플러스·지수·구분자는 허용하지 않습니다. 허용 범위는 선언한 폭과 signedness로 결정합니다. 위의 `u128::MAX`는 JSON과 lowering을 거치면서 그대로 보존됩니다. 음수 unsigned 값이나 범위를 벗어난 값은 wrap하지 않고 오류입니다. Float 값에는 [수치 연산](numeric-operations)에서 설명하는 정확한 폭의 16진 비트 문자열을 사용합니다.

이 AST의 `format_version`과 `semantics_version`은 각각 1입니다. `features`는 빈 배열이어야 합니다. 알 수 없는 필드·버전·기능, 중복된 원본 JSON 키(escape를 풀면 같은 키인 경우 포함), 뒤따르는 추가 값은 `--no-verify`에서도 오류입니다. 라이브러리 진입점은 `ir::lower_ast::interchange::decode`이며 `encode`는 envelope를 출력합니다. `decode`의 기본 원본 크기 한도는 8 MiB이고 `decode_with_limit`으로 한도를 지정합니다. JSON 중첩에도 한도가 있습니다. 먼저 일반 map으로 읽으면 중복 키가 사라질 수 있으므로 원본 decoder를 사용합니다.

[전체 JSON Schema](https://github.com/wavefnd/Whale/blob/master/ir/schema/ast-v1.schema.json)는 형태·필수 필드·variant를 정의합니다. 범위·타입 검사와 중복 키 검사가 추가로 적용됩니다. 스칼라 lowering은 리터럴, 변수·상수, add/sub/mul, 비교, 대입, if/while, return과 break/continue를 지원합니다. 호출·복합 값 표현식은 미지원입니다. `Opaque`는 스키마에서 표현할 수 있지만 lowering에서 거부합니다.

기존의 bare Program에는 envelope를 추가하고 JSON 숫자 리터럴을 정수의 10진 문자열 또는 float 비트 문자열로 바꿔야 합니다. 무버전 입력은 거부합니다. AST와 typed IR의 버전 번호는 지금 둘 다 1이어도 독립적으로 관리됩니다.

### 거부되는 입력과 CLI 복구

다음 완전한 입력을 `invalid.json`으로 저장합니다.

```json
{"format_version":99,"semantics_version":1,"features":[],"program":{"globals":[],"functions":[]}}
```

```sh
cargo run --locked --features socket-cli -- ir lower invalid.json -o rejected.wir
```

```text
Failed to parse socket JSON: unsupported AST format_version 99; expected 1
```

명령은 0이 아닌 상태로 종료하며 새 출력을 만들거나 기존 파일을 덮어쓰지 않습니다. 타입 불일치도 출력 게시 전에 실패합니다. `socket-cli` 없이 빌드한 바이너리는 상태 2로 종료하고 `--features socket-cli`가 포함된 복구 명령을 출력합니다.
