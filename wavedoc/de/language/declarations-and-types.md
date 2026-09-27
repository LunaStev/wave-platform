---
translation_set_id: types
path: language/declarations-and-types
locale: de
group: language
group_order: 2
order: 2
title: 2. Variablen, Typen und Umfang
summary: Lernen Sie lokale Variablen, den Ganzzahlbereich, bool und Bereiche kennen.
---

## Werte namentlich behandeln

Wenn Sie Ihre Preise direkt an mehreren Stellen angeben, müssen Sie beim Ändern des Preises alle finden. Eine Variable ist ein Speicherplatz, der Werten Namen gibt und es Ihnen ermöglicht, sie mithilfe dieser Namen zu lesen und zu ändern. In diesem Kapitel erfahren Sie mehr über den Umfang von Deklarationen, Zuweisungen, Typen und den Umfang sichtbarer Namen innerhalb von Blöcken.

Die folgenden Programme sind jeweils Teil eines separaten main.wave. Speichern Sie ein Beispiel, führen Sie es als `wavec run main.wave` aus und ersetzen Sie es durch das nächste Beispiel.

## Deklaration und Initialisierung

<!-- wave-example: book-variable-first -->
```wave
fun main() {
    var price: i32 = 1200;
    var quantity: i32 = 3;
    println("price={}", price);
    println("quantity={}", quantity);
    println("total={}", price * quantity);
}
```

Ausführungsergebnis:

```text
price=1200
quantity=3
total=3600
```

Lesen Sie die Erklärung in vier Teilen.

|Teil|dieses Beispiel|Rolle|
| --- | --- | --- |
|Deklarationsschlüsselwort| var |Lokale Variable erstellen|
|Namen| price |Bezeichner, der später verwendet werden soll|
|Typ| i32 |Typ und Bereich der zu speichernden Werte|
|Anfangswert| 1200 |Erster zu speichernder Wert|

Der Doppelpunkt vor dem Typ und das Gleichheitszeichen vor dem Anfangswert haben unterschiedliche Rollen. Geben Sie Ihrem Namen eine Bedeutung. In diesem Beispiel ist price der Stückpreis und quantity die Menge. Selbst wenn dasselbe i32 synonym verwendet wird, können falsche Berechnungen ohne grammatikalische Fehler durchgeführt werden.

## Zuweisung ist keine Formel, die Beziehungen aufrechterhält.

Wenn Sie das Berechnungsergebnis in einer Variablen speichern, wird der Wert an dieser Stelle eingegeben. Berechnungen werden nicht gespeichert und später automatisch neu ausgewertet.

<!-- wave-example: book-assignment-snapshot -->
```wave
fun main() {
    var price: i32 = 1200;
    var quantity: i32 = 3;
    var total: i32 = price * quantity;
    quantity = 5;
    println("before recalculation={}", total);
    total = price * quantity;
    println("after recalculation={}", total);
}
```

Ausführungsergebnis:

```text
before recalculation=3600
after recalculation=6000
```

`quantity = 5` schreibt einen neuen Wert in eine bereits vorhandene Variable. Es muss von einer Neudeklaration wie `var quantity` unterschieden werden. `total` ist ebenfalls 3600, bevor es erneut ersetzt wird. Wenn ein Programm Beziehungen zwischen mehreren Variablen aufrechterhalten muss, sollte es so geschrieben sein, dass es Berechnungen durchführt, wenn sich die Beziehungen ändern.

## Berechnen Sie den nächsten Wert aus dem vorherigen Wert

Berechnen Sie zuerst die rechte Seite der Zuweisungsanweisung und schreiben Sie das Ergebnis in den Speicherplatz auf der linken Seite. Im Gegensatz zu Gleichungen in der Mathematik ist `count = count + 1` ein gültiges Update.

<!-- wave-example: book-update-variable -->
```wave
fun main() {
    var count: i32 = 0;
    count = count + 1;
    count += 2;
    count *= 3;
    println("{}", count);
}
```

Ausführungsergebnis:

```text
9
```

Der Werteverlauf ist 0 → 1 → 3 → 9. `+=`, `*=` drücken Berechnung und Speicherung gemeinsam aus. Wenn Sie die Berechnungsreihenfolge aufschreiben, können Sie herausfinden, zu welchem ​​Zeitpunkt Ihre Gedanken anders waren, als die Ergebnisse anders ausfielen als erwartet.

## Breite und Vorzeichen von Ganzzahltypen

Wenn es mit `i` beginnt, wird es zu signed, und wenn es mit `u` beginnt, wird es zu unsigned. Die letzte Zahl ist die Anzahl der Bits. Mit zunehmender Bitzahl erhöht sich auch der darstellbare Bereich und damit auch der Speicherplatz.

