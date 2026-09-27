---
translation_set_id: system-io
path: reference/system-io-network-process
locale: de
group: stdlib
group_order: 1
order: 15
title: Systemfunktionen und -prozesse
summary: Beschreibt die Grenze und Prozesslebensdauer der übergeordneten Schnittstellen API und OS.
---

## Dokumentation nach Funktion

Lesen Sie [fs und io](/docs/de/stdlib/files-io) zum Bearbeiten der Datei, [TCP](/docs/de/stdlib/tcp) zum Verknüpfen und [resolver](/docs/de/stdlib/resolution) zum Nachschlagen der Adresse. Nachfolgend finden Sie den Prozess und die Zugriffsregeln der unteren Ebene OS.

## Prozessgrundlagen API

```text
std::process::core
proc_exit(code: i32) -> !
proc_getpid() -> i64
proc_getppid() -> i64
proc_execve(path: str, argv: ptr<ptr<i8>>, envp: ptr<ptr<i8>>) -> i64
proc_waitpid_raw(pid: i64, status: ptr<i32>, options: i32) -> i64
```

`proc_exit` kehrt nicht zum Rufpunkt zurück. Bitte führen Sie vor dem Herunterfahren alle erforderlichen Datei-/Speicherbereinigungen durch. `proc_execve` unterscheidet sich von einer typischen untergeordneten Erstellungsfunktion, da es bei Erfolg das vorhandene Prozessabbild ersetzt. raw argv/envp muss für die NUL-Beendigung jedes Strings vorbereitet werden, wobei der null-Zeiger das Ende anzeigt.

Die Funktion spawn in `std::process::spawn` verarbeitet das Erstellungsergebnis und die Funktion „await“ verarbeitet den Exit-Status des untergeordneten Elements. Erfolgreiche Erstellung und erfolgreiche Beendigung des Programms sind zwei verschiedene Dinge. Wenn Sie eine Pipe erstellen, müssen Eltern und Kind das nicht verwendete Ende schließen, damit EOF übergeben wird. Wenn Sie darauf warten, dass das untergeordnete Element beendet wird, ohne die Capture-Pipe zu lesen, kann es sein, dass sich der Puffer füllt und aufeinander wartet.

## Portabilität und Low-Level-Ansatz

fork/exec, Dateideskriptor und Windows-Handle sind nicht dieselbe OS-Funktion. Überprüfen Sie die Unterstützung für das ausgewählte Ziel und behandeln Sie unsupported als normalen Fehlerpfad. `std::sys` ist eine OS-spezifische Schnittstelle und verwendet ihre numerischen Flags und ihr Layout nicht von anderen OS.

Wenn Sie direkt mit der externen C-Bibliothek verknüpfen, lesen Sie bitte [FFI](/docs/de/language/modules-imports-and-ffi). Es besteht keine Notwendigkeit, die Funktion libc willkürlich zu deklarieren, um die übergeordnete Funktion std API zu verwenden. Überprüfen Sie zuerst [Ziel- und Linkumgebung](/docs/de/whale/build-link-targets).
