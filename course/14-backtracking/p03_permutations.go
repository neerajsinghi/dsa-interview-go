package backtracking

import "sort"

// Problem: Permutations I & II
// Difficulty: Medium
//
// Part I: return every ordering of a list of distinct integers.
// Part II: the list may contain repeated values; return every DISTINCT
// ordering exactly once. Any order of results is fine.
//
// Examples:
//
//	I:  [1,2,3] -> [1 2 3] [1 3 2] [2 1 3] [2 3 1] [3 1 2] [3 2 1]
//	II: [1,1,2] -> [1 1 2] [1 2 1] [2 1 1]
//
// References: https://leetcode.com/problems/permutations/
//
//	https://leetcode.com/problems/permutations-ii/
//
// Pattern: at each position pick any element not used yet. A used[] array
// marks what the current path already holds; answers are the leaves at depth
// n. (Contrast with subsets/combinations, which use a start index because
// order does not matter there.)
//
// Duplicates: sort, then at each depth skip nums[i] if it equals nums[i-1]
// and nums[i-1] is NOT currently used. That rule forces equal values to be
// placed in their original left-to-right order, so each distinct arrangement
// is produced by exactly one sequence of choices.
//
// Time: O(n · n!). Space: O(n) recursion + used[] besides the output.

// Permute returns all permutations of distinct nums.
func Permute(nums []int) [][]int {
	var res [][]int
	path := make([]int, 0, len(nums))
	used := make([]bool, len(nums))
	var backtrack func()
	backtrack = func() {
		if len(path) == len(nums) {
			res = append(res, clone(path))
			return
		}
		for i, x := range nums {
			if used[i] {
				continue
			}
			used[i] = true
			path = append(path, x)
			backtrack()
			path = path[:len(path)-1]
			used[i] = false
		}
	}
	backtrack()
	return res
}

// PermuteSwap is the in-place alternative: fix position i by swapping each
// candidate into it, recurse on i+1, then swap back. No used[] or path.
func PermuteSwap(nums []int) [][]int {
	a := append([]int(nil), nums...)
	var res [][]int
	var backtrack func(i int)
	backtrack = func(i int) {
		if i == len(a) {
			res = append(res, clone(a))
			return
		}
		for j := i; j < len(a); j++ {
			a[i], a[j] = a[j], a[i]
			backtrack(i + 1)
			a[i], a[j] = a[j], a[i] // undo the swap
		}
	}
	backtrack(0)
	return res
}

// PermuteUnique returns each distinct permutation of nums once.
func PermuteUnique(nums []int) [][]int {
	sorted := append([]int(nil), nums...)
	sort.Ints(sorted)
	var res [][]int
	path := make([]int, 0, len(sorted))
	used := make([]bool, len(sorted))
	var backtrack func()
	backtrack = func() {
		if len(path) == len(sorted) {
			res = append(res, clone(path))
			return
		}
		for i, x := range sorted {
			if used[i] {
				continue
			}
			// Equal values must be used left to right: if the previous copy
			// is unused here, choosing this copy would duplicate that branch.
			if i > 0 && sorted[i] == sorted[i-1] && !used[i-1] {
				continue
			}
			used[i] = true
			path = append(path, x)
			backtrack()
			path = path[:len(path)-1]
			used[i] = false
		}
	}
	backtrack()
	return res
}
