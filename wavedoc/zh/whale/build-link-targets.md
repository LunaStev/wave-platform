---
translation_set_id: build-link-targets
path: whale/build-link-targets
locale: zh
group: whale
group_order: 1
order: 5
title: 构建、链接与目标选项
summary: 描述 emit 可交付成果、输入类型、链接、target/CPU/ABI 和独立构建计划。
---

## emit 输出

```shell
wavec build main.wave --emit=check
wavec build main.wave --emit=ast
wavec build main.wave --emit=ir
wavec build main.wave --emit=bc
wavec build main.wave --emit=asm
wavec build main.wave --emit=obj -o main.o
wavec build main.wave --emit=bin -o app
```

artifact emit 类型为`ast`、`ir`、`bc`、`asm`、`obj`、`bin`。 `check` 是检验控制模式，不是交付类型，不能与其他artifact emit 一起使用。

```shell
wavec print supported-emit-kinds
```

## 输入类型和 link-only

除了 Wave 源之外，编译器还区分 IR、bitcode、assembly、object 和 archive 输入。使用以下命令查询支持列表：

```shell
wavec print supported-input-types
```

要仅链接已创建的object或archive，您可以使用`--input-type`和`--link-only`。

```shell
wavec build module.o --input-type=obj --link-only --emit=bin -o app
```

## 原生链接

```shell
wavec --link=m -L ./lib build main.wave
```

`--link` 添加库，`-L` 添加搜索路径。即使在FFI中声明了一个符号，提供该符号的库也不会自动链接。

## 选择目标

选择运行 OS 和 CPU 的选项。

- `--target <triple>`
- `--cpu <name>`
- `--features <csv>`
- `--abi <name>`
- `--sysroot <path>`

使用以下命令检查主机默认值和支持的目标：

```shell
wavec print host-target
wavec print supported-targets
wavec print target-spec --target <triple>
wavec print cpu-list --target <triple>
wavec print target-features --target <triple>
```

## 支持者

您当前计算机的程序将在不指定target的情况下构建。要选择不同的环境，请将下面的目标名称传递给`--target`。

|OS·环境|建筑学|目标名称|
| --- | --- | --- |
| Linux | amd64 | `x86_64-unknown-linux-gnu` |
| Linux | ARM64 | `aarch64-unknown-linux-gnu` |
| Linux | RISC-V64 | `riscv64-unknown-linux-gnu` |
| Linux | LoongArch64 | `loongarch64-unknown-linux-gnu` |
| macOS | amd64 | `x86_64-apple-darwin` |
| macOS | ARM64 | `aarch64-apple-darwin` |
| Windows | amd64 | `x86_64-pc-windows-msvc` |
| Windows | ARM64 | `aarch64-pc-windows-msvc` |
| FreeBSD | amd64 | `x86_64-unknown-freebsd` |
| WebAssembly |64位| `wasm64-unknown-unknown` |

对于不使用 OS 运行的程序，请使用 `x86_64-unknown-none-elf`、`aarch64-unknown-none-elf` 和 `riscv64-unknown-none-elf`。已安装版本的完整目标列表可以在`wavec print supported-targets`找到。

## RISC-V 64合约

Hosted RISC-V 目标的默认值为 `generic-rv64`、RV64GC、`lp64d` ABI。 Freestanding目标的默认值为`generic-rv64`、RV64IMAC和`lp64`。

```shell
wavec print target-spec --target riscv64-unknown-linux-gnu --format=json
wavec print target-spec --target riscv64-unknown-none-elf --format=json
```

支持 RISC-V CPU `generic`、`generic-rv64`、`rocket-rv64`、`sifive-u74`。 Feature override 为 `m`、`a`、`f`、`d`、`c`、`zicsr`、 `zifencei` 在姓名前面放置一个符号，并用逗号分隔。

```shell
wavec build main.wave \
  --target riscv64-unknown-linux-gnu \
  --features=+m,+a,+f,-d,+c,+zicsr \
  --abi=lp64f
```

RISC-V 验证拒绝不一致的组合。 `d` 需要`f`，`f` 需要`zicsr`。 `lp64`、`lp64f`、`lp64d` 必须与您启用的浮点feature 匹配。如果不直接指定ABI，编译器将从feature派生出ABI。

## 独立式链接

```shell
wavec build kernel.wave \
  --target riscv64-unknown-none-elf \
  --freestanding \
  --entry=_start \
  --linker-script=linker.ld \
  --no-start-files \
  -o kernel.elf
```

`--freestanding` 调整构建设置以避免使用默认库。 `--entry` 指定链接器条目，`--linker-script` 指定脚本，`--no-start-files` 指定主机启动文件排除。

在实际执行之前，您可以使用`--dry-run`检查链接计划。

## Hosted 交叉链接

创建在另一个OS·CPU上运行的程序时，请将目标环境的库路径指定为sysroot。如果您使用单独的链接器，请将路径指定为 `-C linker`。

```shell
wavec build main.wave \
  --target riscv64-unknown-linux-gnu \
  --sysroot /path/to/riscv64-sysroot \
  -C linker=/path/to/target-linker \
  -o app-riscv64
```

要与 sysroot 链接的库与所选的 OS·CPU·ABI 对齐。

## 交叉构建中要检查什么

- target triple 在您的编译器的支持列表中
- sysroot 且链接器与目标ABI 匹配
- 是目标架构的链接库吗？
- CPU feature 对目标 CPU 有效
- 如果是独立的，请检查入口符号和内存布局是否与链接描述文件匹配。

## 运行WebAssembly

wasm64结果用于支持memory64的执行环境。使用文件、时间、输入等外部函数的模块必须连接与该函数对应的hostimport。

创建和连接目标代码的过程可以参见`--dry-run`。实际执行发生在所选的OS或WebAssembly执行环境中。
