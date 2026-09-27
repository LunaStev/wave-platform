---
translation_set_id: console-io-formatting
path: language/console-io-and-formatting
locale: id
group: language
group_order: 2
order: 16
title: Input, output, dan pemformatan konsol
summary: print, println, input Menjelaskan kalimat dan aturan placeholder.
---

## pernyataan masukan/keluaran

Wave menyediakan `print`, `println`, dan `input` sebagai pernyataan input/output konsol.

```wave
fun main() {
    var count: i32 = 0;
    input("{}", count);
    print("count = ");
    println("{}", count);
}
```

Setiap kalimat diakhiri dengan `;`. Argumen pertama harus berupa string literal. Variabel atau string terhitung tidak dapat digunakan sebagai argumen format.

## pengganti

Hanya dua karakter, `{}`, yang merupakan placeholder.

```wave
println("name = {}, score = {}", name, score);
```

Jumlah placeholder dan jumlah ekspresi berikutnya harus sama persis. Jika angkanya berbeda, itu kesalahan tata bahasa.

```wave
println("{} {}", one);
// 오류: 자리표시자 2개, 값 1개
println("plain text", one);
// 오류: 자리표시자 없음, 값이 남음
```

Bentuk kurung kurawal lainnya dibiarkan sebagai teks biasa. Tidak ada placeholder yang diberi nama atau diberi nomor dalam tata bahasa ini.

## print dan println

`print` mencetak teks yang diformat apa adanya, dan `println` menambahkan jeda baris.

```wave
print("loading...");
println("done");
```

Argumen pemformatan menggunakan nilai skalar seperti bilangan bulat, angka floating point, string, dan pointer. Array dan struktur tidak dapat digunakan sebagai argumen pemformatan.

## input Sasaran

`input` menyimpan nilai baca di tujuan, jadi semua ekspresi setelah format harus berupa lokasi yang dapat ditulis.

```wave
var number: i32 = 0;
input("{}", number);
```

Variabel, bidang, dan lokasi penyimpanan yang direferensikan dapat digunakan sebagai target. Literal dan hasil perhitungan tidak dapat dijadikan masukan.

Jika semua nilai input tidak dapat dikonversi ke tipe yang diminta, program akan keluar dengan status kegagalan.

## batas waktu proses

Pernyataan ini menggunakan input dan output konsol dari lingkungan hosted. Dalam lingkungan yang berdiri sendiri, input/output yang disediakan oleh kernel atau perangkat harus didefinisikan sebagai fungsi atau batas FFI.

## Nilai dan rentang masukan

Masukan bool hanya menerima `0` dan `1`. Itu tidak menafsirkan 2 sebagai true atau menerima string `true` sebagai masukan yang sama. Input bilangan bulat harus berada dalam rentang lebar bilangan bulat target. Bilangan bulat 128, 256, 512, dan 1024-bit juga diproses berdasarkan lebar keseluruhan tipenya.

Kesalahan format, di luar jangkauan, sebelum input yang diperlukan EOF gagal. input bawaan bukanlah fungsi yang mengembalikan kegagalan dan masuk kembali, tetapi merupakan fungsi masukan yang menghentikan proses jika terjadi kegagalan. Jika Anda memerlukan pemrosesan input yang dapat dipulihkan, baca byte dengan io dan buat parser terpisah.

[Latihan Kalkulator Masukan](/docs/id/practice/input-calculator) · [Berkas dan io](/docs/id/stdlib/files-io)
