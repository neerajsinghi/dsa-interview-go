package complexity

// Problem 04 — Amortized Append (Easy)
//
// Show why appending n elements to a Go slice costs O(n) in total even though
// an individual append occasionally copies the whole slice.
//
//	CountGrowths(1)    → 1   (cap 0 → 1)
//	CountGrowths(1000) → ~11 (cap roughly doubles each time)
//
// Pattern: amortized analysis. When the backing array is full, append
// allocates a bigger one (about 2× for small slices, ~1.25× for large ones)
// and copies everything over. With geometric growth the copies form a series
// 1 + 2 + 4 + … + n < 2n, so the TOTAL cost of n appends is O(n), i.e.
// O(1) amortized per append. Contrast with inserting at the FRONT, which
// shifts every element on every insert: O(n) each, O(n²) total.
//
//	   len: 1  2  3  4  5  6  7  8  9 ...
//	   cap: 1  2  4  4  8  8  8  8 16 ...
//	copies: 0  1  2  0  4  0  0  0  8 ...   (sum ≤ 2n)
//
//	AppendN:             Time O(n) total, O(1) amortized per append
//	AppendNPreallocated: Time O(n), no regrowth at all
//	PrependN:            Time O(n²)

// CountGrowths appends n ints to a nil slice and counts how many times the
// capacity changed — that is, how many reallocate-and-copy events happened.
func CountGrowths(n int) int {
	var s []int
	growths := 0
	for i := range n {
		before := cap(s)
		s = append(s, i)
		if cap(s) != before {
			growths++
		}
	}
	return growths
}

// AppendN builds [0, 1, ..., n-1] letting append grow the slice.
func AppendN(n int) []int {
	var s []int
	for i := range n {
		s = append(s, i)
	}
	return s
}

// AppendNPreallocated does the same with make(..., 0, n): when you know the
// final size, preallocating skips every regrowth.
func AppendNPreallocated(n int) []int {
	s := make([]int, 0, n)
	for i := range n {
		s = append(s, i)
	}
	return s
}

// PrependN builds [n-1, ..., 1, 0] by inserting each value at index 0.
// Each insert shifts the whole slice right by one.
func PrependN(n int) []int {
	var s []int
	for i := range n {
		s = append(s, 0) // make room at the end
		copy(s[1:], s)   // shift everything right: O(len)
		s[0] = i
	}
	return s
}
