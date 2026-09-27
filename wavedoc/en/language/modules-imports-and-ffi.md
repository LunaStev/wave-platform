---
translation_set_id: modules-ffi
path: language/modules-imports-and-ffi
locale: en
group: language
group_order: 2
order: 11
title: 11. Modules and generic code
summary: Learn about getting public names and explicit type arguments.
---

## Reasons for dividing files

As the program grows, it is easier to find it by grouping related functions together rather than placing all functions in main.wave. Module boundaries determine which names are exposed to other code. Generics are a tool that reuses the same task with different types, independent of file separation.

In this chapter, we create a two-file program and learn module aliases, selection import, and generic functions and structures.

## two file program

Create helpers.wave and main.wave in the same directory.

helpers.wave All:

```wave
pub fun double(value: i32) -> i32 {
    return value * 2;
}
```

main.wave All:

<!-- wave-example: book-module-two-files -->
```wave
import("./helpers")::{
    double
};

fun main() {
    var result: i32 = double(21);

    println("{}", result);
}
```

Execution result:

```text
42
```

Run `wavec run main.wave` in terminal. helpers.wave is also not run separately. The required source is connected via import.

The pub in front of the function in helpers indicates that it can be imported by other modules. Auxiliary functions that do not need to be exposed to the outside world do not need to be made public. Even if you change the internal implementation of a module, you can reduce changes to the code you use by keeping the contract of the public function.

## Baseline for relative paths

`./helpers` is relative to the directory of the source file that created the import sentence. When running the program, separate it from the working directory where the file I/O is written. The steps for finding the file import and the steps for finding input.txt while running are different.

For local import, the `.wave` extension can be omitted. If you divided the directory, write the path `./module` according to the location. Do not confuse local relative paths with paths that fetch the names of package dependencies.

## Select import and alias

The selection import causes only the desired public names to be used directly in the current file. If there is a name conflict or you want to reveal which module the function belongs to, use an alias.

<!-- wave-example: book-module-alias -->
```wave
import("std::string::len" as strings);

fun main() {
    var length: i32 = strings::len("Wave");

    println("{}", length);
}
```

Execution result:

```text
4
```

strings is the module alias defined in this file. `strings::len` uses the name of the module. The dot for field access and `::` for module separation are different notations.

Do not use the option import and the alias import together in one sentence. Whatever your style, use it consistently throughout the file so that the origin of the name is easily readable.

## Standard libraries and packages

The path `std::` points to the standard library. The user publishes API and import the required modules. Not all standard library functions are automatically placed into the current namespace.

External package paths start from the package name. The location of the package is provided by compiler options or the package manager. First learn the boundaries with local modules and then learn how to manage dependencies in [Vex How to use](/docs/en/whale/vex-package-manager).

## Use the same function for each type

The following functions return their input verbatim: For i32 and str, use the type parameter T to avoid writing the same code twice.

<!-- wave-example: book-generic-identity -->
```wave
fun identity<T>(value: T) -> T {
    return value;
}

fun main() {
    var number: i32 = identity<i32>(42);
    var text: str = identity<str>("Wave");

    println("{} {}", number, text);
}
```

Execution result:

```text
42 Wave
```

T is a place to enter the type, not the integer value passed during execution. Call it by specifying the type argument like `<i32>`. Normal user generic functions do not omit type arguments.

identity<str> does not duplicate string bytes by allocating them again. Returns the value as is. The grammar of generics does not change the copy and ownership rules of data.

## Operations required by the generic body

Just because there is a type parameter does not mean that all operations can be used on all types. minimum below should be used as an actual type that can be compared.

<!-- wave-example: book-generic-minimum -->
```wave
fun minimum<T>(left: T, right: T) -> T {
    if (left < right) {
        return left;
    }

    return right;
}

fun main() {
    var small: i32 = minimum<i32>(7, 4);
    var wide: i64 = minimum<i64>(100, 20);

    println("{} {}", small, wide);
}
```

Execution result:

```text
4 20
```

If you change the type argument, `<` and the return used in the text must be of the corresponding type. When reading generic errors, check both the type combination called and the operation required by the function body.

## generic struct

Let's create Pair, which binds two different values.

<!-- wave-example: book-generic-pair -->
```wave
struct Pair<A, B> {
    first: A;
    second: B;
}

fun main() {
    var item: Pair<i32, str> = Pair<i32, str> {
        first: 7,
        second: "seven"
    };

    println("{} {}", item.first, item.second);
}
```

Execution result:

```text
7 seven
```

Pair<i32, str> and Pair<i64, str> are different specific types. The order of type arguments is also meaningful. Read the declaration and generation code to see where the types of first and second are determined.

## Name and contract of API released

When you publish a function, you specify not only the name but also the input unit, return value, failure, and ownership. For example, the loop the caller will write depends on whether read reads the maximum length or the exact length.

pub is the public scope between modules Wave. This is different from export (c), which exports external symbols for other languages ​​to call. You can see a complete example linking the two languages ​​at [See FFI](/docs/en/language/modules-imports-and-ffi).

## Exercise and complete solution

Create a public function square on math.wave and call it as an alias on main.wave to print the power of 3 and 5.

math.wave:

```wave
pub fun square(value: i32) -> i32 {
    return value * value;
}
```

main.wave:

<!-- wave-example: book-module-solution -->
```wave
import("./math" as math);

fun main() {
    println("{} {}", math::square(3), math::square(5));
}
```

Execution result:

```text
9 25
```

If the error is that the function has disappeared, check the import path and pub first. If there is a name conflict, check whether the call has an alias. If it is a type error, check the input of the function and the type of the passed argument. Don't try to solve different problems with one path correction.


## C Import function

```wave
extern(c) fun puts(text: ptr<i8>) -> i32;
```

The name ABI can be followed by the actual symbol name as a string.

```wave
extern(c, "native_symbol") fun local_name(value: i32) -> i32;
```

## Wave Export function

```wave
export(c) fun wave_add(left: i32, right: i32) -> i32 {
    return left + right;
}
```

`extern` and `export` can be used as single functions and blocks. The exported function must have the specific signature ABI and therefore cannot be generic.

## Target condition property

Target condition properties can be attached to top-level items.

```wave
#[target(os="linux", arch="x86_64")]
extern(c) fun platform_call(value: i32) -> i32;
```

The condition keys are `arch`, `os`, `env`, `abi` and the properties apply to the next top-level item.

## Connect with the C function you wrote yourself

This lab is for the native environment with the C compiler. Concatenates an integer function without library allocation or string processing.

`native.c`:

```c
#include <stdint.h>
int32_t native_double(int32_t value) { return value * 2; }
```

`main.wave`:

```wave
extern(c) fun native_double(value: i32) -> i32;

fun main() {
    println("{}", native_double(21));
}
```

Run Linux/macOS from a terminal in the same working directory.

```shell
cc -c native.c -o native.o
wavec build main.wave native.o -o ffi-example
./ffi-example
```

The expected output is `42`. In the developer shell of MSVC in Windows, create object with `cl /c native.c /Fonative.obj` and connect to `wavec build main.wave native.obj -o ffi-example.exe`. The source and target architectures of object must be the same. The examples use only small values. To pass a large value to the C function, the multiplication range of the C page must also be guaranteed separately.

Local file paths start with `./`.
