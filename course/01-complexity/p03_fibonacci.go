package complexity

// Problem 03 — Fibonacci: exponential vs linear (Easy)
//
// Return the n-th Fibonacci number, where F(0) = 0, F(1) = 1 and
// F(n) = F(n-1) + F(n-2). Assume 0 ≤ n ≤ 92 (F(93) overflows int64).
//
//	n = 2  → 1
//	n = 10 → 55
//
// Reference: https://leetcode.com/problems/fibonacci-number/
//
// Pattern: draw the recursion tree. The naive version recomputes the same
// sub-problems over and over — the tree has ~φⁿ ≈ 1.618ⁿ nodes:
//
//	                 F(5)
//	           /             \
//	       F(4)               F(3)
//	      /    \             /    \
//	   F(3)    F(2)       F(2)   F(1)
//	   /  \    /  \       /  \
//	F(2) F(1) F(1) F(0) F(1) F(0)
//	...
//
// Caching each F(k) the first time collapses the tree to a single path of n
// distinct calls (memoization). Computing bottom-up removes the recursion
// stack too.
//
//	FibNaive: Time O(φⁿ), Space O(n) recursion stack
//	FibMemo:  Time O(n),  Space O(n) cache + O(n) stack
//	FibIter:  Time O(n),  Space O(1)

// FibNaive is the direct translation of the definition.
func FibNaive(n int) int {
	if n < 2 {
		return n
	}
	return FibNaive(n-1) + FibNaive(n-2)
}

// FibMemo remembers every answer it computes (top-down DP).
func FibMemo(n int) int {
	memo := make(map[int]int, n+1)
	var fib func(int) int
	fib = func(k int) int {
		if k < 2 {
			return k
		}
		if v, ok := memo[k]; ok {
			return v
		}
		memo[k] = fib(k-1) + fib(k-2)
		return memo[k]
	}
	return fib(n)
}

// FibIter keeps only the last two values (bottom-up DP).
func FibIter(n int) int {
	prev, cur := 0, 1 // F(0), F(1)
	for range n {
		prev, cur = cur, prev+cur
	}
	return prev
}
