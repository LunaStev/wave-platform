---
translation_set_id: ecosystem
path: whale/ecosystem
locale: id
group: whale
group_order: 1
order: 2
title: Komponen rantai alat
summary: Menjelaskan peran rantai alat tingkat rendah yang terpisah, Whale, dan batasan komponen ekosistem Wave.
---

## WhaleIran

Whale adalah rantai alat tingkat rendah yang menangani perakitan dan representasi perantara. Komponen yang menangani rakitan, objek, tautan, dan representasi perantara dirancang agar dapat digunakan kembali di Wave dan alat pembuat kode asli lainnya.

Whale bukan nama untuk keseluruhan lingkungan pengembangan Wave. Tanggung jawab untuk setiap proyek dibagi sebagai berikut:

|proyek|tanggung jawab|
| --- | --- |
| `wavec` |Wave Memeriksa sumber dan membuat file yang dapat dieksekusi.|
| Vex |Mengelola paket Wave, manifest, grafik ketergantungan, lockfile dan pembuatan paket.|
| Whale |Menyediakan komponen independen assembler, object, linker dan IR.|
| Wave `std` |Runtime dan sistem API disediakan sebagai modul sumber Wave.|

## komponen

Whale workspace terdiri dari empat area perpustakaan utama:

- `assembler`: Tokenisasi, AMD64 Parsing/Encoding, section, symbol dan relocation
- `object`: model file objek dan ELF64 writer
- `linker`: Lapisan tautan
- `ir`: Whale IR ketik, builder, keluaran, verifikasi dan opsional frontend socket

Eksekusi `whale` memberi wilayah ini perintah `asm`, `object`, `link`, dan `ir`.

## batas alat

Program Wave dibuat sebagai `wavec`. Saat berhadapan langsung dengan rakitan, file objek, dan IR, gunakan perintah `whale`.

Menginstal Whale tidak mengubah metode build `wavec`. Vex menggunakan `wavec` untuk membuat paket Wave, dan Whale menjalankannya secara langsung dalam tugas yang menangani artefak tingkat rendah.

## Verifikasi yang dapat disampaikan

Saat menautkan artefak Whale ke proses pembangunan Anda, pastikan bahwa object format dan target architecture cocok. symbol dan relocation dapat diperiksa dengan alat independen seperti `readelf` dan `objdump`. Build yang menggunakan IR socket harus menggunakan socket schema dari produsen yang sama dengan Whale.
