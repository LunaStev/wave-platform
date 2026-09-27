---
translation_set_id: practice-file-reader
path: practice/file-reader
locale: en
group: practice
group_order: 4
order: 2
title: Project: Reading a file and counting bytes
summary: Buffer, file I/O, links failure handling and memory freeing.
---

## ready

Create input.txt in the working directory containing `Wave` followed by a single LF newline. The file then contains 5 bytes. With CRLF it contains 6 bytes; a UTF-8 BOM adds further bytes. Check the editor’s file encoding and line endings.

Save the program as `main.wave` and run it as `wavec run main.wave` from the same directory. Relative paths are relative to the running working directory.

<!-- wave-example: file-reader -->
```wave
import("std::buffer::alloc")::{
    Buffer, buffer_init, buffer_free
};
import("std::fs::file")::{
    read_to_end
};

fun main() -> i32 {
    var data: Buffer;
    if (buffer_init(&data, 0) < 0) {
        return 1;
    }

    var count: i64 = read_to_end("input.txt", &data);
    if (count < 0) {
        println("read failed={}", count);
        buffer_free(&data);
        return 2;
    }

    var lines: i64 = 0;
    for (var i: i64 = 0; i < data.len; i += 1) {
        if (deref data.data[i] == 10) {
            lines += 1;
        }
    }

    println("bytes={} LF={}", count, lines);
    if (buffer_free(&data) < 0) {
        return 3;
    }

    return 0;
}
```

Execution result:

```text
bytes=5 LF=1
```

## Actions and Responsibilities

`read_to_end` opens and closes the file, but Buffer release is the caller's responsibility. Turns off on both success and read-failure paths. Since it is data with a number of bytes, we do not assume that it is a string ending with NUL.

LF The number of lines and the number of lines that people think are not always the same. If the last line does not contain LF, it is not included in the LF count for this program. Also check for read failures by renaming the file. Specific error numbers may vary depending on your environment.

## Extended exercises and commentary

If you want to handle very large files, replace it with a fixed size array and `io_read` iteration. Processes only positive return ranges and terminates at 0. You can accumulate byte counts and LF counts without having to keep the entire file in memory. If you opened it yourself, it also closes the descriptor on any exit path.

[See fs and io](/docs/en/stdlib/files-io) · [See Buffer](/docs/en/stdlib/buffer)

## Follow the processing flow

1. Initialize bin Buffer. There is no file content yet.
2. read_to_end reads the file and increases the required space.
3. If the read is successful, the bytes in the range data.len are checked.
4. Whenever we encounter the byte value 10 of LF, we increment lines.
5. Print the results and release Buffer.

data.cap is the reserved storage space and data.len is the valid data length. If you change the repeat condition to cap, bytes that were not in the file will be read, so use len. count is the number of bytes added in this call to read_to_end. This example starts with an empty Buffer, so count and data.len are equal.

## Change input to check

|input.txt Contents| bytes | LF |reason|
| --- | --- | --- | --- |
|empty file| 0 | 0 |No bytes to read|
| `Wave` | 4 | 0 |No final line break|
| `Wave` + LF | 5 | 1 |Data up to the last LF|
| `A` + LF + `B` + LF | 4 | 2 |Count two LF|
| `Wave` + CRLF | 6 | 1 |CR is also 1 byte, but only LF is counted.|

To count files without LF in the last line as a line, add 1 to the number of lines if the file is not empty and the last byte is not 10. You must first check if data.len is 0 before you can access the last element.

## Expand to large files

The current method of archiving the entire content is convenient for later rereading or retrieval of the data. If you only need the number of bytes and the LF count, reusing a fixed-size buffer makes sense.

In the io_read loop in the [Reading a file into a fixed size buffer](/docs/en/stdlib/files-io) example, just count LF by the number of bytes returned. It can be processed with memory equal to the buffer size, not the entire length of the file. If the read returns 0, it's over; if it's negative, it's an error.
