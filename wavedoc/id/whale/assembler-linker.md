---
translation_set_id: whale-assembler-linker
path: whale/assembler-linker
locale: id
group: whale
group_order: 1
order: 11
title: Perakitan kode dan penautan statis
summary: Menjelaskan pengkodean operan, penempatan bagian, pengikatan simbol, dan titik masuk eksekusi.
---

## Majelis dan Objek

Assembler Whale mengubah instruksi AMD64 menjadi byte kode mesin dan informasi relokasi. Objek ELF64 berisi informasi dan bagian/simbol ini. Perakitan dilakukan oleh implementasi Whale sendiri dan tidak memerlukan assembler eksternal.

Objek yang dapat direlokasi mungkin memiliki referensi yang alamat akhirnya belum diketahui. Proses penyelesaian alamat ini adalah tautan. Perakitan yang berhasil tidak boleh menyimpulkan bahwa semua simbol eksternal telah diselesaikan atau bahwa executable telah dibuat.

## Kumpulkan fungsinya

Simpan kode berikut sebagai `answer.asm`.

```asm
section .text
global answer

answer:
    mov eax, 42
    ret
```

```sh
whale asm --amd64 answer.asm -o answer.o
```

Isi `.text` adalah `b8 2a 00 00 00 c3`, dan `mov eax, 42` diikuti oleh `ret`. Objek ELF64 mengekspos `answer`. Ini adalah fungsi yang dapat dipanggil tanpa kode startup proses dan bukan file yang dapat dieksekusi. Instruksi dalam contoh ini dapat diproses oleh assembler saat ini.

## Membangun objek dengan Rust API

Berikut adalah contoh lengkap penulisan byte fungsi yang sama ke dalam peti `object`. Menentukan target keluaran dan simbol global.

```rust
use object::{ObjectFile, ObjectSymbol, SectionKind, SymbolBinding, SymbolVisibility, Target};

fn main() {
    let mut object = ObjectFile::with_target(Target::X86_64WhaleLinux.object_target());
    let text = object.add_section(".text", SectionKind::Text, 16);
    object.sections[text].data = vec![0xb8, 0x2a, 0x00, 0x00, 0x00, 0xc3];
    object.symbols.push(ObjectSymbol {
        name: "answer".into(),
        section_index: Some(text),
        value: 0,
        size: 6,
        binding: SymbolBinding::Global,
        visibility: SymbolVisibility::Default,
    });
    let elf = object.write().unwrap();
    assert_eq!(&elf[..7], b"\x7fELF\x02\x01\x01");
    assert_eq!(u16::from_le_bytes([elf[18], elf[19]]), 62);
    std::fs::write("answer.o", elf).unwrap();
}
```

Assertion memeriksa kelas ELF, urutan byte, pengidentifikasi machine. `value: 0` adalah `.text` di dalam offset dan `size: 6` adalah ukuran byte simbol. Referensi atau rentang bagian yang tidak valid adalah kesalahan serialisasi. Menolak pesanan machine atau byte lainnya tanpa menandainya sebagai AMD64.

## Menafsirkan simbol dua benda

Saat ini peti `linker` menyediakan interpretasi simbol. Contoh yang dapat dieksekusi berikut mendefinisikan `helper` lokal pada dua objek, lalu memeriksa apakah terjadi kesalahan jika nama yang sama diterbitkan dua kali.

```rust
use linker::core::symbol_table::{SymbolKey, SymbolTable};
use object::{
    ObjectFile, ObjectFormat, ObjectSymbol, SectionKind, SymbolBinding, SymbolVisibility,
};

fn input(binding: SymbolBinding) -> ObjectFile {
    let mut object = ObjectFile::new(ObjectFormat::ELF64);
    let text = object.add_section(".text", SectionKind::Text, 1);
    object.sections[text].data = vec![0xc3];
    object.symbols.push(ObjectSymbol {
        name: "helper".into(),
        section_index: Some(text),
        value: 0,
        size: 1,
        binding,
        visibility: SymbolVisibility::Default,
    });
    object
}

fn main() {
    let mut symbols = SymbolTable::new();
    symbols
        .resolve(&[input(SymbolBinding::Local), input(SymbolBinding::Local)])
        .unwrap();
    for object_index in 0..2 {
        let key = SymbolKey::Local {
            object_index,
            name: "helper".into(),
        };
        assert_eq!(symbols.symbols[&key].object_index, Some(object_index));
    }
    let error = symbols
        .resolve(&[input(SymbolBinding::Global), input(SymbolBinding::Global)])
        .unwrap_err();
    assert_eq!(error, "Duplicate global symbol: helper");
    println!("{error}");
}
```

