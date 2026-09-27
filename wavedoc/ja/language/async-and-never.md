---
translation_set_id: language-async-and-never
path: language/async-and-never
locale: ja
group: language
group_order: 2
order: 13
title: 13. 非同期関数と Future
summary: 遅延実行されるFutureとawait、block_onの役割を学びます。
---

## 待つタスクを表現する

ファイル・ソケット・タイマーのように待つ作業では、計算を継続することと完了を待つことを区別する必要があります。非同期関数は、完了する結果をFutureで表します。 asyncを付けたと自動的に新しいスレッドを作成したり、すべての同期呼び出しが非同期に変わるわけではありません。

この章は関数、ポインタ、エラー処理の後に読んでください。例は、`std::task`ランチャーを使用するネイティブプログラムで、各ファイルを`wavec run main.wave`として実行します。

## Futureを作成して結果を得る

<!-- wave-example: book-async-basic -->
```wave
import("std::task" as task);

async fun calculate(value: i64) -> i64 {
    return value * 2;
}

fun main() {
    var pending: Future<i64> = calculate(21);
    var result: i64 = task::block_on(pending);

    task::shutdown();
    println("{}", result);
}
```

実行結果：

```text
42
```

calculateの宣言に少ないi64は完了後に得られる値です。呼び出し自体の結果はFuture<i64>です。一般mainでは、block_onでFutureを駆動し、完了した結果を受け取ります。

すべてのタスクをクリーンアップした後、shutdown を呼び出してエグゼキュータのリソースを解放します。タスクがまだバッファを使用している間はバッファを解放しないでください。また、未完了のタスクを無視しないでください。

## 呼び出しと本文の実行が異なる

非同期関数は遅延実行されます。呼び出した直後に本文を最後まで実行する汎用関数と区別する必要があります。

<!-- wave-example: book-async-lazy -->
```wave
import("std::task" as task);

static entered: i32 = 0;

async fun work() -> i32 {
    entered += 1;
    return 7;
}

fun main() {
    var pending: Future<i32> = work();

    println("before={}", entered);

    var result: i32 = task::block_on(pending);

    task::shutdown();
    println("after={} result={}", entered, result);
}
```

実行結果：

```text
before=0
after=1 result=7
```

Futureを作成したとき、enteredはまだゼロです。実行を駆動した後、本文が実行され、1になります。ただ Future を変数に保存したからといって作業が完了したものとして処理してはならない理由です。

## 非同期関数内で待つ

async関数の中では、awaitで他のFutureの完了を待ちます。 awaitある式の結果は完了値です。

<!-- wave-example: book-async-nested -->
```wave
import("std::task" as task);

async fun twice(value: i32) -> i32 {
    await task::yield_now();
    return value * 2;
}

async fun process() -> i32 {
    var value: i32 = await twice(20);
    return value + 2;
}

fun main() {
    var answer: i32 = task::block_on(process());

    task::shutdown();
    println("{}", answer);
}
```

実行結果：

```text
42
```

processはtwiceのFutureを待ってから完了値に2を加算します。汎用関数の呼び出し結果と同様に値を使用できますが、待機中に実行チャンスを他のタスクに渡すことができます。

yield_nowは協力的に実行機会を譲ります。長い計算ループで一度も譲歩しないと、他の作業の進行が遅くなる可能性があります。非同期は、CPU計算を自動的に並列に分散するデバイスではありません。

## 複数のジョブをスケジュールする

spawnでジョブをスケジュールし、それぞれの結果を待つことができます。 2つのジョブの中間出力順序に依存せずに最終結果を確認する例です。

<!-- wave-example: book-async-spawn -->
```wave
import("std::task" as task);

async fun compute(value: i32) -> i32 {
    await task::yield_now();
    return value * 2;
}

async fun combine() -> i32 {
    var first: Future<i32> = task::spawn(compute(10));
    var second: Future<i32> = task::spawn(compute(20));

    var left: i32 = await first;
    var right: i32 = await second;

    return left + right;
}

fun main() {
    var total: i32 = task::block_on(combine());

    task::shutdown();
    println("{}", total);
}
```

実行結果：

```text
60
```

スケジュールしたタスクごとに結果を待つ場所があります。ジョブを作成してハンドルを忘れるのではなく、誰が完了を確認するかを決める必要があります。 await 順序と内部ジョブが実行されるすべての順序を同じとは考えないでください。

