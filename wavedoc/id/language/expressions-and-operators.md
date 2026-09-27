---
translation_set_id: expressions
path: language/expressions-and-operators
locale: id
group: language
group_order: 2
order: 3
title: 3. Aritmatika, perbandingan, dan konversi
summary: Pelajari urutan penghitungan, pembagian bilangan bulat, operasi bitwise, dan cast.
---

## Lihat hasil penghitungan dan jenis penghitungan secara bersamaan

Ekspresi adalah kode yang menghitung suatu nilai. Nama variabel, literal, pemanggilan fungsi, dan ekspresi yang menghubungkan beberapa nilai dengan operator semuanya merupakan ekspresi. Dalam matematika, meskipun persamaannya terlihat sama, hasilnya akan berbeda tergantung apakah persamaan tersebut bilangan bulat atau bilangan real, dan berapa banyak bit yang dikandungnya.

Dimulai dengan penghitungan sederhana, bab ini memperkenalkan tanda kurung, pembagian, operasi logika, operasi bitwise, dan cast. Setiap contoh adalah file main.wave lengkap yang dapat Anda jalankan dengan `wavec run main.wave`.

## Rentang diapit tanda kurung

<!-- wave-example: book-precedence -->
```wave
fun main() {
    var first: i32 = 2 + 3 * 4;
    var second: i32 = (2 + 3) * 4;

    println("{} {}", first, second);
}
```

Hasil eksekusi:

```text
14 20
```

Perkalian dievaluasi sebelum penjumlahan, jadi ekspresi pertama adalah 2+12. Pada persamaan kedua, kita ambil jumlah dalam tanda kurung yaitu 5 lalu dikalikan dengan 4. Tujuannya agar tanda kurung tidak lebih sedikit. Disarankan untuk menggunakannya agar pembaca dapat dengan mudah memahami ruang lingkup perhitungan.

Meskipun menggunakan operator yang sama beberapa kali, arah rangkaian tetap penting. `20 - 5 - 3` adalah `(20 - 5) - 3`, yaitu 12. `20 - (5 - 3)` adalah 18. Urutan lengkap persisnya ada di [referensi operator](/docs/id/language/expressions-and-operators).

## Pembagian bilangan bulat dan sisanya

Membagi dua bilangan bulat tidak menghasilkan hasil pecahan floating-point. Gunakan pembagian untuk hasil bagi dan operator sisa untuk sisanya.

<!-- wave-example: book-division -->
```wave
fun main() {
    var items: i32 = 17;
    var box_size: i32 = 5;
    var full_boxes: i32 = items / box_size;
    var remaining: i32 = items % box_size;

    println("boxes={} remaining={}", full_boxes, remaining);
    println("negative quotient={}", -17 / 5);
}
```

Hasil eksekusi:

```text
boxes=3 remaining=2
negative quotient=-3
```

Membagi 17 item menjadi kelompok yang terdiri dari 5 orang menghasilkan 3 kelompok lengkap dan 2 item tersisa. Pembagian bertanda terpotong menuju nol, jadi -17/5 adalah -3. Hal ini berbeda dengan pembulatan ke bawah menuju tak terhingga negatif.

Anda tidak dapat membagi dengan 0. signed Nilai yang diperoleh dengan membagi nilai minimum dengan -1 tidak termasuk dalam tipe yang sama. Fungsi yang mengambil masukan ini harus diperiksa sebelum membagi atau menggunakan checked matematika API.

## Titik konversi mengubah hasilnya.

Dua ekspresi berikut keduanya disimpan dalam variabel f64, namun proses penghitungannya berbeda.

<!-- wave-example: book-division-conversion -->
```wave
fun main() {
    var already_divided: f64 = (7 / 2) as f64;
    var floating_division: f64 = (7 as f64) / (2 as f64);

    if (already_divided == 3.0) {
        println("integer division lost the fraction");
    }

    if (floating_division == 3.5) {
        println("floating division kept the fraction");
    }
}
```

Hasil eksekusi:

```text
integer division lost the fraction
floating division kept the fraction
```

Ekspresi pertama melakukan pembagian bilangan bulat untuk mendapatkan 3, lalu mengubahnya menjadi f64. Yang kedua mengonversi operan menjadi f64 sebelum melakukan pembagian floating-point. Memilih tipe yang lebih luas untuk variabel akhir tidak dapat memulihkan informasi yang hilang sebelumnya.

