---
translation_set_id: stdlib-files-io
path: stdlib/files-io
locale: de
group: stdlib
group_order: 1
order: 7
title: fs und io: Dateien und Byteübertragungen
summary: Beschreibt die Dateilebensdauer, den vollständigen Lesevorgang, die teilweise Übertragung und den Status nach einem Fehler.
---

## Datei API auswählen

Die Komfortfunktion von `std::fs::file` empfängt den Weg und führt das notwendige Öffnen und Schließen durch. Funktionen, die einen Deskriptor zurückgeben, müssen vom Aufrufer geschlossen werden.

|Erklärung|Erfolgsergebnisse und Vorsichtsmaßnahmen|
| --- | --- |
| `open_read(path: str) -> i64` |Deskriptor öffnen. Negative Zahlen sind Fehler|
| `create(path: str) -> i64` |Erstellen oder löschen Sie vorhandene Dateiinhalte. Gibt den besitzenden Deskriptor zurück|
| `open_append(path: str) -> i64` |Zum Hinzufügen öffnen oder erstellen|
| `size(path: str) -> i64` |Anzahl der Bytes. Negative Zahlen sind Fehler|
| `read_into(path: str, dst: ptr<u8>, dst_cap: i64) -> i64` |Anzahl der Bytes in der gesamten Datei. Mangelnde Kapazität ist ein Fehler|
| `read_to_end(path: str, dst_buffer: ptr<Buffer>) -> i64` |Fügen Sie die Datei nach dem vorhandenen Buffer hinzu und geben Sie den zusätzlichen Betrag zurück|
| `write(path: str, src: ptr<u8>, len: i64) -> i64` |Anzahl der geschriebenen Bytes, die vorhandenen Inhalt ersetzen|
| `append(path: str, src: ptr<u8>, len: i64) -> i64` |Anzahl der am Ende hinzugefügten Bytes|
| `remove(path: str) -> i64` |Entfernungsstatus. Scheitern ist negativ|

false von `exists(path)` allein kann nicht zwischen fehlenden Dateien und Berechtigungsfehlern unterscheiden. Überprüfen Sie unbedingt das tatsächliche Öffnungsergebnis, da sich der Status zwischen der Prüfung auf Existenz und dem Öffnen ändern kann.

## Niedriges Niveau I/O

```text
std::io::fd
io_read(fd: i64, buf: ptr<u8>, len: i64) -> i64
io_write(fd: i64, buf: ptr<u8>, len: i64) -> i64
io_read_exact(fd: i64, buf: ptr<u8>, len: i64) -> i64
io_write_all(fd: i64, buf: ptr<u8>, len: i64) -> i64
io_close(fd: i64) -> i64
```

Das positive Ergebnis von `io_read` ist die Anzahl der gelesenen Bytes, und 0 in einer Anfrage mit positiver Länge ist EOF. `io_write` kann weniger als angefordert geschrieben werden. Wenn eine vollständige Übertragung erforderlich ist, verwenden Sie die Funktion exact/all. Dennoch gehen wir nicht davon aus, dass der Fehler den externen Status wiederherstellt, da einige Übertragungen möglicherweise bereits vor dem Fehler stattgefunden haben.

`io_read_exact` ist ein Fehler, wenn EOF vor der erforderlichen Länge auftritt. `read_into` gibt `IO_ERR_NO_SPACE` zurück, wenn der Puffer voll ist und möglicherweise bereits einige Bytes geschrieben wurden. Die Lesefunktion hängt NUL nicht automatisch an das Ende der Zeichenfolge an.

## Buffer und Fehlerbehandlung

Bei einem Fehler stellt `read_to_end` die ursprüngliche Länge wieder her, aber ihre Kapazität und Datenadresse haben sich möglicherweise geändert. Der Anrufer muss den Buffer nach Erfolg oder Misserfolg freigeben. APIs zum Schreiben von Dateien garantieren keine atomare Dateiersetzung.

Ab [Übung zum Lesen von Dateien](/docs/de/practice/file-reader) können Sie das Programm von import bis zur Veröffentlichung ausführen. Berücksichtigen Sie die Pfad-/Berechtigungsunterschiede in Linux/macOS/Windows/FreeBSD und die Einschränkungen für zugängliche Verzeichnisse in WASI. Der Deskriptorwert wird nicht direkt als Rohhandle eines anderen OS interpretiert.

## Große Dateien in kleine Puffer einlesen

Vorgänge, bei denen nicht die gesamte Datei im Speicher abgelegt werden muss, können mit festen Puffern und Leseiterationen verarbeitet werden. Das folgende Programm druckt den Inhalt von input.txt und zählt die Gesamtzahl der gelesenen Bytes. Speichern Sie ein `Wave` und ein LF in der Eingabedatei.

