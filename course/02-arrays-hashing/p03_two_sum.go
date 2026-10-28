package arrayshashing

// Problem 03 — Two Sum (Easy)
//
// Given nums and a target, return the indices [i, j] (i < j) of two different
// elements that add up to target, or nil if no such pair exists.
//
//	nums = [2, 7, 11, 15], target = 9 → [0, 1]
//	nums = [3, 3],         target = 6 → [0, 1]
//
// Reference: https://leetcode.com/problems/two-sum/
//
// Pattern: the brute force's inner loop is a SEARCH for the complement
// (target - x). Replace the search with a map lookup from value → index of
// everything seen so far. Checking before inserting stops an element from
// pairing with itself.
//
//	TwoSumBrute: Time O(n²), Space O(1)
//	TwoSum:      Time O(n),  Space O(n)

// TwoSumBrute tries every pair.
func TwoSumBrute(nums []int, target int) []int {
	for i := range nums {
		for j := i + 1; j < len(nums); j++ {
			if nums[i]+nums[j] == target {
				return []int{i, j}
			}
		}
	}
	return nil
}

// TwoSum does one pass with a value → index map.
func TwoSum(nums []int, target int) []int {
	indexOf := make(map[int]int, len(nums))
	for j, x := range nums {
		if i, ok := indexOf[target-x]; ok {
			return []int{i, j}
		}
		indexOf[x] = j // insert AFTER the lookup so x can't pair with itself
	}
	return nil
}
