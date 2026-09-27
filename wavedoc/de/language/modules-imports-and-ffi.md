---
translation_set_id: modules-ffi
path: language/modules-imports-and-ffi
locale: de
group: language
group_order: 2
order: 11
title: 11. Module und generischer Code
summary: Erfahren Sie, wie Sie öffentliche Namen und explizite Typargumente abrufen.
---

## Gründe für die Aufteilung von Dateien

Wenn das Programm wächst, ist es einfacher, es zu finden, indem man verwandte Funktionen gruppiert, anstatt alle Funktionen in main.wave zu platzieren. Modulgrenzen bestimmen, welche Namen für anderen Code verfügbar gemacht werden. Generics sind ein Tool, das dieselbe Aufgabe mit unterschiedlichen Typen wiederverwendet, unabhängig von der Dateitrennung.

In diesem Kapitel erstellen wir ein Programm mit zwei Dateien und lernen Modulaliase, Auswahl import sowie generische Funktionen und Strukturen kennen.

## Zwei-Dateien-Programm

Erstellen Sie helpers.wave und main.wave im selben Verzeichnis.

helpers.wave Alle:

```wave
pub fun double(value: i32) -> i32 {
    return value * 2;
}
```

main.wave Alle:

<!-- wave-example: book-module-two-files -->
```wave
import("./helpers")::{
    double
};

fun main() {
    var result: i32 = double(21);

    println("{}", result);
}
```

Ausführungsergebnis:

```text
42
```

Führen Sie `wavec run main.wave` im Terminal aus. helpers.wave wird auch nicht separat ausgeführt. Die benötigte Quelle wird über import angeschlossen.

Das pub vor der Funktion in helpers zeigt an, dass sie von anderen Modulen importiert werden kann. Hilfsfunktionen, die nicht der Außenwelt zugänglich gemacht werden müssen, müssen nicht öffentlich gemacht werden. Selbst wenn Sie die interne Implementierung eines Moduls ändern, können Sie Änderungen am verwendeten Code reduzieren, indem Sie den Vertrag der öffentlichen Funktion beibehalten.

## Basislinie für relative Pfade

`./helpers` ist relativ zum Verzeichnis der Quelldatei, die den Satz import erstellt hat. Wenn Sie das Programm ausführen, trennen Sie es vom Arbeitsverzeichnis, in das die Datei I/O geschrieben wird. Die Schritte zum Suchen der Datei import und die Schritte zum Suchen von input.txt während der Ausführung sind unterschiedlich.

Für lokales import kann die Erweiterung `.wave` weggelassen werden. Wenn Sie das Verzeichnis geteilt haben, schreiben Sie den Pfad `./module` entsprechend dem Speicherort. Verwechseln Sie lokale relative Pfade nicht mit Pfaden, die die Namen von Paketabhängigkeiten abrufen.

## Wählen Sie import und Alias

Die Auswahl import bewirkt, dass nur die gewünschten öffentlichen Namen direkt in der aktuellen Datei verwendet werden. Wenn ein Namenskonflikt besteht oder Sie offenlegen möchten, zu welchem ​​Modul die Funktion gehört, verwenden Sie einen Alias.

<!-- wave-example: book-module-alias -->
```wave
import("std::string::len" as strings);

fun main() {
    var length: i32 = strings::len("Wave");

    println("{}", length);
}
```

Ausführungsergebnis:

```text
4
```

strings ist der in dieser Datei definierte Modulalias. `strings::len` verwendet den Namen des Moduls. Der Punkt für den Feldzugriff und `::` für die Modultrennung sind unterschiedliche Schreibweisen.

Verwenden Sie die Option import und den Alias import nicht zusammen in einem Satz. Unabhängig von Ihrem Stil sollten Sie ihn in der gesamten Datei einheitlich verwenden, damit der Ursprung des Namens leicht lesbar ist.

## Standardbibliotheken und -pakete

Der Pfad `std::` zeigt auf die Standardbibliothek. Der Benutzer veröffentlicht API und import die erforderlichen Module. Nicht alle Standardbibliotheksfunktionen werden automatisch im aktuellen Namespace platziert.

Externe Paketpfade beginnen mit dem Paketnamen. Der Speicherort des Pakets wird durch Compileroptionen oder den Paketmanager bereitgestellt. Lernen Sie zunächst die Grenzen mit lokalen Modulen kennen und lernen Sie dann in [Vex Verwendung](/docs/de/whale/vex-package-manager), wie Sie Abhängigkeiten verwalten.

## Verwenden Sie für jeden Typ dieselbe Funktion

Die folgenden Funktionen geben ihre Eingabe wörtlich zurück: Verwenden Sie für i32 und str den Typparameter T, um zu vermeiden, dass derselbe Code zweimal geschrieben wird.

<!-- wave-example: book-generic-identity -->
```wave
fun identity<T>(value: T) -> T {
    return value;
}

fun main() {
    var number: i32 = identity<i32>(42);
    var text: str = identity<str>("Wave");

    println("{} {}", number, text);
}
```

Ausführungsergebnis:

```text
42 Wave
```

T ist ein Ort zur Eingabe des Typs, nicht des während der Ausführung übergebenen ganzzahligen Werts. Rufen Sie es auf, indem Sie das Typargument wie `<i32>` angeben. Normale generische Benutzerfunktionen lassen Typargumente nicht aus.

identity<str> dupliziert keine String-Bytes, indem es sie erneut zuweist. Gibt den Wert unverändert zurück. Die Grammatik der Generika ändert nichts an den Kopier- und Eigentumsregeln von Daten.

## Von der generischen Stelle erforderliche Operationen

