---
translation_set_id: compiler
path: getting-started/compiler
locale: ja
group: getting-started
group_order: 1
order: 3
title: コンパイラコマンドリファレンス
summary: wavecコマンド、ビルドパイプライン、出力、ターゲット、診断、依存関係、およびツールクエリについて説明します。
---

## コマンドモデル

`wavec`はコンパイラCLIです。個々の入力を直接コンパイルし、ツールにコンパイラサポート情報を提供し、インストールされている標準ライブラリソースを管理します。

```text
wavec [global-options] <command> [command-options]
```

|コマンド|用途|
| --- | --- |
| `wavec build <input...>` |フラグに従ってチェック、コード生成、リンク、または実行パイプラインを実行します。|
| `wavec check <file>` |`build <file> --emit=check`のエイリアスです。|
| `wavec run <file> [-- <args...>]` |`build <file> --run`のエイリアスで、`--`の後の引数をプログラムに渡します。|
| `wavec print <item>` |ターゲットとツールチェーンのサポート情報を問い合わせます。|
| `wavec install std` |標準ライブラリをインストールします。|
| `wavec update std` |インストールされている標準ライブラリを更新します。|
| `wavec --version` |インストールしたバージョン情報を出力します。|

`wavec --help`でコマンドとオプションの完全なリストを見ることができます。

## ビルド、検査、実行

```shell
wavec build main.wave
wavec check main.wave
wavec run main.wave -- first-argument second-argument
```

`build`はデフォルトで実行ファイルを作成します。 `check`はフロントエンド検査を終えた後停止します。 `run`はバイナリ出力が必要で、共有ライブラリのビルドでは使用できません。

コンパイル、リンク、実行なしで要求を検証し、実行する手順を確認するには、`--dry-run`を使用します。

```shell
wavec build main.wave --target riscv64-unknown-linux-gnu --dry-run
wavec build main.wave --dry-run --error-format=json
```

JSONフォーマットは、Vexのようなビルドツールが使用する安定した統合インターフェースです。

## emitと入力種類

```shell
wavec build main.wave --emit=ast
wavec build main.wave --emit=ir
wavec build main.wave --emit=bc
wavec build main.wave --emit=asm
wavec build main.wave --emit=obj -o main.o
wavec build main.wave --emit=bin -o app
```

出力エミットの種類は `ast`、`ir`、`bc`、`asm`、`obj`、`bin`です。 `check`は制御モードなので、単独で使用する必要があります。パイプラインで許可される出力の種類は、カンマで複数指定できます。

入力種類は`wave`、`ir`、`bc`、`asm`、`obj`、`archive`です。 `--input-type=<kind>`はすべての入力の種類を強制的に指定します。 objectまたはarchive入力のみリンクするときは、バイナリemitと`--link-only`を使用します。

```shell
wavec build module.o --input-type=obj --link-only --emit=bin -o app
```

## 出力位置

|オプション|効果|
| --- | --- |
| `-o <file>` |メイン出力パスを指定します。|
| `--out-dir <dir>` |emit 出力を指定したディレクトリに置きます。|
| `--target-dir <dir>` |中間出力とデフォルト出力ルートを指定します。|

## 最適化と診断出力

```shell
wavec -O2 build main.wave
wavec --debug-wave=tokens,ast build main.wave
```

最適化ステップは、「-O0」、「-O1」、「-O2」、「-O3」、「-Os」、「-Oz」、「-Ofast」です。 `--debug-wave` には `tokens`、`ast`、`ir`、`mc`、`hex`、`all` を使うことができ、いくつかのステップはコンマで結合できます。

## ネイティブリンク

```shell
wavec --link=m -L ./lib build main.wave
wavec build main.wave --shared -o libexample.so
wavec build main.wave --static -o app
wavec build main.wave --pie -o app
```

`--link=<lib>`はネイティブライブラリを追加し、`-L <path>`は検索パスを追加します。リンクモードでは、互換規則に従って`--shared`、`--static`、`--pie`、`--no-pie`を使用します。

バックエンドとリンカー制御オプションは次のとおりです。

- `--target`, `--cpu`, `--features`, `--abi`, `--sysroot`
- `-C linker=<path>`と`-C link-arg=<arg>`
- `-C link-sysroot=<path>`と`-C relocation-model=<model>`
- `-C no-default-libs`

カーネルのようなプリスタンディング出力は、`--freestanding`とともに環境に合った`--entry`、`--linker-script`、`--no-start-files`設定を使用します。

## 外部パッケージの解釈

```shell
wavec --dep-root .vex/deps build main.wave
wavec --dep math=/opt/wave-deps/math build main.wave
```

`--dep-root <dir>`は、外部`package::module`importを検索するルートを追加します。 `--dep <name>=<path>`はパッケージ名を1つのディレクトリに固定します。これはコンパイラ統合ポイントであり、プロジェクトmanifest、依存関係のダウンロードとlockfileはVexが担当します。

## サポート機能の問い合わせ

ターゲットや成果物の種類を使用するツールは、`wavec print`でサポート情報を問い合わせることができます。

```shell
wavec print host-target
wavec print target-spec --format=json
wavec print supported-targets
wavec print supported-input-types
wavec print supported-emit-kinds
wavec print supported-print-items
wavec print cpu-list --target riscv64-unknown-linux-gnu
wavec print target-features --target riscv64-unknown-linux-gnu
wavec print default-linker
wavec print sysroot
wavec print std-path
wavec print dep-search-paths
```

`host`、`default-target`、`target-list`などの項目も問合せできます。構造化出力をサポートする項目は、`--format=json`を受け取ります。

## コンパイラとツールチェーンの境界

`wavec`はソースチェック、コード生成、リンクを担当します。 Vexはパッケージmanifest、依存グラフ、lockfileと再現可能なパッケージビルドを担当します。 Whaleは、独立して実行する低レベルツールチェーンです。

## std パス指定

```shell
wavec --std-root /absolute/path/to/std check main.wave
wavec --std-root /absolute/path/to/std run main.wave
```

指定されたstdパスはインストールパスより優先され、誤ったパスまたは互換性のないstdの場合は失敗します。他のインストールのstdに自動的に置き換えられません。コンパイラに対応するstdを選択してください。

出力 `-o`にはソース・入力ファイルと異なるパスを使用します。 `check`はランタイム動作まで確認しないため、[練習](/docs/ja/practice/input-calculator)では実行結果も確認します。
