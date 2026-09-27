---
translation_set_id: functions
path: language/functions-and-generics
locale: en
group: language
group_order: 2
order: 5
title: 5. Designing and composing functions
summary: Learn parameters, return values, default values, and passing by value.
---

## Starting from repetitive code

Functions are a tool for reducing syntax, but they are also a tool for demarcating tasks. Separating what it takes as input, what it computes, and what results it returns allows you to understand your program in smaller chunks.

In this chapter, we start with a program that writes discount calculations multiple times. Each example is complete in main.wave and run as `wavec run main.wave`. We are not dividing the files yet.

<!-- wave-example: book-function-before -->
```wave
fun main() {
    var first_price: i32 = 2000;
    var first_discount: i32 = first_price * 10 / 100;
    var first_total: i32 = first_price - first_discount;
    var second_price: i32 = 5000;
    var second_discount: i32 = second_price * 10 / 100;
    var second_total: i32 = second_price - second_discount;
    println("{} {}", first_total, second_total);
}
```

Execution result:

```text
1800 4500
```

The two calculations differ only in price and have the same structure. When changing discount rules, you'll need to edit both places. If you change just one side, you will get different results for products that require the same policy.

## Determine input and output

Move redundant calculations to functions. The changed value is received as input price, and the calculated price is returned.

<!-- wave-example: book-function-extract -->
```wave
fun discounted(price: i32) -> i32 {
    var discount: i32 = price * 10 / 100;
    return price - discount;
}

fun main() {
    println("{} {}", discounted(2000), discounted(5000));
}
```

Execution result:

```text
1800 4500
```

The function name is discounted. `price: i32` in parentheses is the parameter declaration and `-> i32` is the result type. The local variable discount in the body is used only within this function.

`discounted(2000)` is an expression that calls a function. The 2000 in parentheses is the argument you actually pass. Since the value returned by the function becomes the result of this calling expression, it can be directly used as an argument for println.

|terminology|code|meaning|
| --- | --- | --- |
|parameter| price |Input name specified when declaring the function|
|factor| 2000 |Value passed when calling|
|return type| i32 |Type of value produced by the call expression|
|return statement| return price - discount |Pass the result and end this call|

## Follow the calling and execution order

Just writing a function declaration in the source does not cause its body to be executed immediately. Executes when the point called in main is reached.

<!-- wave-example: book-function-trace -->
```wave
fun calculate(value: i32) -> i32 {
    println("inside: {}", value);
    return value * 2;
}

fun main() {
    println("before");
    var result: i32 = calculate(7);
    println("after: {}", result);
}
```

Execution result:

```text
before
inside: 7
after: 14
```

The order of progress is the first output of main, the main text of calculate, and the last output of main. When `return` is executed, the call to calculate ends with result 14, and the initialization of result of main is complete.

If multiple function calls are mixed in one expression, and the order of side effects is important, break the calls into separate statements. The examples in this chapter also store the results of calls that need to be traced in local variables.

## multiple parameters

If you also receive a discount rate as input, you can calculate multiple policies with the same function.

<!-- wave-example: book-function-parameters -->
```wave
fun discounted(price: i32, percent: i32) -> i32 {
    return price - price * percent / 100;
}

fun main() {
    var standard: i32 = discounted(2000, 10);
    var special: i32 = discounted(2000, 25);
    println("standard={}", standard);
    println("special={}", special);
}
```

Execution result:

```text
standard=1800
special=1500
```

The order of arguments must match the declaration. If both parameters are i32, it is difficult to distinguish their meaning through type checking alone, even if the order is changed. Clearly define the function name and parameter name, and write the call location in an easy-to-read manner.

This function assumes a small amount and a valid discount rate. Does not handle negative prices, ratios greater than 100, or mid-multiplication overflow. When creating a function, you must describe not only the body but also the input conditions. The later completion program separates the inspection steps.

## default argument

You can provide frequently used values as defaults.

<!-- wave-example: book-function-default -->
```wave
fun discounted(price: i32, percent: i32 = 10) -> i32 {
    return price - price * percent / 100;
}

fun main() {
    println("{}", discounted(2000));
    println("{}", discounted(2000, 25));
}
```

Execution result:

```text
1800
1500
```

The first call omits the second argument and uses 10. The second call uses the specified 25. Default values ​​are left for optional parameters at the end. It is not used as a grammar to write a blank space in order to omit only the first argument.

Changing the default changes the behavior of skip calls. The default values ​​of public functions are also part of the behavior you rely on. This is why calls with specified arguments and calls with omitted arguments are tested separately.

## means passing by value

Passing an integer value separates the value the function receives from the caller's variable storage. Computing a result does not automatically change the caller's variables.

<!-- wave-example: book-function-value -->
```wave
fun next(value: i32) -> i32 {
    return value + 1;
}

fun main() {
    var count: i32 = 4;
    var later: i32 = next(count);
    println("count={} later={}", count, later);
    count = next(count);
    println("count={}", count);
}
```

Execution result:

```text
count=4 later=5
count=5
```

The first call reads count and initializes later. count is still 4. After the second call, the result is assigned to count, so it becomes 5. The return-by-value design makes it obvious to the caller where the data has changed.

If you want to change the original storage space within a function, you can pass a pointer. This is covered in [pointer chapter](/docs/en/language/explicit-memory-type-model). Even if you pass a pointer, you must distinguish between the pointer value itself and the storage space for its address.

## Don't leave paths that don't return

