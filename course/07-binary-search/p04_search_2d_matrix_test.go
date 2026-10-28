package binarysearch

import (
	"math/rand"
	"testing"
)

func TestSearchMatrix(t *testing.T) {
	m := [][]int{{1, 3, 5, 7}, {10, 11, 16, 20}, {23, 30, 34, 60}}
	tests := []struct {
		matrix [][]int
		target int
		want   bool
	}{
		{m, 3, true},
		{m, 13, false},
		{m, 1, true},
		{m, 60, true},
		{m, 0, false},
		{m, 61, false},
		{m, 10, true},
		{[][]int{{5}}, 5, true},
		{[][]int{{1}, {3}, {5}}, 4, false},
		{nil, 1, false},
		{[][]int{{}}, 1, false},
	}
	for _, tt := range tests {
		for name, fn := range map[string]func([][]int, int) bool{
			"brute": SearchMatrixBrute, "flat": SearchMatrix,
		} {
			if got := fn(tt.matrix, tt.target); got != tt.want {
				t.Errorf("%s(%v, %d) = %v, want %v", name, tt.matrix, tt.target, got, tt.want)
			}
		}
	}
}

func TestSearchMatrixRandom(t *testing.T) {
	r := rand.New(rand.NewSource(4))
	for iter := 0; iter < 500; iter++ {
		rows, cols := 1+r.Intn(5), 1+r.Intn(5)
		matrix := make([][]int, rows)
		v := 0
		for i := range matrix {
			matrix[i] = make([]int, cols)
			for j := range matrix[i] {
				v += 1 + r.Intn(3) // strictly increasing in reading order
				matrix[i][j] = v
			}
		}
		target := r.Intn(v + 2)
		if got, want := SearchMatrix(matrix, target), SearchMatrixBrute(matrix, target); got != want {
			t.Fatalf("SearchMatrix(%v, %d) = %v, want %v", matrix, target, got, want)
		}
	}
}
