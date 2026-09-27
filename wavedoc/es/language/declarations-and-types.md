---
translation_set_id: types
path: language/declarations-and-types
locale: es
group: language
group_order: 2
order: 2
title: 2. Variables, tipos y alcance
summary: Aprenda las variables locales, el alcance de los enteros, bool y los alcances.
---

## Tratar los valores por su nombre

Si escribe sus precios directamente en varios lugares, tendrá que encontrarlos todos cuando cambie el precio. Una variable es un espacio de almacenamiento que da nombres a los valores y le permite leerlos y cambiarlos usando esos nombres. En este capítulo, aprenderá sobre el alcance de las declaraciones, asignaciones, tipos y el alcance de los nombres visibles dentro de los bloques.

Los programas siguientes son cada uno parte de un main.wave independiente. Guarde un ejemplo, ejecútelo como `wavec run main.wave` y reemplácelo con el siguiente ejemplo.

## Declaración e inicialización

<!-- wave-example: book-variable-first -->
```wave
fun main() {
    var price: i32 = 1200;
    var quantity: i32 = 3;
    println("price={}", price);
    println("quantity={}", quantity);
    println("total={}", price * quantity);
}
```

Resultado de la ejecución:

```text
price=1200
quantity=3
total=3600
```

Lea la declaración en cuatro partes.

|parte|este ejemplo|papel|
| --- | --- | --- |
|palabra clave de declaración| var |Crear variable local|
|nombre| price |Identificador que se utilizará más adelante|
|tipo| i32 |Tipo y rango de valores a almacenar|
|valor inicial| 1200 |Primer valor a almacenar|

Los dos puntos antes del tipo y el signo igual antes del valor inicial tienen funciones diferentes. Haz que tu nombre tenga significado. En este ejemplo, price es el precio unitario y quantity es la cantidad. Incluso si se usa indistintamente el mismo i32, se pueden realizar cálculos incorrectos sin errores gramaticales.

## La asignación no es una fórmula que mantenga las relaciones.

Si guarda el resultado del cálculo en una variable, se ingresará el valor en ese punto. No recuerda los cálculos y los vuelve a evaluar automáticamente más tarde.

<!-- wave-example: book-assignment-snapshot -->
```wave
fun main() {
    var price: i32 = 1200;
    var quantity: i32 = 3;
    var total: i32 = price * quantity;
    quantity = 5;
    println("before recalculation={}", total);
    total = price * quantity;
    println("after recalculation={}", total);
}
```

Resultado de la ejecución:

```text
before recalculation=3600
after recalculation=6000
```

`quantity = 5` escribe un nuevo valor en una variable ya existente. Debe distinguirse de una redeclaración como `var quantity`. `total` también es 3600 antes de sustituirlo nuevamente. Si un programa debe mantener relaciones entre múltiples variables, debe escribirse para realizar cálculos cuando las relaciones cambien.

## Calcular el siguiente valor a partir del valor anterior

Primero calcule el lado derecho de la declaración de asignación y escriba el resultado en el espacio de almacenamiento de la izquierda. A diferencia de las ecuaciones matemáticas, `count = count + 1` es una actualización válida.

<!-- wave-example: book-update-variable -->
```wave
fun main() {
    var count: i32 = 0;
    count = count + 1;
    count += 2;
    count *= 3;
    println("{}", count);
}
```

Resultado de la ejecución:

```text
9
```

La progresión de valores es 0 → 1 → 3 → 9. `+=`, `*=` expresan cálculo y almacenamiento juntos. Si anota el orden de los cálculos, podrá descubrir en qué etapa sus pensamientos fueron diferentes cuando los resultados son diferentes de lo que esperaba.

## Ancho y signo de tipos enteros.

Si comienza con `i`, se convierte en signed, y si comienza con `u`, se convierte en unsigned. El último número es el número de bits. A medida que aumenta el número de bits, aumenta el rango que se puede expresar y también aumenta el espacio de almacenamiento.

|tipo|valor mínimo|valor máximo|Uso de ejemplo|
| --- | --- | --- | --- |
| i8 | -128 | 127 |pequeño valor con signo|
| u8 | 0 | 255 |un byte|
| i16 | -32768 | 32767 |datos enteros pequeños|
| u16 | 0 | 65535 |Puerto/campo de 16 bits|
| i32 | -2147483648 | 2147483647 |Cálculos comunes con números enteros pequeños|
| u32 | 0 | 4294967295 |campo de bits de 32 bits|

Wave también proporciona enteros con y sin signo de 64, 128, 256, 512 y 1024 bits. Un tipo más amplio no hace que todos los cálculos sean seguros: el resultado aún puede exceder el rango seleccionado. Elija primero el rango que necesita. `isz` y `usz` siguen el ancho de la dirección de destino.

Se debe hacer una distinción entre almacenar un literal grande en un tipo pequeño y descartar bits intencionalmente haciendo cast. La conversión a un tipo más pequeño simplemente para eliminar errores puede cambiar el valor en sí. Las transformaciones se tratan en el próximo capítulo.

## Tipos de punto flotante y bool

`f32`·`f64` son números de coma flotante. A diferencia de los números enteros, pueden representar partes decimales, pero no pueden almacenar todos los números decimales exactamente. Ésta es una de las razones por las que los importes se gestionan en unidades enteras pequeñas.

bool representa verdadero y falso. Puede darle un nombre a su condición guardando los resultados de la comparación de esta manera:

<!-- wave-example: book-named-condition -->
```wave
fun main() {
    var balance: i32 = 5000;
    var cost: i32 = 3600;
    var can_buy: bool = balance >= cost;
    if (can_buy) {
        println("purchase allowed");
    } else {
        println("not enough money");
    }
}
```

Resultado de la ejecución:

