package complexity

// Problem 01 — Linear vs Binary Search (Easy)
//
// Given a slice of ints sorted in ascending order and a target, return the
// index of target, or -1 if it is not present.
//
//	nums = [1, 3, 5, 7, 9], target = 7  →  3
//	nums = [1, 3, 5, 7, 9], target = 4  →  -1
//
// Reference: https://leetcode.com/problems/binary-search/
//
// Pattern: when the input is sorted, each comparison can discard half of the
// remaining candidates. Halving n until it reaches 1 takes log₂ n steps, so a
// million elements need ~20 comparisons instead of a million.
//
//	LinearSearch: Time O(n),     Space O(1)
//	BinarySearch: Time O(log n), Space O(1)

// LinearSearch scans left to right. It does not need sorted input — that is
// its only advantage.
func LinearSearch(nums []int, target int) int {
	for i, x := range nums {
		if x == target {
			return i
		}
	}
	return -1
}

// BinarySearch keeps the invariant "if target exists, it is in nums[lo..hi]".
func BinarySearch(nums []int, target int) int {
	lo, hi := 0, len(nums)-1
	for lo <= hi {
		mid := lo + (hi-lo)/2 // avoids (lo+hi) overflow
		switch {
		case nums[mid] == target:
			return mid
		case nums[mid] < target:
			lo = mid + 1 // target is strictly right of mid
		default:
			hi = mid - 1 // target is strictly left of mid
		}
	}
	return -1
}
