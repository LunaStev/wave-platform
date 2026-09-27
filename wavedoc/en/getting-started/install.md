---
translation_set_id: install
path: getting-started/install
locale: en
group: getting-started
group_order: 1
order: 2
title: Install Wave
summary: Install Wave on Linux, macOS, or Windows and run your first program.
---

## Linux and macOS

Run the following command in a terminal. It installs Wave and the Vex package manager together.

```shell
curl -fsSL https://wave-lang.dev/install.sh | bash -s -- latest
```

After installation, open a new terminal and check the installed version.

```shell
wavec --version
```

## Windows

Run the following commands in PowerShell. They install Wave and the Vex package manager together.

```powershell
irm https://wave-lang.dev/install.ps1 -OutFile install.ps1
powershell -ExecutionPolicy Bypass -File .\install.ps1 -Latest
```

After installation, open a new PowerShell window and check the installed versions.

```powershell
wavec --version
vex --version
```

## Run your first program

Save the following code as `main.wave`.

<!-- wave-example: install-stdlib -->
```wave
import("std::string::len")::{
    len
};

fun main() {
    println("Wave: {} bytes", len("Wave"));
}
```

Run it from the directory where you saved the file.

```shell
wavec run main.wave
```

Output:

```text
Wave: 4 bytes
```

If Wave reports that it cannot find the standard library, install it and run the program again.

```shell
wavec install std
wavec run main.wave
```

[Next: Your first program](/docs/en/language/program-structure) · [Troubleshooting](/docs/en/reference/diagnostics)
