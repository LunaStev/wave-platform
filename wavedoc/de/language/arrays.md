---
translation_set_id: learn-arrays-strings
path: language/arrays
locale: de
group: language
group_order: 2
order: 6
title: 6. Arrays und Iteration
summary: Lernen Sie Arrays fester Größe, Indizierung, Traversierung, Kopieren und Suchen.
---

## Mehrere Werte desselben Typs

Wenn Sie die drei Scores separat als score1, score2 und score3 erstellen, müssen sowohl Deklarationen als auch Berechnungen geändert werden, wenn sich die Zahl ändert. Arrays gruppieren eine bestimmte Anzahl von Elementen desselben Typs. Mithilfe von Schleifen können Sie auf jedes Element dieselben Regeln anwenden.

In diesem Kapitel wird das Erstellen, Indizieren, Ändern, Durchlaufen, Durchsuchen und Aggregieren von Arrays behandelt. Strings unterstützen auch die Indizierung, aber ihre Bedeutung ist unterschiedlich, daher werden Strings im nächsten Kapitel behandelt.

## Geben Sie die Länge als Typ ein

<!-- wave-example: book-array-create -->
```wave
fun main() {
    var scores: array<i32, 3> = [70, 80, 90];

    println("first={}", scores[0]);
    println("second={}", scores[1]);
    println("last={}", scores[2]);
}
```

Ausführungsergebnis:

```text
first=70
second=80
last=90
```

i32 in `array<i32, 3>` ist der Elementtyp und 3 ist die Anzahl der Elemente. Was wir speichern, sind 3 ganze Zahlen. Dies bedeutet nicht, dass die Anzahl der Bytes 3 beträgt. Die Anzahl der Elemente in einem Array-Literal muss mit der deklarierten Länge übereinstimmen.

Indizes beginnen bei 0. Das erste Element ist 0, das letzte Element hat die Länge 1. scores[3] ist ein Zugriff außerhalb des Gültigkeitsbereichs, kein drittes Element.

## Element ändern

<!-- wave-example: book-array-update -->
```wave
fun main() {
    var scores: array<i32, 3> = [70, 80, 90];

    scores[1] = 85;
    scores[0] += 5;

    println("{} {} {}", scores[0], scores[1], scores[2]);
}
```

Ausführungsergebnis:

```text
75 85 90
```

Ändern Sie die Speicherung bestimmter Elemente, ohne das gesamte Array neu zu erstellen. Ein Indexausdruck kann auch das Ergebnis einer Berechnung sein, Sie müssen jedoch sicherstellen, dass der Wert innerhalb eines Bereichs liegt. Bei Verwendung einer externen Eingabe als Index werden sowohl die negative Zahl als auch die Obergrenze überprüft.

## Iterieren über ein Array

<!-- wave-example: book-array-sum -->
```wave
fun main() {
    var scores: array<i32, 4> = [60, 70, 80, 90];
    var total: i32 = 0;

    for (var index: i32 = 0; index < 4; index += 1) {
        total += scores[index];
    }

    println("total={} average={}", total, total / 4);
}
```

Ausführungsergebnis:

```text
total=300 average=75
```

Jede Iteration liest ein Element an einem anderen Index. Die Summenvariable muss außerhalb der Iteration initialisiert werden. Die Initialisierung auf 0 innerhalb des Schleifenkörpers führt zu falschen Ergebnissen, z. B. dass nur das letzte Element übrig bleibt.

Bei der ganzzahligen Division zur Berechnung des Durchschnitts wird der Bruchteil verworfen. Für einen Gleitkomma-Durchschnitt konvertieren Sie die Summe vor der Division. Stellen Sie bei größeren Arrays oder größeren Werten außerdem sicher, dass der Akkumulatortyp die Summe darstellen kann.

## Aggregieren Sie nur einige Elemente

Die Filterung kann durch die Kombination von bedingten Anweisungen und Durchquerung erreicht werden. Hier zählen wir die Anzahl der Elemente mit einer Punktzahl von 80 oder höher.

<!-- wave-example: book-array-filter -->
```wave
fun main() {
    var scores: array<i32, 5> = [60, 80, 90, 75, 100];
    var passed: i32 = 0;

    for (var index: i32 = 0; index < 5; index += 1) {
        if (scores[index] >= 80) {
            passed += 1;
        }
    }

    println("passed={}", passed);
}
```

Ausführungsergebnis:

```text
passed=3
```

Index- und Elementwerte müssen getrennt werden. Beim Untersuchen von `index >= 80` werden Positionen verglichen, nicht Punkte. Beide können i32 sein, daher ist es schwierig, diesen semantischen Fehler allein anhand des Typs zu finden.

