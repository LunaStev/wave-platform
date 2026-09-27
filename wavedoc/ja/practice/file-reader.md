---
translation_set_id: practice-file-reader
path: practice/file-reader
locale: ja
group: practice
group_order: 4
order: 2
title: プロジェクト: ファイルの読み取りとバイト数のカウント
summary: Buffer、ファイルI/O、失敗処理とメモリ解放を接続します。
---

## 準備

作業ディレクトリに、`Wave` の後に 1 つの LF 改行を含む input.txt を作成します。ファイルには 5 バイトが含まれます。 CRLF の場合は 6 バイトが含まれます。 UTF-8 BOM はさらにバイトを追加します。エディターのファイルエンコーディングと行末を確認してください。

プログラムを`main.wave`に保存し、同じディレクトリで`wavec run main.wave`として実行します。相対パスは、実行している作業ディレクトリに基づいています。

<!-- wave-example: file-reader -->
```wave
import("std::buffer::alloc")::{
    Buffer, buffer_init, buffer_free
};
import("std::fs::file")::{
    read_to_end
};

fun main() -> i32 {
    var data: Buffer;
    if (buffer_init(&data, 0) < 0) {
        return 1;
    }

    var count: i64 = read_to_end("input.txt", &data);
    if (count < 0) {
        println("read failed={}", count);
        buffer_free(&data);
        return 2;
    }

    var lines: i64 = 0;
    for (var i: i64 = 0; i < data.len; i += 1) {
        if (deref data.data[i] == 10) {
            lines += 1;
        }
    }

    println("bytes={} LF={}", count, lines);
    if (buffer_free(&data) < 0) {
        return 3;
    }

    return 0;
}
```

実行結果：

```text
bytes=5 LF=1
```

## 動作と責任

`read_to_end`はファイルを開閉しますが、Buffer解除は発信者の責任です。成功と読み取り失敗パスの両方から解放します。バイト数のデータなので、NUL 終了文字列とは仮定しません。

LF個数と人が考える行数は常に同じではありません。最後の行にLFがない場合、このプログラムのLFの数には含まれません。ファイル名を変更して、読み込み失敗も確認してください。具体的なエラー数は環境によって異なる場合があります。

## 拡張練習と解説

非常に大きなファイルを処理するには、固定サイズの配列と`io_read`繰り返しに置き換えます。正の戻り範囲のみを処理し、0から終了します。ファイル全体をメモリに保存しなくても、バイト数とLF数を累積できます。直接開いた場合は、すべての出口パスで記述子も閉じます。

[fsとio参照](/docs/ja/stdlib/files-io) · [Buffer参照](/docs/ja/stdlib/buffer)

## 処理フローに従う

1. 空のBufferを初期化します。まだファイルの内容はありません。
2. read_to_endがファイルを読み込み、必要なスペースを増やします。
3. 読み取りに成功すると、data.len範囲のバイトを調べます。
4. LFのバイト値10に会うたびにlinesを増やします。
5. 結果を出力し、Bufferを解除します。

data.capは確保した記憶領域であり、data.lenは有効なデータ長です。繰り返し条件をcapに置き換えるとファイルになかったバイトまで読み取られるので、lenを使用します。 countは、今回のread_to_end呼び出しで追加したバイト数です。この例は空のBufferから始まるので、countとdata.lenが同じです。

## 入力を変更して確認する

|input.txt内容| bytes | LF |理由|
| --- | --- | --- | --- |
|空のファイル| 0 | 0 |読み込むバイトがない|
| `Wave` | 4 | 0 |最後の改行はありません|
| `Wave` + LF | 5 | 1 |最後のLFまでのデータ|
| `A` + LF + `B` + LF | 4 | 2 |2つのLFを数えます|
| `Wave` + CRLF | 6 | 1 |CRも1バイトですがLFのみ|

最後の行にLFがないファイルも1行に数えるには、ファイルが空でなく最後のバイトが10以外のときに行数に1を加算します。 data.lenが0かどうかを最初に確認する必要があり、最後の要素にアクセスできます。

## 大きなファイルに拡張する

コンテンツ全体を保持する現在の方法は、後でデータを再読み込みまたは検索するときに便利です。バイト数とLF数だけが必要な場合は、固定サイズのバッファを再利用する方法が適しています。

[固定サイズのバッファでファイルを読み込む](/docs/ja/stdlib/files-io)例のio_read反復文の中で、返されたバイト数だけLFを数えます。ファイル全体の長さではなく、バッファサイズと同じメモリで処理できます。読み取りがゼロを返す場合は終了であり、負の場合はエラーです。
