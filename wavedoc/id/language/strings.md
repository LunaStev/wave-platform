---
translation_set_id: learn-strings
path: language/strings
locale: id
group: language
group_order: 2
order: 7
title: 7. String, karakter, dan byte
summary: Bedakan antara string dan char, UTF-8 panjang byte, NUL, penelusuran, dan data biner.
---

## huruf di layar dan byte di memori

Anda melihat huruf di layar, tetapi byte disimpan di memori. Terutama dalam kasus di mana satu karakter memiliki beberapa byte UTF-8, seperti dalam bahasa Korea, mudah untuk membuat kesalahan jika “panjang” dan “jumlah karakter” digunakan secara bergantian.

Bab ini membedakan antara str dan char, berakhiran NUL, escape, lokasi pencarian, dan data biner. Contohnya masing-masing merupakan program yang lengkap.

## String literal dan keluaran

<!-- wave-example: book-string-literals -->
```wave
fun main() {
    var greeting: str = "안녕하세요";

    println("{}", greeting);
    println("line one\nline two");
    println("quote: \"Wave\"");
}
```

Hasil eksekusi:

```text
안녕하세요
line one
line two
quote: "Wave"
```

Karakter biasa dalam tanda kutip ganda dinyatakan sebagai UTF-8. escape menunjukkan byte yang sulit untuk ditulis langsung dari sumbernya. `\n` adalah byte pemisah baris dan tidak mencetak dua karakter, garis miring terbalik dan n. Untuk mencetak garis miring terbalik itu sendiri, gunakan `\\`.

## Panjangnya adalah jumlah byte

<!-- wave-example: book-utf8-length -->
```wave
import("std::string::len")::{
    len
};

fun main() {
    println("ASCII={}", len("Wave"));
    println("Korean={}", len("한"));
    println("mixed={}", len("Wave한"));
}
```

Hasil eksekusi:

```text
ASCII=4
Korean=3
mixed=7
```

`len` menghitung byte sebelum NUL terminasi, bukan karakter yang terlihat. Karakter `한` membutuhkan tiga byte dalam UTF-8. Jumlah karakter, jumlah titik kode Unicode, dan jumlah byte umumnya tidak dapat dipertukarkan. Lebar tampilan juga bergantung pada faktor-faktor seperti font dan penggabungan karakter.

Oleh karena itu, fungsi yang memotong string pada posisi byte sembarang dan menampilkannya di layar harus mempertimbangkan batas Unicode secara terpisah. Tentukan dengan jelas kondisi masukan, apakah itu program yang hanya memproses teks ASCII atau teks umum Unicode.

## NUL Ujung dan Panjang

str menggunakan 0 byte untuk menunjukkan akhir. len tidak menyertakan panjang byte terakhir tersebut. Merupakan kesalahan untuk memasukkan NUL ke dalam string literal. Berikut adalah contoh kesalahan yang disengaja:

```wave
fun main() {
    var text: str = "left\x00right";
}
```

`\xNN` di sumber menentukan satu byte dengan tepat dua digit heksadesimal. `\x41` mewakili byte A yaitu 65. Karena penulisan karakter biasa sebagai UTF-8 dan memasukkan byte sembarang berbeda, tidak semua str yang dapat ditulis sebagai `\xNN` valid UTF-8.

## char tidak berisi seluruh karakter Unicode

char adalah nilai karakter 8-bit yang tidak ditandatangani. Anda dapat menggunakan literal yang mewakili nilai dalam rentang byte tunggal, seperti `'A'`. `'한'` adalah kesalahan karena tidak termasuk dalam kisaran ini. `"한"` adalah str terpisah dengan beberapa UTF-8 byte.

Jangan mencoba untuk selalu memasukkan satu huruf ke dalam satu char. Anda harus terlebih dahulu memutuskan apakah unit yang diperlukan untuk memproses teks adalah byte atau poin kode Unicode.

## perbandingan string

Untuk membandingkan konten string, gunakan fungsi std. Di bawah ini adalah program yang memeriksa kesamaan konten dan perbedaan huruf besar-kecil.

<!-- wave-example: book-string-compare -->
```wave
import("std::string::cmp")::{
    eq, starts_with, ends_with
};

fun main() {
    if (eq("Wave", "Wave")) {
        println("same bytes");
    }

    if (!eq("Wave", "wave")) {
        println("case differs");
    }

    if (starts_with("report.txt", "report") && ends_with("report.txt", ".txt")) {
        println("name matches");
    }
}
```

