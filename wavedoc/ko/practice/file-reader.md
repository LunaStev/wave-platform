---
translation_set_id: practice-file-reader
path: practice/file-reader
locale: ko
group: practice
group_order: 4
order: 2
title: 실습: 파일을 끝까지 읽고 바이트 세기
summary: Buffer, 파일 I/O, 실패 처리와 메모리 해제를 연결합니다.
---

## 준비

작업 디렉터리에 `input.txt`를 만들고 `Wave` 뒤에 LF 줄바꿈 하나를 저장합니다. 이때 파일은 5바이트입니다. CRLF로 저장하면 6바이트, UTF-8 BOM을 붙이면 추가 바이트가 있으므로 편집기의 저장 형식을 확인하십시오.

프로그램을 `main.wave`에 저장하고 같은 디렉터리에서 `wavec run main.wave`로 실행합니다. 상대 경로는 실행하는 작업 디렉터리를 기준으로 합니다.

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

실행 결과:

```text
bytes=5 LF=1
```

## 동작과 책임

`read_to_end`가 파일을 열고 닫지만 Buffer 해제는 호출자의 책임입니다. 성공과 읽기 실패 경로 모두에서 해제합니다. 바이트 수를 가진 데이터이므로 NUL 종료 문자열이라고 가정하지 않습니다.

LF 개수와 사람이 생각하는 줄 수는 항상 같지 않습니다. 마지막 줄에 LF가 없으면 이 프로그램의 LF 개수에는 포함되지 않습니다. 파일 이름을 바꾸어 읽기 실패도 확인하십시오. 구체적인 오류 숫자는 환경에 따라 다를 수 있습니다.

## 확장 연습과 해설

매우 큰 파일을 처리하려면 고정 크기 배열과 `io_read` 반복으로 바꾸십시오. 양수 반환 범위만 처리하고 0에서 종료합니다. 전체 파일을 메모리에 유지하지 않아도 바이트 수와 LF 수를 누적할 수 있습니다. 직접 열었다면 모든 종료 경로에서 디스크립터도 닫습니다.

[fs와 io 참조](/docs/ko/stdlib/files-io) · [Buffer 참조](/docs/ko/stdlib/buffer)

## 처리 흐름을 따라가기

1. 빈 Buffer를 초기화합니다. 아직 파일 내용은 없습니다.
2. read_to_end가 파일을 읽으며 필요한 공간을 늘립니다.
3. 읽기에 성공하면 data.len 범위의 바이트를 검사합니다.
4. LF의 바이트 값 10을 만날 때마다 lines를 늘립니다.
5. 결과를 출력하고 Buffer를 해제합니다.

data.cap은 확보한 저장 공간이고 data.len은 유효한 데이터 길이입니다. 반복 조건을 cap으로 바꾸면 파일에 없던 바이트까지 읽게 되므로 len을 사용합니다. count는 이번 read_to_end 호출에서 추가한 바이트 수입니다. 이 예제는 빈 Buffer에서 시작하므로 count와 data.len이 같습니다.

## 입력을 바꿔 확인하기

| input.txt 내용 | bytes | LF | 이유 |
| --- | --- | --- | --- |
| 빈 파일 | 0 | 0 | 읽을 바이트가 없음 |
| `Wave` | 4 | 0 | 마지막 줄바꿈이 없음 |
| `Wave` + LF | 5 | 1 | 마지막 LF까지 데이터 |
| `A` + LF + `B` + LF | 4 | 2 | 두 LF를 셈 |
| `Wave` + CRLF | 6 | 1 | CR도 1바이트지만 LF만 셈 |

마지막 줄에 LF가 없는 파일도 한 줄로 세려면, 파일이 비어 있지 않고 마지막 바이트가 10이 아닐 때 줄 수에 1을 더합니다. data.len이 0인지 먼저 확인해야 마지막 원소에 접근할 수 있습니다.

## 큰 파일로 확장하기

전체 내용을 보관하는 현재 방식은 나중에 데이터를 다시 읽거나 검색할 때 편리합니다. 바이트 수와 LF 수만 필요하다면 고정 크기 버퍼를 재사용하는 방식이 적합합니다.

[고정 크기 버퍼로 파일 읽기](/docs/ko/stdlib/files-io) 예제의 io_read 반복문 안에서, 반환된 바이트 수만큼 LF를 세면 됩니다. 파일 전체 길이가 아니라 버퍼 크기만큼의 메모리로 처리할 수 있습니다. 읽기가 0을 반환하면 끝이고 음수면 오류입니다.
