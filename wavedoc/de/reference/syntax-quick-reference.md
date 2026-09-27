---
translation_set_id: quick-reference
path: reference/syntax-quick-reference
locale: de
group: reference
group_order: 5
order: 3
title: Syntax-Kurzreferenz
summary: Häufig verwendete Deklarationen, Kontrollfluss, Typen, Zeiger und FFI-Grammatik sind auf einer Seite organisiert.
---

## Erklärung

```wave
var value: i32 = 1;
var next: i32 = value + 1;
const LIMIT: i32 = 64;
static total: i64 = 0;
type Identifier = u64;
```

`var` ist die Region, `const`/`static` sind Deklarationen der obersten Ebene. Lokale Variablen deklarieren explizit ihren Typ.

## Funktionen

```wave
fun max(left: i32, right: i32) -> i32 {
    if (left > right) {
        return left;
    }

    return right;
}
```

## generisch

```wave
fun identity<T>(value: T) -> T {
    return value;
}

var value: i32 = identity<i32>(10);
```

Geben Sie beim Aufruf eines Generic ein Typargument an.

## Struktur und enum

```wave
struct Pair {
    left: i32;
    right: i32;
}

enum Result -> i32 {
    Ok = 0,
    Error,
}
```

## Bedingungen und Schleifen

```wave
if (ready) {
    println("ready");
}

while (count < 10) {
    count += 1;
}

for (var i: i32 = 0; i < 10; i += 1) {
    println("{}", i);
}

match (status) {
    Ready => {
        println("ready");
    }
    0 => {
        println("zero");
    }
    _ => {
        println("other");
    }
}
```

Die Überschriften von `if`, `while`, `for` und `match` verwenden Klammern.

## Arrays und Zeiger

```wave
var values: array<i32, 4> = [1, 2, 3, 4];
var p: ptr<i32> = &values[0];
var first: i32 = deref p;
```

## Konsolen-Ein-/Ausgabe

```wave
print("value = ");
println("{}", value);
input("{}", value);
```

Das erste Argument ist ein String-Literal. Jeder exakte `{}`-Platzhalter erfordert einen darauf folgenden Ausdruck, und das `input`-Ziel muss zuweisbar sein.

## import und öffentliche Gegenstände

```wave
import("std::string::len");
import("./helpers" as helpers);
import("math")::{
    Vector
};

pub fun add(left: i32, right: i32) -> i32 {
    return left + right;
}

pub import("./extra"):: {
    increment
};
```

Der lokale Pfad beginnt mit `./`. Der Alias ​​import gibt den Modulnamen an und der ausgewählte import importiert die erforderlichen öffentlichen Einträge in den Namensraum dieser Datei. `pub import` exportiert die ausgewählten Elemente erneut.

## FFI

```wave
extern(c) fun native_call(value: i32) -> i32;

export(c) fun wave_call(value: i32) -> i32 {
    return value + 1;
}
```

## Bedingtes Zielelement

```wave
#[target(os="linux", arch="riscv64")]
extern(c) fun platform_call(value: i32) -> i32;
```

Unterstützungsbedingungsschlüssel sind `arch`, `os`, `env`, `abi`. Eigenschaften steuern das nächste Element der obersten Ebene.

## Inline-Montage

```wave
var result: i64 = 0;
asm {
    "mv a0, a1"
    in("a1") 7
    out("a0") result
    clobber("memory")
}
```

Anweisungstext und Registernamen sind zielabhängig. Deklarieren Sie alle für den Block erforderlichen Eingänge, Ausgänge und versteckten clobber.

## Quelleninspektion

```shell
wavec build main.wave --emit=check
wavec print supported-targets
wavec print supported-emit-kinds
```

## Lern- und Beispielbereich

Beispiele für lokale Variablen und Anweisungen, die außerhalb der Funktion separat angezeigt werden, sind Codefragmente, die in den Hauptteil der Funktion eingefügt werden. Vollständige Laufbeispiele und Übungen folgen in [Wave Lernprozess](/docs/de/getting-started/overview). Detaillierte Regeln für Speicher und externe Funktionen finden Sie unter [Standardbibliothek](/docs/de/stdlib).
