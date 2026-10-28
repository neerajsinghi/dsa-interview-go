package dp2d

// Problem: 0/1 Knapsack
// Difficulty: Medium
//
// You have n items; item i weighs weights[i] and is worth values[i]. A bag
// holds at most capacity total weight. Each item is taken whole or left
// behind (no fractions, no repeats). Return the best total value.
//
// Examples:
//
//	weights [1,3,4,5], values [1,4,5,7], capacity 7 -> 9  (items 3+4)
//	weights [5], values [10], capacity 4            -> 0
//
// Reference: https://en.wikipedia.org/wiki/Knapsack_problem#0-1_knapsack_problem
//
// Pattern: "take it or leave it" DP. dp[i][c] = best value using the first i
// items with capacity c.
//
//	dp[i][c] = max(dp[i-1][c],                               // skip item i-1
//	               dp[i-1][c-w] + v  if w <= c)              // take it
//
// Row i depends only on row i-1, so a single array works IF capacity is
// iterated from high to low — otherwise dp[c-w] would already include item
// i and you'd silently solve the unbounded knapsack.
//
// Time: O(n·W). Space: O(n·W) table, O(W) rolling. Brute force: O(2^n).

// KnapsackBrute tries every subset by recursion.
func KnapsackBrute(weights, values []int, capacity int) int {
	var best func(i, c int) int
	best = func(i, c int) int {
		if i == len(weights) {
			return 0
		}
		skip := best(i+1, c)
		if weights[i] > c {
			return skip
		}
		return max(skip, values[i]+best(i+1, c-weights[i]))
	}
	return best(0, capacity)
}

// KnapsackMemo is the brute force with a (item, capacity) cache.
func KnapsackMemo(weights, values []int, capacity int) int {
	if capacity < 0 {
		return 0
	}
	memo := newTable(len(weights), capacity+1, -1)
	var best func(i, c int) int
	best = func(i, c int) int {
		if i == len(weights) {
			return 0
		}
		if memo[i][c] != -1 {
			return memo[i][c]
		}
		res := best(i+1, c)
		if weights[i] <= c {
			res = max(res, values[i]+best(i+1, c-weights[i]))
		}
		memo[i][c] = res
		return res
	}
	return best(0, capacity)
}

// KnapsackTab is the 2-D bottom-up table.
func KnapsackTab(weights, values []int, capacity int) int {
	if capacity < 0 {
		return 0
	}
	n := len(weights)
	dp := newTable(n+1, capacity+1, 0)
	for i := 1; i <= n; i++ {
		w, v := weights[i-1], values[i-1]
		for c := 0; c <= capacity; c++ {
			dp[i][c] = dp[i-1][c]
			if w <= c {
				dp[i][c] = max(dp[i][c], dp[i-1][c-w]+v)
			}
		}
	}
	return dp[n][capacity]
}

// KnapsackRolling collapses the table to one row of capacities.
func KnapsackRolling(weights, values []int, capacity int) int {
	if capacity < 0 {
		return 0
	}
	dp := make([]int, capacity+1)
	for i, w := range weights {
		// High → low: dp[c-w] still holds the previous row's value.
		for c := capacity; c >= w; c-- {
			dp[c] = max(dp[c], dp[c-w]+values[i])
		}
	}
	return dp[capacity]
}
