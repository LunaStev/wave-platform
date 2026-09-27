---
translation_set_id: stdlib-resolution
path: stdlib/resolution
locale: en
group: stdlib
group_order: 1
order: 8
title: net.resolve: Address resolution
summary: Describes truncation of numeric addresses, caller-owned name lists, and results.
---

## Select inquiry method

The default for Linux, resolver, takes a numeric IPv4/IPv6 address and a numeric port. Hostname or service name is not passed implicitly to the system DNS. For example, `127.0.0.1` and `8080` are numeric inputs, and `example.com` and `http` are names.

To use an explicitly provided list of names, use `std::net::resolve_table`. You can connect directly to the numeric address or search by registering the address in the name list.

## Basic lookup

```text
std::net::resolve
net_resolve_tcp(host: str, service: str, output: ptr<SocketAddr>, capacity: i64) -> NetResolveResult
net_resolve_udp(host: str, service: str, output: ptr<SocketAddr>, capacity: i64) -> NetResolveResult
net_resolve_host(host: str, service: str, output: ptr<SocketAddr>, capacity: i64) -> NetResolveResult
```

The output array is provided by the caller. The unit of capacity is not bytes, but `SocketAddr` number of elements. You can only search the count with capacity=0 and null output. The internal lookup store is cleaned up before return and does not take ownership of the output array.

|result field|meaning|
| --- | --- |
| `ok` |Whether the query was successful or not|
| `count` |Total number of address results supported|
| `written` |Number of elements actually written to the output array|
| `truncated` |Is the output capacity smaller than the total result?|
| `error.kind` |Normalized classifications such as INVALID, NOT_FOUND, UNSUPPORTED|
| `error.native_code` |Native code to determine cause|

Even with success, written may be 0. Check `ok && written > 0` before using output[0].

## List of names you provide yourself

```text
net_resolve_from_table(host: str, service: str, protocol: i32,
    entries: ptr<NetResolveEntry>, entry_count: i64,
    output: ptr<SocketAddr>, capacity: i64) -> NetResolveResult
```

`NetResolveEntry` in `std::net::resolve_table` has host, service, protocol, address. The host name ASCII is case insensitive and the service name is compared exactly. TCP/UDP/For all protocol queries, constants of the module are used. Items and strings are owned by the caller and are not retained after the call.

Matches are copied in input order and duplicates are maintained. Output space and input storage space must not overlap. If the name is not in the list, it is NOT_FOUND and will not retry with the external DNS.

[TCP How to use](/docs/en/stdlib/tcp) · [TCP Client Practice](/docs/en/practice/tcp-client)

## Look up numeric addresses

Prepare an array to hold the query results and check whether the address is recorded. The example below does not require a server because address lookup does not create the connection itself.

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

Execution result:

```text
address ready
```

To actually connect, pass addresses[0] to the tcp_connect_addr series of functions. Even if you get the address, you can tell at the connection stage whether a server is running on that port.

## Search by name list

A custom address list is useful for linking names based on a configuration file or a program's service list. The following example maps the name api to the local 8080 port.

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

Execution result:

```text
matched=1
```

API and api match in hostname comparison. http and HTTP do not match in this example because they are different in the service name comparison. If you register multiple addresses with the same name, results will appear in the order entered. If you only received part of a small array, read only up to written, and if you need the entire list, prepare space for count and search again.
