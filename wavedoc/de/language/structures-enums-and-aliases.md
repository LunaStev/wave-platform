---
translation_set_id: data-types
path: language/structures-enums-and-aliases
locale: de
group: language
group_order: 2
order: 8
title: 8. Strukturen, Aufzählungen und Varianten
summary: Erfahren Sie mehr über Felder, Strukturinitialisierung und die Rollen von enum und variant.
---

## Beziehungen zwischen Daten als Typen ausdrücken

Wenn der Preis und die Anzahl der Produkte jeweils als Variablen übergeben werden, lässt sich anhand des Codes nur schwer erkennen, ob die beiden Werte zum selben Produkt gehören. Eine Struktur gruppiert verwandte Felder. enum stellt einen benannten Zustand dar und variant speichert für jeden Fall verschiedene Daten zusammen.

Bei den drei Funktionen handelt es sich nicht um Syntaxen, die sich gegenseitig ersetzen. Wählen Sie je nachdem, was Sie ausdrücken möchten.

|etwas zum Ausdruck bringen|auswählen|ja|
| --- | --- | --- |
|Mehrere Felder gleichzeitig vorhanden| struct |Stückpreis und Produktmenge|
|Benannter ganzzahliger Zustand| enum |Warten/Fortfahren/Fertig stellen|
|Jeweils unterschiedliche Daten| variant |Erfolgswert oder Fehler|

## Strukturdeklaration und Wertschöpfung

<!-- wave-example: book-struct-first -->
```wave
struct Product {
    price: i32;
    quantity: i32;
}

fun main() {
    var item: Product = Product {
        price: 1500,
        quantity: 2
    };

    println("{} {}", item.price, item.quantity);
}
```

Ausführungsergebnis:

```text
1500 2
```

Felder in Deklarationen enden mit einem Semikolon, und beim Erstellen von Werten werden Felder und Werte mit Doppelpunkten verbunden und durch Kommas getrennt. Das Definieren des Typs und das Erstellen des tatsächlichen Werts sind zwei verschiedene Schritte. Durch die Angabe des Typs Product wird nicht automatisch Speicherplatz für ein Produkt geschaffen.

Der Zugriff auf das Feld erfolgt als `item.price`. Wenn Sie mehrere item desselben Typs erstellen, können Sie in jedem unterschiedliche Werte speichern.

## Übergabe einer Struktur an eine Funktion

<!-- wave-example: book-struct-function -->
```wave
struct Product {
    price: i32;
    quantity: i32;
}

fun subtotal(item: Product) -> i32 {
    return item.price * item.quantity;
}

fun main() {
    var item: Product = Product {
        price: 1500,
        quantity: 2
    };

    println("{}", subtotal(item));

    item.quantity = 3;
    println("{}", subtotal(item));
}
```

Ausführungsergebnis:

```text
3000
4500
```

Die Funktion erhält als Typ die Beziehung, dass der Stückpreis und die Menge zum gleichen Produkt gehören. Dies ist eine Funktion, die die als Wert übergebene ganzzahlige Feldstruktur liest. Jede Funktion, die den Speicher des Aufrufers ändern möchte, kann so gestaltet werden, dass sie einen Zeiger akzeptiert.

Wenn die Struktur über Zeigerfelder verfügt, wird beim Kopieren des Werts auch die Adresse kopiert. Dies ist keine Funktion zum tiefen Kopieren in separate Zuordnungen. Typen, die Ressourcen wie Dateihandles oder Buffer enthalten, müssen gemeinsam Kopier- und Freigaberegeln angeben.

## Methode und proto

Verwandte Funktionen können in Methodenform gruppiert werden. proto ist eine Methode zum Schreiben der Methode einer Struktur als separaten Block.

<!-- wave-example: book-struct-method -->
```wave
struct Counter {
    value: i32;
}

proto Counter {
    fun current(self: Counter) -> i32 {
        return self.value;
    }
}

fun main() {
    var counter: Counter = Counter {
        value: 3
    };

    println("{}", counter.current());
}
```

Ausführungsergebnis:

```text
3
```

