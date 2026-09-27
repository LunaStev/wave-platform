---
translation_set_id: stdlib-time
path: stdlib/time
locale: de
group: stdlib
group_order: 1
order: 11
title: time: Dauer, Messung und Warten
summary: Erklären Sie den Unterschied zwischen den Einheiten von Duration und realtime·monotonic clock.
---

## Duration

`Duration` in `std::time::duration` hat seconds und nanoseconds. Der normalisierte nanoseconds-Bereich liegt zwischen 0 und 999999999. 1000 Millisekunden entsprechen 1 Sekunde.

```text
time_duration_from_ms(milliseconds: i64) -> Duration
time_duration_from_seconds(seconds: i64) -> Duration
time_duration_checked_add(left: Duration, right: Duration) -> DurationResult
time_duration_to_ns(value: Duration) -> DurationValueResult
```

Die Ergebnisse der checked-Operation und der Ganzzahlkonvertierung umfassen ok und value. Das Ersetzen des breiten Duration durch einen einzelnen i64 Nanosekundenwert liegt möglicherweise außerhalb des zulässigen Bereichs. Überprüfen Sie daher zuerst ok.

## Unterscheiden Sie zwischen Messung und Vision

`std::time::clock` bis `time_now_realtime(tp: ptr<TimeSpec>) -> i64` entsprechen Kalenderzeiten. Verwenden Sie `time_now_monotonic` für die Messung der abgelaufenen Zeit, da sich diese durch Korrekturen der Systemuhr ändern kann. Der Ausgabespeicher wird vom Aufrufer bereitgestellt und liest nur sec/nsec, wenn der Status Erfolg lautet.

```text
std::time::sleep
time_sleep(duration: Duration) -> i64
time_sleep_ms(ms: i64) -> i64
time_sleep_us(us: i64) -> i64
time_sleep_ns(ns: i64) -> i64
time_sleep_seconds(seconds: i64) -> i64
```

Eine negative Wartezeit ist ein Fehler und 0 ist ein sofortiger Erfolg. Nach interruption wird nur noch die verbleibende Zeit bis monotonic deadline abgewartet. Da sich die tatsächliche Weckzeit aufgrund der Planung verzögern kann, wird sie nicht als Funktion verwendet, die eine genaue Ausführungszeit garantiert. Der synchrone Aufruf von sleep innerhalb der Aufgabe async kann den Fortschritt des Executors verhindern. Wählen Sie daher `task::sleep_ms`.

## Beispiel für die Einheitenumrechnung

<!-- wave-example: duration-api -->
```wave
import("std::time::duration")::{
    Duration, DurationValueResult, time_duration_from_ms, time_duration_to_ns
};

fun main() -> i32 {
    var duration: Duration = time_duration_from_ms(1500);
    var value: DurationValueResult = time_duration_to_ns(duration);
    if (!value.ok) {
        return 1;
    }

    println("{} {}", duration.seconds, duration.nanoseconds);
    println("{}", value.value);
    return 0;
}
```

Ausführungsergebnis:

```text
1 500000000
1500000000
```

## Fügen Sie Zeit hinzu und ändern Sie die Einheiten

750 ms und 800 ms ergeben zusammen 1 Sekunde und 550000000 Nanosekunden. Anstatt Sekunden und Nanosekunden separat zu addieren, können Sie mit checked_add Übertrag und Reichweite gemeinsam prüfen.

<!-- wave-example: book-duration-add -->
```wave
import("std::time::duration")::{
    Duration,
    DurationResult,
    DurationValueResult,
    time_duration_from_ms,
    time_duration_checked_add,
    time_duration_to_ms
};

fun main() -> i32 {
    var first: Duration = time_duration_from_ms(750);
    var second: Duration = time_duration_from_ms(800);
    var sum: DurationResult = time_duration_checked_add(first, second);

    if (!sum.ok) {
        return 1;
    }

    var milliseconds: DurationValueResult = time_duration_to_ms(sum.value);

    if (!milliseconds.ok) {
        return 2;
    }

    println("{}s {}ns", sum.value.seconds, sum.value.nanoseconds);
    println("{}ms", milliseconds.value);
    return 0;
}
```

Ausführungsergebnis:

```text
1s 550000000ns
1550ms
```

## negatives Zeitintervall

Duration steht auch für negative Zahlen. -1 ms normalisiert sich zu seconds=-1, nanoseconds=999000000. Die beiden Felder zusammen ergeben einen Wert von minus 1 Millisekunde. Sie sollten nicht allein aufgrund des Blicks auf das Feld nanoseconds beurteilen, dass es sich um eine positive Zahl handelt.

Negative Zahlen sind in Zeitintervallberechnungen gültig, aber die Übergabe einer negativen Zahl an sleep ist ein Fehler. Wenn bei der Berechnung der verbleibenden Wartezeit die Frist bereits abgelaufen ist, wird die nächste Verarbeitung ohne Wartezeit fortgesetzt.

## Uhrzeit auswählen API

|Zweck|auswählen|Was die Ergebnisse bedeuten|
| --- | --- | --- |
|Zeit, die zwischen zwei Zeitpunkten verstrichen ist| monotonic clock |Intervall unabhängig von der visuellen Korrektur des Systems|
|tatsächliche Kalenderzeit| realtime clock |Vom System festgelegte Zeit|
|Warten auf synchrones Programm| `time_sleep_ms` |Der Anrufverlauf befindet sich in der Warteschlange|
|async Warten auf Aufgabe| `task::sleep_ms` |Übergabe der Ausführungsmöglichkeit an eine andere Aufgabe|

Selbst eine Uhr, die Nanosekunden zurückgibt, bedeutet nicht, dass die tatsächliche Messgenauigkeit 1 Nanosekunde beträgt. Messen Sie beim Leistungsvergleich die Gesamtzeit für die mehrfache Wiederholung einer kurzen Aufgabe und verschieben Sie Aufgaben, die nichts mit dem Messziel zu tun haben, wie z. B. Ein-/Ausgabe, aus dem Abschnitt.
