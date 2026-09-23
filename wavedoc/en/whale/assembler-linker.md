---
translation_set_id: whale-assembler-linker
path: whale/assembler-linker
locale: en
group: whale
group_order: 5
order: 7
title: Assembler and static linking
summary: Operand encoding, section layout, symbol binding, and executable entry points.
---

## Assembly and objects

Whale's assembler converts AMD64 instructions into machine bytes and relocation records. ELF64 objects hold those records together with sections and symbols. Assembly is implemented by Whale and does not require an external assembler.

A relocatable object may contain references whose final addresses are not known. Resolving those addresses is a linking operation. Do not treat successful assembly as evidence that all external symbols can be resolved or that the output is executable.

## Literals and memory operands

Preserve a literal's width and signedness until the actual instruction encoding range is checked. A value that does not fit must produce an error instead of silent truncation.

An explicit memory width is required when the instruction does not otherwise determine it. For example, a register operand may supply the width, whereas a memory operand paired only with an immediate can be ambiguous. The assembler must not guess an ambiguous width.

An AMD64 memory operand containing only a symbol defaults to RIP-relative addressing. Explicit rel/abs addressing selects the intended mode. An unknown escape in a literal is an error.

## Sections and alignment

| Section content | Alignment behavior |
| --- | --- |
| Code | Insert NOP instructions |
| Initialized data | Insert zero bytes |
| BSS | Increase logical memory size without adding file payload |

File size and memory size are separate. BSS reserves memory but does not require an equally sized block of stored zeros in the object file. Custom sections carry their attributes, and symbols retain their binding and type information.

## Symbol identity

Functions and variables use separate identities inside the IR. External linking uses explicit `link_name` values supplied by the frontend. Whale preserves those names rather than automatically renaming one of two colliding exports.

Consequently, an internal function and variable may both be named `item`, but exporting both under the same external name can still be an error. Distinct internal namespaces do not create distinct external namespaces automatically.

Object-local symbols are scoped to their input object. Global symbols participate in resolution across objects. A confirmed function/data collision is an error. A NOTYPE symbol remains compatible with inputs that do not provide a more specific type; the absence of a type is not proof that a symbol denotes a function or data.

## Definition selection

| Definitions or references | Result |
| --- | --- |
| Strong and strong | Duplicate-definition error |
| Strong and weak | Select the strong definition |
| Weak and weak | Select the first definition in input order |
| Unresolved strong reference | Link error |
| Unresolved weak reference | Unsupported in the initial static profile; report an error |

Input order is therefore significant when multiple weak definitions exist. A deterministic link uses the supplied input order consistently.

## Static executable output

The static native profile produces ELF ET_EXEC with an explicitly supplied entry point. The entry point is not inferred from a function named `main`. Whale does not automatically insert startup code to call that function.

The profile does not automatically remove sections, fold identical code, or strip symbols. File layout must account separately for bytes stored in the file and memory reserved at runtime.

The complete static executable path is not yet available in the CLI. `whale asm` produces a relocatable object; `whale object` wraps raw bytes in an object. See the [toolchain overview](overview) for availability and the [AMD64 target](amd64-target) for ABI requirements.
