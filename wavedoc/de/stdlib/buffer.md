---
translation_set_id: stdlib-buffer
path: stdlib/buffer
locale: de
group: stdlib
group_order: 1
order: 5
title: buffer: Erweiterbarer Bytespeicher
summary: Buffer Beschreibt Initialisierungs-, Hinzufügungs-, Abfrage-, Kapazitäts- und Freigaberegeln.
---

## Bedeutung von Buffer

`Buffer` in `std::buffer::types` hat `data: ptr<u8>`, `len: i64` und `cap: i64`. len ist die Anzahl der initialisierten und verwendeten Bytes und cap ist die Gesamtzahl der zugewiesenen Bytes. Halten Sie immer `0 <= len <= cap` ein. Die Zeichenfolge NUL garantiert nicht automatisch die Beendigung.

## Basic API

|Modul|Erklärung|Bedeutung|
| --- | --- | --- |
| `std::buffer::alloc` | `buffer_init(out_buf: ptr<Buffer>, capacity: i64) -> i64` |Neues Repository initialisieren. Nicht in bereits besessene Puffer zurückrufen|
|gleiches Modul| `buffer_free(buf: ptr<Buffer>) -> i64` |Freigabe aufheben. Bei Erfolg leer|
|gleiches Modul| `buffer_reserve(buf: ptr<Buffer>, required_cap: i64) -> i64` |Stellen Sie sicher, dass die Mindestkapazität voll ist. len Gepflegt|
|gleiches Modul| `buffer_clear(buf: ptr<Buffer>) -> i64` |Kapazität beibehalten und len=0|
|gleiches Modul| `buffer_resize(buf: ptr<Buffer>, new_len: i64, value: u8) -> i64` |Länge ändern, neue Bytes mit value füllen|
| `std::buffer::write` | `buffer_push(buf: ptr<Buffer>, value: u8) -> i64` |Füge ein Byte hinzu|
|gleiches Modul| `buffer_append(buf: ptr<Buffer>, src: ptr<u8>, size: i64) -> i64` |sizeByte kopieren und hinzufügen|
|gleiches Modul| `buffer_append_str(buf: ptr<Buffer>, s: str) -> i64` |String-Bytes mit Ausnahme von NUL hinzufügen|
| `std::buffer::read` | `buffer_get(buf: ptr<Buffer>, index: i64, out_value: ptr<u8>) -> i64` |Ein Byte im Bereich lesen|

Statusrückgabe API gibt `BUFFER_OK`(0) bei Erfolg zurück. Unterscheiden Sie zwischen den Fehlern INVALID, BOUNDS, OVERFLOW und ALLOC von `std::buffer::error`. Die Nummer wird nicht als OS errno interpretiert. `buffer_new` stellt einen Zuordnungsfehler als leeres Buffer dar. Wenn Sie also zwischen Fehlern unterscheiden müssen, verwenden Sie `buffer_init`.

## Laufbeispiel

<!-- wave-example: buffer-api -->
```wave
import("std::buffer::alloc")::{
    Buffer, buffer_init, buffer_free
};
import("std::buffer::write")::{
    buffer_append_str, buffer_push
};
import("std::buffer::read")::{
    buffer_get
};

fun main() -> i32 {
    var data: Buffer;
    if (buffer_init(&data, 0) < 0) {
        return 1;
    }
    if (buffer_append_str(&data, "Hi") < 0 || buffer_push(&data, 33) < 0) {
        buffer_free(&data);
        return 2;
    }

    var value: u8 = 0;
    if (buffer_get(&data, 2, &value) < 0) {
        buffer_free(&data);
        return 3;
    }

    println("{} {}", data.len, value);
    if (buffer_free(&data) < 0) {
        return 4;
    }

    return 0;
}
```

Ausführungsergebnis:

```text
3 33
```

Speichern Sie es als `main.wave` und führen Sie es als `wavec run main.wave` aus. Eine anfängliche Kapazität von 0 ist kein Fehler, sondern ein gültiger leerer Puffer. Schafft Platz bei der Weiterverarbeitung.

## Lebensdauer und Ausfall

Durch die Vergrößerung des Puffers können sich Daten ändern. Verwenden Sie keine zuvor geliehene Adresse nach einem Vorgang, bei dem es zu einer Neuzuweisung kommen könnte. Durch das Kopieren der Buffer-Struktur wird deren Zuordnung nicht dupliziert. Geben Sie dieser Zuordnung also einen einzigen Eigentümer.

`buffer_get` ändert die Ausgabeargumente nicht, wenn dies fehlschlägt. Andererseits stellt die Komfortfunktion `buffer_at` Fehler auch als 0 dar. Verwenden Sie daher `buffer_get`, um zwischen tatsächlichen 0 Bytes und Fehlern zu unterscheiden. Vermeiden Sie die Erstellung ungültiger len/cap, indem Sie öffentliche Felder direkt ändern.

