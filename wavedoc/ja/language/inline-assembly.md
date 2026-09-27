---
translation_set_id: assembly
path: language/inline-assembly
locale: ja
group: language
group_order: 2
order: 19
title: インラインアセンブリ
summary: asmブロックのコマンド文字列、in/outオペランドとclobber契約を説明します。
---

## asmブロック

`asm`は、ターゲットアーキテクチャのコマンドを直接挿入するための低レベルの文法です。

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

ブロック内の文字列リテラルは、アセンブリコマンドのリストに渡されます。

## 入力と出力

```wave
var result: i64 = 0;
asm {
    "mov rax, rdi"
    in("rdi") 123
    out("rax") result
}
```

- `in("reg") expression`はWaveの値を入力オペランドに連結します。
- `out("reg") target`は、出力値を代入可能なWaveターゲットに書き込みます。
- レジスタ名は文字列または識別子の形式で指定できます。

入力オペランドには変数、整数・文字列リテラル、`&identifier`、`deref identifier`と負数を使用できます。

## clobber

ブロックが明示的な出力以外のレジスタやメモリ状態を変更すると、`clobber(...)`に書き込まれます。

```wave
asm {
    "nop"
    clobber("rax", "rcx", "memory")
}
```

## 使用時に確認する

- 命令文法は、ターゲットアーキテクチャとLLVMインラインアセンブリ契約に適合する必要があります。
- 呼び出し規約上保存する必要があるレジスタをランダムに破壊しないでください。
- メモリを読み書きするブロックは、`memory`を含む必要なclobberを宣言してください。
- 可能であれば、アーキテクチャ固有のasmを小さな関数の後に分離します。

インラインアセンブリの動作と移植性は、言語タイプだけでは保証されません。

## 学習と例の範囲

[プログラム全体で練習する](/docs/ja/getting-started/overview) · [標準ライブラリ](/docs/ja/stdlib)
