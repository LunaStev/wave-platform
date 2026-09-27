---
translation_set_id: build-link-targets
path: whale/build-link-targets
locale: ru
group: whale
group_order: 1
order: 5
title: Параметры сборки, компоновки и целевой платформы
summary: Описывает результаты emit, типы входных данных, ссылки, target/CPU/ABI и план автономной сборки.
---

## emit Выход

```shell
wavec build main.wave --emit=check
wavec build main.wave --emit=ast
wavec build main.wave --emit=ir
wavec build main.wave --emit=bc
wavec build main.wave --emit=asm
wavec build main.wave --emit=obj -o main.o
wavec build main.wave --emit=bin -o app
```

artifact emit Типы: `ast`, `ir`, `bc`, `asm`, `obj`, `bin`. `check` — это режим контроля проверки, а не тип поставки, и его нельзя использовать с другим artifact emit.

```shell
wavec print supported-emit-kinds
```

## Тип входа и link-only

Помимо источника Wave, компилятор различает входы IR, bitcode, assembly, object и archive. Список поддержки запрашивается с помощью следующей команды:

```shell
wavec print supported-input-types
```

Чтобы связать только уже созданные object или archive, вы можете использовать `--input-type` и `--link-only`.

```shell
wavec build module.o --input-type=obj --link-only --emit=bin -o app
```

## собственная ссылка

```shell
wavec --link=m -L ./lib build main.wave
```

`--link` добавляет библиотеку, а `-L` добавляет путь поиска. Даже если символ объявлен в FFI, библиотека, предоставляющая этот символ, не подключается автоматически.

## Выберите цель

Возможность выбрать OS и CPU для запуска.

- `--target <triple>`
- `--cpu <name>`
- `--features <csv>`
- `--abi <name>`
- `--sysroot <path>`

Проверьте настройки хоста по умолчанию и поддерживаемые цели с помощью следующей команды:

```shell
wavec print host-target
wavec print supported-targets
wavec print target-spec --target <triple>
wavec print cpu-list --target <triple>
wavec print target-features --target <triple>
```

## При поддержке

Программа для вашего текущего компьютера будет собрана без указания target. Чтобы выбрать другую среду, передайте указанное ниже имя цели в `--target`.

|OS·Окружающая среда|архитектура|целевое имя|
| --- | --- | --- |
| Linux | amd64 | `x86_64-unknown-linux-gnu` |
| Linux | ARM64 | `aarch64-unknown-linux-gnu` |
| Linux | RISC-V64 | `riscv64-unknown-linux-gnu` |
| Linux | LoongArch64 | `loongarch64-unknown-linux-gnu` |
| macOS | amd64 | `x86_64-apple-darwin` |
| macOS | ARM64 | `aarch64-apple-darwin` |
| Windows | amd64 | `x86_64-pc-windows-msvc` |
| Windows | ARM64 | `aarch64-pc-windows-msvc` |
| FreeBSD | amd64 | `x86_64-unknown-freebsd` |
| WebAssembly |64 бит| `wasm64-unknown-unknown` |

Для программ, которые работают без OS, используйте `x86_64-unknown-none-elf`, `aarch64-unknown-none-elf` и `riscv64-unknown-none-elf`. Полный список целей для установленных версий можно найти по адресу `wavec print supported-targets`.

## RISC-V 64 контракта

Значения по умолчанию для целей Hosted RISC-V: `generic-rv64`, RV64GC, `lp64d` ABI. Значения по умолчанию для цели Freestanding: `generic-rv64`, RV64IMAC и `lp64`.

```shell
wavec print target-spec --target riscv64-unknown-linux-gnu --format=json
wavec print target-spec --target riscv64-unknown-none-elf --format=json
```

Поддерживает RISC-V CPU для `generic`, `generic-rv64`, `rocket-rv64`, `sifive-u74`. Feature override — это `m`, `a`, `f`, `d`, `c`, `zicsr`, `zifencei` Поставьте знак перед именем и разделите его запятой.

```shell
wavec build main.wave \
  --target riscv64-unknown-linux-gnu \
  --features=+m,+a,+f,-d,+c,+zicsr \
  --abi=lp64f
```

RISC-V Проверка отклоняет несовместимые комбинации. Для `d` требуется `f`, а для `f` требуется `zicsr`. `lp64`, `lp64f`, `lp64d` должны соответствовать значению с плавающей запятой feature, которое вы включили. Если вы не укажете ABI напрямую, компилятор выведет ABI из feature.

## Отдельно стоящая ссылка

```shell
wavec build kernel.wave \
  --target riscv64-unknown-none-elf \
  --freestanding \
  --entry=_start \
  --linker-script=linker.ld \
  --no-start-files \
  -o kernel.elf
```

`--freestanding` корректирует настройки сборки, чтобы избежать использования библиотек по умолчанию. `--entry` указывает записи компоновщика, `--linker-script` определяет сценарии, а `--no-start-files` определяет исключения файлов запуска хоста.

Вы можете использовать `--dry-run`, чтобы проверить план связи перед фактическим выполнением.

## Hosted Перекрестная ссылка

При создании программы для запуска на другом OS·CPU укажите путь к библиотеке целевой среды как sysroot. Если вы используете отдельный компоновщик, укажите путь как `-C linker`.

```shell
wavec build main.wave \
  --target riscv64-unknown-linux-gnu \
  --sysroot /path/to/riscv64-sysroot \
  -C linker=/path/to/target-linker \
  -o app-riscv64
```

Библиотеки, которые будут связаны с sysroot, соответствуют выбранному OS·CPU· ABI.

## Что проверить в кросс-билде

- target triple находится в списке поддержки вашего компилятора.
- sysroot и линкер соответствуют целевому ABI
- Поддерживает ли библиотека ссылок целевую архитектуру?
- CPU feature действительно для цели CPU
- Если он стоит отдельно, проверьте, соответствуют ли символ входа и размещение в памяти сценарию компоновщика.

## Запустите WebAssembly

Результаты wasm64 используются в среде выполнения, поддерживающей memory64. Модули, использующие внешние функции, такие как файл, время и ввод, должны подключаться к host import, соответствующим этой функции.

Процесс создания и подключения целевого кода можно найти в `--dry-run`. Фактическое выполнение происходит в выбранной среде выполнения OS или WebAssembly.
