---
translation_set_id: stdlib-tcp
path: stdlib/tcp
locale: id
group: stdlib
group_order: 1
order: 9
title: net.tcp: Koneksi dan transfer
summary: TCP Menjelaskan struktur hasil, tanggung jawab transfer sebagian dan pemutusan hubungan.
---

## Periksa hasil koneksi

```text
std::net::tcp
tcp_connect_addr(addr: SocketAddr) -> NetResult<TcpStream>
tcp_bind_loopback(port: u16) -> NetResult<TcpListener>
tcp_accept(listener: TcpListener) -> NetResult<TcpStream>
tcp_read(stream: TcpStream, buf: ptr<u8>, size: i64) -> i64
tcp_write_all(stream: TcpStream, buf: ptr<u8>, size: i64) -> i64
tcp_read_exact(stream: TcpStream, buf: ptr<u8>, size: i64) -> i64
tcp_close(stream: TcpStream) -> NetError
tcp_close_listener(listener: TcpListener) -> NetError
```

Fungsi tautan `NetResult<T>` memiliki `ok`, `value`, dan `error`. Gunakan value hanya jika berhasil. `NetError` berisi klasifikasi kesalahan yang dinormalisasi dan native_code, dan keberhasilan dapat diperiksa dengan `error.kind == NET_ERROR_NONE`. Konstanta ini diambil dari `std::net::error`.

Aliran yang terhubung dan aliran yang diterima dengan accept masing-masing harus ditutup. Menutup pendengar tidak berarti menutup semua aliran yang telah diterimanya. Jangan tutup setiap salinan struktur.

## TCP tidak mempertahankan batasan pesan

Jangan berasumsi bahwa semua yang pernah Anda tulis akan kembali lagi kepada Anda setelah Anda membacanya. Protokol harus dibatasi oleh awalan panjang, pembatas, atau panjang tetap. Panjang tetap ditangani oleh `tcp_read_exact`, dan aliran dengan panjang yang tidak diketahui ditangani oleh iterasi baca dan pemrosesan EOF.

Angka negatif adalah kesalahan, dan dalam pembacaan panjang positif, 0 adalah ujung dari ujung lainnya. write_all Beberapa data mungkin telah dikirim sebelum kegagalan. Jika Anda mengirimkan ulang konten yang sama dari awal, konten tersebut mungkin terduplikasi.

## Batas waktu dan waktu tunggu

Fungsi default blocking dapat menunggu lama di sisi lain. Untuk program yang memerlukan batasan waktu, pilih seri `tcp_connect_addr_timeout`, `tcp_read_timeout`, `tcp_write_timeout`. Periksa faktor milidetik dan konsekuensi kesalahan, dan rancang juga rute di mana pihak lain tidak merespons. Penggunaan asinkron [task](/docs/id/stdlib/task) bersama dengan jaringan asinkron terkait API.

[pencarian alamat](/docs/id/stdlib/resolution) · [Klien lokal yang sudah selesai](/docs/id/practice/tcp-client)

## Kriteria pemilihan fungsi membaca

|diperlukan tindakan|fungsi|nilai untuk diperiksa|
| --- | --- | --- |
|Baca beberapa data yang tiba| `tcp_read` |Kesalahan negatif, penghentian nol, jumlah byte positif|
|Bacalah dengan panjang tertentu| `tcp_read_exact` |Apakah panjang permintaan sudah dibaca secara lengkap?|
|Kirim konten yang diberikan sampai akhir| `tcp_write_all` |Panjang Permintaan dan Nilai Pengembalian|
|Batas waktu tunggu|timeout Seri|Hasil pengembalian dan kesalahan timeout|

Untuk protokol dengan bidang panjang empat byte yang diikuti oleh isi, pertama-tama baca bidang panjang lengkap, periksa apakah panjangnya tidak melebihi maksimum yang diizinkan, lalu alokasikan ruang untuk isi. Jangan mengalokasikan memori dalam jumlah besar dari panjang yang tidak divalidasi yang disediakan oleh rekan.

## Kapan harus menutup koneksi

Meskipun fungsi baca mengembalikan 0, sumber daya aliran lokal tetap ada. Setelah selesai membaca, hubungi tcp_close. Jika Anda menggunakan jalur pembersihan yang sama bahkan setelah terjadi kesalahan transmisi, Anda tidak akan lupa untuk menutup jalur normal dan gagal.

Listener adalah sumber daya yang menerima koneksi baru, dan stream adalah sumber daya komunikasi yang sudah terhubung. Saat membuat server, Anda mengelola kedua jenis tersebut. Setiap kali accept berhasil, aliran baru dibuat, jadi kami menutupnya setelah memprosesnya, dan menutup pendengar ketika iterasi penerimaan server selesai.

## Cobalah sendiri

[Praktek Klien TCP Lokal](/docs/id/practice/tcp-client) berisi kode lengkap untuk server dan klien. Anda dapat membedakan antara pencarian alamat dan kegagalan koneksi dengan membandingkan kasus ketika server dijalankan terlebih dahulu, ketika tidak ada server, dan ketika nomor port berbeda.
