---
translation_set_id: stdlib-contracts
path: stdlib/contracts
locale: vi
group: stdlib
group_order: 1
order: 2
title: Đọc tài liệu API: Lỗi và quyền sở hữu
summary: Hiểu các đơn vị đối số, cấu trúc kết quả, thành công một phần và vòng đời tài nguyên.
---

## Đọc tuyên bố

Ký hiệu sau đây mô tả việc khai báo hàm và không phải là toàn bộ tệp thực thi.

```text
io_read(fd: i64, buf: ptr<u8>, len: i64) -> i64
```

`fd` là bộ mô tả mở, `buf` là bộ nhớ do người gọi cung cấp và `len` là số byte có sẵn để ghi. `i64` không có nghĩa là độ dài âm là hợp lệ. Giá trị trả về là số byte thực tế được đọc, không phải độ dài yêu cầu, do đó chỉ sử dụng phạm vi được trả về.

## Biểu hiện lỗi khác nhau tùy theo chức năng

|đường|vâng|Phương pháp kiểm tra|
| --- | --- | --- |
|Con trỏ hoặc null| `mem_alloc` |null Truy cập bộ nhớ sau khi kiểm tra|
|số byte hoặc số âm| `io_read` |Lỗi âm, 0 EOF, dữ liệu dương|
|mã trạng thái| `buffer_push` |Lỗi so sánh liên tục với `BUFFER_OK`|
|Thành công và giá trị| `NetResult<T>` |Sau khi kiểm tra `ok`, hãy sử dụng `value`|
|Bao gồm tiến độ một phần| `RandomFillResult` |Kiểm tra `ok`, `written`, `error` cùng nhau.|

Nó chỉ xem xét số lượng lỗi và không so sánh chúng với các hằng số trong các mô-đun khác. Ví dụ: số lỗi env và OS errno không phải là cùng một hệ thống. Lỗi ban đầu trong WASI không nên được hiểu là Linux errno.

## sở hữu và cho thuê

- **Owned**: Khi có được bộ nhớ được phân bổ, tệp đang mở hoặc ổ cắm mở, nó có trách nhiệm gọi bản phát hành/đóng tương ứng.
- **Borrow**: Byte view hoặc bộ đệm được truyền cho hàm đề cập đến bộ nhớ hiện có. Nếu một hàm không chỉ định rằng nó nhận được quyền sở hữu thì hàm gọi sẽ giữ quyền kiểm soát.
- **Đối số đầu ra**: Truyền một không gian lưu trữ hợp lệ nơi kết quả có thể được ghi vào hàm nhận nó, chẳng hạn như `out_value: ptr<T>`. Đảm bảo rằng hợp đồng nêu rõ rằng kết quả chỉ có giá trị nếu thành công.

Sao chép cấu trúc Buffer có thể khiến cả hai bản sao trỏ đến cùng một phân bổ. Không giải phóng từng bản sao riêng biệt. Một con trỏ mượn sẽ trở nên không hợp lệ sau khi việc phân bổ được giải phóng hoặc phân bổ lại. Chuỗi ký tự không phải là bộ đệm có thể ghi.

## Thất bại không có nghĩa là quay trở lại trạng thái trước đó

`io_write_all` có thể bị lỗi sau khi ghi một số byte. Các byte đã được ghi ra bên ngoài sẽ không được trả lại. Mặt khác, việc đọc checked cursor của bytes sẽ giữ nguyên vị trí và giá trị đầu ra nếu thất bại. Những khác biệt này được chỉ định bởi API.

Nếu bạn cũng muốn thực hành xử lý lỗi, hãy tiếp tục với [Trình đọc tập tin](/docs/vi/practice/file-reader) và [Tin nhắn nhị phân](/docs/vi/practice/binary-message).
