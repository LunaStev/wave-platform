---
translation_set_id: data-types
path: language/structures-enums-and-aliases
locale: id
group: language
group_order: 2
order: 8
title: 8. Struktur, enum, dan varian
summary: Pelajari tentang bidang, inisialisasi struktur, dan peran enum dan variant.
---

## Mengekspresikan hubungan antar data sebagai tipe

Jika harga dan jumlah produk masing-masing dilewatkan sebagai variabel, sulit untuk mengetahui hanya dengan melihat kode apakah kedua nilai tersebut milik produk yang sama. Struktur mengelompokkan bidang terkait. enum mewakili status bernama, dan variant menyimpan data berbeda secara bersamaan untuk setiap kasus.

Ketiga fungsi tersebut bukanlah sintaksis yang saling menggantikan. Pilih tergantung pada apa yang ingin Anda ekspresikan.

|sesuatu untuk diungkapkan|pilih|ya|
| --- | --- | --- |
|Beberapa bidang ada secara bersamaan| struct |Harga satuan dan kuantitas produk|
|Dinamakan keadaan bilangan bulat| enum |Tunggu/Lanjutkan/Selesai|
|Data berbeda dalam setiap kasus| variant |nilai keberhasilan atau kesalahan|

## Deklarasi struktur dan penciptaan nilai

<!-- wave-example: book-struct-first -->
```wave
struct Product {
    price: i32;
    quantity: i32;
}

fun main() {
    var item: Product = Product {
        price: 1500,
        quantity: 2
    };

    println("{} {}", item.price, item.quantity);
}
```

Hasil eksekusi:

```text
1500 2
```

Bidang dalam deklarasi diakhiri dengan titik koma, dan saat membuat nilai, bidang dan nilai dihubungkan dengan titik dua dan dipisahkan dengan koma. Mendefinisikan tipe dan membuat nilai sebenarnya adalah dua langkah berbeda. Mendeklarasikan jenis Product tidak secara otomatis membuat ruang penyimpanan untuk satu produk.

Bidang diakses sebagai `item.price`. Jika Anda membuat beberapa item dengan tipe yang sama, Anda dapat menyimpan nilai berbeda di masing-masingnya.

## Meneruskan struktur ke suatu fungsi

<!-- wave-example: book-struct-function -->
```wave
struct Product {
    price: i32;
    quantity: i32;
}

fun subtotal(item: Product) -> i32 {
    return item.price * item.quantity;
}

fun main() {
    var item: Product = Product {
        price: 1500,
        quantity: 2
    };

    println("{}", subtotal(item));

    item.quantity = 3;
    println("{}", subtotal(item));
}
```

Hasil eksekusi:

```text
3000
4500
```

Fungsi tersebut menerima sebagai tipe hubungan bahwa harga satuan dan kuantitas dimiliki oleh produk yang sama. Ini adalah fungsi yang membaca struktur bidang bilangan bulat yang diteruskan sebagai nilai. Fungsi apa pun yang ingin mengubah penyimpanan pemanggil dapat dirancang untuk menerima pointer.

Jika struktur memiliki bidang penunjuk, nilai penyalinan juga akan menyalin alamatnya. Ini bukan fungsi untuk penyalinan mendalam ke alokasi terpisah. Jenis yang berisi sumber daya seperti pegangan file atau Buffer harus menentukan aturan penyalinan dan rilis secara bersamaan.

## Metode dan proto

Fungsi terkait dapat dikelompokkan dalam bentuk metode. proto adalah metode penulisan metode struktur sebagai blok terpisah.

<!-- wave-example: book-struct-method -->
```wave
struct Counter {
    value: i32;
}

proto Counter {
    fun current(self: Counter) -> i32 {
        return self.value;
    }
}

fun main() {
    var counter: Counter = Counter {
        value: 3
    };

    println("{}", counter.current());
}
```

Hasil eksekusi:

```text
3
```

`self: Counter` adalah parameter yang menerima nilai. Menggunakan notasi pemanggilan metode tidak secara otomatis menjadikannya metode yang memodifikasi metode asli. Silakan baca bersama-sama jenis self dan fungsinya dalam teks.

Anda tidak dapat mengharapkan suatu bidang hanya memiliki status valid hanya karena Anda melampirkan metode ke dalamnya. Jika ada kombinasi tidak valid yang dapat dibuat pengguna dengan kolom publik, fungsi Anda harus memeriksanya atau memberikan aturan pembuatan.

## Beri nama negara bagian tersebut dengan enum

<!-- wave-example: book-enum-state -->
```wave
enum State -> i32 {
    Ready = 0,
    Running,
    Finished,
}

fun main() {
    var state: State = State::Ready;

    if (state == State::Ready) {
        println("ready");
    }

    state = State::Running;

    if (state == State::Running) {
        println("running");
    }
}
```

Hasil eksekusi:

```text
ready
running
```

`-> i32` adalah tipe integer yang digunakan dalam ekspresi. Nilai pertama disetel ke 0, dan nilai berikutnya yang dihilangkan adalah 1 lebih besar dari nilai sebelumnya. Daripada hanya membandingkan 0 dan 1 dalam kode Anda, menggunakan State::Ready dan State::Running akan mengungkapkan artinya.

enum Memiliki nama tidak secara otomatis membatasi transisi negara. Aturan seperti apakah mungkin untuk kembali dari Finished ke Running harus diterapkan sebagai fungsi.

## Menghubungkan kasus dan data dengan variant

Jika terdapat nilai hanya jika berhasil dan informasi kesalahan diperlukan jika terjadi kegagalan, nilai tersebut dapat dinyatakan sebagai variant.

