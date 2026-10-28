package dp2d

import (
	"math/rand"
	"testing"
)

func TestMinPathSum(t *testing.T) {
	tests := []struct {
		name string
		grid [][]int
		want int
	}{
		{"classic", [][]int{{1, 3, 1}, {1, 5, 1}, {4, 2, 1}}, 7},
		{"two rows", [][]int{{1, 2, 3}, {4, 5, 6}}, 12},
		{"single", [][]int{{5}}, 5},
		{"one column", [][]int{{1}, {2}, {3}}, 6},
		{"zeros", [][]int{{0, 0}, {0, 0}}, 0},
		{"empty", [][]int{}, 0},
	}
	impls := map[string]func([][]int) int{"memo": MinPathSumMemo, "tab": MinPathSumTab, "rolling": MinPathSumRolling}
	for _, tt := range tests {
		for name, fn := range impls {
			if got := fn(tt.grid); got != tt.want {
				t.Errorf("%s/%s = %d, want %d", tt.name, name, got, tt.want)
			}
		}
	}
}

func TestMinPathSumAgree(t *testing.T) {
	r := rand.New(rand.NewSource(2))
	for iter := 0; iter < 200; iter++ {
		m, n := 1+r.Intn(8), 1+r.Intn(8)
		grid := make([][]int, m)
		for i := range grid {
			grid[i] = make([]int, n)
			for j := range grid[i] {
				grid[i][j] = r.Intn(10)
			}
		}
		a, b, c := MinPathSumMemo(grid), MinPathSumTab(grid), MinPathSumRolling(grid)
		if a != b || b != c {
			t.Fatalf("grid %v: memo=%d tab=%d rolling=%d", grid, a, b, c)
		}
	}
}