Keluaran:

```text
Duplicate global symbol: helper
```

Kedua definisi lokal memiliki kunci terpisah dengan `object_index` yang berbeda. Kedua definisi global tersebut bertentangan. Contoh ini secara langsung menafsirkan simbol-simbol objek yang telah Anda buat di memori. `.o` Membaca file, menerapkan relokasi, atau mengeluarkan file yang dapat dieksekusi tidak dilakukan. Dari kebijakan seleksi definisi lengkap yang dijelaskan di bawah ini, prioritas weak, dll. belum diterapkan.

## Operan literal dan memori

Literal mempertahankan lebar dan penandaannya hingga rentang pengkodean instruksi sebenarnya diperiksa. Nilai yang berada di luar rentang seharusnya menghasilkan kesalahan, bukan terpotong secara diam-diam.

Notasi ukuran eksplisit diperlukan ketika lebar memori tidak dapat ditentukan dari informasi lain dalam perintah. Misalnya, operan register dapat menentukan lebarnya, tetapi hal ini mungkin menjadi ambigu jika hanya ada memori dan nilai langsung. Assembler tidak boleh menebak lebar yang ambigu secara acak.

Operan memori yang terdiri dari satu simbol di AMD64 pada dasarnya adalah RIP-relative. Pilih metode pengalamatan sebagai eksplisit rel/abs. escape yang tidak diketahui secara literal adalah kesalahan.

## Bagian dan perataan

|Bagian Konten|Perilaku penyelarasan|
| --- | --- |
|kode|NOP Sisipkan perintah|
|data yang diinisialisasi|Masukkan 0 byte|
| BSS |Tingkatkan ukuran memori logis tanpa menambahkan file payload|

Ada perbedaan antara ukuran file dan ukuran memori. BSS mencadangkan memori, tetapi tidak memerlukan ukuran yang sama yaitu 0 byte untuk disimpan dalam file objek. Bagian khusus memiliki properti dan simbol yang menjaga pengikatan dan mengetik informasi.

### Payload Cadangan BSS tanpa alokasi

Simpan yang berikut ini sebagai `buffer.asm`. Logika BSS mencadangkan 1 TiB, tetapi tidak menetapkan atau menulis 1 TiB selama perakitan.

```asm
section .bss
global buffer
buffer:
    resb 1099511627776
buffer_end:

section .text
    ret
```

```sh
whale asm --amd64 buffer.asm -o buffer.o
```

Nilai di dalam bagian `buffer` adalah 0, dan nilai `buffer_end` adalah 1099511627776. Header `.bss` adalah `SHT_NOBITS` dengan ukuran tersebut dan tidak ada file payload. Setelah itu, ketika Anda kembali ke `.bss`, lanjutkan menggunakan logika offset. Arahan data dengan nol juga meningkatkan ukuran logis. Nilai awal bukan nol, relokasi yang akan diterapkan dalam BSS, dan perintah dalam BSS ditolak.

Perbedaan yang sama dapat digunakan untuk objek dan linker API.

```rust
use linker::core::layout::Layout;
use object::{ObjectFile, ObjectFormat, SectionKind};

fn main() {
    let mut object = ObjectFile::new(ObjectFormat::ELF64);
    let text = object.add_section(".text", SectionKind::Text, 16);
    object.sections[text].data = vec![0xc3];
    let bss = object.add_section(".bss", SectionKind::Bss, 16);
    object.sections[bss].zero_fill = 1 << 40;
    assert!(object.sections[bss].data.is_empty());

    let elf = object.write_with_limit(4096).unwrap();
    assert!(elf.len() < 4096);
    let layout = Layout::compute(&[object], 0x1000).unwrap();
    assert_eq!(layout.file_size, 1);
    assert_eq!(layout.memory_size, (1 << 40) + 16);
    assert_eq!(layout.sections[bss].memory_address, 0x1010);
    assert_eq!(layout.sections[bss].file_size, 0);
}
```

