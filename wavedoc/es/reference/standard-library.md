---
translation_set_id: standard-library
path: reference/standard-library
locale: es
group: stdlib
group_order: 1
order: 1
title: guía de biblioteca estándar
summary: Cómo encontrar un módulo que se adapte a su propósito y leer los errores y las reglas de propiedad de la función.
---

## Encuentre las funciones que necesita

La biblioteca estándar es import con la ruta `std::module::file`. Incluso si los nombres son similares, las funciones pueden devolver errores de diferentes maneras. Primero lea [API Cómo leer](/docs/es/stdlib/contracts), luego vaya al módulo que necesita en la siguiente tabla.

|lo que quiero hacer|documento|Principal import|
| --- | --- | --- |
|Longitud de cadena/comparación/búsqueda| [string](/docs/es/reference/string-and-bytes) | `std::string::len`, `cmp`, `find`, `trim` |
|Asignación de memoria/copia/tamaño| [mem](/docs/es/reference/memory-and-buffer) | `std::mem::alloc`, `ops`, `layout` |
|Lista de bytes de diferente tamaño| [buffer](/docs/es/stdlib/buffer) | `std::buffer::alloc`, `read`, `write` |
|Lectura/escritura binaria| [bytes](/docs/es/stdlib/bytes) | `std::bytes::types`, `cursor`, `leb128` |
|Archivo·Descriptor I/O| [fs y io](/docs/es/stdlib/files-io) | `std::fs::file`, `std::io::fd` |
|Combinación de ruta/configuración del entorno| [path y env](/docs/es/stdlib/path-env) | `std::path::copy`, `std::env::environ` |
|Medición del tiempo/espera| [time](/docs/es/stdlib/time) | `std::time::duration`, `clock`, `sleep` |
|Búsqueda numérica de lista de direcciones/nombres| [net.resolve](/docs/es/stdlib/resolution) | `std::net::resolve`, `resolve_table` |
|TCP Conexión/Transmisión| [net.tcp](/docs/es/stdlib/tcp) | `std::net::tcp`, `address`, `error` |
|OS Número aleatorio| [random](/docs/es/stdlib/random) | `std::random::fill` |
|Proceso·OS Límite| [función del sistema](/docs/es/reference/system-io-network-process) | `std::process::core`, `spawn`, `std::sys` |
|Ejecución de tareas asincrónicas| [task](/docs/es/stdlib/task) | `std::task` |
|Asistente de Matemáticas/Diagnóstico| [math y debug](/docs/es/stdlib/math-debug) | `std::math::int`, `float`, `std::debug::core` |

## Ejemplo de primer uso

El siguiente programa utiliza una función de std sin descargar un paquete por separado. Guárdelo como `main.wave` y ejecútelo como `wavec run main.wave`.

<!-- wave-example: library-import -->
```wave
import("std::string::len")::{
    len
};

fun main() {
    println("{}", len("Wave"));
}
```

El resultado es `4`. Para entender el mismo ejemplo paso a paso, lea [cuerdas](/docs/es/language/strings).

## Compatible con el compilador std

Confirme la ruta seleccionada con `wavec print std-path`. Cuando utilice std desde otro pago, especifique la ruta como `wavec --std-root /absolute/path/to/std check main.wave`. Si la ruta especificada no es válida o es incompatible, se mostrará un error.

## borde de la plataforma

Distinga entre funciones computacionales como cadena/byte y funciones OS como archivo/socket. Reconocer un objetivo no garantiza que se proporcionen todos los hosts API. Lea las entradas de la plataforma para [Objetivo de soporte](/docs/es/whale/build-link-targets) y cada API juntas. `std::sys` es una interfaz de nivel inferior y los programas portátiles utilizarán primero el módulo de nivel superior.
