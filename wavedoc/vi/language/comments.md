---
translation_set_id: comments
path: language/comments
locale: vi
group: language
group_order: 2
order: 15
title: Bình luận
summary: Mô tả các nhận xét một dòng, nhận xét khối có thể lồng nhau và chẩn đoán nhận xét không được tiết lộ.
---

## bình luận một dòng

Nội dung sau `//` là phần bình luận cho đến hết dòng.

```wave
var count: i32 = 10;
// 현재 요청 수
```

## chặn chú thích

Xử lý khoảng trắng giữa `/*` và `*/` dưới dạng nhận xét khối.

```wave
/* 여러 줄에 걸친
   설명을 작성할 수 있습니다. */
```

Bạn có thể lồng các nhận xét khối khác vào trong một nhận xét khối.

```wave
/* 바깥 주석
   /* 안쪽 주석 */
다시 바깥 주석
*/
```

## Chuỗi và dấu chú thích

`//`, `/*` và `*/` trong chuỗi ký tự và chuỗi ký tự là nội dung chuỗi và không được coi là phần đầu hoặc phần cuối của nhận xét.

```wave
var text: str = "https://wave-lang.dev";
```

## Chú thích khối chưa đóng

Việc không đóng nhận xét chặn bằng `*/` sẽ dẫn đến chẩn đoán `E1002 UnterminatedComment`.

Ngay cả khi tạm thời vô hiệu hóa các khối dài, hãy đảm bảo độ sâu lồng nhau là chính xác.
