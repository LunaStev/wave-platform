---
translation_set_id: standard-library
path: reference/standard-library
locale: en
group: stdlib
group_order: 1
order: 1
title: Standard library guide
summary: How to find a module that suits your purpose and read the errors and ownership rules of the function.
---

## Find the features you need

The standard library is import with the path `std::module::file`. Even if the names are similar, functions may return errors in different ways. First read [API How to read](/docs/en/stdlib/contracts), then go to the module you need in the following table.

|What I want to do|document|Main import|
| --- | --- | --- |
|String length/comparison/search| [string](/docs/en/reference/string-and-bytes) | `std::string::len`, `cmp`, `find`, `trim` |
|Memory allocation/copy/size| [mem](/docs/en/reference/memory-and-buffer) | `std::mem::alloc`, `ops`, `layout` |
|List of bytes of varying size| [buffer](/docs/en/stdlib/buffer) | `std::buffer::alloc`, `read`, `write` |
|Binary read/write| [bytes](/docs/en/stdlib/bytes) | `std::bytes::types`, `cursor`, `leb128` |
|File·Descriptor I/O| [fs and io](/docs/en/stdlib/files-io) | `std::fs::file`, `std::io::fd` |
|Path combination/environment settings| [path and env](/docs/en/stdlib/path-env) | `std::path::copy`, `std::env::environ` |
|Time measurement/waiting| [time](/docs/en/stdlib/time) | `std::time::duration`, `clock`, `sleep` |
|Numerical address/name list search| [net.resolve](/docs/en/stdlib/resolution) | `std::net::resolve`, `resolve_table` |
|TCP Connection/Transmission| [net.tcp](/docs/en/stdlib/tcp) | `std::net::tcp`, `address`, `error` |
|OS Random number| [random](/docs/en/stdlib/random) | `std::random::fill` |
|Process·OS Boundary| [system function](/docs/en/reference/system-io-network-process) | `std::process::core`, `spawn`, `std::sys` |
|Asynchronous task execution| [task](/docs/en/stdlib/task) | `std::task` |
|Mathematics/Diagnosis Assistant| [math and debug](/docs/en/stdlib/math-debug) | `std::math::int`, `float`, `std::debug::core` |

## Example of first use

The program below uses one function from std without downloading a separate package. Save it as `main.wave` and run it as `wavec run main.wave`.

<!-- wave-example: library-import -->
```wave
import("std::string::len")::{
    len
};

fun main() {
    println("{}", len("Wave"));
}
```

The result is `4`. To understand the same example step by step, read [Strings](/docs/en/language/strings).

## Compatible with compiler std

Confirm the selected path with `wavec print std-path`. When using std from another checkout, specify the path as `wavec --std-root /absolute/path/to/std check main.wave`. If the specified path is invalid or incompatible, an error will be displayed.

## platform border

Distinguish between computational functions such as string/byte and OS functions such as file/socket. Recognizing a target does not guarantee that all hosts API will be provided. Read the platform entries for [Support target](/docs/en/whale/build-link-targets) and each API together. `std::sys` is a lower-level interface and portable programs will use the higher-level module first.
