---
translation_set_id: install
path: getting-started/install
locale: vi
group: getting-started
group_order: 1
order: 2
title: Cài đặt Wave
summary: Cài đặt Wave trên Linux, macOS hoặc Windows và chạy chương trình đầu tiên.
---

## Linux và macOS

Chạy lệnh sau trong cửa sổ dòng lệnh. Lệnh này cài đặt Wave cùng trình quản lý gói Vex.

```shell
curl -fsSL https://wave-lang.dev/install.sh | bash -s -- latest
```

Sau khi cài đặt xong, mở cửa sổ dòng lệnh mới và kiểm tra phiên bản.

```shell
wavec --version
```

## Windows

Chạy các lệnh sau trong PowerShell. Các lệnh này cài đặt Wave cùng trình quản lý gói Vex.

```powershell
irm https://wave-lang.dev/install.ps1 -OutFile install.ps1
powershell -ExecutionPolicy Bypass -File .\install.ps1 -Latest
```

Sau khi cài đặt xong, mở cửa sổ PowerShell mới và kiểm tra phiên bản.

```powershell
wavec --version
vex --version
```

## Chạy chương trình đầu tiên

Lưu đoạn mã sau vào tệp `main.wave`.

<!-- wave-example: install-stdlib -->
```wave
import("std::string::len")::{
    len
};

fun main() {
    println("Wave: {} bytes", len("Wave"));
}
```

Chạy chương trình từ thư mục chứa tệp vừa lưu.

```shell
wavec run main.wave
```

Kết quả:

```text
Wave: 4 bytes
```

Nếu xuất hiện thông báo không tìm thấy thư viện chuẩn, hãy cài đặt thư viện rồi chạy lại chương trình.

```shell
wavec install std
wavec run main.wave
```

[Tiếp theo: Chương trình đầu tiên](/docs/vi/language/program-structure) · [Khắc phục sự cố](/docs/vi/reference/diagnostics)
