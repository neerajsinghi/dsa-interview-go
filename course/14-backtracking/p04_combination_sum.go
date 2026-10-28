package backtracking

import "sort"

// Problem: Combination Sum
// Difficulty: Medium
//
// Given distinct positive integers (candidates) and a target, return every
// distinct multiset of candidates that adds up to target. A candidate may be
// used any number of times; [2,2,3] and [3,2,2] count as the same answer.
//
// Examples:
//
//	candidates [2,3,6,7], target 7 -> [2 2 3] [7]
//	candidates [2,3,5],   target 8 -> [2 2 2 2] [2 3 3] [3 5]
//	candidates [2],       target 1 -> []
//
// Reference: https://leetcode.com/problems/combination-sum/
//
// Pattern: backtracking with a start index and reuse. Recursing with i
// (not i+1) lets the same candidate be picked again; never going back to
// smaller indices keeps each multiset in non-decreasing order, so it is
// generated once. Sorting the candidates enables the key prune: once
// candidates[i] exceeds what remains, every later candidate does too — break.
//
//	target 7, [2 3 6 7]:   rem 7
//	                   /   |    |   \
//	               +2(5) +3(4) +6(1) +7(0)✓
//	              /  |  \   |
//	          +2(3) +3(2) ...
//	          /  \     \
//	      +2(1) +3(0)✓  ...
//
// Time: O(n^(T/m)) in the worst case, where m is the smallest candidate
// (the depth is at most T/m). Space: O(T/m) recursion.

// CombinationSum returns all multisets summing to target, each sorted.
func CombinationSum(candidates []int, target int) [][]int {
	cands := append([]int(nil), candidates...)
	sort.Ints(cands)
	var res [][]int
	var path []int
	var backtrack func(start, remaining int)
	backtrack = func(start, remaining int) {
		if remaining == 0 {
			res = append(res, clone(path))
			return
		}
		for i := start; i < len(cands); i++ {
			if cands[i] > remaining {
				break // sorted: nothing further can fit either
			}
			path = append(path, cands[i])
			backtrack(i, remaining-cands[i]) // i, not i+1: reuse allowed
			path = path[:len(path)-1]
		}
	}
	backtrack(0, target)
	return res
}

// countCombinations counts the same multisets with the unbounded-knapsack DP
// (see module 15, Coin Change II). The tests use it as an independent check
// on how many answers CombinationSum should produce.
func countCombinations(candidates []int, target int) int {
	ways := make([]int, target+1)
	ways[0] = 1
	for _, c := range candidates { // coins outer => combinations, not orderings
		for s := c; s <= target; s++ {
			ways[s] += ways[s-c]
		}
	}
	return ways[target]
}
