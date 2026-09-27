---
translation_set_id: learn-allocation
path: language/allocation
locale: es
group: language
group_order: 2
order: 10
title: 10. Asignación de memoria y gestión de recursos.
summary: Obtenga información sobre errores de asignación, inicialización, alcance y liberación.
---

## ¿Cuándo necesitas espacio de almacenamiento dinámico?

Una matriz de tamaño fijo incluye su longitud en su tipo. Utilice la memoria dinámica cuando la cantidad de datos se conozca sólo en tiempo de ejecución, como el tamaño de un archivo o la longitud de una entrada. Libera cada asignación cuando ya no la necesites.

En este capítulo, administrará una pequeña asignación, cambiará su tamaño y luego usará un Buffer. Pasar un puntero es diferente a transferir la propiedad. Identifique qué recursos posee cada función a medida que sigue los ejemplos.

## Asignar, comprobar, utilizar y liberar

<!-- wave-example: book-alloc-lifecycle -->
```wave
import("std::mem::alloc")::{
    mem_alloc_zeroed, mem_free
};

fun main() -> i32 {
    var size: i64 = 4;
    var data: ptr<u8> = mem_alloc_zeroed(size);

    if (data == null) {
        println("allocation failed");
        return 1;
    }

    deref data[0] = 42;
    println("{} {}", deref data[0], deref data[1]);

    var status: i64 = mem_free(data, size);

    if (status < 0) {
        return 2;
    }

    return 0;
}
```

Resultado de la ejecución:

```text
42 0
```

El programa tiene cuatro pasos: solicitar 4 bytes, verificar null, acceder solo al rango válido y liberar la asignación. En caso de éxito, mem_alloc_zeroed inicializa la memoria a cero, por lo que el segundo byte es cero aunque el programa no haya escrito en él.

No asuma ningún contenido inicial para la memoria devuelta por mem_alloc. Inicialice cada región antes de leerla. Un tamaño de asignación cero o negativo devuelve null. Una asignación de tamaño positivo también puede no obtener memoria.

## unidad de tamaño

El argumento de tamaño de una API de asignación de memoria se mide en bytes. Para asignar diez números enteros, multiplique el tamaño del elemento por el número de elementos. Comprueba que esta multiplicación no se desborde.

<!-- wave-example: book-alloc-typed -->
```wave
import("std::mem::alloc")::{
    mem_alloc, mem_free
};
import("std::mem::layout")::{
    size_of
};
import("std::mem::ops")::{
    mem_size_mul_checked
};

fun main() -> i32 {
    var count: i64 = 3;
    var bytes: i64 = 0;
    var item_size: i64 = size_of<i32>() as i64;

    if (mem_size_mul_checked(count, item_size, &bytes) < 0) {
        return 1;
    }

    var raw: ptr<u8> = mem_alloc(bytes);

    if (raw == null) {
        return 2;
    }

    var values: ptr<i32> = raw as ptr<i32>;

    for (var index: i64 = 0; index < count; index += 1) {
        deref values[index] = (index + 1) as i32;
    }

    println("{} {} {}", deref values[0], deref values[1], deref values[2]);

    if (mem_free(raw, bytes) < 0) {
        return 3;
    }

    return 0;
}
```

Resultado de la ejecución:

```text
1 2 3
```

count es un recuento de elementos; bytes es un recuento de bytes. La aritmética de punteros se mueve en unidades de i32, pero para liberar la asignación se requiere su tamaño original en bytes. size_of utiliza el diseño del tipo de destino, haciendo explícita la relación con el tipo de elemento.

Al cambiar el resultado de size_of a i64 para un tipo general muy grande, también se debe considerar el rango de conversión. Aquí usamos i32, que tiene un tamaño conocido.

## Limpie incluso en rutas de falla

Si otra operación falla después de la asignación, libere la memoria antes de regresar antes de tiempo. Una tabla de propiedad ayuda a identificar rutas de limpieza que de otro modo podría pasar por alto.

|paso|Recursos propiedad de|si fallas|
| --- | --- | --- |
|Antes de la asignación|Ninguno|devuélvelo de inmediato|
|Después de una asignación exitosa|data y tamaño original|data Regreso después del lanzamiento|
|Después de una reasignación exitosa|Nueva dirección y nuevo tamaño.|Lanzamiento de nueva dirección|
|Después del lanzamiento|Ninguno|No utilices la dirección anterior.|

