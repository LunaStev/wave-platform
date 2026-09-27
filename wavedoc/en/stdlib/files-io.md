---
translation_set_id: stdlib-files-io
path: stdlib/files-io
locale: en
group: stdlib
group_order: 1
order: 7
title: fs and io: Files and byte transfers
summary: Describes file lifetime, full read, partial transfer and post-failure states.
---

## Select file API

The convenience function of `std::fs::file` receives the path and performs the necessary opening and closing. Functions that return a descriptor must be closed by the caller.

|declaration|Success results and precautions|
| --- | --- |
| `open_read(path: str) -> i64` |Open descriptor. Negative numbers are errors|
| `create(path: str) -> i64` |Create or delete existing file contents. Returns the owning descriptor|
| `open_append(path: str) -> i64` |Open or create for addition|
| `size(path: str) -> i64` |Number of bytes. Negative numbers are errors|
| `read_into(path: str, dst: ptr<u8>, dst_cap: i64) -> i64` |Number of bytes in the entire file. Lack of capacity is an error|
| `read_to_end(path: str, dst_buffer: ptr<Buffer>) -> i64` |Add file after existing Buffer and return additional amount|
| `write(path: str, src: ptr<u8>, len: i64) -> i64` |Number of bytes written replacing existing content|
| `append(path: str, src: ptr<u8>, len: i64) -> i64` |Number of bytes added to the end|
| `remove(path: str) -> i64` |removal status. failure is negative|

false of `exists(path)` alone cannot distinguish between missing files and permission errors. Be sure to check the actual open result, as the state may change between checking for existence and opening.

## Low level I/O

```text
std::io::fd
io_read(fd: i64, buf: ptr<u8>, len: i64) -> i64
io_write(fd: i64, buf: ptr<u8>, len: i64) -> i64
io_read_exact(fd: i64, buf: ptr<u8>, len: i64) -> i64
io_write_all(fd: i64, buf: ptr<u8>, len: i64) -> i64
io_close(fd: i64) -> i64
```

The positive result of `io_read` is the number of bytes read, and 0 in a positive length request is EOF. `io_write` may be written less than requested. If a full transfer is required, use the exact/all function. Still, we do not assume that the failure reverts the external state, as some transfers may have occurred before the error.

`io_read_exact` is an error if EOF is encountered before the required length. `read_into` returns `IO_ERR_NO_SPACE` if the buffer is full, and some bytes may have already been written. The read function does not automatically append NUL to the end of the string.

## Buffer and error handling

On failure, `read_to_end` restores the original len, but its capacity and data address may have changed. The caller must free the Buffer after either success or failure. File-writing APIs do not guarantee atomic file replacement.

From [File reading practice](/docs/en/practice/file-reader), you can run the program from import to release. Consider the path/permission differences in Linux/macOS/Windows/FreeBSD and the accessible directory restrictions in WASI. It does not directly interpret the descriptor value as a raw handle of another OS.

## Reading large files into small buffers

Operations that do not require the entire file to be placed in memory can be handled with fixed buffers and read iterations. The following program prints the contents of input.txt and counts the total number of bytes read. Save one `Wave` and LF in the input file.

<!-- wave-example: book-file-stream -->
```wave
import("std::fs::file")::{open_read};
import("std::io::fd")::{io_read, io_write_all, io_close};
import("std::io::consts")::{IO_STDOUT_FD};

fun main() -> i32 {
    var descriptor: i64 = open_read("input.txt");

    if (descriptor < 0) {
        return 1;
    }

    var buffer: array<u8, 4>;
    var total: i64 = 0;

    while (true) {
        var count: i64 = io_read(descriptor, &buffer[0], 4);

        if (count < 0) {
            io_close(descriptor);
            return 2;
        }

        if (count == 0) {
            break;
        }

        if (io_write_all(IO_STDOUT_FD, &buffer[0], count) < 0) {
            io_close(descriptor);
            return 3;
        }

        total += count;
    }

    if (io_close(descriptor) < 0) {
        return 4;
    }

    println("bytes={}", total);
    return 0;
}
```

Execution result:

```text
Wave
bytes=5
```

The buffer capacity is 4, but the last read may be 1 byte. We always pass the actual count to the output. If you write the entire array, even old, unread bytes can be output.

The program opened descriptor of input.txt, so close it. Standard output is not a newly acquired resource in this function, so it is not arbitrarily closed at the end of the example.

## Full reads and insufficient capacity

read_into receives a fixed storage space that holds the entire file. If space is insufficient, it silently truncates and returns NO_SPACE without success.

<!-- wave-example: book-file-capacity -->
```wave
import("std::fs::file")::{read_into};
import("std::io::consts")::{IO_ERR_NO_SPACE};

fun main() {
    var data: array<u8, 2>;
    var status: i64 = read_into("input.txt", &data[0], 2);

    if (status == IO_ERR_NO_SPACE) {
        println("destination too small");
    }
}
```

Execution result:

```text
destination too small
```

It does not assume that the failed read did not change the destination byte at all. Do not use it as the finished file content, prepare a larger repository or choose the streaming method. Even if you query the size first, the actual read result is the final judgment, as the file may change between query and read.

## Difference between writing and appending files

write replaces the existing content and append is added to the end. The number of bytes to store can be obtained directly from the string length and passed. The NUL at the end of the string is usually not included in the text file contents.

<!-- wave-example: book-file-write-append -->
```wave
import("std::fs::file")::{write, append, size, remove};
import("std::string::len")::{len};

fun main() -> i32 {
    var first: str = "Wave";
    var second: str = " study";

    if (write("output.txt", first as ptr<u8>, len(first) as i64) < 0) {
        return 1;
    }

    if (append("output.txt", second as ptr<u8>, len(second) as i64) < 0) {
        return 2;
    }

    println("bytes={}", size("output.txt"));

    if (remove("output.txt") < 0) {
        return 3;
    }

    return 0;
}
```

Execution result:

```text
bytes=10
```

This example creates, replaces, and finally deletes output.txt in the working directory. Run from the practice directory with no existing files. The actual editor or storage program may require separate storage policies, such as temporary files and replacement.

## API Selection table

|situation|select|
| --- | --- |
|Read small entire file into fixed buffer| read_into |
|Keep entire contents without knowing the size|read_to_end and Buffer|
|Process content in order rather than storing it in its entirety|open_read + io_read repeat|
|Read records of fixed length| io_read_exact |
|Transmit entire byte string| io_write_all |
|Handling already open files|fd function instead of path function|

After selecting a function, check how the buffer, file location, and external data change in the event of a failure.
