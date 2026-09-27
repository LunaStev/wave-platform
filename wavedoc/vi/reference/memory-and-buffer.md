---
translation_set_id: memory-buffer
path: reference/memory-and-buffer
locale: vi
group: stdlib
group_order: 1
order: 4
title: mem: Phân bổ, tái phân bổ và bố cục
summary: Mô tả kích thước tính bằng byte, lỗi phân bổ, ranh giới phân bổ lại và trách nhiệm giải phóng.
---

## Phân bổ và phân bổ

```text
std::mem::alloc
mem_alloc(size: i64) -> ptr<u8>
mem_alloc_zeroed(size: i64) -> ptr<u8>
mem_free(p: ptr<u8>, size: i64) -> i64
mem_realloc(old_ptr: ptr<u8>, old_size: i64, new_size: i64) -> ptr<u8>
```

Kích thước tính bằng byte. Phân bổ có kích thước 0 trở xuống trả về null. Nếu việc phân bổ cũng không thành công đối với kích thước dương, nó có thể là null. Đừng giả sử nội dung ban đầu của `mem_alloc` mà hãy sử dụng `mem_alloc_zeroed` nếu không cần khởi tạo.

Người gọi sở hữu mỗi lần phân bổ thành công và phải vượt qua kích thước ban đầu của nó khi giải phóng nó. `mem_free(null, size)` trả về 0. Một con trỏ không phải null được ghép với kích thước không dương là một lỗi. Không bao giờ truy cập hoặc giải phóng một phân bổ sau khi nó đã được giải phóng.

## Phân loại trong trường hợp tái phân bổ

|yêu cầu|hành động|
| --- | --- |
|kích thước mới là tích cực và thành công|`min(old_size, new_size)` Sao chép byte và phân bổ byte trước đó|
|Phân bổ kích thước dương mới không thành công|null Trả lại, duy trì phân bổ hiện có|
|`old_ptr == null`, kích thước mới tích cực|cư xử như một nhiệm vụ mới|
| `new_size == 0` |Cố gắng giải phóng phân bổ hợp lệ trước đó và trả về null|
|old_size=0 cho kích thước âm hoặc con trỏ hiện có|null Trở về|

Kết quả null từ việc phân bổ lại kích thước 0 không xác định liệu việc giải phóng có thành công hay không. Gọi trực tiếp `mem_free` nếu bạn cần trạng thái của nó. Sau khi phát triển phân bổ, hãy tự khởi tạo vùng mới được thêm vào.

## Ví dụ về việc bảo toàn một con trỏ hiện có

Dưới đây là trường hợp kích thước mới là dương. Lưu nó dưới dạng `main.wave` và chạy nó.

<!-- wave-example: reallocation -->
```wave
import("std::mem::alloc")::{
    mem_alloc, mem_realloc, mem_free
};

fun main() -> i32 {
    var data: ptr<u8> = mem_alloc(4);
    if (data == null) {
        return 1;
    }
    deref data[0] = 7;
    var grown: ptr<u8> = mem_realloc(data, 4, 8);
    if (grown == null) {
        mem_free(data, 4);
        return 2;
    }
    data = grown;
    println("{}", deref data[0]);
    if (mem_free(data, 8) < 0) {
        return 3;
    }

    return 0;
}
```

Kết quả thực hiện:

```text
7
```

Nếu bạn ghi đè lên nó bằng `data = mem_realloc(...)` trước khi xác nhận lỗi, bạn có thể mất địa chỉ hiện có. Nếu việc tái phân bổ thành công, địa chỉ cũ và các con trỏ trỏ tới địa chỉ đó sẽ không được sử dụng.

## Kích thước và sự liên kết của loại mục tiêu

```text
std::mem::layout
size_of<T>() -> u64
align_of<T>() -> u64
```

Cả hai giá trị đều là bố cục của mục tiêu biên dịch, không phải máy tính đang chạy nó. `size_of` bao gồm phần đệm đuôi và không tạo ra hoặc đánh giá giá trị. Khi nhân số phần tử với kích thước của chúng, chúng tôi kiểm tra tình trạng tràn. Bạn có thể sử dụng `mem_size_mul_checked` và `mem_size_add_checked` của `std::mem::ops`.

