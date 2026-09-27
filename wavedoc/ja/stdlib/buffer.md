---
translation_set_id: stdlib-buffer
path: stdlib/buffer
locale: ja
group: stdlib
group_order: 1
order: 5
title: buffer: 拡張可能なバイトストレージ
summary: Buffer初期化、追加、照会、容量、および解放の規則について説明します。
---

## Bufferの意味

`std::buffer::types`の`Buffer`には`data: ptr<u8>`、`len: i64`、`cap: i64`があります。 lenは初期化され使用中のバイト数、capは割り当てられた合計バイト数です。常に`0 <= len <= cap`を維持してください。文字列NUL終了は自動的に保証されません。

## 基本API

|モジュール|宣言|意味|
| --- | --- | --- |
| `std::buffer::alloc` | `buffer_init(out_buf: ptr<Buffer>, capacity: i64) -> i64` |新しいストレージの初期化。すでに所有しているバッファに再呼び出ししない|
|同じモジュール| `buffer_free(buf: ptr<Buffer>) -> i64` |割り当て解除。成功すると空の状態|
|同じモジュール| `buffer_reserve(buf: ptr<Buffer>, required_cap: i64) -> i64` |最小総容量を確保。 len維持|
|同じモジュール| `buffer_clear(buf: ptr<Buffer>) -> i64` |容量を維持し、len=0|
|同じモジュール| `buffer_resize(buf: ptr<Buffer>, new_len: i64, value: u8) -> i64` |長さ変更、新しいバイトをvalueで埋める|
| `std::buffer::write` | `buffer_push(buf: ptr<Buffer>, value: u8) -> i64` |1バイト追加|
|同じモジュール| `buffer_append(buf: ptr<Buffer>, src: ptr<u8>, size: i64) -> i64` |sizeバイトコピーして追加|
|同じモジュール| `buffer_append_str(buf: ptr<Buffer>, s: str) -> i64` |NUL以外の文字列バイトを追加|
| `std::buffer::read` | `buffer_get(buf: ptr<Buffer>, index: i64, out_value: ptr<u8>) -> i64` |範囲内の1バイトの読み取り|

ステータスを返す APIは、成功すると`BUFFER_OK`(0)を返します。 `std::buffer::error`のINVALID、BOUNDS、OVERFLOW、ALLOCエラーを区別してください。数値をOSerrnoと解釈しません。 `buffer_new`は割り当て失敗を空のBufferと表現するため、失敗を区別する必要があるときは`buffer_init`を使用します。

## 実行例

<!-- wave-example: buffer-api -->
```wave
import("std::buffer::alloc")::{
    Buffer, buffer_init, buffer_free
};
import("std::buffer::write")::{
    buffer_append_str, buffer_push
};
import("std::buffer::read")::{
    buffer_get
};

fun main() -> i32 {
    var data: Buffer;
    if (buffer_init(&data, 0) < 0) {
        return 1;
    }
    if (buffer_append_str(&data, "Hi") < 0 || buffer_push(&data, 33) < 0) {
        buffer_free(&data);
        return 2;
    }

    var value: u8 = 0;
    if (buffer_get(&data, 2, &value) < 0) {
        buffer_free(&data);
        return 3;
    }

    println("{} {}", data.len, value);
    if (buffer_free(&data) < 0) {
        return 4;
    }

    return 0;
}
```

実行結果：

```text
3 33
```

`main.wave`で保存し、`wavec run main.wave`で実行します。初期容量 0 は失敗ではなく、有効な空のバッファです。追加の過程でスペースを解放します。

## 寿命と失敗

バッファを拡張するとデータが変更される可能性があります。再割り当ての可能性がある操作の後は、以前に借用したアドレスを使用しないでください。 Buffer 構造をコピーしてもその割り当ては複製されないため、その割り当てに単一の所有者を与えます。

`buffer_get`は、失敗した場合に出力引数を変更しません。一方、利便性関数`buffer_at`はエラーもゼロであるため、実際のゼロバイトと失敗を区別するには`buffer_get`を使用します。公開フィールドを直接変更して無効なlen/capを作成しません。

[メモリーAPI](/docs/ja/reference/memory-and-buffer) · [ファイルをBufferで読む練習](/docs/ja/practice/file-reader)

