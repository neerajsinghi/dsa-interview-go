package heaps

import "slices"

// Problem: Find Median from Data Stream
// Difficulty: Hard
//
// Build a structure that receives integers one at a time (AddNum) and can
// report the median of everything received so far (FindMedian) at any point.
// With an even count the median is the mean of the two middle values.
// FindMedian is only called after at least one AddNum.
//
// Example:
//
//	AddNum(1) AddNum(2)  FindMedian() -> 1.5
//	AddNum(3)            FindMedian() -> 2.0
//
// Reference: https://leetcode.com/problems/find-median-from-data-stream/
//
// Pattern: two heaps split the data in half. A max-heap `low` holds the
// smaller half, a min-heap `high` holds the larger half, and we keep
// len(low) == len(high) or len(low) == len(high)+1. The median is then
// low's top, or the mean of both tops.
//
//	low (max-heap)      high (min-heap)
//	[1 2 3] top=3   |   top=4 [4 5]        median = 3
//
// Time: AddNum O(log n), FindMedian O(1). The sorted-slice version is
// O(n) per AddNum (binary search, then shift). Space: O(n).

// MedianFinder maintains a running median with two heaps.
type MedianFinder struct {
	low  *Heap[int] // max-heap: smaller half
	high *Heap[int] // min-heap: larger half
}

// NewMedianFinder returns an empty MedianFinder.
func NewMedianFinder() *MedianFinder {
	return &MedianFinder{low: NewMax[int](), high: NewMin[int]()}
}

// AddNum inserts v.
func (m *MedianFinder) AddNum(v int) {
	// Always route through low so that low's max ends up <= high's min.
	m.low.Push(v)
	m.high.Push(m.low.Pop())
	// Rebalance: low may be larger by at most one.
	if m.high.Len() > m.low.Len() {
		m.low.Push(m.high.Pop())
	}
}

// FindMedian returns the current median.
func (m *MedianFinder) FindMedian() float64 {
	if m.low.Len() > m.high.Len() {
		return float64(m.low.Peek())
	}
	return (float64(m.low.Peek()) + float64(m.high.Peek())) / 2
}

// SortedMedianFinder is the simple baseline: keep a sorted slice.
type SortedMedianFinder struct {
	xs []int
}

// AddNum inserts v at its sorted position. O(n) due to shifting.
func (s *SortedMedianFinder) AddNum(v int) {
	i, _ := slices.BinarySearch(s.xs, v)
	s.xs = slices.Insert(s.xs, i, v)
}

// FindMedian returns the current median.
func (s *SortedMedianFinder) FindMedian() float64 {
	n := len(s.xs)
	if n%2 == 1 {
		return float64(s.xs[n/2])
	}
	return (float64(s.xs[n/2-1]) + float64(s.xs[n/2])) / 2
}