```text
section  memory address  file bytes       memory bytes
.text    0x1000          1                1
.bss     0x1010          0                1099511627776
```

`Section::zero_fill` adalah jumlah byte tambahan BSS yang tidak disimpan dalam file. Ukuran memori yang diuji adalah `data.len() + zero_fill`. BSS `data` dengan bantalan nol yang ada juga diterima, namun penetapannya dapat dihindari dengan vektor `data` yang kosong. Bagian selain BSS harus `zero_fill == 0`. Ruang penyimpanan fisik 0 tidak menyiratkan status inisialisasi IR.

`Layout::compute` mengembalikan `Result`, yang mencatat korespondensi objek/bagian masukan, perataan, file offset, alamat memori, dan kedua ukuran semua bagian. Aritmatika alamat dan perataan diperiksa, urutan input dipertahankan, dan perataan objek 0 diperlakukan sebagai tidak dibatasi (perataan 1). BSS tidak menggerakkan kursor file. Ini adalah penerapan payload, dengan header yang dapat dieksekusi, load segment dengan hak akses, menerapkan relokasi adalah tugas terpisah.

ELF writer menolak overflow dan pemotongan saat mengurangi lebar bidang. Nomor bagian yang diperluas tidak didukung, dan jumlah total header termasuk tabel yang dibuat dan bagian penataan ulang harus kurang dari `0xff00`. Batas default untuk keluaran berseri adalah 256 MiB. Anda dapat menentukan batas byte termasuk tabel padding· dengan `ObjectFile::write_with_limit` atau `write_elf_with_limit`, dan batas tersebut tidak termasuk ukuran memori BSS yang tidak disimpan dalam file. Ukuran keluaran diperiksa sebelum mengalokasikan vektor byte terakhir.

## Identifikasi simbol

Fungsi dan variabel menggunakan pengidentifikasi berbeda di dalam IR. Koneksi eksternal menggunakan `link_name` yang ditentukan oleh frontend. Whale tidak secara otomatis mengganti nama salah satu simbol publik yang bertentangan, tetap mempertahankan nama yang disebutkan.

Tabel deklarasi IR mencatat referensi `FunctionId` yang diketik dan nilai fungsi eksplisit `link_name`. API assembler/object/linker di bawah ini tetap merupakan antarmuka yang terpisah: emisi asli IR belum membawa identitas tersebut ke tautan akhir. IR verifikasi panggilan saja tidak menetapkan properti end-to-end ini.

Oleh karena itu, mungkin saja fungsi internal dan variabel memiliki nama `item`, namun akan menjadi kesalahan jika mengekspos keduanya dengan nama eksternal yang sama. Memisahkan namespace internal tidak secara otomatis memisahkan namespace eksternal.

Ruang lingkup simbol lokal suatu objek adalah objek masukannya. Simbol global berpartisipasi dalam interpretasi antar objek. Konflik fungsi/data yang dikonfirmasi adalah kesalahan. Simbol NOTYPE menjaga kompatibilitas dengan input yang tidak menyediakan tipe yang lebih spesifik. Fakta bahwa ia tidak memiliki tipe tidak berarti bahwa ia adalah suatu fungsi atau data.

## pilih definisi

|definisi atau referensi|hasil|
| --- | --- |
|Strong dan strong|kesalahan definisi duplikat|
|Strong dan weak|Strong Pilih definisi|
|Weak dan weak|Pilih definisi pertama dalam urutan input|
|Lihat strong belum terselesaikan|kesalahan tautan|
|Lihat belum terselesaikan weak|Kesalahan tidak didukung di profil statis awal|

Jika terdapat beberapa definisi weak, urutan pemasukannya akan mempengaruhi hasil. Anda harus menggunakan urutan masukan yang diteruskan secara konsisten untuk tautan deterministik.

## Output statis yang dapat dieksekusi