Nilai floating-point adalah perkiraan. Dua hasil yang terlihat seperti nilai desimal yang sama mungkin tidak cocok untuk perbandingan persamaan yang tepat. Pilih toleransi yang sesuai dengan unit dan skala permasalahan; satu epsilon tetap tidak cocok untuk setiap perhitungan.

## Perbandingan dilakukan bool

`<`, `<=`, `>`, `>=`, `==`, `!=` periksa hubungannya. Satu tanda sama dengan `=` adalah tugas, dan dua tanda sama dengan `==` adalah perbandingan yang sama.

<!-- wave-example: book-comparisons -->
```wave
fun main() {
    var age: i32 = 20;
    var minimum: i32 = 18;
    var eligible: bool = age >= minimum;

    if (eligible) {
        println("eligible");
    }

    if (age != minimum) {
        println("not exactly the boundary");
    }
}
```

Hasil eksekusi:

```text
eligible
not exactly the boundary
```

“Setidaknya 18” mencakup 18; “lebih besar dari 18” tidak termasuk. Uji 17, 18, dan 19 untuk memeriksa batas ini. Perbandingan tipe campuran bergantung pada penandatanganan dan lebar, jadi mengonversi kedua operan ke tipe yang diinginkan dapat membuat perbandingan lebih jelas.

## Operasi logis dan evaluasi hubung singkat

`&&` memeriksa apakah keduanya benar, `||` memeriksa apakah lebih dari satu yang benar, dan `!` membalik benar/salah. Operasi ini tidak selalu mengeksekusi ekspresi di sisi kanan.

<!-- wave-example: book-short-circuit -->
```wave
fun report() -> bool {
    println("right side evaluated");
    return true;
}

fun main() {
    var enabled: bool = false;

    if (enabled && report()) {
        println("both true");
    }

    if (!enabled || report()) {
        println("at least one true");
    }
}
```

Hasil eksekusi:

```text
at least one true
```

Pada kondisi pertama, enabled salah, jadi tidak perlu melihat ke kanan. Yang kedua, !enabled benar, jadi ruas kanan juga tidak diperlukan. Oleh karena itu, keluaran report tidak pernah muncul.

Anda dapat menggunakan ini untuk memeriksa penyebut sebelum pembagian. Fragmen isi fungsi `if (divisor != 0 && value / divisor > 2) { ... }` tidak melakukan pembagian ketika penyebutnya 0. Namun, pengujian ini saja tidak menyelesaikan batasan lain seperti nilai minimum signed/-1.

`&&` lebih diutamakan daripada `||`. Dalam kebijakan yang kompleks, tunjukkan maksud Anda dalam tanda kurung, seperti `(member && active) || admin`.

## Mempersempit atau Melebarkan Bilangan Bulat

`as` adalah pemeran eksplisit. Bit tingkat tinggi yang dibuang saat mempersempit bilangan bulat tidak dapat diperoleh kembali dengan melebarkannya lagi.

<!-- wave-example: book-cast-chain -->
```wave
fun main() {
    var original: i32 = 300;
    var small: u8 = original as u8;
    var widened: i32 = small as i32;

    println("{} -> {} -> {}", original, small, widened);

    var negative: i8 = -1;
    var signed_value: i32 = negative as i32;
    var same_bits: u8 = negative as u8;

    println("{} {}", signed_value, same_bits);
}
```

Hasil eksekusi:

```text
300 -> 44 -> 44
-1 255
```

8 bit terbawah dari 300 adalah 44. Jika Anda memperluas -1 ke tipe signed, tandanya akan diperluas untuk mempertahankan -1. Jika Anda menafsirkan 8 bit yang sama dengan unsigned, hasilnya adalah 255.

Transformasi yang berupaya mempertahankan nilai dalam rentang dan transformasi yang berupaya memanipulasi bit penyimpanan memiliki tujuan berbeda. Jika Anda memiliki input pengguna, periksa dulu apakah itu rentang tujuan dan konversikan. Kehadiran cast tidak menjamin bahwa nilainya berada dalam kisaran aman.

## Ganti dengan bool

