package prefixsums

// Problem 08 — Range Sum Query 2D, Immutable (Medium)
//
// Build a structure over a fixed matrix that answers many queries "what is the
// sum of the rectangle with top-left (r1, c1) and bottom-right (r2, c2)?"
// (inclusive).
//
//	matrix = [[3, 0, 1],
//	          [5, 6, 3],
//	          [1, 2, 0]]
//	SumRegion(1, 1, 2, 2) → 11   (6 + 3 + 2 + 0)
//	SumRegion(0, 0, 0, 0) → 3
//
// Reference: https://leetcode.com/problems/range-sum-query-2d-immutable/
//
// Pattern: 2-D prefix sums with a padding row and column of zeros.
// P[i][j] = sum of the rectangle matrix[0..i-1][0..j-1]. Inclusion–exclusion:
//
//	P[i][j] = matrix[i-1][j-1] + P[i-1][j] + P[i][j-1] − P[i-1][j-1]
//
//	region(r1,c1,r2,c2) = P[r2+1][c2+1] − P[r1][c2+1] − P[r2+1][c1] + P[r1][c1]
//
//	+---------+---+
//	|    D    | B |      wanted = whole − top − left + corner
//	+---------+---+      (corner D was subtracted twice, so add it back)
//	|    C    | X |
//	+---------+---+
//
//	SumRegionBrute: Time O(rows·cols) per query, Space O(1)
//	NumMatrix:      Time O(rows·cols) build, O(1) per query, Space O(rows·cols)

// SumRegionBrute adds every cell in the rectangle.
func SumRegionBrute(matrix [][]int, r1, c1, r2, c2 int) int {
	sum := 0
	for r := r1; r <= r2; r++ {
		for c := c1; c <= c2; c++ {
			sum += matrix[r][c]
		}
	}
	return sum
}

// NumMatrix answers rectangle-sum queries in O(1).
type NumMatrix struct {
	p [][]int // (rows+1) × (cols+1), first row and column are zero
}

// NewNumMatrix precomputes 2-D prefix sums. An empty matrix (or one with
// empty rows) is allowed; it simply has no valid queries.
func NewNumMatrix(matrix [][]int) NumMatrix {
	rows := len(matrix)
	cols := 0
	if rows > 0 {
		cols = len(matrix[0])
	}
	p := make([][]int, rows+1)
	for i := range p {
		p[i] = make([]int, cols+1) // each row its own slice: no shared backing array
	}
	for i := 1; i <= rows; i++ {
		for j := 1; j <= cols; j++ {
			p[i][j] = matrix[i-1][j-1] + p[i-1][j] + p[i][j-1] - p[i-1][j-1]
		}
	}
	return NumMatrix{p: p}
}

// SumRegion returns the sum of matrix[r1..r2][c1..c2].
func (m NumMatrix) SumRegion(r1, c1, r2, c2 int) int {
	p := m.p
	return p[r2+1][c2+1] - p[r1][c2+1] - p[r2+1][c1] + p[r1][c1]
}
