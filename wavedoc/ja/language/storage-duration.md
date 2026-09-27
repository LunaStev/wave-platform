---
translation_set_id: storage-duration
path: language/storage-duration
locale: ja
group: language
group_order: 2
order: 18
title: 保存期間と可変性
summary: var、constとstaticの範囲と書き込み可能性を区別します。
---

## 宣言による意味

|フォーマット|許容位置|再代入|用途|
| --- | --- | --- | --- |
| `var` |関数・ブロック|可能|一般的な可変ローカル変数|
| `const` |最上位|不可|グローバル定数宣言|
| `static` |最上位|可能|プログラム寿命中に存在する静的ストレージ宣言|

```wave
const PAGE_SIZE: i32 = 4096;
static request_count: i64 = 0;

fun main() {
    var limit: i32 = 4;
    var current: i32 = 0;
    var retries: i32 = 0;

    current += 1;
    retries += 1;
    println("{} {} {}", limit, current, retries);
}
```

## ローカル宣言の規則

```wave
var value: i32 = 1;
value = 2;
```

`var`で宣言したローカル変数に新しい値を代入できます。プログラム全体で使用する定数は、最上位で`const`と宣言します。

## const と static のローカルでの使用

`const`と`static`は最上位宣言です。関数本文と`for`初期化では、地域宣言である`var`を使用します。

## 寿命とポインター

ローカル変数のアドレスは`&`で取得できますが、ポインタが指すリポジトリの実際の有効期間は`ptr<T>`タイプによって追跡されません。ローカルリポジトリのアドレスを関数から引き渡すときは、そのアドレスが引き続き有効であることをプログラム構造で直接保証する必要があります。

## 学習と例の範囲

[プログラム全体で練習する](/docs/ja/getting-started/overview) · [標準ライブラリ](/docs/ja/stdlib)
