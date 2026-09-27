---
translation_set_id: learn-allocation
path: language/allocation
locale: id
group: language
group_order: 2
order: 10
title: 10. Alokasi memori dan manajemen sumber daya
summary: Pelajari tentang kegagalan alokasi, inisialisasi, cakupan, dan pembebasan.
---

## Kapan Anda memerlukan ruang penyimpanan dinamis?

Array berukuran tetap menyertakan panjangnya dalam tipenya. Gunakan memori dinamis ketika jumlah data hanya diketahui saat runtime, seperti ukuran file atau panjang input. Bebaskan setiap alokasi saat Anda tidak lagi membutuhkannya.

Dalam bab ini, Anda akan mengelola alokasi kecil, mengubah ukurannya, dan kemudian menggunakan Buffer. Melewati pointer berbeda dengan mentransfer kepemilikan. Identifikasi sumber daya mana yang dimiliki setiap fungsi saat Anda mengikuti contoh.

## Alokasikan, periksa, gunakan, dan gratis

<!-- wave-example: book-alloc-lifecycle -->
```wave
import("std::mem::alloc")::{
    mem_alloc_zeroed, mem_free
};

fun main() -> i32 {
    var size: i64 = 4;
    var data: ptr<u8> = mem_alloc_zeroed(size);

    if (data == null) {
        println("allocation failed");
        return 1;
    }

    deref data[0] = 42;
    println("{} {}", deref data[0], deref data[1]);

    var status: i64 = mem_free(data, size);

    if (status < 0) {
        return 2;
    }

    return 0;
}
```

Hasil eksekusi:

```text
42 0
```

Program ini memiliki empat langkah: meminta 4 byte, memeriksa null, hanya mengakses rentang yang valid, dan mengosongkan alokasi. Jika berhasil, mem_alloc_zeroed menginisialisasi memori ke nol, sehingga byte kedua adalah nol meskipun program belum menulis ke dalamnya.

Jangan berasumsi ada konten awal untuk memori yang dikembalikan oleh mem_alloc. Inisialisasi setiap wilayah sebelum membacanya. Ukuran alokasi nol atau negatif menghasilkan null. Alokasi ukuran positif juga bisa gagal memperoleh memori.

## satuan ukuran

Argumen ukuran API alokasi memori diukur dalam byte. Untuk mengalokasikan sepuluh bilangan bulat, kalikan ukuran elemen dengan jumlah elemen. Pastikan perkalian ini tidak meluap.

<!-- wave-example: book-alloc-typed -->
```wave
import("std::mem::alloc")::{
    mem_alloc, mem_free
};
import("std::mem::layout")::{
    size_of
};
import("std::mem::ops")::{
    mem_size_mul_checked
};

fun main() -> i32 {
    var count: i64 = 3;
    var bytes: i64 = 0;
    var item_size: i64 = size_of<i32>() as i64;

    if (mem_size_mul_checked(count, item_size, &bytes) < 0) {
        return 1;
    }

    var raw: ptr<u8> = mem_alloc(bytes);

    if (raw == null) {
        return 2;
    }

    var values: ptr<i32> = raw as ptr<i32>;

    for (var index: i64 = 0; index < count; index += 1) {
        deref values[index] = (index + 1) as i32;
    }

    println("{} {} {}", deref values[0], deref values[1], deref values[2]);

    if (mem_free(raw, bytes) < 0) {
        return 3;
    }

    return 0;
}
```

Hasil eksekusi:

```text
1 2 3
```

count adalah jumlah elemen; byte adalah jumlah byte. Aritmatika penunjuk bergerak dalam satuan i32, tetapi membebaskan alokasi memerlukan ukuran aslinya dalam byte. size_of menggunakan tata letak tipe target, membuat hubungan dengan tipe elemen menjadi eksplisit.

Saat mengubah hasil size_of menjadi i64 untuk tipe umum yang sangat besar, rentang konversi juga harus dipertimbangkan. Di sini kita menggunakan i32, yang ukurannya diketahui.

## Bersihkan bahkan di jalur kegagalan

Jika operasi lain gagal setelah alokasi, kosongkan memori sebelum kembali lebih awal. Tabel kepemilikan membantu mengidentifikasi jalur pembersihan yang mungkin Anda lewatkan.

|Langkah|Sumber daya yang dimiliki oleh|Jika Anda gagal|
| --- | --- | --- |
|Sebelum alokasi|Tidak ada|segera kembalikan|
|Setelah alokasi berhasil|data dan ukuran asli|data Kembali setelah rilis|
|Setelah realokasi berhasil|Alamat baru dan ukuran baru|Rilis alamat baru|
|Setelah rilis|Tidak ada|Jangan gunakan alamat lama|

