package binarysearch

// Problem: Binary Search
// Difficulty: Easy
//
// Given an ascending array of distinct integers and a target, return the
// index of the target, or -1 if it is not present. Aim for O(log n).
//
// Examples:
//
//	nums=[-1,0,3,5,9,12], target=9 -> 4
//	nums=[-1,0,3,5,9,12], target=2 -> -1
//
// Reference: https://leetcode.com/problems/binary-search/
//
// Pattern: classic "find exact" binary search on a closed interval
// [lo, hi]. Compare with the middle and discard the half that cannot hold
// the target.
//
// Time: O(log n). Space: O(1).

// Search returns the index of target in sorted nums, or -1.
func Search(nums []int, target int) int {
	lo, hi := 0, len(nums)-1 // closed interval: both ends are candidates
	for lo <= hi {
		mid := lo + (hi-lo)/2
		switch {
		case nums[mid] == target:
			return mid
		case nums[mid] < target:
			lo = mid + 1
		default:
			hi = mid - 1
		}
	}
	return -1
}

// SearchLowerBound is the same search written with the half-open template:
// find the lower bound, then check whether it holds the target.
func SearchLowerBound(nums []int, target int) int {
	i := LowerBound(nums, target)
	if i < len(nums) && nums[i] == target {
		return i
	}
	return -1
}
