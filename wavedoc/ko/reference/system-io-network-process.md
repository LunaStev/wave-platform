---
translation_set_id: system-io
path: reference/system-io-network-process
locale: ko
group: stdlib
group_order: 1
order: 15
title: 시스템 기능과 프로세스
summary: 상위 API와 OS 인터페이스의 경계, 프로세스 수명을 설명합니다.
---

## 기능별 문서

파일을 다루려면 [fs와 io](/docs/ko/stdlib/files-io), 연결하려면 [TCP](/docs/ko/stdlib/tcp), 주소를 조회하려면 [resolver](/docs/ko/stdlib/resolution)를 읽으십시오. 아래는 프로세스 및 더 낮은 수준의 OS 접근 규칙입니다.

## 프로세스 기본 API

```text
std::process::core
proc_exit(code: i32) -> !
proc_getpid() -> i64
proc_getppid() -> i64
proc_execve(path: str, argv: ptr<ptr<i8>>, envp: ptr<ptr<i8>>) -> i64
proc_waitpid_raw(pid: i64, status: ptr<i32>, options: i32) -> i64
```

`proc_exit`은 호출 지점으로 돌아오지 않습니다. 종료 전에 필요한 파일·메모리 정리를 직접 수행하십시오. `proc_execve`는 성공하면 기존 프로세스 이미지를 바꾸므로 일반적인 자식 생성 함수와 다릅니다. raw argv/envp는 끝을 나타내는 null 포인터와 각 문자열의 NUL 종료를 준비해야 합니다.

`std::process::spawn`의 spawn 함수는 생성 결과를, 대기 함수는 자식의 종료 상태를 다룹니다. 생성 성공과 프로그램의 성공 종료는 서로 다릅니다. 파이프를 만들면 부모와 자식이 사용하지 않는 끝을 닫아야 EOF가 전달됩니다. 캡처 파이프를 읽지 않은 채 자식 종료만 기다리면 버퍼가 차서 서로 기다릴 수 있습니다.

## 이식성과 낮은 수준 접근

fork/exec, 파일 디스크립터, Windows 핸들은 동일한 OS 기능이 아닙니다. 선택한 대상의 지원을 확인하고 unsupported를 정상적인 실패 경로로 처리합니다. `std::sys`는 OS별 인터페이스이며 숫자 플래그와 레이아웃을 다른 OS에서 재사용하지 않습니다.

외부 C 라이브러리와 직접 연동할 때는 [FFI](/docs/ko/language/modules-imports-and-ffi)를 읽으십시오. 상위 std API를 쓰기 위해 임의로 libc 함수를 선언할 필요는 없습니다. [타깃과 링크 환경](/docs/ko/whale/build-link-targets)을 먼저 확인합니다.
