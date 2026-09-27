---
translation_set_id: stdlib-time
path: stdlib/time
locale: id
group: stdlib
group_order: 1
order: 11
title: time: Durasi, pengukuran, dan waktu tunggu
summary: Jelaskan perbedaan antara satuan Duration dan realtime·monotonic clock.
---

## Duration

`Duration` di `std::time::duration` memiliki seconds dan nanoseconds. Rentang nanoseconds yang dinormalisasi adalah 0 hingga 999999999. 1000 milidetik sama dengan 1 detik.

```text
time_duration_from_ms(milliseconds: i64) -> Duration
time_duration_from_seconds(seconds: i64) -> Duration
time_duration_checked_add(left: Duration, right: Duration) -> DurationResult
time_duration_to_ns(value: Duration) -> DurationValueResult
```

Operasi checked dan hasil konversi bilangan bulat meliputi ok dan value. Mengganti nilai Duration yang lebar dengan satu nilai nanodetik i64 mungkin di luar jangkauan, jadi periksalah ok terlebih dahulu.

## Bedakan antara pengukuran dan penglihatan

`std::time::clock` hingga `time_now_realtime(tp: ptr<TimeSpec>) -> i64` sesuai dengan waktu kalender. Gunakan `time_now_monotonic` untuk pengukuran waktu yang telah berlalu, karena ini dapat berubah seiring koreksi jam sistem. Penyimpanan keluaran disediakan oleh pemanggil dan hanya membaca sec/nsec jika statusnya berhasil.

```text
std::time::sleep
time_sleep(duration: Duration) -> i64
time_sleep_ms(ms: i64) -> i64
time_sleep_us(us: i64) -> i64
time_sleep_ns(ns: i64) -> i64
time_sleep_seconds(seconds: i64) -> i64
```

Penantian negatif adalah sebuah kesalahan, dan 0 adalah kesuksesan langsung. Setelah interruption, hanya sisa waktu hingga monotonic deadline yang akan ditunggu. Karena waktu bangun sebenarnya mungkin tertunda karena penjadwalan, ini tidak digunakan sebagai fungsi yang menjamin waktu eksekusi yang tepat. Memanggil sleep sinkron dalam tugas async dapat mencegah kemajuan eksekutor, jadi pilih `task::sleep_ms`.

## Contoh konversi satuan

<!-- wave-example: duration-api -->
```wave
import("std::time::duration")::{
    Duration, DurationValueResult, time_duration_from_ms, time_duration_to_ns
};

fun main() -> i32 {
    var duration: Duration = time_duration_from_ms(1500);
    var value: DurationValueResult = time_duration_to_ns(duration);
    if (!value.ok) {
        return 1;
    }

    println("{} {}", duration.seconds, duration.nanoseconds);
    println("{}", value.value);
    return 0;
}
```

Hasil eksekusi:

```text
1 500000000
1500000000
```

## Tambahkan waktu dan ubah satuan

750 ms dan 800 ms berjumlah 1 detik dan 55.0000000 nanodetik. Daripada menambahkan detik dan nanodetik secara terpisah, Anda dapat menggunakan checked_add untuk memeriksa carry dan range secara bersamaan.

<!-- wave-example: book-duration-add -->
```wave
import("std::time::duration")::{
    Duration,
    DurationResult,
    DurationValueResult,
    time_duration_from_ms,
    time_duration_checked_add,
    time_duration_to_ms
};

fun main() -> i32 {
    var first: Duration = time_duration_from_ms(750);
    var second: Duration = time_duration_from_ms(800);
    var sum: DurationResult = time_duration_checked_add(first, second);

    if (!sum.ok) {
        return 1;
    }

    var milliseconds: DurationValueResult = time_duration_to_ms(sum.value);

    if (!milliseconds.ok) {
        return 2;
    }

    println("{}s {}ns", sum.value.seconds, sum.value.nanoseconds);
    println("{}ms", milliseconds.value);
    return 0;
}
```

Hasil eksekusi:

```text
1s 550000000ns
1550ms
```

## interval waktu negatif

Duration juga mewakili angka negatif. -1ms dinormalisasi menjadi seconds=-1, nanoseconds=999000000. Gabungan kedua bidang tersebut bernilai negatif 1 milidetik. Anda tidak boleh menilai bahwa ini adalah angka positif hanya dengan melihat kolom nanoseconds.

Angka negatif valid dalam penghitungan interval waktu, tetapi meneruskan angka negatif ke sleep merupakan kesalahan. Saat menghitung sisa waktu tunggu, jika batas waktu tersebut telah terlewati maka proses selanjutnya akan dilanjutkan tanpa menunggu.

## Pilih waktu API

|tujuan|pilih|Apa arti hasilnya|
| --- | --- | --- |
|Waktu berlalu antara dua titik waktu| monotonic clock |Interval tidak bergantung pada koreksi visual sistem|
|waktu kalender sebenarnya| realtime clock |Waktu yang ditentukan oleh sistem|
|Menunggu program sinkron| `time_sleep_ms` |Alur panggilan antri|
|async Menunggu tugas| `task::sleep_ms` |Meneruskan peluang eksekusi ke tugas lain|

Bahkan jam yang mengembalikan nanodetik tidak berarti presisi pengukuran sebenarnya adalah 1 nanodetik. Saat membandingkan kinerja, ukur total waktu untuk mengulangi tugas singkat beberapa kali, dan pindahkan tugas yang tidak terkait dengan target pengukuran, seperti input/output, keluar dari bagian tersebut.
