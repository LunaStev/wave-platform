---
translation_set_id: functions
path: language/functions-and-generics
locale: id
group: language
group_order: 2
order: 5
title: 5. Merancang dan menyusun fungsi
summary: Pelajari parameter, nilai kembalian, nilai default, dan meneruskan nilai.
---

## Mulai dari kode yang berulang

Fungsi adalah alat untuk mengurangi sintaksis, namun juga merupakan alat untuk membatasi tugas. Memisahkan apa yang diperlukan sebagai masukan, apa yang dihitung, dan hasil apa yang dikembalikan memungkinkan Anda memahami program Anda dalam bagian yang lebih kecil.

Dalam bab ini, kita mulai dengan program yang menulis perhitungan diskon beberapa kali. Setiap contoh diselesaikan di main.wave dan dijalankan sebagai `wavec run main.wave`. Kami belum membagi filenya.

<!-- wave-example: book-function-before -->
```wave
fun main() {
    var first_price: i32 = 2000;
    var first_discount: i32 = first_price * 10 / 100;
    var first_total: i32 = first_price - first_discount;
    var second_price: i32 = 5000;
    var second_discount: i32 = second_price * 10 / 100;
    var second_total: i32 = second_price - second_discount;
    println("{} {}", first_total, second_total);
}
```

Hasil eksekusi:

```text
1800 4500
```

Kedua perhitungan tersebut hanya berbeda pada harga dan memiliki struktur yang sama. Saat mengubah aturan diskon, Anda harus mengedit kedua tempat tersebut. Jika Anda mengubah satu sisi saja, Anda akan mendapatkan hasil berbeda untuk produk yang memerlukan kebijakan yang sama.

## Tentukan masukan dan keluaran

Pindahkan perhitungan berlebihan ke fungsi. Nilai yang diubah diterima sebagai masukan price, dan harga yang dihitung dikembalikan.

<!-- wave-example: book-function-extract -->
```wave
fun discounted(price: i32) -> i32 {
    var discount: i32 = price * 10 / 100;
    return price - discount;
}

fun main() {
    println("{} {}", discounted(2000), discounted(5000));
}
```

Hasil eksekusi:

```text
1800 4500
```

Nama fungsinya adalah discounted. `price: i32` dalam tanda kurung adalah deklarasi parameter dan `-> i32` adalah tipe hasil. Variabel lokal discount di badan hanya digunakan dalam fungsi ini.

`discounted(2000)` adalah ekspresi yang memanggil suatu fungsi. Angka 2000 dalam tanda kurung adalah argumen yang sebenarnya Anda sampaikan. Karena nilai yang dikembalikan oleh fungsi menjadi hasil dari ekspresi pemanggilan ini, maka dapat langsung digunakan sebagai argumen untuk println.

|terminologi|kode|artinya|
| --- | --- | --- |
|parameter| price |Nama masukan ditentukan saat mendeklarasikan fungsi|
|faktor| 2000 |Nilai diteruskan saat menelepon|
|tipe pengembalian| i32 |Jenis nilai yang dihasilkan oleh ekspresi panggilan|
|pernyataan pengembalian| return price - discount |Berikan hasilnya dan akhiri panggilan ini|

## Ikuti perintah pemanggilan dan eksekusi

Hanya menulis deklarasi fungsi di sumbernya tidak menyebabkan isi fungsi tersebut langsung dieksekusi. Dieksekusi ketika titik yang dipanggil di main tercapai.

<!-- wave-example: book-function-trace -->
```wave
fun calculate(value: i32) -> i32 {
    println("inside: {}", value);
    return value * 2;
}

fun main() {
    println("before");
    var result: i32 = calculate(7);
    println("after: {}", result);
}
```

Hasil eksekusi:

```text
before
inside: 7
after: 14
```

Urutan kemajuannya adalah keluaran pertama main, teks utama calculate, dan keluaran terakhir main. Ketika `return` dijalankan, panggilan ke calculate berakhir dengan hasil 14, dan inisialisasi result dari main selesai.

