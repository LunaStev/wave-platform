---
translation_set_id: practice-binary-message
path: practice/binary-message
locale: ja
group: practice
group_order: 4
order: 3
title: プロジェクト: バイナリ メッセージの作成
summary: 明示的なバイト順序とULEB128を使用し、短い入力を拒否します。
---

## メッセージ形式

最初の2バイトにはbig-endianu16タイプ番号、次にULEB128u64値を保存します。構造体メモリをそのままファイルに書き込むと、パディングやバイト順の影響を受けるため、フィールドごとにエンコードします。

`main.wave`で保存して実行します。

<!-- wave-example: binary-message -->
```wave
import("std::bytes::types")::{
    ByteReader, ByteWriter
};
import("std::bytes::cursor")::{
    bytes_reader, bytes_writer, bytes_writer_write_be_u16, bytes_reader_read_be_u16
};
import("std::bytes::leb128")::{
    bytes_writer_write_uleb128_u64, bytes_reader_read_uleb128_u64
};

fun main() -> i32 {
    var data: array<u8, 12>;
    var writer: ByteWriter = bytes_writer(&data[0], 12);
    if (bytes_writer_write_be_u16(&writer, 7) < 0) {
        return 1;
    }
    if (bytes_writer_write_uleb128_u64(&writer, 300) < 0) {
        return 2;
    }

    var reader: ByteReader = bytes_reader(&data[0], writer.position);
    var kind: u16 = 0;
    var value: u64 = 0;
    if (bytes_reader_read_be_u16(&reader, &kind) < 0) {
        return 3;
    }
    if (bytes_reader_read_uleb128_u64(&reader, &value) < 0) {
        return 4;
    }

    println("kind={} value={} bytes={}", kind, value, reader.position);
    var short: ByteReader = bytes_reader(&data[0], 1);
    kind = 99;
    if (bytes_reader_read_be_u16(&short, &kind) >= 0) {
        return 5;
    }

    println("preserved={} {}", short.position, kind);
    return 0;
}
```

実行結果：

```text
kind=7 value=300 bytes=4
preserved=0 99
```

## 確認するポイント

readerは配列全体の容量12ではなく、実際に書いた長さを渡します。初期化されていない後方バイトを入力として読み取らないためです。短い入力読み取りが失敗しても、position=0とkind=99が維持されます。

## 拡張練習と解説

末尾に残るバイトを許可しない形式の場合は、解析完了後、`reader.position == reader.len`を確認してください。長さフィールドを追加するときは、入力の残りのバイトより大きくないことを確認し、長さ+offset計算が範囲を超えないようにする必要があります。

[bytes参照](/docs/ja/stdlib/bytes)

## 実際のバイトを見る

種類番号7はbig-endianu16なので`00 07`です。値300はULEB128から`AC 02`になります。メッセージ全体は次の4バイトです。

```text
00 07 AC 02
───── ─────
종류  값
```

ULEB128は、バイトごとに低い7ビットを値に使用し、高いビットが1の場合、次のバイトが続くことを示します。 300の低い7ビットは44、残りは2です。最初のバイトは、44に連続表示128を加えた172、すなわち0xACです。最後のバイト0x02には連続表示はありません。

## 容量と使用長

u16は2バイト、u64のULEB128は最大10バイトを使用するため、12バイト配列であれば2つのフィールドを保存できます。値300は2バイトのみを使用し、実際のメッセージは4バイトです。ファイルやソケットに転送するときは、配列全体ではなく、writer.positionバイトを送信します。

readerのpositionは現在の読み取り位置です。種類を読んだ後2、値を読んだ後4になります。途中で失敗したときに次のメッセージに進むには、まず失敗したメッセージの境界を知ることができます。読み取り関数が場所を保存するだけでは、メッセージ回復ポリシーが自動的に決まることはありません。

## 境界値の練習

値を0、127、128、16383、16384に置き換えて、エンコードの長さを確認します。 127 から 128 に変わるとき ULEB128 長さは 1 から 2 に、16383 から 16384 に変わるときは 2 から 3 に伸びます。種類フィールドの2バイトも含めて、全長を計算します。

最後のバイトが切り捨てられたメッセージも調べます。元のメッセージ長が4の場合、readerに長さ3を渡します。タイプフィールドまでは読み取られますが、ULEB128の最後のバイトがないため、値の読み取りは失敗するはずです。このとき、ULEB128読み取り直前の位置2と出力値が保持されていることを確認してください。