[Erinnerung API](/docs/de/reference/memory-and-buffer) · [Üben Sie das Lesen einer Datei als Buffer](/docs/de/practice/file-reader)

## Länge und Fassungsvermögen separat beachten

reserve gibt Speicherplatz frei, erhöht aber nicht len. resize ändert die tatsächlich verwendete Länge und initialisiert den erweiterten Teil auf das angegebene Byte. clear setzt nur die verwendete Länge auf 0, sodass die Zuordnung wiederverwendet werden kann.

Speichern Sie das folgende Programm als main.wave und führen Sie es aus. Dabei kommt es nicht auf das exakte Wachstumsmultiplikator der Kapazität an; Es stellt lediglich sicher, dass Sie den Platz haben, den Sie benötigen.

<!-- wave-example: book-buffer-length-capacity -->
```wave
import("std::buffer::alloc")::{
    Buffer,
    buffer_init,
    buffer_reserve,
    buffer_resize,
    buffer_clear,
    buffer_free
};

fun main() -> i32 {
    var data: Buffer;

    if (buffer_init(&data, 0) < 0) {
        return 1;
    }

    if (buffer_reserve(&data, 20) < 0) {
        buffer_free(&data);
        return 2;
    }

    println("reserved length={}", data.len);

    if (buffer_resize(&data, 3, 7) < 0) {
        buffer_free(&data);
        return 3;
    }

    println("resized length={} first={}", data.len, deref data.data[0]);

    if (buffer_clear(&data) < 0) {
        buffer_free(&data);
        return 4;
    }

    println("cleared length={}", data.len);

    if (buffer_free(&data) < 0) {
        return 5;
    }

    return 0;
}
```

Ausführungsergebnis:

```text
reserved length=0
resized length=3 first=7
cleared length=0
```

Beim Erhöhen auf resize haben wir value=7 übergeben, sodass alle drei neuen Bytes, die wir sehen, 7 sind. Nur mit reserve gesicherter Speicherplatz wird nicht als initialisierte Daten gelesen. cap bleibt nach clear bestehen und kann wieder zu demselben Buffer hinzugefügt werden.

## Unterscheiden Sie zwischen Nullbytes und Suchfehlern

buffer_get gibt den Status zurück und schreibt die tatsächlichen Bytes als Ausgabeargumente. Auch wenn die Daten 0 sind, handelt es sich um einen normalen Erfolg.

<!-- wave-example: book-buffer-zero-bounds -->
```wave
import("std::buffer::alloc")::{Buffer, buffer_init, buffer_free};
import("std::buffer::write")::{buffer_push};
import("std::buffer::read")::{buffer_get};
import("std::buffer::error")::{BUFFER_ERR_BOUNDS};

fun main() -> i32 {
    var data: Buffer;

    if (buffer_init(&data, 0) < 0) {
        return 1;
    }

    if (buffer_push(&data, 0) < 0) {
        buffer_free(&data);
        return 2;
    }

    var value: u8 = 99;

    if (buffer_get(&data, 0, &value) < 0) {
        buffer_free(&data);
        return 3;
    }

    println("stored={}", value);
    value = 99;

    if (buffer_get(&data, 1, &value) == BUFFER_ERR_BOUNDS) {
        println("outside, preserved={}", value);
    }

    if (buffer_free(&data) < 0) {
        return 4;
    }

    return 0;
}
```

Ausführungsergebnis:

```text
stored=0
outside, preserved=99
```

Der erste Treffer ist ein Erfolg mit der Anzeige 0, der zweite Treffer ist ein Fehlschlag im Aus. Auch wenn value=99 nach einem Fehler bestehen bleibt, heißt das nicht, dass es sich um den aus dem Puffer gelesenen Wert handelt. Überprüfen Sie unbedingt gemeinsam den Status.

## Übungslösung: Byte-Akkumulation

Um die Zahlen 0 bis 9 hinzuzufügen, wiederholen Sie buffer_push und überprüfen Sie jedes Ergebnis. Speichern Sie die Summe in i64 und lesen Sie nur den Bereich `0 <= index < data.len`. Nachdem wir den Puffer bearbeitet haben, rufen wir buffer_free sowohl auf dem Erfolgs- als auch auf dem Fehlerpfad auf.

<!-- wave-example: book-buffer-sum -->
```wave
import("std::buffer::alloc")::{Buffer, buffer_init, buffer_free};
import("std::buffer::write")::{buffer_push};

fun main() -> i32 {
    var data: Buffer;

    if (buffer_init(&data, 0) < 0) {
        return 1;
    }

    for (var value: i32 = 0; value < 10; value += 1) {
        if (buffer_push(&data, value as u8) < 0) {
            buffer_free(&data);
            return 2;
        }
    }

    var total: i64 = 0;

    for (var index: i64 = 0; index < data.len; index += 1) {
        total += deref data.data[index] as i64;
    }

    println("sum={}", total);

    if (buffer_free(&data) < 0) {
        return 3;
    }

    return 0;
}
```

Ausführungsergebnis:

```text
sum=45
```
