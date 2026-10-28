package twopointers

// Problem 02 — Remove Duplicates from Sorted Array (Easy)
//
// nums is sorted in non-decreasing order. Rearrange it IN PLACE so that its
// first k elements are the distinct values in their original order, and
// return k. What is left beyond index k-1 does not matter.
//
//	[1, 1, 2]                      → k = 2, nums[:2] = [1, 2]
//	[0, 0, 1, 1, 1, 2, 2, 3, 3, 4] → k = 5, nums[:5] = [0, 1, 2, 3, 4]
//
// Reference: https://leetcode.com/problems/remove-duplicates-from-sorted-array/
//
// Pattern: same-direction READ/WRITE pointers. `write` marks the end of the
// deduplicated prefix; `read` scans ahead. Because the input is sorted, a
// value is new exactly when it differs from the last value written.
//
//	nums:  0 0 1 1 1 2
//	       w r            nums[r]==nums[w-1]? skip
//	         w   r        new value 1 → nums[w]=1, w++
//
//	RemoveDuplicatesBrute: Time O(n), Space O(n) (temporary slice)
//	RemoveDuplicates:      Time O(n), Space O(1)

// RemoveDuplicatesBrute collects distinct values in a new slice and copies
// them back.
func RemoveDuplicatesBrute(nums []int) int {
	var distinct []int
	for i, x := range nums {
		if i == 0 || x != nums[i-1] {
			distinct = append(distinct, x)
		}
	}
	return copy(nums, distinct)
}

// RemoveDuplicates compacts in place.
func RemoveDuplicates(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	write := 1 // nums[:write] is deduplicated; nums[0] is always kept
	for read := 1; read < len(nums); read++ {
		if nums[read] != nums[write-1] {
			nums[write] = nums[read]
			write++
		}
	}
	return write
}
