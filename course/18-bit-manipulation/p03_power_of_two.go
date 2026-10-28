package bits

// Problem: Power of Two
// Difficulty: Easy
//
// Decide whether n equals 2^k for some integer k ≥ 0.
//
// Examples:
//
//	1  -> true  (2^0)
//	16 -> true
//	6  -> false
//	0, -8 -> false
//
// Reference: https://leetcode.com/problems/power-of-two/
//
// Pattern: a power of two has exactly ONE set bit, and n & (n-1) clears the
// lowest set bit — so the result is zero exactly when n had a single 1.
// Guard n > 0: 0 has no set bits, and negative numbers are never powers of
// two (math.MinInt has one set bit — the sign bit — and would fool the test).
//
// Variation: power of four = power of two whose bit sits at an even
// position, i.e. n & 0x5555_5555_5555_5555 != 0.
//
// Time: O(1). Space: O(1). Division loop: O(log n).

// IsPowerOfTwoLoop divides by two while even.
func IsPowerOfTwoLoop(n int) bool {
	if n <= 0 {
		return false
	}
	for n%2 == 0 {
		n /= 2
	}
	return n == 1
}

// IsPowerOfTwo is the single-bit test.
func IsPowerOfTwo(n int) bool {
	return n > 0 && n&(n-1) == 0
}

// IsPowerOfFour adds an "even bit position" mask to the power-of-two test.
func IsPowerOfFour(n int) bool {
	const evenBits = 0x5555_5555_5555_5555 // ...0101 0101
	return IsPowerOfTwo(n) && n&evenBits != 0
}
