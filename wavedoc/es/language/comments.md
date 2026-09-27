---
translation_set_id: comments
path: language/comments
locale: es
group: language
group_order: 2
order: 15
title: Comentarios
summary: Describe comentarios de una sola línea, comentarios de bloques anidados y diagnósticos de comentarios no cerrados.
---

## comentario de una linea

El contenido después de `//` es un comentario hasta el final de la línea.

```wave
var count: i32 = 10;
// 현재 요청 수
```

## anotación de bloque

Procese el espacio entre `/*` y `*/` como comentario de bloque.

```wave
/* 여러 줄에 걸친
   설명을 작성할 수 있습니다. */
```

Puede anidar otros comentarios de bloque dentro de un comentario de bloque.

```wave
/* 바깥 주석
   /* 안쪽 주석 */
다시 바깥 주석
*/
```

## Cadenas y marcas de comentarios

`//`, `/*` y `*/` dentro de cadenas y caracteres literales son contenido de cadena y no se tratan como el principio ni el final de un comentario.

```wave
var text: str = "https://wave-lang.dev";
```

## Comentarios de bloque sin cerrar

Si no se cierra el comentario del bloque con `*/`, se producirá el diagnóstico `E1002 UnterminatedComment`.

Incluso cuando desactive temporalmente los bloques largos, asegúrese de que la profundidad de anidamiento sea correcta.
