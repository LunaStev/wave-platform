---
translation_set_id: whale-assembler-linker
path: whale/assembler-linker
locale: zh
group: whale
group_order: 1
order: 11
title: 汇编与静态链接
summary: 描述操作数编码、节放置、符号绑定和执行入口点。
---

## 组件和对象

Whale汇编器将AMD64指令转换为机器代码字节和重定位信息。 ELF64 对象包含此信息和部分/符号。汇编由Whale自己的实现执行，不需要外部汇编器。

可重定位对象可能具有最终地址未知的引用。解析这个地址的过程就是一个链接。成功的汇编不应推断所有外部符号已被解析或可执行文件已被创建。

## 组装函数

将以下代码保存为`answer.asm`。

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

`.text`的内容是`b8 2a 00 00 00 c3`，`mov eax, 42`后面是`ret`。对象 ELF64 公开 `answer`。它是一个没有进程启动代码的可调用函数，不是可执行文件。本例中的指令可以由当前的汇编器处理。

## 使用 Rust API 构造对象

这是将相同功能字节写入板条箱`object`的完整示例。指定输出目标和全局符号。

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

Assertion 检查 ELF 类、字节顺序、machine 标识符。 `value: 0` 是offset 内的`.text`，`size: 6` 是符号的字节大小。无效的节引用或范围是序列化错误。拒绝其他machine或字节顺序，而不将它们标记为AMD64。

## 解释两个物体的符号

目前`linker` crate提供符号解释。以下可执行示例在两个对象上定义本地 `helper`，然后检查如果相同的名称发布两次是否会发生错误。

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

输出：

```text
Duplicate global symbol: helper
```

这两个本地定义具有不同的`object_index`的单独键。这两个全局定义是冲突的。这个例子直接解释了你在内存中构造的对象的符号。 `.o` 不执行读取文件、应用重定位或输出可执行文件。在下面描述的完整定义选择策略中，weak优先级等尚未实施。

## 文字和内存操作数

文字保留其宽度和标牌，直到检查实际指令的编码范围。超出范围的值应该导致错误，而不是被默默地截断。

当无法从命令中的其他信息确定内存宽度时，需要明确的大小符号。例如，寄存器操作数可以确定宽度，但是如果仅存在存储器和立即值，则这可能是不明确的。汇编器不应随机猜测不明确的宽度。

AMD64中由一个符号组成的内存操作数基本上是RIP-relative。选择显式寻址方法rel/abs。文字中未知的 escape 是一个错误。

## 剖面和对齐方式

|栏目内容|对齐行为|
| --- | --- |
|代码|NOP 插入命令|
|初始化数据|插入0字节|
| BSS |增加逻辑内存大小而不添加文件payload|

文件大小和内存大小之间存在区别。 BSS保留内存，但不要求目标文件中存储相同大小的0字节。自定义部分具有维护绑定和类型信息的属性和符号。

### Payload 预留 BSS 未分配

将以下内容另存为`buffer.asm`。逻辑BSS保留1个TiB，但在汇编时不会赋值或写入1个TiB。

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

`buffer`段内的值为0，`buffer_end`的值为1099511627776。`.bss`标头是具有该大小的`SHT_NOBITS`，并且没有文件payload。之后，当您返回`.bss`时，继续使用逻辑offset。带零的数据指令也会增加逻辑大小。非零初始值、BSS内应用的重定位以及BSS内的命令将被拒绝。

同样的区别也可用于对象和链接器API。

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

`Section::zero_fill` 是未存储在文件中的额外 BSS 字节数。测试的内存大小为`data.len() + zero_fill`。现有的零填充 BSS `data` 也被接受，但可以使用空 `data` 向量来避免其分配。除 BSS 之外的部分必须为 `zero_fill == 0`。物理0存储空间并不意味着IR的初始化状态。

`Layout::compute`返回`Result`，记录了输入对象/节的对应关系、对齐方式、文件offset、内存地址以及所有节的大小。检查地址和对齐算法，保留输入顺序，并将对象对齐 0 视为不受约束（对齐 1）。 BSS 不移动文件光标。这是 payload 部署，具有可执行标头，load segment 具有访问权限，应用重定位是一项单独的任务。

ELF writer 拒绝 overflow 以及减小字段宽度时的截断。不支持扩展节号，并且包含创建的表和重排节的表头总数必须小于`0xff00`。序列化输出的默认限制是 256 MiB。您可以使用`ObjectFile::write_with_limit`或`write_elf_with_limit`指定包括padding·表的字节限制，并且该限制不包括未存储在文件中的BSS内存大小。在分配最终字节向量之前检查输出大小。