`self: Counter` ist ein Parameter, der einen Wert erhält. Die Verwendung einer Methodenaufrufnotation macht sie nicht automatisch zu einer Methode, die das Original ändert. Bitte lesen Sie gemeinsam den Typ von self und was er im Text bewirkt.

Sie können nicht erwarten, dass ein Feld nur einen gültigen Status hat, nur weil Sie ihm eine Methode angehängt haben. Wenn es ungültige Kombinationen gibt, die der Benutzer mit öffentlichen Feldern erstellen kann, sollte Ihre Funktion nach diesen suchen oder eine Erstellungsregel bereitstellen.

## Benennen Sie den Staat mit enum

<!-- wave-example: book-enum-state -->
```wave
enum State -> i32 {
    Ready = 0,
    Running,
    Finished,
}

fun main() {
    var state: State = State::Ready;

    if (state == State::Ready) {
        println("ready");
    }

    state = State::Running;

    if (state == State::Running) {
        println("running");
    }
}
```

Ausführungsergebnis:

```text
ready
running
```

`-> i32` ist ein ganzzahliger Typ, der in Ausdrücken verwendet wird. Der erste Wert wird auf 0 gesetzt, und nachfolgende ausgelassene Werte sind um 1 größer als der vorherige Wert. Anstatt einfach 0 und 1 in Ihrem Code zu vergleichen, wird die Verwendung von State::Ready und State::Running die Bedeutung offenbaren.

enum Das Vorhandensein eines Namens schränkt Statusübergänge nicht automatisch ein. Regeln, z. B. ob eine Rückkehr von Finished nach Running möglich ist, müssen als Funktionen implementiert werden.

## Fälle und Daten mit variant verbinden

Wenn es nur im Erfolgsfall einen Wert gibt und im Fehlerfall Fehlerinformationen benötigt werden, kann dieser als variant ausgedrückt werden.

<!-- wave-example: book-variant-result -->
```wave
variant Result {
    Value(i32),
    Error(i32)
}

fun divide(left: i32, right: i32) -> Result {
    if (left < 0 || right <= 0) {
        return Result::Error(1);
    }

    return Result::Value(left / right);
}

fun show(result: Result) {
    match (result) {
        Result::Value(value) => {
            println("value={}", value);
        }
        Result::Error(code) => {
            println("error={}", code);
        }
    }
}

fun main() {
    show(divide(12, 3));
    show(divide(12, 0));
}
```

Ausführungsergebnis:

```text
value=4
error=1
```

Result::Value und Result::Error enthalten jeweils payload. Auch wenn sie denselben Integer-Typ enthalten, werden bestimmte Fälle unterschieden. Der Anrufer prüft den Fall mit match und verwendet payload innerhalb dieses arm.

divide oben ist ein kleines Beispiel, das keine negativen Operanden unterstützt. Da der Eingabebereich spezifiziert ist, sollte er nicht mit einer Funktion verwechselt werden, die alle Grenzen der regulären Division signed abdeckt.

## Wenn payload nicht vorhanden ist

Daten sind nicht in allen Fällen erforderlich. Der Zustand der Wertlosigkeit kann als separater Fall ausgedrückt werden.

<!-- wave-example: book-variant-empty -->
```wave
variant Lookup {
    Found(i32),
    Missing
}

fun main() {
    var result: Lookup = Lookup::Missing;

    match (result) {
        Lookup::Found(value) => {
            println("found={}", value);
        }
        Lookup::Missing => {
            println("missing");
        }
    }
}
```

Ausführungsergebnis:

```text
missing
```

Anstatt eine beliebige ganze Zahl als „none“ zu reservieren, haben wir den Fall Missing verwendet. Egal wie hoch der Erfolgswert ist, die Bedeutung überschneidet sich nicht.

match bis `_` kümmern sich um die übrigen Fälle. Wenn Sie möchten, dass jeder Anrufer erneut überprüft wird, wenn ein neuer Fall hinzugefügt wird, ist es besser, alle Fälle explizit zu trennen. Unabhängig davon, für welche Methode Sie sich entscheiden, stellen Sie sicher, dass keine unverarbeiteten Eingaben vorhanden sind.

