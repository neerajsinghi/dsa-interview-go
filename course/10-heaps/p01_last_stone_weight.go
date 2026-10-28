package heaps

import "slices"

// Problem: Last Stone Weight
// Difficulty: Easy
//
// You have stones with positive weights. Each turn, take the two heaviest
// stones x <= y and smash them: if x == y both vanish, otherwise one stone of
// weight y - x remains. When at most one stone is left, return its weight
// (0 if none remain).
//
// Examples:
//
//	[2,7,4,1,8,1] -> 1   (8,7 -> 1; 4,2 -> 2; 2,1 -> 1; 1,1 -> 0; left: 1)
//	[3,3]         -> 0
//
// Reference: https://leetcode.com/problems/last-stone-weight/
//
// Pattern: max-heap simulation. "Repeatedly take the largest (or smallest)
// and put something back" means a heap.
//
// Time: brute O(n^2 log n) (re-sort each turn); heap O(n log n). Space: O(n).

// LastStoneWeightBrute re-sorts the remaining stones every turn.
func LastStoneWeightBrute(stones []int) int {
	s := slices.Clone(stones)
	for len(s) > 1 {
		slices.Sort(s)
		y, x := s[len(s)-1], s[len(s)-2]
		s = s[:len(s)-2]
		if y != x {
			s = append(s, y-x)
		}
	}
	if len(s) == 0 {
		return 0
	}
	return s[0]
}

// LastStoneWeight uses a max-heap.
func LastStoneWeight(stones []int) int {
	h := FromSlice(slices.Clone(stones), func(a, b int) bool { return a > b }) // O(n) build
	for h.Len() > 1 {
		y, x := h.Pop(), h.Pop() // heaviest, second heaviest
		if y != x {
			h.Push(y - x)
		}
	}
	if h.Len() == 0 {
		return 0
	}
	return h.Peek()
}
