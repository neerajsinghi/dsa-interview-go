package twopointers

// Problem 06 — Container With Most Water (Medium)
//
// heights[i] is the height of a vertical line at x = i. Choose two lines that,
// together with the x-axis, hold the most water. Return that area:
// (j - i) × min(heights[i], heights[j]).
//
//	[1, 8, 6, 2, 5, 4, 8, 3, 7] → 49   (lines at 1 and 8: 7 × min(8, 7))
//	[1, 1]                       → 1
//
// Reference: https://leetcode.com/problems/container-with-most-water/
//
// Pattern: start with the widest container (both ends) and move the SHORTER
// line inwards. Why is that safe? The area is capped by the shorter line, and
// any container that keeps it while getting narrower can only be smaller —
// so the shorter line can be discarded for good.
//
//	ContainerBrute: Time O(n²), Space O(1)
//	MaxArea:        Time O(n),  Space O(1)

// ContainerBrute tries every pair of lines.
func ContainerBrute(heights []int) int {
	best := 0
	for i := range heights {
		for j := i + 1; j < len(heights); j++ {
			best = max(best, (j-i)*min(heights[i], heights[j]))
		}
	}
	return best
}

// MaxArea squeezes from both ends, always dropping the shorter line.
func MaxArea(heights []int) int {
	best := 0
	l, r := 0, len(heights)-1
	for l < r {
		best = max(best, (r-l)*min(heights[l], heights[r]))
		if heights[l] < heights[r] {
			l++
		} else {
			r--
		}
	}
	return best
}
