---
translation_set_id: stdlib-resolution
path: stdlib/resolution
locale: ko
group: stdlib
group_order: 1
order: 8
title: net.resolve: 주소 조회 정책
summary: 숫자 주소, 호출자 소유 이름 목록과 결과의 잘림을 설명합니다.
---

## 조회 방식 선택

Linux의 기본 resolver는 숫자 IPv4/IPv6 주소와 숫자 포트를 받습니다. 호스트명이나 서비스명을 암묵적으로 시스템 DNS에 넘기지 않습니다. 예를 들어 `127.0.0.1`과 `8080`은 숫자 입력이고, `example.com`과 `http`는 이름입니다.

명시적으로 제공한 이름 목록을 사용하려면 `std::net::resolve_table`을 사용합니다. 숫자 주소로 직접 연결하거나, 이름 목록에 주소를 등록해 조회할 수 있습니다.

## 기본 조회

```text
std::net::resolve
net_resolve_tcp(host: str, service: str, output: ptr<SocketAddr>, capacity: i64) -> NetResolveResult
net_resolve_udp(host: str, service: str, output: ptr<SocketAddr>, capacity: i64) -> NetResolveResult
net_resolve_host(host: str, service: str, output: ptr<SocketAddr>, capacity: i64) -> NetResolveResult
```

출력 배열은 호출자가 제공합니다. capacity의 단위는 바이트가 아니라 `SocketAddr` 요소 수입니다. capacity=0과 null output으로 개수만 조회할 수 있습니다. 내부 조회 저장소는 반환 전에 정리되며 출력 배열의 소유권을 가져가지 않습니다.

| 결과 필드 | 의미 |
| --- | --- |
| `ok` | 조회 성공 여부 |
| `count` | 지원하는 주소 결과의 전체 개수 |
| `written` | 출력 배열에 실제 기록한 요소 수 |
| `truncated` | 출력 용량이 전체 결과보다 작은지 |
| `error.kind` | INVALID, NOT_FOUND, UNSUPPORTED 등 정규화된 분류 |
| `error.native_code` | 원인 확인용 네이티브 코드 |

성공이어도 written이 0일 수 있습니다. output[0]을 사용하기 전에 `ok && written > 0`을 확인합니다.

## 직접 제공하는 이름 목록

```text
net_resolve_from_table(host: str, service: str, protocol: i32,
    entries: ptr<NetResolveEntry>, entry_count: i64,
    output: ptr<SocketAddr>, capacity: i64) -> NetResolveResult
```

`std::net::resolve_table`의 `NetResolveEntry`는 host, service, protocol, address를 가집니다. 호스트명은 ASCII 대소문자를 구분하지 않고 서비스명은 정확히 비교합니다. TCP/UDP/모든 프로토콜 조회에는 모듈의 상수를 사용합니다. 항목과 문자열은 호출자가 소유하며 호출 이후 보관하지 않습니다.

일치 항목은 입력 순서대로 복사하고 중복을 유지합니다. 출력 공간과 입력 항목 저장 공간은 겹치지 않아야 합니다. 목록에 이름이 없으면 NOT_FOUND이며 외부 DNS로 재시도하지 않습니다.

[TCP 사용법](/docs/ko/stdlib/tcp) · [TCP 클라이언트 실습](/docs/ko/practice/tcp-client)

## 숫자 주소 조회하기

조회 결과를 담을 배열을 준비하고 기록된 주소가 있는지 확인합니다. 주소 조회는 연결 자체를 만들지 않으므로 아래 예제에는 서버가 필요하지 않습니다.

<!-- wave-example: book-resolve-numeric -->
```wave
import("std::net::address")::{
    SocketAddr
};
import("std::net::resolve")::{
    NetResolveResult,
    net_resolve_tcp
};

fun main() -> i32 {
    var addresses: array<SocketAddr, 4>;
    var result: NetResolveResult = net_resolve_tcp(
        "127.0.0.1",
        "8080",
        &addresses[0],
        4
    );

    if (!result.ok || result.written == 0) {
        println("address unavailable");
        return 1;
    }

    println("address ready");
    return 0;
}
```

실행 결과:

```text
address ready
```

실제로 연결하려면 addresses[0]을 tcp_connect_addr 계열 함수에 전달합니다. 주소를 얻었더라도 해당 포트에서 서버가 실행 중인지는 연결 단계에서 알 수 있습니다.

## 이름 목록으로 조회하기

직접 작성한 주소 목록은 설정 파일이나 프로그램의 서비스 목록을 기반으로 이름을 연결할 때 유용합니다. 다음 예제는 api라는 이름을 로컬 8080 포트에 대응시킵니다.

<!-- wave-example: book-resolve-table -->
```wave
import("std::net::address")::{
    SocketAddr,
    socket_addr_from_v4,
    socket_addr_v4_loopback
};
import("std::net::resolve")::{
    NetResolveResult
};
import("std::net::resolve_table")::{
    NetResolveEntry,
    RESOLVE_TCP,
    net_resolve_from_table
};

fun main() -> i32 {
    var entries: array<NetResolveEntry, 1>;
    entries[0] = NetResolveEntry {
        host: "api",
        service: "http",
        protocol: RESOLVE_TCP,
        address: socket_addr_from_v4(socket_addr_v4_loopback(8080))
    };

    var addresses: array<SocketAddr, 2>;
    var result: NetResolveResult = net_resolve_from_table(
        "API",
        "http",
        RESOLVE_TCP,
        &entries[0],
        1,
        &addresses[0],
        2
    );

    if (!result.ok || result.written != 1) {
        return 1;
    }

    println("matched={}", result.written);
    return 0;
}
```

실행 결과:

```text
matched=1
```

API와 api는 호스트 이름 비교에서 일치합니다. http와 HTTP는 서비스 이름 비교에서 다르므로 이 예제에서는 일치하지 않습니다. 같은 이름으로 주소를 여러 개 등록하면 입력 순서대로 결과가 나옵니다. 작은 배열에 일부만 받았다면 written까지만 읽고, 전체 목록이 필요하면 count에 맞는 공간을 준비해 다시 조회합니다.
