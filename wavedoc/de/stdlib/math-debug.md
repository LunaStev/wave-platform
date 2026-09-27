---
translation_set_id: stdlib-math-debug
path: stdlib/math-debug
locale: de
group: stdlib
group_order: 1
order: 14
title: math und debug: Mathematische Hilfsfunktionen und Diagnose
summary: Verwenden Sie mathematische Funktionen mit Bereichsprüfung und Diagnoseausgabe.
---

## Ganzzahlige Funktion zur Überprüfung des Bereichs

```text
std::math::int
min_i32(a: i32, b: i32) -> i32
max_i32(a: i32, b: i32) -> i32
abs_i32_checked(x: i32) -> MathResult<i32>
clamp_i32_checked(x: i32, lo: i32, hi: i32) -> MathResult<i32>
div_floor_i32_checked(a: i32, b: i32) -> MathResult<i32>
div_ceil_i32_checked(a: i32, b: i32) -> MathResult<i32>
```

`MathResult<T>` enthält `value` und `error`. Importieren Sie die Fehlerkonstanten aus `std::math::result` und verwenden Sie `value` erst nach der Prüfung `error == MATH_ERROR_NONE`. Der Betrag der kleinsten vorzeichenbehafteten Ganzzahl lässt sich nicht im selben Typ darstellen; die geprüfte Betragsfunktion meldet deshalb einen Fehler. `clamp` weist `lo > hi` zurück. Die Division prüft auf einen Divisor von null und auf Bereichsüberschreitung. `floor` und `ceil` runden anders als die Ganzzahldivision, die in Richtung null abschneidet.

## Gleitkommawerte klassifizieren

`is_nan_f64`, `is_infinite_f64` und `is_finite_f64` in `std::math::float` unterscheiden besondere Werte. Die Funktion f32 ist ebenfalls vorhanden. `float_to_bits_f64(value: f64) -> u64` ist eine Funktion zum Erhalten von Speicherbits und unterscheidet sich von der numerischen Konvertierung von `value as u64`. NaN ist nicht einmal sich selbst gleich, daher wird es nicht gegen `value == nan` geprüft.

## Diagnose

`debug_assert(condition: bool, message: str)` in `std::debug::core` wird beendet, nachdem ein falscher Zustand diagnostiziert wurde. Situationen, die normalerweise fehlschlagen, wie z. B. Benutzereingaben, werden mit dem Ergebniswert behandelt, und assert wird bei der Überprüfung interner Programmbedingungen verwendet, die erfüllt sein müssen.

Speichern Sie das Programm unten als `main.wave` und führen Sie es aus.

<!-- wave-example: math-api -->
```wave
import("std::math::int")::{
    min_i32, max_i32
};
import("std::debug::core")::{
    debug_assert
};

fun main() {
    var low: i32 = min_i32(9, 4);
    var high: i32 = max_i32(9, 4);
    debug_assert(low <= high, "invalid range");
    println("{} {}", low, high);
}
```

Ausführungsergebnis:

```text
4 9
```

## Rundungsrichtung für negative Division

Vergleichen Sie, wie man -7 durch 3 dividiert. `/` wird in Richtung 0 gekürzt, um zu -2 zu werden. floor wählt die kleinere Ganzzahl -3 und ceil wählt die größere Ganzzahl -2 aus.

<!-- wave-example: book-math-rounding -->
```wave
import("std::math::int")::{
    div_floor_i32_checked,
    div_ceil_i32_checked,
    abs_i32_checked
};
import("std::math::result")::{
    MathResult,
    MATH_ERROR_NONE,
    MATH_ERROR_OVERFLOW
};

fun main() -> i32 {
    var floor: MathResult<i32> = div_floor_i32_checked(-7, 3);
    var ceil: MathResult<i32> = div_ceil_i32_checked(-7, 3);

    if (floor.error != MATH_ERROR_NONE || ceil.error != MATH_ERROR_NONE) {
        return 1;
    }

    println("truncate={} floor={} ceil={}", -7 / 3, floor.value, ceil.value);

    var absolute: MathResult<i32> = abs_i32_checked(-2147483648);

    if (absolute.error == MATH_ERROR_OVERFLOW) {
        println("absolute value is outside i32");
    }

    return 0;
}
```

Ausführungsergebnis:

```text
truncate=-2 floor=-3 ceil=-2
absolute value is outside i32
```

Wenn Sie nur den Ergebniswert überprüfen, können Sie nicht zwischen dem im Fehlerfall einbezogenen Ersatzwert und dem tatsächlichen Berechnungsergebnis unterscheiden. Befolgen Sie zunächst die Reihenfolge der Überprüfung error. floor ist nützlich, wenn negative Koordinaten in ein Intervall einer bestimmten Größe eingegeben werden, und ceil ist nützlich, wenn die erforderliche Anzahl von Bündeln aufgerundet wird.

## Mit welchen Fehlern sollte ich umgehen?

|Situation|Fehler|Verarbeitungsbeispiel|
| --- | --- | --- |
|Durch Null dividieren| `MATH_ERROR_DIVIDE_BY_ZERO` |Nimmt erneut die Nennereingabe vor|
|Ergebnis kann nicht im Typ gespeichert werden| `MATH_ERROR_OVERFLOW` |Berechnen Sie mit breiterem Typ oder verwerfen Sie die Eingabe|
|Das Minimum ist größer als das Maximum in clamp| `MATH_ERROR_INVALID_ARGUMENT` |Einstellbereich ändern|

assert ist kein Fehlerreparaturtool. Fehler bei Benutzereingaben werden durch bedingte Anweisungen und Rückgabewerte behandelt und nach Abschluss der Berechnung werden die internen Bedingungen, die erfüllt sein müssen, mit debug_assert überprüft.
