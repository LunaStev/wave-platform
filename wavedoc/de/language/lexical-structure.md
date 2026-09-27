---
translation_set_id: lexical
path: language/lexical-structure
locale: de
group: language
group_order: 2
order: 14
title: Lexikalische Struktur
summary: Beschreibt Bezeichner, Literale, Trennzeichen, Schlüsselwörter und Typnamen.
---

## Identifikator

Bezeichner benennen Variablen, Funktionen, Typen und Felder. Der Name unterscheidet zwischen Groß- und Kleinschreibung und kann eine beliebige Kombination aus Buchstaben, Zahlen und `_` enthalten. Im ersten Buchstaben dürfen keine Zahlen verwendet werden. Die Zeichen Unicode können auch in Bezeichnern verwendet werden.

```wave
var request_count: i64 = 0;
var 이름: str = "Wave";
```

In realen Projekten empfiehlt es sich, aus Gründen der Werkzeugkompatibilität und Durchsuchbarkeit eine einheitliche Namenskonvention zu verwenden.

## Sätze und Trennzeichen

Die meisten Deklarations- und Ausdrucksanweisungen enden mit `;`. Anweisungen mit einem Körper, wie Funktionen, bedingte Anweisungen, Schleifenanweisungen und Strukturen, verwenden den Block `{ ... }`.

```wave
var answer: i32 = 42;

fun double(value: i32) -> i32 {
    return value * 2;
}
```

## wörtlich

```wave
var integer: i32 = 42;
var decimal: f64 = 3.14;
var text: str = "Wave";
var letter: char = 'W';
var enabled: bool = true;
var address: ptr<u8> = null;
```

Sie können Ganzzahlen, Gleitkommazahlen, Zeichenfolgen, Zeichen, boolesche Werte und `null`-Literale verwenden. Verwenden Sie `null` für Zeigerwerte.

## Schlüsselwörter und Typnamen

Die wichtigsten Schlüsselwörter, die in der Wave-Grammatik verwendet werden, sind wie folgt.

`pub`, `fun`, `extern`, `export`, `type`, `enum`, `static`, `var`, `deref`, `const`, `if`, `else`, `proto`, `struct`, `while`, `for`, `in`, `out`, `clobber`, `as`, `asm`, `import`, `return`, `continue`, `print`, `input`, `println`, `match`, `break`, `true`, `false`, `null`.

Zu den integrierten Typnamen gehören `bool`, `char`, `byte`, `str`, Ganzzahl- und Gleitkommatypen, `ptr` und `array`. Zeiger werden in der Form `ptr<T>` geschrieben, und Arrays fester Länge werden in der Form `array<T, N>` geschrieben.

## Zeichenfolgen und Zeichen escape

|Notation|Bedeutung|
| --- | --- |
| `\n` |LF Zeilenumbruch|
| `\r` | CR |
| `\t` |Registerkarte|
| `\\` |Backslash|
| `\"` |doppelte Anführungszeichen|
| `\xNN` |Ein Byte, angegeben als genau zwei Hexadezimalziffern|

Allgemeine Zeichenfolgenzeichen werden als UTF-8 gespeichert. Da `\xNN` ein Byte beibehält, gibt es keine Garantie dafür, dass die gesamte Zeichenfolge ein gültiges UTF-8 ist. NUL (einschließlich `\x00`) innerhalb eines String-Literals ist ein Kompilierungsfehler. Für Daten, die Nullen enthalten, verwenden Sie ein Byte-Array und eine Byte-Länge.

`char` Das Literal muss in einen 8-Bit-Wert passen. Zeichen, die diesen Bereich überschreiten, wie z. B. `'한'`, sind Fehler. Sie unterscheidet sich von der Zeichenfolge `"한"`.

LF, CRLF und CR allein in der Quelle werden jeweils als ein logischer Zeilenumbruch behandelt. Hierbei handelt es sich um Regeln zum Quellspeicherort und zur Kommentarbeendigung. Sie bedeuten nicht, dass sie die tatsächlichen Bytes der Dateidaten ändern.

Zu den weiteren Grammatiknamen gehören `variant`, `async` und `await`, und der asynchrone Wert wird als `Future<T>` ausgedrückt. Der obige unabhängige Deklarationsblock `var` ist ein Codefragment innerhalb einer Funktion.

[String-Klasse](/docs/de/language/arrays) · [Kommentare](/docs/de/language/comments)

## Beispiel für vorsätzliches Scheitern

Wenn Sie das Programm unter check ausführen, sollte ein interner NUL-Fehler auftreten. Wenn Sie 0 Bytes benötigen, verwenden Sie das Byte-Array `[97, 0, 98]`.

<!-- wave-example: reject-string-nul -->
```wave
fun main() {
    var text: str = "a\x00b";
}
```

Das folgende Zeichenliteral überschreitet ebenfalls den 8-Bit-Bereich, es handelt sich also um einen Kompilierungsfehler. Um die Zeichenfolge UTF-8 darzustellen, verwenden Sie `str` und doppelte Anführungszeichen.

<!-- wave-example: reject-wide-char -->
```wave
fun main() {
    var letter: char = '한';
}
```
