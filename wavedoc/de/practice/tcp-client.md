---
translation_set_id: practice-tcp-client
path: practice/tcp-client
locale: de
group: practice
group_order: 4
order: 4
title: Projekt: Ein lokaler TCP-Client
summary: Verbindet die Suche nach numerischen Adressen, Verbindungsfehler, Übertragung und Schließen.
---

## Vorbereitung und Umfang

Dieses Programm stellt eine Verbindung zu `127.0.0.1:8080` auf nativem OS her, sendet `ping` und LF und wird beendet. Sie benötigen einen Testserver in einem separaten Terminal. Wenn Python 3 vorhanden ist, empfängt und zeigt der nächste Server bis zu 5 Bytes von einer lokalen Verbindung an.

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

Führen Sie zuerst den Server aus, speichern Sie dann den Client als `main.wave` und führen Sie ihn als `wavec run main.wave` aus. Der Server gibt `b'ping\n'` aus. Wenn der Port bereits verwendet wird, ändern Sie den Port sowohl auf dem Server als auch auf dem Client gemeinsam.

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

Die Erfolgsausgabe des Clients lautet `sent=5`. Wenn der Server nicht läuft, besteht die normale Fehlerbehandlung darin, `connect failed` auszudrucken und mit 2 zu beenden.

## Den Kodex verstehen

Der gültige Bereich des Lookup-Arrays reicht bis written. Dies ist ein kleines Beispiel, das nur die erste Adresse verwendet. Bei Diensten mit mehreren Kandidaten müssen Sie für jeden Kandidaten eine Verbindungsrichtlinie festlegen. Der verbundene Stream wird unabhängig vom Erfolg oder Misserfolg der Übertragung geschlossen. Das Verbindungszeitlimit beträgt 1000 ms. Dies ist kein Beispiel, das ein Zeitlimit für die gesamte Übertragung garantiert.

TCP behält keine Übertragungsgrenzen bei. Aus diesem Grund wurde es so geschrieben, dass der Server mehrmals ausgeführt werden kann recv. Wie Sie eine Adresse anhand des Namens finden, erfahren Sie in [Adresssuche](/docs/de/stdlib/resolution).

## erweiterte Praxis

Lassen Sie den Server die Antwort senden und den Client lesen. Nachdem Sie die Länge der Antwort bestimmt haben, müssen Sie Teillesevorgänge, EOF, Zeitüberschreitungen und Schließungen gemeinsam verarbeiten.

## Rolle der beiden Terminals

Das Server-Terminal wartet auf eine Verbindung im Status listen. Wenn ein Client eine Verbindung herstellt, gibt accept den verbundenen Socket zurück und die Schleife recv sammelt Daten. Das Client-Terminal führt die Adresssuche, die Verbindung, die Übertragung und das Schließen in dieser Reihenfolge durch.

Der Code Python ist der Laborpartnerserver. Wave Es dient zur einfachen Überprüfung der vom Programm übergebenen Bytes und wird nicht in derselben Datei wie der Client-Code abgelegt. Der Server wird nach der Verarbeitung einer Verbindung beendet und startet daher auch neu, bevor der Client erneut ausgeführt wird.

## Warum 5 Bytes senden?

Da `ping` 4 Bytes und LF 1 Byte groß ist, beträgt die Übertragungslänge 5. Das NUL am Ende der Zeichenfolge ist nicht in der Nachricht enthalten. Der Server sammelt 5 Bytes und gibt sie dann aus. Selbst wenn die Pakete in mehreren Blöcken ankommen, wird also das gleiche Ergebnis erzielt.

Wenn tcp_write_all 5 zurückgibt, hat die lokale Übertragungsfunktion die angeforderten Bytes verarbeitet. Dies bestätigt nicht, dass das andere Programm die Nachricht interpretiert und gespeichert hat. Wenn das Protokoll eine Bestätigung des Abschlusses der Verarbeitung erfordert, sendet der Server eine Antwort und der Client liest die Antwort.

## Weg zur Fehlerpraxis

|ändern|Ergebnis|Was zu lernen ist|
| --- | --- | --- |
|Ohne Server laufen|Verbindung fehlgeschlagen, Exit-Code 2|Auch wenn die Adresse gültig ist, existiert der Server möglicherweise nicht|
|Stellen Sie Server- und Client-Ports unterschiedlich ein|Verbindung fehlgeschlagen|Sowohl IP als auch der Port in der Adresse müssen übereinstimmen.|
|Ändern Sie die Größe des Servers recv auf 1|Gleiche 5 Bytes zusammen|TCP Die Leseeinheit unterscheidet sich von der Nachrichteneinheit|
|Aufteilen von Kundentransfers in mehrere Male|in der gleichen Reihenfolge erhalten|Nachrichtengrenzen werden durch das Protokoll bestimmt.|

Der Verbindungs-Timeout von 1000 ms ist eine Einstellung in der Verbindungsphase. Um das Lesen und Schreiben einzuschränken, verwenden Sie im entsprechenden Schritt die Funktion timeout. Wenn für das gesamte Programm eine Zeitbegrenzung gilt, berechnen Sie die verbleibende Zeit und geben Sie diese weiter.
