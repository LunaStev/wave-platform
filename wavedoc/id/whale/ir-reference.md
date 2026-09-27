---
translation_set_id: whale-ir-reference
path: whale/ir-reference
locale: id
group: whale
group_order: 1
order: 6
title: Referensi Whale IR
summary: Menjelaskan jenis, pengidentifikasi, validitas fungsi, urutan evaluasi, dan format pertukaran.
---

## Modul dan Pengidentifikasi

Modul terdiri dari informasi target, definisi global, dan fungsi. Nilai memiliki tipe eksplisit. Frontend menyelesaikan nama, jenis, kelebihan beban, dan generik bahasa sumber dan menghasilkan typed IR.

Fungsi dan variabel global menggunakan namespace internal yang berbeda. Oleh karena itu, fungsi dan variabel dapat memiliki nama yang sama. Pengidentifikasi internal berbeda dari nama koneksi eksternal `link_name`, dan nama eksternal ditentukan oleh frontend. Whale tidak menyelesaikan konflik eksternal dengan membuat nama baru secara otomatis. Silakan merujuk ke [Simbol dan Tautan](assembler-linker).

Setiap definisi nilai memiliki pengidentifikasi. Definisi tidak dapat diduplikasi, dan tipe metadata harus cocok dengan tipe yang ditentukan dalam definisi. Nama saja tidak dapat mengidentifikasi definisi, meskipun deklarasi dengan nama yang sama mengaburkan satu sama lain.

## IR Konfigurasi dan pembacaan

Di bawah ini adalah contoh lengkap Rust yang membangun dan memverifikasi fungsi dengan peti `ir`, lalu menghasilkan keluaran typed IR.

```rust
use ir::{ModuleBuilder, Target, Type};

fn main() {
    let target = Target::lookup("x86_64-whale-linux").unwrap();
    let mut module = ModuleBuilder::new(target.name(), target.data_layout());
    let mut function = module.begin_function("answer", vec![], Type::I32);
    let left = function.const_i32(40);
    let right = function.const_i32(2);
    let answer = function.add(Type::I32, left, right);
    function.ret(Some(answer));
    function.finish();
    let module = module.finish();
    ir::verify_module(&module).unwrap();
    print!("{}", ir::print_module(&module));
}
```

Output printer IR:

```text
module {
  format_version 2
  semantics_version 1
  target "x86_64-whale-linux"
  datalayout { ptr=64, endian=little }

  declare @f0 "answer": whale () -> i32, linkage internal

  fn @answer() -> i32, id @f0 {
  entry:
    %v0: i32 = const i32 40
    %v1: i32 = const i32 2
    %v2: i32 = add i32 %v0, %v1
    ret i32 %v2
  }

}
```

`%v0` dan `%v1` adalah definisi konstanta i32. Gunakan `%v2` yang didefinisikan oleh `add` sebagai nilai kembalian i32 dari fungsi tersebut. Jika Anda mengubah nilai kembalian menjadi Bool, ini akan menyebabkan kesalahan verifikasi karena tidak cocok dengan tanda tangan fungsi. Meskipun kedua operan adalah konstanta, O0 mempertahankan instruksi penjumlahan.

Kode ini adalah keluaran printer sebenarnya, bukan file masukan untuk diteruskan ke pengurai teks. Penguraian teks dan eksekusi IR belum didukung, saat ini modul ini dapat dikonfigurasi sebagai Rust builder.

## mengetik

|mengetik|artinya|
| --- | --- |
| `bool` |Nilai logika false atau true|
| `i1`, `i8`, `i16`, `i32`, `i64`, `i128` |bilangan bulat bertanda dengan lebar bit tertentu|
| `u1`, `u8`, `u16`, `u32`, `u64`, `u128` |bilangan bulat tak bertanda dengan lebar bit tertentu|
| `f16`, `f32`, `f64` |Nilai floating point dengan lebar bit tertentu|
| `ptr<T>` |T Penunjuk untuk mengetik nilai|
| `fnptr<signature>` | Pointer yang dapat dipanggil dengan parameter/tipe hasil yang tepat dan konvensi pemanggilan |
| `array<T, N>` |N elemen bertipe sama|
| `struct{T, ...}` |Bidang struktur yang dipesan|
| `tuple<T, ...>` |elemen tupel yang dipesan|
| `void` |Tidak ada hasil|

