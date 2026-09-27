---
translation_set_id: string-bytes
path: reference/string-and-bytes
locale: ru
group: stdlib
group_order: 1
order: 3
title: string: длина, поиск и диапазоны.
summary: NUL Описывает байтовую единицу завершающей строки API и возвращаемое значение.
---

## Хранение строк и условия аргументов

Аргумент `str` этого модуля должен быть доступным байтом завершения NUL. Длина и индекс поиска указаны в байтах. Обычные символы сохраняются как UTF-8, а поиск по байтам — Unicode без нормализации или посимвольного разделения. Не предполагайте, что возвращаемый индекс является границей символа.

## Сравните с длиной

```text
std::string::len
len(s: str) -> i32
is_empty(s: str) -> bool

std::string::cmp
eq(a: str, b: str) -> bool
cmp(a: str, b: str) -> i32
starts_with(s: str, prefix: str) -> bool
ends_with(s: str, suffix: str) -> bool
```

`len` исключает последний NUL. `cmp` Порядок определяется знаком результата. Возвращаемое значение не интерпретируется как порядок символов Unicode или сортировка по словарю для конкретного языка. Эти функции не выделяют память и не меняют вводимые данные.

## поиск

Получите нужное имя, например `import("std::string::find")::{find, contains, count};`.

|объявление функции|результат|
| --- | --- |
| `find(s: str, needle: str) -> i32` |Место первого матча. -1, если нет, 0, если пусто needle|
| `contains(s: str, needle: str) -> bool` |Включено или нет. Пустой needle равен true|
| `count(s: str, needle: str) -> i32` |Количество непересекающихся совпадений. Бин needle равен 0|
| `find_char(s: str, c: u8) -> i32` |первая позиция байта или -1|
| `rfind_char(s: str, c: u8) -> i32` |Последняя позиция байта или -1|
| `contains_char(s: str, c: u8) -> bool` |Существование этого байта|
| `count_char(s: str, c: u8) -> i32` |количество рассматриваемых байтов|

`c` в имени `*_char` — это байт, а не кодовая точка Unicode. Сам символ NUL в конце строки не включается в цель поиска.

## Диапазон без пробелов

```text
std::string::trim
trim_left_index(s: str) -> i32
trim_right_index(s: str) -> i32
trim_range(s: str, out_start: ptr<i32>, out_end: ptr<i32>)
```

`trim_range` записывает полуоткрытый диапазон `[start, end)` исключая пробелы ASCII в выходной аргумент. Оба указателя вывода должны указывать на записываемые целые числа. Он не изменяет исходный текст и не создает новые строки. Если все пусто, это становится пустым диапазоном.

## Запуск примера

Сохраните его в `main.wave` и запустите `wavec run main.wave`.

<!-- wave-example: string-api -->
```wave
import("std::string::find")::{
    find, count
};
import("std::string::trim")::{
    trim_range
};

fun main() {
    var start: i32 = 0;
    var end: i32 = 0;
    trim_range("  Wave  ", &start, &end);
    println("{} {}", start, end);
    println("{} {}", find("banana", "na"), count("aaaa", "aa"));
}
```

Результат выполнения:

```text
2 6
2 2
```

## Связанные функции

Классификация/преобразование регистра `std::string::ascii` относится к диапазону ASCII. `djb2_32` и `fnv1a_64` из `std::string::hash` не используются для криптографических хэшей или хранения паролей. Для данных, содержащих NUL, используйте [bytes](/docs/ru/stdlib/bytes).

## Шаблоны, которые пересекаются с пустыми поисковыми запросами

Наблюдение краевого поведения функции поиска с фактическими значениями облегчает определение условий вызова.

<!-- wave-example: book-string-empty-search -->
```wave
import("std::string::find")::{find, contains, count};

fun main() {
    println("empty find={}", find("Wave", ""));
    println("empty count={}", count("Wave", ""));
    println("nonoverlapping={}", count("aaaa", "aa"));

    if (contains("Wave", "")) {
        println("empty needle is contained");
    }
}
```

Результат выполнения:

```text
empty find=0
empty count=0
nonoverlapping=2
empty needle is contained
```

Значение успеха find, 0, является первой позицией. Значение успеха 0 для count является результатом отсутствия совпадения или пустого правила поиска. Никакие две ценности не рассматриваются одинаково. Если требуется игнорирование регистра или нормализация Unicode, до и после этого поиска байтов необходимо реализовать отдельные политики.

## trim Копирование диапазона в новую строку

В диапазоне, возвращаемом trim_range, нет нового окончания NUL. При копировании в отдельный пункт назначения зарезервируйте длину + 1 пробел и запишите последний байт непосредственно как 0.

<!-- wave-example: book-trim-copy -->
```wave
import("std::string::trim")::{trim_range};

fun main() -> i32 {
    var original: str = "  Wave  ";
    var start: i32 = 0;
    var end: i32 = 0;
    var destination: array<u8, 16>;

    trim_range(original, &start, &end);

    var length: i32 = end - start;

    if (length >= 16) {
        return 1;
    }

    for (var index: i32 = 0; index < length; index += 1) {
        destination[index] = original[start + index];
    }

    destination[length] = 0;
    println("{}", &destination[0] as str);
    return 0;
}
```

Результат выполнения:

```text
Wave
```

Строка длиной 16 не поместится в это место назначения. Это потому, что вам нужен последний NUL. Даже если длина равна 0, запись destination[0]=0 приводит к появлению допустимой пустой строки. Локальный массив назначения существует до конца main, поэтому мы печатаем внутри него.

## Строка API Порядок использования.

При разработке строкового API укажите, завершается ли его ввод NUL, подсчитывают ли индексы байты и заимствует ли результат исходный диапазон или владеет новым распределением. Заимствованный диапазон зависит от времени жизни источника. Выделенный результат должен указывать, кто его освобождает.

Прочтите [Глава обучения струнам](/docs/ru/language/strings) для получения основных понятий и [bytes](/docs/ru/stdlib/bytes) для получения данных, включая NUL.
