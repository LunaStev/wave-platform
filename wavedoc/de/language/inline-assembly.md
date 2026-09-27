---
translation_set_id: assembly
path: language/inline-assembly
locale: de
group: language
group_order: 2
order: 19
title: Inline-Montage
summary: Beschreibt den Vertrag der Befehlszeichenfolge des Blocks asm, der Operanden in/out und clobber.
---

## asm-Block

`asm` ist eine Low-Level-Syntax zum direkten Einfügen von Anweisungen der Zielarchitektur.

```wave
fun read_value() -> i64 {
    var result: i64 = 0;
    asm {
        "mov rax, 123"
        out("rax") result
    }

    return result;
}
```

String-Literale innerhalb eines Blocks werden als Assembler-Anweisungsliste übergeben.

## Eingabe und Ausgabe

```wave
var result: i64 = 0;
asm {
    "mov rax, rdi"
    in("rdi") 123
    out("rax") result
}
```

- `in("reg") expression` verbindet den Wert Wave mit dem Eingangsoperanden.
- `out("reg") target` schreibt den Ausgangswert in das zuweisbare Ziel Wave.
- Registernamen können als Zeichenfolgen oder Bezeichner geschrieben werden.

Eingabeoperanden können Variablen, Ganzzahl-/String-Literale, `&identifier`, `deref identifier` und negative Zahlen umfassen.

## clobber

Wenn ein Block einen Register- oder Speicherzustand außer einer expliziten Ausgabe ändert, wird dies in `clobber(...)` aufgezeichnet.

```wave
asm {
    "nop"
    clobber("rax", "rcx", "memory")
}
```

## Bei Verwendung prüfen

- Die Befehlssyntax muss mit der Zielarchitektur und dem Inline-Assembly-Vertrag LLVM übereinstimmen.
- Zerstören Sie nicht willkürlich Register, die gemäß der Aufrufkonvention erhalten bleiben müssen.
- Deklarieren Sie für Blöcke, die Speicher lesen oder schreiben, clobber, einschließlich `memory`.
- Wenn möglich, architekturspezifisches asm hinter einer kleinen Funktion isolieren.

Das Verhalten und die Portabilität der Inline-Assembly werden nicht allein durch den Sprachtyp garantiert.

## Lern- und Beispielbereich

[Üben Sie mit dem vollständigen Programm](/docs/de/getting-started/overview) · [Standardbibliothek](/docs/de/stdlib)
