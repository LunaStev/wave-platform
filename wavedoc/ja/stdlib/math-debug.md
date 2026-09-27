---
translation_set_id: stdlib-math-debug
path: stdlib/math-debug
locale: ja
group: stdlib
group_order: 1
order: 14
title: math と debug: 数学ユーティリティと診断
summary: 範囲チェック付きの数学関数と診断出力を使用します。
---

## 範囲を調べる整数関数

```text
std::math::int
min_i32(a: i32, b: i32) -> i32
max_i32(a: i32, b: i32) -> i32
abs_i32_checked(x: i32) -> MathResult<i32>
clamp_i32_checked(x: i32, lo: i32, hi: i32) -> MathResult<i32>
div_floor_i32_checked(a: i32, b: i32) -> MathResult<i32>
div_ceil_i32_checked(a: i32, b: i32) -> MathResult<i32>
```

`MathResult<T>` は `value` と `error` を持ちます。`std::math::result` からエラー定数をインポートし、`error == MATH_ERROR_NONE` を確認してから `value` を使います。符号付き整数の最小値の絶対値は同じ型では表せないため、検査付き絶対値関数はエラーを返します。`clamp` は `lo > hi` を拒否します。除算は除数がゼロの場合と結果が範囲を超える場合を検査します。`floor` と `ceil` の丸め方は、ゼロ方向に小数部分を切り捨てる整数除算とは異なります。

## 浮動小数点値の分類

`std::math::float`の`is_nan_f64`、`is_infinite_f64`、`is_finite_f64`は特殊値を区別します。 f32関数も提供します。 `float_to_bits_f64(value: f64) -> u64`は記憶ビットを取得する関数であり、`value as u64`の数値変換とは異なります。 NaNは自分自身と同じではないので、`value == nan`で検査しません。

## 診断

`std::debug::core`の`debug_assert(condition: bool, message: str)`は偽の状態で診断後終了します。ユーザー入力のように正常に失敗する可能性のある状況は結果値として処理し、必ず成立しなければならないプログラム内部条件を確認するときにassertを使用します。

以下のプログラムを`main.wave`として保存して実行します。

<!-- wave-example: math-api -->
```wave
import("std::math::int")::{
    min_i32, max_i32
};
import("std::debug::core")::{
    debug_assert
};

fun main() {
    var low: i32 = min_i32(9, 4);
    var high: i32 = max_i32(9, 4);
    debug_assert(low <= high, "invalid range");
    println("{} {}", low, high);
}
```

実行結果：

```text
4 9
```

## 負の除算の丸め方向

-7を3で割る方法を比較します。 `/`は0に切り、-2になります。 floorは小さい整数-3、ceilはより大きい整数-2を選択します。

<!-- wave-example: book-math-rounding -->
```wave
import("std::math::int")::{
    div_floor_i32_checked,
    div_ceil_i32_checked,
    abs_i32_checked
};
import("std::math::result")::{
    MathResult,
    MATH_ERROR_NONE,
    MATH_ERROR_OVERFLOW
};

fun main() -> i32 {
    var floor: MathResult<i32> = div_floor_i32_checked(-7, 3);
    var ceil: MathResult<i32> = div_ceil_i32_checked(-7, 3);

    if (floor.error != MATH_ERROR_NONE || ceil.error != MATH_ERROR_NONE) {
        return 1;
    }

    println("truncate={} floor={} ceil={}", -7 / 3, floor.value, ceil.value);

    var absolute: MathResult<i32> = abs_i32_checked(-2147483648);

    if (absolute.error == MATH_ERROR_OVERFLOW) {
        println("absolute value is outside i32");
    }

    return 0;
}
```

実行結果：

```text
truncate=-2 floor=-3 ceil=-2
absolute value is outside i32
```

結果値のみを確認すると、失敗したときに含まれる代替値と実際の計算結果を区別できません。 errorを最初に検査する順序を守ります。負の座標を一定の大きさの区間に入れるときはfloor、必要な束数を上げるときはceilが便利です。

## どのエラーを処理する必要がありますか

|状況|エラー|処理例|
| --- | --- | --- |
|0で割る| `MATH_ERROR_DIVIDE_BY_ZERO` |分母入力を再受信|
|結果をタイプに入れることはできません| `MATH_ERROR_OVERFLOW` |より広いタイプで計算または入力を拒否|
|clampの最上位値が最低値より大きい| `MATH_ERROR_INVALID_ARGUMENT` |設定範囲を変更|

assertはエラーを修復するツールではありません。ユーザー入力のエラーは条件文と戻り値として処理し、計算を終えた後に必ず成立しなければならない内部条件をdebug_assertで確認します。
