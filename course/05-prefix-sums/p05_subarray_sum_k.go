package prefixsums

// Problem 05 — Subarray Sum Equals K (Medium)
//
// Count the contiguous, non-empty subarrays of nums whose elements sum to k.
// nums may contain negative numbers and zeros.
//
//	nums = [1, 1, 1],   k = 2 → 2   ([1,1] twice)
//	nums = [1, -1, 0],  k = 0 → 3   ([1,-1], [0], [1,-1,0])
//
// Reference: https://leetcode.com/problems/subarray-sum-equals-k/
//
// Pattern: PREFIX SUM + HASH MAP. A subarray (i, j] sums to k exactly when
// P[j] − P[i] = k, i.e. P[i] = P[j] − k. Walking left to right, keep a map
// from prefix value → how many times it has occurred; at each j, add the
// count of P[j] − k. Seed the map with {0: 1} for subarrays starting at 0.
//
// A sliding window does NOT work here: with negatives, shrinking the window
// can increase the sum, so there is no monotone condition.
//
//	nums = [1, 1, 1], k = 2
//	j  P[j]  need P[j]-k  seen before   count
//	0   1        -1       {0:1}         0
//	1   2         0       {0:1,1:1}     1
//	2   3         1       {0:1,1:1,2:1} 2
//
//	SubarraySumBrute: Time O(n²), Space O(1)
//	SubarraySum:      Time O(n),  Space O(n)

// SubarraySumBrute extends a running sum from every start.
func SubarraySumBrute(nums []int, k int) int {
	count := 0
	for i := range nums {
		sum := 0
		for j := i; j < len(nums); j++ {
			sum += nums[j]
			if sum == k {
				count++
			}
		}
	}
	return count
}

// SubarraySum counts earlier prefixes equal to prefix − k.
func SubarraySum(nums []int, k int) int {
	seen := map[int]int{0: 1} // the empty prefix
	count, prefix := 0, 0
	for _, x := range nums {
		prefix += x
		count += seen[prefix-k] // missing key → 0
		seen[prefix]++          // insert AFTER the lookup: subarrays must be non-empty
	}
	return count
}
