---
translation_set_id: compiler
path: getting-started/compiler
locale: id
group: getting-started
group_order: 1
order: 3
title: Referensi perintah kompiler
summary: wavec Menjelaskan perintah, membangun pipeline, output, target, diagnostik, tautan ketergantungan, dan kueri alat.
---

## model perintah

`wavec` adalah kompiler CLI. Ini mengkompilasi input individual secara langsung, memberikan informasi dukungan kompiler ke alat, dan mengelola sumber perpustakaan standar yang diinstal.

```text
wavec [global-options] <command> [command-options]
```

|perintah|Gunakan|
| --- | --- |
| `wavec build <input...>` |Bergantung pada flagnya, ia melakukan inspeksi, pembuatan kode, penautan, atau alur eksekusi.|
| `wavec check <file>` |Nama panggilan untuk `build <file> --emit=check`.|
| `wavec run <file> [-- <args...>]` |Ini adalah alias untuk `build <file> --run`, dan argumen setelah `--` diteruskan ke program.|
| `wavec print <item>` |Target kueri dan informasi dukungan rantai alat.|
| `wavec install std` |Instal perpustakaan standar.|
| `wavec update std` |Perbarui perpustakaan standar yang diinstal.|
| `wavec --version` |Mencetak informasi versi yang diinstal.|

Daftar lengkap perintah dan opsi dapat ditemukan di `wavec --help`.

## Bangun, uji, dan jalankan

```shell
wavec build main.wave
wavec check main.wave
wavec run main.wave -- first-argument second-argument
```

`build` membuat file yang dapat dieksekusi secara default. `check` berhenti setelah menyelesaikan pemeriksaan frontend. `run` memerlukan keluaran biner dan tidak dapat digunakan dengan build perpustakaan bersama.

Gunakan `--dry-run` untuk memverifikasi permintaan dan menentukan langkah mana yang harus dijalankan tanpa kompilasi, penautan, atau eksekusi.

```shell
wavec build main.wave --target riscv64-unknown-linux-gnu --dry-run
wavec build main.wave --dry-run --error-format=json
```

Format JSON adalah antarmuka yang stabil dan terpadu yang digunakan oleh alat pembangunan seperti Vex.

## emit dan jenis masukan

```shell
wavec build main.wave --emit=ast
wavec build main.wave --emit=ir
wavec build main.wave --emit=bc
wavec build main.wave --emit=asm
wavec build main.wave --emit=obj -o main.o
wavec build main.wave --emit=bin -o app
```

Jenis keluaran emit adalah `ast`, `ir`, `bc`, `asm`, `obj`, `bin`. `check` adalah mode kontrol dan harus digunakan sendiri. Anda dapat menentukan beberapa tipe keluaran yang diterima alur, dipisahkan dengan koma.

Jenis masukannya adalah `wave`, `ir`, `bc`, `asm`, `obj`, `archive`. `--input-type=<kind>` memaksa jenis semua input ditentukan. Saat menautkan hanya input object atau archive, gunakan biner emit dan `--link-only`.

```shell
wavec build module.o --input-type=obj --link-only --emit=bin -o app
```

## lokasi keluaran

|pilihan|efek|
| --- | --- |
| `-o <file>` |Menentukan jalur keluaran utama.|
| `--out-dir <dir>` |emit Menempatkan output di direktori yang ditentukan.|
| `--target-dir <dir>` |Menentukan rute pengiriman perantara dan utama.|

## Optimalisasi dan keluaran diagnostik

```shell
wavec -O2 build main.wave
wavec --debug-wave=tokens,ast build main.wave
```

Langkah optimasinya adalah `-O0`, `-O1`, `-O2`, `-O3`, `-Os`, `-Oz`, `-Ofast`. `--debug-wave` dapat diikuti oleh `tokens`, `ast`, `ir`, `mc`, `hex`, `all`, dan beberapa langkah dapat digabungkan dengan koma.

## tautan asli

```shell
wavec --link=m -L ./lib build main.wave
wavec build main.wave --shared -o libexample.so
wavec build main.wave --static -o app
wavec build main.wave --pie -o app
```

`--link=<lib>` menambahkan perpustakaan asli dan `-L <path>` menambahkan jalur pencarian. Mode tautan menggunakan `--shared`, `--static`, `--pie`, `--no-pie` sesuai dengan aturan kompatibilitas.

Opsi kontrol backend dan linker meliputi:

- `--target`, `--cpu`, `--features`, `--abi`, `--sysroot`
- `-C linker=<path>` dan `-C link-arg=<arg>`
- `-C link-sysroot=<path>` dan `-C relocation-model=<model>`
- `-C no-default-libs`

Output yang berdiri sendiri seperti kernel menggunakan pengaturan `--freestanding` bersama dengan `--entry`, `--linker-script`, dan `--no-start-files` yang sesuai dengan lingkungan.

## Interpretasi paket eksternal

```shell
wavec --dep-root .vex/deps build main.wave
wavec --dep math=/opt/wave-deps/math build main.wave
```

`--dep-root <dir>` menambahkan root untuk menemukan `package::module` import eksternal. `--dep <name>=<path>` memperbaiki nama paket ke satu direktori. Ini adalah titik integrasi kompiler dan proyek manifest, unduhan ketergantungan dan lockfile ditangani oleh Vex.

## Permintaan fungsi dukungan

Alat yang menggunakan tipe target atau keluaran dapat menanyakan informasi dukungan dengan `wavec print`.

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

Anda juga dapat menanyakan item seperti `host`, `default-target`, dan `target-list`. Item yang mendukung keluaran terstruktur menerima `--format=json`.

## Batas antara kompiler dan rantai alat

`wavec` bertanggung jawab atas pemeriksaan sumber, pembuatan kode, dan penautan. Vex bertanggung jawab untuk membuat paket yang dapat direproduksi dengan paket manifest, grafik ketergantungan, dan lockfile. Whale adalah rantai alat tingkat rendah yang berjalan secara independen.

## std Tentukan jalur

```shell
wavec --std-root /absolute/path/to/std check main.wave
wavec --std-root /absolute/path/to/std run main.wave
```

Jalur std yang ditentukan lebih diutamakan daripada jalur instalasi, dan akan gagal jika jalur tersebut salah atau tidak kompatibel dengan std. Itu tidak secara otomatis menggantikan std dari instalasi lain. Pilih std yang sesuai dengan kompiler Anda.

Untuk keluaran `-o`, jalur yang berbeda digunakan dari file sumber/input. Karena `check` tidak memeriksa operasi runtime, [berlatih](/docs/id/practice/input-calculator) juga memeriksa hasil eksekusi.
