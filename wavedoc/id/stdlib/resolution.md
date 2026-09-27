---
translation_set_id: stdlib-resolution
path: stdlib/resolution
locale: id
group: stdlib
group_order: 1
order: 8
title: net.resolve: Resolusi alamat
summary: Menjelaskan pemotongan alamat numerik, daftar nama milik penelepon, dan hasil.
---

## Pilih metode penyelidikan

Default untuk Linux, resolver, menggunakan alamat numerik IPv4/IPv6 dan port numerik. Nama host atau nama layanan tidak diteruskan secara implisit ke sistem DNS. Misalnya, `127.0.0.1` dan `8080` adalah input numerik, dan `example.com` dan `http` adalah nama.

Untuk menggunakan daftar nama yang diberikan secara eksplisit, gunakan `std::net::resolve_table`. Anda dapat terhubung langsung ke alamat numerik atau mencari dengan mendaftarkan alamat di daftar nama.

## Pencarian dasar

```text
std::net::resolve
net_resolve_tcp(host: str, service: str, output: ptr<SocketAddr>, capacity: i64) -> NetResolveResult
net_resolve_udp(host: str, service: str, output: ptr<SocketAddr>, capacity: i64) -> NetResolveResult
net_resolve_host(host: str, service: str, output: ptr<SocketAddr>, capacity: i64) -> NetResolveResult
```

Array keluaran disediakan oleh pemanggil. Satuan capacity bukanlah byte, melainkan `SocketAddr` jumlah elemen. Anda hanya dapat mencari hitungan dengan capacity=0 dan null output. Penyimpanan pencarian internal dibersihkan sebelum dikembalikan dan tidak mengambil alih kepemilikan larik keluaran.

|bidang hasil|artinya|
| --- | --- |
| `ok` |Apakah kueri berhasil atau tidak|
| `count` |Jumlah total hasil alamat yang didukung|
| `written` |Jumlah elemen yang sebenarnya ditulis ke larik keluaran|
| `truncated` |Apakah kapasitas keluarannya lebih kecil dari hasil totalnya?|
| `error.kind` |Klasifikasi yang dinormalisasi seperti INVALID, NOT_FOUND, UNSUPPORTED|
| `error.native_code` |Kode asli untuk menentukan penyebabnya|

Bahkan jika berhasil, written mungkin 0. Centang `ok && written > 0` sebelum menggunakan output[0].

## Daftar nama yang Anda berikan sendiri

```text
net_resolve_from_table(host: str, service: str, protocol: i32,
    entries: ptr<NetResolveEntry>, entry_count: i64,
    output: ptr<SocketAddr>, capacity: i64) -> NetResolveResult
```

`NetResolveEntry` di `std::net::resolve_table` memiliki host, service, protocol, address. Nama host ASCII tidak peka huruf besar-kecil dan nama layanan dibandingkan dengan tepat. TCP/UDP/Untuk semua kueri protokol, konstanta modul digunakan. Item dan string dimiliki oleh pemanggil dan tidak disimpan setelah panggilan.

Kecocokan disalin dalam urutan masukan dan duplikat dipertahankan. Ruang keluaran dan ruang penyimpanan masukan tidak boleh tumpang tindih. Jika namanya tidak ada dalam daftar, itu adalah NOT_FOUND dan tidak akan mencoba lagi dengan DNS eksternal.

[TCP Cara menggunakan](/docs/id/stdlib/tcp) · [TCP Praktek Klien](/docs/id/practice/tcp-client)

## Cari alamat numerik

Siapkan array untuk menampung hasil kueri dan periksa apakah alamatnya dicatat. Contoh di bawah ini tidak memerlukan server karena pencarian alamat tidak membuat koneksi itu sendiri.

<!-- wave-example: book-resolve-numeric -->
```wave
import("std::net::address")::{
    SocketAddr
};
import("std::net::resolve")::{
    NetResolveResult,
    net_resolve_tcp
};

fun main() -> i32 {
    var addresses: array<SocketAddr, 4>;
    var result: NetResolveResult = net_resolve_tcp(
        "127.0.0.1",
        "8080",
        &addresses[0],
        4
    );

    if (!result.ok || result.written == 0) {
        println("address unavailable");
        return 1;
    }

    println("address ready");
    return 0;
}
```

Hasil eksekusi:

```text
address ready
```

Untuk benar-benar terhubung, teruskan addresses[0] ke rangkaian fungsi tcp_connect_addr. Bahkan jika Anda mendapatkan alamatnya, Anda dapat mengetahui pada tahap koneksi apakah server berjalan pada port tersebut.

## Cari berdasarkan daftar nama

Daftar alamat khusus berguna untuk menghubungkan nama berdasarkan file konfigurasi atau daftar layanan program. Contoh berikut memetakan nama api ke port lokal 8080.

<!-- wave-example: book-resolve-table -->
```wave
import("std::net::address")::{
    SocketAddr,
    socket_addr_from_v4,
    socket_addr_v4_loopback
};
import("std::net::resolve")::{
    NetResolveResult
};
import("std::net::resolve_table")::{
    NetResolveEntry,
    RESOLVE_TCP,
    net_resolve_from_table
};

fun main() -> i32 {
    var entries: array<NetResolveEntry, 1>;
    entries[0] = NetResolveEntry {
        host: "api",
        service: "http",
        protocol: RESOLVE_TCP,
        address: socket_addr_from_v4(socket_addr_v4_loopback(8080))
    };

    var addresses: array<SocketAddr, 2>;
    var result: NetResolveResult = net_resolve_from_table(
        "API",
        "http",
        RESOLVE_TCP,
        &entries[0],
        1,
        &addresses[0],
        2
    );

    if (!result.ok || result.written != 1) {
        return 1;
    }

    println("matched={}", result.written);
    return 0;
}
```

Hasil eksekusi:

```text
matched=1
```

API dan api cocok dalam perbandingan nama host. http dan HTTP tidak cocok dalam contoh ini karena berbeda dalam perbandingan nama layanan. Jika Anda mendaftarkan beberapa alamat dengan nama yang sama, hasilnya akan muncul sesuai urutan yang dimasukkan. Jika Anda hanya menerima sebagian dari array kecil, baca hanya hingga written, dan jika Anda memerlukan seluruh daftar, siapkan ruang untuk count dan cari lagi.
