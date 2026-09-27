---
translation_set_id: quick-reference
path: reference/syntax-quick-reference
locale: es
group: reference
group_order: 5
order: 3
title: Referencia rápida de sintaxis
summary: Las declaraciones de uso frecuente, el flujo de control, los tipos, los punteros y la gramática FFI están organizados en una sola página.
---

## declaración

```wave
var value: i32 = 1;
var next: i32 = value + 1;
const LIMIT: i32 = 64;
static total: i64 = 0;
type Identifier = u64;
```

`var` es la región, `const`/`static` son declaraciones de nivel superior. Las variables locales declaran explícitamente su tipo.

## Funciones

```wave
fun max(left: i32, right: i32) -> i32 {
    if (left > right) {
        return left;
    }

    return right;
}
```

## genérico

```wave
fun identity<T>(value: T) -> T {
    return value;
}

var value: i32 = identity<i32>(10);
```

Al llamar a un genérico, especifique un argumento de tipo.

## Estructura y enum

```wave
struct Pair {
    left: i32;
    right: i32;
}

enum Result -> i32 {
    Ok = 0,
    Error,
}
```

## Condiciones y bucles

```wave
if (ready) {
    println("ready");
}

while (count < 10) {
    count += 1;
}

for (var i: i32 = 0; i < 10; i += 1) {
    println("{}", i);
}

match (status) {
    Ready => {
        println("ready");
    }
    0 => {
        println("zero");
    }
    _ => {
        println("other");
    }
}
```

Los encabezados de `if`, `while`, `for` y `match` usan paréntesis.

## Matrices y punteros

```wave
var values: array<i32, 4> = [1, 2, 3, 4];
var p: ptr<i32> = &values[0];
var first: i32 = deref p;
```

## entrada/salida de consola

```wave
print("value = ");
println("{}", value);
input("{}", value);
```

El primer argumento es una cadena literal. Cada marcador de posición `{}` exacto requiere una expresión a continuación y el destino `input` debe ser asignable.

## import y elementos públicos

```wave
import("std::string::len");
import("./helpers" as helpers);
import("math")::{
    Vector
};

pub fun add(left: i32, right: i32) -> i32 {
    return left + right;
}

pub import("./extra"):: {
    increment
};
```

La ruta local comienza con `./`. El alias import especifica el nombre del módulo y la selección import importa las entradas públicas requeridas al espacio de nombres de este archivo. `pub import` reexporta los elementos seleccionados.

## FFI

```wave
extern(c) fun native_call(value: i32) -> i32;

export(c) fun wave_call(value: i32) -> i32 {
    return value + 1;
}
```

## Elemento condicional de destino

```wave
#[target(os="linux", arch="riscv64")]
extern(c) fun platform_call(value: i32) -> i32;
```

Las claves de condición de soporte son `arch`, `os`, `env`, `abi`. Las propiedades controlan el siguiente elemento de nivel superior.

## Montaje en línea

```wave
var result: i64 = 0;
asm {
    "mv a0, a1"
    in("a1") 7
    out("a0") result
    clobber("memory")
}
```

El texto de las instrucciones y los nombres de los registros dependen del objetivo. Declare todas las entradas, salidas y clobber ocultas requeridas para el bloque.

## inspección de fuente

```shell
wavec build main.wave --emit=check
wavec print supported-targets
wavec print supported-emit-kinds
```

## Rango de aprendizaje y ejemplo

Ejemplos de variables locales y declaraciones que se muestran por separado fuera de la función son fragmentos de código insertados en el cuerpo de la función. A continuación encontrará ejemplos y ejercicios completos de carrera en [Wave Proceso de aprendizaje](/docs/es/getting-started/overview). Consulte [Biblioteca estándar](/docs/es/stdlib) para conocer las reglas detalladas de memoria y funciones externas.
