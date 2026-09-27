---
translation_set_id: data-types
path: language/structures-enums-and-aliases
locale: en
group: language
group_order: 2
order: 8
title: 8. Structs, enums, and variants
summary: Learn about fields, structure initialization, and the roles of enum and variant.
---

## Expressing relationships between data as types

If the price and number of products are each passed as variables, it is difficult to tell by just looking at the code whether the two values belong to the same product. A structure groups related fields. enum represents a named state, and variant stores different data together for each case.

The three functions are not syntaxes that replace each other. Choose depending on what you want to express.

|something to express|select|yes|
| --- | --- | --- |
|Multiple fields existing simultaneously| struct |Unit price and quantity of product|
|Named integer state| enum |Wait/Proceed/Complete|
|Different data in each case| variant |success value or error|

## Structure declaration and value creation

<!-- wave-example: book-struct-first -->
```wave
struct Product {
    price: i32;
    quantity: i32;
}

fun main() {
    var item: Product = Product {
        price: 1500,
        quantity: 2
    };

    println("{} {}", item.price, item.quantity);
}
```

Execution result:

```text
1500 2
```

Fields in declarations end with a semicolon, and when creating values, fields and values are connected with colons and separated by commas. Defining the type and creating the actual value are two different steps. Declaring the type Product does not automatically create storage space for one product.

The field is accessed as `item.price`. If you create multiple item of the same type, you can store different values ​​in each.

## Passing a structure to a function

<!-- wave-example: book-struct-function -->
```wave
struct Product {
    price: i32;
    quantity: i32;
}

fun subtotal(item: Product) -> i32 {
    return item.price * item.quantity;
}

fun main() {
    var item: Product = Product {
        price: 1500,
        quantity: 2
    };

    println("{}", subtotal(item));

    item.quantity = 3;
    println("{}", subtotal(item));
}
```

Execution result:

```text
3000
4500
```

The function receives as a type the relationship that the unit price and quantity belong to the same product. This is a function that reads the integer field structure passed as a value. Any function that wants to modify the caller's storage can be designed to accept a pointer.

If the structure has pointer fields, copying value also copies the address. This is not a function for deep copying to separate allocations. Types containing resources such as file handles or Buffer must specify copy and release rules together.

## Method and proto

Related functions can be grouped in method form. proto is a method of writing the method of a structure as a separate block.

<!-- wave-example: book-struct-method -->
```wave
struct Counter {
    value: i32;
}

proto Counter {
    fun current(self: Counter) -> i32 {
        return self.value;
    }
}

fun main() {
    var counter: Counter = Counter {
        value: 3
    };

    println("{}", counter.current());
}
```

Execution result:

```text
3
```

`self: Counter` is a parameter that receives a value. Using a method call notation does not automatically make it a method that modifies the original. Please read together the type of self and what it does in the text.

You can't expect a field to have only a valid state just because you attached a method to it. If there are invalid combinations that the user can create with public fields, your function should check for them or provide a creation rule.

## Name the state with enum

<!-- wave-example: book-enum-state -->
```wave
enum State -> i32 {
    Ready = 0,
    Running,
    Finished,
}

fun main() {
    var state: State = State::Ready;

    if (state == State::Ready) {
        println("ready");
    }

    state = State::Running;

    if (state == State::Running) {
        println("running");
    }
}
```

Execution result:

```text
ready
running
```

`-> i32` is an integer type used in expressions. The first value is set to 0, and subsequent omitted values ​​are 1 greater than the previous value. Rather than simply comparing 0 and 1 in your code, using State::Ready and State::Running will reveal the meaning.

enum Having a name does not automatically restrict state transitions. Rules such as whether it is possible to return from Finished to Running must be implemented as functions.

## Connecting cases and data with variant

If there is a value only in case of success and error information is needed in case of failure, it can be expressed as variant.

