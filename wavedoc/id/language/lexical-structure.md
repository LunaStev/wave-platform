---
translation_set_id: lexical
path: language/lexical-structure
locale: id
group: language
group_order: 2
order: 14
title: Struktur leksikal
summary: Menjelaskan pengidentifikasi, literal, pembatas, kata kunci, dan nama tipe.
---

## pengidentifikasi

Pengidentifikasi nama variabel, fungsi, tipe, dan bidang. Nama peka huruf besar-kecil dan dapat berisi kombinasi huruf, angka, dan `_` apa pun. Angka tidak dapat digunakan pada huruf pertama. Karakter Unicode juga dapat digunakan dalam pengidentifikasi.

```wave
var request_count: i64 = 0;
var 이름: str = "Wave";
```

Dalam proyek nyata, disarankan untuk menggunakan konvensi penamaan yang konsisten untuk kompatibilitas alat dan kemampuan pencarian.

## Kalimat dan Pemisah

Sebagian besar pernyataan deklarasi dan ekspresi diakhiri dengan `;`. Pernyataan dengan isi, seperti fungsi, pernyataan kondisional, pernyataan loop, dan struktur, menggunakan blok `{ ... }`.

```wave
var answer: i32 = 42;

fun double(value: i32) -> i32 {
    return value * 2;
}
```

## harafiah

```wave
var integer: i32 = 42;
var decimal: f64 = 3.14;
var text: str = "Wave";
var letter: char = 'W';
var enabled: bool = true;
var address: ptr<u8> = null;
```

Anda dapat menggunakan bilangan bulat, angka floating point, string, karakter, boolean, dan literal `null`. Gunakan `null` untuk nilai penunjuk.

## Kata Kunci dan Nama Jenis

Kata kunci utama yang digunakan dalam tata bahasa Wave adalah sebagai berikut.

`pub`, `fun`, `extern`, `export`, `type`, `enum`, `static`, `var`, `deref`, `const`, `if`, `else`, `proto`, `struct`, `while`, `for`, `in`, `out`, `clobber`, `as`, `asm`, `import`, `return`, `continue`, `print`, `input`, `println`, `match`, `break`, `true`, `false`, `null`.

Nama tipe bawaan mencakup `bool`, `char`, `byte`, `str`, tipe integer dan floating point, `ptr` dan `array`. Pointer ditulis dalam bentuk `ptr<T>`, dan array dengan panjang tetap ditulis dalam bentuk `array<T, N>`.

## String dan Karakter escape

|notasi|artinya|
| --- | --- |
| `\n` |LF Pemutusan baris|
| `\r` | CR |
| `\t` |tab|
| `\\` |Garis miring terbalik|
| `\"` |tanda kutip ganda|
| `\xNN` |Satu byte ditentukan tepat sebagai dua digit heksadesimal|

Karakter string umum disimpan sebagai UTF-8. Karena `\xNN` mempertahankan satu byte, tidak ada jaminan bahwa seluruh string adalah UTF-8 yang valid. NUL (termasuk `\x00`) di dalam string literal adalah kesalahan kompilasi. Untuk data yang berisi nol, gunakan array byte dan panjangnya.

`char` Literal harus sesuai dengan nilai 8-bit. Karakter yang melebihi rentang tersebut, seperti `'한'`, adalah kesalahan. Berbeda dengan string `"한"`.

LF, CRLF, dan CR saja di sumber masing-masing diperlakukan sebagai satu baris baru yang logis. Ini adalah aturan tentang lokasi sumber dan penghentian komentar dan tidak berarti bahwa aturan tersebut mengubah byte sebenarnya dari data file.

Nama tata bahasa tambahan mencakup `variant`, `async`, dan `await`, dan nilai asinkron dinyatakan sebagai `Future<T>`. Blok deklarasi independen `var` di atas adalah fragmen kode di dalam suatu fungsi.

[kelas string](/docs/id/language/arrays) · [Komentar](/docs/id/language/comments)

## Contoh kegagalan yang disengaja

Jika Anda menjalankan program di bawah check, kesalahan internal NUL akan terjadi. Jika Anda membutuhkan 0 byte, gunakan array byte `[97, 0, 98]`.

<!-- wave-example: reject-string-nul -->
```wave
fun main() {
    var text: str = "a\x00b";
}
```

Karakter literal di bawah ini juga melebihi rentang 8-bit, jadi ini merupakan kesalahan kompilasi. Untuk mewakili string UTF-8, gunakan `str` dan tanda kutip ganda.

<!-- wave-example: reject-wide-char -->
```wave
fun main() {
    var letter: char = '한';
}
```
