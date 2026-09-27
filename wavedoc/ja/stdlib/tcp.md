---
translation_set_id: stdlib-tcp
path: stdlib/tcp
locale: ja
group: stdlib
group_order: 1
order: 9
title: net.tcp: 接続と転送
summary: TCP結果構造体、部分転送、切断の責任について説明します。
---

## 接続結果の確認

```text
std::net::tcp
tcp_connect_addr(addr: SocketAddr) -> NetResult<TcpStream>
tcp_bind_loopback(port: u16) -> NetResult<TcpListener>
tcp_accept(listener: TcpListener) -> NetResult<TcpStream>
tcp_read(stream: TcpStream, buf: ptr<u8>, size: i64) -> i64
tcp_write_all(stream: TcpStream, buf: ptr<u8>, size: i64) -> i64
tcp_read_exact(stream: TcpStream, buf: ptr<u8>, size: i64) -> i64
tcp_close(stream: TcpStream) -> NetError
tcp_close_listener(listener: TcpListener) -> NetError
```

接続関数の`NetResult<T>`には`ok`、`value`、`error`があります。成功した場合にのみvalueを使用してください。 `NetError`は正規化されたエラー分類とnative_codeを含み、`error.kind == NET_ERROR_NONE`で成功するかどうかを確認できます。この定数は`std::net::error`から取得されます。

接続したストリームとacceptで受け取ったストリームはそれぞれ閉じる必要があります。リスナーを閉じることがすでに受け入れられているすべてのストリームを閉じることを意味するわけではありません。構造体のコピーをそれぞれ閉じないでください。

## TCPはメッセージ境界を保存しません

一度書いた内容が一度読んで全部戻るとは思いません。長さプレフィックス、区切り文字、固定長など、プロトコルの境界を設定する必要があります。固定長は`tcp_read_exact`、長さが不明なストリームは読み取り反復とEOF処理で扱われます。

負の数はエラーであり、正の長さの読み取りでは、0は相手の終了です。 write_all失敗前に一部のデータが送信されることがあります。同じ内容を無条件に最初から再送信すると重複することがあります。

## 待機と時間制限

基本的なblocking関数は相手を長く待つことができます。時間制限が必要なプログラムは、`tcp_connect_addr_timeout`、`tcp_read_timeout`、`tcp_write_timeout`シリーズを選択します。ミリ秒の要因とエラーの結果を確認し、相手が応答しないパスも設計します。非同期は[task](/docs/ja/stdlib/task)と対応する非同期ネットワークAPIを一緒に使用します。

[住所検索](/docs/ja/stdlib/resolution) · [完成したローカルクライアント](/docs/ja/practice/tcp-client)

## 読み取り関数を選択する基準

|必要な動作|機能|確認する値|
| --- | --- | --- |
|到着したデータの一部を読む| `tcp_read` |負のエラー、ゼロ終了、正のバイト数|
|決まった長さだけ読む| `tcp_read_exact` |リクエストの長さをすべて読んだか|
|与えられた内容を最後まで送る| `tcp_write_all` |要求の長さと戻り値|
|待ち時間制限|timeoutシリーズ|戻り結果とtimeoutエラー|

4 バイトの長さフィールドの後に本文があるプロトコルの場合、最初に完全な長さフィールドを読み取り、長さが許容最大値を超えていないことを確認してから、本文用のスペースを割り当てます。ピアから提供された未検証の長さから大量のメモリを割り当てないでください。

## 接続を閉じる時点

読み取り関数がゼロを返しても、ローカルストリームリソースは残ります。読んだ後、tcp_closeを呼び出します。転送エラーの後にも同じクリーンアップパスを使用すると、正常・失敗パスのクローズを欠かせません。

リスナーは新しい接続を受け取るリソースであり、ストリームはすでに接続されている通信リソースです。サーバーを作成するときは、両方の種類を管理します。 acceptが成功するたびに新しいストリームが発生するため、処理後に閉じ、サーバーの受け入れの繰り返しが終わったらリスナーを閉じます。

## 直接実行してみる

[ローカルTCPクライアント実習](/docs/ja/practice/tcp-client)にサーバーとクライアントの完全なコードがあります。サーバーを最初に実行した場合、サーバーが存在しない場合、ポート番号が異なる場合を比較すると、アドレス検索と接続障害を区別できます。
