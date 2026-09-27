---
translation_set_id: install
path: getting-started/install
locale: de
group: getting-started
group_order: 1
order: 2
title: Wave installieren
summary: Installiere Wave unter Linux, macOS oder Windows und führe dein erstes Programm aus.
---

## Linux und macOS

Führe den folgenden Befehl in einem Terminal aus. Er installiert Wave zusammen mit dem Paketmanager Vex.

```shell
curl -fsSL https://wave-lang.dev/install.sh | bash -s -- latest
```

Öffne nach der Installation ein neues Terminal und prüfe die Version.

```shell
wavec --version
```

## Windows

Führe die folgenden Befehle in PowerShell aus. Sie installieren Wave zusammen mit dem Paketmanager Vex.

```powershell
irm https://wave-lang.dev/install.ps1 -OutFile install.ps1
powershell -ExecutionPolicy Bypass -File .\install.ps1 -Latest
```

Öffne nach der Installation ein neues PowerShell-Fenster und prüfe die Versionen.

```powershell
wavec --version
vex --version
```

## Das erste Programm ausführen

Speichere den folgenden Code als `main.wave`.

<!-- wave-example: install-stdlib -->
```wave
import("std::string::len")::{
    len
};

fun main() {
    println("Wave: {} bytes", len("Wave"));
}
```

Führe das Programm in dem Verzeichnis aus, in dem du die Datei gespeichert hast.

```shell
wavec run main.wave
```

Ausgabe:

```text
Wave: 4 bytes
```

Falls die Standardbibliothek nicht gefunden wird, installiere sie und führe das Programm erneut aus.

```shell
wavec install std
wavec run main.wave
```

[Weiter: Dein erstes Programm](/docs/de/language/program-structure) · [Fehlerbehebung](/docs/de/reference/diagnostics)
