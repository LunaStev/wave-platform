---
translation_set_id: diagnostics
path: reference/diagnostics
locale: id
group: reference
group_order: 5
order: 2
title: Pemecahan masalah: Dari instalasi hingga eksekusi
summary: Pisahkan langkah-langkah yang gagal dan persempit penyebabnya dengan informasi yang dapat direproduksi.
---

## Pertama, bedakan tahapan kegagalan

|fenomena yang diamati|Periksa dulu|tindakan selanjutnya|
| --- | --- | --- |
|wavec Perintah tidak ditemukan|PATH dan lokasi file yang dapat dieksekusi|Jalankan dengan jalur absolut dan atur PATH|
|Tidak dapat menemukan file yang diperlukan untuk dijalankan|Apakah ada file yang hilang dari folder instalasi?|Buka paket dan instal kembali seluruh paket|
|std import Gagal| `wavec print std-path` |Korespondensi: Instal std atau tentukan `--std-root`|
|Lokasi sumber dan jenis keluaran kesalahan| `wavec check main.wave` |Perbaiki kesalahan pertama dan periksa lagi|
|Kegagalan pembangunan untuk target OS·CPU lainnya|Ditentukan target dan lingkungan target|[Pengaturan lintas build](/docs/id/whale/build-link-targets) Konfirmasi|
|Eksekusi gagal setelah build berhasil|kode keluar, input, direktori kerja|Lingkungan eksekusi dan pemeriksaan kesalahan API|

## Contoh diagnostik kecil

Berikut keseluruhan programnya, yang sengaja dibuat salah:

```wave
fun main() {
    var count: i32 = 1;
    println("{}", missing);
}
```

`wavec check main.wave` harus menunjuk ke nama yang tidak dideklarasikan missing. Ubah nama variabel menjadi count, lalu periksa dan jalankan kembali. Fokus pada file, lokasi, dan penyebab daripada keseluruhan teks diagnostik. Kesalahan selanjutnya mungkin disebabkan oleh kesalahan awal.

## Ketika sebuah executable gagal

Di shell Linux/macOS, kode keluar diperiksa segera setelah eksekusi sebagai `echo $?`, dan di PowerShell, adalah `$LASTEXITCODE`. Kesalahan input dan `return 1` eksplisit bukan penyebab yang sama. Nilai runtime yang tidak valid untuk jumlah shift atau konversi nyata dapat menyebabkan trap. Lihat [aturan operasi](/docs/id/language/expressions-and-operators).

Jalur file relatif dipengaruhi oleh direktori kerja yang dapat dieksekusi, bukan lokasi file sumber. Jangan perlakukan kegagalan pembacaan file sebagai panjang string 0, periksa kesalahan pengembalian terlebih dahulu. Kegagalan koneksi jaringan diperiksa melalui pencarian alamat, menunggu server, izin, dan batas waktu.

## Informasi diperlukan untuk melaporkan masalah

1. `wavec --version` Output dan perintah yang tepat dijalankan.
2. target. ditentukan secara terpisah dari host OS·arsitektur
3. Sumber kompiler digunakan dengan jalur std yang dipilih.
4. Sumber minimal, masukan, dan file yang diperlukan untuk mereproduksi masalah.
5. Hasil yang diharapkan, hasil aktual, diagnosis dan kode keluar.

Kata sandi, token, dan konten file pribadi dihapus. Jika masalah hilang ketika Anda mengurangi contoh minimal, elemen terakhir yang Anda hilangkan adalah petunjuknya. `--error-format=json` tersedia saat alat mengumpulkan diagnostik.

[Instalasi](/docs/id/getting-started/install) · [perintah kompiler](/docs/id/getting-started/compiler) · [Target dan Tautan](/docs/id/whale/build-link-targets)
