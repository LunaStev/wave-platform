---
translation_set_id: stdlib-bytes
path: stdlib/bytes
locale: de
group: stdlib
group_order: 1
order: 6
title: bytes: Bereiche, Cursor und ULEB128
summary: Beschreibt das Lesen/Schreiben von Bytes mit der Länge view und die Beibehaltung des Status im Fehlerfall.
---

## Wie unterscheidet es sich von einer Saite?

Bytedaten können Nullen enthalten. Übergeben Sie daher einen Zeiger zusammen mit einer Länge. `Bytes` und `BytesMut` sind nicht besitzende Ansichten und sind nur gültig, solange der zugrunde liegende Speicher gültig bleibt. `BytesMut` erfordert beschreibbaren Speicher.

## Erstellen Sie cursor

```text
std::bytes::cursor
bytes_reader(data: ptr<u8>, len: i64) -> ByteReader
bytes_writer(data: ptr<u8>, len: i64) -> ByteWriter
bytes_reader_read_be_u16(reader: ptr<ByteReader>, out_value: ptr<u16>) -> i32
bytes_writer_write_be_u16(writer: ptr<ByteWriter>, value: u16) -> i32
```

Importieren Sie `ByteReader` und `ByteWriter` aus `std::bytes::types`. Ihr Positionsfeld identifiziert den Standort des nächsten Vorgangs. Ihre Länge entspricht der gesamten zugänglichen Bytezahl. Durch das Erstellen eines Cursors wird der zugrunde liegende Speicher nicht kopiert oder zugewiesen.

„be“ ist Big-Endian, „le“ ist Little-Endian. Wenn das Dateiformat Big-Endian ist, wird die Funktion „be“ unabhängig von der Bytereihenfolge der Host-CPU verwendet. Es gibt 16-, 32- und 64-Bit-Lese-/Schreibfunktionen mit/ohne Vorzeichen und Ein-Byte-Funktionen.

## Fehler und Zustandserhaltung

`BYTES_OK` von `std::bytes::errors` ist 0. INVALID zeigt einen ungültigen Bereich an, EOF unzureichende Eingabe, NO_SPACE unzureichende Ausgabekapazität und OVERFLOW einen Wert außerhalb des darstellbaren Bereichs. Überprüfte Cursoroperationen erhöhen die Position erst, nachdem die gesamte Operation erfolgreich war. Bei einem fehlgeschlagenen Lesevorgang bleibt auch der Ausgabewert erhalten.

## ULEB128

```text
std::bytes::leb128
bytes_reader_read_uleb128_u64(reader: ptr<ByteReader>, output_value: ptr<u64>) -> i32
bytes_writer_write_uleb128_u64(writer: ptr<ByteWriter>, value: u64) -> i32
```

ULEB128 speichert eine vorzeichenlose 64-Bit-Ganzzahl in einer variablen Anzahl von Bytes, wobei maximal 10 Bytes verwendet werden. Wenn nicht genügend Platz vorhanden ist, behält der Writer sowohl seine Position als auch die Zielbytes bei. Das Lesegerät unterscheidet eine unvollständige Eingabe von einem Wert größer als u64. Terminierte, nicht minimale Kodierungen werden akzeptiert.

Um eine tatsächliche Nachricht zu erstellen und einen kurzen Eingabefehler anzuzeigen, fahren Sie mit [Binäre Nachrichtenpraxis](/docs/de/practice/binary-message) fort. Versuchen Sie nicht, eine Bytefolge mit Nullen als `str` auszugeben.

## Lesen derselben Bytes in unterschiedlicher Reihenfolge

Die Bytereihenfolge ist die Speicherregel für Zahlen. Wenn Sie die beiden Bytes 1 und 2 als big-endian lesen, ist es 1×256+2, und wenn Sie sie als little-endian lesen, ist es 2×256+1. Wählen Sie basierend auf Netzwerk- oder Dateitypregeln aus.

<!-- wave-example: book-bytes-endian -->
```wave
import("std::bytes::types")::{Bytes, bytes_view};
import("std::bytes::read")::{bytes_read_be_u16, bytes_read_le_u16};

fun main() -> i32 {
    var data: array<u8, 2> = [1, 2];
    var view: Bytes = bytes_view(&data[0], 2);
    var big: u16 = 0;
    var little: u16 = 0;

    if (bytes_read_be_u16(view, 0, &big) < 0) {
        return 1;
    }

    if (bytes_read_le_u16(view, 0, &little) < 0) {
        return 2;
    }

    println("be={} le={}", big, little);
    return 0;
}
```