Jika beberapa pemanggilan fungsi digabungkan dalam satu ekspresi, dan urutan efek sampingnya penting, pisahkan pemanggilan tersebut menjadi pernyataan terpisah. Contoh dalam bab ini juga menyimpan hasil panggilan yang perlu ditelusuri dalam variabel lokal.

## beberapa parameter

Jika Anda juga menerima tingkat diskonto sebagai input, Anda dapat menghitung beberapa polis dengan fungsi yang sama.

<!-- wave-example: book-function-parameters -->
```wave
fun discounted(price: i32, percent: i32) -> i32 {
    return price - price * percent / 100;
}

fun main() {
    var standard: i32 = discounted(2000, 10);
    var special: i32 = discounted(2000, 25);
    println("standard={}", standard);
    println("special={}", special);
}
```

Hasil eksekusi:

```text
standard=1800
special=1500
```

Urutan argumen harus sesuai dengan deklarasi. Jika kedua parameternya adalah i32, sulit untuk membedakan maknanya melalui pemeriksaan tipe saja, meskipun urutannya diubah. Tentukan dengan jelas nama fungsi dan nama parameter, dan tulis lokasi panggilan dengan cara yang mudah dibaca.

Fungsi ini mengasumsikan jumlah kecil dan tingkat diskonto yang valid. Tidak menangani harga negatif, rasio lebih besar dari 100, atau luapan perkalian di tengah-tengah. Saat membuat suatu fungsi, Anda harus mendeskripsikan tidak hanya isi fungsi tetapi juga kondisi inputnya. Program penyelesaian selanjutnya memisahkan langkah-langkah inspeksi.

## argumen bawaan

Anda dapat memberikan nilai yang sering digunakan sebagai default.

<!-- wave-example: book-function-default -->
```wave
fun discounted(price: i32, percent: i32 = 10) -> i32 {
    return price - price * percent / 100;
}

fun main() {
    println("{}", discounted(2000));
    println("{}", discounted(2000, 25));
}
```

Hasil eksekusi:

```text
1800
1500
```

Panggilan pertama menghilangkan argumen kedua dan menggunakan 10. Panggilan kedua menggunakan 25 yang ditentukan. Nilai default dibiarkan untuk parameter opsional di akhir. Tidak digunakan sebagai tata bahasa untuk menulis spasi kosong untuk menghilangkan argumen pertama saja.

Mengubah default akan mengubah perilaku melewatkan panggilan. Nilai default fungsi publik juga merupakan bagian dari perilaku yang Anda andalkan. Inilah sebabnya mengapa panggilan dengan argumen tertentu dan panggilan dengan argumen yang dihilangkan diuji secara terpisah.

## berarti melewati nilai

Melewati nilai bilangan bulat akan memisahkan nilai yang diterima fungsi dari penyimpanan variabel pemanggil. Menghitung hasil tidak secara otomatis mengubah variabel pemanggil.

<!-- wave-example: book-function-value -->
```wave
fun next(value: i32) -> i32 {
    return value + 1;
}

fun main() {
    var count: i32 = 4;
    var later: i32 = next(count);
    println("count={} later={}", count, later);
    count = next(count);
    println("count={}", count);
}
```

Hasil eksekusi:

```text
count=4 later=5
count=5
```

Panggilan pertama berbunyi count dan menginisialisasi later. count masih 4. Setelah panggilan kedua, hasilnya ditetapkan ke count, sehingga menjadi 5. Desain pengembalian berdasarkan nilai memperjelas kepada pemanggil di mana data telah berubah.

Jika Anda ingin mengubah ruang penyimpanan asli dalam suatu fungsi, Anda dapat meneruskan sebuah pointer. Hal ini tercakup dalam [bab penunjuk](/docs/id/language/explicit-memory-type-model). Bahkan jika Anda meneruskan sebuah pointer, Anda harus membedakan antara nilai pointer itu sendiri dan ruang penyimpanan untuk alamatnya.

