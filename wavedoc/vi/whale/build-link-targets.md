---
translation_set_id: build-link-targets
path: whale/build-link-targets
locale: vi
group: whale
group_order: 1
order: 5
title: Tùy chọn biên dịch, liên kết và nền tảng đích
summary: Mô tả emit sản phẩm bàn giao, loại đầu vào, liên kết, target/CPU/ABI và kế hoạch xây dựng độc lập.
---

## emit Đầu ra

```shell
wavec build main.wave --emit=check
wavec build main.wave --emit=ast
wavec build main.wave --emit=ir
wavec build main.wave --emit=bc
wavec build main.wave --emit=asm
wavec build main.wave --emit=obj -o main.o
wavec build main.wave --emit=bin -o app
```

artifact emit Các loại là `ast`, `ir`, `bc`, `asm`, `obj`, `bin`. `check` là chế độ kiểm soát kiểm tra, không phải là loại có thể phân phối và không được sử dụng cùng với chế độ khác artifact emit.

```shell
wavec print supported-emit-kinds
```

## Loại đầu vào và link-only

Ngoài nguồn Wave, trình biên dịch còn phân biệt giữa các đầu vào IR, bitcode, assembly, object và archive. Danh sách hỗ trợ được truy vấn bằng lệnh sau:

```shell
wavec print supported-input-types
```

Để chỉ liên kết object hoặc archive đã được tạo, bạn có thể sử dụng `--input-type` và `--link-only`.

```shell
wavec build module.o --input-type=obj --link-only --emit=bin -o app
```

## liên kết gốc

```shell
wavec --link=m -L ./lib build main.wave
```

`--link` thêm thư viện và `-L` thêm đường dẫn tìm kiếm. Ngay cả khi ký hiệu được khai báo trong FFI, thư viện cung cấp ký hiệu đó sẽ không được liên kết tự động.

## Chọn mục tiêu

Tùy chọn để chọn OS và CPU để chạy.

- `--target <triple>`
- `--cpu <name>`
- `--features <csv>`
- `--abi <name>`
- `--sysroot <path>`

Kiểm tra mặc định của máy chủ và các mục tiêu được hỗ trợ bằng lệnh sau:

```shell
wavec print host-target
wavec print supported-targets
wavec print target-spec --target <triple>
wavec print cpu-list --target <triple>
wavec print target-features --target <triple>
```

## Được hỗ trợ bởi

Chương trình cho máy tính hiện tại của bạn sẽ được xây dựng mà không chỉ định target. Để chọn một môi trường khác, hãy chuyển tên mục tiêu bên dưới cho `--target`.

|OS·Môi trường|kiến trúc|tên mục tiêu|
| --- | --- | --- |
| Linux | amd64 | `x86_64-unknown-linux-gnu` |
| Linux | ARM64 | `aarch64-unknown-linux-gnu` |
| Linux | RISC-V64 | `riscv64-unknown-linux-gnu` |
| Linux | LoongArch64 | `loongarch64-unknown-linux-gnu` |
| macOS | amd64 | `x86_64-apple-darwin` |
| macOS | ARM64 | `aarch64-apple-darwin` |
| Windows | amd64 | `x86_64-pc-windows-msvc` |
| Windows | ARM64 | `aarch64-pc-windows-msvc` |
| FreeBSD | amd64 | `x86_64-unknown-freebsd` |
| WebAssembly |64 bit| `wasm64-unknown-unknown` |

Đối với các chương trình chạy không có OS, hãy sử dụng `x86_64-unknown-none-elf`, `aarch64-unknown-none-elf` và `riscv64-unknown-none-elf`. Bạn có thể tìm thấy danh sách đầy đủ các mục tiêu cho các phiên bản đã cài đặt tại `wavec print supported-targets`.

## RISC-V Hợp đồng 64

Các giá trị mặc định cho mục tiêu Hosted RISC-V là `generic-rv64`, RV64GC, `lp64d` ABI. Các giá trị mặc định cho mục tiêu Freestanding là `generic-rv64`, RV64IMAC và `lp64`.

```shell
wavec print target-spec --target riscv64-unknown-linux-gnu --format=json
wavec print target-spec --target riscv64-unknown-none-elf --format=json
```

Hỗ trợ RISC-V CPU cho `generic`, `generic-rv64`, `rocket-rv64`, `sifive-u74`. Feature override là `m`, `a`, `f`, `d`, `c`, `zicsr`, `zifencei` Đặt ký hiệu trước tên và ngăn cách bằng dấu phẩy.

```shell
wavec build main.wave \
  --target riscv64-unknown-linux-gnu \
  --features=+m,+a,+f,-d,+c,+zicsr \
  --abi=lp64f
```

RISC-V Xác minh từ chối các kết hợp không nhất quán. `d` yêu cầu `f` và `f` yêu cầu `zicsr`. `lp64`, `lp64f`, `lp64d` phải khớp với dấu phẩy động feature bạn đã bật. Nếu bạn không chỉ định trực tiếp ABI, trình biên dịch sẽ lấy ABI từ feature.

## Liên kết độc lập

```shell
wavec build kernel.wave \
  --target riscv64-unknown-none-elf \
  --freestanding \
  --entry=_start \
  --linker-script=linker.ld \
  --no-start-files \
  -o kernel.elf
```

`--freestanding` điều chỉnh cài đặt bản dựng để tránh sử dụng thư viện mặc định. `--entry` chỉ định các mục liên kết, `--linker-script` chỉ định tập lệnh và `--no-start-files` chỉ định loại trừ tệp khởi động máy chủ.

Bạn có thể sử dụng `--dry-run` để kiểm tra sơ đồ liên kết trước khi thực hiện thực tế.

## Hosted Liên kết chéo

Khi tạo một chương trình để chạy trên OS·CPU khác, hãy chỉ định đường dẫn thư viện của môi trường đích là sysroot. Nếu bạn sử dụng một trình liên kết riêng, hãy chỉ định đường dẫn là `-C linker`.

```shell
wavec build main.wave \
  --target riscv64-unknown-linux-gnu \
  --sysroot /path/to/riscv64-sysroot \
  -C linker=/path/to/target-linker \
  -o app-riscv64
```

Các thư viện được liên kết với sysroot được căn chỉnh với OS·CPU·ABI đã chọn.

## Những gì cần kiểm tra trong bản dựng chéo

- target triple nằm trong danh sách hỗ trợ trình biên dịch của bạn
- sysroot và trình liên kết khớp với mục tiêu ABI
- Thư viện liên kết có dành cho kiến trúc đích không?
- CPU feature hợp lệ cho mục tiêu CPU
- Nếu đứng độc lập, hãy kiểm tra xem ký hiệu mục nhập và vị trí bộ nhớ có khớp với tập lệnh trình liên kết hay không.

## Chạy WebAssembly

Kết quả wasm64 được sử dụng trong môi trường thực thi hỗ trợ memory64. Các module sử dụng chức năng bên ngoài như file, thời gian, đầu vào phải kết nối host import tương ứng với chức năng đó.

Quá trình tạo và kết nối mã mục tiêu có thể được tìm thấy trong `--dry-run`. Quá trình thực thi thực tế diễn ra trong môi trường thực thi OS hoặc WebAssembly đã chọn.
