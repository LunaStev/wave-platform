---
translation_set_id: learn-errors
path: language/errors
locale: de
group: language
group_order: 2
order: 12
title: 12. Darstellung und Behebung von Fehlern
summary: Trennen Sie Fehler von Normalwerten und bereinigen Sie Ressourcen von Fehlerpfaden.
---

## Ergebnis der Ausfallgradfunktion

Fehlende Dateien oder Eingaben außerhalb des zulässigen Bereichs kommen in Programmen ganz natürlich vor. Die Behandlung eines Fehlers umfasst mehr als nur das Drucken einer Nachricht. Dabei handelt es sich um den Prozess der Unterscheidung von Fehlern, der Überprüfung des Status der bereits geleisteten Arbeit, der Organisation erworbener Ressourcen und der anschließenden Entscheidung, ob fortgefahren oder beendet werden soll.

In diesem Kapitel beginnen wir mit der Fehlerdarstellung einer kleinen Funktion und fahren mit der Ergebnisstruktur variant, der frühen Rückkehr und der Ressourcenbereinigung fort.

## Fehlermarkierungen dürfen sich nicht mit Erfolgswerten überschneiden

Der Grund dafür, dass -1 bei einer Array-Suche als nicht gefunden verwendet wird, liegt darin, dass der effektive Index 0 oder mehr ist. Andererseits kann bei Berechnungen, bei denen jede ganze Zahl ein normales Ergebnis sein kann, -1, wenn es als Fehler bezeichnet wird, nicht vom Normalwert -1 unterschieden werden.

0 ist auch ein häufig missverstandener Wert. Die Länge des leeren Strings 0, die erste Position 0 und die Anzahl der übertragenen Bytes 0 haben für verschiedene Funktionen unterschiedliche Bedeutungen. Beurteilen Sie den Erfolg nicht nur, weil der Rückgabewert ungleich Null ist.

## Geben Sie gemeinsam Erfolg und Wert zurück

<!-- wave-example: book-error-result -->
```wave
struct Division {
    ok: bool;
    value: i32;
}

fun divide_nonnegative(left: i32, right: i32) -> Division {
    if (left < 0 || right <= 0) {
        return Division {
            ok: false,
            value: 0
        };
    }

    return Division {
        ok: true,
        value: left / right
    };
}

fun main() {
    var result: Division = divide_nonnegative(0, 3);

    if (!result.ok) {
        println("invalid input");
        return;
    }

    println("value={}", result.value);
}
```

Ausführungsergebnis:

```text
value=0
```

Ein normales Ergebnis kann auch 0 sein. Anstatt auf value zu schauen und zu raten, ob es erfolgreich war, überprüfen Sie zuerst ok. Selbst wenn das Feld value im fehlgeschlagenen Ergebnis vorhanden ist, bedeutet dies nicht, dass es sich um den zu verwendenden Wert handelt.

Die Eingaberegeln für dieses Beispiel lauten, dass die linke Seite 0 oder größer und die rechte Seite positiv ist. Der Umfang wird im Funktionsnamen und in der Beschreibung angezeigt, um ihn von der typischen Ganzzahldivision signed zu unterscheiden.

## Unterscheiden Sie zwischen Fehlerursachen

Fügen Sie Fehlerinformationen hinzu, um abhängig von der Fehlerursache zusätzliche Hinweise oder eine Wiederherstellung zu erhalten. Es gibt auch eine Möglichkeit, die Funktion in eine Funktion, die den Eingabebereich prüft, und eine Funktion, die ihn berechnet, zu unterteilen.

<!-- wave-example: book-error-codes -->
```wave
struct Check {
    ok: bool;
    error: i32;
}

fun check_quantity(value: i32) -> Check {
    if (value < 1) {
        return Check {
            ok: false,
            error: 1
        };
    }

    if (value > 100) {
        return Check {
            ok: false,
            error: 2
        };
    }

    return Check {
        ok: true,
        error: 0
    };
}

fun main() {
    var result: Check = check_quantity(101);

    if (!result.ok) {
        match (result.error) {
            1 => {
                println("quantity is too small");
            }
            2 => {
                println("quantity is too large");
            }
            _ => {
                println("invalid quantity");
            }
        }
    }
}
```

