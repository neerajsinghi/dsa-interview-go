package shortestpaths

import (
	"math/rand"
	"testing"
)

func TestSwimInWater(t *testing.T) {
	tests := []struct {
		name string
		grid [][]int
		want int
	}{
		{"two by two", [][]int{{0, 2}, {1, 3}}, 3},
		{"spiral", [][]int{{0, 1, 2, 3, 4}, {24, 23, 22, 21, 5}, {12, 13, 14, 15, 16}, {11, 17, 18, 19, 20}, {10, 9, 8, 7, 6}}, 16},
		{"single cell", [][]int{{0}}, 0},
		{"start is high", [][]int{{3, 0}, {1, 2}}, 3},
	}
	for _, tt := range tests {
		for name, fn := range map[string]func([][]int) int{"dijkstra": SwimInWater, "unionfind": SwimInWaterUnionFind} {
			if got := fn(tt.grid); got != tt.want {
				t.Errorf("%s/%s = %d, want %d", tt.name, name, got, tt.want)
			}
		}
	}
}

func TestSwimInWaterRandom(t *testing.T) {
	r := rand.New(rand.NewSource(50))
	for iter := 0; iter < 200; iter++ {
		n := 1 + r.Intn(6)
		perm := r.Perm(n * n)
		grid := make([][]int, n)
		for i := range grid {
			grid[i] = perm[i*n : (i+1)*n]
		}
		if a, b := SwimInWater(grid), SwimInWaterUnionFind(grid); a != b {
			t.Fatalf("%v: dijkstra=%d unionfind=%d", grid, a, b)
		}
	}
}
