---
translation_set_id: whale-assembler-linker
path: whale/assembler-linker
locale: vi
group: whale
group_order: 1
order: 11
title: Hợp dịch và liên kết tĩnh
summary: Mô tả mã hóa toán hạng, vị trí phần, liên kết ký hiệu và các điểm nhập thực thi.
---

## Hội đồng và đối tượng

Trình biên dịch Whale chuyển đổi lệnh AMD64 thành byte mã máy và thông tin di chuyển. Đối tượng ELF64 chứa thông tin này và các phần/ký hiệu. Việc lắp ráp được thực hiện bằng cách triển khai riêng của Whale và không yêu cầu trình biên dịch mã bên ngoài.

Các đối tượng có thể định vị lại có thể có các tham chiếu chưa biết địa chỉ cuối cùng. Quá trình giải quyết địa chỉ này là một liên kết. Việc lắp ráp thành công không nên suy ra rằng tất cả các ký hiệu bên ngoài đã được giải quyết hoặc một tệp thực thi đã được tạo.

## Lắp ráp chức năng

Lưu mã sau dưới dạng `answer.asm`.

```asm
section .text
global answer

answer:
    mov eax, 42
    ret
```

```sh
whale asm --amd64 answer.asm -o answer.o
```

Nội dung của `.text` là `b8 2a 00 00 00 c3`, và `mov eax, 42` theo sau là `ret`. Đối tượng ELF64 lộ ra `answer`. Đây là một hàm có thể gọi được mà không cần mã khởi động quy trình và không phải là một tệp thực thi. Các hướng dẫn trong ví dụ này có thể được xử lý bởi trình biên dịch mã hiện tại.

## Xây dựng các đối tượng với Rust API

Đây là ví dụ hoàn chỉnh về cách ghi các byte hàm tương tự vào thùng `object`. Chỉ định mục tiêu đầu ra và ký hiệu chung.

```rust
use object::{ObjectFile, ObjectSymbol, SectionKind, SymbolBinding, SymbolVisibility, Target};

fn main() {
    let mut object = ObjectFile::with_target(Target::X86_64WhaleLinux.object_target());
    let text = object.add_section(".text", SectionKind::Text, 16);
    object.sections[text].data = vec![0xb8, 0x2a, 0x00, 0x00, 0x00, 0xc3];
    object.symbols.push(ObjectSymbol {
        name: "answer".into(),
        section_index: Some(text),
        value: 0,
        size: 6,
        binding: SymbolBinding::Global,
        visibility: SymbolVisibility::Default,
    });
    let elf = object.write().unwrap();
    assert_eq!(&elf[..7], b"\x7fELF\x02\x01\x01");
    assert_eq!(u16::from_le_bytes([elf[18], elf[19]]), 62);
    std::fs::write("answer.o", elf).unwrap();
}
```

Assertion kiểm tra lớp ELF, thứ tự byte, mã định danh machine. `value: 0` là `.text` bên trong offset và `size: 6` là kích thước byte của ký hiệu. Tham chiếu phần hoặc phạm vi không hợp lệ là lỗi tuần tự hóa. Từ chối các đơn hàng machine hoặc byte khác mà không đánh dấu chúng là AMD64.

## Giải thích các ký hiệu của hai đối tượng

Hiện tại thùng `linker` cung cấp giải thích ký hiệu. Ví dụ thực thi sau đây xác định `helper` cục bộ trên hai đối tượng, sau đó kiểm tra xem có xảy ra lỗi hay không nếu cùng một tên được xuất bản hai lần.

```rust
use linker::core::symbol_table::{SymbolKey, SymbolTable};
use object::{
    ObjectFile, ObjectFormat, ObjectSymbol, SectionKind, SymbolBinding, SymbolVisibility,
};

fn input(binding: SymbolBinding) -> ObjectFile {
    let mut object = ObjectFile::new(ObjectFormat::ELF64);
    let text = object.add_section(".text", SectionKind::Text, 1);
    object.sections[text].data = vec![0xc3];
    object.symbols.push(ObjectSymbol {
        name: "helper".into(),
        section_index: Some(text),
        value: 0,
        size: 1,
        binding,
        visibility: SymbolVisibility::Default,
    });
    object
}

fn main() {
    let mut symbols = SymbolTable::new();
    symbols
        .resolve(&[input(SymbolBinding::Local), input(SymbolBinding::Local)])
        .unwrap();
    for object_index in 0..2 {
        let key = SymbolKey::Local {
            object_index,
            name: "helper".into(),
        };
        assert_eq!(symbols.symbols[&key].object_index, Some(object_index));
    }
    let error = symbols
        .resolve(&[input(SymbolBinding::Global), input(SymbolBinding::Global)])
        .unwrap_err();
    assert_eq!(error, "Duplicate global symbol: helper");
    println!("{error}");
}
```

