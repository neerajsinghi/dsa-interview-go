package dp2d

// Problem: Burst Balloons
// Difficulty: Hard
//
// A row of balloons carries numbers. Bursting balloon i earns
// left · nums[i] · right, where left/right are its CURRENT neighbours
// (missing neighbours at the ends count as 1). After a burst the
// neighbours become adjacent. Burst all balloons; maximise total coins.
//
// Examples:
//
//	[3,1,5,8] -> 167   (burst 1, 5, 3, 8: 15 + 120 + 24 + 8)
//	[1,5]     -> 10
//
// Reference: https://leetcode.com/problems/burst-balloons/
//
// Pattern: interval DP — pick the LAST balloon to burst. Pad the array with
// 1 on both ends: vals = [1, nums..., 1]. dp[l][r] = best coins from
// bursting everything strictly between l and r. If k is the last one
// burst in (l, r), its neighbours at that moment are exactly l and r, and
// the two sides (l, k) and (k, r) are independent:
//
//	dp[l][r] = max over l<k<r of dp[l][k] + vals[l]·vals[k]·vals[r] + dp[k][r]
//
// Thinking "first burst" fails because the sides stay coupled.
//
// Time: O(n³). Space: O(n²). Brute force over burst orders: O(n!).

// BurstBalloonsBrute tries every burst order. Factorial; tests only.
func BurstBalloonsBrute(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	best := 0
	for i := range nums {
		left, right := 1, 1
		if i > 0 {
			left = nums[i-1]
		}
		if i < len(nums)-1 {
			right = nums[i+1]
		}
		rest := append(append([]int{}, nums[:i]...), nums[i+1:]...)
		best = max(best, left*nums[i]*right+BurstBalloonsBrute(rest))
	}
	return best
}

// padded returns [1, nums..., 1].
func padded(nums []int) []int {
	vals := make([]int, 0, len(nums)+2)
	vals = append(vals, 1)
	vals = append(vals, nums...)
	return append(vals, 1)
}

// BurstBalloonsMemo is top-down over open intervals (l, r).
func BurstBalloonsMemo(nums []int) int {
	vals := padded(nums)
	n := len(vals)
	memo := newTable(n, n, -1)
	var best func(l, r int) int
	best = func(l, r int) int {
		if l+1 == r {
			return 0 // nothing strictly between
		}
		if memo[l][r] != -1 {
			return memo[l][r]
		}
		res := 0
		for k := l + 1; k < r; k++ { // k is burst last in (l, r)
			res = max(res, best(l, k)+vals[l]*vals[k]*vals[r]+best(k, r))
		}
		memo[l][r] = res
		return res
	}
	return best(0, n-1)
}

// BurstBalloonsTab fills intervals in order of increasing width.
func BurstBalloonsTab(nums []int) int {
	vals := padded(nums)
	n := len(vals)
	dp := newTable(n, n, 0)
	for width := 2; width < n; width++ {
		for l := 0; l+width < n; l++ {
			r := l + width
			for k := l + 1; k < r; k++ {
				dp[l][r] = max(dp[l][r], dp[l][k]+vals[l]*vals[k]*vals[r]+dp[k][r])
			}
		}
	}
	return dp[0][n-1]
}
