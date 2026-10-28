package shortestpaths

import (
	"math/rand"
	"testing"
)

func TestMinimumEffortPath(t *testing.T) {
	tests := []struct {
		name    string
		heights [][]int
		want    int
	}{
		{"classic", [][]int{{1, 2, 2}, {3, 8, 2}, {5, 3, 5}}, 2},
		{"along the edge", [][]int{{1, 2, 3}, {3, 8, 4}, {5, 3, 5}}, 1},
		{"flat path exists", [][]int{{1, 2, 1, 1, 1}, {1, 2, 1, 2, 1}, {1, 2, 1, 2, 1}, {1, 2, 1, 2, 1}, {1, 1, 1, 2, 1}}, 0},
		{"single cell", [][]int{{42}}, 0},
		{"single row", [][]int{{1, 10, 6}}, 9},
		{"wall forces climb", [][]int{{1, 100}, {100, 1}}, 99},
	}
	for _, tt := range tests {
		for name, fn := range map[string]func([][]int) int{"dijkstra": MinimumEffortPath, "binarysearch": MinimumEffortPathBinarySearch} {
			if got := fn(tt.heights); got != tt.want {
				t.Errorf("%s/%s = %d, want %d", tt.name, name, got, tt.want)
			}
		}
	}
}

func TestMinimumEffortPathRandom(t *testing.T) {
	r := rand.New(rand.NewSource(47))
	for iter := 0; iter < 200; iter++ {
		m, n := 1+r.Intn(6), 1+r.Intn(6)
		h := make([][]int, m)
		for i := range h {
			h[i] = make([]int, n)
			for j := range h[i] {
				h[i][j] = 1 + r.Intn(20)
			}
		}
		if a, b := MinimumEffortPath(h), MinimumEffortPathBinarySearch(h); a != b {
			t.Fatalf("%v: dijkstra=%d binarysearch=%d", h, a, b)
		}
	}
}