```text
purchase allowed
```

`can_buy` no sigue automáticamente los cambios en balance y cost. Después de intercambiar los dos valores, la comparación se realiza nuevamente si se necesita el estado actual.

## espacio de almacenamiento no inicializado

`var value: i32;` es un formulario que declara solo el espacio de almacenamiento. Se debe escribir un valor válido antes de leer. No asuma que 0 se inserta automáticamente solo porque lo declara. En un curso introductorio, es más fácil de entender si el valor se conoce inmediatamente y se inicializa al mismo tiempo que la declaración.

En las llamadas a bibliotecas que reciben valores como argumentos de salida, hay casos en los que primero se declara un espacio y luego se lee cuando tiene éxito. En ese momento, debe verificar el resultado exitoso de la función. Evite el error de leer resultados no inicializados después de una llamada fallida.

## Rango válido de bloques y nombres.

Un bloque es una región de código encerrada entre llaves. Si declara una nueva variable con el mismo nombre dentro, la nueva variable se usará dentro de ese bloque. Esto se llama shadowing.

<!-- wave-example: book-shadowing -->
```wave
fun main() {
    var value: i32 = 10;

    if (true) {
        var value: i32 = value + 5;
        println("inner={}", value);
        value += 1;
        println("inner changed={}", value);
    }

    println("outer={}", value);
}
```

Resultado de la ejecución:

```text
inner=15
inner changed=16
outer=10
```

La expresión inicial `value + 5` en la declaración interna lee la exterior value. Una vez completada la inicialización de la nueva variable, el interior value es 15. Incluso si cambia el valor interno a 16, el espacio de almacenamiento externo no cambia. Después del bloque, verá value afuera nuevamente.

Por el contrario, si solo ejecuta `value += 1` sin `var` en el bloque interno, se cambiarán las variables visibles existentes. Mire la palabra clave para determinar si se trata de una nueva declaración o de un cambio en un valor existente.

## Variables locales y almacenamiento de nivel superior.

Fuera de la función, puedes usar const y static. const representa un valor constante y static es el espacio de almacenamiento mantenido durante la ejecución. No se pueden declarar en todas partes del mismo modo que las variables locales.

<!-- wave-example: book-storage -->
```wave
const LIMIT: i32 = 3;
static visits: i32 = 0;

fun visit() {
    visits += 1;
    println("visit={}", visits);
}

fun main() {
    visit();
    visit();
    println("limit={}", LIMIT);
}
```

Resultado de la ejecución:

```text
visit=1
visit=2
limit=3
```

Incluso si se llama a visit dos veces, static no se restablecerá a 0 en cada llamada. Por otro lado, si lo declaras como `var visits: i32 = 0;` dentro de una función, inicializará el almacenamiento local cada vez que lo llames. El estado mutable compartido puede dificultar el seguimiento del comportamiento, así que primero considere si se puede resolver con las entradas y salidas de la función.

## errores comunes

- Cuando se confunden declaración y cesión y se vuelve a declarar innecesariamente el mismo nombre.
- Si cree que la variable que almacena el resultado del cálculo sigue automáticamente los cambios en la variable de entrada.
- Si cree que, dado que el tipo es el mismo, las unidades como el número y el número de bytes también son las mismas.
- Cuando se lee sin inicialización o cuando se lee el argumento de salida cuando la función falla.
- Si cree que el nombre de una variable local es visible fuera del bloque.

Cuando busque un nombre por error, verifique la posición de declaración del nombre y el rango de llaves.

## Ejercicio: cálculo de cambios de inventario

El stock inicial es de 20 unidades y se vende dos veces, 3 unidades cada una. Imprima el inventario restante y el total de unidades vendidas. Cada vez que se cambia el inventario, se actualizan las mismas variables y el volumen de ventas también se acumula por separado.

### Solución completa

<!-- wave-example: book-stock-solution -->
```wave
fun main() {
    var stock: i32 = 20;
    var sold: i32 = 0;
    var order: i32 = 3;
    stock -= order;
    sold += order;
    stock -= order;
    sold += order;
    println("stock={} sold={}", stock, sold);
}
```

Resultado de la ejecución:

```text
stock=14 sold=6
```

El inventario y el volumen de ventas deben cambiar juntos. La actualización de cualquiera de ellos rompe la relación entre los valores. Agruparemos el procesamiento de pedidos repetidos mediante funciones de aprendizaje y declaraciones de bucle.


## Tipos de punto entero y flotante

Los tipos de números enteros son los siguientes:

- Firmado: `i8`, `i16`, `i32`, `i64`, `i128`, `i256`, `i512`, `i1024`
- Sin firmar: `u8`, `u16`, `u32`, `u64`, `u128`, `u256`, `u512`, `u1024`
- Entero del tamaño de la dirección: `isz`, `usz`
- Punto flotante: `f32`, `f64`

`isz` es un tipo entero con signo que coincide con el tamaño de la dirección y `usz` es un tipo entero sin signo que coincide con el tamaño de la dirección.

## Otros tipos incorporados

|tipo|uso|
| --- | --- |
| `bool` |`true` o `false`|
| `char` |Un valor de carácter de 8 bits sin signo. No arbitrario Unicode tipo de punto de código|
| `byte` |valor de bytes de 8 bits|
| `str` |Bytes de cadena que terminan en NUL|
| `ptr<T>` |Orientación por puntero `T`|
| `array<T, N>` |Matriz de longitud fija con tipo de elemento `T` y longitud `N`|

También se pueden utilizar estructuras, enumeraciones y alias de tipos definidos por el usuario en ubicaciones de tipos.

`var` es la sintaxis para declarar variables locales. Un alias de tipo es una gramática que expresa el mismo tipo con un nombre que se ajusta al contexto del código.
