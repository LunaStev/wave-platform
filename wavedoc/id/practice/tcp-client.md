---
translation_set_id: practice-tcp-client
path: practice/tcp-client
locale: id
group: practice
group_order: 4
order: 4
title: Proyek: Klien TCP lokal
summary: Menghubungkan pencarian alamat numerik, kegagalan koneksi, transfer dan penutupan.
---

## Persiapan dan Ruang Lingkup

Program ini terhubung ke `127.0.0.1:8080` pada OS asli, mengirim `ping` dan LF, dan keluar. Anda memerlukan server pengujian di terminal terpisah. Jika Python 3 ada, server berikutnya akan menerima dan menampilkan hingga 5 byte dari satu koneksi lokal.

```python
import socket
with socket.socket() as server:
    server.bind(("127.0.0.1", 8080))
    server.listen(1)
    connection, address = server.accept()
    with connection:
        message = b""
        while len(message) < 5:
            chunk = connection.recv(5 - len(message))
            if not chunk:
                break
            message += chunk
        print(repr(message))
```

Jalankan server terlebih dahulu, lalu simpan klien sebagai `main.wave` dan jalankan sebagai `wavec run main.wave`. Server mengeluarkan `b'ping\n'`. Jika port sudah digunakan, ubah port di server dan klien secara bersamaan.

<!-- wave-example: tcp-client -->
```wave
import("std::net::address")::{
    SocketAddr
};
import("std::net::resolve")::{
    NetResolveResult, net_resolve_tcp
};
import("std::net::error")::{
    NetResult, NetError, NET_ERROR_NONE
};
import("std::net::tcp")::{
    TcpStream, tcp_connect_addr_timeout, tcp_write_all, tcp_close
};

fun main() -> i32 {
    var addresses: array<SocketAddr, 4>;
    var resolved: NetResolveResult = net_resolve_tcp("127.0.0.1", "8080", &addresses[0], 4);
    if (!resolved.ok || resolved.written == 0) {
        println("resolve failed");
        return 1;
    }

    var connected: NetResult<TcpStream> = tcp_connect_addr_timeout(addresses[0], 1000);
    if (!connected.ok) {
        println("connect failed");
        return 2;
    }

    var sent: i64 = tcp_write_all(connected.value, "ping\n" as ptr<u8>, 5);
    var closed: NetError = tcp_close(connected.value);
    if (sent != 5 || closed.kind != NET_ERROR_NONE) {
        return 3;
    }

    println("sent=5");
    return 0;
}
```

Output keberhasilan dari klien adalah `sent=5`. Jika server tidak berjalan, penanganan kegagalan yang normal adalah dengan mencetak `connect failed` dan keluar dengan 2.

## Memahami Kode

Rentang array pencarian yang valid hingga written. Ini adalah contoh kecil yang hanya menggunakan alamat pertama. Dalam layanan dengan banyak kandidat, Anda perlu menentukan kebijakan koneksi untuk setiap kandidat. Aliran yang terhubung ditutup terlepas dari keberhasilan atau kegagalan transmisi. Batas waktu koneksi adalah 1000 ms, dan ini bukan contoh yang menjamin batas waktu untuk seluruh transmisi.

TCP tidak mempertahankan batas transmisi. Inilah sebabnya mengapa ini ditulis agar server dapat dijalankan recv beberapa kali. Cara menemukan alamat berdasarkan nama tercakup dalam [pencarian alamat](/docs/id/stdlib/resolution).

## latihan yang diperpanjang

Minta server mengirimkan respons dan minta klien membacanya. Setelah menentukan durasi respons, Anda perlu menangani pembacaan sebagian, EOF, batas waktu, dan penutupan secara bersamaan.

## Peran kedua terminal

Terminal server menunggu koneksi dalam status listen. Saat klien terhubung, accept mengembalikan soket yang terhubung dan loop recv mengumpulkan data. Terminal klien melakukan pencarian alamat, koneksi, transfer, dan penutupan dalam urutan itu.

Kode Python adalah server mitra lab. Wave Digunakan untuk dengan mudah memeriksa byte yang diteruskan oleh program, dan tidak ditempatkan di file yang sama dengan kode klien. Server berakhir setelah menangani satu koneksi, sehingga server juga restart sebelum menjalankan klien lagi.

## Mengapa mengirim 5 byte

Karena `ping` berukuran 4 byte dan LF berukuran 1 byte, panjang transmisinya adalah 5. NUL di akhir string tidak disertakan dalam pesan. Server mengumpulkan 5 byte dan kemudian mengeluarkannya, sehingga meskipun paket tiba dalam beberapa bagian, hasil yang sama akan dihasilkan.

Jika tcp_write_all mengembalikan 5, fungsi transfer lokal telah memproses byte yang diminta. Ini tidak mengkonfirmasi bahwa program lain telah menafsirkan dan menyimpan pesan tersebut. Jika protokol memerlukan konfirmasi penyelesaian pemrosesan, server mengirimkan respons dan klien membaca respons tersebut.

## Praktek Jalan Menuju Kegagalan

|perubahan|hasil|Apa yang harus dipelajari|
| --- | --- | --- |
|Jalankan tanpa server|Koneksi gagal, kode keluar 2|Meskipun alamatnya valid, servernya mungkin tidak ada|
|Atur port server dan klien secara berbeda|Koneksi gagal|Baik IP maupun port di alamat harus cocok.|
|Ubah ukuran server recv menjadi 1|Sama 5 byte bersama-sama|TCP Unit baca berbeda dengan unit pesan|
|Membagi transfer klien menjadi beberapa kali|diterima dalam urutan yang sama|Batasan pesan ditentukan oleh protokol.|

Batas waktu koneksi 1000ms merupakan pengaturan pada fase koneksi. Untuk membatasi membaca dan menulis, gunakan fungsi timeout pada langkah yang relevan, dan jika ada batas waktu yang berlaku untuk keseluruhan program, hitung sisa waktu dan teruskan.
