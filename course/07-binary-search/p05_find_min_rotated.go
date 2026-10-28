package binarysearch

// Problem: Find Minimum in Rotated Sorted Array
// Difficulty: Medium
//
// An ascending array of distinct integers was rotated at some unknown pivot,
// so [0,1,2,4,5,6,7] might now be [4,5,6,7,0,1,2]. Return the minimum element
// in O(log n). The array is non-empty.
//
// Examples:
//
//	[3,4,5,1,2]   -> 1
//	[11,13,15,17] -> 11 (rotated by zero, or by n)
//
// Reference: https://leetcode.com/problems/find-minimum-in-rotated-sorted-array/
//
// Pattern: binary search on a predicate, not a value. "nums[i] <= last
// element" is false on the left (larger) segment and true on the right
// (smaller) segment, which is monotone. The minimum is the first index where it is true.
// Comparing mid with nums[hi] (not nums[lo]) avoids the unrotated edge case.
//
// Time: brute O(n); optimal O(log n). Space: O(1).

// FindMinBrute scans linearly.
func FindMinBrute(nums []int) int {
	m := nums[0]
	for _, v := range nums[1:] {
		m = min(m, v)
	}
	return m
}

// FindMin binary-searches for the rotation point.
func FindMin(nums []int) int {
	lo, hi := 0, len(nums)-1
	for lo < hi {
		mid := lo + (hi-lo)/2
		if nums[mid] > nums[hi] {
			lo = mid + 1 // mid is in the left (larger) segment; min is strictly right
		} else {
			hi = mid // mid is in the right segment; it may be the min itself
		}
	}
	return nums[lo]
}
