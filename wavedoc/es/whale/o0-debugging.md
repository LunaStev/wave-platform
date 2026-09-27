---
translation_set_id: whale-o0-debugging
path: whale/o0-debugging
locale: es
group: whale
group_order: 1
order: 9
title: O0 y depuración
summary: Reglas para preservar cálculos, expresiones constantes, espacio de almacenamiento e información de respuesta de depuración.
---

## modelo de conservación

O0 conserva los cálculos, las variables y el flujo de control del typed IR original para fines de depuración. También mantiene los resultados no utilizados y los bloques estructuralmente inalcanzables. La verificación diagnostica el IR defectuoso y no elimina bloques ni simplifica las operaciones.

Las conversiones necesarias para generar lenguaje de máquina se realizan en un subsubtítulo separado IR. Incluso después de la conversión, se debe mantener la correspondencia con el original ID. La conservación de typed IR no implica una correspondencia uno a uno de todas las IR operaciones e instrucciones de la máquina.

## Conversión no realizada por O0

|conversión|O0 Operación|
| --- | --- |
|en línea|Mantener los límites de llamadas y funciones|
|Convertir Tail-call|Mantener la estructura general de llamada/devolución.|
|eliminar código muerto|Mantener cálculos no utilizados y bloques inalcanzables|
|Contraer constantes de tiempo de ejecución|mantener la operación original|
|Combinar cálculos comunes|Mantenga cálculos separados|
|Reutilización del espacio de almacenamiento de variables locales.|Mantener cada espacio de almacenamiento|
|Omitir puntero de marco|mantener el puntero del marco|
|Cadenas de fusión automática|Mantener objetos de cadena separados|

Por ejemplo, una operación en tiempo de ejecución que suma dos constantes sigue siendo una suma incluso si no se utiliza el resultado. Las sucursales con condiciones constantes también mantienen la estructura de flujo de control original.

## constantes de tiempo de compilación

Las declaraciones de constantes en tiempo de compilación conservan tanto la expresión de inicialización typed como el resultado de la evaluación. Esto es distinto del plegado constante en instrucciones de tiempo de ejecución regulares.

Una declaración con una expresión de inicialización de `1 + 2` queda con una expresión de suma y un resultado de `3`. Este ejemplo ilustra el significado de la expresión y no la gramática declarativa del idioma fuente específico. El resultado se puede utilizar al construir datos estáticos y no reemplaza la expresión de inicialización con una adición en tiempo de ejecución.

También deben identificarse las declaraciones no utilizadas e inalcanzables. Incluso si las declaraciones con el mismo nombre se ocultan entre sí, deben separarse por nombre·declaración ID·referencia. La validación rechaza los resultados almacenados que tengan referencias no válidas, ciclos de dependencia, tipos no válidos o que no coincidan con la expresión original.

## declaración de fuente inalcanzable

Las oraciones posteriores a return·break·continue también permanecen en un bloque desconectado. Debido a esta declaración, no debe cambiar el terminator anterior ni crear una nueva ruta de ejecución. También diagnostica expresiones inválidas en oraciones inalcanzables.

Esta regla de preservación permite examinar la estructura original del programa. Esto no significa que la declaración después de terminator realmente se ejecutará.

## Ejemplo de IR conservado

A continuación se muestra la salida de impresora actual para el módulo que pasó el verificador. Muestra juntos cálculos no utilizados, declaraciones en tiempo de compilación y bloques no vinculados.

```text
module {
  format_version 2
  semantics_version 1
  target "x86_64-whale-linux"
  datalayout { ptr=64, endian=little }

  declare @f0 "preserved": whale () -> void, linkage internal

  fn @preserved() -> void, id @f0 {
  entry:
    %v0: i32 = const i32 1
    %v1: i32 = const i32 2
    %v2: i32 = add i32 %v0, %v1
    %v3: i32 = const_decl "count" add(i32 1, i32 2) => const i32 3
    ret void
  unreachable.cont:
    %v4: i32 = const i32 4
    %v5: i32 = const i32 5
    %v6: i32 = add i32 %v4, %v5
    ret void
  }

}
```

`%v2` permanece como `add` incluso si no tiene uso. `%v3` es un `const_decl` separado, que conserva la expresión `add(i32 1, i32 2)` y el resultado de evaluación 3 juntos. Dado que `unreachable.cont` no tiene borde entrante, podemos verificar la adición de `%v6` sin agregar una ruta de ejecución después de `ret void`.

La verificación conserva estas instrucciones y bloques. Este ejemplo muestra la configuración/verificación/salida de IR y no implica que se proporcione la ejecución de native o la salida de DWARF.

## Información de origen y pila

La interfaz de depuración utiliza información de función, líneas de origen, variables locales predeterminadas e información de marco de llamada de DWARF 5. Incluso si se convierte en una subexpresión, se debe mantener la correspondencia entre la información de función/variable local y el identificador IR original.

El perfil AMD64 conserva el puntero del marco y no utiliza red zone. La información del marco de llamada se utiliza para la inspección de la pila y no implica compatibilidad con la excepción unwinding. trap finaliza la ejecución sin garantizar un destructor o unwinding.

Consulte [Descripción general de la cadena de herramientas](overview) para conocer la disponibilidad de la salida DWARF y la ejecución native. El comportamiento de O1 y optimizaciones superiores está fuera del alcance de esta referencia O0.
