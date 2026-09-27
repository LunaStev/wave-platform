---
translation_set_id: learn-errors
path: language/errors
locale: es
group: language
group_order: 2
order: 12
title: 12. Representar y recuperarse de errores
summary: Separe los errores de los valores normales y limpie los recursos de las rutas de error.
---

## Resultado de la función de grado de falla

Los archivos faltantes o las entradas fuera de límites ocurren naturalmente en los programas. Manejar un error es más que simplemente imprimir un mensaje. Es el proceso de distinguir fallas, verificar el estado del trabajo ya realizado, organizar los recursos adquiridos y luego elegir si continuar o terminar.

En este capítulo, comenzamos con la representación del error de una función pequeña y avanzamos hacia la estructura de resultados, variant, el retorno anticipado y la limpieza de recursos.

## Los marcadores de fracaso no deben superponerse a los valores de éxito.

La razón por la que -1 se utiliza como no encontrado en una búsqueda de matriz es porque el índice efectivo es 0 o más. Por otro lado, en cálculos en los que cualquier número entero puede ser un resultado normal, si -1 se designa como error, no se puede distinguir del valor normal -1.

0 también es un valor comúnmente mal entendido. La longitud de cadena vacía 0, la primera posición 0 y el número de bytes transferidos 0 tienen diferentes significados para diferentes funciones. No juzgue el éxito solo porque el valor de retorno no es cero.

## Devolver el éxito y el valor juntos

<!-- wave-example: book-error-result -->
```wave
struct Division {
    ok: bool;
    value: i32;
}

fun divide_nonnegative(left: i32, right: i32) -> Division {
    if (left < 0 || right <= 0) {
        return Division {
            ok: false,
            value: 0
        };
    }

    return Division {
        ok: true,
        value: left / right
    };
}

fun main() {
    var result: Division = divide_nonnegative(0, 3);

    if (!result.ok) {
        println("invalid input");
        return;
    }

    println("value={}", result.value);
}
```

Resultado de la ejecución:

```text
value=0
```

Un resultado normal también puede ser 0. En lugar de mirar value y adivinar si fue exitoso, verifique ok primero. Incluso si el campo value existe en el resultado fallido, no significa que sea el valor a utilizar.

Las reglas de entrada para este ejemplo son que el lado izquierdo es 0 o mayor y el lado derecho es positivo. El alcance se muestra en el nombre y la descripción de la función para distinguirlo de la típica división entera signed.

## Distinguir entre causas de errores.

Agregue información de error para obtener orientación adicional o recuperación según el motivo del error. También hay una manera de dividir la función en una función que verifica el rango de entrada y una función que lo calcula.

<!-- wave-example: book-error-codes -->
```wave
struct Check {
    ok: bool;
    error: i32;
}

fun check_quantity(value: i32) -> Check {
    if (value < 1) {
        return Check {
            ok: false,
            error: 1
        };
    }

    if (value > 100) {
        return Check {
            ok: false,
            error: 2
        };
    }

    return Check {
        ok: true,
        error: 0
    };
}

fun main() {
    var result: Check = check_quantity(101);

    if (!result.ok) {
        match (result.error) {
            1 => {
                println("quantity is too small");
            }
            2 => {
                println("quantity is too large");
            }
            _ => {
                println("invalid quantity");
            }
        }
    }
}
```

Resultado de la ejecución:

```text
quantity is too large
```

El significado de los números está definido por esta función. No se puede considerar igual que el error número 1 o 2 en otras bibliotecas. En público API, nombrar constantes o tipos de error evita que la persona que llama tenga que memorizar números aleatorios.

## Separe los resultados con variant

La relación de que un valor de éxito y un error no pueden existir al mismo tiempo se puede expresar como variant.

<!-- wave-example: book-error-variant -->
```wave
variant ByteResult {
    Value(u8),
    OutOfRange(i32)
}

fun to_byte(value: i32) -> ByteResult {
    if (value < 0 || value > 255) {
        return ByteResult::OutOfRange(value);
    }

    return ByteResult::Value(value as u8);
}

fun main() {
    var result: ByteResult = to_byte(300);

    match (result) {
        ByteResult::Value(value) => {
            println("byte={}", value);
        }
        ByteResult::OutOfRange(original) => {
            println("out of range: {}", original);
        }
    }
}
```

Resultado de la ejecución:

```text
out of range: 300
```

El error contenía el valor de entrada original. Es más fácil para la persona que llama explicar el problema que simplemente devolver false. La naturaleza de los datos también se tiene en cuenta al no dejar entradas confidenciales como contraseñas o tokens en el registro.

## Haga que las rutas normales sean más fáciles de leer con retornos anticipados

