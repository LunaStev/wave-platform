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

## Constructing and reading IR

This complete Rust example uses the `ir` crate to construct a function, verify it, and print its typed IR:

```rust
use ir::{ModuleBuilder, Target, Type};

fn main() {
    let target = Target::lookup("x86_64-whale-linux").unwrap();
    let mut module = ModuleBuilder::new(target.name(), target.data_layout());
    let mut function = module.begin_function("answer", vec![], Type::I32);
    let left = function.const_i32(40);
    let right = function.const_i32(2);
    let answer = function.add(Type::I32, left, right);
    function.ret(Some(answer));
    function.finish();
    let module = module.finish();
    ir::verify_module(&module).unwrap();
    print!("{}", ir::print_module(&module));
}
```

The printer produces:

```text
module {
  format_version 1
  semantics_version 1
  target "x86_64-whale-linux"
  datalayout { ptr=64, endian=little }

  fn @answer() -> i32 {
  entry:
    %v0: i32 = const i32 40
    %v1: i32 = const i32 2
    %v2: i32 = add i32 %v0, %v1
    ret i32 %v2
  }

}
```

`%v0` and `%v1` define i32 constants. `add` defines `%v2`, which supplies the function's i32 return. Changing the return value to a Bool would violate the signature and fail verification. The addition remains an instruction at O0 even though both operands are constant.

These are actual printer outputs, not input files for a text parser. Text parsing and IR execution are not yet available; the Rust builder is the available way to construct this module.

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

The scalar AST JSON contract below is available. Printed typed IR carries version metadata, but its parser and full round-trip interchange remain unavailable.


### Versioned AST JSON

Save the following as `program.json`. All four envelope fields are required. `program` contains required `globals` and `functions` arrays, which may be empty. Function name, parameters, return type and body are required. Each enum uses either its unit name or a single variant-key object. `VarDecl.init` may be absent or null; other required fields must be present.

```json
{
  "format_version": 1,
  "semantics_version": 1,
  "features": [],
  "program": {
    "globals": [],
    "functions": [
      {
        "name": "answer",
        "parameters": [],
        "return_type": {
          "Int": {
            "bits": 128,
            "signed": false
          }
        },
        "body": [
          {
            "Return": {
              "Lit": {
                "Int": {
                  "bits": 128,
                  "signed": false,
                  "value": "340282366920938463463374607431768211455"
                }
              }
            }
          }
        ]
      }
    ]
  }
}
```

```sh
cargo run --locked --features socket-cli -- ir lower program.json
```

```text
module {
  format_version 1
  semantics_version 1
  target "x86_64-whale-linux"
  datalayout { ptr=64, endian=little }

  fn @answer() -> u128 {
  entry:
    %v0: u128 = const u128 340282366920938463463374607431768211455
    ret u128 %v0
  }

}
```

Integer `value` is a decimal string: optional minus for signed integers, followed by decimal digits, without whitespace, plus, exponent or separators. The declared width/signedness determines the accepted range. `u128::MAX` above survives JSON and lowering exactly. A negative unsigned value or an out-of-range value fails instead of wrapping. Floating values use exact-width hexadecimal storage strings, described in [numeric operations](numeric-operations).

`format_version` is 1 for this AST format; `semantics_version` is 1. `features` must be an empty array. Unknown fields, versions, features, duplicate raw JSON keys (including escaped equivalent keys), and trailing values are errors, even with `--no-verify`. The library entry point is `ir::lower_ast::interchange::decode`; `encode` emits the envelope. `decode` defaults to an 8 MiB source-byte limit; `decode_with_limit` accepts a caller limit. JSON nesting is bounded. Use this raw decoder rather than parsing into a generic map that could already discard duplicate keys.

[The complete JSON Schema](https://github.com/wavefnd/Whale/blob/master/ir/schema/ast-v1.schema.json) specifies shapes, required fields and variants. Range/type checks and duplicate-key detection additionally apply. The scalar lowering subset includes literals, variables/constants, add/sub/mul, comparisons, assignment, if/while, return and break/continue. Call and aggregate expressions are unsupported. `Opaque` is representable in the schema but unsupported by lowering.

Migration requires wrapping old bare Program payloads and replacing numeric JSON literals with decimal integer strings or float bit strings. Old unversioned payloads are rejected. AST and typed IR version numbers are independent even when both currently equal 1.

### Rejected input and CLI recovery

Save this complete input as `invalid.json`:

```json
{"format_version":99,"semantics_version":1,"features":[],"program":{"globals":[],"functions":[]}}
```

```sh
cargo run --locked --features socket-cli -- ir lower invalid.json -o rejected.wir
```

```text
Failed to parse socket JSON: unsupported AST format_version 99; expected 1
```

The command exits nonzero without creating an output or replacing an existing file. A type mismatch also fails before output publication. A binary built without `socket-cli` exits with status 2 and prints a recovery command containing `--features socket-cli`.