Mengonversi bilangan bulat menjadi bool menghasilkan false untuk nol dan true sebaliknya. Ini tidak terpotong ke bit terendah: 2 juga diubah menjadi true. Untuk nilai floating-point, hanya +0,0 dan -0,0 yang dikonversi menjadi false; semua nilai lainnya, termasuk NaN dan tak terhingga, diubah menjadi true.

<!-- wave-example: book-bool-conversion -->
```wave
fun main() {
    var zero: i32 = 0;
    var two: i32 = 2;

    if (!(zero as bool)) {
        println("zero is false");
    }

    if (two as bool) {
        println("two is true");
    }
}
```

Hasil eksekusi:

```text
zero is false
two is true
```

Konversi pointer-ke-bool tidak didukung. Bandingkan penunjuk dengan null secara eksplisit, misalnya dengan `pointer != null`. Apakah alamat non-null aman untuk dibaca adalah pertanyaan tersendiri.

## operasi sedikit

`&`, `|`, `^`, `~` mencakup setiap bit bilangan bulat. Ini dapat digunakan untuk menyatakan izin atau fungsi dalam bit. Di bawah, 1 adalah izin baca dan 2 adalah izin menulis.

<!-- wave-example: book-bit-flags -->
```wave
const READ: u32 = 1;
const WRITE: u32 = 2;

fun main() {
    var permissions: u32 = READ | WRITE;

    if ((permissions & WRITE) != 0) {
        println("write enabled");
    }

    permissions = permissions & ~WRITE;
    println("remaining={}", permissions);
}
```

Hasil eksekusi:

```text
write enabled
remaining=1
```

Tambahkan sedikit dengan OR dan periksa keberadaan bit tertentu dengan AND. Buat topeng dengan hanya bit yang relevan menjadi 0 dengan `~WRITE` dan hapus. Operasi bitwise `&`·`|` adalah operator yang berbeda dari evaluasi hubung singkat `&&`·`||` dari bool.

## Lebar dan jumlah shift

Pergeseran ke kiri memindahkan bit ke kiri dan membuang bit tinggi di luar lebar operan. Pergeseran ke kanan menggunakan ekstensi tanda untuk nilai yang ditandatangani dan ekstensi nol untuk nilai yang tidak ditandatangani. Hasilnya selalu bertipe operan kiri.

<!-- wave-example: book-shifts -->
```wave
fun main() {
    var bits: u8 = 129;
    var shifted: u8 = bits << 1;
    var negative: i32 = -8;
    var positive: u32 = 8;

    println("left={}", shifted);
    println("right={} {}", negative >> 1, positive >> 1);
}
```

Hasil eksekusi:

```text
left=2
right=-4 4
```

Menggeser nilai u8 129 ke kiri sebanyak satu akan membuang bit tertingginya dan menyisakan 2. Nilai shift-count asli harus non-negatif dan kurang dari lebar bit operan kiri. Untuk u8, penghitungan yang valid adalah 0 hingga 7. Penghitungan konstan yang tidak valid merupakan kesalahan waktu kompilasi; hitungan runtime yang tidak valid menyebabkan jebakan.

## Mengubah nilai floating-point menjadi bilangan bulat

Konversi titik mengambang ke bilangan bulat terlebih dahulu terpotong menuju nol, lalu memeriksa rentang bilangan bulat tujuan. Hasil NaN, tak terhingga, dan di luar jangkauan tidak valid. Konversi konstan yang tidak valid menghasilkan kesalahan waktu kompilasi; konversi runtime yang tidak valid menyebabkan jebakan.

Perangkap tidak mengembalikan nilai kesalahan dari fungsi tersebut. Untuk kegagalan konversi yang dapat dipulihkan, rancang antarmuka yang memeriksa rentang sebelum melakukan konversi. Untuk mendapatkan bit-bit yang disimpan dari nilai titik-mengambang, gunakan fungsi konversi bit di `std::math::float` alih-alih menggunakan angka.

## Latihan dan solusi

Tulis fungsi yang membagi 137 won menjadi 50 unit won dan sisanya, dan hanya mengubah bilangan bulat dalam rentang 0 hingga 255 hingga u8. Contoh ini menunjukkan langkah validasi sebagai fungsi terpisah yang menunjukkan di luar jangkauan sebagai -1.

