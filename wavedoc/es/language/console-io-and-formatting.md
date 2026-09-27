---
translation_set_id: console-io-formatting
path: language/console-io-and-formatting
locale: es
group: language
group_order: 2
order: 16
title: Entrada, salida y formato de la consola
summary: print, println, input Describe oraciones y reglas de marcador de posición.
---

## declaración de entrada/salida

Wave proporciona `print`, `println` y `input` como declaraciones de entrada/salida de la consola.

```wave
fun main() {
    var count: i32 = 0;
    input("{}", count);
    print("count = ");
    println("{}", count);
}
```

Cada oración termina con `;`. El primer argumento debe ser una cadena literal. Las variables o cadenas calculadas no se pueden utilizar como argumentos format.

## marcador de posición

Sólo exactamente dos caracteres, `{}`, son marcadores de posición.

```wave
println("name = {}, score = {}", name, score);
```

El número de marcadores de posición y el número de expresiones que siguen deben ser exactamente iguales. Si los números son diferentes, es un error gramatical.

```wave
println("{} {}", one);
// 오류: 자리표시자 2개, 값 1개
println("plain text", one);
// 오류: 자리표시자 없음, 값이 남음
```

Otras formas de llaves se dejan como texto sin formato. No hay marcadores de posición numerados o con nombre en esta gramática.

## print y println

`print` imprime el texto formateado tal cual y `println` agrega un salto de línea.

```wave
print("loading...");
println("done");
```

Los argumentos de formato utilizan valores escalares como números enteros, números de punto flotante, cadenas y punteros. Las matrices y estructuras no se pueden utilizar como argumentos de formato.

## input Objetivo

`input` almacena el valor leído en el destino, por lo que todas las expresiones después de format deben ser ubicaciones en las que se pueda escribir.

```wave
var number: i32 = 0;
input("{}", number);
```

Como objetivos se pueden utilizar variables, campos y ubicaciones de almacenamiento desreferenciadas. Los literales y los resultados de los cálculos no se pueden utilizar como entrada.

Si todos los valores de entrada no se pueden convertir al tipo solicitado, el programa sale con un estado de error.

## límite de tiempo de ejecución

Estas declaraciones utilizan entradas y salidas de la consola del entorno hosted. En un entorno independiente, la entrada/salida proporcionada por el kernel o dispositivo debe definirse como una función o límite FFI.

## Valor de entrada y rango

La entrada bool solo acepta `0` y `1`. No interpreta 2 como true ni acepta la cadena `true` como la misma entrada. La entrada de enteros debe estar dentro del rango del ancho del entero objetivo. También se procesan enteros de 128, 256, 512 y 1024 bits en función del ancho total del tipo.

Error de formato, fuera de rango, antes de que falle la entrada requerida EOF. El input integrado no es una función que devuelve un error y vuelve a ingresar, sino que es una función de entrada que finaliza el proceso en caso de error. Si necesita procesamiento de entrada recuperable, lea los bytes con io y construya un analizador independiente.

[Práctica de calculadora de entrada](/docs/es/practice/input-calculator) · [Archivo y io](/docs/es/stdlib/files-io)
