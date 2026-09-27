---
translation_set_id: learn-errors
path: language/errors
locale: ja
group: language
group_order: 2
order: 12
title: 12. エラーの表示と回復
summary: エラーと正常値を区別し、障害パスからリソースをクリーンアップします。
---

## 失敗度関数の結果

ファイルが存在しないか、入力が範囲外であることは、プログラムで自然に発生します。エラーを処理することは、メッセージを出力するのにとどまりません。失敗を区別し、すでに進行したジョブの状態を確認し、確保したリソースを整理してから続行するか終了するかを選択するプロセスです。

今回の章では、小さな関数の失敗表現から始まり、結果構造体、variant、早期リターンとリソース整理に進化します。

## 失敗標識が成功値と重複してはいけません

配列検索で -1 を検索しなかったのは、有効なインデックスが 0 以上であるためです。一方、どの整数でも正常な結果である可能性がある計算では、-1をエラーに設定すると通常の-1と区別できません。

0もよく誤解する値です。空の文字列長0、最初の位置0、送信したバイト数0は、関数ごとに意味が異なります。戻り値がゼロでないという理由だけで成功と判断しないでください。

## 成功と値を一緒に返す

<!-- wave-example: book-error-result -->
```wave
struct Division {
    ok: bool;
    value: i32;
}

fun divide_nonnegative(left: i32, right: i32) -> Division {
    if (left < 0 || right <= 0) {
        return Division {
            ok: false,
            value: 0
        };
    }

    return Division {
        ok: true,
        value: left / right
    };
}

fun main() {
    var result: Division = divide_nonnegative(0, 3);

    if (!result.ok) {
        println("invalid input");
        return;
    }

    println("value={}", result.value);
}
```

実行結果：

```text
value=0
```

正常な結果もゼロになる可能性があります。 valueを見て、成功するかどうかを推測せずにokを最初に確認します。失敗した結果にvalueフィールドが存在しても、使用する値という意味ではありません。

この例の入力規則は、左がゼロ以上、右が正であることです。関数名と説明に範囲を示し、一般的なsigned整数除算と区別しました。

## エラーの原因も区別する

失敗の理由に応じて他のガイダンスや回復を行うには、エラー情報を追加します。入力範囲を調べる関数と計算する関数を分ける方法もあります。

<!-- wave-example: book-error-codes -->
```wave
struct Check {
    ok: bool;
    error: i32;
}

fun check_quantity(value: i32) -> Check {
    if (value < 1) {
        return Check {
            ok: false,
            error: 1
        };
    }

    if (value > 100) {
        return Check {
            ok: false,
            error: 2
        };
    }

    return Check {
        ok: true,
        error: 0
    };
}

fun main() {
    var result: Check = check_quantity(101);

    if (!result.ok) {
        match (result.error) {
            1 => {
                println("quantity is too small");
            }
            2 => {
                println("quantity is too large");
            }
            _ => {
                println("invalid quantity");
            }
        }
    }
}
```

実行結果：

```text
quantity is too large
```

数値の意味は、この関数によって定義されます。他のライブラリのエラー番号1または2と同じとは見えません。公開 APIでは、エラー定数やタイプに名前を付けると、発呼者が任意の数字を覚える必要はありません。

## variantで結果を分離する

成功値とエラーが同時に存在できないという関係をvariantで表現できます。

<!-- wave-example: book-error-variant -->
```wave
variant ByteResult {
    Value(u8),
    OutOfRange(i32)
}

fun to_byte(value: i32) -> ByteResult {
    if (value < 0 || value > 255) {
        return ByteResult::OutOfRange(value);
    }

    return ByteResult::Value(value as u8);
}

fun main() {
    var result: ByteResult = to_byte(300);

    match (result) {
        ByteResult::Value(value) => {
            println("byte={}", value);
        }
        ByteResult::OutOfRange(original) => {
            println("out of range: {}", original);
        }
    }
}
```

実行結果：

```text
out of range: 300
```

エラーに元の入力値が含まれています。単にfalseだけを返すよりも、呼び出し側が問題を説明するのは簡単です。パスワードやトークンなどの機密入力は、そのままログに残さないようにデータの性格も考慮します。

## 早期戻りで通常のパスを読みやすくする

複数の手順を実行するときは、すべての通常のコードを深いifの中に配置する必要はありません。失敗した場合はすぐに返し、以下は成功パスを続けることができます。