Đầu ra:

```text
Duplicate global symbol: helper
```

Hai định nghĩa cục bộ có các khóa riêng biệt với `object_index` khác nhau. Hai định nghĩa toàn cầu xung đột. Ví dụ này diễn giải trực tiếp các ký hiệu của đối tượng bạn đã tạo trong bộ nhớ. `.o` Việc đọc tệp, áp dụng di chuyển hoặc xuất tệp thực thi không được thực hiện. Trong số các chính sách lựa chọn định nghĩa đầy đủ được mô tả bên dưới, các ưu tiên weak, v.v. vẫn chưa được triển khai.

## Ký tự và toán hạng bộ nhớ

Literal duy trì chiều rộng và biển báo của chúng cho đến khi kiểm tra phạm vi mã hóa của lệnh thực tế. Các giá trị nằm ngoài phạm vi sẽ dẫn đến lỗi thay vì bị cắt bớt một cách âm thầm.

Cần có ký hiệu kích thước rõ ràng khi không thể xác định được độ rộng bộ nhớ từ thông tin khác trong lệnh. Ví dụ: toán hạng thanh ghi có thể xác định độ rộng, nhưng điều này có thể không rõ ràng nếu chỉ có bộ nhớ và giá trị tức thời. Trình biên dịch mã không nên đoán chiều rộng mơ hồ một cách ngẫu nhiên.

Toán hạng bộ nhớ bao gồm một ký hiệu trong AMD64 về cơ bản là RIP-relative. Chọn phương thức đánh địa chỉ rõ ràng rel/abs. escape không xác định theo nghĩa đen là một lỗi.

## Phần và căn chỉnh

|Mục Nội dung|Hành vi căn chỉnh|
| --- | --- |
|mã|NOP Chèn lệnh|
|dữ liệu khởi tạo|Chèn 0 byte|
| BSS |Tăng kích thước bộ nhớ logic mà không cần thêm tệp payload|

Có sự khác biệt giữa kích thước tập tin và kích thước bộ nhớ. BSS dự trữ bộ nhớ, nhưng không yêu cầu lưu trữ cùng kích thước 0 byte trong tệp đối tượng. Các phần tùy chỉnh có các thuộc tính và ký hiệu duy trì thông tin ràng buộc và loại.

### Payload Dự trữ BSS không phân bổ

Lưu phần sau dưới dạng `buffer.asm`. Logic BSS dự trữ 1 TiB, nhưng không gán hoặc ghi 1 TiB trong quá trình lắp ráp.

```asm
section .bss
global buffer
buffer:
    resb 1099511627776
buffer_end:

section .text
    ret
```

```sh
whale asm --amd64 buffer.asm -o buffer.o
```

Giá trị bên trong phần `buffer` là 0 và giá trị của `buffer_end` là 1099511627776. Tiêu đề `.bss` là `SHT_NOBITS` với kích thước đó và không có tệp payload. Sau đó, khi bạn quay lại `.bss`, hãy tiếp tục sử dụng logic offset. Chỉ thị dữ liệu có số 0 cũng làm tăng kích thước logic. Giá trị ban đầu khác 0, việc di chuyển được áp dụng trong BSS và các lệnh trong BSS đều bị từ chối.

Sự khác biệt tương tự có thể được áp dụng cho các đối tượng và trình liên kết API.

```rust
use linker::core::layout::Layout;
use object::{ObjectFile, ObjectFormat, SectionKind};

fn main() {
    let mut object = ObjectFile::new(ObjectFormat::ELF64);
    let text = object.add_section(".text", SectionKind::Text, 16);
    object.sections[text].data = vec![0xc3];
    let bss = object.add_section(".bss", SectionKind::Bss, 16);
    object.sections[bss].zero_fill = 1 << 40;
    assert!(object.sections[bss].data.is_empty());

    let elf = object.write_with_limit(4096).unwrap();
    assert!(elf.len() < 4096);
    let layout = Layout::compute(&[object], 0x1000).unwrap();
    assert_eq!(layout.file_size, 1);
    assert_eq!(layout.memory_size, (1 << 40) + 16);
    assert_eq!(layout.sections[bss].memory_address, 0x1010);
    assert_eq!(layout.sections[bss].file_size, 0);
}
```

