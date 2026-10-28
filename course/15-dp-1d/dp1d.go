// Package dp1d covers one-dimensional dynamic programming: problems whose
// state is a single index or amount, solved four ways — recursive brute
// force, memoized top-down, bottom-up tabulation and space-optimized.
//
// The DP framework used throughout:
//
//  1. State       — what does dp[i] mean, in one sentence?
//  2. Transition  — how is dp[i] built from smaller states?
//  3. Base case   — which states are known without recursion?
//  4. Order       — in which order must states be computed so that every
//     dependency is ready (bottom-up), or just recurse and cache (top-down)?
//  5. Answer      — which state (or combination of states) is the result?
//
// Fibonacci below is the smallest example of all four implementations.
package dp1d

// FibRecursive is the brute force straight from the definition
// F(n) = F(n-1) + F(n-2). It recomputes the same subproblems exponentially
// often: F(n-2) is computed by both F(n) and F(n-1), and so on down.
//
// Time: O(φ^n) ≈ O(1.618^n). Space: O(n) recursion.
func FibRecursive(n int) int {
	if n < 2 {
		return n
	}
	return FibRecursive(n-1) + FibRecursive(n-2)
}

// FibMemo is top-down DP: the same recursion, but each F(i) is computed once
// and cached. Recursion naturally visits states in a valid order.
//
// Time: O(n). Space: O(n) memo + O(n) recursion.
func FibMemo(n int) int {
	memo := make(map[int]int)
	var f func(i int) int
	f = func(i int) int {
		if i < 2 {
			return i
		}
		if v, ok := memo[i]; ok {
			return v
		}
		memo[i] = f(i-1) + f(i-2)
		return memo[i]
	}
	return f(n)
}

// FibTab is bottom-up DP: fill dp[0..n] in increasing order so that
// dp[i-1] and dp[i-2] are always ready. No recursion.
//
// Time: O(n). Space: O(n).
func FibTab(n int) int {
	if n < 2 {
		return n
	}
	dp := make([]int, n+1)
	dp[1] = 1
	for i := 2; i <= n; i++ {
		dp[i] = dp[i-1] + dp[i-2]
	}
	return dp[n]
}

// FibOptimized notices that dp[i] only reads the previous two entries, so
// two variables replace the table. This "rolling" trick applies to most
// 1-D DPs with a fixed look-back window.
//
// Time: O(n). Space: O(1).
func FibOptimized(n int) int {
	prev, cur := 0, 1 // F(0), F(1)
	if n == 0 {
		return 0
	}
	for i := 2; i <= n; i++ {
		prev, cur = cur, prev+cur // tuple assignment: both read old values
	}
	return cur
}
