---
translation_set_id: diagnostics
path: reference/diagnostics
locale: en
group: reference
group_order: 5
order: 2
title: Troubleshooting: From installation to execution
summary: Isolate failed steps and narrow down the cause with reproducible information.
---

## First, distinguish the stages of failure

|observed phenomenon|Check first|next action|
| --- | --- | --- |
|wavec Command not found|PATH and executable file location|Run with absolute path and set PATH|
|Can't find files needed to run|Are any files missing from the installation folder?|Unpack and install the entire package again|
|std import Failed| `wavec print std-path` |Correspondence: Install std or specify `--std-root`|
|Source location and type error output| `wavec check main.wave` |Fix the first error and check again|
|Build failure for other OS·CPU targets|Specified target and target environment|[Cross build settings](/docs/en/whale/build-link-targets) Confirm|
|Execution fails after successful build|exit code, input, working directory|Execution environment and API error check|

## A small diagnostic example

Here's the entire program, which is intentionally incorrect:

```wave
fun main() {
    var count: i32 = 1;
    println("{}", missing);
}
```

`wavec check main.wave` must point to the undeclared name missing. Change the variable name to count, then check and run again. Focus on the file, location, and cause rather than the entire diagnostic text. Subsequent errors may be a result of the initial error.

## When an executable fails

In the Linux/macOS shell, the exit code is checked immediately after execution as `echo $?`, and in PowerShell, it is `$LASTEXITCODE`. Input error and explicit `return 1` are not the same cause. Invalid runtime values ​​for shift count or real conversion can cause trap. Check out [operation rules](/docs/en/language/expressions-and-operators).

Relative file paths are affected by the executable working directory rather than the source file location. Don't treat a file read failure as a string length of 0, check the return error first. Network connection failures are checked through address lookup, server waiting, permissions, and timeouts.

## Information needed to report a problem

1. `wavec --version` Output and exact command executed.
2. target. specified separately from host OS·architecture
3. Compiler source used with the selected std path.
4. Minimal source, input and required files to reproduce the problem.
5. Expected results, actual results, diagnosis and exit codes.

Passwords, tokens, and personal file contents are removed. If the problem goes away when you reduce the minimal example, the last element you removed is the clue. `--error-format=json` is available when the tool collects diagnostics.

[Installation](/docs/en/getting-started/install) · [compiler command](/docs/en/getting-started/compiler) · [Targets and Links](/docs/en/whale/build-link-targets)
