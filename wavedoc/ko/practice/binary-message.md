---
translation_set_id: practice-binary-message
path: practice/binary-message
locale: ko
group: practice
group_order: 4
order: 3
title: 실습: 바이너리 메시지 만들기
summary: 명시적 바이트 순서와 ULEB128을 사용하고 짧은 입력을 거부합니다.
---

## 메시지 형식

첫 2바이트에는 big-endian u16 종류 번호, 다음에는 ULEB128 u64 값을 저장합니다. 구조체 메모리를 그대로 파일에 쓰면 패딩과 바이트 순서에 영향을 받으므로 필드별로 인코딩합니다.

`main.wave`로 저장한 뒤 실행합니다.

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

실행 결과:

```text
kind=7 value=300 bytes=4
preserved=0 99
```

## 확인할 점

reader에는 배열 전체 용량 12가 아니라 실제로 쓴 길이를 전달합니다. 초기화하지 않은 뒤쪽 바이트를 입력으로 읽지 않기 위해서입니다. 짧은 입력 읽기가 실패해도 position=0과 kind=99가 유지됩니다.

## 확장 연습과 해설

끝에 남는 바이트를 허용하지 않는 형식이라면 파싱 완료 후 `reader.position == reader.len`을 검사하십시오. 길이 필드를 추가할 때는 입력의 남은 바이트보다 크지 않은지 확인하고, 길이+offset 계산이 범위를 넘지 않도록 해야 합니다.

[bytes 참조](/docs/ko/stdlib/bytes)

## 실제 바이트 살펴보기

종류 번호 7은 big-endian u16이므로 `00 07`입니다. 값 300은 ULEB128에서 `AC 02`가 됩니다. 전체 메시지는 다음 네 바이트입니다.

```text
00 07 AC 02
───── ─────
종류  값
```

ULEB128은 바이트마다 낮은 7비트를 값에 사용하고, 높은 비트가 1이면 다음 바이트가 이어짐을 나타냅니다. 300의 낮은 7비트는 44이고 나머지는 2입니다. 첫 바이트는 44에 연속 표시 128을 더한 172, 즉 0xAC입니다. 마지막 바이트 0x02에는 연속 표시가 없습니다.

## 용량과 사용 길이

u16은 2바이트, u64의 ULEB128은 최대 10바이트를 사용하므로 12바이트 배열이면 두 필드를 저장할 수 있습니다. 값 300은 2바이트만 사용해 실제 메시지는 4바이트입니다. 파일이나 소켓으로 전송할 때는 배열 전체가 아니라 writer.position 바이트를 보냅니다.

reader의 position은 현재 읽기 위치입니다. 종류를 읽은 뒤 2, 값을 읽은 뒤 4가 됩니다. 중간에 실패했을 때 다음 메시지로 진행하려면 먼저 실패한 메시지의 경계를 알 수 있어야 합니다. 읽기 함수가 위치를 보존해 주는 것만으로 메시지 복구 정책이 자동으로 정해지지는 않습니다.

## 경계값 연습

값을 0, 127, 128, 16383, 16384로 바꿔 인코딩 길이를 확인합니다. 127에서 128로 바뀔 때 ULEB128 길이는 1에서 2로, 16383에서 16384로 바뀔 때는 2에서 3으로 늘어납니다. 종류 필드 2바이트도 포함해 전체 길이를 계산하십시오.

마지막 바이트가 잘린 메시지도 검사합니다. 원래 메시지 길이가 4라면 reader에 길이 3을 전달합니다. 종류 필드까지는 읽히지만 ULEB128의 마지막 바이트가 없어 값 읽기는 실패해야 합니다. 이때 ULEB128 읽기 직전의 위치 2와 출력값이 유지되는지 확인합니다.
