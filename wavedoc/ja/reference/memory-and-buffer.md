---
translation_set_id: memory-buffer
path: reference/memory-and-buffer
locale: ja
group: stdlib
group_order: 1
order: 4
title: mem: 割り当て、再割り当て、およびレイアウト
summary: バイト単位のサイズ、割り当て失敗、再割り当ての境界と解放の責任について説明します。
---

## 割り当てと解放

```text
std::mem::alloc
mem_alloc(size: i64) -> ptr<u8>
mem_alloc_zeroed(size: i64) -> ptr<u8>
mem_free(p: ptr<u8>, size: i64) -> i64
mem_realloc(old_ptr: ptr<u8>, old_size: i64, new_size: i64) -> ptr<u8>
```

サイズはバイト単位です。 0以下のサイズの割り当ては、nullを返します。正のサイズも割り当てに失敗した場合は、nullかもしれません。 `mem_alloc`の初期内容を想定せず、0初期化が必要な場合は`mem_alloc_zeroed`を使用します。

呼び出し元は成功した各割り当てを所有しており、それを解放するときに元のサイズを渡す必要があります。 `mem_free(null, size)` は 0 を返します。null 以外のポインターと正でないサイズの組み合わせはエラーです。割り当てがすでに解放された後は、決してアクセスしたり解放したりしないでください。

## 再割り当ての場合

|リクエスト|アクション|
| --- | --- |
|新しいサイズが正で成功|`min(old_size, new_size)`バイトをコピーして以前の割り当てを解除|
|正のサイズの新しい割り当てに失敗しました|null リターン、既存の割り当てを維持|
|`old_ptr == null`、正の新しいサイズ|新しい割り当てのように動作|
| `new_size == 0` |有効な以前の割り当てを解除しようとし、nullを返す|
|負のサイズまたは既存のポインタにold_size=0|null 返却|

サイズ 0 への再割り当ての結果nullは、解放が成功したかどうかを証明しません。ステータスが必要な場合は、`mem_free` に直接お電話ください。割り当てを拡大した後、新しく追加したリージョンを自分で初期化します。

## 既存のポインタを保存する例

以下は、新しいサイズが正の場合です。 `main.wave`で保存して実行します。

<!-- wave-example: reallocation -->
```wave
import("std::mem::alloc")::{
    mem_alloc, mem_realloc, mem_free
};

fun main() -> i32 {
    var data: ptr<u8> = mem_alloc(4);
    if (data == null) {
        return 1;
    }
    deref data[0] = 7;
    var grown: ptr<u8> = mem_realloc(data, 4, 8);
    if (grown == null) {
        mem_free(data, 4);
        return 2;
    }
    data = grown;
    println("{}", deref data[0]);
    if (mem_free(data, 8) < 0) {
        return 3;
    }

    return 0;
}
```

実行結果：

```text
7
```

失敗確認前に`data = mem_realloc(...)`で上書きすると、既存のアドレスが失われる可能性があります。再割り当てが成功すると、以前のアドレスとその内部を指していたポインタは使用されません。

## ターゲットタイプのサイズと配置

```text
std::mem::layout
size_of<T>() -> u64
align_of<T>() -> u64
```

両方の値は、実行中のコンピュータではなく、コンパイル対象のレイアウトです。 `size_of`にはテールパディングが含まれており、値を作成または評価しません。要素数にサイズを掛けるときは、オーバーフローを調べます。 `std::mem::ops`の`mem_size_mul_checked`と`mem_size_add_checked`を使用できます。

`mem_copy`は重ならない範囲、`mem_move`は重なり合う範囲のコピーに使用します。どちらもポインタだけで実際の割り当て長を見つけることができないため、呼び出し側は範囲を保証する必要があります。サイズが変わるバイトのリストは、[Buffer](/docs/ja/stdlib/buffer)で管理できます。

## レイアウト照会の例

main.waveで保存して実行します。文書で扱うターゲットのi32サイズとソートはそれぞれ4バイトなので、`4 4`を出力します。他のタイプ、特に構造体とポインタの値は、ターゲットごとに確認してください。

<!-- wave-example: layout-api -->
```wave
import("std::mem::layout")::{
    size_of, align_of
};

fun main() {
    println("{} {}", size_of<i32>(), align_of<i32>());
}
```

## サイズ計算でオーバーフローを確認する

割り当て関数に渡す前に、`count * element_size`が有効であることを確認する必要があります。あふれた値で小さなスペースを割り当て、元の数だけ書き込むと境界を外れます。

<!-- wave-example: book-memory-size-check -->
```wave
import("std::mem::ops")::{mem_size_add_checked, mem_size_mul_checked};

fun main() {
    var result: i64 = 99;

    if (mem_size_mul_checked(3, 4, &result) == 0) {
        println("bytes={}", result);
    }

    if (mem_size_add_checked(9223372036854775807, 1, &result) < 0) {
        println("overflow rejected");
    }
}
```

実行結果：

```text
bytes=12
overflow rejected
```

失敗した結果は割り当てサイズとして書き込まれません。一般算術で計算した後、結果が負であるかどうかを見るだけですべてのオーバーフローを検出することもできません。検査が必要なサイズの計算には、最初からchecked関数を使用します。

## 重複するコピーにはmem_move

同じ配列の一部を後ろに移動すると、入力領域と出力領域が重なります。 mem_copyに重なる範囲を渡さず、mem_moveを使用してください。

<!-- wave-example: book-memory-overlap -->
```wave
import("std::mem::ops")::{mem_move};

fun main() {
    var data: array<u8, 5> = [1, 2, 3, 4, 5];

    mem_move(&data[1], &data[0], 4);

    for (var index: i32 = 0; index < 5; index += 1) {
        println("{}", data[index]);
    }
}
```

実行結果：

```text
1
1
2
3
4
```

元の最初の 4 バイトは 1 つ右に移動します。これらを手動で前方にコピーすると、すでに上書きされた値が読み取られ、誤ってすべて 1 が生成される可能性があります。重複を認識する API がコピー方向を処理します。

## 所有権を渡す関数の作成

メモリを返す関数は、成功時に戻りアドレスと解放に必要なサイズを一緒に提供する方が良いです。発信者がそのサイズを推測する必要があると、誤った解放が発生しやすくなります。借りたアドレスを返す関数であれば、呼び出し元が解放しないことと元の寿命を説明します。

関数の境界を越えて、`allocator → owner → deallocation` を追跡できるはずです。ポインタ変数の名前や型によって所有権が自動的に決定されるわけではありません。
