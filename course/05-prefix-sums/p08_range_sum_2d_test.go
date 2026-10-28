package prefixsums

import (
	"math/rand"
	"testing"
)

func TestNumMatrix(t *testing.T) {
	matrix := [][]int{
		{3, 0, 1, 4, 2},
		{5, 6, 3, 2, 1},
		{1, 2, 0, 1, 5},
		{4, 1, 0, 1, 7},
		{1, 0, 3, 0, 5},
	}
	m := NewNumMatrix(matrix)
	tests := []struct {
		r1, c1, r2, c2, want int
	}{
		{2, 1, 4, 3, 8},
		{1, 1, 2, 2, 11},
		{1, 2, 2, 4, 12},
		{0, 0, 0, 0, 3},
		{0, 0, 4, 4, 58},
		{4, 4, 4, 4, 5},
	}
	for _, tc := range tests {
		if got := m.SumRegion(tc.r1, tc.c1, tc.r2, tc.c2); got != tc.want {
			t.Errorf("SumRegion(%d,%d,%d,%d) = %d, want %d", tc.r1, tc.c1, tc.r2, tc.c2, got, tc.want)
		}
		if got := SumRegionBrute(matrix, tc.r1, tc.c1, tc.r2, tc.c2); got != tc.want {
			t.Errorf("SumRegionBrute(%d,%d,%d,%d) = %d, want %d", tc.r1, tc.c1, tc.r2, tc.c2, got, tc.want)
		}
	}
}

func TestNumMatrixShapes(t *testing.T) {
	// Construction must not panic on degenerate shapes.
	NewNumMatrix(nil)
	NewNumMatrix([][]int{{}})

	row := NewNumMatrix([][]int{{1, -2, 3}})
	if got := row.SumRegion(0, 1, 0, 2); got != 1 {
		t.Errorf("single row: got %d, want 1", got)
	}
	col := NewNumMatrix([][]int{{1}, {-2}, {3}})
	if got := col.SumRegion(0, 0, 2, 0); got != 2 {
		t.Errorf("single column: got %d, want 2", got)
	}
}

func TestNumMatrixRandomized(t *testing.T) {
	rng := rand.New(rand.NewSource(8))
	for range 100 {
		rows, cols := 1+rng.Intn(6), 1+rng.Intn(6)
		matrix := make([][]int, rows)
		for r := range matrix {
			matrix[r] = make([]int, cols)
			for c := range matrix[r] {
				matrix[r][c] = rng.Intn(21) - 10
			}
		}
		m := NewNumMatrix(matrix)
		for range 20 {
			r1, c1 := rng.Intn(rows), rng.Intn(cols)
			r2, c2 := r1+rng.Intn(rows-r1), c1+rng.Intn(cols-c1)
			if a, b := SumRegionBrute(matrix, r1, c1, r2, c2), m.SumRegion(r1, c1, r2, c2); a != b {
				t.Fatalf("%v (%d,%d)-(%d,%d): brute=%d prefix=%d", matrix, r1, c1, r2, c2, a, b)
			}
		}
	}
}
