---
translation_set_id: comments
path: language/comments
locale: zh
group: language
group_order: 2
order: 15
title: 评论
summary: 描述单行注释、可嵌套块注释和未封闭注释诊断。
---

## 一行评论

`//`之后的内容是注释，直到行尾。

```wave
var count: i32 = 10;
// 현재 요청 수
```

## 块注释

将`/*`和`*/`之间的空格处理为块注释。

```wave
/* 여러 줄에 걸친
   설명을 작성할 수 있습니다. */
```

您可以在块注释中嵌套其他块注释。

```wave
/* 바깥 주석
   /* 안쪽 주석 */
다시 바깥 주석
*/
```

## 字符串和注释标记

字符串和字符文本中的`//`、`/*`和`*/`是字符串内容，不被视为注释的开头或结尾。

```wave
var text: str = "https://wave-lang.dev";
```

## 未闭合的块注释

未能使用`*/`关闭块注释将导致诊断`E1002 UnterminatedComment`。

即使暂时禁用长块，也要确保嵌套深度正确。
