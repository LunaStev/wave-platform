---
translation_set_id: standard-library
path: reference/standard-library
locale: zh
group: stdlib
group_order: 1
order: 1
title: 标准库指南
summary: 如何找到适合您目的的模块并阅读该函数的错误和所有权规则。
---

## 找到您需要的功能

标准库是import，路径为`std::module::file`。即使名称相似，函数也可能以不同的方式返回错误。首先阅读[API 如何阅读](/docs/zh/stdlib/contracts)，然后转到下表中您需要的模块。

|我想做的事|文件|主要import|
| --- | --- | --- |
|字符串长度/比较/搜索| [string](/docs/zh/reference/string-and-bytes) | `std::string::len`, `cmp`, `find`, `trim` |
|内存分配/复制/大小| [mem](/docs/zh/reference/memory-and-buffer) | `std::mem::alloc`, `ops`, `layout` |
|不同大小的字节列表| [buffer](/docs/zh/stdlib/buffer) | `std::buffer::alloc`, `read`, `write` |
|二进制读/写| [bytes](/docs/zh/stdlib/bytes) | `std::bytes::types`, `cursor`, `leb128` |
|文件·描述符I/O| [fs 和 io](/docs/zh/stdlib/files-io) | `std::fs::file`, `std::io::fd` |
|路径组合/环境设置| [path 和 env](/docs/zh/stdlib/path-env) | `std::path::copy`, `std::env::environ` |
|时间测量/等待| [time](/docs/zh/stdlib/time) | `std::time::duration`, `clock`, `sleep` |
|数字地址/姓名列表搜索| [net.resolve](/docs/zh/stdlib/resolution) | `std::net::resolve`, `resolve_table` |
|TCP 连接/传输| [net.tcp](/docs/zh/stdlib/tcp) | `std::net::tcp`, `address`, `error` |
|OS 随机数| [random](/docs/zh/stdlib/random) | `std::random::fill` |
|流程·OS边界| [系统功能](/docs/zh/reference/system-io-network-process) | `std::process::core`, `spawn`, `std::sys` |
|异步任务执行| [task](/docs/zh/stdlib/task) | `std::task` |
|数学/诊断助理| [math 和 debug](/docs/zh/stdlib/math-debug) | `std::math::int`, `float`, `std::debug::core` |

## 首次使用示例

下面的程序使用std中的一个函数，无需下载单独的包。将其保存为`main.wave`并以`wavec run main.wave`运行。

<!-- wave-example: library-import -->
```wave
import("std::string::len")::{
    len
};

fun main() {
    println("{}", len("Wave"));
}
```

结果是`4`。要逐步理解同一示例，请阅读[弦乐](/docs/zh/language/strings)。

## 兼容编译器std

使用`wavec print std-path`确认所选路径。当从另一个结账中使用std时，请将路径指定为`wavec --std-root /absolute/path/to/std check main.wave`。如果指定的路径无效或不兼容，则会显示错误。

## 平台边框

区分字符串/字节等计算函数和文件/套接字等OS函数。识别目标并不能保证提供所有主机API。一起阅读 [支持对象](/docs/zh/whale/build-link-targets) 和每个 API 的平台条目。 `std::sys`是一个较低级别的接口，可移植程序将首先使用较高级别的模块。
