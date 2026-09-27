---
translation_set_id: compiler
path: getting-started/compiler
locale: zh
group: getting-started
group_order: 1
order: 3
title: 编译命令参考
summary: wavec 描述命令、构建管道、输出、目标、诊断、依赖项链接和工具查询。
---

## 命令模型

`wavec` 是编译器CLI。它直接编译各个输入，为工具提供编译器支持信息，并管理已安装的标准库源。

```text
wavec [global-options] <command> [command-options]
```

|命令|使用|
| --- | --- |
| `wavec build <input...>` |根据标志，它执行检查、代码生成、链接或执行管道。|
| `wavec check <file>` |`build <file> --emit=check` 的昵称。|
| `wavec run <file> [-- <args...>]` |它是`build <file> --run`的别名，`--`后面的参数被传递给程序。|
| `wavec print <item>` |查询目标和工具链支持信息。|
| `wavec install std` |安装标准库。|
| `wavec update std` |更新已安装的标准库。|
| `wavec --version` |打印安装的版本信息。|

命令和选项的完整列表可以在`wavec --help`找到。

## 构建、测试和运行

```shell
wavec build main.wave
wavec check main.wave
wavec run main.wave -- first-argument second-argument
```

`build` 默认创建一个可执行文件。 `check` 完成前端检查后停止。 `run` 需要二进制输出，不能与共享库构建一起使用。

使用 `--dry-run` 验证请求并确定执行哪些步骤，而无需编译、链接或执行。

```shell
wavec build main.wave --target riscv64-unknown-linux-gnu --dry-run
wavec build main.wave --dry-run --error-format=json
```

JSON 格式是Vex 等构建工具使用的稳定、统一的接口。

## emit 和输入类型

```shell
wavec build main.wave --emit=ast
wavec build main.wave --emit=ir
wavec build main.wave --emit=bc
wavec build main.wave --emit=asm
wavec build main.wave --emit=obj -o main.o
wavec build main.wave --emit=bin -o app
```

输出emit类型为`ast`、`ir`、`bc`、`asm`、`obj`、`bin`。 `check` 为控制模式，必须单独使用。您可以指定管道接受的多种输出类型，以逗号分隔。

输入类型为`wave`、`ir`、`bc`、`asm`、`obj`、`archive`。 `--input-type=<kind>` 强制指定所有输入的类型。仅链接 object 或 archive 输入时，请使用二进制文件 emit 和 `--link-only`。

```shell
wavec build module.o --input-type=obj --link-only --emit=bin -o app
```

## 输出位置

|选项|效果|
| --- | --- |
| `-o <file>` |指定主输出路径。|
| `--out-dir <dir>` |emit 将输出放置在指定目录中。|
| `--target-dir <dir>` |指定中间和主要可交付路由。|

## 优化和诊断输出

```shell
wavec -O2 build main.wave
wavec --debug-wave=tokens,ast build main.wave
```

优化步骤为`-O0`、`-O1`、`-O2`、`-O3`、`-Os`、`-Oz`、`-Ofast`。 `--debug-wave`后面可以跟`tokens`、`ast`、`ir`、`mc`、`hex`、`all`，多个步骤可以用逗号组合。

## 原生链接

```shell
wavec --link=m -L ./lib build main.wave
wavec build main.wave --shared -o libexample.so
wavec build main.wave --static -o app
wavec build main.wave --pie -o app
```

`--link=<lib>` 添加本机库，`-L <path>` 添加搜索路径。根据兼容性规则，链接模式使用`--shared`、`--static`、`--pie`、`--no-pie`。

后端和链接器控制选项包括：

- `--target`, `--cpu`, `--features`, `--abi`, `--sysroot`
- `-C linker=<path>` 和 `-C link-arg=<arg>`
- `-C link-sysroot=<path>` 和 `-C relocation-model=<model>`
- `-C no-default-libs`

独立输出（例如内核）使用 `--freestanding` 以及适合环境的 `--entry`、`--linker-script` 和 `--no-start-files` 设置。

## 外部包解读

```shell
wavec --dep-root .vex/deps build main.wave
wavec --dep math=/opt/wave-deps/math build main.wave
```

`--dep-root <dir>` 添加一个根来查找外部`package::module` import。 `--dep <name>=<path>` 将包名固定到一个目录。这是编译器集成点，项目manifest、依赖项下载和lockfile 由Vex 处理。

## 支持功能查询

使用目标或输出类型的工具可以通过`wavec print`查询支持信息。

```shell
wavec print host-target
wavec print target-spec --format=json
wavec print supported-targets
wavec print supported-input-types
wavec print supported-emit-kinds
wavec print supported-print-items
wavec print cpu-list --target riscv64-unknown-linux-gnu
wavec print target-features --target riscv64-unknown-linux-gnu
wavec print default-linker
wavec print sysroot
wavec print std-path
wavec print dep-search-paths
```

您还可以查询`host`、`default-target`和`target-list`等项目。支持结构化输出的项目接收`--format=json`。

## 编译器和工具链之间的边界

`wavec` 负责源码检查、代码生成和链接。 Vex 负责使用包 manifest、依赖图和 lockfile 构建可重现的包。 Whale是一个独立运行的低级工具链。

## std 指定路径

```shell
wavec --std-root /absolute/path/to/std check main.wave
wavec --std-root /absolute/path/to/std run main.wave
```

指定的std路径优先于安装路径，如果路径不正确或与std不兼容，将会失败。它不会自动替换其他安装中的std。选择与您的编译器相对应的std。

对于输出`-o`，使用与源/输入文件不同的路径。由于`check`不检查运行时操作，因此[练习](/docs/zh/practice/input-calculator)也会检查执行结果。
