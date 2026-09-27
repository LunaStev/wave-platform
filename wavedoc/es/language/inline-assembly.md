---
translation_set_id: assembly
path: language/inline-assembly
locale: es
group: language
group_order: 2
order: 19
title: Montaje en línea
summary: Describe el contrato de la cadena de comando del bloque asm, los operandos in/out y clobber.
---

## asm bloque

`asm` es una sintaxis de bajo nivel para insertar directamente instrucciones de la arquitectura de destino.

```wave
fun read_value() -> i64 {
    var result: i64 = 0;
    asm {
        "mov rax, 123"
        out("rax") result
    }

    return result;
}
```

Los literales de cadena dentro de un bloque se pasan como una lista de instrucciones de ensamblaje.

## entrada y salida

```wave
var result: i64 = 0;
asm {
    "mov rax, rdi"
    in("rdi") 123
    out("rax") result
}
```

- `in("reg") expression` conecta el valor Wave al operando de entrada.
- `out("reg") target` escribe el valor de salida en el destino Wave asignable.
- Los nombres de los registros se pueden escribir como cadenas o identificadores.

Los operandos de entrada pueden incluir variables, literales enteros/de cadena, `&identifier`, `deref identifier` y números negativos.

## clobber

Si un bloque cambia un estado de registro o memoria que no sea una salida explícita, se registra en `clobber(...)`.

```wave
asm {
    "nop"
    clobber("rax", "rcx", "memory")
}
```

## Comprobar cuando se utiliza

- La sintaxis de la instrucción debe coincidir con la arquitectura de destino y el contrato de ensamblaje en línea LLVM.
- No destruya arbitrariamente los registros que deben conservarse según la convención de llamada.
- Para bloques que leen o escriben memoria, declare clobber, incluido `memory`.
- Si es posible, aísle el asm específico de la arquitectura detrás de una función pequeña.

El comportamiento y la portabilidad del ensamblador en línea no están garantizados únicamente por el tipo de idioma.

## Rango de aprendizaje y ejemplo

[Practica con el programa completo](/docs/es/getting-started/overview) · [Biblioteca estándar](/docs/es/stdlib)