## 符号识别

函数和变量在IR内使用不同的标识符。外部连接使用前端指定的`link_name`。 Whale 不会自动重命名冲突的公共符号之一，保留其声明的名称。

IR声明表记录类型化的`FunctionId`引用和显式函数`link_name`值。下面的汇编器/对象/链接器 API 仍然是单独的接口：本机 IR 发射尚未将这些身份传递到最终链接。 IR 单独的调用验证并不能建立这种端到端属性。

因此，内部函数和变量都可以具有名称`item`，但使用相同的外部名称公开两者将是错误的。分离内部命名空间不会自动分离外部命名空间。

对象局部符号的范围是其输入对象。全局符号参与对象之间的解释。已确认的功能/数据冲突属于错误。 NOTYPE 符号保持与不提供更具体类型的输入的兼容性。它没有类型这一事实并不意味着它是函数或数据。

## 选择定义

|定义或参考|结果|
| --- | --- |
|Strong 和 strong|重复定义错误|
|Strong 和 weak|Strong 选择定义|
|Weak 和 weak|选择输入顺序中的第一个定义|
|请参阅strong未解决|链接错误|
|查看未解决weak|初始静态配置文件不支持错误|

如果有多个weak定义，它们的输入顺序将影响结果。对于确定性链接，您必须一致地使用传递的输入顺序。

## 静态可执行输出

静态 native 配置文件生成指定入口点的 ELF ET_EXEC。入口点不是从函数名称`main`推断出来的。它也不会自动插入调用该函数的启动代码。

它不会自动删除部分、合并相同的代码或删除符号。文件放置需要单独计算实际保存的字节和运行时保留的内存。

CLI 的完整静态可执行文件创建路径尚不可用。 `whale asm` 创建一个可重定位对象，并且 `whale object` 将原始字节包装到一个对象中。请参阅[工具链概述](overview)了解可用性，ABI和[AMD64 目标](amd64-target)了解要求。


## 可选Wave记录序列化

Linux Whale x86_64 的构建现在可以使用 Wave 实现来固定 ELF64 header·section·symbol·RELA 记录。还提供了默认的 Rust 实现。输出目标选择、对象验证、放置、符号解释和缓冲区分配由Rust处理。选择 Wave 不会添加支持的架构或完整链接器。

安装 Rust、LLVM 21 个开发库、C 链接器和 `ar`，然后从 Whale 存储库构建选定的路径。

```sh
git clone https://github.com/wavefnd/Wave.git /tmp/whale-wave-bootstrap
git -C /tmp/whale-wave-bootstrap checkout --detach 8a465e30aeea4b817d925cdd0e8d08c1bb029c9a
python3 tools/build_wave_elf.py --wave-source /tmp/whale-wave-bootstrap --out-dir /tmp/whale-wave-elf
WHALE_WAVE_ELF_DIR=/tmp/whale-wave-elf cargo build --locked --all-features
WHALE_WAVE_ELF_DIR=/tmp/whale-wave-elf cargo test --locked --workspace --all-features
```

该脚本检查固定revision，拒绝对tracked的更改，构建Wave编译器，使用LLVM创建Wave对象，并将其与archive捆绑在一起以进行静态链接。 `WHALE_WAVE_ELF_DIR` 应该有这个archive。明确请求无效的 archive 或不受支持的主机是构建错误。未指定变量的正常构建不需要 `--all-features` 或 Wave 编译器。链接的 Whale 可执行文件在运行时也不需要 Wave 编译器。尚不支持将 Whale 本身与 bootstrap 交叉编译的路径。

例如，将以下内容保存为 `return.asm` 并将其汇编到生成的二进制文件中：

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

记录ABI携带类型、u64字段指针、字段数量、输出指针和容量。该缓冲区由调用者Rust拥有。它不会跨边界传递赋值所有权或 Rust enum/String/Vec 表达式。 Wave 例程在写入之前检查计数、容量和字段宽度。成功返回状态 0，无效类型/指针/容量返回 1，字段 overflow 返回 2。指针必须指向活动的、大小正确的缓冲区，并且彼此不重叠。原始的C指针本身无法证明这个条件； wrapper满足。该测试将整个 ELF（包括 BSS 和 signed 重定位 addend）与 Rust 路径进行比较。这是使用 Wave/LLVM bootstrap 的部分 Wave 实现，并且不是完全自托管的。
