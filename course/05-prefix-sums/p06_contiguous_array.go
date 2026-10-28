package prefixsums

// Problem 06 — Contiguous Array (Medium)
//
// nums contains only 0s and 1s. Return the length of the longest contiguous
// subarray with an equal number of 0s and 1s.
//
//	[0, 1]          → 2
//	[0, 1, 1, 1, 0, 0] → 6
//	[0, 0, 0]       → 0
//
// Reference: https://leetcode.com/problems/contiguous-array/
//
// Pattern: TRANSFORM, then prefix sums. Treat 0 as −1; "equal counts" becomes
// "sum is 0", i.e. two equal prefix values. For the LONGEST such subarray,
// remember the FIRST index at which each prefix value appeared.
//
//	nums:     0   1   1   1   0   0
//	as ±1:   -1  +1  +1  +1  -1  -1
//	prefix: 0 -1   0   1   2   1   0
//	        ^-------------------------^ prefix 0 at index -1 and 5 → length 6
//
//	FindMaxLengthBrute: Time O(n²), Space O(1)
//	FindMaxLength:      Time O(n),  Space O(n)

// FindMaxLengthBrute tracks the balance from every start.
func FindMaxLengthBrute(nums []int) int {
	best := 0
	for i := range nums {
		balance := 0
		for j := i; j < len(nums); j++ {
			if nums[j] == 1 {
				balance++
			} else {
				balance--
			}
			if balance == 0 {
				best = max(best, j-i+1)
			}
		}
	}
	return best
}

// FindMaxLength stores the first index of each prefix balance.
func FindMaxLength(nums []int) int {
	first := map[int]int{0: -1} // balance 0 "occurs" before the array starts
	best, balance := 0, 0
	for i, x := range nums {
		if x == 1 {
			balance++
		} else {
			balance--
		}
		if j, ok := first[balance]; ok {
			best = max(best, i-j)
		} else {
			first[balance] = i // only the first occurrence: it gives the longest span
		}
	}
	return best
}