Nur weil es einen Typparameter gibt, bedeutet das nicht, dass alle Operationen für alle Typen verwendet werden können. minimum unten sollte als tatsächlicher Typ verwendet werden, der verglichen werden kann.

<!-- wave-example: book-generic-minimum -->
```wave
fun minimum<T>(left: T, right: T) -> T {
    if (left < right) {
        return left;
    }

    return right;
}

fun main() {
    var small: i32 = minimum<i32>(7, 4);
    var wide: i64 = minimum<i64>(100, 20);

    println("{} {}", small, wide);
}
```

Ausführungsergebnis:

```text
4 20
```

Wenn Sie das Typargument ändern, müssen `<` und die im Text verwendete Rückgabe vom entsprechenden Typ sein. Überprüfen Sie beim Lesen allgemeiner Fehler sowohl die aufgerufene Typkombination als auch die vom Funktionskörper geforderte Operation.

## generische Struktur

Erstellen wir Pair, das zwei verschiedene Werte bindet.

<!-- wave-example: book-generic-pair -->
```wave
struct Pair<A, B> {
    first: A;
    second: B;
}

fun main() {
    var item: Pair<i32, str> = Pair<i32, str> {
        first: 7,
        second: "seven"
    };

    println("{} {}", item.first, item.second);
}
```

Ausführungsergebnis:

```text
7 seven
```

Pair<i32, str> und Pair<i64, str> sind verschiedene spezifische Typen. Auch die Reihenfolge der Typargumente ist aussagekräftig. Lesen Sie den Deklarations- und Generierungscode, um zu sehen, wo die Typen von first und second bestimmt werden.

## Name und Vertrag von API veröffentlicht

Wenn Sie eine Funktion veröffentlichen, geben Sie nicht nur den Namen an, sondern auch die Eingabeeinheit, den Rückgabewert, den Fehler und den Besitz. Die Schleife, die der Aufrufer schreibt, hängt beispielsweise davon ab, ob read die maximale Länge oder die genaue Länge liest.

pub ist der öffentliche Bereich zwischen den Modulen Wave. Dies unterscheidet sich von export (c), das externe Symbole exportiert, damit andere Sprachen sie aufrufen können. Ein vollständiges Beispiel zur Verknüpfung der beiden Sprachen finden Sie unter [Siehe FFI](/docs/de/language/modules-imports-and-ffi).

## Übung und vollständige Lösung

Erstellen Sie eine öffentliche Funktion square auf math.wave und rufen Sie sie als Alias auf main.wave auf, um die Potenz von 3 und 5 auszugeben.

math.wave:

```wave
pub fun square(value: i32) -> i32 {
    return value * value;
}
```

main.wave:

<!-- wave-example: book-module-solution -->
```wave
import("./math" as math);

fun main() {
    println("{} {}", math::square(3), math::square(5));
}
```

Ausführungsergebnis:

```text
9 25
```

Wenn der Fehler darin besteht, dass die Funktion verschwunden ist, überprüfen Sie zuerst den Pfad import und pub. Wenn ein Namenskonflikt vorliegt, prüfen Sie, ob der Anruf einen Alias ​​hat. Wenn es sich um einen Typfehler handelt, überprüfen Sie die Eingabe der Funktion und den Typ des übergebenen Arguments. Versuchen Sie nicht, verschiedene Probleme mit einer Pfadkorrektur zu lösen.


## C Importfunktion

```wave
extern(c) fun puts(text: ptr<i8>) -> i32;
```

Dem Namen ABI kann der eigentliche Symbolname als String folgen.

```wave
extern(c, "native_symbol") fun local_name(value: i32) -> i32;
```

## Wave Exportfunktion

```wave
export(c) fun wave_add(left: i32, right: i32) -> i32 {
    return left + right;
}
```

`extern` und `export` können als einzelne Funktionen und Blöcke verwendet werden. Die exportierte Funktion muss die spezifische Signatur ABI haben und darf daher nicht generisch sein.

## Zielbedingungseigenschaft

Zielbedingungseigenschaften können an Elemente der obersten Ebene angehängt werden.

```wave
#[target(os="linux", arch="x86_64")]
extern(c) fun platform_call(value: i32) -> i32;
```

Die Bedingungsschlüssel sind `arch`, `os`, `env`, `abi` und die Eigenschaften gelten für das nächste Element der obersten Ebene.

## Verbinden Sie sich mit der Funktion C, die Sie selbst geschrieben haben

Dieses Labor ist für die native Umgebung mit dem Compiler C gedacht. Verkettet eine Ganzzahlfunktion ohne Bibliothekszuordnung oder Zeichenfolgenverarbeitung.

`native.c`:

```c
#include <stdint.h>
int32_t native_double(int32_t value) { return value * 2; }
```

`main.wave`:

```wave
extern(c) fun native_double(value: i32) -> i32;

fun main() {
    println("{}", native_double(21));
}
```

Führen Sie Linux/macOS von einem Terminal im selben Arbeitsverzeichnis aus.

```shell
cc -c native.c -o native.o
wavec build main.wave native.o -o ffi-example
./ffi-example
```

Die erwartete Ausgabe ist `42`. Erstellen Sie in der Entwickler-Shell von MSVC in Windows object mit `cl /c native.c /Fonative.obj` und verbinden Sie sich mit `wavec build main.wave native.obj -o ffi-example.exe`. Die Quell- und Zielarchitekturen von object müssen identisch sein. In den Beispielen werden nur kleine Werte verwendet. Um einen großen Wert an die Funktion C zu übergeben, muss auch der Multiplikationsbereich der Seite C separat garantiert werden.

Lokale Dateipfade beginnen mit `./`.
