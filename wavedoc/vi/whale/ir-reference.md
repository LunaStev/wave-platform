---
translation_set_id: whale-ir-reference
path: whale/ir-reference
locale: vi
group: whale
group_order: 1
order: 6
title: Tài liệu tham khảo Whale IR
summary: Mô tả các loại, mã định danh, tính hợp lệ của hàm, thứ tự đánh giá và định dạng trao đổi.
---

## Mô-đun và mã định danh

Một mô-đun bao gồm thông tin mục tiêu, định nghĩa chung và các chức năng. Các giá trị có loại rõ ràng. Giao diện người dùng phân giải tên, loại, tình trạng quá tải và tên chung của ngôn ngữ nguồn và tạo ra typed IR.

Các hàm và biến toàn cục sử dụng các không gian tên nội bộ khác nhau. Vì vậy, hàm và biến có thể có cùng tên. Mã định danh bên trong khác với tên kết nối bên ngoài `link_name` và tên bên ngoài do giao diện người dùng chỉ định. Whale không giải quyết xung đột bên ngoài bằng cách tự động tạo tên mới. Vui lòng tham khảo [Biểu tượng và liên kết](assembler-linker).

Mỗi định nghĩa giá trị có một mã định danh. Các định nghĩa không thể trùng lặp và siêu dữ liệu loại phải khớp với loại được chỉ định trong định nghĩa. Chỉ riêng tên không xác định được các định nghĩa, ngay cả khi các khai báo có cùng tên sẽ che khuất lẫn nhau.

## IR Cấu hình và đọc

Dưới đây là một ví dụ Rust hoàn chỉnh giúp xây dựng và xác minh hàm với thùng `ir`, sau đó xuất ra typed IR.

```rust
use ir::{ModuleBuilder, Target, Type};

fn main() {
    let target = Target::lookup("x86_64-whale-linux").unwrap();
    let mut module = ModuleBuilder::new(target.name(), target.data_layout());
    let mut function = module.begin_function("answer", vec![], Type::I32);
    let left = function.const_i32(40);
    let right = function.const_i32(2);
    let answer = function.add(Type::I32, left, right);
    function.ret(Some(answer));
    function.finish();
    let module = module.finish();
    ir::verify_module(&module).unwrap();
    print!("{}", ir::print_module(&module));
}
```

Máy in xuất ra IR:

```text
module {
  format_version 3
  semantics_version 1
  target "x86_64-whale-linux"
  datalayout { ptr=64, endian=little }

  declare @f0 "answer": whale () -> i32, linkage internal

  fn @f0 "answer"() -> i32, entry %b0 {
  %b0 "entry":
    %v0: i32 = const i32 40
    %v1: i32 = const i32 2
    %v2: i32 = add i32 %v0, %v1
    ret i32 %v2
  }

}
```

`%v0` và `%v1` là định nghĩa của hằng số i32. Sử dụng `%v2` được xác định bởi `add` làm giá trị trả về của i32 của hàm. Nếu bạn thay đổi giá trị trả về thành Bool sẽ gây ra lỗi xác minh vì không khớp với chữ ký hàm. Ngay cả khi cả hai toán hạng đều là hằng số, O0 vẫn duy trì lệnh cộng.

Mã này là đầu ra thực tế của máy in, không phải tệp đầu vào để chuyển tới trình phân tích cú pháp văn bản. Phân tích cú pháp văn bản và thực thi IR chưa được hỗ trợ, hiện tại mô-đun này có thể được định cấu hình là Rust builder.

## loại

|loại|ý nghĩa|
| --- | --- |
| `bool` |Giá trị logic false hoặc true|
| `i1`, `i8`, `i16`, `i32`, `i64`, `i128` |số nguyên có dấu của độ rộng bit được chỉ định|
| `u1`, `u8`, `u16`, `u32`, `u64`, `u128` |số nguyên không dấu có độ rộng bit được chỉ định|
| `f16`, `f32`, `f64` |Giá trị dấu phẩy động với độ rộng bit được chỉ định|
| `ptr<T>` |T Con trỏ để gõ giá trị|
| `fnptr<signature>` | Con trỏ có thể gọi được với các loại tham số/kết quả chính xác và quy ước gọi |
| `array<T, N>` |N phần tử cùng loại|
| `struct{T, ...}` |Trường cấu trúc có thứ tự|
| `tuple<T, ...>` |các phần tử tuple được sắp xếp|
| `void` |Không có kết quả|

