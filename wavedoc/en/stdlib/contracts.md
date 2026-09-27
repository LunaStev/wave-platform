---
translation_set_id: stdlib-contracts
path: stdlib/contracts
locale: en
group: stdlib
group_order: 1
order: 2
title: Reading API documentation: Errors and ownership
summary: Understand argument units, result structures, partial success, and resource lifetimes.
---

## Read the declaration

The following notation describes the function declaration and is not the entire executable file.

```text
io_read(fd: i64, buf: ptr<u8>, len: i64) -> i64
```

`fd` is the open descriptor, `buf` is the storage provided by the caller, and `len` is the number of bytes available for writing. `i64` does not mean that negative lengths are valid. The return value is the actual number of bytes read, not the request length, so only the range returned is used.

## Failure expression varies from function to function

|way|yes|Inspection method|
| --- | --- | --- |
|Pointer or null| `mem_alloc` |null Memory access after inspection|
|number of bytes or negative number| `io_read` |Negative error, 0 EOF, positive data|
|status code| `buffer_push` |Error constant comparison with `BUFFER_OK`|
|Success and value| `NetResult<T>` |After examining `ok`, use `value`|
|Partial progress included| `RandomFillResult` |Check `ok`, `written`, `error` together.|

It only looks at the number of errors and does not compare them to constants in other modules. For example, error numbers env and OS errno are not the same system. The original error in WASI should not be interpreted as Linux errno.

## owning and renting

- **Owned**: When an allocated memory, open file, or open socket is acquired, it is responsible for calling the corresponding release/close.
- **Borrow**: The byte view or the buffer passed to the function refers to existing memory. If a function does not specify that it receives ownership, the caller retains control.
- **Output argument**: Passes a valid storage space where the result can be written to a function that receives it, such as `out_value: ptr<T>`. Ensure that the contract states that the result is only valid if successful.

Copying a Buffer structure can leave both copies pointing to the same allocation. Do not free each copy separately. A borrowed pointer becomes invalid after the allocation is freed or reallocated. String literals are not writable buffers.

## Failure does not mean reverting to a previous state

`io_write_all` may fail after writing some bytes. Bytes already written externally will not be returned. On the other hand, reading checked cursor of bytes preserves the position and output value if it fails. These differences are specified by API.

If you also want to practice handling failures, proceed with [File reader](/docs/en/practice/file-reader) and [Binary message](/docs/en/practice/binary-message).
