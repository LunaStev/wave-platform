---
translation_set_id: string-bytes
path: reference/string-and-bytes
locale: id
group: stdlib
group_order: 1
order: 3
title: string: Panjang, pencarian, dan rentang
summary: NUL Menjelaskan unit byte dari string terminasi API dan nilai yang dikembalikan.
---

## Penyimpanan string dan kondisi argumen

Argumen `str` pada modul ini harus berupa byte terminasi NUL yang dapat diakses. Panjang dan indeks pencarian dalam byte. Karakter biasa disimpan sebagai UTF-8, tetapi pencarian byte adalah Unicode tanpa normalisasi atau pemisahan karakter demi karakter. Jangan berasumsi bahwa indeks yang dikembalikan adalah batas karakter.

## Bandingkan dengan panjangnya

```text
std::string::len
len(s: str) -> i32
is_empty(s: str) -> bool

std::string::cmp
eq(a: str, b: str) -> bool
cmp(a: str, b: str) -> i32
starts_with(s: str, prefix: str) -> bool
ends_with(s: str, suffix: str) -> bool
```

`len` tidak termasuk NUL terakhir. `cmp` Urutan dinilai berdasarkan tanda hasilnya. Nilai yang dikembalikan tidak ditafsirkan sebagai urutan karakter Unicode atau pengurutan kamus khusus bahasa. Fungsi-fungsi ini tidak mengalokasikan memori dan tidak mengubah masukannya.

## pencarian

Dapatkan nama yang Anda butuhkan, seperti `import("std::string::find")::{find, contains, count};`.

|deklarasi fungsi|hasil|
| --- | --- |
| `find(s: str, needle: str) -> i32` |Lokasi pertandingan pertama. -1 jika tidak ada, 0 untuk kosong needle|
| `contains(s: str, needle: str) -> bool` |Termasuk atau tidak. Kosong needle adalah true|
| `count(s: str, needle: str) -> i32` |Jumlah kecocokan yang tidak tumpang tindih. Bin needle adalah 0|
| `find_char(s: str, c: u8) -> i32` |posisi pertama byte atau -1|
| `rfind_char(s: str, c: u8) -> i32` |Posisi terakhir byte atau -1|
| `contains_char(s: str, c: u8) -> bool` |Keberadaan byte itu|
| `count_char(s: str, c: u8) -> i32` |jumlah byte yang dimaksud|

`c` dalam nama `*_char` adalah byte, bukan titik kode Unicode. NUL itu sendiri di akhir string tidak termasuk dalam target pencarian.

## Rentang tidak termasuk spasi

```text
std::string::trim
trim_left_index(s: str) -> i32
trim_right_index(s: str) -> i32
trim_range(s: str, out_start: ptr<i32>, out_end: ptr<i32>)
```

`trim_range` menulis rentang semi-terbuka `[start, end)` tidak termasuk spasi ASCII ke argumen keluaran. Kedua penunjuk keluaran harus menunjuk ke bilangan bulat yang dapat ditulis. Itu tidak mengubah teks asli atau membuat string baru. Jika semuanya kosong, maka akan menjadi rentang kosong.

## Contoh berjalan

Simpan ke `main.wave` dan jalankan `wavec run main.wave`.

<!-- wave-example: string-api -->
```wave
import("std::string::find")::{
    find, count
};
import("std::string::trim")::{
    trim_range
};

fun main() {
    var start: i32 = 0;
    var end: i32 = 0;
    trim_range("  Wave  ", &start, &end);
    println("{} {}", start, end);
    println("{} {}", find("banana", "na"), count("aaaa", "aa"));
}
```

Hasil eksekusi:

```text
2 6
2 2
```

## Fitur terkait

Klasifikasi/konversi kasus `std::string::ascii` adalah untuk rentang ASCII. `djb2_32` dan `fnv1a_64` dari `std::string::hash` tidak digunakan untuk hash kriptografi atau penyimpanan kata sandi. Untuk data yang berisi NUL, gunakan [bytes](/docs/id/stdlib/bytes).

## Pola yang tumpang tindih dengan istilah penelusuran kosong

Melihat perilaku tepi fungsi pencarian dengan nilai sebenarnya memudahkan dalam menentukan kondisi pemanggilan.

<!-- wave-example: book-string-empty-search -->
```wave
import("std::string::find")::{find, contains, count};

fun main() {
    println("empty find={}", find("Wave", ""));
    println("empty count={}", count("Wave", ""));
    println("nonoverlapping={}", count("aaaa", "aa"));

    if (contains("Wave", "")) {
        println("empty needle is contained");
    }
}
```

Hasil eksekusi:

```text
empty find=0
empty count=0
nonoverlapping=2
empty needle is contained
```

Nilai keberhasilan find, 0, menempati posisi pertama. Nilai keberhasilan 0 untuk count adalah hasil dari aturan istilah pencarian yang tidak cocok atau kosong. Tidak ada dua nilai yang diperlakukan sama. Jika pengabaian huruf besar-kecil atau normalisasi Unicode diperlukan, kebijakan terpisah harus diterapkan sebelum dan sesudah pencarian byte ini.

## trim Menyalin rentang ke string baru

Tidak ada akhiran baru NUL dalam rentang yang dikembalikan oleh trim_range. Saat menyalin ke tujuan terpisah, sisakan panjang + 1 spasi dan tulis byte terakhir langsung sebagai 0.

<!-- wave-example: book-trim-copy -->
```wave
import("std::string::trim")::{trim_range};

fun main() -> i32 {
    var original: str = "  Wave  ";
    var start: i32 = 0;
    var end: i32 = 0;
    var destination: array<u8, 16>;

    trim_range(original, &start, &end);

    var length: i32 = end - start;

    if (length >= 16) {
        return 1;
    }

    for (var index: i32 = 0; index < length; index += 1) {
        destination[index] = original[start + index];
    }

    destination[length] = 0;
    println("{}", &destination[0] as str);
    return 0;
}
```

Hasil eksekusi:

```text
Wave
```

String dengan panjang 16 tidak akan cocok dengan tujuan ini. Ini karena Anda memerlukan NUL terakhir. Meskipun panjangnya 0, menulis destination[0]=0 menghasilkan string kosong yang valid. Array lokal tujuan hidup hingga akhir main, jadi kami mencetak di dalamnya.

## String API Urutan penggunaan

Saat mendesain API string, tentukan apakah inputnya diakhiri dengan NUL, apakah indeks menghitung byte, dan apakah hasilnya meminjam rentang sumber atau memiliki alokasi baru. Kisaran yang dipinjam bergantung pada masa pakai sumbernya. Hasil yang dialokasikan harus menentukan siapa yang membebaskannya.

Baca [Bab Pembelajaran String](/docs/id/language/strings) untuk konsep dasar dan [bytes](/docs/id/stdlib/bytes) untuk data termasuk NUL.
