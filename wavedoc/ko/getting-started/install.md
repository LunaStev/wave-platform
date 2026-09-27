---
translation_set_id: install
path: getting-started/install
locale: ko
group: getting-started
group_order: 1
order: 2
title: Wave 설치
summary: Linux, macOS, Windows에 Wave를 설치하고 첫 프로그램을 실행합니다.
---

## Linux와 macOS

터미널에서 다음 명령을 실행합니다. Wave와 패키지 관리자 Vex가 함께 설치됩니다.

```shell
curl -fsSL https://wave-lang.dev/install.sh | bash -s -- latest
```

설치가 끝나면 새 터미널을 열고 확인합니다.

```shell
wavec --version
```

## Windows

PowerShell에서 다음 명령을 실행합니다. Wave와 패키지 관리자 Vex가 함께 설치됩니다.

```powershell
irm https://wave-lang.dev/install.ps1 -OutFile install.ps1
powershell -ExecutionPolicy Bypass -File .\install.ps1 -Latest
```

설치가 끝나면 새 PowerShell 창을 열고 확인합니다.

```powershell
wavec --version
vex --version
```

## 첫 실행

다음 내용을 `main.wave`로 저장합니다.

<!-- wave-example: install-stdlib -->
```wave
import("std::string::len")::{
    len
};

fun main() {
    println("Wave: {} bytes", len("Wave"));
}
```

파일을 저장한 디렉터리에서 실행합니다.

```shell
wavec run main.wave
```

실행 결과:

```text
Wave: 4 bytes
```

표준 라이브러리를 찾을 수 없다는 메시지가 나오면 설치한 뒤 다시 실행합니다.

```shell
wavec install std
wavec run main.wave
```

[다음: 첫 프로그램](/docs/ko/language/program-structure) · [문제 해결](/docs/ko/reference/diagnostics)
