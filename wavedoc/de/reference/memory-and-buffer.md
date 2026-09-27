---
translation_set_id: memory-buffer
path: reference/memory-and-buffer
locale: de
group: stdlib
group_order: 1
order: 4
title: mem: Zuweisung, Neuzuweisung und Layout
summary: Beschreibt die Größe in Bytes, Zuweisungsfehler, Neuzuweisungsgrenzen und die Freigabe von Verantwortlichkeiten.
---

## Zuteilung und Freigabe

```text
std::mem::alloc
mem_alloc(size: i64) -> ptr<u8>
mem_alloc_zeroed(size: i64) -> ptr<u8>
mem_free(p: ptr<u8>, size: i64) -> i64
mem_realloc(old_ptr: ptr<u8>, old_size: i64, new_size: i64) -> ptr<u8>
```

Die Größe wird in Bytes angegeben. Eine Zuordnung der Größe 0 oder weniger gibt null zurück. Wenn die Zuweisung auch für positive Größen fehlschlägt, liegt möglicherweise null vor. Gehen Sie nicht vom ursprünglichen Inhalt von `mem_alloc` aus, sondern verwenden Sie `mem_alloc_zeroed`, wenn eine Nullinitialisierung erforderlich ist.

Der Aufrufer ist Eigentümer jeder erfolgreichen Zuweisung und muss beim Freigeben die ursprüngliche Größe übergeben. `mem_free(null, size)` gibt 0 zurück. Ein Nicht-null-Zeiger gepaart mit einer nichtpositiven Größe ist ein Fehler. Greifen Sie niemals auf eine Zuteilung zu oder geben Sie diese frei, nachdem sie bereits freigegeben wurde.

## Einstufung bei Umwidmung

|Anfrage|Aktion|
| --- | --- |
|neue Größe ist positiv und Erfolg|`min(old_size, new_size)` Byte kopieren und vorheriges Byte freigeben|
|Neue Zuordnung positiver Größe schlägt fehl|null Zurück, bestehende Zuordnung beibehalten|
|`old_ptr == null`, positive neue Größe|verhält sich wie eine neue Aufgabe|
| `new_size == 0` |Versucht, eine gültige vorherige Zuordnung freizugeben und gibt null zurück|
|old_size=0 für negative Größe oder vorhandenen Zeiger|null Zurück|

Ein null-Ergebnis aus der Neuzuweisung auf Größe Null stellt nicht fest, ob die Freigabe erfolgreich war. Rufen Sie `mem_free` direkt an, wenn Sie den Status benötigen. Nachdem Sie eine Zuteilung erweitert haben, initialisieren Sie die neu hinzugefügte Region selbst.

## Beispiel für die Beibehaltung eines vorhandenen Zeigers

Unten sehen Sie den Fall, in dem die neue Größe positiv ist. Speichern Sie es als `main.wave` und führen Sie es aus.

<!-- wave-example: reallocation -->
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
    var grown: ptr<u8> = mem_realloc(data, 4, 8);
    if (grown == null) {
        mem_free(data, 4);
        return 2;
    }
    data = grown;
    println("{}", deref data[0]);
    if (mem_free(data, 8) < 0) {
        return 3;
    }

    return 0;
}
```

Ausführungsergebnis:

```text
7
```

Wenn Sie es mit `data = mem_realloc(...)` überschreiben, bevor Sie den Fehler bestätigen, geht möglicherweise die vorhandene Adresse verloren. Wenn die Neuzuweisung erfolgreich ist, werden die alte Adresse und die darauf verweisenden Zeiger nicht verwendet.

## Die Größe und Ausrichtung eines Zieltyps

```text
std::mem::layout
size_of<T>() -> u64
align_of<T>() -> u64
```

Bei beiden Werten handelt es sich um das Layout des Kompilierungsziels, nicht um den Computer, auf dem es ausgeführt wird. `size_of` beinhaltet die Schwanzpolsterung und generiert oder wertet keinen Wert aus. Wenn wir die Anzahl der Elemente mit ihrer Größe multiplizieren, prüfen wir, ob ein Überlauf vorliegt. Sie können `mem_size_mul_checked` und `mem_size_add_checked` von `std::mem::ops` verwenden.

`mem_copy` wird zum Kopieren eines Bereichs verwendet, der sich nicht überlappt, und `mem_move` wird zum Kopieren eines Bereichs verwendet, der sich möglicherweise überlappt. Keiner von beiden kann die tatsächliche Zuordnungslänge allein aus Zeigern bestimmen, daher muss der Aufrufer Grenzen garantieren. Mit [Buffer](/docs/de/stdlib/buffer) kann eine Liste von Bytes unterschiedlicher Größe verwaltet werden.

## Beispiel für eine Layoutabfrage

Speichern Sie es als main.wave und führen Sie es aus. Die Größe und Ausrichtung des im Dokument abgedeckten Ziels, i32, beträgt jeweils 4 Bytes, daher wird `4 4` ausgegeben. Überprüfen Sie die Werte anderer Typen, insbesondere Strukturen und Zeiger, nach Ziel.

<!-- wave-example: layout-api -->
```wave
import("std::mem::layout")::{
    size_of, align_of
};

