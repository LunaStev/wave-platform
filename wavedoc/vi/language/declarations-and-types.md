---
translation_set_id: types
path: language/declarations-and-types
locale: vi
group: language
group_order: 2
order: 2
title: 2. Biến, loại và phạm vi
summary: Tìm hiểu các biến cục bộ, phạm vi số nguyên, bool và phạm vi.
---

## Xử lý giá trị theo tên

Nếu bạn viết giá trực tiếp ở nhiều nơi, bạn sẽ phải tìm thấy tất cả các giá đó khi thay đổi giá. Biến là không gian lưu trữ đặt tên cho các giá trị và cho phép bạn đọc và thay đổi chúng bằng cách sử dụng các tên đó. Trong chương này, bạn sẽ tìm hiểu về phạm vi khai báo, phép gán, kiểu và phạm vi của tên hiển thị trong các khối.

Các chương trình bên dưới là một phần của một main.wave riêng biệt. Lưu một ví dụ, chạy nó dưới dạng `wavec run main.wave` và thay thế bằng ví dụ tiếp theo.

## Khai báo và khởi tạo

<!-- wave-example: book-variable-first -->
```wave
fun main() {
    var price: i32 = 1200;
    var quantity: i32 = 3;
    println("price={}", price);
    println("quantity={}", quantity);
    println("total={}", price * quantity);
}
```

Kết quả thực hiện:

```text
price=1200
quantity=3
total=3600
```

Đọc lời tuyên bố thành bốn phần.

|một phần|ví dụ này|vai trò|
| --- | --- | --- |
|từ khóa khai báo| var |Tạo biến cục bộ|
|tên| price |Mã định danh sẽ được sử dụng sau này|
|loại| i32 |Loại và phạm vi giá trị cần lưu trữ|
|giá trị ban đầu| 1200 |Giá trị đầu tiên được lưu trữ|

Dấu hai chấm trước loại và dấu bằng trước giá trị ban đầu có vai trò khác nhau. Làm cho tên của bạn có ý nghĩa. Trong ví dụ này, price là đơn giá và quantity là số lượng. Ngay cả khi cùng một i32 được sử dụng thay thế cho nhau, vẫn có thể thực hiện các phép tính không chính xác mà không mắc lỗi ngữ pháp.

## Sự phân công không phải là một công thức duy trì các mối quan hệ.

Nếu bạn lưu kết quả tính toán vào một biến thì giá trị tại thời điểm đó sẽ được nhập vào. Nó không nhớ các phép tính và tự động đánh giá lại chúng sau này.

<!-- wave-example: book-assignment-snapshot -->
```wave
fun main() {
    var price: i32 = 1200;
    var quantity: i32 = 3;
    var total: i32 = price * quantity;
    quantity = 5;
    println("before recalculation={}", total);
    total = price * quantity;
    println("after recalculation={}", total);
}
```

Kết quả thực hiện:

```text
before recalculation=3600
after recalculation=6000
```

`quantity = 5` ghi một giá trị mới vào một biến đã tồn tại. Nó phải được phân biệt với việc khai báo lại như `var quantity`. `total` cũng là 3600 trước khi thay lại. Nếu một chương trình phải duy trì mối quan hệ giữa nhiều biến, thì nó phải được viết để thực hiện các phép tính khi mối quan hệ thay đổi.

## Tính giá trị tiếp theo từ giá trị trước đó

Trước tiên hãy tính vế phải của câu lệnh gán và ghi kết quả vào vùng lưu trữ bên trái. Không giống như các phương trình trong toán học, `count = count + 1` là một bản cập nhật hợp lệ.

<!-- wave-example: book-update-variable -->
```wave
fun main() {
    var count: i32 = 0;
    count = count + 1;
    count += 2;
    count *= 3;
    println("{}", count);
}
```

Kết quả thực hiện:

```text
9
```

Tiến trình của các giá trị là 0 → 1 → 3 → 9. `+=`, `*=` thể hiện tính toán và lưu trữ cùng nhau. Nếu bạn viết ra thứ tự tính toán, bạn có thể biết được suy nghĩ của bạn khác ở giai đoạn nào khi kết quả khác với những gì bạn mong đợi.

## Chiều rộng và dấu của các loại số nguyên

Nếu bắt đầu bằng `i`, nó sẽ trở thành signed và nếu bắt đầu bằng `u`, nó sẽ trở thành unsigned. Số cuối cùng là số bit. Khi số lượng bit tăng lên, phạm vi có thể biểu diễn tăng lên và không gian lưu trữ cũng tăng lên.

