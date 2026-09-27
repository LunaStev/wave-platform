---
translation_set_id: ecosystem
path: whale/ecosystem
locale: ja
group: whale
group_order: 1
order: 2
title: ツールチェーンコンポーネント
summary: 別々の低レベルツールチェーンであるWhaleの役割とWaveエコシステムのコンポーネント境界を説明します。
---

## Whaleイラン

Whaleは、アセンブリと中間表現を扱う低レベルツールチェーンです。アセンブリ、オブジェクト、リンク、および中間表現をカバーするコンポーネントは、Waveとは異なるネイティブコード生成ツールで再利用できるように設計されています。

Whaleは、Wave開発環境全体を呼ぶ名前ではありません。各プロジェクトの責任は次のように区別されます。

|プロジェクト|責任|
| --- | --- |
| `wavec` |Wave ソースを調べて実行ファイルを作成します。|
| Vex |Waveパッケージ、manifest、依存グラフ、lockfileとパッケージビルドを管理します。|
| Whale |独立したassembler、object、linkerとIRコンポーネントを提供します。|
| Wave `std` |ランタイムとシステムAPIをWaveソースモジュールとして提供します。|

## コンポーネント

Whale workspaceは、4つの主要なライブラリ領域で構成されています。

- `assembler`：トークン化、AMD64解析・エンコーディング、section、symbolとrelocation
- `object`：オブジェクトファイルモデルとELF64writer
- `linker`：リンク層
- `ir`: Whale IRタイプ、builder、出力、検証とオプション frontend socket

`whale`実行可能ファイルは、この領域を`asm`、`object`、`link`、`ir`コマンドとして提供します。

## ツール境界

Waveプログラムは`wavec`でビルドします。アセンブリ、オブジェクトファイル、IRを直接扱う場合は、`whale`コマンドを使用します。

Whaleを設置しても`wavec`のビルド方式は変わりません。 Vexは`wavec`を使ってWaveパッケージをビルドし、Whaleは低レベルの出力を扱う作業で直接実行します。

## 成果物の検証

Whale出力をビルドプロセスにリンクするときは、objectformatとターゲットarchitectureが一致することを確認してください。 symbolとrelocationは、`readelf`、`objdump`などの独立したツールで検査できます。 IR socketを使用するビルドでは、生産者とWhaleが同じsocketschemaを使用する必要があります。
