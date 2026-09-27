---
translation_set_id: expressions
path: language/expressions-and-operators
locale: ja
group: language
group_order: 2
order: 3
title: 3. 算術、比較、変換
summary: 計算順序、整数除算、ビット演算、castを学びます。
---

## 計算結果と計算タイプを一緒に見る

式は値を計算するコードです。変数名、リテラル、関数呼び出し、複数の値を演算子で連結した式はすべて式です。数学で同じ式のように見えても整数なのか間違いなのか、何ビットなのかによって結果が変わります。

この章では、簡単な計算から始めて、括弧、除算、論理演算、ビット単位の演算、キャストについて説明します。各サンプルは、`wavec run main.wave` で実行できる完全な main.wave ファイルです。

## かっこで囲む範囲

<!-- wave-example: book-precedence -->
```wave
fun main() {
    var first: i32 = 2 + 3 * 4;
    var second: i32 = (2 + 3) * 4;

    println("{} {}", first, second);
}
```

実行結果：

```text
14 20
```

乗算が加算より先に計算され、最初の式は2+12になります。 2番目の式では、括弧内の合計5を求め、次に4を掛けます。かっこを書くのが目標ではありません。読者が計算範囲をわかりやすくするために使用することをお勧めします。

同じ演算子を何度も書くときにも結び付ける方向が重要です。 `20 - 5 - 3`は`(20 - 5) - 3`なので12です。 `20 - (5 - 3)`は18です。正確な完全な順序は[演算子参照](/docs/ja/language/expressions-and-operators)にあります。

## 整数除算と残り

2 つの整数を除算しても、小数浮動小数点の結果は生成されません。商には除算を使用し、剰余には剰余演算子を使用します。

<!-- wave-example: book-division -->
```wave
fun main() {
    var items: i32 = 17;
    var box_size: i32 = 5;
    var full_boxes: i32 = items / box_size;
    var remaining: i32 = items % box_size;

    println("boxes={} remaining={}", full_boxes, remaining);
    println("negative quotient={}", -17 / 5);
}
```

実行結果：

```text
boxes=3 remaining=2
negative quotient=-3
```

17 個のアイテムを 5 個のグループに分割すると、完全なグループが 3 つと残りのアイテムが 2 つ生成されます。符号付き除算はゼロに向かって切り捨てられるため、-17/5 は -3 になります。これは、負の無限大に向かって切り捨てることとは異なります。

0 で割ることはできません。 signed最小値を-1で割った値も同じタイプには入りません。これらの入力を受け取る関数は、分割する前にチェックするか、checked数学APIを使用する必要があります。

## 変換する時点が結果を変える

次の両方の式はf64変数に格納されますが、計算プロセスは異なります。

<!-- wave-example: book-division-conversion -->
```wave
fun main() {
    var already_divided: f64 = (7 / 2) as f64;
    var floating_division: f64 = (7 as f64) / (2 as f64);

    if (already_divided == 3.0) {
        println("integer division lost the fraction");
    }

    if (floating_division == 3.5) {
        println("floating division kept the fraction");
    }
}
```

実行結果：

```text
integer division lost the fraction
floating division kept the fraction
```

最初の式は、整数の除算を実行して 3 を取得し、それを f64 に変換します。 2 番目の関数は、浮動小数点除算を実行する前にオペランドを f64 に変換します。最終変数に対してより幅広い型を選択すると、以前に失われた情報を回復できなくなります。

浮動小数点値は近似値です。同じ 10 進数値のように見える 2 つの結果は、正確な等価比較には適していない可能性があります。問題の単位と規模に適合する許容値を選択してください。 1 つの固定イプシロンがすべての計算に適しているわけではありません。

## 比較はboolになります

`<`、`<=`、`>`、`>=`、`==`、`!=`は関係を調べます。等号一つ`=`は代入であり、二人`==`は同等比較です。

<!-- wave-example: book-comparisons -->
```wave
fun main() {
    var age: i32 = 20;
    var minimum: i32 = 18;
    var eligible: bool = age >= minimum;

    if (eligible) {
        println("eligible");
    }

    if (age != minimum) {
        println("not exactly the boundary");
    }
}
```

実行結果：

```text
eligible
not exactly the boundary
```

「少なくとも 18 人」には 18 人が含まれます。 「18歳以上」は除きます。この境界を確認するには、17、18、19 をテストします。混合型の比較は符号の有無と幅に依存するため、両方のオペランドを目的の型に変換すると比較がより明確になります。

## 論理演算と短絡評価

