package twopointers

import (
	"slices"
)

// Problem 05 — 3Sum (Medium)
//
// Return every distinct triplet [a, b, c] of values from nums (taken at three
// different positions) with a + b + c == 0. Each triplet is sorted
// ascending, the list of triplets is in lexicographic order, and no triplet
// appears twice.
//
//	[-1, 0, 1, 2, -1, -4] → [[-1, -1, 2], [-1, 0, 1]]
//	[0, 0, 0, 0]          → [[0, 0, 0]]
//
// Reference: https://leetcode.com/problems/3sum/
//
// Pattern: SORT, then fix the first element and run Two Sum II (opposite-ends
// two pointers) on the rest. Sorting also makes deduplication easy: skip a
// value equal to the one you just used, at every level.
//
//	sorted: [-4, -1, -1, 0, 1, 2]
//	i=-4:   l→ … ←r   no pair sums to 4
//	i=-1:   [-1, 2] ✓, then [0, 1] ✓
//	i=-1 (again): skipped — same first value
//
//	ThreeSumBrute: Time O(n³), Space O(answers)
//	ThreeSum:      Time O(n²), Space O(1) beyond the output (sorting a copy is O(n))

// ThreeSumBrute tries every triple and deduplicates with a set of sorted
// triplets ([3]int arrays are comparable, so they work as map keys).
func ThreeSumBrute(nums []int) [][]int {
	seen := map[[3]int]bool{}
	res := [][]int{}
	for i := range nums {
		for j := i + 1; j < len(nums); j++ {
			for k := j + 1; k < len(nums); k++ {
				if nums[i]+nums[j]+nums[k] != 0 {
					continue
				}
				t := [3]int{nums[i], nums[j], nums[k]}
				slices.Sort(t[:])
				if !seen[t] {
					seen[t] = true
					res = append(res, t[:])
				}
			}
		}
	}
	slices.SortFunc(res, slices.Compare)
	return res
}

// ThreeSum sorts a copy and runs two pointers for each first element.
func ThreeSum(nums []int) [][]int {
	s := slices.Clone(nums) // don't reorder the caller's slice
	slices.Sort(s)
	res := [][]int{}
	for i := 0; i < len(s)-2; i++ {
		if s[i] > 0 {
			break // smallest value positive ⇒ no sum can reach 0
		}
		if i > 0 && s[i] == s[i-1] {
			continue // same first value ⇒ same triplets
		}
		l, r := i+1, len(s)-1
		for l < r {
			switch sum := s[i] + s[l] + s[r]; {
			case sum < 0:
				l++
			case sum > 0:
				r--
			default:
				res = append(res, []int{s[i], s[l], s[r]})
				l++
				r--
				for l < r && s[l] == s[l-1] {
					l++ // skip duplicate second values
				}
			}
		}
	}
	return res
}
