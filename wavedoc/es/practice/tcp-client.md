---
translation_set_id: practice-tcp-client
path: practice/tcp-client
locale: es
group: practice
group_order: 4
order: 4
title: Proyecto: un cliente TCP local
summary: Conecta búsqueda de direcciones numéricas, falla de conexión, transferencia y cierre.
---

## Preparación y alcance

Este programa se conecta a `127.0.0.1:8080` en OS nativo, envía `ping` y LF, y sale. Necesitará un servidor de prueba en una terminal separada. Si Python 3 está presente, el siguiente servidor recibirá y mostrará hasta 5 bytes de una conexión local.

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

Primero ejecute el servidor, luego guarde el cliente como `main.wave` y ejecútelo como `wavec run main.wave`. El servidor genera `b'ping\n'`. Si el puerto ya está en uso, cámbielo tanto en el servidor como en el cliente juntos.

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

El resultado exitoso del cliente es `sent=5`. Si el servidor no se está ejecutando, el manejo normal de fallas es imprimir `connect failed` y salir con 2.

## Comprender el código

El rango válido de la matriz de búsqueda es hasta written. Este es un pequeño ejemplo que solo usa la primera dirección. En servicios con múltiples candidatos, deberá determinar una política de conexión para cada candidato. La transmisión conectada se cierra independientemente del éxito o fracaso de la transmisión. El tiempo de espera de la conexión es de 1000 ms y este no es un ejemplo que garantice un tiempo de espera para toda la transmisión.

TCP no conserva los límites de transmisión. Por eso se escribió para que el servidor pueda ejecutarse recv varias veces. En [búsqueda de dirección](/docs/es/stdlib/resolution) se explica cómo encontrar una dirección por nombre.

## práctica extendida

Haga que el servidor envíe la respuesta y que el cliente la lea. Después de determinar la duración de la respuesta, debe manejar lecturas parciales, EOF, tiempos de espera y cierres juntos.

## Papel de los dos terminales

El terminal del servidor espera una conexión en el estado listen. Cuando un cliente se conecta, accept devuelve el socket conectado y el bucle recv recopila datos. El terminal de cliente realiza la búsqueda de direcciones, la conexión, la transferencia y el cierre en ese orden.

El código Python es el servidor asociado del laboratorio. Wave Se utiliza para verificar fácilmente los bytes pasados ​​por el programa y no se coloca en el mismo archivo que el código del cliente. El servidor finaliza después de manejar una conexión, por lo que también se reinicia antes de ejecutar el cliente nuevamente.

## ¿Por qué enviar 5 bytes?

Dado que `ping` tiene 4 bytes y LF es 1 byte, la longitud de transmisión es 5. El NUL al final de la cadena no se incluye en el mensaje. El servidor recopila 5 bytes y luego los genera, por lo que incluso si los paquetes llegan en varios fragmentos, se producirá el mismo resultado.

Si tcp_write_all devuelve 5, la función de transferencia local ha procesado los bytes solicitados. Esto no confirma que el otro programa haya interpretado y guardado el mensaje. Si el protocolo requiere confirmación de la finalización del procesamiento, el servidor envía una respuesta y el cliente lee la respuesta.

## Práctica del camino al fracaso

|cambiar|resultado|que aprender|
| --- | --- | --- |
|Ejecutar sin servidor|Error de conexión, código de salida 2|Incluso si la dirección es válida, es posible que el servidor no exista.|
|Configure los puertos del servidor y del cliente de manera diferente|La conexión falló|Tanto IP como el puerto de la dirección deben coincidir.|
|Cambie el tamaño del servidor recv a 1|Los mismos 5 bytes juntos|TCP La unidad de lectura es diferente de la unidad de mensaje|
|Dividir las transferencias de clientes en varias veces|recibido en el mismo orden|Los límites de los mensajes están determinados por el protocolo.|

El tiempo de espera de conexión de 1000 ms es un ajuste en la fase de conexión. Para limitar la lectura y escritura, utilice la función timeout en el paso correspondiente, y si hay un límite de tiempo para aplicar a todo el programa, calcule el tiempo restante y páselo.
