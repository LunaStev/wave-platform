---
translation_set_id: stdlib-bytes
path: stdlib/bytes
locale: id
group: stdlib
group_order: 1
order: 6
title: bytes: Rentang, kursor, dan ULEB128
summary: Menjelaskan byte baca/tulis dengan panjang view dan mempertahankan status jika terjadi kegagalan.
---

## Apa bedanya dengan string?

Data byte bisa berisi nol, jadi berikan pointer bersama dengan panjangnya. `Bytes` dan `BytesMut` merupakan tampilan yang tidak dimiliki dan hanya valid selama penyimpanan yang mendasarinya tetap valid. `BytesMut` memerlukan penyimpanan yang dapat ditulisi.

## Buat cursor

```text
std::bytes::cursor
bytes_reader(data: ptr<u8>, len: i64) -> ByteReader
bytes_writer(data: ptr<u8>, len: i64) -> ByteWriter
bytes_reader_read_be_u16(reader: ptr<ByteReader>, out_value: ptr<u16>) -> i32
bytes_writer_write_be_u16(writer: ptr<ByteWriter>, value: u16) -> i32
```

Impor `ByteReader` dan `ByteWriter` dari `std::bytes::types`. Bidang posisinya mengidentifikasi lokasi operasi berikutnya. Panjangnya adalah jumlah total byte yang dapat diakses. Membangun salah satu kursor tidak menyalin atau mengalokasikan memori yang mendasarinya.

`be` adalah big-endian, `le` adalah little-endian. Jika jenis file adalah big-endian, gunakan fungsi `be` terlepas dari urutan byte host CPU. Ada fungsi baca/tulis 16, 32, dan 64 bit signed/unsigned dan satu byte.

## Kesalahan dan Pelestarian Negara

`BYTES_OK` dari `std::bytes::errors` adalah 0. INVALID menunjukkan rentang yang tidak valid, input EOF tidak mencukupi, kapasitas output NO_SPACE tidak mencukupi, dan OVERFLOW merupakan nilai di luar rentang yang dapat diwakili. Memeriksa posisi gerak maju kursor hanya setelah seluruh operasi berhasil. Pembacaan yang gagal juga mempertahankan nilai keluaran.

## ULEB128

```text
std::bytes::leb128
bytes_reader_read_uleb128_u64(reader: ptr<ByteReader>, output_value: ptr<u64>) -> i32
bytes_writer_write_uleb128_u64(writer: ptr<ByteWriter>, value: u64) -> i32
```

ULEB128 menyimpan bilangan bulat 64-bit yang tidak ditandatangani dalam jumlah byte yang bervariasi, menggunakan paling banyak 10 byte. Jika ruang tidak mencukupi, penulis mempertahankan posisinya dan byte tujuan. Pembaca membedakan masukan yang tidak lengkap dari nilai yang melebihi u64. Dihentikan, pengkodean non-minimal diterima.

Untuk membuat pesan aktual dan melihat kegagalan input singkat, lanjutkan ke [Latihan pesan biner](/docs/id/practice/binary-message). Jangan mencoba mengeluarkan string byte yang berisi angka nol sebagai `str`.

## Membaca byte yang sama dalam urutan berbeda

Urutan byte adalah aturan penyimpanan angka. Jika Anda membaca dua byte 1 dan 2 sebagai big-endian, maka hasilnya adalah 1×256+2, dan jika Anda membacanya sebagai little-endian, maka hasilnya adalah 2×256+1. Pilih berdasarkan aturan jaringan atau jenis file.

<!-- wave-example: book-bytes-endian -->
```wave
import("std::bytes::types")::{Bytes, bytes_view};
import("std::bytes::read")::{bytes_read_be_u16, bytes_read_le_u16};

fun main() -> i32 {
    var data: array<u8, 2> = [1, 2];
    var view: Bytes = bytes_view(&data[0], 2);
    var big: u16 = 0;
    var little: u16 = 0;

    if (bytes_read_be_u16(view, 0, &big) < 0) {
        return 1;
    }

    if (bytes_read_le_u16(view, 0, &little) < 0) {
        return 2;
    }

    println("be={} le={}", big, little);
    return 0;
}
```

