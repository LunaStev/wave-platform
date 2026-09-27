---
translation_set_id: system-io
path: reference/system-io-network-process
locale: es
group: stdlib
group_order: 1
order: 15
title: Funciones y procesos del sistema.
summary: Describe el límite y la vida útil del proceso de las interfaces principales API y OS.
---

## Documentación por función

Lea [fs y io](/docs/es/stdlib/files-io) para manejar el archivo, [TCP](/docs/es/stdlib/tcp) para vincular y [resolver](/docs/es/stdlib/resolution) para buscar la dirección. A continuación se muestran el proceso y las reglas de acceso de nivel inferior OS.

## Conceptos básicos del proceso API

```text
std::process::core
proc_exit(code: i32) -> !
proc_getpid() -> i64
proc_getppid() -> i64
proc_execve(path: str, argv: ptr<ptr<i8>>, envp: ptr<ptr<i8>>) -> i64
proc_waitpid_raw(pid: i64, status: ptr<i32>, options: i32) -> i64
```

`proc_exit` no regresa al pulsador. Realice cualquier limpieza necesaria de archivos/memoria antes del apagado. `proc_execve` se diferencia de una función de creación de hijos típica porque, si tiene éxito, reemplaza la imagen de proceso existente. raw argv/envp debe estar preparado para la terminación NUL de cada cadena, con el puntero null indicando el final.

La función spawn en `std::process::spawn` maneja el resultado de la creación y la función de espera maneja el estado de salida del niño. La creación exitosa y la finalización exitosa del programa son dos cosas diferentes. Cuando crea una tubería, el padre y el hijo deben cerrar el extremo no utilizado para que se pase EOF. Si espera a que el niño salga sin leer el canal de captura, el búfer puede llenarse y esperarse el uno al otro.

## Portabilidad y enfoque de bajo nivel.

fork/exec, el descriptor de archivo y el identificador Windows no son la misma función OS. Verifique el soporte para el objetivo seleccionado y trate unsupported como una ruta de falla normal. `std::sys` es una interfaz específica de OS y no reutiliza sus indicadores numéricos ni su diseño de otros OS.

Al vincular directamente con la biblioteca externa C, lea [FFI](/docs/es/language/modules-imports-and-ffi). No es necesario declarar arbitrariamente la función libc para usar el padre std API. Verifique [Entorno de destino y enlace](/docs/es/whale/build-link-targets) primero.