Hasil eksekusi:

```text
same bytes
case differs
name matches
```

Perbandingan ini membandingkan string byte. Itu tidak secara otomatis melakukan konversi kasus khusus bahasa atau normalisasi Unicode. Bahkan ketika membandingkan nama file, aturan kesetaraan nama file di OS tidak sama dengan perbandingan string sederhana.

## Unit dan kegagalan hasil pencarian

<!-- wave-example: book-string-search -->
```wave
import("std::string::find")::{
    find, count
};

fun main() {
    var first: i32 = find("banana", "na");
    var missing: i32 = find("banana", "xy");
    var matches: i32 = count("aaaa", "aa");

    println("first={} missing={}", first, missing);
    println("matches={}", matches);
}
```

Hasil eksekusi:

```text
first=2 missing=-1
matches=2
```

find mengembalikan posisi pertama atau -1. Indeks 0 juga sukses, jadi dicentang dengan `result >= 0`. count bukanlah posisi, tapi jumlah pertandingan yang tidak tumpang tindih. Di atas, aa adalah nomor 2 karena cocok dengan 0~1 dan 2~3.

Bin needle juga merupakan bagian dari kontrak. find mengembalikan 0, contains mengembalikan true, dan count mengembalikan 0. Jangan berasumsi bahwa hanya karena nama fungsinya ada dalam modul yang sama, bahkan metode pengembaliannya pun sama.

## Menghapus spasi berbeda dengan membuat string baru

trim_range mengembalikan rentang tidak termasuk spasi tanpa mengubah atau menyalin teks asli. Karena kita menerima penunjuk keluaran, pertama-tama kita menyiapkan bilangan bulat untuk menyimpan hasilnya.

<!-- wave-example: book-string-trim-range -->
```wave
import("std::string::trim")::{
    trim_range
};

fun main() {
    var start: i32 = 0;
    var end: i32 = 0;

    trim_range("  Wave  ", &start, &end);

    println("start={} end={} bytes={}", start, end, end - start);
}
```

Hasil eksekusi:

```text
start=2 end=6 bytes=4
```

Kisarannya adalah `[start, end)`. Ini mencakup awal tetapi bukan akhir, jadi panjangnya adalah end-start. Menambahkan start ke alamat awal dokumen asli tidak secara otomatis membuat NUL di lokasi end. Anda perlu membawa jangkauannya secara terpisah atau menyiapkan ruang senar baru.

## Data biner memiliki panjang tersendiri

Data yang mengandung angka nol tidak tercakup dalam aturan akhir string. Ia menggunakan array dan panjang byte.

<!-- wave-example: book-binary-not-string -->
```wave
fun main() {
    var bytes: array<u8, 3> = [65, 0, 66];

    for (var index: i32 = 0; index < 3; index += 1) {
        println("{}", bytes[index]);
    }
}
```

Hasil eksekusi:

```text
65
0
66
```

Angka 0 kedua adalah data aktual. Jika Anda menafsirkan ini sebagai str, maka dianggap berakhir pada 0 pertama dan Anda tidak dapat melihat 66 berikutnya. Sebaliknya, jika Anda mengubah array tanpa NUL menjadi str dengan cast, ada risiko pembacaan di luar array. cast bukan operasi untuk menambahkan byte terminasi.

## Latihan: Memeriksa nama file

Periksa apakah nama file diakhiri dengan `.wave`, dan jika string berisi `test`, keluarkan ke file pengujian. Latihan ini hanya memeriksa pola byte pada nama dan tidak membahas keberadaan file sebenarnya.

### Solusi lengkap

<!-- wave-example: book-string-solution -->
```wave
import("std::string::cmp")::{
    ends_with
};
import("std::string::find")::{
    contains
};

fun classify(name: str) {
    if (!ends_with(name, ".wave")) {
        println("other file");
        return;
    }

    if (contains(name, "test")) {
        println("Wave test file");
    } else {
        println("Wave source file");
    }
}

fun main() {
    classify("main.wave");
    classify("parser_test.wave");
    classify("notes.txt");
}
```

Hasil eksekusi:

```text
Wave source file
Wave test file
other file
```

Terdapat kebijakan terpisah mengenai cara menangani huruf kapital `.WAVE` dan apakah akan menganggapnya sebagai ujian meskipun test disertakan di seluruh jalur. Sekalipun fungsinya kecil, pengoperasiannya hanya dapat dijelaskan secara akurat jika ditentukan input apa yang ditargetkan.
