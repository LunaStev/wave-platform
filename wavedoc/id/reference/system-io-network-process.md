---
translation_set_id: system-io
path: reference/system-io-network-process
locale: id
group: stdlib
group_order: 1
order: 15
title: Fungsi dan proses sistem
summary: Menjelaskan batas dan masa proses antarmuka induk API dan OS.
---

## Dokumentasi berdasarkan fungsi

Baca [fs dan io](/docs/id/stdlib/files-io) untuk menangani file, [TCP](/docs/id/stdlib/tcp) untuk menautkan, dan [resolver](/docs/id/stdlib/resolution) untuk mencari alamat. Di bawah ini adalah proses dan aturan akses tingkat bawah OS.

## Dasar-Dasar Proses API

```text
std::process::core
proc_exit(code: i32) -> !
proc_getpid() -> i64
proc_getppid() -> i64
proc_execve(path: str, argv: ptr<ptr<i8>>, envp: ptr<ptr<i8>>) -> i64
proc_waitpid_raw(pid: i64, status: ptr<i32>, options: i32) -> i64
```

`proc_exit` tidak kembali ke titik panggilan. Silakan lakukan pembersihan file/memori yang diperlukan sebelum mematikan. `proc_execve` berbeda dari fungsi pembuatan anak pada umumnya karena, jika berhasil, ini akan menggantikan gambar proses yang ada. raw argv/envp harus dipersiapkan untuk penghentian NUL setiap string, dengan penunjuk null yang menunjukkan akhir.

Fungsi spawn di `std::process::spawn` menangani hasil pembuatan, dan fungsi menunggu menangani status keluar anak. Keberhasilan pembuatan dan keberhasilan penghentian program adalah dua hal yang berbeda. Saat Anda membuat pipa, induk dan anak harus menutup ujung yang tidak terpakai agar EOF dapat dilewati. Jika Anda menunggu anak keluar tanpa membaca pipa penangkap, buffer mungkin terisi dan menunggu satu sama lain.

## Portabilitas dan pendekatan tingkat rendah

fork/exec, deskriptor file, dan pegangan Windows tidak memiliki fungsi OS yang sama. Verifikasi dukungan untuk target yang dipilih dan perlakukan unsupported sebagai jalur kegagalan normal. `std::sys` adalah antarmuka khusus OS dan tidak menggunakan kembali tanda numerik dan tata letaknya dari OS lainnya.

Saat menghubungkan langsung dengan perpustakaan C eksternal, harap baca [FFI](/docs/id/language/modules-imports-and-ffi). Tidak perlu mendeklarasikan fungsi libc secara sewenang-wenang untuk menggunakan induk std API. Periksa [Lingkungan Target dan Tautan](/docs/id/whale/build-link-targets) terlebih dahulu.
