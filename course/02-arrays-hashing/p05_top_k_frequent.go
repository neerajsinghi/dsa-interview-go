package arrayshashing

import (
	"cmp"
	"slices"
)

// Problem 05 — Top K Frequent Elements (Medium)
//
// Return the k values that occur most often in nums, most frequent first.
// Ties in frequency are broken by smaller value first so the output is
// deterministic. Assume 0 ≤ k ≤ number of distinct values.
//
//	nums = [1,1,1,2,2,3], k = 2 → [1, 2]
//	nums = [4],           k = 1 → [4]
//
// Reference: https://leetcode.com/problems/top-k-frequent-elements/
//
// Pattern: count with a map, then avoid a full O(n log n) sort with BUCKET
// SORT: a frequency is at most n, so make n+1 buckets where bucket[f] holds
// the values seen exactly f times, and walk the buckets from high to low.
//
//	  nums = [1,1,1,2,2,3]
//	counts = {1:3, 2:2, 3:1}
//	bucket:  0:[]  1:[3]  2:[2]  3:[1]  ...  ← read right to left: 1, 2, 3
//
//	TopKFrequentSort: Time O(n + d log d), Space O(n)   (d = distinct values)
//	TopKFrequent:     Time O(n) (+ sorting inside tied buckets), Space O(n)

// TopKFrequentSort counts then sorts distinct values by (count desc, value asc).
func TopKFrequentSort(nums []int, k int) []int {
	counts := Counter(nums)
	vals := make([]int, 0, len(counts))
	for v := range counts {
		vals = append(vals, v)
	}
	slices.SortFunc(vals, func(a, b int) int {
		return cmp.Or(cmp.Compare(counts[b], counts[a]), cmp.Compare(a, b))
	})
	return vals[:k]
}

// TopKFrequent uses buckets indexed by frequency.
func TopKFrequent(nums []int, k int) []int {
	counts := Counter(nums)
	buckets := make([][]int, len(nums)+1) // frequency ranges over 1..n
	for v, c := range counts {
		buckets[c] = append(buckets[c], v)
	}
	res := make([]int, 0, k)
	for f := len(buckets) - 1; f > 0 && len(res) < k; f-- {
		b := buckets[f]
		slices.Sort(b) // map iteration order is random: fix tie order
		for _, v := range b {
			if len(res) == k {
				break
			}
			res = append(res, v)
		}
	}
	return res
}
