// Package prefixsums covers problems that precompute running totals so that
// any range sum becomes a subtraction: sum(nums[l..r]) = P[r+1] − P[l].
package prefixsums

// BuildPrefix returns P with len(nums)+1 entries where P[i] is the sum of
// nums[0..i-1]. The leading 0 (the sum of the empty prefix) removes the
// special case for ranges that start at index 0.
//
//	nums = [3, 1, 4, 1]
//	P    = [0, 3, 4, 8, 9]
func BuildPrefix(nums []int) []int {
	p := make([]int, len(nums)+1)
	for i, x := range nums {
		p[i+1] = p[i] + x
	}
	return p
}

// RangeSum returns sum(nums[l..r]) inclusive, given P from BuildPrefix.
// It assumes 0 ≤ l ≤ r < len(nums).
func RangeSum(p []int, l, r int) int {
	return p[r+1] - p[l]
}

// mod returns x mod k in [0, k) for k > 0.
//
// Go gotcha: the % operator keeps the sign of the dividend, so -3 % 5 == -3,
// not 2. Normalise before using a remainder as a map key or index.
func mod(x, k int) int {
	return ((x % k) + k) % k
}
