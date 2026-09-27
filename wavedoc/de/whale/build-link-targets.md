---
translation_set_id: build-link-targets
path: whale/build-link-targets
locale: de
group: whale
group_order: 1
order: 5
title: Build-, Link- und Zieloptionen
summary: Beschreibt emit-Liefergegenstände, Eingabetypen, Links, target/CPU/ABI und den freistehenden Bauplan.
---

## emit Ausgabe

```shell
wavec build main.wave --emit=check
wavec build main.wave --emit=ast
wavec build main.wave --emit=ir
wavec build main.wave --emit=bc
wavec build main.wave --emit=asm
wavec build main.wave --emit=obj -o main.o
wavec build main.wave --emit=bin -o app
```

artifact emit Die Typen sind `ast`, `ir`, `bc`, `asm`, `obj`, `bin`. `check` ist ein Inspektionskontrollmodus, kein lieferbarer Typ und darf nicht mit den anderen artifact emit verwendet werden.

```shell
wavec print supported-emit-kinds
```

## Eingabetyp und link-only

Zusätzlich zur Quelle Wave unterscheidet der Compiler zwischen den Eingängen IR, bitcode, assembly, object und archive. Die Supportliste wird mit folgendem Befehl abgefragt:

```shell
wavec print supported-input-types
```

Um nur die bereits erstellten object oder archive zu verknüpfen, können Sie `--input-type` und `--link-only` verwenden.

```shell
wavec build module.o --input-type=obj --link-only --emit=bin -o app
```

## nativer Link

```shell
wavec --link=m -L ./lib build main.wave
```

`--link` fügt eine Bibliothek hinzu und `-L` fügt einen Suchpfad hinzu. Selbst wenn ein Symbol in FFI deklariert ist, wird die Bibliothek, die dieses Symbol bereitstellt, nicht automatisch verknüpft.

## Ziel auswählen

Option zur Auswahl von OS und CPU zum Ausführen.

- `--target <triple>`
- `--cpu <name>`
- `--features <csv>`
- `--abi <name>`
- `--sysroot <path>`

Überprüfen Sie die Host-Standardeinstellungen und unterstützten Ziele mit dem folgenden Befehl:

```shell
wavec print host-target
wavec print supported-targets
wavec print target-spec --target <triple>
wavec print cpu-list --target <triple>
wavec print target-features --target <triple>
```

## Unterstützt von

Das Programm für Ihren aktuellen Computer wird ohne Angabe von target erstellt. Um eine andere Umgebung auszuwählen, übergeben Sie den Zielnamen unten an `--target`.

|OS·Umwelt|Architektur|Zielname|
| --- | --- | --- |
| Linux | amd64 | `x86_64-unknown-linux-gnu` |
| Linux | ARM64 | `aarch64-unknown-linux-gnu` |
| Linux | RISC-V64 | `riscv64-unknown-linux-gnu` |
| Linux | LoongArch64 | `loongarch64-unknown-linux-gnu` |
| macOS | amd64 | `x86_64-apple-darwin` |
| macOS | ARM64 | `aarch64-apple-darwin` |
| Windows | amd64 | `x86_64-pc-windows-msvc` |
| Windows | ARM64 | `aarch64-pc-windows-msvc` |
| FreeBSD | amd64 | `x86_64-unknown-freebsd` |
| WebAssembly |64-Bit| `wasm64-unknown-unknown` |

Für Programme, die ohne OS laufen, verwenden Sie `x86_64-unknown-none-elf`, `aarch64-unknown-none-elf` und `riscv64-unknown-none-elf`. Eine vollständige Liste der Ziele für installierte Versionen finden Sie unter `wavec print supported-targets`.

## RISC-V 64 Vertrag

Die Standardwerte für Hosted RISC-V-Ziele sind `generic-rv64`, RV64GC, `lp64d` ABI. Die Standardwerte für das Ziel Freestanding sind `generic-rv64`, RV64IMAC und `lp64`.

```shell
wavec print target-spec --target riscv64-unknown-linux-gnu --format=json
wavec print target-spec --target riscv64-unknown-none-elf --format=json
```

Unterstützt RISC-V CPU für `generic`, `generic-rv64`, `rocket-rv64`, `sifive-u74`. Feature override ist `m`, `a`, `f`, `d`, `c`, `zicsr`, `zifencei` Platzieren Sie ein Zeichen vor dem Namen und trennen Sie es durch ein Komma.

```shell
wavec build main.wave \
  --target riscv64-unknown-linux-gnu \
  --features=+m,+a,+f,-d,+c,+zicsr \
  --abi=lp64f
```

RISC-V Die Überprüfung lehnt inkonsistente Kombinationen ab. `d` erfordert `f` und `f` erfordert `zicsr`. `lp64`, `lp64f`, `lp64d` müssen mit dem von Ihnen aktivierten Gleitkomma feature übereinstimmen. Wenn Sie ABI nicht direkt angeben, leitet der Compiler ABI von feature ab.

## Freistehender Link

```shell
wavec build kernel.wave \
  --target riscv64-unknown-none-elf \
  --freestanding \
  --entry=_start \
  --linker-script=linker.ld \
  --no-start-files \
  -o kernel.elf
```

`--freestanding` passt die Build-Einstellungen an, um die Verwendung von Standardbibliotheken zu vermeiden. `--entry` spezifiziert Linker-Einträge, `--linker-script` spezifiziert Skripts und `--no-start-files` spezifiziert Ausschlüsse von Host-Startdateien.

Mit `--dry-run` können Sie den Linkplan vor der eigentlichen Ausführung prüfen.

## Hosted Cross Link

Wenn Sie ein Programm erstellen, das auf einem anderen OS·CPU ausgeführt werden soll, geben Sie den Bibliothekspfad der Zielumgebung als sysroot an. Wenn Sie einen separaten Linker verwenden, geben Sie den Pfad als `-C linker` an.

```shell
wavec build main.wave \
  --target riscv64-unknown-linux-gnu \
  --sysroot /path/to/riscv64-sysroot \
  -C linker=/path/to/target-linker \
  -o app-riscv64
```

Die mit sysroot zu verknüpfenden Bibliotheken werden mit dem ausgewählten OS·CPU·ABI abgeglichen.

## Was beim Cross Build zu überprüfen ist

- target triple befindet sich in der Supportliste Ihres Compilers
- sysroot und der Linker passen zum Ziel ABI
- Ist die Linkbibliothek für die Zielarchitektur geeignet?
- CPU feature gilt für Ziel CPU
- Überprüfen Sie bei freistehender Installation, ob das Eintragssymbol und die Speicherplatzierung mit dem Linker-Skript übereinstimmen.

## Führen Sie WebAssembly aus

wasm64-Ergebnisse werden in einer Ausführungsumgebung verwendet, die memory64 unterstützt. Module, die externe Funktionen wie Datei, Zeit und Eingabe verwenden, müssen host import entsprechend dieser Funktion verbinden.

Den Prozess zum Erstellen und Verbinden des Zielcodes finden Sie in `--dry-run`. Die eigentliche Ausführung erfolgt in der ausgewählten Ausführungsumgebung OS oder WebAssembly.
