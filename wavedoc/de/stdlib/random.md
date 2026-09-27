---
translation_set_id: stdlib-random
path: stdlib/random
locale: de
group: stdlib
group_order: 1
order: 10
title: random: Füllt Puffer mit der Zufälligkeit des Betriebssystems
summary: OS Füllen Sie den Puffer mit Entropie und behandeln Sie Teilfehler.
---

## API

```text
std::random::fill
random_available() -> bool
random_fill(buffer: ptr<u8>, size: i64) -> RandomFillResult
```

Bei der Größe handelt es sich um eine Byteanzahl, und der Aufrufer stellt den Speicher bereit. `RandomFillResult` enthält ok, geschrieben und Fehler. Bei Erfolg entspricht geschrieben die angeforderte Länge. Bei einem Fehler identifiziert schriftlich das gültige, ausgefüllte Präfix; Verwenden Sie die verbleibenden Bytes nicht als Zufallsdaten.

`random_available` sagt Ihnen, ob die Zufallszahlenfunktion OS unterstützt wird. Der Erfolg einer einzelnen Anfrage wird anhand der Ergebnisse von random_fill überprüft. Verwendet nur OS Entropie und greift bei einem Fehler nicht auf den Zeitwert oder schwach PRNG zurück. size=0 ist auch dann erfolgreich, wenn es zusammen mit null übergeben wird. null ist ein Fehler für negative oder positive Längen.

## Laufbeispiel

<!-- wave-example: random-api -->
```wave
import("std::random::fill")::{
    RandomFillResult, random_fill
};

fun main() -> i32 {
    var bytes: array<u8, 16>;
    var result: RandomFillResult = random_fill(&bytes[0], 16);
    if (!result.ok) {
        println("random error={}", result.error);
        return 1;
    }

    println("filled={}", result.written);
    return 0;
}
```

Ausführungsergebnis:

```text
filled=16
```

Speichern Sie es als `main.wave` und führen Sie es aus. Der Byteinhalt ist jedes Mal unterschiedlich, sodass kein spezifischer Wert erwartet wird. Wenn dies fehlschlägt, überprüfen Sie die Ursache mit result.error. Anstatt zufällige Bytes buchstäblich auszugeben, verwenden Sie bei Bedarf eine separate Codierung.

## Wenn Ihre Anfrage falsch ist

Eine Null-Byte-Anfrage ist erfolgreich, da nichts geschrieben werden muss. Die Übergabe von null mit einer positiven Länge schlägt fehl, da kein Zielpuffer vorhanden ist. Das folgende Programm vergleicht diese Fälle ohne Speicherzuweisung.

<!-- wave-example: book-random-boundaries -->
```wave
import("std::random::fill")::{
    RandomFillResult,
    random_fill
};

fun main() -> i32 {
    var empty: RandomFillResult = random_fill(null, 0);

    if (!empty.ok || empty.written != 0) {
        return 1;
    }

    println("empty request succeeded");

    var invalid: RandomFillResult = random_fill(null, 16);

    if (invalid.ok) {
        return 2;
    }

    println("missing buffer rejected");
    return 0;
}
```

Ausführungsergebnis:

```text
empty request succeeded
missing buffer rejected
```

## Umgang mit teilweise gefüllten Puffern

Wenn 16 Bytes angefordert wurden, aber fehlgeschlagen sind und written=8 zurückgegeben wird, werden nur die ersten 8 Bytes gefüllt. Wenn die Aufgabe darin besteht, einen 16-Byte-Bezeichner zu erstellen, handelt es sich nicht um einen erfolgreichen Bezeichner, daher verwerfen wir das gesamte Ergebnis und melden einen Fehler. Sie sollten die restlichen 8 Bytes nicht mit 0 füllen und es dann als Erfolg betrachten.

Die Zuordnung zufälliger Bytes zu einem Ganzzahlbereich erfordert Sorgfalt. Die Anwendung von `% 10` auf gleichmäßig verteilte u8-Werte macht 0–5 wahrscheinlicher als 6–9, da 256 nicht durch 10 teilbar ist. Um diese Verzerrung zu beseitigen, lehnen Sie die Werte 250–255 ab, zeichnen Sie erneut und wenden Sie die Restoperation nur auf akzeptierte Werte an.

Die Speicherung und Lebensdauer der Zufallsbytes werden vom Aufrufer verwaltet. Bei Verwendung eines Arrays wird dieser innerhalb des Array-Bereichs verarbeitet, bei Verwendung von dynamischem Speicher wird er nach der Verwendung freigegeben. Stattdessen besitzt die Rückgabestruktur nicht den Puffer.
