---
translation_set_id: expressions
path: language/expressions-and-operators
locale: de
group: language
group_order: 2
order: 3
title: 3. Arithmetik, Vergleiche und Umrechnungen
summary: Lernen Sie Rechenreihenfolge, Ganzzahldivision, bitweise Operationen und cast.
---

## Berechnungsergebnisse und Berechnungsarten gemeinsam anzeigen

Ein Ausdruck ist Code, der einen Wert berechnet. Variablennamen, Literale, Funktionsaufrufe und Ausdrücke, die mehrere Werte mit Operatoren verbinden, sind allesamt Ausdrücke. Selbst wenn es in der Mathematik wie die gleiche Gleichung aussieht, ist das Ergebnis unterschiedlich, je nachdem, ob es sich um eine ganze Zahl oder eine reelle Zahl handelt und wie viele Bits sie enthält.

Ausgehend von einfachen Berechnungen werden in diesem Kapitel Klammern, Divisionen, logische Operationen, bitweise Operationen und Umwandlungen vorgestellt. Bei jedem Beispiel handelt es sich um eine vollständige main.wave-Datei, die Sie mit `wavec run main.wave` ausführen können.

## Bereich in Klammern eingeschlossen

<!-- wave-example: book-precedence -->
```wave
fun main() {
    var first: i32 = 2 + 3 * 4;
    var second: i32 = (2 + 3) * 4;

    println("{} {}", first, second);
}
```

Ausführungsergebnis:

```text
14 20
```

Die Multiplikation wird vor der Addition ausgewertet, daher ist der erste Ausdruck 2+12. In der zweiten Gleichung nehmen wir die Summe in Klammern, also 5, und multiplizieren sie dann mit 4. Das Ziel besteht nicht darin, weniger Klammern zu verwenden. Es wird empfohlen, es zu verwenden, damit der Leser den Umfang der Berechnung leicht verstehen kann.

Auch wenn derselbe Operator mehrmals verwendet wird, ist die Richtung der Verkettung wichtig. `20 - 5 - 3` ist `(20 - 5) - 3`, also 12. `20 - (5 - 3)` ist 18. Die genaue vollständige Sequenz befindet sich in [Operator-Referenz](/docs/de/language/expressions-and-operators).

## Ganzzahlige Division und Rest

Die Division zweier Ganzzahlen führt nicht zu einem gebrochenen Gleitkommaergebnis. Verwenden Sie Division für den Quotienten und den Restoperator für den Rest.

<!-- wave-example: book-division -->
```wave
fun main() {
    var items: i32 = 17;
    var box_size: i32 = 5;
    var full_boxes: i32 = items / box_size;
    var remaining: i32 = items % box_size;

    println("boxes={} remaining={}", full_boxes, remaining);
    println("negative quotient={}", -17 / 5);
}
```

Ausführungsergebnis:

```text
boxes=3 remaining=2
negative quotient=-3
```

Wenn man 17 Elemente in 5er-Gruppen aufteilt, erhält man 3 vollständige Gruppen und 2 verbleibende Elemente. Vorzeichenbehaftete Division schneidet in Richtung Null ab, also ist -17/5 -3. Dies unterscheidet sich vom Abrunden in Richtung negativer Unendlichkeit.

Eine Division durch 0 ist nicht möglich. signed Der Wert, der durch Division des Mindestwerts durch -1 erhalten wird, fällt nicht in denselben Typ. Funktionen, die diese Eingaben übernehmen, sollten vor dem Dividieren prüfen oder checked Mathematik API verwenden.

## Der Konvertierungspunkt verändert das Ergebnis.

Die folgenden beiden Ausdrücke werden beide in der Variablen f64 gespeichert, der Berechnungsprozess ist jedoch unterschiedlich.

<!-- wave-example: book-division-conversion -->
```wave
fun main() {
    var already_divided: f64 = (7 / 2) as f64;
    var floating_division: f64 = (7 as f64) / (2 as f64);

    if (already_divided == 3.0) {
        println("integer division lost the fraction");
    }

    if (floating_division == 3.5) {
        println("floating division kept the fraction");
    }
}
```

Ausführungsergebnis:

```text
integer division lost the fraction
floating division kept the fraction
```

Der erste Ausdruck führt eine Ganzzahldivision durch, um 3 zu erhalten, und wandelt sie dann in f64 um. Der zweite konvertiert die Operanden in f64, bevor eine Gleitkommadivision durchgeführt wird. Durch die Auswahl eines breiteren Typs für die endgültige Variable können zuvor verlorene Informationen nicht wiederhergestellt werden.

