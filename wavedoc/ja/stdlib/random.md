---
translation_set_id: stdlib-random
path: stdlib/random
locale: ja
group: stdlib
group_order: 1
order: 10
title: random: OS のランダム性を使用してバッファを埋める
summary: OS エントロピーでバッファを埋め、部分的な障害を処理します。
---

## API

```text
std::random::fill
random_available() -> bool
random_fill(buffer: ptr<u8>, size: i64) -> RandomFillResult
```

size はバイト数であり、呼び出し元がストレージを提供します。 `RandomFillResult` には、OK、書き込み済み、およびエラーが含まれています。成功すると、要求された長さと同じ値が書き込まれます。失敗した場合、write は有効な埋め込まれたプレフィックスを識別します。残りのバイトをランダム データとして使用しないでください。

`random_available`は、OS乱数機能がサポートされているかどうかを示します。個々の要求が成功したかどうかは、random_fillの結果で確認されます。 OSエントロピーのみを使用し、失敗時に時間値や弱いPRNGに置き換えません。 size=0はnullと一緒に渡しても成功します。負の長さまたは正の長さの場合、nullはエラーです。

## 実行例

<!-- wave-example: random-api -->
```wave
import("std::random::fill")::{
    RandomFillResult, random_fill
};

fun main() -> i32 {
    var bytes: array<u8, 16>;
    var result: RandomFillResult = random_fill(&bytes[0], 16);
    if (!result.ok) {
        println("random error={}", result.error);
        return 1;
    }

    println("filled={}", result.written);
    return 0;
}
```

実行結果：

```text
filled=16
```

`main.wave`で保存して実行します。バイトの内容は毎回異なるため、特定の値は予想されません。失敗した場合は、result.errorで原因を確認してください。乱数バイトを文字通り出力しないでください。必要に応じて別のエンコーディングを使用してください。

## リクエストが間違っている場合

ゼロバイトのリクエストは何も書き込む必要がないため成功します。宛先バッファがないため、正の長さで null を渡すと失敗します。次のプログラムは、メモリを割り当てずにこれらのケースを比較します。

<!-- wave-example: book-random-boundaries -->
```wave
import("std::random::fill")::{
    RandomFillResult,
    random_fill
};

fun main() -> i32 {
    var empty: RandomFillResult = random_fill(null, 0);

    if (!empty.ok || empty.written != 0) {
        return 1;
    }

    println("empty request succeeded");

    var invalid: RandomFillResult = random_fill(null, 16);

    if (invalid.ok) {
        return 2;
    }

    println("missing buffer rejected");
    return 0;
}
```

実行結果：

```text
empty request succeeded
missing buffer rejected
```

## 一部のみが埋め込まれたバッファ処理

16バイトを要求しましたが、失敗してwritten=8が返された場合、前の8バイトだけが満たされた状態です。 16バイトの識別子を作成する操作は成功した識別子ではないため、結果全体を破棄して失敗を渡します。残りの8バイトをゼロで埋めた後、成功として処理しないでください。

ランダムなバイトを整数の範囲にマッピングするには注意が必要です。均一に分布した u8 値に `% 10` を適用すると、256 は 10 で割り切れないため、6 ～ 9 より 0 ～ 5 の可能性が高くなります。この偏りを取り除くには、値 250 ～ 255 を拒否し、再度描画し、受け入れられた値にのみ剰余演算を適用します。

乱数バイトの記憶領域と寿命は、呼び出し側によって管理されます。配列を使用すると配列の範囲内で処理され、動的メモリを使用すると使用終了後に解放されます。戻り構造体はバッファを代わりに所有しません。
