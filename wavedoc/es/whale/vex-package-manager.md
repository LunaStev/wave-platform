---
translation_set_id: vex-package-manager
path: whale/vex-package-manager
locale: es
group: whale
group_order: 1
order: 4
title: Gestor de paquetes Vex
summary: Describe proyectos Wave basados en manifest, dependencias de ruta Git, lockfile, compilaciones fuera de línea y límites wavec.
---

## papel

Vex es el administrador de paquetes y la herramienta de compilación para Wave. Vex opera encima de `wavec`. Vex es responsable de la estructura del proyecto y el análisis de dependencia, y `wavec` es responsable de los indicadores del compilador y el proceso de compilación.

El comando Vex se basa en manifest. `vex build`, `vex check` y `vex run` no reciben intencionalmente la bandera raw `wavec`.

## Crear un paquete

```shell
vex init
vex init --lib
```

La aplicación usa `src/main.wave` y la biblioteca usa `src/lib.wave`. La estructura raíz del paquete es la siguiente:

```text
my_project/
├── src/
│   └── main.wave
├── vex.ws
├── vex.lock
└── .vex/
    └── deps/
```

`vex.ws` se convierte en manifest. Vex no utiliza manifest en la extensión `.wson`.

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

## comando de construcción

```shell
vex build [--target <triple>] [--release] [--dry-run] [--locked] [--offline]
vex check [--target <triple>] [--release] [--dry-run] [--locked] [--offline]
vex run   [--target <triple>] [--release] [--dry-run] [--locked] [--offline] [-- <args...>]
```

Las opciones para Vex se mantienen pequeñas. Si necesita control dependiente del compilador, como emit, linker, CPU, ABI o debug, use `wavec` directamente. Cuando necesite utilizar un compilador específico, configure `VEX_WAVEC=/path/to/wavec`.

Los pasos de progreso como `Resolving`, `Fetching`, `Compiling`, `Checking`, `Running`, `Finished` se generan en stderr y la salida del programa se mantiene en stdout.

## Git Dependencia Central

Las dependencias Vex se especifican como `path` locales o Git URL. Una dependencia sólo puede utilizar uno de los dos métodos.

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

La dependencia Git solo puede especificar como máximo uno de `branch`, `tag` o `rev`. Cada raíz de dependencia debe tener su propia `vex.ws`. Vex resuelve recursivamente la dependencia manifest, rechaza identidades de paquetes en conflicto y almacena el Git checkout administrado en `.vex/deps/<name>`.

## lockfile Contrato

El esquema v2 `vex.lock` registra todo el gráfico de dependencia transitiva y el Git commit exacto. Comprométete con manifest. Usando el mismo manifest y un lockfile válido seleccionará el mismo gráfico de dependencia sin seguir branch o tag nuevamente.

Los comandos que requieren dependencias se interpretan automáticamente y se pueden preparar de antemano con los siguientes comandos.

```shell
vex fetch
vex update
vex update math shared_core
```

`vex update` actualiza todos los paquetes Git, o solo los paquetes especificados y los gráficos de transición afectados. Los paquetes bloqueados no relacionados mantendrán commit seleccionado.

## Flujo de trabajo locked y offline

`--locked` prohíbe la creación y modificación de `vex.lock`. Fallará si el archivo no existe, tiene un esquema no compatible o no coincide con el gráfico manifest. Ya fijado en lockfile, commit se puede importar cuando sea necesario.

`--offline` prohíbe todas las operaciones de red Git. Los checkout y commit requeridos ya deberían existir localmente.

```shell
vex fetch --locked
vex build --locked --offline
```

Estos dos comandos son un flujo de trabajo estricto CI. Prepare un commit bloqueado exactamente cuando la red esté disponible y luego compílelo sin cambiar la red o lockfile. dry-run no importa dependencias ni reescribe lockfile.

## Información y configuración del compilador

```shell
vex info
vex setup wavec
vex setup wavec --version <version>
vex --version
```

Vex valida el esquema `wavec` dry-run JSON antes de la compilación real. Los compiladores que no implementen el esquema requerido no ejecutarán un plan desconocido y lo rechazarán con un error de compatibilidad.
