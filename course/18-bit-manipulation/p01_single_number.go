package bits

// Problem: Single Number
// Difficulty: Easy
//
// Every value in a non-empty slice appears exactly twice except one value,
// which appears once. Find it in linear time without extra memory.
//
// Examples:
//
//	[2,2,1]     -> 1
//	[4,1,2,1,2] -> 4
//
// Reference: https://leetcode.com/problems/single-number/
//
// Pattern: XOR cancels pairs. x^x == 0, x^0 == x, and XOR is commutative
// and associative, so XOR-ing everything leaves only the unpaired value.
//
//	4 ^ 1 ^ 2 ^ 1 ^ 2 = 4 ^ (1^1) ^ (2^2) = 4 ^ 0 ^ 0 = 4
//
// Time: O(n). Space: O(1). Map version: O(n) space.

// SingleNumberMap counts occurrences — the obvious first answer.
func SingleNumberMap(nums []int) int {
	seen := make(map[int]int, len(nums))
	for _, x := range nums {
		seen[x]++
	}
	for x, c := range seen {
		if c == 1 {
			return x
		}
	}
	return 0
}

// SingleNumber XORs every element.
func SingleNumber(nums []int) int {
	acc := 0
	for _, x := range nums {
		acc ^= x
	}
	return acc
}
