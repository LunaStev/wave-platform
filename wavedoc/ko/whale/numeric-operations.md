---
translation_set_id: whale-numeric-operations
path: whale/numeric-operations
locale: ko
group: whale
group_order: 5
order: 3
title: 수치 연산
summary: 정수 wrap, checked 연산, 시프트, 형변환 오류와 부동소수점 결과를 설명합니다.
---

## 정수 표현

N비트 정수는 N개의 값 비트를 갖습니다. 부호 없는 정수의 범위는 0부터 2^N − 1까지이고, 부호 있는 정수의 범위는 −2^(N−1)부터 2^(N−1) − 1까지입니다. 연산의 signedness가 비트열의 해석을 결정합니다.

아래 표는 연산 결과를 설명합니다. IR 코드 예제는 현재 프린터의 표현을 사용합니다.

## 덧셈·뺄셈·곱셈

기본 정수 add·sub·mul은 결과의 하위 N비트를 유지합니다. overflow는 trap을 발생시키지 않습니다. checked 연산은 같은 wrap 결과와 함께, 수학적 결과가 해당 signed 또는 unsigned 범위를 벗어났는지를 나타내는 Bool을 반환합니다.

| 연산 | Wrap 결과 | Checked overflow |
| --- | --- | --- |
| u8: 255 + 1 | 0 | true |
| i8: 127 + 1 | −128 | true |
| u8: 0 − 1 | 255 | true |
| i8: 12 × 3 | 36 | false |

overflow에서 실행을 중단하는 언어의 프런트엔드는 checked 연산의 overflow 결과에 명시적인 `trap_if`를 사용해야 합니다. 기본 연산이 소스 언어의 overflow 정책을 암묵적으로 적용하지는 않습니다.

### Wrap과 명시적 검사를 표현한 IR

다음 모듈은 Rust builder로 구성해 검증기를 통과한 현재 프린터 출력입니다. 텍스트 파서와 실행 백엔드는 아직 제공되지 않습니다.

```text
module {
  format_version 1
  semantics_version 1
  target "x86_64-whale-linux"
  datalayout { ptr=64, endian=little }

  fn @add_u8() -> u8 {
  entry:
    %v0: u8 = const u8 255
    %v1: u8 = const u8 1
    %v2: u8 = add u8 %v0, %v1
    ret u8 %v2
  }

  fn @require_no_overflow() -> u8 {
  entry:
    %v3: u8 = const u8 255
    %v4: u8 = const u8 1
    %v5: tuple<u8, bool> = uadd_chk u8 %v3, %v4
    %v6: u8 = extract %v5, 0
    %v7: bool = extract %v5, 1
    trap_if bool %v7, reason="integer overflow"
    ret u8 %v6
  }

}
```

산술 계약에 따르면 `add_u8`은 256을 8비트로 wrap한 0을 반환합니다. `require_no_overflow`는 감긴 결과(`extract ..., 0`)와 Bool overflow 플래그(`extract ..., 1`)를 분리합니다. 플래그가 참이면 명시적인 `trap_if`가 반환 전에 실행을 중단합니다. 이는 정의된 실행 결과이며, 구현된 인터프리터에서 얻은 실행 결과는 아닙니다.

## 나눗셈과 나머지

정수의 0 나눗셈과 0에 대한 나머지는 trap입니다. signed 최솟값을 −1로 나누면 최솟값으로 wrap합니다. 이 경우의 나머지는 0입니다.

| 연산 | 결과 |
| --- | --- |
| i8: −128 / −1 | −128 |
| i8: −128 % −1 | 0 |
| 정수 0 나눗셈 | trap |
| 정수 0에 대한 나머지 | trap |

## 시프트

N비트 값의 시프트에서는 count의 비트열을 unsigned로 해석한 뒤 N으로 나눈 나머지를 사용합니다. count가 0부터 N−1까지의 범위 밖이라는 이유만으로 trap을 발생시키지 않습니다.

8비트 값에서 count 0·8·16은 모두 0비트 시프트입니다. 8비트 count의 비트열 `11111111`은 7비트 시프트가 됩니다. 이 비트열이 signed −1을 나타내더라도 같습니다. unsigned 해석은 나머지 계산보다 먼저 이루어집니다.