Gleitkommawerte sind Näherungswerte. Zwei Ergebnisse, die wie der gleiche Dezimalwert aussehen, sind möglicherweise nicht für einen exakten Gleichheitsvergleich geeignet. Wählen Sie eine Toleranz, die zu den Einheiten und dem Umfang des Problems passt. Ein festes Epsilon ist nicht für jede Berechnung geeignet.

## Vergleich durchgeführt bool

`<`, `<=`, `>`, `>=`, `==`, `!=` Überprüfen Sie die Beziehung. Ein Gleichheitszeichen `=` ist eine Zuweisung und zwei Gleichheitszeichen `==` sind Gleichheitsvergleiche.

<!-- wave-example: book-comparisons -->
```wave
fun main() {
    var age: i32 = 20;
    var minimum: i32 = 18;
    var eligible: bool = age >= minimum;

    if (eligible) {
        println("eligible");
    }

    if (age != minimum) {
        println("not exactly the boundary");
    }
}
```

Ausführungsergebnis:

```text
eligible
not exactly the boundary
```

„Mindestens 18“ schließt 18 ein; „größer als 18“ schließt es aus. Testen Sie 17, 18 und 19, um diese Grenze zu überprüfen. Vergleiche gemischter Typen hängen von Vorzeichen und Breite ab, daher kann die Konvertierung beider Operanden in den beabsichtigten Typ den Vergleich klarer machen.

## Logische Operationen und Kurzschlussauswertung

`&&` prüft, ob beide wahr sind, `||` prüft, ob mehr als einer wahr ist, und `!` dreht wahr/falsch um. Bei dieser Operation wird der Ausdruck auf der rechten Seite nicht immer ausgeführt.

<!-- wave-example: book-short-circuit -->
```wave
fun report() -> bool {
    println("right side evaluated");
    return true;
}

fun main() {
    var enabled: bool = false;

    if (enabled && report()) {
        println("both true");
    }

    if (!enabled || report()) {
        println("at least one true");
    }
}
```

Ausführungsergebnis:

```text
at least one true
```

In der ersten Bedingung ist enabled falsch, es besteht also keine Notwendigkeit, nach rechts zu schauen. Im zweiten Fall ist !enabled wahr, sodass die rechte Seite ebenfalls nicht benötigt wird. Daher erscheint die Ausgabe von report nie.

Damit können Sie den Nenner vor der Division überprüfen. Das Funktionskörperfragment `if (divisor != 0 && value / divisor > 2) { ... }` führt keine Division durch, wenn der Nenner 0 ist. Dieser Test allein löst jedoch keine anderen Grenzen wie signed Mindestwert/-1 auf.

`&&` hat Vorrang vor `||`. Geben Sie in komplexen Richtlinien Ihre Absicht in Klammern an, z. B. `(member && active) || admin`.

## Ganzzahlen verengen oder erweitern

`as` ist eine explizite Besetzung. Die höherwertigen Bits, die beim Einengen einer Ganzzahl verworfen werden, können nicht durch erneutes Erweitern wiederhergestellt werden.

<!-- wave-example: book-cast-chain -->
```wave
fun main() {
    var original: i32 = 300;
    var small: u8 = original as u8;
    var widened: i32 = small as i32;

    println("{} -> {} -> {}", original, small, widened);

    var negative: i8 = -1;
    var signed_value: i32 = negative as i32;
    var same_bits: u8 = negative as u8;

    println("{} {}", signed_value, same_bits);
}
```

Ausführungsergebnis:

```text
300 -> 44 -> 44
-1 255
```

Die unteren 8 Bits von 300 sind 44. Wenn Sie -1 auf den Typ signed erweitern, wird das Vorzeichen erweitert, um -1 beizubehalten. Wenn Sie die gleichen 8 Bits wie unsigned interpretieren, ist es 255.

Transformationen, die darauf abzielen, Werte innerhalb des Bereichs beizubehalten, und Transformationen, die versuchen, Speicherbits zu manipulieren, haben unterschiedliche Zwecke. Wenn Sie eine Benutzereingabe haben, prüfen Sie zunächst, ob es sich um den Zielbereich handelt, und konvertieren Sie ihn. Das Vorhandensein von cast garantiert nicht, dass der Wert in einem sicheren Bereich lag.

## Ersetzen durch bool

