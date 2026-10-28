package complexity

import "slices"

// Problem 02 — Contains Duplicate: O(n²) vs O(n log n) vs O(n) (Easy)
//
// Report whether any value appears at least twice in nums.
//
//	[1, 2, 3, 1] → true
//	[1, 2, 3, 4] → false
//
// Reference: https://leetcode.com/problems/contains-duplicate/
//
// Pattern: the same question answered three ways shows the classic
// time/space trade-off.
//
//	ContainsDuplicateBrute: Time O(n²),       Space O(1)  — compare every pair
//	ContainsDuplicateSort:  Time O(n log n),  Space O(n)  — duplicates become neighbours
//	ContainsDuplicate:      Time O(n),        Space O(n)  — remember what you've seen

// ContainsDuplicateBrute compares each pair once: n(n-1)/2 comparisons.
func ContainsDuplicateBrute(nums []int) bool {
	for i := 0; i < len(nums); i++ {
		for j := i + 1; j < len(nums); j++ {
			if nums[i] == nums[j] {
				return true
			}
		}
	}
	return false
}

// ContainsDuplicateSort sorts a copy so equal values sit next to each other.
// Cloning keeps the caller's slice untouched (the copy is the O(n) space;
// sorting in place would be O(log n) but mutates the input — ask first!).
func ContainsDuplicateSort(nums []int) bool {
	s := slices.Clone(nums)
	slices.Sort(s)
	for i := 1; i < len(s); i++ {
		if s[i] == s[i-1] {
			return true
		}
	}
	return false
}

// ContainsDuplicate trades memory for speed with a hash set.
func ContainsDuplicate(nums []int) bool {
	seen := make(map[int]struct{}, len(nums))
	for _, x := range nums {
		if _, ok := seen[x]; ok {
			return true
		}
		seen[x] = struct{}{}
	}
	return false
}