|loại|giá trị tối thiểu|giá trị tối đa|Ví dụ sử dụng|
| --- | --- | --- | --- |
| i8 | -128 | 127 |giá trị ký nhỏ|
| u8 | 0 | 255 |một byte|
| i16 | -32768 | 32767 |dữ liệu số nguyên nhỏ|
| u16 | 0 | 65535 |Trường cổng/16-bit|
| i32 | -2147483648 | 2147483647 |Các phép tính số nguyên nhỏ thông dụng|
| u32 | 0 | 4294967295 |Trường bit 32 bit|

Wave cũng cung cấp số nguyên có dấu và không dấu gồm 64, 128, 256, 512 và 1024 bit. Loại rộng hơn không làm cho mọi phép tính trở nên an toàn: kết quả vẫn có thể vượt quá phạm vi đã chọn. Chọn phạm vi bạn cần đầu tiên. `isz` và `usz` tuân theo độ rộng địa chỉ đích.

Phải phân biệt giữa việc lưu trữ một chữ lớn trong một loại nhỏ và cố tình loại bỏ các bit bằng cách thực hiện cast. Việc chuyển đổi sang loại nhỏ hơn chỉ để loại bỏ lỗi có thể tự thay đổi giá trị. Các phép biến đổi sẽ được đề cập trong chương tiếp theo.

## Các kiểu dấu phẩy động và bool

`f32`·`f64` là số dấu phẩy động. Không giống như số nguyên, chúng có thể biểu diễn phần thập phân nhưng không thể lưu trữ chính xác tất cả các số thập phân. Đây là một lý do tại sao số tiền được quản lý theo đơn vị số nguyên nhỏ.

bool đại diện cho đúng và sai. Bạn có thể đặt tên cho điều kiện của mình bằng cách lưu kết quả so sánh như thế này:

<!-- wave-example: book-named-condition -->
```wave
fun main() {
    var balance: i32 = 5000;
    var cost: i32 = 3600;
    var can_buy: bool = balance >= cost;
    if (can_buy) {
        println("purchase allowed");
    } else {
        println("not enough money");
    }
}
```

Kết quả thực hiện:

```text
purchase allowed
```

`can_buy` không tự động tuân theo các thay đổi trong balance và cost. Sau khi hai giá trị được hoán đổi, việc so sánh sẽ được thực hiện lại nếu cần trạng thái hiện tại.

## không gian lưu trữ chưa được khởi tạo

`var value: i32;` là biểu mẫu chỉ khai báo dung lượng lưu trữ. Một giá trị hợp lệ phải được viết trước khi đọc. Đừng cho rằng 0 được tự động chèn vào chỉ vì bạn khai báo nó. Trong khóa học giới thiệu, sẽ dễ hiểu hơn nếu giá trị được biết và khởi tạo ngay lập tức cùng lúc với việc khai báo.

Trong các cuộc gọi thư viện nhận các giá trị làm đối số đầu ra, có những trường hợp khoảng trắng được khai báo trước rồi mới đọc khi thành công. Khi đó bạn nên kiểm tra kết quả thành công của hàm. Tránh lỗi đọc kết quả chưa được khởi tạo sau cuộc gọi thất bại.

## Phạm vi hợp lệ của các khối và tên

Khối là một vùng mã được bao bọc trong dấu ngoặc nhọn. Nếu bạn khai báo một biến mới có cùng tên bên trong thì biến mới sẽ được sử dụng bên trong khối đó. Cái này được gọi là shadowing.

<!-- wave-example: book-shadowing -->
```wave
fun main() {
    var value: i32 = 10;

    if (true) {
        var value: i32 = value + 5;
        println("inner={}", value);
        value += 1;
        println("inner changed={}", value);
    }

    println("outer={}", value);
}
```

Kết quả thực hiện:

```text
inner=15
inner changed=16
outer=10
```

Biểu thức ban đầu `value + 5` trong phần khai báo bên trong đọc value bên ngoài. Sau khi quá trình khởi tạo biến mới hoàn tất, bên trong value là 15. Ngay cả khi bạn thay đổi giá trị bên trong thành 16, không gian lưu trữ bên ngoài vẫn không thay đổi. Sau khi chặn, bạn sẽ lại thấy value ở bên ngoài.

Ngược lại, nếu bạn chỉ thực thi `value += 1` mà không có `var` trong khối bên trong, các biến hiển thị hiện có sẽ bị thay đổi. Nhìn vào từ khóa để xác định xem đó là khai báo mới hay thay đổi giá trị hiện có.

## Biến cục bộ và lưu trữ cấp cao nhất

