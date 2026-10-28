package stacks

// Problem: Largest Rectangle in Histogram
// Difficulty: Hard
//
// A histogram has bars of width 1 with the given non-negative heights. Find
// the area of the largest axis-aligned rectangle that fits entirely inside
// the bars.
//
// Examples:
//
//	[2,1,5,6,2,3] -> 10   (bars 5 and 6, height 5, width 2)
//	[2,4]         -> 4
//
// Reference: https://leetcode.com/problems/largest-rectangle-in-histogram/
//
// Pattern: monotonic increasing stack. For each bar, the best rectangle that
// uses its full height stretches left to the previous shorter bar and right to
// the next shorter bar. Keep indices of increasing heights; when a shorter bar
// arrives, it is the "next shorter" for everything taller on the stack, and the
// element beneath each popped bar is its "previous shorter". A sentinel height
// 0 at the end flushes the stack.
//
// Time: brute O(n^2); optimal O(n). Space: O(n).

// LargestRectangleBrute tries every left edge and extends right, tracking the
// minimum height seen.
func LargestRectangleBrute(heights []int) int {
	best := 0
	for i := range heights {
		minH := heights[i]
		for j := i; j < len(heights); j++ {
			minH = min(minH, heights[j])
			best = max(best, minH*(j-i+1))
		}
	}
	return best
}

// LargestRectangle uses a monotonic stack of bar indices.
func LargestRectangle(heights []int) int {
	best := 0
	var st Stack[int] // indices with strictly increasing heights
	for i := 0; i <= len(heights); i++ {
		h := 0 // sentinel bar after the end
		if i < len(heights) {
			h = heights[i]
		}
		for {
			top, ok := st.Peek()
			if !ok || heights[top] < h {
				break
			}
			st.Pop()
			// Bar `top` extends right up to i-1. It extends left to just after
			// the new stack top (the previous shorter bar), or to 0 if none.
			left := -1
			if prev, ok := st.Peek(); ok {
				left = prev
			}
			best = max(best, heights[top]*(i-left-1))
		}
		st.Push(i)
	}
	return best
}
