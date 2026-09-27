---
translation_set_id: data-types
path: language/structures-enums-and-aliases
locale: ja
group: language
group_order: 2
order: 8
title: 8. 構造体、列挙型、およびバリアント
summary: フィールド、構造体の初期化、enumとvariantの役割を学びます。
---

## データ間の関係を型で表現する

商品の価格と数をそれぞれ変数として渡すと、2つの値が同じ商品に属しているかどうかコードだけを見てわかりにくいです。構造体は関連するフィールドを囲みます。 enumは名前付きの状態を表し、variantは毎回別のデータをまとめて保存します。

3つの機能は互いに置き換える文法ではありません。何を表現したいかによって選択します。

|表現するもの|選択|はい|
| --- | --- | --- |
|同時に存在する複数のフィールド| struct |商品の単価と数量|
|名前のある整数の状態| enum |待機・進行・完了|
|場合によっては異なるデータ| variant |成功値またはエラー|

## 構造体宣言と値の生成

<!-- wave-example: book-struct-first -->
```wave
struct Product {
    price: i32;
    quantity: i32;
}

fun main() {
    var item: Product = Product {
        price: 1500,
        quantity: 2
    };

    println("{} {}", item.price, item.quantity);
}
```

実行結果：

```text
1500 2
```

宣言のフィールドはセミコロンで終わり、値を作成するとフィールドと値はコロンで連結され、コンマで区切られます。タイプを定義したものと実際の値を作成することは異なるステップです。 Productというタイプを宣言したと商品一つの収納スペースが自動的にできません。

フィールドは`item.price`のように近づきます。同じタイプのitemを複数作成すると、それぞれ異なる値を保存できます。

## 構造体を関数に渡す

<!-- wave-example: book-struct-function -->
```wave
struct Product {
    price: i32;
    quantity: i32;
}

fun subtotal(item: Product) -> i32 {
    return item.price * item.quantity;
}

fun main() {
    var item: Product = Product {
        price: 1500,
        quantity: 2
    };

    println("{}", subtotal(item));

    item.quantity = 3;
    println("{}", subtotal(item));
}
```

実行結果：

```text
3000
4500
```

関数は単価と数量が同じ商品に属するという関係をタイプとして受け取ります。値として渡された整数フィールド構造体を読み取る関数です。呼び出し元の記憶域を変更したい関数であれば、ポインタを受け取るように設計できます。

構造体にポインタフィールドがある場合、値のコピーはアドレスもコピーします。別の割り当てまで深くコピーする機能ではありません。ファイルハンドルやBufferなどのリソースを含むタイプは、コピールールとリリースルールを一緒に定義する必要があります。

## メソッドとproto

関連関数をメソッドの形で囲むことができます。 protoは、構造体のメソッドを別々のブロックとして作成する方法です。

<!-- wave-example: book-struct-method -->
```wave
struct Counter {
    value: i32;
}

proto Counter {
    fun current(self: Counter) -> i32 {
        return self.value;
    }
}

fun main() {
    var counter: Counter = Counter {
        value: 3
    };

    println("{}", counter.current());
}
```

実行結果：

```text
3
```

`self: Counter`は値を受け取るパラメータです。メソッド呼び出し表記を書くと自動的に元を変更するメソッドになるわけではありません。 selfの種類と本文で行う作業を一緒にお読みください。

メソッドを付けたため、フィールドが有効な状態のみを持つとは期待できません。パブリックフィールドでユーザーが作成できる不適切な組み合わせがある場合は、関数でチェックするか、生成ルールを提供する必要があります。

## enumで状態に名前を付ける

<!-- wave-example: book-enum-state -->
```wave
enum State -> i32 {
    Ready = 0,
    Running,
    Finished,
}

fun main() {
    var state: State = State::Ready;

    if (state == State::Ready) {
        println("ready");
    }

    state = State::Running;

    if (state == State::Running) {
        println("running");
    }
}
```

実行結果：

```text
ready
running
```

`-> i32`は表現に使用する整数型です。最初の値を0に設定し、後の省略値は前の値より1大きいです。コードで単に0と1を比較するよりも、State::ReadyとState::Runningを書くと意味が明らかになります。

enum 名前があると状態遷移が自動的に制限されるわけではありません。 FinishedからRunningに戻ることができるかどうかなどの規則は、関数として実装する必要があります。

## variantでケースとデータを接続する

成功した場合にのみ値があり、失敗したときにエラー情報が必要な場合は、variantで表すことができます。