```text
section  memory address  file bytes       memory bytes
.text    0x1000          1                1
.bss     0x1010          0                1099511627776
```

`Section::zero_fill` là số byte BSS bổ sung không được lưu trữ trong tệp. Kích thước bộ nhớ được kiểm tra là `data.len() + zero_fill`. BSS `data` không đệm hiện có cũng được chấp nhận, nhưng có thể tránh được việc gán nó bằng một vectơ `data` trống. Các phần không phải BSS phải là `zero_fill == 0`. Không gian lưu trữ vật lý 0 không bao hàm trạng thái khởi tạo của IR.

`Layout::compute` trả về `Result`, ghi lại sự tương ứng của đối tượng/phần đầu vào, căn chỉnh, tệp offset, địa chỉ bộ nhớ và cả hai kích thước của tất cả các phần. Địa chỉ và căn chỉnh số học được kiểm tra, thứ tự đầu vào được giữ nguyên và căn chỉnh đối tượng 0 được coi là không bị ràng buộc (căn chỉnh 1). BSS không di chuyển con trỏ tập tin. Đây là quá trình triển khai payload, với tiêu đề thực thi, load segment với quyền truy cập, việc áp dụng di dời là một nhiệm vụ riêng biệt.

ELF writer từ chối overflow và cắt bớt khi giảm độ rộng trường. Số phần mở rộng không được hỗ trợ và tổng số tiêu đề bao gồm bảng đã tạo và các phần sắp xếp lại phải nhỏ hơn `0xff00`. Giới hạn mặc định cho đầu ra được tuần tự hóa là 256 MiB. Bạn có thể chỉ định giới hạn byte bao gồm bảng padding· với `ObjectFile::write_with_limit` hoặc `write_elf_with_limit` và giới hạn này không bao gồm kích thước bộ nhớ BSS không được lưu trữ trong tệp. Kích thước đầu ra được kiểm tra trước khi phân bổ vectơ byte cuối cùng.

## Nhận dạng biểu tượng

Các hàm và biến sử dụng các mã định danh khác nhau bên trong IR. Các kết nối bên ngoài sử dụng `link_name` do giao diện người dùng chỉ định. Whale không tự động đổi tên một trong những biểu tượng công cộng xung đột mà vẫn giữ nguyên tên đã nêu.

Bảng khai báo IR ghi lại các tham chiếu đã nhập `FunctionId` và các giá trị hàm rõ ràng `link_name`. Các API trình biên dịch/đối tượng/trình liên kết bên dưới vẫn là các giao diện riêng biệt: quá trình phát xạ IR gốc chưa mang các danh tính đó đến liên kết cuối cùng. IR riêng việc xác minh cuộc gọi không thiết lập được thuộc tính đầu cuối này.

Do đó, có thể cả hàm bên trong và biến đều có tên `item`, nhưng sẽ là sai sót nếu hiển thị cả hai có cùng tên bên ngoài. Việc tách không gian tên bên trong không tự động tách không gian tên bên ngoài.

Phạm vi của một ký hiệu cục bộ đối tượng là đối tượng đầu vào của nó. Các biểu tượng toàn cầu tham gia giải thích giữa các đối tượng. Xác nhận xung đột chức năng/dữ liệu là lỗi. Biểu tượng NOTYPE duy trì khả năng tương thích với các đầu vào không cung cấp loại cụ thể hơn. Việc nó không có kiểu không có nghĩa nó là một hàm hoặc dữ liệu.

## chọn độ nét

|định nghĩa hoặc tài liệu tham khảo|kết quả|
| --- | --- |
|Strong và strong|lỗi định nghĩa trùng lặp|
|Strong và weak|Strong Chọn định nghĩa|
|Weak và weak|Chọn định nghĩa đầu tiên theo thứ tự đầu vào|
|Xem strong chưa được giải quyết|lỗi liên kết|
|Xem chưa được giải quyết weak|Lỗi không được hỗ trợ trong hồ sơ tĩnh ban đầu|

Nếu có nhiều định nghĩa weak, thứ tự nhập chúng sẽ ảnh hưởng đến kết quả. Bạn phải sử dụng thứ tự đầu vào đã chuyển một cách nhất quán cho các liên kết xác định.

## Đầu ra thực thi tĩnh

