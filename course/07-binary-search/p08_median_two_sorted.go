package binarysearch

import "math"

// Problem: Median of Two Sorted Arrays
// Difficulty: Hard
//
// Given two ascending arrays a and b, return the median of all their values
// combined, in O(log(min(m, n))). For an even total, the median is the mean
// of the two middle values. Returns 0 when both are empty.
//
// Examples:
//
//	a=[1,3],   b=[2]    -> 2.0
//	a=[1,2],   b=[3,4]  -> 2.5
//
// Reference: https://leetcode.com/problems/median-of-two-sorted-arrays/
//
// Pattern: binary search on a partition. Take i elements from a and
// j = half - i from b to form the "left half" of the merged order. The cut is
// correct when every left value <= every right value, i.e.
// a[i-1] <= b[j] and b[j-1] <= a[i]. If a[i-1] is too big, move i left; if
// b[j-1] is too big, move i right. Search over the shorter array.
//
//	a:  a[0] .. a[i-1] | a[i] .. a[m-1]
//	b:  b[0] .. b[j-1] | b[j] .. b[n-1]
//	    \___ left ___/   \__ right __/      len(left) = (m+n+1)/2
//
// Time: brute O(m+n); optimal O(log min(m,n)). Space: O(1).

// FindMedianSortedArraysBrute merges the two arrays and picks the middle.
func FindMedianSortedArraysBrute(a, b []int) float64 {
	merged := make([]int, 0, len(a)+len(b))
	i, j := 0, 0
	for i < len(a) || j < len(b) {
		if j == len(b) || (i < len(a) && a[i] <= b[j]) {
			merged = append(merged, a[i])
			i++
		} else {
			merged = append(merged, b[j])
			j++
		}
	}
	n := len(merged)
	if n == 0 {
		return 0
	}
	if n%2 == 1 {
		return float64(merged[n/2])
	}
	return float64(merged[n/2-1]+merged[n/2]) / 2
}

// FindMedianSortedArrays binary-searches the partition point.
func FindMedianSortedArrays(a, b []int) float64 {
	if len(a) > len(b) {
		a, b = b, a // search the shorter array so j stays in range
	}
	m, n := len(a), len(b)
	if m+n == 0 {
		return 0
	}
	half := (m + n + 1) / 2 // left side gets the extra element when odd

	// at returns xs[k], or ±infinity past the edges, so empty sides never "win".
	at := func(xs []int, k int) int {
		switch {
		case k < 0:
			return math.MinInt
		case k >= len(xs):
			return math.MaxInt
		}
		return xs[k]
	}

	lo, hi := 0, m // i = number of elements taken from a
	for lo <= hi {
		i := lo + (hi-lo)/2
		j := half - i
		aLeft, aRight := at(a, i-1), at(a, i)
		bLeft, bRight := at(b, j-1), at(b, j)
		switch {
		case aLeft > bRight:
			hi = i - 1 // took too many from a
		case bLeft > aRight:
			lo = i + 1 // took too few from a
		default:
			leftMax := max(aLeft, bLeft)
			if (m+n)%2 == 1 {
				return float64(leftMax)
			}
			rightMin := min(aRight, bRight)
			return (float64(leftMax) + float64(rightMin)) / 2
		}
	}
	panic("unreachable: inputs must be sorted")
}
