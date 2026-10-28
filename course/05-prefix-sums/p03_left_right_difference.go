package prefixsums

// Problem 03 — Left and Right Sum Differences (Easy)
//
// For each index i, compute |leftSum(i) − rightSum(i)|, where leftSum is the
// sum of elements strictly before i and rightSum the sum strictly after it.
//
//	[10, 4, 8, 3] → [15, 1, 11, 22]
//	  left:  0, 10, 14, 22      right: 15, 11, 3, 0
//	[1]           → [0]
//
// Reference: https://leetcode.com/problems/left-and-right-sum-differences/
//
// Pattern: PREFIX and SUFFIX sums. A prefix pass gives left sums; a suffix
// sum is total − prefix − self, so no second array is required. The same
// prefix/suffix idea powers "product of array except self" (module 02).
//
//	LeftRightDifferenceBrute: Time O(n²), Space O(1) extra
//	LeftRightDifference:      Time O(n),  Space O(1) extra (output not counted)

// LeftRightDifferenceBrute sums both sides for every index.
func LeftRightDifferenceBrute(nums []int) []int {
	out := make([]int, len(nums))
	for i := range nums {
		left, right := 0, 0
		for j, x := range nums {
			switch {
			case j < i:
				left += x
			case j > i:
				right += x
			}
		}
		out[i] = abs(left - right)
	}
	return out
}

// LeftRightDifference uses a running prefix and the total.
func LeftRightDifference(nums []int) []int {
	total := 0
	for _, x := range nums {
		total += x
	}
	out := make([]int, len(nums))
	left := 0
	for i, x := range nums {
		right := total - left - x
		out[i] = abs(left - right)
		left += x
	}
	return out
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
