package greedy

// Problem: Insert Interval
// Difficulty: Medium
//
// You are given closed intervals that are already sorted by start and do not
// overlap. Insert one new interval, merging where necessary, and return a
// list that is still sorted and non-overlapping.
//
// Examples:
//
//	[1,3] [6,9] + [2,5]                    -> [1,5] [6,9]
//	[1,2] [3,5] [6,7] [8,10] [12,16] + [4,8] -> [1,2] [3,10] [12,16]
//
// Reference: https://leetcode.com/problems/insert-interval/
//
// Pattern: three-phase linear sweep — no sort needed because input is sorted.
//
//	1. copy intervals that end before the new one starts
//	2. absorb every interval that overlaps it (grow new to cover them)
//	3. append the merged interval, then copy the rest
//
// Time: O(n). Space: O(n) output.

// InsertInterval returns a new slice; the input is not modified.
func InsertInterval(intervals []Interval, add Interval) []Interval {
	out := make([]Interval, 0, len(intervals)+1)
	i, n := 0, len(intervals)
	for i < n && intervals[i].End < add.Start { // 1. strictly before
		out = append(out, intervals[i])
		i++
	}
	for i < n && intervals[i].Start <= add.End { // 2. overlapping
		add.Start = min(add.Start, intervals[i].Start)
		add.End = max(add.End, intervals[i].End)
		i++
	}
	out = append(out, add)               // 3. merged interval
	return append(out, intervals[i:]...) //    and everything after
}
