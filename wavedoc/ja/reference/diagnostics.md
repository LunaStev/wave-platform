---
translation_set_id: diagnostics
path: reference/diagnostics
locale: ja
group: reference
group_order: 5
order: 2
title: トラブルシューティング: インストールから実行まで
summary: 失敗したステップを区別し、再現可能な情報で原因を絞り込みます。
---

## まず、失敗フェーズを区別する

|観察した現象|まず確認する|次のアクション|
| --- | --- | --- |
|wavecコマンドが見つかりません|PATHと実行ファイルの場所|絶対パスで実行した後、PATH設定|
|実行に必要なファイルが見つかりません|インストールフォルダからファイルが欠落しているか|パッケージ全体を解凍してインストール|
|std import失敗| `wavec print std-path` |対応 std設置または`--std-root`指定|
|ソース位置とタイプエラー出力| `wavec check main.wave` |最初のエラーから修正して再スキャン|
|その他 OS・CPU 対象のビルド失敗|指定したtargetと対象環境|[クロスビルド設定](/docs/ja/whale/build-link-targets)確認|
|ビルド成功後に実行失敗|終了コード、入力、作業ディレクトリ|実行環境とAPIエラーチェック|

## 小さな診断例

以下は意図的に間違ったプログラム全体です。

```wave
fun main() {
    var count: i32 = 1;
    println("{}", missing);
}
```

`wavec check main.wave`は宣言されていない名前missingを指す必要があります。変数名をcountに修正してから、もう一度確認して実行します。診断フレーズ全体よりファイル・位置・原因に集中してください。後続のエラーは最初のエラーの結果である可能性があります。

## 実行ファイルが失敗したとき

Linux/macOSシェルでは実行直後に`echo $?`、PowerShellでは`$LASTEXITCODE`で終了コードを確認します。入力エラーと明示的`return 1`は同じ原因ではありません。シフト数や実数変換の誤ったランタイム値は、trapを引き起こす可能性があります。 [演算規則](/docs/ja/language/expressions-and-operators)をご確認ください。

相対ファイルパスは、ソースファイルの場所ではなく実行ジョブディレクトリの影響を受けます。ファイルの読み取り失敗を文字列の長さ0として扱うのではなく、戻りエラーを最初に確認してください。ネットワーク接続の失敗は、アドレスのルックアップ、サーバーの待機、権限、およびタイムアウトを分けて確認します。

## 問題報告に必要な情報

1. `wavec --version`出力と実行された正確なコマンド。
2. ホスト OS・アーキテクチャとは別に指定したtarget.
3. 選択したstdパスで使用したコンパイラソース。
4. 問題を再現する最小ソース、入力と必要なファイル。
5. 期待される結果、実際の結果、診断と終了コード。

パスワード・トークン・個人ファイルの内容は削除します。最小の例を減らして問題がなくなると、最後に削除した要素は手がかりです。 `--error-format=json`は、ツールが診断を収集するときに使用できます。

[インストール](/docs/ja/getting-started/install) · [コンパイラコマンド](/docs/ja/getting-started/compiler) · [ターゲットとリンク](/docs/ja/whale/build-link-targets)
