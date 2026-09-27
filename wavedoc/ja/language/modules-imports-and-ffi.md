---
translation_set_id: modules-ffi
path: language/modules-imports-and-ffi
locale: ja
group: language
group_order: 2
order: 11
title: 11. モジュールと汎用コード
summary: 公開名の取得と明示的な型引数を学びます。
---

## ファイルを分割する理由

プログラムが大きくなると、すべての関数をmain.waveに置くよりも関連機能同士を結ぶほうが見つけやすくなります。モジュール境界は、どの名前を別のコードに公開するかを決定します。ジェネリックは、ファイルの分離とは別にタイプが異なる同じ操作を再利用するツールです。

この章では、2つのファイルプログラムを作成し、モジュールエイリアス、選択import、ジェネリック関数と構造体を学びます。

## 2つのファイルプログラム

同じディレクトリにhelpers.waveとmain.waveを作成します。

helpers.wave 全体:

```wave
pub fun double(value: i32) -> i32 {
    return value * 2;
}
```

main.wave 全体:

<!-- wave-example: book-module-two-files -->
```wave
import("./helpers")::{
    double
};

fun main() {
    var result: i32 = double(21);

    println("{}", result);
}
```

実行結果：

```text
42
```

端末で`wavec run main.wave`で実行します。 helpers.waveも別々に実行するものではありません。 importを介して必要なソースが接続されます。

helpersの関数の前に付いているpubは、他のモジュールがインポートできることを示します。外部に公開する必要がない補助関数は公開する必要はありません。モジュールの内部実装を変更しても、公開関数の契約を守ることで、使用するコードの変更を減らすことができます。

## 相対パスの基準

`./helpers`は、import文を作成したソースファイルのディレクトリに基づいています。プログラム実行時にファイルI/Oが書き込む作業ディレクトリと区切ります。 importファイルの検索手順と実行中のinput.txtの検索手順は異なります。

ローカルimportでは、`.wave`拡張子を省略できます。ディレクトリを分割したら、場所に合わせて`./module`パスを作成します。パッケージ依存関係の名前を取得するパスとローカル相対パスを混同しないでください。

## 選択importとエイリアス

選択 importは、必要な公開名のみを現在のファイルで直接使用します。名前の競合がある場合、またはどのモジュールの関数かを明らかにしたい場合は、エイリアスを使用してください。

<!-- wave-example: book-module-alias -->
```wave
import("std::string::len" as strings);

fun main() {
    var length: i32 = strings::len("Wave");

    println("{}", length);
}
```

実行結果：

```text
4
```

stringsは、このファイルで指定されたモジュールエイリアスです。 `strings::len`はそのモジュールの名前を使用します。フィールドアクセスのポイントとモジュールの区別の`::`は異なる表記です。

選択importとエイリアスimportを1つの文にまとめて書きません。どんなスタイルでも、ファイル全体で名前のソースを読みやすくするために一貫して使用します。

## 標準ライブラリとパッケージ

`std::`パスは標準ライブラリを指します。ユーザーは、必要なモジュールの公開APIをimportします。標準ライブラリ関数がすべて自動的に現在の名前空間に入るわけではありません。

外部パッケージパスはパッケージ名から始まります。パッケージの場所は、コンパイラオプションまたはパッケージマネージャによって提供されます。まず、ローカルモジュールに境界をつけた後、[Vex 使い方](/docs/ja/whale/vex-package-manager)で依存関係を管理する方法を学びます。

## 型別に同じ関数を書く

次の関数は入力をそのまま返します。 i32とstrのために同じコードを2回書き込まないように、タイプパラメータTを使用してください。

<!-- wave-example: book-generic-identity -->
```wave
fun identity<T>(value: T) -> T {
    return value;
}

fun main() {
    var number: i32 = identity<i32>(42);
    var text: str = identity<str>("Wave");

    println("{} {}", number, text);
}
```

実行結果：

```text
42 Wave
```

Tは、実行中に渡す整数値ではなく、型を入れる桁です。 `<i32>`のように型引数を明示して呼び出します。一般ユーザージェネリック関数では型引数を省略しません。

identity<str>は文字列バイトを新しく割り当てて複製しません。値をそのまま返します。ジェネリックという文法がデータのコピー・所有権規則を変えません。

## ジェネリック本文が要求する演算

タイプパラメーターがあると、すべての操作がすべてのタイプに使用できるわけではありません。以下のminimumは、比較可能な実際のタイプとして使用する必要があります。