`&&`はどちらも真であるか、`||`は一つ以上真であるか検査し、`!`は真・偽を覆します。この操作は右式を常に実行しません。

<!-- wave-example: book-short-circuit -->
```wave
fun report() -> bool {
    println("right side evaluated");
    return true;
}

fun main() {
    var enabled: bool = false;

    if (enabled && report()) {
        println("both true");
    }

    if (!enabled || report()) {
        println("at least one true");
    }
}
```

実行結果：

```text
at least one true
```

最初の条件では、enabledが偽であるため、右側を見る必要はありません。 2番目では、!enabledが真なので、やはり右は必要ありません。そのため、reportの出力は一度も表示されません。

これを使用して除算前に分母を確認できます。関数本文の断片`if (divisor != 0 && value / divisor > 2) { ... }`では、分母が0のときに除算を行いません。ただし、この検査だけでsigned最小値/-1のような他の境界まで解決されるわけではありません。

`&&`が`||`よりも優先されます。複合ポリシーでは、`(member && active) || admin`のように括弧で意図を明らかにします。

## 整数を絞り込むか広げる

`as`は明示的な型変換です。整数を絞り込むと捨てる上位ビットは再び広がっても回復しません。

<!-- wave-example: book-cast-chain -->
```wave
fun main() {
    var original: i32 = 300;
    var small: u8 = original as u8;
    var widened: i32 = small as i32;

    println("{} -> {} -> {}", original, small, widened);

    var negative: i8 = -1;
    var signed_value: i32 = negative as i32;
    var same_bits: u8 = negative as u8;

    println("{} {}", signed_value, same_bits);
}
```

実行結果：

```text
300 -> 44 -> 44
-1 255
```

300の低い8ビットは44です。 -1をsignedタイプに広げると符号を拡張して-1を維持します。同じ8ビットをunsignedと解釈すると255です。

範囲が正しい値を保存したい変換と記憶ビットを扱おうとする変換は目的が異なる。ユーザー入力を受け取ったら、まず目的地の範囲であることを確認して変換します。 castがあるため、値が安全な範囲であったことは保証されません。

## boolに置き換える

整数を bool に変換すると、ゼロの場合は false、それ以外の場合は true が得られます。これは最下位ビットまで切り捨てられません。2 も true に変換されます。浮動小数点値の場合、+0.0 と -0.0 のみが false に変換されます。 NaN と無限大を含む他のすべての値は、true に変換されます。

<!-- wave-example: book-bool-conversion -->
```wave
fun main() {
    var zero: i32 = 0;
    var two: i32 = 2;

    if (!(zero as bool)) {
        println("zero is false");
    }

    if (two as bool) {
        println("two is true");
    }
}
```

実行結果：

```text
zero is false
two is true
```

ポインタからboolへの変換はサポートされていません。ポインタを null と明示的に比較します (たとえば、`pointer != null`)。 null 以外のアドレスを安全に読み取れるかどうかは別の問題です。

## ビット演算

`&`、`|`、`^`、`~`は整数の各ビットを扱います。権限や機能をビットで表現する場合に使用できます。以下では、1 が読み取り権限、2 が書き込み権限です。

<!-- wave-example: book-bit-flags -->
```wave
const READ: u32 = 1;
const WRITE: u32 = 2;

fun main() {
    var permissions: u32 = READ | WRITE;

    if ((permissions & WRITE) != 0) {
        println("write enabled");
    }

    permissions = permissions & ~WRITE;
    println("remaining={}", permissions);
}
```

実行結果：

```text
write enabled
remaining=1
```

ORでビットを追加し、ANDで特定のビットがあるかどうかを確認します。 `~WRITE`で対応するビットのみ0のマスクを作成して消去します。ビット演算の`&`・`|`は、boolの段落評価`&&`・`||`とは異なる演算子です。

## シフトの幅と回数

左シフトはビットを左に移動し、オペランドの幅外の上位ビットを破棄します。右シフトでは、符号付きの値には符号拡張が使用され、符号なしの値にはゼロ拡張が使用されます。結果には常に左オペランドの型が含まれます。

<!-- wave-example: book-shifts -->
```wave
fun main() {
    var bits: u8 = 129;
    var shifted: u8 = bits << 1;
    var negative: i32 = -8;
    var positive: u32 = 8;

    println("left={}", shifted);
    println("right={} {}", negative >> 1, positive >> 1);
}
```

実行結果：

```text
left=2
right=-4 4
```

