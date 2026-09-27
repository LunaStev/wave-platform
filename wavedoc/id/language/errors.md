---
translation_set_id: learn-errors
path: language/errors
locale: id
group: language
group_order: 2
order: 12
title: 12. Mewakili dan memulihkan kesalahan
summary: Pisahkan kesalahan dari nilai normal dan bersihkan sumber daya dari jalur kegagalan.
---

## Hasil fungsi derajat kegagalan

File yang hilang atau input di luar batas terjadi secara alami dalam program. Menangani kesalahan lebih dari sekedar mencetak pesan. Ini adalah proses membedakan kegagalan, memeriksa status pekerjaan yang telah dilakukan, mengatur sumber daya yang diperoleh, dan kemudian memilih apakah akan melanjutkan atau mengakhiri.

Dalam bab ini, kita mulai dengan representasi kegagalan fungsi kecil dan berlanjut ke struktur hasil, variant, pengembalian awal, dan pembersihan sumber daya.

## Penanda kegagalan tidak boleh tumpang tindih dengan nilai keberhasilan

Alasan -1 digunakan sebagai tidak ditemukan dalam pencarian array adalah karena indeks efektifnya adalah 0 atau lebih. Di sisi lain, dalam perhitungan di mana bilangan bulat apa pun bisa menjadi hasil normal, jika -1 ditetapkan sebagai kesalahan, maka tidak dapat dibedakan dari nilai normal -1.

0 juga merupakan nilai yang sering disalahpahami. Panjang string kosong 0, posisi pertama 0, dan jumlah byte yang ditransfer 0 memiliki arti berbeda untuk fungsi berbeda. Jangan menilai kesuksesan hanya karena nilai pengembaliannya bukan nol.

## Kembalikan kesuksesan dan nilai bersama-sama

<!-- wave-example: book-error-result -->
```wave
struct Division {
    ok: bool;
    value: i32;
}

fun divide_nonnegative(left: i32, right: i32) -> Division {
    if (left < 0 || right <= 0) {
        return Division {
            ok: false,
            value: 0
        };
    }

    return Division {
        ok: true,
        value: left / right
    };
}

fun main() {
    var result: Division = divide_nonnegative(0, 3);

    if (!result.ok) {
        println("invalid input");
        return;
    }

    println("value={}", result.value);
}
```

Hasil eksekusi:

```text
value=0
```

Hasil normalnya mungkin juga 0. Daripada melihat value dan menebak-nebak apakah berhasil, periksa dulu ok. Meskipun kolom value ada pada hasil yang gagal, bukan berarti itu adalah nilai yang akan digunakan.

Aturan masukan untuk contoh ini adalah ruas kiri bernilai 0 atau lebih besar dan ruas kanan bernilai positif. Ruang lingkupnya ditampilkan dalam nama fungsi dan deskripsi untuk membedakannya dari pembagian bilangan bulat signed pada umumnya.

## Membedakan penyebab kesalahan

Tambahkan informasi kesalahan untuk panduan tambahan atau pemulihan tergantung pada alasan kegagalan. Ada juga cara untuk membagi fungsi menjadi fungsi yang memeriksa rentang masukan dan fungsi yang menghitungnya.

<!-- wave-example: book-error-codes -->
```wave
struct Check {
    ok: bool;
    error: i32;
}

fun check_quantity(value: i32) -> Check {
    if (value < 1) {
        return Check {
            ok: false,
            error: 1
        };
    }

    if (value > 100) {
        return Check {
            ok: false,
            error: 2
        };
    }

    return Check {
        ok: true,
        error: 0
    };
}

fun main() {
    var result: Check = check_quantity(101);

    if (!result.ok) {
        match (result.error) {
            1 => {
                println("quantity is too small");
            }
            2 => {
                println("quantity is too large");
            }
            _ => {
                println("invalid quantity");
            }
        }
    }
}
```

Hasil eksekusi:

```text
quantity is too large
```

Arti angka ditentukan oleh fungsi ini. Tidak bisa dianggap sama dengan kesalahan nomor 1 atau 2 di perpustakaan lain. Di publik API, penamaan konstanta atau jenis kesalahan membuat penelepon tidak perlu menghafal nomor acak.

## Pisahkan hasil dengan variant

Hubungan antara nilai keberhasilan dan kesalahan tidak dapat terjadi secara bersamaan dapat dinyatakan sebagai variant.

<!-- wave-example: book-error-variant -->
```wave
variant ByteResult {
    Value(u8),
    OutOfRange(i32)
}

fun to_byte(value: i32) -> ByteResult {
    if (value < 0 || value > 255) {
        return ByteResult::OutOfRange(value);
    }

    return ByteResult::Value(value as u8);
}

fun main() {
    var result: ByteResult = to_byte(300);

    match (result) {
        ByteResult::Value(value) => {
            println("byte={}", value);
        }
        ByteResult::OutOfRange(original) => {
            println("out of range: {}", original);
        }
    }
}
```

Hasil eksekusi:

```text
out of range: 300
```

Kesalahan berisi nilai masukan asli. Lebih mudah bagi penelepon untuk menjelaskan masalahnya daripada sekadar mengembalikan false. Sifat data juga dipertimbangkan dengan tidak meninggalkan input sensitif seperti kata sandi atau token di log.

## Jadikan rute normal lebih mudah dibaca dengan pengembalian lebih awal

