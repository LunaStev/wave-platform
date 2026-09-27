---
translation_set_id: stdlib-path-env
path: stdlib/path-env
locale: es
group: stdlib
group_order: 1
order: 12
title: path y env: Rutas y configuración del entorno
summary: Lee la ruta y las variables de entorno en el búfer de la persona que llama e identifica errores de capacidad.
---

## combinación de caminos

```text
std::path::copy
path_join2(dst: ptr<u8>, dst_cap: i32, left: str, right: str) -> i32
path_basename_copy(dst: ptr<u8>, dst_cap: i32, path: str) -> i32
path_dirname_copy(dst: ptr<u8>, dst_cap: i32, path: str) -> i32
```

La capacidad incluye el último espacio NUL. Un resultado exitoso es cualquier longitud excluyendo NUL, el fracaso es -1. Solo en caso de éxito usaremos el destino como una cadena. Estas funciones operan en cadenas de ruta y no verifican la existencia de archivos ni los permisos de acceso. La combinación de rutas por sí sola no evita los escapes de directorios ni verifica la identidad de los archivos reales.

## variable de entorno

```text
std::env::environ
env_get(name: str, dst: ptr<u8>, dst_cap: i64) -> i64
env_exists(name: str) -> bool
env_get_i64(name: str) -> EnvResult<i64>
env_get_i32(name: str) -> EnvResult<i32>
```

`env_get` devuelve la longitud excluyendo NUL en caso de éxito. El búfer de llamadas debe poder contener hasta NUL. Un valor vacío se diferencia de un error sin clave porque puede tener éxito con una longitud de 0.

Obtenga y distinga errores de `std::env::consts` a NOT_FOUND, NO_SPACE, INVALID_KEY, READ, SOURCE_INCOMPLETE, NO_MEMORY. No trate la falta de búfer como si faltara una clave. Para buscar números, marque ok en el resultado y luego use value. No trate automáticamente el contenido de las variables de entorno como configuraciones confiables; comprobar su alcance y tipo.

El siguiente ejemplo combina los directorios data y los nombres de archivos input.txt.

<!-- wave-example: path-api -->
```wave
import("std::path::copy")::{
    path_join2
};

fun main() -> i32 {
    var output: array<u8, 64>;
    var length: i32 = path_join2(&output[0], 64, "data", "input.txt");
    if (length < 0) {
        return 1;
    }

    println("{}", &output[0] as str);
    return 0;
}
```

Resultado de la ejecución:

```text
data/input.txt
```

## Dividir directorios y nombres de archivos

El siguiente ejemplo copia una ruta dividida en dos búferes. No es necesario que el archivo original exista realmente.

<!-- wave-example: book-path-parts -->
```wave
import("std::path::copy")::{
    path_basename_copy,
    path_dirname_copy
};

fun main() -> i32 {
    var directory: array<u8, 64>;
    var filename: array<u8, 64>;
    var directory_length: i32 = path_dirname_copy(&directory[0], 64, "data/report.txt");
    var filename_length: i32 = path_basename_copy(&filename[0], 64, "data/report.txt");

    if (directory_length < 0 || filename_length < 0) {
        return 1;
    }

    println("directory={}", &directory[0] as str);
    println("filename={}", &filename[0] as str);
    return 0;
}
```

Resultado de la ejecución:

```text
directory=data
filename=report.txt
```

Ambas reservas permanecerán vigentes hasta el final de main. `as str` lee el mismo búfer que una cadena sin asignar una nueva cadena. Por lo tanto, si modifica el búfer, la cadena leída en esa dirección también cambiará.

## Establecer preferencias predeterminadas

Las variables de entorno son configuraciones pasadas fuera del programa. Al leer una configuración numérica, marque "¿Se puede leer como un número entero?" y "¿Está dentro del rango permitido por este programa?"

<!-- wave-example: book-env-setting -->
```wave
import("std::env::environ")::{
    EnvResult,
    env_get_i32
};

fun main() -> i32 {
    var setting: EnvResult<i32> = env_get_i32("WAVE_EXAMPLE_WORKERS");
    var workers: i32 = 4;

    if (setting.ok) {
        if (setting.value < 1 || setting.value > 32) {
            println("workers must be between 1 and 32");
            return 1;
        }

        workers = setting.value;
    }

    println("workers={}", workers);
    return 0;
}
```

Si `WAVE_EXAMPLE_WORKERS` no está presente o no se puede leer como un número entero, se utiliza el valor predeterminado de 4. Si se establece un número entero entre 1 y 32, se utiliza ese valor, y si es un número entero fuera del rango, termina con un error.

Linux/macOS:

```shell
WAVE_EXAMPLE_WORKERS=8 wavec run main.wave
```

PowerShell:

```powershell
$env:WAVE_EXAMPLE_WORKERS = "8"
wavec run main.wave
```

En ambos casos, se genera `workers=8`. El ejemplo anterior selecciona una política predeterminada simple. Si esta es una configuración obligatoria, trate los errores de búsqueda numérica como errores en lugar de reemplazarlos con valores predeterminados. Si necesita distinguir entre clave faltante, búfer insuficiente y error de lectura, utilice las constantes env_get y ENV_ERR_*.