Profil native statis menghasilkan ELF ET_EXEC yang menentukan titik masuk. Titik masuk tidak disimpulkan dari nama fungsi `main`. Itu juga tidak secara otomatis memasukkan kode startup yang memanggil fungsi itu.

Itu tidak secara otomatis menghapus bagian, menggabungkan kode yang identik, atau menghapus simbol. Penempatan file memerlukan penghitungan terpisah atas byte yang sebenarnya Anda simpan dan memori yang Anda cadangan pada waktu proses.

Jalur pembuatan statis lengkap yang dapat dieksekusi untuk CLI belum tersedia. `whale asm` membuat objek yang dapat direlokasi, dan `whale object` membungkus byte mentah menjadi sebuah objek. Silakan lihat [Ikhtisar rantai alat](overview) untuk ketersediaan, ABI dan [AMD64 Sasaran](amd64-target) untuk persyaratan.


## Serialisasi rekaman Wave yang dapat dipilih

Linux Whale build untuk x86_64 sekarang dapat menggunakan implementasi Wave untuk rekaman ELF64 header·section·symbol·RELA yang diperbaiki. Implementasi default Rust juga disediakan. Pemilihan target keluaran, verifikasi objek, penempatan, interpretasi simbol, dan alokasi buffer ditangani oleh Rust. Memilih Wave tidak menambahkan arsitektur yang didukung atau linker lengkap.

Instal Rust, LLVM 21 perpustakaan pengembangan, C linker, dan `ar`, lalu buat jalur yang dipilih dari repositori Whale.

```sh
git clone https://github.com/wavefnd/Wave.git /tmp/whale-wave-bootstrap
git -C /tmp/whale-wave-bootstrap checkout --detach 8a465e30aeea4b817d925cdd0e8d08c1bb029c9a
python3 tools/build_wave_elf.py --wave-source /tmp/whale-wave-bootstrap --out-dir /tmp/whale-wave-elf
WHALE_WAVE_ELF_DIR=/tmp/whale-wave-elf cargo build --locked --all-features
WHALE_WAVE_ELF_DIR=/tmp/whale-wave-elf cargo test --locked --workspace --all-features
```

Skrip memeriksa fiksasi revision, menolak perubahan ke tracked, membuat kompiler Wave, membuat objek Wave dengan LLVM, dan menggabungkannya dengan archive untuk tautan statis. `WHALE_WAVE_ELF_DIR` seharusnya memiliki archive ini. Meminta archive yang tidak valid atau host yang tidak didukung secara eksplisit adalah kesalahan build. Build normal tanpa variabel yang ditentukan tidak memerlukan kompiler `--all-features` atau Wave. Kompiler Wave juga tidak diperlukan saat runtime untuk executable Whale yang tertaut. Jalur untuk melakukan kompilasi silang Whale dengan bootstrap belum didukung.

Misalnya, simpan yang berikut ini sebagai `return.asm` dan kumpulkan ke dalam biner yang dihasilkan:

```asm
section .text
global entry
entry:
    ret
```

```sh
target/debug/whale asm --amd64 return.asm -o return.o
```

```text
ELF class: ELF64
machine: AMD64 (62)
.text bytes: c3
```

Catatan ABI berisi tipe, penunjuk bidang u64, jumlah bidang, penunjuk keluaran, dan kapasitas. Buffer dimiliki oleh pemanggil Rust. Itu tidak melewati kepemilikan penugasan atau ekspresi Rust enum/String/Vec melintasi batas. Rutin Wave memeriksa jumlah, kapasitas, dan lebar lapangan sebelum menulis. Sukses mengembalikan status 0, tipe/penunjuk/kapasitas tidak valid mengembalikan 1, dan bidang overflow mengembalikan 2. Penunjuk harus menunjuk ke buffer yang hidup dan berukuran benar dan tidak tumpang tindih satu sama lain. Pointer mentah C saja tidak dapat membuktikan kondisi ini; wrapper memenuhinya. Tes ini membandingkan seluruh ELF, termasuk BSS dan signed relokasi addend, ke jalur Rust. Ini adalah implementasi parsial Wave yang menggunakan Wave/LLVM bootstrap dan tidak sepenuhnya dihosting sendiri.
