---
translation_set_id: practice-binary-message
path: practice/binary-message
locale: en
group: practice
group_order: 4
order: 3
title: Project: Creating a binary message
summary: Use explicit byte ordering and ULEB128 and reject short input.
---

## message format

The first 2 bytes store the type numbers big-endian u16, and then the values ULEB128 u64. If you write structure memory to a file as is, it will be affected by padding and byte order, so encode it by field.

Save it as `main.wave` and run it.

<!-- wave-example: binary-message -->
```wave
import("std::bytes::types")::{
    ByteReader, ByteWriter
};
import("std::bytes::cursor")::{
    bytes_reader, bytes_writer, bytes_writer_write_be_u16, bytes_reader_read_be_u16
};
import("std::bytes::leb128")::{
    bytes_writer_write_uleb128_u64, bytes_reader_read_uleb128_u64
};

fun main() -> i32 {
    var data: array<u8, 12>;
    var writer: ByteWriter = bytes_writer(&data[0], 12);
    if (bytes_writer_write_be_u16(&writer, 7) < 0) {
        return 1;
    }
    if (bytes_writer_write_uleb128_u64(&writer, 300) < 0) {
        return 2;
    }

    var reader: ByteReader = bytes_reader(&data[0], writer.position);
    var kind: u16 = 0;
    var value: u64 = 0;
    if (bytes_reader_read_be_u16(&reader, &kind) < 0) {
        return 3;
    }
    if (bytes_reader_read_uleb128_u64(&reader, &value) < 0) {
        return 4;
    }

    println("kind={} value={} bytes={}", kind, value, reader.position);
    var short: ByteReader = bytes_reader(&data[0], 1);
    kind = 99;
    if (bytes_reader_read_be_u16(&short, &kind) >= 0) {
        return 5;
    }

    println("preserved={} {}", short.position, kind);
    return 0;
}
```

Execution result:

```text
kind=7 value=300 bytes=4
preserved=0 99
```

## What to check

In reader, we pass the actual written length, not the total array capacity of 12. This is to avoid reading uninitialized trailing bytes as input. Even if a short input read fails, position=0 and kind=99 are maintained.

## Extended exercises and commentary

If the format does not allow extra bytes at the end, check `reader.position == reader.len` after parsing is complete. When adding a length field, make sure it is no larger than the remaining bytes of the input, and that the calculation of length+offset does not exceed the range.

[See bytes](/docs/en/stdlib/bytes)

## Look at the actual bytes

Type number 7 is big-endian u16 and therefore `00 07`. The value 300 becomes ULEB128 to `AC 02`. The entire message is the following four bytes:

```text
00 07 AC 02
───── ─────
종류  값
```

ULEB128 uses the low 7 bits for each byte for the value, and if the high bit is 1, it indicates that the next byte follows. The low 7 bits of 300 are 44 and the remainder is 2. The first byte is 44 plus 128 consecutive marks, or 172, or 0xAC. There is no continuation mark in the last byte 0x02.

## Capacity and length of use

u16 uses 2 bytes, and ULEB128 of u64 uses up to 10 bytes, so a 12-byte array can store both fields. A value of 300 uses only 2 bytes, making the actual message 4 bytes. When sending to a file or socket, you send the byte writer.position rather than the entire array.

position in reader is the current read position. After reading the type, it becomes 2, and after reading the value, it becomes 4. In order to proceed to the next message when an error occurs, the boundary of the failed message must be known. A message recovery policy is not automatically established simply because the read function preserves the location.

## Boundary value exercise

Change the values to 0, 127, 128, 16383, 16384 to determine the encoding length. When changing from 127 to 128, the length of ULEB128 increases from 1 to 2, and when changing from 16383 to 16384, the length increases from 2 to 3. Calculate the total length, including the 2 bytes of the type field.

Messages with the last byte truncated are also checked. If the original message length was 4, we pass length 3 to reader. The type field is read, but reading the value must fail because the last byte of ULEB128 is missing. At this time, check whether position 2 and the output value just before reading ULEB128 are maintained.
