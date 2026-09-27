---
translation_set_id: practice-tcp-client
path: practice/tcp-client
locale: ja
group: practice
group_order: 4
order: 4
title: プロジェクト: ローカル TCP クライアント
summary: 番号アドレス検索、接続失敗、転送とクローズを接続します。
---

## 準備と範囲

このプログラムはネイティブOSの`127.0.0.1:8080`に接続し、`ping`とLFを転送して終了します。別の端末にテストサーバーが必要です。 Python3がある場合、次のサーバーは1つのローカル接続で最大5バイトを受信して​​表示します。

```python
import socket
with socket.socket() as server:
    server.bind(("127.0.0.1", 8080))
    server.listen(1)
    connection, address = server.accept()
    with connection:
        message = b""
        while len(message) < 5:
            chunk = connection.recv(5 - len(message))
            if not chunk:
                break
            message += chunk
        print(repr(message))
```

サーバーを最初に実行し、クライアントを`main.wave`として保存し、`wavec run main.wave`として実行します。サーバーは`b'ping\n'`を出力します。ポートがすでに使用されている場合は、サーバーとクライアントの両方のポートを一緒に置き換えます。

<!-- wave-example: tcp-client -->
```wave
import("std::net::address")::{
    SocketAddr
};
import("std::net::resolve")::{
    NetResolveResult, net_resolve_tcp
};
import("std::net::error")::{
    NetResult, NetError, NET_ERROR_NONE
};
import("std::net::tcp")::{
    TcpStream, tcp_connect_addr_timeout, tcp_write_all, tcp_close
};

fun main() -> i32 {
    var addresses: array<SocketAddr, 4>;
    var resolved: NetResolveResult = net_resolve_tcp("127.0.0.1", "8080", &addresses[0], 4);
    if (!resolved.ok || resolved.written == 0) {
        println("resolve failed");
        return 1;
    }

    var connected: NetResult<TcpStream> = tcp_connect_addr_timeout(addresses[0], 1000);
    if (!connected.ok) {
        println("connect failed");
        return 2;
    }

    var sent: i64 = tcp_write_all(connected.value, "ping\n" as ptr<u8>, 5);
    var closed: NetError = tcp_close(connected.value);
    if (sent != 5 || closed.kind != NET_ERROR_NONE) {
        return 3;
    }

    println("sent=5");
    return 0;
}
```

クライアントの成功出力は`sent=5`です。サーバを実行していない場合は、`connect failed`を出力して2で終了するのが正常な失敗処理です。

## コードを理解する

照会配列の有効範囲は、writtenまでです。最初のアドレスのみを使用する小さな例であり、複数の候補を持つサービスでは候補固有の接続ポリシーを設定する必要があります。接続したストリームは転送成功・失敗とは無関係に閉じます。接続タイムアウトは1000msであり、伝送全体のタイムアウトを保証する例ではありません。

TCPは伝送境界を保存しません。サーバーも何度もrecvできるように作成した理由です。名前で住所を見つける方法は[住所検索](/docs/ja/stdlib/resolution)で説明されています。

## 拡張練習

サーバーが応答を送信できるようにし、クライアントから読んでみてください。応答の長さを決定した後、部分読み取り、EOF、タイムアウトとクローズを一緒に処理する必要があります。

## 両端末の役割

サーバー端末は、listen状態で接続を待ちます。クライアントが接続すると、acceptが接続されたソケットを返し、recvイテレーションステートメントがデータを集めます。クライアント端末は、アドレスの照会、接続、送信、クローズを順番に実行します。

Pythonコードは実習相手サーバーです。 Waveプログラムが渡したバイトを簡単に確認しようとし、クライアントコードなどのファイルに入れません。サーバーは1つの接続を処理してからシャットダウンするため、クライアントを再実行する前にサーバーも再起動します。

## 5バイトを送る理由

`ping`は4バイト、LFは1バイトなので、転送長は5です。文字列の末尾のNULはメッセージに含まれません。サーバーは5バイトを集めた後に出力するので、パケットが複数回に分けて到着しても同じ結果が出ます。

tcp_write_allが5を返すと、ローカル転送関数が要求したバイトを処理します。相手プログラムがメッセージを解釈して保存まで完了したという確認ではありません。処理完了確認が必要なプロトコルの場合、サーバーは応答を送信し、クライアントがその応答を読み取るようにします。

## 失敗パスの練習

|変更|結果|学ぶポイント|
| --- | --- | --- |
|サーバーなしで実行|接続失敗、終了コード2|アドレスが有効であってもサーバーが存在しない可能性がある|
|サーバーとクライアントポートの異なる設定|接続失敗|住所のIPとポートの両方が正しいこと|
|サーバーのrecvサイズを1に変更|同じ5バイトが集まります|TCP 読み取り単位はメッセージ単位と異なる|
|クライアント転送を複数回に分ける|同じ順序で受信|メッセージ境界はプロトコルによって決まります|

接続タイムアウト1000msは接続フェーズの設定です。読み書きも制限するには、該当ステップのtimeout関数を使用し、プログラム全体に適用する制限時間があれば残り時間を計算して渡します。
