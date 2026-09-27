---
translation_set_id: practice-tcp-client
path: practice/tcp-client
locale: zh
group: practice
group_order: 4
order: 4
title: 项目：本地TCP客户端
summary: 连接数字地址查找、连接失败、传输和关闭。
---

## 准备工作和范围

该程序连接到本机OS上的`127.0.0.1:8080`，发送`ping`和LF，然后退出。您将需要在单独的终端中安装测试服务器。如果存在 Python 3，则下一台服务器将从一个本地连接接收并显示最多 5 个字节。

```python
import socket
with socket.socket() as server:
    server.bind(("127.0.0.1", 8080))
    server.listen(1)
    connection, address = server.accept()
    with connection:
        message = b""
        while len(message) < 5:
            chunk = connection.recv(5 - len(message))
            if not chunk:
                break
            message += chunk
        print(repr(message))
```

先运行服务端，然后将客户端保存为`main.wave`，运行为`wavec run main.wave`。服务器输出`b'ping\n'`。如果该端口已被使用，请同时更改服务器和客户端上的端口。

<!-- wave-example: tcp-client -->
```wave
import("std::net::address")::{
    SocketAddr
};
import("std::net::resolve")::{
    NetResolveResult, net_resolve_tcp
};
import("std::net::error")::{
    NetResult, NetError, NET_ERROR_NONE
};
import("std::net::tcp")::{
    TcpStream, tcp_connect_addr_timeout, tcp_write_all, tcp_close
};

fun main() -> i32 {
    var addresses: array<SocketAddr, 4>;
    var resolved: NetResolveResult = net_resolve_tcp("127.0.0.1", "8080", &addresses[0], 4);
    if (!resolved.ok || resolved.written == 0) {
        println("resolve failed");
        return 1;
    }

    var connected: NetResult<TcpStream> = tcp_connect_addr_timeout(addresses[0], 1000);
    if (!connected.ok) {
        println("connect failed");
        return 2;
    }

    var sent: i64 = tcp_write_all(connected.value, "ping\n" as ptr<u8>, 5);
    var closed: NetError = tcp_close(connected.value);
    if (sent != 5 || closed.kind != NET_ERROR_NONE) {
        return 3;
    }

    println("sent=5");
    return 0;
}
```

客户端成功输出为`sent=5`。如果服务器没有运行，正常的故障处理是打印 `connect failed` 并以 2 退出。

## 理解代码

查找数组的有效范围最大为written。这是一个仅使用第一个地址的小示例。在具有多个候选者的服务中，您需要为每个候选者确定连接策略。无论传输成功还是失败，连接的流都会关闭。连接超时为1000ms，这并不是保证整个传输超时的示例。

TCP 不保留传输边界。这就是为什么编写它以便服务器可以多次运行recv。 [地址查找](/docs/zh/stdlib/resolution) 介绍了如何通过名称查找地址。

## 扩展练习

让服务器发送响应并让客户端读取它。确定响应的长度后，您需要同时处理部分读取、EOF、超时和关闭。

## 两个终端的作用

服务器终端在状态listen等待连接。当客户端连接时，accept返回已连接的套接字，recv循环收集数据。客户端按照地址查找、连接、传输、关闭的顺序进行。

Python 代码是实验室合作伙伴服务器。 Wave 用于方便检查程序传递的字节，不与客户端代码放在同一个文件中。服务器在处理一个连接后终止，因此它也会在再次运行客户端之前重新启动。

## 为什么要发送5个字节

由于`ping`为4字节，LF为1字节，因此传输长度为5。字符串末尾的NUL不包含在消息中。服务器收集 5 个字节然后输出它们，因此即使数据包以多个块到达，也会产生相同的结果。

如果tcp_write_all返回5，则本地传输函数已经处理了请求的字节。这并不确认其他程序已解释并保存该消息。如果协议需要确认处理完成，则服务器发送响应，客户端读取响应。

## 失败之路实践

|改变|结果|学什么|
| --- | --- | --- |
|无需服务器即可运行|连接失败，退出代码 2|即使地址有效，服务器也可能不存在|
|设置服务器和客户端端口不同|连接失败|IP 和地址中的端口必须匹配。|
|将服务器recv的大小更改为1|相同的 5 个字节在一起|TCP 读取单位与消息单位不同|
|将客户端传输分成多次|以相同顺序收到|消息边界由协议确定。|

连接超时1000ms是连接阶段的设置。要限制读写，请在相关步骤中使用timeout函数，如果有适用于整个程序的时间限制，请计算剩余时间并将其传递。
