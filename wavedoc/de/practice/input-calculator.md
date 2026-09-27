---
translation_set_id: practice-input-calculator
path: practice/input-calculator
locale: de
group: practice
group_order: 4
order: 1
title: Projekt: Ein Rechner, der Eingaben validiert
summary: Verknüpfen Sie Eingaben, Grenzprüfungen, Funktionen und Exit-Codes.
---

## Ziele und Maßnahmen

Geben Sie die Menge und den Stückpreis ein und berechnen Sie die Gesamtsumme. Geben Sie zwei durch ein Leerzeichen oder Zeilenumbruch getrennte Ganzzahlen ein. In diesem Beispiel werden nur Mengen von 1 bis 1000 und Einzelpreise von 0 bis 100000 akzeptiert, sodass die Berechnung innerhalb des Multiplikationsbereichs i32 erfolgt.

Speichern Sie es unter `main.wave`, führen Sie `wavec run main.wave` aus und geben Sie dann `3 1200` ein. Die vom Programm ausgegebenen Ergebnisse sind wie folgt. Die Sichtbarkeit der im Terminal eingegebenen Zeichen erfolgt unabhängig von der Programmausgabe.

<!-- wave-example: input-calculator -->
```wave
fun total(quantity: i32, price: i32) -> i32 {
    return quantity * price;
}

fun main() -> i32 {
    var quantity: i32 = 0;
    var price: i32 = 0;
    input("{} {}", quantity, price);
    if (quantity < 1 || quantity > 1000 || price < 0 || price > 100000) {
        println("out of range");
        return 1;
    }

    println("total={}", total(quantity, price));
    return 0;
}
```

Ausführungsergebnis:

```text
total=3600
```

## Überprüfen Sie auch auf Fehler

Wenn Sie `0 1200` eingeben, wird `out of range` erwartet und der Exit-Code 1 wird beendet. Fehler beim numerischen Parsen in `input` und die Überprüfung des Arbeitsumfangs des Programms sind zwei verschiedene Schritte. Nicht-numerisches Token/Typ überschreitet Bereich/vor erforderlicher Eingabe EOF ist ein Eingabefehler. Bei dieser integrierten Eingabe handelt es sich nicht um eine Schnittstelle, die einen Fehler zurückgibt und eine erneute Eingabe erfordert. Wenn Sie einen wiederherstellbaren Parser benötigen, konfigurieren Sie den Überprüfungsprozess selbst, indem Sie die Bytes mit [io](/docs/de/stdlib/files-io) lesen.

## Erweiterte Übungen und Kommentare

Nehmen Sie den Abzinsungssatz als dritte Eingabe und prüfen Sie, ob er zwischen 0 und 100 liegt. Um große Zwischenmultiplikationen zu vermeiden, müssen Sie den Berechnungsbereich mit i64 erweitern und den Bereich überprüfen, wenn Sie das Ergebnis eingrenzen. Wenn Sie einfach den letzten Ergebnistyp erweitern, wurden möglicherweise bereits Zwischenberechnungen für schmale Typen durchgeführt.

[Konsole I/O](/docs/de/language/console-io-and-formatting) · [Weiter: Dateien lesen](/docs/de/practice/file-reader)

## Warum zuerst den Berechnungsumfang festlegen?

Die größten Eingaben sind Menge 1000 und Stückpreis 100000. Die Multiplikation der beiden Werte ergibt 100000000, liegt also im Bereich i32. Diese Bereichsprüfung ist erforderlich, damit das Multiplikationsergebnis der Funktion total unverändert verwendet werden kann.

main, das Eingaben empfängt, ist für Eingaben und Fehlermeldungen verantwortlich, und total ist nur für Berechnungen verantwortlich. Auch wenn Sie später auf das Lesen von Aufträgen aus einer Datei umsteigen, können Sie die Berechnungsfunktion weiterhin nutzen.

## Komplette Rabattberechnung

Sie können die Gesamtsumme mit 100 multiplizieren, indem Sie den Rabattprozentsatz hinzufügen. Umrechnen, um Zwischenberechnungen als i64 durchzuführen, und dann den Abzinsungssatz anwenden. Da der Bruchteil der Ganzzahldivision verworfen wird, wird in diesem Beispiel der Betrag nach dem Rabatt auf eine Ganzzahl gekürzt.

<!-- wave-example: book-calculator-discount -->
```wave
fun discounted_total(quantity: i32, price: i32, discount: i32) -> i64 {
    var subtotal: i64 = (quantity as i64) * (price as i64);
    var remaining: i64 = 100 - (discount as i64);
    return subtotal * remaining / 100;
}

fun main() -> i32 {
    var quantity: i32 = 0;
    var price: i32 = 0;
    var discount: i32 = 0;

    input("{} {} {}", quantity, price, discount);

    if (quantity < 1 || quantity > 1000) {
        println("invalid quantity");
        return 1;
    }

    if (price < 0 || price > 100000) {
        println("invalid price");
        return 2;
    }

    if (discount < 0 || discount > 100) {
        println("invalid discount");
        return 3;
    }

    println("total={}", discounted_total(quantity, price, discount));
    return 0;
}
```

Ausführungsergebnis:

```text
total=3240
```

Das obige Ergebnis ist, wenn `3 1200 10` eingegeben wird. Die Ausgabe beträgt 3240, was 10 % vom ursprünglichen Gesamtwert von 3600 abgezogen wird.

|Eingabe|erwartetes Ergebnis|Pfad zur Überprüfung|
| --- | --- | --- |
| `3 1200 0` | `total=3600` |kein Rabatt|
| `3 1200 100` | `total=0` |Voller Rabatt|
| `3 1200 101` | `invalid discount` |Der Diskontsatz überschreitet den Bereich|
| `1000 100000 0` | `total=100000000` |maximale Eingabe|
| `0 1200 10` | `invalid quantity` |Menge unterhalb des Bereichs|

## nächste Übung

Versuchen Sie, die Funktion so zu ändern, dass Dezimalbeträge gerundet werden. In diesem Programm, das nur positive Beträge akzeptiert, wird durch Addition von 50 vor Division durch 100 auf ganze Zahlen gerundet. Wenn Sie 99 für 1 und einen Abzinsungssatz von 50 eingeben, können Sie das Schnittergebnis von 49 und das Rundungsergebnis von 50 vergleichen.