Ausführungsergebnis:

```text
quantity is too large
```

Die Bedeutung der Zahlen wird durch diese Funktion definiert. Es kann nicht als dasselbe wie Fehler Nr. 1 oder 2 in anderen Bibliotheken angesehen werden. Im öffentlichen API erspart die Benennung von Fehlerkonstanten oder -typen dem Aufrufer das Merken von Zufallszahlen.

## Ergebnisse mit variant trennen

Die Beziehung, dass ein Erfolgswert und ein Fehler nicht gleichzeitig bestehen können, kann als variant ausgedrückt werden.

<!-- wave-example: book-error-variant -->
```wave
variant ByteResult {
    Value(u8),
    OutOfRange(i32)
}

fun to_byte(value: i32) -> ByteResult {
    if (value < 0 || value > 255) {
        return ByteResult::OutOfRange(value);
    }

    return ByteResult::Value(value as u8);
}

fun main() {
    var result: ByteResult = to_byte(300);

    match (result) {
        ByteResult::Value(value) => {
            println("byte={}", value);
        }
        ByteResult::OutOfRange(original) => {
            println("out of range: {}", original);
        }
    }
}
```

Ausführungsergebnis:

```text
out of range: 300
```

Der Fehler enthielt den ursprünglichen Eingabewert. Für den Anrufer ist es einfacher, das Problem zu erklären, als einfach false zurückzusenden. Auch die Art der Daten wird berücksichtigt, indem sensible Eingaben wie Passwörter oder Token nicht im Protokoll verbleiben.

## Erleichtern Sie die Lesbarkeit normaler Routen durch frühzeitige Rückkehr

Bei einer mehrstufigen Prüfung ist es nicht erforderlich, den gesamten guten Code tief in if einzufügen. Wenn dies fehlschlägt, können Sie sofort zurückkehren und den unten aufgeführten Erfolgspfad fortsetzen.

<!-- wave-example: book-error-early-return -->
```wave
fun show_total(quantity: i32, price: i32) {
    if (quantity < 1 || quantity > 100) {
        println("invalid quantity");
        return;
    }

    if (price < 0 || price > 100000) {
        println("invalid price");
        return;
    }

    var total: i32 = quantity * price;
    println("total={}", total);
}

fun main() {
    show_total(3, 1200);
    show_total(0, 1200);
    show_total(3, -1);
}
```

Ausführungsergebnis:

```text
total=3600
invalid quantity
invalid price
```

Nach bestandener Prüfung können Sie davon profitieren, dass quantity und price im vorgegebenen Bereich liegen. Der Bereich wurde so eingestellt, dass auch die Zwischenmultiplikation in den Bereich i32 fällt. Wenn Sie eine vorzeitige Rückgabe hinzufügen, sollten Sie auch prüfen, ob Sie zu diesem Zeitpunkt bereits über Ressourcen verfügen.

## Speicherbereinigung von Fehlerpfaden

<!-- wave-example: book-error-cleanup -->
```wave
import("std::buffer::alloc")::{
    Buffer, buffer_init, buffer_free
};
import("std::buffer::write")::{
    buffer_append_str
};

fun make_message(allowed: bool) -> i32 {
    var data: Buffer;

    if (buffer_init(&data, 16) < 0) {
        return 1;
    }

    if (!allowed) {
        buffer_free(&data);
        return 2;
    }

    if (buffer_append_str(&data, "ready") < 0) {
        buffer_free(&data);
        return 3;
    }

    println("bytes={}", data.len);

    if (buffer_free(&data) < 0) {
        return 4;
    }

    return 0;
}

fun main() {
    println("status={}", make_message(false));
}
```