`bool`, `i1`, dan `u1` adalah tipe yang berbeda. signed `i1` mewakili −1 dan 0, dan unsigned `u1` mewakili 0 dan 1. Bilangan bulat 1 bukan merupakan kondisi logika implisit. Cabang bersyarat, kondisi Select, `trap_if` memerlukan operan Bool.

Ukuran penyimpanan tidak hanya ditentukan oleh jumlah bit dalam nilainya dan mengikuti [tata letak sasaran](memory-model). Misalnya, nilai `i1` adalah 1 bit, tetapi memakan setidaknya 1 byte di memori.

## Fungsi dan Panggilan

Fungsi ini menentukan semua parameter, tipe hasil, konvensi pemanggilan, dan linkage. Panggilan langsung dan tidak langsung harus sesuai dengan tanda tangan penelepon. Panggilan void tidak membuahkan hasil ID. Panggilan ke nonvoid mempertahankan definisi hasil meskipun O0 tidak menggunakan hasilnya.

Pengembaliannya harus sesuai dengan jenis hasil fungsi. Pengembalian void tidak membawa nilai, dan pengembalian nonvoid membawa nilai dari tipe hasil yang dideklarasikan.

### Deklarasi, identitas dan panggilan

`Module.declarations` mencatat `FunctionId` setiap fungsi, nama, tanda tangan lengkap, tautan, dan nama tautan eksternal. Definisi mengacu pada identitas ini; parameter dan tipe kembalian harus sesuai dengan deklarasinya. Deklarasi berulang yang identik diselesaikan ke ID yang sama melalui `declare_function`; konflik dan definisi duplikat adalah kesalahan. Deklarasi internal memerlukan isi modul. Deklarasi eksternal mungkin belum terselesaikan hingga tertaut, atau memiliki badan yang diekspor. Fungsi internal tidak memiliki `link_name`; fungsi eksternal memerlukan nama kosong yang eksplisit tanpa NUL. Dua deklarasi fungsi yang berbeda tidak dapat mengklaim nama eksternal yang sama. Global dan fungsi masih menggunakan namespace internal yang terpisah.

Daftarkan deklarasi sebelum membuat badan dengan `begin_declared_function` untuk mendukung panggilan penerusan dan rekursi. `begin_function` tetap menjadi kemudahan untuk fungsi Whale internal baru. API `declare_function`, `begin_declared_function`, `function_addr`, `null_function` dan `call` yang dicentang akan menghasilkan `Result`; panggilan yang ditolak tidak menambahkan instruksi atau mengalokasikan ID hasilnya.

Program Rust lengkap berikut ini mendeklarasikan fungsi eksternal, mengambil alamat yang diketiknya, dan melakukan panggilan langsung dan tidak langsung:

```rust
use ir::{Callee, CallingConvention, DataLayout, FunctionSignature, Linkage, ModuleBuilder, Type};

fn main() {
    let mut module = ModuleBuilder::new("x86_64-whale-linux", DataLayout::default_64bit_le());
    let signature = FunctionSignature {
        params: vec![Type::I32], ret: Type::I32,
        convention: CallingConvention::SysV64, variadic: false,
    };
    let identity = module.declare_function(
        "identity", signature, Linkage::External, Some("identity_i32".into()),
    ).unwrap();
    let mut function = module.begin_function("answer", vec![], Type::I32);
    let input = function.const_i32(42);
    let callback = function.function_addr(identity).unwrap();
    // The direct call's result remains defined even though it is unused.
    function.call(Callee::Direct(identity), vec![input]).unwrap();
    let result = function.call(Callee::Indirect(callback), vec![input]).unwrap().unwrap();
    function.ret(Some(result));
    function.finish();
    let module = module.finish();
    ir::verify_module(&module).unwrap();
    print!("{}", ir::print_module(&module));
}
```

```text
module {
  format_version 2
  semantics_version 1
  target "x86_64-whale-linux"
  datalayout { ptr=64, endian=little }

  declare @f0 "identity": sysv64 (i32) -> i32, linkage external, link_name "identity_i32"
  declare @f1 "answer": whale () -> i32, linkage internal

  fn @answer() -> i32, id @f1 {
  entry:
    %v0: i32 = const i32 42
    %v1: fnptr<sysv64 (i32) -> i32> = function_addr @f0
    %v2: i32 = call sysv64 i32 @f0(%v0)
    %v3: i32 = call sysv64 i32 indirect %v1(%v0)
    ret i32 %v3
  }

}
```