<!-- wave-example: book-variant-result -->
```wave
variant Result {
    Value(i32),
    Error(i32)
}

fun divide(left: i32, right: i32) -> Result {
    if (left < 0 || right <= 0) {
        return Result::Error(1);
    }

    return Result::Value(left / right);
}

fun show(result: Result) {
    match (result) {
        Result::Value(value) => {
            println("value={}", value);
        }
        Result::Error(code) => {
            println("error={}", code);
        }
    }
}

fun main() {
    show(divide(12, 3));
    show(divide(12, 0));
}
```

実行結果：

```text
value=4
error=1
```

Result::ValueとResult::Errorはそれぞれpayloadを含む場合です。同じ整数型を入れてもどちらの場合か区別されます。発信者はmatchでケースを確認し、対応するarmの中でpayloadを使用します。

上記divideは負のオペランドをサポートしない小さな例です。入力範囲を指定したので、一般的なsigned除算のすべての境界をカバーする関数と混同しないでください。

## payloadがない場合

すべての場合にデータが必要なわけではありません。値のない状態を別々の場合として表現できます。

<!-- wave-example: book-variant-empty -->
```wave
variant Lookup {
    Found(i32),
    Missing
}

fun main() {
    var result: Lookup = Lookup::Missing;

    match (result) {
        Lookup::Found(value) => {
            println("found={}", value);
        }
        Lookup::Missing => {
            println("missing");
        }
    }
}
```

実行結果：

```text
missing
```

任意の整数1つを「なし」として予約するのではなく、Missingという場合を使用しました。成功値がどんな整数であっても意味が重ならない。

matchから`_`は残りのケースを処理します。新しいケースを追加したときに各呼び出し元を再確認したい場合は、すべてのケースを明示的に分割する方がよいでしょう。どの方法を選択しても、未処理の入力がないことを確認してください。

## 構造体とvariantを一緒に書く

異なるデータを選択することはvariant、ある場合に属する複数のフィールドを囲むことは構造体で表現できます。たとえば、注文処理の結果が成功した場合はレシート構造、失敗した場合はエラー番号を含むように設計できます。

値の中に他の値が含まれていても、メモリ寿命ルールは消えません。 variantにポインタを入れた場合、そのポインタが有効か、誰が解除するかは別途定めます。外部ファイル形式で保存する場合でも、構造体メモリをそのままダンプせず、フィールドごとのエンコーディングを設定する必要があります。

## 演習と解答例

0〜100のスコアのみを許可する判定結果を作成してください。有効なスコアはGrade（score）、残りはInvalidで表します。

<!-- wave-example: book-model-solution -->
```wave
variant CheckedScore {
    Grade(i32),
    Invalid
}

fun check_score(score: i32) -> CheckedScore {
    if (score < 0 || score > 100) {
        return CheckedScore::Invalid;
    }

    return CheckedScore::Grade(score);
}

fun main() {
    var result: CheckedScore = check_score(87);

    match (result) {
        CheckedScore::Grade(score) => {
            println("accepted={}", score);
        }
        CheckedScore::Invalid => {
            println("invalid score");
        }
    }
}
```

実行結果：

```text
accepted=87
```

入力を-1、0、100、101に変更して境界を確認します。成功と失敗は同じ整数空間を共有しないため、呼び出し側で誤差値を誤って平均計算に追加するリスクを減らすことができます。


## どの表現を選ぶか

|データの形状|適切な表現|はい|
| --- | --- | --- |
|同じタイプの値を複数|配列|スコア10個|
|互いに関連する複数のフィールド|構造体|名前とスコア|
|名前付きステータス値| enum | Ready, Running, Stopped |
|状態ごとに異なる追加データ| variant | Value(i32), Error(str) |
|既存の型への文脈上の名前|タイプエイリアス| UserId = u64 |

データ構造を選択するときは、保存する値だけでなく、どのような誤った状態を表現できるかを考えます。成功可否と値・エラーフィールドの両方を持つ構造体は間違った組み合わせを作ることができますが、variantは場合別payloadに区分できます。

## メモリ配置と外部データ

構造体のメモリには、フィールド間の位置合わせを調整するための空きスペースが含まれています。フィールドサイズを単に加算した値は、構造体全体のサイズと必ずしも同じではありません。サイズとソートを知る必要がある場合は、[mem レイアウト関数](/docs/ja/reference/memory-and-buffer)を使用してください。

ファイルまたはネットワークメッセージは、[bytes](/docs/ja/stdlib/bytes)関数でフィールドを順番にエンコードすることで、バイトの順序と長さを明確に決定できます。他の言語と構造体を渡すときは、[FFI](/docs/ja/language/modules-imports-and-ffi)の外部宣言とターゲットABIを合わせます。

[構造体の学習と練習](/docs/ja/language/structures-enums-and-aliases) · [variant](/docs/ja/language/variants)
