---
translation_set_id: whale-amd64-target
path: whale/amd64-target
locale: en
group: whale
group_order: 5
order: 6
title: AMD64 target and ABI
summary: Linux AMD64 target properties, calling boundaries, and native feature coverage.
---

## Target identity

The native profile uses the target identifier `x86_64-whale-linux`.

| Property | Value |
| --- | --- |
| Operating system | Linux |
| Instruction set | AMD64 |
| Byte order | Little endian |
| Native address width | 64 bits |
| Object format | ELF64 |
| C calling convention | SysV AMD64 ABI |
| Static executable format | ELF ET_EXEC |

A build host and an output target are different concepts. Running Whale on another host does not imply support for that host's instruction set or object format. See [feature availability](overview) for the implemented compilation paths.

## Calls and signatures

A supported call must have an explicit signature and calling convention. Unsupported signatures are errors; a backend must not approximate them by dropping arguments or substituting a different representation.

Signature support distinguishes basic integers and pointers, f32/f64, aggregates and wide values, and variadic arguments. A type appearing in the IR type system does not by itself imply ABI support for passing or returning that type. Check each required category before selecting a backend.

For example, a backend accepting an i32 return type must still reject an unsupported aggregate-return signature. It cannot use the scalar return convention merely because part of the aggregate fits in a scalar register.

## Internal calls and C boundaries

Native pointer addresses occupy 64 bits. Tracked pointers also carry shadow metadata that must survive copies, storage, arguments, and return values.

The C ABI and the internal metadata-passing convention are distinct. Calls across a C boundary require explicit adapters. Passing a numerical address through the C ABI must not silently claim that allocation identity, lifetime, bounds, or access permissions were preserved.

## Stack frames

Keep the frame pointer and do not use the red zone. At O0, local-variable storage is not reused between distinct variables. Call-frame debug information describes the native frame and its relation to the original program.

These rules support debugging and do not establish exception-unwinding support. See [O0 and debugging](o0-debugging).

## Profile limits

The first O0 native profile does not include:

- O1 and higher optimization, vectorization, or LTO.
- Shared memory and atomic operations.
- Exception unwinding, async functions, or coroutines.
- Garbage collection or language-specific ownership enforcement.
- Arbitrary ownership transfer of external memory.
- Full inline-assembly constraint support.
- Dynamic linking, TLS, additional instruction sets, or additional object formats.

A request outside the supported profile must be rejected rather than silently replaced with another feature. These limits describe the native compilation path; they do not remove independently available toolchain components.
