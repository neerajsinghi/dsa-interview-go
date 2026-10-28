package dp1d

// Problem: Climbing Stairs
// Difficulty: Easy
//
// A staircase has n steps. Each move climbs either 1 or 2 steps. In how many
// distinct sequences of moves can you reach the top?
//
// Examples:
//
//	n=2 -> 2   (1+1, 2)
//	n=3 -> 3   (1+1+1, 1+2, 2+1)
//	n=5 -> 8
//
// Reference: https://leetcode.com/problems/climbing-stairs/
//
// Pattern: count paths by the LAST move. Every way to reach step i ends with
// a 1-step from i-1 or a 2-step from i-2, and those two groups don't
// overlap, so ways(i) = ways(i-1) + ways(i-2) — Fibonacci in disguise.
//
//	State:      ways(i) = number of move sequences that land exactly on step i
//	Transition: ways(i) = ways(i-1) + ways(i-2)
//	Base:       ways(0) = 1 (the empty sequence), ways(1) = 1
//	Order:      i increasing
//
// Time: O(2^n) brute force; O(n) for the DP versions.
// Space: O(n) memo/table; O(1) optimized.

// ClimbStairsRecursive is the brute force.
func ClimbStairsRecursive(n int) int {
	if n <= 1 {
		return 1
	}
	return ClimbStairsRecursive(n-1) + ClimbStairsRecursive(n-2)
}

// ClimbStairsMemo caches each subproblem in a slice (-1 = not computed).
func ClimbStairsMemo(n int) int {
	memo := make([]int, n+1)
	for i := range memo {
		memo[i] = -1
	}
	var ways func(i int) int
	ways = func(i int) int {
		if i <= 1 {
			return 1
		}
		if memo[i] == -1 {
			memo[i] = ways(i-1) + ways(i-2)
		}
		return memo[i]
	}
	return ways(n)
}

// ClimbStairsTab fills the table bottom-up.
func ClimbStairsTab(n int) int {
	dp := make([]int, max(n+1, 2))
	dp[0], dp[1] = 1, 1
	for i := 2; i <= n; i++ {
		dp[i] = dp[i-1] + dp[i-2]
	}
	return dp[n]
}

// ClimbStairs keeps only the last two values.
func ClimbStairs(n int) int {
	twoBack, oneBack := 1, 1 // ways(0), ways(1)
	for i := 2; i <= n; i++ {
		twoBack, oneBack = oneBack, oneBack+twoBack
	}
	return oneBack
}
