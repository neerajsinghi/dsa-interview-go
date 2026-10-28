package complexity

// Problem 05 — Space Complexity of Recursion (Easy)
//
// Sum the integers in nums, once recursively and once with a loop, and
// observe that they differ in SPACE, not time.
//
//	[1, 2, 3, 4] → 10
//	[]           → 0
//
// Pattern: every pending recursive call keeps a stack frame alive. A
// recursion that goes n levels deep uses O(n) space even if it allocates
// nothing on the heap. Go's goroutine stacks grow dynamically, so you won't
// crash at depth 10⁵ the way C or Python might — but the memory is still
// used, and an interviewer will expect you to count it.
//
//	SumRecursive: Time O(n), Space O(n) call stack
//	SumIterative: Time O(n), Space O(1)
//	SumHalves:    Time O(n), Space O(log n) call stack (depth halves each level)

// SumRecursive: sum(nums) = nums[0] + sum(nums[1:]). Depth n.
// Re-slicing nums[1:] is O(1) in Go — it shares the backing array.
func SumRecursive(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	return nums[0] + SumRecursive(nums[1:])
}

// SumIterative uses one accumulator: constant extra space.
func SumIterative(nums []int) int {
	total := 0
	for _, x := range nums {
		total += x
	}
	return total
}

// SumHalves splits in two like merge sort. The tree has n leaves (so O(n)
// work) but is only log₂ n deep (so O(log n) stack).
func SumHalves(nums []int) int {
	switch len(nums) {
	case 0:
		return 0
	case 1:
		return nums[0]
	}
	mid := len(nums) / 2
	return SumHalves(nums[:mid]) + SumHalves(nums[mid:])
}
