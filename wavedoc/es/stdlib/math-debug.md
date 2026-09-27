---
translation_set_id: stdlib-math-debug
path: stdlib/math-debug
locale: es
group: stdlib
group_order: 1
order: 14
title: math y debug: Utilidades matemáticas y diagnóstico
summary: Utilice funciones matemáticas con verificación de rango y salida de diagnóstico.
---

## Función entera para verificar el rango

```text
std::math::int
min_i32(a: i32, b: i32) -> i32
max_i32(a: i32, b: i32) -> i32
abs_i32_checked(x: i32) -> MathResult<i32>
clamp_i32_checked(x: i32, lo: i32, hi: i32) -> MathResult<i32>
div_floor_i32_checked(a: i32, b: i32) -> MathResult<i32>
div_ceil_i32_checked(a: i32, b: i32) -> MathResult<i32>
```

`MathResult<T>` contiene `value` y `error`. Importe las constantes de error de `std::math::result` y use `value` solo después de comprobar `error == MATH_ERROR_NONE`. El valor absoluto del mínimo entero con signo no se puede representar con el mismo tipo, por lo que la función de valor absoluto con comprobación informa de un error. `clamp` rechaza `lo > hi`. La división comprueba si el divisor es cero o si el resultado desborda el rango. `floor` y `ceil` redondean de forma distinta a la división entera, que trunca hacia cero.

## Clasificación de valores de punto flotante

`is_nan_f64`, `is_infinite_f64` y `is_finite_f64` en `std::math::float` distinguen valores especiales. También se proporciona la función f32. `float_to_bits_f64(value: f64) -> u64` es una función para obtener bits de almacenamiento y es diferente de la conversión numérica de `value as u64`. NaN ni siquiera es igual a sí mismo, por lo que no se compara con `value == nan`.

## diagnostico

`debug_assert(condition: bool, message: str)` en `std::debug::core` finaliza después de diagnosticar una condición falsa. Las situaciones que normalmente pueden fallar, como la entrada del usuario, se manejan con el valor del resultado y assert se usa al verificar las condiciones internas del programa que deben cumplirse.

Guarde el programa a continuación como `main.wave` y ejecútelo.

<!-- wave-example: math-api -->
```wave
import("std::math::int")::{
    min_i32, max_i32
};
import("std::debug::core")::{
    debug_assert
};

fun main() {
    var low: i32 = min_i32(9, 4);
    var high: i32 = max_i32(9, 4);
    debug_assert(low <= high, "invalid range");
    println("{} {}", low, high);
}
```

Resultado de la ejecución:

```text
4 9
```

## Dirección de redondeo para división negativa

Compare cómo dividir -7 entre 3. `/` se trunca hacia 0 para convertirse en -2. floor selecciona el entero más pequeño -3 y ceil selecciona el entero más grande -2.

<!-- wave-example: book-math-rounding -->
```wave
import("std::math::int")::{
    div_floor_i32_checked,
    div_ceil_i32_checked,
    abs_i32_checked
};
import("std::math::result")::{
    MathResult,
    MATH_ERROR_NONE,
    MATH_ERROR_OVERFLOW
};

fun main() -> i32 {
    var floor: MathResult<i32> = div_floor_i32_checked(-7, 3);
    var ceil: MathResult<i32> = div_ceil_i32_checked(-7, 3);

    if (floor.error != MATH_ERROR_NONE || ceil.error != MATH_ERROR_NONE) {
        return 1;
    }

    println("truncate={} floor={} ceil={}", -7 / 3, floor.value, ceil.value);

    var absolute: MathResult<i32> = abs_i32_checked(-2147483648);

    if (absolute.error == MATH_ERROR_OVERFLOW) {
        println("absolute value is outside i32");
    }

    return 0;
}
```

Resultado de la ejecución:

```text
truncate=-2 floor=-3 ceil=-2
absolute value is outside i32
```

Si solo comprueba el valor del resultado, no podrá distinguir entre el valor sustitutivo incluido en caso de fallo y el resultado del cálculo real. Siga primero el orden de verificación error. floor es útil cuando se colocan coordenadas negativas en un intervalo de cierto tamaño, y ceil es útil cuando se redondea el número requerido de paquetes.

## ¿Qué errores debo manejar?

|situación|error|Ejemplo de procesamiento|
| --- | --- | --- |
|dividir por cero| `MATH_ERROR_DIVIDE_BY_ZERO` |Toma de nuevo el denominador|
|El resultado no se puede almacenar en tipo| `MATH_ERROR_OVERFLOW` |Calcule con un tipo más amplio o rechace la entrada|
|El mínimo es mayor que el máximo en clamp| `MATH_ERROR_INVALID_ARGUMENT` |Modificar rango de configuración|

assert no es una herramienta de reparación de errores. Los errores en la entrada del usuario se manejan mediante declaraciones condicionales y valores de retorno, y después de completar el cálculo, las condiciones internas que deben cumplirse se verifican con debug_assert.
