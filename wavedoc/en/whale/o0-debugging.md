---
translation_set_id: whale-o0-debugging
path: whale/o0-debugging
locale: en
group: whale
group_order: 5
order: 5
title: O0 and debugging
summary: Preservation of source computations, constant expressions, storage, and debug correspondence.
---

## Preservation model

O0 preserves the original typed IR's computations, variables, and control flow for debugging. Unused results and structurally unreachable blocks remain present. Verification diagnoses invalid IR without deleting blocks or simplifying operations.

Required machine transformations take place in a separate lower IR. The mapping to original IDs must survive those transformations. Preserving typed IR does not require every IR operation to correspond to exactly one machine instruction.

## Transformations excluded at O0

| Transformation | O0 behavior |
| --- | --- |
| Inlining | Keep the call and function boundary |
| Tail-call transformation | Keep the ordinary call/return structure |
| Dead-code elimination | Retain unused computations and unreachable blocks |
| Runtime constant folding | Keep the original operation |
| Common-subexpression elimination | Keep distinct computations |
| Local-variable storage reuse | Keep separate local storage |
| Frame-pointer omission | Retain the frame pointer |
| Automatic string merging | Keep distinct string objects |

For example, a runtime addition of two constant operands remains an addition even if the result is unused. A branch with a constant condition retains its original control-flow structure.

## Compile-time constants

A compile-time constant declaration retains both its typed initializer expression and its evaluated result. This is different from folding an ordinary runtime instruction.

For a declaration whose initializer is `1 + 2`, the stored information includes the addition expression and the result `3`. This example describes the expression, not source-language declaration syntax. The result can be used to construct static data without turning the initializer into runtime arithmetic.

Unused and unreachable declarations remain identifiable. Names, declaration IDs, and references must continue to distinguish shadowed declarations. Verification rejects invalid references, dependency cycles, invalid types, and a cached result that disagrees with its expression.

## Unreachable source statements

Statements following return, break, or continue remain represented in disconnected blocks. Their presence must not change the preceding terminator or create a new executable path. Invalid expressions in unreachable statements still produce diagnostics.

Retaining such statements is useful for inspecting the original program structure. It does not cause them to execute after the terminator.

## Preserved IR example

The following module passes verification and shows the current printer representation of an unused computation, a compile-time declaration, and a disconnected block:

```text
module {
  format_version 1
  semantics_version 1
  target "x86_64-whale-linux"
  datalayout { ptr=64, endian=little }

  fn @preserved() -> void {
  entry:
    %v0: i32 = const i32 1
    %v1: i32 = const i32 2
    %v2: i32 = add i32 %v0, %v1
    %v3: i32 = const_decl "count" add(i32 1, i32 2) => const i32 3
    ret void
  unreachable.cont:
    %v4: i32 = const i32 4
    %v5: i32 = const i32 5
    %v6: i32 = add i32 %v4, %v5
    ret void
  }

}
```

`%v2` remains an `add` even though it has no uses. `%v3` is a separate `const_decl`: it retains `add(i32 1, i32 2)` and the evaluated result 3. The `unreachable.cont` block has no incoming edge, so its `%v6` addition remains inspectable without adding an execution path after `ret void`.

Verification preserves these instructions and blocks. This example demonstrates IR construction, verification, and printing; it does not imply available native execution or DWARF emission.

## Source and stack information

The debug interface uses DWARF 5 for function records, source lines, basic local variables, and call-frame information. Function and local-variable information must remain associated with original IR identities through lowering.

The AMD64 profile retains frame pointers and does not use the red zone. Call-frame information supports stack inspection; it does not imply support for exception unwinding. Traps terminate execution without guaranteeing destructors or unwinding.

DWARF emission and native execution availability are listed in the [toolchain overview](overview). O1 and higher optimization behavior is outside this O0 reference.
