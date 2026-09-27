---
translation_set_id: practice-binary-message
path: practice/binary-message
locale: de
group: practice
group_order: 4
order: 3
title: Projekt: Erstellen einer Binärnachricht
summary: Verwenden Sie die explizite Bytereihenfolge und ULEB128 und lehnen Sie kurze Eingaben ab.
---

## Nachrichtenformat

Die ersten 2 Bytes speichern die Typnummern big-endian u16 und dann die Werte ULEB128 u64. Wenn Sie den Strukturspeicher so wie er ist in eine Datei schreiben, wird er von der Auffüllung und der Bytereihenfolge beeinflusst. Codieren Sie ihn daher nach Feld.

Speichern Sie es als `main.wave` und führen Sie es aus.

<!-- wave-example: binary-message -->
```wave
import("std::bytes::types")::{
    ByteReader, ByteWriter
};
import("std::bytes::cursor")::{
    bytes_reader, bytes_writer, bytes_writer_write_be_u16, bytes_reader_read_be_u16
};
import("std::bytes::leb128")::{
    bytes_writer_write_uleb128_u64, bytes_reader_read_uleb128_u64
};

fun main() -> i32 {
    var data: array<u8, 12>;
    var writer: ByteWriter = bytes_writer(&data[0], 12);
    if (bytes_writer_write_be_u16(&writer, 7) < 0) {
        return 1;
    }
    if (bytes_writer_write_uleb128_u64(&writer, 300) < 0) {
        return 2;
    }

    var reader: ByteReader = bytes_reader(&data[0], writer.position);
    var kind: u16 = 0;
    var value: u64 = 0;
    if (bytes_reader_read_be_u16(&reader, &kind) < 0) {
        return 3;
    }
    if (bytes_reader_read_uleb128_u64(&reader, &value) < 0) {
        return 4;
    }

    println("kind={} value={} bytes={}", kind, value, reader.position);
    var short: ByteReader = bytes_reader(&data[0], 1);
    kind = 99;
    if (bytes_reader_read_be_u16(&short, &kind) >= 0) {
        return 5;
    }

    println("preserved={} {}", short.position, kind);
    return 0;
}
```

Ausführungsergebnis:

```text
kind=7 value=300 bytes=4
preserved=0 99
```

## Was zu überprüfen ist

In reader übergeben wir die tatsächliche geschriebene Länge, nicht die gesamte Array-Kapazität von 12. Dadurch soll vermieden werden, dass nicht initialisierte nachfolgende Bytes als Eingabe gelesen werden. Selbst wenn ein kurzer Eingabelesevorgang fehlschlägt, bleiben position=0 und kind=99 erhalten.

## Erweiterte Übungen und Kommentare

Wenn das Format am Ende keine zusätzlichen Bytes zulässt, überprüfen Sie `reader.position == reader.len`, nachdem die Analyse abgeschlossen ist. Stellen Sie beim Hinzufügen eines Längenfelds sicher, dass es nicht größer als die verbleibenden Bytes der Eingabe ist und dass die Berechnung von Länge+offset den Bereich nicht überschreitet.

[Siehe bytes](/docs/de/stdlib/bytes)

## Schauen Sie sich die tatsächlichen Bytes an

Typnummer 7 ist big-endian u16 und daher `00 07`. Der Wert 300 wird zu ULEB128 bis `AC 02`. Die gesamte Nachricht besteht aus den folgenden vier Bytes:

```text
00 07 AC 02
───── ─────
종류  값
```

ULEB128 verwendet die niedrigen 7 Bits für jedes Byte für den Wert, und wenn das hohe Bit 1 ist, zeigt es an, dass das nächste Byte folgt. Die unteren 7 Bits von 300 sind 44 und der Rest ist 2. Das erste Byte ist 44 plus 128 aufeinanderfolgende Markierungen oder 172 oder 0xAC. Im letzten Byte 0x02 gibt es keine Fortsetzungsmarke.

## Kapazität und Nutzungsdauer

u16 verwendet 2 Bytes und ULEB128 von u64 verwendet bis zu 10 Bytes, sodass ein 12-Byte-Array beide Felder speichern kann. Ein Wert von 300 verwendet nur 2 Bytes, sodass die eigentliche Nachricht 4 Bytes groß ist. Beim Senden an eine Datei oder einen Socket senden Sie das Byte writer.position und nicht das gesamte Array.

position in reader ist die aktuelle Leseposition. Nach dem Lesen des Typs wird er zu 2 und nach dem Lesen des Werts zu 4. Um bei Auftreten eines Fehlers mit der nächsten Nachricht fortzufahren, muss die Grenze der fehlgeschlagenen Nachricht bekannt sein. Eine Nachrichtenwiederherstellungsrichtlinie wird nicht automatisch erstellt, nur weil die Lesefunktion den Speicherort beibehält.

## Grenzwertübung

Ändern Sie die Werte auf 0, 127, 128, 16383, 16384, um die Codierungslänge zu bestimmen. Beim Wechsel von 127 auf 128 erhöht sich die Länge von ULEB128 von 1 auf 2, und beim Wechsel von 16383 auf 16384 erhöht sich die Länge von 2 auf 3. Berechnen Sie die Gesamtlänge, einschließlich der 2 Bytes des Typfelds.

Nachrichten, bei denen das letzte Byte abgeschnitten ist, werden ebenfalls überprüft. Wenn die ursprüngliche Nachrichtenlänge 4 war, übergeben wir Länge 3 an reader. Das Typfeld wird gelesen, aber das Lesen des Werts muss fehlschlagen, da das letzte Byte von ULEB128 fehlt. Überprüfen Sie zu diesem Zeitpunkt, ob Position 2 und der Ausgangswert kurz vor dem Ablesen von ULEB128 beibehalten werden.
