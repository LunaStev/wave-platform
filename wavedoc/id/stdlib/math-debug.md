---
translation_set_id: stdlib-math-debug
path: stdlib/math-debug
locale: id
group: stdlib
group_order: 1
order: 14
title: math dan debug: Utilitas matematika dan diagnostik
summary: Gunakan fungsi matematika dengan pemeriksaan jangkauan dan keluaran diagnostik.
---

## Fungsi integer untuk memeriksa jangkauan

```text
std::math::int
min_i32(a: i32, b: i32) -> i32
max_i32(a: i32, b: i32) -> i32
abs_i32_checked(x: i32) -> MathResult<i32>
clamp_i32_checked(x: i32, lo: i32, hi: i32) -> MathResult<i32>
div_floor_i32_checked(a: i32, b: i32) -> MathResult<i32>
div_ceil_i32_checked(a: i32, b: i32) -> MathResult<i32>
```

`MathResult<T>` berisi `value` dan `error`. Impor konstanta kesalahan dari `std::math::result`, lalu gunakan `value` hanya setelah memeriksa `error == MATH_ERROR_NONE`. Nilai absolut bilangan bulat bertanda terkecil tidak dapat direpresentasikan dengan tipe yang sama, sehingga fungsi nilai absolut dengan pemeriksaan melaporkan kegagalan. `clamp` menolak `lo > hi`. Pembagian memeriksa pembagi nol dan hasil yang melampaui rentang. `floor` dan `ceil` membulatkan secara berbeda dari pembagian bilangan bulat, yang membuang bagian pecahan menuju nol.

## Mengklasifikasikan nilai floating-point

`is_nan_f64`, `is_infinite_f64`, dan `is_finite_f64` di `std::math::float` membedakan nilai-nilai khusus. Fungsi f32 juga disediakan. `float_to_bits_f64(value: f64) -> u64` adalah fungsi untuk mendapatkan bit penyimpanan dan berbeda dari konversi numerik `value as u64`. NaN bahkan tidak sama dengan dirinya sendiri, sehingga tidak dicocokkan dengan `value == nan`.

## diagnosis

`debug_assert(condition: bool, message: str)` di `std::debug::core` berakhir setelah mendiagnosis kondisi yang salah. Situasi yang biasanya mungkin gagal, seperti masukan pengguna, ditangani dengan nilai hasil, dan assert digunakan saat memeriksa kondisi internal program yang harus dipenuhi.

Simpan program di bawah ini sebagai `main.wave` dan jalankan.

<!-- wave-example: math-api -->
```wave
import("std::math::int")::{
    min_i32, max_i32
};
import("std::debug::core")::{
    debug_assert
};

fun main() {
    var low: i32 = min_i32(9, 4);
    var high: i32 = max_i32(9, 4);
    debug_assert(low <= high, "invalid range");
    println("{} {}", low, high);
}
```

Hasil eksekusi:

```text
4 9
```

## Arah pembulatan untuk pembagian negatif

Bandingkan cara membagi -7 dengan 3. `/` terpotong menuju 0 menjadi -2. floor memilih bilangan bulat yang lebih kecil -3, dan ceil memilih bilangan bulat yang lebih besar -2.

<!-- wave-example: book-math-rounding -->
```wave
import("std::math::int")::{
    div_floor_i32_checked,
    div_ceil_i32_checked,
    abs_i32_checked
};
import("std::math::result")::{
    MathResult,
    MATH_ERROR_NONE,
    MATH_ERROR_OVERFLOW
};

fun main() -> i32 {
    var floor: MathResult<i32> = div_floor_i32_checked(-7, 3);
    var ceil: MathResult<i32> = div_ceil_i32_checked(-7, 3);

    if (floor.error != MATH_ERROR_NONE || ceil.error != MATH_ERROR_NONE) {
        return 1;
    }

    println("truncate={} floor={} ceil={}", -7 / 3, floor.value, ceil.value);

    var absolute: MathResult<i32> = abs_i32_checked(-2147483648);

    if (absolute.error == MATH_ERROR_OVERFLOW) {
        println("absolute value is outside i32");
    }

    return 0;
}
```

Hasil eksekusi:

```text
truncate=-2 floor=-3 ceil=-2
absolute value is outside i32
```

Jika Anda hanya memeriksa nilai hasil, Anda tidak dapat membedakan antara nilai pengganti yang disertakan jika terjadi kegagalan dan hasil penghitungan sebenarnya. Ikuti urutan pengecekan error terlebih dahulu. floor berguna saat memasukkan koordinat negatif ke dalam interval dengan ukuran tertentu, dan ceil berguna saat membulatkan jumlah bundel yang diperlukan.

## Kesalahan apa yang harus saya tangani?

|situasi|kesalahan|Contoh pemrosesan|
| --- | --- | --- |
|Bagilah dengan nol| `MATH_ERROR_DIVIDE_BY_ZERO` |Mengambil masukan penyebut lagi|
|Hasil tidak dapat disimpan dalam tipe| `MATH_ERROR_OVERFLOW` |Hitung dengan tipe yang lebih luas atau input tolak|
|Minimum lebih besar dari maksimum di clamp| `MATH_ERROR_INVALID_ARGUMENT` |Ubah rentang pengaturan|

assert bukan alat perbaikan kesalahan. Kesalahan dalam input pengguna ditangani menggunakan pernyataan kondisional dan nilai kembalian, dan setelah menyelesaikan perhitungan, kondisi internal yang harus dipenuhi diperiksa dengan debug_assert.