## Jangan tinggalkan jalan yang tidak akan kembali

Fungsi yang mengembalikan nilai harus memberikan hasil pada semua jalur yang diperlukan. Itu tidak meninggalkan jalur terakhir seperti ini:

```wave
// 의도적으로 잘못된 함수 조각

fun positive(value: i32) -> i32 {
    if (value > 0) {
        return value;
    }
}
```

Jika value kurang dari atau sama dengan 0, tidak ada nilai tetap yang dapat dikembalikan. Setelah Anda menetapkan aturan yang Anda inginkan, Anda perlu menuliskan semua rute Anda.

<!-- wave-example: book-function-paths -->
```wave
fun positive_or_zero(value: i32) -> i32 {
    if (value > 0) {
        return value;
    }

    return 0;
}

fun main() {
    println("{} {} {}", positive_or_zero(-2), positive_or_zero(0), positive_or_zero(8));
}
```

Hasil eksekusi:

```text
0 0 8
```

Panggilan yang pertama kali mengeksekusi return tidak dilanjutkan ke return di bawah. return terakhir dicapai hanya jika kondisinya salah. Kami memeriksa aturan untuk tiga kasus: positif, 0, dan negatif.

## Berfungsi tanpa hasil

Jika Anda hanya melakukan operasi seperti pencetakan, Anda dapat menghilangkan jenis pengembalian. Fungsi tanpa hasil juga dapat dihentikan lebih awal dengan `return;`.

<!-- wave-example: book-function-void -->
```wave
fun show_positive(value: i32) {
    if (value <= 0) {
        return;
    }

    println("positive={}", value);
}

fun main() {
    show_positive(-3);
    show_positive(6);
}
```

Hasil eksekusi:

```text
positive=6
```

Panggilan pertama kembali tanpa mencetak apa pun. Cetakan panggilan kedua: “Tidak ada hasil” dan “Tidak kembali ke titik panggilan” berbeda. Fungsi yang tidak kembali, seperti fungsi yang menghentikan suatu proses, diklasifikasikan berdasarkan tipe kembaliannya `!`.

## Menyusun program lengkap dengan banyak fungsi

Sekarang kami membagi validasi input, perhitungan, dan output ke dalam fungsi yang berbeda. Karena kami telah membatasi kisaran harga, perkalian antara dalam contoh ini berada dalam kisaran i32.

<!-- wave-example: book-function-program -->
```wave
fun valid_order(price: i32, quantity: i32, percent: i32) -> bool {
    return price >= 0 && price <= 100000
        && quantity >= 1 && quantity <= 100
        && percent >= 0 && percent <= 100;
}

fun discounted_unit(price: i32, percent: i32) -> i32 {
    return price - price * percent / 100;
}

fun order_total(price: i32, quantity: i32, percent: i32) -> i32 {
    var unit: i32 = discounted_unit(price, percent);
    return unit * quantity;
}

fun show_order(price: i32, quantity: i32, percent: i32) {
    if (!valid_order(price, quantity, percent)) {
        println("invalid order");
        return;
    }

    var total: i32 = order_total(price, quantity, percent);
    println("total={}", total);
}

fun main() {
    show_order(2000, 3, 10);
    show_order(2000, 0, 10);
    show_order(2000, 3, 120);
}
```

Hasil eksekusi:

```text
total=5400
invalid order
invalid order
```

Jawab satu pertanyaan untuk setiap fungsi. valid_order menentukan apakah input diperbolehkan, discounted_unit menentukan berapa besar satu diskon, order_total menentukan berapa total harga, dan show_order menentukan apa yang akan ditampilkan.

Fungsi yang lebih kecil belum tentu lebih baik. Memberi nama pada setiap ekspresi dapat mempersulit pelacakannya. Pisahkan aturan ketika aturan tersebut bermakna untuk digunakan kembali di tempat lain atau ketika ada aturan yang perlu dijelaskan dan diverifikasi secara independen.

