package binarysearch

// Problem: Find First and Last Position of Element in Sorted Array
// Difficulty: Medium
//
// Given an ascending array that may contain duplicates, return the first and
// last index at which target appears, or [-1, -1] if it doesn't appear.
// Aim for O(log n).
//
// Examples:
//
//	nums=[5,7,7,8,8,10], target=8 -> [3,4]
//	nums=[5,7,7,8,8,10], target=6 -> [-1,-1]
//
// Reference: https://leetcode.com/problems/find-first-and-last-position-of-element-in-sorted-array/
//
// Pattern: two boundary searches. first = lower bound (first >= target);
// last = upper bound (first > target) minus one. Two O(log n) searches beat
// "find any copy, then walk outward", which degrades to O(n) on many duplicates.
//
// Time: brute O(n); optimal O(log n). Space: O(1).

// SearchRangeBrute scans linearly.
func SearchRangeBrute(nums []int, target int) [2]int {
	res := [2]int{-1, -1}
	for i, v := range nums {
		if v == target {
			if res[0] == -1 {
				res[0] = i
			}
			res[1] = i
		}
	}
	return res
}

// SearchRange uses lower and upper bound.
func SearchRange(nums []int, target int) [2]int {
	first := LowerBound(nums, target)
	if first == len(nums) || nums[first] != target {
		return [2]int{-1, -1}
	}
	last := UpperBound(nums, target) - 1
	return [2]int{first, last}
}
