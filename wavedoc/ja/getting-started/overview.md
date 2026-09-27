---
translation_set_id: overview
path: getting-started/overview
locale: ja
group: getting-started
group_order: 1
order: 1
title: Wave のドキュメントと学習ガイド
summary: Wave をインストールから実践的なプログラムまで段階的に学習し、言語ルールや標準ライブラリ API を調べます。
---

## このガイドで Wave を学びましょう

Wave ソース コードの作成、コンパイルと実行、結果の確認を学びます。プログラミングが初めての場合は、以下の順序に従ってください。別の言語を知っている場合は、各章の例を実行し、そのルールと境界ケースをすでに知っている言語と比較してください。

## 学習パス

|ステップ|章|何を学ぶか|
| --- | --- | --- |
|セットアップ| [インストール](/docs/ja/getting-started/install) |コンパイラと標準ライブラリを準備し、それらが実行されることを確認します。|
| 1 | [初めてのプログラム](/docs/ja/language/program-structure) |ソースファイルを作成、確認、実行し、終了コードを理解する|
| 2 | [変数と型](/docs/ja/language/declarations-and-types) |値を保存し、必要な範囲のタイプを選択します|
| 3 | [演算子と変換](/docs/ja/language/expressions-and-operators) |評価順序と型変換の結果を説明する|
| 4 | [条件とループ](/docs/ja/language/control-flow) |条件で分岐し、ループでデータを処理する|
| 5 | [関数](/docs/ja/language/functions-and-generics) |繰り返される操作を関数に抽出する|
| 6 | [配列](/docs/ja/language/arrays) |インデックスによって要素にアクセスし、配列を反復処理します。|
| 7 | [文字列](/docs/ja/language/strings) |文字とバイトを区別し、エスケープと文字列の長さを理解する|
| 8 | [構造体とバリアント](/docs/ja/language/structures-enums-and-aliases) |関連データをグループ化し、成功と失敗を表す|
| 9 | [ポインタとライフタイム](/docs/ja/language/explicit-memory-type-model) |アドレスを通じて元の値を変更し、その有効期間を管理します|
| 10 | [ダイナミックメモリ](/docs/ja/language/allocation) |割り当て失敗の処理とメモリの解放|
| 11 | [モジュールとジェネリック](/docs/ja/language/modules-imports-and-ffi) |ファイル間でコードを分割し、異なるタイプの関数を再利用する|
| 12 | [エラー処理](/docs/ja/language/errors) |結果を確認し、失敗した場合はリソースをクリーンアップする|
| 13 | [非同期コードの概要](/docs/ja/language/async-and-never) |Future を作成し、完了するまで待ちます|

## 知識を実践に移す

核となる章の後に、[入力計算機](/docs/ja/practice/input-calculator)、[ファイル リーダー](/docs/ja/practice/file-reader)、[バイナリ メッセージ](/docs/ja/practice/binary-message)、[TCP クライアント](/docs/ja/practice/tcp-client) を構築します。各プロジェクトで入力が成功したケースと失敗したケースの両方をテストします。

## 3 つのドキュメント タブ

- **Wave**: ガイド付き言語コースと実践的なプロジェクトを順番に進めていきます。
- **[標準ライブラリ](/docs/ja/stdlib)**: 各モジュールの API、戻り値、エラー、所有権ルール、およびプラットフォーム要件。
- **[Whale](/docs/ja/whale)**: ビルドとリンク、パッケージ管理、コマンドの使用法、および低レベルのツールチェーン。

例では、完全なプログラムと関数内に属するスニペットを区別します。ターミナルで `wavec` コマンドを実行し、`wave` コード ブロックを `.wave` ファイルに保存します。入力と出力は別々に示されています。標準入力を読み取る例では、入力する内容を指定します。

## 行き詰まったとき

[トラブルシューティング](/docs/ja/reference/diagnostics) を使用して、インストール、ソースのチェック、リンク、および実行の問題を区別します。言語規則は [構文クイック リファレンス](/docs/ja/reference/syntax-quick-reference)、コマンドは [コンパイラー リファレンス](/docs/ja/getting-started/compiler)、API は [標準ライブラリ ガイド](/docs/ja/reference/standard-library) で調べてください。
