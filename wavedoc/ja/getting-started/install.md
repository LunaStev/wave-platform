---
translation_set_id: install
path: getting-started/install
locale: ja
group: getting-started
group_order: 1
order: 2
title: Waveのインストール
summary: Linux、macOS、WindowsにWaveをインストールし、最初のプログラムを実行します。
---

## LinuxとmacOS

ターミナルで次のコマンドを実行します。WaveとパッケージマネージャーのVexがインストールされます。

```shell
curl -fsSL https://wave-lang.dev/install.sh | bash -s -- latest
```

インストールが完了したら、新しいターミナルを開いてバージョンを確認します。

```shell
wavec --version
```

## Windows

PowerShellで次のコマンドを実行します。WaveとパッケージマネージャーのVexがインストールされます。

```powershell
irm https://wave-lang.dev/install.ps1 -OutFile install.ps1
powershell -ExecutionPolicy Bypass -File .\install.ps1 -Latest
```

インストールが完了したら、新しいPowerShellウィンドウを開いてバージョンを確認します。

```powershell
wavec --version
vex --version
```

## 最初の実行

次のコードを`main.wave`として保存します。

<!-- wave-example: install-stdlib -->
```wave
import("std::string::len")::{
    len
};

fun main() {
    println("Wave: {} bytes", len("Wave"));
}
```

ファイルを保存したディレクトリで実行します。

```shell
wavec run main.wave
```

実行結果：

```text
Wave: 4 bytes
```

標準ライブラリが見つからないというメッセージが表示された場合は、インストールしてから再度実行します。

```shell
wavec install std
wavec run main.wave
```

[次へ：最初のプログラム](/docs/ja/language/program-structure) · [トラブルシューティング](/docs/ja/reference/diagnostics)
