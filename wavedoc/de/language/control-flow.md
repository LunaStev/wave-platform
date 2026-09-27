---
translation_set_id: control-flow
path: language/control-flow
locale: de
group: language
group_order: 2
order: 4
title: 4. Bedingungen, Schleifen und Grenzwerte
summary: Lernen Sie if, for, while und die Bandbreite sich wiederholender Aussagen.
---

## Wählen Sie den Ausführungspfad aus

Das Programm im vorherigen Kapitel führte die Anweisungen von oben nach unten aus. Ein echtes Programm muss je nach Eingabe und Zustand unterschiedliche Dinge tun. Bedingte Anweisungen wählen einen Ausführungspfad aus und Schleifen wenden dieselbe Regel auf mehrere Werte an.

Speichern Sie jedes Beispiel in main.wave und führen Sie es aus. Notieren Sie beim Lesen des Codes die aktuellen Variablenwerte, die als nächstes zu testenden Bedingungen und die Reihenfolge der auszuführenden Anweisungen auf einem Blatt Papier. Es ist wichtiger zu üben, dem Ablauf zu folgen, als sich die Ergebnisse zu merken.

## if und else

Bedingungen werden in Klammern geschrieben und der Text wird in geschweifte Klammern eingeschlossen. Ersetzen Sie im folgenden Beispiel balance durch 500 oder 2000, um zu bestimmen, welcher Zweig ausgeführt wird.

<!-- wave-example: book-if-else -->
```wave
fun main() {
    var balance: i32 = 2000;
    var price: i32 = 1200;

    if (balance >= price) {
        balance -= price;
        println("bought, balance={}", balance);
    } else {
        println("not enough money");
    }
}
```

Ausführungsergebnis:

```text
bought, balance=800
```

Es werden nicht beide Blöcke ausgeführt. Wenn die Bedingung wahr ist, wird der erste Block ausgeführt. Wenn die Bedingung falsch ist, wird der Block else ausgeführt. Käufe sind auch dann erlaubt, wenn balance gleich price ist, also habe ich `>=` geschrieben. Wenn Sie es in `>` ändern, ändert sich der Vorgang um den gleichen Betrag.

## Reihenfolge mehrerer Bedingungen

Sie können Bedingungen mit else if verketten. Da wir zuerst einen erfüllten Zweig von oben ausführen, müssen wir überlegen, ob wir zuerst die größere Grenze oder die kleinere Grenze überprüfen möchten.

<!-- wave-example: book-grade -->
```wave
fun grade(score: i32) {
    if (score < 0 || score > 100) {
        println("invalid");
    } else if (score >= 90) {
        println("A");
    } else if (score >= 80) {
        println("B");
    } else {
        println("C");
    }
}

fun main() {
    grade(95);
    grade(80);
    grade(79);
    grade(101);
}
```

Ausführungsergebnis:

```text
A
B
C
invalid
```

Ungültige Noten werden zunächst abgelehnt und anschließend bewertet. Wenn Sie `score >= 80` am Anfang setzen, werden Sie nie den Zweig A erreichen, da 95 Grad in diesen Zweig hineingehen. Überprüfen Sie die Reihenfolge der Bedingungen sowie die Richtigkeit jeder einzelnen Bedingung.

## Ändern Sie keine Werte in bedingten Ausdrücken

Zuweisungen, zusammengesetzte Zuweisungen sowie Inkrementierungs- oder Dekrementierungsoperationen sind in if-, while- oder for-Bedingungen nicht zulässig. Verwenden Sie zum Vergleich `==`. Um einen Wert zu aktualisieren und ihn dann zu testen, schreiben Sie zwei separate Anweisungen.

Richtige Form eines Fragments innerhalb einer Funktion:

```wave
value = read_value();

if (value == expected) {
    println("matched");
}
```

read_value und expected sind in diesem Fragment nicht definiert, daher ist es nicht das vollständige Programm, das so ausgeführt wird, wie es ist. Die hier gezeigte Regel lautet „Vergleichen nach Statusänderung“.

## while: Solange die Bedingungen gelten

<!-- wave-example: book-while-countdown -->
```wave
fun main() {
    var remaining: i32 = 3;

    while (remaining > 0) {
        println("{}", remaining);
        remaining -= 1;
    }

    println("finished at {}", remaining);
}
```

Ausführungsergebnis:

```text
3
2
1
finished at 0
```

Vor dem Betreten des Körpers werden die Bedingungen überprüft. Wenn remaining von Anfang an 0 ist, wird der Körper nie ausgeführt. Wenn die Dekrementierung am Ende des Körpers weggelassen wird, bleibt die Bedingung wahr und die Schleife wird nicht beendet.

Überprüfen Sie nach dem Schreiben der Schleife „Was bringt sie näher an die Beendigungsbedingung?“ Wenn es sich um eine Schleife handelt, die auf Eingaben wartet, spielen Eingabeänderungen oder EOF ihre Rolle, und wenn es sich um eine numerische Schleife handelt, spielt die Indexaktualisierung ihre Rolle.

## for: Initialisierung/Bedingungen/Aktualisierung

for drückt die drei für die Wiederholung notwendigen Teile aus.

<!-- wave-example: book-for-sum -->
```wave
fun main() {
    var total: i32 = 0;

    for (var number: i32 = 1; number <= 5; number += 1) {
        total += number;
    }

    println("sum={}", total);
}
```