<!-- wave-example: book-variant-result -->
```wave
variant Result {
    Value(i32),
    Error(i32)
}

fun divide(left: i32, right: i32) -> Result {
    if (left < 0 || right <= 0) {
        return Result::Error(1);
    }

    return Result::Value(left / right);
}

fun show(result: Result) {
    match (result) {
        Result::Value(value) => {
            println("value={}", value);
        }
        Result::Error(code) => {
            println("error={}", code);
        }
    }
}

fun main() {
    show(divide(12, 3));
    show(divide(12, 0));
}
```

Hasil eksekusi:

```text
value=4
error=1
```

Result::Value dan Result::Error masing-masing berisi payload. Meskipun keduanya berisi tipe bilangan bulat yang sama, kasus tertentu dapat dibedakan. Penelepon memeriksa kasus dengan match dan menggunakan payload di dalamnya arm.

divide di atas adalah contoh kecil yang tidak mendukung operan negatif. Karena rentang masukan sudah ditentukan, jangan bingung dengan fungsi yang mencakup semua batasan pembagian signed reguler.

## Jika payload tidak ada

Data tidak diperlukan dalam semua kasus. Keadaan tidak bernilai dapat dinyatakan sebagai kasus tersendiri.

<!-- wave-example: book-variant-empty -->
```wave
variant Lookup {
    Found(i32),
    Missing
}

fun main() {
    var result: Lookup = Lookup::Missing;

    match (result) {
        Lookup::Found(value) => {
            println("found={}", value);
        }
        Lookup::Missing => {
            println("missing");
        }
    }
}
```

Hasil eksekusi:

```text
missing
```

Daripada menyimpan satu bilangan bulat sembarang sebagai “tidak ada”, kami menggunakan kasus Missing. Apapun nilai kesuksesannya, maknanya tidak tumpang tindih.

match hingga `_` menangani kasus yang tersisa. Jika Anda ingin setiap penelepon ditinjau kembali ketika kasus baru ditambahkan, lebih baik pisahkan semua kasus secara eksplisit. Apapun metode yang Anda pilih, pastikan tidak ada input yang belum diproses.

## Menggunakan struktur dan variant bersama-sama

Memilih data yang berbeda dapat dinyatakan sebagai variant, dan mengelompokkan beberapa bidang milik satu kasus dapat dinyatakan sebagai sebuah struktur. Misalnya, jika hasil pemrosesan pesanan berhasil, maka dapat dirancang untuk memuat struktur penerimaan, dan jika gagal, dapat dirancang untuk memuat nomor kesalahan.

Aturan seumur hidup memori tidak hilang meskipun nilainya berisi nilai lain. Jika sebuah pointer disimpan di variant, apakah pointer tersebut valid dan siapa yang akan melepaskannya ditentukan secara terpisah. Bahkan ketika menyimpan dalam format file eksternal, Anda harus menentukan pengkodean untuk setiap bidang daripada membuang memori struktur apa adanya.

## Latihan dan solusi lengkap

Buat keputusan yang hanya menerima skor antara 0 dan 100. Skor yang valid dinyatakan sebagai Grade(score), sisanya dinyatakan sebagai Invalid.

<!-- wave-example: book-model-solution -->
```wave
variant CheckedScore {
    Grade(i32),
    Invalid
}

fun check_score(score: i32) -> CheckedScore {
    if (score < 0 || score > 100) {
        return CheckedScore::Invalid;
    }

    return CheckedScore::Grade(score);
}

fun main() {
    var result: CheckedScore = check_score(87);

    match (result) {
        CheckedScore::Grade(score) => {
            println("accepted={}", score);
        }
        CheckedScore::Invalid => {
            println("invalid score");
        }
    }
}
```

Hasil eksekusi:

```text
accepted=87
```

Ubah input menjadi -1, 0, 100, 101 untuk memeriksa batasnya. Karena keberhasilan dan kegagalan tidak berbagi ruang bilangan bulat yang sama, pemanggil mengurangi risiko penambahan nilai kesalahan secara tidak sengaja ke penghitungan rata-rata.


## Ekspresi mana yang harus saya pilih?

|bentuk data|ekspresi yang cocok|ya|
| --- | --- | --- |
|Beberapa nilai dari tipe yang sama|Himpunan|10 poin|
|Beberapa bidang terkait satu sama lain|struktur|nama dan skor|
|Nilai negara bernama| enum | Ready, Running, Stopped |
|Data tambahan yang bervariasi menurut negara bagian| variant | Value(i32), Error(str) |
|Nama kontekstual untuk tipe yang sudah ada|ketik alias| UserId = u64 |

Saat memilih struktur data, pertimbangkan tidak hanya nilai yang ingin Anda simpan, tetapi juga status salah yang dapat Anda wakili. Struktur dengan bidang keberhasilan/kegagalan dan nilai/kesalahan dapat menghasilkan kombinasi yang salah, namun variant dapat dibedakan menjadi payload untuk setiap kasus.

## Penempatan memori dan data eksternal

Memori struktur mungkin berisi ruang kosong untuk memastikan keselarasan antar bidang. Menambahkan ukuran bidang saja tidak selalu sama dengan ukuran total struktur. Jika Anda perlu mengetahui ukuran dan perataannya, gunakan [mem Fungsi tata letak](/docs/id/reference/memory-and-buffer).

Untuk file atau pesan jaringan, urutan dan panjang byte dapat ditentukan dengan jelas dengan mengkodekan bidang secara berurutan dengan fungsi [bytes](/docs/id/stdlib/bytes). Saat meneruskan struktur ke bahasa lain, sejajarkan deklarasi eksternal [FFI](/docs/id/language/modules-imports-and-ffi) dengan target ABI.

[Struktur pembelajaran dan praktik](/docs/id/language/structures-enums-and-aliases) · [variant](/docs/id/language/variants)
