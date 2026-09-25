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

## Assemble a function

Save this as `answer.asm`:

```asm
section .text
global answer

answer:
    mov eax, 42
    ret
```

```sh
whale asm --amd64 answer.asm -o answer.o
```

The `.text` payload is `b8 2a 00 00 00 c3`: `mov eax, 42` followed by `ret`. The ELF64 object exports `answer`. It is a callable function without process startup code, not an executable. This example uses instructions accepted by the current assembler.

## Construct an object with the Rust API

This complete example uses the `object` crate to write the same function bytes with an explicit output target and global symbol:

```rust
use object::{ObjectFile, ObjectSymbol, SectionKind, SymbolBinding, SymbolVisibility, Target};

fn main() {
    let mut object = ObjectFile::with_target(Target::X86_64WhaleLinux.object_target());
    let text = object.add_section(".text", SectionKind::Text, 16);
    object.sections[text].data = vec![0xb8, 0x2a, 0x00, 0x00, 0x00, 0xc3];
    object.symbols.push(ObjectSymbol {
        name: "answer".into(),
        section_index: Some(text),
        value: 0,
        size: 6,
        binding: SymbolBinding::Global,
        visibility: SymbolVisibility::Default,
    });
    let elf = object.write().unwrap();
    assert_eq!(&elf[..7], b"\x7fELF\x02\x01\x01");
    assert_eq!(u16::from_le_bytes([elf[18], elf[19]]), 62);
    std::fs::write("answer.o", elf).unwrap();
}
```

The assertions inspect ELF class, byte order, and machine identity. `value: 0` is an offset within `.text`, and `size: 6` is the symbol's byte extent. An invalid section reference or extent fails serialization. A different machine or byte order also fails rather than being labeled as AMD64.

## Resolve symbols from two objects

The `linker` crate currently exposes symbol resolution. The following executable example gives two objects a local `helper`, then checks that exporting the same name twice produces an error:

```rust
use linker::core::symbol_table::{SymbolKey, SymbolTable};
use object::{
    ObjectFile, ObjectFormat, ObjectSymbol, SectionKind, SymbolBinding, SymbolVisibility,
};

fn input(binding: SymbolBinding) -> ObjectFile {
    let mut object = ObjectFile::new(ObjectFormat::ELF64);
    let text = object.add_section(".text", SectionKind::Text, 1);
    object.sections[text].data = vec![0xc3];
    object.symbols.push(ObjectSymbol {
        name: "helper".into(),
        section_index: Some(text),
        value: 0,
        size: 1,
        binding,
        visibility: SymbolVisibility::Default,
    });
    object
}

fn main() {
    let mut symbols = SymbolTable::new();
    symbols
        .resolve(&[input(SymbolBinding::Local), input(SymbolBinding::Local)])
        .unwrap();
    for object_index in 0..2 {
        let key = SymbolKey::Local {
            object_index,
            name: "helper".into(),
        };
        assert_eq!(symbols.symbols[&key].object_index, Some(object_index));
    }
    let error = symbols
        .resolve(&[input(SymbolBinding::Global), input(SymbolBinding::Global)])
        .unwrap_err();
    assert_eq!(error, "Duplicate global symbol: helper");
    println!("{error}");
}
```

Output:

```text
Duplicate global symbol: helper
```

The two local definitions have separate keys through `object_index`. The two global definitions collide. This example resolves model objects directly; it does not read `.o` files, apply relocations, or emit an executable. The complete definition-selection policy described below, including weak precedence, is not yet implemented.

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

### Reserving BSS without allocating its payload

Save this as `buffer.asm`. It reserves one TiB of logical BSS; it does not allocate or write one TiB while assembling:

```asm
section .bss
global buffer
buffer:
    resb 1099511627776
buffer_end:

section .text
    ret
```

```sh
whale asm --amd64 buffer.asm -o buffer.o
```

`buffer` has section-relative value 0 and `buffer_end` has value 1099511627776. The `.bss` header is `SHT_NOBITS` with that size and no file payload. Returning to `.bss` later continues its logical offset. Zero data directives also increase the logical size; nonzero initializers, relocations into BSS, and instructions in BSS are rejected.

The same distinction is available through the object and linker APIs:

