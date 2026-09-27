---
translation_set_id: stdlib-task
path: stdlib/task
locale: id
group: stdlib
group_order: 1
order: 13
title: task: Menjalankan Future dan membersihkan sumber daya
summary: Menjelaskan konsumsi tunggal, eksekusi, dan pembatalan setelah pembersihan operasi asinkron.
---

## Dasar API

Diimpor sebagai `import("std::task" as task);`.

```text
block_on<T>(future: Future<T>) -> T
spawn<T>(future: Future<T>) -> Future<T>
cancel<T>(future: Future<T>) -> bool
yield_now() -> Future<void>
sleep_ms(milliseconds: i64) -> Future<void>
cancel_and_join<T>(future: Future<T>) -> Future<void>
shutdown()
```

`task::block_on(pending)` berjalan hingga Future selesai dan mengembalikan hasilnya. Jenis hasil ditentukan dari Future yang dilewati.

## aturan hidup

Future adalah pegangan konsumsi satu kali. Saya tidak menganggap menyalin nilai sebagai dua operasi independen. Anda juga bertanggung jawab untuk menunggu hasil atau membatalkan/mengatur tugas yang dijadwalkan dengan `spawn`. Future sudah dikonsumsi dengan `await` dan `block_on` tidak akan dikonsumsi lagi.

membatalkan permintaan pembatalan; hal ini tidak menjamin bahwa pembersihan telah selesai. Tunggu penyelesaian yang diperlukan sebelum melepaskan sumber daya. Jangan mengosongkan memori yang dipinjam oleh I/O asinkron saat tugas masih dapat mengaksesnya. Hubungi `shutdown` setelah tugas selesai atau selesai pembersihan pembatalan.

## pelaksanaan kolaboratif

Komputasi yang panjang dan panggilan blocking yang sinkron dapat memperlambat kemajuan keseluruhan pelaksana. Penantian asinkron dengan yield menghasilkan peluang eksekusi. Pemblokiran I/O tidak menjadi asinkron hanya karena berada di dalam fungsi async.

Lihat urutan eksekusi dan pembersihan dalam program lengkap di [Pengantar kode asinkron](/docs/id/language/async-and-never).

## Hasilkan dan tunggu sampai selesai

Program berikut menghasilkan eksekusi di tengah operasi dan mengembalikan hasilnya. Karena yield bukan penghentian fungsi, kode setelah await dijalankan terus menerus.

<!-- wave-example: book-task-yield -->
```wave
import("std::task" as task);

async fun calculate() -> i32 {
    println("started");
    await task::yield_now();
    println("resumed");
    return 42;
}

fun main() -> i32 {
    var pending: Future<i32> = calculate();
    var result: i32 = task::block_on(pending);

    println("result={}", result);
    task::shutdown();
    return 0;
}
```

Hasil eksekusi:

```text
started
resumed
result=42
```

Tidak ada operasi lain dalam contoh ini, sehingga urutan keluarannya konsisten. Dalam program yang memiliki banyak tugas spawn, tugas lain dapat dilanjutkan pada titik yield, sehingga tidak bergantung pada urutan keluaran tugas yang berbeda.

## Urutan pengorganisasian tugas

1. Siapkan ruang penyimpanan dan sumber daya yang Anda perlukan untuk pekerjaan Anda.
2. Buat Future dan jalankan sebagai await, block_on, atau spawn.
3. Jika Anda membutuhkan hasil, tunggu hingga selesai.
4. Jika Anda membatalkan tugas yang sedang berjalan, tunggu hingga tugas tersebut selesai.
5. Membersihkan buffer, file, dan soket yang dipinjam oleh tugas.
6. Jika tidak ada pekerjaan tersisa, hubungi shutdown.

Future Meninggalkan cakupan suatu variabel adalah satu hal, dan membersihkan pekerjaan dengan aman adalah hal lain. Khususnya, jika Anda meneruskan alamat array fungsi-lokal ke tugas async, tugas tersebut harus selesai menggunakan array tersebut sebelum fungsi kembali.

## async Pembagian fungsi dan fungsi umum

Perhitungan murni dapat dipisahkan menjadi fungsi reguler. Lampirkan async ke fungsi yang perlu menyatakan menunggu, dan tunggu hingga selesai dengan await di dalam fungsi tersebut. Hanya membungkus fungsi sinkron yang membutuhkan waktu lama, seperti membaca file, dengan fungsi async tidak memberikan kesempatan pada tugas lain untuk dijalankan.

Anda dapat membandingkan saat membuat dan menjalankan Future di [pembelajaran asinkron](/docs/id/language/async-and-never).
