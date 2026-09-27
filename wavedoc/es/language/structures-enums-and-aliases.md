---
translation_set_id: data-types
path: language/structures-enums-and-aliases
locale: es
group: language
group_order: 2
order: 8
title: 8. Estructuras, enumeraciones y variantes
summary: Obtenga información sobre los campos, la inicialización de estructuras y las funciones de enum y variant.
---

## Expresar relaciones entre datos como tipos.

Si el precio y la cantidad de productos se pasan como variables, es difícil saber con solo mirar el código si los dos valores pertenecen al mismo producto. Una estructura agrupa campos relacionados. enum representa un estado con nombre y variant almacena diferentes datos juntos para cada caso.

Las tres funciones no son sintaxis que se reemplazan entre sí. Elige dependiendo de lo que quieras expresar.

|algo que expresar|seleccionar|si|
| --- | --- | --- |
|Múltiples campos existentes simultáneamente| struct |Precio unitario y cantidad de producto.|
|Estado entero con nombre| enum |Esperar/Continuar/Completar|
|Datos diferentes en cada caso| variant |valor de éxito o error|

## Declaración de estructura y creación de valor.

<!-- wave-example: book-struct-first -->
```wave
struct Product {
    price: i32;
    quantity: i32;
}

fun main() {
    var item: Product = Product {
        price: 1500,
        quantity: 2
    };

    println("{} {}", item.price, item.quantity);
}
```

Resultado de la ejecución:

```text
1500 2
```

Los campos en las declaraciones terminan con punto y coma, y al crear valores, los campos y valores se conectan con dos puntos y se separan por comas. Definir el tipo y crear el valor real son dos pasos diferentes. Declarar el tipo Product no crea automáticamente espacio de almacenamiento para un producto.

Se accede al campo como `item.price`. Si crea varios item del mismo tipo, puede almacenar valores diferentes en cada uno.

## Pasar una estructura a una función.

<!-- wave-example: book-struct-function -->
```wave
struct Product {
    price: i32;
    quantity: i32;
}

fun subtotal(item: Product) -> i32 {
    return item.price * item.quantity;
}

fun main() {
    var item: Product = Product {
        price: 1500,
        quantity: 2
    };

    println("{}", subtotal(item));

    item.quantity = 3;
    println("{}", subtotal(item));
}
```

Resultado de la ejecución:

```text
3000
4500
```

La función recibe como tipo la relación de que el precio unitario y la cantidad pertenecen al mismo producto. Esta es una función que lee la estructura del campo entero pasada como valor. Cualquier función que quiera modificar el almacenamiento de la persona que llama puede diseñarse para aceptar un puntero.

Si la estructura tiene campos de puntero, copiar el valor también copia la dirección. Esta no es una función para realizar copias profundas en asignaciones separadas. Los tipos que contienen recursos como identificadores de archivos o Buffer deben especificar reglas de copia y liberación juntas.

## Método y proto

Las funciones relacionadas se pueden agrupar en forma de método. proto es un método para escribir el método de una estructura como un bloque separado.

<!-- wave-example: book-struct-method -->
```wave
struct Counter {
    value: i32;
}

proto Counter {
    fun current(self: Counter) -> i32 {
        return self.value;
    }
}

fun main() {
    var counter: Counter = Counter {
        value: 3
    };

    println("{}", counter.current());
}
```

Resultado de la ejecución:

```text
3
```

`self: Counter` es un parámetro que recibe un valor. El uso de una notación de llamada a un método no lo convierte automáticamente en un método que modifica el original. Lean juntos el tipo de self y lo que hace en el texto.

No se puede esperar que un campo tenga solo un estado válido solo porque le adjuntó un método. Si hay combinaciones no válidas que el usuario puede crear con campos públicos, su función debe verificarlas o proporcionar una regla de creación.

## Nombra el estado con enum

<!-- wave-example: book-enum-state -->
```wave
enum State -> i32 {
    Ready = 0,
    Running,
    Finished,
}

fun main() {
    var state: State = State::Ready;

    if (state == State::Ready) {
        println("ready");
    }

    state = State::Running;

    if (state == State::Running) {
        println("running");
    }
}
```

Resultado de la ejecución:

```text
ready
running
```

`-> i32` es un tipo entero utilizado en expresiones. El primer valor se establece en 0 y los valores omitidos posteriores son 1 mayor que el valor anterior. En lugar de simplemente comparar 0 y 1 en su código, usar State::Ready y State::Running revelará el significado.

enum Tener un nombre no restringe automáticamente las transiciones de estado. Reglas como si es posible regresar de Finished a Running deben implementarse como funciones.

## Conectando casos y datos con variant

Si hay un valor solo en caso de éxito y se necesita información de error en caso de falla, se puede expresar como variant.

