package greedy

// Problem: Maximum Subarray (Kadane's algorithm)
// Difficulty: Medium
//
// Given a non-empty slice of integers (negatives allowed), return the
// largest sum of any contiguous, non-empty run of elements.
//
// Examples:
//
//	[-2,1,-3,4,-1,2,1,-5,4] -> 6   ([4,-1,2,1])
//	[-3,-1,-2]              -> -1  (must pick at least one element)
//
// Reference: https://leetcode.com/problems/maximum-subarray/
//
// Pattern: greedy / 1-D DP. Walking left to right, the best run ending at i
// either extends the best run ending at i-1 or starts fresh at i — a
// negative running sum can only hurt what follows, so drop it.
//
//	cur  = max(x, cur + x)
//	best = max(best, cur)
//
// Time: O(n). Space: O(1). Brute force: O(n²).

// MaxSubarrayBrute tries every start and extends to every end.
func MaxSubarrayBrute(nums []int) int {
	best := nums[0]
	for i := range nums {
		sum := 0
		for j := i; j < len(nums); j++ {
			sum += nums[j]
			best = max(best, sum)
		}
	}
	return best
}

// MaxSubarray is Kadane's algorithm. It panics on an empty slice, matching
// the problem's "non-empty" precondition.
func MaxSubarray(nums []int) int {
	cur, best := nums[0], nums[0]
	for _, x := range nums[1:] {
		cur = max(x, cur+x) // extend, or restart at x
		best = max(best, cur)
	}
	return best
}

// MaxSubarrayRange also returns the inclusive bounds [lo, hi] of one best run.
func MaxSubarrayRange(nums []int) (sum, lo, hi int) {
	cur, start := nums[0], 0
	sum, lo, hi = nums[0], 0, 0
	for i := 1; i < len(nums); i++ {
		if cur < 0 { // a negative prefix never helps; restart here
			cur, start = nums[i], i
		} else {
			cur += nums[i]
		}
		if cur > sum {
			sum, lo, hi = cur, start, i
		}
	}
	return sum, lo, hi
}
