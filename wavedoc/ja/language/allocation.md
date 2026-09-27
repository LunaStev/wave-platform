---
translation_set_id: learn-allocation
path: language/allocation
locale: ja
group: language
group_order: 2
order: 10
title: 10. メモリ割り当てとリソース管理
summary: 割り当て失敗、初期化、有効範囲、および解放を学びます。
---

## いつ動的ストレージスペースが必要ですか

固定サイズの配列には、その型にその長さが含まれます。ファイル サイズや入力長など、データ量が実行時にのみ判明する場合は、動的メモリを使用します。必要がなくなったら、各割り当てを解放します。

この章では、小さな割り当てを管理し、サイズを変更して、Buffer を使用します。ポインタを渡すことは、所有権を譲渡することとは異なります。例に従って、各関数が所有するリソースを特定します。

## 割り当て、確認、使用、解放

<!-- wave-example: book-alloc-lifecycle -->
```wave
import("std::mem::alloc")::{
    mem_alloc_zeroed, mem_free
};

fun main() -> i32 {
    var size: i64 = 4;
    var data: ptr<u8> = mem_alloc_zeroed(size);

    if (data == null) {
        println("allocation failed");
        return 1;
    }

    deref data[0] = 42;
    println("{} {}", deref data[0], deref data[1]);

    var status: i64 = mem_free(data, size);

    if (status < 0) {
        return 2;
    }

    return 0;
}
```

実行結果：

```text
42 0
```

プログラムには 4 つのステップがあります。4 バイトの要求、null のチェック、有効な範囲のみへのアクセス、および割り当ての解放です。成功すると、mem_alloc_zeroed はメモリをゼロに初期化するため、プログラムが書き込まれていなくても 2 番目のバイトはゼロになります。

mem_alloc によって返されるメモリの初期内容を想定しないでください。各領域を読み取る前に初期化してください。ゼロまたは負の割り当てサイズは、null を返します。正のサイズの割り当てでもメモリの取得に失敗する可能性があります。

## サイズの単位

メモリ割り当て API のサイズ引数はバイト単位で測定されます。 10 個の整数を割り当てるには、要素のサイズに要素の数を掛けます。この乗算がオーバーフローしないことを確認してください。

<!-- wave-example: book-alloc-typed -->
```wave
import("std::mem::alloc")::{
    mem_alloc, mem_free
};
import("std::mem::layout")::{
    size_of
};
import("std::mem::ops")::{
    mem_size_mul_checked
};

fun main() -> i32 {
    var count: i64 = 3;
    var bytes: i64 = 0;
    var item_size: i64 = size_of<i32>() as i64;

    if (mem_size_mul_checked(count, item_size, &bytes) < 0) {
        return 1;
    }

    var raw: ptr<u8> = mem_alloc(bytes);

    if (raw == null) {
        return 2;
    }

    var values: ptr<i32> = raw as ptr<i32>;

    for (var index: i64 = 0; index < count; index += 1) {
        deref values[index] = (index + 1) as i32;
    }

    println("{} {} {}", deref values[0], deref values[1], deref values[2]);

    if (mem_free(raw, bytes) < 0) {
        return 3;
    }

    return 0;
}
```

実行結果：

```text
1 2 3
```

count は要素数です。 bytes はバイト数です。ポインター演算は i32 の単位で移動しますが、割り当てを解放するには元のサイズ (バイト単位) が必要です。 size_of はターゲット タイプ レイアウトを使用し、要素タイプとの関係を明示します。

一般的な非常に大きなタイプにsize_of結果をi64に変更する場合は、その変換範囲も考慮する必要があります。ここではサイズが知られているi32を使用します。

## 失敗パスでもクリーンアップ

割り当て後に別の操作が失敗した場合は、早期に戻る前にメモリを解放してください。所有権テーブルは、見逃してしまう可能性のあるクリーンアップ パスを特定するのに役立ちます。