<!-- wave-example: book-expression-solution -->
```wave
fun checked_byte(value: i32) -> i32 {
    if (value < 0 || value > 255) {
        return -1;
    }

    var small: u8 = value as u8;
    return small as i32;
}

fun main() {
    var money: i32 = 137;

    println("coins={} remainder={}", money / 50, money % 50);
    println("{} {} {}", checked_byte(0), checked_byte(255), checked_byte(256));
}
```

Hasil eksekusi:

```text
coins=2 remainder=37
0 255 -1
```

Alasan mengapa -1 dapat digunakan sebagai indikator kegagalan adalah karena rentang keberhasilannya adalah 0 hingga 255. Jika bilangan bulat mana pun dapat menjadi nilai keberhasilan, diperlukan representasi hasil yang berbeda. Kami melanjutkan desain ini di bab penanganan kesalahan nanti.


## prioritas

Prioritas operator adalah sebagai berikut, dimulai dari yang tertinggi:

1. Ekspresi dasar dan akses postfix: pemanggilan fungsi, akses field, pengindeksan, postfix `++`·`--`
2. Operasi unary: `!`, `~`, `&`, `deref`, awalan `++`·`--`, unary `+`·`-`
3. `as` ketik konversi
4. `*`, `/`, `%`
5. `+`, `-`
6. `<<`, `>>`
7. `<`, `<=`, `>`, `>=`
8. `==`, `!=`
9. Sedikit `&`
10. Sedikit `^`
11. Sedikit `|`
12. `&&`
13. `||`
14. Tugas dan tugas majemuk

Penugasan rantai digabungkan dari kanan. Saat menggabungkan berbagai jenis operator, gunakan tanda kurung untuk menyatakan dengan jelas urutan evaluasi.

## Subyek yang bisa diganti

Penugasan, `++`, dan `--` memerlukan ekspresi yang menunjukkan lokasi penyimpanan, seperti variabel, bidang, elemen array, atau penunjuk dereferensi. Menulis ke `const` tidak diperbolehkan.

## bergeser

Hasil pergeseran selalu bertipe operan kiri; tipe operan yang tepat tidak memperluas komputasi. Pergeseran kiri membuang bit tinggi melebihi lebar tersebut. Pergeseran kanan tanda-perluas nilai yang ditandatangani dan perpanjang nol nilai yang tidak ditandatangani.

Jumlah shift harus berupa bilangan bulat yang nilai aslinya memenuhi `0 <= n < LHS bit width`. Itu diperiksa sebelum pemotongan apa pun ke tipe yang lebih kecil. Penghitungan konstan yang tidak valid adalah kesalahan waktu kompilasi; hitungan runtime yang tidak valid menyebabkan jebakan.

## Konversi yang melibatkan nilai bool dan floating-point

Konversi bilangan bulat ke bool menghasilkan false untuk nol dan true untuk setiap nilai lainnya. Konversi floating-point-ke-bool menghasilkan false hanya untuk +0,0 dan -0,0; NaN dan tak terhingga positif atau negatif menghasilkan true. Konversi pointer-ke-bool tidak didukung: bandingkan secara eksplisit dengan `pointer != null`.

Konversi titik mengambang ke bilangan bulat terpotong menuju nol dan kemudian memeriksa rentang tujuan. Hasil NaN, tak terhingga, dan di luar jangkauan tidak valid. Konversi konstan yang tidak valid adalah kesalahan waktu kompilasi; konversi runtime yang tidak valid menyebabkan jebakan. Jebakan bukanlah pengembalian kesalahan yang dapat dipulihkan.

`&&` dan `||` adalah evaluasi hubung singkat. Efek samping dari operan kanan yang tidak dijalankan tidak terjadi. Anda dapat memeriksa hasilnya dengan nilai kecil di [kelas aritmatika](/docs/id/language/expressions-and-operators).

## shift yang gagal dengan sengaja

Jumlah pergeseran untuk nilai 8-bit harus 0 sampai 7. Program di bawah ini harus ditolak sebelum dieksekusi.

<!-- wave-example: reject-shift-count -->
```wave
fun main() {
    var value: u8 = 1;
    var result: u8 = value << 8;
}
```
