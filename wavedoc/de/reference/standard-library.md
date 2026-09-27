---
translation_set_id: standard-library
path: reference/standard-library
locale: de
group: stdlib
group_order: 1
order: 1
title: Standard-Bibliotheksführer
summary: So finden Sie ein für Ihren Zweck geeignetes Modul und lesen die Fehler und Eigentumsregeln der Funktion.
---

## Finden Sie die Funktionen, die Sie benötigen

Die Standardbibliothek ist import mit dem Pfad `std::module::file`. Auch wenn die Namen ähnlich sind, können Funktionen auf unterschiedliche Weise Fehler zurückgeben. Lesen Sie zuerst [API Wie man liest](/docs/de/stdlib/contracts) und gehen Sie dann in der folgenden Tabelle zu dem Modul, das Sie benötigen.

|Was ich tun möchte|Dokument|Haupt import|
| --- | --- | --- |
|Stringlänge/Vergleich/Suche| [string](/docs/de/reference/string-and-bytes) | `std::string::len`, `cmp`, `find`, `trim` |
|Speicherzuordnung/Kopie/Größe| [mem](/docs/de/reference/memory-and-buffer) | `std::mem::alloc`, `ops`, `layout` |
|Liste von Bytes unterschiedlicher Größe| [buffer](/docs/de/stdlib/buffer) | `std::buffer::alloc`, `read`, `write` |
|Binäres Lesen/Schreiben| [bytes](/docs/de/stdlib/bytes) | `std::bytes::types`, `cursor`, `leb128` |
|Datei·Deskriptor I/O| [fs und io](/docs/de/stdlib/files-io) | `std::fs::file`, `std::io::fd` |
|Pfadkombination/Umgebungseinstellungen| [path und env](/docs/de/stdlib/path-env) | `std::path::copy`, `std::env::environ` |
|Zeitmessung/Warten| [time](/docs/de/stdlib/time) | `std::time::duration`, `clock`, `sleep` |
|Numerische Suche nach Adressen/Namenslisten| [net.resolve](/docs/de/stdlib/resolution) | `std::net::resolve`, `resolve_table` |
|TCP Verbindung/Übertragung| [net.tcp](/docs/de/stdlib/tcp) | `std::net::tcp`, `address`, `error` |
|OS Zufallszahl| [random](/docs/de/stdlib/random) | `std::random::fill` |
|Prozess·OS Grenze| [Systemfunktion](/docs/de/reference/system-io-network-process) | `std::process::core`, `spawn`, `std::sys` |
|Asynchrone Aufgabenausführung| [task](/docs/de/stdlib/task) | `std::task` |
|Mathematik-/Diagnoseassistent| [math und debug](/docs/de/stdlib/math-debug) | `std::math::int`, `float`, `std::debug::core` |

## Beispiel für den ersten Einsatz

Das folgende Programm verwendet eine Funktion von std, ohne ein separates Paket herunterzuladen. Speichern Sie es als `main.wave` und führen Sie es als `wavec run main.wave` aus.

<!-- wave-example: library-import -->
```wave
import("std::string::len")::{
    len
};

fun main() {
    println("{}", len("Wave"));
}
```

Das Ergebnis ist `4`. Um das gleiche Beispiel Schritt für Schritt zu verstehen, lesen Sie [Saiten](/docs/de/language/strings).

## Kompatibel mit Compiler std

Bestätigen Sie den ausgewählten Pfad mit `wavec print std-path`. Wenn Sie std von einer anderen Kasse aus verwenden, geben Sie den Pfad als `wavec --std-root /absolute/path/to/std check main.wave` an. Wenn der angegebene Pfad ungültig oder inkompatibel ist, wird ein Fehler angezeigt.

## Bahnsteiggrenze

Unterscheiden Sie zwischen Rechenfunktionen wie String/Byte und OS-Funktionen wie Datei/Socket. Das Erkennen eines Ziels garantiert nicht, dass alle Hosts API bereitgestellt werden. Lesen Sie die Plattformeinträge für [Unterstützungsziel](/docs/de/whale/build-link-targets) und jeweils API gemeinsam. `std::sys` ist eine Schnittstelle auf niedrigerer Ebene und tragbare Programme verwenden zuerst das Modul auf höherer Ebene.
