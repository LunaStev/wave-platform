---
translation_set_id: stdlib-files-io
path: stdlib/files-io
locale: id
group: stdlib
group_order: 1
order: 7
title: fs dan io: Transfer file dan byte
summary: Menjelaskan masa pakai file, pembacaan penuh, transfer sebagian, dan status pasca-kegagalan.
---

## Pilih file API

Fungsi kenyamanan `std::fs::file` menerima jalur dan melakukan pembukaan dan penutupan yang diperlukan. Fungsi yang mengembalikan deskriptor harus ditutup oleh pemanggil.

|deklarasi|Hasil sukses dan tindakan pencegahan|
| --- | --- |
| `open_read(path: str) -> i64` |Buka deskriptor. Angka negatif adalah kesalahan|
| `create(path: str) -> i64` |Membuat atau menghapus konten file yang ada. Mengembalikan deskriptor pemilik|
| `open_append(path: str) -> i64` |Buka atau buat untuk penambahan|
| `size(path: str) -> i64` |Jumlah byte. Angka negatif adalah kesalahan|
| `read_into(path: str, dst: ptr<u8>, dst_cap: i64) -> i64` |Jumlah byte di seluruh file. Kurangnya kapasitas adalah sebuah kesalahan|
| `read_to_end(path: str, dst_buffer: ptr<Buffer>) -> i64` |Tambahkan file setelah Buffer yang ada dan kembalikan jumlah tambahan|
| `write(path: str, src: ptr<u8>, len: i64) -> i64` |Jumlah byte yang ditulis menggantikan konten yang ada|
| `append(path: str, src: ptr<u8>, len: i64) -> i64` |Jumlah byte ditambahkan ke akhir|
| `remove(path: str) -> i64` |status penghapusan. kegagalan itu negatif|

false dari `exists(path)` saja tidak dapat membedakan antara file yang hilang dan kesalahan izin. Pastikan untuk memeriksa hasil pembukaan sebenarnya, karena status dapat berubah antara pemeriksaan keberadaan dan pembukaan.

## Tingkat rendah I/O

```text
std::io::fd
io_read(fd: i64, buf: ptr<u8>, len: i64) -> i64
io_write(fd: i64, buf: ptr<u8>, len: i64) -> i64
io_read_exact(fd: i64, buf: ptr<u8>, len: i64) -> i64
io_write_all(fd: i64, buf: ptr<u8>, len: i64) -> i64
io_close(fd: i64) -> i64
```

Hasil positif dari `io_read` adalah jumlah byte yang dibaca, dan 0 dalam permintaan panjang positif adalah EOF. `io_write` boleh ditulis kurang dari yang diminta. Jika transfer penuh diperlukan, gunakan fungsi exact/all. Namun, kami tidak berasumsi bahwa kegagalan tersebut akan mengembalikan keadaan eksternal, karena beberapa transfer mungkin telah terjadi sebelum kesalahan tersebut.

`io_read_exact` adalah kesalahan jika EOF ditemukan sebelum panjang yang dibutuhkan. `read_into` mengembalikan `IO_ERR_NO_SPACE` jika buffer penuh, dan beberapa byte mungkin sudah ditulis. Fungsi baca tidak secara otomatis menambahkan NUL ke akhir string.

## Buffer dan penanganan kesalahan

Jika gagal, `read_to_end` memulihkan lensa asli, namun kapasitas dan alamat datanya mungkin telah berubah. Penelepon harus membebaskan Buffer setelah berhasil atau gagal. API penulisan file tidak menjamin penggantian file atom.

Dari [Latihan membaca file](/docs/id/practice/file-reader), Anda dapat menjalankan program dari import hingga rilis. Pertimbangkan perbedaan jalur/izin di Linux/macOS/Windows/FreeBSD dan batasan direktori yang dapat diakses di WASI. Itu tidak secara langsung menafsirkan nilai deskriptor sebagai pegangan mentah dari OS lainnya.

## Membaca file besar ke dalam buffer kecil

Operasi yang tidak memerlukan seluruh file ditempatkan di memori dapat ditangani dengan buffer tetap dan iterasi baca. Program berikut mencetak konten input.txt dan menghitung jumlah total byte yang dibaca. Simpan satu `Wave` dan LF di file masukan.

