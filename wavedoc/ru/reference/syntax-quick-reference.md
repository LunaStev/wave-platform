---
translation_set_id: quick-reference
path: reference/syntax-quick-reference
locale: ru
group: reference
group_order: 5
order: 3
title: Краткий справочник по синтаксису
summary: Часто используемые объявления, поток управления, типы, указатели и грамматика FFI организованы на одной странице.
---

## декларация

```wave
var value: i32 = 1;
var next: i32 = value + 1;
const LIMIT: i32 = 64;
static total: i64 = 0;
type Identifier = u64;
```

`var` — регион, `const`/`static` — объявления верхнего уровня. Локальные переменные явно объявляют свой тип.

## Функции

```wave
fun max(left: i32, right: i32) -> i32 {
    if (left > right) {
        return left;
    }

    return right;
}
```

## универсальный

```wave
fun identity<T>(value: T) -> T {
    return value;
}

var value: i32 = identity<i32>(10);
```

При вызове универсального метода укажите аргумент типа.

## Структура и enum

```wave
struct Pair {
    left: i32;
    right: i32;
}

enum Result -> i32 {
    Ok = 0,
    Error,
}
```

## Условия и циклы

```wave
if (ready) {
    println("ready");
}

while (count < 10) {
    count += 1;
}

for (var i: i32 = 0; i < 10; i += 1) {
    println("{}", i);
}

match (status) {
    Ready => {
        println("ready");
    }
    0 => {
        println("zero");
    }
    _ => {
        println("other");
    }
}
```

В заголовках `if`, `while`, `for` и `match` используются круглые скобки.

## Массивы и указатели

```wave
var values: array<i32, 4> = [1, 2, 3, 4];
var p: ptr<i32> = &values[0];
var first: i32 = deref p;
```

## консольный ввод/вывод

```wave
print("value = ");
println("{}", value);
input("{}", value);
```

Первый аргумент — строковый литерал. За каждым точным заполнителем `{}` должно следовать выражение, а цель `input` должна быть назначаемой.

## import и общедоступные предметы

```wave
import("std::string::len");
import("./helpers" as helpers);
import("math")::{
    Vector
};

pub fun add(left: i32, right: i32) -> i32 {
    return left + right;
}

pub import("./extra"):: {
    increment
};
```

Локальный путь начинается с `./`. Псевдоним import определяет имя модуля, а выбор import импортирует необходимые общедоступные записи в пространство имен этого файла. `pub import` повторно экспортирует выбранные элементы.

## FFI

```wave
extern(c) fun native_call(value: i32) -> i32;

export(c) fun wave_call(value: i32) -> i32 {
    return value + 1;
}
```

## Целевой условный элемент

```wave
#[target(os="linux", arch="riscv64")]
extern(c) fun platform_call(value: i32) -> i32;
```

Ключи условий поддержки: `arch`, `os`, `env`, `abi`. Свойства управляют следующим элементом верхнего уровня.

## Линейная сборка

```wave
var result: i64 = 0;
asm {
    "mv a0, a1"
    in("a1") 7
    out("a0") result
    clobber("memory")
}
```

Текст инструкций и имена регистров зависят от цели. Объявите все входы, выходы и скрытые clobber, необходимые для блока.

## проверка источника

```shell
wavec build main.wave --emit=check
wavec print supported-targets
wavec print supported-emit-kinds
```

## Обучение и примеры

Примерами локальных переменных и операторов, отдельно отображаемых вне функции, являются фрагменты кода, вставленные в тело функции. Полные примеры бега и упражнения приведены в [Wave Процесс обучения](/docs/ru/getting-started/overview). Пожалуйста, проверьте [Стандартная библиотека](/docs/ru/stdlib) для получения подробных правил использования памяти и внешних функций.
