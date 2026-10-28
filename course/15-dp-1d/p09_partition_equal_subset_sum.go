package dp1d

// Problem: Partition Equal Subset Sum
// Difficulty: Medium
//
// Given positive integers, decide whether they can be split into two groups
// with equal sums.
//
// Examples:
//
//	[1, 5, 11, 5] -> true   ({1 5 5} and {11})
//	[1, 2, 3, 5]  -> false  (total 11 is odd)
//	[1, 2, 5]     -> false  (total 8, but no subset sums to 4)
//
// Reference: https://leetcode.com/problems/partition-equal-subset-sum/
//
// Pattern: 0/1 knapsack, reduced to "subset sum". Equal halves exist iff the
// total is even and some subset sums to target = total/2.
//
//	State:      can[s] = some subset of the items seen so far sums to s
//	Transition: for each item x: can[s] = can[s] OR can[s-x]
//	Base:       can[0] = true (the empty subset)
//	Order:      items outer; s from target DOWN to x
//
// Why downward? Each item may be used at most once. Iterating s downward
// means can[s-x] still holds the value from BEFORE this item was considered.
// Iterating upward would let the same item be added again and again — that
// is the unbounded knapsack of Coin Change, a different problem.
//
//	nums [1 5 11 5], target 11
//	after 1:  sums {0 1}
//	after 5:  sums {0 1 5 6}
//	after 11: sums {0 1 5 6 11}  -> can[11] = true
//
// Time: O(2^n) brute force; O(n · target) DP. Space: O(n · target) memo,
// O(target) for the 1-D table.

// CanPartitionRecursive decides take/skip for each item.
func CanPartitionRecursive(nums []int) bool {
	total := sum(nums)
	if total%2 != 0 {
		return false
	}
	var can func(i, remaining int) bool
	can = func(i, remaining int) bool {
		if remaining == 0 {
			return true
		}
		if i == len(nums) || remaining < 0 {
			return false
		}
		return can(i+1, remaining-nums[i]) || can(i+1, remaining)
	}
	return can(0, total/2)
}

// CanPartitionMemo caches can(i, remaining): 0 unknown, 1 true, 2 false.
func CanPartitionMemo(nums []int) bool {
	total := sum(nums)
	if total%2 != 0 {
		return false
	}
	target := total / 2
	memo := make([][]int8, len(nums))
	for i := range memo {
		memo[i] = make([]int8, target+1)
	}
	var can func(i, remaining int) bool
	can = func(i, remaining int) bool {
		if remaining == 0 {
			return true
		}
		if i == len(nums) || remaining < 0 {
			return false
		}
		if memo[i][remaining] == 0 {
			memo[i][remaining] = 2
			if can(i+1, remaining-nums[i]) || can(i+1, remaining) {
				memo[i][remaining] = 1
			}
		}
		return memo[i][remaining] == 1
	}
	return can(0, target)
}

// CanPartitionTab is the 2-D table: dp[i][s] = some subset of nums[:i] sums to s.
func CanPartitionTab(nums []int) bool {
	total := sum(nums)
	if total%2 != 0 {
		return false
	}
	target := total / 2
	dp := make([][]bool, len(nums)+1)
	for i := range dp {
		dp[i] = make([]bool, target+1)
		dp[i][0] = true
	}
	for i := 1; i <= len(nums); i++ {
		x := nums[i-1]
		for s := 1; s <= target; s++ {
			dp[i][s] = dp[i-1][s] || (s >= x && dp[i-1][s-x])
		}
	}
	return dp[len(nums)][target]
}

// CanPartition is the 1-D space-optimized version.
func CanPartition(nums []int) bool {
	total := sum(nums)
	if total%2 != 0 {
		return false
	}
	target := total / 2
	can := make([]bool, target+1)
	can[0] = true
	for _, x := range nums {
		for s := target; s >= x; s-- { // downward: each item used at most once
			can[s] = can[s] || can[s-x]
		}
		if can[target] {
			return true // early exit
		}
	}
	return can[target]
}

func sum(nums []int) int {
	total := 0
	for _, x := range nums {
		total += x
	}
	return total
}
