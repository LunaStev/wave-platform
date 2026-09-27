---
translation_set_id: control-flow
path: language/control-flow
locale: ja
group: language
group_order: 2
order: 4
title: 4. 条件、ループ、境界値
summary: if、for、whileと繰り返し文の範囲を学びます。
---

## 実行パスを選択する

前の章のプログラムは文を上から下に実行しました。実際のプログラムは、入力と状態によって異なる作業を行う必要があります。条件文は実行パスを選択し、反復文は同じ規則を複数の値に適用します。

例はそれぞれmain.waveに保存して実行します。コードを読むときは、現在の変数値、次に調べる条件、実行する文の順序を紙に書き留めてください。結果を覚えるよりも、流れに沿った練習が重要です。

## ifとelse

条件は括弧で囲み、本文は中括弧で囲みます。次の例では、balanceを500または2000に置き換えて、どのブランチが実行されているかを確認してください。

<!-- wave-example: book-if-else -->
```wave
fun main() {
    var balance: i32 = 2000;
    var price: i32 = 1200;

    if (balance >= price) {
        balance -= price;
        println("bought, balance={}", balance);
    } else {
        println("not enough money");
    }
}
```

実行結果：

```text
bought, balance=800
```

両方のブロックを実行するわけではありません。条件が真の場合は最初のブロック、偽の場合はelseブロックを実行します。 balanceがpriceと同じ場合でも購入を許可するので、`>=`を書きました。 `>`に変更すると、同じ金額で動作が異なります。

## 複数の条件の順序

else ifで条件をつなぎます。上記で最初に満足する四半期を実行するので、大きな境界から検査するか、小さな境界から検査するかを考える必要があります。

<!-- wave-example: book-grade -->
```wave
fun grade(score: i32) {
    if (score < 0 || score > 100) {
        println("invalid");
    } else if (score >= 90) {
        println("A");
    } else if (score >= 80) {
        println("B");
    } else {
        println("C");
    }
}

fun main() {
    grade(95);
    grade(80);
    grade(79);
    grade(101);
}
```

実行結果：

```text
A
B
C
invalid
```

無効なスコアを最初に拒否してから評価を分割します。 `score >= 80`を一番前に置くと95度その分岐に入るので、A分岐に到達できません。各条件の正確さだけでなく、条件の順序も確認してください。

## 条件式で値を変更しない

if、while、または for 条件では、代入、複合代入、およびインクリメントまたはデクリメント演算は許可されません。比較には`==`を使用します。値を更新してテストするには、2 つの別々のステートメントを作成します。

関数内部フラグメントの正しい形式：

```wave
value = read_value();

if (value == expected) {
    println("matched");
}
```

read_valueとexpectedはこの部分で定義されていないため、そのまま実行するプログラム全体ではありません。ここに示す規則は「状態変更後の比較」です。

## while:条件が維持されている間

<!-- wave-example: book-while-countdown -->
```wave
fun main() {
    var remaining: i32 = 3;

    while (remaining > 0) {
        println("{}", remaining);
        remaining -= 1;
    }

    println("finished at {}", remaining);
}
```

実行結果：

```text
3
2
1
finished at 0
```

条件は本文に入る前にチェックします。最初からremainingが0の場合、本文を一度も実行しません。本文の最後の減少が欠けている場合、条件は真であり、繰り返しは終了しません。

繰り返しステートメントを作成したら、「何が終了条件に近づくのですか？」を確認してください。入力を待つループなら入力変化やEOFが、数値ループならインデックス更新がその役割をします。

## for:初期化・条件・更新

forは、繰り返しに必要な3つの部分を集めて表現します。

<!-- wave-example: book-for-sum -->
```wave
fun main() {
    var total: i32 = 0;

    for (var number: i32 = 1; number <= 5; number += 1) {
        total += number;
    }

    println("sum={}", total);
}
```

実行結果：

```text
sum=15
```

