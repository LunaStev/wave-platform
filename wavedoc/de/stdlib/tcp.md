---
translation_set_id: stdlib-tcp
path: stdlib/tcp
locale: de
group: stdlib
group_order: 1
order: 9
title: net.tcp: Verbindungen und Übertragungen
summary: TCP Beschreibt die Ergebnisstruktur, die Teilübertragungs- und Sperrverantwortung.
---

## Überprüfen Sie die Verbindungsergebnisse

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

Die Verknüpfungsfunktion `NetResult<T>` hat `ok`, `value` und `error`. Verwenden Sie value nur bei Erfolg. `NetError` enthält die normalisierte Fehlerklassifizierung und native_code, der Erfolg kann mit `error.kind == NET_ERROR_NONE` überprüft werden. Diese Konstante stammt aus `std::net::error`.

Der verbundene Stream und der mit accept empfangene Stream müssen jeweils geschlossen werden. Das Schließen eines Listeners bedeutet nicht, dass alle Streams geschlossen werden, die er bereits akzeptiert hat. Schließen Sie nicht jede Kopie der Struktur.

## TCP behält die Nachrichtengrenzen nicht bei

Gehen Sie nicht davon aus, dass Ihnen alles, was Sie einmal geschrieben haben, wieder einfällt, wenn Sie es gelesen haben. Das Protokoll muss durch ein Längenpräfix, ein Trennzeichen oder eine feste Länge begrenzt werden. Feste Längen werden von `tcp_read_exact` verarbeitet, und Streams unbekannter Länge werden von Leseiterationen und EOF-Verarbeitung verarbeitet.

Negative Zahlen sind Fehler, und bei Lesevorgängen mit positiver Länge ist 0 das Ende des anderen Endes. write_all Einige Daten wurden möglicherweise vor dem Fehler übertragen. Wenn Sie denselben Inhalt von Anfang an erneut übertragen, kann es zu Duplikaten kommen.

## Wartezeiten und Fristen

Die Standardfunktion blocking kann lange auf die andere Seite warten. Für Programme, die zeitliche Einschränkungen erfordern, wählen Sie die Serien `tcp_connect_addr_timeout`, `tcp_read_timeout`, `tcp_write_timeout`. Überprüfen Sie den Millisekundenfaktor und die Fehlerfolgen und entwerfen Sie auch eine Route, bei der die Gegenseite nicht reagiert. Asynchron verwendet [task](/docs/de/stdlib/task) zusammen mit dem entsprechenden asynchronen Netzwerk API.

[Adresssuche](/docs/de/stdlib/resolution) · [Lokaler Client abgeschlossen](/docs/de/practice/tcp-client)

## Kriterien für die Auswahl einer Lesefunktion

|Handlungsbedarf|Funktion|Wert zu überprüfen|
| --- | --- | --- |
|Lesen Sie einige der angekommenen Daten| `tcp_read` |Negativer Fehler, Nullterminierung, positive Anzahl von Bytes|
|Lesen Sie eine bestimmte Länge| `tcp_read_exact` |Wurde die Anfragelänge vollständig gelesen?|
|Senden Sie den gegebenen Inhalt bis zum Ende| `tcp_write_all` |Anforderungslänge und Rückgabewert|
|Wartezeitlimit|timeout Serie|Ergebnis und Fehler zurückgeben timeout|

Bei einem Protokoll mit einem vier Byte langen Längenfeld, gefolgt von einem Textkörper, lesen Sie zunächst das vollständige Längenfeld, überprüfen Sie, ob die Länge das zulässige Maximum nicht überschreitet, und weisen Sie dann Platz für den Textkörper zu. Weisen Sie keine großen Speichermengen aus einer vom Peer bereitgestellten, nicht validierten Länge zu.

## Wann die Verbindung geschlossen werden soll

Auch wenn die Lesefunktion 0 zurückgibt, bleiben lokale Stream-Ressourcen erhalten. Rufen Sie nach Abschluss der Lektüre tcp_close an. Wenn Sie auch nach einem Übertragungsfehler denselben Bereinigungspfad verwenden, vergessen Sie nicht, den normalen und den fehlgeschlagenen Pfad zu schließen.

Ein Listener ist eine Ressource, die neue Verbindungen empfängt, und ein Stream ist eine Kommunikationsressource, die bereits verbunden ist. Beim Erstellen eines Servers verwalten Sie beide Typen. Jedes Mal, wenn accept erfolgreich ist, wird ein neuer Stream erstellt. Daher schließen wir ihn nach der Verarbeitung und schließen den Listener, wenn die Akzeptanziteration des Servers abgeschlossen ist.

## Probieren Sie es selbst aus

[Lokale TCP Kundenpraxis](/docs/de/practice/tcp-client) enthält den vollständigen Code für den Server und den Client. Sie können zwischen Adresssuche und Verbindungsfehler unterscheiden, indem Sie Fälle vergleichen, in denen der Server zuerst ausgeführt wird, wenn kein Server vorhanden ist und wenn die Portnummer unterschiedlich ist.
