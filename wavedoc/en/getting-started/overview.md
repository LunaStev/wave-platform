---
translation_set_id: overview
path: getting-started/overview
locale: en
group: getting-started
group_order: 1
order: 1
title: Wave documentation and learning guide
summary: Learn Wave step by step, from installation to practical programs, and look up language rules and standard library APIs.
---

## Learn Wave with this guide

Learn to write Wave source code, compile and run it, and check the results. If you are new to programming, follow the sequence below. If you know another language, run each chapter’s examples and compare its rules and boundary cases with what you already know.

## Learning path

|Step|Chapter|What you will learn|
| --- | --- | --- |
|Setup| [Installation](/docs/en/getting-started/install) |Prepare the compiler and standard library and verify that they run|
| 1 | [Your first program](/docs/en/language/program-structure) |Create, check, and run a source file and understand exit codes|
| 2 | [Variables and types](/docs/en/language/declarations-and-types) |Store values and choose types with the required range|
| 3 | [Operators and conversions](/docs/en/language/expressions-and-operators) |Explain evaluation order and the results of type conversions|
| 4 | [Conditions and loops](/docs/en/language/control-flow) |Branch on conditions and process data with loops|
| 5 | [Functions](/docs/en/language/functions-and-generics) |Extract repeated operations into functions|
| 6 | [Arrays](/docs/en/language/arrays) |Access elements by index and iterate over an array|
| 7 | [Strings](/docs/en/language/strings) |Distinguish characters from bytes and understand escapes and string length|
| 8 | [Structs and variants](/docs/en/language/structures-enums-and-aliases) |Group related data and represent success and failure|
| 9 | [Pointers and lifetimes](/docs/en/language/explicit-memory-type-model) |Modify the original value through its address and manage its lifetime|
| 10 | [Dynamic memory](/docs/en/language/allocation) |Handle allocation failures and free memory|
| 11 | [Modules and generics](/docs/en/language/modules-imports-and-ffi) |Split code across files and reuse functions with different types|
| 12 | [Error handling](/docs/en/language/errors) |Check results and clean up resources on failure|
| 13 | [Introduction to asynchronous code](/docs/en/language/async-and-never) |Create a Future and wait for it to finish|

## Put your knowledge into practice

After the core chapters, build an [input calculator](/docs/en/practice/input-calculator), a [file reader](/docs/en/practice/file-reader), a [binary message](/docs/en/practice/binary-message), and a [TCP client](/docs/en/practice/tcp-client). Test both successful input and failure cases in each project.

## Three documentation tabs

- **Wave**: A guided language course and practical projects to follow in order.
- **[Standard library](/docs/en/stdlib)**: Each module’s APIs, return values, errors, ownership rules, and platform requirements.
- **[Whale](/docs/en/whale)**: Building and linking, package management, command usage, and the low-level toolchain.

Examples distinguish complete programs from snippets that belong inside a function. Run `wavec` commands in a terminal and save `wave` code blocks in `.wave` files. Input and output are shown separately; examples that read standard input specify what to enter.

## When you get stuck

Use [Troubleshooting](/docs/en/reference/diagnostics) to distinguish installation, source checking, linking, and execution problems. Look up language rules in the [syntax quick reference](/docs/en/reference/syntax-quick-reference), commands in the [compiler reference](/docs/en/getting-started/compiler), and APIs in the [standard library guide](/docs/en/reference/standard-library).
