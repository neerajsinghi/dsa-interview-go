package greedy

import "math"

// Problem: Non-overlapping Intervals
// Difficulty: Medium
//
// Return the fewest intervals you must delete so the remaining ones do not
// overlap. Here intervals that merely touch ([1,2] and [2,3]) are fine.
//
// Examples:
//
//	[1,2] [2,3] [3,4] [1,3] -> 1  (drop [1,3])
//	[1,2] [1,2] [1,2]       -> 2
//
// Reference: https://leetcode.com/problems/non-overlapping-intervals/
//
// Pattern: interval scheduling. Removing the fewest = keeping the most.
// Sort by END and greedily keep each interval that starts at or after the
// last kept end. Exchange argument: the interval that ends first leaves the
// most room; swapping it into any optimal schedule keeps it valid.
//
// Time: O(n log n). Space: O(n) for the sorted copy.

// EraseOverlapIntervalsBrute tries every subset to keep (n ≤ ~15). Tests only.
func EraseOverlapIntervalsBrute(in []Interval) int {
	n, bestKeep := len(in), 0
	for mask := 0; mask < 1<<n; mask++ {
		var kept []Interval
		for i := 0; i < n; i++ {
			if mask&(1<<i) != 0 {
				kept = append(kept, in[i])
			}
		}
		kept = sortedByStart(kept)
		ok := true
		for i := 1; i < len(kept); i++ {
			if kept[i].Start < kept[i-1].End {
				ok = false
				break
			}
		}
		if ok {
			bestKeep = max(bestKeep, len(kept))
		}
	}
	return n - bestKeep
}

// EraseOverlapIntervals is the earliest-end-first greedy.
func EraseOverlapIntervals(in []Interval) int {
	kept, lastEnd := 0, math.MinInt
	for _, iv := range sortedByEnd(in) {
		if iv.Start >= lastEnd { // compatible with everything kept so far
			kept++
			lastEnd = iv.End
		}
	}
	return len(in) - kept
}
