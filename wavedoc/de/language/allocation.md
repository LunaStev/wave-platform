---
translation_set_id: learn-allocation
path: language/allocation
locale: de
group: language
group_order: 2
order: 10
title: 10. Speicherzuweisung und Ressourcenverwaltung
summary: Erfahren Sie mehr über Zuordnungsfehler, Initialisierung, Umfang und Freigabe.
---

## Wann benötigen Sie dynamischen Speicherplatz?

Ein Array fester Größe schließt seine Länge in seinen Typ ein. Verwenden Sie dynamischen Speicher, wenn die Datenmenge erst zur Laufzeit bekannt ist, beispielsweise eine Dateigröße oder eine Eingabelänge. Geben Sie jede Zuteilung frei, wenn Sie sie nicht mehr benötigen.

In diesem Kapitel verwalten Sie eine kleine Zuteilung, ändern ihre Größe und verwenden dann einen Buffer. Das Übergeben eines Zeigers unterscheidet sich vom Übertragen des Eigentums. Identifizieren Sie, welche Ressourcen jede Funktion besitzt, während Sie den Beispielen folgen.

## Zuordnen, prüfen, nutzen und freigeben

<!-- wave-example: book-alloc-lifecycle -->
```wave
import("std::mem::alloc")::{
    mem_alloc_zeroed, mem_free
};

fun main() -> i32 {
    var size: i64 = 4;
    var data: ptr<u8> = mem_alloc_zeroed(size);

    if (data == null) {
        println("allocation failed");
        return 1;
    }

    deref data[0] = 42;
    println("{} {}", deref data[0], deref data[1]);

    var status: i64 = mem_free(data, size);

    if (status < 0) {
        return 2;
    }

    return 0;
}
```

Ausführungsergebnis:

```text
42 0
```

Das Programm besteht aus vier Schritten: 4 Bytes anfordern, auf null prüfen, nur auf den gültigen Bereich zugreifen und die Zuordnung freigeben. Bei Erfolg initialisiert mem_alloc_zeroed den Speicher auf Null, sodass das zweite Byte Null ist, obwohl das Programm nicht darauf geschrieben hat.

Gehen Sie nicht davon aus, dass der von mem_alloc zurückgegebene Anfangsinhalt für den Speicher vorhanden ist. Initialisieren Sie jede Region, bevor Sie sie lesen. Eine Zuordnungsgröße von Null oder negativ gibt null zurück. Eine Zuweisung positiver Größe kann auch dazu führen, dass kein Speicher abgerufen wird.

## Größeneinheit

Das Größenargument einer Speicherzuweisungs-API wird in Bytes gemessen. Um zehn Ganzzahlen zuzuweisen, multiplizieren Sie die Elementgröße mit der Anzahl der Elemente. Stellen Sie sicher, dass diese Multiplikation nicht überläuft.

<!-- wave-example: book-alloc-typed -->
```wave
import("std::mem::alloc")::{
    mem_alloc, mem_free
};
import("std::mem::layout")::{
    size_of
};
import("std::mem::ops")::{
    mem_size_mul_checked
};

fun main() -> i32 {
    var count: i64 = 3;
    var bytes: i64 = 0;
    var item_size: i64 = size_of<i32>() as i64;

    if (mem_size_mul_checked(count, item_size, &bytes) < 0) {
        return 1;
    }

    var raw: ptr<u8> = mem_alloc(bytes);

    if (raw == null) {
        return 2;
    }

    var values: ptr<i32> = raw as ptr<i32>;

    for (var index: i64 = 0; index < count; index += 1) {
        deref values[index] = (index + 1) as i32;
    }

    println("{} {} {}", deref values[0], deref values[1], deref values[2]);

    if (mem_free(raw, bytes) < 0) {
        return 3;
    }

    return 0;
}
```

Ausführungsergebnis:

```text
1 2 3
```

count ist eine Elementanzahl; Bytes ist eine Byteanzahl. Die Zeigerarithmetik bewegt sich in Einheiten von i32, aber um die Zuweisung freizugeben, ist ihre ursprüngliche Größe in Bytes erforderlich. size_of verwendet das Zieltyp-Layout und macht die Beziehung zum Elementtyp explizit.

Wenn Sie das Ergebnis von size_of in i64 für einen allgemein sehr großen Typ ändern, muss auch der Konvertierungsbereich berücksichtigt werden. Hier verwenden wir i32, das eine bekannte Größe hat.

## Bereinigen Sie auch Fehlerpfade

Wenn ein anderer Vorgang nach der Zuweisung fehlschlägt, geben Sie den Speicher frei, bevor Sie vorzeitig zurückkehren. Eine Eigentümertabelle hilft dabei, Bereinigungspfade zu identifizieren, die Ihnen andernfalls entgehen würden.

|Schritt|Ressourcen im Besitz von|Wenn Sie scheitern|
| --- | --- | --- |
|Vor der Zuteilung|Keine|gib es sofort zurück|
|Nach erfolgreicher Zuordnung|data und Originalgröße|data Rückgabe nach Freigabe|
|Nach erfolgreicher Umverteilung|Neue Adresse und neue Größe|Neue Adressveröffentlichung|
|Nach der Veröffentlichung|Keine|Verwenden Sie keine alte Adresse|

