package dp1d

// Problem: Min Cost Climbing Stairs
// Difficulty: Easy
//
// cost[i] is the price of stepping on stair i. You may start on stair 0 or
// stair 1, and after paying for a stair you climb 1 or 2 stairs. The top is
// just past the last stair (index n). Return the cheapest way to reach it.
//
// Examples:
//
//	[10, 15, 20]                          -> 15  (start at 1, jump 2 to the top)
//	[1, 100, 1, 1, 1, 100, 1, 1, 100, 1]  -> 6
//
// Reference: https://leetcode.com/problems/min-cost-climbing-stairs/
//
// Pattern: minimum over the last move. To stand on position i you arrived
// from i-1 (paying cost[i-1]) or from i-2 (paying cost[i-2]).
//
//	State:      best(i) = min cost to stand on position i (nothing paid for i yet)
//	Transition: best(i) = min(best(i-1) + cost[i-1], best(i-2) + cost[i-2])
//	Base:       best(0) = best(1) = 0 (free starting points)
//	Answer:     best(n)
//
// Time: O(2^n) brute force; O(n) DP. Space: O(n) table; O(1) optimized.

// MinCostClimbingStairsRecursive is the brute force.
func MinCostClimbingStairsRecursive(cost []int) int {
	var best func(i int) int
	best = func(i int) int {
		if i <= 1 {
			return 0
		}
		return min(best(i-1)+cost[i-1], best(i-2)+cost[i-2])
	}
	return best(len(cost))
}

// MinCostClimbingStairsMemo caches best(i).
func MinCostClimbingStairsMemo(cost []int) int {
	n := len(cost)
	memo := make([]int, n+1)
	done := make([]bool, n+1) // separate flag: 0 is a legitimate cached value
	var best func(i int) int
	best = func(i int) int {
		if i <= 1 {
			return 0
		}
		if !done[i] {
			memo[i] = min(best(i-1)+cost[i-1], best(i-2)+cost[i-2])
			done[i] = true
		}
		return memo[i]
	}
	return best(n)
}

// MinCostClimbingStairsTab fills best[0..n] bottom-up.
func MinCostClimbingStairsTab(cost []int) int {
	n := len(cost)
	dp := make([]int, max(n+1, 2)) // dp[0] = dp[1] = 0
	for i := 2; i <= n; i++ {
		dp[i] = min(dp[i-1]+cost[i-1], dp[i-2]+cost[i-2])
	}
	return dp[n]
}

// MinCostClimbingStairs rolls the table into two variables.
func MinCostClimbingStairs(cost []int) int {
	twoBack, oneBack := 0, 0
	for i := 2; i <= len(cost); i++ {
		twoBack, oneBack = oneBack, min(oneBack+cost[i-1], twoBack+cost[i-2])
	}
	return oneBack
}
