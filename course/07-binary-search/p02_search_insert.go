package binarysearch

// Problem: Search Insert Position
// Difficulty: Easy
//
// Given an ascending array of distinct integers and a target, return the
// index of the target if present; otherwise return the index where it would
// have to be inserted to keep the array sorted. Aim for O(log n).
//
// Examples:
//
//	nums=[1,3,5,6], target=5 -> 2
//	nums=[1,3,5,6], target=2 -> 1
//	nums=[1,3,5,6], target=7 -> 4
//
// Reference: https://leetcode.com/problems/search-insert-position/
//
// Pattern: this is exactly the lower bound, the first index whose value is
// >= target. Both "found" and "insert here" are answered by the same index.
//
// Time: O(log n). Space: O(1).

// SearchInsert returns the index of target or its insertion point.
func SearchInsert(nums []int, target int) int {
	lo, hi := 0, len(nums) // hi = len(nums) is a valid answer (append at end)
	for lo < hi {
		mid := lo + (hi-lo)/2
		if nums[mid] < target {
			lo = mid + 1
		} else {
			hi = mid
		}
	}
	return lo
}
