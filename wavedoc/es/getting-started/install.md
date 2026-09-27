---
translation_set_id: install
path: getting-started/install
locale: es
group: getting-started
group_order: 1
order: 2
title: Instalar Wave
summary: Instala Wave en Linux, macOS o Windows y ejecuta tu primer programa.
---

## Linux y macOS

Ejecuta el siguiente comando en una terminal. Instala Wave junto con el gestor de paquetes Vex.

```shell
curl -fsSL https://wave-lang.dev/install.sh | bash -s -- latest
```

Cuando termine la instalación, abre una terminal nueva y comprueba la versión.

```shell
wavec --version
```

## Windows

Ejecuta los siguientes comandos en PowerShell. Instalan Wave junto con el gestor de paquetes Vex.

```powershell
irm https://wave-lang.dev/install.ps1 -OutFile install.ps1
powershell -ExecutionPolicy Bypass -File .\install.ps1 -Latest
```

Cuando termine la instalación, abre una ventana nueva de PowerShell y comprueba las versiones.

```powershell
wavec --version
vex --version
```

## Primera ejecución

Guarda el siguiente código como `main.wave`.

<!-- wave-example: install-stdlib -->
```wave
import("std::string::len")::{
    len
};

fun main() {
    println("Wave: {} bytes", len("Wave"));
}
```

Ejecuta el programa desde el directorio donde guardaste el archivo.

```shell
wavec run main.wave
```

Resultado:

```text
Wave: 4 bytes
```

Si aparece un mensaje que indica que no se encuentra la biblioteca estándar, instálala y vuelve a ejecutar el programa.

```shell
wavec install std
wavec run main.wave
```

[Siguiente: Tu primer programa](/docs/es/language/program-structure) · [Solución de problemas](/docs/es/reference/diagnostics)
