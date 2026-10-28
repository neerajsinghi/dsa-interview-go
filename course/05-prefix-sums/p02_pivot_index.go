package prefixsums

// Problem 02 — Find Pivot Index (Easy)
//
// The pivot index is the position where the sum of everything strictly to its
// left equals the sum of everything strictly to its right (an empty side sums
// to 0). Return the leftmost pivot index, or -1 if there is none.
//
//	[1, 7, 3, 6, 5, 6] → 3   (1+7+3 = 11 = 5+6)
//	[2, 1, -1]         → 0   (left is empty = 0 = 1 + -1)
//
// Reference: https://leetcode.com/problems/find-pivot-index/
//
// Pattern: with the total known, the right sum is derivable from the left
// sum: right = total − left − nums[i]. One running prefix is enough.
//
//	PivotIndexBrute: Time O(n²), Space O(1)
//	PivotIndex:      Time O(n),  Space O(1)

// PivotIndexBrute sums both sides for every index.
func PivotIndexBrute(nums []int) int {
	for i := range nums {
		left, right := 0, 0
		for _, x := range nums[:i] {
			left += x
		}
		for _, x := range nums[i+1:] {
			right += x
		}
		if left == right {
			return i
		}
	}
	return -1
}

// PivotIndex keeps a running left sum.
func PivotIndex(nums []int) int {
	total := 0
	for _, x := range nums {
		total += x
	}
	left := 0
	for i, x := range nums {
		if left == total-left-x {
			return i
		}
		left += x
	}
	return -1
}
