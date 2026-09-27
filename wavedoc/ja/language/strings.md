---
translation_set_id: learn-strings
path: language/strings
locale: ja
group: language
group_order: 2
order: 7
title: 7. 文字列、文字、バイト
summary: 文字列とchar、UTF-8バイト長、NUL、検索とバイナリデータを区別します。
---

## 画面の文字とメモリのバイト

画面には文字が表示されますが、メモリにはバイトが格納されます。特に、韓国語のように1文字が複数のUTF-8バイトの場合、「長さ」と「数番目の文字」を同じ意味で使用すると間違いが発生しやすいです。

この章では、strとcharの違い、NUL終了、escape、検索位置とバイナリデータを区別します。例はそれぞれプログラム全体です。

## 文字列リテラルと出力

<!-- wave-example: book-string-literals -->
```wave
fun main() {
    var greeting: str = "안녕하세요";

    println("{}", greeting);
    println("line one\nline two");
    println("quote: \"Wave\"");
}
```

実行結果：

```text
안녕하세요
line one
line two
quote: "Wave"
```

二重引用符の中の一般的な文字はUTF-8で表されます。 escapeは、ソースから直接書き込むのが難しいバイトを示します。 `\n`は改行1バイトで、逆スラッシュとnの2文字を出力するものではありません。逆スラッシュ自体を出力するには、`\\`を使用します。

## 長さはバイト数

<!-- wave-example: book-utf8-length -->
```wave
import("std::string::len")::{
    len
};

fun main() {
    println("ASCII={}", len("Wave"));
    println("Korean={}", len("한"));
    println("mixed={}", len("Wave한"));
}
```

実行結果：

```text
ASCII=4
Korean=3
mixed=7
```

`len` は、目に見えない文字である終端の NUL の前のバイトをカウントします。文字 `한` は、UTF-8 では 3 バイトを必要とします。文字数、Unicode コード ポイント数、およびバイト数は、一般に互換性がありません。表示幅はフォントや文字の組み合わせなどによっても異なります。

したがって、任意のバイト位置で文字列を切り取り、画面に出力する機能は、Unicode境界を別々に考慮する必要があります。 ASCIIのみ処理するプログラムなのか一般Unicodeテキストなのか、入力条件を明確に決めてください。

## NUL 終了と長さ

strは終了を示す0バイトを使用します。 lenはその最後のバイトを長さに含めません。文字列リテラル内にNULを入れるのはエラーです。以下は意図的なエラーの例です。

```wave
fun main() {
    var text: str = "left\x00right";
}
```

ソースの`\xNN`は、正確に2桁の16進数で1バイトを指定します。 `\x41`は65バイトAを表します。一般文字をUTF-8で書くのとランダムな1バイトを入れることは異なりますので、`\xNN`を書くことができるすべてのstrが有効なUTF-8であるわけではありません。

## charはUnicode文字全体を含まない

charは符号なし8ビット文字値です。 `'A'`のように、1バイト範囲の値を表すリテラルを使用できます。 `'한'`はこの範囲に入らないのでエラーです。 `"한"`は複数のUTF-8バイトを持つstrとは別個です。

1文字を常にchar1に入れようとしないでください。テキストを処理するために必要な単位がバイトであるかUnicodeコードポイントであるかを最初に決定する必要があります。

## 文字列比較

文字列内容の比較には、stdの関数を使用します。以下は、同じ内容と大文字と小文字の違いを調べるプログラムです。

<!-- wave-example: book-string-compare -->
```wave
import("std::string::cmp")::{
    eq, starts_with, ends_with
};

fun main() {
    if (eq("Wave", "Wave")) {
        println("same bytes");
    }

    if (!eq("Wave", "wave")) {
        println("case differs");
    }

    if (starts_with("report.txt", "report") && ends_with("report.txt", ".txt")) {
        println("name matches");
    }
}
```

実行結果：

```text
same bytes
case differs
name matches
```

この比較はバイト列を比較します。言語別の大文字と小文字の変換や、Unicodeの正規化は自動的に行われません。ファイル名を比較するときも、OSのファイル名同等性ルールと単純な文字列比較が同じではありません。

## 検索結果の単位と失敗

<!-- wave-example: book-string-search -->
```wave
import("std::string::find")::{
    find, count
};

fun main() {
    var first: i32 = find("banana", "na");
    var missing: i32 = find("banana", "xy");
    var matches: i32 = count("aaaa", "aa");

    println("first={} missing={}", first, missing);
    println("matches={}", matches);
}
```

実行結果：

```text
first=2 missing=-1
matches=2
```

findは最初の位置または-1を返します。インデックス0も成功なので、`result >= 0`でチェックします。 countは位置ではなく、重ならない一致数です。上でaaは0～1と2～3に一致して2回です。

空のneedleも契約の一部です。 findは0、containsはtrue、countは0を返します。関数名が同じモジュールにあると返す方法まで同じだとは思わないでください。

## 空白の削除は新しい文字列の生成とは異なります

trim_rangeは、原文を変更またはコピーせずに空白以外の範囲を返します。出力ポインタを受け取るので、結果を保存する整数を最初に準備します。

<!-- wave-example: book-string-trim-range -->
```wave
import("std::string::trim")::{
    trim_range
};

fun main() {
    var start: i32 = 0;
    var end: i32 = 0;

    trim_range("  Wave  ", &start, &end);

    println("start={} end={} bytes={}", start, end, end - start);
}
```

実行結果：

```text
start=2 end=6 bytes=4
```

範囲は`[start, end)`です。開始は含まれ、終了は含まれないので、長さはend-startです。原文の先頭アドレスにstartを加えると自動的にend位置にNULが生じるわけではありません。範囲を別々に持ち歩くか、新しい文字列スペースを準備する必要があります。

## バイナリデータには別の長さ

0 を含むデータは、文字列終了ルールとして扱われません。バイト配列と長さを使用します。

<!-- wave-example: book-binary-not-string -->
```wave
fun main() {
    var bytes: array<u8, 3> = [65, 0, 66];

    for (var index: i32 = 0; index < 3; index += 1) {
        println("{}", bytes[index]);
    }
}
```

実行結果：

```text
65
0
66
```

2番目のゼロは実際のデータです。これをstrと解釈すると最初の0で終わったものとして扱い、後の66を見ることができません。逆にNULのない配列をstrにcastすると、配列の外まで読む危険があります。 castは終了バイトを追加する操作ではありません。

## 演習：ファイル名の確認

ファイル名が`.wave`で終わっているか確認し、文字列に`test`が含まれている場合はテストファイルに出力してください。この演習は名前のバイトパターンのみをチェックし、実際のファイルの存在については説明しません。

### 解答例

<!-- wave-example: book-string-solution -->
```wave
import("std::string::cmp")::{
    ends_with
};
import("std::string::find")::{
    contains
};

fun classify(name: str) {
    if (!ends_with(name, ".wave")) {
        println("other file");
        return;
    }

    if (contains(name, "test")) {
        println("Wave test file");
    } else {
        println("Wave source file");
    }
}

fun main() {
    classify("main.wave");
    classify("parser_test.wave");
    classify("notes.txt");
}
```

実行結果：

```text
Wave source file
Wave test file
other file
```

大文字`.WAVE`はどのように処理するのか、パス全体にtestが入っていてもテストで見るかなどは別途ポリシーです。小さい関数でも、どの入力を対象とするかを決めなければ、動作を正確に説明できます。
