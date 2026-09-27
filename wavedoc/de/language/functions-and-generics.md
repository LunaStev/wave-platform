---
translation_set_id: functions
path: language/functions-and-generics
locale: de
group: language
group_order: 2
order: 5
title: 5. Entwerfen und Verfassen von Funktionen
summary: Lernen Sie Parameter, Rückgabewerte, Standardwerte und die Übergabe von Werten.
---

## Ausgehend von repetitivem Code

Funktionen sind ein Werkzeug zur Reduzierung der Syntax, aber auch ein Werkzeug zur Abgrenzung von Aufgaben. Durch die Trennung dessen, was als Eingabe verwendet wird, was es berechnet und welche Ergebnisse es zurückgibt, können Sie Ihr Programm in kleineren Teilen verstehen.

In diesem Kapitel beginnen wir mit einem Programm, das Rabattberechnungen mehrmals schreibt. Jedes Beispiel ist in main.wave vollständig und wird als `wavec run main.wave` ausgeführt. Wir teilen die Dateien noch nicht auf.

<!-- wave-example: book-function-before -->
```wave
fun main() {
    var first_price: i32 = 2000;
    var first_discount: i32 = first_price * 10 / 100;
    var first_total: i32 = first_price - first_discount;
    var second_price: i32 = 5000;
    var second_discount: i32 = second_price * 10 / 100;
    var second_total: i32 = second_price - second_discount;
    println("{} {}", first_total, second_total);
}
```

Ausführungsergebnis:

```text
1800 4500
```

Die beiden Kalkulationen unterscheiden sich lediglich im Preis und sind gleich aufgebaut. Wenn Sie Rabattregeln ändern, müssen Sie beide Stellen bearbeiten. Wenn Sie nur eine Seite ändern, erhalten Sie unterschiedliche Ergebnisse für Produkte, für die dieselbe Richtlinie erforderlich ist.

## Bestimmen Sie Input und Output

Verschieben Sie redundante Berechnungen in Funktionen. Der geänderte Wert wird als Eingabe price empfangen und der berechnete Preis zurückgegeben.

<!-- wave-example: book-function-extract -->
```wave
fun discounted(price: i32) -> i32 {
    var discount: i32 = price * 10 / 100;
    return price - discount;
}

fun main() {
    println("{} {}", discounted(2000), discounted(5000));
}
```

Ausführungsergebnis:

```text
1800 4500
```

Der Funktionsname ist discounted. `price: i32` in Klammern ist die Parameterdeklaration und `-> i32` ist der Ergebnistyp. Die lokale Variable discount im Rumpf wird nur innerhalb dieser Funktion verwendet.

`discounted(2000)` ist ein Ausdruck, der eine Funktion aufruft. Die 2000 in Klammern ist das Argument, das Sie tatsächlich übergeben. Da der von der Funktion zurückgegebene Wert das Ergebnis dieses aufrufenden Ausdrucks wird, kann er direkt als Argument für println verwendet werden.

|Terminologie|Code|Bedeutung|
| --- | --- | --- |
|Parameter| price |Eingabename, der bei der Deklaration der Funktion angegeben wurde|
|Faktor| 2000 |Beim Aufruf übergebener Wert|
|Rückgabetyp| i32 |Typ des vom Aufrufausdruck erzeugten Werts|
|return-Anweisung| return price - discount |Übergeben Sie das Ergebnis und beenden Sie diesen Anruf|

## Befolgen Sie die Aufruf- und Ausführungsreihenfolge

Das bloße Schreiben einer Funktionsdeklaration in die Quelle führt nicht dazu, dass ihr Hauptteil sofort ausgeführt wird. Wird ausgeführt, wenn der in main aufgerufene Punkt erreicht ist.

<!-- wave-example: book-function-trace -->
```wave
fun calculate(value: i32) -> i32 {
    println("inside: {}", value);
    return value * 2;
}

fun main() {
    println("before");
    var result: i32 = calculate(7);
    println("after: {}", result);
}
```

Ausführungsergebnis:

```text
before
inside: 7
after: 14
```

Die Reihenfolge des Fortschritts ist die erste Ausgabe von main, der Haupttext von calculate und die letzte Ausgabe von main. Wenn `return` ausgeführt wird, endet der Aufruf von calculate mit Ergebnis 14 und die Initialisierung von result von main ist abgeschlossen.

