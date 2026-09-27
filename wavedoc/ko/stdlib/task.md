---
translation_set_id: stdlib-task
path: stdlib/task
locale: ko
group: stdlib
group_order: 1
order: 13
title: task: Future 실행과 정리
summary: 비동기 작업의 단일 소비, 실행과 취소 후 정리를 설명합니다.
---

## 기본 API

`import("std::task" as task);`로 가져옵니다.

```text
block_on<T>(future: Future<T>) -> T
spawn<T>(future: Future<T>) -> Future<T>
cancel<T>(future: Future<T>) -> bool
yield_now() -> Future<void>
sleep_ms(milliseconds: i64) -> Future<void>
cancel_and_join<T>(future: Future<T>) -> Future<void>
shutdown()
```

`task::block_on(pending)`은 Future가 완료될 때까지 실행하고 결과를 반환합니다. 결과 타입은 전달한 Future에서 결정됩니다.

## 수명 규칙

Future는 한 번 소비하는 핸들입니다. 값을 복사해 독립된 두 작업이라고 생각하지 않습니다. `spawn`으로 예약한 작업도 결과를 기다리거나 취소·정리할 책임이 있습니다. `await`와 `block_on`으로 이미 소비한 Future를 다시 소비하지 않습니다.

cancel은 취소 요청이고 자원을 안전하게 정리할 때는 취소 완료를 기다리는 경로가 필요합니다. 비동기 I/O에 빌려 준 메모리를 작업이 접근 중인 상태에서 해제하지 않습니다. 작업을 완료하거나 취소한 뒤 `shutdown`을 호출합니다.

## 협력적 실행

긴 계산과 동기 blocking 호출은 실행기 전체 진행을 늦출 수 있습니다. yield와 비동기 대기는 실행 기회를 양보합니다. 블로킹 I/O가 단지 async 함수 안에 있다는 이유로 비동기가 되지는 않습니다.

[비동기 입문](/docs/ko/language/async-and-never)의 전체 프로그램에서 실행과 정리 순서를 확인하십시오.

## 양보와 완료 대기

다음 프로그램은 작업 중간에 실행 기회를 양보하고 결과를 돌려줍니다. yield는 함수 종료가 아니므로 await 뒤의 코드가 이어서 실행됩니다.

<!-- wave-example: book-task-yield -->
```wave
import("std::task" as task);

async fun calculate() -> i32 {
    println("started");
    await task::yield_now();
    println("resumed");
    return 42;
}

fun main() -> i32 {
    var pending: Future<i32> = calculate();
    var result: i32 = task::block_on(pending);

    println("result={}", result);
    task::shutdown();
    return 0;
}
```

실행 결과:

```text
started
resumed
result=42
```

이 예제에는 다른 작업이 없어 출력 순서가 일정합니다. 여러 작업을 spawn한 프로그램에서는 yield 지점에서 다른 작업이 진행할 수 있으므로 서로 다른 작업의 출력 순서에 의존하지 않습니다.

## 작업을 정리하는 순서

1. 작업에 필요한 저장 공간과 자원을 준비합니다.
2. Future를 만들고 await, block_on 또는 spawn으로 실행합니다.
3. 결과가 필요하면 완료까지 기다립니다.
4. 실행 중인 작업을 취소했다면 그 작업이 정리될 때까지 기다립니다.
5. 작업이 빌린 버퍼와 파일·소켓을 정리합니다.
6. 남은 작업이 없는 상태에서 shutdown을 호출합니다.

Future 변수의 범위를 벗어나는 것과 작업이 안전하게 정리되는 것은 별개입니다. 특히 함수 지역 배열의 주소를 async 작업에 전달했다면, 함수가 반환되기 전에 작업이 그 배열 사용을 마쳐야 합니다.

## async 함수와 일반 함수 나누기

순수 계산은 일반 함수로 분리해도 됩니다. 대기를 표현해야 하는 함수에 async를 붙이고, 그 함수 안에서 await로 완료를 기다립니다. 파일 읽기처럼 오래 걸리는 동기 함수를 async 함수로 감싸는 것만으로 다른 작업에 실행 기회를 주지는 않습니다.

[비동기 학습](/docs/ko/language/async-and-never)에서 Future를 만드는 시점과 실행하는 시점을 비교할 수 있습니다.
