---
translation_set_id: stdlib-resolution
path: stdlib/resolution
locale: zh
group: stdlib
group_order: 1
order: 8
title: net.resolve: 地址解析
summary: 描述数字地址、调用者拥有的姓名列表和结果的截断。
---

## 选择查询方式

Linux、resolver 的默认值采用数字 IPv4/IPv6 地址和数字端口。主机名或服务名称不会隐式传递到系统DNS。例如，`127.0.0.1`和`8080`是数字输入，`example.com`和`http`是名称。

要使用明确提供的名称列表，请使用`std::net::resolve_table`。您可以直接连接到数字地址，也可以通过在姓名列表中注册地址来进行搜索。

## 基本查找

```text
std::net::resolve
net_resolve_tcp(host: str, service: str, output: ptr<SocketAddr>, capacity: i64) -> NetResolveResult
net_resolve_udp(host: str, service: str, output: ptr<SocketAddr>, capacity: i64) -> NetResolveResult
net_resolve_host(host: str, service: str, output: ptr<SocketAddr>, capacity: i64) -> NetResolveResult
```

输出数组由调用者提供。 capacity的单位不是字节，而是`SocketAddr`元素的数量。您只能搜索capacity=0 和null output 的计数。内部查找存储在返回之前被清理，并且不取得输出数组的所有权。

|结果字段|意义|
| --- | --- |
| `ok` |查询是否成功|
| `count` |支持的地址结果总数|
| `written` |实际写入输出数组的元素数量|
| `truncated` |输出能力是否小于总结果？|
| `error.kind` |标准化分类，例如 INVALID、NOT_FOUND、UNSUPPORTED|
| `error.native_code` |用于确定原因的本机代码|

即使成功，written也可能为0。在使用output[0]之前检查`ok && written > 0`。

## 您自己提供的姓名列表

```text
net_resolve_from_table(host: str, service: str, protocol: i32,
    entries: ptr<NetResolveEntry>, entry_count: i64,
    output: ptr<SocketAddr>, capacity: i64) -> NetResolveResult
```

`std::net::resolve_table` 中的`NetResolveEntry` 有 host、service、protocol、address。主机名ASCII不区分大小写，准确比较服务名。 TCP/UDP/对于所有协议查询，都使用模块的常量。项目和字符串归调用者所有，调用后不会保留。

匹配项按输入顺序复制并保留重复项。输出空间和输入存储空间不得重叠。如果该名称不在列表中，则为 NOT_FOUND，并且不会使用外部 DNS 重试。

[TCP 使用方法](/docs/zh/stdlib/tcp) · [TCP 客户实践](/docs/zh/practice/tcp-client)

## 查找数字地址

准备一个数组来保存查询结果，检查是否记录了地址。下面的示例不需要服务器，因为地址查找本身不会创建连接。

<!-- wave-example: book-resolve-numeric -->
```wave
import("std::net::address")::{
    SocketAddr
};
import("std::net::resolve")::{
    NetResolveResult,
    net_resolve_tcp
};

fun main() -> i32 {
    var addresses: array<SocketAddr, 4>;
    var result: NetResolveResult = net_resolve_tcp(
        "127.0.0.1",
        "8080",
        &addresses[0],
        4
    );

    if (!result.ok || result.written == 0) {
        println("address unavailable");
        return 1;
    }

    println("address ready");
    return 0;
}
```

执行结果：

```text
address ready
```

要实际连接，请将addresses[0]传递给tcp_connect_addr系列函数。即使您获得了地址，您也可以在连接阶段判断服务器是否正在该端口上运行。

## 按姓名列表搜索

自定义地址列表对于基于配置文件或程序的服务列表链接名称非常有用。以下示例将名称 api 映射到本地 8080 端口。

<!-- wave-example: book-resolve-table -->
```wave
import("std::net::address")::{
    SocketAddr,
    socket_addr_from_v4,
    socket_addr_v4_loopback
};
import("std::net::resolve")::{
    NetResolveResult
};
import("std::net::resolve_table")::{
    NetResolveEntry,
    RESOLVE_TCP,
    net_resolve_from_table
};

fun main() -> i32 {
    var entries: array<NetResolveEntry, 1>;
    entries[0] = NetResolveEntry {
        host: "api",
        service: "http",
        protocol: RESOLVE_TCP,
        address: socket_addr_from_v4(socket_addr_v4_loopback(8080))
    };

    var addresses: array<SocketAddr, 2>;
    var result: NetResolveResult = net_resolve_from_table(
        "API",
        "http",
        RESOLVE_TCP,
        &entries[0],
        1,
        &addresses[0],
        2
    );

    if (!result.ok || result.written != 1) {
        return 1;
    }

    println("matched={}", result.written);
    return 0;
}
```

执行结果：

```text
matched=1
```

API 和 api 在主机名比较中匹配。在此示例中，http 和 HTTP 不匹配，因为它们在服务名称比较中不同。如果您使用相同的名称注册多个地址，结果将按输入的顺序显示。如果您只收到一个小数组的一部分，则只读取written，如果您需要整个列表，请为count准备空间并再次搜索。
