---
translation_set_id: types
path: language/declarations-and-types
locale: id
group: language
group_order: 2
order: 2
title: 2. Variabel, jenis, dan ruang lingkup
summary: Pelajari variabel lokal, cakupan bilangan bulat, bool, dan cakupan.
---

## Memperlakukan nilai berdasarkan nama

Jika Anda menulis harga secara langsung di beberapa tempat, Anda harus menemukan semuanya saat Anda mengubah harga. Variabel adalah ruang penyimpanan yang memberi nama pada nilai dan memungkinkan Anda membaca dan mengubahnya menggunakan nama tersebut. Dalam bab ini, Anda akan mempelajari tentang cakupan deklarasi, penugasan, tipe, dan cakupan nama yang terlihat dalam blok.

Program di bawah ini masing-masing merupakan bagian dari main.wave yang terpisah. Simpan satu contoh, jalankan sebagai `wavec run main.wave`, dan ganti dengan contoh berikutnya.

## Deklarasi dan inisialisasi

<!-- wave-example: book-variable-first -->
```wave
fun main() {
    var price: i32 = 1200;
    var quantity: i32 = 3;
    println("price={}", price);
    println("quantity={}", quantity);
    println("total={}", price * quantity);
}
```

Hasil eksekusi:

```text
price=1200
quantity=3
total=3600
```

Baca deklarasi dalam empat bagian.

|bagian|contoh ini|peran|
| --- | --- | --- |
|kata kunci deklarasi| var |Buat variabel lokal|
|nama| price |Pengenal yang akan digunakan nanti|
|mengetik| i32 |Ketik dan rentang nilai yang akan disimpan|
|nilai awal| 1200 |Nilai pertama yang disimpan|

Titik dua sebelum tipe dan tanda sama dengan sebelum nilai awal mempunyai peran yang berbeda. Jadikan nama Anda bermakna. Dalam contoh ini, price adalah harga satuan dan quantity adalah kuantitas. Meskipun i32 yang sama digunakan secara bergantian, penghitungan yang salah dapat dilakukan tanpa kesalahan tata bahasa.

## Penugasan bukanlah formula yang menjaga hubungan.

Jika Anda menyimpan hasil perhitungan dalam suatu variabel, maka nilai pada titik tersebut akan dimasukkan. Ia tidak mengingat perhitungan dan secara otomatis mengevaluasinya kembali nanti.

<!-- wave-example: book-assignment-snapshot -->
```wave
fun main() {
    var price: i32 = 1200;
    var quantity: i32 = 3;
    var total: i32 = price * quantity;
    quantity = 5;
    println("before recalculation={}", total);
    total = price * quantity;
    println("after recalculation={}", total);
}
```

Hasil eksekusi:

```text
before recalculation=3600
after recalculation=6000
```

`quantity = 5` menulis nilai baru ke variabel yang sudah ada. Harus dibedakan dengan deklarasi ulang seperti `var quantity`. `total` juga 3600 sebelum diganti lagi. Jika suatu program harus memelihara hubungan antara beberapa variabel, program tersebut harus ditulis untuk melakukan perhitungan ketika hubungan tersebut berubah.

## Hitung nilai berikutnya dari nilai sebelumnya

Hitung terlebih dahulu sisi kanan pernyataan penugasan dan tuliskan hasilnya ke ruang penyimpanan di sebelah kiri. Berbeda dengan persamaan dalam matematika, `count = count + 1` adalah pembaruan yang valid.

<!-- wave-example: book-update-variable -->
```wave
fun main() {
    var count: i32 = 0;
    count = count + 1;
    count += 2;
    count *= 3;
    println("{}", count);
}
```

Hasil eksekusi:

```text
9
```

Perkembangan nilai adalah 0 → 1 → 3 → 9. `+=`, `*=` perhitungan ekspres dan penyimpanan bersama. Jika Anda menuliskan urutan perhitungannya, Anda dapat mengetahui pada tahap mana pemikiran Anda berbeda ketika hasilnya berbeda dari yang Anda harapkan.

## Lebar dan tanda tipe integer

Jika dimulai dengan `i`, menjadi signed, dan jika dimulai dengan `u`, menjadi unsigned. Angka terakhir adalah jumlah bit. Ketika jumlah bit bertambah, rentang yang dapat dinyatakan bertambah dan ruang penyimpanan juga bertambah.

