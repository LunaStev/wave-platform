---
translation_set_id: stdlib-path-env
path: stdlib/path-env
locale: ru
group: stdlib
group_order: 1
order: 12
title: path и env: Пути и настройки окружения
summary: Считывает переменные пути и среды в буфере вызывающего объекта и определяет ошибки емкости.
---

## комбинация путей

```text
std::path::copy
path_join2(dst: ptr<u8>, dst_cap: i32, left: str, right: str) -> i32
path_basename_copy(dst: ptr<u8>, dst_cap: i32, path: str) -> i32
path_dirname_copy(dst: ptr<u8>, dst_cap: i32, path: str) -> i32
```

Емкость включает последнее пространство NUL. Результатом успеха является любая длина, за исключением NUL, неудачей является -1. Только в случае успеха мы будем использовать пункт назначения в виде строки. Эти функции работают со строками пути и не проверяют существование файла или права доступа. Объединение путей само по себе не предотвращает экранирование каталогов и не проверяет подлинность реальных файлов.

## переменная среды

```text
std::env::environ
env_get(name: str, dst: ptr<u8>, dst_cap: i64) -> i64
env_exists(name: str) -> bool
env_get_i64(name: str) -> EnvResult<i64>
env_get_i32(name: str) -> EnvResult<i32>
```

`env_get` возвращает длину, исключая NUL в случае успеха. Буфер вызывающего объекта должен иметь возможность хранить до NUL. Пустое значение отличается от ошибки без ключа, поскольку оно может иметь длину 0.

Получите и распознайте ошибки от `std::env::consts` до NOT_FOUND, NO_SPACE, INVALID_KEY, READ, SOURCE_INCOMPLETE, NO_MEMORY. Не рассматривайте выход из буфера как отсутствующий ключ. Чтобы найти числа, отметьте ok в результате, а затем используйте value. Не рассматривайте автоматически содержимое переменных среды как доверенные параметры; проверьте их область действия и тип.

В следующем примере сочетаются каталоги data и имена файлов input.txt.

<!-- wave-example: path-api -->
```wave
import("std::path::copy")::{
    path_join2
};

fun main() -> i32 {
    var output: array<u8, 64>;
    var length: i32 = path_join2(&output[0], 64, "data", "input.txt");
    if (length < 0) {
        return 1;
    }

    println("{}", &output[0] as str);
    return 0;
}
```

Результат выполнения:

```text
data/input.txt
```

## Разделить каталоги и имена файлов

В следующем примере копируется путь, разделенный на два буфера. Исходный файл не обязательно должен существовать.

<!-- wave-example: book-path-parts -->
```wave
import("std::path::copy")::{
    path_basename_copy,
    path_dirname_copy
};

fun main() -> i32 {
    var directory: array<u8, 64>;
    var filename: array<u8, 64>;
    var directory_length: i32 = path_dirname_copy(&directory[0], 64, "data/report.txt");
    var filename_length: i32 = path_basename_copy(&filename[0], 64, "data/report.txt");

    if (directory_length < 0 || filename_length < 0) {
        return 1;
    }

    println("directory={}", &directory[0] as str);
    println("filename={}", &filename[0] as str);
    return 0;
}
```

Результат выполнения:

```text
directory=data
filename=report.txt
```

Оба буфера будут действовать до конца main. `as str` читает тот же буфер, что и строка, без выделения новой строки. Следовательно, если вы измените буфер, строка, считанная по этому адресу, также изменится.

## Установить настройки по умолчанию

Переменные среды — это настройки, передаваемые вне программы. При чтении числового параметра установите флажок «Можно ли его прочитать как целое число?» и «Находится ли это в пределах, разрешенных этой программой?»

<!-- wave-example: book-env-setting -->
```wave
import("std::env::environ")::{
    EnvResult,
    env_get_i32
};

fun main() -> i32 {
    var setting: EnvResult<i32> = env_get_i32("WAVE_EXAMPLE_WORKERS");
    var workers: i32 = 4;

    if (setting.ok) {
        if (setting.value < 1 || setting.value > 32) {
            println("workers must be between 1 and 32");
            return 1;
        }

        workers = setting.value;
    }

    println("workers={}", workers);
    return 0;
}
```

Если `WAVE_EXAMPLE_WORKERS` отсутствует или не может быть прочитан как целое число, используется значение по умолчанию — 4. Если установлено целое число от 1 до 32, используется это значение, а если это целое число вне диапазона, оно завершается с ошибкой.

Linux/macOS:

```shell
WAVE_EXAMPLE_WORKERS=8 wavec run main.wave
```

PowerShell:

```powershell
$env:WAVE_EXAMPLE_WORKERS = "8"
wavec run main.wave
```

В обоих случаях выводится `workers=8`. В приведенном выше примере выбрана простая политика по умолчанию. Если это обязательный параметр, рассматривайте ошибки числового поиска как ошибки, а не заменяйте их значениями по умолчанию. Если вам нужно отличить отсутствующий ключ, недостаточный буфер и ошибку чтения, используйте константы env_get и ENV_ERR_*.
