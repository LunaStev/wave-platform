---
translation_set_id: build-link-targets
path: whale/build-link-targets
locale: es
group: whale
group_order: 1
order: 5
title: Opciones de compilación, enlace y destino
summary: Describe emit entregables, tipos de entrada, enlaces, target/CPU/ABI y el plan de construcción independiente.
---

## emit Salida

```shell
wavec build main.wave --emit=check
wavec build main.wave --emit=ast
wavec build main.wave --emit=ir
wavec build main.wave --emit=bc
wavec build main.wave --emit=asm
wavec build main.wave --emit=obj -o main.o
wavec build main.wave --emit=bin -o app
```

artifact emit Los tipos son `ast`, `ir`, `bc`, `asm`, `obj`, `bin`. `check` es un modo de control de inspección, no un tipo entregable, y no debe usarse con el otro artifact emit.

```shell
wavec print supported-emit-kinds
```

## Tipo de entrada y link-only

Además de la fuente Wave, el compilador distingue entre las entradas IR, bitcode, assembly, object y archive. La lista de soporte se consulta con el siguiente comando:

```shell
wavec print supported-input-types
```

Para vincular solo los object o archive ya creados, puede usar `--input-type` y `--link-only`.

```shell
wavec build module.o --input-type=obj --link-only --emit=bin -o app
```

## enlace nativo

```shell
wavec --link=m -L ./lib build main.wave
```

`--link` agrega una biblioteca y `-L` agrega una ruta de búsqueda. Incluso si se declara un símbolo en FFI, la biblioteca que proporciona ese símbolo no se vincula automáticamente.

## Seleccionar objetivo

Opción para seleccionar OS y CPU para ejecutar.

- `--target <triple>`
- `--cpu <name>`
- `--features <csv>`
- `--abi <name>`
- `--sysroot <path>`

Verifique los valores predeterminados del host y los destinos admitidos con el siguiente comando:

```shell
wavec print host-target
wavec print supported-targets
wavec print target-spec --target <triple>
wavec print cpu-list --target <triple>
wavec print target-features --target <triple>
```

## Apoyado por

El programa para su computadora actual se creará sin especificar target. Para seleccionar un entorno diferente, pase el nombre de destino a continuación a `--target`.

|OS·Medio ambiente|arquitectura|nombre del objetivo|
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
| WebAssembly |64 bits| `wasm64-unknown-unknown` |

Para programas que se ejecutan sin OS, utilice `x86_64-unknown-none-elf`, `aarch64-unknown-none-elf` y `riscv64-unknown-none-elf`. Puede encontrar una lista completa de destinos para las versiones instaladas en `wavec print supported-targets`.

## RISC-V 64 contrato

Los valores predeterminados para los objetivos Hosted RISC-V son `generic-rv64`, RV64GC, `lp64d` ABI. Los valores predeterminados para el objetivo Freestanding son `generic-rv64`, RV64IMAC y `lp64`.

```shell
wavec print target-spec --target riscv64-unknown-linux-gnu --format=json
wavec print target-spec --target riscv64-unknown-none-elf --format=json
```

Admite RISC-V CPU para `generic`, `generic-rv64`, `rocket-rv64`, `sifive-u74`. Feature override es `m`, `a`, `f`, `d`, `c`, `zicsr`, `zifencei` Coloque un cartel delante del nombre y sepárelo con una coma.

```shell
wavec build main.wave \
  --target riscv64-unknown-linux-gnu \
  --features=+m,+a,+f,-d,+c,+zicsr \
  --abi=lp64f
```

RISC-V La verificación rechaza combinaciones inconsistentes. `d` requiere `f` y `f` requiere `zicsr`. `lp64`, `lp64f`, `lp64d` deben coincidir con el punto flotante feature que ha habilitado. Si no especifica ABI directamente, el compilador derivará ABI de feature.

## Enlace independiente

```shell
wavec build kernel.wave \
  --target riscv64-unknown-none-elf \
  --freestanding \
  --entry=_start \
  --linker-script=linker.ld \
  --no-start-files \
  -o kernel.elf
```

`--freestanding` ajusta la configuración de compilación para evitar el uso de bibliotecas predeterminadas. `--entry` especifica entradas del vinculador, `--linker-script` especifica scripts y `--no-start-files` especifica exclusiones de archivos de inicio del host.

Puede utilizar `--dry-run` para comprobar el plan de enlace antes de la ejecución real.

## Hosted Enlace cruzado

Al crear un programa para ejecutar en otro OS·CPU, especifique la ruta de la biblioteca del entorno de destino como sysroot. Si utiliza un vinculador independiente, especifique la ruta como `-C linker`.

```shell
wavec build main.wave \
  --target riscv64-unknown-linux-gnu \
  --sysroot /path/to/riscv64-sysroot \
  -C linker=/path/to/target-linker \
  -o app-riscv64
```

Las bibliotecas que se vincularán con sysroot están alineadas con el OS·CPU·ABI.

## Qué comprobar en la construcción cruzada

- target triple está en la lista de soporte de su compilador
- sysroot y el vinculador coinciden con el objetivo ABI
- ¿La biblioteca de enlaces es para la arquitectura de destino?
- CPU feature es válido para el objetivo CPU
- Si es independiente, verifique si el símbolo de entrada y la ubicación de la memoria coinciden con el script del vinculador.

## Ejecute WebAssembly

Los resultados de wasm64 se utilizan en un entorno de ejecución que admite memory64. Los módulos que utilizan funciones externas como archivo, hora y entrada deben conectarse al host import correspondiente a esa función.

El proceso de creación y conexión del código de destino se puede encontrar en `--dry-run`. La ejecución real tiene lugar en el entorno de ejecución OS o WebAssembly seleccionado.
