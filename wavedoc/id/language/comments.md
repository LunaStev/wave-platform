---
translation_set_id: comments
path: language/comments
locale: id
group: language
group_order: 2
order: 15
title: Komentar
summary: Menjelaskan komentar satu baris, komentar blok bersarang, dan diagnostik komentar tidak tertutup.
---

## komentar satu baris

Konten setelah `//` adalah komentar hingga akhir baris.

```wave
var count: i32 = 10;
// 현재 요청 수
```

## blok anotasi

Proses spasi antara `/*` dan `*/` sebagai komentar blok.

```wave
/* 여러 줄에 걸친
   설명을 작성할 수 있습니다. */
```

Anda dapat menyarangkan komentar blok lainnya di dalam komentar blok.

```wave
/* 바깥 주석
   /* 안쪽 주석 */
다시 바깥 주석
*/
```

## String dan Tanda Komentar

`//`, `/*`, dan `*/` dalam literal string dan karakter adalah konten string dan tidak diperlakukan sebagai awal atau akhir komentar.

```wave
var text: str = "https://wave-lang.dev";
```

## Komentar blok yang tidak ditutup

Kegagalan menutup komentar blok dengan `*/` akan mengakibatkan diagnosis `E1002 UnterminatedComment`.

Meskipun menonaktifkan sementara blok panjang, pastikan kedalaman sarangnya benar.