A function that returns a value must provide results on all paths required. It doesn't leave the final path like this:

```wave
// 의도적으로 잘못된 함수 조각

fun positive(value: i32) -> i32 {
    if (value > 0) {
        return value;
    }
}
```

If value is less than or equal to 0, there is no set value to return. Once you have established your intended rules, you need to write out all your routes.

<!-- wave-example: book-function-paths -->
```wave
fun positive_or_zero(value: i32) -> i32 {
    if (value > 0) {
        return value;
    }

    return 0;
}

fun main() {
    println("{} {} {}", positive_or_zero(-2), positive_or_zero(0), positive_or_zero(8));
}
```

Execution result:

```text
0 0 8
```

The call that first executes return does not proceed to return below. The last return is reached only when the condition is false. We checked the rule for three cases: positive, 0, and negative.

## Function with no result

If you only perform operations such as printing, you can omit the return type. A function without a result can also be terminated early with `return;`.

<!-- wave-example: book-function-void -->
```wave
fun show_positive(value: i32) {
    if (value <= 0) {
        return;
    }

    println("positive={}", value);
}

fun main() {
    show_positive(-3);
    show_positive(6);
}
```

Execution result:

```text
positive=6
```

The first call returns without printing anything. The second call prints: “No result” and “Does not return to call point” are different. Functions that do not return, such as functions that terminate a process, are classified by their return type `!`.

## Composing a complete program with multiple functions

Now we split the input validation, calculation, and output into different functions. Since we have restricted the price range, the intermediate multiplication in this example is within the range i32.

<!-- wave-example: book-function-program -->
```wave
fun valid_order(price: i32, quantity: i32, percent: i32) -> bool {
    return price >= 0 && price <= 100000
        && quantity >= 1 && quantity <= 100
        && percent >= 0 && percent <= 100;
}

fun discounted_unit(price: i32, percent: i32) -> i32 {
    return price - price * percent / 100;
}

fun order_total(price: i32, quantity: i32, percent: i32) -> i32 {
    var unit: i32 = discounted_unit(price, percent);
    return unit * quantity;
}

fun show_order(price: i32, quantity: i32, percent: i32) {
    if (!valid_order(price, quantity, percent)) {
        println("invalid order");
        return;
    }

    var total: i32 = order_total(price, quantity, percent);
    println("total={}", total);
}

fun main() {
    show_order(2000, 3, 10);
    show_order(2000, 0, 10);
    show_order(2000, 3, 120);
}
```

Execution result:

```text
total=5400
invalid order
invalid order
```

Answer one question for each function. valid_order determines whether input is allowed, discounted_unit determines how much one discount is, order_total determines how much the total price is, and show_order determines what to show.

Smaller functions are not necessarily better. Giving each expression a name can make it difficult to keep track. Separate rules when they are meaningful to be reused elsewhere or when there are rules to be explained and verified independently.

## practice problems

1. Write maximum, which returns the larger of two integers.
2. Write a function clamp that returns a value between two boundaries. In the solution below, low <= high is set as the calling condition.
3. Create a function that adds tax to a value and combine it with the order sum function. Determine the range and integer cutting points first.

### Solution: Function with boundary

<!-- wave-example: book-function-solution -->
```wave
fun maximum(left: i32, right: i32) -> i32 {
    if (left > right) {
        return left;
    }

    return right;
}

fun clamp(value: i32, low: i32, high: i32) -> i32 {
    if (value < low) {
        return low;
    }
    if (value > high) {
        return high;
    }

    return value;
}

fun main() {
    println("max={}", maximum(7, 4));
    println("{} {} {}", clamp(-3, 0, 10), clamp(6, 0, 10), clamp(20, 0, 10));
}
```

Execution result:

```text
max=7
0 6 10
```

clamp has three paths: less than range, within range, and greater than range. Check by manually adding the boundary values ​​0 and 10 as well. If you want to apply up to low > high, you must decide how to express failure. [Error handling chapter](/docs/en/language/errors) addresses this issue using result structures.


## recursive call

A function can call itself. Recursion requires an exit condition where no more calls are made, and a process in which each call gets closer to that condition.

<!-- wave-example: book-ref-recursion -->
```wave
fun factorial(value: i32) -> i32 {
    if (value <= 1) {
        return 1;
    }

    return value * factorial(value - 1);
}

fun main() {
    println("{}", factorial(5));
}
```

Execution result:

```text
120
```

The calculation of 5 leads to `5 * factorial(4)`, which returns 1 when it is reached. The return value is passed sequentially to the previous call, resulting in 120. This function is an example to explain small positive integers. With large inputs, you need to consider the result range and call depth. You can avoid the problem of increased call depth by writing the same operation in a loop.

## Common errors

|phenomenon|Check|
| --- | --- |
|Error with insufficient or too many arguments|Number of parameters and optional default values|
|Return type does not match|return Expression type and function declaration|
|Not returning from a specific path|Does it return even if the condition is false?|
|Missing generic type argument|After the function name `<Type>`|
|The original is changed after a pointer function is called.|Is it a function that only reads the value or a function that modifies it?|

Externally exported functions, such as `export(c)`, use specific signatures. Generic functions themselves cannot be exported with an external calling convention. `ptr<T>` and `array<T, N>` are the language's built-in memory types and distinguish them from user generic structure declarations.

[function learning](/docs/en/language/functions-and-generics) · [Learning modules and generics](/docs/en/language/modules-imports-and-ffi) · [FFI](/docs/en/language/modules-imports-and-ffi)
