---
translation_set_id: stdlib-tcp
path: stdlib/tcp
locale: en
group: stdlib
group_order: 1
order: 9
title: net.tcp: Connections and transfers
summary: TCP Describes the result structure, partial transfer and disconnection responsibilities.
---

## Check connection results

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

The link function `NetResult<T>` has `ok`, `value`, and `error`. Use value only if successful. `NetError` contains the normalized error classification and native_code, and success can be checked with `error.kind == NET_ERROR_NONE`. This constant is taken from `std::net::error`.

The connected stream and the stream received with accept must be closed respectively. Closing a listener does not mean closing all streams it has already accepted. Do not close each copy of the structure.

## TCP does not preserve message boundaries

Don't assume that everything you write once will come back to you once you read it. The protocol must be delimited by a length prefix, delimiter, or fixed length. Fixed lengths are handled by `tcp_read_exact`, and streams of unknown length are handled by read iterations and EOF processing.

Negative numbers are errors, and in positive length reads, 0 is the end of the other end. write_all Some data may have been transmitted before failure. If you retransmit the same content from the beginning, it may be duplicated.

## Waiting and time limits

The default blocking function can wait a long time for the other side. For programs that require time restrictions, select the `tcp_connect_addr_timeout`, `tcp_read_timeout`, `tcp_write_timeout` series. Check the millisecond factor and error consequences, and also design a route where the other side does not respond. Asynchronous uses [task](/docs/en/stdlib/task) together with the corresponding asynchronous network API.

[address lookup](/docs/en/stdlib/resolution) · [Completed local client](/docs/en/practice/tcp-client)

## Criteria for selecting a reading function

|action required|function|value to check|
| --- | --- | --- |
|Read some of the arrived data| `tcp_read` |Negative error, zero termination, positive number of bytes|
|Read a certain length| `tcp_read_exact` |Has the request length been read in full?|
|Send given content to the end| `tcp_write_all` |Request Length and Return Value|
|Wait time limit|timeout Series|Return result and error timeout|

For a protocol with a four-byte length field followed by a body, first read the complete length field, check that the length does not exceed the allowed maximum, and then allocate space for the body. Do not allocate large amounts of memory from an unvalidated length supplied by the peer.

## When to close the connection

Even if the read function returns 0, local stream resources remain. After finishing reading, call tcp_close. If you use the same cleanup path even after a transmission error, you will not forget to close the normal and failed paths.

A listener is a resource that receives new connections, and a stream is a communication resource that is already connected. When creating a server, you manage both types. Each time accept succeeds, a new stream is created, so we close it after processing it, and close the listener when the server's acceptance iteration is finished.

## Try it yourself

[Local TCP Client Practice](/docs/en/practice/tcp-client) contains the complete code for the server and client. You can distinguish between address lookup and connection failure by comparing cases where the server is run first, when there is no server, and when the port number is different.
