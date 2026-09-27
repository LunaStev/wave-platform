---
translation_set_id: learn-arrays-strings
path: language/arrays
locale: ja
group: language
group_order: 2
order: 6
title: 6. 配列と反復
summary: 固定サイズの配列、インデックス付け、トラバーサル、コピー、検索について学びます。
---

## 同じ種類の複数の値

3つのスコアをscore1、score2、score3で別々に作成する場合は、個数が変わったときに宣言と計算の両方を変更する必要があります。配列は同じ型の要素を指定された数だけ束ねます。繰り返しステートメントを使用すると、要素ごとに同じルールを適用できます。

この章では、配列の作成、インデックス付け、変更、反復、検索、集計について説明します。文字列もインデックス付けをサポートしていますが、その意味は異なるため、文字列については次の章で説明します。

## タイプに長さを書く

<!-- wave-example: book-array-create -->
```wave
fun main() {
    var scores: array<i32, 3> = [70, 80, 90];

    println("first={}", scores[0]);
    println("second={}", scores[1]);
    println("last={}", scores[2]);
}
```

実行結果：

```text
first=70
second=80
last=90
```

`array<i32, 3>`のi32は要素タイプで、3は要素数です。保存するのは整数3つです。バイト数3という意味ではありません。配列リテラルの要素数は、宣言された長さと一致する必要があります。

インデックスは0から始まります。最初の要素は0、最後の要素は長さ-1です。 scores[3]は3番目の要素ではなく、有効範囲外のアプローチです。

## 要素の変更

<!-- wave-example: book-array-update -->
```wave
fun main() {
    var scores: array<i32, 3> = [70, 80, 90];

    scores[1] = 85;
    scores[0] += 5;

    println("{} {} {}", scores[0], scores[1], scores[2]);
}
```

実行結果：

```text
75 85 90
```

配列全体を再作成せずに特定の要素の記憶領域を変更します。インデックス式も計算結果である可能性がありますが、その値が範囲内でないことを確認する必要があります。外部入力をインデックスとして書き込むときは、負の値と上限の両方を調べます。

## 配列の反復処理

<!-- wave-example: book-array-sum -->
```wave
fun main() {
    var scores: array<i32, 4> = [60, 70, 80, 90];
    var total: i32 = 0;

    for (var index: i32 = 0; index < 4; index += 1) {
        total += scores[index];
    }

    println("total={} average={}", total, total / 4);
}
```

実行結果：

```text
total=300 average=75
```

繰り返しごとに異なるインデックスの要素を読み込みます。合計変数は繰り返しの外で初期化する必要があります。繰り返しボディ内で毎回ゼロに初期化すると、最後の要素だけが残る式の誤った結果が得られます。

平均の計算に使用される整数の除算では、小数部分が切り捨てられます。浮動小数点平均の場合は、除算する前に合計を変換します。より大きな配列またはより大きな値の場合は、アキュムレータ タイプが合計を表現できることも確認してください。

## 一部の要素のみを集計する

条件文と巡回を組み合わせてフィルタリングできます。ここでは、80点以上の要素の数を数えます。

<!-- wave-example: book-array-filter -->
```wave
fun main() {
    var scores: array<i32, 5> = [60, 80, 90, 75, 100];
    var passed: i32 = 0;

    for (var index: i32 = 0; index < 5; index += 1) {
        if (scores[index] >= 80) {
            passed += 1;
        }
    }

    println("passed={}", passed);
}
```

実行結果：

```text
passed=3
```

インデックスと要素の値は区別する必要があります。 `index >= 80`を調べると、スコアではなく位置を比較します。どちらもi32であるため、タイプだけでこの意味エラーを見つけるのは難しいです。

## 最初の一致位置を探す

見つからない結果をどのように表すかを最初に決定します。この例では、有効なインデックスは0〜4なので、-1を失敗のカバーとして使用します。

<!-- wave-example: book-array-find -->
```wave
fun main() {
    var values: array<i32, 5> = [8, 3, 8, 1, 5];
    var target: i32 = 8;
    var found: i32 = -1;

    for (var index: i32 = 0; index < 5; index += 1) {
        if (values[index] == target) {
            found = index;
            break;
        }
    }

    if (found >= 0) {
        println("found at {}", found);
    } else {
        println("not found");
    }
}
```

実行結果：

```text
found at 0
```

最初の位置0も正常な結果です。 `found > 0`で成功を調べると、最初の要素が見つからなかったと誤解します。 boolでcastして成功を判断するのも同じ理由で間違っています。

breakを取り除くと、後の一致がfoundを上書きするため、最後の一致位置が得られます。 1つの文が関数の契約を変えることができるので、「探す」という説明も最初の位置か最後の位置かを具体的に書く必要があります。

## 配列要素のコピー

配列の値を別の記憶領域にコピーするには、要素ごとに読み取って代入できます。コピーが完了した後、片方の整数要素を変更しても、もう一方の整数要素は変わりません。

<!-- wave-example: book-array-copy -->
```wave
fun main() {
    var original: array<i32, 3> = [1, 2, 3];
    var copied: array<i32, 3>;

    for (var index: i32 = 0; index < 3; index += 1) {
        copied[index] = original[index];
    }

    copied[0] = 99;

    println("original={}", original[0]);
    println("copied={}", copied[0]);
}
```

実行結果：

```text
original=1
copied=99
```

要素がポインタの場合、要素をコピーするとアドレスもコピーされます。それらが指す個別のメモリを複製することはありません。この区別は、所有権を管理する際に重要になります。

## 初期化と有効範囲

初期値なしで宣言された配列のすべての要素を読み取ることができるとは思いません。一部しか記録していない場合は、実際に初期化した数を別々に管理する必要があります。ライブラリの読み取り関数が返す長さがバッファ全体の容量より小さい可能性がある理由も同じです。

配列の長さが型に含まれているため、実行中は任意に伸びません。サイズが大きくなるバイトのリストは、`Buffer`などの動的ストレージを使用します。配列の長さを変更するには、型と初期値、巡回上限とその長さに依存する計算を一緒に確認する必要があります。

## 練習：最大値と場所

配列 `[4, 9, 2, 9, 1]` で、最大値とその値が最初に表示される位置を取得します。すべての要素が負である可能性がある汎用関数である場合は、最大値をゼロに初期化しないでください。

### 解答例

<!-- wave-example: book-array-solution -->
```wave
fun main() {
    var values: array<i32, 5> = [4, 9, 2, 9, 1];
    var maximum: i32 = values[0];
    var position: i32 = 0;

    for (var index: i32 = 1; index < 5; index += 1) {
        if (values[index] > maximum) {
            maximum = values[index];
            position = index;
        }
    }

    println("max={} first={}", maximum, position);
}
```

実行結果：

```text
max=9 first=1
```

最初の要素を初期基準として、2番目から比較します。 `>`なので、同じ最コメント値が再び出ても位置を変えません。 `>=`に変更すると最後の位置になります。長さがゼロになる可能性があるインターフェースの場合は、最初の要素を読み取る前に空の入力を処理する必要があります。