|mengetik|nilai minimum|nilai maksimum|Contoh penggunaan|
| --- | --- | --- | --- |
| i8 | -128 | 127 |nilai bertanda kecil|
| u8 | 0 | 255 |satu byte|
| i16 | -32768 | 32767 |data bilangan bulat kecil|
| u16 | 0 | 65535 |Bidang port/16-bit|
| i32 | -2147483648 | 2147483647 |Perhitungan bilangan bulat kecil yang umum|
| u32 | 0 | 4294967295 |bidang 32-bit|

Wave juga menyediakan bilangan bulat bertanda dan tidak bertanda 64, 128, 256, 512, dan 1024 bit. Tipe yang lebih luas tidak membuat setiap penghitungan aman: hasilnya masih bisa melebihi rentang yang dipilih. Pilih rentang yang Anda perlukan terlebih dahulu. `isz` dan `usz` mengikuti lebar alamat target.

Perbedaan harus dibuat antara menyimpan literal besar dalam tipe kecil dan sengaja membuang bit dengan melakukan cast. Mengonversi ke tipe yang lebih kecil hanya untuk menghilangkan kesalahan dapat mengubah nilainya sendiri. Transformasi dibahas dalam bab berikutnya.

## Tipe floating-point dan bool

`f32`·`f64` adalah angka floating point. Tidak seperti bilangan bulat, bilangan bulat dapat mewakili bagian desimal, namun bilangan desimal tidak dapat menyimpan semua bilangan desimal secara tepat. Inilah salah satu alasan mengapa jumlah dikelola dalam satuan bilangan bulat kecil.

bool mewakili benar dan salah. Anda dapat memberi nama pada kondisi Anda dengan menyimpan hasil perbandingan seperti ini:

<!-- wave-example: book-named-condition -->
```wave
fun main() {
    var balance: i32 = 5000;
    var cost: i32 = 3600;
    var can_buy: bool = balance >= cost;
    if (can_buy) {
        println("purchase allowed");
    } else {
        println("not enough money");
    }
}
```

Hasil eksekusi:

```text
purchase allowed
```

`can_buy` tidak otomatis mengikuti perubahan balance dan cost. Setelah kedua nilai ditukar, perbandingan dilakukan kembali jika diperlukan keadaan saat ini.

## ruang penyimpanan yang tidak diinisialisasi

`var value: i32;` adalah formulir yang hanya menyatakan ruang penyimpanan. Nilai yang valid harus ditulis sebelum dibaca. Jangan berasumsi bahwa 0 otomatis dimasukkan hanya karena Anda mendeklarasikannya. Dalam kursus pengantar, lebih mudah dipahami jika nilainya segera diketahui dan diinisialisasi bersamaan dengan deklarasi.

Dalam panggilan perpustakaan yang menerima nilai sebagai argumen keluaran, ada kasus di mana spasi dideklarasikan terlebih dahulu dan kemudian dibaca ketika berhasil. Pada saat itu, Anda harus memeriksa hasil keberhasilan fungsi tersebut. Hindari kesalahan membaca keluaran yang tidak diinisialisasi setelah panggilan gagal.

## Rentang blok dan nama yang valid

Blok adalah wilayah kode yang diapit kurung kurawal. Jika Anda mendeklarasikan variabel baru dengan nama yang sama di dalamnya, variabel baru tersebut akan digunakan di dalam blok tersebut. Ini disebut shadowing.

<!-- wave-example: book-shadowing -->
```wave
fun main() {
    var value: i32 = 10;

    if (true) {
        var value: i32 = value + 5;
        println("inner={}", value);
        value += 1;
        println("inner changed={}", value);
    }

    println("outer={}", value);
}
```

Hasil eksekusi:

```text
inner=15
inner changed=16
outer=10
```

Ekspresi awal `value + 5` dalam deklarasi bagian dalam membaca bagian luar value. Setelah inisialisasi variabel baru selesai, bagian dalam value adalah 15. Bahkan jika Anda mengubah nilai bagian dalam menjadi 16, ruang penyimpanan bagian luar tidak berubah. Setelah blok, Anda akan melihat value di luar lagi.

Sebaliknya, jika Anda hanya menjalankan `value += 1` tanpa `var` di blok dalam, variabel yang terlihat yang ada akan diubah. Lihatlah kata kunci untuk menentukan apakah itu merupakan deklarasi baru atau perubahan pada nilai yang sudah ada.

