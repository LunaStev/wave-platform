---
translation_set_id: build-link-targets
path: whale/build-link-targets
locale: id
group: whale
group_order: 1
order: 5
title: Opsi build, penautan, dan target
summary: Menjelaskan emit kiriman, jenis masukan, tautan, target/CPU/ABI dan rencana pembangunan yang berdiri sendiri.
---

## emit Keluaran

```shell
wavec build main.wave --emit=check
wavec build main.wave --emit=ast
wavec build main.wave --emit=ir
wavec build main.wave --emit=bc
wavec build main.wave --emit=asm
wavec build main.wave --emit=obj -o main.o
wavec build main.wave --emit=bin -o app
```

artifact emit Jenisnya adalah `ast`, `ir`, `bc`, `asm`, `obj`, `bin`. `check` adalah mode kontrol inspeksi, bukan tipe pengiriman, dan tidak untuk digunakan dengan artifact emit lainnya.

```shell
wavec print supported-emit-kinds
```

## Jenis masukan dan link-only

Selain sumber Wave, kompiler membedakan antara input IR, bitcode, assembly, object dan archive. Daftar dukungan ditanyakan dengan perintah berikut:

```shell
wavec print supported-input-types
```

Untuk menautkan hanya object atau archive yang sudah dibuat, Anda dapat menggunakan `--input-type` dan `--link-only`.

```shell
wavec build module.o --input-type=obj --link-only --emit=bin -o app
```

## tautan asli

```shell
wavec --link=m -L ./lib build main.wave
```

`--link` menambahkan perpustakaan dan `-L` menambahkan jalur pencarian. Meskipun sebuah simbol dideklarasikan di FFI, perpustakaan yang menyediakan simbol tersebut tidak ditautkan secara otomatis.

## Pilih sasaran

Pilihan untuk memilih OS dan CPU untuk dijalankan.

- `--target <triple>`
- `--cpu <name>`
- `--features <csv>`
- `--abi <name>`
- `--sysroot <path>`

Periksa default host dan target yang didukung dengan perintah berikut:

```shell
wavec print host-target
wavec print supported-targets
wavec print target-spec --target <triple>
wavec print cpu-list --target <triple>
wavec print target-features --target <triple>
```

## Didukung oleh

Program untuk komputer Anda saat ini akan dibuat tanpa menentukan target. Untuk memilih lingkungan yang berbeda, teruskan nama target di bawah ini ke `--target`.

|OS·Lingkungan|arsitektur|nama sasaran|
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
| WebAssembly |64 sedikit| `wasm64-unknown-unknown` |

Untuk program yang berjalan tanpa OS, gunakan `x86_64-unknown-none-elf`, `aarch64-unknown-none-elf`, dan `riscv64-unknown-none-elf`. Daftar lengkap target untuk versi yang diinstal dapat ditemukan di `wavec print supported-targets`.

## RISC-V 64 kontrak

Nilai default untuk target Hosted RISC-V adalah `generic-rv64`, RV64GC, `lp64d` ABI. Nilai default untuk target Freestanding adalah `generic-rv64`, RV64IMAC, dan `lp64`.

```shell
wavec print target-spec --target riscv64-unknown-linux-gnu --format=json
wavec print target-spec --target riscv64-unknown-none-elf --format=json
```

Mendukung RISC-V CPU untuk `generic`, `generic-rv64`, `rocket-rv64`, `sifive-u74`. Feature override adalah `m`, `a`, `f`, `d`, `c`, `zicsr`, `zifencei` Beri tanda di depan nama dan pisahkan dengan koma.

```shell
wavec build main.wave \
  --target riscv64-unknown-linux-gnu \
  --features=+m,+a,+f,-d,+c,+zicsr \
  --abi=lp64f
```

RISC-V Verifikasi menolak kombinasi yang tidak konsisten. `d` memerlukan `f`, dan `f` memerlukan `zicsr`. `lp64`, `lp64f`, `lp64d` harus cocok dengan floating point feature yang telah Anda aktifkan. Jika Anda tidak menentukan ABI secara langsung, kompiler akan mengambil ABI dari feature.

## Tautan Berdiri Bebas

```shell
wavec build kernel.wave \
  --target riscv64-unknown-none-elf \
  --freestanding \
  --entry=_start \
  --linker-script=linker.ld \
  --no-start-files \
  -o kernel.elf
```

`--freestanding` menyesuaikan pengaturan build untuk menghindari penggunaan perpustakaan default. `--entry` menentukan entri linker, `--linker-script` menentukan skrip, dan `--no-start-files` menentukan pengecualian file startup host.

Anda dapat menggunakan `--dry-run` untuk memeriksa rencana tautan sebelum eksekusi sebenarnya.

## Hosted Tautan Silang

Saat membuat program untuk dijalankan di OS·CPU lain, tentukan jalur perpustakaan lingkungan target sebagai sysroot. Jika Anda menggunakan tautan terpisah, tentukan jalurnya sebagai `-C linker`.

```shell
wavec build main.wave \
  --target riscv64-unknown-linux-gnu \
  --sysroot /path/to/riscv64-sysroot \
  -C linker=/path/to/target-linker \
  -o app-riscv64
```

Perpustakaan yang akan ditautkan dengan sysroot disejajarkan dengan OS·CPU·ABI yang dipilih.

## Apa yang harus diperiksa di cross build

- target triple ada dalam daftar dukungan kompiler Anda
- sysroot dan linkernya sesuai dengan target ABI
- Apakah perpustakaan tautan untuk arsitektur target?
- CPU feature valid untuk target CPU
- Jika berdiri bebas, periksa apakah simbol entri dan penempatan memori cocok dengan skrip linker.

## Jalankan WebAssembly

Hasil wasm64 digunakan dalam lingkungan eksekusi yang mendukung memory64. Modul yang menggunakan fungsi eksternal seperti file, waktu, dan input harus terhubung host import yang sesuai dengan fungsi tersebut.

Proses pembuatan dan penyambungan kode target dapat dilihat di `--dry-run`. Eksekusi sebenarnya terjadi di lingkungan eksekusi OS atau WebAssembly yang dipilih.
