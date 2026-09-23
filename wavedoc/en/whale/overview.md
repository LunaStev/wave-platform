---
translation_set_id: whale-overview
path: whale/overview
locale: en
group: whale
group_order: 5
order: 1
title: Whale documentation
summary: Toolchain components, available commands, and reference manuals.
---

## Introduction

Whale is a general-purpose compiler toolchain for language implementers and compiler tools. It provides a typed intermediate representation (IR), an AMD64 assembler, an object-file library, and linker infrastructure. The components are exposed as Rust libraries and through the `whale` command.

The IR specifies the meaning of computations independently of their machine encoding. The assembler encodes machine instructions and emits relocatable objects. The object library represents sections, symbols, and relocations. The linker resolves references between objects and lays out an executable.

## Using this manual

| Reference | Contents |
| --- | --- |
| [IR reference](ir-reference) | Types, values, functions, control flow, validation, and interchange |
| [Numeric operations](numeric-operations) | Integer arithmetic, shifts, conversions, and floating point |
| [Memory model](memory-model) | Initialization, pointer validity, address calculation, layout, and strings |
| [O0 and debugging](o0-debugging) | Preserved computations, variable storage, and debug information |
| [AMD64 target](amd64-target) | Target identity, calling conventions, and native feature limits |
| [Assembler and linker](assembler-linker) | Assembly operands, sections, symbols, and static linking |

These references specify Whale semantics. Feature availability is listed below; describing an operation in the reference does not make it available in every build.

## Feature availability

| Component | Available interface | Limitations |
| --- | --- | --- |
| Assembler | AMD64 assembly to ELF64 relocatable objects | Instruction and directive coverage is incomplete |
| Object library | Object construction and ELF64 serialization | A relocatable object is not an executable |
| IR | Construction, printing, verification, and scalar AST lowering | Experimental API; strict versioned interchange and text parsing are incomplete |
| Linker | Symbol resolution infrastructure | Complete relocation application and executable output are unavailable |
| Execution and debugging | Semantic requirements described in this manual | IR interpretation, end-to-end native code generation, tracked-memory runtime checks, and DWARF emission are not yet available |

The no-undefined-behavior rules apply to verified IR and tracked memory. The experimental implementation does not yet provide the complete runtime enforcement described by the memory and execution references.

## Build and assemble

Build with Rust 1.86.0 or newer:

```sh
git clone https://github.com/wavefnd/Whale.git
cd Whale
cargo build --release --locked
```

Save this assembly as `answer.asm`:

```asm
section .text
global answer

answer:
    mov eax, 42
    ret
```

Create an ELF64 object:

```sh
./target/release/whale asm --amd64 answer.asm -o answer.o
```

Assembly uses Whale's own assembler; no external assembler is required. The output contains a callable function, not process startup code.

To enable the experimental AST-to-IR command, build with `--features socket-cli`. See the [command reference](/docs/en/toolchain/whale-cli) for commands available in the CLI.

## Library use and diagnostics

Validate IR before execution or code generation. Input errors and builder misuse are reported as structured errors; a library input error must not terminate the host process or overwrite existing content. Resource limits must be configurable for tools that accept untrusted input.

Artifact generation is required to be deterministic for the same toolchain version, input, target, and settings. Distribution metadata identifies the version, commit, and available features. CI coverage includes accepted and rejected inputs, traps, O0 preservation, round trips, and native semantics as those interfaces become available. See [Contributing to Whale](https://github.com/wavefnd/Whale/blob/master/CONTRIBUTING.md) for development and validation commands.
