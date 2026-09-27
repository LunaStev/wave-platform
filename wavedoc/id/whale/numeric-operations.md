---
translation_set_id: whale-numeric-operations
path: whale/numeric-operations
locale: id
group: whale
group_order: 1
order: 7
title: Operasi numerik
summary: Menjelaskan operasi integer wrap, checked, shift, kesalahan konversi tipe, dan hasil floating point.
---

## representasi bilangan bulat

Bit integer N memiliki bit nilai N. Bilangan bulat tak bertanda berkisar dari 0 hingga 2^N − 1, dan bilangan bulat bertanda berkisar dari −2^(N−1) hingga 2^(N−1) − 1. Operasi signedness menentukan interpretasi string bit.

Tabel di bawah ini menjelaskan hasil perhitungannya. IR Contoh kode menggunakan representasi printer saat ini.

## Penjumlahan, pengurangan, perkalian

Bilangan bulat dasar add·sub·mul menampung bit hasil N yang lebih rendah. overflow tidak menyebabkan trap. Operasi checked mengembalikan hasil wrap yang sama, serta Bool, yang menunjukkan apakah hasil matematika berada di luar rentang signed atau unsigned yang sesuai.

|operasi|Wrap Hasil| Checked overflow |
| --- | --- | --- |
| u8: 255 + 1 | 0 | true |
| i8: 127 + 1 | −128 | true |
| u8: 0 − 1 | 255 | true |
| i8: 12 × 3 | 36 | false |

Frontend untuk bahasa yang menghentikan eksekusi di overflow harus menggunakan `trap_if` eksplisit untuk overflow hasil operasi checked. Operasi default tidak secara implisit menerapkan kebijakan overflow bahasa sumber.

### Wrap dan IR menyatakan pemeriksaan eksplisit

Modul berikut dikonfigurasi sebagai Rust builder dan merupakan keluaran printer saat ini yang telah lolos verifikasi. Pengurai teks dan backend eksekusi belum tersedia.

```text
module {
  format_version 2
  semantics_version 1
  target "x86_64-whale-linux"
  datalayout { ptr=64, endian=little }

  declare @f0 "add_u8": whale () -> u8, linkage internal
  declare @f1 "require_no_overflow": whale () -> u8, linkage internal

  fn @add_u8() -> u8, id @f0 {
  entry:
    %v0: u8 = const u8 255
    %v1: u8 = const u8 1
    %v2: u8 = add u8 %v0, %v1
    ret u8 %v2
  }

  fn @require_no_overflow() -> u8, id @f1 {
  entry:
    %v3: u8 = const u8 255
    %v4: u8 = const u8 1
    %v5: tuple<u8, bool> = uadd_chk u8 %v3, %v4
    %v6: u8 = extract %v5, 0
    %v7: bool = extract %v5, 1
    trap_if bool %v7, reason="integer overflow"
    ret u8 %v6
  }

}
```

Berdasarkan kontrak aritmatika, `add_u8` membungkus 256 ke hasil 8-bit 0. `require_no_overflow` mengekstrak hasil yang dibungkus dengan `extract ..., 0` dan tanda luapan Bool dengan `extract ..., 1`. Jika benderanya adalah true, `trap_if` menghentikan eksekusi sebelum pengembalian. Ini adalah hasil yang ditentukan, bukan keluaran dari penerjemah yang diimplementasikan.

## Pembagian dan sisanya

Pembagian bilangan bulat atau sisa dengan nol menyebabkan jebakan. Membagi nilai minimum yang ditandatangani dengan −1 akan membungkus nilai minimum. Sisanya dalam hal ini adalah nol.

|operasi|hasil|
| --- | --- |
| i8: −128 / −1 | −128 |
| i8: −128 % −1 | 0 |
|Bagilah dengan bilangan bulat 0| trap |
|sisa bilangan bulat 0| trap |

## bergeser

Dalam menggeser nilai bit N, string bit count ditafsirkan sebagai unsigned dan sisanya dibagi dengan N digunakan. Itu tidak menghasilkan trap hanya karena count berada di luar rentang dari 0 hingga N−1.

Dalam nilai 8-bit, count 0·8·16 semuanya merupakan pergeseran 0-bit. String bit `11111111` dari 8 bit count digeser sebanyak 7 bit. Hal ini tetap sama meskipun string bit ini mewakili signed −1. unsigned Analisis dilakukan sebelum penghitungan lainnya.

Bahasa sumber yang menolak count negatif atau berlebihan harus memberi tanda centang eksplisit sebelum pergeseran.

