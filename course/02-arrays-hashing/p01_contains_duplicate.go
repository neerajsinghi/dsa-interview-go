package arrayshashing

// Problem 01 — Contains Duplicate (Easy)
//
// Given a slice of integers, report whether any value occurs more than once.
//
//	[1, 2, 3, 1] → true
//	[1, 2, 3, 4] → false
//
// Reference: https://leetcode.com/problems/contains-duplicate/
//
// Pattern: "have I seen this before?" is a hash-set membership test. Insert
// as you go and stop at the first repeat.
//
//	ContainsDuplicateBrute: Time O(n²), Space O(1)
//	ContainsDuplicate:      Time O(n),  Space O(n)

// ContainsDuplicateBrute compares every pair.
func ContainsDuplicateBrute(nums []int) bool {
	for i := range nums {
		for j := i + 1; j < len(nums); j++ {
			if nums[i] == nums[j] {
				return true
			}
		}
	}
	return false
}

// ContainsDuplicate keeps a set of values seen so far.
func ContainsDuplicate(nums []int) bool {
	seen := make(Set[int], len(nums))
	for _, x := range nums {
		if seen.Has(x) {
			return true
		}
		seen.Add(x)
	}
	return false
}