## Finden Sie den ersten Spielort

Entscheiden Sie zunächst, wie die nicht gefundenen Ergebnisse angezeigt werden sollen. In diesem Beispiel sind gültige Indizes 0 bis 4, daher verwenden wir -1 als Fehlermarkierung.

<!-- wave-example: book-array-find -->
```wave
fun main() {
    var values: array<i32, 5> = [8, 3, 8, 1, 5];
    var target: i32 = 8;
    var found: i32 = -1;

    for (var index: i32 = 0; index < 5; index += 1) {
        if (values[index] == target) {
            found = index;
            break;
        }
    }

    if (found >= 0) {
        println("found at {}", found);
    } else {
        println("not found");
    }
}
```

Ausführungsergebnis:

```text
found at 0
```

Auch die erste Position 0 ist ein normales Ergebnis. Wenn Sie mit `found > 0` den Erfolg prüfen, werden Sie fälschlicherweise das erste Element nicht gefunden. Aus dem gleichen Grund ist es falsch, den Erfolg anhand von bool als cast zu beurteilen.

Wenn Sie break entfernen, überschreiben nachfolgende Treffer found und geben Ihnen die Position des letzten Treffers zurück. Da eine einzelne Anweisung den Vertrag einer Funktion ändern kann, muss in der Beschreibung von „search“ auch ausdrücklich darauf hingewiesen werden, ob es an der ersten oder letzten Position steht.

## Array-Elemente kopieren

Um die Werte eines Arrays auf einen anderen Speicherplatz zu kopieren, können Sie diese Element für Element auslesen und zuweisen. Selbst wenn Sie nach dem Kopieren ein ganzzahliges Element ändern, ändert sich das andere ganzzahlige Element nicht.

<!-- wave-example: book-array-copy -->
```wave
fun main() {
    var original: array<i32, 3> = [1, 2, 3];
    var copied: array<i32, 3>;

    for (var index: i32 = 0; index < 3; index += 1) {
        copied[index] = original[index];
    }

    copied[0] = 99;

    println("original={}", original[0]);
    println("copied={}", copied[0]);
}
```

Ausführungsergebnis:

```text
original=1
copied=99
```

Wenn es sich bei den Elementen um Zeiger handelt, werden beim Kopieren deren Adressen kopiert. Der separate Speicher, auf den sie verweisen, wird nicht dupliziert. Diese Unterscheidung ist bei der Verwaltung des Eigentums wichtig.

## Initialisierung und Wirkbereich

Es wird nicht davon ausgegangen, dass alle Elemente eines Arrays, das ohne Anfangswert deklariert wurde, gelesen werden können. Wird nur eine Teilnummer erfasst, muss die tatsächlich initialisierte Nummer separat verwaltet werden. Aus demselben Grund kann die von der Bibliothekslesefunktion zurückgegebene Länge kleiner sein als die gesamte Pufferkapazität.

Da die Array-Länge im Typ enthalten ist, wächst sie während der Ausführung nicht willkürlich an. Listen von Bytes, deren Größe zunimmt, verwenden dynamischen Speicher wie `Buffer`. Das Ändern der Länge eines Arrays erfordert die Berücksichtigung des Typs, des Anfangswerts, der Obergrenze des Durchlaufs und der Berechnungen, die von dieser Länge abhängen.

## Übung: Maxima und Orte

Finden Sie den Maximalwert und die Position, an der er zum ersten Mal im Array `[4, 9, 2, 9, 1]` erscheint. Wenn es sich um eine reguläre Funktion handelt, bei der alle Elemente negativ sein können, sollte der Maximalwert nicht auf 0 initialisiert werden.

### Vollständige Lösung

<!-- wave-example: book-array-solution -->
```wave
fun main() {
    var values: array<i32, 5> = [4, 9, 2, 9, 1];
    var maximum: i32 = values[0];
    var position: i32 = 0;

    for (var index: i32 = 1; index < 5; index += 1) {
        if (values[index] > maximum) {
            maximum = values[index];
            position = index;
        }
    }

    println("max={} first={}", maximum, position);
}
```

Ausführungsergebnis:

```text
max=9 first=1
```

Nehmen Sie das erste Element als Ausgangsreferenz und vergleichen Sie es mit dem zweiten. Seit `>` ändert sich die Position nicht, auch wenn der gleiche Maximalwert erneut auftritt. Ändern Sie es in `>=`, um die letzte Position zu übernehmen. Jede Schnittstelle, die die Länge Null haben kann, muss eine leere Eingabe verarbeiten, bevor das erste Element gelesen wird.
