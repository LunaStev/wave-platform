---
translation_set_id: functions
path: language/functions-and-generics
locale: ja
group: language
group_order: 2
order: 5
title: 5. 機能の設計と構成
summary: パラメータ、戻り値、デフォルト値、および値渡しを学びます。
---

## 繰り返しコードから出発する

関数は文法を減らすためのツールでもありますが、作業の境界を定めるツールです。入力として何を受け取り、何を計算し、どの結果を返すかを分離すると、プログラムを小さな単位として理解できます。

この章では、割引計算を複数回作成するプログラムから出発します。各例はmain.wave全体で、`wavec run main.wave`で実行します。まだファイルを分割していません。

<!-- wave-example: book-function-before -->
```wave
fun main() {
    var first_price: i32 = 2000;
    var first_discount: i32 = first_price * 10 / 100;
    var first_total: i32 = first_price - first_discount;
    var second_price: i32 = 5000;
    var second_discount: i32 = second_price * 10 / 100;
    var second_total: i32 = second_price - second_discount;
    println("{} {}", first_total, second_total);
}
```

実行結果：

```text
1800 4500
```

両方の計算は価格のみが異なり、構造が同じです。割引ルールを変更するときは、両方の場所を変更する必要があります。一方だけを変えると、同じポリシーを適用する必要がある商品で異なる結果が得られます。

## 入力と出力を決める

重複した計算を関数に移動します。変わる値はpriceという入力で受け取り、計算した価格は返します。

<!-- wave-example: book-function-extract -->
```wave
fun discounted(price: i32) -> i32 {
    var discount: i32 = price * 10 / 100;
    return price - discount;
}

fun main() {
    println("{} {}", discounted(2000), discounted(5000));
}
```

実行結果：

```text
1800 4500
```

関数名はdiscountedです。括弧内の`price: i32`はパラメータ宣言で、`-> i32`は結果型です。本文のローカル変数discountは、この関数内でのみ使用されます。

`discounted(2000)`は関数を呼び出す式です。括弧内2000は実際に渡す引数です。関数が返した値がこの呼び出し式の結果になるため、printlnの引数として直接使用できます。

|用語|コード|意味|
| --- | --- | --- |
|パラメータ| price |関数を宣言するときに指定した入力名|
|引数| 2000 |呼び出し時に渡した値|
|戻りタイプ| i32 |呼び出し式が生成する値の型|
|戻り文| return price - discount |結果を渡して今回の呼び出しを終了|

## 呼び出しと実行順序に従う

関数宣言をソースに書き込むだけで、本文はすぐには実行されません。 mainから呼び出されたポイントに達したときに実行します。

<!-- wave-example: book-function-trace -->
```wave
fun calculate(value: i32) -> i32 {
    println("inside: {}", value);
    return value * 2;
}

fun main() {
    println("before");
    var result: i32 = calculate(7);
    println("after: {}", result);
}
```

実行結果：

```text
before
inside: 7
after: 14
```

進行順序は、mainの最初の出力、calculate本文、mainの最後の出力です。 `return`が実行されると、calculate呼び出しが結果14で終了し、mainのresult初期化が完了します。

複数の関数呼び出しが1つの式に混在していて副作用の順序が重要な場合は、呼び出しを別の文に分割してください。この章の例も、トレースを必要とする呼び出しの結果をローカル変数に保存します。

## 複数のパラメータ

割引率も入力として受け取ると、同じ関数で複数のポリシーを計算できます。

<!-- wave-example: book-function-parameters -->
```wave
fun discounted(price: i32, percent: i32) -> i32 {
    return price - price * percent / 100;
}

fun main() {
    var standard: i32 = discounted(2000, 10);
    var special: i32 = discounted(2000, 25);
    println("standard={}", standard);
    println("special={}", special);
}
```

実行結果：

```text
standard=1800
special=1500
```

引数の順序が宣言と一致する必要があります。両方のパラメータがi32の場合、順序を変更しても型チェックだけで意味を区別することは困難です。関数名とパラメータ名を明確に定義して呼び出す場所も読みやすくします。

