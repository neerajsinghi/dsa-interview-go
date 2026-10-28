package graphs

import (
	"math/rand"
	"testing"
)

func TestOrangesRotting(t *testing.T) {
	tests := []struct {
		name string
		grid [][]int
		want int
	}{
		{"classic", [][]int{{2, 1, 1}, {1, 1, 0}, {0, 1, 1}}, 4},
		{"unreachable", [][]int{{2, 1, 1}, {0, 1, 1}, {1, 0, 1}}, -1},
		{"no fresh", [][]int{{0, 2}}, 0},
		{"nothing at all", [][]int{{0}}, 0},
		{"fresh but no rotten", [][]int{{1}}, -1},
		{"two sources meet in the middle", [][]int{{2, 1, 1, 1, 2}}, 2},
		{"empty", nil, 0},
	}
	for _, tt := range tests {
		for name, fn := range map[string]func([][]int) int{"bfs": OrangesRotting, "simulate": OrangesRottingSimulate} {
			if got := fn(tt.grid); got != tt.want {
				t.Errorf("%s/%s = %d, want %d", tt.name, name, got, tt.want)
			}
		}
	}
}

func TestOrangesRottingRandom(t *testing.T) {
	r := rand.New(rand.NewSource(25))
	for iter := 0; iter < 300; iter++ {
		m, n := 1+r.Intn(6), 1+r.Intn(6)
		grid := make([][]int, m)
		for i := range grid {
			grid[i] = make([]int, n)
			for j := range grid[i] {
				grid[i][j] = []int{0, 1, 1, 1, 2}[r.Intn(5)]
			}
		}
		if a, b := OrangesRotting(grid), OrangesRottingSimulate(grid); a != b {
			t.Fatalf("%v: bfs=%d simulate=%d", grid, a, b)
		}
	}
}
