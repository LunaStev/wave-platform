---
translation_set_id: standard-library
path: reference/standard-library
locale: ja
group: stdlib
group_order: 1
order: 1
title: 標準ライブラリガイド
summary: 目的に合ったモジュールを見つけ、関数のエラーと所有権ルールを読み取る方法です。
---

## 必要な機能を探す

標準ライブラリは`std::module::file`パスでimportします。名前が似ていても、関数のエラーを返す方法は異なる場合があります。最初は[API 読み方](/docs/ja/stdlib/contracts)を読み、次の表から必要なモジュールに移動します。

|やりたいこと|文書|主なimport|
| --- | --- | --- |
|文字列長・比較・検索| [string](/docs/ja/reference/string-and-bytes) | `std::string::len`, `cmp`, `find`, `trim` |
|メモリ割り当て・コピー・サイズ| [mem](/docs/ja/reference/memory-and-buffer) | `std::mem::alloc`, `ops`, `layout` |
|サイズが変わるバイトリスト| [buffer](/docs/ja/stdlib/buffer) | `std::buffer::alloc`, `read`, `write` |
|バイナリの読み書き| [bytes](/docs/ja/stdlib/bytes) | `std::bytes::types`, `cursor`, `leb128` |
|ファイル・ディスクリプタ I/O| [fsとio](/docs/ja/stdlib/files-io) | `std::fs::file`, `std::io::fd` |
|経路組み合わせ・環境設定| [pathとenv](/docs/ja/stdlib/path-env) | `std::path::copy`, `std::env::environ` |
|時間測定・待機| [time](/docs/ja/stdlib/time) | `std::time::duration`, `clock`, `sleep` |
|数字アドレス・名前リスト照会| [net.resolve](/docs/ja/stdlib/resolution) | `std::net::resolve`, `resolve_table` |
|TCP 接続・転送| [net.tcp](/docs/ja/stdlib/tcp) | `std::net::tcp`, `address`, `error` |
|OS乱数| [random](/docs/ja/stdlib/random) | `std::random::fill` |
|プロセス・OS 境界| [システム機能](/docs/ja/reference/system-io-network-process) | `std::process::core`, `spawn`, `std::sys` |
|非同期ジョブの実行| [task](/docs/ja/stdlib/task) | `std::task` |
|数学・診断ヘルパー| [mathとdebug](/docs/ja/stdlib/math-debug) | `std::math::int`, `float`, `std::debug::core` |

## 初めて使用する例

以下のプログラムは別途パッケージをダウンロードせず、stdの関数を1つ使用します。 `main.wave`として保存し、`wavec run main.wave`として実行します。

<!-- wave-example: library-import -->
```wave
import("std::string::len")::{
    len
};

fun main() {
    println("{}", len("Wave"));
}
```

結果は`4`です。同じ例を段階的に理解するには、[文字列](/docs/ja/language/strings)を読んでください。

## コンパイラに合ったstd

`wavec print std-path`で選択したパスを確認してください。他のチェックアウトのstdを使用する場合は、`wavec --std-root /absolute/path/to/std check main.wave`のようにパスを指定します。指定されたパスが無効または互換性がない場合は、エラーを表示します。

## プラットフォーム境界

文字列・バイトのような計算機能とファイル・ソケットのようなOS機能を区別してください。ターゲットを認識しても、すべてのホストAPIが提供されるわけではありません。 [サポートターゲット](/docs/ja/whale/build-link-targets)と各APIのプラットフォーム項目を一緒に読みます。 `std::sys`は低レベルのインタフェースであり、移植可能なプログラムは最初に親モジュールを使用します。
