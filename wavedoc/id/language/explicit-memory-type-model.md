---
translation_set_id: memory-model
path: language/explicit-memory-type-model
locale: id
group: language
group_order: 2
order: 9
title: 9. Petunjuk, pengubah nilai, dan masa hidup
summary: Pelajari alamat, dereferensi, modifikasi nilai melalui pointer, dan pointer yang menggantung.
---

## Bedakan antara nilai dan lokasi penyimpanan

Bilangan bulat 42 dan alamat tempat bilangan bulat tersebut disimpan mempunyai nilai yang berbeda. Penunjuk menunjuk ke lokasi penyimpanan. Melewati alamat memungkinkan fungsi membaca atau mengubah penyimpanan pemanggil.

Bab ini mencakup pengambilan alamat, dereferensi, modifikasi nilai asli, aritmatika penunjuk, dan masa pakai. Alokasi dinamis dibahas pada bab berikutnya. Mulailah dengan alamat variabel lokal dan elemen array.

## Mengambil alamat dan dereferensi

<!-- wave-example: book-pointer-first -->
```wave
fun main() {
    var value: i32 = 42;
    var address: ptr<i32> = &value;

    println("value={}", value);
    println("through pointer={}", deref address);

    deref address = 99;
    println("changed={}", value);
}
```

Hasil eksekusi:

```text
value=42
through pointer=42
changed=99
```

`&value` mendapatkan alamatnya, `deref address` membaca atau menulis nilai alamat itu, dan seterusnya. Daripada menyimpan 99 di address, 99 ditulis ke bilangan bulat yang ditunjukkan oleh address. Variabel address sendiri terus mengarah ke value.

`ptr<i32>` adalah tipe penunjuk untuk mengakses penyimpanan i32. Tipe ini tidak mencatat panjangnya atau menyediakan deallokasi otomatis.

## Mengubah penunjuk itu sendiri

<!-- wave-example: book-pointer-reassign -->
```wave
fun main() {
    var first: i32 = 10;
    var second: i32 = 20;
    var selected: ptr<i32> = &first;

    selected = &second;
    deref selected = 25;

    println("{} {}", first, second);
}
```

Hasil eksekusi:

```text
10 25
```

`selected = &second` menyimpan alamat lain dalam variabel penunjuk. Nilai first tidak berubah. Kemudian, jika Anda menulis nilainya sebagai deref, second akan berubah. Memisahkan “perubahan alamat” dan “perubahan nilai melalui alamat” menjadi kalimat terpisah akan mengurangi kebingungan.

## Membuat suatu fungsi mengubah aslinya

<!-- wave-example: book-pointer-increment -->
```wave
fun increment(value: ptr<i32>) {
    deref value = deref value + 1;
}

fun main() {
    var count: i32 = 4;

    increment(&count);
    increment(&count);

    println("{}", count);
}
```

Hasil eksekusi:

```text
6
```

Fungsi tersebut menerima alamat hitungan, bukan nilainya, 4. Mengubah penyimpanan itu juga mengubah jumlah pemanggil. Fungsi ini memerlukan alamat i32 yang valid dan dapat ditulis. Melewati null melanggar persyaratan ini.

Setiap fungsi menentukan apakah ia menerima null. Jika tidak, penelepon harus memberikan alamat yang valid. Jika ya, fungsi tersebut harus menyertakan jalur yang menangani null.

## Berfungsi untuk memproses null

<!-- wave-example: book-pointer-null -->
```wave
fun try_increment(value: ptr<i32>) -> bool {
    if (value == null) {
        return false;
    }

    deref value = deref value + 1;
    return true;
}

fun main() {
    var count: i32 = 7;

    if (!try_increment(null)) {
        println("no value");
    }

    if (try_increment(&count)) {
        println("count={}", count);
    }
}
```

Hasil eksekusi:

```text
no value
count=8
```

Cek null hanya menangani tidak adanya alamat. Mengonversi angka non-null sembarang menjadi pointer tidak menghasilkan memori yang valid. Membaca dan menulis juga memerlukan masa berlaku yang valid, ukuran yang memadai, penyelarasan yang benar, dan izin akses yang sesuai.

## Alamat array dan aritmatika penunjuk berdasarkan elemen

<!-- wave-example: book-pointer-array -->
```wave
fun main() {
    var values: array<i32, 3> = [10, 20, 30];
    var first: ptr<i32> = &values[0];
    var second: ptr<i32> = first + 1;

    println("{}", deref first);
    println("{}", deref second);
    println("{}", deref first[2]);
}
```

Hasil eksekusi:

```text
10
20
30
```

Menambahkan 1 ke pointer akan memindahkannya satu elemen dari tipe target. Elemen berikutnya di i32 dan elemen berikutnya di u8 memiliki jumlah shift byte yang berbeda. Jika Anda mengalikan `first + 1` lagi dengan ukuran jenis dan menambahkannya, maka akan berpindah ke posisi yang tidak diinginkan.

Pengindeksan pointer juga harus dilakukan dalam rentang yang valid. first tidak mengingat panjang array 3 itu sendiri, jadi ketika meneruskan rentang ke fungsi, ia menggunakan formulir yang menerima penunjuk dan panjangnya.

## Melewati rentang baca ke suatu fungsi

