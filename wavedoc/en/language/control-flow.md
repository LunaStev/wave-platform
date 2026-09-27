---
translation_set_id: control-flow
path: language/control-flow
locale: en
group: language
group_order: 2
order: 4
title: 4. Conditions, loops, and boundary values
summary: Learn if, for, while and the range of repetitive statements.
---

## Select execution path

The program in the previous chapter executed statements from top to bottom. A real program must do different things depending on the input and state. Conditional statements select an execution path, and loops apply the same rule to multiple values.

Save each example in main.wave and run it. When reading code, write down on a piece of paper the current variable values, the conditions to be tested next, and the order of statements to be executed. It is more important to practice following the flow than to memorize the results.

## if and else

Conditions are written in parentheses, and the text is enclosed in curly brackets. In the following example, replace balance with 500 or 2000 to determine which branch is executed.

<!-- wave-example: book-if-else -->
```wave
fun main() {
    var balance: i32 = 2000;
    var price: i32 = 1200;

    if (balance >= price) {
        balance -= price;
        println("bought, balance={}", balance);
    } else {
        println("not enough money");
    }
}
```

Execution result:

```text
bought, balance=800
```

It's not executing both blocks. If the condition is true, the first block is executed. If the condition is false, the else block is executed. Purchases are allowed even when balance is equal to price, so I wrote `>=`. If you change it to `>`, the operation will change for the same amount.

## Order of multiple conditions

You can concatenate conditions with else if. Since we execute one satisfied branch from above first, we need to consider whether we want to check the larger boundary or the smaller boundary first.

<!-- wave-example: book-grade -->
```wave
fun grade(score: i32) {
    if (score < 0 || score > 100) {
        println("invalid");
    } else if (score >= 90) {
        println("A");
    } else if (score >= 80) {
        println("B");
    } else {
        println("C");
    }
}

fun main() {
    grade(95);
    grade(80);
    grade(79);
    grade(101);
}
```

Execution result:

```text
A
B
C
invalid
```

Invalid scores are first rejected and then graded. If you put `score >= 80` at the beginning, you will never reach the A branch because 95 degrees goes into that branch. Check the order of the conditions as well as the correctness of each condition.

## Do not change values in conditional expressions

Assignment, compound assignment, and increment or decrement operations are not allowed in if, while, or for conditions. Use `==` for comparison. To update a value and then test it, write two separate statements.

Correct form of a fragment inside a function:

```wave
value = read_value();

if (value == expected) {
    println("matched");
}
```

read_value and expected are not defined in this fragment, so it is not the full program that runs as is. The rule shown here is “Compare after state change.”

## while: While conditions hold

<!-- wave-example: book-while-countdown -->
```wave
fun main() {
    var remaining: i32 = 3;

    while (remaining > 0) {
        println("{}", remaining);
        remaining -= 1;
    }

    println("finished at {}", remaining);
}
```

Execution result:

```text
3
2
1
finished at 0
```

Conditions are checked before entering the body. If remaining is 0 from the beginning, the body is never executed. If the decrement at the end of the body is omitted, the condition remains true and the loop does not end.

After writing the loop, check “What brings it closer to the termination condition?” If it is a loop waiting for input, input change or EOF plays its role, and if it is a numeric loop, index update plays its role.

## for: Initialization/Conditions/Update

for expresses the three parts necessary for repetition.

<!-- wave-example: book-for-sum -->
```wave
fun main() {
    var total: i32 = 0;

    for (var number: i32 = 1; number <= 5; number += 1) {
        total += number;
    }

    println("sum={}", total);
}
```

Execution result:

```text
sum=15
```

1. Initialize number to 1. This step is one time.
2. Checks number <= 5. If false, end the iteration.
3. Add number to total in the text.
4. Increase number by 1 and return to condition checking.

It is not assumed that repetition variables declared within for can be used after repetition. If your design requires a value after iteration, declare it outside and make the initialization location clear.

## Inclusion and Exclusion Boundaries

The natural sum of 1 to 5 is `<= 5`. On the other hand, the index of an array of length 5 should use `< 5`. This is because array indices start at 0 and end at 4.

Do not confuse “run five times” with “up to and including the value 5.” You can determine the number of repetitions by looking at the start and end values ​​together. Cases where the input is empty and has only one element are good for detecting boundary mistakes.

## continue and break

continue skips the rest of this iteration, and break ends the nearest iteration.

<!-- wave-example: book-loop-control -->
```wave
fun main() {
    var total: i32 = 0;

    for (var number: i32 = 1; number <= 10; number += 1) {
        if (number == 3) {
            continue;
        }

        if (number == 6) {
            break;
        }

        total += number;
    }

    println("{}", total);
}
```

Execution result:

```text
12
```

The numbers in the total are 1, 2, 4, and 5. If continue is encountered in for, the renewal process will proceed. Since while does not have a separate update expression like for, you must be careful not to omit any state changes required before continue.

If loops are nested, one break does not complete all loops. If you need to stop at multiple stages, enclose the work in a function and indicate your intention by using return or by checking the termination condition in the outer iteration as well.

## Divide the case by match

You can use match when comparing multiple cases of the same value. The body of each arm is a block.

<!-- wave-example: book-match-number -->
```wave
fun describe(status: i32) {
    match (status) {
        200 => {
            println("ok");
        }
        404 => {
            println("missing");
        }
        _ => {
            println("other");
        }
    }
}

fun main() {
    describe(200);
    describe(404);
    describe(500);
}
```

Execution result:

```text
ok
missing
other
```

`_` is the pattern that handles the rest. Do not place duplicates within the same match. variant, which has different data types depending on the value, is covered in [Data model chapter](/docs/en/language/structures-enums-and-aliases).

## Complete example: Counting numbers that meet conditions

Find the number and sum of even numbers from 1 to 10. Since count and sum are different information, they are accumulated as variables.

<!-- wave-example: book-loop-statistics -->
```wave
fun main() {
    var count: i32 = 0;
    var total: i32 = 0;

    for (var number: i32 = 1; number <= 10; number += 1) {
        if (number % 2 == 0) {
            count += 1;
            total += number;
        }
    }

    println("count={} total={}", count, total);
}
```

Execution result:

```text
count=5 total=30
```

The even numbers are 2, 4, 6, 8, and 10, so the number is 5 and the sum is 30. Even if the expression is short, it is easy to verify the iteration boundary if you first check the result in a small range that can be obtained by hand.

## Exercise and complete solution

Add only multiples of 3 from 1 to 20, but do not add any values that add up to more than 30. We need to distinguish between “adding and then checking to see if it is over” and “checking to see if it is over and then adding.”

<!-- wave-example: book-loop-solution -->
```wave
fun main() {
    var total: i32 = 0;

    for (var number: i32 = 1; number <= 20; number += 1) {
        if (number % 3 != 0) {
            continue;
        }

        if (total + number > 30) {
            break;
        }

        total += number;
    }

    println("{}", total);
}
```

Execution result:

```text
30
```

The values 3+6+9+12 sum to 30, so the next value, 15, is not added. These small inputs are safe, but for large integers the check `total + number` can itself overflow. Writing a check does not automatically handle every boundary case.
