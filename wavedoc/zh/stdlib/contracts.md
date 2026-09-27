---
translation_set_id: stdlib-contracts
path: stdlib/contracts
locale: zh
group: stdlib
group_order: 1
order: 2
title: 阅读 API 文档：错误和所有权
summary: 了解参数单元、结果结构、部分成功和资源生命周期。
---

## 阅读声明

以下符号描述了函数声明，而不是整个可执行文件。

```text
io_read(fd: i64, buf: ptr<u8>, len: i64) -> i64
```

`fd`是打开描述符，`buf`是调用者提供的存储空间，`len`是可写入的字节数。 `i64` 并不意味着负长度有效。返回值是实际读取的字节数，而不是请求长度，因此只使用返回的范围。

## 失败表达式因函数而异

|方式|是的|检查方法|
| --- | --- | --- |
|指针或null| `mem_alloc` |null 检查后的内存访问|
|字节数或负数| `io_read` |负误差，0 EOF，正数据|
|状态码| `buffer_push` |与`BUFFER_OK`的误差常数比较|
|成功与价值| `NetResult<T>` |检查`ok`后，使用`value`|
|包括部分进展| `RandomFillResult` |一起检查`ok`、`written`、`error`。|

它只查看错误数量，不会将它们与其他模块中的常量进行比较。例如，错误号 env 和 OS errno 不是同一个系统。 WASI 中的原始错误不应解释为 Linux errno。

## 拥有和租赁

- **Owned**：当获取已分配的内存、打开的文件或打开的套接字时，负责调用相应的释放/关闭。
- **借用**：传递给函数的字节view或缓冲区引用现有内存。如果函数未指定它接收所有权，则调用者保留控制权。
- **输出参数**：传递一个有效的存储空间，结果可以写入到接收它的函数中，例如`out_value: ptr<T>`。确保合同规定结果仅在成功时才有效。

复制 Buffer 结构可以使两个副本都指向相同的分配。不要单独释放每个副本。释放或重新分配分配后，借用的指针将变得无效。字符串文字不是可写缓冲区。

## 失败并不意味着恢复到以前的状态

`io_write_all` 写入一些字节后可能会失败。已经从外部写入的字节将不会被返回。另一方面，读取bytes的checkedcursor，如果失败则保留位置和输出值。这些差异由 API 指定。

如果您还想练习处理故障，请继续[文件阅读器](/docs/zh/practice/file-reader)和[二进制消息](/docs/zh/practice/binary-message)。
