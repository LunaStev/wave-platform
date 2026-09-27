---
translation_set_id: stdlib-tcp
path: stdlib/tcp
locale: ko
group: stdlib
group_order: 1
order: 9
title: net.tcp: 연결과 전송
summary: TCP 결과 구조체, 부분 전송과 연결 해제 책임을 설명합니다.
---

## 연결 결과 검사

```text
std::net::tcp
tcp_connect_addr(addr: SocketAddr) -> NetResult<TcpStream>
tcp_bind_loopback(port: u16) -> NetResult<TcpListener>
tcp_accept(listener: TcpListener) -> NetResult<TcpStream>
tcp_read(stream: TcpStream, buf: ptr<u8>, size: i64) -> i64
tcp_write_all(stream: TcpStream, buf: ptr<u8>, size: i64) -> i64
tcp_read_exact(stream: TcpStream, buf: ptr<u8>, size: i64) -> i64
tcp_close(stream: TcpStream) -> NetError
tcp_close_listener(listener: TcpListener) -> NetError
```

연결 함수의 `NetResult<T>`는 `ok`, `value`, `error`를 가집니다. 성공한 경우에만 value를 사용합니다. `NetError`는 정규화된 오류 분류와 native_code를 담으며 `error.kind == NET_ERROR_NONE`으로 성공 여부를 확인할 수 있습니다. 이 상수는 `std::net::error`에서 가져옵니다.

연결한 스트림과 accept로 받은 스트림은 각각 닫아야 합니다. 리스너를 닫는 것이 이미 수락한 모든 스트림을 닫는다는 뜻은 아닙니다. 구조체 복사본을 각각 닫지 마십시오.

## TCP는 메시지 경계를 보존하지 않음

한 번 쓴 내용이 한 번 읽기로 전부 돌아온다고 가정하지 않습니다. 길이 접두사, 구분자 또는 고정 길이처럼 프로토콜의 경계를 정해야 합니다. 고정 길이는 `tcp_read_exact`, 길이가 알려지지 않은 스트림은 읽기 반복과 EOF 처리로 다룹니다.

음수는 오류이고, 양수 길이 읽기에서 0은 상대편의 종료입니다. write_all 실패 전에 일부 데이터가 전송될 수 있습니다. 같은 내용을 무조건 처음부터 재전송하면 중복될 수 있습니다.

## 대기와 시간 제한

기본 blocking 함수는 상대편을 오래 기다릴 수 있습니다. 시간 제한이 필요한 프로그램은 `tcp_connect_addr_timeout`, `tcp_read_timeout`, `tcp_write_timeout` 계열을 선택합니다. 밀리초 인자와 오류 결과를 확인하고, 상대편이 응답하지 않는 경로도 설계하십시오. 비동기는 [task](/docs/ko/stdlib/task)와 대응하는 비동기 네트워크 API를 함께 사용합니다.

[주소 조회](/docs/ko/stdlib/resolution) · [완성된 로컬 클라이언트](/docs/ko/practice/tcp-client)

## 읽기 함수를 고르는 기준

| 필요한 동작 | 함수 | 확인할 값 |
| --- | --- | --- |
| 도착한 데이터 일부 읽기 | `tcp_read` | 음수 오류, 0 종료, 양수 바이트 수 |
| 정해진 길이만큼 읽기 | `tcp_read_exact` | 요청 길이를 모두 읽었는지 |
| 주어진 내용을 끝까지 보내기 | `tcp_write_all` | 요청 길이와 반환값 |
| 대기 시간 제한 | timeout 계열 | 반환 결과와 timeout 오류 |

예를 들어 4바이트 길이 필드 다음에 본문이 오는 프로토콜은, 먼저 길이 필드를 끝까지 읽고 허용하는 최대 길이인지 검사한 뒤 본문 공간을 확보합니다. 길이가 확인되기 전에 상대가 준 값만 믿고 큰 메모리를 할당하지 않습니다.

## 연결을 닫는 시점

읽기 함수가 0을 반환해도 로컬 스트림 자원은 남아 있습니다. 읽기를 마친 뒤 tcp_close를 호출합니다. 전송 오류 뒤에도 같은 정리 경로를 사용하면 정상·실패 경로의 닫기를 빠뜨리지 않습니다.

리스너는 새 연결을 받는 자원이고 스트림은 이미 연결된 통신 자원입니다. 서버를 만들 때는 두 종류를 각각 관리합니다. accept가 성공할 때마다 새 스트림이 생기므로 처리 후 닫고, 서버의 수락 반복이 끝나면 리스너를 닫습니다.

## 직접 실행해 보기

[로컬 TCP 클라이언트 실습](/docs/ko/practice/tcp-client)에 서버와 클라이언트의 전체 코드가 있습니다. 서버를 먼저 실행한 경우, 서버가 없는 경우, 포트 번호가 다른 경우를 비교하면 주소 조회와 연결 실패를 구분할 수 있습니다.
