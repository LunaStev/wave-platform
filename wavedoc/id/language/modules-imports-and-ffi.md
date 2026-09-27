---
translation_set_id: modules-ffi
path: language/modules-imports-and-ffi
locale: id
group: language
group_order: 2
order: 11
title: 11. Modul dan kode generik
summary: Pelajari cara mendapatkan nama publik dan argumen tipe eksplisit.
---

## Alasan untuk membagi file

Seiring berkembangnya program, akan lebih mudah untuk menemukannya dengan mengelompokkan fungsi-fungsi terkait daripada menempatkan semua fungsi di main.wave. Batasan modul menentukan nama mana yang diekspos ke kode lain. Generik adalah alat yang menggunakan kembali tugas yang sama dengan tipe berbeda, terlepas dari pemisahan file.

Dalam bab ini, kita membuat program dua file dan mempelajari alias modul, pilihan import, serta fungsi dan struktur umum.

## dua file program

Buat helpers.wave dan main.wave di direktori yang sama.

helpers.wave Semua:

```wave
pub fun double(value: i32) -> i32 {
    return value * 2;
}
```

main.wave Semua:

<!-- wave-example: book-module-two-files -->
```wave
import("./helpers")::{
    double
};

fun main() {
    var result: i32 = double(21);

    println("{}", result);
}
```

Hasil eksekusi:

```text
42
```

Jalankan `wavec run main.wave` di terminal. helpers.wave juga tidak dijalankan secara terpisah. Sumber yang diperlukan terhubung melalui import.

Tanda pub di depan fungsi di helpers menunjukkan bahwa fungsi tersebut dapat diimpor oleh modul lain. Fungsi-fungsi tambahan yang tidak perlu diekspos ke dunia luar tidak perlu dipublikasikan. Bahkan jika Anda mengubah implementasi internal modul, Anda dapat mengurangi perubahan pada kode yang Anda gunakan dengan mempertahankan kontrak fungsi publik.

## Dasar untuk jalur relatif

`./helpers` relatif terhadap direktori file sumber yang membuat kalimat import. Saat menjalankan program, pisahkan dari direktori kerja tempat file I/O ditulis. Langkah mencari file import dan langkah mencari input.txt saat dijalankan berbeda-beda.

Untuk import lokal, ekstensi `.wave` dapat dihilangkan. Jika Anda membagi direktori, tulis path `./module` sesuai lokasinya. Jangan bingung antara jalur relatif lokal dengan jalur yang mengambil nama dependensi paket.

## Pilih import dan alias

Pemilihan import menyebabkan hanya nama publik yang diinginkan yang digunakan langsung di file saat ini. Jika ada konflik nama atau Anda ingin mengetahui modul mana yang memiliki fungsi tersebut, gunakan alias.

<!-- wave-example: book-module-alias -->
```wave
import("std::string::len" as strings);

fun main() {
    var length: i32 = strings::len("Wave");

    println("{}", length);
}
```

Hasil eksekusi:

```text
4
```

strings adalah alias modul yang ditentukan dalam file ini. `strings::len` menggunakan nama modul. Titik untuk akses lapangan dan `::` untuk pemisahan modul merupakan notasi yang berbeda.

Jangan gunakan opsi import dan alias import bersamaan dalam satu kalimat. Apapun gaya Anda, gunakan secara konsisten di seluruh file sehingga asal nama mudah dibaca.

## Perpustakaan dan paket standar

Jalur `std::` menunjuk ke perpustakaan standar. Pengguna menerbitkan API dan import modul yang diperlukan. Tidak semua fungsi perpustakaan standar secara otomatis ditempatkan ke dalam namespace saat ini.

Jalur paket eksternal dimulai dari nama paket. Lokasi paket disediakan oleh opsi kompiler atau manajer paket. Pertama-tama pelajari batasannya dengan modul lokal, lalu pelajari cara mengelola dependensi di [Vex Cara menggunakan](/docs/id/whale/vex-package-manager).

## Gunakan fungsi yang sama untuk setiap jenis

Fungsi berikut mengembalikan inputnya kata demi kata: Untuk i32 dan str, gunakan parameter tipe T untuk menghindari penulisan kode yang sama dua kali.

<!-- wave-example: book-generic-identity -->
```wave
fun identity<T>(value: T) -> T {
    return value;
}

fun main() {
    var number: i32 = identity<i32>(42);
    var text: str = identity<str>("Wave");

    println("{} {}", number, text);
}
```

Hasil eksekusi:

```text
42 Wave
```

T adalah tempat untuk memasukkan tipe, bukan nilai integer yang diteruskan selama eksekusi. Sebut saja dengan menentukan argumen tipe seperti `<i32>`. Fungsi umum pengguna normal tidak menghilangkan argumen tipe.

identity<str> tidak menduplikasi byte string dengan mengalokasikannya lagi. Mengembalikan nilai apa adanya. Tata bahasa generik tidak mengubah aturan penyalinan dan kepemilikan data.

## Operasi yang diperlukan oleh badan generik

Hanya karena terdapat parameter tipe tidak berarti semua operasi dapat digunakan pada semua tipe. minimum di bawah ini harus digunakan sebagai tipe sebenarnya yang dapat dibandingkan.

