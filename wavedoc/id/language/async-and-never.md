---
translation_set_id: language-async-and-never
path: language/async-and-never
locale: id
group: language
group_order: 2
order: 13
title: 13. Fungsi asinkron dan Future
summary: Pelajari peran eksekusi tertunda Future, await, dan block_on.
---

## Mengekspresikan tugas menunggu

Untuk operasi menunggu seperti file, soket, dan pengatur waktu, perlu dibedakan antara melanjutkan perhitungan dan menunggu selesai. Fungsi asinkron menyatakan hasil yang harus diselesaikan sebagai Future. Menambahkan async tidak secara otomatis membuat thread baru atau mengubah semua panggilan sinkron menjadi asinkron.

Baca bab ini setelah Fungsi, Petunjuk, dan Penanganan Kesalahan. Contohnya adalah program asli yang menggunakan peluncur `std::task` dan menjalankan setiap file sebagai `wavec run main.wave`.

## Buat Future dan terima hasilnya

<!-- wave-example: book-async-basic -->
```wave
import("std::task" as task);

async fun calculate(value: i64) -> i64 {
    return value * 2;
}

fun main() {
    var pending: Future<i64> = calculate(21);
    var result: i64 = task::block_on(pending);

    task::shutdown();
    println("{}", result);
}
```

Hasil eksekusi:

```text
42
```

i64 yang tertulis pada deklarasi calculate adalah nilai yang diperoleh setelah selesai. Hasil panggilannya sendiri adalah Future<i64>. Dalam main normal, jalankan Future dengan block_on dan terima hasil lengkapnya.

Panggil shutdown setelah membersihkan semua tugas untuk melepaskan sumber daya pelaksana. Jangan kosongkan buffer saat tugas masih menggunakannya, dan jangan abaikan tugas yang belum selesai.

## Doa dan eksekusi tubuh berbeda

Fungsi asinkron dijalankan dengan lambat. Itu harus dibedakan dari fungsi biasa yang mengeksekusi tubuhnya segera setelah dipanggil.

<!-- wave-example: book-async-lazy -->
```wave
import("std::task" as task);

static entered: i32 = 0;

async fun work() -> i32 {
    entered += 1;
    return 7;
}

fun main() {
    var pending: Future<i32> = work();

    println("before={}", entered);

    var result: i32 = task::block_on(pending);

    task::shutdown();
    println("after={} result={}", entered, result);
}
```

Hasil eksekusi:

```text
before=0
after=1 result=7
```

Ketika Future dibuat, entered masih 0. Setelah menjalankan eksekusi, isi dieksekusi dan menjadi 1. Inilah sebabnya Anda tidak boleh menganggap tugas selesai hanya karena Future disimpan dalam variabel.

## Menunggu di dalam fungsi asinkron

Di dalam fungsi async, await menunggu penyelesaian Future lainnya. Hasil dari ekspresi await adalah nilai penyelesaian.

<!-- wave-example: book-async-nested -->
```wave
import("std::task" as task);

async fun twice(value: i32) -> i32 {
    await task::yield_now();
    return value * 2;
}

async fun process() -> i32 {
    var value: i32 = await twice(20);
    return value + 2;
}

fun main() {
    var answer: i32 = task::block_on(process());

    task::shutdown();
    println("{}", answer);
}
```

Hasil eksekusi:

```text
42
```

process menunggu Future dari twice lalu menambahkan 2 pada nilai penyelesaian. Anda dapat menggunakan nilai yang mirip dengan hasil panggilan ke fungsi biasa, namun sambil menunggu, Anda dapat meneruskan peluang eksekusi ke tugas lain.

yield_now menghasilkan peluang eksekusi kooperatif. Kegagalan untuk memberikan satu kali pun dalam loop komputasi yang panjang dapat memperlambat tugas lainnya. Asynchronous CPU bukanlah perangkat untuk memparalelkan dan mendistribusikan komputasi secara otomatis.

## Jadwalkan banyak tugas

Anda dapat menjadwalkan tugas dengan spawn dan menunggu hasilnya masing-masing. Ini adalah contoh pemeriksaan hasil akhir tanpa bergantung pada urutan keluaran antara kedua operasi.

<!-- wave-example: book-async-spawn -->
```wave
import("std::task" as task);

async fun compute(value: i32) -> i32 {
    await task::yield_now();
    return value * 2;
}

async fun combine() -> i32 {
    var first: Future<i32> = task::spawn(compute(10));
    var second: Future<i32> = task::spawn(compute(20));

    var left: i32 = await first;
    var right: i32 = await second;

    return left + right;
}

fun main() {
    var total: i32 = task::block_on(combine());

    task::shutdown();
    println("{}", total);
}
```

Hasil eksekusi:

```text
60
```

Setiap tugas yang Anda jadwalkan memiliki tempat menunggu hasilnya. Daripada membuat tugas dan melupakan penanganannya, Anda perlu memutuskan siapa yang akan memeriksa penyelesaiannya. await Jangan menganggap urutan dan urutan pelaksanaan operasi internal sama.

## Konsumsi Future sekali

Future diperlakukan sebagai pegangan konsumsi tunggal. Menyalin Future yang sama tidak membuatnya menunggu seperti dua tugas berbeda. Jangan block_on atau await lagi untuk Future yang sudah selesai.

