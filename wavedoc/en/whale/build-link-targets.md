---
translation_set_id: build-link-targets
path: whale/build-link-targets
locale: en
group: whale
group_order: 1
order: 5
title: Building, linking, and target options
summary: Describes emit deliverables, input types, links, target/CPU/ABI and the freestanding build plan.
---

## emit Output

```shell
wavec build main.wave --emit=check
wavec build main.wave --emit=ast
wavec build main.wave --emit=ir
wavec build main.wave --emit=bc
wavec build main.wave --emit=asm
wavec build main.wave --emit=obj -o main.o
wavec build main.wave --emit=bin -o app
```

artifact emit The types are `ast`, `ir`, `bc`, `asm`, `obj`, `bin`. `check` is an inspection control mode, not a deliverable type, and is not to be used with the other artifact emit.

```shell
wavec print supported-emit-kinds
```

## Input type and link-only

In addition to the Wave source, the compiler distinguishes between IR, bitcode, assembly, object and archive inputs. The support list is queried with the following command:

```shell
wavec print supported-input-types
```

To link only the already created object or archive, you can use `--input-type` and `--link-only`.

```shell
wavec build module.o --input-type=obj --link-only --emit=bin -o app
```

## native link

```shell
wavec --link=m -L ./lib build main.wave
```

`--link` adds a library and `-L` adds a search path. Even if a symbol is declared in FFI, the library providing that symbol is not automatically linked.

## Select target

Option to select OS and CPU to run.

- `--target <triple>`
- `--cpu <name>`
- `--features <csv>`
- `--abi <name>`
- `--sysroot <path>`

Check host defaults and supported targets with the following command:

```shell
wavec print host-target
wavec print supported-targets
wavec print target-spec --target <triple>
wavec print cpu-list --target <triple>
wavec print target-features --target <triple>
```

## Supported by

The program for your current computer will be built without specifying target. To select a different environment, pass the target name below to `--target`.

|OS·Environment|architecture|target name|
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
| WebAssembly |64 bit| `wasm64-unknown-unknown` |

For programs that run without OS, use `x86_64-unknown-none-elf`, `aarch64-unknown-none-elf`, and `riscv64-unknown-none-elf`. A complete list of targets for installed versions can be found at `wavec print supported-targets`.

## RISC-V 64 contract

The default values for Hosted RISC-V targets are `generic-rv64`, RV64GC, `lp64d` ABI. The default values ​​for the Freestanding target are `generic-rv64`, RV64IMAC, and `lp64`.

```shell
wavec print target-spec --target riscv64-unknown-linux-gnu --format=json
wavec print target-spec --target riscv64-unknown-none-elf --format=json
```

Supports RISC-V CPU for `generic`, `generic-rv64`, `rocket-rv64`, `sifive-u74`. Feature override is `m`, `a`, `f`, `d`, `c`, `zicsr`, `zifencei` Place a sign in front of the name and separate it with a comma.

```shell
wavec build main.wave \
  --target riscv64-unknown-linux-gnu \
  --features=+m,+a,+f,-d,+c,+zicsr \
  --abi=lp64f
```

RISC-V Verification rejects inconsistent combinations. `d` requires `f`, and `f` requires `zicsr`. `lp64`, `lp64f`, `lp64d` must match the floating point feature you have enabled. If you do not specify ABI directly, the compiler will derive ABI from feature.

## Freestanding Link

```shell
wavec build kernel.wave \
  --target riscv64-unknown-none-elf \
  --freestanding \
  --entry=_start \
  --linker-script=linker.ld \
  --no-start-files \
  -o kernel.elf
```

`--freestanding` adjusts build settings to avoid using default libraries. `--entry` specifies linker entries, `--linker-script` specifies scripts, and `--no-start-files` specifies host startup file exclusions.

You can use `--dry-run` to check the link plan before actual execution.

## Hosted Cross Link

When creating a program to run on another OS·CPU, specify the target environment's library path as sysroot. If you use a separate linker, specify the path as `-C linker`.

```shell
wavec build main.wave \
  --target riscv64-unknown-linux-gnu \
  --sysroot /path/to/riscv64-sysroot \
  -C linker=/path/to/target-linker \
  -o app-riscv64
```

The libraries to be linked with sysroot are aligned with the selected OS·CPU·ABI.

## What to check in cross build

- target triple is in your compiler's support list
- sysroot and the linker match the target ABI
- Is the link library for the target architecture?
- CPU feature is valid for target CPU
- If freestanding, check whether the entry symbol and memory placement match the linker script.

## Run WebAssembly

wasm64 results are used in an execution environment that supports memory64. Modules that use external functions such as file, time, and input must connect host import corresponding to that function.

The process of creating and connecting the target code can be found in `--dry-run`. Actual execution takes place in the selected OS or WebAssembly execution environment.
