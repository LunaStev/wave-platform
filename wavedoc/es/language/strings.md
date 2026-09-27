---
translation_set_id: learn-strings
path: language/strings
locale: es
group: language
group_order: 2
order: 7
title: 7. Cadenas, caracteres y bytes
summary: Distinga entre cadena y char, UTF-8 longitud de bytes, NUL, búsqueda y datos binarios.
---

## letras en pantalla y bytes en memoria

Ves letras en la pantalla, pero los bytes se almacenan en la memoria. Especialmente en los casos en los que un carácter tiene varios bytes UTF-8, como en coreano, es fácil cometer un error si "longitud" y "número de caracteres" se usan indistintamente.

Este capítulo distingue entre str y char, terminando en NUL, escape, ubicación de búsqueda y datos binarios. Cada uno de los ejemplos es un programa completo.

## Literales de cadena y salida

<!-- wave-example: book-string-literals -->
```wave
fun main() {
    var greeting: str = "안녕하세요";

    println("{}", greeting);
    println("line one\nline two");
    println("quote: \"Wave\"");
}
```

Resultado de la ejecución:

```text
안녕하세요
line one
line two
quote: "Wave"
```

Los caracteres regulares entre comillas dobles se expresan como UTF-8. escape indica bytes que son difíciles de escribir directamente desde la fuente. `\n` es un byte de salto de línea y no imprime dos caracteres, una barra invertida y n. Para imprimir la barra invertida, utilice `\\`.

## La longitud es el número de bytes.

<!-- wave-example: book-utf8-length -->
```wave
import("std::string::len")::{
    len
};

fun main() {
    println("ASCII={}", len("Wave"));
    println("Korean={}", len("한"));
    println("mixed={}", len("Wave한"));
}
```

Resultado de la ejecución:

```text
ASCII=4
Korean=3
mixed=7
```

`len` cuenta bytes antes del NUL final, caracteres no visibles. El carácter `한` ocupa tres bytes en UTF-8. Los recuentos de caracteres, los recuentos de puntos de código Unicode y los recuentos de bytes generalmente no son intercambiables. El ancho de visualización también depende de factores como las fuentes y la combinación de caracteres.

Por lo tanto, las funciones que cortan una cadena en una posición de byte arbitraria y la muestran en la pantalla deben considerar por separado el límite Unicode. Defina claramente las condiciones de entrada, ya sea un programa que solo procesa texto ASCII o general Unicode.

## NUL Extremo y Longitud

str utiliza un byte 0 para indicar el final. len no incluye ese último byte en la longitud. Es un error poner NUL dentro de una cadena literal. Aquí hay un ejemplo de un error intencional:

```wave
fun main() {
    var text: str = "left\x00right";
}
```

`\xNN` en la fuente especifica un byte con exactamente dos dígitos hexadecimales. `\x41` representa el byte A, que es 65. Debido a que escribir un carácter normal como UTF-8 e insertar un byte arbitrario son diferentes, no todos los str que se pueden escribir como `\xNN` son válidos UTF-8.

## char no contiene el carácter completo Unicode

char es un valor de carácter de 8 bits sin signo. Puede utilizar literales que representen valores en un rango de un solo byte, como `'A'`. `'한'` es un error porque no se encuentra dentro de este rango. `"한"` es un str separado con múltiples UTF-8 bytes.

No intentes poner siempre una letra en una char. Primero debe decidir si las unidades necesarias para procesar el texto son bytes o puntos de código Unicode.

## comparación de cadenas

Para comparar el contenido de la cadena, use la función std. A continuación se muestra un programa que busca contenido idéntico y diferencias entre casos.

<!-- wave-example: book-string-compare -->
```wave
import("std::string::cmp")::{
    eq, starts_with, ends_with
};

fun main() {
    if (eq("Wave", "Wave")) {
        println("same bytes");
    }

    if (!eq("Wave", "wave")) {
        println("case differs");
    }

    if (starts_with("report.txt", "report") && ends_with("report.txt", ".txt")) {
        println("name matches");
    }
}
```

Resultado de la ejecución:

