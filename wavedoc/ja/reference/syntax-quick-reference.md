---
translation_set_id: quick-reference
path: reference/syntax-quick-reference
locale: ja
group: reference
group_order: 5
order: 3
title: 構文クイックリファレンス
summary: よく使う宣言、制御フロー、型、ポインタとFFI文法を1ページにまとめます。
---

## 宣言

```wave
var value: i32 = 1;
var next: i32 = value + 1;
const LIMIT: i32 = 64;
static total: i64 = 0;
type Identifier = u64;
```

`var`は地域、`const`/`static`は最上位宣言です。ローカル変数は型を明示的に宣言します。

## 関数

```wave
fun max(left: i32, right: i32) -> i32 {
    if (left > right) {
        return left;
    }

    return right;
}
```

## ジェネリック

```wave
fun identity<T>(value: T) -> T {
    return value;
}

var value: i32 = identity<i32>(10);
```

ジェネリックコールは型引数を指定します。

## 構造体とenum

```wave
struct Pair {
    left: i32;
    right: i32;
}

enum Result -> i32 {
    Ok = 0,
    Error,
}
```

## 条件とループ

```wave
if (ready) {
    println("ready");
}

while (count < 10) {
    count += 1;
}

for (var i: i32 = 0; i < 10; i += 1) {
    println("{}", i);
}

match (status) {
    Ready => {
        println("ready");
    }
    0 => {
        println("zero");
    }
    _ => {
        println("other");
    }
}
```

`if`、`while`、`for`、`match`のヘッダーは括弧を使用します。

## 配列とポインタ

```wave
var values: array<i32, 4> = [1, 2, 3, 4];
var p: ptr<i32> = &values[0];
var first: i32 = deref p;
```

## コンソール入出力

```wave
print("value = ");
println("{}", value);
input("{}", value);
```

最初の引数は文字列リテラルです。正確な`{}`プレースホルダーごとに次の式が1つずつ必要であり、`input`対象は代入可能でなければなりません。

## importと公開アイテム

```wave
import("std::string::len");
import("./helpers" as helpers);
import("math")::{
    Vector
};

pub fun add(left: i32, right: i32) -> i32 {
    return left + right;
}

pub import("./extra"):: {
    increment
};
```

ローカルパスは`./`で始まります。エイリアスimportはモジュール名を指定し、選択importは必要なパブリックエントリをこのファイルの名前空間にインポートします。 `pub import`は選択した項目を再エクスポートします。

## FFI

```wave
extern(c) fun native_call(value: i32) -> i32;

export(c) fun wave_call(value: i32) -> i32 {
    return value + 1;
}
```

## 対象条件付き項目

```wave
#[target(os="linux", arch="riscv64")]
extern(c) fun platform_call(value: i32) -> i32;
```

サポート条件キーは`arch`、`os`、`env`、`abi`です。プロパティはすぐに次の最上位項目を制御します。

## インラインアセンブリ

```wave
var result: i64 = 0;
asm {
    "mv a0, a1"
    in("a1") 7
    out("a0") result
    clobber("memory")
}
```

コマンドテキストとレジスタ名はターゲットに依存します。ブロックに必要なすべての入力、出力、および隠れたclobberを宣言します。

## ソースチェック

```shell
wavec build main.wave --emit=check
wavec print supported-targets
wavec print supported-emit-kinds
```

## 学習と例の範囲

関数の外に別に表示したローカル変数・文例は、関数本文に入れるコードの断片です。完全な実行例と練習は[Wave学習コース](/docs/ja/getting-started/overview)に続きます。メモリと外部機能の詳細規則については、[標準ライブラリ](/docs/ja/stdlib)を確認してください。
