---
translation_set_id: practice-binary-message
path: practice/binary-message
locale: zh
group: practice
group_order: 4
order: 3
title: 项目：创建二进制消息
summary: 使用显式字节排序和 ULEB128 并拒绝短输入。
---

## 消息格式

前 2 个字节存储类型号 big-endian u16，然后是值 ULEB128 u64。如果将结构内存按原样写入文件，它将受到填充和字节顺序的影响，因此请按字段对其进行编码。

将其保存为`main.wave`并运行。

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

执行结果：

```text
kind=7 value=300 bytes=4
preserved=0 99
```

## 要检查什么

在reader中，我们传递实际写入的长度，而不是数组的总容量12。这是为了避免读取未初始化的尾随字节作为输入。即使短输入读取失败，也会保持position=0 和kind=99。

## 扩展练习和评论

如果格式不允许末尾有多余字节，请在解析完成后检查`reader.position == reader.len`。添加长度字段时，请确保其不大于输入的剩余字节，并且长度+offset的计算不超出范围。

[参见bytes](/docs/zh/stdlib/bytes)

## 查看实际字节数

类型编号 7 是 big-endian u16，因此是 `00 07`。值 300 变为 ULEB128 到 `AC 02`。整个消息由以下四个字节组成：

```text
00 07 AC 02
───── ─────
종류  값
```

ULEB128 每个字节使用低 7 位作为值，如果高位为 1，则表示接下来是下一个字节。 300的低7位是44，余数是2。第一个字节是44加上128个连续标记，或者172，或者0xAC。最后一个字节0x02没有连续标记。

## 容量和使用时间

u16使用2个字节，u64的ULEB128最多使用10个字节，因此一个12字节的数组可以存储这两个字段。值 300 仅使用 2 个字节，使得实际消息为 4 个字节。当发送到文件或套接字时，您发送字节writer.position而不是整个数组。

reader中的position是当前读取的位置。读取类型后变为2，读取值后变为4。为了在发生错误时继续处理下一条消息，必须知道失败消息的边界。消息恢复策略不会仅仅因为读取函数保留位置而自动建立。

## 边界值练习

将值更改为 0、127、128、16383、16384 以确定编码长度。当从127变为128时，ULEB128的长度从1增加到2，从16383变为16384时，长度从2增加到3。计算总长度，包括类型字段的2个字节。

最后一个字节被截断的消息也会被检查。如果原始消息长度为 4，我们将长度 3 传递给reader。读取了类型字段，但读取值一定会失败，因为ULEB128的最后一个字节丢失了。此时，检查位置2和读取ULEB128之前的输出值是否保持。
