---
translation_set_id: control-flow
path: language/control-flow
locale: id
group: language
group_order: 2
order: 4
title: 4. Kondisi, loop, dan nilai batas
summary: Pelajari if, for, while dan rentang pernyataan yang berulang.
---

## Pilih jalur eksekusi

Program pada bab sebelumnya mengeksekusi pernyataan dari atas ke bawah. Program nyata harus melakukan hal yang berbeda tergantung pada input dan statusnya. Pernyataan kondisional memilih jalur eksekusi, dan loop menerapkan aturan yang sama ke beberapa nilai.

Simpan setiap contoh di main.wave dan jalankan. Saat membaca kode, tuliskan di selembar kertas nilai variabel saat ini, kondisi yang akan diuji selanjutnya, dan urutan pernyataan yang akan dieksekusi. Lebih penting berlatih mengikuti arus daripada menghafal hasilnya.

## if dan else

Ketentuan ditulis dalam tanda kurung, dan teks diapit tanda kurung kurawal. Pada contoh berikut, ganti balance dengan 500 atau 2000 untuk menentukan cabang mana yang dijalankan.

<!-- wave-example: book-if-else -->
```wave
fun main() {
    var balance: i32 = 2000;
    var price: i32 = 1200;

    if (balance >= price) {
        balance -= price;
        println("bought, balance={}", balance);
    } else {
        println("not enough money");
    }
}
```

Hasil eksekusi:

```text
bought, balance=800
```

Itu tidak mengeksekusi kedua blok. Jika kondisinya benar, blok pertama akan dieksekusi. Jika kondisinya salah, blok else dijalankan. Pembelian diperbolehkan meskipun balance sama dengan price, jadi saya menulis `>=`. Jika Anda mengubahnya menjadi `>`, operasi akan berubah dengan jumlah yang sama.

## Urutan beberapa kondisi

Anda dapat menggabungkan ketentuan dengan else if. Karena kita mengeksekusi satu cabang yang terpenuhi dari atas terlebih dahulu, kita perlu mempertimbangkan apakah kita ingin memeriksa batas yang lebih besar atau batas yang lebih kecil terlebih dahulu.

<!-- wave-example: book-grade -->
```wave
fun grade(score: i32) {
    if (score < 0 || score > 100) {
        println("invalid");
    } else if (score >= 90) {
        println("A");
    } else if (score >= 80) {
        println("B");
    } else {
        println("C");
    }
}

fun main() {
    grade(95);
    grade(80);
    grade(79);
    grade(101);
}
```

Hasil eksekusi:

```text
A
B
C
invalid
```

Skor yang tidak valid ditolak terlebih dahulu dan kemudian dinilai. Jika Anda meletakkan `score >= 80` di awal, Anda tidak akan pernah mencapai cabang A karena sudutnya adalah 95 derajat. Periksa urutan ketentuan serta kebenaran setiap ketentuan.

## Jangan mengubah nilai dalam ekspresi kondisional

Operasi penugasan, penugasan gabungan, dan penambahan atau pengurangan tidak diperbolehkan dalam kondisi if, while, atau for. Gunakan `==` untuk perbandingan. Untuk memperbarui nilai dan kemudian mengujinya, tulis dua pernyataan terpisah.

Bentuk fragmen yang benar di dalam suatu fungsi:

```wave
value = read_value();

if (value == expected) {
    println("matched");
}
```

read_value dan expected tidak didefinisikan dalam fragmen ini, jadi ini bukan program lengkap yang berjalan sebagaimana adanya. Aturan yang ditampilkan di sini adalah “Bandingkan setelah perubahan status.”

## while: Selama kondisi berlaku

<!-- wave-example: book-while-countdown -->
```wave
fun main() {
    var remaining: i32 = 3;

    while (remaining > 0) {
        println("{}", remaining);
        remaining -= 1;
    }

    println("finished at {}", remaining);
}
```

Hasil eksekusi:

```text
3
2
1
finished at 0
```

Kondisi diperiksa sebelum masuk ke dalam tubuh. Jika remaining bernilai 0 dari awal, isi tidak akan pernah dieksekusi. Jika pengurangan pada akhir isi dihilangkan, kondisi tetap benar dan perulangan tidak berakhir.

Setelah menulis perulangan, centang “Apa yang membuatnya mendekati kondisi terminasi?” Jika ini adalah perulangan yang menunggu masukan, perubahan masukan atau EOF memainkan perannya, dan jika ini adalah perulangan numerik, pembaruan indeks memainkan perannya.

## for: Inisialisasi/Kondisi/Pembaruan

for mengungkapkan tiga bagian yang diperlukan untuk pengulangan.

<!-- wave-example: book-for-sum -->
```wave
fun main() {
    var total: i32 = 0;

    for (var number: i32 = 1; number <= 5; number += 1) {
        total += number;
    }

    println("sum={}", total);
}
```

Hasil eksekusi:

```text
sum=15
```

