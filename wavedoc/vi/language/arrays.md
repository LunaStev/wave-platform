---
translation_set_id: learn-arrays-strings
path: language/arrays
locale: vi
group: language
group_order: 2
order: 6
title: 6. Mảng và phép lặp
summary: Tìm hiểu các mảng có kích thước cố định, lập chỉ mục, truyền tải, sao chép và tìm kiếm.
---

## Nhiều giá trị cùng loại

Nếu bạn tạo ba điểm riêng biệt dưới dạng score1, score2 và score3 thì cả phần khai báo và phép tính đều phải được thay đổi khi số thay đổi. Mảng nhóm một số phần tử cùng loại lại với nhau. Bằng cách sử dụng vòng lặp, bạn có thể áp dụng các quy tắc tương tự cho từng phần tử.

Chương này bao gồm việc tạo, lập chỉ mục, sửa đổi, lặp lại, tìm kiếm và tổng hợp các mảng. Chuỗi cũng hỗ trợ lập chỉ mục, nhưng ý nghĩa của nó khác nhau, vì vậy chuỗi sẽ được đề cập trong chương tiếp theo.

## Nhập độ dài vào loại

<!-- wave-example: book-array-create -->
```wave
fun main() {
    var scores: array<i32, 3> = [70, 80, 90];

    println("first={}", scores[0]);
    println("second={}", scores[1]);
    println("last={}", scores[2]);
}
```

Kết quả thực hiện:

```text
first=70
second=80
last=90
```

i32 trong `array<i32, 3>` là loại phần tử và 3 là số phần tử. Những gì chúng tôi lưu trữ là 3 số nguyên. Điều đó không có nghĩa là số byte là 3. Số phần tử trong một mảng phải khớp với độ dài đã khai báo.

Các chỉ mục bắt đầu từ 0. Phần tử đầu tiên là 0, phần tử cuối cùng là length-1. scores[3] là quyền truy cập ngoài phạm vi, không phải là thành phần thứ ba.

## thay đổi phần tử

<!-- wave-example: book-array-update -->
```wave
fun main() {
    var scores: array<i32, 3> = [70, 80, 90];

    scores[1] = 85;
    scores[0] += 5;

    println("{} {} {}", scores[0], scores[1], scores[2]);
}
```

Kết quả thực hiện:

```text
75 85 90
```

Thay đổi việc lưu trữ các phần tử cụ thể mà không tạo lại toàn bộ mảng. Biểu thức chỉ mục cũng có thể là kết quả của phép tính nhưng bạn phải đảm bảo rằng giá trị nằm trong một phạm vi. Khi sử dụng đầu vào bên ngoài làm chỉ mục, cả số âm và giới hạn trên đều được kiểm tra.

## Lặp lại trên một mảng

<!-- wave-example: book-array-sum -->
```wave
fun main() {
    var scores: array<i32, 4> = [60, 70, 80, 90];
    var total: i32 = 0;

    for (var index: i32 = 0; index < 4; index += 1) {
        total += scores[index];
    }

    println("total={} average={}", total, total / 4);
}
```

Kết quả thực hiện:

```text
total=300 average=75
```

Mỗi lần lặp lại đọc một phần tử ở một chỉ mục khác nhau. Biến tổng phải được khởi tạo bên ngoài vòng lặp. Việc khởi tạo nó thành 0 mỗi lần trong thân vòng lặp sẽ tạo ra kết quả không chính xác, chẳng hạn như chỉ để lại phần tử cuối cùng.

Phép chia số nguyên được sử dụng để tính trung bình sẽ loại bỏ phần phân số. Đối với giá trị trung bình của dấu phẩy động, hãy chuyển đổi tổng trước khi chia. Đối với các mảng lớn hơn hoặc giá trị lớn hơn, hãy đảm bảo rằng loại bộ tích lũy có thể biểu thị tổng.

## Chỉ tổng hợp một số phần tử

Việc lọc có thể đạt được bằng cách kết hợp các câu lệnh có điều kiện và truyền tải. Ở đây chúng tôi đếm số phần tử có điểm từ 80 trở lên.

<!-- wave-example: book-array-filter -->
```wave
fun main() {
    var scores: array<i32, 5> = [60, 80, 90, 75, 100];
    var passed: i32 = 0;

    for (var index: i32 = 0; index < 5; index += 1) {
        if (scores[index] >= 80) {
            passed += 1;
        }
    }

    println("passed={}", passed);
}
```

Kết quả thực hiện:

```text
passed=3
```

Giá trị chỉ mục và phần tử phải được tách biệt. Kiểm tra `index >= 80` so sánh vị trí chứ không phải điểm số. Cả hai đều có thể là i32, vì vậy rất khó tìm ra lỗi ngữ nghĩa này nếu chỉ dựa vào loại.

