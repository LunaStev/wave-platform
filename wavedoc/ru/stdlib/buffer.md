---
translation_set_id: stdlib-buffer
path: stdlib/buffer
locale: ru
group: stdlib
group_order: 1
order: 5
title: buffer: расширяемое хранилище байтов
summary: Buffer Описывает правила инициализации, добавления, запроса, емкости и освобождения.
---

## Значение Buffer

`Buffer` в `std::buffer::types` имеет `data: ptr<u8>`, `len: i64` и `cap: i64`. len — количество инициализированных и используемых байтов, а cap — общее количество выделенных байтов. Всегда сохраняйте `0 <= len <= cap`. Строка NUL не гарантирует автоматического завершения.

## Базовый API

|модуль|декларация|смысл|
| --- | --- | --- |
| `std::buffer::alloc` | `buffer_init(out_buf: ptr<Buffer>, capacity: i64) -> i64` |Инициализируйте новый репозиторий. Не вызывать в уже принадлежащие буферы|
|тот же модуль| `buffer_free(buf: ptr<Buffer>) -> i64` |Освободить. Пусто в случае успеха|
|тот же модуль| `buffer_reserve(buf: ptr<Buffer>, required_cap: i64) -> i64` |Обеспечьте минимальную полную мощность. len Поддерживается|
|тот же модуль| `buffer_clear(buf: ptr<Buffer>) -> i64` |Поддерживать емкость и len=0|
|тот же модуль| `buffer_resize(buf: ptr<Buffer>, new_len: i64, value: u8) -> i64` |Измените длину, заполните новые байты с помощью value|
| `std::buffer::write` | `buffer_push(buf: ptr<Buffer>, value: u8) -> i64` |добавить один байт|
|тот же модуль| `buffer_append(buf: ptr<Buffer>, src: ptr<u8>, size: i64) -> i64` |sizeСкопировать и добавить байт|
|тот же модуль| `buffer_append_str(buf: ptr<Buffer>, s: str) -> i64` |Добавьте строковые байты, исключая NUL.|
| `std::buffer::read` | `buffer_get(buf: ptr<Buffer>, index: i64, out_value: ptr<u8>) -> i64` |Чтение одного байта в диапазоне|

Возврат статуса API возвращает `BUFFER_OK`(0) в случае успеха. Отличайте ошибки INVALID, BOUNDS, OVERFLOW и ALLOC от `std::buffer::error`. Число не интерпретируется как OS errno. `buffer_new` представляет ошибку выделения в виде пустого Buffer, поэтому, когда вам нужно различать ошибки, используйте `buffer_init`.

## Запуск примера

<!-- wave-example: buffer-api -->
```wave
import("std::buffer::alloc")::{
    Buffer, buffer_init, buffer_free
};
import("std::buffer::write")::{
    buffer_append_str, buffer_push
};
import("std::buffer::read")::{
    buffer_get
};

fun main() -> i32 {
    var data: Buffer;
    if (buffer_init(&data, 0) < 0) {
        return 1;
    }
    if (buffer_append_str(&data, "Hi") < 0 || buffer_push(&data, 33) < 0) {
        buffer_free(&data);
        return 2;
    }

    var value: u8 = 0;
    if (buffer_get(&data, 2, &value) < 0) {
        buffer_free(&data);
        return 3;
    }

    println("{} {}", data.len, value);
    if (buffer_free(&data) < 0) {
        return 4;
    }

    return 0;
}
```

Результат выполнения:

```text
3 33
```

Сохраните его как `main.wave` и запустите как `wavec run main.wave`. Начальная емкость, равная 0, — это не сбой, а действительный пустой буфер. Освобождает место при дальнейшей обработке.

## Продолжительность жизни и неудача

Увеличение буфера может изменить данные. Не используйте ранее заимствованный адрес после операции, которая может привести к перераспределению. Копирование структуры Buffer не дублирует ее выделение, поэтому назначьте этому выделению одного владельца.

`buffer_get` не меняет выходные аргументы в случае сбоя. С другой стороны, удобная функция `buffer_at` также представляет ошибки как 0, поэтому используйте `buffer_get`, чтобы различать фактические 0 байтов и ошибки. Избегайте создания недействительного len/cap путем непосредственного изменения общедоступных полей.