```text
same bytes
case differs
name matches
```

Esta comparación compara cadenas de bytes. No realiza automáticamente la conversión de casos específicos del idioma ni la normalización Unicode. Incluso al comparar nombres de archivos, las reglas de igualdad de nombres de archivos en OS no son las mismas que las comparaciones simples de cadenas.

## Unidades y fallos de los resultados de búsqueda.

<!-- wave-example: book-string-search -->
```wave
import("std::string::find")::{
    find, count
};

fun main() {
    var first: i32 = find("banana", "na");
    var missing: i32 = find("banana", "xy");
    var matches: i32 = count("aaaa", "aa");

    println("first={} missing={}", first, missing);
    println("matches={}", matches);
}
```

Resultado de la ejecución:

```text
first=2 missing=-1
matches=2
```

find devuelve la primera posición o -1. El índice 0 también es un éxito, por lo que se comprueba con `result >= 0`. count no es una posición, sino el número de coincidencias que no se superponen. Arriba, aa es el número 2 ya que coincide con 0~1 y 2~3.

Bin needle también forma parte del contrato. find devuelve 0, contains devuelve true y count devuelve 0. No asuma que solo porque el nombre de la función está en el mismo módulo, incluso el método de retorno es el mismo.

## Eliminar espacios es diferente a crear una nueva cadena

trim_range devuelve el rango excluyendo espacios sin modificar ni copiar el texto original. Como recibimos un puntero de salida, primero preparamos un número entero para almacenar el resultado.

<!-- wave-example: book-string-trim-range -->
```wave
import("std::string::trim")::{
    trim_range
};

fun main() {
    var start: i32 = 0;
    var end: i32 = 0;

    trim_range("  Wave  ", &start, &end);

    println("start={} end={} bytes={}", start, end, end - start);
}
```

Resultado de la ejecución:

```text
start=2 end=6 bytes=4
```

El rango es `[start, end)`. Incluye el inicio pero no el final, por lo que su longitud es end-start. Agregar start a la dirección inicial del original no crea automáticamente NUL en la ubicación end. Deberá transportar la estufa por separado o preparar un nuevo espacio para cuerdas.

## Los datos binarios tienen una longitud separada.

Los datos que contienen ceros no están cubiertos por las reglas de fin de cadena. Utiliza matrices de bytes y longitudes.

<!-- wave-example: book-binary-not-string -->
```wave
fun main() {
    var bytes: array<u8, 3> = [65, 0, 66];

    for (var index: i32 = 0; index < 3; index += 1) {
        println("{}", bytes[index]);
    }
}
```

Resultado de la ejecución:

```text
65
0
66
```

El segundo 0 son los datos reales. Si interpreta esto como str, se considera que termina en el primer 0 y no puede ver el siguiente 66. Por el contrario, si cambia una matriz sin NUL a str con cast, existe el riesgo de leer más allá de la matriz. cast no es una operación para agregar un byte de terminación.

## Ejercicio: examinar nombres de archivos

Compruebe si el nombre del archivo termina con `.wave` y, si la cadena contiene `test`, envíelo a un archivo de prueba. Este ejercicio sólo comprueba patrones de bytes en el nombre y no aborda la existencia real del archivo.

### Solución completa

<!-- wave-example: book-string-solution -->
```wave
import("std::string::cmp")::{
    ends_with
};
import("std::string::find")::{
    contains
};

fun classify(name: str) {
    if (!ends_with(name, ".wave")) {
        println("other file");
        return;
    }

    if (contains(name, "test")) {
        println("Wave test file");
    } else {
        println("Wave source file");
    }
}

fun main() {
    classify("main.wave");
    classify("parser_test.wave");
    classify("notes.txt");
}
```

Resultado de la ejecución:

```text
Wave source file
Wave test file
other file
```

Existen políticas separadas sobre cómo manejar la letra mayúscula `.WAVE` y si considerarla como una prueba incluso si test está incluida en toda la ruta. Incluso si se trata de una función pequeña, su funcionamiento sólo se puede describir con precisión si se determina a qué entrada se dirige.
