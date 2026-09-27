---
translation_set_id: vex-package-manager
path: whale/vex-package-manager
locale: de
group: whale
group_order: 1
order: 4
title: Vex-Paketmanager
summary: Beschreibt manifest-basierte Wave-Projekte, Git·Pfadabhängigkeiten, lockfile, Offline-Builds und wavec-Grenzen.
---

## Rolle

Vex ist der Paketmanager und Build-Tool für Wave. Vex arbeitet zusätzlich zu `wavec`. Vex ist für die Projektstruktur und Abhängigkeitsanalyse verantwortlich und `wavec` ist für Compiler-Flags und Kompilierungspipeline verantwortlich.

Der Befehl Vex basiert auf manifest. `vex build`, `vex check` und `vex run` erhalten absichtlich nicht die Flagge raw `wavec`.

## Erstellen Sie ein Paket

```shell
vex init
vex init --lib
```

Die Anwendung verwendet `src/main.wave` und die Bibliothek verwendet `src/lib.wave`. Die Paketstammstruktur ist wie folgt:

```text
my_project/
├── src/
│   └── main.wave
├── vex.ws
├── vex.lock
└── .vex/
    └── deps/
```

`vex.ws` wird zu manifest. Vex verwendet manifest nicht in der Erweiterung `.wson`.

```wson
{
    name = "my_project",
    version = 0.1.0,
    lib = false,
    description = "my_project Project",
    author = "unknown",
    license = "Unknown",
    dependencies = []
}
```

## Build-Befehl

```shell
vex build [--target <triple>] [--release] [--dry-run] [--locked] [--offline]
vex check [--target <triple>] [--release] [--dry-run] [--locked] [--offline]
vex run   [--target <triple>] [--release] [--dry-run] [--locked] [--offline] [-- <args...>]
```

Die Optionen für Vex sind klein gehalten. Wenn Sie eine vom Compiler abhängige Steuerung benötigen, z. B. emit, linker, CPU, ABI oder debug, verwenden Sie `wavec` direkt. Wenn Sie einen bestimmten Compiler verwenden müssen, legen Sie `VEX_WAVEC=/path/to/wavec` fest.

Fortschrittsschritte wie `Resolving`, `Fetching`, `Compiling`, `Checking`, `Running`, `Finished` werden in stderr ausgegeben und die Programmausgabe wird unter beibehalten stdout.

## Git Zentrale Abhängigkeit

Vex Abhängigkeiten werden als lokale `path` oder Git URL angegeben. Eine Abhängigkeit kann nur eine der beiden Methoden verwenden.

```wson
{
    name = "app",
    version = 0.1.0,
    dependencies = [
        { name = "local_math", path = "../local_math" },
        { name = "remote_math", git = "https://github.com/example/math.git", tag = "v0.1.0" }
    ]
}
```

Die Git-Abhängigkeit kann höchstens eines von `branch`, `tag` oder `rev` angeben. Jeder Abhängigkeitsstamm muss seinen eigenen `vex.ws` haben. Vex löst die Abhängigkeit manifest rekursiv auf, weist widersprüchliche Paketidentitäten zurück und speichert die verwalteten Git checkout in `.vex/deps/<name>`.

## lockfile Vertrag

Das Schema v2 `vex.lock` zeichnet den gesamten transitiven Abhängigkeitsgraphen und den genauen Git commit auf. Commit mit manifest. Wenn Sie dasselbe manifest und ein gültiges lockfile verwenden, wird dasselbe Abhängigkeitsdiagramm ausgewählt, ohne erneut branch oder tag zu folgen.

Befehle, die Abhängigkeiten erfordern, werden automatisch interpretiert und können mit den folgenden Befehlen vorab vorbereitet werden.

```shell
vex fetch
vex update
vex update math shared_core
```

`vex update` aktualisiert alle Git Pakete oder nur die angegebenen Pakete und betroffenen Übergangsdiagramme. Bei nicht verwandten gesperrten Paketen bleibt commit ausgewählt.

## locked und offline Arbeitsablauf

`--locked` verbietet die Erstellung und Änderung von `vex.lock`. Es schlägt fehl, wenn die Datei nicht existiert, ein nicht unterstütztes Schema hat oder nicht mit dem Diagramm manifest übereinstimmt. commit ist bereits an lockfile angeheftet und kann bei Bedarf importiert werden.

`--offline` verbietet alle Git Netzwerkoperationen. Die erforderlichen checkout und commit sollten bereits lokal vorhanden sein.

```shell
vex fetch --locked
vex build --locked --offline
```

Bei diesen beiden Befehlen handelt es sich um einen strengen CI-Workflow. Bereiten Sie ein gesperrtes commit genau dann vor, wenn das Netzwerk verfügbar ist, und kompilieren Sie es dann, ohne das Netzwerk oder lockfile zu ändern. dry-run importiert keine Abhängigkeiten und schreibt lockfile nicht neu.

## Compiler-Einstellungen und Informationen

```shell
vex info
vex setup wavec
vex setup wavec --version <version>
vex --version
```

Vex validiert das Schema `wavec` dry-run JSON vor dem eigentlichen Build. Compiler, die das erforderliche Schema nicht implementieren, werden mit einem unbekannten Plan nicht ausgeführt und lehnen ihn mit einem Kompatibilitätsfehler ab.