|Typ|Mindestwert|Maximalwert|Beispielanwendung|
| --- | --- | --- | --- |
| i8 | -128 | 127 |kleiner vorzeichenbehafteter Wert|
| u8 | 0 | 255 |ein Byte|
| i16 | -32768 | 32767 |kleine ganzzahlige Daten|
| u16 | 0 | 65535 |Port/16-Bit-Feld|
| i32 | -2147483648 | 2147483647 |Gängige kleine Ganzzahlberechnungen|
| u32 | 0 | 4294967295 |32-Bit-Bitfeld|

Wave stellt auch vorzeichenbehaftete und vorzeichenlose Ganzzahlen mit 64, 128, 256, 512 und 1024 Bit bereit. Ein breiterer Typ macht nicht jede Berechnung sicher: Das Ergebnis kann immer noch den ausgewählten Bereich überschreiten. Wählen Sie zunächst den Bereich aus, den Sie benötigen. `isz` und `usz` folgen der Zieladressbreite.

Es muss zwischen dem Speichern eines großen Literals in einem kleinen Typ und dem absichtlichen Verwerfen von Bits durch Ausführen von cast unterschieden werden. Die Konvertierung in einen kleineren Typ, nur um Fehler zu beseitigen, kann den Wert selbst ändern. Transformationen werden im nächsten Kapitel behandelt.

## Gleitkommatypen und Bool

`f32`·`f64` sind Gleitkommazahlen. Im Gegensatz zu ganzen Zahlen können sie Dezimalteile darstellen, aber nicht alle Dezimalzahlen exakt speichern. Dies ist ein Grund, warum Beträge in kleinen ganzzahligen Einheiten verwaltet werden.

bool steht für wahr und falsch. Sie können Ihrer Bedingung einen Namen geben, indem Sie die Vergleichsergebnisse wie folgt speichern:

<!-- wave-example: book-named-condition -->
```wave
fun main() {
    var balance: i32 = 5000;
    var cost: i32 = 3600;
    var can_buy: bool = balance >= cost;
    if (can_buy) {
        println("purchase allowed");
    } else {
        println("not enough money");
    }
}
```

Ausführungsergebnis:

```text
purchase allowed
```

`can_buy` folgt nicht automatisch den Änderungen in balance und cost. Nachdem die beiden Werte ausgetauscht wurden, wird der Vergleich erneut durchgeführt, wenn der aktuelle Status benötigt wird.

## nicht initialisierter Speicherplatz

`var value: i32;` ist ein Formular, das nur den Speicherplatz deklariert. Vor dem Lesen muss ein gültiger Wert geschrieben werden. Gehen Sie nicht davon aus, dass 0 automatisch eingefügt wird, nur weil Sie es deklarieren. In einem Einführungskurs ist es einfacher zu verstehen, ob der Wert sofort bekannt ist und gleichzeitig mit der Deklaration initialisiert wird.

Bei Bibliotheksaufrufen, die Werte als Ausgabeargumente erhalten, gibt es Fälle, in denen zuerst ein Leerzeichen deklariert und dann bei Erfolg gelesen wird. Zu diesem Zeitpunkt sollten Sie das Erfolgsergebnis der Funktion überprüfen. Vermeiden Sie den Fehler, nach einem fehlgeschlagenen Aufruf nicht initialisierte Ausgaben zu lesen.

## Gültiger Bereich von Blöcken und Namen

Ein Block ist ein in geschweifte Klammern eingeschlossener Codebereich. Wenn Sie eine neue Variable mit demselben Namen darin deklarieren, wird die neue Variable in diesem Block verwendet. Dies wird als shadowing bezeichnet.

<!-- wave-example: book-shadowing -->
```wave
fun main() {
    var value: i32 = 10;

    if (true) {
        var value: i32 = value + 5;
        println("inner={}", value);
        value += 1;
        println("inner changed={}", value);
    }

    println("outer={}", value);
}
```

Ausführungsergebnis:

```text
inner=15
inner changed=16
outer=10
```

Der anfängliche Ausdruck `value + 5` in der inneren Deklaration liest den äußeren value. Nachdem die Initialisierung der neuen Variablen abgeschlossen ist, ist der innere value 15. Auch wenn Sie den inneren Wert auf 16 ändern, ändert sich der äußere Speicherplatz nicht. Nach der Sperrung sehen Sie draußen wieder value.

Wenn Sie umgekehrt nur `value += 1` ohne `var` im inneren Block ausführen, werden die vorhandenen sichtbaren Variablen geändert. Sehen Sie sich das Schlüsselwort an, um festzustellen, ob es sich um eine neue Deklaration oder eine Änderung eines vorhandenen Werts handelt.

