---
translation_set_id: stdlib-task
path: stdlib/task
locale: vi
group: stdlib
group_order: 1
order: 13
title: task: Chạy Future và dọn dẹp tài nguyên
summary: Mô tả mức tiêu thụ, thực thi và hủy bỏ sau khi dọn dẹp các hoạt động không đồng bộ.
---

## Cơ bản API

Đã nhập dưới dạng `import("std::task" as task);`.

```text
block_on<T>(future: Future<T>) -> T
spawn<T>(future: Future<T>) -> Future<T>
cancel<T>(future: Future<T>) -> bool
yield_now() -> Future<void>
sleep_ms(milliseconds: i64) -> Future<void>
cancel_and_join<T>(future: Future<T>) -> Future<void>
shutdown()
```

`task::block_on(pending)` chạy cho đến khi Future hoàn thành và trả về kết quả. Loại kết quả được xác định từ Future đã đạt.

## quy luật cuộc sống

Future là loại xử lý tiêu dùng một lần. Tôi không nghĩ việc sao chép các giá trị là hai thao tác độc lập. Bạn cũng có trách nhiệm chờ kết quả hoặc hủy/sắp xếp các nhiệm vụ đã lên lịch với `spawn`. Future đã được tiêu thụ cùng với `await` và `block_on` sẽ không được tiêu thụ nữa.

hủy bỏ yêu cầu hủy bỏ; bản thân nó không đảm bảo rằng việc dọn dẹp đã hoàn tất. Đợi hoàn thành theo yêu cầu trước khi giải phóng tài nguyên. Không giải phóng bộ nhớ được mượn bởi I/O không đồng bộ trong khi tác vụ vẫn có thể truy cập vào bộ nhớ đó. Hãy gọi `shutdown` sau khi nhiệm vụ đã hoàn thành hoặc việc dọn dẹp đã hoàn tất.

## thực hiện hợp tác

Tính toán dài và lệnh gọi blocking đồng bộ có thể làm chậm tiến độ chung của người thực thi. Chờ không đồng bộ với yield mang lại cơ hội thực thi. Việc chặn I/O không trở nên không đồng bộ chỉ vì nó nằm trong hàm async.

Hãy xem trình tự thực hiện và dọn dẹp trong chương trình đầy đủ tại [Giới thiệu về mã không đồng bộ](/docs/vi/language/async-and-never).

## Năng suất và chờ hoàn thành

Chương trình sau đây thực thi ở giữa một thao tác và trả về một kết quả. Vì yield không phải là kết thúc hàm nên mã sau await được thực thi liên tục.

<!-- wave-example: book-task-yield -->
```wave
import("std::task" as task);

async fun calculate() -> i32 {
    println("started");
    await task::yield_now();
    println("resumed");
    return 42;
}

fun main() -> i32 {
    var pending: Future<i32> = calculate();
    var result: i32 = task::block_on(pending);

    println("result={}", result);
    task::shutdown();
    return 0;
}
```

Kết quả thực hiện:

```text
started
resumed
result=42
```

Không có thao tác nào khác trong ví dụ này nên thứ tự đầu ra là nhất quán. Trong một chương trình có nhiều tác vụ spawn, các tác vụ khác có thể được tiến hành tại điểm yield nên không phụ thuộc vào thứ tự đầu ra của các tác vụ khác nhau.

## Thứ tự tổ chức nhiệm vụ

1. Chuẩn bị không gian lưu trữ và tài nguyên bạn cần cho công việc của mình.
2. Tạo Future và chạy dưới dạng await, block_on hoặc spawn.
3. Nếu bạn cần kết quả, hãy đợi cho đến khi hoàn thành.
4. Nếu bạn hủy một tác vụ đang chạy, hãy đợi tác vụ đó được dọn sạch.
5. Dọn dẹp bộ đệm, tập tin và ổ cắm mà tác vụ mượn.
6. Không còn việc gì nữa, hãy gọi shutdown.

Future Rời khỏi phạm vi của một biến là một chuyện và việc dọn dẹp công việc một cách an toàn lại là một chuyện khác. Đặc biệt, nếu bạn chuyển địa chỉ của mảng hàm cục bộ cho một tác vụ async, thì tác vụ đó phải hoàn thành việc sử dụng mảng đó trước khi hàm trả về.

## async Hàm chia và hàm tổng quát

Các phép tính thuần túy có thể được tách thành các hàm thông thường. Đính kèm async vào hàm cần thể hiện sự chờ đợi và chờ hoàn thành với await trong hàm đó. Chỉ gói một hàm đồng bộ mất nhiều thời gian, chẳng hạn như đọc tệp, bằng hàm async sẽ không cho các tác vụ khác cơ hội thực thi.

Bạn có thể so sánh khi tạo và chạy Future trong [học tập không đồng bộ](/docs/vi/language/async-and-never).