Ngoài chức năng, bạn có thể sử dụng const và static. const đại diện cho một giá trị không đổi và static là không gian lưu trữ được duy trì trong quá trình thực thi. Chúng không thể được khai báo ở mọi nơi giống như các biến cục bộ.

<!-- wave-example: book-storage -->
```wave
const LIMIT: i32 = 3;
static visits: i32 = 0;

fun visit() {
    visits += 1;
    println("visit={}", visits);
}

fun main() {
    visit();
    visit();
    println("limit={}", LIMIT);
}
```

Kết quả thực hiện:

```text
visit=1
visit=2
limit=3
```

Ngay cả khi visit được gọi hai lần, static sẽ không được đặt lại về 0 trong mỗi cuộc gọi. Mặt khác, nếu bạn khai báo nó là `var visits: i32 = 0;` bên trong một hàm, nó sẽ khởi tạo bộ nhớ cục bộ mỗi khi bạn gọi nó. Trạng thái có thể thay đổi được chia sẻ có thể khiến hành vi khó theo dõi, vì vậy trước tiên hãy xem xét liệu hành vi đó có thể được giải quyết bằng đầu vào và đầu ra của hàm hay không.

## những sai lầm phổ biến

- Khi khai báo và gán bị nhầm lẫn và khai báo lại cùng một tên là không cần thiết.
- Nếu bạn cho rằng biến lưu trữ kết quả tính toán sẽ tự động tuân theo những thay đổi trong biến đầu vào.
- Nếu bạn nghĩ rằng vì loại giống nhau nên các đơn vị như số lượng và số byte cũng giống nhau.
- Khi đọc mà không khởi tạo hoặc khi đối số đầu ra được đọc khi hàm bị lỗi.
- Nếu bạn cho rằng tên của biến cục bộ hiển thị bên ngoài khối.

Khi tìm tên bị lỗi, hãy kiểm tra vị trí khai báo và phạm vi dấu ngoặc nhọn của tên đó.

## Bài tập: Tính toán thay đổi hàng tồn kho

Số lượng ban đầu là 20 chiếc và bán hai lần, mỗi lần 3 chiếc. In ra hàng tồn kho còn lại và tổng số đơn vị đã bán. Bất cứ khi nào hàng tồn kho thay đổi, các biến tương tự sẽ được cập nhật và khối lượng bán hàng cũng được tích lũy riêng.

### Lời giải đầy đủ

<!-- wave-example: book-stock-solution -->
```wave
fun main() {
    var stock: i32 = 20;
    var sold: i32 = 0;
    var order: i32 = 3;
    stock -= order;
    sold += order;
    stock -= order;
    sold += order;
    println("stock={} sold={}", stock, sold);
}
```

Kết quả thực hiện:

```text
stock=14 sold=6
```

Lượng hàng tồn kho và doanh số bán hàng phải thay đổi cùng nhau. Việc cập nhật một trong hai sẽ phá vỡ mối quan hệ giữa các giá trị. Chúng ta sẽ nhóm việc xử lý đơn hàng lặp đi lặp lại bằng cách học các hàm và câu lệnh vòng lặp.


## Các loại số nguyên và dấu phẩy động

Các kiểu số nguyên như sau:

- Đã ký: `i8`, `i16`, `i32`, `i64`, `i128`, `i256`, `i512`, `i1024`
- Chưa ký: `u8`, `u16`, `u32`, `u64`, `u128`, `u256`, `u512`, `u1024`
- Số nguyên kích thước địa chỉ: `isz`, `usz`
- Điểm nổi: `f32`, `f64`

`isz` là loại số nguyên có dấu phù hợp với kích thước địa chỉ và `usz` là loại số nguyên không dấu phù hợp với kích thước địa chỉ.

## Các loại tích hợp khác

|loại|sử dụng|
| --- | --- |
| `bool` |`true` hoặc `false`|
| `char` |Giá trị ký tự 8 bit không dấu. Không tùy ý loại điểm mã Unicode|
| `byte` |Giá trị byte 8 bit|
| `str` |Chuỗi byte kết thúc bằng NUL|
| `ptr<T>` |Nhắm mục tiêu theo con trỏ `T`|
| `array<T, N>` |Mảng có độ dài cố định với loại phần tử `T` và độ dài `N`|

Các cấu trúc, bảng liệt kê và bí danh loại do người dùng xác định cũng có thể được sử dụng ở các vị trí loại.

`var` là cú pháp khai báo biến cục bộ. Bí danh loại là một ngữ pháp biểu thị cùng loại với tên phù hợp với ngữ cảnh của mã.
