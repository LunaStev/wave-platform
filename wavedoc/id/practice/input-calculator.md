---
translation_set_id: practice-input-calculator
path: practice/input-calculator
locale: id
group: practice
group_order: 4
order: 1
title: Proyek: Kalkulator yang memvalidasi input
summary: Kaitkan input, pemeriksaan batas, fungsi, dan kode keluar.
---

## Tujuan dan Tindakan

Masukkan kuantitas dan harga satuan dan hitung totalnya. Masukkan dua bilangan bulat yang dipisahkan oleh spasi atau pemisah baris. Contoh ini hanya menerima kuantitas 1 hingga 1000 dan harga satuan 0 hingga 100000, sehingga dihitung dalam rentang perkalian i32.

Simpan ke `main.wave`, jalankan `wavec run main.wave`, lalu masukkan `3 1200`. Hasil yang dihasilkan oleh program ini adalah sebagai berikut. Visibilitas karakter yang dimasukkan di terminal terpisah dari keluaran program.

<!-- wave-example: input-calculator -->
```wave
fun total(quantity: i32, price: i32) -> i32 {
    return quantity * price;
}

fun main() -> i32 {
    var quantity: i32 = 0;
    var price: i32 = 0;
    input("{} {}", quantity, price);
    if (quantity < 1 || quantity > 1000 || price < 0 || price > 100000) {
        println("out of range");
        return 1;
    }

    println("total={}", total(quantity, price));
    return 0;
}
```

Hasil eksekusi:

```text
total=3600
```

## Periksa juga kegagalannya

Jika Anda memasukkan `0 1200`, ia mengharapkan `out of range` dan kode keluar 1. Kegagalan penguraian numerik di `input` dan memeriksa ruang lingkup kerja program adalah dua langkah yang berbeda. Token/tipe non-numerik melebihi rentang/sebelum input yang diperlukan EOF merupakan kegagalan input. Masukan bawaan ini bukan antarmuka yang mengembalikan kesalahan dan memerlukan masukan ulang. Jika Anda memerlukan parser yang dapat dipulihkan, konfigurasikan sendiri proses verifikasi dengan membaca byte menggunakan [io](/docs/id/stdlib/files-io).

## Latihan dan komentar yang diperluas

Ambil tingkat diskonto sebagai masukan ketiga dan periksa apakah nilainya 0 hingga 100. Untuk menghindari perkalian menengah yang besar, Anda perlu memperluas rentang penghitungan dengan i64 dan memeriksa rentang tersebut saat mempersempit hasilnya. Jika Anda hanya memperluas jenis hasil terakhir, penghitungan perantara mungkin sudah dilakukan pada jenis hasil yang sempit.

[Konsol I/O](/docs/id/language/console-io-and-formatting) · [Berikutnya: Membaca file](/docs/id/practice/file-reader)

## Mengapa menentukan ruang lingkup penghitungan terlebih dahulu?

Input terbesar adalah kuantitas 1000 dan harga satuan 100000. Mengalikan kedua nilai tersebut adalah 100000000, sehingga berada dalam kisaran i32. Pengecekan rentang ini diperlukan agar hasil perkalian fungsi total dapat digunakan apa adanya.

main, yang menerima masukan, bertanggung jawab atas masukan dan pesan kesalahan, dan total hanya bertanggung jawab atas penghitungan. Bahkan jika nanti Anda mengubah perintah membaca dari suatu file, Anda masih dapat menggunakan fungsi penghitungan.

## Perhitungan diskon lengkap

Anda dapat mengalikan totalnya dengan 100 dengan menambahkan persentase diskon. Konversikan untuk melakukan penghitungan perantara sebagai i64 lalu terapkan tingkat diskonto. Karena bagian pecahan dari pembagian bilangan bulat dibuang, contoh ini memotong jumlah setelah diskon menjadi bilangan bulat.

<!-- wave-example: book-calculator-discount -->
```wave
fun discounted_total(quantity: i32, price: i32, discount: i32) -> i64 {
    var subtotal: i64 = (quantity as i64) * (price as i64);
    var remaining: i64 = 100 - (discount as i64);
    return subtotal * remaining / 100;
}

fun main() -> i32 {
    var quantity: i32 = 0;
    var price: i32 = 0;
    var discount: i32 = 0;

    input("{} {} {}", quantity, price, discount);

    if (quantity < 1 || quantity > 1000) {
        println("invalid quantity");
        return 1;
    }

    if (price < 0 || price > 100000) {
        println("invalid price");
        return 2;
    }

    if (discount < 0 || discount > 100) {
        println("invalid discount");
        return 3;
    }

    println("total={}", discounted_total(quantity, price, discount));
    return 0;
}
```

Hasil eksekusi:

```text
total=3240
```

Hasil di atas adalah ketika `3 1200 10` dimasukkan. Outputnya adalah 3240, yang merupakan pengurangan 10% dari total awal 3600.

|masukan|hasil yang diharapkan|jalur untuk memeriksa|
| --- | --- | --- |
| `3 1200 0` | `total=3600` |tidak ada diskon|
| `3 1200 100` | `total=0` |Diskon penuh|
| `3 1200 101` | `invalid discount` |Tingkat diskon melebihi kisaran|
| `1000 100000 0` | `total=100000000` |masukan maksimal|
| `0 1200 10` | `invalid quantity` |Kuantitas di bawah kisaran|

## latihan selanjutnya

Coba ubah fungsinya menjadi membulatkan angka desimal. Dalam program ini, yang hanya menerima jumlah positif, penjumlahan 50 sebelum membaginya dengan 100 akan dibulatkan menjadi bilangan bulat. Jika Anda memasukkan 99 untuk 1 dan tingkat diskonto 50, Anda dapat membandingkan hasil pemotongan 49 dan hasil pembulatan 50.
