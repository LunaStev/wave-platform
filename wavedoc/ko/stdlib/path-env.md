---
translation_set_id: stdlib-path-env
path: stdlib/path-env
locale: ko
group: stdlib
group_order: 1
order: 12
title: path와 env: 경로·환경 설정
summary: 호출자 버퍼에 경로와 환경변수를 읽고 용량 오류를 구분합니다.
---

## 경로 조합

```text
std::path::copy
path_join2(dst: ptr<u8>, dst_cap: i32, left: str, right: str) -> i32
path_basename_copy(dst: ptr<u8>, dst_cap: i32, path: str) -> i32
path_dirname_copy(dst: ptr<u8>, dst_cap: i32, path: str) -> i32
```

용량에는 마지막 NUL 공간을 포함합니다. 성공 결과는 NUL을 제외한 길이이며 실패는 -1입니다. 성공했을 때만 목적지를 문자열로 사용합니다. 이 함수들은 경로 문자열을 다루며 파일 존재나 접근 권한을 검사하지 않습니다. 경로를 조합하는 것만으로 디렉터리 탈출을 차단하거나 실제 파일의 동일성을 확인할 수 없습니다.

## 환경변수

```text
std::env::environ
env_get(name: str, dst: ptr<u8>, dst_cap: i64) -> i64
env_exists(name: str) -> bool
env_get_i64(name: str) -> EnvResult<i64>
env_get_i32(name: str) -> EnvResult<i32>
```

`env_get`은 성공 시 NUL을 제외한 길이를 반환합니다. 호출자 버퍼는 NUL까지 담을 수 있어야 합니다. 빈 값은 길이 0으로 성공할 수 있으므로 키가 없는 오류와 다릅니다.

`std::env::consts`에서 NOT_FOUND, NO_SPACE, INVALID_KEY, READ, SOURCE_INCOMPLETE, NO_MEMORY 오류를 가져와 구분합니다. 버퍼 부족을 키 없음으로 취급하지 않습니다. 숫자 조회는 결과의 ok를 확인한 뒤 value를 사용합니다. 환경변수의 내용을 신뢰할 수 있는 설정으로 자동 취급하지 말고 범위와 형식을 검사하십시오.

다음 예제는 data 디렉터리와 input.txt 파일 이름을 조합합니다.

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

실행 결과:

```text
data/input.txt
```

## 디렉터리와 파일 이름 나누기

다음 예제는 경로를 두 버퍼에 나누어 복사합니다. 원본 파일이 실제로 존재할 필요는 없습니다.

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

실행 결과:

```text
directory=data
filename=report.txt
```

버퍼 두 개는 main이 끝날 때까지 유효합니다. `as str`은 새 문자열을 할당하지 않고 같은 버퍼를 문자열로 읽습니다. 따라서 버퍼를 수정하면 해당 주소로 읽는 문자열도 바뀝니다.

## 환경 설정의 기본값 정하기

환경변수는 프로그램 밖에서 전달하는 설정입니다. 숫자 설정을 읽을 때는 “정수로 읽을 수 있는가”와 “이 프로그램에서 허용한 범위인가”를 각각 검사합니다.

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

`WAVE_EXAMPLE_WORKERS`가 없거나 정수로 읽을 수 없으면 기본값 4를 사용합니다. 1~32 사이의 정수가 설정되어 있으면 그 값을 사용하고, 범위 밖 정수면 오류로 종료합니다.

Linux/macOS:

```shell
WAVE_EXAMPLE_WORKERS=8 wavec run main.wave
```

PowerShell:

```powershell
$env:WAVE_EXAMPLE_WORKERS = "8"
wavec run main.wave
```

두 경우 모두 `workers=8`이 출력됩니다. 위 예제는 간단한 기본값 정책을 선택한 것입니다. 필수 설정이라면 숫자 조회 실패를 기본값으로 대체하지 말고 오류로 처리합니다. 키 없음, 버퍼 부족, 읽기 실패를 구분해야 한다면 env_get과 ENV_ERR_* 상수를 사용합니다.
