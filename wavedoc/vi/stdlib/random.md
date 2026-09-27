---
translation_set_id: stdlib-random
path: stdlib/random
locale: vi
group: stdlib
group_order: 1
order: 10
title: random: Làm đầy bộ đệm với tính ngẫu nhiên của hệ điều hành
summary: OS Đổ entropy vào bộ đệm và xử lý các lỗi một phần.
---

## API

```text
std::random::fill
random_available() -> bool
random_fill(buffer: ptr<u8>, size: i64) -> RandomFillResult
```

size là số byte và người gọi cung cấp dung lượng lưu trữ. `RandomFillResult` chứa ok, viết và lỗi. Khi thành công, văn bản bằng độ dài được yêu cầu. Khi thất bại, văn bản xác định tiền tố hợp lệ, được điền; không sử dụng các byte còn lại làm dữ liệu ngẫu nhiên.

`random_available` cho bạn biết liệu chức năng số ngẫu nhiên OS có được hỗ trợ hay không. Sự thành công của một yêu cầu riêng lẻ được kiểm tra bằng kết quả của random_fill. Chỉ sử dụng OS entropy và không giảm giá trị theo thời gian hoặc yếu PRNG khi thất bại. size=0 sẽ thành công ngay cả khi được thông qua cùng với null. null là lỗi về độ dài âm hoặc dương.

## Ví dụ đang chạy

<!-- wave-example: random-api -->
```wave
import("std::random::fill")::{
    RandomFillResult, random_fill
};

fun main() -> i32 {
    var bytes: array<u8, 16>;
    var result: RandomFillResult = random_fill(&bytes[0], 16);
    if (!result.ok) {
        println("random error={}", result.error);
        return 1;
    }

    println("filled={}", result.written);
    return 0;
}
```

Kết quả thực hiện:

```text
filled=16
```

Lưu nó dưới dạng `main.wave` và chạy nó. Nội dung byte mỗi lần khác nhau nên không có giá trị cụ thể nào được mong đợi. Nếu không thành công, hãy kiểm tra nguyên nhân bằng result.error. Thay vì xuất ra các byte ngẫu nhiên theo nghĩa đen, hãy sử dụng mã hóa riêng nếu cần.

## Nếu yêu cầu của bạn không chính xác

Yêu cầu 0 byte thành công vì không cần phải viết gì. Việc truyền null với độ dài dương không thành công do không có bộ đệm đích. Chương trình sau đây so sánh các trường hợp này mà không cấp phát bộ nhớ.

<!-- wave-example: book-random-boundaries -->
```wave
import("std::random::fill")::{
    RandomFillResult,
    random_fill
};

fun main() -> i32 {
    var empty: RandomFillResult = random_fill(null, 0);

    if (!empty.ok || empty.written != 0) {
        return 1;
    }

    println("empty request succeeded");

    var invalid: RandomFillResult = random_fill(null, 16);

    if (invalid.ok) {
        return 2;
    }

    println("missing buffer rejected");
    return 0;
}
```

Kết quả thực hiện:

```text
empty request succeeded
missing buffer rejected
```

## Xử lý bộ đệm được lấp đầy một phần

Nếu 16 byte được yêu cầu nhưng không thành công và written=8 được trả về thì chỉ 8 byte đầu tiên được lấp đầy. Nếu nhiệm vụ là tạo mã định danh 16 byte thì đó không phải là mã định danh thành công nên chúng tôi loại bỏ toàn bộ kết quả và báo cáo lỗi. Bạn không nên điền 8 byte còn lại bằng 0 rồi coi như thành công.

Ánh xạ các byte ngẫu nhiên vào một phạm vi số nguyên đòi hỏi phải cẩn thận. Việc áp dụng `% 10` cho các giá trị u8 được phân bố đồng đều làm cho 0–5 có nhiều khả năng hơn 6–9, vì 256 không chia hết cho 10. Để loại bỏ sai lệch này, hãy loại bỏ các giá trị 250–255, vẽ lại và chỉ áp dụng thao tác còn lại cho các giá trị được chấp nhận.

Việc lưu trữ và thời gian tồn tại của các byte ngẫu nhiên được quản lý bởi người gọi. Khi sử dụng mảng, nó được xử lý trong phạm vi của mảng và khi sử dụng bộ nhớ động, nó sẽ được giải phóng sau khi sử dụng. Thay vào đó, cấu trúc trả về không sở hữu bộ đệm.