Das Konvertieren einer Ganzzahl in bool ergibt false für Null und andernfalls true. Dies wird nicht auf das niedrigste Bit gekürzt: 2 wird auch in true konvertiert. Bei Gleitkommawerten werden nur +0,0 und -0,0 in false konvertiert; Alle anderen Werte, einschließlich NaN und Unendlichkeiten, werden in true umgewandelt.

<!-- wave-example: book-bool-conversion -->
```wave
fun main() {
    var zero: i32 = 0;
    var two: i32 = 2;

    if (!(zero as bool)) {
        println("zero is false");
    }

    if (two as bool) {
        println("two is true");
    }
}
```

Ausführungsergebnis:

```text
zero is false
two is true
```

Zeiger-zu-bool-Konvertierungen werden nicht unterstützt. Vergleichen Sie den Zeiger explizit mit null, beispielsweise mit `pointer != null`. Ob eine Nicht-null-Adresse sicher lesbar ist, ist eine separate Frage.

## Bitoperationen

`&`, `|`, `^`, `~` decken jedes Bit der Ganzzahl ab. Es kann verwendet werden, um Berechtigungen oder Funktionen in Bits auszudrücken. Unten ist 1 die Leseberechtigung und 2 die Schreibberechtigung.

<!-- wave-example: book-bit-flags -->
```wave
const READ: u32 = 1;
const WRITE: u32 = 2;

fun main() {
    var permissions: u32 = READ | WRITE;

    if ((permissions & WRITE) != 0) {
        println("write enabled");
    }

    permissions = permissions & ~WRITE;
    println("remaining={}", permissions);
}
```

Ausführungsergebnis:

```text
write enabled
remaining=1
```

Fügen Sie mit OR ein Bit hinzu und prüfen Sie mit AND, ob ein bestimmtes Bit vorhanden ist. Erstellen Sie mit `~WRITE` eine Maske, bei der nur die relevanten Bits 0 sind, und löschen Sie sie. Die bitweisen Operationen `&`·`|` sind andere Operatoren als die Kurzschlussauswertung `&&`·`||` von bool.

## Breite und Anzahl der Schichten

Eine Linksverschiebung verschiebt Bits nach links und verwirft hohe Bits außerhalb der Breite des Operanden. Bei einer Rechtsverschiebung wird die Vorzeichenerweiterung für vorzeichenbehaftete Werte und die Nullerweiterung für vorzeichenlose Werte verwendet. Das Ergebnis hat immer den Typ des linken Operanden.

<!-- wave-example: book-shifts -->
```wave
fun main() {
    var bits: u8 = 129;
    var shifted: u8 = bits << 1;
    var negative: i32 = -8;
    var positive: u32 = 8;

    println("left={}", shifted);
    println("right={} {}", negative >> 1, positive >> 1);
}
```

Ausführungsergebnis:

```text
left=2
right=-4 4
```

Durch das Verschieben des u8-Werts 129 um eins nach links wird sein höchstes Bit verworfen und es bleibt 2 übrig. Der ursprüngliche Verschiebungszählwert muss nicht negativ und kleiner als die Bitbreite des linken Operanden sein. Für u8 sind gültige Zählwerte 0 bis 7. Ein ungültiger konstanter Zählwert ist ein Fehler bei der Kompilierung; Eine ungültige Laufzeitzählung führt zu einem Trap.

## Konvertieren von Gleitkommawerten in Ganzzahlen

Bei der Konvertierung von Gleitkommazahlen in Ganzzahlen wird zunächst in Richtung Null gekürzt und dann der Ganzzahlbereich des Ziels überprüft. NaN-, Unendlich- und außerhalb des Bereichs liegende Ergebnisse sind ungültig. Ungültige Konstantenkonvertierungen führen zu einem Fehler bei der Kompilierung. Ungültige Laufzeitkonvertierungen verursachen einen Trap.

Ein Trap gibt keinen Fehlerwert von der Funktion zurück. Entwerfen Sie für behebbare Konvertierungsfehler eine Schnittstelle, die den Bereich vor der Konvertierung überprüft. Um die gespeicherten Bits eines Gleitkommawerts zu erhalten, verwenden Sie die Bitkonvertierungsfunktionen in `std::math::float` anstelle einer numerischen Umwandlung.

## Übung und Lösung

Schreiben Sie eine Funktion, die 137 Won als Wechselgeld in 50 Won-Einheiten und den Rest dividiert und nur ganze Zahlen im Bereich von 0 bis 255 in u8 ändert. Dieses Beispiel demonstriert den Validierungsschritt als separate Funktion, die außerhalb des Bereichs als -1 angibt.

