---
translation_set_id: compiler
path: getting-started/compiler
locale: de
group: getting-started
group_order: 1
order: 3
title: Compiler-Befehlsreferenz
summary: wavec Beschreibt Befehle, Build-Pipeline, Ausgabe, Ziele, Diagnose, Abhängigkeitsverknüpfung und Toolabfragen.
---

## Befehlsmodell

`wavec` ist Compiler CLI. Es kompiliert einzelne Eingaben direkt, stellt Tools Compiler-Unterstützungsinformationen zur Verfügung und verwaltet installierte Standardbibliotheksquellen.

```text
wavec [global-options] <command> [command-options]
```

|Befehl|Benutzen|
| --- | --- |
| `wavec build <input...>` |Abhängig von den Flags führt es Inspektionen, Codegenerierung, Verknüpfungen oder Ausführungspipelines durch.|
| `wavec check <file>` |Spitzname für `build <file> --emit=check`.|
| `wavec run <file> [-- <args...>]` |Es ist ein Alias für `build <file> --run` und das Argument nach `--` wird an das Programm übergeben.|
| `wavec print <item>` |Fragt Informationen zur Ziel- und Toolchain-Unterstützung ab.|
| `wavec install std` |Installieren Sie die Standardbibliothek.|
| `wavec update std` |Aktualisieren Sie installierte Standardbibliotheken.|
| `wavec --version` |Druckt Informationen zur installierten Version.|

Eine vollständige Liste der Befehle und Optionen finden Sie unter `wavec --help`.

## Erstellen, testen und ausführen

```shell
wavec build main.wave
wavec check main.wave
wavec run main.wave -- first-argument second-argument
```

`build` erstellt standardmäßig eine ausführbare Datei. `check` stoppt nach Abschluss der Frontend-Prüfung. `run` erfordert eine Binärausgabe und kann nicht mit Shared-Library-Builds verwendet werden.

Verwenden Sie `--dry-run`, um eine Anfrage zu überprüfen und zu bestimmen, welche Schritte ausgeführt werden sollen, ohne zu kompilieren, zu verknüpfen oder auszuführen.

```shell
wavec build main.wave --target riscv64-unknown-linux-gnu --dry-run
wavec build main.wave --dry-run --error-format=json
```

Das JSON-Format ist eine stabile, einheitliche Schnittstelle, die von Build-Tools wie Vex verwendet wird.

## emit und Eingabetyp

```shell
wavec build main.wave --emit=ast
wavec build main.wave --emit=ir
wavec build main.wave --emit=bc
wavec build main.wave --emit=asm
wavec build main.wave --emit=obj -o main.o
wavec build main.wave --emit=bin -o app
```

Die Ausgabetypen emit sind `ast`, `ir`, `bc`, `asm`, `obj`, `bin`. `check` ist ein Steuermodus und muss allein verwendet werden. Sie können mehrere Ausgabetypen angeben, die eine Pipeline akzeptiert, durch Kommas getrennt.

Die Eingabetypen sind `wave`, `ir`, `bc`, `asm`, `obj`, `archive`. `--input-type=<kind>` erzwingt die Angabe des Typs aller Eingänge. Wenn Sie nur die Eingänge object oder archive verknüpfen, verwenden Sie die Binärdateien emit und `--link-only`.

```shell
wavec build module.o --input-type=obj --link-only --emit=bin -o app
```

## Ausgabeort

|Optionen|Wirkung|
| --- | --- |
| `-o <file>` |Gibt den primären Ausgabepfad an.|
| `--out-dir <dir>` |emit Platziert die Ausgabe im angegebenen Verzeichnis.|
| `--target-dir <dir>` |Gibt Zwischen- und primäre Lieferrouten an.|

## Optimierung und Diagnoseausgabe

```shell
wavec -O2 build main.wave
wavec --debug-wave=tokens,ast build main.wave
```

Die Optimierungsschritte sind `-O0`, `-O1`, `-O2`, `-O3`, `-Os`, `-Oz`, `-Ofast`. Auf `--debug-wave` können `tokens`, `ast`, `ir`, `mc`, `hex`, `all` folgen, und mehrere Schritte können mit Kommas kombiniert werden.

## nativer Link

```shell
wavec --link=m -L ./lib build main.wave
wavec build main.wave --shared -o libexample.so
wavec build main.wave --static -o app
wavec build main.wave --pie -o app
```

`--link=<lib>` fügt native Bibliotheken hinzu und `-L <path>` fügt Suchpfade hinzu. Verbindungsmodi verwenden `--shared`, `--static`, `--pie`, `--no-pie` entsprechend den Kompatibilitätsregeln.

Zu den Backend- und Linker-Steuerungsoptionen gehören:

- `--target`, `--cpu`, `--features`, `--abi`, `--sysroot`
- `-C linker=<path>` und `-C link-arg=<arg>`
- `-C link-sysroot=<path>` und `-C relocation-model=<model>`
- `-C no-default-libs`

Freistehende Ausgaben wie Kernel verwenden `--freestanding` zusammen mit den für die Umgebung geeigneten Einstellungen `--entry`, `--linker-script` und `--no-start-files`.

## Externe Paketinterpretation

```shell
wavec --dep-root .vex/deps build main.wave
wavec --dep math=/opt/wave-deps/math build main.wave
```

`--dep-root <dir>` fügt eine Wurzel hinzu, um externe `package::module` import zu finden. `--dep <name>=<path>` legt den Paketnamen auf ein Verzeichnis fest. Dies ist der Compiler-Integrationspunkt und Projekte manifest, Abhängigkeitsdownloads und lockfile werden von Vex verwaltet.

## Unterstützungsfunktionsabfrage

Tools, die den Ziel- oder Ausgabetyp verwenden, können Supportinformationen mit `wavec print` abfragen.

```shell
wavec print host-target
wavec print target-spec --format=json
wavec print supported-targets
wavec print supported-input-types
wavec print supported-emit-kinds
wavec print supported-print-items
wavec print cpu-list --target riscv64-unknown-linux-gnu
wavec print target-features --target riscv64-unknown-linux-gnu
wavec print default-linker
wavec print sysroot
wavec print std-path
wavec print dep-search-paths
```

Sie können auch Elemente wie `host`, `default-target` und `target-list` abfragen. Elemente, die eine strukturierte Ausgabe unterstützen, erhalten `--format=json`.

## Grenze zwischen Compiler und Toolchain

`wavec` ist für die Quellprüfung, Codegenerierung und Verknüpfung verantwortlich. Vex ist für die Erstellung reproduzierbarer Pakete mit Paket manifest, dem Abhängigkeitsdiagramm und lockfile verantwortlich. Whale ist eine Low-Level-Toolchain, die unabhängig läuft.

## std Pfad angeben

```shell
wavec --std-root /absolute/path/to/std check main.wave
wavec --std-root /absolute/path/to/std run main.wave
```

Der angegebene Pfad std hat Vorrang vor dem Installationspfad und schlägt fehl, wenn der Pfad falsch oder mit std nicht kompatibel ist. Es ersetzt nicht automatisch std aus anderen Installationen. Wählen Sie std aus, das Ihrem Compiler entspricht.

Für die Ausgabe `-o` wird ein anderer Pfad als die Quell-/Eingabedatei verwendet. Da `check` den Laufzeitvorgang nicht prüft, prüft [üben](/docs/de/practice/input-calculator) auch die Ausführungsergebnisse.
