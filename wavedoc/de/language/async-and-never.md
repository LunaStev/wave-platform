---
translation_set_id: language-async-and-never
path: language/async-and-never
locale: de
group: language
group_order: 2
order: 13
title: 13. Asynchrone Funktionen und Future
summary: Lernen Sie die Rolle der verzögerten Ausführung Future, await und block_on kennen.
---

## Warteaufgaben ausdrücken

Bei Wartevorgängen wie Dateien, Sockets und Timern muss zwischen fortlaufenden Berechnungen und dem Warten auf den Abschluss unterschieden werden. Die asynchrone Funktion drückt das zu vervollständigende Ergebnis als Future aus. Durch das Hinzufügen von async wird nicht automatisch ein neuer Thread erstellt oder alle synchronen Aufrufe in asynchrone geändert.

Lesen Sie dieses Kapitel nach Funktionen, Zeiger und Fehlerbehandlung. Das Beispiel ist ein natives Programm, das den Launcher `std::task` verwendet und jede Datei als `wavec run main.wave` ausführt.

## Erstellen Sie Future und erhalten Sie Ergebnisse

<!-- wave-example: book-async-basic -->
```wave
import("std::task" as task);

async fun calculate(value: i64) -> i64 {
    return value * 2;
}

fun main() {
    var pending: Future<i64> = calculate(21);
    var result: i64 = task::block_on(pending);

    task::shutdown();
    println("{}", result);
}
```

Ausführungsergebnis:

```text
42
```

i64, geschrieben in der Erklärung von calculate, ist der nach Abschluss erhaltene Wert. Das Ergebnis des Anrufs selbst ist Future<i64>. Führen Sie im normalen main Future mit block_on aus und erhalten Sie das vollständige Ergebnis.

Rufen Sie Shutdown auf, nachdem Sie alle Aufgaben bereinigt haben, um die Ressourcen des Executors freizugeben. Geben Sie keinen Puffer frei, solange er von einer Aufgabe noch verwendet wird, und ignorieren Sie nicht abgeschlossene Aufgaben nicht.

## Aufruf und Körperausführung sind unterschiedlich

Asynchrone Funktionen werden träge ausgeführt. Sie muss von einer regulären Funktion unterschieden werden, die ihren Rumpf unmittelbar nach dem Aufruf ausführt.

<!-- wave-example: book-async-lazy -->
```wave
import("std::task" as task);

static entered: i32 = 0;

async fun work() -> i32 {
    entered += 1;
    return 7;
}

fun main() {
    var pending: Future<i32> = work();

    println("before={}", entered);

    var result: i32 = task::block_on(pending);

    task::shutdown();
    println("after={} result={}", entered, result);
}
```

Ausführungsergebnis:

```text
before=0
after=1 result=7
```

Als Future erstellt wurde, war entered noch 0. Nach der Steuerung der Ausführung wird der Körper ausgeführt und wird 1. Aus diesem Grund sollten Sie die Aufgabe nicht als abgeschlossen betrachten, nur weil Future in einer Variablen gespeichert ist.

## Warten innerhalb einer asynchronen Funktion

Innerhalb der Funktion async wartet await auf den Abschluss eines weiteren Future. Das Ergebnis eines Ausdrucks await ist ein Vervollständigungswert.

<!-- wave-example: book-async-nested -->
```wave
import("std::task" as task);

async fun twice(value: i32) -> i32 {
    await task::yield_now();
    return value * 2;
}

async fun process() -> i32 {
    var value: i32 = await twice(20);
    return value + 2;
}

fun main() {
    var answer: i32 = task::block_on(process());

    task::shutdown();
    println("{}", answer);
}
```

Ausführungsergebnis:

```text
42
```

process wartet auf Future von twice und addiert dann 2 zum Abschlusswert. Sie können den Wert ähnlich wie das Ergebnis eines Aufrufs einer regulären Funktion verwenden, aber während Sie warten, können Sie die Gelegenheit zur Ausführung an eine andere Aufgabe übergeben.

yield_now bietet kooperative Ausführungsmöglichkeiten. Wenn in einer langen Berechnungsschleife auch nur ein einziges Mal versäumt wird, andere Aufgaben zu verlangsamen, kann dies zu einer Verlangsamung führen. Asynchronous CPU ist kein Gerät zum automatischen Parallelisieren und Verteilen von Berechnungen.

## Planen Sie mehrere Aufgaben

Sie können Aufgaben mit spawn planen und auf die jeweiligen Ergebnisse warten. Dies ist ein Beispiel für die Überprüfung des Endergebnisses, ohne sich auf die Zwischenausgabereihenfolge der beiden Vorgänge zu verlassen.

<!-- wave-example: book-async-spawn -->
```wave
import("std::task" as task);

async fun compute(value: i32) -> i32 {
    await task::yield_now();
    return value * 2;
}

async fun combine() -> i32 {
    var first: Future<i32> = task::spawn(compute(10));
    var second: Future<i32> = task::spawn(compute(20));

    var left: i32 = await first;
    var right: i32 = await second;

    return left + right;
}

fun main() {
    var total: i32 = task::block_on(combine());

    task::shutdown();
    println("{}", total);
}
```

Ausführungsergebnis:

```text
60
```

Jede von Ihnen geplante Aufgabe hat einen Ort, an dem sie auf Ergebnisse wartet. Anstatt eine Aufgabe zu erstellen und ihr Handle zu vergessen, müssen Sie entscheiden, wer ihre Erledigung überprüft. await Betrachten Sie die Reihenfolge und Reihenfolge, in der interne Vorgänge ausgeführt werden, nicht als gleich.

## Verbrauchen Sie Future einmal

