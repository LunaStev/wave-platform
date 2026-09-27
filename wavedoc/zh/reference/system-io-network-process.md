---
translation_set_id: system-io
path: reference/system-io-network-process
locale: zh
group: stdlib
group_order: 1
order: 15
title: 系统功能及流程
summary: 描述父 API 和 OS 接口的边界和进程生命周期。
---

## 按功能分类的文档

读取[fs 和 io](/docs/zh/stdlib/files-io)来处理文件，读取[TCP](/docs/zh/stdlib/tcp)来链接，读取[resolver](/docs/zh/stdlib/resolution)来查找地址。以下是流程和较低级别OS访问规则。

## 工艺基础知识API

```text
std::process::core
proc_exit(code: i32) -> !
proc_getpid() -> i64
proc_getppid() -> i64
proc_execve(path: str, argv: ptr<ptr<i8>>, envp: ptr<ptr<i8>>) -> i64
proc_waitpid_raw(pid: i64, status: ptr<i32>, options: i32) -> i64
```

`proc_exit` 不返回呼叫点。请在关机前执行任何必要的文件/内存清理。 `proc_execve` 与典型的子创建函数不同，因为如果成功，它将替换现有的过程映像。 raw argv/envp 必须为每个字符串的NUL 终止做好准备，null 指针指示结束。

`std::process::spawn`中的spawn函数处理创建结果，await函数处理子进程的退出状态。程序的成功创建和成功终止是两件不同的事情。创建管道时，父管道和子管道必须关闭未使用的一端，以便传递 EOF。如果等待子进程退出而不读取捕获管道，则缓冲区可能会填满并互相等待。

## 可移植性和低级方法

fork/exec、文件描述符和Windows句柄不是同一个OS函数。验证对所选目标的支持并将 unsupported 视为正常故障路径。 `std::sys` 是OS 特定的接口，不会重用其他 OS 的数字标志和布局。

直接链接外部C库时，请阅读[FFI](/docs/zh/language/modules-imports-and-ffi)。无需任意声明函数libc来使用父函数stdAPI。首先检查[目标和链接环境](/docs/zh/whale/build-link-targets)。
