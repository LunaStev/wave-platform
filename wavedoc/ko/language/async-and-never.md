---
translation_set_id: language-async-and-never
path: language/async-and-never
locale: ko
group: language
group_order: 2
order: 13
title: 13. 비동기 함수와 Future
summary: 지연 실행되는 Future와 await, block_on의 역할을 배웁니다.
---

## 기다리는 작업을 표현하기

파일·소켓·타이머처럼 기다리는 작업에서는 계산을 계속하는 것과 완료를 기다리는 것을 구분할 필요가 있습니다. 비동기 함수는 완료될 결과를 Future로 표현합니다. async를 붙였다고 자동으로 새 스레드를 만들거나 모든 동기 호출이 비동기로 바뀌지는 않습니다.

이 장은 함수, 포인터와 오류 처리 다음에 읽으십시오. 예제는 `std::task` 실행기를 사용하는 네이티브 프로그램이며 각 파일을 `wavec run main.wave`로 실행합니다.

## Future를 만들고 결과 받기

<!-- wave-example: book-async-basic -->
```wave
import("std::task" as task);

async fun calculate(value: i64) -> i64 {
    return value * 2;
}

fun main() {
    var pending: Future<i64> = calculate(21);
    var result: i64 = task::block_on(pending);

    task::shutdown();
    println("{}", result);
}
```

실행 결과:

```text
42
```

calculate의 선언에 적은 i64는 완료한 뒤 얻는 값입니다. 호출 자체의 결과는 Future<i64>입니다. 일반 main에서는 block_on으로 Future를 구동하고 완료된 결과를 받습니다.

shutdown은 작업을 모두 정리한 뒤 실행기 자원을 해제하는 단계입니다. 아직 접근 중인 버퍼를 먼저 해제하거나 작업이 남은 상태를 무시하지 않습니다.

## 호출과 본문 실행이 다름

비동기 함수는 지연 실행됩니다. 호출한 즉시 본문을 끝까지 수행하는 일반 함수와 구분해야 합니다.

<!-- wave-example: book-async-lazy -->
```wave
import("std::task" as task);

static entered: i32 = 0;

async fun work() -> i32 {
    entered += 1;
    return 7;
}

fun main() {
    var pending: Future<i32> = work();

    println("before={}", entered);

    var result: i32 = task::block_on(pending);

    task::shutdown();
    println("after={} result={}", entered, result);
}
```

실행 결과:

```text
before=0
after=1 result=7
```

Future를 만들었을 때 entered는 아직 0입니다. 실행을 구동한 뒤 본문이 수행되어 1이 됩니다. 단지 Future를 변수에 저장했다고 작업이 완료된 것으로 처리하면 안 되는 이유입니다.

## 비동기 함수 안에서 기다리기

async 함수 안에서는 await로 다른 Future의 완료를 기다립니다. await한 식의 결과는 완료값입니다.

<!-- wave-example: book-async-nested -->
```wave
import("std::task" as task);

async fun twice(value: i32) -> i32 {
    await task::yield_now();
    return value * 2;
}

async fun process() -> i32 {
    var value: i32 = await twice(20);
    return value + 2;
}

fun main() {
    var answer: i32 = task::block_on(process());

    task::shutdown();
    println("{}", answer);
}
```

실행 결과:

```text
42
```

process는 twice의 Future를 기다린 다음 완료값에 2를 더합니다. 일반 함수의 호출 결과와 비슷하게 값을 사용할 수 있지만 기다리는 동안 실행 기회를 다른 작업에 넘길 수 있습니다.

yield_now는 협력적으로 실행 기회를 양보합니다. 긴 계산 루프에서 한 번도 양보하지 않으면 다른 작업 진행이 늦어질 수 있습니다. 비동기는 CPU 계산을 자동으로 병렬 분산하는 장치가 아닙니다.

## 여러 작업 예약하기

spawn으로 작업을 예약하고 각각의 결과를 기다릴 수 있습니다. 두 작업의 중간 출력 순서에 의존하지 않고 최종 결과를 확인하는 예제입니다.

<!-- wave-example: book-async-spawn -->
```wave
import("std::task" as task);

async fun compute(value: i32) -> i32 {
    await task::yield_now();
    return value * 2;
}

async fun combine() -> i32 {
    var first: Future<i32> = task::spawn(compute(10));
    var second: Future<i32> = task::spawn(compute(20));

    var left: i32 = await first;
    var right: i32 = await second;

    return left + right;
}

fun main() {
    var total: i32 = task::block_on(combine());

    task::shutdown();
    println("{}", total);
}
```

실행 결과:

```text
60
```

예약한 작업마다 결과를 기다리는 위치가 있습니다. 작업을 만들고 핸들을 잊어버리는 대신 누가 완료를 확인할지 정해야 합니다. await 순서와 내부 작업이 실행되는 모든 순서를 같은 것으로 생각하지 마십시오.

## Future는 한 번 소비하기

