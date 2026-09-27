---
translation_set_id: memory-model
path: language/explicit-memory-type-model
locale: ja
group: language
group_order: 2
order: 9
title: 9. ポインター、値の変更、および有効期間
summary: アドレス、逆参照、ポインタを介した値の変更、およびダングリング ポインタについて学習します。
---

## 値と保存場所を区別する

整数４２とその整数が格納されるアドレスは異なる値である。ポインタは保存場所を指します。アドレスを渡すと、関数は呼び出し元の記憶領域を読み取ったり変更したりできます。

この章では、アドレスの取得、逆参照、元の値の変更、ポインタの演算、有効期間について説明します。動的割り当てについては次の章で説明します。ローカル変数と配列要素のアドレスから始めます。

## アドレスの取得と逆参照

<!-- wave-example: book-pointer-first -->
```wave
fun main() {
    var value: i32 = 42;
    var address: ptr<i32> = &value;

    println("value={}", value);
    println("through pointer={}", deref address);

    deref address = 99;
    println("changed={}", value);
}
```

実行結果：

```text
value=42
through pointer=42
changed=99
```

`&value`はアドレスを取得し、`deref address`はそのアドレスの値を読み書きする式です。 addressに99を保存したのではなく、addressが指す整数に99を書きました。変数 address自体は継続valueを指します。

`ptr<i32>` は、i32 ストレージにアクセスするためのポインター型です。この型は長さを記録したり、自動割り当て解除を提供したりしません。

## ポインタ自体を置き換える

<!-- wave-example: book-pointer-reassign -->
```wave
fun main() {
    var first: i32 = 10;
    var second: i32 = 20;
    var selected: ptr<i32> = &first;

    selected = &second;
    deref selected = 25;

    println("{} {}", first, second);
}
```

実行結果：

```text
10 25
```

`selected = &second`はポインタ変数に異なるアドレスを格納します。 firstの値は変更しません。その後、derefに値を書き込むとsecondが変わります。 「住所の変更」と「住所による値の変更」を文章別に区別すると混乱が減ります。

## 関数にオリジナルを置き換えさせる

<!-- wave-example: book-pointer-increment -->
```wave
fun increment(value: ptr<i32>) {
    deref value = deref value + 1;
}

fun main() {
    var count: i32 = 4;

    increment(&count);
    increment(&count);

    println("{}", count);
}
```

実行結果：

```text
6
```

この関数は、count の値 4 ではなく、count のアドレスを受け取ります。そのストレージを変更すると、呼び出し元の count も変更されます。この関数には、有効な書き込み可能な i32 アドレスが必要です。 null を渡すと、この要件に違反します。

各関数は、null を受け入れるかどうかを定義します。そうでない場合、呼び出し元は有効なアドレスを提供する必要があります。その場合、関数には null を処理するパスが含まれている必要があります。

## nullを処理する関数

<!-- wave-example: book-pointer-null -->
```wave
fun try_increment(value: ptr<i32>) -> bool {
    if (value == null) {
        return false;
    }

    deref value = deref value + 1;
    return true;
}

fun main() {
    var count: i32 = 7;

    if (!try_increment(null)) {
        println("no value");
    }

    if (try_increment(&count)) {
        println("count={}", count);
    }
}
```

実行結果：

```text
no value
count=8
```

null チェックでは、住所がない場合のみ処理されます。 null 以外の任意の数値をポインターに変換しても、有効なメモリーは作成されません。読み取りと書き込みには、有効な有効期間、十分なサイズ、正しい位置合わせ、および適切なアクセス許可も必要です。

## 配列アドレスと要素ごとのポインタ演算

<!-- wave-example: book-pointer-array -->
```wave
fun main() {
    var values: array<i32, 3> = [10, 20, 30];
    var first: ptr<i32> = &values[0];
    var second: ptr<i32> = first + 1;

    println("{}", deref first);
    println("{}", deref second);
    println("{}", deref first[2]);
}
```

実行結果：

```text
10
20
30
```

ポインタに 1 を加えると、ターゲット型の要素だけ移動します。 i32の次の要素とu8の次の要素は移動バイト数が異なります。 `first + 1`に再びタイプサイズを掛けて加算すると不要な位置に移動します。

ポインタインデックスも有効範囲内でなければなりません。 firstは配列長3を自分で覚えていないので、関数に範囲を渡すときは、ポインタと長さを一緒に受け取る型を使用します。

## 読み取る範囲を関数に渡す