`bool`, `i1` và `u1` là các loại khác nhau. signed `i1` đại diện cho −1 và 0, và unsigned `u1` đại diện cho 0 và 1. Số nguyên 1 không phải là điều kiện logic ẩn. Nhánh có điều kiện, điều kiện Select, `trap_if` yêu cầu toán hạng Bool.

Kích thước lưu trữ không chỉ được xác định bởi số bit trong giá trị và tuân theo [bố trí mục tiêu](memory-model). Ví dụ: giá trị của `i1` là 1 bit nhưng chiếm ít nhất 1 byte trong bộ nhớ.

## Chức năng và cuộc gọi

Hàm chỉ định tất cả các tham số, loại kết quả, quy ước gọi và linkage. Cuộc gọi trực tiếp và gián tiếp phải trùng với chữ ký của người gọi. Cuộc gọi void không có kết quả ID. Lệnh gọi tới nonvoid sẽ giữ nguyên định nghĩa kết quả ngay cả khi O0 không sử dụng kết quả.

Kết quả trả về phải khớp với loại kết quả của hàm. Trả về void không mang giá trị và trả về nonvoid mang giá trị của loại kết quả được khai báo.

### Tuyên bố, nhận dạng và kêu gọi

`Module.declarations` ghi lại `FunctionId`, tên, chữ ký đầy đủ, liên kết và tên liên kết bên ngoài của từng chức năng. Một định nghĩa đề cập đến danh tính này; kiểu tham số và trả về phải khớp với khai báo của nó. Các nội dung khai báo lặp lại giống hệt nhau sẽ chuyển thành cùng một ID thông qua `declare_function`; xung đột và định nghĩa trùng lặp là lỗi. Một khai báo nội bộ cần có phần thân trong mô-đun. Một khai báo bên ngoài có thể chưa được giải quyết cho đến khi liên kết hoặc có nội dung được xuất. Các chức năng bên trong không có `link_name`; các hàm bên ngoài yêu cầu tên không trống rõ ràng và không có NUL. Hai khai báo hàm riêng biệt không thể có cùng tên bên ngoài. Toàn cầu và hàm vẫn sử dụng các không gian tên nội bộ riêng biệt.

Đăng ký các khai báo trước khi xây dựng các phần thân với `begin_declared_function` để hỗ trợ các lệnh gọi chuyển tiếp và đệ quy. `begin_function` vẫn là sự tiện lợi cho chức năng Whale nội bộ mới. Các API đã kiểm tra `declare_function`, `begin_declared_function`, `function_addr`, `null_function` và `call` trả về `Result`; cuộc gọi bị từ chối không nối thêm lệnh hoặc phân bổ ID kết quả của nó.

Chương trình Rust hoàn chỉnh sau đây khai báo một hàm bên ngoài, lấy địa chỉ đã nhập của nó và thực hiện cả lệnh gọi trực tiếp và gián tiếp:

```rust
use ir::{Callee, CallingConvention, DataLayout, FunctionSignature, Linkage, ModuleBuilder, Type};

fn main() {
    let mut module = ModuleBuilder::new("x86_64-whale-linux", DataLayout::default_64bit_le());
    let signature = FunctionSignature {
        params: vec![Type::I32], ret: Type::I32,
        convention: CallingConvention::SysV64, variadic: false,
    };
    let identity = module.declare_function(
        "identity", signature, Linkage::External, Some("identity_i32".into()),
    ).unwrap();
    let mut function = module.begin_function("answer", vec![], Type::I32);
    let input = function.const_i32(42);
    let callback = function.function_addr(identity).unwrap();
    // The direct call's result remains defined even though it is unused.
    function.call(Callee::Direct(identity), vec![input]).unwrap();
    let result = function.call(Callee::Indirect(callback), vec![input]).unwrap().unwrap();
    function.ret(Some(result));
    function.finish();
    let module = module.finish();
    ir::verify_module(&module).unwrap();
    print!("{}", ir::print_module(&module));
}
```

```text
module {
  format_version 3
  semantics_version 1
  target "x86_64-whale-linux"
  datalayout { ptr=64, endian=little }

  declare @f0 "identity": sysv64 (i32) -> i32, linkage external, link_name "identity_i32"
  declare @f1 "answer": whale () -> i32, linkage internal

  fn @f1 "answer"() -> i32, entry %b0 {
  %b0 "entry":
    %v0: i32 = const i32 42
    %v1: fnptr<sysv64 (i32) -> i32> = function_addr @f0
    %v2: i32 = call sysv64 i32 @f0(%v0)
    %v3: i32 = call sysv64 i32 indirect %v1(%v0)
    ret i32 %v3
  }

}
```