fun main() {
    println("{} {}", size_of<i32>(), align_of<i32>());
}
```

## Überprüfen Sie die Größenberechnung auf einen Überlauf

Sie müssen prüfen, ob `count * element_size` gültig ist, bevor Sie es an die Zuweisungsfunktion übergeben. Wenn Sie dem Überlaufwert einen kleinen Raum zuweisen und so viel wie die ursprüngliche Zahl schreiben, wird der Wert überschritten.

<!-- wave-example: book-memory-size-check -->
```wave
import("std::mem::ops")::{mem_size_add_checked, mem_size_mul_checked};

fun main() {
    var result: i64 = 99;

    if (mem_size_mul_checked(3, 4, &result) == 0) {
        println("bytes={}", result);
    }

    if (mem_size_add_checked(9223372036854775807, 1, &result) < 0) {
        println("overflow rejected");
    }
}
```

Ausführungsergebnis:

```text
bytes=12
overflow rejected
```

Fehlerergebnisse werden nicht in die Zuordnungsgröße geschrieben. Sie können nicht alle Überläufe erkennen, indem Sie sie einfach mit regulärer Arithmetik berechnen und prüfen, ob das Ergebnis negativ ist. Um die zu prüfende Größe zu berechnen, verwenden Sie von Anfang an die Funktion checked.

## Für überlappende Kopien: mem_move

Beim Verschieben eines Teils desselben Arrays nach hinten überlappen sich die Eingabe- und Ausgabebereiche. Übergeben Sie keine überlappenden Bereiche an mem_copy, sondern verwenden Sie mem_move.

<!-- wave-example: book-memory-overlap -->
```wave
import("std::mem::ops")::{mem_move};

fun main() {
    var data: array<u8, 5> = [1, 2, 3, 4, 5];

    mem_move(&data[1], &data[0], 4);

    for (var index: i32 = 0; index < 5; index += 1) {
        println("{}", data[index]);
    }
}
```

Ausführungsergebnis:

```text
1
1
2
3
4
```

Die ursprünglichen ersten vier Bytes verschieben sich um eine Position nach rechts. Durch manuelles Vorwärtskopieren können bereits überschriebene Werte gelesen werden, wodurch versehentlich nur Einsen entstehen. Eine überlappungsfähige API übernimmt die Kopierrichtung für Sie.

## Schreiben Sie eine Funktion, die den Besitz übergibt

Für Funktionen, die Speicher zurückgeben, ist es am besten, bei Erfolg eine Rücksprungadresse und die für die Freigabe erforderliche Größe anzugeben. Wenn der Anrufer die Größe erraten muss, besteht die Gefahr einer falschen Freigabe. Wenn eine Funktion eine geliehene Adresse zurückgibt, sollte diese vom Aufrufer nicht freigegeben werden und beschreibt die Lebensdauer des Originals.

Über Funktionsgrenzen hinweg sollten Sie in der Lage sein, `allocator → owner → deallocation` zu verfolgen. Der Name oder Typ einer Zeigervariablen bestimmt nicht automatisch den Eigentümer.
