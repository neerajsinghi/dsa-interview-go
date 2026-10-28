package heaps

import (
	"math/rand/v2"
	"slices"
)

// Problem: Kth Largest Element in an Array
// Difficulty: Medium
//
// Return the k-th largest value in an unsorted array, counting duplicates
// separately (so in [3,3,1] the 2nd largest is 3). 1 <= k <= len(nums).
//
// Examples:
//
//	nums=[3,2,1,5,6,4], k=2       -> 5
//	nums=[3,2,3,1,2,4,5,5,6], k=4 -> 4
//
// Reference: https://leetcode.com/problems/kth-largest-element-in-an-array/
//
// Pattern: two classic answers.
//   - Size-k min-heap: keep the k largest seen so far. Its top is the
//     smallest of them, i.e. the k-th largest. Great for streams and k << n.
//   - Quickselect: partition like quicksort but recurse into one side only.
//     Expected O(n) with a random pivot. Three-way partitioning keeps it fast
//     when there are many duplicates.
//
// Time: sort O(n log n); heap O(n log k); quickselect O(n) expected, O(n^2)
// worst. Space: O(1) sort-in-place aside, O(k) heap, O(1) quickselect.

// FindKthLargestSort sorts a copy. It's the simplest correct answer and the baseline.
func FindKthLargestSort(nums []int, k int) int {
	s := slices.Clone(nums)
	slices.Sort(s)
	return s[len(s)-k]
}

// FindKthLargestHeap keeps a min-heap of the k largest elements.
func FindKthLargestHeap(nums []int, k int) int {
	h := NewMin[int]()
	for _, v := range nums {
		h.Push(v)
		if h.Len() > k {
			h.Pop() // evict the smallest; it can't be in the top k
		}
	}
	return h.Peek()
}

// FindKthLargestQuickselect uses randomized three-way quickselect on a copy.
func FindKthLargestQuickselect(nums []int, k int) int {
	a := slices.Clone(nums)
	target := len(a) - k // index of the answer in ascending order
	lo, hi := 0, len(a)-1
	for {
		pivot := a[lo+rand.IntN(hi-lo+1)] // global rand/v2 is safe for concurrent use
		// Dutch national flag: [lo,lt) < pivot, [lt,i) == pivot, (gt,hi] > pivot.
		lt, i, gt := lo, lo, hi
		for i <= gt {
			switch {
			case a[i] < pivot:
				a[lt], a[i] = a[i], a[lt]
				lt++
				i++
			case a[i] > pivot:
				a[i], a[gt] = a[gt], a[i]
				gt-- // don't advance i: the swapped-in value is unexamined
			default:
				i++
			}
		}
		switch {
		case target < lt:
			hi = lt - 1
		case target > gt:
			lo = gt + 1
		default:
			return pivot // target falls in the == pivot block
		}
	}
}
