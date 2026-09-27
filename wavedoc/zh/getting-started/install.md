---
translation_set_id: install
path: getting-started/install
locale: zh
group: getting-started
group_order: 1
order: 2
title: 安装 Wave
summary: 在 Linux、macOS 或 Windows 上安装 Wave，并运行第一个程序。
---

## Linux 和 macOS

在终端中运行以下命令。它会同时安装 Wave 和 Vex 包管理器。

```shell
curl -fsSL https://wave-lang.dev/install.sh | bash -s -- latest
```

安装完成后，打开新的终端并检查版本。

```shell
wavec --version
```

## Windows

在 PowerShell 中运行以下命令。它们会同时安装 Wave 和 Vex 包管理器。

```powershell
irm https://wave-lang.dev/install.ps1 -OutFile install.ps1
powershell -ExecutionPolicy Bypass -File .\install.ps1 -Latest
```

安装完成后，打开新的 PowerShell 窗口并检查版本。

```powershell
wavec --version
vex --version
```

## 首次运行

将以下代码保存为 `main.wave`。

<!-- wave-example: install-stdlib -->
```wave
import("std::string::len")::{
    len
};

fun main() {
    println("Wave: {} bytes", len("Wave"));
}
```

在保存文件的目录中运行程序。

```shell
wavec run main.wave
```

运行结果：

```text
Wave: 4 bytes
```

如果提示找不到标准库，请先安装标准库，然后再次运行程序。

```shell
wavec install std
wavec run main.wave
```

[下一步：第一个程序](/docs/zh/language/program-structure) · [故障排查](/docs/zh/reference/diagnostics)
