package twopointers

// Problem 04 — Sort Colors / Dutch National Flag (Medium)
//
// nums contains only 0s, 1s and 2s. Sort it in place in a single pass using
// O(1) extra space, without calling a library sort.
//
//	[2, 0, 2, 1, 1, 0] → [0, 0, 1, 1, 2, 2]
//	[2, 0, 1]          → [0, 1, 2]
//
// Reference: https://leetcode.com/problems/sort-colors/
//
// Pattern: THREE pointers partitioning the slice into four regions:
//
//	[0 … low)   all 0s
//	[low … mid) all 1s
//	[mid … high] unknown, still to examine
//	(high … n)  all 2s
//
// Look at nums[mid]: a 0 is swapped down to `low` (both advance), a 1 just
// advances mid, a 2 is swapped up to `high` (only high retreats — the value
// swapped in from high is unexamined, so mid must look at it next).
//
//	SortColorsCounting: Time O(n) two passes, Space O(1)
//	SortColors:         Time O(n) one pass,   Space O(1)

// SortColorsCounting counts each colour and rewrites the slice.
func SortColorsCounting(nums []int) {
	var count [3]int
	for _, x := range nums {
		count[x]++
	}
	i := 0
	for color, c := range count {
		for range c {
			nums[i] = color
			i++
		}
	}
}

// SortColors is Dijkstra's one-pass Dutch national flag partition.
func SortColors(nums []int) {
	low, mid, high := 0, 0, len(nums)-1
	for mid <= high {
		switch nums[mid] {
		case 0:
			nums[low], nums[mid] = nums[mid], nums[low]
			low++
			mid++
		case 1:
			mid++
		case 2:
			nums[mid], nums[high] = nums[high], nums[mid]
			high-- // don't advance mid: the swapped-in value is unexamined
		}
	}
}
