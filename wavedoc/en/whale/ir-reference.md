---
translation_set_id: whale-ir-reference
path: whale/ir-reference
locale: en
group: whale
group_order: 5
order: 2
title: Whale IR reference
summary: Types, identities, well-formed functions, evaluation order, and interchange rules.
---

## Modules and identities

A module contains target information, global definitions, and functions. Values have explicit types. Frontends resolve source-language names, types, overloads, and generics before producing typed IR.

Functions and global variables occupy separate internal namespaces. A function and a variable may therefore share a name. Internal identity is distinct from an external `link_name`; external names are supplied explicitly by the frontend. Whale does not resolve external collisions by inventing new names. See [symbols and linking](assembler-linker).

Each value definition has an identity. Definitions must be unique, and their type metadata must agree with their declared types. A name is not a substitute for identity when declarations shadow each other.

## Types

| Type | Meaning |
| --- | --- |
| `bool` | Logical false or true |
| `i1`, `i8`, `i16`, `i32`, `i64`, `i128` | Signed integers of the indicated bit width |
| `u1`, `u8`, `u16`, `u32`, `u64`, `u128` | Unsigned integers of the indicated bit width |
| `f16`, `f32`, `f64` | Floating-point values of the indicated bit width |
| `ptr<T>` | Pointer to a value of type T |
| `array<T, N>` | N elements of the same type |
| `struct{T, ...}` | Ordered structure fields |
| `tuple<T, ...>` | Ordered tuple elements |
| `void` | Absence of a result |

`bool`, `i1`, and `u1` are different types. Signed `i1` represents −1 and 0; unsigned `u1` represents 0 and 1. An integer 1 is not implicitly a Boolean condition. Branches, Select conditions, and `trap_if` require Bool operands.

Storage size is determined by the [target layout](memory-model), not solely by the number of value bits. For example, `i1` has one value bit but occupies at least one byte in memory.

## Functions and calls

A function specifies its complete parameter and result types, calling convention, and linkage. Direct and indirect calls must match the callable signature. A void call has no result ID. A nonvoid call defines a result even when that result is unused at O0.

A return must agree with the function result type. Void returns carry no value; nonvoid returns carry a value of the declared result type.

## Blocks and value availability

Every block has a unique identity and exactly one terminator. Branch destinations must belong to the same function. The entry block must exist and has no incoming edges or phi instructions. To form a loop, branch from entry to a separate loop header.

On executable paths, a value definition must dominate its ordinary uses: every path from entry to the use must pass through the definition. Within one block, the definition precedes the use. Physical block storage order does not establish dominance.

Consider entry branching to either left or right, followed by a join. A value defined only in left cannot be used as an ordinary value at the join, because the path through right has no definition. Use a phi with an input from each predecessor instead.

Unreachable blocks remain in the module. Verification still checks their identities, types, operands, and branch structure. A definition in an unreachable block cannot supply an ordinary value on a reachable path.

## Phi instructions

Phi instructions appear before all non-phi instructions in a block. A phi has exactly one input per distinct predecessor block. Each input value must have the phi's type and be available at the end of its corresponding predecessor.

Repeated edges from one predecessor require one input, not one input per edge. A loop phi can refer to a value computed on its backedge even when that block is stored later in the module. Missing, duplicate, unrelated, or incorrectly typed inputs are verification errors.

## Evaluation and selection

The Whale AST evaluates the call target and subexpressions in the specified left-to-right order. Frontends express short-circuit operations with control-flow branches.

Select chooses between values that have already been computed. It does not suppress evaluation of either input. For example, selecting a safe value does not prevent a trap in the computation of the other input. Put a potentially trapping computation inside a conditional block when it must not execute on the other path.

## Validation and traps

Malformed IR is rejected by verification. Runtime violations of defined execution conditions produce traps; `undef` and `poison` are not permitted values. Builder misuse, duplicate definitions, and attempts to add a second terminator must return structured errors.

A trap contains a reason, source location, and IR ID. It stops subsequent execution. Native execution terminates the program; the interpreter API returns a Trap error. Earlier side effects remain, but buffer flushing, destructors, and stack unwinding are not guaranteed.

The guarantee covers verified IR and tracked memory. External C, raw addresses, and inline assembly have separate contracts; violations beyond those boundaries are not guaranteed to be detected. See the [memory model](memory-model).

## Interchange and textual representation

AST and typed IR use separate format versions and a common semantics version. Readers reject missing or unknown versions, unknown fields or features, and duplicate JSON keys. Producers must not rely on a reader silently ignoring an unsupported property.

Integers carry a bit width, signedness, and a textual numeric value. Floating-point constants carry a width and an exact bit pattern. A text IR round trip must preserve names, IDs, types, constants, ordering, attributes, and metadata. Whitespace and comment placement need not survive the round trip.

These are interchange requirements, not a JSON schema or a complete textual grammar. Consult [feature availability](overview) before selecting a parser or serializer.
