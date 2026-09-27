---
translation_set_id: vex-package-manager
path: whale/vex-package-manager
locale: vi
group: whale
group_order: 1
order: 4
title: Trình quản lý gói Vex
summary: Mô tả các dự án dựa trên manifest, các dự án Wave, Git·các phần phụ thuộc của đường dẫn, lockfile, các bản dựng ngoại tuyến và ranh giới wavec.
---

## vai trò

Vex là công cụ xây dựng và quản lý gói dành cho Wave. Vex hoạt động dựa trên `wavec`. Vex chịu trách nhiệm về cấu trúc dự án và phân tích phụ thuộc, còn `wavec` chịu trách nhiệm về cờ trình biên dịch và quy trình biên dịch.

Lệnh Vex dựa trên manifest. `vex build`, `vex check` và `vex run` cố tình không nhận cờ raw `wavec`.

## Tạo một gói

```shell
vex init
vex init --lib
```

Ứng dụng sử dụng `src/main.wave` và thư viện sử dụng `src/lib.wave`. Cấu trúc gốc của gói như sau:

```text
my_project/
├── src/
│   └── main.wave
├── vex.ws
├── vex.lock
└── .vex/
    └── deps/
```

`vex.ws` trở thành manifest. Vex không sử dụng manifest trong phần mở rộng `.wson`.

```wson
{
    name = "my_project",
    version = 0.1.0,
    lib = false,
    description = "my_project Project",
    author = "unknown",
    license = "Unknown",
    dependencies = []
}
```

## lệnh xây dựng

```shell
vex build [--target <triple>] [--release] [--dry-run] [--locked] [--offline]
vex check [--target <triple>] [--release] [--dry-run] [--locked] [--offline]
vex run   [--target <triple>] [--release] [--dry-run] [--locked] [--offline] [-- <args...>]
```

Các tùy chọn cho Vex được giữ ở mức nhỏ. Nếu bạn cần điều khiển phụ thuộc vào trình biên dịch, chẳng hạn như emit, linker, CPU, ABI hoặc debug, hãy sử dụng trực tiếp `wavec`. Khi bạn cần sử dụng một trình biên dịch cụ thể, hãy đặt `VEX_WAVEC=/path/to/wavec`.

Các bước tiến trình như `Resolving`, `Fetching`, `Compiling`, `Checking`, `Running`, `Finished` được xuất ra ở stderr và đầu ra chương trình được duy trì ở stdout.

## Git Phụ thuộc trung tâm

Các phần phụ thuộc Vex được chỉ định là `path` hoặc Git URL cục bộ. Một phần phụ thuộc chỉ có thể sử dụng một trong hai phương pháp.

```wson
{
    name = "app",
    version = 0.1.0,
    dependencies = [
        { name = "local_math", path = "../local_math" },
        { name = "remote_math", git = "https://github.com/example/math.git", tag = "v0.1.0" }
    ]
}
```

Phần phụ thuộc Git chỉ có thể chỉ định tối đa một trong số `branch`, `tag` hoặc `rev`. Mỗi gốc phụ thuộc phải có `vex.ws` riêng. Vex giải quyết đệ quy phần phụ thuộc manifest, từ chối nhận dạng gói xung đột và lưu trữ Git checkout được quản lý trong `.vex/deps/<name>`.

## lockfile Hợp đồng

Lược đồ v2 `vex.lock` ghi lại toàn bộ biểu đồ phụ thuộc bắc cầu và chính xác Git commit. Cam kết với manifest. Việc sử dụng cùng một manifest và một lockfile hợp lệ sẽ chọn cùng một biểu đồ phụ thuộc mà không cần tuân theo branch hoặc tag lần nữa.

Các lệnh yêu cầu phụ thuộc sẽ được diễn giải tự động và có thể được chuẩn bị trước bằng các lệnh sau.

```shell
vex fetch
vex update
vex update math shared_core
```

`vex update` cập nhật tất cả các gói Git hoặc chỉ các gói được chỉ định và biểu đồ chuyển tiếp bị ảnh hưởng. Các gói bị khóa không liên quan sẽ tiếp tục chọn commit.

## quy trình làm việc locked và offline

`--locked` nghiêm cấm việc tạo và sửa đổi `vex.lock`. Nó sẽ không thành công nếu tệp không tồn tại, có lược đồ không được hỗ trợ hoặc không khớp với biểu đồ manifest. Đã ghim vào lockfile, commit có thể nhập khi cần.

`--offline` cấm tất cả các hoạt động mạng Git. checkout và commit bắt buộc phải tồn tại cục bộ.

```shell
vex fetch --locked
vex build --locked --offline
```

Hai lệnh này là quy trình làm việc nghiêm ngặt CI. Chuẩn bị khóa commit chính xác khi mạng khả dụng, sau đó biên dịch nó mà không thay đổi mạng hoặc lockfile. dry-run không nhập phần phụ thuộc hoặc viết lại lockfile.

## Thông tin và cài đặt trình biên dịch

```shell
vex info
vex setup wavec
vex setup wavec --version <version>
vex --version
```

Vex xác thực lược đồ `wavec` dry-run JSON trước khi xây dựng thực tế. Các trình biên dịch không triển khai lược đồ được yêu cầu sẽ không thực thi với một kế hoạch không xác định và sẽ từ chối nó với lỗi tương thích.
