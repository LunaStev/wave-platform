---
translation_set_id: language-async-and-never
path: language/async-and-never
locale: vi
group: language
group_order: 2
order: 13
title: 13. Hàm bất đồng bộ và Future
summary: Tìm hiểu vai trò của việc trì hoãn thực thi Future, await và block_on.
---

## Thể hiện nhiệm vụ chờ đợi

Đối với các hoạt động chờ đợi như tệp, ổ cắm và bộ hẹn giờ, cần phân biệt giữa việc tiếp tục tính toán và chờ hoàn thành. Hàm không đồng bộ biểu thị kết quả cần hoàn thành là Future. Việc thêm async không tự động tạo chuỗi mới hoặc thay đổi tất cả các lệnh gọi đồng bộ thành không đồng bộ.

Đọc chương này sau Hàm, Con trỏ và Xử lý lỗi. Ví dụ này là một chương trình gốc sử dụng trình khởi chạy `std::task` và chạy từng tệp dưới dạng `wavec run main.wave`.

## Tạo Future và nhận kết quả

<!-- wave-example: book-async-basic -->
```wave
import("std::task" as task);

async fun calculate(value: i64) -> i64 {
    return value * 2;
}

fun main() {
    var pending: Future<i64> = calculate(21);
    var result: i64 = task::block_on(pending);

    task::shutdown();
    println("{}", result);
}
```

Kết quả thực hiện:

```text
42
```

i64 ghi trong khai báo calculate là giá trị thu được sau khi hoàn thành. Kết quả của cuộc gọi là Future<i64>. Trong main bình thường, hãy chạy Future với block_on và nhận kết quả hoàn thành.

Tắt lệnh gọi sau khi dọn dẹp tất cả các tác vụ để giải phóng tài nguyên của người thực thi. Không giải phóng bộ đệm khi tác vụ vẫn đang sử dụng nó và không bỏ qua các tác vụ chưa hoàn thành.

## Việc gọi và thực thi nội dung là khác nhau

Các hàm không đồng bộ thực thi một cách lười biếng. Nó phải được phân biệt với một hàm thông thường thực thi phần thân của nó ngay sau khi gọi nó.

<!-- wave-example: book-async-lazy -->
```wave
import("std::task" as task);

static entered: i32 = 0;

async fun work() -> i32 {
    entered += 1;
    return 7;
}

fun main() {
    var pending: Future<i32> = work();

    println("before={}", entered);

    var result: i32 = task::block_on(pending);

    task::shutdown();
    println("after={} result={}", entered, result);
}
```

Kết quả thực hiện:

```text
before=0
after=1 result=7
```

Khi Future được tạo, entered vẫn bằng 0. Sau khi thực thi điều khiển, nội dung được thực thi và trở thành 1. Đây là lý do tại sao bạn không nên coi nhiệm vụ là đã hoàn thành chỉ vì Future được lưu trữ trong một biến.

## Chờ đợi bên trong một hàm không đồng bộ

Bên trong hàm async, await chờ hoàn thành một Future khác. Kết quả của biểu thức await là giá trị hoàn thành.

<!-- wave-example: book-async-nested -->
```wave
import("std::task" as task);

async fun twice(value: i32) -> i32 {
    await task::yield_now();
    return value * 2;
}

async fun process() -> i32 {
    var value: i32 = await twice(20);
    return value + 2;
}

fun main() {
    var answer: i32 = task::block_on(process());

    task::shutdown();
    println("{}", answer);
}
```

Kết quả thực hiện:

```text
42
```

process đợi Future của twice rồi thêm 2 vào giá trị hoàn thành. Bạn có thể sử dụng giá trị tương tự như kết quả của lệnh gọi đến một hàm thông thường, nhưng trong khi chờ đợi, bạn có thể chuyển cơ hội thực thi sang một tác vụ khác.

yield_now mang lại cơ hội thực hiện hợp tác. Việc không đưa ra dù chỉ một lần trong một vòng tính toán dài có thể làm chậm các tác vụ khác. Không đồng bộ CPU không phải là thiết bị để tự động song song hóa và phân phối các phép tính.

## Lên lịch nhiều nhiệm vụ

Bạn có thể lên lịch các nhiệm vụ bằng spawn và chờ kết quả tương ứng. Đây là ví dụ về kiểm tra kết quả cuối cùng mà không dựa vào thứ tự đầu ra trung gian của hai thao tác.

<!-- wave-example: book-async-spawn -->
```wave
import("std::task" as task);

async fun compute(value: i32) -> i32 {
    await task::yield_now();
    return value * 2;
}

async fun combine() -> i32 {
    var first: Future<i32> = task::spawn(compute(10));
    var second: Future<i32> = task::spawn(compute(20));

    var left: i32 = await first;
    var right: i32 = await second;

    return left + right;
}

fun main() {
    var total: i32 = task::block_on(combine());

    task::shutdown();
    println("{}", total);
}
```

Kết quả thực hiện:

```text
60
```

Mỗi nhiệm vụ bạn lên lịch đều có một nơi để chờ kết quả. Thay vì tạo một nhiệm vụ và quên mất phần xử lý của nó, bạn cần quyết định ai sẽ kiểm tra việc hoàn thành nhiệm vụ đó. await Đừng coi trình tự và thứ tự thực hiện các hoạt động nội bộ là giống nhau.

## Tiêu thụ Future một lần

Future được coi là một tay cầm tiêu thụ duy nhất. Việc sao chép cùng một Future không khiến nó phải chờ như hai tác vụ khác nhau. Không block_on hoặc await lặp lại đối với Future đã được hoàn thành.

