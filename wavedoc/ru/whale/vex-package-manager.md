---
translation_set_id: vex-package-manager
path: whale/vex-package-manager
locale: ru
group: whale
group_order: 1
order: 4
title: Менеджер пакетов Vex
summary: Описывает проекты на основе manifest, Wave, зависимости путей Git·, lockfile, автономные сборки и границы wavec.
---

## роль

Vex — это менеджер пакетов и инструмент сборки для Wave. Vex работает поверх `wavec`. Vex отвечает за структуру проекта и анализ зависимостей, а `wavec` отвечает за флаги компилятора и конвейер компиляции.

Команда Vex основана на manifest. `vex build`, `vex check` и `vex run` намеренно не получают флаг raw `wavec`.

## Создать пакет

```shell
vex init
vex init --lib
```

Приложение использует `src/main.wave`, а библиотека — `src/lib.wave`. Корневая структура пакета выглядит следующим образом:

```text
my_project/
├── src/
│   └── main.wave
├── vex.ws
├── vex.lock
└── .vex/
    └── deps/
```

`vex.ws` становится manifest. Vex не использует manifest в расширении `.wson`.

```wson
{
    name = "my_project",
    version = 0.1.0,
    lib = false,
    description = "my_project Project",
    author = "unknown",
    license = "Unknown",
    dependencies = []
}
```

## команда сборки

```shell
vex build [--target <triple>] [--release] [--dry-run] [--locked] [--offline]
vex check [--target <triple>] [--release] [--dry-run] [--locked] [--offline]
vex run   [--target <triple>] [--release] [--dry-run] [--locked] [--offline] [-- <args...>]
```

Опции для Vex остаются небольшими. Если вам нужен элемент управления, зависящий от компилятора, например emit, linker, CPU, ABI или debug, используйте `wavec` напрямую. Если вам нужно использовать конкретный компилятор, установите `VEX_WAVEC=/path/to/wavec`.

Такие шаги выполнения, как `Resolving`, `Fetching`, `Compiling`, `Checking`, `Running`, `Finished`, выводятся в stderr, а выходные данные программы сохраняются на stdout.

## Git Центральная зависимость

Зависимости Vex указываются как локальные `path` или Git URL. Зависимость может использовать только один из двух методов.

```wson
{
    name = "app",
    version = 0.1.0,
    dependencies = [
        { name = "local_math", path = "../local_math" },
        { name = "remote_math", git = "https://github.com/example/math.git", tag = "v0.1.0" }
    ]
}
```

Зависимость Git может указывать не более одного из `branch`, `tag` или `rev`. Каждый корень зависимости должен иметь свой собственный `vex.ws`. Vex рекурсивно разрешает зависимость manifest, отклоняет конфликтующие идентификаторы пакетов и сохраняет управляемые Git checkout в `.vex/deps/<name>`.

## lockfile Контракт

Схема v2 `vex.lock` записывает весь граф транзитивных зависимостей и точный Git commit. Зафиксируйте с помощью manifest. Использование того же manifest и допустимого lockfile выберет тот же график зависимостей без повторного следования branch или tag.

Команды, требующие зависимостей, автоматически интерпретируются и могут быть подготовлены заранее с помощью следующих команд.

```shell
vex fetch
vex update
vex update math shared_core
```

`vex update` обновляет все пакеты Git или только указанные пакеты и затронутые графы переходов. Несвязанные заблокированные пакеты сохранят выбранный commit.

## Рабочий процесс locked и offline

`--locked` запрещает создание и изменение `vex.lock`. Ошибка завершится, если файл не существует, имеет неподдерживаемую схему или не соответствует графику manifest. Уже прикрепленные к lockfile, commit можно импортировать при необходимости.

`--offline` запрещает все Git сетевые операции. Требуемые checkout и commit уже должны существовать локально.

```shell
vex fetch --locked
vex build --locked --offline
```

Эти две команды представляют собой строгий рабочий процесс CI. Подготовьте заблокированный commit именно тогда, когда сеть доступна, а затем скомпилируйте его без изменения сети или lockfile. dry-run не импортирует зависимости и не переписывает lockfile.

## Настройки и информация компилятора

```shell
vex info
vex setup wavec
vex setup wavec --version <version>
vex --version
```

Vex проверяет схему `wavec` dry-run JSON перед фактической сборкой. Компиляторы, которые не реализуют требуемую схему, не будут выполняться с неизвестным планом и отклонят его с ошибкой совместимости.
