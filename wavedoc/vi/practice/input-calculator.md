---
translation_set_id: practice-input-calculator
path: practice/input-calculator
locale: vi
group: practice
group_order: 4
order: 1
title: Dự án: Một máy tính xác nhận đầu vào
summary: Liên kết đầu vào, kiểm tra giới hạn, chức năng và mã thoát.
---

## Mục tiêu và hành động

Nhập số lượng, đơn giá và tính tổng. Nhập hai số nguyên cách nhau bằng dấu cách hoặc ngắt dòng. Ví dụ này chỉ chấp nhận số lượng từ 1 đến 1000 và đơn giá từ 0 đến 100000 nên được tính trong phạm vi nhân i32.

Lưu nó vào `main.wave`, chạy `wavec run main.wave`, sau đó nhập `3 1200`. Kết quả đầu ra của chương trình như sau. Khả năng hiển thị của các ký tự được nhập trong thiết bị đầu cuối tách biệt với đầu ra của chương trình.

<!-- wave-example: input-calculator -->
```wave
fun total(quantity: i32, price: i32) -> i32 {
    return quantity * price;
}

fun main() -> i32 {
    var quantity: i32 = 0;
    var price: i32 = 0;
    input("{} {}", quantity, price);
    if (quantity < 1 || quantity > 1000 || price < 0 || price > 100000) {
        println("out of range");
        return 1;
    }

    println("total={}", total(quantity, price));
    return 0;
}
```

Kết quả thực hiện:

```text
total=3600
```

## Đồng thời kiểm tra lỗi

Nếu bạn đặt `0 1200`, nó sẽ yêu cầu `out of range` và mã thoát 1. Lỗi phân tích số trong `input` và việc kiểm tra phạm vi công việc của chương trình là hai bước khác nhau. Mã thông báo/loại không phải số vượt quá phạm vi/trước khi yêu cầu đầu vào EOF là lỗi đầu vào. Đầu vào tích hợp này không phải là giao diện trả về lỗi và yêu cầu nhập lại. Nếu bạn cần một trình phân tích cú pháp có thể phục hồi, hãy tự định cấu hình quy trình xác minh bằng cách đọc byte bằng [io](/docs/vi/stdlib/files-io).

## Bài tập và bình luận mở rộng

Lấy tỷ lệ chiết khấu làm đầu vào thứ ba và kiểm tra xem nó có phải là 0 đến 100 hay không. Để tránh các phép nhân trung gian lớn, bạn cần mở rộng phạm vi tính toán bằng i64 và kiểm tra phạm vi khi thu hẹp kết quả. Nếu bạn chỉ mở rộng loại kết quả cuối cùng, các phép tính trung gian có thể đã được thực hiện trên loại hẹp.

[Bảng điều khiển I/O](/docs/vi/language/console-io-and-formatting) · [Tiếp theo: Đọc tệp](/docs/vi/practice/file-reader)

## Tại sao phải quyết định phạm vi tính toán trước?

Đầu vào lớn nhất là số lượng 1000 và đơn giá 100000. Nhân hai giá trị là 100000000 nên nằm trong khoảng i32. Việc kiểm tra phạm vi này là cần thiết để có thể sử dụng kết quả nhân của hàm total.

main, nhận đầu vào, chịu trách nhiệm về thông báo đầu vào và lỗi, còn total chỉ chịu trách nhiệm tính toán. Ngay cả khi sau này bạn chuyển sang đọc lệnh từ tệp, bạn vẫn có thể sử dụng chức năng tính toán.

## Tính toán chiết khấu hoàn chỉnh

Bạn có thể nhân tổng số tiền với 100 bằng cách cộng phần trăm chiết khấu. Chuyển đổi để thực hiện các phép tính trung gian dưới dạng i64 rồi áp dụng tỷ lệ chiết khấu. Vì phần phân số của phép chia số nguyên bị loại bỏ nên ví dụ này sẽ cắt bớt số tiền sau khi chiết khấu thành số nguyên.

<!-- wave-example: book-calculator-discount -->
```wave
fun discounted_total(quantity: i32, price: i32, discount: i32) -> i64 {
    var subtotal: i64 = (quantity as i64) * (price as i64);
    var remaining: i64 = 100 - (discount as i64);
    return subtotal * remaining / 100;
}

fun main() -> i32 {
    var quantity: i32 = 0;
    var price: i32 = 0;
    var discount: i32 = 0;

    input("{} {} {}", quantity, price, discount);

    if (quantity < 1 || quantity > 1000) {
        println("invalid quantity");
        return 1;
    }

    if (price < 0 || price > 100000) {
        println("invalid price");
        return 2;
    }

    if (discount < 0 || discount > 100) {
        println("invalid discount");
        return 3;
    }

    println("total={}", discounted_total(quantity, price, discount));
    return 0;
}
```

Kết quả thực hiện:

```text
total=3240
```

Kết quả trên là khi nhập `3 1200 10`. Đầu ra là 3240, được trừ 10% so với tổng số 3600 ban đầu.

|đầu vào|kết quả mong đợi|đường dẫn để kiểm tra|
| --- | --- | --- |
| `3 1200 0` | `total=3600` |không giảm giá|
| `3 1200 100` | `total=0` |Giảm giá đầy đủ|
| `3 1200 101` | `invalid discount` |Tỷ lệ chiết khấu vượt quá phạm vi|
| `1000 100000 0` | `total=100000000` |đầu vào tối đa|
| `0 1200 10` | `invalid quantity` |Số lượng dưới phạm vi|

## bài tập tiếp theo

Hãy thử thay đổi hàm thành số thập phân làm tròn. Trong chương trình chỉ chấp nhận số dương này, việc cộng 50 trước khi chia cho 100 sẽ làm tròn thành số nguyên. Khi bạn nhập 99 cho 1 và tỷ lệ chiết khấu là 50, bạn có thể so sánh kết quả cắt của 49 và kết quả làm tròn của 50.
