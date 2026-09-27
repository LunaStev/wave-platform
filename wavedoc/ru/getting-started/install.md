---
translation_set_id: install
path: getting-started/install
locale: ru
group: getting-started
group_order: 1
order: 2
title: Установка Wave
summary: Установите Wave в Linux, macOS или Windows и запустите первую программу.
---

## Linux и macOS

Выполните следующую команду в терминале. Она установит Wave вместе с менеджером пакетов Vex.

```shell
curl -fsSL https://wave-lang.dev/install.sh | bash -s -- latest
```

После установки откройте новый терминал и проверьте версию.

```shell
wavec --version
```

## Windows

Выполните следующие команды в PowerShell. Они установят Wave вместе с менеджером пакетов Vex.

```powershell
irm https://wave-lang.dev/install.ps1 -OutFile install.ps1
powershell -ExecutionPolicy Bypass -File .\install.ps1 -Latest
```

После установки откройте новое окно PowerShell и проверьте версии.

```powershell
wavec --version
vex --version
```

## Первый запуск

Сохраните следующий код в файле `main.wave`.

<!-- wave-example: install-stdlib -->
```wave
import("std::string::len")::{
    len
};

fun main() {
    println("Wave: {} bytes", len("Wave"));
}
```

Запустите программу из каталога, в котором сохранён файл.

```shell
wavec run main.wave
```

Результат:

```text
Wave: 4 bytes
```

Если появится сообщение о том, что стандартная библиотека не найдена, установите её и запустите программу снова.

```shell
wavec install std
wavec run main.wave
```

[Далее: Первая программа](/docs/ru/language/program-structure) · [Решение проблем](/docs/ru/reference/diagnostics)
