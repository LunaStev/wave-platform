---
translation_set_id: build-link-targets
path: whale/build-link-targets
locale: ja
group: whale
group_order: 1
order: 5
title: ビルド、リンク、ターゲットの指定
summary: emit出力、入力タイプ、リンク、target/CPU/ABIとプレスタンディングビルドプランについて説明します。
---

## emit 出力

```shell
wavec build main.wave --emit=check
wavec build main.wave --emit=ast
wavec build main.wave --emit=ir
wavec build main.wave --emit=bc
wavec build main.wave --emit=asm
wavec build main.wave --emit=obj -o main.o
wavec build main.wave --emit=bin -o app
```

artifact emit の種類は `ast`、`ir`、`bc`、`asm`、`obj`、`bin` です。 `check`は出力の種類ではなく検査制御モードであり、他のアーティファクト emitと一緒に使用しません。

```shell
wavec print supported-emit-kinds
```

## 入力種類とlink-only

コンパイラは、Waveソースに加えて、IR、bitcode、assembly、objectとarchive入力を区別します。サポートリストは次のコマンドで問い合わせます。

```shell
wavec print supported-input-types
```

すでに作成されているobjectまたはarchiveのみリンクするには、`--input-type`と`--link-only`を使用できます。

```shell
wavec build module.o --input-type=obj --link-only --emit=bin -o app
```

## ネイティブリンク

```shell
wavec --link=m -L ./lib build main.wave
```

`--link`はライブラリを追加し、`-L`は検索パスを追加します。 FFIでシンボルを宣言しても、そのシンボルを提供するライブラリが自動的にリンクされるわけではありません。

## ターゲットを選択

実行するOSとCPUを選択するオプションです。

- `--target <triple>`
- `--cpu <name>`
- `--features <csv>`
- `--abi <name>`
- `--sysroot <path>`

ホストのデフォルトとサポート対象は、次のコマンドで確認します。

```shell
wavec print host-target
wavec print supported-targets
wavec print target-spec --target <triple>
wavec print cpu-list --target <triple>
wavec print target-features --target <triple>
```

## サポート対象

現在のコンピュータに合ったプログラムは、targetを指定せずにビルドします。別の環境を選択するには、以下のターゲット名を`--target`に渡します。

|OS・環境|アーキテクチャ|ターゲット名|
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
| WebAssembly |64ビット| `wasm64-unknown-unknown` |

OSなしで実行するプログラムには`x86_64-unknown-none-elf`、`aarch64-unknown-none-elf`、`riscv64-unknown-none-elf`を使用します。インストールされているバージョンの完全なターゲットのリストは、`wavec print supported-targets`で確認できます。

## RISC-V64契約

Hosted RISC-V 対象のデフォルト値は`generic-rv64`、RV64GC、`lp64d`ABIです。 Freestandingターゲットのデフォルト値は`generic-rv64`、RV64IMAC、`lp64`です。

```shell
wavec print target-spec --target riscv64-unknown-linux-gnu --format=json
wavec print target-spec --target riscv64-unknown-none-elf --format=json
```

サポートRISC-V CPUは、「generic」、「generic-rv64」、「rocket-rv64」、「sifive-u74」です。 Feature override は、`m`、`a`、`f`、`d`、`c`、`zicsr`、`zifencei` という名前の前に符号を付けてカンマで区切ります。

```shell
wavec build main.wave \
  --target riscv64-unknown-linux-gnu \
  --features=+m,+a,+f,-d,+c,+zicsr \
  --abi=lp64f
```

RISC-V検証は一貫性のない組み合わせを拒否します。 `d`には`f`が必要、`f`には`zicsr`が必要です。 `lp64`、`lp64f`、`lp64d`はアクティブな浮動小数点featureと一致する必要があります。 ABIを直接指定しないと、コンパイラはfeatureからABIを派生します。

## フリースタンディングリンク

```shell
wavec build kernel.wave \
  --target riscv64-unknown-none-elf \
  --freestanding \
  --entry=_start \
  --linker-script=linker.ld \
  --no-start-files \
  -o kernel.elf
```

`--freestanding`は、デフォルトライブラリを使用しない側にビルド設定を調整します。 `--entry`はリンカエントリを、`--linker-script`はスクリプトを指定し、`--no-start-files`はホスト起動ファイルの除外を指定します。

実際の実行前にリンク計画を確認するときは、`--dry-run`を使用できます。

## Hosted クロスリンク

他のOS・CPUで実行するプログラムを作成する場合は、対象環境のライブラリパスをsysrootに指定します。別のリンカーを使用する場合は、`-C linker`でパスを指定します。

```shell
wavec build main.wave \
  --target riscv64-unknown-linux-gnu \
  --sysroot /path/to/riscv64-sysroot \
  -C linker=/path/to/target-linker \
  -o app-riscv64
```

sysrootとリンクするライブラリは、選択したOS・CPU・ABIに合わせます。

## クロスビルドで確認する

- target tripleがコンパイラのサポートリストにあるか
- sysrootとリンカーが対象ABIに合うか
- リンクライブラリがターゲットアーキテクチャ用かどうか
- CPU featureが対象CPUで有効か
- プリスタンディングの場合、エントリシンボルとメモリ配置がリンカスクリプトと一致するかどうか

## WebAssembly実行

wasm64結果はmemory64をサポートする実行環境で使用されます。ファイル・時間・入力などの外部機能を使用するモジュールは、その機能に対応するhostimportを接続する必要があります。

ターゲットコードを作成して接続するプロセスは、`--dry-run`で確認できます。実際の実行は、選択したOSまたはWebAssembly実行環境で行われます。