<!-- wave-example: book-generic-minimum -->
```wave
fun minimum<T>(left: T, right: T) -> T {
    if (left < right) {
        return left;
    }

    return right;
}

fun main() {
    var small: i32 = minimum<i32>(7, 4);
    var wide: i64 = minimum<i64>(100, 20);

    println("{} {}", small, wide);
}
```

Hasil eksekusi:

```text
4 20
```

Jika Anda mengubah argumen tipe, `<` dan pengembalian yang digunakan dalam teks harus bertipe sesuai. Saat membaca kesalahan umum, periksa kombinasi tipe yang dipanggil dan operasi yang diperlukan oleh badan fungsi.

## struktur generik

Mari buat Pair, yang mengikat dua nilai berbeda.

<!-- wave-example: book-generic-pair -->
```wave
struct Pair<A, B> {
    first: A;
    second: B;
}

fun main() {
    var item: Pair<i32, str> = Pair<i32, str> {
        first: 7,
        second: "seven"
    };

    println("{} {}", item.first, item.second);
}
```

Hasil eksekusi:

```text
7 seven
```

Pair<i32, str> dan Pair<i64, str> adalah tipe spesifik yang berbeda. Urutan argumen tipe juga bermakna. Baca deklarasi dan kode pembuatan untuk melihat di mana jenis first dan second ditentukan.

## Nama dan kontrak API dirilis

Saat Anda menerbitkan suatu fungsi, Anda tidak hanya menentukan nama tetapi juga unit input, nilai kembalian, kegagalan, dan kepemilikan. Misalnya, loop yang akan ditulis oleh pemanggil bergantung pada apakah read membaca panjang maksimum atau panjang persisnya.

pub adalah ruang lingkup publik antar modul Wave. Ini berbeda dengan export (c), yang mengekspor simbol eksternal untuk dipanggil oleh bahasa lain. Anda dapat melihat contoh lengkap menghubungkan kedua bahasa tersebut di [Lihat FFI](/docs/id/language/modules-imports-and-ffi).

## Latihan dan solusi lengkap

Buat fungsi publik square pada math.wave dan beri nama sebagai alias pada main.wave untuk mencetak pangkat 3 dan 5.

math.wave:

```wave
pub fun square(value: i32) -> i32 {
    return value * value;
}
```

main.wave:

<!-- wave-example: book-module-solution -->
```wave
import("./math" as math);

fun main() {
    println("{} {}", math::square(3), math::square(5));
}
```

Hasil eksekusi:

```text
9 25
```

Jika kesalahannya adalah fungsi tersebut hilang, periksa jalur import dan pub terlebih dahulu. Jika ada konflik nama, periksa apakah panggilan tersebut memiliki alias. Jika ini adalah kesalahan tipe, periksa input fungsi dan tipe argumen yang diteruskan. Jangan mencoba menyelesaikan masalah yang berbeda dengan satu koreksi jalur.


## C Fungsi impor

```wave
extern(c) fun puts(text: ptr<i8>) -> i32;
```

Nama ABI dapat diikuti dengan nama simbol sebenarnya sebagai string.

```wave
extern(c, "native_symbol") fun local_name(value: i32) -> i32;
```

## Wave Fungsi ekspor

```wave
export(c) fun wave_add(left: i32, right: i32) -> i32 {
    return left + right;
}
```

`extern` dan `export` dapat digunakan sebagai fungsi tunggal dan blok. Fungsi yang diekspor harus memiliki tanda tangan spesifik ABI dan oleh karena itu tidak boleh bersifat generik.

## Properti kondisi target

Properti kondisi target dapat dilampirkan ke item tingkat atas.

```wave
#[target(os="linux", arch="x86_64")]
extern(c) fun platform_call(value: i32) -> i32;
```

Kunci ketentuannya adalah `arch`, `os`, `env`, `abi` dan properti berlaku untuk item tingkat atas berikutnya.

## Hubungkan dengan fungsi C yang Anda tulis sendiri

Lab ini ditujukan untuk lingkungan asli dengan compiler C. Menggabungkan fungsi bilangan bulat tanpa alokasi perpustakaan atau pemrosesan string.

`native.c`:

```c
#include <stdint.h>
int32_t native_double(int32_t value) { return value * 2; }
```

`main.wave`:

```wave
extern(c) fun native_double(value: i32) -> i32;

fun main() {
    println("{}", native_double(21));
}
```

Jalankan Linux/macOS dari terminal di direktori kerja yang sama.

```shell
cc -c native.c -o native.o
wavec build main.wave native.o -o ffi-example
./ffi-example
```

Output yang diharapkan adalah `42`. Di shell pengembang MSVC di Windows, buat object dengan `cl /c native.c /Fonative.obj` dan sambungkan ke `wavec build main.wave native.obj -o ffi-example.exe`. Arsitektur sumber dan target object harus sama. Contohnya hanya menggunakan nilai kecil. Untuk meneruskan nilai besar ke fungsi C, rentang perkalian halaman C juga harus dijamin secara terpisah.

Jalur file lokal dimulai dengan `./`.