<!-- wave-example: book-file-stream -->
```wave
import("std::fs::file")::{open_read};
import("std::io::fd")::{io_read, io_write_all, io_close};
import("std::io::consts")::{IO_STDOUT_FD};

fun main() -> i32 {
    var descriptor: i64 = open_read("input.txt");

    if (descriptor < 0) {
        return 1;
    }

    var buffer: array<u8, 4>;
    var total: i64 = 0;

    while (true) {
        var count: i64 = io_read(descriptor, &buffer[0], 4);

        if (count < 0) {
            io_close(descriptor);
            return 2;
        }

        if (count == 0) {
            break;
        }

        if (io_write_all(IO_STDOUT_FD, &buffer[0], count) < 0) {
            io_close(descriptor);
            return 3;
        }

        total += count;
    }

    if (io_close(descriptor) < 0) {
        return 4;
    }

    println("bytes={}", total);
    return 0;
}
```

Ausführungsergebnis:

```text
Wave
bytes=5
```

Die Pufferkapazität beträgt 4, aber der letzte Lesevorgang kann 1 Byte umfassen. Wir übergeben immer den tatsächlichen count an die Ausgabe. Wenn Sie das gesamte Array schreiben, können auch alte, ungelesene Bytes ausgegeben werden.

Das Programm hat descriptor von input.txt geöffnet, also schließen Sie es. Die Standardausgabe ist in dieser Funktion keine neu erworbene Ressource und wird daher am Ende des Beispiels nicht willkürlich geschlossen.

## Vollständige Lesevorgänge und unzureichende Kapazität

read_into erhält einen festen Speicherplatz, der die gesamte Datei enthält. Wenn der Platz nicht ausreicht, wird der Text stillschweigend abgeschnitten und NO_SPACE ohne Erfolg zurückgegeben.

<!-- wave-example: book-file-capacity -->
```wave
import("std::fs::file")::{read_into};
import("std::io::consts")::{IO_ERR_NO_SPACE};

fun main() {
    var data: array<u8, 2>;
    var status: i64 = read_into("input.txt", &data[0], 2);

    if (status == IO_ERR_NO_SPACE) {
        println("destination too small");
    }
}
```

Ausführungsergebnis:

```text
destination too small
```

Es wird nicht davon ausgegangen, dass der fehlgeschlagene Lesevorgang das Zielbyte überhaupt nicht geändert hat. Verwenden Sie es nicht als fertigen Dateiinhalt, bereiten Sie ein größeres Repository vor oder wählen Sie die Methode streaming. Selbst wenn Sie zuerst die Größe abfragen, ist das tatsächliche Leseergebnis die endgültige Beurteilung, da sich die Datei zwischen Abfrage und Lesevorgang ändern kann.

## Unterschied zwischen Schreiben und Anhängen von Dateien

write ersetzt den vorhandenen Inhalt und append wird am Ende hinzugefügt. Die Anzahl der zu speichernden Bytes kann direkt aus der Stringlänge ermittelt und übergeben werden. Das NUL am Ende der Zeichenfolge ist normalerweise nicht im Inhalt der Textdatei enthalten.

<!-- wave-example: book-file-write-append -->
```wave
import("std::fs::file")::{write, append, size, remove};
import("std::string::len")::{len};

fun main() -> i32 {
    var first: str = "Wave";
    var second: str = " study";

    if (write("output.txt", first as ptr<u8>, len(first) as i64) < 0) {
        return 1;
    }

    if (append("output.txt", second as ptr<u8>, len(second) as i64) < 0) {
        return 2;
    }

    println("bytes={}", size("output.txt"));

    if (remove("output.txt") < 0) {
        return 3;
    }

    return 0;
}
```

Ausführungsergebnis:

```text
bytes=10
```

In diesem Beispiel wird output.txt im Arbeitsverzeichnis erstellt, ersetzt und schließlich gelöscht. Vom Übungsverzeichnis aus ausführen, ohne dass Dateien vorhanden sind. Für den eigentlichen Editor oder das Speicherprogramm sind möglicherweise separate Speicherrichtlinien erforderlich, z. B. für temporäre Dateien und Ersatz.

## API Auswahltabelle

|Situation|auswählen|
| --- | --- |
|Kleine gesamte Datei in festen Puffer einlesen| read_into |
|Behalten Sie den gesamten Inhalt, ohne die Größe zu kennen|read_to_end und Buffer|
|Verarbeiten Sie Inhalte der Reihe nach, anstatt sie vollständig zu speichern|open_read + io_read wiederholen|
|Datensätze mit fester Länge lesen| io_read_exact |
|Gesamte Bytefolge übertragen| io_write_all |
|Umgang mit bereits geöffneten Dateien|fd-Funktion statt Pfadfunktion|

Überprüfen Sie nach Auswahl einer Funktion, wie sich Puffer, Dateispeicherort und externe Daten im Fehlerfall ändern.