```rust
use linker::core::layout::Layout;
use object::{ObjectFile, ObjectFormat, SectionKind};

fn main() {
    let mut object = ObjectFile::new(ObjectFormat::ELF64);
    let text = object.add_section(".text", SectionKind::Text, 16);
    object.sections[text].data = vec![0xc3];
    let bss = object.add_section(".bss", SectionKind::Bss, 16);
    object.sections[bss].zero_fill = 1 << 40;
    assert!(object.sections[bss].data.is_empty());

    let elf = object.write_with_limit(4096).unwrap();
    assert!(elf.len() < 4096);
    let layout = Layout::compute(&[object], 0x1000).unwrap();
    assert_eq!(layout.file_size, 1);
    assert_eq!(layout.memory_size, (1 << 40) + 16);
    assert_eq!(layout.sections[bss].memory_address, 0x1010);
    assert_eq!(layout.sections[bss].file_size, 0);
}
```

```text
section  memory address  file bytes       memory bytes
.text    0x1000          1                1
.bss     0x1010          0                1099511627776
```

`Section::zero_fill` counts additional unbacked BSS bytes. Its checked memory size is `data.len() + zero_fill`. Existing all-zero BSS `data` remains accepted, but callers can avoid that allocation with an empty `data` vector. Non-BSS sections require `zero_fill == 0`. Physical zero storage does not establish IR initialization state.

`Layout::compute` returns `Result` and records an input object/section mapping, alignment, file offset, memory address and both sizes for every section. It checks address and alignment arithmetic, keeps input order, and treats object alignment 0 as no constraint (alignment 1). BSS never advances the file cursor. These are payload placements; executable headers, permission-bearing load segments and relocation application remain separate work.

ELF writers reject overflow and field-width truncation. Extended section numbering is unsupported: the total header count, including generated tables and relocation sections, must be below `0xff00`. The default serialized-output limit is 256 MiB. `ObjectFile::write_with_limit` or `write_elf_with_limit` accepts a byte budget, including padding and tables; sparse BSS memory size does not count against it. Output size is validated before allocating the final byte vector.

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


## Optional Wave record serializer

A Linux x86_64 build of Whale can use a Wave implementation for fixed ELF64 header, section, symbol and RELA records. The default Rust implementation remains available. Output target selection, object validation, layout, symbol resolution and buffer allocation stay in Rust; selecting Wave does not add an architecture or a complete linker.

Build the optional path from the Whale repository with Rust, LLVM 21 development libraries, a C linker and `ar` installed:

```sh
git clone https://github.com/wavefnd/Wave.git /tmp/whale-wave-bootstrap
git -C /tmp/whale-wave-bootstrap checkout --detach 8a465e30aeea4b817d925cdd0e8d08c1bb029c9a
python3 tools/build_wave_elf.py --wave-source /tmp/whale-wave-bootstrap --out-dir /tmp/whale-wave-elf
WHALE_WAVE_ELF_DIR=/tmp/whale-wave-elf cargo build --locked --all-features
WHALE_WAVE_ELF_DIR=/tmp/whale-wave-elf cargo test --locked --workspace --all-features
```

The script verifies the pinned revision and rejects tracked modifications, rebuilds the Wave compiler, emits the Wave object through LLVM, and archives it for static linking. `WHALE_WAVE_ELF_DIR` must contain that archive; a requested but invalid archive or unsupported host fails the build. Without the variable, ordinary builds including `--all-features` require no Wave compiler. The linked Whale executable needs no Wave compiler at runtime. Cross-compiling Whale through this bootstrap is not supported yet.

For example, assemble this as `return.asm` with the resulting binary:

```asm
section .text
global entry
entry:
    ret
```

```sh
target/debug/whale asm --amd64 return.asm -o return.o
```

```text
ELF class: ELF64
machine: AMD64 (62)
.text bytes: c3
```

The record ABI passes a kind, pointer to u64 fields, field count, output pointer and capacity. Buffers belong to the Rust caller. No allocation ownership or Rust enum/String/Vec representation crosses the boundary. The Wave routine validates count, capacity and field widths before writing; it returns status 0 on success, 1 for invalid shape/pointers/capacity and 2 for field overflow. Pointers must designate live, correctly sized, non-overlapping buffers; raw C pointers cannot prove those conditions. The wrapper supplies these invariants. Tests compare complete ELF files against the Rust path, including BSS and signed relocation addends. This is a partial Wave implementation with a Wave/LLVM bootstrap, not complete self-hosting.