<!-- wave-example: book-error-early-return -->
```wave
fun show_total(quantity: i32, price: i32) {
    if (quantity < 1 || quantity > 100) {
        println("invalid quantity");
        return;
    }

    if (price < 0 || price > 100000) {
        println("invalid price");
        return;
    }

    var total: i32 = quantity * price;
    println("total={}", total);
}

fun main() {
    show_total(3, 1200);
    show_total(0, 1200);
    show_total(3, -1);
}
```

実行結果：

```text
total=3600
invalid quantity
invalid price
```

検査に合格した後は、quantityとpriceが定められた範囲内であるという事実を利用できます。範囲は中間乗算度i32に入るように決めました。早期返却を追加するときは、その時点ですでに所有しているリソースがあるかどうかを確認する必要があります。

## 障害パスのメモリクリーンアップ

<!-- wave-example: book-error-cleanup -->
```wave
import("std::buffer::alloc")::{
    Buffer, buffer_init, buffer_free
};
import("std::buffer::write")::{
    buffer_append_str
};

fun make_message(allowed: bool) -> i32 {
    var data: Buffer;

    if (buffer_init(&data, 16) < 0) {
        return 1;
    }

    if (!allowed) {
        buffer_free(&data);
        return 2;
    }

    if (buffer_append_str(&data, "ready") < 0) {
        buffer_free(&data);
        return 3;
    }

    println("bytes={}", data.len);

    if (buffer_free(&data) < 0) {
        return 4;
    }

    return 0;
}

fun main() {
    println("status={}", make_message(false));
}
```

実行結果：

```text
status=2
```

データを出力しない故障パスでもBufferを解除します。すべての障害を1つのreturnに変えるよりも、各支店で何を所有しているかを明確に管理することが重要です。

整理自体が失敗するAPIもあります。元のジョブのエラーとクリーンアップエラーをどのように保存するかを決定します。この小さな例では、ジョブエラー番号が優先されます。より大きなプログラムは両方を記録することができます。

## 部分的な成功は自動的にキャンセルされません

ファイルにいくつかのバイトを書き込んだ後に書き込みが失敗すると、すでに書き込んだバイトは消えません。ネットワークでも、相手が一部のデータを受け取った可能性があります。同じ操作を最初から繰り返すと、重複記録が発生する可能性があります。

逆に、checkedbytecursor読み取りは失敗したときの位置と出力値を保存します。これらの関数は、入力をさらに受け取った後、同じ場所で再試行できます。 「失敗すると何も変わらない」というルールをすべてのAPIに適用せず、その関数のドキュメントを確認します。

## 回復可能なエラーとトラップ

無効なファイル パスまたは無効なユーザー入力は、呼び出し元が回復できるようにエラー値として報告できます。無効な実行時シフト カウントまたは浮動小数点から整数への変換によって発生するトラップは、別のメカニズムです。

引き続き実行する必要があるプログラムは、危険な操作の前に入力を確認する必要があります。 assertも、ユーザー入力の正常な失敗を処理する手段として乱用しません。ユーザーに再入力する機会を与える必要がある場合は、戻り結果を介して制御フローを続けます。

## 演習と解答例

配列要素をインデックスで読み取り、インデックスが負の場合、または長さ以上の場合に失敗する関数を作成します。ゼロが有効な要素値になる可能性があるため、結果構造体を使用します。

<!-- wave-example: book-error-solution -->
```wave
struct Lookup {
    ok: bool;
    value: i32;
}

fun at(values: ptr<i32>, length: i32, index: i32) -> Lookup {
    if (index < 0 || index >= length) {
        return Lookup {
            ok: false,
            value: 0
        };
    }

    return Lookup {
        ok: true,
        value: deref values[index]
    };
}

fun main() {
    var values: array<i32, 3> = [0, 10, 20];
    var first: Lookup = at(&values[0], 3, 0);
    var outside: Lookup = at(&values[0], 3, 3);

    if (first.ok) {
        println("value={}", first.value);
    }

    if (!outside.ok) {
        println("out of bounds");
    }
}
```

実行結果：

```text
value=0
out of bounds
```

ポインタとlengthが実際に読み取ることができる配列を表すのは呼び出し元の条件です。インデックスチェックだけで任意のアドレスまで安全にする関数ではありません。関数が責任を負う検査と、呼び出し元が保証する条件を分けて読んでください。
