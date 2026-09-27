---
translation_set_id: quick-reference
path: reference/syntax-quick-reference
locale: id
group: reference
group_order: 5
order: 3
title: Referensi cepat sintaks
summary: Deklarasi yang sering digunakan, alur kontrol, tipe, pointer, dan tata bahasa FFI disusun dalam satu halaman.
---

## deklarasi

```wave
var value: i32 = 1;
var next: i32 = value + 1;
const LIMIT: i32 = 64;
static total: i64 = 0;
type Identifier = u64;
```

`var` adalah wilayahnya, `const`/`static` adalah deklarasi tingkat atas. Variabel lokal secara eksplisit mendeklarasikan tipenya.

## Fungsi

```wave
fun max(left: i32, right: i32) -> i32 {
    if (left > right) {
        return left;
    }

    return right;
}
```

## generik

```wave
fun identity<T>(value: T) -> T {
    return value;
}

var value: i32 = identity<i32>(10);
```

Saat memanggil generik, tentukan argumen tipe.

## Struktur dan enum

```wave
struct Pair {
    left: i32;
    right: i32;
}

enum Result -> i32 {
    Ok = 0,
    Error,
}
```

## Kondisi dan loop

```wave
if (ready) {
    println("ready");
}

while (count < 10) {
    count += 1;
}

for (var i: i32 = 0; i < 10; i += 1) {
    println("{}", i);
}

match (status) {
    Ready => {
        println("ready");
    }
    0 => {
        println("zero");
    }
    _ => {
        println("other");
    }
}
```

Header `if`, `while`, `for`, dan `match` menggunakan tanda kurung.

## Array dan pointer

```wave
var values: array<i32, 4> = [1, 2, 3, 4];
var p: ptr<i32> = &values[0];
var first: i32 = deref p;
```

## masukan/keluaran konsol

```wave
print("value = ");
println("{}", value);
input("{}", value);
```

Argumen pertama adalah string literal. Setiap placeholder `{}` yang tepat memerlukan ekspresi yang mengikutinya, dan target `input` harus dapat ditetapkan.

## import dan barang-barang umum

```wave
import("std::string::len");
import("./helpers" as helpers);
import("math")::{
    Vector
};

pub fun add(left: i32, right: i32) -> i32 {
    return left + right;
}

pub import("./extra"):: {
    increment
};
```

Jalur lokal dimulai dengan `./`. Alias ​​import menentukan nama modul, dan pilihan import mengimpor entri publik yang diperlukan ke dalam namespace file ini. `pub import` mengekspor ulang item yang dipilih.

## FFI

```wave
extern(c) fun native_call(value: i32) -> i32;

export(c) fun wave_call(value: i32) -> i32 {
    return value + 1;
}
```

## Targetkan item bersyarat

```wave
#[target(os="linux", arch="riscv64")]
extern(c) fun platform_call(value: i32) -> i32;
```

Kunci kondisi dukungan adalah `arch`, `os`, `env`, `abi`. Properti mengontrol item tingkat atas berikutnya.

## Perakitan sebaris

```wave
var result: i64 = 0;
asm {
    "mv a0, a1"
    in("a1") 7
    out("a0") result
    clobber("memory")
}
```

Teks instruksi dan nama register bergantung pada target. Deklarasikan semua input, output, dan clobber tersembunyi yang diperlukan untuk blok tersebut.

## pemeriksaan sumber

```shell
wavec build main.wave --emit=check
wavec print supported-targets
wavec print supported-emit-kinds
```

## Rentang Pembelajaran dan Contoh

Contoh variabel lokal dan pernyataan yang ditampilkan secara terpisah di luar fungsi adalah potongan kode yang dimasukkan ke dalam isi fungsi. Selesaikan contoh lari dan latihan ikuti di [Wave Proses Pembelajaran](/docs/id/getting-started/overview). Silakan periksa [Perpustakaan standar](/docs/id/stdlib) untuk mengetahui aturan rinci tentang memori dan fungsi eksternal.