## Futureは一度消費する

Futureは単一消費ハンドルとして扱われます。同じFutureをコピーして、2つの異なるタスクのように待ちません。既に完了を受けているFutureを再度block_onまたはawaitしません。

複数の場所で同じ結果が必要な場合は、Futureを複数回消費するのではなく、完了した値を保存し、その値のコピー・共有規則に従って渡します。値内にポインタや所有リソースがあるかどうかを確認する必要があります。

## タイマーと同期待機の違い

async関数内で待つときは、`await task::sleep_ms(...)`を使用できます。同期 sleep を呼び出すと、現在の実行フローをブロックし、ランチャーの他のタスクの進行にも影響を与える可能性があります。

待ち時間が正確に要求されたミリ秒に等しいとは期待しません。スケジューリングやその他の作業によっては、遅く目覚めることがあります。時間制限を実装するときは、経過時間を測定するclockとdeadlineを使用し、毎回元の全体の時間を再び待つ方法と区別します。

## バッファ寿命とキャンセル

非同期 I/O に渡されるバッファは、呼び出し元の関数が中断されている間でも有効なままでなければなりません。クリーンアップの完了またはキャンセルの前にそれを解放または再割り当てすると、操作に無効なアドレスが残る可能性があります。

キャンセル要求とジョブクリーンアップの完了は、同じ瞬間であると判断することはできません。 cancelシリーズAPIの結果を確認し、リソースを解放する前に必要な完了待ちを行います。詳細呼び出し規則については、[task参照](/docs/ja/stdlib/task)をお読みください。

## 一般的な誤解

|考える|実際に確認すること|
| --- | --- |
|async呼び出しをしたので終わった|Futureを実際に駆動して完了したか|
|async関数内のすべての呼び出しは非同期です|呼び出したAPIが同期か非同期か|
|Futureコピーはジョブの複製|同じハンドルを重複して消費していないか|
|キャンセルしたのですぐにバッファ解除可能|キャンセル後の作業の整理まで終了しましたか|
|中間出力順序は常に固定|結果に必要な順序だけを明示的に待つか|

## 演習と解答例

3つの非同期関数を連続して待つパイプラインを作成します。 2倍の計算の後に5を加えた結果を返します。

<!-- wave-example: book-async-solution -->
```wave
import("std::task" as task);

async fun read_value() -> i32 {
    return 10;
}

async fun transform(value: i32) -> i32 {
    await task::yield_now();
    return value * 2;
}

async fun pipeline() -> i32 {
    var initial: i32 = await read_value();
    var changed: i32 = await transform(initial);

    return changed + 5;
}

fun main() {
    var result: i32 = task::block_on(pipeline());

    task::shutdown();
    println("{}", result);
}
```

実行結果：

```text
25
```

この例は、意図的に逐次的な依存関係です。 transformはread_valueの結果が必要なので、単にすべてspawnとは関係がなくなりません。独立した作業と結果を必要とする作業を区別することは、非同期設計の出発点です。

基礎学習を終えたら、[ファイルの読み方](/docs/ja/practice/file-reader)と[TCP実習](/docs/ja/practice/tcp-client)で実際の外部資源と接続してください。

## voidとnever

戻り値の型を省略した汎用関数は、値なしで呼び出し点に戻ることができます。 neverタイプは`!`と表記し、呼び出し点に正常に戻らないという意味です。プロセス終了関数が代表的です。

宣言を説明する例：

```text
fun log(message: str)              // 작업 후 돌아옴, 결과값 없음
fun stop(code: i32) -> !           // 정상 반환하지 않음
async fun work() -> i64            // 호출 결과는 Future<i64>
```

neverを一般的な保存値にしようとしません。戻らないと宣言した関数が正常な戻りパスを持つように書かないでください。終了する前にリソースのクリーンアップが必要な場合は、呼び出し元が最初に実行する必要があります。

## 返さない関数の完全な例

main.waveで保存して実行すると、出力なしで終了コード0で終了します。 stopは呼び出し元に返されないため、`-> !`として宣言します。

<!-- wave-example: never-function -->
```wave
import("std::process::core")::{
    proc_exit
};

fun stop() -> ! {
    proc_exit(0);
}

fun main() {
    stop();
}
```
