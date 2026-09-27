---
translation_set_id: learn-arrays-strings
path: language/arrays
locale: id
group: language
group_order: 2
order: 6
title: 6. Array dan iterasi
summary: Pelajari array berukuran tetap, pengindeksan, traversal, penyalinan, dan pencarian.
---

## Beberapa nilai dari tipe yang sama

Jika Anda membuat tiga skor secara terpisah sebagai score1, score2, dan score3, deklarasi dan perhitungan harus diubah ketika angkanya berubah. Array mengelompokkan sejumlah elemen dengan tipe yang sama. Dengan menggunakan loop, Anda bisa menerapkan aturan yang sama ke setiap elemen.

Bab ini mencakup pembuatan, pengindeksan, modifikasi, iterasi, pencarian, dan agregasi array. String juga mendukung pengindeksan, namun maknanya berbeda, sehingga string dibahas pada bab berikutnya.

## Masukkan panjang dalam tipe

<!-- wave-example: book-array-create -->
```wave
fun main() {
    var scores: array<i32, 3> = [70, 80, 90];

    println("first={}", scores[0]);
    println("second={}", scores[1]);
    println("last={}", scores[2]);
}
```

Hasil eksekusi:

```text
first=70
second=80
last=90
```

i32 di `array<i32, 3>` adalah jenis elemen dan 3 adalah jumlah elemen. Yang kita simpan adalah 3 bilangan bulat. Ini tidak berarti jumlah byte adalah 3. Jumlah elemen dalam literal array harus sesuai dengan panjang yang dideklarasikan.

Indeks dimulai dari 0. Elemen pertama adalah 0, elemen terakhir panjangnya-1. scores[3] adalah akses di luar cakupan, bukan elemen ketiga.

## elemen perubahan

<!-- wave-example: book-array-update -->
```wave
fun main() {
    var scores: array<i32, 3> = [70, 80, 90];

    scores[1] = 85;
    scores[0] += 5;

    println("{} {} {}", scores[0], scores[1], scores[2]);
}
```

Hasil eksekusi:

```text
75 85 90
```

Ubah penyimpanan elemen tertentu tanpa membuat ulang seluruh array. Ekspresi indeks juga bisa menjadi hasil penghitungan, namun Anda harus memastikan bahwa nilainya berada dalam rentang tertentu. Saat menggunakan input eksternal sebagai indeks, angka negatif dan batas atas akan diperiksa.

## Iterasi pada array

<!-- wave-example: book-array-sum -->
```wave
fun main() {
    var scores: array<i32, 4> = [60, 70, 80, 90];
    var total: i32 = 0;

    for (var index: i32 = 0; index < 4; index += 1) {
        total += scores[index];
    }

    println("total={} average={}", total, total / 4);
}
```

Hasil eksekusi:

```text
total=300 average=75
```

Setiap iterasi membaca elemen pada indeks yang berbeda. Variabel jumlah harus diinisialisasi di luar iterasi. Menginisialisasinya ke 0 setiap kali dalam badan perulangan akan menghasilkan hasil yang salah, seperti hanya menyisakan elemen terakhir.

Pembagian bilangan bulat yang digunakan untuk menghitung rata-rata membuang bagian pecahan. Untuk rata-rata floating-point, konversikan jumlahnya sebelum membagi. Untuk array yang lebih besar atau nilai yang lebih besar, pastikan juga bahwa jenis akumulator dapat mewakili jumlahnya.

## Menggabungkan hanya beberapa elemen

Pemfilteran dapat dicapai dengan menggabungkan pernyataan kondisional dan traversal. Di sini kita menghitung jumlah elemen dengan skor 80 atau lebih tinggi.

<!-- wave-example: book-array-filter -->
```wave
fun main() {
    var scores: array<i32, 5> = [60, 80, 90, 75, 100];
    var passed: i32 = 0;

    for (var index: i32 = 0; index < 5; index += 1) {
        if (scores[index] >= 80) {
            passed += 1;
        }
    }

    println("passed={}", passed);
}
```

Hasil eksekusi:

```text
passed=3
```

Nilai indeks dan elemen harus dipisahkan. Memeriksa `index >= 80` membandingkan posisi, bukan skor. Keduanya bisa berupa i32, sehingga sulit untuk menemukan kesalahan semantik ini berdasarkan jenisnya saja.