## Variabel lokal dan penyimpanan tingkat atas

Di luar fungsi tersebut, Anda dapat menggunakan const dan static. const mewakili nilai konstan dan static adalah ruang penyimpanan yang dipertahankan selama eksekusi. Mereka tidak dapat dideklarasikan di mana pun dengan cara yang sama seperti variabel lokal.

<!-- wave-example: book-storage -->
```wave
const LIMIT: i32 = 3;
static visits: i32 = 0;

fun visit() {
    visits += 1;
    println("visit={}", visits);
}

fun main() {
    visit();
    visit();
    println("limit={}", LIMIT);
}
```

Hasil eksekusi:

```text
visit=1
visit=2
limit=3
```

Meskipun visit dipanggil dua kali, static tidak akan disetel ulang ke 0 pada setiap panggilan. Di sisi lain, jika Anda mendeklarasikannya sebagai `var visits: i32 = 0;` di dalam suatu fungsi, penyimpanan lokal akan diinisialisasi setiap kali Anda memanggilnya. Status bersama yang dapat diubah dapat membuat perilaku sulit dilacak, jadi pertama-tama pertimbangkan apakah hal ini dapat diselesaikan dengan input dan output fungsi.

## kesalahan umum

- Ketika deklarasi dan penugasan membingungkan dan nama yang sama dideklarasikan lagi jika tidak diperlukan.
- Jika menurut Anda variabel yang menyimpan hasil perhitungan secara otomatis mengikuti perubahan pada variabel masukan.
- Jika menurut Anda karena tipenya sama, maka satuan seperti jumlah dan jumlah byte juga sama.
- Saat membaca tanpa inisialisasi atau saat argumen keluaran dibaca saat fungsi gagal.
- Jika menurut Anda nama variabel lokal terlihat di luar blok.

Saat mencari nama yang error, periksa posisi deklarasi nama dan rentang kurung kurawalnya.

## Latihan: Menghitung Perubahan Persediaan

Stok awal 20 unit dan terjual dua kali, masing-masing 3 unit. Cetak sisa persediaan dan total unit terjual. Setiap kali persediaan berubah, variabel yang sama diperbarui dan volume penjualan juga diakumulasikan secara terpisah.

### Solusi lengkap

<!-- wave-example: book-stock-solution -->
```wave
fun main() {
    var stock: i32 = 20;
    var sold: i32 = 0;
    var order: i32 = 3;
    stock -= order;
    sold += order;
    stock -= order;
    sold += order;
    println("stock={} sold={}", stock, sold);
}
```

Hasil eksekusi:

```text
stock=14 sold=6
```

Persediaan dan volume penjualan harus berubah secara bersamaan. Memperbarui salah satunya akan memutus hubungan antar nilai. Kami akan mengelompokkan pemrosesan pesanan berulang berdasarkan fungsi pembelajaran dan pernyataan loop.


## Tipe integer dan floating point

Tipe integer adalah sebagai berikut:

- Ditandatangani: `i8`, `i16`, `i32`, `i64`, `i128`, `i256`, `i512`, `i1024`
- Tidak bertanda tangan: `u8`, `u16`, `u32`, `u64`, `u128`, `u256`, `u512`, `u1024`
- Bilangan bulat ukuran alamat: `isz`, `usz`
- Titik mengambang: `f32`, `f64`

`isz` adalah tipe integer bertanda yang cocok dengan ukuran alamat, dan `usz` adalah tipe integer tak bertanda yang cocok dengan ukuran alamat.

## Tipe bawaan lainnya

|mengetik|Gunakan|
| --- | --- |
| `bool` |`true` atau `false`|
| `char` |Nilai karakter 8-bit yang tidak ditandatangani. Bukan jenis titik kode Unicode sembarangan|
| `byte` |Nilai byte 8-bit|
| `str` |String byte diakhiri dengan NUL|
| `ptr<T>` |Penargetan penunjuk `T`|
| `array<T, N>` |Array dengan panjang tetap dengan tipe elemen `T` dan panjang `N`|

Struktur, enumerasi, dan alias tipe yang ditentukan pengguna juga dapat digunakan di lokasi tipe.

`var` adalah sintaks untuk mendeklarasikan variabel lokal. Alias ​​tipe adalah tata bahasa yang mengekspresikan tipe yang sama dengan nama yang sesuai dengan konteks kode.