`Callee::Direct(FunctionId)` giải quyết thông qua bảng khai báo; `Callee::Indirect(ValueId)` yêu cầu giá trị `Type::FnPtr(FunctionSignature)`. Chữ ký bao gồm tất cả các loại tham số, loại kết quả và `CallingConvention::{Whale, SysV64}`. Nó được giữ lại thông qua các bản sao, lưu trữ, tham số, trả về, phi và chọn. Con trỏ dữ liệu và giá trị số nguyên không thể gọi được. Các kiểu ép kiểu liên quan đến kiểu con trỏ hàm bị từ chối; thay đổi chú thích loại không thể thay đổi chữ ký có thể gọi được. Một con trỏ hàm có bộ lưu trữ địa chỉ 64-bit trên mục tiêu này; bản thân điều này không triển khai siêu dữ liệu bóng thời gian chạy.

Tính chất, loại đối số/kết quả chính xác, sự hiện diện của ID kết quả và quy ước gọi hàm phải khớp. Không có chuyển đổi tiềm ẩn. Người được gọi gián tiếp phải thống trị cuộc gọi giống như các đối số của nó. `variadic: true`, tham số void và SysV64 chữ ký tham số/kết quả tổng hợp bị từ chối. Whale chữ ký tổng hợp có thể được thể hiện trong IR; Phân loại ABI gốc và phát lệnh gọi máy vẫn chưa có sẵn cho cả hai quy ước.

`null_function(signature)` đại diện cho một con trỏ hàm null được gõ. Gọi nó là gõ đúng IR với bẫy thời gian chạy bắt buộc trước khi vào callee. Mục tiêu không có giá trị không hợp lệ, hết hạn hoặc không tương thích với chữ ký đã kiểm tra cũng phải bẫy. Các hoạt động kiểm tra thời gian chạy và quản lý thời gian gọi lại nước ngoài này đang chờ lớp trình thông dịch/thực thi gốc; trình xác minh thành công không có nghĩa là các địa chỉ bên ngoài tùy ý là an toàn.

### Dạng lời gọi trong AST

Đây là các đoạn biểu thức bên trong chương trình định dạng AST 2:

```json
{"Call":{"callee":{"Direct":"increment"},"args":[{"Lit":{"Int":{"bits":32,"signed":true,"value":"41"}}}]}}
```

```json
{"Call":{"callee":{"Indirect":{"FunctionRef":"increment"}},"args":[{"Lit":{"Int":{"bits":32,"signed":true,"value":"41"}}}]}}
```

`Direct` và `FunctionRef` sử dụng không gian tên hàm ngay cả khi một biến có cùng tên. `Indirect` đánh giá biểu thức của nó trước, sau đó đánh giá các đối số từ trái sang phải. Lệnh gọi void có giá trị dưới dạng `ExprStmt` nhưng không hợp lệ dưới dạng giá trị khởi tạo biến, đối số, toán hạng hoặc giá trị trả về. Các lệnh gọi và tham chiếu hàm không phải là các biểu thức hằng số tại thời điểm biên dịch. `NullFunction` lấy đối tượng chữ ký với các trường `params`, `ret`, `convention` và `variadic`.

