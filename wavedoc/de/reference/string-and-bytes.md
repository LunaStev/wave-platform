---
translation_set_id: string-bytes
path: reference/string-and-bytes
locale: de
group: stdlib
group_order: 1
order: 3
title: string: Länge, Suche und Bereiche
summary: NUL Beschreibt die Byte-Einheit der Abschlusszeichenfolge API und den Rückgabewert.
---

## String-Speicher- und Argumentbedingungen

Das Argument `str` für dieses Modul muss ein zugängliches Abschlussbyte NUL sein. Länge und Suchindex werden in Bytes angegeben. Reguläre Zeichen werden als UTF-8 gespeichert, die Bytesuche erfolgt jedoch als Unicode ohne Normalisierung oder Zeichen-für-Zeichen-Aufteilung. Gehen Sie nicht davon aus, dass der zurückgegebene Index eine Zeichengrenze ist.

## Vergleichen Sie mit der Länge

```text
std::string::len
len(s: str) -> i32
is_empty(s: str) -> bool

std::string::cmp
eq(a: str, b: str) -> bool
cmp(a: str, b: str) -> i32
starts_with(s: str, prefix: str) -> bool
ends_with(s: str, suffix: str) -> bool
```

`len` schließt das letzte NUL aus. `cmp` Die Reihenfolge wird anhand des Vorzeichens des Ergebnisses beurteilt. Der Rückgabewert wird nicht als Unicode Zeichenreihenfolge oder sprachspezifische Wörterbuchsortierung interpretiert. Diese Funktionen belegen keinen Speicher und ändern ihre Eingaben nicht.

## suchen

Holen Sie sich den Namen, den Sie brauchen, z. B. `import("std::string::find")::{find, contains, count};`.

|Funktionsdeklaration|Ergebnis|
| --- | --- |
| `find(s: str, needle: str) -> i32` |Ort des ersten Spiels. -1 wenn nicht vorhanden, 0 für leer needle|
| `contains(s: str, needle: str) -> bool` |Enthalten oder nicht. Leer needle ist true|
| `count(s: str, needle: str) -> i32` |Anzahl nicht überlappender Übereinstimmungen. Bin needle ist 0|
| `find_char(s: str, c: u8) -> i32` |erste Position des Bytes oder -1|
| `rfind_char(s: str, c: u8) -> i32` |Letzte Position des Bytes oder -1|
| `contains_char(s: str, c: u8) -> bool` |Existenz dieses Bytes|
| `count_char(s: str, c: u8) -> i32` |die Anzahl der betreffenden Bytes|

`c` im Namen `*_char` ist ein Byte, kein Unicode-Codepunkt. Der NUL selbst am Ende der Zeichenfolge ist nicht im Suchziel enthalten.

## Bereich ohne Leerzeichen

```text
std::string::trim
trim_left_index(s: str) -> i32
trim_right_index(s: str) -> i32
trim_range(s: str, out_start: ptr<i32>, out_end: ptr<i32>)
```

`trim_range` schreibt den halboffenen Bereich `[start, end)` ohne Leerzeichen ASCII in das Ausgabeargument. Beide Ausgabezeiger müssen auf beschreibbare Ganzzahlen zeigen. Der ursprüngliche Text wird nicht geändert und es werden keine neuen Zeichenfolgen erstellt. Wenn alles leer ist, wird es zu einem leeren Bereich.

## Laufbeispiel

Speichern Sie es unter `main.wave` und führen Sie `wavec run main.wave` aus.

<!-- wave-example: string-api -->
```wave
import("std::string::find")::{
    find, count
};
import("std::string::trim")::{
    trim_range
};

fun main() {
    var start: i32 = 0;
    var end: i32 = 0;
    trim_range("  Wave  ", &start, &end);
    println("{} {}", start, end);
    println("{} {}", find("banana", "na"), count("aaaa", "aa"));
}
```

Ausführungsergebnis:

```text
2 6
2 2
```

## Verwandte Funktionen

Die Klassifizierung/Fallkonvertierung von `std::string::ascii` gilt für den Bereich ASCII. `djb2_32` und `fnv1a_64` von `std::string::hash` werden nicht für kryptografische Hashes oder Passwortspeicherung verwendet. Für Daten, die NUL enthalten, verwenden Sie [bytes](/docs/de/stdlib/bytes).

## Muster, die sich mit leeren Suchbegriffen überschneiden

Wenn Sie das Kantenverhalten der Suchfunktion mit tatsächlichen Werten betrachten, können Sie die Aufrufbedingungen leichter bestimmen.

<!-- wave-example: book-string-empty-search -->
```wave
import("std::string::find")::{find, contains, count};

fun main() {
    println("empty find={}", find("Wave", ""));
    println("empty count={}", count("Wave", ""));
    println("nonoverlapping={}", count("aaaa", "aa"));

    if (contains("Wave", "")) {
        println("empty needle is contained");
    }
}
```

Ausführungsergebnis:

```text
empty find=0
empty count=0
nonoverlapping=2
empty needle is contained
```

Der Erfolgswert von find, 0, ist die erste Position. Ein Erfolgswert von 0 für count ist das Ergebnis einer Regel ohne Übereinstimmung oder eines leeren Suchbegriffs. Keine zwei Werte werden gleich behandelt. Wenn das Ignorieren der Groß-/Kleinschreibung oder eine Unicode-Normalisierung erforderlich ist, müssen vor und nach dieser Bytesuche separate Richtlinien implementiert werden.

## trim Bereich in neue Zeichenfolge kopieren

Es gibt keine neue Endung NUL in dem von trim_range zurückgegebenen Bereich. Beim Kopieren an ein separates Ziel reservieren Sie Länge + 1 Platz und schreiben Sie das letzte Byte direkt als 0.

<!-- wave-example: book-trim-copy -->
```wave
import("std::string::trim")::{trim_range};

fun main() -> i32 {
    var original: str = "  Wave  ";
    var start: i32 = 0;
    var end: i32 = 0;
    var destination: array<u8, 16>;

    trim_range(original, &start, &end);

    var length: i32 = end - start;

    if (length >= 16) {
        return 1;
    }

    for (var index: i32 = 0; index < length; index += 1) {
        destination[index] = original[start + index];
    }

    destination[length] = 0;
    println("{}", &destination[0] as str);
    return 0;
}
```

Ausführungsergebnis:

```text
Wave
```

Eine Zeichenfolge der Länge 16 passt nicht in dieses Ziel. Dies liegt daran, dass Sie den letzten NUL benötigen. Selbst wenn die Länge 0 ist, führt das Schreiben von destination[0]=0 zu einer gültigen leeren Zeichenfolge. Das lokale Zielarray bleibt bis zum Ende von main bestehen, daher drucken wir darin.

## Zeichenfolge API Reihenfolge der Verwendung

Geben Sie beim Entwerfen einer Zeichenfolgen-API an, ob ihre Eingabe NUL-terminiert ist, ob Indizes Bytes zählen und ob das Ergebnis einen Quellbereich übernimmt oder eine neue Zuordnung besitzt. Ein geliehener Bereich hängt von der Lebensdauer der Quelle ab. Ein zugewiesenes Ergebnis muss angeben, wer es freigibt.

Lesen Sie [Kapitel zum String-Lernen](/docs/de/language/strings) für grundlegende Konzepte und [bytes](/docs/de/stdlib/bytes) für Daten einschließlich NUL.
