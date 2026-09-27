---
translation_set_id: vex-package-manager
path: whale/vex-package-manager
locale: zh
group: whale
group_order: 1
order: 4
title: Vex 包管理器
summary: 描述基于 manifest Wave 项目、Git·路径依赖项、lockfile、离线构建和 wavec 边界。
---

## 角色

Vex 是Wave 的包管理器和构建工具。 Vex 在 `wavec` 之上运行。 Vex负责项目结构和依赖分析，`wavec`负责编译器标志和编译管道。

Vex 命令基于manifest。 `vex build`、`vex check` 和 `vex run` 故意不接收 raw `wavec` 标志。

## 创建一个包

```shell
vex init
vex init --lib
```

应用程序使用`src/main.wave`，库使用`src/lib.wave`。包根结构如下：

```text
my_project/
├── src/
│   └── main.wave
├── vex.ws
├── vex.lock
└── .vex/
    └── deps/
```

`vex.ws` 变为 manifest。 Vex 在 `.wson` 扩展中不使用 manifest。

```wson
{
    name = "my_project",
    version = 0.1.0,
    lib = false,
    description = "my_project Project",
    author = "unknown",
    license = "Unknown",
    dependencies = []
}
```

## 构建命令

```shell
vex build [--target <triple>] [--release] [--dry-run] [--locked] [--offline]
vex check [--target <triple>] [--release] [--dry-run] [--locked] [--offline]
vex run   [--target <triple>] [--release] [--dry-run] [--locked] [--offline] [-- <args...>]
```

Vex 的选项保持较小。如果您需要依赖于编译器的控制，例如emit、linker、CPU、ABI或debug，请直接使用`wavec`。当需要使用特定编译器时，设置`VEX_WAVEC=/path/to/wavec`。

`Resolving`、`Fetching`、`Compiling`、`Checking`、`Running`、`Finished`等进度步骤在stderr中输出，程序输出维持在stdout。

## Git 中央依赖

Vex 依赖项指定为本地 `path` 或 Git URL。依赖项只能使用这两种方法之一。

```wson
{
    name = "app",
    version = 0.1.0,
    dependencies = [
        { name = "local_math", path = "../local_math" },
        { name = "remote_math", git = "https://github.com/example/math.git", tag = "v0.1.0" }
    ]
}
```

Git 依赖项最多只能指定 `branch`、`tag` 或 `rev` 之一。每个依赖根必须有自己的`vex.ws`。 Vex递归地解析依赖关系manifest，拒绝冲突的包标识，并将托管的Gitcheckout存储在`.vex/deps/<name>`中。

## lockfile 合同

模式v2`vex.lock`记录了整个传递依赖图和确切的Gitcommit。使用 manifest 进行承诺。使用相同的manifest和有效的lockfile将选择相同的依赖关系图，而无需再次遵循branch或tag。

需要依赖的命令会自动解释，可以通过以下命令提前准备好。

```shell
vex fetch
vex update
vex update math shared_core
```

`vex update` 更新所有 Git 包，或仅更新指定的包和受影响的转换图。不相关的锁定包将保持commit处于选中状态。

## locked 和 offline 工作流程

`--locked` 禁止创建和修改`vex.lock`。如果文件不存在、具有不受支持的架构或与manifest图不匹配，它将失败。已固定到lockfile，需要时可以导入commit。

`--offline` 禁止所有Git 网络操作。所需的checkout和commit应该已经存在于本地。

```shell
vex fetch --locked
vex build --locked --offline
```

这两个命令是严格的CI工作流程。正好在网络可用的时候准备一个锁定的commit，然后在不改变网络或lockfile的情况下编译它。 dry-run 不导入依赖项或重写lockfile。

## 编译器设置和信息

```shell
vex info
vex setup wavec
vex setup wavec --version <version>
vex --version
```

Vex 在实际构建之前验证 `wavec` dry-run JSON 架构。未实现所需架构的编译器将不会以未知计划执行，并会因兼容性错误而拒绝它。