## konversi tipe

|konversi|artinya|
| --- | --- |
| Zero extension |Tingkatkan lebar dengan mengisi bit tingkat tinggi dengan 0|
| Sign extension |Tingkatkan lebarnya dengan menduplikasi bit tanda|
|sedikit pemotongan|Simpan hanya bit tingkat rendah yang sesuai dengan lebar tujuan|
|Kalahkan penafsiran ulang|Menafsirkan string bit yang sama dengan tipe yang berbeda|
|Konversi numerik tanpa kerugian|Jika tidak dapat diungkapkan sambil mempertahankan nilai numeriknya, trap|

Misalnya, mengubah 8-bit `11111111` menjadi 16-bit zero extension menjadi `0000000011111111`, dan sign extension menjadi `1111111111111111`. Sekalipun bit inputnya sama, operasinya berbeda.

konversi float-to-int terpotong menuju nol, lalu periksa rentang bilangan bulat. NaN dan ketidakterbatasan menyebabkan jebakan. Saat mengonversi ke i8, 127,9 menjadi 127, sedangkan 128,0 menjebak. Bool dikonversi ke bilangan bulat 0 atau 1. i1 yang ditandatangani tidak dapat mewakili 1, sehingga tidak dapat menjadi tujuan konversi ini.

Mengonversi alamat menjadi bilangan bulat tidak memulihkan akses penunjuk ke bilangan bulat tersebut. Silakan merujuk ke [validitas penunjuk](memory-model).

## aritmatika titik mengambang

Nilai floating point memiliki string bit yang tepat f16·f32·f64. Operasi ini membulatkan ke nilai terdekat pada lebar yang dinyatakan dan, jika tepat di tengah, menggunakan ties-to-even, yang memilih nilai dengan bit paling signifikan dari digit signifikan.

Operasi defaultnya adalah fast-math, implisit FMA, yang tidak mengizinkan penanganan nol paksa pada nilai kecil. Perkalian yang diikuti penjumlahan mempertahankan langkah pembulatan masing-masing dan backend tidak boleh menggabungkannya secara implisit ke dalam satu operasi.

Hasil operasi numerik bisa berupa NaN atau tak terhingga. NaN dari operasi matematika dan perubahan lebar dinormalisasi menjadi satu jumlah tetap NaN tenang per lebar. Simpan dan salin, pertahankan bit NaN asli. Oleh karena itu, perilakunya berbeda saat meneruskan payload NaN ke memori tanpa operasi aritmatika dan saat menghitungnya.

Itu tidak mengekspos tanda status floating-point. Meskipun aritmatika titik-mengambang dasar memungkinkan hasil NaN·tak terhingga, konversi float→int menerapkan aturan trap di atas.


### Simpan konstanta yang tepat

Gunakan `FloatBits` variant atau string bit heksadesimal apa pun dengan lebar yang tepat. Kesetaraan nilai yang disimpan dibandingkan dengan string bit yang berisi negatif 0 dan NaN payload. Verifikator menolak jika lebar tipe IR dan payload berbeda.

```rust
use ir::{FloatBits, ModuleBuilder, Target, Type};
fn main() {
    let bits = FloatBits::parse(32, "0xffc01234").unwrap();
    assert_eq!(bits, FloatBits::F32(0xffc01234));
    let target = Target::X86_64WhaleLinux;
    let mut module = ModuleBuilder::new(target.name(), target.data_layout());
    let mut function = module.begin_function("payload", vec![], Type::F32);
    let value = function.const_float_bits(Type::F32, bits);
    function.ret(Some(value));
    function.finish();
    let module = module.finish();
    ir::verify_module(&module).unwrap();
    assert!(ir::print_module(&module).contains("const f32 0xffc01234"));
    println!("{}", bits);
}
```

```text
0xffc01234
```

String bit f16, f32, atau f64 dimulai dengan `0x`, diikuti dengan tepat 4, 8, atau 16 digit heksadesimal. `0x80000000` mewakili f32 nol negatif; `0x7f800000` mewakili ketidakterbatasan positif. `const_float` mengonversi nilai host f64 secara numerik; gunakan `const_float_bits` untuk mempertahankan bit aslinya. Representasi penyimpanan yang tepat tidak berarti backend eksekusi floating-point yang lengkap. Aritmatika waktu kompilasi masih menggunakan perantara host f64, sehingga kontrak pembulatan penuh untuk setiap lebar yang dinyatakan belum diterapkan.
