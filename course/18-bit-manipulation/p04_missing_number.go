package bits

// Problem: Missing Number
// Difficulty: Easy
//
// A slice holds n distinct numbers taken from 0..n, so exactly one value in
// that range is absent. Return it.
//
// Examples:
//
//	[3,0,1]             -> 2
//	[9,6,4,2,3,5,7,0,1] -> 8
//	[0]                 -> 1
//
// Reference: https://leetcode.com/problems/missing-number/
//
// Pattern: XOR pairing. XOR all indices 0..n together with all values;
// every present number appears twice and cancels, leaving the missing one.
// (The Gauss-sum version n(n+1)/2 − Σ works too but can overflow in
// fixed-width languages; XOR never overflows.)
//
// Time: O(n). Space: O(1). Lookup-table version: O(n) space.

// MissingNumberSeen marks which values exist.
func MissingNumberSeen(nums []int) int {
	seen := make([]bool, len(nums)+1)
	for _, x := range nums {
		seen[x] = true
	}
	for v, ok := range seen {
		if !ok {
			return v
		}
	}
	return -1 // unreachable for valid input
}

// MissingNumberSum uses the arithmetic series formula.
func MissingNumberSum(nums []int) int {
	n := len(nums)
	sum := n * (n + 1) / 2
	for _, x := range nums {
		sum -= x
	}
	return sum
}

// MissingNumber XORs indices and values.
func MissingNumber(nums []int) int {
	acc := len(nums) // the index loop below covers 0..n-1; n is added here
	for i, x := range nums {
		acc ^= i ^ x
	}
	return acc
}
