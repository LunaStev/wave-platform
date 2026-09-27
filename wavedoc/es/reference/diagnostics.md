---
translation_set_id: diagnostics
path: reference/diagnostics
locale: es
group: reference
group_order: 5
order: 2
title: Solución de problemas: desde la instalación hasta la ejecución
summary: Aísle los pasos fallidos y limite la causa con información reproducible.
---

## Primero, distinguir las etapas del fracaso.

|fenómeno observado|comprobar primero|próxima acción|
| --- | --- | --- |
|wavec Comando no encontrado|PATH y ubicación del archivo ejecutable|Ejecute con ruta absoluta y configure PATH|
|No puedo encontrar los archivos necesarios para ejecutar|¿Falta algún archivo en la carpeta de instalación?|Desempaquete e instale el paquete completo nuevamente|
|std import Error| `wavec print std-path` |Correspondencia: Instale std o especifique `--std-root`|
|Ubicación de origen y salida de error de tipo| `wavec check main.wave` |Corrija el primer error y verifique nuevamente|
|Error de compilación para otros objetivos OS·CPU|Especificado target y entorno de destino|[Configuración de compilación cruzada](/docs/es/whale/build-link-targets) Confirmar|
|La ejecución falla después de una compilación exitosa|código de salida, entrada, directorio de trabajo|Entorno de ejecución y comprobación de errores API|

## Un pequeño ejemplo de diagnóstico

Aquí está el programa completo, que es intencionalmente incorrecto:

```wave
fun main() {
    var count: i32 = 1;
    println("{}", missing);
}
```

`wavec check main.wave` debe apuntar al nombre no declarado missing. Cambie el nombre de la variable a count, luego verifique y ejecute nuevamente. Concéntrese en el archivo, la ubicación y la causa en lugar de todo el texto de diagnóstico. Los errores posteriores pueden ser el resultado del error inicial.

## Cuando falla un ejecutable

En el shell Linux/macOS, el código de salida se verifica inmediatamente después de la ejecución como `echo $?`, y en PowerShell, es `$LASTEXITCODE`. El error de entrada y `return 1` explícito no son la misma causa. Los valores de tiempo de ejecución no válidos para el recuento de turnos o la conversión real pueden causar trap. Consulte [reglas de operación](/docs/es/language/expressions-and-operators).

Las rutas de archivo relativas se ven afectadas por el directorio de trabajo ejecutable y no por la ubicación del archivo fuente. No trate un error de lectura de archivo como una longitud de cadena de 0, verifique primero el error de retorno. Las fallas de conexión de red se verifican mediante la búsqueda de direcciones, la espera del servidor, los permisos y los tiempos de espera.

## Información necesaria para informar un problema

1. `wavec --version` Salida y comando exacto ejecutado.
2. target. especificado por separado del host OS·arquitectura
3. Fuente del compilador utilizada con la ruta std seleccionada.
4. Fuente, entrada y archivos mínimos necesarios para reproducir el problema.
5. Resultados esperados, resultados reales, diagnóstico y códigos de salida.

Se eliminan las contraseñas, los tokens y el contenido de los archivos personales. Si el problema desaparece cuando reduce el ejemplo mínimo, el último elemento que eliminó es la pista. `--error-format=json` está disponible cuando la herramienta recopila diagnósticos.

[Instalación](/docs/es/getting-started/install) · [comando del compilador](/docs/es/getting-started/compiler) · [Objetivos y enlaces](/docs/es/whale/build-link-targets)