[Память API](/docs/ru/reference/memory-and-buffer) · [Потренируйтесь читать файл как Buffer.](/docs/ru/practice/file-reader)

## Учитывайте длину и мощность отдельно.

reserve освобождает место для хранения, но не увеличивает len. resize изменяет фактическую используемую длину и инициализирует расширенную часть указанным байтом. clear устанавливает только используемую длину в 0, что позволяет повторно использовать выделение.

Сохраните следующую программу как main.wave и запустите ее. Он не зависит от точного коэффициента роста мощности; это просто гарантирует, что у вас есть необходимое пространство.

<!-- wave-example: book-buffer-length-capacity -->
```wave
import("std::buffer::alloc")::{
    Buffer,
    buffer_init,
    buffer_reserve,
    buffer_resize,
    buffer_clear,
    buffer_free
};

fun main() -> i32 {
    var data: Buffer;

    if (buffer_init(&data, 0) < 0) {
        return 1;
    }

    if (buffer_reserve(&data, 20) < 0) {
        buffer_free(&data);
        return 2;
    }

    println("reserved length={}", data.len);

    if (buffer_resize(&data, 3, 7) < 0) {
        buffer_free(&data);
        return 3;
    }

    println("resized length={} first={}", data.len, deref data.data[0]);

    if (buffer_clear(&data) < 0) {
        buffer_free(&data);
        return 4;
    }

    println("cleared length={}", data.len);

    if (buffer_free(&data) < 0) {
        return 5;
    }

    return 0;
}
```

Результат выполнения:

```text
reserved length=0
resized length=3 first=7
cleared length=0
```

При увеличении до resize мы передали value=7, поэтому все три новых байта, которые мы видим, равны 7. Пространство, защищенное только reserve, не считывается как инициализированные данные. cap останется после clear и может быть добавлен снова к тому же Buffer.

## Различайте нулевые байты и ошибки поиска.

buffer_get возвращает статус и записывает фактические байты в качестве выходных аргументов. Даже если данные равны 0, это нормальный успех.

<!-- wave-example: book-buffer-zero-bounds -->
```wave
import("std::buffer::alloc")::{Buffer, buffer_init, buffer_free};
import("std::buffer::write")::{buffer_push};
import("std::buffer::read")::{buffer_get};
import("std::buffer::error")::{BUFFER_ERR_BOUNDS};

fun main() -> i32 {
    var data: Buffer;

    if (buffer_init(&data, 0) < 0) {
        return 1;
    }

    if (buffer_push(&data, 0) < 0) {
        buffer_free(&data);
        return 2;
    }

    var value: u8 = 99;

    if (buffer_get(&data, 0, &value) < 0) {
        buffer_free(&data);
        return 3;
    }

    println("stored={}", value);
    value = 99;

    if (buffer_get(&data, 1, &value) == BUFFER_ERR_BOUNDS) {
        println("outside, preserved={}", value);
    }

    if (buffer_free(&data) < 0) {
        return 4;
    }

    return 0;
}
```

Результат выполнения:

```text
stored=0
outside, preserved=99
```

Первое попадание означает успех, равный 0, второе попадание — неудачный выход за пределы поля. Даже если после сбоя остается value=99, это не означает, что это значение, прочитанное из буфера. Обязательно вместе проверьте статус.

## Практическое решение: накопление байтов

Чтобы сложить числа от 0 до 9, повторите buffer_push и проверьте каждый результат. Сохраните сумму в i64 и читайте только диапазон `0 <= index < data.len`. После обработки буфера мы вызываем buffer_free как для успешного, так и для неудачного пути.

<!-- wave-example: book-buffer-sum -->
```wave
import("std::buffer::alloc")::{Buffer, buffer_init, buffer_free};
import("std::buffer::write")::{buffer_push};

fun main() -> i32 {
    var data: Buffer;

    if (buffer_init(&data, 0) < 0) {
        return 1;
    }

    for (var value: i32 = 0; value < 10; value += 1) {
        if (buffer_push(&data, value as u8) < 0) {
            buffer_free(&data);
            return 2;
        }
    }

    var total: i64 = 0;

    for (var index: i64 = 0; index < data.len; index += 1) {
        total += deref data.data[index] as i64;
    }

    println("sum={}", total);

    if (buffer_free(&data) < 0) {
        return 3;
    }

    return 0;
}
```

Результат выполнения:

```text
sum=45
```
