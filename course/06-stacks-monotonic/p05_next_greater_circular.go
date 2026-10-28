package stacks

// Problem: Next Greater Element (circular)
// Difficulty: Medium
//
// The array is circular: after the last element you wrap around to the first.
// For every element, find the first element strictly greater than it when
// walking forward (wrapping at most once). Use -1 if there is none.
//
// Examples:
//
//	[1,2,1]     -> [2,-1,2]   (the last 1 wraps around to find 2)
//	[5,4,3,2,1] -> [-1,5,5,5,5]
//
// Reference: https://leetcode.com/problems/next-greater-element-ii/
// (The non-circular "Next Greater Element I" is the same loop run once.)
//
// Pattern: monotonic stack + walk the array twice. Indices 0..2n-1 visit every
// element twice (i % n); the second pass resolves elements whose answer lies
// before them. Only the first pass pushes, so the stack stays O(n).
//
// Time: brute O(n^2); optimal O(n). Space: O(n).

// NextGreaterCircularBrute walks up to n-1 steps forward from each element.
func NextGreaterCircularBrute(nums []int) []int {
	n := len(nums)
	ans := make([]int, n)
	for i := range nums {
		ans[i] = -1
		for step := 1; step < n; step++ {
			if v := nums[(i+step)%n]; v > nums[i] {
				ans[i] = v
				break
			}
		}
	}
	return ans
}

// NextGreaterCircular uses a monotonic decreasing stack over two passes.
func NextGreaterCircular(nums []int) []int {
	n := len(nums)
	ans := make([]int, n)
	for i := range ans {
		ans[i] = -1
	}
	var st Stack[int] // indices with no greater element found yet
	for i := 0; i < 2*n; i++ {
		v := nums[i%n]
		for {
			j, ok := st.Peek()
			if !ok || nums[j] >= v {
				break
			}
			st.Pop()
			ans[j] = v
		}
		if i < n { // second pass only resolves, never pushes
			st.Push(i)
		}
	}
	return ans
}
