---
translation_set_id: stdlib-bytes
path: stdlib/bytes
locale: ja
group: stdlib
group_order: 1
order: 6
title: bytes: 範囲、カーソル、ULEB128
summary: 長さを持つバイトviewと失敗時の状態を保存する読み書きを説明します。
---

## 文字列と異なる点

バイト データにはゼロを含めることができるため、長さとともにポインタを渡します。 `Bytes` と `BytesMut` は非所有ビューであり、基になるストレージが有効な間のみ有効です。 `BytesMut` には書き込み可能なストレージが必要です。

## cursor作り

```text
std::bytes::cursor
bytes_reader(data: ptr<u8>, len: i64) -> ByteReader
bytes_writer(data: ptr<u8>, len: i64) -> ByteWriter
bytes_reader_read_be_u16(reader: ptr<ByteReader>, out_value: ptr<u16>) -> i32
bytes_writer_write_be_u16(writer: ptr<ByteWriter>, value: u16) -> i32
```

`std::bytes::types` から `ByteReader` と `ByteWriter` をインポートします。それらの位置フィールドは、次の操作の場所を識別します。それらの長さは、アクセス可能な合計バイト数です。どちらのカーソルを構築しても、基礎となるメモリのコピーや割り当ては行われません。

`be`はbig-endian、`le`はlittle-endianです。ファイル形式がbig-endianの場合、ホストCPUのバイト順に関係なく、`be`関数を使用します。 16・32・64ビット signed/unsigned 読み出し・書き込みと1バイト関数があります。

## エラーと状態の保存

`std::bytes::errors`の`BYTES_OK`は0です。INVALIDは無効な範囲、EOFは入力が不十分、NO_SPACEは出力容量が不十分、OVERFLOWは表現可能な範囲外の値を示します。チェックされたカーソル操作は、操作全体が成功した後にのみ位置を進めます。失敗した読み取りでも出力値は保持されます。

## ULEB128

```text
std::bytes::leb128
bytes_reader_read_uleb128_u64(reader: ptr<ByteReader>, output_value: ptr<u64>) -> i32
bytes_writer_write_uleb128_u64(writer: ptr<ByteWriter>, value: u64) -> i32
```

ULEB128 は、最大 10 バイトを使用して、符号なし 64 ビット整数を可変バイト数で格納します。スペースが不十分な場合、ライターはその位置と宛先バイトの両方を保存します。リーダーは、不完全な入力とu64を超える値を区別します。終了した非最小エンコーディングは受け入れられます。

実際のメッセージを作成して短い入力失敗を確認するには、[バイナリメッセージの実践](/docs/ja/practice/binary-message)に進みます。 0を含むバイト列を`str`として出力しないでください。

## 同じバイトを異なる順序で読み取る

バイト順序は数値の格納規則です。 1、2という2バイトをbig-endianで読むと1×256+2、little-endianで読むと2×256+1です。ネットワークまたはファイル形式の規則に従って選択します。

<!-- wave-example: book-bytes-endian -->
```wave
import("std::bytes::types")::{Bytes, bytes_view};
import("std::bytes::read")::{bytes_read_be_u16, bytes_read_le_u16};

fun main() -> i32 {
    var data: array<u8, 2> = [1, 2];
    var view: Bytes = bytes_view(&data[0], 2);
    var big: u16 = 0;
    var little: u16 = 0;

    if (bytes_read_be_u16(view, 0, &big) < 0) {
        return 1;
    }

    if (bytes_read_le_u16(view, 0, &little) < 0) {
        return 2;
    }

    println("be={} le={}", big, little);
    return 0;
}
```

実行結果：

```text
be=258 le=513
```

viewは配列を借ります。個別の割り当てやコピーがないため、配列が有効な間のみ使用できます。 offsetはバイト単位で、16ビット読み取りにはその位置から2バイトが必要です。

## offset APIとcursor API選択

指定されたフィールド位置を直接読み取る形式には、offset引数を受け取るread/write関数が楽です。前のフィールド長によって次の位置が変わるストリームには、positionを持つcursorが楽です。

2つを混ぜるときは、cursor.positionとは別にoffsetのどれが基準かを明確にします。同じ場所を2回追加したり、読み取りに成功しなかったため、次の場所に移動する間違いを避けてください。

## 短い入力でステータスを確認する

<!-- wave-example: book-bytes-short-read -->
```wave
import("std::bytes::types")::{ByteReader};
import("std::bytes::cursor")::{bytes_reader, bytes_reader_read_be_u16};
import("std::bytes::errors")::{BYTES_ERROR_EOF};

fun main() {
    var data: array<u8, 1> = [1];
    var reader: ByteReader = bytes_reader(&data[0], 1);
    var value: u16 = 99;
    var status: i32 = bytes_reader_read_be_u16(&reader, &value);

    if (status == BYTES_ERROR_EOF) {
        println("position={} value={}", reader.position, value);
    }
}
```

実行結果：

```text
position=0 value=99
```

必要な2バイトがないため、出力値と位置が維持されます。この性質は、入力をさらに取得した後に同じフィールドの読み取りを再試行する設計に役立ちます。ただし、viewが指すストレージスペースを再割り当てした場合は、アドレスも更新する必要があります。

## ULEB128の境界

0〜127は1バイト、128からはより多くのバイトを使用します。各バイトの上位ビットは、後にデータが続くかどうかを示します。 u64を超える値または長すぎる連続入力はOVERFLOWであり、単に入力が少ないEOFとは異なります。

容量不足 writerが状態を保存しているか直接確認してください。

<!-- wave-example: book-leb128-space -->
```wave
import("std::bytes::types")::{ByteWriter};
import("std::bytes::cursor")::{bytes_writer};
import("std::bytes::leb128")::{bytes_writer_write_uleb128_u64};
import("std::bytes::errors")::{BYTES_ERROR_NO_SPACE};

fun main() {
    var data: array<u8, 1> = [85];
    var writer: ByteWriter = bytes_writer(&data[0], 1);
    var status: i32 = bytes_writer_write_uleb128_u64(&writer, 128);

    if (status == BYTES_ERROR_NO_SPACE) {
        println("position={} byte={}", writer.position, data[0]);
    }
}
```

実行結果：

```text
position=0 byte=85
```

128は2バイトが必要ですが、スペースは1つだけです。失敗後、最初のバイト85も残ります。エラーの種類と状態の保存を一緒にチェックすることは、単純なroundtrip成功チェックよりも境界をよく説明します。

## メッセージパーサーの作成順序

1. 固定ヘッダーを読み、種類とバージョンを確認します。
2. 長さを読み、残りの入力範囲と比較します。
3. 必要なデータのみをviewまたは別のバッファに渡します。
4. メッセージ全体を要求する形式であれば、残りのバイトも調べます。
5. EOFと誤った形式のエラーを区別して発信者に渡します。

[バイナリメッセージの実践](/docs/ja/practice/binary-message)でフィールドを連結して1つのプログラムを作成できます。
