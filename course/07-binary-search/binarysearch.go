// Package binarysearch covers binary search over sorted data and over a
// monotone answer space ("binary search on the answer").
//
// Everything here uses one convention: the half-open interval [lo, hi). The
// loop runs while lo < hi, the answer is lo when the loop exits, and the loop
// never needs a separate "found" check. Once you learn one template well you
// can derive every variant (first/last position, insert position, minimum
// feasible speed, ...) from it.
//
// The standard library has the same tools:
//
//	sort.Search(n, func(i int) bool)            // first i in [0,n) where f(i) is true
//	slices.BinarySearch(xs, target) (int, bool) // lower bound + found flag
//	slices.BinarySearchFunc(xs, target, cmp)
//
// Use them in production. In an interview, write the loop yourself unless the
// interviewer says otherwise, because the loop is the part they are testing.
package binarysearch

import "cmp"

// LowerBound returns the first index i such that xs[i] >= target, or len(xs)
// if every element is smaller. xs must be sorted ascending. It is also the
// index where target would be inserted to keep xs sorted.
func LowerBound[T cmp.Ordered](xs []T, target T) int {
	lo, hi := 0, len(xs) // answer lies in [lo, hi]
	for lo < hi {
		mid := lo + (hi-lo)/2 // never overflows, unlike (lo+hi)/2
		if xs[mid] < target {
			lo = mid + 1 // mid is too small: answer is to the right
		} else {
			hi = mid // mid could be the answer: keep it
		}
	}
	return lo
}

// UpperBound returns the first index i such that xs[i] > target, or len(xs).
// UpperBound(xs, t) - LowerBound(xs, t) is the number of copies of t.
func UpperBound[T cmp.Ordered](xs []T, target T) int {
	lo, hi := 0, len(xs)
	for lo < hi {
		mid := lo + (hi-lo)/2
		if xs[mid] <= target { // the only difference from LowerBound
			lo = mid + 1
		} else {
			hi = mid
		}
	}
	return lo
}

// FirstTrue returns the smallest x in [lo, hi) for which pred(x) is true,
// or hi if pred is false everywhere. pred must be monotone: once true, it
// stays true (false false ... false true true ... true).
//
// This is the "binary search on the answer" template. It is sort.Search with
// an explicit range.
func FirstTrue(lo, hi int, pred func(int) bool) int {
	for lo < hi {
		mid := lo + (hi-lo)/2
		if pred(mid) {
			hi = mid
		} else {
			lo = mid + 1
		}
	}
	return lo
}
