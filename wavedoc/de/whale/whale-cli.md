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

## Text-IR prüfen und ausgeben

Der Standardbuild liest und prüft typed IR format 3. Speichern Sie das vollständige Beispiel aus der [IR-Referenz](ir-reference) als `answer.wir`. `print` prüft vor der kanonischen Ausgabe und bewahrt bei Fehlern eine bestehende Datei. Es führt IR nicht aus und erzeugt keinen nativen Code.

```shell
whale ir verify answer.wir
whale ir print answer.wir -o canonical.wir
```

## Optional IR socket

AST JSON mit `ir lower` benötigt die Feature `socket-cli`. Text-IR mit `verify` und `print` benötigt sie nicht.

```shell
cargo run -p whale --features socket-cli -- ir lower program.json
cargo run -p whale --features socket-cli -- ir lower program.json -o program.wir
```

`ir lower` liest JSON von Whale socket schema, wandelt es in Whale IR um und überprüft das Modul. Der Text IR wird in den Pfad stdout bzw. `-o` ausgegeben. `--target <triple>` ersetzt die Zielzeichenfolge und `--no-verify` lässt die Validierung aus.

Bauen Sie mit `socket-cli`, um `ir lower` zu nutzen. Socket-JSON-Produzenten und Whale müssen dieselbe AST schema version verwenden.


## Interpreter für skalare Ganzzahlen

Der Standardbuild führt auch skalare Ganzzahl-/Bool-IR aus. `--function @fN` ist erforderlich; wiederholtes `--arg` übergibt exakte Dezimalliterale oder true/false für Bool. `--max-steps` zählt Anweisungen und terminators, standardmäßig 1,000,000. run akzeptiert weder -o noch --no-verify. Speichern Sie die vollständige Schleife aus der [IR-Referenz](ir-reference) als `swap-loop.wir`:

```shell
whale ir run swap-loop.wir --function @f7 --arg 3 --max-steps 100
```

```text
i32 22
```