1. numberを1に初期化します。このステップは一度です。
2. number<= 5を調べます。偽の場合は繰り返しを終了します。
3. 本文でtotalにnumberを加えます。
4. numberを1増やして条件チェックに戻ります。

forの中で宣言した反復変数は、反復の後にも使用できるとは仮定しません。繰り返し後の値が必要な設計の場合は、外で宣言して初期化位置を明確にしてください。

## 包含境界と除外境界

1から5までの和には`<= 5`が自然です。一方、長さが 5 の配列のインデックスには `< 5` を使用する必要があります。配列インデックスは 0 から始まり 4 で終わるからです。

「5回実行」と「5という値まで含める」を混同しないでください。開始値と終了比較を一緒に見ると、繰り返し回数がわかります。入力が空の場合と要素が1つだけの場合は、境界の間違いを見つけることをお勧めします。

## continueとbreak

continueは今回の繰り返しの残りをスキップし、breakは最も近い繰り返し文を終了します。

<!-- wave-example: book-loop-control -->
```wave
fun main() {
    var total: i32 = 0;

    for (var number: i32 = 1; number <= 10; number += 1) {
        if (number == 3) {
            continue;
        }

        if (number == 6) {
            break;
        }

        total += number;
    }

    println("{}", total);
}
```

実行結果：

```text
12
```

合計に入った数は1、2、4、5です。 forでcontinueに会うと更新式に進みます。 whileにはforのような別途更新式がないため、continue以前に必要な状態変更を欠かさないように注意する必要があります。

反復文が入れ子になっている場合、break1つがすべての反復を終了するわけではありません。複数のフェーズで中断する必要がある場合は、関数を使用してタスクを束ねてreturnを使用するか、外側の繰り返しでも終了条件を確認するように意図を表します。

## matchでケースを分割する

同じ値の複数のケースを比較する場合、matchを使用できます。各armの本文はブロックです。

<!-- wave-example: book-match-number -->
```wave
fun describe(status: i32) {
    match (status) {
        200 => {
            println("ok");
        }
        404 => {
            println("missing");
        }
        _ => {
            println("other");
        }
    }
}

fun main() {
    describe(200);
    describe(404);
    describe(500);
}
```

実行結果：

```text
ok
missing
other
```

`_`は残りを処理するパターンです。同じmatchの中に重複して置きません。値によってデータの種類も変わるvariantは[データモデルの章](/docs/ja/language/structures-enums-and-aliases)で扱います。

## 完成例：条件に合った数の集計

1から10までのうち、偶数の数と合計を求めます。数と合計は異なる情報なので、それぞれ変数として累積します。

<!-- wave-example: book-loop-statistics -->
```wave
fun main() {
    var count: i32 = 0;
    var total: i32 = 0;

    for (var number: i32 = 1; number <= 10; number += 1) {
        if (number % 2 == 0) {
            count += 1;
            total += number;
        }
    }

    println("count={} total={}", count, total);
}
```

実行結果：

```text
count=5 total=30
```

偶数は2、4、6、8、10なので、数は5、合計は30です。式が短くても結果を手で求められる小さな範囲で最初に確認すると、繰り返し境界を検証しやすいです。

## 演習と解答例

1から20までのうち3の倍数だけ加算しますが、合計が30を超える値は加算せずに終了してください。 「もっと後に乗ったかどうかを確認する」と「超えていることを確認してからプラス」を区別する必要があります。

<!-- wave-example: book-loop-solution -->
```wave
fun main() {
    var total: i32 = 0;

    for (var number: i32 = 1; number <= 20; number += 1) {
        if (number % 3 != 0) {
            continue;
        }

        if (total + number > 30) {
            break;
        }

        total += number;
    }

    println("{}", total);
}
```

実行結果：

```text
30
```

値 3+6+9+12 の合計は 30 になるため、次の値 15 は加算されません。これらの小さな入力は安全ですが、大きな整数の場合、チェック `total + number` 自体がオーバーフローする可能性があります。小切手を書いても、すべての境界ケースが自動的に処理されるわけではありません。
