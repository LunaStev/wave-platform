---
translation_set_id: stdlib-bytes
path: stdlib/bytes
locale: en
group: stdlib
group_order: 1
order: 6
title: bytes: Ranges, cursors, and ULEB128
summary: Describes read/write bytes with length view and preserving state in case of failure.
---

## How is it different from a string?

Byte data can contain zero, so pass a pointer together with a length. `Bytes` and `BytesMut` are non-owning views and are valid only while the underlying storage remains valid. `BytesMut` requires writable storage.

## Create cursor

```text
std::bytes::cursor
bytes_reader(data: ptr<u8>, len: i64) -> ByteReader
bytes_writer(data: ptr<u8>, len: i64) -> ByteWriter
bytes_reader_read_be_u16(reader: ptr<ByteReader>, out_value: ptr<u16>) -> i32
bytes_writer_write_be_u16(writer: ptr<ByteWriter>, value: u16) -> i32
```

Import `ByteReader` and `ByteWriter` from `std::bytes::types`. Their position field identifies the next operation’s location. Their length is the total accessible byte count. Constructing either cursor does not copy or allocate the underlying memory.

`be` is big-endian, `le` is little-endian. If the file type is big-endian, use the `be` function regardless of the byte order of the host CPU. There are 16, 32, and 64 bit signed/unsigned read/write and one-byte functions.

## Errors and State Preservation

`BYTES_OK` from `std::bytes::errors` is 0. INVALID indicates an invalid range, EOF insufficient input, NO_SPACE insufficient output capacity, and OVERFLOW a value outside the representable range. Checked cursor operations advance position only after the entire operation succeeds. A failed read also preserves the output value.

## ULEB128

```text
std::bytes::leb128
bytes_reader_read_uleb128_u64(reader: ptr<ByteReader>, output_value: ptr<u64>) -> i32
bytes_writer_write_uleb128_u64(writer: ptr<ByteWriter>, value: u64) -> i32
```

ULEB128 stores an unsigned 64-bit integer in a variable number of bytes, using at most 10 bytes. If there is insufficient space, the writer preserves both its position and the destination bytes. The reader distinguishes incomplete input from a value exceeding u64. Terminated, non-minimal encodings are accepted.

To create an actual message and see a short input failure, proceed to [Binary message practice](/docs/en/practice/binary-message). Do not attempt to output a byte string containing zeros as `str`.

## Reading the same bytes in different orders

Byte order is the storage rule for numbers. If you read the two bytes 1 and 2 as big-endian, it is 1×256+2, and if you read them as little-endian, it is 2×256+1. Select based on network or file type rules.

<!-- wave-example: book-bytes-endian -->
```wave
import("std::bytes::types")::{Bytes, bytes_view};
import("std::bytes::read")::{bytes_read_be_u16, bytes_read_le_u16};

fun main() -> i32 {
    var data: array<u8, 2> = [1, 2];
    var view: Bytes = bytes_view(&data[0], 2);
    var big: u16 = 0;
    var little: u16 = 0;

    if (bytes_read_be_u16(view, 0, &big) < 0) {
        return 1;
    }

    if (bytes_read_le_u16(view, 0, &little) < 0) {
        return 2;
    }

    println("be={} le={}", big, little);
    return 0;
}
```

Execution result:

```text
be=258 le=513
```

view borrows an array. Since there is no separate allocation or copy, the array can only be used while it is valid. offset is in bytes and reading 16 bits requires 2 bytes from that position.

## Select offset API and cursor API

The read/write function, which takes the offset argument, is convenient for formats that directly read a specified field position. For streams where the next position depends on the length of the preceding field, cursor with position is convenient.

When mixing the two, make it clear which is the standard: cursor.position or separate offset. Avoid the mistake of adding the same location twice or moving to the next location without a successful read.

## Check status from short input

<!-- wave-example: book-bytes-short-read -->
```wave
import("std::bytes::types")::{ByteReader};
import("std::bytes::cursor")::{bytes_reader, bytes_reader_read_be_u16};
import("std::bytes::errors")::{BYTES_ERROR_EOF};

fun main() {
    var data: array<u8, 1> = [1];
    var reader: ByteReader = bytes_reader(&data[0], 1);
    var value: u16 = 99;
    var status: i32 = bytes_reader_read_be_u16(&reader, &value);

    if (status == BYTES_ERROR_EOF) {
        println("position={} value={}", reader.position, value);
    }
}
```

Execution result:

```text
position=0 value=99
```

There are no two bytes needed, so the output value and location are preserved. This property is useful in designs that retry reading the same field after acquiring more input. However, if the storage space pointed to by view has been reallocated, the address must also be updated.

## Boundary of ULEB128

0~127 uses one byte, 128 onwards uses more bytes. The high-order bit of each byte indicates whether data follows. A value outside of u64 or a continuous input that is too long is OVERFLOW, which is different from EOF, which simply has less input.

Directly check if Capacity Insufficient writer preserves its state.

<!-- wave-example: book-leb128-space -->
```wave
import("std::bytes::types")::{ByteWriter};
import("std::bytes::cursor")::{bytes_writer};
import("std::bytes::leb128")::{bytes_writer_write_uleb128_u64};
import("std::bytes::errors")::{BYTES_ERROR_NO_SPACE};

fun main() {
    var data: array<u8, 1> = [85];
    var writer: ByteWriter = bytes_writer(&data[0], 1);
    var status: i32 = bytes_writer_write_uleb128_u64(&writer, 128);

    if (status == BYTES_ERROR_NO_SPACE) {
        println("position={} byte={}", writer.position, data[0]);
    }
}
```

Execution result:

```text
position=0 byte=85
```

128 requires two bytes but only one space. After a failure, the first byte 85 also remains. Checking error types and state preservation together describes the boundary better than a simple roundtrip success check.

## Message parser creation sequence

1. Read the fixed header and check the type and version.
2. Read the length and compare it to the remaining input range.
3. Pass only the necessary data to view or a separate buffer.
4. If the format requires the entire message, the extra bytes are also checked.
5. Distinguishes between EOF and invalid format errors and passes them on to the caller.

You can create one program by connecting fields in [Binary message practice](/docs/en/practice/binary-message).
