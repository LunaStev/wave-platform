---
translation_set_id: compiler
path: getting-started/compiler
locale: vi
group: getting-started
group_order: 1
order: 3
title: Tham chiếu lệnh biên dịch
summary: wavec Mô tả các lệnh, xây dựng quy trình, đầu ra, mục tiêu, chẩn đoán, liên kết phụ thuộc và truy vấn công cụ.
---

## mô hình lệnh

`wavec` là trình biên dịch CLI. Nó biên dịch trực tiếp các đầu vào riêng lẻ, cung cấp thông tin hỗ trợ trình biên dịch cho các công cụ và quản lý các nguồn thư viện tiêu chuẩn đã cài đặt.

```text
wavec [global-options] <command> [command-options]
```

|lệnh|sử dụng|
| --- | --- |
| `wavec build <input...>` |Tùy thuộc vào các cờ, nó thực hiện các đường dẫn kiểm tra, tạo mã, liên kết hoặc thực thi.|
| `wavec check <file>` |Biệt danh của `build <file> --emit=check`.|
| `wavec run <file> [-- <args...>]` |Đó là bí danh cho `build <file> --run` và đối số sau `--` được truyền vào chương trình.|
| `wavec print <item>` |Truy vấn mục tiêu và thông tin hỗ trợ chuỗi công cụ.|
| `wavec install std` |Cài đặt thư viện chuẩn.|
| `wavec update std` |Cập nhật các thư viện chuẩn đã cài đặt.|
| `wavec --version` |In thông tin phiên bản đã cài đặt.|

Bạn có thể tìm thấy danh sách đầy đủ các lệnh và tùy chọn tại `wavec --help`.

## Xây dựng, thử nghiệm và chạy

```shell
wavec build main.wave
wavec check main.wave
wavec run main.wave -- first-argument second-argument
```

`build` tạo một tệp thực thi theo mặc định. `check` dừng sau khi hoàn tất kiểm tra giao diện người dùng. `run` yêu cầu đầu ra nhị phân và không thể sử dụng với các bản dựng thư viện dùng chung.

Sử dụng `--dry-run` để xác minh yêu cầu và xác định các bước cần thực hiện mà không cần biên dịch, liên kết hoặc thực thi.

```shell
wavec build main.wave --target riscv64-unknown-linux-gnu --dry-run
wavec build main.wave --dry-run --error-format=json
```

Định dạng JSON là giao diện thống nhất, ổn định được sử dụng bởi các công cụ xây dựng như Vex.

## emit và loại đầu vào

```shell
wavec build main.wave --emit=ast
wavec build main.wave --emit=ir
wavec build main.wave --emit=bc
wavec build main.wave --emit=asm
wavec build main.wave --emit=obj -o main.o
wavec build main.wave --emit=bin -o app
```

Các loại đầu ra emit là `ast`, `ir`, `bc`, `asm`, `obj`, `bin`. `check` là chế độ điều khiển và phải được sử dụng một mình. Bạn có thể chỉ định nhiều loại đầu ra mà một đường dẫn chấp nhận, được phân tách bằng dấu phẩy.

Các loại đầu vào là `wave`, `ir`, `bc`, `asm`, `obj`, `archive`. `--input-type=<kind>` buộc loại tất cả đầu vào phải được chỉ định. Khi chỉ liên kết các đầu vào object hoặc archive, hãy sử dụng các tệp nhị phân emit và `--link-only`.

```shell
wavec build module.o --input-type=obj --link-only --emit=bin -o app
```

## vị trí đầu ra

|tùy chọn|hiệu ứng|
| --- | --- |
| `-o <file>` |Chỉ định đường dẫn đầu ra chính.|
| `--out-dir <dir>` |emit Đặt đầu ra vào thư mục đã chỉ định.|
| `--target-dir <dir>` |Chỉ định các tuyến đường phân phối trung gian và chính.|

## Đầu ra tối ưu hóa và chẩn đoán

```shell
wavec -O2 build main.wave
wavec --debug-wave=tokens,ast build main.wave
```

Các bước tối ưu hóa là `-O0`, `-O1`, `-O2`, `-O3`, `-Os`, `-Oz`, `-Ofast`. `--debug-wave` có thể được theo sau bởi `tokens`, `ast`, `ir`, `mc`, `hex`, `all` và có thể kết hợp nhiều bước bằng dấu phẩy.

## liên kết gốc

```shell
wavec --link=m -L ./lib build main.wave
wavec build main.wave --shared -o libexample.so
wavec build main.wave --static -o app
wavec build main.wave --pie -o app
```

`--link=<lib>` thêm thư viện gốc và `-L <path>` thêm đường dẫn tìm kiếm. Chế độ liên kết sử dụng `--shared`, `--static`, `--pie`, `--no-pie` theo quy tắc tương thích.

Các tùy chọn kiểm soát phụ trợ và trình liên kết bao gồm:

- `--target`, `--cpu`, `--features`, `--abi`, `--sysroot`
- `-C linker=<path>` và `-C link-arg=<arg>`
- `-C link-sysroot=<path>` và `-C relocation-model=<model>`
- `-C no-default-libs`

Các đầu ra độc lập như hạt nhân sử dụng `--freestanding` cùng với các cài đặt `--entry`, `--linker-script` và `--no-start-files` phù hợp với môi trường.

## Giải thích gói bên ngoài

```shell
wavec --dep-root .vex/deps build main.wave
wavec --dep math=/opt/wave-deps/math build main.wave
```

`--dep-root <dir>` thêm root để tìm `package::module` import bên ngoài. `--dep <name>=<path>` sửa tên gói vào một thư mục. Đây là điểm tích hợp trình biên dịch và các dự án manifest, các lượt tải xuống phụ thuộc và lockfile được xử lý bởi Vex.

## Hỗ trợ truy vấn chức năng

Các công cụ sử dụng loại mục tiêu hoặc đầu ra có thể truy vấn thông tin hỗ trợ bằng `wavec print`.

```shell
wavec print host-target
wavec print target-spec --format=json
wavec print supported-targets
wavec print supported-input-types
wavec print supported-emit-kinds
wavec print supported-print-items
wavec print cpu-list --target riscv64-unknown-linux-gnu
wavec print target-features --target riscv64-unknown-linux-gnu
wavec print default-linker
wavec print sysroot
wavec print std-path
wavec print dep-search-paths
```

Bạn cũng có thể truy vấn các mục như `host`, `default-target` và `target-list`. Các mục hỗ trợ đầu ra có cấu trúc nhận được `--format=json`.

## Ranh giới giữa trình biên dịch và chuỗi công cụ

`wavec` chịu trách nhiệm kiểm tra nguồn, tạo mã và liên kết. Vex chịu trách nhiệm xây dựng các gói có thể tái tạo với gói manifest, biểu đồ phụ thuộc và lockfile. Whale là chuỗi công cụ cấp thấp chạy độc lập.

## std Chỉ định đường dẫn

```shell
wavec --std-root /absolute/path/to/std check main.wave
wavec --std-root /absolute/path/to/std run main.wave
```

Đường dẫn std được chỉ định sẽ được ưu tiên hơn đường dẫn cài đặt và sẽ thất bại nếu đường dẫn không chính xác hoặc không tương thích với std. Nó không tự động thay thế std từ các cài đặt khác. Chọn std tương ứng với trình biên dịch của bạn.

Đối với đầu ra `-o`, một đường dẫn khác được sử dụng với tệp nguồn/đầu vào. Vì `check` không kiểm tra hoạt động trong thời gian chạy nên [luyện tập](/docs/vi/practice/input-calculator) cũng kiểm tra kết quả thực thi.