Jika Anda memerlukan hasil yang sama di beberapa tempat, daripada menggunakan Future beberapa kali, simpan nilai yang telah selesai dan teruskan sesuai dengan aturan salin/bagikan untuk nilai tersebut. Anda juga harus memeriksa apakah ada pointer atau sumber daya yang dimiliki dalam nilai tersebut.

## Perbedaan antara pengatur waktu dan menunggu sinkron

async Saat menunggu dalam suatu fungsi, Anda dapat menggunakan `await task::sleep_ms(...)`. Memanggil sinkron sleep memblokir aliran eksekusi saat ini, yang juga dapat mempengaruhi kemajuan tugas lain di pelaksana.

Saya tidak berharap waktu tunggu persis dengan jumlah milidetik yang diminta. Tergantung pada jadwal dan tugas lainnya, Anda mungkin bangun terlambat. Saat menerapkan batas waktu, kami menggunakan clock dan deadline untuk mengukur waktu yang telah berlalu, dibandingkan menunggu waktu penuh semula setiap saat.

## Buffer Seumur Hidup dan Pembatalan

Buffer yang diteruskan ke I/O asinkron harus tetap valid meskipun fungsi pemanggilan ditangguhkan. Membebaskan atau mengalokasikannya kembali sebelum penyelesaian atau pembatalan pembersihan dapat mengakibatkan operasi dengan alamat yang tidak valid.

Tidak dapat diasumsikan bahwa permintaan pembatalan dan penyelesaian tugas terjadi pada waktu yang bersamaan. Periksa hasil seri cancel API dan tunggu penyelesaian yang diperlukan sebelum melepaskan sumber daya. Silakan baca [Lihat task](/docs/id/stdlib/task) untuk mengetahui aturan panggilan secara rinci.

## kesalahpahaman umum

|berpikir|sebenarnya memeriksa|
| --- | --- |
|async Saya menelepon dan semuanya berakhir.|Apakah Anda benar-benar menjalankan dan menyelesaikan Future?|
|async Semua panggilan dalam fungsi ini tidak sinkron|Apakah yang disebut API sinkron atau asinkron?|
|Future Salin tugas duplikat|Apakah Anda mengonsumsi pegangan yang sama berulang kali?|
|Karena Anda membatalkannya, Anda dapat segera melepaskan buffernya.|Apakah Anda sudah selesai mengatur pekerjaan setelah pembatalan?|
|Urutan keluaran antara selalu tetap|Apakah secara eksplisit hanya menunggu pesanan yang diperlukan untuk hasilnya?|

## Latihan dan solusi lengkap

Buat alur yang menunggu tiga fungsi asinkron berturut-turut. Mengembalikan hasil penggandaan dan penambahan 5.

<!-- wave-example: book-async-solution -->
```wave
import("std::task" as task);

async fun read_value() -> i32 {
    return 10;
}

async fun transform(value: i32) -> i32 {
    await task::yield_now();
    return value * 2;
}

async fun pipeline() -> i32 {
    var initial: i32 = await read_value();
    var changed: i32 = await transform(initial);

    return changed + 5;
}

fun main() {
    var result: i32 = task::block_on(pipeline());

    task::shutdown();
    println("{}", result);
}
```

Hasil eksekusi:

```text
25
```

Contoh ini sengaja merupakan ketergantungan berurutan. transform membutuhkan hasil read_value, jadi melakukan semua spawn tidak akan membuat hubungan hilang. Titik awal dari desain asinkron adalah perbedaan antara operasi independen dan operasi yang memerlukan hasil.

Setelah Anda menyelesaikan dasar-dasarnya, sambungkan ke sumber daya eksternal nyata dengan [Latihan membaca file](/docs/id/practice/file-reader) dan [TCP Latihan](/docs/id/practice/tcp-client).

## void dan never

Fungsi reguler yang menghilangkan tipe kembalian dapat kembali ke titik panggilan tanpa nilai. Tipe never ditulis sebagai `!`, artinya tidak kembali ke titik panggilan secara normal. Contoh yang representatif adalah fungsi penghentian proses.

Contoh yang mengilustrasikan deklarasi:

```text
fun log(message: str)              // 작업 후 돌아옴, 결과값 없음
fun stop(code: i32) -> !           // 정상 반환하지 않음
async fun work() -> i64            // 호출 결과는 Future<i64>
```

Tidak ada upaya untuk menjadikan never sebagai nilai tersimpan yang umum. Jangan menulis fungsi yang dideklarasikan sebagai non-return agar memiliki jalur kembali normal. Jika pembersihan sumber daya diperlukan sebelum keluar, pemanggil harus melakukannya terlebih dahulu.

## Contoh lengkap fungsi yang tidak kembali

Jika Anda menyimpannya sebagai main.wave dan menjalankannya, kode keluarnya 0 tanpa keluaran. stop tidak kembali ke pemanggil, sehingga dinyatakan sebagai `-> !`.

<!-- wave-example: never-function -->
```wave
import("std::process::core")::{
    proc_exit
};

fun stop() -> ! {
    proc_exit(0);
}

fun main() {
    stop();
}
```
