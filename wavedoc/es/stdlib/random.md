---
translation_set_id: stdlib-random
path: stdlib/random
locale: es
group: stdlib
group_order: 1
order: 10
title: random: Llenar buffers con aleatoriedad del sistema operativo
summary: OS Llene el búfer con entropía y maneje fallas parciales.
---

## API

```text
std::random::fill
random_available() -> bool
random_fill(buffer: ptr<u8>, size: i64) -> RandomFillResult
```

El tamaño es un recuento de bytes y la persona que llama proporciona el almacenamiento. `RandomFillResult` contiene ok, escrito y error. En caso de éxito, escrito es igual a la longitud solicitada. En caso de error, escrito identifica el prefijo válido y completo; no utilice los bytes restantes como datos aleatorios.

`random_available` le indica si se admite la función de número aleatorio OS. El éxito de una solicitud individual se verifica mediante los resultados de random_fill. Utiliza solo entropía OS y no recurre al valor de tiempo ni a PRNG débil en caso de falla. size=0 tendrá éxito incluso si se pasa junto con null. null es un error para longitudes negativas o positivas.

## Ejemplo de ejecución

<!-- wave-example: random-api -->
```wave
import("std::random::fill")::{
    RandomFillResult, random_fill
};

fun main() -> i32 {
    var bytes: array<u8, 16>;
    var result: RandomFillResult = random_fill(&bytes[0], 16);
    if (!result.ok) {
        println("random error={}", result.error);
        return 1;
    }

    println("filled={}", result.written);
    return 0;
}
```

Resultado de la ejecución:

```text
filled=16
```

Guárdelo como `main.wave` y ejecútelo. El contenido de los bytes es diferente cada vez, por lo que no se espera ningún valor específico. Si falla, verifique la causa con result.error. En lugar de generar literalmente bytes aleatorios, utilice una codificación separada si es necesario.

## Si su solicitud es incorrecta

Una solicitud de cero bytes tiene éxito porque no es necesario escribir nada. Pasar null con una longitud positiva falla porque no hay un búfer de destino. El siguiente programa compara estos casos sin asignar memoria.

<!-- wave-example: book-random-boundaries -->
```wave
import("std::random::fill")::{
    RandomFillResult,
    random_fill
};

fun main() -> i32 {
    var empty: RandomFillResult = random_fill(null, 0);

    if (!empty.ok || empty.written != 0) {
        return 1;
    }

    println("empty request succeeded");

    var invalid: RandomFillResult = random_fill(null, 16);

    if (invalid.ok) {
        return 2;
    }

    println("missing buffer rejected");
    return 0;
}
```

Resultado de la ejecución:

```text
empty request succeeded
missing buffer rejected
```

## Manejo de buffers parcialmente llenos

Si se solicitaron 16 bytes pero fallaron y se devuelve written=8, solo se llenan los primeros 8 bytes. Si la tarea es crear un identificador de 16 bytes, no es un identificador exitoso, por lo que descartamos todo el resultado e informamos un error. No debe llenar los 8 bytes restantes con 0 y luego tratarlo como un éxito.

La asignación de bytes aleatorios a un rango de números enteros requiere cuidado. La aplicación de `% 10` a valores de u8 distribuidos uniformemente hace que 0–5 sea más probable que 6–9, porque 256 no es divisible por 10. Para eliminar este sesgo, rechace los valores 250–255, dibuje nuevamente y aplique la operación restante solo a los valores aceptados.

La persona que llama gestiona el almacenamiento y la vida útil de los bytes aleatorios. Cuando se usa una matriz, se procesa dentro del alcance de la matriz y cuando se usa memoria dinámica, se libera después de su uso. En cambio, la estructura de retorno no es propietaria del búfer.
