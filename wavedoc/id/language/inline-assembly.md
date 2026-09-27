---
translation_set_id: assembly
path: language/inline-assembly
locale: id
group: language
group_order: 2
order: 19
title: Perakitan sebaris
summary: Menjelaskan kontrak string perintah blok asm, operan in/out, dan clobber.
---

## asm blok

`asm` adalah sintaks tingkat rendah untuk memasukkan instruksi arsitektur target secara langsung.

```wave
fun read_value() -> i64 {
    var result: i64 = 0;
    asm {
        "mov rax, 123"
        out("rax") result
    }

    return result;
}
```

String literal dalam satu blok diteruskan sebagai daftar instruksi perakitan.

## masukan dan keluaran

```wave
var result: i64 = 0;
asm {
    "mov rax, rdi"
    in("rdi") 123
    out("rax") result
}
```

- `in("reg") expression` menyambungkan nilai Wave ke operan input.
- `out("reg") target` menulis nilai keluaran ke target Wave yang dapat ditetapkan.
- Nama register dapat ditulis sebagai string atau pengidentifikasi.

Operan masukan dapat mencakup variabel, literal integer/string, `&identifier`, `deref identifier`, dan angka negatif.

## clobber

Jika sebuah blok mengubah status register atau memori selain output eksplisit, maka blok tersebut dicatat di `clobber(...)`.

```wave
asm {
    "nop"
    clobber("rax", "rcx", "memory")
}
```

## Periksa saat menggunakan

- Sintaks instruksi harus sesuai dengan arsitektur target dan kontrak perakitan inline LLVM.
- Jangan sembarangan memusnahkan register yang harus disimpan sesuai dengan konvensi pemanggilan.
- Untuk blok yang membaca atau menulis memori, deklarasikan clobber, termasuk `memory`.
- Jika memungkinkan, pisahkan asm khusus arsitektur di belakang fungsi kecil.

Perilaku dan portabilitas perakitan inline tidak dijamin hanya berdasarkan jenis bahasa saja.

## Rentang Pembelajaran dan Contoh

[Berlatihlah dengan program penuh](/docs/id/getting-started/overview) · [Perpustakaan standar](/docs/id/stdlib)
