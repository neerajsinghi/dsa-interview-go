package slidingwindow

// Problem 01 — Maximum Sum Subarray of Size K (Easy)
//
// Return the largest sum of any k consecutive elements of nums. If k is not
// in 1..len(nums), report ok = false.
//
//	nums = [2, 1, 5, 1, 3, 2], k = 3 → 9   (5 + 1 + 3)
//	nums = [-1, -2, -3],       k = 2 → -3  (-1 + -2)
//
// Reference: https://leetcode.com/problems/maximum-average-subarray-i/ (same idea, averaged)
//
// Pattern: FIXED-SIZE window. Moving the window one step right adds one
// element and removes one, so the new sum is old + nums[right] - nums[right-k].
//
//	[2, 1, 5, 1, 3, 2]
//	[2  1  5]           sum 8
//	   [1  5  1]        8 + 1 - 2 = 7
//	      [5  1  3]     7 + 3 - 1 = 9  ← best
//	         [1  3  2]  9 + 2 - 5 = 6
//
//	MaxSumSubarrayBrute: Time O(n·k), Space O(1)
//	MaxSumSubarray:      Time O(n),   Space O(1)

// MaxSumSubarrayBrute sums every window from scratch.
func MaxSumSubarrayBrute(nums []int, k int) (best int, ok bool) {
	if k <= 0 || k > len(nums) {
		return 0, false
	}
	for start := 0; start+k <= len(nums); start++ {
		sum := 0
		for i := start; i < start+k; i++ {
			sum += nums[i]
		}
		if start == 0 || sum > best {
			best = sum
		}
	}
	return best, true
}

// MaxSumSubarray slides the window, adjusting the sum in O(1) per step.
func MaxSumSubarray(nums []int, k int) (best int, ok bool) {
	if k <= 0 || k > len(nums) {
		return 0, false
	}
	sum := 0
	for i := range k {
		sum += nums[i]
	}
	best = sum // seed with the first window, NOT 0 — all values may be negative
	for right := k; right < len(nums); right++ {
		sum += nums[right] - nums[right-k]
		best = max(best, sum)
	}
	return best, true
}