この関数は、小さな金額と有効な割引率を前提としています。負の価格、100より大きい割合、中間乗算オーバーフローを処理しません。関数を作成するときは、本文だけでなく入力条件も説明する必要があります。後の完成プログラムは検査ステップを分離します。

## 基本因子

よく使用する値をデフォルト値として指定できます。

<!-- wave-example: book-function-default -->
```wave
fun discounted(price: i32, percent: i32 = 10) -> i32 {
    return price - price * percent / 100;
}

fun main() {
    println("{}", discounted(2000));
    println("{}", discounted(2000, 25));
}
```

実行結果：

```text
1800
1500
```

最初の呼び出しは2番目の引数を省略して10を使用します。 2 番目の呼び出しは、明示した 25 を使用します。デフォルト値は、後ろの省略可能なパラメータに置きます。最初の引数だけを欠かすために空白を書く文法として使用しません。

デフォルト値を変更すると、省略呼び出しの動作が変わります。公開関数のデフォルト値も、ユーザーが依存する動作の一部です。引数を指定した呼び出しと省略した呼び出しをそれぞれテストする理由です。

## 値として渡すという意味

整数値を渡すと、関数が受け取る値と呼び出し元の変数記憶領域が区別されます。結果を計算したとしても、呼び出し側の変数は自動的に変更されません。

<!-- wave-example: book-function-value -->
```wave
fun next(value: i32) -> i32 {
    return value + 1;
}

fun main() {
    var count: i32 = 4;
    var later: i32 = next(count);
    println("count={} later={}", count, later);
    count = next(count);
    println("count={}", count);
}
```

実行結果：

```text
count=4 later=5
count=5
```

最初の呼び出しはcountを読み、laterを初期化します。 countはまだ4です。 2回目の呼び出しの後に結果をcountに代入したので、5になります。値を返す設計は、データがどこで変更されるかを呼び出し側で明らかにします。

元の記憶領域を関数内で置き換えたい場合は、ポインタを渡すことができます。これは[ポインター章](/docs/ja/language/explicit-memory-type-model)でカバーされています。ポインタを渡しても、ポインタ値自体とそのアドレスの記憶領域は区別する必要があります。

## 返さないパスを残さない

値を返す関数は、必要なすべてのパスで結果を提供する必要があります。次のように最後のパスを残しません。

```wave
// 의도적으로 잘못된 함수 조각

fun positive(value: i32) -> i32 {
    if (value > 0) {
        return value;
    }
}
```

valueが0以下の場合、どの値を返すかは決まっていません。意図したルールを設定し、すべてのパスを作成する必要があります。

<!-- wave-example: book-function-paths -->
```wave
fun positive_or_zero(value: i32) -> i32 {
    if (value > 0) {
        return value;
    }

    return 0;
}

fun main() {
    println("{} {} {}", positive_or_zero(-2), positive_or_zero(0), positive_or_zero(8));
}
```

実行結果：

```text
0 0 8
```

最初のreturnを実行した呼び出しでは、以下のreturnに進みません。条件が偽のときだけ最後のreturnに達します。正、0、負の3つのケースでルールを確認しました。

## 結果のない関数

出力などの操作のみを実行する場合は、戻り型を省略できます。結果のない関数も`return;`で早く終了できます。

<!-- wave-example: book-function-void -->
```wave
fun show_positive(value: i32) {
    if (value <= 0) {
        return;
    }

    println("positive={}", value);
}

fun main() {
    show_positive(-3);
    show_positive(6);
}
```

実行結果：

```text
positive=6
```

最初の呼び出しは何も出力せずに戻ります。 2番目の呼び出しは出力します。 「結果値がない」と「呼び出し元に戻らない」は異なります。プロセスを終了する関数のように戻らない関数は、戻り型`!`で区切ります。

## 複数の関数でコンプリートプログラムを構成する

これで、入力検証、計算、出力をさまざまな関数に分割します。価格範囲が制限されているため、この例の中間乗算はi32の範囲内です。

