---
translation_set_id: stdlib-resolution
path: stdlib/resolution
locale: ja
group: stdlib
group_order: 1
order: 8
title: net.resolve: アドレス解決
summary: 番号アドレス、発信者所有名のリスト、および結果の切り捨てについて説明します。
---

## 検索方法の選択

Linuxのデフォルトresolverは番号IPv4/IPv6アドレスと番号ポートを受け取ります。ホスト名やサービス名を暗黙的にシステムDNSに渡しません。たとえば、`127.0.0.1`と`8080`は数値入力で、`example.com`と`http`は名前です。

明示的に提供された名前のリストを使用するには、`std::net::resolve_table`を使用します。番号アドレスに直接接続したり、名前リストにアドレスを登録して照会することができます。

## 基本検索

```text
std::net::resolve
net_resolve_tcp(host: str, service: str, output: ptr<SocketAddr>, capacity: i64) -> NetResolveResult
net_resolve_udp(host: str, service: str, output: ptr<SocketAddr>, capacity: i64) -> NetResolveResult
net_resolve_host(host: str, service: str, output: ptr<SocketAddr>, capacity: i64) -> NetResolveResult
```

出力配列は呼び出し元によって提供されます。 capacityの単位はバイトではなく、`SocketAddr`要素の数です。 capacity=0とnull outputで個数のみ照会できます。内部ルックアップストアは戻り前にクリーンアップされ、出力配列の所有権を取得しません。

|結果フィールド|意味|
| --- | --- |
| `ok` |検索が成功したかどうか|
| `count` |サポートする住所結果の総数|
| `written` |出力配列に実際に書き込んだ要素数|
| `truncated` |出力容量が全体の結果より小さいかどうか|
| `error.kind` |INVALID、NOT_FOUND、UNSUPPORTEDなど正規化された分類|
| `error.native_code` |原因確認用ネイティブコード|

成功してもwrittenが0になることがあります。 output[0]を使用する前に`ok && written > 0`を確認してください。

## 直接提供する名前のリスト

```text
net_resolve_from_table(host: str, service: str, protocol: i32,
    entries: ptr<NetResolveEntry>, entry_count: i64,
    output: ptr<SocketAddr>, capacity: i64) -> NetResolveResult
```

`std::net::resolve_table`の`NetResolveEntry`にはhost、service、protocol、addressがあります。ホスト名はASCII大文字と小文字を区別せず、サービス名は正確に比較します。 TCP/UDP/すべてのプロトコルルックアップにはモジュールの定数を使用します。エントリと文字列は呼び出し元が所有し、呼び出し後に保持しません。

一致は入力順にコピーし、重複を保持します。出力スペースと入力項目の保管スペースは重複してはいけません。リストに名前がない場合はNOT_FOUNDであり、外部DNSに再試行しません。

[TCP 使い方](/docs/ja/stdlib/tcp) · [TCP クライアント実習](/docs/ja/practice/tcp-client)

## 番号アドレスの検索

検索結果を格納する配列を準備し、記録されたアドレスがあることを確認します。アドレスルックアップは接続自体を作成しないため、以下の例ではサーバーは必要ありません。

<!-- wave-example: book-resolve-numeric -->
```wave
import("std::net::address")::{
    SocketAddr
};
import("std::net::resolve")::{
    NetResolveResult,
    net_resolve_tcp
};

fun main() -> i32 {
    var addresses: array<SocketAddr, 4>;
    var result: NetResolveResult = net_resolve_tcp(
        "127.0.0.1",
        "8080",
        &addresses[0],
        4
    );

    if (!result.ok || result.written == 0) {
        println("address unavailable");
        return 1;
    }

    println("address ready");
    return 0;
}
```

実行結果：

```text
address ready
```

実際に接続するには、addresses[0]をtcp_connect_addr系列関数に渡します。アドレスを取得しても、そのポートでサーバーが稼働しているかどうかは、接続フェーズで確認できます。

## 名前リストで照会する

自分で作成したアドレスのリストは、設定ファイルやプログラムのサービスリストに基づいて名前を関連付ける場合に便利です。次の例では、apiという名前をローカル8080ポートに対応させています。

<!-- wave-example: book-resolve-table -->
```wave
import("std::net::address")::{
    SocketAddr,
    socket_addr_from_v4,
    socket_addr_v4_loopback
};
import("std::net::resolve")::{
    NetResolveResult
};
import("std::net::resolve_table")::{
    NetResolveEntry,
    RESOLVE_TCP,
    net_resolve_from_table
};

fun main() -> i32 {
    var entries: array<NetResolveEntry, 1>;
    entries[0] = NetResolveEntry {
        host: "api",
        service: "http",
        protocol: RESOLVE_TCP,
        address: socket_addr_from_v4(socket_addr_v4_loopback(8080))
    };

    var addresses: array<SocketAddr, 2>;
    var result: NetResolveResult = net_resolve_from_table(
        "API",
        "http",
        RESOLVE_TCP,
        &entries[0],
        1,
        &addresses[0],
        2
    );

    if (!result.ok || result.written != 1) {
        return 1;
    }

    println("matched={}", result.written);
    return 0;
}
```

実行結果：

```text
matched=1
```

APIとapiはホスト名比較で一致します。 httpとHTTPはサービス名の比較で異なるため、この例では一致しません。同じ名前で複数のアドレスを登録すると、入力順に結果が表示されます。小さな配列で一部しか受け取っていない場合はwrittenまでを読んで、完全なリストが必要な場合はcountに合ったスペースを用意して再度照会します。
