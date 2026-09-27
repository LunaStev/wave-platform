---
translation_set_id: stdlib-path-env
path: stdlib/path-env
locale: id
group: stdlib
group_order: 1
order: 12
title: path dan env: Jalur dan pengaturan lingkungan
summary: Membaca variabel jalur dan lingkungan di buffer pemanggil dan mengidentifikasi kesalahan kapasitas.
---

## kombinasi jalur

```text
std::path::copy
path_join2(dst: ptr<u8>, dst_cap: i32, left: str, right: str) -> i32
path_basename_copy(dst: ptr<u8>, dst_cap: i32, path: str) -> i32
path_dirname_copy(dst: ptr<u8>, dst_cap: i32, path: str) -> i32
```

Kapasitas mencakup ruang NUL terakhir. Hasil sukses adalah berapa pun panjangnya kecuali NUL, kegagalannya adalah -1. Hanya jika berhasil kita akan menggunakan tujuan sebagai string. Fungsi-fungsi ini beroperasi pada string jalur dan tidak memeriksa keberadaan file atau izin akses. Menggabungkan jalur saja tidak mencegah pelolosan direktori atau memverifikasi identitas file sebenarnya.

## variabel lingkungan

```text
std::env::environ
env_get(name: str, dst: ptr<u8>, dst_cap: i64) -> i64
env_exists(name: str) -> bool
env_get_i64(name: str) -> EnvResult<i64>
env_get_i32(name: str) -> EnvResult<i32>
```

`env_get` mengembalikan panjang tidak termasuk NUL jika berhasil. Buffer pemanggil harus mampu menampung hingga NUL. Nilai kosong berbeda dengan kesalahan tanpa kunci karena dapat berhasil dengan panjang 0.

Dapatkan dan bedakan kesalahan dari `std::env::consts` hingga NOT_FOUND, NO_SPACE, INVALID_KEY, READ, SOURCE_INCOMPLETE, NO_MEMORY. Jangan perlakukan buffer sebagai kunci yang hilang. Untuk mencari nomor, centang ok di hasilnya lalu gunakan value. Jangan secara otomatis memperlakukan konten variabel lingkungan sebagai pengaturan tepercaya; periksa cakupan dan jenisnya.

Contoh berikut menggabungkan direktori data dan nama file input.txt.

<!-- wave-example: path-api -->
```wave
import("std::path::copy")::{
    path_join2
};

fun main() -> i32 {
    var output: array<u8, 64>;
    var length: i32 = path_join2(&output[0], 64, "data", "input.txt");
    if (length < 0) {
        return 1;
    }

    println("{}", &output[0] as str);
    return 0;
}
```

Hasil eksekusi:

```text
data/input.txt
```

## Pisahkan direktori dan nama file

Contoh berikut menyalin jalur yang dibagi menjadi dua buffer. File asli tidak harus benar-benar ada.

<!-- wave-example: book-path-parts -->
```wave
import("std::path::copy")::{
    path_basename_copy,
    path_dirname_copy
};

fun main() -> i32 {
    var directory: array<u8, 64>;
    var filename: array<u8, 64>;
    var directory_length: i32 = path_dirname_copy(&directory[0], 64, "data/report.txt");
    var filename_length: i32 = path_basename_copy(&filename[0], 64, "data/report.txt");

    if (directory_length < 0 || filename_length < 0) {
        return 1;
    }

    println("directory={}", &directory[0] as str);
    println("filename={}", &filename[0] as str);
    return 0;
}
```

Hasil eksekusi:

```text
directory=data
filename=report.txt
```

Kedua buffer akan tetap berlaku hingga akhir main. `as str` membaca buffer yang sama dengan string tanpa mengalokasikan string baru. Oleh karena itu, jika Anda mengubah buffer, string yang dibaca ke alamat tersebut juga akan berubah.

## Tetapkan preferensi default

Variabel lingkungan adalah pengaturan yang diteruskan di luar program. Saat membaca pengaturan numerik, centang “Dapatkah dibaca sebagai bilangan bulat?” dan “Apakah berada dalam rentang yang diperbolehkan oleh program ini?”

<!-- wave-example: book-env-setting -->
```wave
import("std::env::environ")::{
    EnvResult,
    env_get_i32
};

fun main() -> i32 {
    var setting: EnvResult<i32> = env_get_i32("WAVE_EXAMPLE_WORKERS");
    var workers: i32 = 4;

    if (setting.ok) {
        if (setting.value < 1 || setting.value > 32) {
            println("workers must be between 1 and 32");
            return 1;
        }

        workers = setting.value;
    }

    println("workers={}", workers);
    return 0;
}
```

Jika `WAVE_EXAMPLE_WORKERS` tidak ada atau tidak dapat dibaca sebagai bilangan bulat, nilai default 4 akan digunakan. Jika bilangan bulat antara 1 dan 32 ditetapkan, nilai tersebut digunakan, dan jika bilangan bulat di luar rentang, maka akan berakhir dengan kesalahan.

Linux/macOS:

```shell
WAVE_EXAMPLE_WORKERS=8 wavec run main.wave
```

PowerShell:

```powershell
$env:WAVE_EXAMPLE_WORKERS = "8"
wavec run main.wave
```

Dalam kedua kasus tersebut, `workers=8` adalah keluaran. Contoh di atas memilih kebijakan default yang sederhana. Jika ini merupakan pengaturan yang diperlukan, perlakukan kegagalan pencarian numerik sebagai kesalahan daripada menggantinya dengan nilai default. Jika Anda perlu membedakan antara kunci yang hilang, buffer yang tidak mencukupi, dan kegagalan baca, gunakan konstanta env_get dan ENV_ERR_*.
