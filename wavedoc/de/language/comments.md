---
translation_set_id: comments
path: language/comments
locale: de
group: language
group_order: 2
order: 15
title: Kommentare
summary: Beschreibt einzeilige Kommentare, verschachtelbare Blockkommentare und die Diagnose nicht geschlossener Kommentare.
---

## einzeiliger Kommentar

Der Inhalt nach `//` ist ein Kommentar bis zum Zeilenende.

```wave
var count: i32 = 10;
// 현재 요청 수
```

## Blockanmerkung

Verarbeiten Sie das Leerzeichen zwischen `/*` und `*/` als Blockkommentar.

```wave
/* 여러 줄에 걸친
   설명을 작성할 수 있습니다. */
```

Sie können andere Blockkommentare innerhalb eines Blockkommentars verschachteln.

```wave
/* 바깥 주석
   /* 안쪽 주석 */
다시 바깥 주석
*/
```

## Zeichenfolgen und Kommentarzeichen

`//`, `/*` und `*/` innerhalb von Zeichenfolgen- und Zeichenliteralen sind Zeichenfolgeninhalte und werden nicht als Anfang oder Ende eines Kommentars behandelt.

```wave
var text: str = "https://wave-lang.dev";
```

## Nicht geschlossene Blockkommentare

Wenn der Blockkommentar nicht mit `*/` geschlossen wird, wird die Diagnose `E1002 UnterminatedComment` ausgegeben.

Auch wenn Sie lange Blöcke vorübergehend deaktivieren, achten Sie auf die richtige Verschachtelungstiefe.