Ausführungsergebnis:

```text
sum=15
```

1. Initialisieren Sie number auf 1. Dieser Schritt ist einmalig.
2. Prüft number <= 5. Wenn falsch, beenden Sie die Iteration.
3. Fügen Sie number zu total im Text hinzu.
4. Erhöhen Sie number um 1 und kehren Sie zur Zustandsprüfung zurück.

Es wird nicht davon ausgegangen, dass in for deklarierte Wiederholungsvariablen nach der Wiederholung verwendet werden können. Wenn Ihr Design nach der Iteration einen Wert erfordert, deklarieren Sie ihn außerhalb und machen Sie den Initialisierungsort klar.

## Einschluss- und Ausschlussgrenzen

Die natürliche Summe von 1 bis 5 ist `<= 5`. Andererseits sollte der Index eines Arrays der Länge 5 `< 5` verwenden. Dies liegt daran, dass Array-Indizes bei 0 beginnen und bei 4 enden.

Verwechseln Sie „fünfmal ausführen“ nicht mit „bis einschließlich dem Wert 5“. Die Anzahl der Wiederholungen können Sie ermitteln, indem Sie die Start- und Endwerte gemeinsam betrachten. Fälle, in denen die Eingabe leer ist und nur ein Element enthält, eignen sich gut zum Erkennen von Grenzfehlern.

## continue und break

continue überspringt den Rest dieser Iteration und break beendet die nächste Iteration.

<!-- wave-example: book-loop-control -->
```wave
fun main() {
    var total: i32 = 0;

    for (var number: i32 = 1; number <= 10; number += 1) {
        if (number == 3) {
            continue;
        }

        if (number == 6) {
            break;
        }

        total += number;
    }

    println("{}", total);
}
```

Ausführungsergebnis:

```text
12
```

Die Gesamtzahl beträgt 1, 2, 4 und 5. Wenn continue in for gefunden wird, wird der Erneuerungsprozess fortgesetzt. Da while keinen separaten Aktualisierungsausdruck wie for hat, müssen Sie darauf achten, keine vor continue erforderlichen Statusänderungen auszulassen.

Wenn Schleifen verschachtelt sind, schließt ein break nicht alle Schleifen ab. Wenn Sie in mehreren Phasen anhalten müssen, schließen Sie die Arbeit in eine Funktion ein und geben Sie Ihre Absicht an, indem Sie return verwenden oder auch die Beendigungsbedingung in der äußeren Iteration überprüfen.

## Teilen Sie den Fall durch match

Sie können match verwenden, wenn Sie mehrere Fälle mit demselben Wert vergleichen. Der Körper jedes arm ist ein Block.

<!-- wave-example: book-match-number -->
```wave
fun describe(status: i32) {
    match (status) {
        200 => {
            println("ok");
        }
        404 => {
            println("missing");
        }
        _ => {
            println("other");
        }
    }
}

fun main() {
    describe(200);
    describe(404);
    describe(500);
}
```

Ausführungsergebnis:

```text
ok
missing
other
```

`_` ist das Muster, das den Rest erledigt. Platzieren Sie keine Duplikate innerhalb desselben match. variant, das je nach Wert unterschiedliche Datentypen hat, wird in [Kapitel Datenmodell](/docs/de/language/structures-enums-and-aliases) behandelt.

## Vollständiges Beispiel: Zahlen zählen, die Bedingungen erfüllen

Finden Sie die Anzahl und Summe der geraden Zahlen von 1 bis 10. Da es sich bei Anzahl und Summe um unterschiedliche Informationen handelt, werden sie als Variablen akkumuliert.

<!-- wave-example: book-loop-statistics -->
```wave
fun main() {
    var count: i32 = 0;
    var total: i32 = 0;

    for (var number: i32 = 1; number <= 10; number += 1) {
        if (number % 2 == 0) {
            count += 1;
            total += number;
        }
    }

    println("count={} total={}", count, total);
}
```

Ausführungsergebnis:

```text
count=5 total=30
```

Die geraden Zahlen sind 2, 4, 6, 8 und 10, also ist die Zahl 5 und die Summe ist 30. Auch wenn der Ausdruck kurz ist, lässt sich die Iterationsgrenze leicht überprüfen, wenn Sie das Ergebnis zunächst in einem kleinen Bereich überprüfen, der manuell ermittelt werden kann.

## Übung und vollständige Lösung

Addieren Sie nur Vielfache von 3 von 1 bis 20, aber addieren Sie keine Werte, die mehr als 30 ergeben. Wir müssen zwischen „addieren und dann prüfen, ob es vorbei ist“ und „prüfen, ob es vorbei ist und dann addieren“ unterscheiden.

<!-- wave-example: book-loop-solution -->
```wave
fun main() {
    var total: i32 = 0;

    for (var number: i32 = 1; number <= 20; number += 1) {
        if (number % 3 != 0) {
            continue;
        }

        if (total + number > 30) {
            break;
        }

        total += number;
    }

    println("{}", total);
}
```

Ausführungsergebnis:

```text
30
```

Die Werte 3+6+9+12 ergeben in der Summe 30, sodass der nächste Wert, 15, nicht addiert wird. Diese kleinen Eingaben sind sicher, aber bei großen ganzen Zahlen kann die Prüfung `total + number` selbst überlaufen. Beim Schreiben eines Schecks wird nicht automatisch jeder Grenzfall behandelt.
