---
translation_set_id: stdlib-path-env
path: stdlib/path-env
locale: vi
group: stdlib
group_order: 1
order: 12
title: path và env: Đường dẫn và cấu hình môi trường
summary: Đọc các biến môi trường và đường dẫn trong bộ đệm người gọi và xác định lỗi dung lượng.
---

## kết hợp đường dẫn

```text
std::path::copy
path_join2(dst: ptr<u8>, dst_cap: i32, left: str, right: str) -> i32
path_basename_copy(dst: ptr<u8>, dst_cap: i32, path: str) -> i32
path_dirname_copy(dst: ptr<u8>, dst_cap: i32, path: str) -> i32
```

Dung lượng bao gồm khoảng trống NUL cuối cùng. Kết quả thành công là bất kỳ độ dài nào ngoại trừ NUL, lỗi là -1. Chỉ khi thành công, chúng ta mới sử dụng đích đến làm chuỗi. Các hàm này hoạt động trên chuỗi đường dẫn và không kiểm tra sự tồn tại của tệp hoặc quyền truy cập. Chỉ kết hợp các đường dẫn không ngăn chặn việc thoát thư mục hoặc xác minh danh tính của các tệp thực tế.

## biến môi trường

```text
std::env::environ
env_get(name: str, dst: ptr<u8>, dst_cap: i64) -> i64
env_exists(name: str) -> bool
env_get_i64(name: str) -> EnvResult<i64>
env_get_i32(name: str) -> EnvResult<i32>
```

`env_get` trả về độ dài không bao gồm NUL nếu thành công. Bộ đệm người gọi phải có khả năng chứa tới NUL. Giá trị trống khác với lỗi không cần khóa vì nó có thể thành công với độ dài 0.

Nhận và phân biệt các lỗi từ `std::env::consts` đến NOT_FOUND, NO_SPACE, INVALID_KEY, READ, SOURCE_INCOMPLETE, NO_MEMORY. Đừng coi hết bộ đệm là thiếu khóa. Để tra cứu các số, hãy đánh dấu ok trong kết quả rồi sử dụng value. Không tự động coi nội dung của các biến môi trường là cài đặt đáng tin cậy; kiểm tra phạm vi và loại của họ.

Ví dụ sau kết hợp các thư mục data và tên tệp input.txt.

<!-- wave-example: path-api -->
```wave
import("std::path::copy")::{
    path_join2
};

fun main() -> i32 {
    var output: array<u8, 64>;
    var length: i32 = path_join2(&output[0], 64, "data", "input.txt");
    if (length < 0) {
        return 1;
    }

    println("{}", &output[0] as str);
    return 0;
}
```

Kết quả thực hiện:

```text
data/input.txt
```

## Tách thư mục và tên tập tin

Ví dụ sau sao chép một đường dẫn được chia thành hai bộ đệm. Tệp gốc không cần phải thực sự tồn tại.

<!-- wave-example: book-path-parts -->
```wave
import("std::path::copy")::{
    path_basename_copy,
    path_dirname_copy
};

fun main() -> i32 {
    var directory: array<u8, 64>;
    var filename: array<u8, 64>;
    var directory_length: i32 = path_dirname_copy(&directory[0], 64, "data/report.txt");
    var filename_length: i32 = path_basename_copy(&filename[0], 64, "data/report.txt");

    if (directory_length < 0 || filename_length < 0) {
        return 1;
    }

    println("directory={}", &directory[0] as str);
    println("filename={}", &filename[0] as str);
    return 0;
}
```

Kết quả thực hiện:

```text
directory=data
filename=report.txt
```

Cả hai bộ đệm sẽ vẫn có hiệu lực cho đến hết main. `as str` đọc cùng một bộ đệm dưới dạng một chuỗi mà không phân bổ chuỗi mới. Vì vậy, nếu bạn sửa đổi bộ đệm thì chuỗi được đọc tới địa chỉ đó cũng sẽ thay đổi.

## Đặt tùy chọn mặc định

Biến môi trường là các cài đặt được truyền bên ngoài chương trình. Khi đọc cài đặt số, hãy chọn “Có thể đọc dưới dạng số nguyên không?” và “Nó có nằm trong phạm vi cho phép của chương trình này không?”

<!-- wave-example: book-env-setting -->
```wave
import("std::env::environ")::{
    EnvResult,
    env_get_i32
};

fun main() -> i32 {
    var setting: EnvResult<i32> = env_get_i32("WAVE_EXAMPLE_WORKERS");
    var workers: i32 = 4;

    if (setting.ok) {
        if (setting.value < 1 || setting.value > 32) {
            println("workers must be between 1 and 32");
            return 1;
        }

        workers = setting.value;
    }

    println("workers={}", workers);
    return 0;
}
```

Nếu `WAVE_EXAMPLE_WORKERS` không xuất hiện hoặc không thể đọc dưới dạng số nguyên thì giá trị mặc định là 4 sẽ được sử dụng. Nếu một số nguyên nằm trong khoảng từ 1 đến 32 được đặt thì giá trị đó sẽ được sử dụng và nếu đó là số nguyên nằm ngoài phạm vi thì giá trị đó sẽ kết thúc với một lỗi.

Linux/macOS:

```shell
WAVE_EXAMPLE_WORKERS=8 wavec run main.wave
```

PowerShell:

```powershell
$env:WAVE_EXAMPLE_WORKERS = "8"
wavec run main.wave
```

Trong cả hai trường hợp, `workers=8` là đầu ra. Ví dụ trên chọn một chính sách mặc định đơn giản. Nếu đây là cài đặt bắt buộc, hãy coi lỗi tra cứu số là lỗi thay vì thay thế chúng bằng giá trị mặc định. Nếu bạn cần phân biệt giữa khóa bị thiếu, bộ đệm không đủ và lỗi đọc, hãy sử dụng hằng số env_get và ENV_ERR_*.
