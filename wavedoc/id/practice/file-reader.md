---
translation_set_id: practice-file-reader
path: practice/file-reader
locale: id
group: practice
group_order: 4
order: 2
title: Proyek: Membaca file dan menghitung byte
summary: Buffer, file I/O, penanganan kegagalan tautan dan pembebasan memori.
---

## siap

Buat input.txt di direktori kerja yang berisi `Wave` diikuti dengan satu baris baru LF. File tersebut kemudian berisi 5 byte. Dengan CRLF berisi 6 byte; BOM UTF-8 menambahkan byte lebih lanjut. Periksa pengkodean file editor dan akhir baris.

Simpan program sebagai `main.wave` dan jalankan sebagai `wavec run main.wave` dari direktori yang sama. Jalur relatif bersifat relatif terhadap direktori kerja yang sedang berjalan.

<!-- wave-example: file-reader -->
```wave
import("std::buffer::alloc")::{
    Buffer, buffer_init, buffer_free
};
import("std::fs::file")::{
    read_to_end
};

fun main() -> i32 {
    var data: Buffer;
    if (buffer_init(&data, 0) < 0) {
        return 1;
    }

    var count: i64 = read_to_end("input.txt", &data);
    if (count < 0) {
        println("read failed={}", count);
        buffer_free(&data);
        return 2;
    }

    var lines: i64 = 0;
    for (var i: i64 = 0; i < data.len; i += 1) {
        if (deref data.data[i] == 10) {
            lines += 1;
        }
    }

    println("bytes={} LF={}", count, lines);
    if (buffer_free(&data) < 0) {
        return 3;
    }

    return 0;
}
```

Hasil eksekusi:

```text
bytes=5 LF=1
```

## Tindakan dan Tanggung Jawab

`read_to_end` membuka dan menutup file, tetapi rilis Buffer adalah tanggung jawab penelepon. Menonaktifkan jalur keberhasilan dan kegagalan baca. Karena ini adalah data dengan jumlah byte, kami tidak berasumsi bahwa ini adalah string yang diakhiri dengan NUL.

LF Jumlah garis dan jumlah garis yang dikira orang tidak selalu sama. Jika baris terakhir tidak berisi LF, maka tidak termasuk dalam hitungan LF untuk program ini. Periksa juga kegagalan baca dengan mengganti nama file. Nomor kesalahan tertentu mungkin berbeda-beda, bergantung pada lingkungan Anda.

## Latihan dan komentar yang diperluas

Jika Anda ingin menangani file yang sangat besar, gantilah dengan array berukuran tetap dan iterasi `io_read`. Hanya memproses rentang pengembalian positif dan berakhir pada 0. Anda dapat mengumpulkan jumlah byte dan jumlah LF tanpa harus menyimpan seluruh file di memori. Jika Anda membukanya sendiri, itu juga akan menutup deskriptor di jalur keluar mana pun.

[Lihat fs dan io](/docs/id/stdlib/files-io) · [Lihat Buffer](/docs/id/stdlib/buffer)

## Ikuti alur pemrosesan

1. Inisialisasi bin Buffer. Belum ada konten file.
2. read_to_end membaca file dan menambah ruang yang diperlukan.
3. Jika pembacaan berhasil, byte dalam kisaran data.len dicentang.
4. Setiap kali kita menemukan nilai byte 10 dari LF, kita menambah lines.
5. Cetak hasilnya dan lepaskan Buffer.

data.cap adalah ruang penyimpanan yang dicadangkan dan data.len adalah panjang data yang valid. Jika Anda mengubah kondisi pengulangan menjadi cap, byte yang tidak ada dalam file akan dibaca, jadi gunakan len. count adalah jumlah byte yang ditambahkan dalam panggilan ini ke read_to_end. Contoh ini dimulai dengan Buffer kosong, jadi count dan data.len sama.

## Ubah masukan untuk memeriksa

|input.txt Isi| bytes | LF |alasan|
| --- | --- | --- | --- |
|berkas kosong| 0 | 0 |Tidak ada byte untuk dibaca|
| `Wave` | 4 | 0 |Tidak ada jeda baris terakhir|
| `Wave` + LF | 5 | 1 |Data hingga LF terakhir|
| `A` + LF + `B` + LF | 4 | 2 |Hitung dua LF|
| `Wave` + CRLF | 6 | 1 |CR juga 1 byte, tetapi hanya LF yang dihitung.|

Untuk menghitung file tanpa LF di baris terakhir sebagai satu baris, tambahkan 1 ke jumlah baris jika file tidak kosong dan byte terakhir bukan 10. Anda harus memeriksa terlebih dahulu apakah data.len adalah 0 sebelum Anda dapat mengakses elemen terakhir.

## Perluas ke file besar

Metode pengarsipan seluruh konten saat ini memudahkan untuk membaca ulang atau mengambil data nanti. Jika Anda hanya membutuhkan jumlah byte dan jumlah LF, menggunakan kembali buffer berukuran tetap adalah hal yang masuk akal.

Dalam perulangan io_read pada contoh [Membaca file ke dalam buffer ukuran tetap](/docs/id/stdlib/files-io), hitung saja LF berdasarkan jumlah byte yang dikembalikan. Itu dapat diproses dengan memori yang sama dengan ukuran buffer, bukan seluruh panjang file. Jika pembacaan menghasilkan 0, maka selesai; jika negatif, itu kesalahan.
