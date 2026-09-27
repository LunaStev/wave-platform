---
translation_set_id: ecosystem
path: whale/ecosystem
locale: zh
group: whale
group_order: 1
order: 2
title: 工具链组件
summary: 描述单独的低级工具链Whale的作用，以及Wave生态系统的组件边界。
---

## Whale伊朗

Whale 是一个处理组装和中间表示的低级工具链。处理程序集、对象、链接和中间表示的组件被设计为可在 Wave 和其他本机代码生成工具中重用。

Whale 不是整个Wave 开发环境的名称。每个项目的职责划分如下：

|项目|责任|
| --- | --- |
| `wavec` |Wave 检查源代码并创建可执行文件。|
| Vex |管理 Wave 包、manifest、依赖关系图、lockfile 和包构建。|
| Whale |提供独立的assembler、object、linker、IR组件。|
| Wave `std` |运行时和系统 API 作为 Wave 源模块提供。|

## 组件

Whale workspace 由四个主要库区组成：

- `assembler`：标记化、AMD64 解析/编码、section、symbol 和 relocation
- `object`：目标文件模型和 ELF64 writer
- `linker`：链路层
- `ir`：Whale IR 类型、builder、输出、验证和可选 frontend socket

`whale` 可执行文件为该区域提供命令 `asm`、`object`、`link` 和 `ir`。

## 工具边框

Wave 程序构建为 `wavec`。直接处理程序集、对象文件和IR时，请使用`whale`命令。

安装Whale不会改变`wavec`的构建方法。 Vex使用`wavec`构建Wave包，Whale直接在处理低级工件的任务中运行它。

## 可交付成果验证

将 Whale 工件链接到构建过程时，请确保 object format 和目标 architecture 匹配。 symbol和relocation可以使用`readelf`和`objdump`等独立工具进行检查。使用 IR socket 的构建必须使用与 Whale 来自同一生产商的 socket schema。