`mem_copy` được sử dụng để sao chép một phạm vi không trùng nhau và `mem_move` được sử dụng để sao chép một phạm vi có thể trùng lặp. Cả hai đều không thể xác định độ dài phân bổ thực tế chỉ từ con trỏ, vì vậy người gọi phải đảm bảo giới hạn. Bạn có thể quản lý danh sách các byte có kích thước khác nhau bằng [Buffer](/docs/vi/stdlib/buffer).

## Ví dụ về truy vấn bố cục

Lưu nó dưới dạng main.wave và chạy nó. Kích thước và căn chỉnh của mục tiêu được đề cập trong tài liệu, i32, mỗi mục có 4 byte, vì vậy `4 4` là đầu ra. Kiểm tra giá trị của các loại khác, đặc biệt là cấu trúc và con trỏ, theo mục tiêu.

<!-- wave-example: layout-api -->
```wave
import("std::mem::layout")::{
    size_of, align_of
};

fun main() {
    println("{} {}", size_of<i32>(), align_of<i32>());
}
```

## Kiểm tra tính toán kích thước tràn

Bạn cần kiểm tra xem `count * element_size` có hợp lệ hay không trước khi chuyển nó sang hàm gán. Nếu bạn phân bổ một khoảng trống nhỏ với giá trị tràn và viết nhiều bằng số ban đầu, nó sẽ vượt quá giới hạn.

<!-- wave-example: book-memory-size-check -->
```wave
import("std::mem::ops")::{mem_size_add_checked, mem_size_mul_checked};

fun main() {
    var result: i64 = 99;

    if (mem_size_mul_checked(3, 4, &result) == 0) {
        println("bytes={}", result);
    }

    if (mem_size_add_checked(9223372036854775807, 1, &result) < 0) {
        println("overflow rejected");
    }
}
```

Kết quả thực hiện:

```text
bytes=12
overflow rejected
```

Kết quả lỗi không được ghi vào kích thước phân bổ. Bạn không thể phát hiện tất cả các lần tràn chỉ bằng cách tính toán chúng bằng số học thông thường và xem kết quả có âm hay không. Để tính toán kích thước cần kiểm tra, hãy sử dụng hàm checked ngay từ đầu.

## Đối với các bản sao chồng chéo, mem_move

Khi di chuyển một phần của cùng một mảng về phía sau, vùng đầu vào và đầu ra sẽ chồng lên nhau. Không chuyển các phạm vi chồng chéo sang mem_copy, hãy sử dụng mem_move.

<!-- wave-example: book-memory-overlap -->
```wave
import("std::mem::ops")::{mem_move};

fun main() {
    var data: array<u8, 5> = [1, 2, 3, 4, 5];

    mem_move(&data[1], &data[0], 4);

    for (var index: i32 = 0; index < 5; index += 1) {
        println("{}", data[index]);
    }
}
```

Kết quả thực hiện:

```text
1
1
2
3
4
```

Bốn byte đầu tiên ban đầu di chuyển sang phải một vị trí. Sao chép chúng về phía trước theo cách thủ công có thể đọc các giá trị đã bị ghi đè, vô tình tạo ra tất cả các giá trị đó. API nhận biết chồng chéo sẽ xử lý hướng sao chép cho bạn.

## Viết hàm chuyển quyền sở hữu

Đối với các hàm trả về bộ nhớ, tốt nhất nên cung cấp địa chỉ trả về khi thành công và kích thước cần thiết để giải phóng nó. Nếu người gọi phải đoán kích thước thì rất dễ xảy ra lỗi phát hành. Nếu một hàm trả về một địa chỉ đã mượn thì nó sẽ không được giải phóng bởi người gọi và mô tả thời gian tồn tại của địa chỉ gốc.

Trên khắp các ranh giới chức năng, bạn sẽ có thể theo dõi `allocator → owner → deallocation`. Tên hoặc kiểu của biến con trỏ không tự động xác định quyền sở hữu.