## Tìm vị trí trận đấu đầu tiên

Đầu tiên, quyết định cách hiển thị kết quả không tìm thấy. Trong ví dụ này, các chỉ số hợp lệ là từ 0 đến 4, vì vậy chúng tôi sử dụng -1 làm điểm đánh dấu lỗi.

<!-- wave-example: book-array-find -->
```wave
fun main() {
    var values: array<i32, 5> = [8, 3, 8, 1, 5];
    var target: i32 = 8;
    var found: i32 = -1;

    for (var index: i32 = 0; index < 5; index += 1) {
        if (values[index] == target) {
            found = index;
            break;
        }
    }

    if (found >= 0) {
        println("found at {}", found);
    } else {
        println("not found");
    }
}
```

Kết quả thực hiện:

```text
found at 0
```

Vị trí đầu tiên 0 cũng là kết quả bình thường. Nếu bạn kiểm tra thành công với `found > 0`, bạn sẽ bị nhầm là không tìm thấy phần tử đầu tiên. Vì lý do tương tự, việc đánh giá thành công bằng cách sử dụng bool thay vì cast là sai.

Nếu bạn xóa break, các kết quả phù hợp tiếp theo sẽ ghi đè found, mang lại cho bạn vị trí của kết quả phù hợp cuối cùng. Vì một câu lệnh có thể thay đổi hợp đồng của một hàm nên phần mô tả “tìm kiếm” cũng phải được viết cụ thể về việc nó ở vị trí đầu tiên hay cuối cùng.

## Sao chép các phần tử mảng

Để sao chép các giá trị của một mảng sang một không gian lưu trữ khác, bạn có thể đọc và gán chúng theo từng phần tử. Ngay cả khi bạn thay đổi một phần tử số nguyên sau khi sao chép, phần tử số nguyên còn lại vẫn không thay đổi.

<!-- wave-example: book-array-copy -->
```wave
fun main() {
    var original: array<i32, 3> = [1, 2, 3];
    var copied: array<i32, 3>;

    for (var index: i32 = 0; index < 3; index += 1) {
        copied[index] = original[index];
    }

    copied[0] = 99;

    println("original={}", original[0]);
    println("copied={}", copied[0]);
}
```

Kết quả thực hiện:

```text
original=1
copied=99
```

Nếu các phần tử là con trỏ, việc sao chép chúng sẽ sao chép địa chỉ của chúng. Nó không trùng lặp bộ nhớ riêng biệt mà chúng trỏ tới. Sự khác biệt này rất quan trọng khi quản lý quyền sở hữu.

## Khởi tạo và phạm vi hiệu quả

Nó không cho rằng tất cả các phần tử của một mảng được khai báo không có giá trị ban đầu đều có thể đọc được. Nếu chỉ ghi một phần số thì số khởi tạo thực tế phải được quản lý riêng. Đây cũng là lý do tại sao độ dài được hàm đọc thư viện trả về có thể nhỏ hơn toàn bộ dung lượng bộ đệm.

Vì độ dài mảng được chứa trong kiểu nên nó không tăng tùy ý trong quá trình thực thi. Danh sách byte tăng kích thước sử dụng bộ lưu trữ động, chẳng hạn như `Buffer`. Việc thay đổi độ dài của một mảng yêu cầu xem xét loại, giá trị ban đầu, giới hạn trên của đường truyền và các phép tính phụ thuộc vào độ dài đó.

## Bài tập: Mức tối đa và Vị trí

Tìm giá trị lớn nhất và vị trí mà nó xuất hiện lần đầu trong mảng `[4, 9, 2, 9, 1]`. Nếu đó là một hàm thông thường trong đó tất cả các phần tử có thể âm thì giá trị tối đa không được khởi tạo thành 0.

### Lời giải đầy đủ

<!-- wave-example: book-array-solution -->
```wave
fun main() {
    var values: array<i32, 5> = [4, 9, 2, 9, 1];
    var maximum: i32 = values[0];
    var position: i32 = 0;

    for (var index: i32 = 1; index < 5; index += 1) {
        if (values[index] > maximum) {
            maximum = values[index];
            position = index;
        }
    }

    println("max={} first={}", maximum, position);
}
```

Kết quả thực hiện:

```text
max=9 first=1
```

Lấy phần tử đầu tiên làm tham chiếu ban đầu và so sánh với phần tử thứ hai. Kể từ `>`, vị trí không thay đổi ngay cả khi giá trị tối đa tương tự xuất hiện lại. Thay đổi nó thành `>=` để trở thành vị trí cuối cùng. Bất kỳ giao diện nào có độ dài bằng 0 đều phải xử lý đầu vào trống trước khi đọc phần tử đầu tiên.