Ausführungsergebnis:

```text
status=2
```

Geben Sie Buffer auch auf Fehlerpfaden frei, die keine Daten ausgeben. Anstatt jeden Fehler in einen einzigen return umzuwandeln, ist es wichtig, klar zu verwalten, was in jeder Filiale vorhanden ist.

Es gibt auch API, bei dem die Bereinigung selbst fehlschlägt. Entscheiden Sie, wie Sie Fehler bewahren und Fehler in der Originalarbeit bereinigen. Dieses kleine Beispiel gibt zunächst die Aufgabenfehlernummer zurück. Größere Programme können beide getrennt aufzeichnen.

## Teilerfolge werden nicht automatisch annulliert

Wenn Sie einige Bytes in eine Datei schreiben und der Schreibvorgang dann fehlschlägt, gehen die bereits geschriebenen Bytes nicht verloren. Möglicherweise hat auch das andere Ende des Netzwerks einige Daten empfangen. Das Wiederholen derselben Aufgabe von Anfang an kann zu doppelten Datensätzen führen.

Umgekehrt bleiben beim Lesen von checked byte cursor die Position und der Ausgabewert erhalten, wenn dies fehlschlägt. Diese Funktionen können mehr Eingaben empfangen und es am selben Ort erneut versuchen. Anstatt die Regel „Wenn es fehlschlägt, ändert sich nichts“ auf jede API anzuwenden, prüfen Sie die Dokumentation für diese Funktion.

## Behebbare Fehler und Fallen

Ein ungültiger Dateipfad oder eine ungültige Benutzereingabe können als Fehlerwert gemeldet werden, damit der Anrufer eine Wiederherstellung durchführen kann. Ein Trap, der durch eine ungültige Laufzeitverschiebungszählung oder eine Gleitkomma-in-Ganzzahl-Konvertierung verursacht wird, ist ein anderer Mechanismus.

Programme, die weiterhin ausgeführt werden müssen, müssen ihre Eingaben vor gefährlichen Vorgängen überprüfen. assert ist auch nicht dazu gedacht, als Mittel zur Behandlung ordnungsgemäßer Fehler bei Benutzereingaben missbraucht zu werden. Wenn Sie dem Benutzer die Möglichkeit geben müssen, seine Eingaben erneut einzugeben, geben Sie Ergebnisse zurück, um den Kontrollfluss fortzusetzen.

## Übung und vollständige Lösung

Schreiben Sie eine Funktion, die ein Array-Element nach Index liest und fehlschlägt, wenn der Index negativ oder größer oder gleich der Länge ist. Verwenden Sie eine Ergebnisstruktur, da Null ein gültiger Elementwert sein kann.

<!-- wave-example: book-error-solution -->
```wave
struct Lookup {
    ok: bool;
    value: i32;
}

fun at(values: ptr<i32>, length: i32, index: i32) -> Lookup {
    if (index < 0 || index >= length) {
        return Lookup {
            ok: false,
            value: 0
        };
    }

    return Lookup {
        ok: true,
        value: deref values[index]
    };
}

fun main() {
    var values: array<i32, 3> = [0, 10, 20];
    var first: Lookup = at(&values[0], 3, 0);
    var outside: Lookup = at(&values[0], 3, 3);

    if (first.ok) {
        println("value={}", first.value);
    }

    if (!outside.ok) {
        println("out of bounds");
    }
}
```

Ausführungsergebnis:

```text
value=0
out of bounds
```

Es ist die Bedingung des Aufrufers, dass der Zeiger und length das tatsächlich lesbare Array darstellen. Dies ist keine Funktion, die zufällige Adressen allein durch die Überprüfung des Index sicher macht. Bitte lesen Sie separat, für welche Prüfungen die Funktion zuständig ist und welche Bedingungen der Aufrufer gewährleisten muss.
