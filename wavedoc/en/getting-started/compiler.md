---
translation_set_id: compiler
path: getting-started/compiler
locale: en
group: getting-started
group_order: 1
order: 3
title: Compiler command reference
summary: wavec Describes commands, build pipeline, output, targets, diagnostics, dependency linking, and tool queries.
---

## command model

`wavec` is compiler CLI. It compiles individual inputs directly, provides compiler support information to tools, and manages installed standard library sources.

```text
wavec [global-options] <command> [command-options]
```

|command|Use|
| --- | --- |
| `wavec build <input...>` |Depending on the flags, it performs inspection, code generation, linking, or execution pipelines.|
| `wavec check <file>` |Nickname for `build <file> --emit=check`.|
| `wavec run <file> [-- <args...>]` |It is an alias for `build <file> --run`, and the argument after `--` is passed to the program.|
| `wavec print <item>` |Queries target and toolchain support information.|
| `wavec install std` |Install the standard library.|
| `wavec update std` |Update installed standard libraries.|
| `wavec --version` |Prints installed version information.|

A full list of commands and options can be found at `wavec --help`.

## Build, test and run

```shell
wavec build main.wave
wavec check main.wave
wavec run main.wave -- first-argument second-argument
```

`build` creates an executable file by default. `check` stops after completing the frontend check. `run` requires binary output and cannot be used with shared library builds.

Use `--dry-run` to verify a request and determine which steps to execute without compiling, linking, or executing.

```shell
wavec build main.wave --target riscv64-unknown-linux-gnu --dry-run
wavec build main.wave --dry-run --error-format=json
```

The JSON format is a stable, unified interface used by build tools such as Vex.

## emit and input type

```shell
wavec build main.wave --emit=ast
wavec build main.wave --emit=ir
wavec build main.wave --emit=bc
wavec build main.wave --emit=asm
wavec build main.wave --emit=obj -o main.o
wavec build main.wave --emit=bin -o app
```

Output emit types are `ast`, `ir`, `bc`, `asm`, `obj`, `bin`. `check` is a control mode and must be used alone. You can specify multiple output types that a pipeline accepts, separated by commas.

The input types are `wave`, `ir`, `bc`, `asm`, `obj`, `archive`. `--input-type=<kind>` forces the type of all inputs to be specified. When linking only the object or archive inputs, use the binaries emit and `--link-only`.

```shell
wavec build module.o --input-type=obj --link-only --emit=bin -o app
```

## output location

|options|effect|
| --- | --- |
| `-o <file>` |Specifies the primary output path.|
| `--out-dir <dir>` |emit Places the output in the specified directory.|
| `--target-dir <dir>` |Specifies intermediate and primary deliverable routes.|

## Optimization and diagnostic output

```shell
wavec -O2 build main.wave
wavec --debug-wave=tokens,ast build main.wave
```

The optimization steps are `-O0`, `-O1`, `-O2`, `-O3`, `-Os`, `-Oz`, `-Ofast`. `--debug-wave` can be followed by `tokens`, `ast`, `ir`, `mc`, `hex`, `all`, and multiple steps can be combined with commas.

## native link

```shell
wavec --link=m -L ./lib build main.wave
wavec build main.wave --shared -o libexample.so
wavec build main.wave --static -o app
wavec build main.wave --pie -o app
```

`--link=<lib>` adds native libraries and `-L <path>` adds search paths. Link modes use `--shared`, `--static`, `--pie`, `--no-pie` according to compatibility rules.

Backend and linker control options include:

- `--target`, `--cpu`, `--features`, `--abi`, `--sysroot`
- `-C linker=<path>` and `-C link-arg=<arg>`
- `-C link-sysroot=<path>` and `-C relocation-model=<model>`
- `-C no-default-libs`

Freestanding outputs such as kernels use `--freestanding` along with `--entry`, `--linker-script`, and `--no-start-files` settings appropriate for the environment.

## External package interpretation

```shell
wavec --dep-root .vex/deps build main.wave
wavec --dep math=/opt/wave-deps/math build main.wave
```

`--dep-root <dir>` adds a root to find external `package::module` import. `--dep <name>=<path>` fixes the package name to one directory. This is the compiler integration point and projects manifest, dependency downloads and lockfile are handled by Vex.

## Support function query

Tools that use the target or output type can query support information with `wavec print`.

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

You can also query items like `host`, `default-target`, and `target-list`. Items that support structured output receive `--format=json`.

## Boundary between compiler and toolchain

`wavec` is responsible for source inspection, code generation, and linking. Vex is responsible for building reproducible packages with package manifest, the dependency graph, and lockfile. Whale is a low-level toolchain that runs independently.

## std Specify path

```shell
wavec --std-root /absolute/path/to/std check main.wave
wavec --std-root /absolute/path/to/std run main.wave
```

The specified std path takes precedence over the installation path, and will fail if the path is incorrect or incompatible with std. It does not automatically replace std from other installations. Select std that corresponds to your compiler.

For output `-o`, a different path is used than the source/input file. Since `check` does not check the runtime operation, [practice](/docs/en/practice/input-calculator) also checks the execution results.