Future wird als einzelnes Verbrauchshandle behandelt. Das Kopieren desselben Future lässt es nicht wie zwei verschiedene Aufgaben warten. Führen Sie block_on oder await nicht erneut für Future aus, das bereits abgeschlossen wurde.

Wenn Sie dasselbe Ergebnis an mehreren Stellen benötigen, speichern Sie den vollständigen Wert und geben Sie ihn gemäß den Kopier-/Freigaberegeln für diesen Wert weiter, anstatt Future mehrmals zu verwenden. Sie sollten auch prüfen, ob der Wert Zeiger oder eigene Ressourcen enthält.

## Unterschied zwischen Timer und synchronem Warten

async Wenn Sie innerhalb einer Funktion warten, können Sie `await task::sleep_ms(...)` verwenden. Der Aufruf von synchron sleep blockiert den aktuellen Ausführungsfluss, was sich auch auf den Fortschritt anderer Aufgaben im Executor auswirken kann.

Ich erwarte nicht, dass die Wartezeit genau der angeforderten Anzahl von Millisekunden entspricht. Abhängig von Ihrer Terminplanung und anderen Aufgaben kann es sein, dass Sie spät aufstehen. Bei der Implementierung von Timeouts verwenden wir clock und deadline, um die verstrichene Zeit zu messen, anstatt jedes Mal erneut auf die ursprüngliche Vollzeit zu warten.

## Pufferlebensdauer und Stornierung

Ein an asynchrone E/A übergebener Puffer muss auch dann gültig bleiben, wenn die aufrufende Funktion angehalten ist. Die Freigabe oder Neuzuweisung vor Abschluss oder Abbruch der Bereinigung kann dazu führen, dass der Vorgang eine ungültige Adresse erhält.

Es kann nicht davon ausgegangen werden, dass die Stornierungsanfrage und die Aufgabenerledigung gleichzeitig erfolgen. Überprüfen Sie die Ergebnisse der Serie cancel API und warten Sie auf den erforderlichen Abschluss, bevor Sie Ressourcen freigeben. Bitte lesen Sie [Siehe task](/docs/de/stdlib/task) für detaillierte Anrufregeln.

## häufiges Missverständnis

|denke|tatsächlich überprüfen|
| --- | --- |
|async Ich habe angerufen und es ist vorbei.|Haben Sie Future tatsächlich ausgeführt und abgeschlossen?|
|async Alle Aufrufe innerhalb der Funktion sind asynchron|Ist der aufgerufene API synchron oder asynchron?|
|Future Kopieren dupliziert eine Aufgabe|Benutzen Sie immer wieder denselben Griff?|
|Da Sie es abgebrochen haben, können Sie den Puffer sofort freigeben.|Haben Sie die Arbeit nach der Absage fertig organisiert?|
|Die Zwischenausgabereihenfolge ist immer festgelegt|Wartet es explizit nur auf die für das Ergebnis erforderliche Reihenfolge?|

## Übung und vollständige Lösung

Erstellen Sie eine Pipeline, die auf drei asynchrone Funktionen hintereinander wartet. Gibt das Ergebnis der Verdoppelung und Addition von 5 zurück.

<!-- wave-example: book-async-solution -->
```wave
import("std::task" as task);

async fun read_value() -> i32 {
    return 10;
}

async fun transform(value: i32) -> i32 {
    await task::yield_now();
    return value * 2;
}

async fun pipeline() -> i32 {
    var initial: i32 = await read_value();
    var changed: i32 = await transform(initial);

    return changed + 5;
}

fun main() {
    var result: i32 = task::block_on(pipeline());

    task::shutdown();
    println("{}", result);
}
```

Ausführungsergebnis:

```text
25
```

Bei diesem Beispiel handelt es sich bewusst um eine sequentielle Abhängigkeit. transform erfordert das Ergebnis von read_value, daher verschwindet die Beziehung nicht, wenn man einfach alles spawn ausführt. Der Ausgangspunkt des asynchronen Designs ist die Unterscheidung zwischen unabhängigen Operationen und Operationen, die Ergebnisse erfordern.

Sobald Sie die Grundlagen vervollständigt haben, stellen Sie mit [Übung zum Lesen von Dateien](/docs/de/practice/file-reader) und [TCP Üben](/docs/de/practice/tcp-client) eine Verbindung zu echten externen Ressourcen her.

## void und never

Eine reguläre Funktion, die einen Rückgabetyp weglässt, kehrt möglicherweise ohne Wert zum Aufrufpunkt zurück. Der Typ never wird als `!` geschrieben, was bedeutet, dass er nicht normal zum Notrufpunkt zurückkehrt. Ein repräsentatives Beispiel ist die Prozessbeendigungsfunktion.

Beispiel zur Veranschaulichung der Deklaration:

```text
fun log(message: str)              // 작업 후 돌아옴, 결과값 없음
fun stop(code: i32) -> !           // 정상 반환하지 않음
async fun work() -> i64            // 호출 결과는 Future<i64>
```

Es gibt keinen Versuch, never zu einem gemeinsamen gespeicherten Wert zu machen. Schreiben Sie keine Funktionen, die als nicht rückkehrend deklariert sind, um einen normalen Rückkehrpfad zu haben. Wenn vor dem Beenden eine Ressourcenbereinigung erforderlich ist, muss der Aufrufer dies zuerst tun.

## Vollständiges Beispiel einer Funktion, die nicht zurückkehrt

Wenn Sie es als main.wave speichern und ausführen, endet es mit dem Exit-Code 0 ohne Ausgabe. stop kehrt nicht zum Anrufer zurück, daher wird es als `-> !` deklariert.

<!-- wave-example: never-function -->
```wave
import("std::process::core")::{
    proc_exit
};

fun stop() -> ! {
    proc_exit(0);
}

fun main() {
    stop();
}
```
