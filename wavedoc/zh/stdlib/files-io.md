---
translation_set_id: stdlib-files-io
path: stdlib/files-io
locale: zh
group: stdlib
group_order: 1
order: 7
title: fs 和 io：文件和字节传输
summary: 描述文件生命周期、完全读取、部分传输和故障后状态。
---

## 选择文件API

`std::fs::file` 的便利功能接收路径并执行必要的打开和关闭。返回描述符的函数必须由调用者关闭。

|声明|成功结果及注意事项|
| --- | --- |
| `open_read(path: str) -> i64` |打开描述符。负数是错误|
| `create(path: str) -> i64` |创建或删除现有文件内容。返回所属描述符|
| `open_append(path: str) -> i64` |打开或创建以进行添加|
| `size(path: str) -> i64` |字节数。负数是错误|
| `read_into(path: str, dst: ptr<u8>, dst_cap: i64) -> i64` |整个文件中的字节数。容量不足是一个错误|
| `read_to_end(path: str, dst_buffer: ptr<Buffer>) -> i64` |在现有Buffer之后添加文件并返回额外金额|
| `write(path: str, src: ptr<u8>, len: i64) -> i64` |写入替换现有内容的字节数|
| `append(path: str, src: ptr<u8>, len: i64) -> i64` |添加到末尾的字节数|
| `remove(path: str) -> i64` |移除状态。失败是消极的|

false 或 `exists(path)` 单独无法区分丢失文件和权限错误。请务必检查实际的打开结果，因为在检查存在和打开之间状态可能会发生变化。

## 低电平I/O

```text
std::io::fd
io_read(fd: i64, buf: ptr<u8>, len: i64) -> i64
io_write(fd: i64, buf: ptr<u8>, len: i64) -> i64
io_read_exact(fd: i64, buf: ptr<u8>, len: i64) -> i64
io_write_all(fd: i64, buf: ptr<u8>, len: i64) -> i64
io_close(fd: i64) -> i64
```

`io_read`的正结果是读取的字节数，正长度请求中的0为EOF。 `io_write` 可能写得比要求的少。如果需要完整传输，请使用exact/all功能。尽管如此，我们并不假设故障会恢复外部状态，因为在错误发生之前可能已经发生了一些传输。

如果在所需长度之前遇到 EOF，则`io_read_exact` 是错误。如果缓冲区已满，并且某些字节可能已被写入，则`read_into`返回`IO_ERR_NO_SPACE`。读取函数不会自动将 NUL 添加到字符串末尾。

## Buffer 和错误处理

失败时，`read_to_end`恢复原始len，但其容量和数据地址可能已更改。无论成功还是失败，调用者都必须释放Buffer。文件写入 API 不保证原子文件替换。

从[文件阅读练习](/docs/zh/practice/file-reader)开始，可以运行import的程序进行发布。考虑Linux/macOS/Windows/FreeBSD中的路径/权限差异以及WASI中的可访问目录限制。它不会直接将描述符值解释为另一个OS的原始句柄。

## 将大文件读入小缓冲区

不需要将整个文件放入内存的操作可以使用固定缓冲区和读取迭代来处理。以下程序打印input.txt的内容并计算读取的总字节数。在输入文件中保存一个`Wave`和LF。

<!-- wave-example: book-file-stream -->
```wave
import("std::fs::file")::{open_read};
import("std::io::fd")::{io_read, io_write_all, io_close};
import("std::io::consts")::{IO_STDOUT_FD};

fun main() -> i32 {
    var descriptor: i64 = open_read("input.txt");

    if (descriptor < 0) {
        return 1;
    }

    var buffer: array<u8, 4>;
    var total: i64 = 0;

    while (true) {
        var count: i64 = io_read(descriptor, &buffer[0], 4);

        if (count < 0) {
            io_close(descriptor);
            return 2;
        }

        if (count == 0) {
            break;
        }

        if (io_write_all(IO_STDOUT_FD, &buffer[0], count) < 0) {
            io_close(descriptor);
            return 3;
        }

        total += count;
    }

    if (io_close(descriptor) < 0) {
        return 4;
    }

    println("bytes={}", total);
    return 0;
}
```

执行结果：

```text
Wave
bytes=5
```

缓冲区容量为4，但最后读取的可能是1个字节。我们总是将实际的 count 传递给输出。如果写入整个数组，甚至可以输出旧的、未读的字节。

该程序打开了input.txt的descriptor，因此将其关闭。标准输出并不是该函数中新获取的资源，因此在示例结束时不会随意关闭。

## 读取已满，容量不足

read_into 接收保存整个文件的固定存储空间。如果空间不足，它会默默地截断并返回NO_SPACE，但不会成功。

<!-- wave-example: book-file-capacity -->
```wave
import("std::fs::file")::{read_into};
import("std::io::consts")::{IO_ERR_NO_SPACE};

fun main() {
    var data: array<u8, 2>;
    var status: i64 = read_into("input.txt", &data[0], 2);

    if (status == IO_ERR_NO_SPACE) {
        println("destination too small");
    }
}
```

执行结果：

```text
destination too small
```

它并不假设失败的读取根本没有改变目标字节。不要将其用作完成的文件内容，准备更大的存储库或选择streaming方法。即使先查询大小，实际读取的结果才是最终判断，因为查询和读取之间文件可能会发生变化。

## 写入文件和附加文件之间的区别

write 替换现有内容，并在末尾添加 append。可以直接从字符串长度获取要存储的字节数并传递。字符串末尾的NUL通常不包含在文本文件内容中。

<!-- wave-example: book-file-write-append -->
```wave
import("std::fs::file")::{write, append, size, remove};
import("std::string::len")::{len};

fun main() -> i32 {
    var first: str = "Wave";
    var second: str = " study";

    if (write("output.txt", first as ptr<u8>, len(first) as i64) < 0) {
        return 1;
    }

    if (append("output.txt", second as ptr<u8>, len(second) as i64) < 0) {
        return 2;
    }

    println("bytes={}", size("output.txt"));

    if (remove("output.txt") < 0) {
        return 3;
    }

    return 0;
}
```

执行结果：

```text
bytes=10
```

此示例在工作目录中创建、替换并最终删除 output.txt。从不存在任何文件的练习目录运行。实际的编辑器或存储程序可能需要单独的存储策略，例如临时文件和替换。

## API选型表

|情况|选择|
| --- | --- |
|将整个小文件读入固定缓冲区| read_into |
|在不知道大小的情况下保留全部内容|read_to_end 和 Buffer|
|按顺序处理内容而不是完整存储内容|open_read + io_read 重复|
|读取固定长度的记录| io_read_exact |
|传输整个字节串| io_write_all |
|处理已经打开的文件|fd 函数代替路径函数|

选择功能后，检查发生故障时缓冲区、文件位置和外部数据如何变化。