|ステップ|所有リソース|失敗した場合|
| --- | --- | --- |
|割り当て前|なし|すぐに返す|
|割り当て成功後|dataとオリジナルサイズ|data解除後に返却|
|再割り当て成功後|新しい住所と新しいサイズ|新しいアドレスを解除|
|解除後|なし|以前の住所の使用を禁止|

ポインタ変数を上書きして元のアドレスを失うと、割り当てを解放するために必要な情報も失われます。これによりメモリ リークが発生します。逆に、2 人の所有者を通じて同じ割り当てを解放すると、二重解放が発生します。

## サイズを増やす再割り当て

<!-- wave-example: book-alloc-grow -->
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
    var next: ptr<u8> = mem_realloc(data, 4, 8);

    if (next == null) {
        mem_free(data, 4);
        return 2;
    }

    data = next;
    deref data[4] = 9;
    println("{} {}", deref data[0], deref data[4]);

    if (mem_free(data, 8) < 0) {
        return 3;
    }

    return 0;
}
```

実行結果：

```text
7 9
```

まず結果を next に格納します。より大きなブロックの割り当てが失敗した場合でも、元のデータは有効なままであり、解放することができます。成功すると、古い割り当てが解放され、新しいアドレスを使用する必要があります。新しく追加された領域を読み取る前に初期化します。

この例の新しいサイズは正です。 new_size=0 のリクエストは、代わりに既存の割り当てを解放しようとし、null を返します。したがって、null という結果は、古い割り当てがまだ有効であることを必ずしも意味するわけではありません。解放が成功したかどうかを確認する必要がある場合は、mem_free に直接呼び出してください。

## 借りたポインタを再確認する必要がある時点

data 内部を指すポインタを保存しておき、再割り当て後に使用するのは間違っています。新しいdataの住所が変わる可能性があるからです。内部位置が必要な場合は、アドレスの代わりにoffsetを保存し、成功後に新しいdataに基づいて再計算できます。

メモリの解放または再割り当ては、メモリを借用しているコードにも影響します。別の操作がまだそのメモリを使用しているかどうかを確認してください。非同期操作に渡されるバッファは、操作が終了するまで有効なままでなければなりません。

## バイトリストにはBuffer

長さが頻繁に変わるバイトリストを直接再割り当てしながら管理するには、lenとcap、拡張失敗、サイズ計算の両方を処理する必要があります。 stdのBufferはこのような作業を結び付けます。

<!-- wave-example: book-buffer-progressive -->
```wave
import("std::buffer::alloc")::{
    Buffer, buffer_init, buffer_free
};
import("std::buffer::write")::{
    buffer_append_str
};

fun main() -> i32 {
    var message: Buffer;

    if (buffer_init(&message, 0) < 0) {
        return 1;
    }

    if (buffer_append_str(&message, "Hello") < 0) {
        buffer_free(&message);
        return 2;
    }

    if (buffer_append_str(&message, ", Wave") < 0) {
        buffer_free(&message);
        return 3;
    }

    println("bytes={}", message.len);

    if (buffer_free(&message) < 0) {
        return 4;
    }

    return 0;
}
```

実行結果：

```text
bytes=11
```

len は使用中のバイト数です。 cap は割り当てられた容量です。データを追加すると、必要に応じて割り当てが増加します。 Buffer を使用しても、発信者のそれを解放する責任が免除されるわけではありません。

buffer_append_strは文字列の末尾にNULを追加しません。したがって、message.dataをすぐにstrに出力しないでください。バイト列は、長さを受け取るI/O関数として出力するか、明示的に文字列表現を構成します。

## 演習と解き方

0から9までのバイトをBufferに順番に追加し、合計を求めます。各追加が失敗した場合は解放する必要があり、読み取りはlenの範囲でのみ行われます。完全なプールと境界の失敗は、[Buffer 使い方](/docs/ja/stdlib/buffer)の例に従って確認できます。

磁気コードで割り当て・再割り当て・解除呼び出しに表示をしてみてください。成功した割り当てごとに誰が所有し、どのパスから解放するかを説明できる必要があります。
