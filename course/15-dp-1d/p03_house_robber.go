package dp1d

// Problem: House Robber I & II
// Difficulty: Medium
//
// Houses along a street hold non-negative amounts of money. You may rob any
// set of houses as long as no two ADJACENT houses are both robbed. Part I:
// the street is a line. Part II: the street is a circle, so the first and
// last houses are also adjacent. Return the most money you can take.
//
// Examples:
//
//	I:  [1,2,3,1]   -> 4    (houses 0 and 2)
//	I:  [2,7,9,3,1] -> 12   (houses 0, 2 and 4)
//	II: [2,3,2]     -> 3    (houses 0 and 2 are now neighbours)
//	II: [1,2,3,1]   -> 4
//
// References: https://leetcode.com/problems/house-robber/
//
//	https://leetcode.com/problems/house-robber-ii/
//
// Pattern: take-or-skip. For house i you either skip it (keep the best of
// the first i houses) or rob it (its money plus the best of the first i-1).
//
//	State:      best(i) = max money from houses 0..i-1
//	Transition: best(i) = max(best(i-1), best(i-2) + nums[i-1])
//	Base:       best(0) = 0, best(1) = nums[0]
//
// Circle: house 0 and house n-1 can't both be robbed, so the answer is the
// better of two LINE problems — houses 0..n-2 and houses 1..n-1. (n = 1 is
// special: the lone house is fine to rob.)
//
// Time: O(2^n) brute force; O(n) DP. Space: O(1) optimized.

// RobRecursive is the brute force for part I.
func RobRecursive(nums []int) int {
	var best func(i int) int // best of the first i houses
	best = func(i int) int {
		if i <= 0 {
			return 0
		}
		return max(best(i-1), best(i-2)+nums[i-1])
	}
	return best(len(nums))
}

// RobMemo caches best(i).
func RobMemo(nums []int) int {
	memo := make([]int, len(nums)+1)
	for i := range memo {
		memo[i] = -1
	}
	var best func(i int) int
	best = func(i int) int {
		if i <= 0 {
			return 0
		}
		if memo[i] == -1 {
			memo[i] = max(best(i-1), best(i-2)+nums[i-1])
		}
		return memo[i]
	}
	return best(len(nums))
}

// RobTab fills best[0..n] bottom-up.
func RobTab(nums []int) int {
	n := len(nums)
	if n == 0 {
		return 0
	}
	dp := make([]int, n+1)
	dp[1] = nums[0]
	for i := 2; i <= n; i++ {
		dp[i] = max(dp[i-1], dp[i-2]+nums[i-1])
	}
	return dp[n]
}

// Rob is the O(1)-space version of part I.
func Rob(nums []int) int {
	skip, take := 0, 0 // best(i-2), best(i-1)
	for _, x := range nums {
		skip, take = take, max(take, skip+x)
	}
	return take
}

// RobCircular solves part II with two linear passes.
func RobCircular(nums []int) int {
	if len(nums) == 1 {
		return nums[0]
	}
	if len(nums) == 0 {
		return 0
	}
	return max(Rob(nums[:len(nums)-1]), Rob(nums[1:]))
}

// RobCircularBrute tries every subset with no two neighbours on the circle.
// O(2^n · n); small inputs only.
func RobCircularBrute(nums []int) int {
	n := len(nums)
	best := 0
	for mask := 0; mask < 1<<n; mask++ {
		adjacent := mask&(mask<<1) != 0 // two neighbouring bits set
		if n > 1 && mask&1 != 0 && mask&(1<<(n-1)) != 0 {
			adjacent = true // first and last houses touch
		}
		if adjacent {
			continue
		}
		sum := 0
		for i := 0; i < n; i++ {
			if mask&(1<<i) != 0 {
				sum += nums[i]
			}
		}
		best = max(best, sum)
	}
	return best
}
