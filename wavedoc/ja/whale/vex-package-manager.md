---
translation_set_id: vex-package-manager
path: whale/vex-package-manager
locale: ja
group: whale
group_order: 1
order: 4
title: Vex パッケージマネージャー
summary: manifestベースWaveプロジェクト、Git・パス依存性、lockfile、オフラインビルドとwavec境界を説明します。
---

## 役割

VexはWaveのパッケージマネージャーでありビルドツールです。 Vexは`wavec`の上で動作します。プロジェクト構造と依存性の解釈はVexが担当し、コンパイラフラグとコンパイルパイプラインは`wavec`が担当します。

Vexコマンドはmanifestベースです。 `vex build`、`vex check`、`vex run`はraw`wavec`フラグを意図的に受け取りません。

## パッケージの作成

```shell
vex init
vex init --lib
```

アプリケーションは`src/main.wave`、ライブラリは`src/lib.wave`を使用します。パッケージのルート構造は次のとおりです。

```text
my_project/
├── src/
│   └── main.wave
├── vex.ws
├── vex.lock
└── .vex/
    └── deps/
```

`vex.ws`がmanifestです。 Vexは`.wson`拡張子のmanifestを使用しません。

```wson
{
    name = "my_project",
    version = 0.1.0,
    lib = false,
    description = "my_project Project",
    author = "unknown",
    license = "Unknown",
    dependencies = []
}
```

## ビルドコマンド

```shell
vex build [--target <triple>] [--release] [--dry-run] [--locked] [--offline]
vex check [--target <triple>] [--release] [--dry-run] [--locked] [--offline]
vex run   [--target <triple>] [--release] [--dry-run] [--locked] [--offline] [-- <args...>]
```

Vexのオプションは小さく保たれます。 emit、linker、CPU、ABIまたはdebugのようにコンパイラに依存する制御が必要な場合は、`wavec`を直接使用してください。特定のコンパイラを使用する場合は、`VEX_WAVEC=/path/to/wavec`を設定してください。

`Resolving`、`Fetching`、`Compiling`、`Checking`、`Running`、`Finished`のような進行段階はstderr stdoutに保持されます。

## Git中心依存性

Vex依存性は、ローカル`path`またはGitURLで指定します。 1つの依存関係には、2つの方法のうち1つのみを使用できます。

```wson
{
    name = "app",
    version = 0.1.0,
    dependencies = [
        { name = "local_math", path = "../local_math" },
        { name = "remote_math", git = "https://github.com/example/math.git", tag = "v0.1.0" }
    ]
}
```

Git依存性には、`branch`、`tag`、`rev`のうち最大1つのみ指定できます。すべての依存ルートには独自の`vex.ws`が必要です。 Vexは依存性manifestを再帰的に解釈し、競合するパッケージのアイデンティティを拒否し、管理するGitcheckoutを`.vex/deps/<name>`に保存します。

## lockfile契約

スキーマ v2 `vex.lock`は、全遷移依存性グラフと正確なGitcommitを記録します。 manifestでコミットしてください。同じmanifestと有効なlockfileを使用すると、branchやtagに戻りません。同じ依存グラフを選択します。

依存性が必要なコマンドは自動的に解釈され、次のコマンドで事前に準備することもできます。

```shell
vex fetch
vex update
vex update math shared_core
```

`vex update`はすべてのGitパッケージを更新するか、指定されたパッケージと影響を受ける遷移グラフのみを更新します。関連のないロックされたパッケージは、選択されたcommitを維持します。

## lockedとofflineワークフロー

`--locked`は`vex.lock`の作成と変更を禁止します。ファイルが存在しないか、サポートされていないスキーマであるか、manifestグラフと一致しないと失敗します。 lockfileにすでに固定されているcommitは必要に応じてインポートできます。

`--offline`はすべてのGitネットワーク操作を禁止します。必要なcheckoutとcommitはローカルにすでに存在している必要があります。

```shell
vex fetch --locked
vex build --locked --offline
```

これら2つのコマンドは厳密なCIワークフローです。ネットワークを書くことができるときに正確にロックされたcommitを準備し、次にネットワークやlockfile変更なしでコンパイルします。 dry-runは依存関係を取得したり、lockfileを書き換えたりしません。

## コンパイラの設定と情報

```shell
vex info
vex setup wavec
vex setup wavec --version <version>
vex --version
```

Vexは、実際のビルド前に`wavec`dry-runJSONスキーマを検証します。必要なスキーマを実装していないコンパイラは、不明な計画で実行せず、互換性エラーで拒否します。
