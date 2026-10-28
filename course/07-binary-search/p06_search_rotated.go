package binarysearch

// Problem: Search in Rotated Sorted Array
// Difficulty: Medium
//
// An ascending array of distinct integers has been rotated at an unknown
// pivot. Return the index of target, or -1, in O(log n).
//
// Examples:
//
//	nums=[4,5,6,7,0,1,2], target=0 -> 4
//	nums=[4,5,6,7,0,1,2], target=3 -> -1
//
// Reference: https://leetcode.com/problems/search-in-rotated-sorted-array/
//
// Pattern: at least one half is always sorted. Cutting a rotated array at
// mid leaves one half that is fully sorted. Check whether target lies inside
// that sorted half's [first, last] range; if so search it, otherwise search the
// other half.
//
// Time: brute O(n); optimal O(log n). Space: O(1).

// SearchRotatedBrute scans linearly.
func SearchRotatedBrute(nums []int, target int) int {
	for i, v := range nums {
		if v == target {
			return i
		}
	}
	return -1
}

// SearchRotated finds target in a rotated sorted array of distinct values.
func SearchRotated(nums []int, target int) int {
	lo, hi := 0, len(nums)-1
	for lo <= hi {
		mid := lo + (hi-lo)/2
		if nums[mid] == target {
			return mid
		}
		if nums[lo] <= nums[mid] { // left half [lo, mid] is sorted
			if nums[lo] <= target && target < nums[mid] {
				hi = mid - 1
			} else {
				lo = mid + 1
			}
		} else { // right half [mid, hi] is sorted
			if nums[mid] < target && target <= nums[hi] {
				lo = mid + 1
			} else {
				hi = mid - 1
			}
		}
	}
	return -1
}
