---
translation_set_id: stdlib-path-env
path: stdlib/path-env
locale: zh
group: stdlib
group_order: 1
order: 12
title: path 和 env: 路径与环境设置
summary: 读取调用者缓冲区中的路径和环境变量并识别容量错误。
---

## 路径组合

```text
std::path::copy
path_join2(dst: ptr<u8>, dst_cap: i32, left: str, right: str) -> i32
path_basename_copy(dst: ptr<u8>, dst_cap: i32, path: str) -> i32
path_dirname_copy(dst: ptr<u8>, dst_cap: i32, path: str) -> i32
```

容量包括最后一个NUL空间。成功结果是除NUL之外的任何长度，失败为-1。只有成功后我们才会使用目的地作为字符串。这些函数对路径字符串进行操作，并且不检查文件是否存在或访问权限。单独组合路径并不能防止目录转义或验证实际文件的身份。

## 环境变量

```text
std::env::environ
env_get(name: str, dst: ptr<u8>, dst_cap: i64) -> i64
env_exists(name: str) -> bool
env_get_i64(name: str) -> EnvResult<i64>
env_get_i32(name: str) -> EnvResult<i32>
```

成功时，`env_get`返回不包括NUL的长度。调用者缓冲区必须能够容纳NUL。空值与无密钥错误不同，因为它可以以长度 0 成功。

获取并区分从`std::env::consts`到NOT_FOUND、NO_SPACE、INVALID_KEY、READ、SOURCE_INCOMPLETE、NO_MEMORY的错误。不要将缓冲区外视为丢失密钥。要查找数字，请检查结果中的ok，然后使用value。不要自动将环境变量的内容视为可信设置；检查它们的范围和类型。

以下示例组合了目录 data 和文件名 input.txt。

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

执行结果：

```text
data/input.txt
```

## 分割目录和文件名

以下示例将路径复制为两个缓冲区。原始文件不需要实际存在。

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

执行结果：

```text
directory=data
filename=report.txt
```

两个缓冲区将一直有效，直到 main 结束。 `as str` 读取与字符串相同的缓冲区，而不分配新字符串。因此，如果修改缓冲区，读取到该地址的字符串也会改变。

## 设置默认首选项

环境变量是在程序外部传递的设置。读取数字设置时，请检查“能否将其读取为整数？”以及“是否在这个程序允许的范围内？”

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

如果 `WAVE_EXAMPLE_WORKERS` 不存在或无法读取为整数，则使用默认值 4。如果设置了 1 到 32 之间的整数，则使用该值，如果它是该范围之外的整数，则终止并出现错误。

Linux/macOS:

```shell
WAVE_EXAMPLE_WORKERS=8 wavec run main.wave
```

PowerShell:

```powershell
$env:WAVE_EXAMPLE_WORKERS = "8"
wavec run main.wave
```

在这两种情况下，都会输出`workers=8`。上面的示例选择了一个简单的默认策略。如果这是必需的设置，请将数字查找失败视为错误，而不是将其替换为默认值。如果需要区分丢失密钥、缓冲区不足和读取失败，请使用 env_get 和 ENV_ERR_* 常量。