Menimpa variabel pointer dan kehilangan alamat asli juga menghilangkan informasi yang diperlukan untuk mengosongkan alokasi. Hal ini menyebabkan kebocoran memori. Sebaliknya, pelepasan alokasi yang sama melalui dua pemilik menyebabkan terjadinya kebebasan ganda.

## Realokasi untuk menambah ukuran

<!-- wave-example: book-alloc-grow -->
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
    var next: ptr<u8> = mem_realloc(data, 4, 8);

    if (next == null) {
        mem_free(data, 4);
        return 2;
    }

    data = next;
    deref data[4] = 9;
    println("{} {}", deref data[0], deref data[4]);

    if (mem_free(data, 8) < 0) {
        return 3;
    }

    return 0;
}
```

Hasil eksekusi:

```text
7 9
```

Pertama simpan hasilnya di berikutnya. Jika alokasi blok yang lebih besar gagal, data asli tetap valid dan masih dapat dibebaskan. Jika berhasil, alokasi lama akan dibebaskan dan alamat baru harus digunakan. Inisialisasi wilayah yang baru ditambahkan sebelum membacanya.

Ukuran baru dalam contoh ini adalah positif. Permintaan dengan new_size=0 malah mencoba mengosongkan alokasi yang ada dan mengembalikan null. Oleh karena itu, hasil null tidak selalu berarti alokasi lama masih valid. Hubungi mem_free secara langsung ketika Anda perlu memeriksa apakah pembebasan berhasil.

## Kapan Harus Memeriksa Ulang Petunjuk yang Dipinjam

data Menyimpan penunjuk internal dan menggunakannya setelah alokasi ulang adalah tindakan yang salah. Ini karena alamat data yang baru mungkin berbeda. Jika Anda memerlukan lokasi internal, Anda dapat menyimpan offset alih-alih alamat dan menghitung ulang berdasarkan data baru setelah berhasil.

Membebaskan atau mengalokasikan kembali memori juga mempengaruhi kode yang meminjamnya. Periksa apakah operasi lain masih menggunakan memori itu. Buffer yang diteruskan ke operasi asinkron harus tetap valid hingga operasi selesai.

## Daftar byte berisi Buffer

Mengelola daftar byte dengan panjang yang sering berubah saat melakukan realokasi secara manual memerlukan penanganan len dan cap, kegagalan ekspansi, dan perhitungan ukuran. Buffer dari std menggabungkan operasi ini.

<!-- wave-example: book-buffer-progressive -->
```wave
import("std::buffer::alloc")::{
    Buffer, buffer_init, buffer_free
};
import("std::buffer::write")::{
    buffer_append_str
};

fun main() -> i32 {
    var message: Buffer;

    if (buffer_init(&message, 0) < 0) {
        return 1;
    }

    if (buffer_append_str(&message, "Hello") < 0) {
        buffer_free(&message);
        return 2;
    }

    if (buffer_append_str(&message, ", Wave") < 0) {
        buffer_free(&message);
        return 3;
    }

    println("bytes={}", message.len);

    if (buffer_free(&message) < 0) {
        return 4;
    }

    return 0;
}
```

Hasil eksekusi:

```text
bytes=11
```

len adalah jumlah byte yang digunakan; cap adalah kapasitas yang dialokasikan. Menambahkan data akan meningkatkan alokasi bila diperlukan. Menggunakan Buffer tidak menghilangkan tanggung jawab penelepon untuk membebaskannya.

buffer_append_str tidak menambahkan NUL di akhir string. Oleh karena itu, message.data tidak boleh ditampilkan secara langsung sebagai str. Byte dikeluarkan dengan fungsi I/O, yang mengambil panjang, atau secara eksplisit membuat representasi string.

## Latihan dan pendekatan penyelesaian

Tambahkan byte 0 hingga 9 satu per satu ke Buffer dan dapatkan jumlahnya. Itu harus dibebaskan ketika setiap penambahan gagal, dan pembacaan hanya dilakukan dalam rentang len. Solusi lengkap dan kegagalan batas dapat diperiksa dengan mengikuti contoh di [Buffer Cara menggunakan](/docs/id/stdlib/buffer).

Coba tandai panggilan alokasi, alokasi ulang, dan dealokasi dalam kode Anda. Untuk setiap alokasi yang berhasil, Anda harus dapat menjelaskan siapa pemiliknya dan jalur mana yang membebaskannya.