Ausführungsergebnis:

```text
be=258 le=513
```

view leiht sich ein Array aus. Da es keine separate Zuordnung oder Kopie gibt, kann das Array nur verwendet werden, solange es gültig ist. offset wird in Bytes angegeben und das Lesen von 16 Bits erfordert 2 Bytes ab dieser Position.

## Wählen Sie offset API und cursor API

Die Funktion read/write, die das Argument offset verwendet, ist praktisch für Formate, die eine angegebene Feldposition direkt lesen. Für Streams, bei denen die nächste Position von der Länge des vorhergehenden Felds abhängt, ist cursor mit position praktisch.

Machen Sie beim Mischen der beiden deutlich, welches der Standard ist: cursor.position oder getrennt offset. Vermeiden Sie den Fehler, denselben Standort zweimal hinzuzufügen oder ohne erfolgreichen Lesevorgang zum nächsten Standort zu wechseln.

## Überprüfen Sie den Status anhand einer kurzen Eingabe

<!-- wave-example: book-bytes-short-read -->
```wave
import("std::bytes::types")::{ByteReader};
import("std::bytes::cursor")::{bytes_reader, bytes_reader_read_be_u16};
import("std::bytes::errors")::{BYTES_ERROR_EOF};

fun main() {
    var data: array<u8, 1> = [1];
    var reader: ByteReader = bytes_reader(&data[0], 1);
    var value: u16 = 99;
    var status: i32 = bytes_reader_read_be_u16(&reader, &value);

    if (status == BYTES_ERROR_EOF) {
        println("position={} value={}", reader.position, value);
    }
}
```

Ausführungsergebnis:

```text
position=0 value=99
```

Da keine zwei Bytes erforderlich sind, bleiben der Ausgabewert und der Speicherort erhalten. Diese Eigenschaft ist in Designs nützlich, die erneut versuchen, dasselbe Feld zu lesen, nachdem weitere Eingaben erfasst wurden. Wenn jedoch der Speicherplatz, auf den view verweist, neu zugewiesen wurde, muss auch die Adresse aktualisiert werden.

## Grenze von ULEB128

0~127 verwendet ein Byte, ab 128 werden mehr Bytes verwendet. Das höherwertige Bit jedes Bytes gibt an, ob Daten folgen. Ein Wert außerhalb von u64 oder eine kontinuierliche Eingabe, die zu lang ist, ist OVERFLOW, was sich von EOF unterscheidet, das einfach weniger Eingaben hat.

Überprüfen Sie direkt, ob „Kapazität nicht ausreichend“ writer seinen Status beibehält.

<!-- wave-example: book-leb128-space -->
```wave
import("std::bytes::types")::{ByteWriter};
import("std::bytes::cursor")::{bytes_writer};
import("std::bytes::leb128")::{bytes_writer_write_uleb128_u64};
import("std::bytes::errors")::{BYTES_ERROR_NO_SPACE};

fun main() {
    var data: array<u8, 1> = [85];
    var writer: ByteWriter = bytes_writer(&data[0], 1);
    var status: i32 = bytes_writer_write_uleb128_u64(&writer, 128);

    if (status == BYTES_ERROR_NO_SPACE) {
        println("position={} byte={}", writer.position, data[0]);
    }
}
```

Ausführungsergebnis:

```text
position=0 byte=85
```

128 erfordert zwei Bytes, aber nur ein Leerzeichen. Nach einem Ausfall bleibt auch das erste Byte 85 bestehen. Die gemeinsame Überprüfung von Fehlertypen und Zustandserhaltung beschreibt die Grenze besser als eine einfache roundtrip-Erfolgsprüfung.

## Reihenfolge der Nachrichtenparser-Erstellung

1. Lesen Sie den festen Header und überprüfen Sie den Typ und die Version.
2. Lesen Sie die Länge ab und vergleichen Sie sie mit dem verbleibenden Eingabebereich.
3. Übergeben Sie nur die notwendigen Daten an view oder einen separaten Puffer.
4. Wenn das Format die gesamte Nachricht erfordert, werden auch die zusätzlichen Bytes überprüft.
5. Unterscheidet zwischen EOF und ungültigen Formatfehlern und gibt diese an den Aufrufer weiter.

Sie können ein Programm erstellen, indem Sie Felder in [Binäre Nachrichtenpraxis](/docs/de/practice/binary-message) verbinden.