<!-- wave-example: book-pointer-range -->
```wave
fun sum(values: ptr<i32>, count: i32) -> i32 {
    var total: i32 = 0;

    for (var index: i32 = 0; index < count; index += 1) {
        total += deref values[index];
    }

    return total;
}

fun main() {
    var values: array<i32, 4> = [2, 4, 6, 8];

    println("first two={}", sum(&values[0], 2));
    println("all={}", sum(&values[0], 4));
}
```

Hasil eksekusi:

```text
first two=6
all=20
```

Satuan count adalah jumlah elemen. Fungsi ini mengharuskan penelepon memiliki count yang dapat dibaca i32. Melewati panjang yang lebih besar dari array sebenarnya melanggar kontrak. Bahkan jika Anda menggunakan kombinasi ptr<u8>·i64 yang sama, API, Anda harus memeriksa dokumentasi apakah panjangnya dalam byte atau dalam elemen.

## Umur: Berapa lama alamat tersebut valid?

Variabel lokal digunakan selama masa panggilan dan blok. Jika Anda mengembalikan alamat variabel lokal di dalam fungsi agar pemanggil dapat membacanya nanti, ruang penyimpanan tersebut mungkin telah mencapai akhir masa pakainya.

Inilah bagian dari desain buruk yang tidak boleh diterapkan:

```wave
fun invalid_address() -> ptr<i32> {
    var local: i32 = 42;

    return &local;
}
```

Jika satu nilai diperlukan, ia akan mengembalikan i32. Jika perlu menulis ke ruang penyimpanan yang disediakan oleh pemanggil, dibutuhkan sebuah pointer sebagai input. Jika Anda memerlukan penyimpanan terpisah untuk disimpan di luar panggilan, alokasikan secara eksplisit dan serahkan tanggung jawab untuk mengosongkannya.

## Dua petunjuk menunjuk ke ruang penyimpanan yang sama

Menyalin penunjuk akan membuat nama lain yang menunjuk ke alamat yang sama. Tidak menduplikasi memori.

<!-- wave-example: book-pointer-alias -->
```wave
fun main() {
    var value: i32 = 1;
    var first: ptr<i32> = &value;
    var second: ptr<i32> = first;

    deref second = 9;

    println("{} {}", value, deref first);
}
```

Hasil eksekusi:

```text
9 9
```

Perubahan hasil melalui second juga dapat dilihat melalui first. Setelah memori sumber dibebaskan, kedua pointer menjadi tidak dapat digunakan. Menetapkan null ke satu variabel penunjuk tidak secara otomatis mengubah salinan lainnya.

## Latihan: Tukarkan dua bilangan bulat

Tulis fungsi yang mengambil dua alamat i32 dan menukar nilainya. Nilai pertama harus disimpan dalam variabel sementara sebelum ditimpa. Periksa apakah nilainya dipertahankan meskipun Anda memberikan alamat yang sama dua kali.

### Solusi lengkap

<!-- wave-example: book-pointer-swap -->
```wave
fun swap(left: ptr<i32>, right: ptr<i32>) {
    var saved: i32 = deref left;

    deref left = deref right;
    deref right = saved;
}

fun main() {
    var first: i32 = 3;
    var second: i32 = 8;

    swap(&first, &second);
    println("{} {}", first, second);

    swap(&first, &first);
    println("{}", first);
}
```

Hasil eksekusi:

```text
8 3
8
```

Fungsi ini juga memerlukan kedua alamat untuk menunjuk ke penyimpanan bilangan bulat yang valid dan dapat ditulis. Untuk menangani null, tambahkan hasil yang menunjukkan keberhasilan atau kegagalan, seperti pada try_increment.


## **Wave Explicit Memory Type Model**

Desain penunjuk Wave didasarkan pada **Wave Explicit Memory Type Model**. Model ini mendefinisikan pointer dan array sebagai tipe memori eksplisit pada tingkat bahasa, bukan trik sintaksis atau abstraksi perpustakaan.

`ptr<T>` adalah tipe yang menunjuk ke alamat memori yang menyimpan nilai `T`, dan `array<T, N>` adalah tipe memori dengan panjang tetap yang menyimpan nilai `N` dari `T` secara berurutan. Oleh karena itu, struktur pointer dan array ditampilkan sebagaimana adanya dalam argumen fungsi, nilai kembalian, bidang struktur, dan tipe lainnya.

## null

```wave
var buffer: ptr<u8> = null;
if (buffer == null) {
    println("no buffer");
}
```

`null` adalah nilai penunjuk yang tidak menunjuk ke alamat memori yang valid. `null` hanya dapat ditetapkan ke tipe `ptr<T>` dan tidak dapat digunakan sebagai nilai integer, Boolean, atau array.

Fungsi alokasi atau pencarian dapat mengembalikan `null` ketika tidak ada hasil. Periksa `null` sebelum melakukan dereferensi hasil tersebut. Dereferensi penunjuk `null` tidak mengakses penyimpanan yang valid.

## konversi penunjuk

Saat Anda perlu mengubah alamat atau representasi penunjuk lainnya, gunakan `as`.

```wave
var raw: i64 = 0;
var p: ptr<u8> = raw as ptr<u8>;
```

Gunakan konversi antara bilangan bulat dan pointer hanya pada batas tingkat rendah, dan pertimbangkan lebar alamat platform target dan ABI.
