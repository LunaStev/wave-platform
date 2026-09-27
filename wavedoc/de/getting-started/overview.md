---
translation_set_id: overview
path: getting-started/overview
locale: de
group: getting-started
group_order: 1
order: 1
title: Wave-Dokumentation und Lernleitfaden
summary: Lernen Sie Wave Schritt für Schritt kennen, von der Installation bis hin zu praktischen Programmen, und suchen Sie nach Sprachregeln und Standardbibliotheks-APIs.
---

## Lernen Sie Wave mit diesem Leitfaden

Lernen Sie, Wave-Quellcode zu schreiben, ihn zu kompilieren, auszuführen und die Ergebnisse zu überprüfen. Wenn Sie neu in der Programmierung sind, befolgen Sie die nachstehende Reihenfolge. Wenn Sie eine andere Sprache beherrschen, führen Sie die Beispiele jedes Kapitels durch und vergleichen Sie die Regeln und Grenzfälle mit dem, was Sie bereits wissen.

## Lernpfad

|Schritt|Kapitel|Was Sie lernen werden|
| --- | --- | --- |
|Einrichtung| [Installation](/docs/de/getting-started/install) |Bereiten Sie den Compiler und die Standardbibliothek vor und überprüfen Sie, ob sie ausgeführt werden|
| 1 | [Ihr erstes Programm](/docs/de/language/program-structure) |Erstellen, prüfen und führen Sie eine Quelldatei aus und verstehen Sie Exit-Codes|
| 2 | [Variablen und Typen](/docs/de/language/declarations-and-types) |Speichern Sie Werte und wählen Sie Typen mit dem erforderlichen Bereich aus|
| 3 | [Operatoren und Konvertierungen](/docs/de/language/expressions-and-operators) |Erklären Sie die Auswertungsreihenfolge und die Ergebnisse von Typkonvertierungen|
| 4 | [Bedingungen und Schleifen](/docs/de/language/control-flow) |Verzweigen Sie auf Bedingungen und verarbeiten Sie Daten mit Schleifen|
| 5 | [Funktionen](/docs/de/language/functions-and-generics) |Extrahieren Sie wiederholte Vorgänge in Funktionen|
| 6 | [Arrays](/docs/de/language/arrays) |Greifen Sie über den Index auf Elemente zu und durchlaufen Sie ein Array|
| 7 | [Saiten](/docs/de/language/strings) |Unterscheiden Sie Zeichen von Bytes und verstehen Sie Escapezeichen und Zeichenfolgenlänge|
| 8 | [Strukturen und Varianten](/docs/de/language/structures-enums-and-aliases) |Gruppieren Sie verwandte Daten und stellen Sie Erfolg und Misserfolg dar|
| 9 | [Zeiger und Lebensdauern](/docs/de/language/explicit-memory-type-model) |Ändern Sie den ursprünglichen Wert über seine Adresse und verwalten Sie seine Lebensdauer|
| 10 | [Dynamisches Gedächtnis](/docs/de/language/allocation) |Behandeln Sie Zuordnungsfehler und geben Sie Speicher frei|
| 11 | [Module und Generika](/docs/de/language/modules-imports-and-ffi) |Teilen Sie Code auf Dateien auf und verwenden Sie Funktionen mit unterschiedlichen Typen wieder|
| 12 | [Fehlerbehandlung](/docs/de/language/errors) |Überprüfen Sie die Ergebnisse und bereinigen Sie die Ressourcen bei Fehlern|
| 13 | [Einführung in asynchronen Code](/docs/de/language/async-and-never) |Erstellen Sie eine Zukunft und warten Sie, bis sie fertig ist|

## Setzen Sie Ihr Wissen in die Praxis um

Erstellen Sie nach den Kernkapiteln einen [Eingaberechner] (/docs/de/practice/input-calculator), einen [Dateileser] (/docs/de/practice/file-reader), eine [Binärnachricht] (/docs/de/practice/binary-message) und einen [TCP-Client] (/docs/de/practice/tcp-client). Testen Sie in jedem Projekt sowohl erfolgreiche Eingabe- als auch Fehlerfälle.

## Drei Dokumentationsregisterkarten

- **Wave**: Ein geführter Sprachkurs und praktische Projekte, die der Reihe nach durchgeführt werden können.
- **[Standardbibliothek](/docs/de/stdlib)**: Die APIs, Rückgabewerte, Fehler, Eigentumsregeln und Plattformanforderungen jedes Moduls.
- **[Whale](/docs/de/whale)**: Erstellen und Verknüpfen, Paketverwaltung, Befehlsverwendung und die Low-Level-Toolchain.

Beispiele unterscheiden komplette Programme von Snippets, die in eine Funktion gehören. Führen Sie `wavec`-Befehle in einem Terminal aus und speichern Sie `wave`-Codeblöcke in `.wave`-Dateien. Eingabe und Ausgabe werden getrennt dargestellt; Beispiele, die die Standardeingabe lesen, geben an, was eingegeben werden soll.

## Wenn du nicht weiterkommst

Verwenden Sie [Fehlerbehebung](/docs/de/reference/diagnostics), um Installations-, Quellprüfungs-, Verknüpfungs- und Ausführungsprobleme zu unterscheiden. Suchen Sie nach Sprachregeln in der [Syntax-Kurzreferenz](/docs/de/reference/syntax-quick-reference), Befehlen in der [Compiler-Referenz](/docs/de/getting-started/compiler) und APIs im [Standardbibliothekshandbuch](/docs/de/reference/standard-library).