<!-- wave-example: book-function-program -->
```wave
fun valid_order(price: i32, quantity: i32, percent: i32) -> bool {
    return price >= 0 && price <= 100000
        && quantity >= 1 && quantity <= 100
        && percent >= 0 && percent <= 100;
}

fun discounted_unit(price: i32, percent: i32) -> i32 {
    return price - price * percent / 100;
}

fun order_total(price: i32, quantity: i32, percent: i32) -> i32 {
    var unit: i32 = discounted_unit(price, percent);
    return unit * quantity;
}

fun show_order(price: i32, quantity: i32, percent: i32) {
    if (!valid_order(price, quantity, percent)) {
        println("invalid order");
        return;
    }

    var total: i32 = order_total(price, quantity, percent);
    println("total={}", total);
}

fun main() {
    show_order(2000, 3, 10);
    show_order(2000, 0, 10);
    show_order(2000, 3, 120);
}
```

実行結果：

```text
total=5400
invalid order
invalid order
```

関数ごとに1つの質問に答えます。 valid_orderは入力が許可されるか、discounted_unitは1つの割引価格がいくらか、order_totalは全体価格がいくらか、show_orderは何を見せるかを担当します。

小さな関数が無条件に良いわけではありません。一式ごとに名前を付けると、むしろ追跡が難しくなることがあります。他の場所で再使用する意味があるか、独立して説明・検証する規則があるときに分離します。

## 練習問題

1. 2 つの整数のうち大きい方を返す maximum を作成します。
2. 2 つの境界の間の値を返す clamp 関数を作成します。以下のプールでは、low<=highを呼び出し条件として設定します。
3. 値に税を加算する関数を作成し、注文合計関数と組み合わせます。範囲と整数切削の時点を最初に決めます。

### プール：境界を持つ関数

<!-- wave-example: book-function-solution -->
```wave
fun maximum(left: i32, right: i32) -> i32 {
    if (left > right) {
        return left;
    }

    return right;
}

fun clamp(value: i32, low: i32, high: i32) -> i32 {
    if (value < low) {
        return low;
    }
    if (value > high) {
        return high;
    }

    return value;
}

fun main() {
    println("max={}", maximum(7, 4));
    println("{} {} {}", clamp(-3, 0, 10), clamp(6, 0, 10), clamp(20, 0, 10));
}
```

実行結果：

```text
max=7
0 6 10
```

clampは、範囲より小さい、範囲内、範囲より大きい3つのパスがあります。境界値0と10も直接追加して確認してください。 low> highの場合までサポートするには、失敗をどのように表現するかを決める必要があります。 [エラー処理の章](/docs/ja/language/errors)で結果構造体を使用してこの問題を扱います。


## 再帰呼び出し

関数は自分自身を呼び出すことができます。再帰では、もはや呼び出さない終了条件と、各呼び出しがその条件に近づく過程が必要です。

<!-- wave-example: book-ref-recursion -->
```wave
fun factorial(value: i32) -> i32 {
    if (value <= 1) {
        return 1;
    }

    return value * factorial(value - 1);
}

fun main() {
    println("{}", factorial(5));
}
```

実行結果：

```text
120
```

5の計算は`5 * factorial(4)`に続き、1に達すると1を返します。戻り値が前の呼び出しに順番に渡され、120になります。この関数は、小さな正の整数を説明するための例です。大きな入力では、結果の範囲と呼び出しの深さを考慮する必要があります。繰り返しステートメントで同じタスクを作成すると、呼び出しの深さが増えるという問題を回避できます。

## 頻繁に発生するエラー

|現象|確認する|
| --- | --- |
|引数が足りないか多いというエラー|パラメータ数と省略可能なデフォルト値|
|戻り型が合わない|return 式の型と関数宣言|
|特定のパスから返されない|条件が偽になるまで返すか|
|ジェネリック型引数がありません|関数名の後 `<Type>`|
|ポインタ関数の呼び出しの後にソースが変わります|値を読み取るだけの関数か変更する関数か|

`export(c)`のように外部にエクスポートする関数は、具体的な署名を使用します。ジェネリック関数自体を外部呼び出し規約にエクスポートすることはできません。 `ptr<T>`と`array<T, N>`は言語の組み込みメモリタイプであり、ユーザージェネリック構造体宣言と区別します。

[関数学習](/docs/ja/language/functions-and-generics) · [モジュールとジェネリック学習](/docs/ja/language/modules-imports-and-ffi) · [FFI](/docs/ja/language/modules-imports-and-ffi)
