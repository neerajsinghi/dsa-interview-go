package backtracking

import "sort"

// Problem: Subsets I & II
// Difficulty: Medium
//
// Part I: given distinct integers, return every subset (the power set).
// Part II: the integers may repeat; return every subset exactly once —
// {1,2} built from the first 2 and {1,2} built from the second 2 count as
// the same subset. Any order is fine.
//
// Examples:
//
//	I:  [1,2,3]  -> [] [1] [2] [3] [1 2] [1 3] [2 3] [1 2 3]
//	II: [1,2,2]  -> [] [1] [2] [1 2] [2 2] [1 2 2]
//
// References: https://leetcode.com/problems/subsets/
//
//	https://leetcode.com/problems/subsets-ii/
//
// Pattern: "every node of the recursion tree is an answer". Unlike
// permutations, we record the current path on EVERY call, not only at the
// leaves. Passing a start index means each element is only considered after
// the ones before it, which avoids generating {2,1} as well as {1,2}.
//
// Duplicates: sort first so equal values are adjacent; then, at a given
// depth, skip an element equal to the previous one at that same depth
// (i > start && nums[i] == nums[i-1]). The first copy already explored every
// subset that the later copy would start.
//
//	[1,2,2] sorted:      []
//	                  /    \      \
//	                [1]    [2]    [2] <- skipped (same value, same depth)
//	               /   \     \
//	           [1 2]  [1 2]x [2 2]
//	             |
//	          [1 2 2]
//
// Time: O(n · 2^n). Space: O(n) recursion besides the output.

// Subsets enumerates the power set by backtracking.
func Subsets(nums []int) [][]int {
	var res [][]int
	var path []int
	var backtrack func(start int)
	backtrack = func(start int) {
		res = append(res, clone(path)) // every partial path is a subset
		for i := start; i < len(nums); i++ {
			path = append(path, nums[i])
			backtrack(i + 1)
			path = path[:len(path)-1]
		}
	}
	backtrack(0)
	return res
}

// SubsetsBitmask is the iterative alternative: the bits of mask 0..2^n-1
// say which elements to include. Simple, but only works for n < 64.
func SubsetsBitmask(nums []int) [][]int {
	n := len(nums)
	res := make([][]int, 0, 1<<n)
	for mask := 0; mask < 1<<n; mask++ {
		sub := []int{}
		for i := 0; i < n; i++ {
			if mask&(1<<i) != 0 {
				sub = append(sub, nums[i])
			}
		}
		res = append(res, sub)
	}
	return res
}

// SubsetsWithDup returns each distinct subset once.
func SubsetsWithDup(nums []int) [][]int {
	sorted := append([]int(nil), nums...) // don't reorder the caller's slice
	sort.Ints(sorted)
	var res [][]int
	var path []int
	var backtrack func(start int)
	backtrack = func(start int) {
		res = append(res, clone(path))
		for i := start; i < len(sorted); i++ {
			if i > start && sorted[i] == sorted[i-1] {
				continue // same value already tried at this depth
			}
			path = append(path, sorted[i])
			backtrack(i + 1)
			path = path[:len(path)-1]
		}
	}
	backtrack(0)
	return res
}
