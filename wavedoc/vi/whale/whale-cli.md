---
translation_set_id: whale-cli
path: whale/whale-cli
locale: vi
group: whale
group_order: 1
order: 3
title: Tài liệu tham khảo lệnh Whale
summary: Mô tả Whale assembler, object wrapper, đầu ra chẩn đoán và các lệnh IR tùy chọn.
---

## Whale Xây dựng

Trong kho lưu trữ Whale, hãy chạy:

```shell
cargo build --release
```

Một tệp thực thi cấp cao nhất có bốn nhóm lệnh:

```text
whale asm [--amd64 | --aarch64] <input> -o <output>
whale object <input> -o <output>
whale ir <subcommand> [options]
```

## AMD64 assembler

```shell
whale asm --amd64 input.asm -o output.o
```

AMD64 assembler nhận đường dẫn `.o` làm đầu ra và ELF64 relocatable chứa section, symbol và relocation Tạo object.

Bật đầu ra chẩn đoán chi tiết bằng `--debug-whale`.

```shell
whale asm --amd64 input.asm -o output.o \
  --debug-whale --token --ast --bytes --dump-hex --stats
```

Cờ chẩn đoán bao gồm `--token`, `--ast`, `--bytes`, `--dump-hex`, `--dump-bin`, `--dump-json` và `--stats`. `--trace` in tiến trình xử lý.

## Object wrapper

```shell
whale object input.bin -o output.o
```

Lệnh `object` đặt các byte thô vào phần ELF64 `.text` và thêm biểu tượng `start` toàn cầu ở độ lệch 0. Lệnh này bao bọc mã máy thô trong tệp đối tượng ELF.

## Tùy chọn IR socket

Lệnh `ir` chỉ được đưa vào khi Whale được xây dựng bằng tính năng `socket-cli`.

```shell
cargo run -p whale --features socket-cli -- ir lower program.json
cargo run -p whale --features socket-cli -- ir lower program.json -o program.wir
```

`ir lower` đọc JSON của Whale socket schema, chuyển đổi nó thành Whale IR và xác minh mô-đun. Văn bản IR được xuất ra đường dẫn stdout hoặc `-o`. `--target <triple>` thay thế chuỗi đích và `--no-verify` bỏ qua xác thực.

Để sử dụng lệnh IR, Whale phải được xây dựng bằng tính năng `socket-cli`. Socket JSON nhà sản xuất và Whale phải sử dụng cùng socket schema version.