Hasil eksekusi:

```text
be=258 le=513
```

view meminjam array. Karena tidak ada alokasi atau penyalinan terpisah, array hanya dapat digunakan selama valid. offset dalam satuan byte dan membaca 16 bit memerlukan 2 byte dari posisi tersebut.

## Pilih offset API dan cursor API

Fungsi read/write, yang menggunakan argumen offset, cocok untuk format yang secara langsung membaca posisi bidang tertentu. Untuk aliran yang posisi berikutnya bergantung pada panjang bidang sebelumnya, cursor dengan position akan lebih mudah digunakan.

Saat mencampur keduanya, jelaskan mana standarnya: cursor.position atau pisahkan offset. Hindari kesalahan menambahkan lokasi yang sama dua kali atau berpindah ke lokasi berikutnya tanpa keberhasilan pembacaan.

## Periksa status dari input singkat

<!-- wave-example: book-bytes-short-read -->
```wave
import("std::bytes::types")::{ByteReader};
import("std::bytes::cursor")::{bytes_reader, bytes_reader_read_be_u16};
import("std::bytes::errors")::{BYTES_ERROR_EOF};

fun main() {
    var data: array<u8, 1> = [1];
    var reader: ByteReader = bytes_reader(&data[0], 1);
    var value: u16 = 99;
    var status: i32 = bytes_reader_read_be_u16(&reader, &value);

    if (status == BYTES_ERROR_EOF) {
        println("position={} value={}", reader.position, value);
    }
}
```

Hasil eksekusi:

```text
position=0 value=99
```

Tidak diperlukan dua byte, sehingga nilai keluaran dan lokasi dipertahankan. Properti ini berguna dalam desain yang mencoba membaca kembali bidang yang sama setelah memperoleh lebih banyak masukan. Namun, jika ruang penyimpanan yang ditunjukkan oleh view telah dialokasikan kembali, alamatnya juga harus diperbarui.

## Batas ULEB128

0~127 menggunakan satu byte, 128 dan seterusnya menggunakan lebih banyak byte. Bit tingkat tinggi setiap byte menunjukkan apakah data mengikuti. Nilai di luar u64 atau masukan kontinu yang terlalu panjang adalah OVERFLOW, yang berbeda dengan EOF yang hanya memiliki masukan lebih sedikit.

Periksa langsung apakah Kapasitas Tidak Memadai writer mempertahankan statusnya.

<!-- wave-example: book-leb128-space -->
```wave
import("std::bytes::types")::{ByteWriter};
import("std::bytes::cursor")::{bytes_writer};
import("std::bytes::leb128")::{bytes_writer_write_uleb128_u64};
import("std::bytes::errors")::{BYTES_ERROR_NO_SPACE};

fun main() {
    var data: array<u8, 1> = [85];
    var writer: ByteWriter = bytes_writer(&data[0], 1);
    var status: i32 = bytes_writer_write_uleb128_u64(&writer, 128);

    if (status == BYTES_ERROR_NO_SPACE) {
        println("position={} byte={}", writer.position, data[0]);
    }
}
```

Hasil eksekusi:

```text
position=0 byte=85
```

128 membutuhkan dua byte tetapi hanya satu spasi. Setelah kegagalan, byte pertama 85 juga tetap ada. Memeriksa jenis kesalahan dan pelestarian status bersama-sama menggambarkan batas lebih baik daripada pemeriksaan keberhasilan roundtrip yang sederhana.

## Urutan pembuatan parser pesan

1. Baca header tetap dan periksa jenis dan versinya.
2. Baca panjangnya dan bandingkan dengan rentang masukan yang tersisa.
3. Hanya meneruskan data yang diperlukan ke view atau buffer terpisah.
4. Jika format memerlukan keseluruhan pesan, byte tambahan juga akan diperiksa.
5. Membedakan antara EOF dan kesalahan format yang tidak valid dan meneruskannya ke pemanggil.

Anda dapat membuat satu program dengan menghubungkan kolom di [Latihan pesan biner](/docs/id/practice/binary-message).
