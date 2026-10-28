package twopointers

// Problem 07 — Trapping Rain Water (Hard)
//
// heights describes an elevation map of bars, each 1 unit wide. Return how
// many units of rain water are trapped between the bars after it rains.
//
//	[0, 1, 0, 2, 1, 0, 1, 3, 2, 1, 2, 1] → 6
//	[4, 2, 0, 3, 2, 5]                   → 9
//
//	              █
//	      █ ≈ ≈ ≈ █ █ ≈ █
//	  █ ≈ █ █ ≈ █ █ █ █ █ █      ≈ = trapped water (6 units)
//	0 1 0 2 1 0 1 3 2 1 2 1
//
// Reference: https://leetcode.com/problems/trapping-rain-water/
//
// Pattern: water above bar i = min(tallest bar to the left, tallest to the
// right) − heights[i]. The brute force recomputes both maxima for every i.
// Precomputing prefix/suffix maxima makes it O(n). Two pointers drop the
// arrays: whichever side has the smaller running max is the bottleneck for
// its current bar — the other side is known to have something at least as
// tall — so that bar can be settled and its pointer moved.
//
//	TrapBrute:     Time O(n²), Space O(1)
//	TrapPrefixMax: Time O(n),  Space O(n)
//	Trap:          Time O(n),  Space O(1)

// TrapBrute scans left and right from every bar.
func TrapBrute(heights []int) int {
	total := 0
	for i := range heights {
		leftMax, rightMax := 0, 0
		for j := 0; j <= i; j++ {
			leftMax = max(leftMax, heights[j])
		}
		for j := i; j < len(heights); j++ {
			rightMax = max(rightMax, heights[j])
		}
		total += min(leftMax, rightMax) - heights[i]
	}
	return total
}

// TrapPrefixMax precomputes the tallest bar on each side of every index.
func TrapPrefixMax(heights []int) int {
	n := len(heights)
	if n == 0 {
		return 0
	}
	leftMax := make([]int, n)  // max of heights[0..i]
	rightMax := make([]int, n) // max of heights[i..n-1]
	leftMax[0] = heights[0]
	for i := 1; i < n; i++ {
		leftMax[i] = max(leftMax[i-1], heights[i])
	}
	rightMax[n-1] = heights[n-1]
	for i := n - 2; i >= 0; i-- {
		rightMax[i] = max(rightMax[i+1], heights[i])
	}
	total := 0
	for i := range heights {
		total += min(leftMax[i], rightMax[i]) - heights[i]
	}
	return total
}

// Trap uses two pointers and two running maxima.
func Trap(heights []int) int {
	l, r := 0, len(heights)-1
	leftMax, rightMax := 0, 0
	total := 0
	for l <= r {
		leftMax = max(leftMax, heights[l])
		rightMax = max(rightMax, heights[r])
		if leftMax <= rightMax {
			// Right side has a wall ≥ leftMax, so bar l's water level is leftMax.
			total += leftMax - heights[l]
			l++
		} else {
			total += rightMax - heights[r]
			r--
		}
	}
	return total
}
