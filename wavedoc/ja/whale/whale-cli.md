---
translation_set_id: whale-cli
path: whale/whale-cli
locale: ja
group: whale
group_order: 1
order: 3
title: Whale コマンドリファレンス
summary: Whaleassembler、objectwrapper、診断出力とオプションのIRコマンドについて説明します。
---

## Whaleビルド

Whaleリポジトリで以下を実行します。

```shell
cargo build --release
```

最上位の実行可能ファイルには4つのコマンドシリーズがあります。

```text
whale asm [--amd64 | --aarch64] <input> -o <output>
whale object <input> -o <output>
whale ir <subcommand> [options]
```

## AMD64 assembler

```shell
whale asm --amd64 input.asm -o output.o
```

AMD64 assembler は `.o` パスを出力として受け取り、section、symbol、および relocation を含む ELF64 relocatable オブジェクトを作成します。

詳細診断出力は`--debug-whale`でオンにします。

```shell
whale asm --amd64 input.asm -o output.o \
  --debug-whale --token --ast --bytes --dump-hex --stats
```

診断フラグには`--token`、`--ast`、`--bytes`、`--dump-hex`、`--dump-bin`、`--dump-json`、`--stats`があります。 `--trace`は処理プロセスを出力します。

## Object wrapper

```shell
whale object input.bin -o output.o
```

`object` コマンドは、生のバイトを ELF64 `.text` セクションに配置し、グローバル `start` シンボルをオフセット 0 に追加します。生のマシン コードを ELF オブジェクト ファイルにラップします。

## テキスト IR の検証と出力

標準ビルドで format 3 typed IR を読み取り・検証できます。[IR リファレンス](ir-reference)の完全な例を `answer.wir` に保存してください。`print` は検証後に標準形式で出力し、失敗時は既存ファイルを保護します。IR 実行や native コード生成は行いません。

```shell
whale ir verify answer.wir
whale ir print answer.wir -o canonical.wir
```

## オプション IR socket

AST JSON の `ir lower` には `socket-cli` feature が必要です。テキスト IR の `verify` と `print` には不要です。

```shell
cargo run -p whale --features socket-cli -- ir lower program.json
cargo run -p whale --features socket-cli -- ir lower program.json -o program.wir
```

`ir lower`はWhalesocketschemaのJSONを読み、WhaleIRテキストIRはstdoutまたは`-o`パスに出力されます。 `--target <triple>`はターゲット文字列を置き換え、`--no-verify`は検証を省略します。

`ir lower` を使うには `socket-cli` でビルドしてください。Socket JSON の生成側と Whale は同じ AST schema version を使う必要があります。
