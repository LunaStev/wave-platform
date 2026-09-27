---
translation_set_id: stdlib-path-env
path: stdlib/path-env
locale: en
group: stdlib
group_order: 1
order: 12
title: path and env: Paths and environment settings
summary: Reads the path and environment variables in the caller buffer and identifies capacity errors.
---

## path combination

```text
std::path::copy
path_join2(dst: ptr<u8>, dst_cap: i32, left: str, right: str) -> i32
path_basename_copy(dst: ptr<u8>, dst_cap: i32, path: str) -> i32
path_dirname_copy(dst: ptr<u8>, dst_cap: i32, path: str) -> i32
```

Capacity includes the last NUL space. A success result is any length excluding NUL, failure is -1. Only on success will we use the destination as a string. These functions operate on path strings and do not check file existence or access permissions. Combining paths alone does not prevent directory escapes or verify the identity of actual files.

## environment variable

```text
std::env::environ
env_get(name: str, dst: ptr<u8>, dst_cap: i64) -> i64
env_exists(name: str) -> bool
env_get_i64(name: str) -> EnvResult<i64>
env_get_i32(name: str) -> EnvResult<i32>
```

`env_get` returns the length excluding NUL on success. The caller buffer must be able to hold up to NUL. An empty value differs from a keyless error because it can succeed with length 0.

Get and distinguish errors from `std::env::consts` to NOT_FOUND, NO_SPACE, INVALID_KEY, READ, SOURCE_INCOMPLETE, NO_MEMORY. Don't treat out of buffer as missing key. To look up numbers, check ok in the result and then use value. Don't automatically treat the contents of environment variables as trusted settings; check their scope and type.

The following example combines the directories data and the file names input.txt.

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

Execution result:

```text
data/input.txt
```

## Split directories and file names

The following example copies a path split into two buffers. The original file does not need to actually exist.

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

Execution result:

```text
directory=data
filename=report.txt
```

Both buffers will remain in effect until the end of main. `as str` reads the same buffer as a string without allocating a new string. Therefore, if you modify the buffer, the string read to that address will also change.

## Set default preferences

Environment variables are settings passed outside of the program. When reading a numeric setting, check “Can it be read as an integer?” and “Is it within the range allowed by this program?”

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

If `WAVE_EXAMPLE_WORKERS` is not present or cannot be read as an integer, the default value of 4 is used. If an integer between 1 and 32 is set, that value is used, and if it is an integer outside the range, it terminates with an error.

Linux/macOS:

```shell
WAVE_EXAMPLE_WORKERS=8 wavec run main.wave
```

PowerShell:

```powershell
$env:WAVE_EXAMPLE_WORKERS = "8"
wavec run main.wave
```

In both cases, `workers=8` is output. The example above selects a simple default policy. If this is a required setting, treat numeric lookup failures as errors rather than replacing them with default values. If you need to distinguish between missing key, insufficient buffer, and read failure, use the env_get and ENV_ERR_* constants.
