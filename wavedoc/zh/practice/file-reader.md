---
translation_set_id: practice-file-reader
path: practice/file-reader
locale: zh
group: practice
group_order: 4
order: 2
title: 项目：读取文件并计算字节数
summary: Buffer，文件I/O，链接故障处理和内存释放。
---

## 准备好了

在工作目录中创建 input.txt，其中包含 `Wave`，后跟单个 LF 换行符。该文件包含 5 个字节。 CRLF 包含 6 个字节； UTF-8 BOM 添加更多字节。检查编辑器的文件编码和行结尾。

将程序保存为`main.wave`，并从同一目录中以`wavec run main.wave` 运行。相对路径是相对于正在运行的工作目录的。

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

执行结果：

```text
bytes=5 LF=1
```

## 行动和责任

`read_to_end` 打开和关闭文件，但 Buffer 释放是调用者的责任。关闭成功和读取失败路径。由于它是多个字节的数据，因此我们不假设它是以NUL结尾的字符串。

LF 行数和人们想象的行数并不总是一样的。如果最后一行不包含 LF，则不包含在该程序的 LF 计数中。还可以通过重命名文件来检查读取失败。具体错误号可能会因您的环境而异。

## 扩展练习和评论

如果要处理非常大的文件，请将其替换为固定大小的数组和`io_read`迭代。仅处理正返回范围并在 0 处终止。您可以累积字节计数和 LF 计数，而无需将整个文件保留在内存中。如果您自己打开它，它还会关闭任何退出路径上的描述符。

[参见 fs 和 io](/docs/zh/stdlib/files-io) · [参见Buffer](/docs/zh/stdlib/buffer)

## 遵循处理流程

1. 初始化 bin Buffer。还没有文件内容。
2. read_to_end 读取文件并增加所需空间。
3. 如果读取成功，则检查data.len范围内的字节。
4. 每当我们遇到LF的字节值10时，我们就会递增lines。
5. 打印结果并发布Buffer。

data.cap 为预留存储空间，data.len 为有效数据长度。如果将重复条件更改为cap，则会读取文件中没有的字节，因此请使用len。 count 是本次调用read_to_end 添加的字节数。此示例以空的 Buffer 开头，因此 count 和 data.len 相等。

## 更改输入以检查

|input.txt 内容| bytes | LF |原因|
| --- | --- | --- | --- |
|空文件| 0 | 0 |没有字节可读取|
| `Wave` | 4 | 0 |最后没有换行|
| `Wave` + LF | 5 | 1 |截至最后的数据LF|
| `A` + LF + `B` + LF | 4 | 2 |数二LF|
| `Wave` + CRLF | 6 | 1 |CR也是1个字节，但是只计算LF。|

要将最后一行没有 LF 的文件算作一行，如果文件不为空且最后一个字节不为 10，则行数加 1。必须先检查 data.len 是否为 0，然后才能访问最后一个元素。

## 扩展到大文件

目前将全部内容存档的方法方便以后重新阅读或检索数据。如果您只需要字节数和 LF 计数，则重用固定大小的缓冲区是有意义的。

在 [将文件读入固定大小的缓冲区](/docs/zh/stdlib/files-io) 示例中的 io_read 循环中，只需按返回的字节数计算 LF 即可。它可以使用等于缓冲区大小的内存进行处理，而不是文件的整个长度。如果读取返回0，则结束；如果为负，则表示错误。
