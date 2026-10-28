package dp1d

// Problem: Maximum Product Subarray
// Difficulty: Medium
//
// Given a non-empty list of integers (which may be negative or zero), find
// the contiguous, non-empty subarray with the largest product and return
// that product.
//
// Examples:
//
//	[2, 3, -2, 4]  -> 6    ([2 3])
//	[-2, 0, -1]    -> 0
//	[-2, 3, -4]    -> 24   (the two negatives cancel)
//
// Reference: https://leetcode.com/problems/maximum-product-subarray/
//
// Pattern: Kadane's algorithm, tracking TWO states. For sums, the best
// subarray ending at i only depends on the best ending at i-1. For products
// a negative number flips the order: the most negative product so far can
// become the largest after multiplying by a negative. So keep both
//
//	hi(i) = max product of a subarray ending at i
//	lo(i) = min product of a subarray ending at i
//
//	hi(i) = max(x, hi(i-1)·x, lo(i-1)·x)
//	lo(i) = min(x, hi(i-1)·x, lo(i-1)·x)      where x = nums[i]
//
// Starting fresh at x (the first term) handles zeros: after a 0, both hi and
// lo restart from the next element.
//
//	x:    2   3   -2   4
//	hi:   2   6   -2   4      answer = max(hi) = 6
//	lo:   2   3  -12 -48
//
// Time: O(n²) brute force; O(n) DP. Space: O(1).

// MaxProductBrute tries every subarray, extending the product as it goes.
func MaxProductBrute(nums []int) int {
	best := nums[0]
	for i := range nums {
		p := 1
		for j := i; j < len(nums); j++ {
			p *= nums[j]
			best = max(best, p)
		}
	}
	return best
}

// MaxProduct runs the two-state Kadane scan.
func MaxProduct(nums []int) int {
	hi, lo, best := nums[0], nums[0], nums[0]
	for _, x := range nums[1:] {
		// Compute both from the OLD hi and lo before overwriting either.
		hi, lo = max(x, hi*x, lo*x), min(x, hi*x, lo*x)
		best = max(best, hi)
	}
	return best
}
