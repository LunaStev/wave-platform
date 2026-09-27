---
translation_set_id: learn-strings
path: language/strings
locale: de
group: language
group_order: 2
order: 7
title: 7. Zeichenfolgen, Zeichen und Bytes
summary: Unterscheiden Sie zwischen String und char, UTF-8 Bytelänge, NUL, Such- und Binärdaten.
---

## Buchstaben auf dem Bildschirm und Bytes im Speicher

Sie sehen Buchstaben auf dem Bildschirm, aber Bytes werden im Speicher gespeichert. Besonders in Fällen, in denen ein Zeichen mehrere UTF-8 Bytes hat, wie zum Beispiel im Koreanischen, kann es leicht zu Fehlern kommen, wenn „Länge“ und „Anzahl der Zeichen“ austauschbar verwendet werden.

In diesem Kapitel wird zwischen str und char, Endung NUL, escape, Suchort und Binärdaten unterschieden. Bei den Beispielen handelt es sich jeweils um ein vollständiges Programm.

## String-Literale und Ausgabe

<!-- wave-example: book-string-literals -->
```wave
fun main() {
    var greeting: str = "안녕하세요";

    println("{}", greeting);
    println("line one\nline two");
    println("quote: \"Wave\"");
}
```

Ausführungsergebnis:

```text
안녕하세요
line one
line two
quote: "Wave"
```

Reguläre Zeichen in doppelten Anführungszeichen werden als UTF-8 ausgedrückt. escape weist auf Bytes hin, die nur schwer direkt aus der Quelle geschrieben werden können. `\n` ist ein Zeilenumbruchbyte und druckt nicht zwei Zeichen, einen Backslash und n. Um den Backslash selbst zu drucken, verwenden Sie `\\`.

## Die Länge ist die Anzahl der Bytes

<!-- wave-example: book-utf8-length -->
```wave
import("std::string::len")::{
    len
};

fun main() {
    println("ASCII={}", len("Wave"));
    println("Korean={}", len("한"));
    println("mixed={}", len("Wave한"));
}
```

Ausführungsergebnis:

```text
ASCII=4
Korean=3
mixed=7
```

`len` zählt Bytes vor dem abschließenden NUL, nicht sichtbare Zeichen. Das Zeichen `한` benötigt in UTF-8 drei Bytes. Zeichenzahlen, Unicode-Codepunktzahlen und Bytezahlen sind im Allgemeinen nicht austauschbar. Die Anzeigebreite hängt auch von Faktoren wie Schriftarten und der Kombination von Zeichen ab.

Daher müssen Funktionen, die eine Zeichenfolge an einer beliebigen Byte-Position ausschneiden und auf dem Bildschirm anzeigen, die Unicode-Grenze separat berücksichtigen. Definieren Sie klar die Eingabebedingungen, ob es sich um ein Programm handelt, das nur ASCII oder allgemeinen Unicode-Text verarbeitet.

## NUL Ende und Länge

str verwendet ein 0-Byte, um das Ende anzuzeigen. len berücksichtigt nicht das letzte Byte in der Länge. Es ist ein Fehler, NUL in ein String-Literal einzufügen. Hier ist ein Beispiel für einen absichtlichen Fehler:

```wave
fun main() {
    var text: str = "left\x00right";
}
```

`\xNN` in der Quelle gibt ein Byte mit genau zwei Hexadezimalziffern an. `\x41` stellt das Byte A dar, das 65 ist. Da das Schreiben eines regulären Zeichens als UTF-8 und das Einfügen eines beliebigen Bytes unterschiedlich sind, sind nicht alle str, die als `\xNN` geschrieben werden können, gültig UTF-8.

## char enthält nicht das gesamte Zeichen Unicode

char ist ein vorzeichenloser 8-Bit-Zeichenwert. Sie können Literale verwenden, die Werte in einem einzelnen Bytebereich darstellen, z. B. `'A'`. `'한'` ist ein Fehler, da es nicht in diesen Bereich fällt. `"한"` ist ein separates str mit mehreren UTF-8 Bytes.

Versuchen Sie nicht, immer einen Buchstaben in einen char zu setzen. Sie müssen zunächst entscheiden, ob die zur Textverarbeitung erforderlichen Einheiten Bytes oder Unicode-Codepunkte sind.

## String-Vergleich

Um String-Inhalte zu vergleichen, verwenden Sie die Funktion std. Nachfolgend finden Sie ein Programm, das nach identischen Inhalten und Groß-/Kleinschreibungsunterschieden sucht.

<!-- wave-example: book-string-compare -->
```wave
import("std::string::cmp")::{
    eq, starts_with, ends_with
};

fun main() {
    if (eq("Wave", "Wave")) {
        println("same bytes");
    }

    if (!eq("Wave", "wave")) {
        println("case differs");
    }

    if (starts_with("report.txt", "report") && ends_with("report.txt", ".txt")) {
        println("name matches");
    }
}
```

