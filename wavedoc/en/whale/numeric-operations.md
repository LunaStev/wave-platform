---
translation_set_id: whale-numeric-operations
path: whale/numeric-operations
locale: en
group: whale
group_order: 5
order: 3
title: Numeric operations
summary: Integer wrapping, checked arithmetic, shifts, conversion errors, and floating-point results.
---

## Integer representation

An N-bit integer has N value bits. Unsigned integers range from 0 to 2^N − 1; signed integers range from −2^(N−1) to 2^(N−1) − 1. The signedness of an operation determines how the bit pattern is interpreted.

The examples below describe operation results; they are not textual IR syntax.

## Addition, subtraction, and multiplication

Basic integer add, sub, and mul keep the low N bits of the result. Overflow does not trap. Checked arithmetic returns the same wrapped result together with a Bool indicating whether the mathematical result exceeded the signed or unsigned range of the operation.

| Operation | Wrapped result | Checked overflow |
| --- | --- | --- |
| u8: 255 + 1 | 0 | true |
| i8: 127 + 1 | −128 | true |
| u8: 0 − 1 | 255 | true |
| i8: 12 × 3 | 36 | false |

A frontend that requires overflow to terminate execution must use checked arithmetic and an explicit `trap_if` on the overflow result. Basic arithmetic does not inherit the source language's overflow policy implicitly.

## Division and remainder

Integer division and remainder by zero trap. Signed division of the minimum representable value by −1 produces that minimum value by wrapping. The corresponding remainder is zero.

| Operation | Result |
| --- | --- |
| i8: −128 / −1 | −128 |
| i8: −128 % −1 | 0 |
| Integer division by 0 | trap |
| Integer remainder by 0 | trap |

## Shifts

For an N-bit value, interpret the shift-count bit pattern as unsigned and reduce it modulo N. A count outside the range 0 through N−1 does not itself trap.

For an 8-bit value, counts 0, 8, and 16 all select a shift of zero. An 8-bit count with bit pattern `11111111` selects a shift of 7, including when that pattern represents signed −1. This unsigned interpretation occurs before the modulo operation.

A source language that rejects excessive or negative counts must express that policy with explicit checks before the shift.

## Conversions

| Conversion | Meaning |
| --- | --- |
| Zero extension | Increase width by adding zero high bits |
| Sign extension | Increase width by replicating the sign bit |
| Bit truncation | Retain the low bits at the destination width |
| Bit reinterpretation | Interpret the same bits using another type |
| Lossless numerical conversion | Preserve the numerical value; trap when it cannot be represented |

For example, extending the bit pattern `11111111` from 8 to 16 bits by zero extension produces `0000000011111111`. Sign extension produces `1111111111111111`. These are different operations even when the source bits are identical.

Float-to-integer conversion truncates toward zero, then checks the integer range. NaN and infinity trap. For conversion to i8, 127.9 produces 127, while 128.0 traps. Bool converts to integer 0 or 1, except that signed i1 is not a permitted destination for this conversion because it cannot represent 1.

Converting an address to an integer does not preserve a right to recover pointer access permissions from that integer. See [pointer validity](memory-model).

## Floating-point arithmetic

Floating-point values have exact f16, f32, or f64 bit patterns. Arithmetic rounds at the declared width using round-to-nearest, ties-to-even: a halfway result selects the representable value with an even least-significant significand bit.

Default arithmetic does not permit fast math, implicit fused multiply-add, or flushing small values to zero. A multiply followed by an add retains its separate rounding steps; the backend must not combine them implicitly.

Numerical operations may produce NaN or infinity. NaNs produced by numerical operations or width conversions are normalized to one fixed positive quiet NaN per width. Storage and copying instead preserve the original NaN bits. This distinction matters when moving a NaN payload through memory without performing arithmetic on it.

Floating-point status flags are not exposed. Float-to-integer conversion has the trapping rules above even though default floating-point arithmetic permits NaN and infinity results.
