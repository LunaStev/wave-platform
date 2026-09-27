---
translation_set_id: stdlib-contracts
path: stdlib/contracts
locale: de
group: stdlib
group_order: 1
order: 2
title: API-Dokumentation lesen: Fehler und Eigentum
summary: Verstehen Sie Argumenteinheiten, Ergebnisstrukturen, Teilerfolg und Ressourcenlebensdauer.
---

## Lesen Sie die Erklärung

Die folgende Notation beschreibt die Funktionsdeklaration und bezieht sich nicht auf die gesamte ausführbare Datei.

```text
io_read(fd: i64, buf: ptr<u8>, len: i64) -> i64
```

`fd` ist der offene Deskriptor, `buf` ist der vom Aufrufer bereitgestellte Speicher und `len` ist die Anzahl der zum Schreiben verfügbaren Bytes. `i64` bedeutet nicht, dass negative Längen gültig sind. Der Rückgabewert ist die tatsächliche Anzahl der gelesenen Bytes, nicht die Anforderungslänge, sodass nur der zurückgegebene Bereich verwendet wird.

## Der Fehlerausdruck variiert von Funktion zu Funktion

|Weg|ja|Inspektionsmethode|
| --- | --- | --- |
|Zeiger oder null| `mem_alloc` |null Speicherzugriff nach Inspektion|
|Anzahl der Bytes oder negative Zahl| `io_read` |Negativer Fehler, 0 EOF, positive Daten|
|Statuscode| `buffer_push` |Fehlerkonstantenvergleich mit `BUFFER_OK`|
|Erfolg und Wert| `NetResult<T>` |Nachdem Sie `ok` untersucht haben, verwenden Sie `value`|
|Teilfortschritt inklusive| `RandomFillResult` |Kreuzen Sie `ok`, `written`, `error` gemeinsam an.|

Es betrachtet nur die Anzahl der Fehler und vergleicht sie nicht mit Konstanten in anderen Modulen. Beispielsweise sind die Fehlernummern env und OS errno nicht dasselbe System. Der ursprüngliche Fehler in WASI sollte nicht als Linux errno interpretiert werden.

## besitzen und vermieten

- **Owned**: Wenn ein zugewiesener Speicher, eine offene Datei oder ein offener Socket erworben wird, ist er für den Aufruf der entsprechenden Freigabe/Schließung verantwortlich.
- **Borrow**: Das Byte view oder der an die Funktion übergebene Puffer bezieht sich auf vorhandenen Speicher. Wenn eine Funktion nicht angibt, dass sie Eigentümer wird, behält der Aufrufer die Kontrolle.
- **Ausgabeargument**: Übergibt einen gültigen Speicherplatz, in dem das Ergebnis in eine Funktion geschrieben werden kann, die es empfängt, z. B. `out_value: ptr<T>`. Stellen Sie sicher, dass im Vertrag festgelegt ist, dass das Ergebnis nur im Erfolgsfall gültig ist.

Das Kopieren einer Buffer-Struktur kann dazu führen, dass beide Kopien auf dieselbe Zuordnung verweisen. Geben Sie nicht jede Kopie einzeln frei. Ein geliehener Zeiger wird ungültig, nachdem die Zuweisung freigegeben oder neu zugewiesen wurde. String-Literale sind keine beschreibbaren Puffer.

## Ein Scheitern bedeutet nicht die Rückkehr zu einem früheren Zustand

`io_write_all` kann nach dem Schreiben einiger Bytes fehlschlagen. Bereits extern geschriebene Bytes werden nicht zurückgegeben. Andererseits bleiben beim Lesen von checked cursor von bytes die Position und der Ausgabewert erhalten, wenn dies fehlschlägt. Diese Unterschiede werden durch API spezifiziert.

Wenn Sie auch den Umgang mit Fehlern üben möchten, fahren Sie mit [Dateileser](/docs/de/practice/file-reader) und [Binäre Nachricht](/docs/de/practice/binary-message) fort.
