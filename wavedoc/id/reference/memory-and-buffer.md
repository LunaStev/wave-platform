---
translation_set_id: memory-buffer
path: reference/memory-and-buffer
locale: id
group: stdlib
group_order: 1
order: 4
title: mem: Alokasi, realokasi, dan tata letak
summary: Menjelaskan ukuran dalam byte, kegagalan alokasi, batas realokasi, dan pembebasan tanggung jawab.
---

## Alokasi dan dealokasi

```text
std::mem::alloc
mem_alloc(size: i64) -> ptr<u8>
mem_alloc_zeroed(size: i64) -> ptr<u8>
mem_free(p: ptr<u8>, size: i64) -> i64
mem_realloc(old_ptr: ptr<u8>, old_size: i64, new_size: i64) -> ptr<u8>
```

Ukurannya dalam byte. Alokasi ukuran 0 atau kurang akan menghasilkan null. Jika alokasi juga gagal untuk ukuran positif, mungkin null. Jangan asumsikan konten awal `mem_alloc`, tetapi gunakan `mem_alloc_zeroed` jika inisialisasi nol diperlukan.

Penelepon memiliki setiap alokasi yang berhasil dan harus meneruskan ukuran aslinya saat mengosongkannya. `mem_free(null, size)` mengembalikan 0. Pointer non-null yang dipasangkan dengan ukuran nonpositif adalah kesalahan. Jangan pernah mengakses atau mengosongkan alokasi setelah alokasi tersebut telah dikosongkan.

## Klasifikasi jika terjadi realokasi

|permintaan|tindakan|
| --- | --- |
|ukuran baru positif dan sukses|`min(old_size, new_size)` Salin byte dan batalkan alokasi byte sebelumnya|
|Alokasi baru dengan ukuran positif gagal|null Kembali, pertahankan alokasi yang ada|
|`old_ptr == null`, ukuran baru yang positif|berperilaku seperti tugas baru|
| `new_size == 0` |Upaya untuk membebaskan alokasi dan pengembalian sebelumnya yang valid null|
|old_size=0 untuk ukuran negatif atau penunjuk yang ada|null Kembali|

Hasil null dari realokasi ke ukuran nol tidak menentukan apakah pembebasan berhasil. Hubungi `mem_free` secara langsung jika Anda memerlukan statusnya. Setelah menambah alokasi, inisialisasi sendiri wilayah yang baru ditambahkan.

## Contoh melestarikan pointer yang sudah ada

Di bawah ini adalah kasus dimana ukuran barunya positif. Simpan sebagai `main.wave` dan jalankan.

<!-- wave-example: reallocation -->
```wave
import("std::mem::alloc")::{
    mem_alloc, mem_realloc, mem_free
};

fun main() -> i32 {
    var data: ptr<u8> = mem_alloc(4);
    if (data == null) {
        return 1;
    }
    deref data[0] = 7;
    var grown: ptr<u8> = mem_realloc(data, 4, 8);
    if (grown == null) {
        mem_free(data, 4);
        return 2;
    }
    data = grown;
    println("{}", deref data[0]);
    if (mem_free(data, 8) < 0) {
        return 3;
    }

    return 0;
}
```

Hasil eksekusi:

```text
7
```

Jika Anda menimpanya dengan `data = mem_realloc(...)` sebelum mengonfirmasi kegagalan, Anda mungkin kehilangan alamat yang ada. Jika realokasi berhasil, alamat lama dan petunjuk yang menunjuk ke sana tidak digunakan.

## Ukuran dan keselarasan jenis target

```text
std::mem::layout
size_of<T>() -> u64
align_of<T>() -> u64
```

Kedua nilai tersebut merupakan tata letak target kompilasi, bukan komputer yang menjalankannya. `size_of` menyertakan bantalan ekor dan tidak menghasilkan atau mengevaluasi nilai. Saat mengalikan jumlah elemen dengan ukurannya, kami memeriksa overflow. Anda dapat menggunakan `mem_size_mul_checked` dan `mem_size_add_checked` dari `std::mem::ops`.

