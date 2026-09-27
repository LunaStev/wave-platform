---
translation_set_id: overview
path: getting-started/overview
locale: ko
group: getting-started
group_order: 1
order: 1
title: Wave 문서와 학습 안내
summary: 설치부터 실전 프로그램까지 순서대로 배우고 언어와 표준 라이브러리 규칙을 찾아봅니다.
---

## 문서로 Wave 배우기

Wave 소스를 작성하고, 컴파일해서 실행하고, 결과를 확인하는 과정을 함께 익힙니다. 처음 프로그래밍을 배우는 독자는 아래 순서대로 진행하십시오. 다른 언어를 알고 있다면 각 장의 예제를 실행한 뒤 같은 장의 규칙과 경계값 설명에서 차이를 확인할 수 있습니다.

## 학습 순서

| 단계 | 읽을 문서 | 끝내고 할 수 있는 일 |
| --- | --- | --- |
| 준비 | [설치](/docs/ko/getting-started/install) | 컴파일러와 std를 준비하고 실행 확인 |
| 1 | [첫 프로그램](/docs/ko/language/program-structure) | 파일 작성, 검사, 실행, 종료 코드 이해 |
| 2 | [변수와 타입](/docs/ko/language/declarations-and-types) | 값을 저장하고 범위에 맞는 타입 선택 |
| 3 | [연산과 변환](/docs/ko/language/expressions-and-operators) | 계산 순서와 형 변환의 결과 설명 |
| 4 | [조건과 반복](/docs/ko/language/control-flow) | 조건에 따라 분기하고 데이터를 반복 처리 |
| 5 | [함수](/docs/ko/language/functions-and-generics) | 반복되는 작업을 함수로 분리 |
| 6 | [배열](/docs/ko/language/arrays) | 인덱스로 원소를 읽고 반복 처리 |
| 7 | [문자열](/docs/ko/language/strings) | 문자와 바이트, escape와 문자열 길이 이해 |
| 8 | [구조체와 variant](/docs/ko/language/structures-enums-and-aliases) | 관련 데이터를 묶고 성공·실패 구분 |
| 9 | [포인터와 수명](/docs/ko/language/explicit-memory-type-model) | 주소로 원본을 변경하고 유효 기간 관리 |
| 10 | [동적 메모리](/docs/ko/language/allocation) | 할당 실패를 처리하고 메모리 해제 |
| 11 | [모듈과 제네릭](/docs/ko/language/modules-imports-and-ffi) | 파일을 나누고 타입별 함수 재사용 |
| 12 | [오류 처리](/docs/ko/language/errors) | 결과를 검사하고 실패해도 자원 정리 |
| 13 | [비동기 입문](/docs/ko/language/async-and-never) | Future를 만들고 실행 완료까지 대기 |

## 실습으로 연결하기

기초 과정을 끝내면 [입력 계산기](/docs/ko/practice/input-calculator), [파일 읽기](/docs/ko/practice/file-reader), [바이너리 메시지](/docs/ko/practice/binary-message), [TCP 클라이언트](/docs/ko/practice/tcp-client)를 만듭니다. 실습마다 정상 입력과 실패 상황을 함께 확인합니다.

## 문서의 세 탭

- **Wave**: 처음부터 순서대로 읽는 언어 학습 과정과 실전 프로젝트입니다.
- **[표준 라이브러리](/docs/ko/stdlib)**: 모듈별 API, 반환값, 오류, 소유권과 플랫폼 조건입니다.
- **[Whale](/docs/ko/whale)**: 빌드·링크, 패키지 관리, 명령 사용법과 저수준 툴체인 문서입니다.

학습 예제는 전체 프로그램과 함수 내부 코드 조각을 구분합니다. `wavec`는 터미널에서 실행하고, `wave` 코드 블록은 `.wave` 파일에 저장합니다. 입력과 출력은 같은 것이 아닙니다. 표준 입력이 필요한 예제는 입력할 값도 따로 표시합니다.

## 막혔을 때

[문제 해결](/docs/ko/reference/diagnostics)에서 설치, 소스 검사, 링크, 실행 단계를 구분하십시오. 정확한 규칙은 [언어 요약](/docs/ko/reference/syntax-quick-reference), 명령은 [컴파일러 참조](/docs/ko/getting-started/compiler), 함수는 [표준 라이브러리 안내](/docs/ko/reference/standard-library)에서 찾습니다.
