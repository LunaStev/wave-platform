---
translation_set_id: stdlib-files-io
path: stdlib/files-io
locale: ja
group: stdlib
group_order: 1
order: 7
title: fs と io: ファイルとバイト転送
summary: ファイルの寿命、完全な読み取り、部分的な転送、および失敗後の状態を説明します。
---

## ファイルAPI選択

`std::fs::file`の便宜関数は経路を受けて必要なオープン・クローズを行います。記述子を返す関数は呼び出し元が閉じる必要があります。

|宣言|成功結果と注意点|
| --- | --- |
| `open_read(path: str) -> i64` |オープン記述子。負の数はエラーです|
| `create(path: str) -> i64` |作成または既存のファイルの内容を消去します。所有記述子を返す|
| `open_append(path: str) -> i64` |追加用のオープンまたは作成|
| `size(path: str) -> i64` |バイト数。負の数はエラーです|
| `read_into(path: str, dst: ptr<u8>, dst_cap: i64) -> i64` |ファイル全体のバイト数。容量不足はエラーです|
| `read_to_end(path: str, dst_buffer: ptr<Buffer>) -> i64` |既存のBufferの後にファイルを追加して追加量を返す|
| `write(path: str, src: ptr<u8>, len: i64) -> i64` |既存の内容を置き換えて書き込んだバイト数|
| `append(path: str, src: ptr<u8>, len: i64) -> i64` |最後に追加したバイト数|
| `remove(path: str) -> i64` |除去状態。失敗は負|

`exists(path)`のfalseだけではないファイルと権限エラーを区別できません。存在を確認した後、開いている間も状態が変わる可能性があるため、実際に開く結果を必ず確認してください。

## 低レベル I/O

```text
std::io::fd
io_read(fd: i64, buf: ptr<u8>, len: i64) -> i64
io_write(fd: i64, buf: ptr<u8>, len: i64) -> i64
io_read_exact(fd: i64, buf: ptr<u8>, len: i64) -> i64
io_write_all(fd: i64, buf: ptr<u8>, len: i64) -> i64
io_close(fd: i64) -> i64
```

`io_read`の正の結果は読み取られたバイト数であり、正の長さ要求では、0はEOFです。 `io_write`は要求より少なく書くことができます。フル転送が必要な場合は、exact/all関数を使用してください。それでも、エラーの前にいくつかの転送が発生する可能性があるため、障害が外部状態を元に戻すとは想定されません。

`io_read_exact`は、必要な長さの前にEOFに会うとエラーです。 `read_into`は、バッファが足りない場合は`IO_ERR_NO_SPACE`を返し、一部のバイトがすでに書き込まれている可能性があります。読み込み関数は、文字列の末尾のNULを自動的に貼り付けません。

## Bufferとエラー処理

障害が発生すると、`read_to_end` は元の len を復元しますが、その容量とデータ アドレスは変更されている可能性があります。呼び出し元は、成功または失敗の後に、Buffer を解放する必要があります。ファイル書き込み API は、アトミックなファイル置換を保証しません。

[ファイルの読み方](/docs/ja/practice/file-reader)からimportから解除までのプログラムを実行できます。 Linux/macOS/Windows/FreeBSDのパス・権限の違いとWASIのアクセス可能なディレクトリ制限を考慮してください。記述子の値を他のOSの生のハンドルとして直接解釈しません。

## 大きなファイルを小さなバッファに読み込む

ファイル全体をメモリーに入れる必要のない操作は、固定バッファーと読み取り反復で処理できます。次のプログラムは、input.txtの内容を出力し、読み込んだ全バイト数を数えます。入力ファイルには、`Wave`とLFのいずれかを保存します。

<!-- wave-example: book-file-stream -->
```wave
import("std::fs::file")::{open_read};
import("std::io::fd")::{io_read, io_write_all, io_close};
import("std::io::consts")::{IO_STDOUT_FD};

fun main() -> i32 {
    var descriptor: i64 = open_read("input.txt");

    if (descriptor < 0) {
        return 1;
    }

    var buffer: array<u8, 4>;
    var total: i64 = 0;

    while (true) {
        var count: i64 = io_read(descriptor, &buffer[0], 4);

        if (count < 0) {
            io_close(descriptor);
            return 2;
        }

        if (count == 0) {
            break;
        }

        if (io_write_all(IO_STDOUT_FD, &buffer[0], count) < 0) {
            io_close(descriptor);
            return 3;
        }

        total += count;
    }

    if (io_close(descriptor) < 0) {
        return 4;
    }

    println("bytes={}", total);
    return 0;
}
```

実行結果：

```text
Wave
bytes=5
```

バッファ容量は4ですが、最後の読み出しは1バイトです。出力には常に実際のcountを渡します。配列全体を書き込むと、未読の古いバイトまで出力できます。

プログラムが開いたのはinput.txtのdescriptorなのでそれを閉じます。標準出力は、この関数から新しく取得した所有リソースではないため、例の最後ではランダムに閉じません。

## 完全な読み取りと容量不足

read_intoは、ファイル全体を含む固定ストレージスペースを受け取ります。スペースが足りない場合は、静かにカットして成功せず、NO_SPACEを返します。

<!-- wave-example: book-file-capacity -->
```wave
import("std::fs::file")::{read_into};
import("std::io::consts")::{IO_ERR_NO_SPACE};

fun main() {
    var data: array<u8, 2>;
    var status: i64 = read_into("input.txt", &data[0], 2);

    if (status == IO_ERR_NO_SPACE) {
        println("destination too small");
    }
}
```

実行結果：

```text
destination too small
```

失敗した読み取りが宛先バイトをまったく変更しなかったとは想定しません。完成したファイルの内容として使用するのではなく、より大きなストレージを準備するか、streaming方法を選択してください。サイズを最初に照会しても、照会と読み取りの間でファイルが変更される可能性があるため、実際の読み取り結果は最終的な判断基準です。

## ファイルの書き込みと追加の違い

writeは既存の内容を置き換え、appendは最後に追加します。格納するバイト数は、文字列の長さから直接取得して渡すことができます。文字列の末尾のNULは通常テキストファイルの内容には含まれません。

<!-- wave-example: book-file-write-append -->
```wave
import("std::fs::file")::{write, append, size, remove};
import("std::string::len")::{len};

fun main() -> i32 {
    var first: str = "Wave";
    var second: str = " study";

    if (write("output.txt", first as ptr<u8>, len(first) as i64) < 0) {
        return 1;
    }

    if (append("output.txt", second as ptr<u8>, len(second) as i64) < 0) {
        return 2;
    }

    println("bytes={}", size("output.txt"));

    if (remove("output.txt") < 0) {
        return 3;
    }

    return 0;
}
```

実行結果：

```text
bytes=10
```

この例では、作業ディレクトリのoutput.txtを作成・置換し、最後に削除します。既存のファイルがない練習ディレクトリで実行します。実際のエディタまたは保存プログラムには、一時ファイルや置換などの別々の保存ポリシーが必要な場合があります。

## API選択表

|状況|選択|
| --- | --- |
|小さなファイル全体を固定バッファに読み込む| read_into |
|サイズがわからず、内容全体を保管|read_to_endとBuffer|
|コンテンツ全体をアーカイブせずに順番に処理|open_read+io_read繰り返し|
|指定された長さのレコードを読む| io_read_exact |
|バイト列全体の転送| io_write_all |
|すでに開いているファイルを扱う|パス関数の代わりにfd関数|

関数選択後は、その失敗でバッファ・ファイル位置・外部データがどのように変わるかを確認します。