1. Inisialisasi number ke 1. Langkah ini hanya dilakukan satu kali.
2. Centang number <= 5. Jika salah, akhiri iterasi.
3. Tambahkan number ke total di teks.
4. Tingkatkan number sebanyak 1 dan kembali ke pemeriksaan kondisi.

Variabel pengulangan yang dideklarasikan dalam for tidak diasumsikan dapat digunakan setelah pengulangan. Jika desain Anda memerlukan nilai setelah iterasi, deklarasikan nilai tersebut di luar dan perjelas lokasi inisialisasi.

## Batasan Inklusi dan Eksklusi

Jumlah alami dari 1 sampai 5 adalah `<= 5`. Di sisi lain, indeks array dengan panjang 5 harus menggunakan `< 5`. Hal ini karena indeks array dimulai dari 0 dan berakhir pada 4.

Jangan bingung antara “jalankan lima kali” dengan “hingga dan termasuk nilai 5.” Anda dapat menentukan jumlah pengulangan dengan melihat nilai awal dan akhir secara bersamaan. Kasus dimana masukannya kosong dan hanya memiliki satu elemen bagus untuk mendeteksi kesalahan batas.

## continue dan break

continue melewatkan sisa iterasi ini, dan break mengakhiri iterasi terdekat.

<!-- wave-example: book-loop-control -->
```wave
fun main() {
    var total: i32 = 0;

    for (var number: i32 = 1; number <= 10; number += 1) {
        if (number == 3) {
            continue;
        }

        if (number == 6) {
            break;
        }

        total += number;
    }

    println("{}", total);
}
```

Hasil eksekusi:

```text
12
```

Jumlah totalnya adalah 1, 2, 4, dan 5. Jika continue ditemukan di for, proses perpanjangan akan dilanjutkan. Karena while tidak memiliki ekspresi pembaruan terpisah seperti for, Anda harus berhati-hati untuk tidak menghilangkan perubahan status apa pun yang diperlukan sebelum continue.

Jika loop disarangkan, satu break tidak menyelesaikan semua loop. Jika Anda perlu berhenti pada beberapa tahap, sertakan pekerjaan dalam suatu fungsi dan tunjukkan niat Anda dengan menggunakan return atau dengan memeriksa juga kondisi penghentian pada iterasi luar.

## Bagi kasusnya dengan match

Anda dapat menggunakan match saat membandingkan beberapa kasus dengan nilai yang sama. Badan setiap arm adalah sebuah blok.

<!-- wave-example: book-match-number -->
```wave
fun describe(status: i32) {
    match (status) {
        200 => {
            println("ok");
        }
        404 => {
            println("missing");
        }
        _ => {
            println("other");
        }
    }
}

fun main() {
    describe(200);
    describe(404);
    describe(500);
}
```

Hasil eksekusi:

```text
ok
missing
other
```

`_` adalah pola yang menangani sisanya. Jangan letakkan duplikat di dalam match yang sama. variant, yang memiliki tipe data berbeda bergantung pada nilainya, tercakup dalam [Bab model data](/docs/id/language/structures-enums-and-aliases).

## Contoh Lengkap : Menghitung bilangan yang memenuhi syarat

Temukan bilangan dan jumlah bilangan genap dari 1 sampai 10. Karena hitungan dan jumlah merupakan informasi yang berbeda, keduanya diakumulasikan sebagai variabel.

<!-- wave-example: book-loop-statistics -->
```wave
fun main() {
    var count: i32 = 0;
    var total: i32 = 0;

    for (var number: i32 = 1; number <= 10; number += 1) {
        if (number % 2 == 0) {
            count += 1;
            total += number;
        }
    }

    println("count={} total={}", count, total);
}
```

Hasil eksekusi:

```text
count=5 total=30
```

Bilangan genap adalah 2, 4, 6, 8, dan 10, jadi bilangan tersebut adalah 5 dan jumlahnya adalah 30. Sekalipun ekspresinya pendek, mudah untuk memverifikasi batas iterasi jika Anda terlebih dahulu memeriksa hasilnya dalam kisaran kecil yang dapat diperoleh dengan tangan.

## Latihan dan solusi lengkap

Tambahkan kelipatan 3 saja dari 1 hingga 20, tetapi jangan tambahkan nilai apa pun yang jumlahnya lebih dari 30. Kita perlu membedakan antara “menambah lalu memeriksa apakah sudah selesai” dan “memeriksa apakah sudah selesai lalu menambahkan.”

<!-- wave-example: book-loop-solution -->
```wave
fun main() {
    var total: i32 = 0;

    for (var number: i32 = 1; number <= 20; number += 1) {
        if (number % 3 != 0) {
            continue;
        }

        if (total + number > 30) {
            break;
        }

        total += number;
    }

    println("{}", total);
}
```

Hasil eksekusi:

```text
30
```

Nilai 3+6+9+12 berjumlah 30, jadi nilai berikutnya, 15, tidak ditambahkan. Masukan kecil ini aman, tetapi untuk bilangan bulat besar, cek `total + number` dapat meluap dengan sendirinya. Menulis cek tidak secara otomatis menangani setiap kasus batas.
