package bits

// Problem: Sum of Two Integers (without + or -)
// Difficulty: Medium
//
// Add two signed 32-bit integers without using the + or - operators.
//
// Examples:
//
//	1 + 2   -> 3
//	-2 + 3  -> 1
//	MaxInt32 + 1 -> MinInt32 (wraps, like hardware)
//
// Reference: https://leetcode.com/problems/sum-of-two-integers/
//
// Pattern: a half adder, repeated. XOR adds without carrying; AND finds the
// positions that carry, and << 1 moves each carry to the next column. Loop
// until no carry remains.
//
//	  a = 0101 (5)        b = 0011 (3)
//	  a^b = 0110          (a&b)<<1 = 0010
//	  a^b = 0100          (a&b)<<1 = 0100
//	  a^b = 0000          (a&b)<<1 = 1000
//	  a^b = 1000 (8)      carry 0 -> done
//
// Go makes this easy: signed integers are two's complement, overflow wraps
// (it is defined, not UB), and a left shift on a signed value is allowed.
// With fixed-width int32 the carry falls off the top after ≤ 32 rounds,
// so negative inputs terminate too — no masking needed (unlike Python).
//
// Time: O(32) iterations at most. Space: O(1).

// GetSum adds a and b with XOR/AND/shift only.
func GetSum(a, b int32) int32 {
	for b != 0 {
		carry := (a & b) << 1
		a ^= b // sum without carries
		b = carry
	}
	return a
}

// GetDifference computes a - b as a + (-b), where -b == ^b + 1 in two's
// complement (and the +1 is itself done with GetSum).
func GetDifference(a, b int32) int32 {
	return GetSum(a, GetSum(^b, 1))
}
