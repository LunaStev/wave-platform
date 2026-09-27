---
translation_set_id: whale-cli
path: whale/whale-cli
locale: es
group: whale
group_order: 1
order: 3
title: Referencia de comandos de Whale
summary: Describe los Whale assembler, object wrapper, la salida de diagnóstico y los comandos opcionales IR.
---

## Whale Construir

En el repositorio Whale, ejecute:

```shell
cargo build --release
```

Un archivo ejecutable de nivel superior tiene cuatro familias de comandos:

```text
whale asm [--amd64 | --aarch64] <input> -o <output>
whale object <input> -o <output>
whale ir <subcommand> [options]
```

## AMD64 assembler

```shell
whale asm --amd64 input.asm -o output.o
```

AMD64 assembler recibe la ruta `.o` como salida, y ELF64 relocatable que contiene section, symbol y relocation Crear object.

Active la salida de diagnóstico detallada con `--debug-whale`.

```shell
whale asm --amd64 input.asm -o output.o \
  --debug-whale --token --ast --bytes --dump-hex --stats
```

Los indicadores de diagnóstico incluyen `--token`, `--ast`, `--bytes`, `--dump-hex`, `--dump-bin`, `--dump-json` y `--stats`. `--trace` imprime el progreso del procesamiento.

## Object wrapper

```shell
whale object input.bin -o output.o
```

El comando `object` coloca bytes sin procesar en una sección ELF64 `.text` y agrega un símbolo `start` global en el desplazamiento 0. Envuelve el código de máquina sin procesar en un archivo de objeto ELF.

## Opcional IR socket

El comando `ir` solo se incluye cuando Whale se construye con la función `socket-cli`.

```shell
cargo run -p whale --features socket-cli -- ir lower program.json
cargo run -p whale --features socket-cli -- ir lower program.json -o program.wir
```

`ir lower` lee JSON de Whale socket schema, lo convierte a Whale IR y verifica el módulo. El texto IR se envía a la ruta stdout o `-o`. `--target <triple>` reemplaza la cadena de destino y `--no-verify` omite la validación.

Para utilizar el comando IR, Whale debe estar integrado con la función `socket-cli`. Socket JSON productores y Whale deben utilizar el mismo socket schema version.
