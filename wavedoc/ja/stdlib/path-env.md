---
translation_set_id: stdlib-path-env
path: stdlib/path-env
locale: ja
group: stdlib
group_order: 1
order: 12
title: path と env: パスと環境設定
summary: 呼び出し側バッファーにパスと環境変数を読み取り、容量エラーを区別します。
---

## パスの組み合わせ

```text
std::path::copy
path_join2(dst: ptr<u8>, dst_cap: i32, left: str, right: str) -> i32
path_basename_copy(dst: ptr<u8>, dst_cap: i32, path: str) -> i32
path_dirname_copy(dst: ptr<u8>, dst_cap: i32, path: str) -> i32
```

容量は最後のNULスペースを含みます。成功結果はNUL以外の長さで、失敗は-1です。成功した場合のみ、宛先を文字列として使用します。これらの関数はパス文字列を扱い、ファイルの存在やアクセス権をチェックしません。パスを組み合わせるだけでディレクトリのエスケープをブロックしたり、実際のファイルの同一性を確認したりすることはできません。

## 環境変数

```text
std::env::environ
env_get(name: str, dst: ptr<u8>, dst_cap: i64) -> i64
env_exists(name: str) -> bool
env_get_i64(name: str) -> EnvResult<i64>
env_get_i32(name: str) -> EnvResult<i32>
```

`env_get`は成功時にNUL以外の長さを返します。呼び出し側バッファーは、NULまで収めることができなければなりません。空の値は長さ0で成功する可能性があるため、キーレスエラーとは異なります。

`std::env::consts`からNOT_FOUND、NO_SPACE、INVALID_KEY、READ、SOURCE_INCOMPLETE、NO_MEMORYバッファ不足をキーなしとして扱いません。数値照会は結果のokを確認し、valueを使用します。環境変数の内容を信頼できる設定として自動的に扱うのではなく、範囲と形式を確認してください。

次の例では、dataディレクトリとinput.txtファイル名を組み合わせています。

<!-- wave-example: path-api -->
```wave
import("std::path::copy")::{
    path_join2
};

fun main() -> i32 {
    var output: array<u8, 64>;
    var length: i32 = path_join2(&output[0], 64, "data", "input.txt");
    if (length < 0) {
        return 1;
    }

    println("{}", &output[0] as str);
    return 0;
}
```

実行結果：

```text
data/input.txt
```

## ディレクトリとファイル名を分割する

次の例では、パスを 2 つのバッファに分割してコピーします。元のファイルが実際に存在する必要はありません。

<!-- wave-example: book-path-parts -->
```wave
import("std::path::copy")::{
    path_basename_copy,
    path_dirname_copy
};

fun main() -> i32 {
    var directory: array<u8, 64>;
    var filename: array<u8, 64>;
    var directory_length: i32 = path_dirname_copy(&directory[0], 64, "data/report.txt");
    var filename_length: i32 = path_basename_copy(&filename[0], 64, "data/report.txt");

    if (directory_length < 0 || filename_length < 0) {
        return 1;
    }

    println("directory={}", &directory[0] as str);
    println("filename={}", &filename[0] as str);
    return 0;
}
```

実行結果：

```text
directory=data
filename=report.txt
```

2つのバッファは、mainの終わりまで有効です。 `as str`は、新しい文字列を割り当てずに同じバッファを文字列として読み込みます。したがって、バッファを変更すると、そのアドレスに読み込まれる文字列も変わります。

## 環境設定のデフォルト値を設定する

環境変数は、プログラム外で渡す設定です。数値設定を読むときは、「整数で読むことができるか」と「このプログラムで許可した範囲か」をそれぞれ調べます。

<!-- wave-example: book-env-setting -->
```wave
import("std::env::environ")::{
    EnvResult,
    env_get_i32
};

fun main() -> i32 {
    var setting: EnvResult<i32> = env_get_i32("WAVE_EXAMPLE_WORKERS");
    var workers: i32 = 4;

    if (setting.ok) {
        if (setting.value < 1 || setting.value > 32) {
            println("workers must be between 1 and 32");
            return 1;
        }

        workers = setting.value;
    }

    println("workers={}", workers);
    return 0;
}
```

`WAVE_EXAMPLE_WORKERS`がない場合、または整数で読み取れない場合は、デフォルト値4を使用してください。 1〜32の整数が設定されている場合はその値を使用し、範囲外の整数面エラーで終了します。

Linux/macOS:

```shell
WAVE_EXAMPLE_WORKERS=8 wavec run main.wave
```

PowerShell:

```powershell
$env:WAVE_EXAMPLE_WORKERS = "8"
wavec run main.wave
```

どちらの場合も、`workers=8`が出力されます。上記の例は、単純なデフォルトポリシーを選択したものです。必須設定の場合は、数値検索の失敗をデフォルト値に置き換えるのではなく、エラーとして扱います。キーなし、バッファ不足、読み取り失敗を区別する必要がある場合は、env_getとENV_ERR_*定数を使用してください。