Nếu bạn cần cùng một kết quả ở nhiều nơi, thay vì sử dụng Future nhiều lần, hãy lưu trữ giá trị đã hoàn thành và chuyển giá trị đó theo quy tắc sao chép/chia sẻ cho giá trị đó. Bạn cũng nên kiểm tra xem có con trỏ hoặc tài nguyên được sở hữu trong giá trị hay không.

## Sự khác biệt giữa hẹn giờ và chờ đồng bộ

async Khi chờ trong một hàm, bạn có thể sử dụng `await task::sleep_ms(...)`. Việc gọi đồng bộ sleep chặn luồng thực thi hiện tại, điều này cũng có thể ảnh hưởng đến tiến độ của các tác vụ khác trong trình thực thi.

Tôi không mong đợi thời gian chờ đợi chính xác bằng số mili giây được yêu cầu. Tùy thuộc vào lịch trình và các nhiệm vụ khác của bạn, bạn có thể thức dậy muộn. Khi triển khai thời gian chờ, chúng tôi sử dụng clock và deadline để đo thời gian đã trôi qua, thay vì phải đợi toàn bộ thời gian ban đầu mỗi lần.

## Thời gian tồn tại của bộ đệm và hủy bỏ

Bộ đệm được chuyển tới I/O không đồng bộ phải vẫn hợp lệ ngay cả khi chức năng gọi bị tạm dừng. Việc giải phóng hoặc phân bổ lại nó trước khi hoàn thành hoặc hủy bỏ việc dọn dẹp có thể khiến hoạt động có địa chỉ không hợp lệ.

Không thể giả định rằng yêu cầu hủy và hoàn thành nhiệm vụ là cùng một lúc. Kiểm tra kết quả của loạt cancel API và chờ hoàn thành cần thiết trước khi phát hành tài nguyên. Vui lòng đọc [Xem task](/docs/vi/stdlib/task) để biết quy tắc gọi chi tiết.

## sự hiểu lầm phổ biến

|nghĩ|thực sự kiểm tra|
| --- | --- |
|async Tôi đã gọi và mọi chuyện đã kết thúc.|Bạn có thực sự đã chạy và hoàn thành Future không?|
|async Tất cả lệnh gọi trong hàm đều không đồng bộ|Cái được gọi là API đồng bộ hay không đồng bộ?|
|Future Sao chép trùng lặp một nhiệm vụ|Bạn có đang sử dụng cùng một tay cầm nhiều lần không?|
|Vì bạn đã hủy nó nên bạn có thể giải phóng bộ đệm ngay lập tức.|Bạn đã sắp xếp xong công việc sau khi hủy chưa?|
|Thứ tự đầu ra trung gian luôn cố định|Nó có rõ ràng chỉ chờ thứ tự cần thiết cho kết quả không?|

## Bài tập và lời giải đầy đủ

Tạo một đường dẫn chờ ba hàm không đồng bộ liên tiếp. Trả về kết quả của việc nhân đôi và cộng 5.

<!-- wave-example: book-async-solution -->
```wave
import("std::task" as task);

async fun read_value() -> i32 {
    return 10;
}

async fun transform(value: i32) -> i32 {
    await task::yield_now();
    return value * 2;
}

async fun pipeline() -> i32 {
    var initial: i32 = await read_value();
    var changed: i32 = await transform(initial);

    return changed + 5;
}

fun main() {
    var result: i32 = task::block_on(pipeline());

    task::shutdown();
    println("{}", result);
}
```

Kết quả thực hiện:

```text
25
```

Ví dụ này cố tình phụ thuộc tuần tự. transform yêu cầu kết quả của read_value nên chỉ cần thực hiện tất cả spawn sẽ không khiến mối quan hệ mất đi. Điểm khởi đầu của thiết kế không đồng bộ là sự phân biệt giữa các hoạt động độc lập và các hoạt động yêu cầu kết quả.

Sau khi bạn đã hoàn thành những điều cơ bản, hãy kết nối với các tài nguyên thực bên ngoài bằng [Luyện đọc file](/docs/vi/practice/file-reader) và [TCP Luyện tập](/docs/vi/practice/tcp-client).

## void và never

Một hàm thông thường bỏ qua kiểu trả về có thể quay trở lại điểm gọi mà không có giá trị. Loại never được viết là `!`, nghĩa là nó không quay lại điểm gọi một cách bình thường. Một ví dụ điển hình là chức năng kết thúc quá trình.

Ví dụ minh họa việc khai báo:

```text
fun log(message: str)              // 작업 후 돌아옴, 결과값 없음
fun stop(code: i32) -> !           // 정상 반환하지 않음
async fun work() -> i64            // 호출 결과는 Future<i64>
```

Không có nỗ lực nào để biến never thành một giá trị được lưu trữ chung. Không viết các hàm được khai báo là không quay trở lại để có đường dẫn trả về bình thường. Nếu cần phải dọn sạch tài nguyên trước khi thoát, thì người gọi phải thực hiện việc đó trước.

## Ví dụ đầy đủ về hàm không trả về

Nếu bạn lưu nó dưới dạng main.wave và chạy nó, nó sẽ kết thúc bằng mã thoát 0 mà không có đầu ra. stop không trả về cho người gọi nên được khai báo là `-> !`.

<!-- wave-example: never-function -->
```wave
import("std::process::core")::{
    proc_exit
};

fun stop() -> ! {
    proc_exit(0);
}

fun main() {
    stop();
}
```
