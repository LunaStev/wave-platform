---
translation_set_id: install
path: getting-started/install
locale: id
group: getting-started
group_order: 1
order: 2
title: Menginstal Wave
summary: Instal Wave di Linux, macOS, atau Windows, lalu jalankan program pertama Anda.
---

## Linux dan macOS

Jalankan perintah berikut di terminal. Perintah ini menginstal Wave beserta pengelola paket Vex.

```shell
curl -fsSL https://wave-lang.dev/install.sh | bash -s -- latest
```

Setelah instalasi selesai, buka terminal baru dan periksa versinya.

```shell
wavec --version
```

## Windows

Jalankan perintah berikut di PowerShell. Perintah ini menginstal Wave beserta pengelola paket Vex.

```powershell
irm https://wave-lang.dev/install.ps1 -OutFile install.ps1
powershell -ExecutionPolicy Bypass -File .\install.ps1 -Latest
```

Setelah instalasi selesai, buka jendela PowerShell baru dan periksa versinya.

```powershell
wavec --version
vex --version
```

## Menjalankan program pertama

Simpan kode berikut sebagai `main.wave`.

<!-- wave-example: install-stdlib -->
```wave
import("std::string::len")::{
    len
};

fun main() {
    println("Wave: {} bytes", len("Wave"));
}
```

Jalankan program dari direktori tempat Anda menyimpan berkas tersebut.

```shell
wavec run main.wave
```

Hasil:

```text
Wave: 4 bytes
```

Jika muncul pesan bahwa pustaka standar tidak ditemukan, instal pustaka tersebut lalu jalankan kembali programnya.

```shell
wavec install std
wavec run main.wave
```

[Berikutnya: Program pertama](/docs/id/language/program-structure) · [Pemecahan masalah](/docs/id/reference/diagnostics)
