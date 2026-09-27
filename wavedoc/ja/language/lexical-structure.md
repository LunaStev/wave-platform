---
translation_set_id: lexical
path: language/lexical-structure
locale: ja
group: language
group_order: 2
order: 14
title: 語彙構造
summary: 識別子、リテラル、区切り文字、キーワード、タイプ名を記述します。
---

## 識別子

識別子は変数、関数、型、フィールドなどに名前を付けます。名前は大文字と小文字を区別し、文字、数字、`_`を組み合わせることができます。最初の文字には数字を使用できません。 Unicode文字も識別子に使用できます。

```wave
var request_count: i64 = 0;
var 이름: str = "Wave";
```

実際のプロジェクトでは、ツールの互換性と検索の利便性のために一貫した名前ルールを設定して使用することをお勧めします。

## 文章と区切り文字

ほとんどの宣言と式の文は`;`で終わります。関数・条件文・繰り返し文・構造体のように本文を持つ構文は `{ ... }`ブロックを使用します。

```wave
var answer: i32 = 42;

fun double(value: i32) -> i32 {
    return value * 2;
}
```

## リテラル

```wave
var integer: i32 = 42;
var decimal: f64 = 3.14;
var text: str = "Wave";
var letter: char = 'W';
var enabled: bool = true;
var address: ptr<u8> = null;
```

整数・浮動小数点・文字列・文字・ブルリオンと`null`リテラルを使用できます。 `null`はポインタ値に使用してください。

## キーワードとタイプ名

Wave 文法で使用する主なキーワードは次のとおりです。

`pub`, `fun`, `extern`, `export`, `type`, `enum`, `static`, `var`, `deref`, `const`, `if`, `else`, `proto`, `struct`, `while`, `for`, `in`, `out`, `clobber`, `as`, `asm`, `import`, `return`, `continue`, `print`, `input`, `println`, `match`, `break`, `true`, `false`, `null`.

内蔵タイプ名には`bool`、`char`、`byte`、`str`、整数・浮動小数点タイプ、`ptr`と`array`ポインタは`ptr<T>`、固定長配列は`array<T, N>`の形で少なくなります。

## 文字列と文字 escape

|表記|意味|
| --- | --- |
| `\n` |LF 改行|
| `\r` | CR |
| `\t` |タブ|
| `\\` |逆スラッシュ|
| `\"` |二重引用符|
| `\xNN` |正確に2桁の16進数で指定された1バイト|

通常の文字列文字はUTF-8として保存されます。 `\xNN`は1バイトを保持するため、文字列全体が有効なUTF-8であるという保証はありません。文字列リテラルの内部NUL（`\x00`を含む）はコンパイルエラーです。 0 を含むデータにはバイト配列と長さを使用します。

`char`リテラルは8ビット値にする必要があります。 `'한'`のように、その範囲を超える文字はエラーです。文字列`"한"`とは異なります。

ソースのLF、CRLF、単独CRは、それぞれ1つの論理改行として扱います。これはソースの場所とコメントの終了に関する規則であり、ファイルデータの実際のバイトを変更するという意味ではありません。

追加の文法名として`variant`、`async`、`await`があり、非同期値は`Future<T>`で表されます。上記の独立した`var`宣言ブロックは、関数の内部コードの断片です。

[文字列クラス](/docs/ja/language/arrays) · [コメント](/docs/ja/language/comments)

## 意図的に失敗する例

以下のプログラムをcheckにすると、内部NULエラーが発生するはずです。 0バイトが必要な場合は、`[97, 0, 98]`バイト配列を使用してください。

<!-- wave-example: reject-string-nul -->
```wave
fun main() {
    var text: str = "a\x00b";
}
```

以下の文字リテラルも8ビット範囲を超えるため、コンパイルエラーです。 UTF-8文字列を表すには、`str`と二重引用符を使用します。

<!-- wave-example: reject-wide-char -->
```wave
fun main() {
    var letter: char = '한';
}
```
