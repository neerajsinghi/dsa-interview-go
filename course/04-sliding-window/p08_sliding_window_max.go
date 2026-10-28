package slidingwindow

// Problem 08 — Sliding Window Maximum (Hard)
//
// For every window of k consecutive elements of nums (left to right), report
// the window's maximum. Return nil if k is not in 1..len(nums).
//
//	nums = [1, 3, -1, -3, 5, 3, 6, 7], k = 3 → [3, 3, 5, 5, 6, 7]
//	nums = [9, 8, 7],                  k = 1 → [9, 8, 7]
//
// Reference: https://leetcode.com/problems/sliding-window-maximum/
//
// Pattern: MONOTONIC DEQUE of indices whose values are strictly decreasing
// from front to back. The front is always the current maximum.
//   - Before pushing i, pop smaller-or-equal values from the back: they are
//     older AND not bigger than nums[i], so they can never be a maximum again.
//   - Pop the front when it slides out of the window (index ≤ i-k).
//
//	i  x   deque (values)   window max
//	0  1   [1]
//	1  3   [3]              (1 popped)
//	2  -1  [3, -1]          3
//	3  -3  [3, -1, -3]      3
//	4  5   [5]              5  (3 expired; -1, -3 popped)
//	5  3   [5, 3]           5
//	6  6   [6]              6
//	7  7   [7]              7
//
//	MaxSlidingWindowBrute: Time O(n·k), Space O(1) extra
//	MaxSlidingWindow:      Time O(n) — each index pushed and popped once, Space O(k)

// MaxSlidingWindowBrute scans each window.
func MaxSlidingWindowBrute(nums []int, k int) []int {
	if k <= 0 || k > len(nums) {
		return nil
	}
	out := make([]int, 0, len(nums)-k+1)
	for start := 0; start+k <= len(nums); start++ {
		m := nums[start]
		for _, x := range nums[start+1 : start+k] {
			m = max(m, x)
		}
		out = append(out, m)
	}
	return out
}

// MaxSlidingWindow uses a monotonic deque of indices.
func MaxSlidingWindow(nums []int, k int) []int {
	if k <= 0 || k > len(nums) {
		return nil
	}
	out := make([]int, 0, len(nums)-k+1)
	var dq Deque // indices; nums[dq...] strictly decreasing
	for i, x := range nums {
		if dq.Len() > 0 && dq.Front() <= i-k {
			dq.PopFront() // front index fell out of the window
		}
		for dq.Len() > 0 && nums[dq.Back()] <= x {
			dq.PopBack() // dominated: older and not larger
		}
		dq.PushBack(i)
		if i >= k-1 {
			out = append(out, nums[dq.Front()])
		}
	}
	return out
}
