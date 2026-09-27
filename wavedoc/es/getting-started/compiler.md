---
translation_set_id: compiler
path: getting-started/compiler
locale: es
group: getting-started
group_order: 1
order: 3
title: Referencia de comandos del compilador
summary: wavec Describe comandos, canalización de compilación, resultados, objetivos, diagnósticos, vinculación de dependencias y consultas de herramientas.
---

## modelo de comando

`wavec` es el compilador CLI. Compila entradas individuales directamente, proporciona información de soporte del compilador a las herramientas y administra las fuentes de biblioteca estándar instaladas.

```text
wavec [global-options] <command> [command-options]
```

|comando|uso|
| --- | --- |
| `wavec build <input...>` |Dependiendo de las banderas, realiza inspección, generación de código, vinculación o canalizaciones de ejecución.|
| `wavec check <file>` |Apodo para `build <file> --emit=check`.|
| `wavec run <file> [-- <args...>]` |Es un alias de `build <file> --run` y el argumento posterior a `--` se pasa al programa.|
| `wavec print <item>` |Consulta información de soporte de destino y cadena de herramientas.|
| `wavec install std` |Instale la biblioteca estándar.|
| `wavec update std` |Actualice las bibliotecas estándar instaladas.|
| `wavec --version` |Imprime la información de la versión instalada.|

Puede encontrar una lista completa de comandos y opciones en `wavec --help`.

## Construir, probar y ejecutar

```shell
wavec build main.wave
wavec check main.wave
wavec run main.wave -- first-argument second-argument
```

`build` crea un archivo ejecutable de forma predeterminada. `check` se detiene después de completar la verificación del frontend. `run` requiere salida binaria y no se puede utilizar con compilaciones de bibliotecas compartidas.

Utilice `--dry-run` para verificar una solicitud y determinar qué pasos ejecutar sin compilar, vincular o ejecutar.

```shell
wavec build main.wave --target riscv64-unknown-linux-gnu --dry-run
wavec build main.wave --dry-run --error-format=json
```

El formato JSON es una interfaz unificada y estable utilizada por herramientas de compilación como Vex.

## emit y tipo de entrada

```shell
wavec build main.wave --emit=ast
wavec build main.wave --emit=ir
wavec build main.wave --emit=bc
wavec build main.wave --emit=asm
wavec build main.wave --emit=obj -o main.o
wavec build main.wave --emit=bin -o app
```

Los tipos de salida emit son `ast`, `ir`, `bc`, `asm`, `obj`, `bin`. `check` es un modo de control y debe usarse solo. Puede especificar varios tipos de salida que acepta una canalización, separados por comas.

Los tipos de entrada son `wave`, `ir`, `bc`, `asm`, `obj`, `archive`. `--input-type=<kind>` fuerza a especificar el tipo de todas las entradas. Al vincular solo las entradas object o archive, use los binarios emit y `--link-only`.

```shell
wavec build module.o --input-type=obj --link-only --emit=bin -o app
```

## ubicación de salida

|opciones|efecto|
| --- | --- |
| `-o <file>` |Especifica la ruta de salida principal.|
| `--out-dir <dir>` |emit Coloca la salida en el directorio especificado.|
| `--target-dir <dir>` |Especifica rutas de entrega intermedias y primarias.|

## Optimización y resultados de diagnóstico

```shell
wavec -O2 build main.wave
wavec --debug-wave=tokens,ast build main.wave
```

Los pasos de optimización son `-O0`, `-O1`, `-O2`, `-O3`, `-Os`, `-Oz`, `-Ofast`. `--debug-wave` puede ir seguido de `tokens`, `ast`, `ir`, `mc`, `hex`, `all`, y se pueden combinar varios pasos con comas.

## enlace nativo

```shell
wavec --link=m -L ./lib build main.wave
wavec build main.wave --shared -o libexample.so
wavec build main.wave --static -o app
wavec build main.wave --pie -o app
```

`--link=<lib>` agrega bibliotecas nativas y `-L <path>` agrega rutas de búsqueda. Los modos de enlace utilizan `--shared`, `--static`, `--pie`, `--no-pie` según las reglas de compatibilidad.

Las opciones de control de backend y vinculador incluyen:

- `--target`, `--cpu`, `--features`, `--abi`, `--sysroot`
- `-C linker=<path>` y `-C link-arg=<arg>`
- `-C link-sysroot=<path>` y `-C relocation-model=<model>`
- `-C no-default-libs`

Las salidas independientes, como los núcleos, utilizan `--freestanding` junto con las configuraciones `--entry`, `--linker-script` y `--no-start-files` apropiadas para el entorno.

## Interpretación del paquete externo

```shell
wavec --dep-root .vex/deps build main.wave
wavec --dep math=/opt/wave-deps/math build main.wave
```

`--dep-root <dir>` agrega una raíz para buscar `package::module` import externo. `--dep <name>=<path>` fija el nombre del paquete en un directorio. Este es el punto de integración del compilador y los proyectos manifest, las descargas de dependencias y lockfile son manejados por Vex.

## Consulta de función de soporte

Las herramientas que utilizan el tipo de destino o de salida pueden consultar información de soporte con `wavec print`.

```shell
wavec print host-target
wavec print target-spec --format=json
wavec print supported-targets
wavec print supported-input-types
wavec print supported-emit-kinds
wavec print supported-print-items
wavec print cpu-list --target riscv64-unknown-linux-gnu
wavec print target-features --target riscv64-unknown-linux-gnu
wavec print default-linker
wavec print sysroot
wavec print std-path
wavec print dep-search-paths
```

También puede consultar elementos como `host`, `default-target` y `target-list`. Los elementos que admiten salida estructurada reciben `--format=json`.

## Límite entre compilador y cadena de herramientas

`wavec` es responsable de la inspección del código fuente, la generación de código y la vinculación. Vex es responsable de crear paquetes reproducibles con el paquete manifest, el gráfico de dependencia y lockfile. Whale es una cadena de herramientas de bajo nivel que se ejecuta de forma independiente.

## std Especificar ruta

```shell
wavec --std-root /absolute/path/to/std check main.wave
wavec --std-root /absolute/path/to/std run main.wave
```

La ruta std especificada tiene prioridad sobre la ruta de instalación y fallará si la ruta es incorrecta o incompatible con std. No reemplaza automáticamente std de otras instalaciones. Seleccione std que corresponda a su compilador.

Para la salida `-o`, se utiliza una ruta diferente a la del archivo fuente/entrada. Dado que `check` no verifica la operación en tiempo de ejecución, [practica](/docs/es/practice/input-calculator) también verifica los resultados de la ejecución.
