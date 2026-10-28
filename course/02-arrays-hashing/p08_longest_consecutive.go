package arrayshashing

import "slices"

// Problem 08 — Longest Consecutive Sequence (Medium)
//
// Given an unsorted slice of integers, return the length of the longest run
// of consecutive values (x, x+1, x+2, …) that all appear in it. Order in the
// slice doesn't matter and duplicates count once.
//
//	[100, 4, 200, 1, 3, 2]            → 4   (1, 2, 3, 4)
//	[0, 3, 7, 2, 5, 8, 4, 6, 0, 1]    → 9   (0 … 8)
//
// Reference: https://leetcode.com/problems/longest-consecutive-sequence/
//
// Pattern: put everything in a set, then only start counting at the BEGINNING
// of a run — a value x whose predecessor x-1 is absent. Each value is then
// visited at most twice (once in the outer loop, once while extending a run),
// so the total is O(n) even though there is a nested loop.
//
//	set: {1, 2, 3, 4, 100, 200}
//	  1: 0 absent → count 1,2,3,4 → 4
//	  2: 1 present → skip (not a start)
//	100: 99 absent → count 100 → 1
//
//	LongestConsecutiveSort: Time O(n log n), Space O(n)
//	LongestConsecutive:     Time O(n),       Space O(n)

// LongestConsecutiveSort sorts a copy and scans for runs.
func LongestConsecutiveSort(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	s := slices.Clone(nums)
	slices.Sort(s)
	best, run := 1, 1
	for i := 1; i < len(s); i++ {
		switch s[i] - s[i-1] {
		case 0:
			// duplicate: run unchanged
		case 1:
			run++
			best = max(best, run)
		default:
			run = 1
		}
	}
	return best
}

// LongestConsecutive uses a set and only expands from run starts.
func LongestConsecutive(nums []int) int {
	set := NewSet(nums)
	best := 0
	for x := range set { // iterate the set, not nums, so duplicates are skipped
		if set.Has(x - 1) {
			continue // x is inside a run that starts earlier
		}
		length := 1
		for set.Has(x + length) {
			length++
		}
		best = max(best, length)
	}
	return best
}
