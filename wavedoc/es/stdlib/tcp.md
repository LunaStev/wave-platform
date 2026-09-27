---
translation_set_id: stdlib-tcp
path: stdlib/tcp
locale: es
group: stdlib
group_order: 1
order: 9
title: net.tcp: Conexiones y transferencias
summary: TCP Describe la estructura de resultados, transferencia parcial y responsabilidades de desconexión.
---

## Verificar resultados de conexión

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

La función de enlace `NetResult<T>` tiene `ok`, `value` y `error`. Utilice value solo si tiene éxito. `NetError` contiene la clasificación de error normalizado y native_code, y el éxito se puede verificar con `error.kind == NET_ERROR_NONE`. Esta constante se toma de `std::net::error`.

La transmisión conectada y la transmisión recibida con accept deben cerrarse respectivamente. Cerrar un oyente no significa cerrar todas las transmisiones que ya aceptó. No cierre cada copia de la estructura.

## TCP no conserva los límites de los mensajes

No asumas que todo lo que escribas una vez volverá a ti una vez que lo leas. El protocolo debe estar delimitado por un prefijo de longitud, delimitador o longitud fija. Las longitudes fijas se manejan mediante `tcp_read_exact`, y los flujos de longitud desconocida se manejan mediante iteraciones de lectura y procesamiento EOF.

Los números negativos son errores y, en lecturas de longitud positiva, 0 es el final del otro extremo. write_all Es posible que algunos datos se hayan transmitido antes del fallo. Si retransmites el mismo contenido desde el principio, es posible que se duplique.

## Espera y límites de tiempo

La función predeterminada blocking puede esperar mucho tiempo por el otro lado. Para programas que requieren restricciones de tiempo, seleccione las series `tcp_connect_addr_timeout`, `tcp_read_timeout`, `tcp_write_timeout`. Verifique el factor de milisegundos y las consecuencias del error, y también diseñe una ruta donde el otro lado no responda. Asíncrono utiliza [task](/docs/es/stdlib/task) junto con la correspondiente red asíncrona API.

[búsqueda de dirección](/docs/es/stdlib/resolution) · [Cliente local completado](/docs/es/practice/tcp-client)

## Criterios para seleccionar una función de lectura

|acción requerida|función|valor a comprobar|
| --- | --- | --- |
|Lea algunos de los datos llegados.| `tcp_read` |Error negativo, terminación cero, número positivo de bytes|
|Leer una cierta longitud| `tcp_read_exact` |¿Se ha leído en su totalidad la extensión de la solicitud?|
|Enviar el contenido dado hasta el final.| `tcp_write_all` |Longitud de la solicitud y valor de retorno|
|límite de tiempo de espera|Serie timeout|Devuelve resultado y error timeout|

Para un protocolo con un campo de longitud de cuatro bytes seguido de un cuerpo, primero lea el campo de longitud completo, verifique que la longitud no exceda el máximo permitido y luego asigne espacio para el cuerpo. No asigne grandes cantidades de memoria a partir de una longitud no validada proporcionada por el par.

## Cuando cerrar la conexión

Incluso si la función de lectura devuelve 0, los recursos del flujo local permanecen. Después de terminar de leer, llame al tcp_close. Si utiliza la misma ruta de limpieza incluso después de un error de transmisión, no olvidará cerrar las rutas normales y fallidas.

Un oyente es un recurso que recibe nuevas conexiones y una transmisión es un recurso de comunicación que ya está conectado. Al crear un servidor, administra ambos tipos. Cada vez que accept tiene éxito, se crea una nueva secuencia, por lo que la cerramos después de procesarla y cerramos el oyente cuando finaliza la iteración de aceptación del servidor.

## Pruébalo tú mismo

[Local TCP Práctica del cliente](/docs/es/practice/tcp-client) contiene el código completo para el servidor y el cliente. Puede distinguir entre búsqueda de direcciones y falla de conexión comparando los casos en los que el servidor se ejecuta primero, cuando no hay ningún servidor y cuando el número de puerto es diferente.
