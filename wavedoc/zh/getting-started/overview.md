---
translation_set_id: overview
path: getting-started/overview
locale: zh
group: getting-started
group_order: 1
order: 1
title: Wave 文档和学习指南
summary: 一步步学习Wave，从安装到实际程序，查找语言规则和标准库API。
---

## 通过本指南学习 Wave

学习编写Wave源代码，编译并运行它，并检查结果。如果您是编程新手，请按照以下顺序操作。如果您了解另一种语言，请运行每一章的示例，并将其规则和边界情况与您已经知道的进行比较。

## 学习路径

|步骤|章|你将学到什么|
| --- | --- | --- |
|设置| [安装](/docs/zh/getting-started/install) |准备编译器和标准库并验证它们是否运行|
| 1 | [你的第一个程序](/docs/zh/language/program-structure) |创建、检查和运行源文件并了解退出代码|
| 2 | [变量和类型](/docs/zh/language/declarations-and-types) |存储值并选择具有所需范围的类型|
| 3 | [运算符和转换](/docs/zh/language/expressions-and-operators) |解释评估顺序和类型转换的结果|
| 4 | [条件和循环](/docs/zh/language/control-flow) |根据条件进行分支并使用循环处理数据|
| 5 | [函数](/docs/zh/language/functions-and-generics) |将重复的操作提取到函数中|
| 6 | [数组](/docs/zh/language/arrays) |按索引访问元素并迭代数组|
| 7 | [弦乐](/docs/zh/language/strings) |区分字符和字节并了解转义符和字符串长度|
| 8 | [结构和变体](/docs/zh/language/structures-enums-and-aliases) |对相关数据进行分组并表示成功和失败|
| 9 | [指针和生命周期](/docs/zh/language/explicit-memory-type-model) |通过其地址修改原始值并管理其生命周期|
| 10 | [动态内存](/docs/zh/language/allocation) |处理分配失败并释放内存|
| 11 | [模块和泛型](/docs/zh/language/modules-imports-and-ffi) |跨文件拆分代码并重用不同类型的函数|
| 12 | [错误处理](/docs/zh/language/errors) |检查结果并在失败时清理资源|
| 13 | [异步代码简介](/docs/zh/language/async-and-never) |创建一个 Future 并等待它完成|

## 将您的知识付诸实践

核心章节之后，构建一个[输入计算器](/docs/zh/practice/input-calculator)、一个[文件阅读器](/docs/zh/practice/file-reader)、一个[二进制消息](/docs/zh/practice/binary-message)和一个[TCP客户端](/docs/zh/practice/tcp-client)。在每个项目中测试成功的输入和失败的案例。

## 三个文档选项卡

- **Wave**：有指导的语言课程和按顺序进行的实践项目。
- **[标准库](/docs/zh/stdlib)**：每个模块的API、返回值、错误、所有权规则和平台要求。
- **[Whale](/docs/zh/whale)**：构建和链接、包管理、命令使用和低级工具链。

示例将完整的程序与属于函数内部的片段区分开来。在终端中运行`wavec`命令并将`wave`代码块保存在`.wave`文件中。输入和输出分开显示；读取标准输入的示例指定要输入的内容。

## 当你陷入困境时

使用[疑难解答](/docs/zh/reference/diagnostics)来区分安装、源检查、链接和执行问题。在[语法快速参考](/docs/zh/reference/syntax-quick-reference)中查找语言规则，在[编译器参考](/docs/zh/getting-started/compiler)中查找命令，在[标准库指南](/docs/zh/reference/standard-library)中查找API。
