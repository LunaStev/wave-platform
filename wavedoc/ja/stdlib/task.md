---
translation_set_id: stdlib-task
path: stdlib/task
locale: ja
group: stdlib
group_order: 1
order: 13
title: task: Future の実行と後処理
summary: 非同期ジョブの単一消費、実行およびキャンセル後のクリーンアップについて説明します。
---

## 基本API

`import("std::task" as task);`にインポートします。

```text
block_on<T>(future: Future<T>) -> T
spawn<T>(future: Future<T>) -> Future<T>
cancel<T>(future: Future<T>) -> bool
yield_now() -> Future<void>
sleep_ms(milliseconds: i64) -> Future<void>
cancel_and_join<T>(future: Future<T>) -> Future<void>
shutdown()
```

`task::block_on(pending)`はFutureが完了するまで実行し、結果を返します。結果のタイプは、渡されたFutureで決定されます。

## 寿命ルール

Futureは一度消費するハンドルです。値をコピーして独立した2つの操作だとは思いません。 `spawn`で予約した作業も結果を待つかキャンセル・整理する責任があります。 `await`と`block_on`ですでに消費しているFutureを再消費しません。

キャンセルはキャンセルを要求します。それ自体ではクリーンアップが完了したことを保証するものではありません。リソースを解放する前に、必要な完了を待ってください。タスクがまだアクセスできる間は、非同期 I/O によって借用されたメモリを解放しないでください。タスクが完了またはキャンセルのクリーンアップが完了したら、`shutdown` に電話します。

## 連携実行

長い計算と同期blocking呼び出しは、ランチャー全体の進行を遅らせることができます。 yieldと非同期待機は実行機会を譲ります。ブロッキングI/Oが単にasync関数内にあるため、非同期にはなりません。

[非同期コードの概要](/docs/ja/language/async-and-never)のプログラム全体で実行とクリーンアップの順序を確認してください。

## 譲歩と完了待ち

次のプログラムは、ジョブの途中で実行の機会を譲り、結果を返します。 yieldは関数終了ではないため、awaitの後のコードが続いて実行されます。

<!-- wave-example: book-task-yield -->
```wave
import("std::task" as task);

async fun calculate() -> i32 {
    println("started");
    await task::yield_now();
    println("resumed");
    return 42;
}

fun main() -> i32 {
    var pending: Future<i32> = calculate();
    var result: i32 = task::block_on(pending);

    println("result={}", result);
    task::shutdown();
    return 0;
}
```

実行結果：

```text
started
resumed
result=42
```

この例には他の操作がないため、出力順序は一定です。複数のジョブをspawnしたプログラムでは、yieldポイントで異なるジョブが進行する可能性があるため、異なるジョブの出力順序に依存しません。

## ジョブを整理する順序

1. 作業に必要なストレージスペースとリソースを準備します。
2. Futureを作成し、await、block_onまたはspawnとして実行します。
3. 結果が必要な場合は完了まで待ちます。
4. 実行中のジョブをキャンセルした場合は、そのジョブがクリーンアップされるまで待ちます。
5. ジョブが借りたバッファとファイル・ソケットをまとめます。
6. 残りの作業がない状態でshutdownを呼び出します。

Future変数の範囲外と作業が安全に整理されるのは別々です。特に、関数ローカル配列のアドレスをasyncタスクに渡した場合は、関数が返される前にタスクがその配列の使用を終了する必要があります。

## async 関数と一般関数の分割

純粋な計算は通常の関数に分割できます。待機を表現しなければならない関数にasyncを付け、その関数の中でawaitで完了を待ちます。ファイル読み込みのように長くかかる同期関数をasync関数で包むだけで、他のタスクに実行の機会を与えません。

[非同期学習](/docs/ja/language/async-and-never)からFutureを作成する時点と実行する時点を比較できます。