Wenn mehrere Funktionsaufrufe in einem Ausdruck gemischt sind und die Reihenfolge der Nebenwirkungen wichtig ist, teilen Sie die Aufrufe in separate Anweisungen auf. Die Beispiele in diesem Kapitel speichern auch die Ergebnisse von Aufrufen, die verfolgt werden müssen, in lokalen Variablen.

## mehrere Parameter

Wenn Sie zusätzlich einen Rabattsatz als Eingabe erhalten, können Sie mehrere Policen mit derselben Funktion berechnen.

<!-- wave-example: book-function-parameters -->
```wave
fun discounted(price: i32, percent: i32) -> i32 {
    return price - price * percent / 100;
}

fun main() {
    var standard: i32 = discounted(2000, 10);
    var special: i32 = discounted(2000, 25);
    println("standard={}", standard);
    println("special={}", special);
}
```

Ausführungsergebnis:

```text
standard=1800
special=1500
```

Die Reihenfolge der Argumente muss mit der Deklaration übereinstimmen. Wenn beide Parameter i32 sind, ist es schwierig, ihre Bedeutung allein durch Typprüfung zu unterscheiden, selbst wenn die Reihenfolge geändert wird. Definieren Sie den Funktionsnamen und den Parameternamen klar und schreiben Sie den Aufrufort gut lesbar.

Diese Funktion geht von einem kleinen Betrag und einem gültigen Rabattsatz aus. Behandelt keine negativen Preise, Verhältnisse größer als 100 oder einen Überlauf in der Mitte der Multiplikation. Beim Erstellen einer Funktion müssen Sie nicht nur den Körper, sondern auch die Eingabebedingungen beschreiben. Das spätere Abschlussprogramm trennt die Prüfschritte.

## Standardargument

Sie können häufig verwendete Werte als Standardwerte angeben.

<!-- wave-example: book-function-default -->
```wave
fun discounted(price: i32, percent: i32 = 10) -> i32 {
    return price - price * percent / 100;
}

fun main() {
    println("{}", discounted(2000));
    println("{}", discounted(2000, 25));
}
```

Ausführungsergebnis:

```text
1800
1500
```

Der erste Aufruf lässt das zweite Argument weg und verwendet 10. Der zweite Aufruf verwendet die angegebenen 25. Für optionale Parameter werden am Ende Standardwerte belassen. Als Grammatik wird es nicht verwendet, ein Leerzeichen zu schreiben, um nur das erste Argument wegzulassen.

Durch Ändern der Standardeinstellung wird das Verhalten von Skip-Aufrufen geändert. Die Standardwerte öffentlicher Funktionen sind ebenfalls Teil des Verhaltens, auf das Sie sich verlassen. Aus diesem Grund werden Aufrufe mit angegebenen Argumenten und Aufrufe mit weggelassenen Argumenten separat getestet.

## bedeutet Wertübergabe

Durch die Übergabe eines ganzzahligen Werts wird der von der Funktion empfangene Wert vom Variablenspeicher des Aufrufers getrennt. Durch die Berechnung eines Ergebnisses werden die Variablen des Aufrufers nicht automatisch geändert.

<!-- wave-example: book-function-value -->
```wave
fun next(value: i32) -> i32 {
    return value + 1;
}

fun main() {
    var count: i32 = 4;
    var later: i32 = next(count);
    println("count={} later={}", count, later);
    count = next(count);
    println("count={}", count);
}
```

Ausführungsergebnis:

```text
count=4 later=5
count=5
```

Der erste Aufruf liest count und initialisiert later. count ist immer noch 4. Nach dem zweiten Aufruf wird das Ergebnis count zugewiesen, sodass es zu 5 wird. Das Return-by-Value-Design macht es für den Aufrufer offensichtlich, wo sich die Daten geändert haben.

Wenn Sie den ursprünglichen Speicherplatz innerhalb einer Funktion ändern möchten, können Sie einen Zeiger übergeben. Dies wird in [Zeigerkapitel](/docs/de/language/explicit-memory-type-model) behandelt. Auch wenn Sie einen Zeiger übergeben, müssen Sie zwischen dem Zeigerwert selbst und dem Speicherplatz für seine Adresse unterscheiden.

