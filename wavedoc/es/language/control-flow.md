---
translation_set_id: control-flow
path: language/control-flow
locale: es
group: language
group_order: 2
order: 4
title: 4. Condiciones, bucles y valores límite.
summary: Aprenda if, for, while y la variedad de declaraciones repetitivas.
---

## Seleccionar ruta de ejecución

El programa del capítulo anterior ejecutó declaraciones de arriba a abajo. Un programa real debe hacer cosas diferentes según la entrada y el estado. Las declaraciones condicionales seleccionan una ruta de ejecución y los bucles aplican la misma regla a varios valores.

Guarde cada ejemplo en main.wave y ejecútelo. Al leer el código, escriba en una hoja de papel los valores de las variables actuales, las condiciones que se probarán a continuación y el orden de las declaraciones que se ejecutarán. Es más importante practicar siguiendo la corriente que memorizar los resultados.

## if y else

Las condiciones están escritas entre paréntesis y el texto entre llaves. En el siguiente ejemplo, reemplace balance con 500 o 2000 para determinar qué rama se ejecuta.

<!-- wave-example: book-if-else -->
```wave
fun main() {
    var balance: i32 = 2000;
    var price: i32 = 1200;

    if (balance >= price) {
        balance -= price;
        println("bought, balance={}", balance);
    } else {
        println("not enough money");
    }
}
```

Resultado de la ejecución:

```text
bought, balance=800
```

No está ejecutando ambos bloques. Si la condición es verdadera, se ejecuta el primer bloque. Si la condición es falsa, se ejecuta el bloque else. Se permiten compras incluso cuando balance es igual a price, así que escribí `>=`. Si lo cambias a `>`, la operación cambiará por el mismo importe.

## Orden de múltiples condiciones

Puede concatenar condiciones con else if. Dado que primero ejecutamos una rama satisfecha desde arriba, debemos considerar si queremos verificar primero el límite más grande o el límite más pequeño.

<!-- wave-example: book-grade -->
```wave
fun grade(score: i32) {
    if (score < 0 || score > 100) {
        println("invalid");
    } else if (score >= 90) {
        println("A");
    } else if (score >= 80) {
        println("B");
    } else {
        println("C");
    }
}

fun main() {
    grade(95);
    grade(80);
    grade(79);
    grade(101);
}
```

Resultado de la ejecución:

```text
A
B
C
invalid
```

Las puntuaciones no válidas primero se rechazan y luego se califican. Si pones `score >= 80` al principio, nunca llegarás a la rama A porque 95 grados entran en esa rama. Verifique el orden de las condiciones así como la exactitud de cada condición.

## No cambiar valores en expresiones condicionales.

Las operaciones de asignación, asignación compuesta y de incremento o decremento no están permitidas en condiciones if, while o for. Utilice `==` para comparar. Para actualizar un valor y luego probarlo, escriba dos declaraciones separadas.

Forma correcta de un fragmento dentro de una función:

```wave
value = read_value();

if (value == expected) {
    println("matched");
}
```

read_value y expected no están definidos en este fragmento, por lo que no es el programa completo el que se ejecuta tal cual. La regla que se muestra aquí es "Comparar después del cambio de estado".

## while: Mientras se mantengan las condiciones

<!-- wave-example: book-while-countdown -->
```wave
fun main() {
    var remaining: i32 = 3;

    while (remaining > 0) {
        println("{}", remaining);
        remaining -= 1;
    }

    println("finished at {}", remaining);
}
```

Resultado de la ejecución:

```text
3
2
1
finished at 0
```

Se verifican las condiciones antes de ingresar al cuerpo. Si remaining es 0 desde el principio, el cuerpo nunca se ejecuta. Si se omite la disminución al final del cuerpo, la condición sigue siendo verdadera y el ciclo no termina.

Después de escribir el bucle, marque "¿Qué lo acerca a la condición de terminación?" Si es un bucle esperando entrada, el cambio de entrada o EOF desempeña su función, y si es un bucle numérico, la actualización del índice desempeña su función.

## for: Inicialización/Condiciones/Actualización

for expresa las tres partes necesarias para la repetición.

<!-- wave-example: book-for-sum -->
```wave
fun main() {
    var total: i32 = 0;

    for (var number: i32 = 1; number <= 5; number += 1) {
        total += number;
    }

    println("sum={}", total);
}
```

