---
translation_set_id: practice-tcp-client
path: practice/tcp-client
locale: ko
group: practice
group_order: 4
order: 4
title: 실습: 로컬 TCP 클라이언트
summary: 숫자 주소 조회, 연결 실패, 전송과 닫기를 연결합니다.
---

## 준비와 범위

이 프로그램은 네이티브 OS의 `127.0.0.1:8080`에 연결해 `ping`과 LF를 전송하고 종료합니다. 별도 터미널에 테스트 서버가 필요합니다. Python 3이 있다면 다음 서버는 로컬 연결 한 개에서 최대 5바이트를 받아 표시합니다.

```python
import socket
with socket.socket() as server:
    server.bind(("127.0.0.1", 8080))
    server.listen(1)
    connection, address = server.accept()
    with connection:
        message = b""
        while len(message) < 5:
            chunk = connection.recv(5 - len(message))
            if not chunk:
                break
            message += chunk
        print(repr(message))
```

서버를 먼저 실행한 뒤 클라이언트를 `main.wave`로 저장해 `wavec run main.wave`로 실행합니다. 서버는 `b'ping\n'`을 출력합니다. 포트가 이미 사용 중이면 서버와 클라이언트 양쪽의 포트를 함께 바꾸십시오.

<!-- wave-example: tcp-client -->
```wave
import("std::net::address")::{
    SocketAddr
};
import("std::net::resolve")::{
    NetResolveResult, net_resolve_tcp
};
import("std::net::error")::{
    NetResult, NetError, NET_ERROR_NONE
};
import("std::net::tcp")::{
    TcpStream, tcp_connect_addr_timeout, tcp_write_all, tcp_close
};

fun main() -> i32 {
    var addresses: array<SocketAddr, 4>;
    var resolved: NetResolveResult = net_resolve_tcp("127.0.0.1", "8080", &addresses[0], 4);
    if (!resolved.ok || resolved.written == 0) {
        println("resolve failed");
        return 1;
    }

    var connected: NetResult<TcpStream> = tcp_connect_addr_timeout(addresses[0], 1000);
    if (!connected.ok) {
        println("connect failed");
        return 2;
    }

    var sent: i64 = tcp_write_all(connected.value, "ping\n" as ptr<u8>, 5);
    var closed: NetError = tcp_close(connected.value);
    if (sent != 5 || closed.kind != NET_ERROR_NONE) {
        return 3;
    }

    println("sent=5");
    return 0;
}
```

클라이언트의 성공 출력은 `sent=5`입니다. 서버를 실행하지 않으면 `connect failed`를 출력하고 2로 종료하는 것이 정상적인 실패 처리입니다.

## 코드 이해하기

조회 배열의 유효 범위는 written까지입니다. 첫 주소만 사용하는 작은 예제이며 여러 후보가 있는 서비스에서는 후보별 연결 정책을 정해야 합니다. 연결한 스트림은 전송 성공·실패와 무관하게 닫습니다. 연결 제한 시간은 1000ms이고, 전송 전체의 제한 시간을 보장하는 예제는 아닙니다.

TCP는 전송 경계를 보존하지 않습니다. 서버도 여러 번 recv할 수 있도록 작성한 이유입니다. 이름으로 주소를 찾는 방법은 [주소 조회](/docs/ko/stdlib/resolution)에서 다룹니다.

## 확장 연습

서버가 응답을 보내도록 바꾸고 클라이언트에서 읽어 보십시오. 응답 길이를 정한 뒤 부분 읽기, EOF, 제한 시간과 닫기를 함께 처리해야 합니다.

## 두 터미널의 역할

서버 터미널은 listen 상태에서 연결을 기다립니다. 클라이언트가 접속하면 accept가 연결된 소켓을 반환하고 recv 반복문이 데이터를 모읍니다. 클라이언트 터미널은 주소 조회, 연결, 전송, 닫기를 순서대로 수행합니다.

Python 코드는 실습 상대 서버입니다. Wave 프로그램이 전달한 바이트를 쉽게 확인하려고 사용하며, 클라이언트 코드와 같은 파일에 넣지 않습니다. 서버는 연결 하나를 처리한 뒤 종료하므로 클라이언트를 다시 실행하기 전에 서버도 다시 시작합니다.

## 5바이트를 보내는 이유

`ping`은 4바이트이고 LF가 1바이트이므로 전송 길이는 5입니다. 문자열 끝의 NUL은 메시지에 포함하지 않습니다. 서버는 5바이트를 모은 뒤 출력하므로 패킷이 여러 번에 나누어 도착해도 같은 결과가 나옵니다.

tcp_write_all이 5를 반환하면 로컬 전송 함수가 요청한 바이트를 처리한 것입니다. 상대 프로그램이 메시지를 해석하고 저장까지 마쳤다는 확인은 아닙니다. 처리 완료 확인이 필요한 프로토콜이라면 서버가 응답을 보내고 클라이언트가 그 응답을 읽도록 정합니다.

## 실패 경로 연습

| 변경 | 결과 | 배울 점 |
| --- | --- | --- |
| 서버 없이 실행 | 연결 실패, 종료 코드 2 | 주소가 유효해도 서버가 없을 수 있음 |
| 서버와 클라이언트 포트 다르게 설정 | 연결 실패 | 주소의 IP와 포트가 모두 맞아야 함 |
| 서버의 recv 크기를 1로 변경 | 같은 5바이트가 모임 | TCP 읽기 단위는 메시지 단위와 다름 |
| 클라이언트 전송을 여러 번으로 나누기 | 같은 순서로 수신 | 메시지 경계는 프로토콜이 정함 |

연결 제한 시간 1000ms는 연결 단계의 설정입니다. 읽기·쓰기도 제한하려면 해당 단계의 timeout 함수를 사용하고, 프로그램 전체에 적용할 제한 시간이 있다면 남은 시간을 계산해 전달합니다.