u8 値 129 を左に 1 シフトすると、その最上位ビットが破棄され、2 が残ります。元のシフトカウント値は非負であり、左オペランドのビット幅より小さくなければなりません。 u8 の場合、有効なカウントは 0 ～ 7 です。無効な定数カウントはコンパイル時エラーです。無効な実行時間カウントによりトラップが発生します。

## 浮動小数点値を整数に変換する

浮動小数点から整数への変換では、まずゼロに向かって切り捨てられ、次に変換先の整数の範囲がチェックされます。 NaN、無限大、範囲外の結果は無効です。無効な定数変換ではコンパイル時エラーが発生します。無効なランタイム変換によりトラップが発生します。

トラップは関数からエラー値を返しません。回復可能な変換エラーの場合は、変換前に範囲をチェックするインターフェイスを設計します。浮動小数点値の格納されたビットを取得するには、数値キャストの代わりに、`std::math::float` のビット変換関数を使用します。

## 演習と解答

残高137ウォンを50ウォン単位と残りに分けて、0～255の範囲の整数のみをu8に置き換える関数を作成します。この例では、範囲外を -1 と表示する別の関数を使用して検証手順を示します。

<!-- wave-example: book-expression-solution -->
```wave
fun checked_byte(value: i32) -> i32 {
    if (value < 0 || value > 255) {
        return -1;
    }

    var small: u8 = value as u8;
    return small as i32;
}

fun main() {
    var money: i32 = 137;

    println("coins={} remainder={}", money / 50, money % 50);
    println("{} {} {}", checked_byte(0), checked_byte(255), checked_byte(256));
}
```

実行結果：

```text
coins=2 remainder=37
0 255 -1
```

-1 を失敗表紙として書けるのは、成功範囲が 0～255 だからです。すべての整数が成功値である場合は、別の結果表現が必要です。後のエラー処理の章でこの設計を続けます。


## 優先順位

演算子の優先順位は、高い順序から次のようになります。

1. 基本式と後位アプローチ：関数呼び出し、フィールドアクセス、インデックス付け、後位`++`・`--`
2. 単項演算：`!`、`~`、`&`、`deref`、電位`++`・`--`、単項`+`・`-`
3. `as`型変換
4. `*`, `/`, `%`
5. `+`, `-`
6. `<<`, `>>`
7. `<`, `<=`, `>`, `>=`
8. `==`, `!=`
9. ビット`&`
10. ビット`^`
11. ビット`|`
12. `&&`
13. `||`
14. 代入と複合代入

連鎖代入は右から結合します。複数の演算子を混在させる場合は、括弧で計算順序を明確に表現してください。

## 代入可能な対象

代入、`++`、および `--` には、変数、フィールド、配列要素、逆参照ポインターなどの格納場所を示す式が必要です。 `const` への書き込みは許可されていません。

## シフト

シフトの結果は常に左オペランドの型を持ちます。右側のオペランドの型によって計算の範囲が広がることはありません。左シフトでは、その幅を超える上位ビットが破棄されます。右シフトは符号付き値を符号拡張し、符号なし値をゼロ拡張します。

シフト数は、元の値が `0 <= n < LHS bit width` を満たす整数でなければなりません。これは、より小さい型に切り詰められる前にチェックされます。無効な定数カウントはコンパイル時エラーです。無効な実行時間カウントによりトラップが発生します。

## ブール値と浮動小数点値を含む変換

整数からboolへの変換では、ゼロの場合はfalse、それ以外のすべての値の場合はtrueが生成されます。浮動小数点からboolへの変換では、+0.0と-0.0の場合のみfalseが生成されます。 NaN と正または負の無限大は true を生成します。 bool へのポインターのキャストはサポートされていません。`pointer != null` と明示的に比較してください。

浮動小数点から整数への変換では、ゼロに向かって切り捨てられてから、変換先の範囲がチェックされます。 NaN、無限大、範囲外の結果は無効です。無効な定数変換はコンパイル時エラーとなります。無効なランタイム変換によりトラップが発生します。トラップは回復可能なエラーリターンではありません。

`&&`と`||`は短絡評価します。実行されない右オペランドの副作用は発生しません。 [演算クラス](/docs/ja/language/expressions-and-operators)では、小さい値で結果を確認できます。

## 意図的に失敗するシフト

8ビット値の移動数は0〜7でなければなりません。以下のプログラムは実行前に拒否する必要があります。

<!-- wave-example: reject-shift-count -->
```wave
fun main() {
    var value: u8 = 1;
    var result: u8 = value << 8;
}
```
