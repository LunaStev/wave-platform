---
translation_set_id: standard-library
path: reference/standard-library
locale: id
group: stdlib
group_order: 1
order: 1
title: Panduan perpustakaan standar
summary: Cara menemukan modul yang sesuai dengan tujuan Anda dan membaca kesalahan serta aturan kepemilikan fungsi tersebut.
---

## Temukan fitur yang Anda perlukan

Pustaka standarnya adalah import dengan jalur `std::module::file`. Meskipun namanya serupa, fungsi mungkin mengembalikan kesalahan dengan cara yang berbeda. Baca terlebih dahulu [API Cara membaca](/docs/id/stdlib/contracts), lalu masuk ke modul yang Anda perlukan pada tabel berikut.

|Apa yang ingin saya lakukan|dokumen|Utama import|
| --- | --- | --- |
|Panjang string/perbandingan/pencarian| [string](/docs/id/reference/string-and-bytes) | `std::string::len`, `cmp`, `find`, `trim` |
|Alokasi memori/salinan/ukuran| [mem](/docs/id/reference/memory-and-buffer) | `std::mem::alloc`, `ops`, `layout` |
|Daftar byte dengan berbagai ukuran| [buffer](/docs/id/stdlib/buffer) | `std::buffer::alloc`, `read`, `write` |
|Baca/tulis biner| [bytes](/docs/id/stdlib/bytes) | `std::bytes::types`, `cursor`, `leb128` |
|File·Deskriptor I/O| [fs dan io](/docs/id/stdlib/files-io) | `std::fs::file`, `std::io::fd` |
|Kombinasi jalur/pengaturan lingkungan| [path dan env](/docs/id/stdlib/path-env) | `std::path::copy`, `std::env::environ` |
|Pengukuran waktu/menunggu| [time](/docs/id/stdlib/time) | `std::time::duration`, `clock`, `sleep` |
|Pencarian daftar alamat/nama numerik| [net.resolve](/docs/id/stdlib/resolution) | `std::net::resolve`, `resolve_table` |
|TCP Koneksi/Transmisi| [net.tcp](/docs/id/stdlib/tcp) | `std::net::tcp`, `address`, `error` |
|OS Nomor acak| [random](/docs/id/stdlib/random) | `std::random::fill` |
|Proses·OS Batas| [fungsi sistem](/docs/id/reference/system-io-network-process) | `std::process::core`, `spawn`, `std::sys` |
|Eksekusi tugas asinkron| [task](/docs/id/stdlib/task) | `std::task` |
|Asisten Matematika/Diagnosis| [math dan debug](/docs/id/stdlib/math-debug) | `std::math::int`, `float`, `std::debug::core` |

## Contoh penggunaan pertama

Program di bawah ini menggunakan satu fungsi dari std tanpa mengunduh paket terpisah. Simpan sebagai `main.wave` dan jalankan sebagai `wavec run main.wave`.

<!-- wave-example: library-import -->
```wave
import("std::string::len")::{
    len
};

fun main() {
    println("{}", len("Wave"));
}
```

Hasilnya adalah `4`. Untuk memahami contoh yang sama langkah demi langkah, baca [string](/docs/id/language/strings).

## Kompatibel dengan kompiler std

Konfirmasikan jalur yang dipilih dengan `wavec print std-path`. Saat menggunakan std dari checkout lain, tentukan jalurnya sebagai `wavec --std-root /absolute/path/to/std check main.wave`. Jika jalur yang ditentukan tidak valid atau tidak kompatibel, kesalahan akan ditampilkan.

## perbatasan platform

Bedakan antara fungsi komputasi seperti string/byte dan fungsi OS seperti file/socket. Mengenali target tidak menjamin bahwa semua host API akan disediakan. Baca entri platform untuk [Sasaran dukungan](/docs/id/whale/build-link-targets) dan masing-masing API secara bersamaan. `std::sys` adalah antarmuka tingkat rendah dan program portabel akan menggunakan modul tingkat lebih tinggi terlebih dahulu.
