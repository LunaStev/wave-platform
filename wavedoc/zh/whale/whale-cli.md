---
translation_set_id: whale-cli
path: whale/whale-cli
locale: zh
group: whale
group_order: 1
order: 3
title: Whale 命令参考
summary: 描述 Whale assembler、object wrapper、诊断输出和可选 IR 命令。
---

## Whale 构建

在 Whale 存储库中，运行：

```shell
cargo build --release
```

顶级可执行文件有四个命令系列：

```text
whale asm [--amd64 | --aarch64] <input> -o <output>
whale object <input> -o <output>
whale ir <subcommand> [options]
```

## AMD64 assembler

```shell
whale asm --amd64 input.asm -o output.o
```

AMD64 assembler 接收 `.o` 路径作为输出，并且 ELF64 relocatable 包含 section、symbol 和 relocation 创建object。

使用 `--debug-whale` 打开详细诊断输出。

```shell
whale asm --amd64 input.asm -o output.o \
  --debug-whale --token --ast --bytes --dump-hex --stats
```

诊断标志包括`--token`、`--ast`、`--bytes`、`--dump-hex`、`--dump-bin`、`--dump-json`和`--stats`。 `--trace` 打印处理进度。

## Object wrapper

```shell
whale object input.bin -o output.o
```

`object` 命令将原始字节放入 ELF64 `.text` 部分，并在偏移量 0 处添加全局 `start` 符号。它将原始机器代码包装在 ELF 对象文件中。

## 可选 IR socket

仅当使用 `socket-cli` 功能构建 Whale 时，才包含 `ir` 命令。

```shell
cargo run -p whale --features socket-cli -- ir lower program.json
cargo run -p whale --features socket-cli -- ir lower program.json -o program.wir
```

`ir lower`读取Whalesocketschema的JSON，将其转换为WhaleIR，并验证模块。文本IR输出到路径stdout或`-o`。 `--target <triple>` 替换目标字符串，`--no-verify` 省略验证。

要使用 IR 命令，Whale 必须使用 `socket-cli` 功能构建。 Socket JSON 生产者和 Whale 必须使用相同的 socket schema version。