`Callee::Direct(FunctionId)` diselesaikan melalui tabel deklarasi; `Callee::Indirect(ValueId)` memerlukan nilai `Type::FnPtr(FunctionSignature)`. Tanda tangan mencakup semua tipe parameter, tipe hasil, dan `CallingConvention::{Whale, SysV64}`. Itu disimpan melalui salinan, penyimpanan, parameter, pengembalian, phi dan pilih. Pointer data dan nilai integer tidak dapat dipanggil. Pemeran yang melibatkan tipe penunjuk fungsi ditolak; mengubah anotasi tipe tidak dapat mengubah tanda tangan yang dapat dipanggil. Penunjuk fungsi memiliki penyimpanan alamat 64-bit pada target ini; ini sendiri tidak mengimplementasikan metadata bayangan runtime.

Arity, tipe argumen/hasil yang tepat, keberadaan ID hasil, dan konvensi pemanggilan harus cocok. Tidak ada konversi implisit. Pihak yang dipanggil tidak langsung harus mendominasi panggilan tersebut seperti halnya argumennya. `variadic: true`, parameter bertipe void dan SysV64 parameter agregat/tanda tangan hasil ditolak. Whale tanda tangan agregat dapat direpresentasikan dalam IR; klasifikasi ABI asli dan emisi panggilan mesin belum tersedia untuk kedua konvensi tersebut.

`null_function(signature)` mewakili penunjuk fungsi nol yang diketik. Memanggilnya diketik dengan baik IR dengan runtime trap yang diperlukan sebelum memasukkan callee. Target nonnull yang tidak valid, kedaluwarsa, atau tidak kompatibel dengan tanda tangan yang dicentang juga harus dijebak. Pemeriksaan runtime dan manajemen seumur hidup panggilan balik asing ini menunggu lapisan penerjemah/eksekusi asli; keberhasilan pemverifikasi tidak berarti alamat eksternal yang sewenang-wenang aman.

### Bentuk pemanggilan AST

Ini adalah fragmen ekspresi di dalam program format 2 AST:

```json
{"Call":{"callee":{"Direct":"increment"},"args":[{"Lit":{"Int":{"bits":32,"signed":true,"value":"41"}}}]}}
```

```json
{"Call":{"callee":{"Indirect":{"FunctionRef":"increment"}},"args":[{"Lit":{"Int":{"bits":32,"signed":true,"value":"41"}}}]}}
```

`Direct` dan `FunctionRef` menggunakan namespace fungsi meskipun variabel memiliki nama yang sama. `Indirect` mengevaluasi ekspresinya terlebih dahulu, lalu mengevaluasi argumen dari kiri ke kanan. Panggilan batal valid sebagai `ExprStmt`, tetapi tidak sebagai penginisialisasi variabel, argumen, operan, atau nilai yang dikembalikan. Panggilan dan referensi fungsi bukanlah ekspresi konstanta numerik pada waktu kompilasi. `NullFunction` mengambil objek tanda tangan dengan bidang `params`, `ret`, `convention` dan `variadic`.

