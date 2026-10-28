package bits

import mbits "math/bits"

// Problem: Counting Bits
// Difficulty: Easy
//
// For every integer i from 0 to n inclusive, report how many 1 bits i has.
// Aim for a single linear pass rather than popcounting each number.
//
// Examples:
//
//	n=2 -> [0 1 1]
//	n=5 -> [0 1 1 2 1 2]
//
// Reference: https://leetcode.com/problems/counting-bits/
//
// Pattern: DP on bits — reuse the answer for a smaller number.
//
//	ans[i] = ans[i>>1] + (i & 1)     // drop the last bit, add it back
//	ans[i] = ans[i&(i-1)] + 1        // drop the lowest set bit
//
// Both i>>1 and i&(i-1) are < i, so they're already computed.
//
// Time: O(n). Space: O(n) for the output.

// CountBitsNaive popcounts every number.
func CountBitsNaive(n int) []int {
	ans := make([]int, n+1)
	for i := range ans {
		ans[i] = mbits.OnesCount(uint(i))
	}
	return ans
}

// CountBitsShift uses the "half plus last bit" recurrence.
func CountBitsShift(n int) []int {
	ans := make([]int, n+1)
	for i := 1; i <= n; i++ {
		ans[i] = ans[i>>1] + i&1
	}
	return ans
}

// CountBits uses the "drop the lowest set bit" recurrence.
func CountBits(n int) []int {
	ans := make([]int, n+1)
	for i := 1; i <= n; i++ {
		ans[i] = ans[i&(i-1)] + 1
	}
	return ans
}
