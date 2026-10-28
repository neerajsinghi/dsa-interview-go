package arrayshashing

// Problem 06 — Product of Array Except Self (Medium)
//
// Return a slice out where out[i] is the product of every element of nums
// except nums[i]. Do it without division (which also breaks on zeros).
// Assume products fit in an int.
//
//	[1, 2, 3, 4]    → [24, 12, 8, 6]
//	[-1, 1, 0, -3]  → [0, 0, 3, 0]
//
// Reference: https://leetcode.com/problems/product-of-array-except-self/
//
// Pattern: prefix/suffix products. out[i] = (product of everything LEFT of i)
// × (product of everything RIGHT of i). Fill out with left products in one
// pass, then sweep from the right multiplying by a running suffix product.
//
//	nums:    1    2    3    4
//	left:    1    1    2    6     (product of nums[0..i-1])
//	right:  24   12    4    1     (product of nums[i+1..n-1])
//	out:    24   12    8    6     left × right
//
//	ProductExceptSelfBrute: Time O(n²), Space O(1) extra
//	ProductExceptSelf:      Time O(n),  Space O(1) extra (output not counted)

// ProductExceptSelfBrute multiplies everything but i, for each i.
func ProductExceptSelfBrute(nums []int) []int {
	out := make([]int, len(nums))
	for i := range nums {
		p := 1
		for j, x := range nums {
			if j != i {
				p *= x
			}
		}
		out[i] = p
	}
	return out
}

// ProductExceptSelf uses the output slice for the prefix pass.
func ProductExceptSelf(nums []int) []int {
	n := len(nums)
	out := make([]int, n)
	left := 1
	for i := 0; i < n; i++ {
		out[i] = left // product of nums[0..i-1]
		left *= nums[i]
	}
	right := 1
	for i := n - 1; i >= 0; i-- {
		out[i] *= right // times product of nums[i+1..n-1]
		right *= nums[i]
	}
	return out
}
