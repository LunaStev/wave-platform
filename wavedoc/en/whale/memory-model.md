---
translation_set_id: whale-memory-model
path: whale/memory-model
locale: en
group: whale
group_order: 5
order: 4
title: Memory model
summary: Tracked allocations, initialized reads, pointer arithmetic, layout, and string storage.
---

## Tracked allocations

A tracked pointer associates an address with an allocation identity, generation, bounds, offset, and access permissions. The memory model also tracks allocation lifetime and initialization state. Access must satisfy the tracked conditions; a violation traps.

The generation distinguishes lifetimes even when a physical address is reused. An address alone does not establish that a pointer is valid or that the caller owns access to its storage.

Native addresses remain 64 bits. Shadow metadata accompanies pointers through copying, storage, calls, and returns. The initial native memory scope consists of trackable stack and global allocations. C boundaries require explicit adapters; arbitrary transfer of external-memory ownership is outside this scope.

## Initialization and reads

Declaring storage does not initialize its value. Check initialization for the byte range actually read. Reading an uninitialized part of a value traps. Do not replace the read with zero or an unspecified value.

Value-read checks exclude padding bytes. For example, after every field of a structure has been initialized, an uninitialized gap inserted for field alignment does not by itself make a value read invalid.

A memory copy propagates initialization state along with the bytes. Copying uninitialized storage does not turn it into initialized storage. A later value read at the destination is checked in the same way as a read at the source.

Physical zeros in BSS do not, on their own, establish that an IR variable has been initialized.

## Pointer arithmetic and comparison

Address-calculation overflow traps. A pointer one element past an allocation may be formed, but it must not be used to access memory.

Pointer equality uses allocation identity; the numerical address alone does not establish identity. Ordering and subtraction between pointers belonging to different allocations trap. Reconstructing an address from an integer does not restore access permissions.

### GEP

GEP computes an address in element and field units. It does not load a value.

The first index offsets in units of the base pointer's pointee type. Subsequent indices select array elements or structure/tuple fields. A structure or tuple field index is a compile-time field ordinal, not a byte offset. The selected type determines the result pointer type.

For a base of type `ptr<array<i32, 4>>`, indices `[0, 2]` select the third i32 element of that array and produce `ptr<i32>`. A first index of 1 instead advances by one entire four-element array. These examples describe indices, not textual instruction syntax.

Pointer arithmetic on zero-sized elements is rejected in the initial native target. Calculating an address does not remove the lifetime, bounds, initialization, or permission checks on a later access.

## Data layout

The output target determines sizes, alignments, field offsets, and array strides. Structure and tuple fields retain their declared order. Do not derive target layout from the machine on which the compiler happens to run.

| Value | Storage rule |
| --- | --- |
| Bool, signed i1, unsigned u1 | At least one byte |
| Empty structure or tuple | Size zero, alignment one |
| Array | Element stride determined by target layout |
| Structure or tuple | Ordered fields with target-required alignment |

Alignment in completed IR is a nonzero power of two. Resolve automatic alignment before producing that IR. Packed layouts, unions, and bitfields are unsupported in this profile and must be rejected.

## Strings and C boundaries

A string is an immutable byte sequence with an explicit length. UTF-8 is the default encoding. Embedded NUL bytes are allowed; there is no implicit terminating NUL. O0 does not automatically merge equal string objects.

The three-byte sequence `A`, NUL, `B` therefore has length three. It cannot be passed through the explicit C-string conversion, which rejects embedded NUL. A frontend must not silently truncate it to `A`.

External C, raw addresses, and inline assembly are separate contract boundaries. Runtime checking of tracked memory does not guarantee detection of every invalid action performed by external code.
