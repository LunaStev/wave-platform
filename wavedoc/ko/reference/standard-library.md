---
translation_set_id: standard-library
path: reference/standard-library
locale: ko
group: stdlib
group_order: 1
order: 1
title: 표준 라이브러리 길잡이
summary: 목적에 맞는 모듈을 찾고 함수의 오류와 소유권 규칙을 읽는 방법입니다.
---

## 필요한 기능 찾기

표준 라이브러리는 `std::모듈::파일` 경로로 import합니다. 이름이 비슷해도 함수의 오류 반환 방식은 다를 수 있습니다. 처음에는 [API 읽는 법](/docs/ko/stdlib/contracts)을 읽고, 다음 표에서 필요한 모듈로 이동하십시오.

| 하고 싶은 일 | 문서 | 주요 import |
| --- | --- | --- |
| 문자열 길이·비교·검색 | [string](/docs/ko/reference/string-and-bytes) | `std::string::len`, `cmp`, `find`, `trim` |
| 메모리 할당·복사·크기 | [mem](/docs/ko/reference/memory-and-buffer) | `std::mem::alloc`, `ops`, `layout` |
| 크기가 변하는 바이트 목록 | [buffer](/docs/ko/stdlib/buffer) | `std::buffer::alloc`, `read`, `write` |
| 바이너리 읽기·쓰기 | [bytes](/docs/ko/stdlib/bytes) | `std::bytes::types`, `cursor`, `leb128` |
| 파일·디스크립터 I/O | [fs와 io](/docs/ko/stdlib/files-io) | `std::fs::file`, `std::io::fd` |
| 경로 조합·환경 설정 | [path와 env](/docs/ko/stdlib/path-env) | `std::path::copy`, `std::env::environ` |
| 시간 측정·대기 | [time](/docs/ko/stdlib/time) | `std::time::duration`, `clock`, `sleep` |
| 숫자 주소·이름 목록 조회 | [net.resolve](/docs/ko/stdlib/resolution) | `std::net::resolve`, `resolve_table` |
| TCP 연결·전송 | [net.tcp](/docs/ko/stdlib/tcp) | `std::net::tcp`, `address`, `error` |
| OS 난수 | [random](/docs/ko/stdlib/random) | `std::random::fill` |
| 프로세스·OS 경계 | [시스템 기능](/docs/ko/reference/system-io-network-process) | `std::process::core`, `spawn`, `std::sys` |
| 비동기 작업 실행 | [task](/docs/ko/stdlib/task) | `std::task` |
| 수학·진단 도우미 | [math와 debug](/docs/ko/stdlib/math-debug) | `std::math::int`, `float`, `std::debug::core` |

## 처음 사용하는 예

아래 프로그램은 별도 패키지를 다운로드하지 않고 std의 함수 하나를 사용합니다. `main.wave`로 저장하여 `wavec run main.wave`로 실행합니다.

<!-- wave-example: library-import -->
```wave
import("std::string::len")::{
    len
};

fun main() {
    println("{}", len("Wave"));
}
```

결과는 `4`입니다. 같은 예제를 단계별로 이해하려면 [문자열](/docs/ko/language/strings)을 읽으십시오.

## 컴파일러와 맞는 std

`wavec print std-path`로 선택한 경로를 확인합니다. 다른 체크아웃의 std를 사용할 때는 `wavec --std-root /absolute/path/to/std check main.wave`처럼 경로를 명시합니다. 명시한 경로가 유효하지 않거나 호환되지 않으면 오류를 표시합니다.

## 플랫폼 경계

문자열·바이트 같은 계산 기능과 파일·소켓 같은 OS 기능을 구분하십시오. 타깃을 인식한다고 모든 호스트 API가 제공되는 것은 아닙니다. [지원 타깃](/docs/ko/whale/build-link-targets)과 각 API의 플랫폼 항목을 함께 읽습니다. `std::sys`는 더 낮은 수준의 인터페이스이며 이식 가능한 프로그램은 먼저 상위 모듈을 사용합니다.