No es necesario colocar todo el código bueno en el interior de if al realizar una verificación de varios pasos. Si falla, puede regresar inmediatamente y continuar el camino exitoso a continuación.

<!-- wave-example: book-error-early-return -->
```wave
fun show_total(quantity: i32, price: i32) {
    if (quantity < 1 || quantity > 100) {
        println("invalid quantity");
        return;
    }

    if (price < 0 || price > 100000) {
        println("invalid price");
        return;
    }

    var total: i32 = quantity * price;
    println("total={}", total);
}

fun main() {
    show_total(3, 1200);
    show_total(0, 1200);
    show_total(3, -1);
}
```

Resultado de la ejecución:

```text
total=3600
invalid quantity
invalid price
```

Después de pasar la prueba, puede aprovechar el hecho de que quantity y price están dentro del rango especificado. El rango se estableció de modo que la multiplicación intermedia también esté dentro de i32. Al agregar una devolución anticipada, también debe verificar si ya posee algún recurso en ese momento.

## Limpieza de memoria de rutas de falla

<!-- wave-example: book-error-cleanup -->
```wave
import("std::buffer::alloc")::{
    Buffer, buffer_init, buffer_free
};
import("std::buffer::write")::{
    buffer_append_str
};

fun make_message(allowed: bool) -> i32 {
    var data: Buffer;

    if (buffer_init(&data, 16) < 0) {
        return 1;
    }

    if (!allowed) {
        buffer_free(&data);
        return 2;
    }

    if (buffer_append_str(&data, "ready") < 0) {
        buffer_free(&data);
        return 3;
    }

    println("bytes={}", data.len);

    if (buffer_free(&data) < 0) {
        return 4;
    }

    return 0;
}

fun main() {
    println("status={}", make_message(false));
}
```

Resultado de la ejecución:

```text
status=2
```

Libere Buffer incluso en rutas de falla que no generen datos. En lugar de convertir cada falla en un solo return, es importante administrar claramente lo que se posee en cada sucursal.

También hay API donde falla la limpieza. Decida cómo preservará los errores y los errores de limpieza en el trabajo original. Este pequeño ejemplo primero devuelve el número de error de la tarea. Los programas más grandes pueden grabar ambos por separado.

## Los éxitos parciales no se cancelan automáticamente

Si escribe algunos bytes en un archivo y luego la escritura falla, los bytes ya escritos no se pierden. Es posible que el otro extremo de la red también haya recibido algunos datos. Repetir la misma tarea desde el principio puede generar registros duplicados.

Por el contrario, la lectura checked byte cursor conserva la posición y el valor de salida cuando falla. Estas funciones pueden recibir más entradas e intentarlo nuevamente en la misma ubicación. En lugar de aplicar la regla "si falla, nada cambia" a cada API, consulte la documentación de esa función.

## Errores y trampas recuperables

Una ruta de archivo no válida o una entrada de usuario no válida se pueden informar como un valor de error para que la persona que llama pueda recuperarse. Una trampa causada por un recuento de cambios de tiempo de ejecución no válido o una conversión de punto flotante a entero es un mecanismo diferente.

Los programas que deben continuar ejecutándose deben inspeccionar sus entradas antes de operaciones peligrosas. Tampoco se pretende abusar de assert como medio para manejar fallas elegantes en la entrada del usuario. Si necesita darle al usuario la oportunidad de volver a ingresar su entrada, devuelva los resultados para continuar con el flujo de control.

## Ejercicio y solución completa

Escriba una función que lea un elemento de matriz por índice y falle cuando el índice sea negativo o mayor o igual a la longitud. Utilice una estructura de resultados porque cero puede ser un valor de elemento válido.

<!-- wave-example: book-error-solution -->
```wave
struct Lookup {
    ok: bool;
    value: i32;
}

fun at(values: ptr<i32>, length: i32, index: i32) -> Lookup {
    if (index < 0 || index >= length) {
        return Lookup {
            ok: false,
            value: 0
        };
    }

    return Lookup {
        ok: true,
        value: deref values[index]
    };
}

fun main() {
    var values: array<i32, 3> = [0, 10, 20];
    var first: Lookup = at(&values[0], 3, 0);
    var outside: Lookup = at(&values[0], 3, 3);

    if (first.ok) {
        println("value={}", first.value);
    }

    if (!outside.ok) {
        println("out of bounds");
    }
}
```

Resultado de la ejecución:

```text
value=0
out of bounds
```

La condición de la persona que llama es que el puntero y length representen la matriz legible real. Esta no es una función que hace que las direcciones aleatorias sean seguras con solo verificar el índice. Lea por separado de qué comprobaciones es responsable la función y qué condiciones debe garantizar la persona que llama.
