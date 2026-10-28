package greedy

// Problem: Merge Intervals
// Difficulty: Medium
//
// Given closed intervals in any order, merge every group that overlaps
// (touching end-to-start counts as overlapping) and return the merged list
// sorted by start.
//
// Examples:
//
//	[1,3] [2,6] [8,10] [15,18] -> [1,6] [8,10] [15,18]
//	[1,4] [4,5]                -> [1,5]
//
// Reference: https://leetcode.com/problems/merge-intervals/
//
// Pattern: sort by start, then sweep. After sorting, an interval can only
// overlap the LAST merged interval — anything earlier ended before it.
//
//	sorted:  [1,3] [2,6] [8,10] [15,18]
//	         └──┬──┘
//	          [1,6]       [8,10] [15,18]
//
// Time: O(n log n) for the sort. Space: O(n) for the output (and the copy).

// MergeIntervalsBrute repeatedly fuses any overlapping pair until none
// remain. O(n³); tests only.
func MergeIntervalsBrute(in []Interval) []Interval {
	cur := append([]Interval{}, in...)
	for changed := true; changed; {
		changed = false
	outer:
		for i := range cur {
			for j := i + 1; j < len(cur); j++ {
				a, b := cur[i], cur[j]
				if a.Start <= b.End && b.Start <= a.End {
					cur[i] = Interval{min(a.Start, b.Start), max(a.End, b.End)}
					cur = append(cur[:j], cur[j+1:]...)
					changed = true
					break outer
				}
			}
		}
	}
	return sortedByStart(cur)
}

// MergeIntervals sorts a copy by start and sweeps once.
func MergeIntervals(in []Interval) []Interval {
	if len(in) == 0 {
		return nil
	}
	sorted := sortedByStart(in)
	out := []Interval{sorted[0]}
	for _, iv := range sorted[1:] {
		last := &out[len(out)-1]
		if iv.Start <= last.End { // overlap (or touch): extend
			last.End = max(last.End, iv.End) // max: [1,10] swallows [2,3]
		} else {
			out = append(out, iv)
		}
	}
	return out
}