## Struktur und variant zusammen verwenden

Die Auswahl verschiedener Daten kann als variant ausgedrückt werden, und die Gruppierung mehrerer zu einem Fall gehörender Felder kann als Struktur ausgedrückt werden. Wenn das Ergebnis der Auftragsabwicklung beispielsweise erfolgreich ist, kann es so gestaltet werden, dass es eine Belegstruktur enthält, und wenn es sich um einen Fehler handelt, kann es so gestaltet werden, dass es eine Fehlernummer enthält.

Die Speicherlebensdauerregel verschwindet nicht, auch wenn der Wert andere Werte enthält. Wenn ein Zeiger in variant gespeichert ist, wird separat festgelegt, ob der Zeiger gültig ist und wer ihn freigibt. Auch beim Speichern in einem externen Dateiformat müssen Sie für jedes Feld eine Codierung angeben, anstatt den Strukturspeicher unverändert zu speichern.

## Übung und vollständige Lösung

Erstellen Sie ein Urteil, das nur Werte zwischen 0 und 100 akzeptiert. Gültige Werte werden als Grade(score) ausgedrückt, der Rest wird als Invalid ausgedrückt.

<!-- wave-example: book-model-solution -->
```wave
variant CheckedScore {
    Grade(i32),
    Invalid
}

fun check_score(score: i32) -> CheckedScore {
    if (score < 0 || score > 100) {
        return CheckedScore::Invalid;
    }

    return CheckedScore::Grade(score);
}

fun main() {
    var result: CheckedScore = check_score(87);

    match (result) {
        CheckedScore::Grade(score) => {
            println("accepted={}", score);
        }
        CheckedScore::Invalid => {
            println("invalid score");
        }
    }
}
```

Ausführungsergebnis:

```text
accepted=87
```

Ändern Sie die Eingabe auf -1, 0, 100, 101, um die Grenzen zu überprüfen. Da sich Erfolg und Misserfolg nicht denselben ganzzahligen Raum teilen, verringert der Aufrufer das Risiko, versehentlich Fehlerwerte zur Durchschnittsberechnung hinzuzufügen.


## Welchen Ausdruck soll ich wählen?

|Form der Daten|passenden Ausdruck|ja|
| --- | --- | --- |
|Mehrere Werte desselben Typs|Array|10 Punkte|
|Mehrere miteinander verbundene Felder|Struktur|Name und Punktzahl|
|Benannter Zustandswert| enum | Ready, Running, Stopped |
|Zusätzliche Daten, die je nach Bundesland variieren| variant | Value(i32), Error(str) |
|Kontextueller Name für einen vorhandenen Typ|Typalias| UserId = u64 |

Berücksichtigen Sie bei der Auswahl einer Datenstruktur nicht nur die Werte, die Sie speichern möchten, sondern auch, welche falschen Zustände Sie darstellen können. Eine Struktur mit den Feldern „Erfolg/Misserfolg“ und „Wert/Fehler“ kann zu einer falschen Kombination führen, aber variant kann in jedem Fall in payload unterschieden werden.

## Speicherplatzierung und externe Daten

Der Speicher der Struktur kann leeren Speicherplatz enthalten, um die Ausrichtung zwischen Feldern sicherzustellen. Die einfache Addition der Feldgrößen ergibt nicht immer die Gesamtgröße der Struktur. Wenn Sie die Größe und Ausrichtung kennen müssen, verwenden Sie [mem Layoutfunktion](/docs/de/reference/memory-and-buffer).

Bei Dateien oder Netzwerknachrichten können die Byte-Reihenfolge und -Länge eindeutig bestimmt werden, indem die Felder mit der Funktion [bytes](/docs/de/stdlib/bytes) in der richtigen Reihenfolge kodiert werden. Wenn Sie Strukturen an andere Sprachen übergeben, richten Sie die externe Deklaration von [FFI](/docs/de/language/modules-imports-and-ffi) am Ziel ABI aus.

[Strukturieren Sie Lernen und Üben](/docs/de/language/structures-enums-and-aliases) · [variant](/docs/de/language/variants)
