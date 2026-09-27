---
translation_set_id: practice-binary-message
path: practice/binary-message
locale: id
group: practice
group_order: 4
order: 3
title: Proyek: Membuat pesan biner
summary: Gunakan pengurutan byte eksplisit dan ULEB128 dan tolak masukan singkat.
---

## format pesan

2 byte pertama menyimpan nomor tipe big-endian u16, dan kemudian nilai ULEB128 u64. Jika Anda menulis memori struktur ke file apa adanya, itu akan dipengaruhi oleh padding dan urutan byte, jadi kodekan berdasarkan bidang.

Simpan sebagai `main.wave` dan jalankan.

<!-- wave-example: binary-message -->
```wave
import("std::bytes::types")::{
    ByteReader, ByteWriter
};
import("std::bytes::cursor")::{
    bytes_reader, bytes_writer, bytes_writer_write_be_u16, bytes_reader_read_be_u16
};
import("std::bytes::leb128")::{
    bytes_writer_write_uleb128_u64, bytes_reader_read_uleb128_u64
};

fun main() -> i32 {
    var data: array<u8, 12>;
    var writer: ByteWriter = bytes_writer(&data[0], 12);
    if (bytes_writer_write_be_u16(&writer, 7) < 0) {
        return 1;
    }
    if (bytes_writer_write_uleb128_u64(&writer, 300) < 0) {
        return 2;
    }

    var reader: ByteReader = bytes_reader(&data[0], writer.position);
    var kind: u16 = 0;
    var value: u64 = 0;
    if (bytes_reader_read_be_u16(&reader, &kind) < 0) {
        return 3;
    }
    if (bytes_reader_read_uleb128_u64(&reader, &value) < 0) {
        return 4;
    }

    println("kind={} value={} bytes={}", kind, value, reader.position);
    var short: ByteReader = bytes_reader(&data[0], 1);
    kind = 99;
    if (bytes_reader_read_be_u16(&short, &kind) >= 0) {
        return 5;
    }

    println("preserved={} {}", short.position, kind);
    return 0;
}
```

Hasil eksekusi:

```text
kind=7 value=300 bytes=4
preserved=0 99
```

## Apa yang harus diperiksa

Di reader, kami meneruskan panjang tertulis sebenarnya, bukan total kapasitas array sebesar 12. Hal ini untuk menghindari pembacaan byte tambahan yang tidak diinisialisasi sebagai input. Bahkan jika pembacaan input singkat gagal, position=0 dan kind=99 tetap dipertahankan.

## Latihan dan komentar yang diperluas

Jika format tidak mengizinkan byte tambahan di akhir, centang `reader.position == reader.len` setelah penguraian selesai. Saat menambahkan bidang panjang, pastikan ukurannya tidak lebih besar dari byte input yang tersisa, dan penghitungan panjang+offset tidak melebihi rentang.

[Lihat bytes](/docs/id/stdlib/bytes)

## Lihatlah byte sebenarnya

Tipe nomor 7 adalah big-endian u16 dan karenanya `00 07`. Nilai 300 menjadi ULEB128 hingga `AC 02`. Keseluruhan pesan berukuran empat byte berikut:

```text
00 07 AC 02
───── ─────
종류  값
```

ULEB128 menggunakan 7 bit rendah untuk setiap byte untuk nilainya, dan jika bit tinggi adalah 1, ini menunjukkan bahwa byte berikutnya mengikuti. 7 bit rendah dari 300 adalah 44 dan sisanya adalah 2. Byte pertama adalah 44 ditambah 128 tanda berturut-turut, atau 172, atau 0xAC. Tidak ada tanda kelanjutan pada byte terakhir 0x02.

## Kapasitas dan lama penggunaan

u16 menggunakan 2 byte, dan ULEB128 dari u64 menggunakan hingga 10 byte, sehingga array 12 byte dapat menyimpan kedua bidang. Nilai 300 hanya menggunakan 2 byte, sehingga pesan sebenarnya menjadi 4 byte. Saat mengirim ke file atau soket, Anda mengirim byte writer.position daripada seluruh array.

position di reader adalah posisi baca saat ini. Setelah dibaca tipenya menjadi 2, dan setelah dibaca nilainya menjadi 4. Untuk melanjutkan ke pesan berikutnya ketika terjadi kesalahan, batas pesan yang gagal harus diketahui. Kebijakan pemulihan pesan tidak secara otomatis dibuat hanya karena fungsi baca mempertahankan lokasinya.

## Latihan nilai batas

Ubah nilainya menjadi 0, 127, 128, 16383, 16384 untuk menentukan panjang pengkodean. Saat mengubah dari 127 menjadi 128, panjang ULEB128 bertambah dari 1 menjadi 2, dan ketika diubah dari 16383 menjadi 16384, panjangnya bertambah dari 2 menjadi 3. Hitung panjang total, termasuk 2 byte bidang tipe.

Pesan dengan byte terakhir terpotong juga diperiksa. Jika panjang pesan asli adalah 4, kami meneruskan panjang 3 ke reader. Bidang tipe sudah dibaca, tetapi pembacaan nilainya harus gagal karena byte terakhir ULEB128 hilang. Pada saat ini, periksa apakah posisi 2 dan nilai keluaran sebelum pembacaan ULEB128 dipertahankan.
