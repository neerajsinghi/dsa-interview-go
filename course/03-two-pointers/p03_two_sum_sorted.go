package twopointers

// Problem 03 — Two Sum II: Input Array Is Sorted (Medium)
//
// nums is sorted ascending. Return the 0-based indices [i, j] with i < j of two
// elements that sum to target, or nil if none exist. Use O(1) extra space.
// (LeetCode's version is 1-based and guarantees exactly one answer.)
//
//	nums = [2, 7, 11, 15], target = 9 → [0, 1]
//	nums = [-1, 0],        target = -1 → [0, 1]
//
// Reference: https://leetcode.com/problems/two-sum-ii-input-array-is-sorted/
//
// Pattern: opposite-ends pointers on sorted data. If the sum is too small,
// the only way to grow it is to move the LEFT pointer right; if too big, move
// the RIGHT pointer left. Each move discards a value that cannot be in any
// solution, so we never miss the answer.
//
//	[2, 7, 11, 15], target 9
//	 l          r   2+15=17 > 9 → r--
//	 l      r       2+11=13 > 9 → r--
//	 l  r           2+7 = 9     → done
//
//	TwoSumSortedBrute: Time O(n²), Space O(1)
//	TwoSumSorted:      Time O(n),  Space O(1)

// TwoSumSortedBrute ignores sortedness and tries every pair.
func TwoSumSortedBrute(nums []int, target int) []int {
	for i := range nums {
		for j := i + 1; j < len(nums); j++ {
			if nums[i]+nums[j] == target {
				return []int{i, j}
			}
		}
	}
	return nil
}

// TwoSumSorted squeezes from both ends.
func TwoSumSorted(nums []int, target int) []int {
	l, r := 0, len(nums)-1
	for l < r {
		switch sum := nums[l] + nums[r]; {
		case sum == target:
			return []int{l, r}
		case sum < target:
			l++ // need a bigger sum
		default:
			r-- // need a smaller sum
		}
	}
	return nil
}
