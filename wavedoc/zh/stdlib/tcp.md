---
translation_set_id: stdlib-tcp
path: stdlib/tcp
locale: zh
group: stdlib
group_order: 1
order: 9
title: net.tcp:连接和传输
summary: TCP 描述结果结构、部分传输和断开责任。
---

## 检查连接结果

```text
std::net::tcp
tcp_connect_addr(addr: SocketAddr) -> NetResult<TcpStream>
tcp_bind_loopback(port: u16) -> NetResult<TcpListener>
tcp_accept(listener: TcpListener) -> NetResult<TcpStream>
tcp_read(stream: TcpStream, buf: ptr<u8>, size: i64) -> i64
tcp_write_all(stream: TcpStream, buf: ptr<u8>, size: i64) -> i64
tcp_read_exact(stream: TcpStream, buf: ptr<u8>, size: i64) -> i64
tcp_close(stream: TcpStream) -> NetError
tcp_close_listener(listener: TcpListener) -> NetError
```

链接函数`NetResult<T>`有`ok`、`value`、`error`。仅在成功时才使用value。 `NetError`包含标准化错误分类和native_code，可以通过`error.kind == NET_ERROR_NONE`检查是否成功。该常数取自`std::net::error`。

连接的流和通过accept接收的流必须分别关闭。关闭监听器并不意味着关闭它已经接受的所有流。不要关闭该结构的每个副本。

## TCP 不保留消息边界

不要以为你曾经写过的所有内容在你读过之后都会回想起你。协议必须由长度前缀、分隔符或固定长度分隔。固定长度由`tcp_read_exact`处理，未知长度的流由读取迭代和EOF处理处理。

负数是错误，正长度读取中，0是另一端的末尾。 write_all 部分数据可能在失败前已传输完毕。如果从头开始重新传输相同的内容，则可能会出现重复。

## 等待和时间限制

默认的blocking功能可以让对方等待较长时间。对于需要时间限制的节目，请选择`tcp_connect_addr_timeout`、`tcp_read_timeout`、`tcp_write_timeout`系列。检查毫秒因素和错误后果，同时设计对方不响应的路线。异步使用[task](/docs/zh/stdlib/task)以及相应的异步网络API。

[地址查找](/docs/zh/stdlib/resolution) · [本地客户端完成](/docs/zh/practice/tcp-client)

## 选择阅读功能的标准

|需要采取行动|功能|要检查的值|
| --- | --- | --- |
|读取一些到达的数据| `tcp_read` |负错误、零终止、正字节数|
|读取一定长度| `tcp_read_exact` |请求长度是否已完整读取？|
|将给定内容发送到末尾| `tcp_write_all` |请求长度和返回值|
|等待时间限制|timeout系列|返回结果及错误timeout|

对于四字节长度字段后跟正文的协议，首先读取完整的长度字段，检查长度是否超过允许的最大值，然后为正文分配空间。不要从对等方提供的未经验证的长度中分配大量内存。

## 何时关闭连接

即使read函数返回0，本地流资源仍然保留。阅读完毕后，拨打tcp_close。如果即使在传输错误之后也使用相同的清理路径，您就不会忘记关闭正常和失败的路径。

监听器是接收新连接的资源，流是已经连接的通信资源。创建服务器时，您可以管理这两种类型。每次accept成功时，都会创建一个新流，因此我们在处理完它后将其关闭，并在服务器接受迭代完成时关闭监听器。

## 自己尝试一下

[本地TCP客户实践](/docs/zh/practice/tcp-client) 包含服务器和客户端的完整代码。您可以通过比较服务器首先运行、没有服务器以及端口号不同的情况来区分地址查找和连接失败。
