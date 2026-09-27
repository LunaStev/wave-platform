---
translation_set_id: overview
path: getting-started/overview
locale: es
group: getting-started
group_order: 1
order: 1
title: Documentación de Wave y guía de aprendizaje.
summary: Aprenda Wave paso a paso, desde la instalación hasta los programas prácticos, y busque reglas de lenguaje y API de biblioteca estándar.
---

## Aprende Wave con esta guía

Aprenda a escribir el código fuente Wave, compílelo, ejecútelo y compruebe los resultados. Si es nuevo en la programación, siga la secuencia a continuación. Si conoce otro idioma, ejecute los ejemplos de cada capítulo y compare sus reglas y casos límite con lo que ya sabe.

## Camino de aprendizaje

|paso|Capítulo|lo que aprenderás|
| --- | --- | --- |
|Configuración| [Instalación](/docs/es/getting-started/install) |Prepare el compilador y la biblioteca estándar y verifique que se ejecuten.|
| 1 | [Tu primer programa](/docs/es/language/program-structure) |Cree, verifique y ejecute un archivo fuente y comprenda los códigos de salida|
| 2 | [Variables y tipos](/docs/es/language/declarations-and-types) |Almacenar valores y elegir tipos con el rango requerido|
| 3 | [Operadores y conversiones](/docs/es/language/expressions-and-operators) |Explicar el orden de evaluación y los resultados de las conversiones de tipos.|
| 4 | [Condiciones y bucles](/docs/es/language/control-flow) |Bifurcación en condiciones y datos de proceso con bucles|
| 5 | [Funciones](/docs/es/language/functions-and-generics) |Extraer operaciones repetidas en funciones|
| 6 | [matrices](/docs/es/language/arrays) |Acceda a elementos por índice e itere sobre una matriz|
| 7 | [cuerdas](/docs/es/language/strings) |Distinguir caracteres de bytes y comprender los escapes y la longitud de las cadenas.|
| 8 | [Estructuras y variantes](/docs/es/language/structures-enums-and-aliases) |Agrupa datos relacionados y representa el éxito y el fracaso.|
| 9 | [Punteros y vidas](/docs/es/language/explicit-memory-type-model) |Modifique el valor original a través de su dirección y administre su vida útil.|
| 10 | [memoria dinámica](/docs/es/language/allocation) |Manejar fallas de asignación y liberar memoria|
| 11 | [Módulos y genéricos.](/docs/es/language/modules-imports-and-ffi) |Divida el código entre archivos y reutilice funciones con diferentes tipos|
| 12 | [Manejo de errores](/docs/es/language/errors) |Verifique los resultados y limpie los recursos en caso de falla|
| 13 | [Introducción al código asincrónico](/docs/es/language/async-and-never) |Crea un futuro y espera a que termine.|

## Pon en práctica tus conocimientos

Después de los capítulos principales, cree una [calculadora de entrada](/docs/es/practice/input-calculator), un [lector de archivos](/docs/es/practice/file-reader), un [mensaje binario](/docs/es/practice/binary-message) y un [cliente TCP](/docs/es/practice/tcp-client). Pruebe tanto los casos de entrada exitosos como los de fracaso en cada proyecto.

## Tres pestañas de documentación

- **Wave**: Un curso de idiomas guiado y proyectos prácticos a seguir en orden.
- **[Biblioteca estándar](/docs/es/stdlib)**: API, valores de retorno, errores, reglas de propiedad y requisitos de plataforma de cada módulo.
- **[Whale](/docs/es/whale)**: Creación y vinculación, administración de paquetes, uso de comandos y cadena de herramientas de bajo nivel.

Los ejemplos distinguen programas completos de fragmentos que pertenecen dentro de una función. Ejecute los comandos `wavec` en una terminal y guarde los bloques de código `wave` en archivos `.wave`. La entrada y la salida se muestran por separado; Los ejemplos que leen la entrada estándar especifican qué ingresar.

## Cuando te quedas atascado

Utilice [Solución de problemas](/docs/es/reference/diagnostics) para distinguir los problemas de instalación, verificación de código fuente, vinculación y ejecución. Busque reglas de lenguaje en la [referencia rápida de sintaxis](/docs/es/reference/syntax-quick-reference), comandos en la [referencia del compilador](/docs/es/getting-started/compiler) y API en la [guía de biblioteca estándar](/docs/es/reference/standard-library).
