---
translation_set_id: whale-cli
path: whale/whale-cli
locale: ko
group: whale
group_order: 1
order: 3
title: Whale 명령 참조
summary: Whale assembler, object wrapper, 진단 출력과 선택적 IR 명령을 설명합니다.
---

## Whale 빌드

Whale 저장소에서 다음을 실행합니다.

```shell
cargo build --release
```

최상위 실행 파일에는 네 가지 명령 계열이 있습니다.

```text
whale asm [--amd64 | --aarch64] <input> -o <output>
whale object <input> -o <output>
whale ir <subcommand> [options]
```

## AMD64 assembler

```shell
whale asm --amd64 input.asm -o output.o
```

AMD64 assembler는 `.o` 경로를 출력으로 받고, section, symbol과 relocation을 담은 ELF64 relocatable object를 만듭니다.

상세 진단 출력은 `--debug-whale`로 켭니다.

```shell
whale asm --amd64 input.asm -o output.o \
  --debug-whale --token --ast --bytes --dump-hex --stats
```

진단 플래그에는 `--token`, `--ast`, `--bytes`, `--dump-hex`, `--dump-bin`, `--dump-json`, `--stats`가 있습니다. `--trace`는 처리 과정을 출력합니다.

## Object wrapper

```shell
whale object input.bin -o output.o
```

`object` 명령은 raw byte를 ELF64 `.text` section에 넣고 offset 0에 전역 `start` symbol을 추가합니다. 이 명령은 raw code를 ELF object로 감싸는 용도입니다.

## 텍스트 IR 검증과 출력

기본 빌드에서 형식 3 typed IR을 읽고 검증할 수 있습니다. [IR 참조](ir-reference)의 전체 예제를 `answer.wir`로 저장하세요. `print`는 검증 후 표준 IR을 출력하며, 오류 시 기존 파일을 보존합니다. 실행이나 native 코드 생성은 수행하지 않습니다.

```shell
whale ir verify answer.wir
whale ir print answer.wir -o canonical.wir
```

## 선택적 IR socket

AST JSON을 읽는 `ir lower`에는 `socket-cli` feature가 필요합니다. 텍스트 IR의 `verify`·`print`에는 필요하지 않습니다.

```shell
cargo run -p whale --features socket-cli -- ir lower program.json
cargo run -p whale --features socket-cli -- ir lower program.json -o program.wir
```

`ir lower`는 Whale socket schema의 JSON을 읽고 Whale IR로 변환한 뒤 모듈을 검증합니다. 텍스트 IR은 stdout 또는 `-o` 경로에 출력됩니다. `--target <triple>`은 대상 문자열을 바꾸고 `--no-verify`는 검증을 생략합니다.

`ir lower`를 사용하려면 `socket-cli`로 빌드해야 합니다. Socket JSON 생산자와 Whale은 같은 AST schema version을 사용해야 합니다.
