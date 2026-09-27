---
translation_set_id: system-io
path: reference/system-io-network-process
locale: ja
group: stdlib
group_order: 1
order: 15
title: システムの機能とプロセス
summary: 上位APIとOSインタフェースの境界、プロセス寿命を説明します。
---

## 機能別文書

ファイルを扱うには[fsとio](/docs/ja/stdlib/files-io)、接続するには[TCP](/docs/ja/stdlib/tcp)、アドレスを調べるには[resolver](/docs/ja/stdlib/resolution)を読んでください。以下は、プロセスと低レベルのOSアクセスルールです。

## プロセス基本 API

```text
std::process::core
proc_exit(code: i32) -> !
proc_getpid() -> i64
proc_getppid() -> i64
proc_execve(path: str, argv: ptr<ptr<i8>>, envp: ptr<ptr<i8>>) -> i64
proc_waitpid_raw(pid: i64, status: ptr<i32>, options: i32) -> i64
```

`proc_exit`は呼び出し点に戻りません。終了前に必要なファイル・メモリの整理を直接行ってください。 `proc_execve`は、成功すると既存のプロセスイメージを置き換えるため、一般的な子生成関数とは異なります。 rawargv/envpは終了を示すnullポインタと各文字列のNUL終了を準備する必要があります。

`std::process::spawn`のspawn関数は生成結果を、待機関数は子の終了状態を扱います。生成の成功とプログラムの成功の終了は異なります。パイプを作成すると、親と子が未使用の端を閉じる必要があり、EOFが渡されます。キャプチャパイプを読み取らずに子の終了だけを待つと、バッファは次々に待つことができます。

## 移植性と低レベルアプローチ

fork/exec、ファイル記述子、Windowsハンドルは同じOS機能ではありません。選択したターゲットのサポートを確認し、unsupportedを通常の失敗パスとして処理します。 `std::sys`はOS固有のインターフェースであり、数値フラグとレイアウトを他のOSでは再利用しません。

外部Cライブラリと直接連動するときは、[FFI](/docs/ja/language/modules-imports-and-ffi)をお読みください。上位stdAPIを書くために任意にlibc関数を宣言する必要はありません。 [ターゲットとリンク環境](/docs/ja/whale/build-link-targets)を先に確認してください。