<!-- wave-example: book-variant-result -->
```wave
variant Result {
    Value(i32),
    Error(i32)
}

fun divide(left: i32, right: i32) -> Result {
    if (left < 0 || right <= 0) {
        return Result::Error(1);
    }

    return Result::Value(left / right);
}

fun show(result: Result) {
    match (result) {
        Result::Value(value) => {
            println("value={}", value);
        }
        Result::Error(code) => {
            println("error={}", code);
        }
    }
}

fun main() {
    show(divide(12, 3));
    show(divide(12, 0));
}
```

Resultado de la ejecución:

```text
value=4
error=1
```

Result::Value y Result::Error contienen cada uno payload. Incluso si contienen el mismo tipo de entero, se distinguen ciertos casos. La persona que llama verifica el caso con match y usa payload dentro de ese arm.

divide arriba es un pequeño ejemplo que no admite operandos negativos. Debido a que se especifica el rango de entrada, no debe confundirse con una función que cubre todos los límites de la división regular signed.

## Si payload no está presente

No se requieren datos en todos los casos. El estado de falta de valor se puede expresar como un caso aparte.

<!-- wave-example: book-variant-empty -->
```wave
variant Lookup {
    Found(i32),
    Missing
}

fun main() {
    var result: Lookup = Lookup::Missing;

    match (result) {
        Lookup::Found(value) => {
            println("found={}", value);
        }
        Lookup::Missing => {
            println("missing");
        }
    }
}
```

Resultado de la ejecución:

```text
missing
```

En lugar de reservar un entero arbitrario como "ninguno", utilizamos el caso Missing. No importa cuál sea el valor del éxito, el significado no se superpone.

match a `_` manejan los casos restantes. Si desea que cada persona que llama sea revisada nuevamente cuando se agrega un nuevo caso, es mejor separar todos los casos explícitamente. Cualquiera que sea el método que elija, asegúrese de que no haya entradas sin procesar.

## Usando estructura y variant juntos

La selección de diferentes datos se puede expresar como variant, y la agrupación de varios campos que pertenecen a un caso se puede expresar como una estructura. Por ejemplo, si el resultado del procesamiento del pedido es exitoso, se puede diseñar para que contenga una estructura de recibo, y si es un error, se puede diseñar para que contenga un número de error.

La regla de duración de la memoria no desaparece incluso si el valor contiene otros valores. Si un puntero se almacena en variant, si el puntero es válido y quién lo liberará se determina por separado. Incluso cuando guarde en un formato de archivo externo, debe especificar la codificación para cada campo en lugar de volcar la memoria de la estructura tal como está.

## Ejercicio y solución completa

Cree un veredicto que solo acepte puntuaciones entre 0 y 100. Las puntuaciones válidas se expresan como Grade(score), el resto se expresan como Invalid.

<!-- wave-example: book-model-solution -->
```wave
variant CheckedScore {
    Grade(i32),
    Invalid
}

fun check_score(score: i32) -> CheckedScore {
    if (score < 0 || score > 100) {
        return CheckedScore::Invalid;
    }

    return CheckedScore::Grade(score);
}

fun main() {
    var result: CheckedScore = check_score(87);

    match (result) {
        CheckedScore::Grade(score) => {
            println("accepted={}", score);
        }
        CheckedScore::Invalid => {
            println("invalid score");
        }
    }
}
```

Resultado de la ejecución:

```text
accepted=87
```

Cambie la entrada a -1, 0, 100, 101 para verificar los límites. Dado que el éxito y el fracaso no comparten el mismo espacio entero, la persona que llama reduce el riesgo de agregar accidentalmente valores de error al cálculo promedio.


## ¿Qué expresión debo elegir?

|forma de datos|expresión adecuada|si|
| --- | --- | --- |
|Múltiples valores del mismo tipo|matriz|10 puntos|
|Múltiples campos relacionados entre sí|estructura|nombre y puntuación|
|Valor de estado nombrado| enum | Ready, Running, Stopped |
|Datos adicionales que varían según el estado| variant | Value(i32), Error(str) |
|Nombre contextual para un tipo existente|escriba alias| UserId = u64 |

Al elegir una estructura de datos, considere no solo los valores que desea almacenar, sino también los estados falsos que puede representar. Una estructura con campos de éxito/fracaso y valor/error puede crear una combinación incorrecta, pero variant se puede diferenciar en payload para cada caso.

## Ubicación de la memoria y datos externos.

La memoria de la estructura puede contener espacios vacíos para garantizar la alineación entre campos. La simple suma de los tamaños de los campos no siempre equivale al tamaño total de la estructura. Cuando necesite conocer el tamaño y la alineación, utilice [mem Función de diseño](/docs/es/reference/memory-and-buffer).

Para archivos o mensajes de red, el orden de los bytes y la longitud se pueden determinar claramente codificando los campos en orden con la función [bytes](/docs/es/stdlib/bytes). Al pasar estructuras a otros idiomas, alinee la declaración externa de [FFI](/docs/es/language/modules-imports-and-ffi) con el objetivo ABI.

[Estructurar el aprendizaje y la práctica.](/docs/es/language/structures-enums-and-aliases) · [variant](/docs/es/language/variants)
