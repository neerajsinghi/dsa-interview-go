package dp2d

import (
	"math/rand"
	"testing"
)

func TestUniquePaths(t *testing.T) {
	tests := []struct {
		name string
		grid [][]int
		want int
	}{
		{"single open cell", [][]int{{0}}, 1},
		{"single blocked cell", [][]int{{1}}, 0},
		{"3x3 empty", [][]int{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}}, 6},
		{"centre blocked", [][]int{{0, 0, 0}, {0, 1, 0}, {0, 0, 0}}, 2},
		{"walled off", [][]int{{0, 1}, {1, 0}}, 0},
		{"start blocked", [][]int{{1, 0}, {0, 0}}, 0},
		{"end blocked", [][]int{{0, 0}, {0, 1}}, 0},
		{"one row", [][]int{{0, 0, 0, 0}}, 1},
		{"one row blocked", [][]int{{0, 1, 0}}, 0},
		{"empty", [][]int{}, 0},
	}
	impls := map[string]func([][]int) int{
		"brute": UniquePathsBrute, "memo": UniquePathsMemo,
		"tab": UniquePathsTab, "rolling": UniquePathsRolling,
	}
	for _, tt := range tests {
		for name, fn := range impls {
			if got := fn(tt.grid); got != tt.want {
				t.Errorf("%s/%s = %d, want %d", tt.name, name, got, tt.want)
			}
		}
	}
}

func TestUniquePathsEmptyGrid(t *testing.T) {
	tests := []struct{ m, n, want int }{{1, 1, 1}, {3, 7, 28}, {3, 2, 3}, {0, 5, 0}, {10, 10, 48620}}
	for _, tt := range tests {
		if got := UniquePaths(tt.m, tt.n); got != tt.want {
			t.Errorf("UniquePaths(%d,%d) = %d, want %d", tt.m, tt.n, got, tt.want)
		}
	}
}

func TestUniquePathsRandom(t *testing.T) {
	r := rand.New(rand.NewSource(16))
	for iter := 0; iter < 300; iter++ {
		m, n := 1+r.Intn(6), 1+r.Intn(6)
		grid := make([][]int, m)
		for i := range grid {
			grid[i] = make([]int, n)
			for j := range grid[i] {
				if r.Intn(5) == 0 {
					grid[i][j] = 1
				}
			}
		}
		want := UniquePathsBrute(grid)
		if a, b, c := UniquePathsMemo(grid), UniquePathsTab(grid), UniquePathsRolling(grid); a != want || b != want || c != want {
			t.Fatalf("grid %v: brute=%d memo=%d tab=%d rolling=%d", grid, want, a, b, c)
		}
	}
}
