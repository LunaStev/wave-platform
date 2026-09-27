---
translation_set_id: practice-file-reader
path: practice/file-reader
locale: de
group: practice
group_order: 4
order: 2
title: Projekt: Eine Datei lesen und Bytes zählen
summary: Buffer, Datei I/O, verknüpft Fehlerbehandlung und Speicherfreigabe.
---

## fertig

Erstellen Sie input.txt im Arbeitsverzeichnis, das `Wave` enthält, gefolgt von einer einzelnen LF-Neuzeile. Die Datei enthält dann 5 Bytes. Bei CRLF sind es 6 Bytes; Ein UTF-8-BOM fügt weitere Bytes hinzu. Überprüfen Sie die Dateikodierung und Zeilenenden des Editors.

Speichern Sie das Programm als `main.wave` und führen Sie es als `wavec run main.wave` im selben Verzeichnis aus. Relative Pfade beziehen sich auf das laufende Arbeitsverzeichnis.

<!-- wave-example: file-reader -->
```wave
import("std::buffer::alloc")::{
    Buffer, buffer_init, buffer_free
};
import("std::fs::file")::{
    read_to_end
};

fun main() -> i32 {
    var data: Buffer;
    if (buffer_init(&data, 0) < 0) {
        return 1;
    }

    var count: i64 = read_to_end("input.txt", &data);
    if (count < 0) {
        println("read failed={}", count);
        buffer_free(&data);
        return 2;
    }

    var lines: i64 = 0;
    for (var i: i64 = 0; i < data.len; i += 1) {
        if (deref data.data[i] == 10) {
            lines += 1;
        }
    }

    println("bytes={} LF={}", count, lines);
    if (buffer_free(&data) < 0) {
        return 3;
    }

    return 0;
}
```

Ausführungsergebnis:

```text
bytes=5 LF=1
```

## Handlungen und Verantwortlichkeiten

`read_to_end` öffnet und schließt die Datei, die Freigabe von Buffer liegt jedoch in der Verantwortung des Anrufers. Wird sowohl auf Erfolgs- als auch auf Lesefehlerpfaden deaktiviert. Da es sich um Daten mit einer Anzahl von Bytes handelt, gehen wir nicht davon aus, dass es sich um eine Zeichenfolge handelt, die mit NUL endet.

LF Die Anzahl der Zeilen und die Anzahl der Zeilen, die Menschen denken, sind nicht immer gleich. Wenn die letzte Zeile nicht LF enthält, wird sie nicht in die LF-Zählung für dieses Programm einbezogen. Überprüfen Sie außerdem, ob Lesefehler vorliegen, indem Sie die Datei umbenennen. Spezifische Fehlernummern können je nach Umgebung variieren.

## Erweiterte Übungen und Kommentare

Wenn Sie sehr große Dateien verarbeiten möchten, ersetzen Sie sie durch ein Array fester Größe und eine `io_read`-Iteration. Verarbeitet nur positive Rückgabebereiche und endet bei 0. Sie können Byte-Zählungen und LF-Zählungen akkumulieren, ohne die gesamte Datei im Speicher behalten zu müssen. Wenn Sie es selbst geöffnet haben, wird der Deskriptor auch auf jedem Exit-Pfad geschlossen.

[Siehe fs und io](/docs/de/stdlib/files-io) · [Siehe Buffer](/docs/de/stdlib/buffer)

## Folgen Sie dem Verarbeitungsablauf

1. Bin Buffer initialisieren. Es ist noch kein Dateiinhalt vorhanden.
2. read_to_end liest die Datei und erhöht den benötigten Speicherplatz.
3. Bei erfolgreichem Lesevorgang werden die Bytes im Bereich data.len überprüft.
4. Immer wenn wir auf den Bytewert 10 von LF stoßen, erhöhen wir lines.
5. Drucken Sie die Ergebnisse aus und geben Sie Buffer frei.

data.cap ist der reservierte Speicherplatz und data.len ist die gültige Datenlänge. Wenn Sie die Wiederholungsbedingung in cap ändern, werden Bytes gelesen, die nicht in der Datei enthalten sind. Verwenden Sie daher len. count ist die Anzahl der Bytes, die in diesem Aufruf zu read_to_end hinzugefügt werden. Dieses Beispiel beginnt mit einem leeren Buffer, daher sind count und data.len gleich.

## Ändern Sie die Eingabe, um sie zu überprüfen

|input.txt Inhalt| bytes | LF |Grund|
| --- | --- | --- | --- |
|leere Datei| 0 | 0 |Keine Bytes zum Lesen|
| `Wave` | 4 | 0 |Kein letzter Zeilenumbruch|
| `Wave` + LF | 5 | 1 |Daten bis zum letzten LF|
| `A` + LF + `B` + LF | 4 | 2 |Zähle zwei LF|
| `Wave` + CRLF | 6 | 1 |CR ist ebenfalls 1 Byte, es wird aber nur LF gezählt.|

Um Dateien ohne LF in der letzten Zeile als Zeile zu zählen, addieren Sie 1 zur Anzahl der Zeilen, wenn die Datei nicht leer ist und das letzte Byte nicht 10 ist. Sie müssen zunächst prüfen, ob data.len 0 ist, bevor Sie auf das letzte Element zugreifen können.

## Auf große Dateien erweitern

Die aktuelle Methode zur Archivierung des gesamten Inhalts ist praktisch für ein späteres erneutes Lesen oder Abrufen der Daten. Wenn Sie nur die Anzahl der Bytes und die Anzahl LF benötigen, ist die Wiederverwendung eines Puffers mit fester Größe sinnvoll.

Zählen Sie in der Schleife io_read im Beispiel [Einlesen einer Datei in einen Puffer fester Größe](/docs/de/stdlib/files-io) einfach LF mit der Anzahl der zurückgegebenen Bytes. Es kann mit Speicher verarbeitet werden, der der Puffergröße entspricht, nicht der gesamten Länge der Datei. Wenn der Lesevorgang 0 zurückgibt, ist er beendet. Wenn es negativ ist, ist es ein Fehler.