음수나 과도한 count를 거부하는 소스 언어는 시프트 앞에 명시적인 검사를 넣어야 합니다.

## 형변환

| 변환 | 의미 |
| --- | --- |
| Zero extension | 상위 비트를 0으로 채워 폭을 늘림 |
| Sign extension | 부호 비트를 복제해 폭을 늘림 |
| 비트 절단 | 목적지 폭에 해당하는 하위 비트만 유지 |
| 비트 재해석 | 같은 비트열을 다른 타입으로 해석 |
| 손실 없는 수치 변환 | 수치를 보존하며 표현할 수 없으면 trap |

예를 들어 8비트 `11111111`을 16비트로 zero extension하면 `0000000011111111`, sign extension하면 `1111111111111111`입니다. 입력 비트가 같아도 서로 다른 연산입니다.

float→int는 0 방향으로 절삭한 뒤 정수 범위를 검사합니다. NaN과 무한대는 trap입니다. i8로 변환할 때 127.9는 127이 되고, 128.0은 trap입니다. Bool은 정수 0 또는 1로 변환합니다. signed i1은 1을 표현할 수 없으므로 이 변환의 목적지로 사용할 수 없습니다.

주소를 정수로 변환해도 그 정수로 포인터 접근 권한을 복구할 수 있는 것은 아닙니다. [포인터 유효성](memory-model)을 참고하세요.

## 부동소수점 연산

부동소수점 값은 정확한 f16·f32·f64 비트열을 갖습니다. 연산은 선언된 폭에서 가장 가까운 값으로 반올림하고, 정확히 중간인 경우 유효숫자의 최하위 비트가 짝수인 값을 선택하는 ties-to-even을 사용합니다.

기본 연산은 fast-math, 암묵적 FMA, 작은 값의 강제 0 처리를 허용하지 않습니다. 곱셈 다음 덧셈은 각각의 반올림 단계를 유지하며 백엔드가 암묵적으로 하나의 연산으로 합쳐서는 안 됩니다.

수치 연산의 결과는 NaN이나 무한대일 수 있습니다. 수치 연산과 폭 변경에서 나오는 NaN은 폭별로 하나의 고정된 양의 quiet NaN으로 정규화합니다. 저장과 복사는 원래 NaN 비트를 보존합니다. 따라서 NaN payload를 산술 연산 없이 메모리로 전달하는 경우와 계산하는 경우의 동작이 다릅니다.

부동소수점 상태 플래그는 노출하지 않습니다. 기본 부동소수점 산술이 NaN·무한대 결과를 허용하더라도 float→int 변환에는 위의 trap 규칙을 적용합니다.


### 정확한 상수 저장

`FloatBits` variant 또는 정확한 폭의 16진 비트 문자열을 사용합니다. 저장 값의 동등성은 음수 0과 NaN payload를 포함한 비트열로 비교합니다. 검증기는 IR 타입과 payload의 폭이 다르면 거부합니다.

```rust
use ir::{FloatBits, ModuleBuilder, Target, Type};
fn main() {
    let bits = FloatBits::parse(32, "0xffc01234").unwrap();
    assert_eq!(bits, FloatBits::F32(0xffc01234));
    let target = Target::X86_64WhaleLinux;
    let mut module = ModuleBuilder::new(target.name(), target.data_layout());
    let mut function = module.begin_function("payload", vec![], Type::F32);
    let value = function.const_float_bits(Type::F32, bits);
    function.ret(Some(value));
    function.finish();
    let module = module.finish();
    ir::verify_module(&module).unwrap();
    assert!(ir::print_module(&module).contains("const f32 0xffc01234"));
    println!("{}", bits);
}
```

```text
0xffc01234
```

f16/f32/f64 문자열은 `0x` 뒤에 정확히 4/8/16개의 16진 숫자를 사용합니다. `0x80000000`은 f32 음수 0이고 `0x7f800000`은 양의 무한대입니다. `const_float`은 host f64에서 수치 변환하는 편의 API이며 원본 비트 보존에는 `const_float_bits`를 사용합니다. 정확한 저장 표현이 float 연산 백엔드의 완성을 의미하지는 않습니다. 컴파일 시점 산술은 여전히 host f64 중간값을 사용하므로 선언된 폭의 전체 반올림 계약은 미완성입니다.
