---
translation_set_id: console-io-formatting
path: language/console-io-and-formatting
locale: vi
group: language
group_order: 2
order: 16
title: Đầu vào, đầu ra và định dạng của bảng điều khiển
summary: print, println, input Mô tả các câu và quy tắc giữ chỗ.
---

## câu lệnh đầu vào/đầu ra

Wave cung cấp `print`, `println` và `input` dưới dạng câu lệnh đầu vào/đầu ra của bảng điều khiển.

```wave
fun main() {
    var count: i32 = 0;
    input("{}", count);
    print("count = ");
    println("{}", count);
}
```

Mỗi câu kết thúc bằng `;`. Đối số đầu tiên phải là một chuỗi ký tự. Không thể sử dụng các biến hoặc chuỗi được tính toán làm đối số format.

## phần giữ chỗ

Chỉ có chính xác hai ký tự, `{}`, là phần giữ chỗ.

```wave
println("name = {}, score = {}", name, score);
```

Số lượng phần giữ chỗ và số lượng biểu thức theo sau phải hoàn toàn giống nhau. Nếu các số khác nhau thì đó là lỗi ngữ pháp.

```wave
println("{} {}", one);
// 오류: 자리표시자 2개, 값 1개
println("plain text", one);
// 오류: 자리표시자 없음, 값이 남음
```

Các dạng dấu ngoặc nhọn khác được để ở dạng văn bản thuần túy. Không có phần giữ chỗ được đặt tên hoặc đánh số trong ngữ pháp này.

## print và println

`print` in nguyên văn bản đã định dạng và `println` thêm ngắt dòng.

```wave
print("loading...");
println("done");
```

Đối số định dạng sử dụng các giá trị vô hướng như số nguyên, số dấu phẩy động, chuỗi và con trỏ. Mảng và cấu trúc không thể được sử dụng làm đối số định dạng.

## input Mục tiêu

`input` lưu trữ giá trị đọc ở đích, vì vậy tất cả các biểu thức sau format phải là vị trí có thể ghi.

```wave
var number: i32 = 0;
input("{}", number);
```

Các biến, trường và vị trí lưu trữ không được tham chiếu có thể được sử dụng làm mục tiêu. Chữ và kết quả tính toán không thể được sử dụng làm đầu vào.

Nếu tất cả các giá trị đầu vào không thể được chuyển đổi sang loại được yêu cầu thì chương trình sẽ thoát với trạng thái lỗi.

## ranh giới thời gian chạy

Các câu lệnh này sử dụng đầu vào và đầu ra của bảng điều khiển từ môi trường hosted. Trong môi trường độc lập, đầu vào/đầu ra do hạt nhân hoặc thiết bị cung cấp phải được xác định là một hàm hoặc ranh giới FFI.

## Giá trị và phạm vi đầu vào

Đầu vào bool chỉ chấp nhận `0` và `1`. Nó không hiểu 2 là true hoặc chấp nhận chuỗi `true` làm cùng một đầu vào. Đầu vào số nguyên phải nằm trong phạm vi độ rộng số nguyên mục tiêu. Các số nguyên 128, 256, 512 và 1024-bit cũng được xử lý dựa trên chiều rộng tổng thể của loại.

Lỗi định dạng, nằm ngoài phạm vi, trước khi yêu cầu đầu vào EOF không thành công. input tích hợp không phải là hàm trả về lỗi và nhập lại mà là hàm đầu vào chấm dứt quá trình khi xảy ra lỗi. Nếu bạn cần xử lý đầu vào có thể phục hồi, hãy đọc các byte bằng io và tạo một trình phân tích cú pháp riêng.

[Thực hành tính toán đầu vào](/docs/vi/practice/input-calculator) · [Tệp và io](/docs/vi/stdlib/files-io)
