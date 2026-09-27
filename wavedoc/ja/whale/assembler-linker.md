---
translation_set_id: whale-assembler-linker
path: whale/assembler-linker
locale: ja
group: whale
group_order: 1
order: 11
title: アセンブリと静的リンク
summary: オペランドエンコーディング、セクション配置、シンボルバインディング、および実行エントリポイントについて説明します。
---

## アセンブリとオブジェクト

WhaleアセンブラはAMD64コマンドを機械語バイトと再配置情報に変換します。 ELF64 オブジェクトはこの情報とセクション・シンボルを含みます。アセンブルはWhale独自の実装で行われ、外部アセンブラは必要ありません。

再配置可能オブジェクトには、最終アドレスがまだわからない参照がある場合があります。このアドレスを解決するプロセスはリンクです。アセンブルが成功したとすべての外部シンボルを解決したか、実行可能ファイルを作成したと判断してはいけません。

## 関数をアセンブルする

次のコードを`answer.asm`として保存します。

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

`.text`内容は`b8 2a 00 00 00 c3`で、`mov eax, 42`の後に`ret`が続きます。 ELF64オブジェクトは`answer`を公開します。プロセスの開始コードがない呼び出し可能な関数であり、実行可能ファイルではありません。この例のコマンドは、現在のアセンブラで処理できます。

## Rust APIでオブジェクトを構成する

以下は、`object`クレートで同じ関数バイトを書き込む完全な例です。出力ターゲットとグローバルシンボルを指定します。

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

AssertionはELFクラス、バイト順、machine識別子をチェックします。 `value: 0`は`.text`内部offsetであり、`size: 6`はシンボルのバイトサイズです。無効なセクション参照または範囲は直列化エラーです。他のmachineやバイトシーケンスもAMD64で表示せずに拒否します。

## 2つのオブジェクトのシンボルを解釈する

現在、`linker`クレートはシンボル解釈を提供します。次の実行可能な例では、2つのオブジェクトにそれぞれローカル`helper`を定義し、同じ名前を2回公開するとエラーが発生することを確認します。

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

出力：

```text
Duplicate global symbol: helper
```

2つのローカル定義は、`object_index`が異なる別々のキーを持ちます。両方のグローバル定義は競合します。この例では、メモリーに構成されたオブジェクトのシンボルを直接解釈します。 `.o`ファイルの読み取り、再配置の適用、実行ファイルの出力は行われません。以下に説明する完全な定義選択ポリシーのうち、weak優先順位などはまだ実装されていません。

## リテラルとメモリオペランド

リテラルは、実際のコマンドのエンコード範囲を調べるまで幅と符号を保持します。範囲に入らない値は静かに切り捨てられず、エラーになるはずです。

命令の他の情報でメモリ幅を決定できない場合は、明示的なサイズ表記が必要です。たとえば、レジ​​スタオペランドは幅を指定できますが、メモリと即値がある場合はあいまいです。アセンブラがあいまいな幅を任意に推測してはいけません。

AMD64でシンボル一つからなるメモリオペランドは基本的にRIP-relativeです。明示的なrel/absでアドレス指定方法を選択します。リテラル内の不明なescapeはエラーです。

## セクションとアライメント

|セクション内容|整列動作|
| --- | --- |
|コード|NOP コマンド挿入|
|初期化されたデータ|0バイト挿入|
| BSS |ファイルpayloadを追加せずに論理メモリサイズを増やす|

ファイルサイズとメモリサイズは区別されます。 BSSはメモリを予約しますが、同じサイズの0バイトをオブジェクトファイルに保存する必要はありません。カスタムセクションは属性を持ち、シンボルはバインディングとタイプ情報を保持します。

### Payload 割り当てなし BSS 予約する

以下を`buffer.asm`として保存します。ロジックBSS1TiBを予約しますが、アセンブル中に1TiBを割り当てたり記録したりしません。

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

`buffer`のセクション内部値は0で、`buffer_end`の値は1099511627776です。 `.bss`ヘッダーはそのサイズの`SHT_NOBITS`であり、ファイルpayloadはありません。以降、`.bss`に戻ると、論理offsetを引き続き使用します。ゼロのデータディレクティブも論理サイズを増やします。ゼロ以外の初期値、BSS内に適用する再配置、BSS内のコマンドは拒否します。

オブジェクトとリンカーAPIでも同じ区切りを使用できます。

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

`Section::zero_fill`は、ファイルに保存しない追加のBSSバイト数です。検査されたメモリサイズは`data.len() + zero_fill`です。既存のゼロで埋められたBSS`data`も許可しますが、空の`data`ベクトルでその割り当てを避けることができます。 BSS以外のセクションは`zero_fill == 0`でなければなりません。物理的なゼロストレージスペースがIRの初期化状態を意味するわけではありません。

`Layout::compute`は`Result`を返し、すべてのセクションの入力オブジェクト・セクション対応、ソート、ファイルoffset、メモリアドレスと2つのサイズを記録します。アドレス・整列算術を検査し、入力順序を維持し、オブジェクト整列 0 は制約なし (整列 1) として処理します。 BSSはファイルカーソルを移動しません。これは payload バッチで、実行可能ファイルのヘッダー、アクセス権を持つ load segment、再配置の適用は別のタスクです。

