---
translation_set_id: lexical
path: language/lexical-structure
locale: es
group: language
group_order: 2
order: 14
title: Estructura léxica
summary: Describe identificadores, literales, delimitadores, palabras clave y nombres de tipos.
---

## identificador

Los identificadores nombran variables, funciones, tipos y campos. El nombre distingue entre mayúsculas y minúsculas y puede contener cualquier combinación de letras, números y `_`. No se pueden utilizar números en la primera letra. Los caracteres Unicode también se pueden utilizar en identificadores.

```wave
var request_count: i64 = 0;
var 이름: str = "Wave";
```

En proyectos reales, se recomienda utilizar una convención de nomenclatura coherente para la compatibilidad de herramientas y la capacidad de búsqueda.

## Oraciones y Separadores

La mayoría de las declaraciones y expresiones terminan con `;`. Las declaraciones con cuerpo, como funciones, declaraciones condicionales, declaraciones de bucle y estructuras, utilizan el bloque `{ ... }`.

```wave
var answer: i32 = 42;

fun double(value: i32) -> i32 {
    return value * 2;
}
```

## literal

```wave
var integer: i32 = 42;
var decimal: f64 = 3.14;
var text: str = "Wave";
var letter: char = 'W';
var enabled: bool = true;
var address: ptr<u8> = null;
```

Puede utilizar números enteros, números de punto flotante, cadenas, caracteres, valores booleanos y literales `null`. Utilice `null` para valores de puntero.

## Palabras clave y nombres de tipos

Las principales palabras clave utilizadas en la gramática Wave son las siguientes.

`pub`, `fun`, `extern`, `export`, `type`, `enum`, `static`, `var`, `deref`, `const`, `if`, `else`, `proto`, `struct`, `while`, `for`, `in`, `out`, `clobber`, `as`, `asm`, `import`, `return`, `continue`, `print`, `input`, `println`, `match`, `break`, `true`, `false`, `null`.

Los nombres de tipos integrados incluyen `bool`, `char`, `byte`, `str`, tipos de punto flotante y entero, `ptr` y `array`. Los punteros se escriben en la forma `ptr<T>` y las matrices de longitud fija se escriben en la forma `array<T, N>`.

## Cadenas y caracteres escape

|notación|significado|
| --- | --- |
| `\n` |LF Salto de línea|
| `\r` | CR |
| `\t` |pestaña|
| `\\` |Barra invertida|
| `\"` |comillas dobles|
| `\xNN` |Un byte especificado como exactamente dos dígitos hexadecimales|

Los caracteres de cadena generales se guardan como UTF-8. Dado que `\xNN` conserva un byte, no hay garantía de que toda la cadena sea un UTF-8 válido. NUL (incluido `\x00`) dentro de una cadena literal es un error de compilación. Para datos que contienen ceros, utilice una matriz de bytes y una longitud.

`char` El literal debe caber en un valor de 8 bits. Los caracteres que exceden ese rango, como `'한'`, son errores. Es diferente de la cadena `"한"`.

LF, CRLF y CR solos en el código fuente se tratan como un salto de línea lógico. Estas son reglas sobre la ubicación del origen y la terminación de comentarios y no significan que cambien los bytes reales de los datos del archivo.

Los nombres gramaticales adicionales incluyen `variant`, `async` y `await`, y el valor asincrónico se expresa como `Future<T>`. El bloque de declaración independiente `var` anterior es un fragmento de código dentro de una función.

[clase de cadena](/docs/es/language/arrays) · [Comentarios](/docs/es/language/comments)

## Ejemplo de fracaso intencional

Si ejecuta el programa siguiente check, debería producirse un error interno NUL. Si necesita 0 bytes, utilice la matriz de bytes `[97, 0, 98]`.

<!-- wave-example: reject-string-nul -->
```wave
fun main() {
    var text: str = "a\x00b";
}
```

El carácter literal a continuación también excede el rango de 8 bits, por lo que se trata de un error de compilación. Para representar la cadena UTF-8, utilice `str` y comillas dobles.

<!-- wave-example: reject-wide-char -->
```wave
fun main() {
    var letter: char = '한';
}
```
