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

## Opsional IR socket

Perintah `ir` hanya disertakan ketika Whale dibuat dengan fitur `socket-cli`.

```shell
cargo run -p whale --features socket-cli -- ir lower program.json
cargo run -p whale --features socket-cli -- ir lower program.json -o program.wir
```

`ir lower` membaca JSON dari Whale socket schema, mengubahnya menjadi Whale IR, dan memverifikasi modul. Teks IR dikeluarkan ke jalur stdout atau `-o`. `--target <triple>` menggantikan string target dan `--no-verify` menghilangkan validasi.

Untuk menggunakan perintah IR, Whale harus dibuat dengan fitur `socket-cli`. Socket JSON produsen dan Whale harus menggunakan socket schema version yang sama.