<!-- wave-example: book-expression-solution -->
```wave
fun checked_byte(value: i32) -> i32 {
    if (value < 0 || value > 255) {
        return -1;
    }

    var small: u8 = value as u8;
    return small as i32;
}

fun main() {
    var money: i32 = 137;

    println("coins={} remainder={}", money / 50, money % 50);
    println("{} {} {}", checked_byte(0), checked_byte(255), checked_byte(256));
}
```

Ausführungsergebnis:

```text
coins=2 remainder=37
0 255 -1
```

Der Grund, warum -1 als Fehlerindikator verwendet werden kann, liegt darin, dass der Erfolgsbereich zwischen 0 und 255 liegt. Wenn eine beliebige Ganzzahl ein Erfolgswert sein kann, ist eine andere Ergebnisdarstellung erforderlich. Wir setzen diesen Entwurf später im Kapitel zur Fehlerbehandlung fort.


## Priorität

Die Priorität der Operatoren ist wie folgt, beginnend mit der höchsten:

1. Grundlegende Ausdrücke und Postfix-Zugriff: Funktionsaufrufe, Feldzugriff, Indizierung, Postfix `++`·`--`
2. Unäre Operationen: `!`, `~`, `&`, `deref`, Präfix `++`·`--`, unäre `+`·`-`
3. `as` Typkonvertierung
4. `*`, `/`, `%`
5. `+`, `-`
6. `<<`, `>>`
7. `<`, `<=`, `>`, `>=`
8. `==`, `!=`
9. Bit `&`
10. Bit `^`
11. Bit `|`
12. `&&`
13. `||`
14. Zuordnung und zusammengesetzte Zuordnung

Die Kettenbelegung erfolgt von rechts. Wenn Sie verschiedene Arten von Operatoren mischen, verwenden Sie Klammern, um die Reihenfolge der Auswertung klar auszudrücken.

## Subjekte, die ersetzt werden können

Zuweisung, `++` und `--` erfordern einen Ausdruck, der einen Speicherort angibt, z. B. eine Variable, ein Feld, ein Array-Element oder einen dereferenzierten Zeiger. Das Schreiben auf einen `const` ist nicht erlaubt.

## Verschiebung

Das Ergebnis einer Verschiebung hat immer den Typ des linken Operanden; Der Typ des richtigen Operanden erweitert die Berechnung nicht. Linksverschiebungen verwerfen hohe Bits über diese Breite hinaus. Nach rechts verschiebt sich vorzeichenerweiternde vorzeichenbehaftete Werte und nullerweiternde vorzeichenlose Werte.

Die Schichtanzahl muss eine Ganzzahl sein, deren ursprünglicher Wert `0 <= n < LHS bit width` erfüllt. Es wird vor jeder Kürzung auf einen kleineren Typ überprüft. Eine ungültige Konstantenanzahl ist ein Fehler bei der Kompilierung. Eine ungültige Laufzeitzählung führt zu einem Trap.

## Konvertierungen mit Bool- und Gleitkommawerten

Die Ganzzahl-in-bool-Konvertierung erzeugt false für Null und true für jeden anderen Wert. Die Konvertierung von Gleitkomma zu bool erzeugt false nur für +0,0 und -0,0; NaN und positive oder negative Unendlichkeit ergeben true. Zeiger-zu-bool-Umwandlungen werden nicht unterstützt: explizit mit `pointer != null` vergleichen.

Bei der Konvertierung von Gleitkommazahlen in Ganzzahlen wird in Richtung Null gekürzt und anschließend der Zielbereich überprüft. NaN-, Unendlich- und außerhalb des Bereichs liegende Ergebnisse sind ungültig. Ungültige Konstantenkonvertierungen sind Fehler bei der Kompilierung. Ungültige Laufzeitkonvertierungen verursachen einen Trap. Eine Falle ist keine behebbare Fehlerrückgabe.

`&&` und `||` sind Kurzschlussauswertungen. Die Nebenwirkungen des nicht ausgeführten rechten Operanden treten nicht auf. Die Ergebnisse mit kleinen Werten können Sie in [Rechenunterricht](/docs/de/language/expressions-and-operators) überprüfen.

## Schicht, die absichtlich scheitert

Die Anzahl der Verschiebungen für einen 8-Bit-Wert muss 0 bis 7 betragen. Die folgenden Programme sollten vor der Ausführung abgelehnt werden.

<!-- wave-example: reject-shift-count -->
```wave
fun main() {
    var value: u8 = 1;
    var result: u8 = value << 8;
}
```