## Verlasse keine Wege, die nicht zurückkehren

Eine Funktion, die einen Wert zurückgibt, muss auf allen erforderlichen Pfaden Ergebnisse liefern. Es verlässt den endgültigen Pfad nicht so:

```wave
// 의도적으로 잘못된 함수 조각

fun positive(value: i32) -> i32 {
    if (value > 0) {
        return value;
    }
}
```

Wenn value kleiner oder gleich 0 ist, gibt es keinen zurückzugebenden eingestellten Wert. Sobald Sie Ihre beabsichtigten Regeln festgelegt haben, müssen Sie alle Ihre Routen aufschreiben.

<!-- wave-example: book-function-paths -->
```wave
fun positive_or_zero(value: i32) -> i32 {
    if (value > 0) {
        return value;
    }

    return 0;
}

fun main() {
    println("{} {} {}", positive_or_zero(-2), positive_or_zero(0), positive_or_zero(8));
}
```

Ausführungsergebnis:

```text
0 0 8
```

Der Aufruf, der zuerst return ausführt, fährt nicht mit return weiter unten fort. Der letzte return wird nur erreicht, wenn die Bedingung falsch ist. Wir haben die Regel für drei Fälle überprüft: positiv, 0 und negativ.

## Funktion ohne Ergebnis

Wenn Sie nur Vorgänge wie Drucken ausführen, können Sie den Rückgabetyp weglassen. Eine Funktion ohne Ergebnis kann auch mit `return;` vorzeitig beendet werden.

<!-- wave-example: book-function-void -->
```wave
fun show_positive(value: i32) {
    if (value <= 0) {
        return;
    }

    println("positive={}", value);
}

fun main() {
    show_positive(-3);
    show_positive(6);
}
```

Ausführungsergebnis:

```text
positive=6
```

Der erste Anruf kehrt zurück, ohne etwas zu drucken. Beim zweiten Anruf wird Folgendes ausgegeben: „Kein Ergebnis“ und „Kehrt nicht zum Anrufpunkt zurück“ sind unterschiedlich. Funktionen, die nicht zurückkehren, wie z. B. Funktionen, die einen Prozess beenden, werden nach ihrem Rückgabetyp `!` klassifiziert.

## Erstellen eines vollständigen Programms mit mehreren Funktionen

Jetzt teilen wir die Eingabevalidierung, Berechnung und Ausgabe in verschiedene Funktionen auf. Da wir die Preisspanne eingeschränkt haben, liegt die Zwischenmultiplikation in diesem Beispiel im Bereich i32.

<!-- wave-example: book-function-program -->
```wave
fun valid_order(price: i32, quantity: i32, percent: i32) -> bool {
    return price >= 0 && price <= 100000
        && quantity >= 1 && quantity <= 100
        && percent >= 0 && percent <= 100;
}

fun discounted_unit(price: i32, percent: i32) -> i32 {
    return price - price * percent / 100;
}

fun order_total(price: i32, quantity: i32, percent: i32) -> i32 {
    var unit: i32 = discounted_unit(price, percent);
    return unit * quantity;
}

fun show_order(price: i32, quantity: i32, percent: i32) {
    if (!valid_order(price, quantity, percent)) {
        println("invalid order");
        return;
    }

    var total: i32 = order_total(price, quantity, percent);
    println("total={}", total);
}

fun main() {
    show_order(2000, 3, 10);
    show_order(2000, 0, 10);
    show_order(2000, 3, 120);
}
```

Ausführungsergebnis:

```text
total=5400
invalid order
invalid order
```

Beantworten Sie für jede Funktion eine Frage. valid_order legt fest, ob eine Eingabe zulässig ist, discounted_unit bestimmt, wie hoch ein Rabatt ist, order_total legt fest, wie hoch der Gesamtpreis ist, und show_order legt fest, was angezeigt werden soll.

Kleinere Funktionen sind nicht unbedingt besser. Wenn Sie jedem Ausdruck einen Namen geben, kann es schwierig sein, den Überblick zu behalten. Separate Regeln, wenn es sinnvoll ist, sie an anderer Stelle wiederzuverwenden, oder wenn es Regeln gibt, die unabhängig erklärt und überprüft werden müssen.

