---
translation_set_id: stdlib-buffer
path: stdlib/buffer
locale: id
group: stdlib
group_order: 1
order: 5
title: buffer: Penyimpanan byte yang dapat dikembangkan
summary: Buffer Menjelaskan aturan inisialisasi, penambahan, penyelidikan, kapasitas dan rilis.
---

## Arti Buffer

`Buffer` di `std::buffer::types` memiliki `data: ptr<u8>`, `len: i64`, dan `cap: i64`. len adalah jumlah byte yang diinisialisasi dan digunakan, dan cap adalah jumlah total byte yang dialokasikan. Selalu pertahankan `0 <= len <= cap`. String NUL tidak secara otomatis menjamin penghentian.

## Dasar API

|modul|deklarasi|artinya|
| --- | --- | --- |
| `std::buffer::alloc` | `buffer_init(out_buf: ptr<Buffer>, capacity: i64) -> i64` |Inisialisasi repositori baru. Jangan mengingat ke dalam buffer yang sudah dimiliki|
|modul yang sama| `buffer_free(buf: ptr<Buffer>) -> i64` |Batalkan alokasi. Kosongkan jika berhasil|
|modul yang sama| `buffer_reserve(buf: ptr<Buffer>, required_cap: i64) -> i64` |Pastikan kapasitas penuh minimum. len Terpelihara|
|modul yang sama| `buffer_clear(buf: ptr<Buffer>) -> i64` |Pertahankan kapasitas dan len=0|
|modul yang sama| `buffer_resize(buf: ptr<Buffer>, new_len: i64, value: u8) -> i64` |Ubah panjangnya, isi byte baru dengan value|
| `std::buffer::write` | `buffer_push(buf: ptr<Buffer>, value: u8) -> i64` |tambahkan satu byte|
|modul yang sama| `buffer_append(buf: ptr<Buffer>, src: ptr<u8>, size: i64) -> i64` |sizeSalin dan tambahkan byte|
|modul yang sama| `buffer_append_str(buf: ptr<Buffer>, s: str) -> i64` |Tambahkan byte string kecuali NUL|
| `std::buffer::read` | `buffer_get(buf: ptr<Buffer>, index: i64, out_value: ptr<u8>) -> i64` |Baca satu byte dalam jangkauan|

Pengembalian status API mengembalikan `BUFFER_OK`(0) jika berhasil. Bedakan antara kesalahan INVALID, BOUNDS, OVERFLOW, dan ALLOC dari `std::buffer::error`. Angka tersebut tidak diartikan sebagai OS errno. `buffer_new` mewakili kegagalan alokasi sebagai Buffer kosong, jadi ketika Anda perlu membedakan kegagalan, gunakan `buffer_init`.

## Contoh berjalan

<!-- wave-example: buffer-api -->
```wave
import("std::buffer::alloc")::{
    Buffer, buffer_init, buffer_free
};
import("std::buffer::write")::{
    buffer_append_str, buffer_push
};
import("std::buffer::read")::{
    buffer_get
};

fun main() -> i32 {
    var data: Buffer;
    if (buffer_init(&data, 0) < 0) {
        return 1;
    }
    if (buffer_append_str(&data, "Hi") < 0 || buffer_push(&data, 33) < 0) {
        buffer_free(&data);
        return 2;
    }

    var value: u8 = 0;
    if (buffer_get(&data, 2, &value) < 0) {
        buffer_free(&data);
        return 3;
    }

    println("{} {}", data.len, value);
    if (buffer_free(&data) < 0) {
        return 4;
    }

    return 0;
}
```

Hasil eksekusi:

```text
3 33
```

Simpan sebagai `main.wave` dan jalankan sebagai `wavec run main.wave`. Kapasitas awal 0 bukanlah kegagalan, namun buffer kosong yang valid. Mengosongkan ruang selama pemrosesan lebih lanjut.

## Umur dan Kegagalan

Menumbuhkan buffer dapat mengubah data. Jangan gunakan alamat yang dipinjam sebelumnya setelah operasi yang dapat direalokasi. Menyalin struktur Buffer tidak menduplikasi alokasinya, jadi berikan alokasi tersebut kepada satu pemilik.

`buffer_get` tidak mengubah argumen keluaran jika gagal. Di sisi lain, fungsi kemudahan `buffer_at` juga mewakili kesalahan sebagai 0, jadi gunakan `buffer_get` untuk membedakan antara 0 byte aktual dan kegagalan. Hindari membuat len/cap yang tidak valid dengan langsung mengubah kolom publik.

[Memori API](/docs/id/reference/memory-and-buffer) · [Berlatih membaca file sebagai Buffer](/docs/id/practice/file-reader)