Cấu hình native tĩnh tạo ra ELF ET_EXEC chỉ định điểm vào. Điểm vào không được suy ra từ tên hàm `main`. Nó cũng không tự động chèn mã khởi động gọi hàm đó.

Nó không tự động loại bỏ các phần, hợp nhất mã giống hệt nhau hoặc loại bỏ các ký hiệu. Vị trí tệp yêu cầu tính toán riêng biệt số byte bạn thực sự lưu và bộ nhớ bạn dự trữ trong thời gian chạy.

Hiện chưa có đường dẫn tạo tệp thực thi tĩnh hoàn chỉnh cho CLI. `whale asm` tạo một đối tượng có thể định vị lại và `whale object` gói các byte thô vào một đối tượng. Vui lòng tham khảo [Tổng quan về chuỗi công cụ](overview) để biết tình trạng sẵn có, ABI và [AMD64 Mục tiêu](amd64-target) để biết yêu cầu.


## Tuần tự hóa bản ghi Wave có thể lựa chọn

Các bản dựng Linux Whale cho x86_64 hiện có thể sử dụng triển khai Wave cho các bản ghi tiêu đề ELF64 cố định · phần · biểu tượng ·RELA. Triển khai Rust mặc định cũng được cung cấp. Việc lựa chọn mục tiêu đầu ra, xác minh đối tượng, vị trí, giải thích ký hiệu và phân bổ bộ đệm được xử lý bởi Rust. Việc chọn Wave không thêm các kiến ​​trúc được hỗ trợ hoặc trình liên kết đầy đủ.

Cài đặt Rust, LLVM 21 thư viện phát triển, trình liên kết C và `ar`, sau đó xây dựng đường dẫn đã chọn từ kho lưu trữ Whale.

```sh
git clone https://github.com/wavefnd/Wave.git /tmp/whale-wave-bootstrap
git -C /tmp/whale-wave-bootstrap checkout --detach 8a465e30aeea4b817d925cdd0e8d08c1bb029c9a
python3 tools/build_wave_elf.py --wave-source /tmp/whale-wave-bootstrap --out-dir /tmp/whale-wave-elf
WHALE_WAVE_ELF_DIR=/tmp/whale-wave-elf cargo build --locked --all-features
WHALE_WAVE_ELF_DIR=/tmp/whale-wave-elf cargo test --locked --workspace --all-features
```

Tập lệnh kiểm tra sự cố định revision, từ chối thay đổi thành tracked, xây dựng trình biên dịch Wave, tạo đối tượng Wave với LLVM và gói nó với archive để liên kết tĩnh. `WHALE_WAVE_ELF_DIR` nên có cái này archive. Việc yêu cầu rõ ràng archive không hợp lệ hoặc máy chủ không được hỗ trợ là lỗi xây dựng. Các bản dựng thông thường không có biến nào được chỉ định không yêu cầu trình biên dịch `--all-features` cũng như Wave. Trình biên dịch Wave cũng không cần thiết trong thời gian chạy đối với tệp thực thi Whale được liên kết. Đường dẫn biên dịch chéo Whale với bootstrap chưa được hỗ trợ.

Ví dụ: lưu thông tin sau dưới dạng `return.asm` và tập hợp nó thành tệp nhị phân được tạo:

```asm
section .text
global entry
entry:
    ret
```

```sh
target/debug/whale asm --amd64 return.asm -o return.o
```

```text
ELF class: ELF64
machine: AMD64 (62)
.text bytes: c3
```

Bản ghi ABI mang loại, con trỏ trường u64, số trường, con trỏ đầu ra và dung lượng. Bộ đệm thuộc sở hữu của người gọi Rust. Nó không chuyển quyền sở hữu chuyển nhượng hoặc biểu thức Rust enum/String/Vec qua các ranh giới. Quy trình Wave kiểm tra số lượng, dung lượng và độ rộng trường trước khi ghi. Thành công trả về trạng thái 0, loại/con trỏ/dung lượng không hợp lệ trả về 1 và trường overflow trả về 2. Con trỏ phải trỏ đến các bộ đệm còn hoạt động và có kích thước chính xác và không chồng chéo lên nhau. Chỉ riêng con trỏ C thô không thể chứng minh điều kiện này; wrapper thỏa mãn điều đó. Quá trình kiểm tra so sánh toàn bộ ELF, bao gồm BSS và signed di dời addend, với đường dẫn Rust. Đây là quá trình triển khai một phần Wave sử dụng Wave/LLVM bootstrap và không tự lưu trữ hoàn toàn.