Durch das Überschreiben einer Zeigervariablen und den Verlust der ursprünglichen Adresse gehen auch die Informationen verloren, die zum Freigeben der Zuweisung erforderlich sind. Dies führt zu einem Speicherverlust. Umgekehrt führt die Freigabe derselben Zuteilung durch zwei Eigentümer zu einer doppelten Freigabe.

## Neuzuordnung zur Vergrößerung

<!-- wave-example: book-alloc-grow -->
```wave
import("std::mem::alloc")::{
    mem_alloc, mem_realloc, mem_free
};

fun main() -> i32 {
    var data: ptr<u8> = mem_alloc(4);

    if (data == null) {
        return 1;
    }

    deref data[0] = 7;
    var next: ptr<u8> = mem_realloc(data, 4, 8);

    if (next == null) {
        mem_free(data, 4);
        return 2;
    }

    data = next;
    deref data[4] = 9;
    println("{} {}", deref data[0], deref data[4]);

    if (mem_free(data, 8) < 0) {
        return 3;
    }

    return 0;
}
```

Ausführungsergebnis:

```text
7 9
```

Speichern Sie das Ergebnis zunächst in next. Wenn die Zuweisung des größeren Blocks fehlschlägt, bleiben die ursprünglichen Daten gültig und können weiterhin freigegeben werden. Bei Erfolg wird die alte Zuordnung freigegeben und die neue Adresse muss verwendet werden. Initialisieren Sie die neu hinzugefügte Region, bevor Sie sie lesen.

Die neue Größe in diesem Beispiel ist positiv. Eine Anfrage mit new_size=0 versucht stattdessen, die vorhandene Zuordnung freizugeben und gibt null zurück. Daher bedeutet ein null-Ergebnis nicht immer, dass die alte Zuordnung noch gültig ist. Rufen Sie mem_free direkt an, wenn Sie überprüfen möchten, ob die Freigabe erfolgreich war.

## Wann sollten geliehene Zeiger erneut überprüft werden?

data Es ist falsch, einen internen Zeiger zu speichern und ihn nach der Neuzuweisung zu verwenden. Dies liegt daran, dass die Adresse des neuen data möglicherweise anders ist. Wenn Sie einen internen Standort benötigen, können Sie anstelle der Adresse offset hinterlegen und nach Erfolg anhand des neuen data neu berechnen.

Das Freigeben oder Neuzuweisen von Speicher wirkt sich auch auf Code aus, der ihn ausgeliehen hat. Überprüfen Sie, ob dieser Speicher noch von einem anderen Vorgang verwendet wird. Ein an einen asynchronen Vorgang übergebener Puffer muss gültig bleiben, bis der Vorgang abgeschlossen ist.

## Die Byteliste enthält Buffer

Das Verwalten einer Liste von Bytes mit häufig wechselnden Längen bei manueller Neuzuweisung erfordert die Behandlung von len und cap, Erweiterungsfehlern und Größenberechnungen. Buffer von std bündelt diese Vorgänge.

<!-- wave-example: book-buffer-progressive -->
```wave
import("std::buffer::alloc")::{
    Buffer, buffer_init, buffer_free
};
import("std::buffer::write")::{
    buffer_append_str
};

fun main() -> i32 {
    var message: Buffer;

    if (buffer_init(&message, 0) < 0) {
        return 1;
    }

    if (buffer_append_str(&message, "Hello") < 0) {
        buffer_free(&message);
        return 2;
    }

    if (buffer_append_str(&message, ", Wave") < 0) {
        buffer_free(&message);
        return 3;
    }

    println("bytes={}", message.len);

    if (buffer_free(&message) < 0) {
        return 4;
    }

    return 0;
}
```

Ausführungsergebnis:

```text
bytes=11
```

len ist die Anzahl der verwendeten Bytes; cap ist die zugewiesene Kapazität. Durch das Anhängen von Daten wird die Zuordnung bei Bedarf erhöht. Die Verwendung von Buffer entbindet den Anrufer nicht von der Verantwortung, es freizugeben.

buffer_append_str fügt NUL nicht am Ende der Zeichenfolge hinzu. Daher sollte message.data nicht direkt als str ausgegeben werden. Bytes werden mit der Funktion I/O ausgegeben, die die Länge übernimmt oder explizit eine String-Darstellung erstellt.

## Übung und Lösungsansatz

Addieren Sie die Bytes 0 bis 9 nacheinander zu Buffer und erhalten Sie die Summe. Es muss freigegeben werden, wenn jedes Anhängen fehlschlägt, und Lesevorgänge werden nur im Bereich len durchgeführt. Die vollständige Lösung und der Grenzfehler können anhand des Beispiels in [Buffer Verwendung](/docs/de/stdlib/buffer) überprüft werden.

Versuchen Sie, Zuordnungs-, Neuzuweisungs- und Freigabeaufrufe in Ihrem Code zu markieren. Für jede erfolgreiche Zuweisung müssen Sie beschreiben können, wem sie gehört und über welchen Pfad sie freigegeben wird.