## 長さと容量を別々に観察する

reserveは保存スペースを確保しますが、lenを増やすことはありません。 resizeは実際の使用長を変え、伸びた部分を指定したバイトに初期化します。 clearは使用長さのみゼロにして割り当てを再利用します。

次のプログラムをmain.waveとして保存して実行します。容量の正確な成長倍数に依存せずに必要なスペースを確保していますが、確認してください。

<!-- wave-example: book-buffer-length-capacity -->
```wave
import("std::buffer::alloc")::{
    Buffer,
    buffer_init,
    buffer_reserve,
    buffer_resize,
    buffer_clear,
    buffer_free
};

fun main() -> i32 {
    var data: Buffer;

    if (buffer_init(&data, 0) < 0) {
        return 1;
    }

    if (buffer_reserve(&data, 20) < 0) {
        buffer_free(&data);
        return 2;
    }

    println("reserved length={}", data.len);

    if (buffer_resize(&data, 3, 7) < 0) {
        buffer_free(&data);
        return 3;
    }

    println("resized length={} first={}", data.len, deref data.data[0]);

    if (buffer_clear(&data) < 0) {
        buffer_free(&data);
        return 4;
    }

    println("cleared length={}", data.len);

    if (buffer_free(&data) < 0) {
        return 5;
    }

    return 0;
}
```

実行結果：

```text
reserved length=0
resized length=3 first=7
cleared length=0
```

resizeに増やしたときにvalue=7を渡したので、新しく見える3バイトはすべて7です。 reserveで確保するだけのスペースを初期化されたデータのように読み取らない。 clearの後ろにもcapは保持され、同じBufferに再度追加できます。

## 0バイトと照会失敗を区別する

buffer_getはステータスを返し、実際のバイトは出力引数として書き込みます。データが0の場合も正常成功です。

<!-- wave-example: book-buffer-zero-bounds -->
```wave
import("std::buffer::alloc")::{Buffer, buffer_init, buffer_free};
import("std::buffer::write")::{buffer_push};
import("std::buffer::read")::{buffer_get};
import("std::buffer::error")::{BUFFER_ERR_BOUNDS};

fun main() -> i32 {
    var data: Buffer;

    if (buffer_init(&data, 0) < 0) {
        return 1;
    }

    if (buffer_push(&data, 0) < 0) {
        buffer_free(&data);
        return 2;
    }

    var value: u8 = 99;

    if (buffer_get(&data, 0, &value) < 0) {
        buffer_free(&data);
        return 3;
    }

    println("stored={}", value);
    value = 99;

    if (buffer_get(&data, 1, &value) == BUFFER_ERR_BOUNDS) {
        println("outside, preserved={}", value);
    }

    if (buffer_free(&data) < 0) {
        return 4;
    }

    return 0;
}
```

実行結果：

```text
stored=0
outside, preserved=99
```

最初の照会はゼロを読み取る成功で、2番目の照会は範囲外の失敗です。失敗後にvalue=99が保持されても、それがバッファから読み取られた値であるという意味ではありません。必ず状態を一緒に確認してください。

## 練習プール: バイト累積

0から9まで追加するには、buffer_pushを繰り返して各結果を確認します。合計はi64に保存し、`0 <= index < data.len`範囲のみを読みます。バッファを扱った後、成功パスと失敗パスの両方でbuffer_freeを呼び出します。

<!-- wave-example: book-buffer-sum -->
```wave
import("std::buffer::alloc")::{Buffer, buffer_init, buffer_free};
import("std::buffer::write")::{buffer_push};

fun main() -> i32 {
    var data: Buffer;

    if (buffer_init(&data, 0) < 0) {
        return 1;
    }

    for (var value: i32 = 0; value < 10; value += 1) {
        if (buffer_push(&data, value as u8) < 0) {
            buffer_free(&data);
            return 2;
        }
    }

    var total: i64 = 0;

    for (var index: i64 = 0; index < data.len; index += 1) {
        total += deref data.data[index] as i64;
    }

    println("sum={}", total);

    if (buffer_free(&data) < 0) {
        return 3;
    }

    return 0;
}
```

実行結果：

```text
sum=45
```