Future는 단일 소비 핸들로 다룹니다. 같은 Future를 복사해서 서로 다른 두 작업처럼 기다리지 않습니다. 이미 완료를 받은 Future를 다시 block_on하거나 await하지 않습니다.

여러 곳에서 같은 결과가 필요하면 Future를 여러 번 소비하는 대신 완료된 값을 저장하고 그 값의 복사·공유 규칙에 따라 전달합니다. 값 안에 포인터나 소유 자원이 있는지도 함께 확인해야 합니다.

## 타이머와 동기 대기의 차이

async 함수 안에서 기다릴 때는 `await task::sleep_ms(...)`를 사용할 수 있습니다. 동기 sleep을 호출하면 현재 실행 흐름을 막아 실행기의 다른 작업 진행에도 영향을 줄 수 있습니다.

기다린 시간이 정확히 요청한 밀리초와 같다고 기대하지 않습니다. 스케줄링과 다른 작업에 따라 늦게 깨어날 수 있습니다. 시간 제한을 구현할 때는 경과 시간을 측정하는 clock과 deadline을 사용하며, 매번 원래 전체 시간을 다시 기다리는 방식과 구분합니다.

## 버퍼 수명과 취소

비동기 I/O가 버퍼 주소를 받았다면 호출한 함수가 잠시 멈춘 동안에도 그 버퍼는 살아 있어야 합니다. 완료나 취소 정리가 끝나기 전에 해제·재할당하면 작업이 더 이상 유효하지 않은 주소에 접근할 수 있습니다.

취소 요청과 작업 정리 완료는 같은 순간이라고 단정할 수 없습니다. cancel 계열 API의 결과를 확인하고, 자원을 해제하기 전에 필요한 완료 대기를 수행합니다. 상세 호출 규칙은 [task 참조](/docs/ko/stdlib/task)를 읽으십시오.

## 흔한 오해

| 생각 | 실제로 확인할 것 |
| --- | --- |
| async 호출을 했으니 끝났음 | Future를 실제로 구동하고 완료했는가 |
| async 함수 안의 모든 호출은 비동기 | 호출한 API가 동기인지 비동기인지 |
| Future 복사는 작업 복제 | 동일 핸들을 중복 소비하고 있지 않은가 |
| 취소했으니 즉시 버퍼 해제 가능 | 취소 후 작업 정리까지 끝났는가 |
| 중간 출력 순서는 항상 고정 | 결과에 필요한 순서만 명시적으로 기다리는가 |

## 연습과 전체 풀이

비동기 함수 세 개를 연속으로 기다리는 파이프라인을 만드십시오. 두 배 계산 뒤 5를 더한 결과를 반환합니다.

<!-- wave-example: book-async-solution -->
```wave
import("std::task" as task);

async fun read_value() -> i32 {
    return 10;
}

async fun transform(value: i32) -> i32 {
    await task::yield_now();
    return value * 2;
}

async fun pipeline() -> i32 {
    var initial: i32 = await read_value();
    var changed: i32 = await transform(initial);

    return changed + 5;
}

fun main() {
    var result: i32 = task::block_on(pipeline());

    task::shutdown();
    println("{}", result);
}
```

실행 결과:

```text
25
```

이 예제는 의도적으로 순차적인 의존 관계입니다. transform은 read_value의 결과가 필요하므로 단순히 모두 spawn한다고 그 관계가 없어지지 않습니다. 독립적인 작업과 결과가 필요한 작업을 구분하는 것이 비동기 설계의 출발점입니다.

기초 학습을 마쳤다면 [파일 읽기 실습](/docs/ko/practice/file-reader)과 [TCP 실습](/docs/ko/practice/tcp-client)으로 실제 외부 자원과 연결하십시오.

## void와 never

반환 타입을 생략한 일반 함수는 값 없이 호출 지점으로 돌아올 수 있습니다. never 타입은 `!`로 표기하며 호출 지점으로 정상적으로 돌아오지 않는다는 뜻입니다. 프로세스 종료 함수가 대표적입니다.

선언을 설명하는 예:

```text
fun log(message: str)              // 작업 후 돌아옴, 결과값 없음
fun stop(code: i32) -> !           // 정상 반환하지 않음
async fun work() -> i64            // 호출 결과는 Future<i64>
```

never를 일반적인 저장 값으로 만들려고 하지 않습니다. 돌아오지 않는다고 선언한 함수가 정상 반환 경로를 가지도록 작성하지 마십시오. 종료하기 전에 자원 정리가 필요하다면 호출자가 먼저 수행해야 합니다.

## 반환하지 않는 함수의 전체 예제

main.wave로 저장해 실행하면 출력 없이 종료 코드 0으로 끝납니다. stop은 호출자에게 반환하지 않으므로 `-> !`로 선언합니다.

<!-- wave-example: never-function -->
```wave
import("std::process::core")::{
    proc_exit
};

fun stop() -> ! {
    proc_exit(0);
}

fun main() {
    stop();
}
```
