---
translation_set_id: stdlib-resolution
path: stdlib/resolution
locale: ru
group: stdlib
group_order: 1
order: 8
title: net.resolve: Разрешение адресов
summary: Описывает усечение числовых адресов, списков имен вызывающих абонентов и результатов.
---

## Выберите метод запроса

По умолчанию для Linux, resolver принимает числовой адрес IPv4/IPv6 и числовой порт. Имя хоста или имя службы не передаются системе неявно DNS. Например, `127.0.0.1` и `8080` — это числовые входные данные, а `example.com` и `http` — имена.

Чтобы использовать явно предоставленный список имен, используйте `std::net::resolve_table`. Вы можете подключиться напрямую к числовому адресу или выполнить поиск, зарегистрировав адрес в списке имен.

## Базовый поиск

```text
std::net::resolve
net_resolve_tcp(host: str, service: str, output: ptr<SocketAddr>, capacity: i64) -> NetResolveResult
net_resolve_udp(host: str, service: str, output: ptr<SocketAddr>, capacity: i64) -> NetResolveResult
net_resolve_host(host: str, service: str, output: ptr<SocketAddr>, capacity: i64) -> NetResolveResult
```

Выходной массив предоставляется вызывающей стороной. Единицей capacity являются не байты, а `SocketAddr` количество элементов. Вы можете осуществлять поиск по счетчику только с помощью capacity=0 и null output. Внутреннее хранилище поиска очищается перед возвратом и не становится владельцем выходного массива.

|поле результата|смысл|
| --- | --- |
| `ok` |Был ли запрос успешным или нет|
| `count` |Общее количество поддерживаемых результатов адресов|
| `written` |Количество элементов, фактически записанных в выходной массив|
| `truncated` |Выходная мощность меньше общего результата?|
| `error.kind` |Нормализованные классификации, такие как INVALID, NOT_FOUND, UNSUPPORTED.|
| `error.native_code` |Собственный код для определения причины|

Даже в случае успеха written может быть 0. Проверьте `ok && written > 0` перед использованием output[0].

## Список имен, которые вы предоставляете сами

```text
net_resolve_from_table(host: str, service: str, protocol: i32,
    entries: ptr<NetResolveEntry>, entry_count: i64,
    output: ptr<SocketAddr>, capacity: i64) -> NetResolveResult
```

`NetResolveEntry` в `std::net::resolve_table` имеет host, service, protocol, address. Имя хоста ASCII не чувствительно к регистру, и имя службы сравнивается точно. TCP/UDP/Для всех запросов протокола используются константы модуля. Элементы и строки принадлежат вызывающему объекту и не сохраняются после вызова.

Совпадения копируются в порядке ввода, дубликаты сохраняются. Выходное пространство и входное пространство хранения не должны перекрываться. Если имени нет в списке, это NOT_FOUND и повторная попытка с внешним DNS не будет.

[TCP Как использовать](/docs/ru/stdlib/tcp) · [TCP Клиентская практика](/docs/ru/practice/tcp-client)

## Поиск числовых адресов

Подготовьте массив для хранения результатов запроса и проверьте, записан ли адрес. В приведенном ниже примере не требуется сервер, поскольку поиск адреса не создает само соединение.

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

Результат выполнения:

```text
address ready
```

Для фактического подключения передайте addresses[0] в серию функций tcp_connect_addr. Даже если вы получите адрес, на этапе подключения вы сможете определить, работает ли сервер на этом порту.

## Поиск по списку имен

Пользовательский список адресов полезен для связывания имен на основе файла конфигурации или списка служб программы. В следующем примере имя api сопоставляется с локальным портом 8080.

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

Результат выполнения:

```text
matched=1
```

API и api совпадают при сравнении имен хостов. http и HTTP не совпадают в этом примере, поскольку они различаются при сравнении имен служб. Если вы зарегистрируете несколько адресов с одним и тем же именем, результаты будут отображаться в указанном порядке. Если вы получили только часть небольшого массива, читайте только до written, а если вам нужен весь список, подготовьте место для count и повторите поиск.
