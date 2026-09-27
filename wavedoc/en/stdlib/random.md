---
translation_set_id: stdlib-random
path: stdlib/random
locale: en
group: stdlib
group_order: 1
order: 10
title: random: Filling buffers with OS randomness
summary: OS Fill buffer with entropy and handle partial failures.
---

## API

```text
std::random::fill
random_available() -> bool
random_fill(buffer: ptr<u8>, size: i64) -> RandomFillResult
```

size is a byte count, and the caller provides the storage. `RandomFillResult` contains ok, written, and error. On success, written equals the requested length. On failure, written identifies the valid, filled prefix; do not use the remaining bytes as random data.

`random_available` tells you whether OS random number function is supported. The success of an individual request is checked by the results of random_fill. Uses only OS entropy and does not fall back on time value or weak PRNG on failure. size=0 will succeed even if passed along with null. null is an error for negative or positive lengths.

## Running example

<!-- wave-example: random-api -->
```wave
import("std::random::fill")::{
    RandomFillResult, random_fill
};

fun main() -> i32 {
    var bytes: array<u8, 16>;
    var result: RandomFillResult = random_fill(&bytes[0], 16);
    if (!result.ok) {
        println("random error={}", result.error);
        return 1;
    }

    println("filled={}", result.written);
    return 0;
}
```

Execution result:

```text
filled=16
```

Save it as `main.wave` and run it. The byte contents are different each time, so no specific value is expected. If it fails, check the cause with result.error. Rather than output random bytes literally, use separate encoding if necessary.

## If your request is incorrect

A zero-byte request succeeds because nothing needs to be written. Passing null with a positive length fails because there is no destination buffer. The following program compares these cases without allocating memory.

<!-- wave-example: book-random-boundaries -->
```wave
import("std::random::fill")::{
    RandomFillResult,
    random_fill
};

fun main() -> i32 {
    var empty: RandomFillResult = random_fill(null, 0);

    if (!empty.ok || empty.written != 0) {
        return 1;
    }

    println("empty request succeeded");

    var invalid: RandomFillResult = random_fill(null, 16);

    if (invalid.ok) {
        return 2;
    }

    println("missing buffer rejected");
    return 0;
}
```

Execution result:

```text
empty request succeeded
missing buffer rejected
```

## Handling partially filled buffers

If 16 bytes were requested but failed and written=8 is returned, only the first 8 bytes are filled. If the task is to create a 16-byte identifier, it is not a successful identifier, so we discard the entire result and report a failure. You should not fill the remaining 8 bytes with 0 and then treat it as success.

Mapping random bytes to an integer range requires care. Applying `% 10` to uniformly distributed u8 values makes 0–5 more likely than 6–9, because 256 is not divisible by 10. To remove this bias, reject values 250–255, draw again, and apply the remainder operation only to accepted values.

The storage and lifetime of the random bytes are managed by the caller. When using an array, it is processed within the scope of the array, and when using dynamic memory, it is freed after use. The return structure does not own the buffer instead.
