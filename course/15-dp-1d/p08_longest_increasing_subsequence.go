package dp1d

import "sort"

// Problem: Longest Increasing Subsequence
// Difficulty: Medium
//
// Return the length of the longest subsequence (not necessarily contiguous)
// whose values are STRICTLY increasing.
//
// Examples:
//
//	[10, 9, 2, 5, 3, 7, 101, 18] -> 4   (2 3 7 18, or 2 5 7 101, …)
//	[0, 1, 0, 3, 2, 3]           -> 4
//	[7, 7, 7, 7]                 -> 1
//
// Reference: https://leetcode.com/problems/longest-increasing-subsequence/
//
// Pattern 1 — O(n²) DP, "ends at i":
//
//	State:      lis[i] = length of the longest increasing subsequence ENDING at i
//	Transition: lis[i] = 1 + max(lis[j]) over j < i with nums[j] < nums[i]
//	Base:       lis[i] = 1 (the element alone)
//	Answer:     max over all i (the best one can end anywhere)
//
// Pattern 2 — O(n log n) patience sorting. Keep tails[k] = the smallest
// possible tail value of an increasing subsequence of length k+1. tails is
// always sorted, so for each x binary-search the first tail ≥ x:
//   - none: x extends the longest subsequence — append it;
//   - otherwise: replace that tail with x (same length, smaller tail, which
//     can only help later elements).
//
// len(tails) is the answer. (tails itself is NOT an actual subsequence.)
//
//	x:      10   9    2    5    3    7    101   18
//	tails: [10] [9]  [2]  [2 5][2 3][2 3 7][2 3 7 101][2 3 7 18]   -> 4
//
// Time: O(2^n) brute force; O(n²) DP; O(n log n) patience. Space: O(n).

// LengthOfLISRecursive decides take/skip for each element given the last
// value taken (prev = -1 means nothing taken yet).
func LengthOfLISRecursive(nums []int) int {
	var best func(i, prev int) int
	best = func(i, prev int) int {
		if i == len(nums) {
			return 0
		}
		res := best(i+1, prev) // skip nums[i]
		if prev == -1 || nums[prev] < nums[i] {
			res = max(res, 1+best(i+1, i)) // take nums[i]
		}
		return res
	}
	return best(0, -1)
}

// LengthOfLISMemo caches best(i, prev); prev+1 shifts -1 to index 0.
func LengthOfLISMemo(nums []int) int {
	n := len(nums)
	memo := make([][]int, n)
	for i := range memo {
		memo[i] = make([]int, n+1)
		for j := range memo[i] {
			memo[i][j] = -1
		}
	}
	var best func(i, prev int) int
	best = func(i, prev int) int {
		if i == n {
			return 0
		}
		if memo[i][prev+1] != -1 {
			return memo[i][prev+1]
		}
		res := best(i+1, prev)
		if prev == -1 || nums[prev] < nums[i] {
			res = max(res, 1+best(i+1, i))
		}
		memo[i][prev+1] = res
		return res
	}
	return best(0, -1)
}

// LengthOfLISQuadratic is the classic O(n²) bottom-up table.
func LengthOfLISQuadratic(nums []int) int {
	lis := make([]int, len(nums))
	best := 0
	for i := range nums {
		lis[i] = 1
		for j := 0; j < i; j++ {
			if nums[j] < nums[i] {
				lis[i] = max(lis[i], lis[j]+1)
			}
		}
		best = max(best, lis[i])
	}
	return best
}

// LengthOfLIS is the O(n log n) patience-sorting solution.
func LengthOfLIS(nums []int) int {
	var tails []int
	for _, x := range nums {
		// First tail >= x. Using >= (not >) is what makes it STRICTLY
		// increasing: an equal value replaces, it does not extend.
		k := sort.SearchInts(tails, x)
		if k == len(tails) {
			tails = append(tails, x)
		} else {
			tails[k] = x
		}
	}
	return len(tails)
}