[Contoh JSON lengkap](https://github.com/wavefnd/Whale/blob/master/ir/tests/fixtures/ast-v2-calls.json) menyimpan panggilan balik dan memanggilnya sebelum panggilan eksternal. Turunkan dengan:

```sh
cargo run --locked --features socket-cli -- ir lower ir/tests/fixtures/ast-v2-calls.json
```

[yang diharapkan IR](https://github.com/wavefnd/Whale/blob/master/ir/tests/fixtures/calls-v2.wir) diperiksa dalam pengujian penurunan. Identitas fungsi dan nama tautan direpresentasikan pada batas IR; melestarikannya melalui pembuatan objek asli dan menghubungkannya masih merupakan pekerjaan terpisah.

## Ketersediaan blok dan nilai

Setiap blok memiliki pengidentifikasi unik dan tepat satu terminator. Target cabang harus memiliki fungsi yang sama. Blok entri harus ada dan tidak boleh memiliki tepi depan dan phi. Saat membuat loop, lakukan cabang dari blok entri ke header loop terpisah.

Definisi nilai di jalur yang dapat dieksekusi harus mengatur titik penggunaannya. Artinya semua jalur dari titik masuk ke titik penggunaan harus melewati definisi tersebut. Di blok yang sama, definisi harus didahulukan sebelum penggunaan. Urutan penyimpanan blok tidak menentukan dominasi.

Asumsikan titik masuk bercabang ke left atau right lalu bergabung di join. Nilai yang ditentukan hanya di left tidak dapat digunakan sebagai nilai umum di join. Hal ini karena jalur yang melewati right tidak terdefinisi. Nilai dari setiap blok sebelumnya harus digabungkan menjadi phi, yang diterima sebagai masukan.

Blok yang tidak dapat dijangkau juga disimpan dalam modul. Verifikator terus-menerus memeriksa pengidentifikasi, jenis, operan, dan struktur cabang blok. Definisi blok yang tidak dapat dijangkau tidak dapat memberikan nilai pada penggunaan normal jalur yang dapat dijangkau.

## Perintah Phi

phi ditempatkan sebelum semua perintah reguler di blok. Tepat satu masukan diperlukan untuk setiap blok sebelumnya yang berbeda. Nilai input harus bertipe phi dan harus tersedia di akhir blok sebelumnya yang sesuai.

Bahkan jika ada beberapa sisi dalam satu blok sebelumnya, hanya ada satu masukan. Loop phi dapat merujuk pada nilai blok yang muncul kemudian dalam urutan penyimpanan modul selama itu adalah nilai yang dihitung pada tepi berulang. Blok sebelumnya yang hilang, terduplikasi, atau tidak relevan dan input yang salah diketik merupakan kesalahan validasi.

## Evaluasi dan Seleksi

Whale AST mengevaluasi target panggilan dan subekspresi dari kiri di tempat yang ditentukan. Frontend mengungkapkan evaluasi hubung singkat sebagai cabang aliran kontrol.

Select memilih salah satu nilai yang sudah dihitung. Itu tidak menghilangkan perhitungan input mana pun. Misalnya, meskipun Anda memilih nilai aman, Anda tidak dapat menghindari trap ditemui saat menghitung input lainnya. Perhitungan yang perlu dijalankan hanya pada jalur tertentu harus ditempatkan di dalam blok bersyarat.

## Departemen Verifikasi trap

IR yang tidak valid akan ditolak pada tahap verifikasi. Pelanggaran kondisi eksekusi ditangani dengan trap yang ditentukan, dan `undef` dan `poison` bukanlah nilai yang dapat diterima. Penyalahgunaan builder, definisi duplikat, penambahan terminator kedua harus dikembalikan sebagai kesalahan struktural.

trap berisi alasan·lokasi sumber·IR ID dan kemudian menghentikan eksekusi. Jalankan native menghentikan program, dan juru bahasa API mengembalikan kesalahan Trap. Mempertahankan efek samping sebelumnya, tetapi tidak menjamin buffer flush·panggilan destruktor·stack unwinding.

Jaminan ini berlaku untuk IR terverifikasi dan memori penelusuran. Eksternal C·alamat mentah·perakitan inline memiliki kontrak terpisah dan tidak selalu mendeteksi pelanggaran di luar batasannya. Silakan merujuk ke [Model memori](memory-model).

## Format Pertukaran dan Representasi Teks

AST dan typed IR masing-masing menggunakan format version dan semantics version yang umum. Pembaca harus menolak kunci JSON yang tidak berversi/tidak diketahui versi/bidang/fungsi/duplikat. Konstruktor tidak boleh berasumsi bahwa properti yang tidak didukung akan diabaikan secara diam-diam.

Bilangan bulat diteruskan sebagai nomor string lebar bit·signedness·. Konstanta floating point dilewatkan sebagai string bit lebar dan tepat. Teks round-trip di IR harus mempertahankan nama·ID·tipe·konstan·urutan·properti·metadata. Spasi dan penempatan komentar tidak dapat dipertahankan.

Anda dapat menggunakan kontrak skalar AST JSON di bawah. Keluaran typed IR menyertakan informasi versi, namun pertukaran bolak-balik penuh dengan parser teks belum didukung.


### Versi yang ditentukan AST JSON

Simpan yang berikut ini sebagai `program.json`. Keempat bidang amplop wajib diisi. `program` berisi array `declarations`, `globals` dan `functions` yang diperlukan, yang mungkin kosong. Nama fungsi, parameter, tipe kembalian, isi, `convention` dan `linkage` wajib diisi. `link_name` mungkin tidak ada/null untuk fungsi internal dan harus berupa string yang tidak kosong tanpa NUL untuk fungsi eksternal. Setiap enum menggunakan nama unitnya atau objek kunci varian tunggal. Varian unit juga menerima objek bernilai null, seperti `{"Void":null}`; encoder mengeluarkan nama unit `"Void"`. `VarDecl.init` mungkin tidak ada atau nol; bidang wajib lainnya harus ada.

```json
{
  "format_version": 2,
  "semantics_version": 1,
  "features": [],
  "program": {
    "globals": [],
    "functions": [
      {
        "name": "answer",
        "parameters": [],
        "return_type": {
          "Int": {
            "bits": 128,
            "signed": false
          }
        },
        "body": [
          {
            "Return": {
              "Lit": {
                "Int": {
                  "bits": 128,
                  "signed": false,
                  "value": "340282366920938463463374607431768211455"
                }
              }
            }
          }
        ],
        "convention": "Whale",
        "linkage": "Internal",
        "link_name": null
      }
    ],
    "declarations": []
  }
}
```

```sh
cargo run --locked --features socket-cli -- ir lower program.json
```

```text
module {
  format_version 2
  semantics_version 1
  target "x86_64-whale-linux"
  datalayout { ptr=64, endian=little }

  declare @f0 "answer": whale () -> u128, linkage internal

  fn @answer() -> u128, id @f0 {
  entry:
    %v0: u128 = const u128 340282366920938463463374607431768211455
    ret u128 %v0
  }

}
```

Bilangan bulat `value` adalah string desimal. Signed Angka digunakan setelah minus opsional dari bilangan bulat, dan spasi, plus, eksponen, dan pemisah tidak diperbolehkan. Kisaran yang diperbolehkan ditentukan oleh lebar yang dinyatakan dan signedness. `u128::MAX` di atas dipertahankan sebagaimana adanya melalui JSON dan lowering. Nilai unsigned negatif atau nilai di luar rentang bukan wrap dan merupakan kesalahan. Nilai Float menggunakan string bit heksadesimal dengan lebar persis seperti yang dijelaskan dalam [Operasi numerik](numeric-operations).

`format_version` adalah 2 untuk format AST ini; `semantics_version` adalah 1. `features` harus berupa array kosong. Bidang, versi, fitur, duplikat kunci JSON mentah yang tidak diketahui (termasuk kunci setara yang di-escape), dan nilai akhir merupakan kesalahan, bahkan dengan `--no-verify`. Titik masuk perpustakaan adalah `ir::lower_ast::interchange::decode`; `encode` mengeluarkan amplop. `decode` defaultnya adalah batas byte sumber 8 MiB; `decode_with_limit` menerima batas penelepon. JSON sarang dibatasi. Gunakan dekoder mentah ini daripada menguraikannya ke dalam peta umum yang sudah dapat membuang kunci duplikat.

[Skema JSON lengkap](https://github.com/wavefnd/Whale/blob/master/ir/schema/ast-v2.schema.json) menentukan bentuk, bidang yang wajib diisi, dan varian. Pemeriksaan rentang/jenis dan deteksi kunci duplikat juga berlaku. Subset penurun skalar mencakup literal, variabel/konstanta, tambah/sub/mul, perbandingan, penugasan, jika/sementara, kembali dan putus/lanjutkan. Referensi fungsi, panggilan langsung dan panggilan tidak langsung didukung; ekspresi agregat tidak didukung. `Opaque` dapat diwakili dalam skema tetapi tidak didukung oleh penurunan.

Migrasi memerlukan pembungkusan muatan Program lama yang kosong dan mengganti literal numerik JSON dengan string bilangan bulat desimal atau string bit float. Payload lama yang tidak berversi ditolak. Muatan format 1 harus dimigrasikan ke format 2: tambahkan `program.declarations` (array kosong jika tidak digunakan) dan eksplisit `convention`/`linkage` pada definisi. AST dan nomor versi IR yang diketik bersifat independen; keduanya sekarang menjadi 2, dengan semantik versi 1.

### Masukan yang ditolak dan pemulihan CLI

Simpan masukan lengkap berikut sebagai `invalid.json`.

```json
{
  "format_version": 99,
  "semantics_version": 1,
  "features": [],
  "program": {
    "globals": [],
    "functions": [],
    "declarations": []
  }
}
```

```sh
cargo run --locked --features socket-cli -- ir lower invalid.json -o rejected.wir
```

```text
Failed to parse socket JSON: unsupported AST format_version 99; expected 2
```

Perintah keluar dengan status bukan nol dan tidak membuat keluaran baru atau menimpa file yang sudah ada. Jenis ketidakcocokan juga akan gagal sebelum menerbitkan keluaran. Biner dibuat tanpa keluar `socket-cli` dengan status 2 dan mengeluarkan perintah pemulihan yang berisi `--features socket-cli`.
