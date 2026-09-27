---
translation_set_id: whale-cli
path: whale/whale-cli
locale: de
group: whale
group_order: 1
order: 3
title: Whale-Befehlsreferenz
summary: Beschreibt die Diagnoseausgabe Whale assembler, object wrapper und die optionalen IR-Befehle.
---

## Whale Bauen

Führen Sie im Whale-Repository Folgendes aus:

```shell
cargo build --release
```

Eine ausführbare Datei der obersten Ebene verfügt über vier Befehlsfamilien:

```text
whale asm [--amd64 | --aarch64] <input> -o <output>
whale object <input> -o <output>
whale ir <subcommand> [options]
```

## AMD64 assembler

```shell
whale asm --amd64 input.asm -o output.o
```

AMD64 assembler erhält den Pfad `.o` als Ausgabe und ELF64 relocatable enthält section, symbol und relocation Erstellen object.

Schalten Sie die detaillierte Diagnoseausgabe mit `--debug-whale` ein.

```shell
whale asm --amd64 input.asm -o output.o \
  --debug-whale --token --ast --bytes --dump-hex --stats
```

Zu den Diagnoseflags gehören `--token`, `--ast`, `--bytes`, `--dump-hex`, `--dump-bin`, `--dump-json` und `--stats`. `--trace` druckt den Verarbeitungsfortschritt.

## Object wrapper

```shell
whale object input.bin -o output.o
```

Der Befehl `object` platziert Rohbytes in einem ELF64-Abschnitt `.text` und fügt ein globales `start`-Symbol am Offset 0 hinzu. Er umschließt rohen Maschinencode in einer ELF-Objektdatei.

## Optional IR socket

Der Befehl `ir` ist nur enthalten, wenn Whale mit der Funktion `socket-cli` erstellt wird.

```shell
cargo run -p whale --features socket-cli -- ir lower program.json
cargo run -p whale --features socket-cli -- ir lower program.json -o program.wir
```

`ir lower` liest JSON von Whale socket schema, wandelt es in Whale IR um und überprüft das Modul. Der Text IR wird in den Pfad stdout bzw. `-o` ausgegeben. `--target <triple>` ersetzt die Zielzeichenfolge und `--no-verify` lässt die Validierung aus.

Um den Befehl IR verwenden zu können, muss Whale mit der Funktion `socket-cli` erstellt werden. Socket JSON Hersteller und Whale müssen dasselbe socket schema version verwenden.