`mem_copy` digunakan untuk menyalin rentang yang tidak tumpang tindih, dan `mem_move` digunakan untuk menyalin rentang yang mungkin tumpang tindih. Tidak ada yang bisa menentukan panjang alokasi sebenarnya hanya dari pointer, jadi pemanggil harus menjamin batasan. Daftar byte dengan berbagai ukuran dapat dikelola dengan [Buffer](/docs/id/stdlib/buffer).

## Contoh kueri tata letak

Simpan sebagai main.wave dan jalankan. Ukuran dan penyelarasan target yang tercakup dalam dokumen, i32, masing-masing 4 byte, jadi `4 4` adalah keluarannya. Periksa nilai tipe lain, terutama struktur dan pointer, berdasarkan target.

<!-- wave-example: layout-api -->
```wave
import("std::mem::layout")::{
    size_of, align_of
};

fun main() {
    println("{} {}", size_of<i32>(), align_of<i32>());
}
```

## Periksa overflow dalam perhitungan ukuran

Anda perlu memeriksa apakah `count * element_size` valid sebelum meneruskannya ke fungsi penugasan. Jika Anda mengalokasikan spasi kecil dengan nilai overflow dan menulis sebanyak angka aslinya, maka akan keluar batas.

<!-- wave-example: book-memory-size-check -->
```wave
import("std::mem::ops")::{mem_size_add_checked, mem_size_mul_checked};

fun main() {
    var result: i64 = 99;

    if (mem_size_mul_checked(3, 4, &result) == 0) {
        println("bytes={}", result);
    }

    if (mem_size_add_checked(9223372036854775807, 1, &result) < 0) {
        println("overflow rejected");
    }
}
```

Hasil eksekusi:

```text
bytes=12
overflow rejected
```

Hasil kegagalan tidak ditulis ke ukuran alokasi. Anda tidak dapat mendeteksi semua luapan hanya dengan menghitungnya dengan aritmatika reguler dan melihat apakah hasilnya negatif. Untuk menghitung ukuran yang perlu diperiksa, gunakan fungsi checked dari awal.

## Untuk salinan yang tumpang tindih, mem_move

Saat memindahkan bagian dari array yang sama ke belakang, area input dan output tumpang tindih. Jangan meneruskan rentang yang tumpang tindih ke mem_copy, gunakan mem_move.

<!-- wave-example: book-memory-overlap -->
```wave
import("std::mem::ops")::{mem_move};

fun main() {
    var data: array<u8, 5> = [1, 2, 3, 4, 5];

    mem_move(&data[1], &data[0], 4);

    for (var index: i32 = 0; index < 5; index += 1) {
        println("{}", data[index]);
    }
}
```

Hasil eksekusi:

```text
1
1
2
3
4
```

Empat byte pertama yang asli berpindah satu posisi ke kanan. Menyalinnya secara manual dapat membaca nilai yang telah ditimpa, tanpa sengaja menghasilkan semuanya. API yang tumpang tindih menangani arah penyalinan untuk Anda.

## Tulis fungsi yang melewati kepemilikan

Untuk fungsi yang mengembalikan memori, yang terbaik adalah memberikan alamat pengirim jika berhasil dan ukuran yang diperlukan untuk mengosongkannya. Jika penelepon harus menebak ukurannya, rawan salah pelepasan. Jika suatu fungsi mengembalikan alamat yang dipinjam, maka fungsi tersebut tidak boleh dibebaskan oleh pemanggil dan menjelaskan masa pakai alamat aslinya.

Melintasi batas fungsi, Anda seharusnya dapat melacak `allocator → owner → deallocation`. Nama atau tipe variabel penunjuk tidak secara otomatis menentukan kepemilikan.