<!-- wave-example: book-generic-minimum -->
```wave
fun minimum<T>(left: T, right: T) -> T {
    if (left < right) {
        return left;
    }

    return right;
}

fun main() {
    var small: i32 = minimum<i32>(7, 4);
    var wide: i64 = minimum<i64>(100, 20);

    println("{} {}", small, wide);
}
```

実行結果：

```text
4 20
```

型引数を変更すると、本文で使用する`<`と戻りがその型で成立する必要があります。ジェネリックエラーを読み取るときは、呼び出した型の組み合わせと関数本体が要求する操作を一緒に確認します。

## ジェネリック構造体

2つの異なる値を囲むPairを作成します。

<!-- wave-example: book-generic-pair -->
```wave
struct Pair<A, B> {
    first: A;
    second: B;
}

fun main() {
    var item: Pair<i32, str> = Pair<i32, str> {
        first: 7,
        second: "seven"
    };

    println("{} {}", item.first, item.second);
}
```

実行結果：

```text
7 seven
```

Pair<i32, str>とPair<i64, str>は異なる具体的なタイプです。型引数の順序も意味があります。 firstとsecondのタイプがどこで決定されるか宣言と生成コードを連結して読んでください。

## 公開APIの名前と契約

関数を公開するときは、名前だけでなく、入力単位、戻り値、失敗と所有権を一緒に決めます。たとえば、readが最大長を読み取るか正確な長さを読み取るかに応じて、呼び出し元が作成する反復文が異なります。

pubはWaveモジュール間の公開範囲です。他の言語が呼び出す外部シンボルをエクスポートするexport（c）とは異なる機能です。 [FFI参照](/docs/ja/language/modules-imports-and-ffi)では、2つの言語を結ぶ完成例を見ることができます。

## 演習と解答例

math.waveに公開関数squareを作成し、main.waveからエイリアスとして呼び出し、3と5の2乗を出力します。

math.wave:

```wave
pub fun square(value: i32) -> i32 {
    return value * value;
}
```

main.wave:

<!-- wave-example: book-module-solution -->
```wave
import("./math" as math);

fun main() {
    println("{} {}", math::square(3), math::square(5));
}
```

実行結果：

```text
9 25
```

関数が消えたというエラーがある場合は、importパスとpubを最初に確認してください。名前の競合の場合は、別名がある呼び出しであることを確認してください。型エラーの場合は、関数の入力と渡された引数型を確認してください。さまざまな問題を1つのパス修正で解決しようとしないでください。


## C関数のインポート

```wave
extern(c) fun puts(text: ptr<i8>) -> i32;
```

ABI名前の後には、実際のシンボル名を文字列として指定できます。

```wave
extern(c, "native_symbol") fun local_name(value: i32) -> i32;
```

## Wave関数のエクスポート

```wave
export(c) fun wave_add(left: i32, right: i32) -> i32 {
    return left + right;
}
```

`extern`と`export`は、単一の関数とブロック形式で使用できます。エクスポートする関数には、具体的なABIシグネチャが必要なので、ジェネリックにすることはできません。

## ターゲット条件属性

最上位項目にはターゲット条件属性を付けることができます。

```wave
#[target(os="linux", arch="x86_64")]
extern(c) fun platform_call(value: i32) -> i32;
```

条件キーは`arch`、`os`、`env`、`abi`で、属性はすぐ次の最上位項目に適用されます。

## 自分で書いたC関数にリンクする

この実習は、Cコンパイラがあるネイティブ環境用です。ライブラリの割り当てや文字列の処理なしで整数関数を連結します。

`native.c`:

```c
#include <stdint.h>
int32_t native_double(int32_t value) { return value * 2; }
```

`main.wave`:

```wave
extern(c) fun native_double(value: i32) -> i32;

fun main() {
    println("{}", native_double(21));
}
```

Linux/macOSで同じ作業ディレクトリの端末として実行します。

```shell
cc -c native.c -o native.o
wavec build main.wave native.o -o ffi-example
./ffi-example
```

期待出力は`42`です。 WindowsのMSVC開発者シェルでは、`cl /c native.c /Fonative.obj`でobjectを作成し、`wavec build main.wave native.obj -o ffi-example.exe`に接続します。ソースとobjectのターゲットアーキテクチャは同じでなければなりません。例では小さい値のみを使用します。 C関数に大きな値を渡すには、C側の乗算の範囲も別々に保証する必要があります。

ローカルファイルパスは`./`で始まります。
