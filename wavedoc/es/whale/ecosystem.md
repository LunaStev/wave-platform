---
translation_set_id: ecosystem
path: whale/ecosystem
locale: es
group: whale
group_order: 1
order: 2
title: Componentes de la cadena de herramientas
summary: Describe la función de una cadena de herramientas de bajo nivel separada, Whale, y los límites de los componentes del ecosistema Wave.
---

## WhaleIrán

Whale es una cadena de herramientas de bajo nivel que maneja ensamblajes y representaciones intermedias. Los componentes que manejan ensamblajes, objetos, enlaces y representaciones intermedias están diseñados para ser reutilizables en Wave y otras herramientas de generación de código nativo.

Whale no es el nombre de todo el entorno de desarrollo Wave. Las responsabilidades de cada proyecto se dividen de la siguiente manera:

|proyecto|responsabilidad|
| --- | --- |
| `wavec` |Wave Examina la fuente y crea un archivo ejecutable.|
| Vex |Gestiona el paquete Wave, manifest, el gráfico de dependencia, lockfile y las compilaciones de paquetes.|
| Whale |Proporciona componentes independientes assembler, object, linker y IR.|
| Wave `std` |El tiempo de ejecución y el sistema API se proporcionan como módulos fuente Wave.|

## componente

Whale workspace consta de cuatro áreas de biblioteca principales:

- `assembler`: Tokenización, AMD64 Análisis/Codificación, section, symbol y relocation
- `object`: modelo de archivo de objeto y ELF64 writer
- `linker`: Capa de enlace
- `ir`: Whale IR tipo, builder, salida, verificación y opcional frontend socket

El ejecutable `whale` proporciona a esta región los comandos `asm`, `object`, `link` y `ir`.

## borde de herramienta

El programa Wave se construye como `wavec`. Cuando trabaje directamente con ensamblados, archivos de objetos y IR, utilice el comando `whale`.

La instalación de Whale no cambia el método de compilación de `wavec`. Vex usa `wavec` para compilar el paquete Wave y Whale lo ejecuta directamente en una tarea que maneja artefactos de bajo nivel.

## Verificación entregable

Al vincular el artefacto Whale a su proceso de compilación, asegúrese de que object format y el objetivo architecture coincidan. symbol y relocation se pueden verificar con herramientas independientes como `readelf` y `objdump`. Las compilaciones que usan IR socket deben usar socket schema del mismo productor que Whale.
