package bits

// Problem: Subsets via Bitmask
// Difficulty: Medium
//
// Given distinct integers, list every subset (the power set, 2^n entries).
// Elements inside a subset keep their original order.
//
// Examples:
//
//	[1,2,3] -> [] [1] [2] [1 2] [3] [1 3] [2 3] [1 2 3]
//	[]      -> []
//
// Reference: https://leetcode.com/problems/subsets/
//
// Pattern: each integer mask in [0, 2^n) IS a subset — bit i set means
// "include nums[i]". Counting 0..2^n-1 enumerates them all, no recursion.
//
//	mask  bits  subset
//	 5    101   {nums[0], nums[2]}
//
// Follow-up (Gosper's hack): step through only the masks with exactly k
// bits set, in increasing order, to list k-element combinations.
//
// Time: O(n·2^n). Space: O(n·2^n) for the output. n must stay small (≤ ~20).

// Subsets enumerates masks 0 .. 2^n-1.
func Subsets(nums []int) [][]int {
	n := len(nums)
	out := make([][]int, 0, 1<<n)
	for mask := 0; mask < 1<<n; mask++ {
		subset := []int{}
		for i := 0; i < n; i++ {
			if mask>>i&1 == 1 {
				subset = append(subset, nums[i])
			}
		}
		out = append(out, subset)
	}
	return out
}

// Combinations returns every k-subset using Gosper's hack to jump directly
// from one k-bit mask to the next larger one.
func Combinations(nums []int, k int) [][]int {
	n := len(nums)
	if k < 0 || k > n {
		return nil
	}
	if k == 0 {
		return [][]int{{}}
	}
	var out [][]int
	for mask := 1<<k - 1; mask < 1<<n; {
		subset := make([]int, 0, k)
		for i := 0; i < n; i++ {
			if mask>>i&1 == 1 {
				subset = append(subset, nums[i])
			}
		}
		out = append(out, subset)
		// Gosper: move the lowest block of 1s' top bit up by one and pack
		// the remaining bits of that block at the bottom.
		low := mask & -mask      // lowest set bit
		ripple := mask + low     // carry clears the lowest block, sets the next 0
		ones := mask ^ ripple    // the block that changed (plus one extra bit)
		ones = (ones >> 2) / low // shift it down to bit 0, dropping two bits
		mask = ripple | ones
	}
	return out
}