<!-- wave-example: book-pointer-range -->
```wave
fun sum(values: ptr<i32>, count: i32) -> i32 {
    var total: i32 = 0;

    for (var index: i32 = 0; index < count; index += 1) {
        total += deref values[index];
    }

    return total;
}

fun main() {
    var values: array<i32, 4> = [2, 4, 6, 8];

    println("first two={}", sum(&values[0], 2));
    println("all={}", sum(&values[0], 4));
}
```

実行結果：

```text
first two=6
all=20
```

countの単位は要素数です。この関数は、countの読みやすいi32があるという条件を呼び出し元に要求します。実際の配列より大きな長さを超えると、契約に違反します。同じptr<u8>・i64組み合わせを使用するAPIでも、長さがバイト数か要素数かをドキュメントで確認する必要があります。

## 寿命：住所がいつまで有効か

ローカル変数は、その呼び出しとブロックの有効期間内で使用されます。関数内のローカル変数アドレスを返して呼び出し元が後で読み取ると、その記憶領域の寿命が終わっている可能性があります。

以下は、実行してはならない誤った設計の一部です。

```wave
fun invalid_address() -> ptr<i32> {
    var local: i32 = 42;

    return &local;
}
```

1つの値が必要な場合は、i32を返します。発信者が用意した記憶領域に書き込む必要がある場合は、ポインタを入力として受け取ります。呼び出し後も保持する別の記憶領域が必要な場合は、明示的に割り当てて解放責任を渡します。

## 同じ記憶領域を指す2つのポインタ

ポインタをコピーすると、同じアドレスを指す名前がもう1つ生成されます。メモリを複製しません。

<!-- wave-example: book-pointer-alias -->
```wave
fun main() {
    var value: i32 = 1;
    var first: ptr<i32> = &value;
    var second: ptr<i32> = first;

    deref second = 9;

    println("{} {}", value, deref first);
}
```

実行結果：

```text
9 9
```

secondで変えた結果がfirstを通しても見えます。元のメモリが解放されると、両方のポインタは使用できません。あるポインタ変数にnullを代入しても、別のコピーまで自動的には変わりません。

## 練習：2つの整数交換

2つのi32アドレスを受け取り、値を交換する関数を作成します。最初の値を上書きする前に一時変数に保存する必要があります。同じアドレスを2回渡しても値が保持されていることを確認してください。

### 解答例

<!-- wave-example: book-pointer-swap -->
```wave
fun swap(left: ptr<i32>, right: ptr<i32>) {
    var saved: i32 = deref left;

    deref left = deref right;
    deref right = saved;
}

fun main() {
    var first: i32 = 3;
    var second: i32 = 8;

    swap(&first, &second);
    println("{} {}", first, second);

    swap(&first, &first);
    println("{}", first);
}
```

実行結果：

```text
8 3
8
```

この関数では、両方のアドレスが有効な書き込み可能な整数ストレージを指す必要もあります。 null を処理するには、try_increment のように成功または失敗を示す結果を追加します。


## **Wave Explicit Memory Type Model**

Waveのポインタ設計は、**Wave Explicit Memory Type Model**に基づいています。このモデルは、ポインタと配列を文法的なトリックやライブラリの抽象化ではなく、言語レベルの明示的なメモリタイプとして定義します。

`ptr<T>`は`T`値を保存したメモリアドレスを指すタイプで、`array<T, N>`は`T`値`N`個を連続して保存する固定長メモリタイプです。そのため、関数引数、戻り値、構造体フィールドなどの型の中でも、ポインタと配列の構造がそのまま現れます。

## null

```wave
var buffer: ptr<u8> = null;
if (buffer == null) {
    println("no buffer");
}
```

`null`は、有効なメモリアドレスを指さないポインタ値です。 `null`は`ptr<T>`タイプにのみ代入でき、整数、ブール値、または配列値としては使用できません。

割り当て関数または検索関数は、結果がない場合に `null` を返すことがあります。このような結果を参照解除する前に、`null` を確認してください。 `null` ポインターを逆参照すると、有効なストレージにアクセスできません。

## ポインタ変換

アドレスや他のポインタ表現を置き換える必要がある場合は、`as`を使用してください。

```wave
var raw: i64 = 0;
var p: ptr<u8> = raw as ptr<u8>;
```

整数とポインタの間の変換は低レベルの境界でのみ使用し、ターゲットプラットフォームのアドレス幅とABIを考慮してください。