<!-- wave-example: book-variant-result -->
```wave
variant Result {
    Value(i32),
    Error(i32)
}

fun divide(left: i32, right: i32) -> Result {
    if (left < 0 || right <= 0) {
        return Result::Error(1);
    }

    return Result::Value(left / right);
}

fun show(result: Result) {
    match (result) {
        Result::Value(value) => {
            println("value={}", value);
        }
        Result::Error(code) => {
            println("error={}", code);
        }
    }
}

fun main() {
    show(divide(12, 3));
    show(divide(12, 0));
}
```

Execution result:

```text
value=4
error=1
```

Result::Value and Result::Error each contain payload. Even if they contain the same integer type, certain cases are distinguished. The caller checks the case with match and uses payload within that arm.

divide above is a small example that does not support negative operands. Because the input range is specified, it should not be confused with a function that covers all boundaries of the regular signed division.

## If payload is not present

Data is not required in all cases. The state of no value can be expressed as a separate case.

<!-- wave-example: book-variant-empty -->
```wave
variant Lookup {
    Found(i32),
    Missing
}

fun main() {
    var result: Lookup = Lookup::Missing;

    match (result) {
        Lookup::Found(value) => {
            println("found={}", value);
        }
        Lookup::Missing => {
            println("missing");
        }
    }
}
```

Execution result:

```text
missing
```

Instead of reserving one arbitrary integer as “none”, we used the case Missing. No matter what the success value is, the meaning does not overlap.

match through `_` handle the remaining cases. If you want each caller to be reviewed again when a new case is added, it is better to separate all cases explicitly. Whichever method you choose, make sure there is no unprocessed input.

## Using structure and variant together

Selecting different data can be expressed as variant, and grouping multiple fields belonging to one case can be expressed as a structure. For example, if the order processing result is success, it can be designed to contain a receipt structure, and if it is a failure, it can be designed to contain an error number.

The memory lifetime rule does not disappear even if the value contains other values. If a pointer is stored in variant, whether the pointer is valid and who will release it are determined separately. Even when saving in an external file format, you must specify encoding for each field rather than dumping the structure memory as is.

## Exercise and complete solution

Create a verdict that only accepts scores between 0 and 100. Valid scores are expressed as Grade(score), the rest are expressed as Invalid.

<!-- wave-example: book-model-solution -->
```wave
variant CheckedScore {
    Grade(i32),
    Invalid
}

fun check_score(score: i32) -> CheckedScore {
    if (score < 0 || score > 100) {
        return CheckedScore::Invalid;
    }

    return CheckedScore::Grade(score);
}

fun main() {
    var result: CheckedScore = check_score(87);

    match (result) {
        CheckedScore::Grade(score) => {
            println("accepted={}", score);
        }
        CheckedScore::Invalid => {
            println("invalid score");
        }
    }
}
```

Execution result:

```text
accepted=87
```

Change the input to -1, 0, 100, 101 to check the boundaries. Since success and failure do not share the same integer space, the caller reduces the risk of accidentally adding error values ​​to the average calculation.


## Which expression should I choose?

|shape of data|suitable expression|yes|
| --- | --- | --- |
|Multiple values of the same type|Array|10 points|
|Multiple fields related to each other|struct|name and score|
|Named state value| enum | Ready, Running, Stopped |
|Additional data that varies by state| variant | Value(i32), Error(str) |
|Contextual name for an existing type|type alias| UserId = u64 |

When choosing a data structure, consider not only the values you want to store, but also what false states you can represent. A structure with both success/failure and value/error fields can create an incorrect combination, but variant can be differentiated into payload for each case.

## Memory placement and external data

The structure's memory may contain empty space to ensure alignment between fields. Simply adding the field sizes does not always equal the total size of the structure. When you need to know the size and alignment, use [mem Layout function](/docs/en/reference/memory-and-buffer).

For files or network messages, the byte order and length can be clearly determined by encoding the fields in order with the [bytes](/docs/en/stdlib/bytes) function. When passing structures to other languages, align the external declaration of [FFI](/docs/en/language/modules-imports-and-ffi) with the target ABI.

[Structure learning and practice](/docs/en/language/structures-enums-and-aliases) · [variant](/docs/en/language/variants)
