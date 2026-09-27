---
translation_set_id: diagnostics
path: reference/diagnostics
locale: vi
group: reference
group_order: 5
order: 2
title: Xử lý sự cố: Từ cài đặt đến thực thi
summary: Cô lập các bước thất bại và thu hẹp nguyên nhân bằng thông tin có thể tái tạo.
---

## Đầu tiên, phân biệt các giai đoạn thất bại

|hiện tượng quan sát được|Kiểm tra trước|hành động tiếp theo|
| --- | --- | --- |
|wavec Không tìm thấy lệnh|PATH và vị trí tệp thực thi|Chạy với đường dẫn tuyệt đối và đặt PATH|
|Không thể tìm thấy tập tin cần thiết để chạy|Có tập tin nào bị thiếu trong thư mục cài đặt không?|Giải nén và cài đặt lại toàn bộ gói|
|std import Không thành công| `wavec print std-path` |Tương ứng: Cài đặt std hoặc chỉ định `--std-root`|
|Vị trí nguồn và đầu ra lỗi loại| `wavec check main.wave` |Sửa lỗi đầu tiên và kiểm tra lại|
|Xây dựng lỗi cho các mục tiêu OS·CPU khác|Đã chỉ định target và môi trường mục tiêu|[Cài đặt xây dựng chéo](/docs/vi/whale/build-link-targets) Xác nhận|
|Thực thi không thành công sau khi xây dựng thành công|mã thoát, đầu vào, thư mục làm việc|Môi trường thực thi và kiểm tra lỗi API|

## Một ví dụ chẩn đoán nhỏ

Đây là toàn bộ chương trình, cố ý không chính xác:

```wave
fun main() {
    var count: i32 = 1;
    println("{}", missing);
}
```

`wavec check main.wave` phải trỏ đến tên chưa được khai báo missing. Đổi tên biến thành count, sau đó kiểm tra và chạy lại. Tập trung vào tệp, vị trí và nguyên nhân thay vì toàn bộ văn bản chẩn đoán. Những lỗi tiếp theo có thể là kết quả của lỗi ban đầu.

## Khi một chương trình thực thi bị lỗi

Trong shell Linux/macOS, mã thoát được kiểm tra ngay sau khi thực thi dưới dạng `echo $?` và trong PowerShell, đó là `$LASTEXITCODE`. Lỗi đầu vào và lỗi `return 1` rõ ràng không phải là cùng một nguyên nhân. Các giá trị thời gian chạy không hợp lệ đối với số ca hoặc chuyển đổi thực có thể gây ra trap. Hãy xem [quy tắc hoạt động](/docs/vi/language/expressions-and-operators).

Đường dẫn tệp tương đối bị ảnh hưởng bởi thư mục làm việc thực thi thay vì vị trí tệp nguồn. Đừng coi lỗi đọc tệp là độ dài chuỗi bằng 0, trước tiên hãy kiểm tra lỗi trả về. Lỗi kết nối mạng được kiểm tra thông qua tra cứu địa chỉ, chờ máy chủ, quyền và thời gian chờ.

## Thông tin cần thiết để báo cáo sự cố

1. `wavec --version` Đầu ra và lệnh chính xác được thực thi.
2. target. được chỉ định riêng biệt với máy chủ OS·kiến trúc
3. Nguồn trình biên dịch được sử dụng với đường dẫn std đã chọn.
4. Các tệp nguồn, đầu vào và yêu cầu tối thiểu để tái tạo sự cố.
5. Kết quả mong đợi, kết quả thực tế, chẩn đoán và mã thoát.

Mật khẩu, mã thông báo và nội dung tệp cá nhân sẽ bị xóa. Nếu vấn đề không còn nữa khi bạn rút gọn ví dụ tối thiểu thì phần tử cuối cùng bạn loại bỏ chính là manh mối. `--error-format=json` khả dụng khi công cụ thu thập chẩn đoán.

[Cài đặt](/docs/vi/getting-started/install) · [lệnh biên dịch](/docs/vi/getting-started/compiler) · [Mục tiêu và liên kết](/docs/vi/whale/build-link-targets)