Ausführungsergebnis:

```text
same bytes
case differs
name matches
```

Dieser Vergleich vergleicht Byte-Strings. Eine sprachspezifische Groß-/Kleinschreibung oder Unicode-Normalisierung wird nicht automatisch durchgeführt. Auch beim Vergleich von Dateinamen sind die Dateinamen-Gleichheitsregeln in OS nicht mit einfachen Zeichenfolgenvergleichen identisch.

## Einheiten und Ausfälle von Suchergebnissen

<!-- wave-example: book-string-search -->
```wave
import("std::string::find")::{
    find, count
};

fun main() {
    var first: i32 = find("banana", "na");
    var missing: i32 = find("banana", "xy");
    var matches: i32 = count("aaaa", "aa");

    println("first={} missing={}", first, missing);
    println("matches={}", matches);
}
```

Ausführungsergebnis:

```text
first=2 missing=-1
matches=2
```

find gibt die erste Position oder -1 zurück. Index 0 ist ebenfalls ein Erfolg, daher wird er mit `result >= 0` überprüft. count ist keine Position, sondern die Anzahl der nicht überlappenden Übereinstimmungen. Oben ist aa die Nummer 2, da sie mit 0~1 und 2~3 übereinstimmt.

Auch die Tonne needle ist Vertragsbestandteil. find gibt 0 zurück, contains gibt true zurück und count gibt 0 zurück. Gehen Sie nicht davon aus, dass nur weil sich der Funktionsname im selben Modul befindet, auch die Rückgabemethode dieselbe ist.

## Das Entfernen von Leerzeichen unterscheidet sich vom Erstellen einer neuen Zeichenfolge

trim_range gibt den Bereich ohne Leerzeichen zurück, ohne den Originaltext zu ändern oder zu kopieren. Da wir einen Ausgabezeiger erhalten, bereiten wir zunächst eine Ganzzahl vor, um das Ergebnis zu speichern.

<!-- wave-example: book-string-trim-range -->
```wave
import("std::string::trim")::{
    trim_range
};

fun main() {
    var start: i32 = 0;
    var end: i32 = 0;

    trim_range("  Wave  ", &start, &end);

    println("start={} end={} bytes={}", start, end, end - start);
}
```

Ausführungsergebnis:

```text
start=2 end=6 bytes=4
```

Der Bereich ist `[start, end)`. Es umfasst den Anfang, aber nicht das Ende, daher beträgt seine Länge end-start. Durch das Hinzufügen von start zur Startadresse des Originals wird nicht automatisch NUL am Standort end erstellt. Sie müssen den Bereich separat transportieren oder einen neuen String-Bereich vorbereiten.

## Binärdaten haben eine eigene Länge

Daten, die Nullen enthalten, fallen nicht unter die End-of-String-Regeln. Es verwendet Byte-Arrays und Längen.

<!-- wave-example: book-binary-not-string -->
```wave
fun main() {
    var bytes: array<u8, 3> = [65, 0, 66];

    for (var index: i32 = 0; index < 3; index += 1) {
        println("{}", bytes[index]);
    }
}
```

Ausführungsergebnis:

```text
65
0
66
```

Die zweite 0 sind die tatsächlichen Daten. Wenn Sie dies als str interpretieren, wird es so behandelt, als würde es bei der ersten 0 enden, und Sie können die folgende 66 nicht sehen. Wenn Sie umgekehrt ein Array ohne NUL in str mit cast ändern, besteht die Gefahr, dass über das Array hinaus gelesen wird. cast ist keine Operation zum Hinzufügen eines Abschlussbytes.

## Übung: Dateinamen untersuchen

Überprüfen Sie, ob der Dateiname mit `.wave` endet und ob die Zeichenfolge `test` enthält, geben Sie ihn in eine Testdatei aus. Bei dieser Übung wird nur nach Bytemustern im Namen gesucht und nicht auf die tatsächliche Dateiexistenz eingegangen.

### Vollständige Lösung

<!-- wave-example: book-string-solution -->
```wave
import("std::string::cmp")::{
    ends_with
};
import("std::string::find")::{
    contains
};

fun classify(name: str) {
    if (!ends_with(name, ".wave")) {
        println("other file");
        return;
    }

    if (contains(name, "test")) {
        println("Wave test file");
    } else {
        println("Wave source file");
    }
}

fun main() {
    classify("main.wave");
    classify("parser_test.wave");
    classify("notes.txt");
}
```

Ausführungsergebnis:

```text
Wave source file
Wave test file
other file
```

Es gibt separate Richtlinien zum Umgang mit dem Großbuchstaben `.WAVE` und zur Frage, ob er als Test betrachtet werden soll, auch wenn test im gesamten Pfad enthalten ist. Selbst wenn es sich um eine kleine Funktion handelt, kann ihre Funktionsweise nur dann genau beschrieben werden, wenn bestimmt wird, auf welche Eingabe sie abzielt.
