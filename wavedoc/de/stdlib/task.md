---
translation_set_id: stdlib-task
path: stdlib/task
locale: de
group: stdlib
group_order: 1
order: 13
title: task: Future ausführen und Ressourcen freigeben
summary: Beschreibt den einzelnen Verbrauch, die Ausführung und den Abbruch nach der Bereinigung asynchroner Vorgänge.
---

## Basic API

Importiert als `import("std::task" as task);`.

```text
block_on<T>(future: Future<T>) -> T
spawn<T>(future: Future<T>) -> Future<T>
cancel<T>(future: Future<T>) -> bool
yield_now() -> Future<void>
sleep_ms(milliseconds: i64) -> Future<void>
cancel_and_join<T>(future: Future<T>) -> Future<void>
shutdown()
```

`task::block_on(pending)` wird ausgeführt, bis Future abgeschlossen ist und das Ergebnis zurückgibt. Der Ergebnistyp wird aus dem übergebenen Future ermittelt.

## Lebensregeln

Future ist ein einmaliger Verbrauchsgriff. Ich betrachte das Kopieren von Werten nicht als zwei unabhängige Vorgänge. Sie sind auch dafür verantwortlich, auf Ergebnisse zu warten oder mit `spawn` geplante Aufgaben abzusagen/zu organisieren. Future wurde bereits mit `await` verbraucht und `block_on` wird nicht erneut verbraucht.

Stornierungsanfragen stornieren; Dies allein garantiert nicht, dass die Bereinigung abgeschlossen ist. Warten Sie auf den erforderlichen Abschluss, bevor Sie Ressourcen freigeben. Geben Sie keinen durch asynchrone E/A geliehenen Speicher frei, solange eine Aufgabe noch darauf zugreifen kann. Rufen Sie `shutdown` auf, nachdem die Aufgaben abgeschlossen sind oder die Bereinigung des Abbruchs abgeschlossen ist.

## gemeinsame Ausführung

Lange Berechnungen und synchrone blocking-Aufrufe können den Gesamtfortschritt des Executors verlangsamen. Asynchrones Warten mit yield ergibt Ausführungsmöglichkeit. Die Blockierung I/O wird nicht asynchron, nur weil sie sich innerhalb einer Funktion async befindet.

Schauen Sie sich die Ausführungs- und Bereinigungssequenz im vollständigen Programm in [Einführung in asynchronen Code](/docs/de/language/async-and-never) an.

## Geben Sie nach und warten Sie auf die Fertigstellung

Das folgende Programm führt die Ausführung mitten in einer Operation aus und gibt ein Ergebnis zurück. Da yield kein Funktionsabbruch ist, wird der Code nach await kontinuierlich ausgeführt.

<!-- wave-example: book-task-yield -->
```wave
import("std::task" as task);

async fun calculate() -> i32 {
    println("started");
    await task::yield_now();
    println("resumed");
    return 42;
}

fun main() -> i32 {
    var pending: Future<i32> = calculate();
    var result: i32 = task::block_on(pending);

    println("result={}", result);
    task::shutdown();
    return 0;
}
```

Ausführungsergebnis:

```text
started
resumed
result=42
```

In diesem Beispiel gibt es keine anderen Vorgänge, daher ist die Ausgabereihenfolge konsistent. In einem Programm mit mehreren Aufgaben spawn können andere Aufgaben am Punkt yield fortfahren, sodass es nicht auf die Ausgabereihenfolge der verschiedenen Aufgaben ankommt.

## Reihenfolge der Organisationsaufgaben

1. Bereiten Sie den Speicherplatz und die Ressourcen vor, die Sie für Ihre Arbeit benötigen.
2. Erstellen Sie Future und führen Sie es als await, block_on oder spawn aus.
3. Wenn Sie Ergebnisse benötigen, warten Sie bis zum Abschluss.
4. Wenn Sie eine laufende Aufgabe abgebrochen haben, warten Sie, bis sie bereinigt ist.
5. Bereinigt von der Aufgabe geliehene Puffer, Dateien und Sockets.
6. Wenn keine Arbeit mehr übrig ist, rufen Sie shutdown an.

Future Den Gültigkeitsbereich einer Variablen zu verlassen ist eine Sache, und die Arbeit sicher zu bereinigen ist eine andere. Insbesondere wenn Sie die Adresse eines funktionslokalen Arrays an eine Aufgabe async übergeben, muss die Aufgabe die Verwendung dieses Arrays beenden, bevor die Funktion zurückkehrt.

## async Teilungsfunktionen und allgemeine Funktionen

Reine Berechnungen können in reguläre Funktionen zerlegt werden. Hängen Sie async an die Funktion an, die das Warten ausdrücken muss, und warten Sie mit await innerhalb dieser Funktion auf den Abschluss. Das bloße Umschließen einer synchronen Funktion, die viel Zeit in Anspruch nimmt, z. B. das Lesen einer Datei, mit der Funktion async gibt anderen Aufgaben keine Chance zur Ausführung.

Sie können vergleichen, wenn Sie Future in [asynchrones Lernen](/docs/de/language/async-and-never) erstellen und ausführen.
