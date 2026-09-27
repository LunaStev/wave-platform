---
translation_set_id: stdlib-random
path: stdlib/random
locale: id
group: stdlib
group_order: 1
order: 10
title: random: Mengisi buffer dengan keacakan OS
summary: OS Isi buffer dengan entropi dan tangani kegagalan sebagian.
---

## API

```text
std::random::fill
random_available() -> bool
random_fill(buffer: ptr<u8>, size: i64) -> RandomFillResult
```

size adalah jumlah byte, dan pemanggil menyediakan penyimpanan. `RandomFillResult` berisi ok, tertulis, dan kesalahan. Jika berhasil, tulisan sama dengan panjang yang diminta. Jika gagal, tulisan mengidentifikasi awalan yang valid dan terisi; jangan gunakan sisa byte sebagai data acak.

`random_available` memberi tahu Anda apakah fungsi angka acak OS didukung. Keberhasilan permintaan individu diperiksa oleh hasil random_fill. Hanya menggunakan entropi OS dan tidak mundur pada nilai waktu atau lemah PRNG pada kegagalan. size=0 akan berhasil meskipun diteruskan dengan null. null adalah kesalahan untuk panjang negatif atau positif.

## Contoh berjalan

<!-- wave-example: random-api -->
```wave
import("std::random::fill")::{
    RandomFillResult, random_fill
};

fun main() -> i32 {
    var bytes: array<u8, 16>;
    var result: RandomFillResult = random_fill(&bytes[0], 16);
    if (!result.ok) {
        println("random error={}", result.error);
        return 1;
    }

    println("filled={}", result.written);
    return 0;
}
```

Hasil eksekusi:

```text
filled=16
```

Simpan sebagai `main.wave` dan jalankan. Isi byte berbeda setiap saat, jadi tidak ada nilai spesifik yang diharapkan. Jika gagal, periksa penyebabnya dengan result.error. Daripada mengeluarkan byte acak secara harfiah, gunakan pengkodean terpisah jika perlu.

## Jika permintaan Anda salah

Permintaan byte nol berhasil karena tidak ada yang perlu ditulis. Melewati null dengan panjang positif gagal karena tidak ada buffer tujuan. Program berikut membandingkan kasus-kasus ini tanpa mengalokasikan memori.

<!-- wave-example: book-random-boundaries -->
```wave
import("std::random::fill")::{
    RandomFillResult,
    random_fill
};

fun main() -> i32 {
    var empty: RandomFillResult = random_fill(null, 0);

    if (!empty.ok || empty.written != 0) {
        return 1;
    }

    println("empty request succeeded");

    var invalid: RandomFillResult = random_fill(null, 16);

    if (invalid.ok) {
        return 2;
    }

    println("missing buffer rejected");
    return 0;
}
```

Hasil eksekusi:

```text
empty request succeeded
missing buffer rejected
```

## Menangani buffer yang terisi sebagian

Jika 16 byte diminta tetapi gagal dan written=8 dikembalikan, hanya 8 byte pertama yang terisi. Jika tugasnya adalah membuat pengidentifikasi 16-byte, itu bukan pengidentifikasi yang berhasil, jadi kami membuang seluruh hasil dan melaporkan kegagalan. Anda tidak boleh mengisi sisa 8 byte dengan 0 dan kemudian menganggapnya sukses.

Memetakan byte acak ke rentang bilangan bulat memerlukan kehati-hatian. Menerapkan `% 10` pada nilai u8 yang terdistribusi merata akan membuat 0–5 lebih mungkin terjadi dibandingkan 6–9, karena 256 tidak habis dibagi 10. Untuk menghilangkan bias ini, tolak nilai 250–255, gambar lagi, dan terapkan operasi sisanya hanya pada nilai yang diterima.

Penyimpanan dan masa pakai byte acak dikelola oleh pemanggil. Saat menggunakan array, itu diproses dalam lingkup array, dan saat menggunakan memori dinamis, memori tersebut dibebaskan setelah digunakan. Struktur pengembalian tidak memiliki buffer.