Resultado de la ejecución:

```text
sum=15
```

1. Inicialice number en 1. Este paso es una sola vez.
2. Comprueba number <= 5. Si es falso, finaliza la iteración.
3. Agregue number a total en el texto.
4. Aumente number en 1 y vuelva a la verificación de condición.

No se supone que las variables de repetición declaradas en for puedan usarse después de la repetición. Si su diseño requiere un valor después de la iteración, declarelo afuera y aclare la ubicación de inicialización.

## Límites de inclusión y exclusión

La suma natural de 1 a 5 es `<= 5`. Por otro lado, el índice de una matriz de longitud 5 debería usar `< 5`. Esto se debe a que los índices de una matriz comienzan en 0 y terminan en 4.

No confunda "ejecutar cinco veces" con "hasta el valor 5 inclusive". Puede determinar el número de repeticiones observando los valores inicial y final juntos. Los casos en los que la entrada está vacía y tiene un solo elemento son buenos para detectar errores de límites.

## continue y break

continue omite el resto de esta iteración y break finaliza la iteración más cercana.

<!-- wave-example: book-loop-control -->
```wave
fun main() {
    var total: i32 = 0;

    for (var number: i32 = 1; number <= 10; number += 1) {
        if (number == 3) {
            continue;
        }

        if (number == 6) {
            break;
        }

        total += number;
    }

    println("{}", total);
}
```

Resultado de la ejecución:

```text
12
```

Los números en el total son 1, 2, 4 y 5. Si se encuentra continue en for, el proceso de renovación continuará. Dado que while no tiene una expresión de actualización separada como for, debe tener cuidado de no omitir ningún cambio de estado requerido antes de continue.

Si los bucles están anidados, uno break no completa todos los bucles. Si necesita detenerse en varias etapas, incluya el trabajo en una función e indique su intención usando return o verificando también la condición de terminación en la iteración externa.

## Divida el caso por match

Puede utilizar match al comparar varios casos del mismo valor. El cuerpo de cada arm es un bloque.

<!-- wave-example: book-match-number -->
```wave
fun describe(status: i32) {
    match (status) {
        200 => {
            println("ok");
        }
        404 => {
            println("missing");
        }
        _ => {
            println("other");
        }
    }
}

fun main() {
    describe(200);
    describe(404);
    describe(500);
}
```

Resultado de la ejecución:

```text
ok
missing
other
```

`_` es el patrón que se encarga del resto. No coloque duplicados dentro del mismo match. variant, que tiene diferentes tipos de datos según el valor, se trata en [Capítulo del modelo de datos](/docs/es/language/structures-enums-and-aliases).

## Ejemplo completo: Contar números que cumplen condiciones

Encuentra el número y la suma de números pares del 1 al 10. Dado que el recuento y la suma son información diferente, se acumulan como variables.

<!-- wave-example: book-loop-statistics -->
```wave
fun main() {
    var count: i32 = 0;
    var total: i32 = 0;

    for (var number: i32 = 1; number <= 10; number += 1) {
        if (number % 2 == 0) {
            count += 1;
            total += number;
        }
    }

    println("count={} total={}", count, total);
}
```

Resultado de la ejecución:

```text
count=5 total=30
```

Los números pares son 2, 4, 6, 8 y 10, por lo que el número es 5 y la suma es 30. Incluso si la expresión es corta, es fácil verificar el límite de la iteración si primero verifica el resultado en un rango pequeño que se puede obtener a mano.

## Ejercicio y solución completa

Sume solo múltiplos de 3 del 1 al 20, pero no agregue ningún valor que sume más de 30. Necesitamos distinguir entre “sumar y luego verificar para ver si terminó” y “verificar si terminó y luego sumar”.

<!-- wave-example: book-loop-solution -->
```wave
fun main() {
    var total: i32 = 0;

    for (var number: i32 = 1; number <= 20; number += 1) {
        if (number % 3 != 0) {
            continue;
        }

        if (total + number > 30) {
            break;
        }

        total += number;
    }

    println("{}", total);
}
```

Resultado de la ejecución:

```text
30
```

Los valores 3+6+9+12 suman 30, por lo que el siguiente valor, 15, no se suma. Estas pequeñas entradas son seguras, pero para números enteros grandes la verificación `total + number` puede desbordarse. Emitir un cheque no maneja automáticamente todos los casos límite.
