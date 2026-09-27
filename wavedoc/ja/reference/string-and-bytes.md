---
translation_set_id: string-bytes
path: reference/string-and-bytes
locale: ja
group: stdlib
group_order: 1
order: 3
title: string: 長さ、検索、範囲
summary: NUL終了文字列のバイト単位APIと戻り値を説明します。
---

## 文字列の格納と引数条件

このモジュールの`str`引数は、アクセス可能なNUL終了バイト列でなければなりません。長さと検索インデックスはバイト単位です。一般文字はUTF-8として保存されますが、バイト検索はUnicode正規化や文字単位の分割は行いません。返されたインデックスが文字境界であると仮定しないでください。

## 長さと比較

```text
std::string::len
len(s: str) -> i32
is_empty(s: str) -> bool

std::string::cmp
eq(a: str, b: str) -> bool
cmp(a: str, b: str) -> i32
starts_with(s: str, prefix: str) -> bool
ends_with(s: str, suffix: str) -> bool
```

`len`は最後のNULを除外します。 `cmp`結果の符号で順序を判断します。戻り値を Unicode 文字順や言語別の事前ソートとして解釈しません。これらの関数はメモリを割り当てずに入力を変更しません。

## 検索

`import("std::string::find")::{find, contains, count};`のように必要な名前を取得します。

|関数宣言|結果|
| --- | --- |
| `find(s: str, needle: str) -> i32` |最初の一致位置。なければ-1、空のneedleは0|
| `contains(s: str, needle: str) -> bool` |含めるかどうか。空 needleはtrue|
| `count(s: str, needle: str) -> i32` |重複しない一致数。空needleは0|
| `find_char(s: str, c: u8) -> i32` |バイトの最初の位置または-1|
| `rfind_char(s: str, c: u8) -> i32` |バイトの最後の位置または-1|
| `contains_char(s: str, c: u8) -> bool` |そのバイトの存在|
| `count_char(s: str, c: u8) -> i32` |対応するバイト数|

`*_char`という名前の`c`は、Unicodeコードポイントではなく1バイトです。文字列の末尾のNUL自体は検索対象には含まれません。

## スペースを除く範囲

```text
std::string::trim
trim_left_index(s: str) -> i32
trim_right_index(s: str) -> i32
trim_range(s: str, out_start: ptr<i32>, out_end: ptr<i32>)
```

`trim_range`は、ASCII空白を除いた半開範囲`[start, end)`を出力係数に記録します。両方の出力ポインタは書き込み可能な整数を指す必要があります。原文を変更したり、新しい文字列を作成したりしません。すべて空白の場合は空の範囲になります。

## 実行例

`main.wave`に保存し、`wavec run main.wave`を実行します。

<!-- wave-example: string-api -->
```wave
import("std::string::find")::{
    find, count
};
import("std::string::trim")::{
    trim_range
};

fun main() {
    var start: i32 = 0;
    var end: i32 = 0;
    trim_range("  Wave  ", &start, &end);
    println("{} {}", start, end);
    println("{} {}", find("banana", "na"), count("aaaa", "aa"));
}
```

実行結果：

```text
2 6
2 2
```

## 関連機能

`std::string::ascii`の分類・大小文字変換はASCII範囲用です。 `std::string::hash`の`djb2_32`と`fnv1a_64`を暗号ハッシュやパスワード保存には使用しません。 NULを含むデータには[bytes](/docs/ja/stdlib/bytes)を使用してください。

## 空のクエリと重なるパターン

検索関数のエッジの振る舞いを実際の値として調べると、呼び出し条件を決めるのは簡単です。

<!-- wave-example: book-string-empty-search -->
```wave
import("std::string::find")::{find, contains, count};

fun main() {
    println("empty find={}", find("Wave", ""));
    println("empty count={}", count("Wave", ""));
    println("nonoverlapping={}", count("aaaa", "aa"));

    if (contains("Wave", "")) {
        println("empty needle is contained");
    }
}
```

実行結果：

```text
empty find=0
empty count=0
nonoverlapping=2
empty needle is contained
```

findの成功値0は最初の位置です。 countの成功値0は一致しないか、空のクエリルールの結果です。両方の値を同じ条件で処理しません。大文字と小文字を無視するか、Unicode正規化が必要な場合は、このバイト検索の前後に別々のポリシーを実装する必要があります。

## trim 範囲を新しい文字列にコピーする

trim_rangeが返す範囲には終了NULが新しくなりません。別の宛先にコピーする場合は、長さ+1のスペースを確保し、最後のバイトを直接0に書き込みます。

<!-- wave-example: book-trim-copy -->
```wave
import("std::string::trim")::{trim_range};

fun main() -> i32 {
    var original: str = "  Wave  ";
    var start: i32 = 0;
    var end: i32 = 0;
    var destination: array<u8, 16>;

    trim_range(original, &start, &end);

    var length: i32 = end - start;

    if (length >= 16) {
        return 1;
    }

    for (var index: i32 = 0; index < length; index += 1) {
        destination[index] = original[start + index];
    }

    destination[length] = 0;
    println("{}", &destination[0] as str);
    return 0;
}
```

実行結果：

```text
Wave
```

長さ16の文字列はこの宛先には入りません。最後のNULまで必要だからです。長さが 0 であっても destination[0]=0 を記録すると有効な空文字列になります。目的地エリア配列はmainの終わりまで生きているので、その中から出力します。

## 文字列API使用順序

文字列 API を設計するときは、入力が NUL で終了するかどうか、インデックスがバイト数をカウントするかどうか、結果がソース範囲を借用するか新しい割り当てを所有するかを指定します。借用範囲はソースの存続期間によって異なります。割り当てられた結果では、誰がそれを解放するかを指定する必要があります。

基礎概念は[文字列学習の章](/docs/ja/language/strings)、NULを含むデータは[bytes](/docs/ja/stdlib/bytes)を読んでください。
