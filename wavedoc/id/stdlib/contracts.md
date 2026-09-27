---
translation_set_id: stdlib-contracts
path: stdlib/contracts
locale: id
group: stdlib
group_order: 1
order: 2
title: Membaca dokumentasi API: Kesalahan dan kepemilikan
summary: Pahami unit argumen, struktur hasil, keberhasilan parsial, dan masa pakai sumber daya.
---

## Baca deklarasinya

Notasi berikut menjelaskan deklarasi fungsi dan bukan keseluruhan file yang dapat dieksekusi.

```text
io_read(fd: i64, buf: ptr<u8>, len: i64) -> i64
```

`fd` adalah deskriptor terbuka, `buf` adalah penyimpanan yang disediakan oleh pemanggil, dan `len` adalah jumlah byte yang tersedia untuk penulisan. `i64` tidak berarti panjang negatif itu valid. Nilai yang dikembalikan adalah jumlah byte sebenarnya yang dibaca, bukan panjang permintaan, jadi hanya rentang yang dikembalikan yang digunakan.

## Ekspresi kegagalan bervariasi dari satu fungsi ke fungsi lainnya

|cara|ya|Metode inspeksi|
| --- | --- | --- |
|Penunjuk atau null| `mem_alloc` |null Akses memori setelah pemeriksaan|
|jumlah byte atau angka negatif| `io_read` |Kesalahan negatif, 0 EOF, data positif|
|kode status| `buffer_push` |Perbandingan konstan kesalahan dengan `BUFFER_OK`|
|Sukses dan bernilai| `NetResult<T>` |Setelah memeriksa `ok`, gunakan `value`|
|Termasuk kemajuan parsial| `RandomFillResult` |Centang `ok`, `written`, `error` bersama-sama.|

Itu hanya melihat jumlah kesalahan dan tidak membandingkannya dengan konstanta di modul lain. Misalnya, nomor kesalahan env dan OS errno bukanlah sistem yang sama. Kesalahan asli di WASI tidak boleh ditafsirkan sebagai Linux errno.

## memiliki dan menyewakan

- **Milik**: Ketika memori yang dialokasikan, file terbuka, atau soket terbuka diperoleh, ia bertanggung jawab untuk memanggil rilis/penutupan yang sesuai.
- **Pinjam**: Byte view atau buffer yang diteruskan ke fungsi mengacu pada memori yang ada. Jika suatu fungsi tidak menentukan bahwa ia menerima kepemilikan, pemanggil tetap memegang kendali.
- **Argumen keluaran**: Melewati ruang penyimpanan yang valid di mana hasilnya dapat ditulis ke fungsi yang menerimanya, seperti `out_value: ptr<T>`. Pastikan kontrak menyatakan bahwa hasilnya hanya valid jika berhasil.

Menyalin struktur Buffer dapat membuat kedua salinan mengarah ke alokasi yang sama. Jangan gratiskan setiap salinan secara terpisah. Pointer yang dipinjam menjadi tidak valid setelah alokasi dibebaskan atau dialokasikan kembali. Literal string bukanlah buffer yang dapat ditulis.

## Kegagalan bukan berarti kembali ke keadaan sebelumnya

`io_write_all` mungkin gagal setelah menulis beberapa byte. Byte yang sudah ditulis secara eksternal tidak akan dikembalikan. Di sisi lain, pembacaan checked cursor dari bytes mempertahankan posisi dan nilai keluaran jika gagal. Perbedaan ini ditentukan oleh API.

Jika Anda juga ingin berlatih menangani kegagalan, lanjutkan dengan [Pembaca file](/docs/id/practice/file-reader) dan [Pesan biner](/docs/id/practice/binary-message).
