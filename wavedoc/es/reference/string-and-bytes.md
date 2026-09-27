---
translation_set_id: string-bytes
path: reference/string-and-bytes
locale: es
group: stdlib
group_order: 1
order: 3
title: string: longitud, búsqueda y rangos
summary: NUL Describe la unidad de bytes de la cadena de terminación API y el valor de retorno.
---

## Condiciones de almacenamiento y argumentos de cadenas.

El argumento `str` de este módulo debe ser un byte de terminación NUL accesible. La longitud y el índice de búsqueda están en bytes. Los caracteres normales se almacenan como UTF-8, pero la búsqueda de bytes es Unicode sin normalización ni división carácter por carácter. No asuma que el índice devuelto es un límite de caracteres.

## Comparar con longitud

```text
std::string::len
len(s: str) -> i32
is_empty(s: str) -> bool

std::string::cmp
eq(a: str, b: str) -> bool
cmp(a: str, b: str) -> i32
starts_with(s: str, prefix: str) -> bool
ends_with(s: str, suffix: str) -> bool
```

`len` excluye el último NUL. `cmp` El orden se juzga por el signo del resultado. El valor devuelto no se interpreta como el orden de los caracteres Unicode ni como una clasificación de diccionario específica del idioma. Estas funciones no asignan memoria y no cambian su entrada.

## buscar

Obtenga el nombre que necesita, como `import("std::string::find")::{find, contains, count};`.

|declaración de función|resultado|
| --- | --- |
| `find(s: str, needle: str) -> i32` |Ubicación del primer partido. -1 si no está presente, 0 para vacío needle|
| `contains(s: str, needle: str) -> bool` |Incluido o no. needle vacío es true|
| `count(s: str, needle: str) -> i32` |Número de coincidencias que no se superponen. Bin needle es 0|
| `find_char(s: str, c: u8) -> i32` |primera posición del byte o -1|
| `rfind_char(s: str, c: u8) -> i32` |Última posición del byte o -1|
| `contains_char(s: str, c: u8) -> bool` |Existencia de ese byte|
| `count_char(s: str, c: u8) -> i32` |el número de bytes en cuestión|

`c` en el nombre `*_char` es un byte, no un punto de código Unicode. El NUL al final de la cadena no está incluido en el objetivo de búsqueda.

## Rango excluyendo espacios

```text
std::string::trim
trim_left_index(s: str) -> i32
trim_right_index(s: str) -> i32
trim_range(s: str, out_start: ptr<i32>, out_end: ptr<i32>)
```

`trim_range` escribe el rango semiabierto `[start, end)` excluyendo espacios ASCII en el argumento de salida. Ambos punteros de salida deben apuntar a números enteros grabables. No modifica el texto original ni crea nuevas cadenas. Si todo está en blanco, se convierte en un rango vacío.

## Ejemplo de ejecución

Guárdelo en `main.wave` y ejecute `wavec run main.wave`.

<!-- wave-example: string-api -->
```wave
import("std::string::find")::{
    find, count
};
import("std::string::trim")::{
    trim_range
};

fun main() {
    var start: i32 = 0;
    var end: i32 = 0;
    trim_range("  Wave  ", &start, &end);
    println("{} {}", start, end);
    println("{} {}", find("banana", "na"), count("aaaa", "aa"));
}
```

Resultado de la ejecución:

```text
2 6
2 2
```

## Funciones relacionadas

La conversión de clasificación/caso de `std::string::ascii` es para el rango ASCII. `djb2_32` y `fnv1a_64` de `std::string::hash` no se utilizan para hashes criptográficos ni para almacenamiento de contraseñas. Para datos que contengan NUL, utilice [bytes](/docs/es/stdlib/bytes).

## Patrones que se superponen con términos de búsqueda vacíos

Ver el comportamiento de los bordes de la función de búsqueda con valores reales facilita la determinación de las condiciones de llamada.

<!-- wave-example: book-string-empty-search -->
```wave
import("std::string::find")::{find, contains, count};

fun main() {
    println("empty find={}", find("Wave", ""));
    println("empty count={}", count("Wave", ""));
    println("nonoverlapping={}", count("aaaa", "aa"));

    if (contains("Wave", "")) {
        println("empty needle is contained");
    }
}
```

Resultado de la ejecución:

```text
empty find=0
empty count=0
nonoverlapping=2
empty needle is contained
```

El valor de éxito de find, 0, es la primera posición. Un valor de éxito de 0 para count es el resultado de una regla de término de búsqueda vacía o sin coincidencias. No hay dos valores que reciban el mismo trato. Si se requiere ignorar casos o normalizar Unicode, se deben implementar políticas separadas antes y después de esta búsqueda de bytes.

## trim Copiando rango a una nueva cadena

No hay ningún final nuevo NUL en el rango devuelto por trim_range. Al copiar a un destino separado, reserve longitud + 1 espacio y escriba el último byte directamente como 0.

<!-- wave-example: book-trim-copy -->
```wave
import("std::string::trim")::{trim_range};

fun main() -> i32 {
    var original: str = "  Wave  ";
    var start: i32 = 0;
    var end: i32 = 0;
    var destination: array<u8, 16>;

    trim_range(original, &start, &end);

    var length: i32 = end - start;

    if (length >= 16) {
        return 1;
    }

    for (var index: i32 = 0; index < length; index += 1) {
        destination[index] = original[start + index];
    }

    destination[length] = 0;
    println("{}", &destination[0] as str);
    return 0;
}
```

Resultado de la ejecución:

```text
Wave
```

Una cadena de longitud 16 no encajará en este destino. Esto se debe a que necesitas el último NUL. Incluso si la longitud es 0, escribir destination[0]=0 da como resultado una cadena vacía válida. La matriz local de destino vive hasta el final de main, por lo que imprimimos dentro de ella.

## Cadena API Orden de uso

Al diseñar una API de cadena, especifique si su entrada termina en NUL, si los índices cuentan bytes y si el resultado toma prestado un rango de origen o posee una nueva asignación. Un rango prestado depende de la vida útil de la fuente. Un resultado asignado debe especificar quién lo libera.

Lea [Capítulo de aprendizaje de cuerdas](/docs/es/language/strings) para conocer conceptos básicos y [bytes](/docs/es/stdlib/bytes) para datos que incluyen NUL.
