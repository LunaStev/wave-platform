---
translation_set_id: stdlib-files-io
path: stdlib/files-io
locale: ko
group: stdlib
group_order: 1
order: 7
title: fs와 io: 파일과 바이트 전송
summary: 파일 수명, 전체 읽기, 부분 전송과 실패 후 상태를 설명합니다.
---

## 파일 API 선택

`std::fs::file`의 편의 함수는 경로를 받아 필요한 열기·닫기를 수행합니다. 디스크립터를 반환하는 함수는 호출자가 닫아야 합니다.

| 선언 | 성공 결과와 주의점 |
| --- | --- |
| `open_read(path: str) -> i64` | 열린 디스크립터. 음수는 오류 |
| `create(path: str) -> i64` | 생성 또는 기존 파일 내용 지움. 소유 디스크립터 반환 |
| `open_append(path: str) -> i64` | 추가용 열기 또는 생성 |
| `size(path: str) -> i64` | 바이트 수. 음수는 오류 |
| `read_into(path: str, dst: ptr<u8>, dst_cap: i64) -> i64` | 전체 파일의 바이트 수. 용량 부족은 오류 |
| `read_to_end(path: str, dst_buffer: ptr<Buffer>) -> i64` | 기존 Buffer 뒤에 파일을 추가하고 추가량 반환 |
| `write(path: str, src: ptr<u8>, len: i64) -> i64` | 기존 내용을 대체하여 쓴 바이트 수 |
| `append(path: str, src: ptr<u8>, len: i64) -> i64` | 끝에 추가한 바이트 수 |
| `remove(path: str) -> i64` | 제거 상태. 실패는 음수 |

`exists(path)`의 false만으로 없는 파일과 권한 오류를 구별할 수 없습니다. 존재를 검사한 뒤 여는 사이에도 상태가 바뀔 수 있으므로 실제 열기 결과를 반드시 검사합니다.

## 낮은 수준 I/O

```text
std::io::fd
io_read(fd: i64, buf: ptr<u8>, len: i64) -> i64
io_write(fd: i64, buf: ptr<u8>, len: i64) -> i64
io_read_exact(fd: i64, buf: ptr<u8>, len: i64) -> i64
io_write_all(fd: i64, buf: ptr<u8>, len: i64) -> i64
io_close(fd: i64) -> i64
```

`io_read`의 양수 결과는 읽은 바이트 수이고, 양수 길이 요청에서 0은 EOF입니다. `io_write`는 요청보다 적게 쓸 수 있습니다. 전체 전송이 필요하면 exact/all 함수를 사용합니다. 그래도 오류 전에 일부 전송이 일어날 수 있으므로 실패가 외부 상태를 되돌린다고 가정하지 않습니다.

`io_read_exact`는 필요한 길이 전에 EOF를 만나면 오류입니다. `read_into`는 버퍼가 모자라면 `IO_ERR_NO_SPACE`를 반환하며 일부 바이트가 이미 기록되었을 수 있습니다. 읽기 함수는 문자열 끝의 NUL을 자동으로 붙이지 않습니다.

## Buffer와 오류 처리

`read_to_end`는 실패 시 원래 len을 복원합니다. 그러나 용량과 data 주소까지 이전과 같다고 가정하지 않습니다. 성공·실패와 무관하게 호출자가 Buffer를 해제합니다. 파일을 쓰는 API는 원자적 파일 교체를 보장하지 않습니다.

[파일 읽기 실습](/docs/ko/practice/file-reader)에서 import부터 해제까지 이어진 프로그램을 실행할 수 있습니다. Linux/macOS/Windows/FreeBSD의 경로·권한 차이와 WASI의 접근 가능한 디렉터리 제한을 고려하십시오. 디스크립터 값을 다른 OS의 원시 핸들로 직접 해석하지 않습니다.

## 큰 파일을 작은 버퍼로 읽기

전체 파일을 메모리에 올릴 필요가 없는 작업은 고정 버퍼와 읽기 반복으로 처리할 수 있습니다. 다음 프로그램은 input.txt의 내용을 출력하고 읽은 전체 바이트 수를 셉니다. 입력 파일에는 `Wave`와 LF 하나를 저장합니다.

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

실행 결과:

```text
Wave
bytes=5
```

버퍼 용량은 4지만 마지막 읽기는 1바이트일 수 있습니다. 출력에는 항상 실제 count를 전달합니다. 배열 전체를 쓰면 읽지 않은 오래된 바이트까지 출력할 수 있습니다.

프로그램이 연 것은 input.txt의 descriptor이므로 그것을 닫습니다. 표준 출력은 이 함수에서 새로 얻은 소유 자원이 아니므로 예제 끝에서 임의로 닫지 않습니다.

## 전체 읽기와 용량 부족

read_into는 파일 전체를 담는 고정 저장 공간을 받습니다. 공간이 모자라면 조용히 잘라서 성공하지 않고 NO_SPACE를 반환합니다.

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

실행 결과:

```text
destination too small
```

실패한 읽기가 목적지 바이트를 전혀 변경하지 않았다고 가정하지 않습니다. 완성된 파일 내용으로 사용하지 말고 더 큰 저장소를 준비하거나 streaming 방식을 선택합니다. 크기를 먼저 조회해도 조회와 읽기 사이에 파일이 바뀔 수 있으므로 실제 읽기 결과가 최종 판단 기준입니다.

## 파일 쓰기와 추가의 차이

write는 기존 내용을 대체하고 append는 끝에 추가합니다. 저장할 바이트 수는 문자열 길이에서 직접 구해 전달할 수 있습니다. 문자열 끝의 NUL은 보통 텍스트 파일 내용에 포함하지 않습니다.

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

실행 결과:

```text
bytes=10
```

이 예제는 작업 디렉터리의 output.txt를 생성·대체하고 마지막에 삭제합니다. 기존 파일이 없는 연습 디렉터리에서 실행하십시오. 실제 편집기나 저장 프로그램에는 임시 파일과 교체 같은 별도 저장 정책이 필요할 수 있습니다.

## API 선택표

| 상황 | 선택 |
| --- | --- |
| 작은 전체 파일을 고정 버퍼로 읽기 | read_into |
| 크기를 모르고 전체 내용을 보관 | read_to_end와 Buffer |
| 내용 전체를 보관하지 않고 순서대로 처리 | open_read + io_read 반복 |
| 정해진 길이의 레코드를 읽기 | io_read_exact |
| 바이트열 전체 전송 | io_write_all |
| 이미 열린 파일을 다루기 | 경로 함수 대신 fd 함수 |

함수 선택 후에는 해당 실패에서 버퍼·파일 위치·외부 데이터가 어떻게 바뀌는지도 확인합니다.
