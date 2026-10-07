---
translation_set_id: whale-cli
path: whale/whale-cli
locale: en
group: whale
group_order: 1
order: 3
title: Whale command reference
summary: Describes the Whale assembler, object wrapper, diagnostic output, and the optional IR commands.
---

## Whale Build

In the Whale repository, run:

```shell
cargo build --release
```

A top-level executable file has four families of commands:

```text
whale asm [--amd64 | --aarch64] <input> -o <output>
whale object <input> -o <output>
whale ir <subcommand> [options]
```

## AMD64 assembler

```shell
whale asm --amd64 input.asm -o output.o
```

AMD64 assembler receives the `.o` path as output, and ELF64 relocatable containing section, symbol and relocation Create object.

Turn on detailed diagnostic output with `--debug-whale`.

```shell
whale asm --amd64 input.asm -o output.o \
  --debug-whale --token --ast --bytes --dump-hex --stats
```

Diagnostic flags include `--token`, `--ast`, `--bytes`, `--dump-hex`, `--dump-bin`, `--dump-json`, and `--stats`. `--trace` prints the processing progress.

## Object wrapper

```shell
whale object input.bin -o output.o
```

The `object` command places raw bytes in an ELF64 `.text` section and adds a global `start` symbol at offset 0. It wraps raw machine code in an ELF object file.

## Typed text IR verification and printing

The default build reads and verifies format 3 typed IR. Save the complete example in the [IR reference](ir-reference) as `answer.wir`. `print` verifies before canonical printing and preserves existing output on failure. It does not execute IR or generate native code.

```shell
whale ir verify answer.wir
whale ir print answer.wir -o canonical.wir
```

## Optional IR socket

AST JSON `ir lower` requires the `socket-cli` feature. Text IR `verify` and `print` do not.

```shell
cargo run -p whale --features socket-cli -- ir lower program.json
cargo run -p whale --features socket-cli -- ir lower program.json -o program.wir
```

`ir lower` reads JSON of Whale socket schema, converts it to Whale IR, and verifies the module. The text IR is output to the path stdout or `-o`. `--target <triple>` replaces the target string and `--no-verify` omits validation.

Build with `socket-cli` to use `ir lower`. Socket JSON producers and Whale must use the same AST schema version.
