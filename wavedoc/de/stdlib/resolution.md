---
translation_set_id: stdlib-resolution
path: stdlib/resolution
locale: de
group: stdlib
group_order: 1
order: 8
title: net.resolve: Adressauflösung
summary: Beschreibt die Kürzung numerischer Adressen, Listen mit anrufereigenen Namen und Ergebnissen.
---

## Anfragemethode auswählen

Der Standardwert für Linux, resolver, nimmt eine numerische IPv4/IPv6-Adresse und einen numerischen Port an. Hostname oder Dienstname werden nicht implizit an das System übergeben DNS. Beispielsweise sind `127.0.0.1` und `8080` numerische Eingaben und `example.com` und `http` sind Namen.

Um eine explizit bereitgestellte Namensliste zu verwenden, verwenden Sie `std::net::resolve_table`. Sie können eine direkte Verbindung zur numerischen Adresse herstellen oder eine Suche durchführen, indem Sie die Adresse in der Namensliste registrieren.

## Einfache Suche

```text
std::net::resolve
net_resolve_tcp(host: str, service: str, output: ptr<SocketAddr>, capacity: i64) -> NetResolveResult
net_resolve_udp(host: str, service: str, output: ptr<SocketAddr>, capacity: i64) -> NetResolveResult
net_resolve_host(host: str, service: str, output: ptr<SocketAddr>, capacity: i64) -> NetResolveResult
```

Das Ausgabearray wird vom Aufrufer bereitgestellt. Die Einheit von capacity sind nicht Bytes, sondern `SocketAddr` Anzahl der Elemente. Sie können die Zählung nur mit capacity=0 und null output durchsuchen. Der interne Suchspeicher wird vor der Rückgabe bereinigt und übernimmt nicht den Besitz des Ausgabearrays.

|Ergebnisfeld|Bedeutung|
| --- | --- |
| `ok` |Ob die Abfrage erfolgreich war oder nicht|
| `count` |Gesamtzahl der unterstützten Adressergebnisse|
| `written` |Anzahl der tatsächlich in das Ausgabearray geschriebenen Elemente|
| `truncated` |Ist die Ausbringungskapazität kleiner als das Gesamtergebnis?|
| `error.kind` |Normalisierte Klassifizierungen wie INVALID, NOT_FOUND, UNSUPPORTED|
| `error.native_code` |Nativer Code zur Ermittlung der Ursache|

Auch bei Erfolg kann written 0 sein. Überprüfen Sie `ok && written > 0`, bevor Sie output[0] verwenden.

## Liste der Namen, die Sie selbst angeben

```text
net_resolve_from_table(host: str, service: str, protocol: i32,
    entries: ptr<NetResolveEntry>, entry_count: i64,
    output: ptr<SocketAddr>, capacity: i64) -> NetResolveResult
```

`NetResolveEntry` in `std::net::resolve_table` hat host, service, protocol, address. Beim Hostnamen ASCII wird die Groß-/Kleinschreibung nicht beachtet und der Dienstname wird genau verglichen. TCP/UDP/Für alle Protokollabfragen werden Konstanten des Moduls verwendet. Elemente und Zeichenfolgen sind Eigentum des Aufrufers und werden nach dem Aufruf nicht beibehalten.

Übereinstimmungen werden in der Eingabereihenfolge kopiert und Duplikate bleiben erhalten. Ausgaberaum und Eingabespeicherplatz dürfen sich nicht überschneiden. Wenn der Name nicht in der Liste enthalten ist, lautet er NOT_FOUND und es wird kein erneuter Versuch mit dem externen DNS durchgeführt.

[TCP Verwendung](/docs/de/stdlib/tcp) · [TCP Kundenpraxis](/docs/de/practice/tcp-client)

## Suchen Sie nach numerischen Adressen

Bereiten Sie ein Array für die Abfrageergebnisse vor und prüfen Sie, ob die Adresse aufgezeichnet wird. Für das folgende Beispiel ist kein Server erforderlich, da die Adresssuche nicht die Verbindung selbst herstellt.

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

Ausführungsergebnis:

```text
address ready
```

Um tatsächlich eine Verbindung herzustellen, übergeben Sie addresses[0] an die Funktionsreihe tcp_connect_addr. Selbst wenn Sie die Adresse erhalten, können Sie bereits beim Verbindungsaufbau erkennen, ob an diesem Port ein Server läuft.

## Suche nach Namensliste

Eine benutzerdefinierte Adressliste ist nützlich, um Namen basierend auf einer Konfigurationsdatei oder der Dienstliste eines Programms zu verknüpfen. Das folgende Beispiel ordnet den Namen api dem lokalen 8080-Port zu.

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

Ausführungsergebnis:

```text
matched=1
```

API und api stimmen im Hostnamen-Vergleich überein. http und HTTP stimmen in diesem Beispiel nicht überein, da sie sich im Dienstnamenvergleich unterscheiden. Wenn Sie mehrere Adressen mit demselben Namen registrieren, werden die Ergebnisse in der eingegebenen Reihenfolge angezeigt. Wenn Sie nur einen Teil eines kleinen Arrays erhalten haben, lesen Sie nur bis written. Wenn Sie die gesamte Liste benötigen, bereiten Sie Platz für count vor und suchen Sie erneut.
