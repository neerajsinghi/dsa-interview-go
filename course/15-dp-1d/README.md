# 15 · Dynamic programming (1-D)

## What you'll learn

- A repeatable **DP framework**: state → transition → base case → order → answer.
- The four-step refinement every DP interview rewards: **brute-force recursion → memoized top-down → bottom-up table → space-optimized**.
- The core 1-D families: **Fibonacci-style counting**, **take-or-skip** (House Robber), **unbounded vs 0/1 knapsack** (Coin Change vs Partition), **prefix segmentation** (Word Break, Decode Ways), **ends-at-i** (LIS), and **two-state Kadane** (Max Product).
- Why **loop order** and **loop direction** silently change what a table computes.

## Recognize the pattern

- The question asks for a **count** ("how many ways"), an **optimum** ("min cost", "max money", "longest"), or **feasibility** ("can it be split") — not the list of all solutions (that's backtracking).
- A brute-force recursion exists, and you can see it **solving the same subproblem many times** (draw two levels of the call tree and look for repeats).
- The decision at step `i` only depends on a **small summary** of the past (an index, an amount, "was the previous house robbed").

## The idea

Five questions, answered in order, give you any DP:

| Step | Question | Climbing Stairs |
|---|---|---|
| **State** | What does `dp[i]` mean, in one sentence? | ways to land exactly on step `i` |
| **Transition** | How is `dp[i]` built from smaller states? | `dp[i] = dp[i-1] + dp[i-2]` (last move was 1 or 2) |
| **Base case** | What is known without recursion? | `dp[0] = dp[1] = 1` |
| **Order** | Which order makes every dependency ready? | `i` increasing |
| **Answer** | Which state is the result? | `dp[n]` |

**Why it works — overlapping subproblems.** The naive recursion for `fib(5)` recomputes `fib(3)` twice and `fib(2)` three times:

```
                 fib(5)
              /          \
         fib(4)          fib(3)      <- computed again
         /    \          /    \
    fib(3)  fib(2)   fib(2)  fib(1)
    /   \
fib(2) fib(1)
```

Caching each answer collapses the tree into a line of `n` states — exponential becomes linear.

**Top-down vs bottom-up.**

| | Top-down (memo) | Bottom-up (table) |
|---|---|---|
| How | recursion + cache | loop over states in dependency order |
| Pros | write it straight from the recurrence; only visits reachable states | no recursion depth limit; easy to roll into O(1) space |
| Cons | recursion overhead; needs a "not computed" sentinel | you must work out the order yourself |

**Space optimization.** If `dp[i]` only reads `dp[i-1]` and `dp[i-2]`, keep two variables instead of the array (`prev, cur = cur, prev+cur`).

## Go template

```go
// Top-down
memo := make([]int, n+1)
for i := range memo {
	memo[i] = -1 // sentinel: not computed
}
var solve func(i int) int
solve = func(i int) int {
	if i <= 1 {
		return base(i)
	}
	if memo[i] == -1 {
		memo[i] = combine(solve(i-1), solve(i-2))
	}
	return memo[i]
}

// Bottom-up, rolled into two variables
prev, cur := base(0), base(1)
for i := 2; i <= n; i++ {
	prev, cur = cur, combine(cur, prev)
}
```

All four styles side by side, on Fibonacci, are in [dp1d.go](dp1d.go).

## Complexity

DP time = **(number of states) × (work per transition)**; space = number of states you keep.

| Problem | States | Per state | Time | Space (optimized) |
|---|---|---|---|---|
| Climbing stairs / House robber / Decode ways | `n` | `O(1)` | `O(n)` | `O(1)` |
| Coin change (I & II) | `amount` | `O(coins)` | `O(coins · amount)` | `O(amount)` |
| Word break | `n` | `O(L)` lookbacks | `O(n · L)` | `O(n)` |
| LIS | `n` | `O(n)` / `O(log n)` | `O(n²)` / `O(n log n)` | `O(n)` |
| Partition equal subset | `target` per item | `O(1)` | `O(n · target)` | `O(target)` |

## Common pitfalls

- **Sentinel collisions.** If `0` is a legitimate answer, initialise the memo to `-1` (or use a separate `done []bool`, as Min Cost Climbing Stairs does). For Coin Change, `-1` already means "impossible", so "not computed" needs `-2`.
- **Overflow in "infinity".** `math.MaxInt + 1` wraps negative. Use a large-but-safe constant (`1 << 62`) when you add to it.
- **Unbounded vs 0/1 knapsack direction.** Iterating amounts *upward* lets an item be reused (Coin Change); *downward* uses each item at most once (Partition). Getting this backwards is the most common 1-D DP bug.
- **Combinations vs permutations.** Coins in the outer loop counts combinations; amounts in the outer loop counts ordered sequences. Know which one the question wants.
- **Tuple assignment is your friend.** `prev, cur = cur, prev+cur` reads both old values before writing — no temp variable, no ordering bug. The same goes for `hi, lo = max(...), min(...)` in Max Product.
- **Off-by-one in the state definition.** Decide up front whether `dp[i]` covers `s[:i]` (prefix of length `i`) or `s[i]` (element `i`), and keep it consistent; most bugs come from mixing the two.
- **Recursion in Go has no tail-call optimization.** Memoized recursion over `n = 10^5` is fine (goroutine stacks grow), but bottom-up avoids the question entirely.
- **Map vs slice memo.** `map[int]int` is flexible but much slower than a slice for dense integer states. Prefer slices when the state space is a known range.
- **Min/max builtins.** Go 1.21+ has `min` and `max` for any ordered type, including more than two arguments: `max(x, hi*x, lo*x)`.

## Problems

| # | Problem | Difficulty | Key idea |
|---|---|---|---|
| 01 | [Climbing Stairs](p01_climbing_stairs.go) | Easy | Split by last move → Fibonacci |
| 02 | [Min Cost Climbing Stairs](p02_min_cost_climbing_stairs.go) | Easy | `min` over the last move; top is index `n` |
| 03 | [House Robber I & II](p03_house_robber.go) | Medium | Take-or-skip; circle = best of two lines |
| 04 | [Decode Ways](p04_decode_ways.go) | Medium | Stairs with validity checks on 1- and 2-digit pieces |
| 05 | [Coin Change I & II](p05_coin_change.go) | Medium | Unbounded knapsack; coins-outer counts combinations |
| 06 | [Word Break](p06_word_break.go) | Medium | Prefix DP: `ok[i] = ok[j] && s[j:i] ∈ dict` |
| 07 | [Maximum Product Subarray](p07_max_product_subarray.go) | Medium | Kadane with both max and min (negatives flip) |
| 08 | [Longest Increasing Subsequence](p08_longest_increasing_subsequence.go) | Medium | Ends-at-i `O(n²)`; patience sorting `O(n log n)` |
| 09 | [Partition Equal Subset Sum](p09_partition_equal_subset_sum.go) | Medium | 0/1 knapsack to `total/2`, iterate sums downward |

## Practice tip

For each problem, write the five framework answers as a comment *before* any code. Then implement the four versions in order — each one is a small edit of the previous — and run the agreement tests:

```bash
go test ./course/15-dp-1d/ -run TestCoinChange -v
```

Every problem's randomized test checks your memo, table and optimized versions against the brute-force recursion on small inputs.

## Next

[16 · Dynamic programming (2-D)](../16-dp-2d/README.md) — grids, two strings, and intervals.
