---
translation_set_id: stdlib-bytes
path: stdlib/bytes
locale: zh
group: stdlib
group_order: 1
order: 6
title: bytes:范围、游标和 ULEB128
summary: 描述长度为view的读/写字节以及在失败时保留状态。
---

## 它与字符串有何不同？

字节数据可以包含零，因此将指针与长度一起传递。 `Bytes` 和 `BytesMut` 是非拥有视图，仅当底层存储保持有效时才有效。 `BytesMut` 需要可写存储。

## 创建cursor

```text
std::bytes::cursor
bytes_reader(data: ptr<u8>, len: i64) -> ByteReader
bytes_writer(data: ptr<u8>, len: i64) -> ByteWriter
bytes_reader_read_be_u16(reader: ptr<ByteReader>, out_value: ptr<u16>) -> i32
bytes_writer_write_be_u16(writer: ptr<ByteWriter>, value: u16) -> i32
```

从 `std::bytes::types` 导入 `ByteReader` 和 `ByteWriter`。他们的位置字段标识了下一个操作的位置。它们的长度是可访问的总字节数。构造游标不会复制或分配底层内存。

`be` 是big-endian，`le` 是little-endian。如果文件类型为big-endian，则无论主机CPU的字节顺序如何，都使用`be`函数。有 16 位、32 位和 64 位signed/unsigned 读/写和一字节函数。

## 错误和状态保存

`BYTES_OK`～`std::bytes::errors`为0。INVALID表示无效范围，EOF输入不足，NO_SPACE输出能力不足，OVERFLOW表示超出可表示范围的值。仅在整个操作成功后，检查的游标操作才会前进位置。读取失败也会保留输出值。

## ULEB128

```text
std::bytes::leb128
bytes_reader_read_uleb128_u64(reader: ptr<ByteReader>, output_value: ptr<u64>) -> i32
bytes_writer_write_uleb128_u64(writer: ptr<ByteWriter>, value: u64) -> i32
```

ULEB128 以可变字节数存储无符号 64 位整数，最多使用 10 个字节。如果空间不足，写入器将保留其位置和目标字节。阅读器将不完整的输入与超过u64的值区分开来。接受终止的非最小编码。

要创建实际消息并查看短暂的输入失败，请继续执行[二进制消息练习](/docs/zh/practice/binary-message)。不要尝试将包含零的字节字符串输出为 `str`。

## 以不同的顺序读取相同的字节

字节顺序是数字的存储规则。 1、2这两个字节如果读为big-endian，就是1×256+2，如果读为little-endian，就是2×256+1。根据网络或文件类型规则进行选择。

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

执行结果：

```text
be=258 le=513
```

view 借用一个数组。由于没有单独的分配或复制，因此数组只能在有效时使用。 offset 以字节为单位，读取 16 位需要从该位置开始 2 个字节。

## 选择 offset API 和 cursor API

read/write 函数采用 offset 参数，对于直接读取指定字段位置的格式很方便。对于下一个位置取决于前一个字段长度的流，cursor 和 position 很方便。

两者混合时，要明确哪个是标准：cursor.position，还是单独的offset。避免两次添加相同位置或移动到下一个位置而没有成功读取的错误。

## 检查短输入状态

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

执行结果：

```text
position=0 value=99
```

不需要两个字节，因此输出值和位置被保留。此属性对于在获取更多输入后重试读取同一字段的设计非常有用。但是，如果view指向的存储空间已被重新分配，则该地址也必须更新。

## ULEB128的边界

0~127 使用一个字节，128 以后使用更多字节。每个字节的高位指示是否有数据跟随。 u64之外的值或连续输入太长的是OVERFLOW，它与EOF不同，后者只是输入较少。

直接检查容量不足writer是否保留其状态。

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

执行结果：

```text
position=0 byte=85
```

128需要两个字节但只有一个空格。发生故障后，第一个字节 85 仍保留。一起检查错误类型和状态保存比简单的 roundtrip 成功检查更好地描述了边界。

## 消息解析器创建顺序

1. 读取固定标头并检查类型和版本。
2. 读取长度并将其与剩余输入范围进行比较。
3. 仅将必要的数据传递到view或单独的缓冲区。
4. 如果格式需要整个消息，则还会检查额外的字节。
5. 区分 EOF 和无效格式错误并将其传递给调用者。

您可以通过连接[二进制消息练习](/docs/zh/practice/binary-message)中的字段来创建一个程序。
