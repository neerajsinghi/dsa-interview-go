package prefixsums

// Problem 01 — Range Sum Query, Immutable (Easy)
//
// Build a structure over a fixed slice nums that answers many queries
// "what is nums[l] + … + nums[r]?" (inclusive, 0 ≤ l ≤ r < n).
//
//	nums = [-2, 0, 3, -5, 2, -1]
//	SumRange(0, 2) → 1     (-2 + 0 + 3)
//	SumRange(2, 5) → -1    (3 - 5 + 2 - 1)
//
// Reference: https://leetcode.com/problems/range-sum-query-immutable/
//
// Pattern: pay O(n) once to build prefix sums, then every query is one
// subtraction. Worth it as soon as there is more than a handful of queries.
//
//	   nums:    -2   0   3  -5   2  -1
//	      P: 0  -2  -2   1  -4  -2  -3
//	SumRange(2, 5) = P[6] − P[2] = -3 − (-2) = -1
//
//	SumRangeBrute: Time O(n) per query, Space O(1)
//	NumArray:      Time O(n) build, O(1) per query, Space O(n)

// SumRangeBrute adds the range directly.
func SumRangeBrute(nums []int, l, r int) int {
	sum := 0
	for _, x := range nums[l : r+1] {
		sum += x
	}
	return sum
}

// NumArray answers range-sum queries in O(1).
type NumArray struct {
	prefix []int
}

// NewNumArray precomputes prefix sums. Later changes to nums are not seen:
// the prefix slice is a separate allocation.
func NewNumArray(nums []int) NumArray {
	return NumArray{prefix: BuildPrefix(nums)}
}

// SumRange returns nums[l] + … + nums[r].
func (a NumArray) SumRange(l, r int) int {
	return RangeSum(a.prefix, l, r)
}