Tidak perlu memasukkan semua kode bagus jauh di dalam if saat melakukan pemeriksaan multi-langkah. Jika gagal, Anda dapat segera kembali dan melanjutkan jalur sukses di bawah ini.

<!-- wave-example: book-error-early-return -->
```wave
fun show_total(quantity: i32, price: i32) {
    if (quantity < 1 || quantity > 100) {
        println("invalid quantity");
        return;
    }

    if (price < 0 || price > 100000) {
        println("invalid price");
        return;
    }

    var total: i32 = quantity * price;
    println("total={}", total);
}

fun main() {
    show_total(3, 1200);
    show_total(0, 1200);
    show_total(3, -1);
}
```

Hasil eksekusi:

```text
total=3600
invalid quantity
invalid price
```

Setelah lulus tes, Anda dapat memanfaatkan fakta bahwa quantity dan price berada dalam rentang yang ditentukan. Kisarannya diatur sehingga perkalian antara juga berada dalam i32. Saat menambahkan pengembalian awal, Anda juga harus memeriksa apakah Anda sudah memiliki sumber daya pada saat itu.

## Pembersihan memori dari jalur kegagalan

<!-- wave-example: book-error-cleanup -->
```wave
import("std::buffer::alloc")::{
    Buffer, buffer_init, buffer_free
};
import("std::buffer::write")::{
    buffer_append_str
};

fun make_message(allowed: bool) -> i32 {
    var data: Buffer;

    if (buffer_init(&data, 16) < 0) {
        return 1;
    }

    if (!allowed) {
        buffer_free(&data);
        return 2;
    }

    if (buffer_append_str(&data, "ready") < 0) {
        buffer_free(&data);
        return 3;
    }

    println("bytes={}", data.len);

    if (buffer_free(&data) < 0) {
        return 4;
    }

    return 0;
}

fun main() {
    println("status={}", make_message(false));
}
```

Hasil eksekusi:

```text
status=2
```

Lepaskan Buffer bahkan pada jalur kegagalan yang tidak mengeluarkan data. Daripada mengubah setiap kegagalan menjadi satu return, penting untuk mengelola dengan jelas apa yang dimiliki di setiap cabang.

Ada juga API yang pembersihannya sendiri gagal. Putuskan bagaimana Anda akan mempertahankan kesalahan dan pembersihan kesalahan dalam karya aslinya. Contoh kecil ini pertama-tama mengembalikan nomor kesalahan tugas. Program yang lebih besar dapat merekam keduanya secara terpisah.

## Keberhasilan sebagian tidak secara otomatis dibatalkan

Jika Anda menulis beberapa byte ke suatu file dan kemudian penulisannya gagal, byte yang sudah ditulis tidak akan hilang. Ujung lain jaringan mungkin juga menerima beberapa data. Mengulangi tugas yang sama dari awal dapat mengakibatkan rekaman duplikat.

Sebaliknya, membaca checked byte cursor mempertahankan posisi dan nilai keluaran ketika gagal. Fungsi-fungsi ini dapat menerima lebih banyak masukan dan mencoba lagi di lokasi yang sama. Daripada menerapkan aturan “jika gagal, tidak ada perubahan” pada setiap API, periksa dokumentasi untuk fungsi tersebut.

## Kesalahan dan jebakan yang dapat dipulihkan

Jalur file yang tidak valid atau input pengguna yang tidak valid dapat dilaporkan sebagai nilai kesalahan sehingga pemanggil dapat memulihkannya. Jebakan yang disebabkan oleh jumlah pergeseran runtime yang tidak valid atau konversi titik-mengambang ke bilangan bulat adalah mekanisme yang berbeda.

Program yang harus terus dijalankan harus memeriksa masukannya sebelum melakukan operasi yang berbahaya. assert juga tidak dimaksudkan untuk disalahgunakan sebagai cara menangani kegagalan masukan pengguna. Jika Anda perlu memberi pengguna kesempatan untuk memasukkan kembali masukannya, kembalikan hasil untuk melanjutkan alur kontrol.

## Latihan dan solusi lengkap

Tulis fungsi yang membaca elemen array berdasarkan indeks dan gagal ketika indeksnya negatif atau lebih besar atau sama dengan panjangnya. Gunakan struktur hasil karena nol dapat menjadi nilai elemen yang valid.

<!-- wave-example: book-error-solution -->
```wave
struct Lookup {
    ok: bool;
    value: i32;
}

fun at(values: ptr<i32>, length: i32, index: i32) -> Lookup {
    if (index < 0 || index >= length) {
        return Lookup {
            ok: false,
            value: 0
        };
    }

    return Lookup {
        ok: true,
        value: deref values[index]
    };
}

fun main() {
    var values: array<i32, 3> = [0, 10, 20];
    var first: Lookup = at(&values[0], 3, 0);
    var outside: Lookup = at(&values[0], 3, 3);

    if (first.ok) {
        println("value={}", first.value);
    }

    if (!outside.ok) {
        println("out of bounds");
    }
}
```

Hasil eksekusi:

```text
value=0
out of bounds
```

Ini adalah kondisi pemanggil bahwa penunjuk dan length mewakili array sebenarnya yang dapat dibaca. Ini bukan fungsi yang membuat alamat acak aman hanya dengan memeriksa indeks. Harap baca secara terpisah pemeriksaan apa yang menjadi tanggung jawab fungsi tersebut dan kondisi apa yang harus dijamin oleh penelepon.
