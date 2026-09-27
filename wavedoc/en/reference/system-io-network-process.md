---
translation_set_id: system-io
path: reference/system-io-network-process
locale: en
group: stdlib
group_order: 1
order: 15
title: System functions and processes
summary: Describes the boundary and process lifetime of the parent API and OS interfaces.
---

## Documentation by function

Read [fs and io](/docs/en/stdlib/files-io) to handle the file, [TCP](/docs/en/stdlib/tcp) to link, and [resolver](/docs/en/stdlib/resolution) to look up the address. Below are the process and lower level OS access rules.

## Process Basics API

```text
std::process::core
proc_exit(code: i32) -> !
proc_getpid() -> i64
proc_getppid() -> i64
proc_execve(path: str, argv: ptr<ptr<i8>>, envp: ptr<ptr<i8>>) -> i64
proc_waitpid_raw(pid: i64, status: ptr<i32>, options: i32) -> i64
```

`proc_exit` does not return to the call point. Please perform any necessary file/memory cleanup before shutdown. `proc_execve` differs from a typical child creation function because, if successful, it replaces the existing process image. raw argv/envp must be prepared for the NUL termination of each string, with the null pointer indicating the end.

The spawn function in `std::process::spawn` handles the creation result, and the await function handles the exit status of the child. Successful creation and successful termination of the program are two different things. When you create a pipe, the parent and child must close the unused end so that EOF is passed. If you wait for the child to exit without reading the capture pipe, the buffer may fill up and wait for each other.

## Portability and low-level approach

fork/exec, file descriptor, and Windows handle are not the same OS function. Verify support for the selected target and treat unsupported as a normal failure path. `std::sys` is a OS-specific interface and does not reuse its numeric flags and layout from other OS.

When linking directly with the external C library, please read [FFI](/docs/en/language/modules-imports-and-ffi). There is no need to arbitrarily declare the function libc to use the parent std API. Check [Target and Link Environment](/docs/en/whale/build-link-targets) first.