[Ví dụ hoàn chỉnh JSON](https://github.com/wavefnd/Whale/blob/master/ir/tests/fixtures/ast-v2-calls.json) lưu trữ lệnh gọi lại và gọi nó trước lệnh gọi bên ngoài. Hạ nó xuống bằng:

```sh
cargo run --locked --features socket-cli -- ir lower ir/tests/fixtures/ast-v2-calls.json
```

[Dự kiến ​​IR](https://github.com/wavefnd/Whale/blob/master/ir/tests/fixtures/calls-v3.wir) của nó đã được kiểm tra trong các thử nghiệm hạ thấp. Danh tính hàm và tên liên kết được thể hiện ở ranh giới IR; bảo tồn chúng thông qua việc tạo và liên kết đối tượng gốc vẫn là công việc riêng biệt.

## Sự sẵn có của các khối và giá trị

Mỗi khối có một mã định danh duy nhất và chính xác là một terminator. Các mục tiêu nhánh phải thuộc cùng một chức năng. Khối nhập phải tồn tại và không được có cạnh đầu và phi. Khi tạo một vòng lặp, hãy phân nhánh từ khối đầu vào tới một tiêu đề vòng lặp riêng biệt.

Định nghĩa về giá trị trong đường dẫn thực thi sẽ chi phối điểm sử dụng thông thường của nó. Điều này có nghĩa là tất cả các đường dẫn từ điểm vào đến điểm sử dụng đều phải đi qua định nghĩa đó. Trong cùng một khối, định nghĩa phải có trước khi sử dụng. Thứ tự lưu trữ của các khối không quyết định sự thống trị.

Giả sử rằng điểm vào phân nhánh thành left hoặc right rồi nối tại join. Các giá trị chỉ được xác định trong left không thể được sử dụng làm giá trị chung trong join. Điều này là do đường dẫn đi qua right không được xác định. Các giá trị của mỗi khối trước phải được kết hợp thành phi, được nhận làm đầu vào.

Các khối không thể truy cập cũng được bảo tồn trong mô-đun. Trình xác minh liên tục kiểm tra mã định danh, loại, toán hạng và cấu trúc nhánh của khối. Định nghĩa về khối không thể truy cập không thể cung cấp giá trị cho việc sử dụng thông thường đường dẫn có thể truy cập.

## lệnh Phi

phi được đặt trước tất cả các lệnh thông thường trong khối. Cần có chính xác một đầu vào cho mỗi khối trước đó. Giá trị đầu vào phải thuộc loại phi và phải có sẵn ở cuối khối tương ứng trước đó.

Ngay cả khi có nhiều cạnh trong một khối trước đó thì cũng chỉ có một đầu vào. Vòng lặp phi có thể tham chiếu đến giá trị của khối xuất hiện sau trong thứ tự lưu trữ mô-đun miễn là đó là giá trị được tính ở cạnh lặp lại. Các khối trước bị thiếu, trùng lặp hoặc không liên quan và đầu vào được nhập không chính xác là lỗi xác thực.

## Đánh giá và lựa chọn

Whale AST đánh giá mục tiêu cuộc gọi và biểu thức con từ bên trái nơi chúng được chỉ định. Giao diện người dùng thể hiện đánh giá ngắn mạch như một nhánh luồng điều khiển.

Select chọn một trong các giá trị đã được tính toán. Nó không bỏ qua việc tính toán của một trong hai đầu vào. Ví dụ: ngay cả khi bạn chọn một giá trị an toàn, bạn không thể tránh gặp phải trap khi tính toán các đầu vào khác. Các phép tính chỉ cần được thực hiện trên một đường dẫn cụ thể phải được đặt bên trong một khối có điều kiện.

## Phòng Kiểm định trap

IR không hợp lệ sẽ bị từ chối trong giai đoạn xác minh. Các vi phạm điều kiện thực thi được xử lý bằng các giá trị được xác định là trap và `undef` và `poison` đều không được chấp nhận. Việc sử dụng sai builder, định nghĩa trùng lặp, bổ sung terminator thứ hai sẽ được trả về dưới dạng lỗi cấu trúc.

trap chứa lý do·vị trí nguồn·IR ID và sau đó dừng thực thi. Chạy native sẽ kết thúc chương trình và trình thông dịch API trả về lỗi Trap. Giữ nguyên các tác dụng phụ trước đó, nhưng không đảm bảo bộ đệm flush·gọi hàm hủy·ngăn xếp unwinding.

Bảo đảm này áp dụng cho IR đã được xác minh và bộ nhớ theo dõi. Bên ngoài C·địa chỉ thô·tổ hợp nội tuyến có hợp đồng riêng và không phải lúc nào cũng phát hiện ra các vi phạm bên ngoài ranh giới của nó. Vui lòng tham khảo [Mô hình bộ nhớ](memory-model).

## Định dạng trao đổi và trình bày văn bản

AST và typed IR sử dụng format version tương ứng và semantics version chung tương ứng. Người đọc nên từ chối khóa JSON chưa được phiên bản/không xác định. Các nhà xây dựng không nên cho rằng các thuộc tính không được hỗ trợ sẽ bị bỏ qua một cách âm thầm.

Các số nguyên được truyền dưới dạng độ rộng bit·signedness·số chuỗi. Các hằng số dấu phẩy động được truyền dưới dạng chuỗi bit có chiều rộng và chính xác. Văn bản round-trip trong IR phải giữ nguyên tên·ID·loại·hằng·chuỗi·thuộc tính·siêu dữ liệu. Không gian và vị trí bình luận không được bảo tồn.

Bạn có thể sử dụng các hợp đồng vô hướng AST JSON bên dưới. Đầu ra typed IR bao gồm thông tin phiên bản nhưng chưa hỗ trợ trao đổi khứ hồi đầy đủ với trình phân tích cú pháp văn bản.


### Định danh được in và tên trong dấu ngoặc kép

Typed IR format 3 in hàm dưới dạng `@fN`, giá trị toàn cục là `@gN`, giá trị là `%vN`, khối là `%bN`. ID hàm và toàn cục thuộc mô-đun; ID giá trị và khối thuộc hàm chứa chúng. Các ID được cung cấp được giữ nguyên, kể cả khoảng trống trong số thứ tự. Tên trong ngoặc kép chỉ là chú thích, không dùng để phân giải tham chiếu. Hàm ghi rõ `entry %bN`, độc lập với thứ tự lưu các khối.

Mô-đun hoàn chỉnh sau đã được kiểm tra và in qua API IR Rust. Cả hai khối nhánh đều có tên `"branch"`; ID phân biệt định nghĩa và đầu vào phi.

```text
module {
  format_version 3
  semantics_version 1
  target "x86_64-whale-linux"
  datalayout { ptr=64, endian=little }

  declare @f0 "choose": whale (bool) -> i32, linkage internal

  fn @f0 "choose"(%v0 "condition": bool) -> i32, entry %b0 {
  %b0 "entry":
    cbr bool %v0, label %b1, label %b2
  %b1 "branch":
    %v1: i32 = const i32 1
    br label %b3
  %b2 "branch":
    %v2: i32 = const i32 2
    br label %b3
  %b3 "join":
    %v99: i32 = phi i32 [ %v1, %b1 ], [ %v2, %b2 ]
    ret i32 %v99
  }

}
```

`%v0` được định nghĩa trong danh sách tham số. `%b1` và `%b2` vẫn khác nhau dù cùng tên; phi chỉ rõ từng khối tiền nhiệm bằng ID. Nhánh và đích switch dùng cùng cú pháp ID khối. Bộ in không đánh số lại `%v99` được chỉ định rõ ràng.

Mọi trường tên và chuỗi đều dùng dấu ngoặc kép: đích, tên hàm, toàn cục, tham số và khối, tên liên kết ngoài, tên khai báo hằng và lý do trap. Unicode có thể in được giữ nguyên. Các escape là `\"`, `\\`, `\n`, `\r`, `\t`, `\0` và `\u{hex}` với chữ số thập lục phân viết thường cho các ký tự điều khiển khác cùng U+2028/U+2029. Tên chứa xuống dòng, tab, dấu ngoặc kép, dấu gạch chéo ngược và tiếng Hàn vẫn được in trong một bản ghi.

```text
"line\ncolumn\tquote\"slash\\한글"
```

Định nghĩa đầu ra format 2 phải chuyển sang ID rõ ràng, ID tham số, tên trong ngoặc kép và tham chiếu khối vào. Chỉ cú pháp typed IR thay đổi; AST JSON format 2 và semantics version 1 vẫn giữ nguyên. Bộ phân tích văn bản và bộ đọc round-trip chưa có.

### Phiên bản được chỉ định AST JSON

Lưu phần sau dưới dạng `program.json`. Tất cả bốn trường phong bì là bắt buộc. `program` chứa các mảng `declarations`, `globals` và `functions` bắt buộc, có thể trống. Tên hàm, tham số, kiểu trả về, nội dung, `convention` và `linkage` là bắt buộc. `link_name` có thể vắng mặt/null đối với các hàm bên trong và phải là một chuỗi khác trống không có NUL đối với các hàm bên ngoài. Mỗi enum sử dụng tên đơn vị của nó hoặc một đối tượng khóa biến thể duy nhất. Các biến thể đơn vị cũng chấp nhận đối tượng có giá trị null, chẳng hạn như `{"Void":null}`; bộ mã hóa phát ra tên đơn vị `"Void"`. `VarDecl.init` có thể vắng mặt hoặc vô hiệu; các trường bắt buộc khác phải có mặt.

```json
{
  "format_version": 2,
  "semantics_version": 1,
  "features": [],
  "program": {
    "globals": [],
    "functions": [
      {
        "name": "answer",
        "parameters": [],
        "return_type": {
          "Int": {
            "bits": 128,
            "signed": false
          }
        },
        "body": [
          {
            "Return": {
              "Lit": {
                "Int": {
                  "bits": 128,
                  "signed": false,
                  "value": "340282366920938463463374607431768211455"
                }
              }
            }
          }
        ],
        "convention": "Whale",
        "linkage": "Internal",
        "link_name": null
      }
    ],
    "declarations": []
  }
}
```

```sh
cargo run --locked --features socket-cli -- ir lower program.json
```

```text
module {
  format_version 3
  semantics_version 1
  target "x86_64-whale-linux"
  datalayout { ptr=64, endian=little }

  declare @f0 "answer": whale () -> u128, linkage internal

  fn @f0 "answer"() -> u128, entry %b0 {
  %b0 "entry":
    %v0: u128 = const u128 340282366920938463463374607431768211455
    ret u128 %v0
  }

}
```

Số nguyên `value` là một chuỗi thập phân. Signed Một số được sử dụng sau dấu trừ tùy chọn của một số nguyên và không được phép có dấu cách, dấu cộng, số mũ và dấu phân cách. Phạm vi cho phép được xác định bởi chiều rộng đã khai báo và signedness. `u128::MAX` ở trên được giữ nguyên thông qua JSON và lowering. Các giá trị unsigned âm hoặc các giá trị ngoài phạm vi không phải là wrap và là lỗi. Giá trị Float sử dụng chuỗi bit thập lục phân có chiều rộng chính xác được mô tả trong [Các phép toán số](numeric-operations).

`format_version` là 2 cho định dạng AST này; `semantics_version` là 1. `features` phải là một mảng trống. Các trường, phiên bản, tính năng không xác định, khóa JSON thô trùng lặp (bao gồm cả các khóa tương đương đã thoát) và các giá trị theo sau đều là lỗi, ngay cả với `--no-verify`. Điểm vào thư viện là `ir::lower_ast::interchange::decode`; `encode` phát ra phong bì. `decode` mặc định có giới hạn 8 byte nguồn MiB; `decode_with_limit` chấp nhận giới hạn người gọi. JSON lồng nhau bị giới hạn. Sử dụng bộ giải mã thô này thay vì phân tích cú pháp thành một bản đồ chung có thể loại bỏ các khóa trùng lặp.

[Lược đồ JSON hoàn chỉnh](https://github.com/wavefnd/Whale/blob/master/ir/schema/ast-v2.schema.json) chỉ định hình dạng, trường bắt buộc và biến thể. Việc kiểm tra phạm vi/loại và phát hiện khóa trùng lặp cũng được áp dụng. Tập hợp con hạ thấp vô hướng bao gồm hằng số, biến/hằng, cộng/phụ/mul, so sánh, gán, if/while, trả về và ngắt/tiếp tục. Hỗ trợ tham chiếu chức năng, gọi trực tiếp và gọi gián tiếp; biểu thức tổng hợp không được hỗ trợ. `Opaque` có ​​thể biểu thị trong lược đồ nhưng không được hỗ trợ bằng cách hạ thấp.

Di chuyển yêu cầu bọc Program cũ chưa có envelope và thay số JSON bằng chuỗi số nguyên thập phân hoặc chuỗi bit số thực. Đầu vào không có phiên bản bị từ chối. Format 1 phải chuyển sang format 2, thêm `program.declarations` (mảng rỗng nếu không dùng) và `convention`/`linkage` rõ ràng trong định nghĩa. Các phiên bản độc lập: AST format 2, typed IR format 3 và semantics version 1.

### Đầu vào bị từ chối và khôi phục CLI

Lưu thông tin đầu vào hoàn chỉnh sau đây dưới dạng `invalid.json`.

```json
{
  "format_version": 99,
  "semantics_version": 1,
  "features": [],
  "program": {
    "globals": [],
    "functions": [],
    "declarations": []
  }
}
```

```sh
cargo run --locked --features socket-cli -- ir lower invalid.json -o rejected.wir
```

```text
Failed to parse socket JSON: unsupported AST format_version 99; expected 2
```

Lệnh thoát với trạng thái khác 0 và không tạo đầu ra mới hoặc ghi đè lên các tệp hiện có. Loại không khớp cũng sẽ không thành công trước khi xuất bản đầu ra. Các tệp nhị phân được tạo không có `socket-cli` thoát với trạng thái 2 và xuất ra lệnh khôi phục có chứa `--features socket-cli`.
