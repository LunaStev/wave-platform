---
translation_set_id: stdlib-contracts
path: stdlib/contracts
locale: es
group: stdlib
group_order: 1
order: 2
title: Lectura de documentación API: errores y propiedad
summary: Comprenda las unidades de argumentos, las estructuras de resultados, el éxito parcial y la vida útil de los recursos.
---

## Leer la declaración

La siguiente notación describe la declaración de función y no es el archivo ejecutable completo.

```text
io_read(fd: i64, buf: ptr<u8>, len: i64) -> i64
```

`fd` es el descriptor abierto, `buf` es el almacenamiento proporcionado por la persona que llama y `len` es el número de bytes disponibles para escribir. `i64` no significa que las longitudes negativas sean válidas. El valor de retorno es el número real de bytes leídos, no la longitud de la solicitud, por lo que solo se utiliza el rango devuelto.

## La expresión de falla varía de una función a otra.

|camino|si|Método de inspección|
| --- | --- | --- |
|Puntero o null| `mem_alloc` |null Acceso a memoria después de la inspección|
|número de bytes o número negativo| `io_read` |Error negativo, 0 EOF, datos positivos|
|código de estado| `buffer_push` |Error de comparación constante con `BUFFER_OK`|
|Éxito y valor| `NetResult<T>` |Después de examinar `ok`, use `value`|
|Progreso parcial incluido.| `RandomFillResult` |Marque `ok`, `written`, `error` juntos.|

Solo analiza la cantidad de errores y no los compara con constantes de otros módulos. Por ejemplo, los números de error env y OS errno no son el mismo sistema. El error original en WASI no debe interpretarse como Linux errno.

## poseer y alquilar

- **Propiedad**: cuando se adquiere una memoria asignada, un archivo abierto o un socket abierto, es responsable de llamar a la liberación/cierre correspondiente.
- **Préstamo**: El byte view o el búfer pasado a la función se refiere a la memoria existente. Si una función no especifica que recibe propiedad, la persona que llama conserva el control.
- **Argumento de salida**: pasa un espacio de almacenamiento válido donde el resultado se puede escribir en una función que lo recibe, como `out_value: ptr<T>`. Asegúrese de que el contrato establezca que el resultado solo es válido si tiene éxito.

Copiar una estructura Buffer puede dejar ambas copias apuntando a la misma asignación. No libere cada copia por separado. Un puntero prestado deja de ser válido después de que la asignación se libera o reasigna. Los literales de cadena no son buffers grabables.

## Fracasar no significa volver a un estado anterior

`io_write_all` puede fallar después de escribir algunos bytes. Los bytes ya escritos externamente no se devolverán. Por otro lado, la lectura checked cursor de bytes conserva la posición y el valor de salida si falla. Estas diferencias se especifican en API.

Si también desea practicar el manejo de fallas, continúe con [Lector de archivos](/docs/es/practice/file-reader) y [mensaje binario](/docs/es/practice/binary-message).