## Temukan lokasi pertandingan pertama

Pertama, putuskan bagaimana menampilkan hasil yang tidak ditemukan. Dalam contoh ini, indeks yang valid adalah 0 hingga 4, jadi kami menggunakan -1 sebagai penanda kegagalan.

<!-- wave-example: book-array-find -->
```wave
fun main() {
    var values: array<i32, 5> = [8, 3, 8, 1, 5];
    var target: i32 = 8;
    var found: i32 = -1;

    for (var index: i32 = 0; index < 5; index += 1) {
        if (values[index] == target) {
            found = index;
            break;
        }
    }

    if (found >= 0) {
        println("found at {}", found);
    } else {
        println("not found");
    }
}
```

Hasil eksekusi:

```text
found at 0
```

Posisi pertama 0 juga merupakan hasil normal. Jika Anda memeriksa keberhasilan dengan `found > 0`, Anda akan keliru karena tidak menemukan elemen pertama. Untuk alasan yang sama, menilai keberhasilan dengan menggunakan bool sebagai cast adalah salah.

Jika Anda menghapus break, pertandingan berikutnya akan menimpa found, sehingga memberi Anda posisi pertandingan terakhir. Karena satu pernyataan dapat mengubah kontrak suatu fungsi, maka uraian “pencarian” juga harus ditulis secara spesifik apakah ia berada di posisi pertama atau terakhir.

## Menyalin elemen array

Untuk menyalin nilai array ke ruang penyimpanan lain, Anda dapat membaca dan menetapkannya elemen demi elemen. Meskipun Anda mengubah satu elemen bilangan bulat setelah menyalin, elemen bilangan bulat lainnya tidak berubah.

<!-- wave-example: book-array-copy -->
```wave
fun main() {
    var original: array<i32, 3> = [1, 2, 3];
    var copied: array<i32, 3>;

    for (var index: i32 = 0; index < 3; index += 1) {
        copied[index] = original[index];
    }

    copied[0] = 99;

    println("original={}", original[0]);
    println("copied={}", copied[0]);
}
```

Hasil eksekusi:

```text
original=1
copied=99
```

Jika elemennya adalah pointer, menyalinnya akan menyalin alamatnya. Itu tidak menduplikasi memori terpisah yang mereka tunjuk. Perbedaan ini penting ketika mengelola kepemilikan.

## Inisialisasi dan jangkauan efektif

Itu tidak berasumsi bahwa semua elemen array yang dideklarasikan tanpa nilai awal dapat dibaca. Jika hanya sebagian nomor yang dicatat, nomor yang diinisialisasi sebenarnya harus dikelola secara terpisah. Ini adalah alasan yang sama mengapa panjang yang dikembalikan oleh fungsi baca perpustakaan mungkin lebih kecil dari seluruh kapasitas buffer.

Karena panjang array terkandung dalam tipenya, ia tidak bertambah secara sembarangan selama eksekusi. Daftar byte yang ukurannya bertambah menggunakan penyimpanan dinamis seperti `Buffer`. Mengubah panjang array memerlukan pertimbangan jenis, nilai awal, batas atas traversal, dan perhitungan yang bergantung pada panjang tersebut.

## Latihan: Maksimum dan Lokasi

Temukan nilai maksimum dan posisi kemunculannya pertama kali dalam larik `[4, 9, 2, 9, 1]`. Jika ini adalah fungsi reguler yang semua elemennya mungkin negatif, nilai maksimum tidak boleh diinisialisasi ke 0.

### Solusi lengkap

<!-- wave-example: book-array-solution -->
```wave
fun main() {
    var values: array<i32, 5> = [4, 9, 2, 9, 1];
    var maximum: i32 = values[0];
    var position: i32 = 0;

    for (var index: i32 = 1; index < 5; index += 1) {
        if (values[index] > maximum) {
            maximum = values[index];
            position = index;
        }
    }

    println("max={} first={}", maximum, position);
}
```

Hasil eksekusi:

```text
max=9 first=1
```

Ambil elemen pertama sebagai acuan awal dan bandingkan dengan elemen kedua. Sejak `>`, posisinya tidak berubah meskipun nilai maksimum yang sama muncul kembali. Ubah menjadi `>=` untuk menjadi posisi terakhir. Antarmuka apa pun yang panjangnya nol harus memproses input kosong sebelum membaca elemen pertama.
