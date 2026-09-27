---
translation_set_id: stdlib-bytes
path: stdlib/bytes
locale: ko
group: stdlib
group_order: 1
order: 6
title: bytes: 범위·cursor·ULEB128
summary: 길이를 가진 바이트 view와 실패 시 상태를 보존하는 읽기·쓰기를 설명합니다.
---

## 문자열과 다른 점

바이트 데이터는 0을 포함할 수 있으므로 포인터와 길이를 함께 전달합니다. `Bytes`와 `BytesMut`는 메모리를 소유하지 않는 view입니다. 원본이 살아 있는 동안만 사용할 수 있습니다. `BytesMut`에는 실제로 쓰기 가능한 저장소를 제공해야 합니다.

## cursor 만들기

```text
std::bytes::cursor
bytes_reader(data: ptr<u8>, len: i64) -> ByteReader
bytes_writer(data: ptr<u8>, len: i64) -> ByteWriter
bytes_reader_read_be_u16(reader: ptr<ByteReader>, out_value: ptr<u16>) -> i32
bytes_writer_write_be_u16(writer: ptr<ByteWriter>, value: u16) -> i32
```

`ByteReader`와 `ByteWriter`는 `std::bytes::types`에서 가져옵니다. 두 타입의 position은 다음 작업 위치입니다. 길이는 전체 접근 가능한 바이트 수이며, 생성 자체는 메모리를 복사하거나 할당하지 않습니다.

`be`는 big-endian, `le`는 little-endian입니다. 파일 형식이 big-endian이면 호스트 CPU의 바이트 순서와 관계없이 `be` 함수를 사용합니다. 16·32·64비트 signed/unsigned 읽기·쓰기와 한 바이트 함수가 있습니다.

## 오류와 상태 보존

`std::bytes::errors`의 `BYTES_OK`는 0입니다. 잘못된 범위는 INVALID, 읽을 데이터 부족은 EOF, 쓸 공간 부족은 NO_SPACE, 표현 범위 초과는 OVERFLOW입니다. checked cursor 작업은 전체 작업이 성공해야 position을 진행합니다. 실패한 읽기는 출력값도 유지합니다.

## ULEB128

```text
std::bytes::leb128
bytes_reader_read_uleb128_u64(reader: ptr<ByteReader>, output_value: ptr<u64>) -> i32
bytes_writer_write_uleb128_u64(writer: ptr<ByteWriter>, value: u64) -> i32
```

부호 없는 64비트 정수를 가변 길이로 저장하며 최대 10바이트를 사용합니다. writer는 필요한 공간이 없으면 위치와 목적지 바이트를 바꾸지 않습니다. reader는 끝나지 않은 입력과 u64 범위 초과를 구분합니다. 종료되는 비최소 인코딩은 허용합니다.

실제 메시지를 만들고 짧은 입력 실패를 확인하려면 [바이너리 메시지 실습](/docs/ko/practice/binary-message)을 진행하십시오. 0을 포함한 바이트열을 `str`로 출력하려고 하지 마십시오.

## 같은 바이트를 다른 순서로 읽기

바이트 순서는 숫자의 저장 규칙입니다. 1, 2라는 두 바이트를 big-endian으로 읽으면 1×256+2이고 little-endian으로 읽으면 2×256+1입니다. 네트워크나 파일 형식의 규칙을 따라 선택합니다.

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

실행 결과:

```text
be=258 le=513
```

view는 배열을 빌립니다. 별도의 할당이나 복사가 없으므로 배열이 유효한 동안에만 사용할 수 있습니다. offset은 바이트 단위이며 16비트 읽기에는 해당 위치부터 2바이트가 필요합니다.

## offset API와 cursor API 선택

정해진 필드 위치를 직접 읽는 형식에는 offset 인자를 받는 read/write 함수가 편합니다. 앞의 필드 길이에 따라 다음 위치가 달라지는 스트림에는 position을 가진 cursor가 편합니다.

둘을 섞을 때는 cursor.position과 별도 offset 중 어떤 것이 기준인지 분명히 합니다. 같은 위치를 두 번 더하거나 읽기를 성공하지 않았는데 다음 위치로 이동하는 실수를 피하십시오.

## 짧은 입력에서 상태 확인하기

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

실행 결과:

```text
position=0 value=99
```

필요한 2바이트가 없으므로 출력값과 위치가 유지됩니다. 이 성질은 입력을 더 확보한 뒤 같은 필드 읽기를 다시 시도하는 설계에 유용합니다. 다만 view가 가리키는 저장 공간을 재할당했다면 주소도 갱신해야 합니다.

## ULEB128의 경계

0~127은 한 바이트, 128부터는 더 많은 바이트를 사용합니다. 각 바이트의 상위 비트는 뒤에 데이터가 이어지는지 표시합니다. u64를 벗어나는 값이나 너무 긴 연속 입력은 OVERFLOW이며 단순히 입력이 덜 온 EOF와 다릅니다.

용량 부족 writer가 상태를 보존하는지 직접 확인합니다.

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

실행 결과:

```text
position=0 byte=85
```

128은 두 바이트가 필요한데 공간은 하나뿐입니다. 실패 후 첫 바이트 85도 남습니다. 오류 종류와 상태 보존을 함께 확인하는 것이 단순한 roundtrip 성공 검사보다 경계를 잘 설명합니다.

## 메시지 파서 작성 순서

1. 고정 헤더를 읽고 종류와 버전을 검사합니다.
2. 길이를 읽고 남은 입력 범위와 비교합니다.
3. 필요한 데이터만 view 또는 별도 버퍼로 전달합니다.
4. 전체 메시지를 요구하는 형식이면 남는 바이트도 검사합니다.
5. EOF와 잘못된 형식의 오류를 구분해 호출자에게 전달합니다.

[바이너리 메시지 실습](/docs/ko/practice/binary-message)에서 필드를 연결해 하나의 프로그램을 만들 수 있습니다.
