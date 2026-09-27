---
translation_set_id: standard-library
path: reference/standard-library
locale: ru
group: stdlib
group_order: 1
order: 1
title: Стандартное руководство по библиотеке
summary: Как найти модуль, соответствующий вашим целям и прочитать ошибки и правила владения функцией.
---

## Найдите нужные вам функции

Стандартная библиотека — import с путем `std::module::file`. Даже если имена похожи, функции могут возвращать ошибки по-разному. Сначала прочтите [API Как читать](/docs/ru/stdlib/contracts), затем перейдите к нужному модулю в следующей таблице.

|Что я хочу сделать|документ|Главный import|
| --- | --- | --- |
|Длина строки/сравнение/поиск| [string](/docs/ru/reference/string-and-bytes) | `std::string::len`, `cmp`, `find`, `trim` |
|Распределение памяти/копирование/размер| [mem](/docs/ru/reference/memory-and-buffer) | `std::mem::alloc`, `ops`, `layout` |
|Список байтов разного размера| [buffer](/docs/ru/stdlib/buffer) | `std::buffer::alloc`, `read`, `write` |
|Двоичное чтение/запись| [bytes](/docs/ru/stdlib/bytes) | `std::bytes::types`, `cursor`, `leb128` |
|Дескриптор файла I/O| [fs и io](/docs/ru/stdlib/files-io) | `std::fs::file`, `std::io::fd` |
|Настройки комбинации пути/среды| [path и env](/docs/ru/stdlib/path-env) | `std::path::copy`, `std::env::environ` |
|Измерение времени/ожидание| [time](/docs/ru/stdlib/time) | `std::time::duration`, `clock`, `sleep` |
|Числовой поиск по списку адресов/имен| [net.resolve](/docs/ru/stdlib/resolution) | `std::net::resolve`, `resolve_table` |
|TCP Соединение/передача| [net.tcp](/docs/ru/stdlib/tcp) | `std::net::tcp`, `address`, `error` |
|OS Случайное число| [random](/docs/ru/stdlib/random) | `std::random::fill` |
|Процесс·OS Граница| [системная функция](/docs/ru/reference/system-io-network-process) | `std::process::core`, `spawn`, `std::sys` |
|Асинхронное выполнение задач| [task](/docs/ru/stdlib/task) | `std::task` |
|Ассистент по математике/диагностике| [math и debug](/docs/ru/stdlib/math-debug) | `std::math::int`, `float`, `std::debug::core` |

## Пример первого использования

В приведенной ниже программе используется одна функция из std без загрузки отдельного пакета. Сохраните его как `main.wave` и запустите как `wavec run main.wave`.

<!-- wave-example: library-import -->
```wave
import("std::string::len")::{
    len
};

fun main() {
    println("{}", len("Wave"));
}
```

Результат: `4`. Чтобы понять тот же пример шаг за шагом, прочитайте [Струны](/docs/ru/language/strings).

## Совместим с компилятором std.

Подтвердите выбранный путь с помощью `wavec print std-path`. При использовании std из другой кассы укажите путь как `wavec --std-root /absolute/path/to/std check main.wave`. Если указанный путь недействителен или несовместим, будет отображена ошибка.

## граница платформы

Различайте вычислительные функции, такие как строка/байт, и функции OS, такие как файл/сокет. Распознавание цели не гарантирует, что будут предоставлены все хосты API. Прочтите записи платформы для [Цель поддержки](/docs/ru/whale/build-link-targets) и каждого API вместе. `std::sys` — это интерфейс нижнего уровня, и переносимые программы сначала будут использовать модуль более высокого уровня.
