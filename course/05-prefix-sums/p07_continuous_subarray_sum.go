package prefixsums

// Problem 07 — Continuous Subarray Sum (Medium)
//
// Given nums and k ≥ 1, report whether some contiguous subarray of length at
// least 2 sums to a multiple of k (0 counts as a multiple). This version also
// accepts negative numbers.
//
//	nums = [23, 2, 4, 6, 7], k = 6 → true   ([2, 4] = 6)
//	nums = [23, 2, 6, 4, 7], k = 13 → false
//
// Reference: https://leetcode.com/problems/continuous-subarray-sum/
//
// Pattern: prefix sums MODULO k. If P[i] ≡ P[j] (mod k) then the subarray
// between them sums to a multiple of k. Store the first index of each
// remainder and check that the span is at least 2.
//
// Go gotcha: -7 % 6 == -1 in Go, but the matching non-negative remainder is
// 5. Normalise with ((x % k) + k) % k, or negative and positive prefixes with
// the same true remainder will be missed.
//
//	nums = [23, 2, 4], k = 6
//	prefix mod 6: (start) 0, 5, 1, 5  → remainder 5 at index 0 and 2: span 2 ✓
//
//	CheckSubarraySumBrute: Time O(n²), Space O(1)
//	CheckSubarraySum:      Time O(n),  Space O(min(n, k))

// CheckSubarraySumBrute tries every subarray of length ≥ 2.
func CheckSubarraySumBrute(nums []int, k int) bool {
	for i := range nums {
		sum := nums[i]
		for j := i + 1; j < len(nums); j++ {
			sum += nums[j]
			if mod(sum, k) == 0 {
				return true
			}
		}
	}
	return false
}

// CheckSubarraySum remembers the first index of each prefix remainder.
func CheckSubarraySum(nums []int, k int) bool {
	first := map[int]int{0: -1} // empty prefix has remainder 0, "at" index -1
	prefix := 0
	for i, x := range nums {
		prefix = mod(prefix+x, k) // keep it reduced: no overflow, valid key
		if j, ok := first[prefix]; ok {
			if i-j >= 2 {
				return true
			}
		} else {
			first[prefix] = i
		}
	}
	return false
}
