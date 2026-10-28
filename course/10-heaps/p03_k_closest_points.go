package heaps

import "slices"

// Problem: K Closest Points to Origin
// Difficulty: Medium
//
// Given points on a plane as [x, y] pairs and an integer k, return the k
// points nearest to (0, 0) by Euclidean distance, in any order. If several
// points tie at the cut-off distance, any of them may be returned.
//
// Examples:
//
//	points=[[1,3],[-2,2]], k=1        -> [[-2,2]]   (dist² 8 < 10)
//	points=[[3,3],[5,-1],[-2,4]], k=2 -> [[3,3],[-2,4]]
//
// Reference: https://leetcode.com/problems/k-closest-points-to-origin/
//
// Pattern: top-k smallest uses a size-k MAX-heap. The heap holds the best k
// so far with the worst of them on top, ready to be evicted when something
// closer arrives. Compare squared distances: no sqrt, no floats.
//
// Time: sort O(n log n); heap O(n log k). Space: O(n) / O(k).

func dist2(p []int) int { return p[0]*p[0] + p[1]*p[1] }

// KClosestSort sorts by distance and takes the first k.
func KClosestSort(points [][]int, k int) [][]int {
	s := slices.Clone(points)
	slices.SortFunc(s, func(a, b []int) int { return dist2(a) - dist2(b) })
	return s[:k]
}

// KClosest keeps a max-heap (by distance) of size k.
func KClosest(points [][]int, k int) [][]int {
	h := New(func(a, b []int) bool { return dist2(a) > dist2(b) }) // farthest on top
	for _, p := range points {
		h.Push(p)
		if h.Len() > k {
			h.Pop()
		}
	}
	out := make([][]int, 0, h.Len())
	for h.Len() > 0 {
		out = append(out, h.Pop())
	}
	return out
}
