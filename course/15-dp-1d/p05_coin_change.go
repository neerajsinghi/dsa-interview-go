package dp1d

// Problem: Coin Change I & II
// Difficulty: Medium
//
// You have unlimited coins of each given denomination. Part I: return the
// FEWEST coins that add up to amount, or -1 if it can't be done. Part II:
// return the NUMBER OF DISTINCT COMBINATIONS (order doesn't matter: 1+2 and
// 2+1 are the same) that add up to amount.
//
// Examples:
//
//	I:  coins [1,2,5], amount 11 -> 3   (5+5+1)
//	I:  coins [2],     amount 3  -> -1
//	II: coins [1,2,5], amount 5  -> 4   (5, 2+2+1, 2+1+1+1, 1+1+1+1+1)
//	II: coins [2],     amount 3  -> 0
//
// References: https://leetcode.com/problems/coin-change/
//
//	https://leetcode.com/problems/coin-change-ii/
//
// Pattern: unbounded knapsack over the amount.
//
//	Part I   State:      fewest(a) = min coins summing to a
//	         Transition: fewest(a) = 1 + min over coins c ≤ a of fewest(a-c)
//	         Base:       fewest(0) = 0
//
//	Part II  State:      ways[a] after processing the first k coin types
//	         Transition: ways[a] += ways[a-c]   for the k-th coin c
//	         Base:       ways[0] = 1
//
// LOOP ORDER IS THE WHOLE TRICK in part II. With coins in the OUTER loop,
// every combination is built in coin order, so each multiset is counted once.
// With amounts in the outer loop you count ORDERED sequences instead
// (1+2 and 2+1 separately) — that is a different problem
// ("Combination Sum IV"). See CountOrderings below for the contrast.
//
// Time: O(len(coins) · amount) for the DP versions. Space: O(amount).

// CoinChangeRecursive is the brute force for part I.
func CoinChangeRecursive(coins []int, amount int) int {
	var fewest func(a int) int // returns -1 if a can't be made
	fewest = func(a int) int {
		if a == 0 {
			return 0
		}
		best := -1
		for _, c := range coins {
			if c <= a {
				if sub := fewest(a - c); sub != -1 && (best == -1 || sub+1 < best) {
					best = sub + 1
				}
			}
		}
		return best
	}
	return fewest(amount)
}

// CoinChangeMemo caches fewest(a). -2 marks "not computed", -1 "impossible".
func CoinChangeMemo(coins []int, amount int) int {
	memo := make([]int, amount+1)
	for i := range memo {
		memo[i] = -2
	}
	var fewest func(a int) int
	fewest = func(a int) int {
		if a == 0 {
			return 0
		}
		if memo[a] != -2 {
			return memo[a]
		}
		best := -1
		for _, c := range coins {
			if c <= a {
				if sub := fewest(a - c); sub != -1 && (best == -1 || sub+1 < best) {
					best = sub + 1
				}
			}
		}
		memo[a] = best
		return best
	}
	return fewest(amount)
}

// CoinChange is the bottom-up table for part I (already O(amount) space).
func CoinChange(coins []int, amount int) int {
	const unreachable = 1 << 62 // larger than any real answer, safe to +1
	dp := make([]int, amount+1)
	for a := 1; a <= amount; a++ {
		dp[a] = unreachable
		for _, c := range coins {
			if c <= a && dp[a-c]+1 < dp[a] {
				dp[a] = dp[a-c] + 1
			}
		}
	}
	if dp[amount] >= unreachable {
		return -1
	}
	return dp[amount]
}

// ChangeRecursive is the brute force for part II: at coin i, either use it
// (stay on i, since reuse is allowed) or move past it for good.
func ChangeRecursive(amount int, coins []int) int {
	var ways func(i, a int) int
	ways = func(i, a int) int {
		if a == 0 {
			return 1
		}
		if i == len(coins) || a < 0 {
			return 0
		}
		return ways(i, a-coins[i]) + ways(i+1, a)
	}
	return ways(0, amount)
}

// ChangeMemo caches ways(i, a) in a 2-D table.
func ChangeMemo(amount int, coins []int) int {
	memo := make([][]int, len(coins))
	for i := range memo {
		memo[i] = make([]int, amount+1)
		for a := range memo[i] {
			memo[i][a] = -1
		}
	}
	var ways func(i, a int) int
	ways = func(i, a int) int {
		if a == 0 {
			return 1
		}
		if i == len(coins) || a < 0 {
			return 0
		}
		if memo[i][a] == -1 {
			memo[i][a] = ways(i, a-coins[i]) + ways(i+1, a)
		}
		return memo[i][a]
	}
	return ways(0, amount)
}

// ChangeTab is the 2-D bottom-up table: dp[i][a] = combinations of a using
// only coins[i:].
func ChangeTab(amount int, coins []int) int {
	n := len(coins)
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, amount+1)
		dp[i][0] = 1
	}
	for i := n - 1; i >= 0; i-- {
		for a := 1; a <= amount; a++ {
			dp[i][a] = dp[i+1][a] // skip coin i
			if a >= coins[i] {
				dp[i][a] += dp[i][a-coins[i]] // use coin i (again)
			}
		}
	}
	return dp[0][amount]
}

// Change is the 1-D space-optimized part II. Coins outer, amounts inner,
// ascending so each coin may be reused.
func Change(amount int, coins []int) int {
	ways := make([]int, amount+1)
	ways[0] = 1
	for _, c := range coins {
		for a := c; a <= amount; a++ {
			ways[a] += ways[a-c]
		}
	}
	return ways[amount]
}

// CountOrderings swaps the loops and therefore counts ORDERED sequences of
// coins (1+2 and 2+1 are different). Shown only to contrast with Change.
func CountOrderings(amount int, coins []int) int {
	ways := make([]int, amount+1)
	ways[0] = 1
	for a := 1; a <= amount; a++ {
		for _, c := range coins {
			if c <= a {
				ways[a] += ways[a-c]
			}
		}
	}
	return ways[amount]
}