## Übungsprobleme

1. Schreiben Sie maximum, das die größere von zwei ganzen Zahlen zurückgibt.
2. Schreiben Sie eine Funktion clamp, die einen Wert zwischen zwei Grenzen zurückgibt. In der folgenden Lösung ist low <= high als Aufrufbedingung festgelegt.
3. Erstellen Sie eine Funktion, die einem Wert Steuern hinzufügt, und kombinieren Sie sie mit der Bestellsummenfunktion. Bestimmen Sie zunächst den Bereich und die ganzzahligen Schnittpunkte.

### Lösung: Funktion mit Rand

<!-- wave-example: book-function-solution -->
```wave
fun maximum(left: i32, right: i32) -> i32 {
    if (left > right) {
        return left;
    }

    return right;
}

fun clamp(value: i32, low: i32, high: i32) -> i32 {
    if (value < low) {
        return low;
    }
    if (value > high) {
        return high;
    }

    return value;
}

fun main() {
    println("max={}", maximum(7, 4));
    println("{} {} {}", clamp(-3, 0, 10), clamp(6, 0, 10), clamp(20, 0, 10));
}
```

Ausführungsergebnis:

```text
max=7
0 6 10
```

clamp hat drei Pfade: kleiner als Reichweite, innerhalb Reichweite und größer als Reichweite. Überprüfen Sie dies, indem Sie auch die Grenzwerte 0 und 10 manuell hinzufügen. Wenn Sie sich bis low > high bewerben möchten, müssen Sie entscheiden, wie Sie das Scheitern äußern. [Kapitel zur Fehlerbehandlung](/docs/de/language/errors) behebt dieses Problem mithilfe von Ergebnisstrukturen.


## rekursiver Aufruf

Eine Funktion kann sich selbst aufrufen. Die Rekursion erfordert eine Beendigungsbedingung, bei der keine weiteren Aufrufe erfolgen, und einen Prozess, bei dem jeder Aufruf dieser Bedingung näher kommt.

<!-- wave-example: book-ref-recursion -->
```wave
fun factorial(value: i32) -> i32 {
    if (value <= 1) {
        return 1;
    }

    return value * factorial(value - 1);
}

fun main() {
    println("{}", factorial(5));
}
```

Ausführungsergebnis:

```text
120
```

Die Berechnung von 5 führt zu `5 * factorial(4)`, das bei Erreichen 1 zurückgibt. Der Rückgabewert wird sequentiell zum vorherigen Aufruf übergeben, was zu 120 führt. Diese Funktion ist ein Beispiel zur Erläuterung kleiner positiver Ganzzahlen. Bei großen Eingaben müssen Sie den Ergebnisbereich und die Aufruftiefe berücksichtigen. Sie können das Problem einer erhöhten Aufruftiefe vermeiden, indem Sie denselben Vorgang in eine Schleife schreiben.

## Häufige Fehler

|Phänomen|Überprüfen|
| --- | --- |
|Fehler mit unzureichenden oder zu vielen Argumenten|Anzahl der Parameter und optionale Standardwerte|
|Der Rückgabetyp stimmt nicht überein|return Ausdruckstyp und Funktionsdeklaration|
|Nicht von einem bestimmten Weg zurückkehren|Wird es zurückgegeben, auch wenn die Bedingung falsch ist?|
|Fehlendes generisches Typargument|Nach dem Funktionsnamen `<Type>`|
|Das Original wird nach dem Aufruf einer Zeigerfunktion geändert.|Handelt es sich um eine Funktion, die den Wert nur liest, oder um eine Funktion, die ihn ändert?|

Extern exportierte Funktionen wie `export(c)` verwenden bestimmte Signaturen. Generische Funktionen selbst können nicht mit einer externen Aufrufkonvention exportiert werden. `ptr<T>` und `array<T, N>` sind die integrierten Speichertypen der Sprache und unterscheiden sie von generischen Strukturdeklarationen des Benutzers.

[Funktionslernen](/docs/de/language/functions-and-generics) · [Lernmodule und Generika](/docs/de/language/modules-imports-and-ffi) · [FFI](/docs/de/language/modules-imports-and-ffi)
