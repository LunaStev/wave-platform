---
translation_set_id: stdlib-resolution
path: stdlib/resolution
locale: es
group: stdlib
group_order: 1
order: 8
title: net.resolve: Resolución de direcciones
summary: Describe el truncamiento de direcciones numéricas, listas de nombres propiedad de la persona que llama y resultados.
---

## Seleccione el método de consulta

El valor predeterminado para Linux, resolver, toma una dirección numérica IPv4/IPv6 y un puerto numérico. El nombre de host o de servicio no se pasa implícitamente al sistema DNS. Por ejemplo, `127.0.0.1` y `8080` son entradas numéricas y `example.com` y `http` son nombres.

Para utilizar una lista de nombres proporcionada explícitamente, utilice `std::net::resolve_table`. Puede conectarse directamente a la dirección numérica o buscar registrando la dirección en la lista de nombres.

## Búsqueda básica

```text
std::net::resolve
net_resolve_tcp(host: str, service: str, output: ptr<SocketAddr>, capacity: i64) -> NetResolveResult
net_resolve_udp(host: str, service: str, output: ptr<SocketAddr>, capacity: i64) -> NetResolveResult
net_resolve_host(host: str, service: str, output: ptr<SocketAddr>, capacity: i64) -> NetResolveResult
```

La matriz de salida la proporciona la persona que llama. La unidad de capacity no es bytes, sino `SocketAddr` número de elementos. Solo puedes buscar el recuento con capacity=0 y null output. El almacén de búsqueda interno se limpia antes de la devolución y no toma posesión de la matriz de salida.

|campo de resultado|significado|
| --- | --- |
| `ok` |Si la consulta fue exitosa o no|
| `count` |Número total de resultados de direcciones admitidos|
| `written` |Número de elementos realmente escritos en la matriz de salida|
| `truncated` |¿La capacidad de producción es menor que el resultado total?|
| `error.kind` |Clasificaciones normalizadas como INVALID, NOT_FOUND, UNSUPPORTED|
| `error.native_code` |Código nativo para determinar la causa.|

Incluso si tiene éxito, written puede ser 0. Verifique `ok && written > 0` antes de usar output[0].

## Lista de nombres que usted mismo proporciona

```text
net_resolve_from_table(host: str, service: str, protocol: i32,
    entries: ptr<NetResolveEntry>, entry_count: i64,
    output: ptr<SocketAddr>, capacity: i64) -> NetResolveResult
```

`NetResolveEntry` en `std::net::resolve_table` tiene host, service, protocol, address. El nombre de host ASCII no distingue entre mayúsculas y minúsculas y el nombre del servicio se compara exactamente. TCP/UDP/Para todas las consultas de protocolo, se utilizan constantes del módulo. Los elementos y las cadenas son propiedad de la persona que llama y no se conservan después de la llamada.

Las coincidencias se copian en el orden de entrada y se mantienen los duplicados. El espacio de salida y el espacio de almacenamiento de entrada no deben superponerse. Si el nombre no está en la lista, es NOT_FOUND y no volverá a intentarlo con el DNS externo.

[TCP Cómo utilizar](/docs/es/stdlib/tcp) · [TCP Práctica del cliente](/docs/es/practice/tcp-client)

## Buscar direcciones numéricas

Prepare una matriz para contener los resultados de la consulta y verifique si la dirección está registrada. El siguiente ejemplo no requiere un servidor porque la búsqueda de direcciones no crea la conexión en sí.

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

Resultado de la ejecución:

```text
address ready
```

Para conectarse realmente, pase addresses[0] a la serie de funciones tcp_connect_addr. Incluso si obtiene la dirección, puede saber en la etapa de conexión si se está ejecutando un servidor en ese puerto.

## Buscar por lista de nombres

Una lista de direcciones personalizada es útil para vincular nombres basados en un archivo de configuración o la lista de servicios de un programa. El siguiente ejemplo asigna el nombre api al puerto 8080 local.

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

Resultado de la ejecución:

```text
matched=1
```

API y api coinciden en la comparación de nombres de host. http y HTTP no coinciden en este ejemplo porque son diferentes en la comparación de nombres de servicios. Si registra varias direcciones con el mismo nombre, los resultados aparecerán en el orden ingresado. Si solo recibió parte de una pequeña matriz, lea solo hasta written, y si necesita la lista completa, prepare espacio para count y busque nuevamente.
