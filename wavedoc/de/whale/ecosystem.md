---
translation_set_id: ecosystem
path: whale/ecosystem
locale: de
group: whale
group_order: 1
order: 2
title: Toolchain-Komponenten
summary: Beschreibt die Rolle einer separaten Low-Level-Toolchain, Whale, und die Komponentengrenzen des Wave-Ökosystems.
---

## WhaleIran

Whale ist eine Low-Level-Toolchain, die Montage- und Zwischendarstellungen verarbeitet. Die Komponenten, die Baugruppen, Objekte, Links und Zwischendarstellungen verarbeiten, sind so konzipiert, dass sie in Wave und anderen Tools zur nativen Codegenerierung wiederverwendbar sind.

Whale ist nicht der Name für die gesamte Entwicklungsumgebung Wave. Die Verantwortlichkeiten für jedes Projekt sind wie folgt aufgeteilt:

|Projekt|Verantwortung|
| --- | --- |
| `wavec` |Wave Untersucht die Quelle und erstellt eine ausführbare Datei.|
| Vex |Verwaltet Wave-Paket, manifest, Abhängigkeitsdiagramm, lockfile und Paket-Builds.|
| Whale |Bietet unabhängige assembler-, object-, linker- und IR-Komponenten.|
| Wave `std` |Laufzeit und System API werden als Wave Quellmodule bereitgestellt.|

## Komponente

Whale workspace besteht aus vier Hauptbibliotheksbereichen:

- `assembler`: Tokenisierung, AMD64 Parsing/Encoding, section, symbol und relocation
- `object`: Objektdateimodell und ELF64 writer
- `linker`: Verbindungsschicht
- `ir`: Whale IR Typ, builder, Ausgabe, Verifizierung und optional frontend socket

Die ausführbare Datei `whale` stellt dieser Region die Befehle `asm`, `object`, `link` und `ir` zur Verfügung.

## Werkzeugrand

Das Programm Wave wird als `wavec` erstellt. Wenn Sie direkt mit Baugruppen, Objektdateien und IR arbeiten, verwenden Sie den Befehl `whale`.

Durch die Installation von Whale wird die Build-Methode von `wavec` nicht geändert. Vex verwendet `wavec`, um das Paket Wave zu erstellen, und Whale führt es direkt in einer Aufgabe aus, die Artefakte auf niedriger Ebene verarbeitet.

## Lieferverifizierung

Stellen Sie beim Verknüpfen des Artefakts Whale mit Ihrem Build-Prozess sicher, dass object format und das Ziel architecture übereinstimmen. symbol und relocation können mit unabhängigen Tools wie `readelf` und `objdump` überprüft werden. Builds, die IR socket verwenden, müssen socket schema vom selben Hersteller wie Whale verwenden.