ELF writerはoverflowとフィールド幅を減らしたときの切り捨てを拒否します。拡張セクション番号は未サポートであり、生成された表・再配置セクションまでのヘッダーの総数は、`0xff00`未満でなければなりません。直列化された出力のデフォルト制限は256MiBです。 `ObjectFile::write_with_limit`または`write_elf_with_limit`でpadding・表を含むバイト制限を指定でき、ファイルに保存しないBSSメモリサイズは制限に含まれません。最終バイトベクトルを割り当てる前に出力サイズを確認します。

## シンボル識別

関数と変数は、IR内で異なる識別子を使用します。外部接続は、フロントエンドが指定した`link_name`を使用します。 Whaleは、競合する公開シンボルの1つの名前を自動的に変更せずに、指定された名前を保持します。

IR 宣言テーブルには、型指定された `FunctionId` 参照と明示的な関数 `link_name` の値が記録されます。以下のアセンブラー/オブジェクト/リンカー API は別個のインターフェイスのままです。ネイティブ IR エミッションは、これらの ID をまだ最終リンクまで伝えていません。 IR 通話検証だけでは、このエンドツーエンドの特性は確立されません。

したがって、内部関数と変数の名前の両方が`item`であることは可能ですが、両方を同じ外部名で公開するとエラーになる可能性があります。内部ネームスペースが分離されていても、外部ネームスペースも自動的に分離されません。

オブジェクトローカルシンボルの範囲は対応する入力オブジェクトです。グローバルシンボルはオブジェクト間の解釈に参加します。確認された関数・データ衝突はエラーです。 NOTYPEシンボルは、より具体的なタイプを提供しない入力との互換性を維持します。タイプがないという事実だけでは関数やデータとは断定しません。

## 定義の選択

|定義または参照|結果|
| --- | --- |
|Strongとstrong|重複定義エラー|
|Strongとweak|Strong定義の選択|
|Weakとweak|入力順序で最初の定義を選択|
|未解決 strong参照|リンクエラー|
|未解決 weak参照|初期静的プロファイルでは、未サポートエラー|

複数のweak定義がある場合、入力順序は結果に影響します。決定的なリンクのために渡される入力シーケンスを一貫して使用する必要があります。

## 静的実行可能ファイルの出力

静的nativeプロファイルは、エントリポイントを指定したELFET_EXECを作成します。 `main`という関数名からエントリポイントを推論しません。その関数を呼び出す開始コードも自動挿入されません。

セクションの削除、同じコードの結合、シンボルの削除は自動的には行われません。ファイル配置は、実際に格納するバイトと実行時に予約するメモリを別々に計算する必要があります。

CLIの完全な静的実行可能ファイル生成パスはまだ提供されていません。 `whale asm`は再配置可能なオブジェクトを作成し、`whale object`は生のバイトをオブジェクトにラップします。提供の有無は[ツールチェーンの概要](overview)、ABI要件は[AMD64ターゲット](amd64-target)をご覧ください。


## 選択可能なWaveレコードのシリアライゼーション

Linux x86_64用 Whale ビルドは固定 ELF64ヘッダ・セクション・シンボル・RELAレコードにWave実装を使用できます。基本的なRust実装も提供されています。出力ターゲットの選択、オブジェクトの検証、配置、シンボルの解釈、およびバッファの割り当ては、Rustが担当します。 Waveを選択すると、サポートアーキテクチャや完全なリンカが追加されるわけではありません。

Rust、LLVM21開発ライブラリ、Cリンカーと`ar`をインストールし、Whaleリポジトリから選択パスを構築します。

```sh
git clone https://github.com/wavefnd/Wave.git /tmp/whale-wave-bootstrap
git -C /tmp/whale-wave-bootstrap checkout --detach 8a465e30aeea4b817d925cdd0e8d08c1bb029c9a
python3 tools/build_wave_elf.py --wave-source /tmp/whale-wave-bootstrap --out-dir /tmp/whale-wave-elf
WHALE_WAVE_ELF_DIR=/tmp/whale-wave-elf cargo build --locked --all-features
WHALE_WAVE_ELF_DIR=/tmp/whale-wave-elf cargo test --locked --workspace --all-features
```

スクリプトは固定リビジョンを確認し、tracked 変更を拒否し、Wave コンパイラをビルドした後、LLVM で Wave オブジェクトを生成し、静的リンク用 archive で囲みます。 `WHALE_WAVE_ELF_DIR`にはこのアーカイブが必要です。誤ったアーカイブまたは未サポートのホストで明示的に要求するとビルドエラーです。変数を指定しない一般的なビルドでは、 `--all-features`でもWaveコンパイラは必要ありません。リンクされたWhale実行可能ファイルの実行時点でもWaveコンパイラは必要ありません。このブートストラップでWhale自体をクロスコンパイルするパスはまだ未サポートです。

たとえば、次を`return.asm`として保存し、生成されたバイナリにアセンブルします。

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

レコードABIは種類、u64フィールドポインタ、フィールド数、出力ポインタ、容量を渡します。バッファはRust呼び出し側が所有します。割り当ての所有権またはRustenum/String/Vec表現を境界を越えて渡しません。 Wave ルーチンは記録前に個数・容量・フィールド幅を検査します。成功は状態0、誤った型・ポインタ・容量は1、フィールドoverflowは2を返します。ポインタは正しいサイズで生きていて、互いに重ならないバッファを指す必要があります。生のCポインタだけではこの条件を証明することはできず、wrapperが条件を満たします。テストは、BSSとsigned再配置addendまでを含む完全なELFをRustパスと比較します。これは、Wave/LLVMbootstrapを使用する部分Wave実装であり、完全なセルフホスティングではありません。
