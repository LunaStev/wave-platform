---
translation_set_id: diagnostics
path: reference/diagnostics
locale: ko
group: reference
group_order: 5
order: 2
title: 문제 해결: 설치부터 실행까지
summary: 실패한 단계를 구분하고 재현 가능한 정보로 원인을 좁힙니다.
---

## 먼저 실패 단계를 구분하기

| 관찰한 현상 | 먼저 확인할 것 | 다음 조치 |
| --- | --- | --- |
| wavec 명령을 찾지 못함 | PATH와 실행 파일 위치 | 절대 경로로 실행한 뒤 PATH 설정 |
| 실행에 필요한 파일을 찾지 못함 | 설치 폴더에서 파일이 빠졌는지 | 패키지 전체를 다시 풀어 설치 |
| std import 실패 | `wavec print std-path` | 대응 std 설치 또는 `--std-root` 지정 |
| 소스 위치와 타입 오류 출력 | `wavec check main.wave` | 첫 오류부터 고치고 다시 검사 |
| 다른 OS·CPU 대상의 빌드 실패 | 지정한 target과 대상 환경 | [크로스 빌드 설정](/docs/ko/whale/build-link-targets) 확인 |
| 빌드 성공 후 실행 실패 | 종료 코드, 입력, 작업 디렉터리 | 실행 환경과 API 오류 검사 |

## 작은 진단 예제

다음은 의도적으로 잘못된 전체 프로그램입니다.

```wave
fun main() {
    var count: i32 = 1;
    println("{}", missing);
}
```

`wavec check main.wave`는 선언되지 않은 이름 missing을 가리켜야 합니다. 변수 이름을 count로 고친 뒤 다시 검사하고 실행합니다. 진단 문구 전체보다 파일·위치·원인에 집중하십시오. 뒤따르는 오류가 최초 오류의 결과일 수 있습니다.

## 실행 파일이 실패할 때

Linux/macOS 셸에서는 실행 직후 `echo $?`, PowerShell에서는 `$LASTEXITCODE`로 종료 코드를 확인합니다. 입력 오류와 명시적 `return 1`은 같은 원인이 아닙니다. 시프트 횟수나 실수 변환의 잘못된 런타임 값은 trap을 일으킬 수 있습니다. [연산 규칙](/docs/ko/language/expressions-and-operators)을 확인하십시오.

상대 파일 경로는 소스 파일 위치가 아니라 실행 작업 디렉터리에 영향을 받습니다. 파일 읽기 실패를 문자열 길이 0으로 처리하지 말고 반환 오류를 먼저 확인합니다. 네트워크 연결 실패는 주소 조회, 서버 대기, 권한과 제한 시간을 나누어 확인합니다.

## 문제 보고에 필요한 정보

1. `wavec --version` 출력과 실행한 정확한 명령.
2. 호스트 OS·아키텍처와 별도로 지정한 target.
3. 선택한 std 경로와 함께 사용한 컴파일러 출처.
4. 문제를 재현하는 최소 소스, 입력과 필요한 파일.
5. 예상 결과, 실제 결과, 진단과 종료 코드.

비밀번호·토큰·개인 파일 내용은 제거합니다. 최소 예제를 줄일 때 문제가 없어지면 마지막으로 제거한 요소가 단서입니다. `--error-format=json`은 도구가 진단을 수집할 때 사용할 수 있습니다.

[설치](/docs/ko/getting-started/install) · [컴파일러 명령](/docs/ko/getting-started/compiler) · [타깃과 링크](/docs/ko/whale/build-link-targets)
