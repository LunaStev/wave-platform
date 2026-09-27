---
translation_set_id: comments
path: language/comments
locale: en
group: language
group_order: 2
order: 15
title: Comments
summary: Describes single-line comments, nestable block comments, and unclosed comment diagnostics.
---

## one line comment

The content after `//` is a comment until the end of the line.

```wave
var count: i32 = 10;
// 현재 요청 수
```

## block annotation

Process the space between `/*` and `*/` as a block comment.

```wave
/* 여러 줄에 걸친
   설명을 작성할 수 있습니다. */
```

You can nest other block comments within a block comment.

```wave
/* 바깥 주석
   /* 안쪽 주석 */
다시 바깥 주석
*/
```

## Strings and Comment Marks

`//`, `/*`, and `*/` within string and character literals are string content and are not treated as the beginning or end of a comment.

```wave
var text: str = "https://wave-lang.dev";
```

## Unclosed block comments

Failure to close the block comment with `*/` will result in the diagnosis `E1002 UnterminatedComment`.

Even when temporarily disabling long blocks, make sure the nesting depth is correct.
