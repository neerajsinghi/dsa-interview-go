package slidingwindow

// Problem 03 — Minimum Size Subarray Sum (Medium)
//
// nums holds POSITIVE integers. Return the length of the shortest contiguous
// subarray whose sum is at least target, or 0 if there is none.
//
//	target = 7, nums = [2, 3, 1, 2, 4, 3] → 2   ([4, 3])
//	target = 11, nums = [1, 1, 1, 1]      → 0
//
// Reference: https://leetcode.com/problems/minimum-size-subarray-sum/
//
// Pattern: VARIABLE-SIZE window. Grow the right edge until the window is
// valid (sum ≥ target), then shrink from the left while it stays valid,
// recording the length each time. This works because values are positive:
// growing only increases the sum, shrinking only decreases it. (With negative
// numbers, use prefix sums instead — see module 05.)
//
//	MinSubArrayLenBrute: Time O(n²), Space O(1)
//	MinSubArrayLen:      Time O(n),  Space O(1) — each index enters and leaves once

// MinSubArrayLenBrute extends a sum from every start.
func MinSubArrayLenBrute(target int, nums []int) int {
	best := 0
	for i := range nums {
		sum := 0
		for j := i; j < len(nums); j++ {
			sum += nums[j]
			if sum >= target {
				if best == 0 || j-i+1 < best {
					best = j - i + 1
				}
				break // longer windows from i can't be shorter
			}
		}
	}
	return best
}

// MinSubArrayLen grows right, shrinks left.
func MinSubArrayLen(target int, nums []int) int {
	best := 0
	sum, left := 0, 0
	for right, x := range nums {
		sum += x
		for sum >= target { // window valid: try to shrink it
			if best == 0 || right-left+1 < best {
				best = right - left + 1
			}
			sum -= nums[left]
			left++
		}
	}
	return best
}
