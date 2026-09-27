---
translation_set_id: practice-input-calculator
path: practice/input-calculator
locale: ja
group: practice
group_order: 4
order: 1
title: プロジェクト: 入力を検証する計算機
summary: 入力、範囲チェック、関数と終了コードをリンクします。
---

## 目標と実行

数量と単価を入力して合計を計算します。 2つの整数をスペースまたは改行で区切って入力します。この例は、数量 1～1000、単価 0～100000 のみを受け取るので、 i32 乗算範囲内で計算します。

`main.wave`に保存して`wavec run main.wave`を実行し、`3 1200`と入力します。プログラムが出力する結果は以下の通りです。端末で入力した文字が見えるのは、プログラム出力とは別です。

<!-- wave-example: input-calculator -->
```wave
fun total(quantity: i32, price: i32) -> i32 {
    return quantity * price;
}

fun main() -> i32 {
    var quantity: i32 = 0;
    var price: i32 = 0;
    input("{} {}", quantity, price);
    if (quantity < 1 || quantity > 1000 || price < 0 || price > 100000) {
        println("out of range");
        return 1;
    }

    println("total={}", total(quantity, price));
    return 0;
}
```

実行結果：

```text
total=3600
```

## 失敗も確認する

`0 1200`を入れると`out of range`と終了コード1を期待します。 `input`の数値解析の失敗とプログラムの作業範囲のチェックは異なる段階です。数字ではなくトークン・タイプ範囲超過・必要な入力前 EOFは入力失敗です。この組み込み入力は、エラーを返して再入力するインターフェイスではありません。回復可能なパーサが必要な場合は、[io](/docs/ja/stdlib/files-io)でバイトを読み取って検証プロセスを直接設定します。

## 拡張練習と解説

割引率を3番目の入力として受け取り、0〜100であることを確認してください。大きな中間乗算を避けるには、i64で計算範囲を広げ、結果を絞り込むときに範囲を確認する必要があります。単純に最後の結果タイプだけを広げると、中間計算がすでに狭いタイプで行われている可能性があります。

[コンソールI/O](/docs/ja/language/console-io-and-formatting) · [次条：ファイルを読む](/docs/ja/practice/file-reader)

## 計算範囲を最初に決める理由

最大の入力は数量1000と単価100000です。 2つの値を掛けると100000000なので、i32の範囲内に入ります。この範囲チェックがなければ、total関数の乗算結果をそのまま使用できます。

入力を受け取るmainは入力とエラーメッセージを担当し、totalは計算のみを担当します。後でファイルから注文を読み取るように変更しても、計算関数はそのまま使用できます。

## 割引計算を完了する

割引率を追加すると、合計に100を掛けることができます。中間計算からi64で実行するように変換し、割引率を適用します。整数除算の小数部は破棄されるため、この例では割引後の金額を整数単位で切り捨てます。

<!-- wave-example: book-calculator-discount -->
```wave
fun discounted_total(quantity: i32, price: i32, discount: i32) -> i64 {
    var subtotal: i64 = (quantity as i64) * (price as i64);
    var remaining: i64 = 100 - (discount as i64);
    return subtotal * remaining / 100;
}

fun main() -> i32 {
    var quantity: i32 = 0;
    var price: i32 = 0;
    var discount: i32 = 0;

    input("{} {} {}", quantity, price, discount);

    if (quantity < 1 || quantity > 1000) {
        println("invalid quantity");
        return 1;
    }

    if (price < 0 || price > 100000) {
        println("invalid price");
        return 2;
    }

    if (discount < 0 || discount > 100) {
        println("invalid discount");
        return 3;
    }

    println("total={}", discounted_total(quantity, price, discount));
    return 0;
}
```

実行結果：

```text
total=3240
```

上記の結果は、`3 1200 10`を入力したときです。元の合計3600から10％を引いた3240が出力されます。

|入力|予想結果|確認するパス|
| --- | --- | --- |
| `3 1200 0` | `total=3600` |割引なし|
| `3 1200 100` | `total=0` |全額割引|
| `3 1200 101` | `invalid discount` |割引率範囲超過|
| `1000 100000 0` | `total=100000000` |最大入力|
| `0 1200 10` | `invalid quantity` |数量範囲を満たさない|

## 次の練習

関数を変更して少数の金額を四捨五入してみてください。正の金額のみを受け取るこのプログラムでは、100で割る前に50を加算すると整数単位の丸めになります。 1つに99、割引率50を入力した場合、切削結果49と丸め結果50を比較できます。