## Lokale Variablen und Speicher auf oberster Ebene

Außerhalb der Funktion können Sie const und static verwenden. const stellt einen konstanten Wert dar und static ist der während der Ausführung beibehaltene Speicherplatz. Sie können nicht überall wie lokale Variablen deklariert werden.

<!-- wave-example: book-storage -->
```wave
const LIMIT: i32 = 3;
static visits: i32 = 0;

fun visit() {
    visits += 1;
    println("visit={}", visits);
}

fun main() {
    visit();
    visit();
    println("limit={}", LIMIT);
}
```

Ausführungsergebnis:

```text
visit=1
visit=2
limit=3
```

Selbst wenn visit zweimal aufgerufen wird, wird static nicht bei jedem Aufruf auf 0 zurückgesetzt. Wenn Sie es hingegen innerhalb einer Funktion als `var visits: i32 = 0;` deklarieren, wird der lokale Speicher bei jedem Aufruf initialisiert. Ein gemeinsam genutzter veränderlicher Zustand kann die Verfolgung des Verhaltens erschweren. Überlegen Sie daher zunächst, ob es mit den Ein- und Ausgängen der Funktion aufgelöst werden kann.

## häufige Fehler

- Wenn Deklaration und Zuweisung verwechselt werden und derselbe Name unnötigerweise erneut deklariert wird.
- Wenn Sie glauben, dass die Variable, die das Berechnungsergebnis speichert, automatisch Änderungen in der Eingabevariablen folgt.
- Wenn Sie denken, dass aufgrund des gleichen Typs auch die Einheiten wie Anzahl und Anzahl der Bytes gleich sind.
- Beim Lesen ohne Initialisierung oder beim Lesen des Ausgabearguments, wenn die Funktion fehlschlägt.
- Wenn Sie denken, dass der Name einer lokalen Variablen außerhalb des Blocks sichtbar ist.

Wenn Sie fehlerhaft nach einem Namen suchen, überprüfen Sie die Deklarationsposition und den Klammerbereich des Namens.

## Übung: Bestandsveränderungen berechnen

Der Anfangsbestand beträgt 20 Einheiten und es werden zweimal je 3 Einheiten verkauft. Drucken Sie den Restbestand und die Gesamtzahl der verkauften Einheiten aus. Bei jeder Bestandsänderung werden dieselben Variablen aktualisiert und auch die Verkaufsmengen werden separat kumuliert.

### Vollständige Lösung

<!-- wave-example: book-stock-solution -->
```wave
fun main() {
    var stock: i32 = 20;
    var sold: i32 = 0;
    var order: i32 = 3;
    stock -= order;
    sold += order;
    stock -= order;
    sold += order;
    println("stock={} sold={}", stock, sold);
}
```

Ausführungsergebnis:

```text
stock=14 sold=6
```

Lagerbestand und Verkaufsvolumen müssen sich gemeinsam verändern. Durch die Aktualisierung eines dieser Werte wird die Beziehung zwischen den Werten unterbrochen. Wir werden die wiederholte Auftragsverarbeitung durch Lernfunktionen und Schleifenanweisungen gruppieren.


## Ganzzahl- und Gleitkommatypen

Ganzzahltypen sind wie folgt:

- Unterzeichnet: `i8`, `i16`, `i32`, `i64`, `i128`, `i256`, `i512`, `i1024`
- Ohne Vorzeichen: `u8`, `u16`, `u32`, `u64`, `u128`, `u256`, `u512`, `u1024`
- Ganzzahlige Adressgröße: `isz`, `usz`
- Gleitkomma: `f32`, `f64`

`isz` ist ein vorzeichenbehafteter Ganzzahltyp, der der Adressgröße entspricht, und `usz` ist ein vorzeichenloser Ganzzahltyp, der der Adressgröße entspricht.

## Andere eingebaute Typen

|Typ|Benutzen|
| --- | --- |
| `bool` |`true` oder `false`|
| `char` |Ein vorzeichenloser 8-Bit-Zeichenwert. Nicht willkürlicher Codepunkttyp Unicode|
| `byte` |8-Bit-Byte-Wert|
| `str` |String-Bytes, die mit NUL enden|
| `ptr<T>` |Pointer-Targeting `T`|
| `array<T, N>` |Array fester Länge mit Elementtyp `T` und Länge `N`|

An Typspeicherorten können auch benutzerdefinierte Strukturen, Aufzählungen und Typaliase verwendet werden.

`var` ist die Syntax zum Deklarieren lokaler Variablen. Ein Typalias ist eine Grammatik, die denselben Typ mit einem Namen ausdrückt, der zum Kontext des Codes passt.