<!-- wave-example: book-file-stream -->
```wave
import("std::fs::file")::{open_read};
import("std::io::fd")::{io_read, io_write_all, io_close};
import("std::io::consts")::{IO_STDOUT_FD};

fun main() -> i32 {
    var descriptor: i64 = open_read("input.txt");

    if (descriptor < 0) {
        return 1;
    }

    var buffer: array<u8, 4>;
    var total: i64 = 0;

    while (true) {
        var count: i64 = io_read(descriptor, &buffer[0], 4);

        if (count < 0) {
            io_close(descriptor);
            return 2;
        }

        if (count == 0) {
            break;
        }

        if (io_write_all(IO_STDOUT_FD, &buffer[0], count) < 0) {
            io_close(descriptor);
            return 3;
        }

        total += count;
    }

    if (io_close(descriptor) < 0) {
        return 4;
    }

    println("bytes={}", total);
    return 0;
}
```

Hasil eksekusi:

```text
Wave
bytes=5
```

Kapasitas buffer adalah 4, namun pembacaan terakhir mungkin 1 byte. Kami selalu meneruskan count aktual ke output. Jika Anda menulis seluruh array, bahkan byte lama yang belum dibaca pun dapat dihasilkan.

Program membuka descriptor dari input.txt, jadi tutuplah. Output standar bukanlah sumber daya yang baru diperoleh dalam fungsi ini, sehingga tidak ditutup secara sembarangan di akhir contoh.

## Bacaan penuh dan kapasitas tidak mencukupi

read_into menerima ruang penyimpanan tetap yang menampung seluruh file. Jika ruang tidak mencukupi, ia akan terpotong secara diam-diam dan mengembalikan NO_SPACE tanpa hasil.

<!-- wave-example: book-file-capacity -->
```wave
import("std::fs::file")::{read_into};
import("std::io::consts")::{IO_ERR_NO_SPACE};

fun main() {
    var data: array<u8, 2>;
    var status: i64 = read_into("input.txt", &data[0], 2);

    if (status == IO_ERR_NO_SPACE) {
        println("destination too small");
    }
}
```

Hasil eksekusi:

```text
destination too small
```

Itu tidak berasumsi bahwa pembacaan yang gagal tidak mengubah byte tujuan sama sekali. Jangan menggunakannya sebagai konten file yang sudah jadi, siapkan repositori yang lebih besar atau pilih metode streaming. Meskipun Anda menanyakan ukurannya terlebih dahulu, hasil pembacaan sebenarnya adalah penilaian akhir, karena file dapat berubah antara kueri dan pembacaan.

## Perbedaan antara menulis dan menambahkan file

write menggantikan konten yang ada dan append ditambahkan di akhir. Jumlah byte yang akan disimpan dapat diperoleh langsung dari panjang string dan diteruskan. Tanda NUL di akhir string biasanya tidak disertakan dalam konten file teks.

<!-- wave-example: book-file-write-append -->
```wave
import("std::fs::file")::{write, append, size, remove};
import("std::string::len")::{len};

fun main() -> i32 {
    var first: str = "Wave";
    var second: str = " study";

    if (write("output.txt", first as ptr<u8>, len(first) as i64) < 0) {
        return 1;
    }

    if (append("output.txt", second as ptr<u8>, len(second) as i64) < 0) {
        return 2;
    }

    println("bytes={}", size("output.txt"));

    if (remove("output.txt") < 0) {
        return 3;
    }

    return 0;
}
```

Hasil eksekusi:

```text
bytes=10
```

Contoh ini membuat, mengganti, dan akhirnya menghapus output.txt di direktori kerja. Jalankan dari direktori latihan tanpa file yang ada. Editor atau program penyimpanan sebenarnya mungkin memerlukan kebijakan penyimpanan terpisah, seperti file sementara dan penggantian.

## API Tabel pemilihan

|situasi|pilih|
| --- | --- |
|Baca seluruh file kecil ke dalam buffer tetap| read_into |
|Simpan seluruh isinya tanpa mengetahui ukurannya|read_to_end dan Buffer|
|Memproses konten secara berurutan daripada menyimpannya secara keseluruhan|open_read + io_read ulangi|
|Baca catatan dengan panjang tetap| io_read_exact |
|Mengirimkan seluruh string byte| io_write_all |
|Menangani file yang sudah terbuka|Fungsi fd alih-alih fungsi jalur|

Setelah memilih fungsi, periksa bagaimana buffer, lokasi file, dan data eksternal berubah jika terjadi kegagalan.
