package binarysearch

// Problem: Search a 2D Matrix
// Difficulty: Medium
//
// An m×n matrix has each row sorted ascending, and the first value of every
// row is greater than the last value of the row above. So read row by row,
// the whole matrix is one sorted sequence. Report whether target is in it, in
// O(log(m·n)).
//
// Example:
//
//	[[1, 3, 5, 7],
//	 [10,11,16,20],
//	 [23,30,34,60]]   target=3 -> true, target=13 -> false
//
// Reference: https://leetcode.com/problems/search-a-2d-matrix/
//
// Pattern: index mapping. Treat the matrix as a virtual flat array of length
// m·n; flat index k lives at row k/n, column k%n. Then it is plain binary search.
// (For matrices that are only sorted per row *and* per column, "Search a 2D
// Matrix II", use the O(m+n) staircase walk from the top-right corner instead.)
//
// Time: brute O(m·n); optimal O(log(m·n)). Space: O(1).

// SearchMatrixBrute checks every cell.
func SearchMatrixBrute(matrix [][]int, target int) bool {
	for _, row := range matrix {
		for _, v := range row {
			if v == target {
				return true
			}
		}
	}
	return false
}

// SearchMatrix binary-searches the matrix as one flattened sorted array.
func SearchMatrix(matrix [][]int, target int) bool {
	if len(matrix) == 0 || len(matrix[0]) == 0 {
		return false
	}
	m, n := len(matrix), len(matrix[0])
	lo, hi := 0, m*n-1
	for lo <= hi {
		mid := lo + (hi-lo)/2
		v := matrix[mid/n][mid%n] // flat index -> (row, col)
		switch {
		case v == target:
			return true
		case v < target:
			lo = mid + 1
		default:
			hi = mid - 1
		}
	}
	return false
}
