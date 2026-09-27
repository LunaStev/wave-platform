---
translation_set_id: compiler
path: getting-started/compiler
locale: ru
group: getting-started
group_order: 1
order: 3
title: Справочник команд компилятора
summary: wavec Описывает команды, конвейер сборки, выходные данные, цели, диагностику, связывание зависимостей и запросы инструментов.
---

## командная модель

`wavec` — компилятор CLI. Он напрямую компилирует отдельные входные данные, предоставляет инструментам информацию о поддержке компилятора и управляет установленными источниками стандартной библиотеки.

```text
wavec [global-options] <command> [command-options]
```

|команда|Использование|
| --- | --- |
| `wavec build <input...>` |В зависимости от флагов он выполняет проверку, генерацию кода, связывание или выполнение конвейеров.|
| `wavec check <file>` |Никнейм для `build <file> --emit=check`.|
| `wavec run <file> [-- <args...>]` |Это псевдоним для `build <file> --run`, и аргумент после `--` передается программе.|
| `wavec print <item>` |Запрашивает информацию о поддержке целевого объекта и цепочки инструментов.|
| `wavec install std` |Установите стандартную библиотеку.|
| `wavec update std` |Обновите установленные стандартные библиотеки.|
| `wavec --version` |Распечатывает информацию об установленной версии.|

Полный список команд и опций можно найти по адресу `wavec --help`.

## Сборка, тестирование и запуск

```shell
wavec build main.wave
wavec check main.wave
wavec run main.wave -- first-argument second-argument
```

`build` по умолчанию создает исполняемый файл. `check` останавливается после завершения проверки интерфейса. `run` требует двоичного вывода и не может использоваться со сборками общих библиотек.

Используйте `--dry-run`, чтобы проверить запрос и определить, какие шаги следует выполнить без компиляции, связывания или выполнения.

```shell
wavec build main.wave --target riscv64-unknown-linux-gnu --dry-run
wavec build main.wave --dry-run --error-format=json
```

Формат JSON — это стабильный унифицированный интерфейс, используемый такими инструментами сборки, как Vex.

## emit и тип входа

```shell
wavec build main.wave --emit=ast
wavec build main.wave --emit=ir
wavec build main.wave --emit=bc
wavec build main.wave --emit=asm
wavec build main.wave --emit=obj -o main.o
wavec build main.wave --emit=bin -o app
```

Типы вывода emit: `ast`, `ir`, `bc`, `asm`, `obj`, `bin`. `check` — это режим управления, который следует использовать отдельно. Вы можете указать несколько типов вывода, которые принимает конвейер, через запятую.

Типы ввода: `wave`, `ir`, `bc`, `asm`, `obj`, `archive`. `--input-type=<kind>` принудительно указывает тип всех входных данных. При связывании только входов object или archive используйте двоичные файлы emit и `--link-only`.

```shell
wavec build module.o --input-type=obj --link-only --emit=bin -o app
```

## выходное местоположение

|варианты|эффект|
| --- | --- |
| `-o <file>` |Указывает основной путь вывода.|
| `--out-dir <dir>` |emit Помещает выходные данные в указанный каталог.|
| `--target-dir <dir>` |Определяет промежуточные и основные маршруты доставки.|

## Оптимизация и диагностический вывод

```shell
wavec -O2 build main.wave
wavec --debug-wave=tokens,ast build main.wave
```

Шаги оптимизации: `-O0`, `-O1`, `-O2`, `-O3`, `-Os`, `-Oz`, `-Ofast`. За `--debug-wave` могут следовать `tokens`, `ast`, `ir`, `mc`, `hex`, `all`, а несколько шагов можно объединять запятыми.

## собственная ссылка

```shell
wavec --link=m -L ./lib build main.wave
wavec build main.wave --shared -o libexample.so
wavec build main.wave --static -o app
wavec build main.wave --pie -o app
```

`--link=<lib>` добавляет собственные библиотеки, а `-L <path>` добавляет пути поиска. В режимах связи используются `--shared`, `--static`, `--pie`, `--no-pie` в соответствии с правилами совместимости.

Опции управления серверной частью и компоновщиком включают в себя:

- `--target`, `--cpu`, `--features`, `--abi`, `--sysroot`
- `-C linker=<path>` и `-C link-arg=<arg>`
- `-C link-sysroot=<path>` и `-C relocation-model=<model>`
- `-C no-default-libs`

Автономные выходные данные, такие как ядра, используют `--freestanding` вместе с настройками `--entry`, `--linker-script` и `--no-start-files`, соответствующими среде.

## Интерпретация внешней упаковки

```shell
wavec --dep-root .vex/deps build main.wave
wavec --dep math=/opt/wave-deps/math build main.wave
```

`--dep-root <dir>` добавляет корень для поиска внешних `package::module` import. `--dep <name>=<path>` фиксирует имя пакета в одном каталоге. Это точка интеграции компилятора, а проекты manifest, загрузки зависимостей и lockfile обрабатываются Vex.

## Поддержка функции запроса

Инструменты, использующие целевой или выходной тип, могут запрашивать информацию о поддержке с помощью `wavec print`.

```shell
wavec print host-target
wavec print target-spec --format=json
wavec print supported-targets
wavec print supported-input-types
wavec print supported-emit-kinds
wavec print supported-print-items
wavec print cpu-list --target riscv64-unknown-linux-gnu
wavec print target-features --target riscv64-unknown-linux-gnu
wavec print default-linker
wavec print sysroot
wavec print std-path
wavec print dep-search-paths
```

Вы также можете запрашивать такие элементы, как `host`, `default-target` и `target-list`. Элементы, поддерживающие структурированный вывод, получают `--format=json`.

## Граница между компилятором и набором инструментов

`wavec` отвечает за проверку исходного кода, генерацию кода и связывание. Vex отвечает за создание воспроизводимых пакетов с пакетом manifest, графом зависимостей и lockfile. Whale — это набор инструментов низкого уровня, который работает независимо.

## std Укажите путь

```shell
wavec --std-root /absolute/path/to/std check main.wave
wavec --std-root /absolute/path/to/std run main.wave
```

Указанный путь std имеет приоритет над путем установки и завершится ошибкой, если путь неверен или несовместим с std. Он не заменяет автоматически std из других установок. Выберите std, соответствующий вашему компилятору.

Для вывода `-o` используется путь, отличный от исходного/входного файла. Поскольку `check` не проверяет операцию во время выполнения, [практика](/docs/ru/practice/input-calculator) также проверяет результаты выполнения.
