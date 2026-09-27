---
translation_set_id: build-link-targets
path: whale/build-link-targets
locale: ko
group: whale
group_order: 1
order: 5
title: 빌드, 링크와 대상 옵션
summary: emit 산출물, 입력 종류, 링크, target/CPU/ABI와 프리스탠딩 빌드 계획을 설명합니다.
---

## emit 산출물

```shell
wavec build main.wave --emit=check
wavec build main.wave --emit=ast
wavec build main.wave --emit=ir
wavec build main.wave --emit=bc
wavec build main.wave --emit=asm
wavec build main.wave --emit=obj -o main.o
wavec build main.wave --emit=bin -o app
```

artifact emit 종류는 `ast`, `ir`, `bc`, `asm`, `obj`, `bin`입니다. `check`는 산출물 종류가 아니라 검사 제어 모드이며 다른 artifact emit과 함께 사용하지 않습니다.

```shell
wavec print supported-emit-kinds
```

## 입력 종류와 link-only

컴파일러는 Wave 소스 외에도 IR, bitcode, assembly, object와 archive 입력을 구분합니다. 지원 목록은 다음 명령으로 질의합니다.

```shell
wavec print supported-input-types
```

이미 만들어진 object나 archive만 링크하려면 `--input-type`과 `--link-only`를 사용할 수 있습니다.

```shell
wavec build module.o --input-type=obj --link-only --emit=bin -o app
```

## 네이티브 링크

```shell
wavec --link=m -L ./lib build main.wave
```

`--link`는 라이브러리를 추가하고 `-L`은 검색 경로를 추가합니다. FFI에서 심볼을 선언했더라도 해당 심볼을 제공하는 라이브러리가 자동으로 링크되는 것은 아닙니다.

## 대상 선택

실행할 OS와 CPU를 선택하는 옵션입니다.

- `--target <triple>`
- `--cpu <name>`
- `--features <csv>`
- `--abi <name>`
- `--sysroot <path>`

호스트 기본값과 지원 대상은 다음 명령으로 확인합니다.

```shell
wavec print host-target
wavec print supported-targets
wavec print target-spec --target <triple>
wavec print cpu-list --target <triple>
wavec print target-features --target <triple>
```

## 지원 대상

현재 컴퓨터에 맞는 프로그램은 target을 지정하지 않고 빌드합니다. 다른 환경을 선택하려면 아래 대상 이름을 `--target`에 전달합니다.

| OS·환경 | 아키텍처 | 대상 이름 |
| --- | --- | --- |
| Linux | amd64 | `x86_64-unknown-linux-gnu` |
| Linux | ARM64 | `aarch64-unknown-linux-gnu` |
| Linux | RISC-V64 | `riscv64-unknown-linux-gnu` |
| Linux | LoongArch64 | `loongarch64-unknown-linux-gnu` |
| macOS | amd64 | `x86_64-apple-darwin` |
| macOS | ARM64 | `aarch64-apple-darwin` |
| Windows | amd64 | `x86_64-pc-windows-msvc` |
| Windows | ARM64 | `aarch64-pc-windows-msvc` |
| FreeBSD | amd64 | `x86_64-unknown-freebsd` |
| WebAssembly | 64비트 | `wasm64-unknown-unknown` |

OS 없이 실행하는 프로그램에는 `x86_64-unknown-none-elf`, `aarch64-unknown-none-elf`, `riscv64-unknown-none-elf`를 사용합니다. 설치된 버전의 전체 대상 목록은 `wavec print supported-targets`로 확인할 수 있습니다.

## RISC-V 64 계약

Hosted RISC-V 대상의 기본값은 `generic-rv64`, RV64GC, `lp64d` ABI입니다. Freestanding 대상의 기본값은 `generic-rv64`, RV64IMAC, `lp64`입니다.

```shell
wavec print target-spec --target riscv64-unknown-linux-gnu --format=json
wavec print target-spec --target riscv64-unknown-none-elf --format=json
```

지원 RISC-V CPU는 `generic`, `generic-rv64`, `rocket-rv64`, `sifive-u74`입니다. Feature override는 `m`, `a`, `f`, `d`, `c`, `zicsr`, `zifencei` 이름 앞에 부호를 붙여 쉼표로 구분합니다.

```shell
wavec build main.wave \
  --target riscv64-unknown-linux-gnu \
  --features=+m,+a,+f,-d,+c,+zicsr \
  --abi=lp64f
```

RISC-V 검증은 일관되지 않은 조합을 거부합니다. `d`에는 `f`가 필요하고 `f`에는 `zicsr`가 필요합니다. `lp64`, `lp64f`, `lp64d`는 활성화한 부동소수점 feature와 일치해야 합니다. ABI를 직접 지정하지 않으면 컴파일러가 feature에서 ABI를 유도합니다.

## 프리스탠딩 링크

```shell
wavec build kernel.wave \
  --target riscv64-unknown-none-elf \
  --freestanding \
  --entry=_start \
  --linker-script=linker.ld \
  --no-start-files \
  -o kernel.elf
```

`--freestanding`은 기본 라이브러리를 사용하지 않는 쪽으로 빌드 설정을 조정합니다. `--entry`는 링커 엔트리를, `--linker-script`는 스크립트를, `--no-start-files`는 호스트 시작 파일 제외를 지정합니다.

실제 실행 전 링크 계획을 확인할 때는 `--dry-run`을 사용할 수 있습니다.

## Hosted 크로스 링크

다른 OS·CPU에서 실행할 프로그램을 만들 때는 대상 환경의 라이브러리 경로를 sysroot로 지정합니다. 별도의 링커를 사용할 경우 `-C linker`로 경로를 지정합니다.

```shell
wavec build main.wave \
  --target riscv64-unknown-linux-gnu \
  --sysroot /path/to/riscv64-sysroot \
  -C linker=/path/to/target-linker \
  -o app-riscv64
```

sysroot와 링크할 라이브러리는 선택한 OS·CPU·ABI에 맞춥니다.

## 크로스 빌드에서 확인할 것

- target triple이 컴파일러의 지원 목록에 있는지
- sysroot와 링커가 대상 ABI에 맞는지
- 링크 라이브러리가 대상 아키텍처용인지
- CPU feature가 대상 CPU에서 유효한지
- 프리스탠딩이면 엔트리 심볼과 메모리 배치가 링커 스크립트와 일치하는지

## WebAssembly 실행

wasm64 결과물은 memory64를 지원하는 실행 환경에서 사용합니다. 파일·시간·입력 같은 외부 기능을 사용하는 모듈은 그 기능에 해당하는 host import를 연결해야 합니다.

대상 코드를 만들고 연결하는 과정은 `--dry-run`으로 확인할 수 있습니다. 실제 실행은 선택한 OS 또는 WebAssembly 실행 환경에서 진행합니다.
