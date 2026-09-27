---
translation_set_id: vex-package-manager
path: whale/vex-package-manager
locale: id
group: whale
group_order: 1
order: 4
title: Pengelola paket Vex
summary: Menjelaskan proyek manifest berbasis Wave, dependensi jalur Git·, lockfile, pembangunan offline, dan batasan wavec.
---

## peran

Vex adalah manajer paket dan alat pembuatan untuk Wave. Vex beroperasi di atas `wavec`. Vex bertanggung jawab atas struktur proyek dan analisis ketergantungan, dan `wavec` bertanggung jawab atas flag compiler dan pipeline kompilasi.

Perintah Vex didasarkan pada manifest. `vex build`, `vex check`, dan `vex run` sengaja tidak menerima bendera raw `wavec`.

## Buat paket

```shell
vex init
vex init --lib
```

Aplikasi menggunakan `src/main.wave` dan perpustakaan menggunakan `src/lib.wave`. Struktur root paket adalah sebagai berikut:

```text
my_project/
├── src/
│   └── main.wave
├── vex.ws
├── vex.lock
└── .vex/
    └── deps/
```

`vex.ws` menjadi manifest. Vex tidak menggunakan manifest di ekstensi `.wson`.

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

## perintah membangun

```shell
vex build [--target <triple>] [--release] [--dry-run] [--locked] [--offline]
vex check [--target <triple>] [--release] [--dry-run] [--locked] [--offline]
vex run   [--target <triple>] [--release] [--dry-run] [--locked] [--offline] [-- <args...>]
```

Opsi untuk Vex dibuat kecil. Jika Anda memerlukan kontrol yang bergantung pada kompiler, seperti emit, linker, CPU, ABI, atau debug, gunakan `wavec` secara langsung. Saat Anda perlu menggunakan kompiler tertentu, setel `VEX_WAVEC=/path/to/wavec`.

Langkah-langkah kemajuan seperti `Resolving`, `Fetching`, `Compiling`, `Checking`, `Running`, `Finished` dikeluarkan di stderr, dan keluaran program Dipertahankan pada stdout.

## Git Ketergantungan Pusat

Dependensi Vex ditentukan sebagai `path` lokal atau Git URL. Ketergantungan hanya dapat menggunakan salah satu dari dua metode tersebut.

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

Ketergantungan Git hanya dapat menentukan paling banyak satu dari `branch`, `tag`, atau `rev`. Setiap akar ketergantungan harus memiliki `vex.ws` sendiri. Vex menyelesaikan ketergantungan manifest secara rekursif, menolak identitas paket yang bertentangan, dan menyimpan Git checkout yang dikelola di `.vex/deps/<name>`.

## lockfile Kontrak

Skema v2 `vex.lock` mencatat seluruh grafik ketergantungan transitif dan Git commit yang tepat. Berkomitmen dengan manifest. Menggunakan manifest yang sama dan lockfile yang valid akan memilih grafik ketergantungan yang sama tanpa mengikuti branch atau tag lagi.

Perintah yang memerlukan dependensi diinterpretasikan secara otomatis dan dapat disiapkan terlebih dahulu dengan perintah berikut.

```shell
vex fetch
vex update
vex update math shared_core
```

`vex update` memperbarui semua paket Git, atau hanya paket tertentu dan grafik transisi yang terpengaruh. Paket terkunci yang tidak terkait akan tetap memilih commit.

## Alur kerja locked dan offline

`--locked` melarang pembuatan dan modifikasi `vex.lock`. Ini akan gagal jika file tidak ada, memiliki skema yang tidak didukung, atau tidak cocok dengan grafik manifest. Sudah disematkan ke lockfile, commit dapat diimpor bila diperlukan.

`--offline` melarang semua operasi jaringan Git. checkout dan commit yang diperlukan seharusnya sudah ada secara lokal.

```shell
vex fetch --locked
vex build --locked --offline
```

Kedua perintah ini adalah alur kerja CI yang ketat. Siapkan commit yang terkunci tepat saat jaringan tersedia, lalu kompilasi tanpa mengubah jaringan atau lockfile. dry-run tidak mengimpor dependensi atau menulis ulang lockfile.

## Pengaturan dan informasi kompiler

```shell
vex info
vex setup wavec
vex setup wavec --version <version>
vex --version
```

Vex memvalidasi skema `wavec` dry-run JSON sebelum pembuatan sebenarnya. Kompiler yang tidak mengimplementasikan skema yang diperlukan tidak akan mengeksekusi rencana yang tidak diketahui dan akan menolaknya dengan kesalahan kompatibilitas.
