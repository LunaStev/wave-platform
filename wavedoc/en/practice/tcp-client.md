---
translation_set_id: practice-tcp-client
path: practice/tcp-client
locale: en
group: practice
group_order: 4
order: 4
title: Project: A local TCP client
summary: Connects numeric address lookup, connection failure, transfer and close.
---

## Preparation and Scope

This program connects to `127.0.0.1:8080` on native OS, sends `ping` and LF, and exits. You will need a test server in a separate terminal. If Python 3 is present, the next server will receive and display up to 5 bytes from one local connection.

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

Run the server first, then save the client as `main.wave` and run it as `wavec run main.wave`. The server outputs `b'ping\n'`. If the port is already in use, change the port on both server and client together.

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

The success output from the client is `sent=5`. If the server is not running, normal failure handling is to print `connect failed` and exit with 2.

## Understanding the Code

The valid range of the lookup array is up to written. This is a small example that only uses the first address. In services with multiple candidates, you will need to determine a connection policy for each candidate. The connected stream is closed regardless of transmission success or failure. The connection timeout is 1000ms, and this is not an example that guarantees a timeout for the entire transmission.

TCP does not preserve transmission boundaries. This is why it was written so that the server can be run recv multiple times. How to find an address by name is covered in [address lookup](/docs/en/stdlib/resolution).

## extended practice

Have the server send the response and have the client read it. After determining the length of the response, you need to handle partial reads, EOF, timeouts and closes together.

## Role of the two terminals

The server terminal waits for a connection in status listen. When a client connects, accept returns the connected socket and recv loop collects data. The client terminal performs address lookup, connection, transfer, and close in that order.

Python code is the lab partner server. Wave It is used to easily check the bytes passed by the program, and is not placed in the same file as the client code. The server terminates after handling one connection, so it also restarts before running the client again.

## Why send 5 bytes

Since `ping` is 4 bytes and LF is 1 byte, the transmission length is 5. The NUL at the end of the string is not included in the message. The server collects 5 bytes and then outputs them, so even if the packets arrive in multiple chunks, the same result will be produced.

If tcp_write_all returns 5, the local transfer function has processed the requested bytes. This does not confirm that the other program has interpreted and saved the message. If the protocol requires confirmation of completion of processing, the server sends a response and the client reads the response.

## Path to Failure Practice

|change|result|What to learn|
| --- | --- | --- |
|Run without a server|Connection failed, exit code 2|Even if the address is valid, the server may not exist|
|Set server and client ports differently|Connection failed|Both IP and the port in the address must match.|
|Change the size of server recv to 1|Same 5 bytes together|TCP Read unit is different from message unit|
|Splitting client transfers into multiple times|received in the same order|Message boundaries are determined by the protocol.|

The connection timeout of 1000ms is a setting in the connection phase. To limit reading and writing, use the timeout function in the relevant step, and if there is a time limit to apply to the entire program, calculate the remaining time and pass it on.