## soal latihan

1. Tulis maximum, yang mengembalikan dua bilangan bulat yang lebih besar.
2. Tulis fungsi clamp yang mengembalikan nilai antara dua batas. Dalam solusi di bawah ini, low <= high ditetapkan sebagai kondisi pemanggilan.
3. Buat fungsi yang menambahkan pajak ke suatu nilai dan gabungkan dengan fungsi jumlah pesanan. Tentukan rentang dan titik potong bilangan bulat terlebih dahulu.

### Solusi: Fungsi dengan batas

<!-- wave-example: book-function-solution -->
```wave
fun maximum(left: i32, right: i32) -> i32 {
    if (left > right) {
        return left;
    }

    return right;
}

fun clamp(value: i32, low: i32, high: i32) -> i32 {
    if (value < low) {
        return low;
    }
    if (value > high) {
        return high;
    }

    return value;
}

fun main() {
    println("max={}", maximum(7, 4));
    println("{} {} {}", clamp(-3, 0, 10), clamp(6, 0, 10), clamp(20, 0, 10));
}
```

Hasil eksekusi:

```text
max=7
0 6 10
```

clamp memiliki tiga jalur: kurang dari jangkauan, dalam jangkauan, dan lebih besar dari jangkauan. Periksa dengan menambahkan secara manual nilai batas 0 dan 10 juga. Jika Anda ingin menerapkan hingga low > high, Anda harus memutuskan bagaimana menyatakan kegagalan. [Bab penanganan kesalahan](/docs/id/language/errors) mengatasi masalah ini menggunakan struktur hasil.


## panggilan rekursif

Suatu fungsi dapat memanggil dirinya sendiri. Rekursi memerlukan kondisi keluar di mana tidak ada lagi panggilan yang dilakukan, dan proses di mana setiap panggilan semakin mendekati kondisi tersebut.

<!-- wave-example: book-ref-recursion -->
```wave
fun factorial(value: i32) -> i32 {
    if (value <= 1) {
        return 1;
    }

    return value * factorial(value - 1);
}

fun main() {
    println("{}", factorial(5));
}
```

Hasil eksekusi:

```text
120
```

Perhitungan 5 mengarah ke `5 * factorial(4)`, yang mengembalikan 1 ketika tercapai. Nilai yang dikembalikan diteruskan secara berurutan ke panggilan sebelumnya, menghasilkan 120. Fungsi ini adalah contoh untuk menjelaskan bilangan bulat positif kecil. Dengan input yang besar, Anda perlu mempertimbangkan rentang hasil dan kedalaman panggilan. Anda dapat menghindari masalah peningkatan kedalaman panggilan dengan menulis operasi yang sama dalam satu lingkaran.

## Kesalahan umum

|fenomena|Periksa|
| --- | --- |
|Kesalahan dengan argumen yang tidak mencukupi atau terlalu banyak|Jumlah parameter dan nilai default opsional|
|Jenis pengembalian tidak cocok|return Jenis ekspresi dan deklarasi fungsi|
|Bukan kembali dari jalur tertentu|Apakah ia kembali meskipun kondisinya salah?|
|Argumen tipe generik tidak ada|Setelah nama fungsi `<Type>`|
|Dokumen asli diubah setelah fungsi penunjuk dipanggil.|Apakah ini fungsi yang hanya membaca nilai atau fungsi yang mengubahnya?|

Fungsi yang diekspor secara eksternal, seperti `export(c)`, menggunakan tanda tangan tertentu. Fungsi generik itu sendiri tidak dapat diekspor dengan konvensi pemanggilan eksternal. `ptr<T>` dan `array<T, N>` adalah tipe memori bawaan bahasa dan membedakannya dari deklarasi struktur generik pengguna.

[pembelajaran fungsi](/docs/id/language/functions-and-generics) · [Modul pembelajaran dan generik](/docs/id/language/modules-imports-and-ffi) · [FFI](/docs/id/language/modules-imports-and-ffi)