Al sobrescribir una variable de puntero y perder la dirección original también se pierde la información necesaria para liberar la asignación. Esto provoca una pérdida de memoria. Por el contrario, liberar la misma asignación a través de dos propietarios provoca una doble liberación.

## Reasignación para aumentar el tamaño

<!-- wave-example: book-alloc-grow -->
```wave
import("std::mem::alloc")::{
    mem_alloc, mem_realloc, mem_free
};

fun main() -> i32 {
    var data: ptr<u8> = mem_alloc(4);

    if (data == null) {
        return 1;
    }

    deref data[0] = 7;
    var next: ptr<u8> = mem_realloc(data, 4, 8);

    if (next == null) {
        mem_free(data, 4);
        return 2;
    }

    data = next;
    deref data[4] = 9;
    println("{} {}", deref data[0], deref data[4]);

    if (mem_free(data, 8) < 0) {
        return 3;
    }

    return 0;
}
```

Resultado de la ejecución:

```text
7 9
```

Primero almacene el resultado en el siguiente. Si falla la asignación del bloque más grande, los datos originales siguen siendo válidos y aún se pueden liberar. Si tiene éxito, la asignación anterior se libera y se debe utilizar la nueva dirección. Inicialice la región recién agregada antes de leerla.

El nuevo tamaño en este ejemplo es positivo. En cambio, una solicitud con new_size=0 intenta liberar la asignación existente y devuelve null. Por lo tanto, un resultado null no siempre significa que la asignación anterior sigue siendo válida. Llame a mem_free directamente cuando necesite verificar si la liberación se realizó correctamente.

## Cuándo volver a comprobar los indicadores prestados

data Es incorrecto guardar un puntero interno y usarlo después de la reasignación. Esto se debe a que la dirección del nuevo data puede ser diferente. Si necesita una ubicación interna, puede almacenar offset en lugar de la dirección y volver a calcularla según el nuevo data después del éxito.

Liberar o reasignar memoria también afecta al código que la ha tomado prestada. Compruebe si otra operación todavía está utilizando esa memoria. Un búfer pasado a una operación asincrónica debe seguir siendo válido hasta que finalice la operación.

## La lista de bytes contiene Buffer

Administrar una lista de bytes con longitudes que cambian con frecuencia mientras se reasigna manualmente requiere manejar tanto len como cap, fallas de expansión y cálculos de tamaño. Buffer de std agrupa estas operaciones.

<!-- wave-example: book-buffer-progressive -->
```wave
import("std::buffer::alloc")::{
    Buffer, buffer_init, buffer_free
};
import("std::buffer::write")::{
    buffer_append_str
};

fun main() -> i32 {
    var message: Buffer;

    if (buffer_init(&message, 0) < 0) {
        return 1;
    }

    if (buffer_append_str(&message, "Hello") < 0) {
        buffer_free(&message);
        return 2;
    }

    if (buffer_append_str(&message, ", Wave") < 0) {
        buffer_free(&message);
        return 3;
    }

    println("bytes={}", message.len);

    if (buffer_free(&message) < 0) {
        return 4;
    }

    return 0;
}
```

Resultado de la ejecución:

```text
bytes=11
```

len es el número de bytes en uso; cap es la capacidad asignada. Agregar datos aumenta la asignación cuando es necesario. El uso de Buffer no elimina la responsabilidad de la persona que llama de liberarlo.

buffer_append_str no agrega NUL al final de la cadena. Por lo tanto, message.data no debe generarse directamente como str. Los bytes se generan con la función I/O, que toma la longitud o construye explícitamente una representación de cadena.

## Ejercicio y enfoque

Agregue los bytes del 0 al 9 uno por uno a Buffer y obtenga la suma. Debe liberarse cuando falla cada adición y las lecturas solo se realizan en el rango len. La solución completa y la falla de límites se pueden verificar siguiendo el ejemplo en [Buffer Cómo utilizar](/docs/es/stdlib/buffer).

Intente marcar llamadas de asignación, reasignación y desasignación en su código. Para cada asignación exitosa, debe poder describir quién es el propietario y qué ruta la libera.