## Amati panjang dan kapasitas secara terpisah

reserve mengosongkan ruang penyimpanan, tetapi tidak menambah len. resize mengubah panjang sebenarnya yang digunakan dan menginisialisasi bagian yang diperluas ke byte yang ditentukan. clear hanya menyetel panjang yang digunakan ke 0, sehingga alokasi dapat digunakan kembali.

Simpan program berikut sebagai main.wave dan jalankan. Hal ini tidak bergantung pada pertumbuhan kapasitas yang berlipat ganda; itu hanya memastikan bahwa Anda memiliki ruang yang Anda butuhkan.

<!-- wave-example: book-buffer-length-capacity -->
```wave
import("std::buffer::alloc")::{
    Buffer,
    buffer_init,
    buffer_reserve,
    buffer_resize,
    buffer_clear,
    buffer_free
};

fun main() -> i32 {
    var data: Buffer;

    if (buffer_init(&data, 0) < 0) {
        return 1;
    }

    if (buffer_reserve(&data, 20) < 0) {
        buffer_free(&data);
        return 2;
    }

    println("reserved length={}", data.len);

    if (buffer_resize(&data, 3, 7) < 0) {
        buffer_free(&data);
        return 3;
    }

    println("resized length={} first={}", data.len, deref data.data[0]);

    if (buffer_clear(&data) < 0) {
        buffer_free(&data);
        return 4;
    }

    println("cleared length={}", data.len);

    if (buffer_free(&data) < 0) {
        return 5;
    }

    return 0;
}
```

Hasil eksekusi:

```text
reserved length=0
resized length=3 first=7
cleared length=0
```

Ketika meningkat menjadi resize, kita meneruskan value=7, jadi ketiga byte baru yang kita lihat adalah 7. Spasi yang hanya diamankan dengan reserve tidak dibaca sebagai data yang diinisialisasi. cap akan tetap ada setelah clear dan dapat ditambahkan lagi ke Buffer yang sama.

## Bedakan antara nol byte dan kegagalan pencarian

buffer_get mengembalikan status dan menulis byte sebenarnya sebagai argumen keluaran. Sekalipun datanya 0, itu adalah keberhasilan yang normal.

<!-- wave-example: book-buffer-zero-bounds -->
```wave
import("std::buffer::alloc")::{Buffer, buffer_init, buffer_free};
import("std::buffer::write")::{buffer_push};
import("std::buffer::read")::{buffer_get};
import("std::buffer::error")::{BUFFER_ERR_BOUNDS};

fun main() -> i32 {
    var data: Buffer;

    if (buffer_init(&data, 0) < 0) {
        return 1;
    }

    if (buffer_push(&data, 0) < 0) {
        buffer_free(&data);
        return 2;
    }

    var value: u8 = 99;

    if (buffer_get(&data, 0, &value) < 0) {
        buffer_free(&data);
        return 3;
    }

    println("stored={}", value);
    value = 99;

    if (buffer_get(&data, 1, &value) == BUFFER_ERR_BOUNDS) {
        println("outside, preserved={}", value);
    }

    if (buffer_free(&data) < 0) {
        return 4;
    }

    return 0;
}
```

Hasil eksekusi:

```text
stored=0
outside, preserved=99
```

Pukulan pertama adalah keberhasilan yang bernilai 0, pukulan kedua adalah kegagalan di luar batas. Bahkan jika value=99 tetap ada setelah kegagalan, itu tidak berarti bahwa itu adalah nilai yang dibaca dari buffer. Pastikan untuk memeriksa statusnya bersama-sama.

## Solusi latihan: Akumulasi byte

Untuk menjumlahkan angka 0 hingga 9, ulangi buffer_push dan periksa setiap hasilnya. Simpan jumlahnya di i64 dan baca hanya rentangnya `0 <= index < data.len`. Setelah menangani buffer, kami memanggil buffer_free pada jalur berhasil dan gagal.

<!-- wave-example: book-buffer-sum -->
```wave
import("std::buffer::alloc")::{Buffer, buffer_init, buffer_free};
import("std::buffer::write")::{buffer_push};

fun main() -> i32 {
    var data: Buffer;

    if (buffer_init(&data, 0) < 0) {
        return 1;
    }

    for (var value: i32 = 0; value < 10; value += 1) {
        if (buffer_push(&data, value as u8) < 0) {
            buffer_free(&data);
            return 2;
        }
    }

    var total: i64 = 0;

    for (var index: i64 = 0; index < data.len; index += 1) {
        total += deref data.data[index] as i64;
    }

    println("sum={}", total);

    if (buffer_free(&data) < 0) {
        return 3;
    }

    return 0;
}
```

Hasil eksekusi:

```text
sum=45
```
