---
translation_set_id: whale-cli
path: whale/whale-cli
locale: id
group: whale
group_order: 1
order: 3
title: Referensi perintah Whale
summary: Menjelaskan Whale assembler, object wrapper, keluaran diagnostik, dan perintah opsional IR.
---

## Whale Bangun

Di repositori Whale, jalankan:

```shell
cargo build --release
```

File yang dapat dieksekusi tingkat atas memiliki empat kelompok perintah:

```text
whale asm [--amd64 | --aarch64] <input> -o <output>
whale object <input> -o <output>
whale ir <subcommand> [options]
```

## AMD64 assembler

```shell
whale asm --amd64 input.asm -o output.o
```

AMD64 assembler menerima jalur `.o` sebagai keluaran, dan ELF64 relocatable berisi section, symbol dan relocation Buat object.

Aktifkan keluaran diagnostik terperinci dengan `--debug-whale`.

```shell
whale asm --amd64 input.asm -o output.o \
  --debug-whale --token --ast --bytes --dump-hex --stats
```

Tanda diagnostik mencakup `--token`, `--ast`, `--bytes`, `--dump-hex`, `--dump-bin`, `--dump-json`, dan `--stats`. `--trace` mencetak kemajuan pemrosesan.

## Object wrapper

```shell
whale object input.bin -o output.o
```

Perintah `object` menempatkan byte mentah di bagian ELF64 `.text` dan menambahkan simbol global `start` pada offset 0. Perintah ini membungkus kode mesin mentah dalam file objek ELF.

## Verifikasi dan pencetakan IR teks

Build bawaan membaca dan memverifikasi typed IR format 3. Simpan contoh lengkap di [referensi IR](ir-reference) sebagai `answer.wir`. `print` memverifikasi sebelum mencetak bentuk kanonis dan mempertahankan file lama jika gagal. Perintah ini tidak menjalankan IR atau menghasilkan kode native.

```shell
whale ir verify answer.wir
whale ir print answer.wir -o canonical.wir
```

## Opsional IR socket

AST JSON `ir lower` memerlukan feature `socket-cli`. IR teks `verify` dan `print` tidak memerlukannya.

```shell
cargo run -p whale --features socket-cli -- ir lower program.json
cargo run -p whale --features socket-cli -- ir lower program.json -o program.wir
```

`ir lower` membaca JSON dari Whale socket schema, mengubahnya menjadi Whale IR, dan memverifikasi modul. Teks IR dikeluarkan ke jalur stdout atau `-o`. `--target <triple>` menggantikan string target dan `--no-verify` menghilangkan validasi.

Build dengan `socket-cli` untuk memakai `ir lower`. Produsen Socket JSON dan Whale harus memakai AST schema version yang sama.
