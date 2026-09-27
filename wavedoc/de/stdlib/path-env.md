---
translation_set_id: stdlib-path-env
path: stdlib/path-env
locale: de
group: stdlib
group_order: 1
order: 12
title: path und env: Pfade und Umgebungseinstellungen
summary: Liest die Pfad- und Umgebungsvariablen im Aufrufpuffer und identifiziert Kapazitätsfehler.
---

## Pfadkombination

```text
std::path::copy
path_join2(dst: ptr<u8>, dst_cap: i32, left: str, right: str) -> i32
path_basename_copy(dst: ptr<u8>, dst_cap: i32, path: str) -> i32
path_dirname_copy(dst: ptr<u8>, dst_cap: i32, path: str) -> i32
```

Die Kapazität umfasst den letzten NUL-Platz. Ein Erfolgsergebnis hat eine beliebige Länge außer NUL, ein Fehler ist -1. Nur bei Erfolg verwenden wir das Ziel als String. Diese Funktionen arbeiten mit Pfadzeichenfolgen und prüfen weder das Vorhandensein von Dateien noch die Zugriffsberechtigungen. Das Kombinieren von Pfaden allein verhindert weder Verzeichnis-Escapes noch überprüft es die Identität tatsächlicher Dateien.

## Umgebungsvariable

```text
std::env::environ
env_get(name: str, dst: ptr<u8>, dst_cap: i64) -> i64
env_exists(name: str) -> bool
env_get_i64(name: str) -> EnvResult<i64>
env_get_i32(name: str) -> EnvResult<i32>
```

`env_get` gibt bei Erfolg die Länge ohne NUL zurück. Der Anruferpuffer muss bis NUL reichen. Ein leerer Wert unterscheidet sich von einem schlüssellosen Fehler dadurch, dass er mit der Länge 0 erfolgreich sein kann.

Ermitteln und unterscheiden Sie Fehler von `std::env::consts` bis NOT_FOUND, NO_SPACE, INVALID_KEY, READ, SOURCE_INCOMPLETE, NO_MEMORY. Behandeln Sie einen Puffermangel nicht als fehlenden Schlüssel. Um Zahlen nachzuschlagen, markieren Sie ok im Ergebnis und verwenden Sie dann value. Behandeln Sie den Inhalt von Umgebungsvariablen nicht automatisch als vertrauenswürdige Einstellungen. Überprüfen Sie deren Umfang und Art.

Das folgende Beispiel kombiniert die Verzeichnisse data und die Dateinamen input.txt.

<!-- wave-example: path-api -->
```wave
import("std::path::copy")::{
    path_join2
};

fun main() -> i32 {
    var output: array<u8, 64>;
    var length: i32 = path_join2(&output[0], 64, "data", "input.txt");
    if (length < 0) {
        return 1;
    }

    println("{}", &output[0] as str);
    return 0;
}
```

Ausführungsergebnis:

```text
data/input.txt
```

## Verzeichnisse und Dateinamen aufteilen

Das folgende Beispiel kopiert einen Pfad, der in zwei Puffer aufgeteilt ist. Die Originaldatei muss nicht tatsächlich existieren.

<!-- wave-example: book-path-parts -->
```wave
import("std::path::copy")::{
    path_basename_copy,
    path_dirname_copy
};

fun main() -> i32 {
    var directory: array<u8, 64>;
    var filename: array<u8, 64>;
    var directory_length: i32 = path_dirname_copy(&directory[0], 64, "data/report.txt");
    var filename_length: i32 = path_basename_copy(&filename[0], 64, "data/report.txt");

    if (directory_length < 0 || filename_length < 0) {
        return 1;
    }

    println("directory={}", &directory[0] as str);
    println("filename={}", &filename[0] as str);
    return 0;
}
```

Ausführungsergebnis:

```text
directory=data
filename=report.txt
```

Beide Puffer bleiben bis zum Ende von main in Kraft. `as str` liest denselben Puffer als String, ohne einen neuen String zuzuweisen. Wenn Sie den Puffer ändern, ändert sich daher auch die an diese Adresse gelesene Zeichenfolge.

## Legen Sie Standardeinstellungen fest

Umgebungsvariablen sind Einstellungen, die außerhalb des Programms übergeben werden. Aktivieren Sie beim Lesen einer numerischen Einstellung die Option „Kann sie als Ganzzahl gelesen werden?“ und „Liegt es innerhalb des von diesem Programm zugelassenen Bereichs?“

<!-- wave-example: book-env-setting -->
```wave
import("std::env::environ")::{
    EnvResult,
    env_get_i32
};

fun main() -> i32 {
    var setting: EnvResult<i32> = env_get_i32("WAVE_EXAMPLE_WORKERS");
    var workers: i32 = 4;

    if (setting.ok) {
        if (setting.value < 1 || setting.value > 32) {
            println("workers must be between 1 and 32");
            return 1;
        }

        workers = setting.value;
    }

    println("workers={}", workers);
    return 0;
}
```

Wenn `WAVE_EXAMPLE_WORKERS` nicht vorhanden ist oder nicht als Ganzzahl gelesen werden kann, wird der Standardwert 4 verwendet. Wenn eine Ganzzahl zwischen 1 und 32 festgelegt ist, wird dieser Wert verwendet. Wenn es sich um eine Ganzzahl außerhalb des Bereichs handelt, wird der Vorgang mit einem Fehler beendet.

Linux/macOS:

```shell
WAVE_EXAMPLE_WORKERS=8 wavec run main.wave
```

PowerShell:

```powershell
$env:WAVE_EXAMPLE_WORKERS = "8"
wavec run main.wave
```

In beiden Fällen wird `workers=8` ausgegeben. Im obigen Beispiel wird eine einfache Standardrichtlinie ausgewählt. Wenn dies eine erforderliche Einstellung ist, behandeln Sie numerische Suchfehler als Fehler, anstatt sie durch Standardwerte zu ersetzen. Wenn Sie zwischen fehlendem Schlüssel, unzureichendem Puffer und Lesefehler unterscheiden müssen, verwenden Sie die Konstanten env_get und ENV_ERR_*.
