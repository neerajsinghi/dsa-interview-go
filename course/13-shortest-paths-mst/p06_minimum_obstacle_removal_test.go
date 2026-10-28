package shortestpaths

import (
	"math/rand"
	"testing"
)

func TestMinimumObstacles(t *testing.T) {
	tests := []struct {
		name string
		grid [][]int
		want int
	}{
		{"must break two", [][]int{{0, 1, 1}, {1, 1, 0}, {1, 1, 0}}, 2},
		{"path around", [][]int{{0, 1, 0, 0, 0}, {0, 1, 0, 1, 0}, {0, 0, 0, 1, 0}}, 0},
		{"single cell", [][]int{{0}}, 0},
		{"one wall column", [][]int{{0, 1, 0}, {0, 1, 0}}, 1},
		{"single row", [][]int{{0, 1, 1, 0}}, 2},
	}
	for _, tt := range tests {
		for name, fn := range map[string]func([][]int) int{"01bfs": MinimumObstacles, "dijkstra": MinimumObstaclesDijkstra} {
			if got := fn(tt.grid); got != tt.want {
				t.Errorf("%s/%s = %d, want %d", tt.name, name, got, tt.want)
			}
		}
	}
}

func TestMinimumObstaclesRandom(t *testing.T) {
	r := rand.New(rand.NewSource(49))
	for iter := 0; iter < 300; iter++ {
		m, n := 1+r.Intn(7), 1+r.Intn(7)
		grid := make([][]int, m)
		for i := range grid {
			grid[i] = make([]int, n)
			for j := range grid[i] {
				grid[i][j] = r.Intn(2)
			}
		}
		grid[0][0], grid[m-1][n-1] = 0, 0
		if a, b := MinimumObstacles(grid), MinimumObstaclesDijkstra(grid); a != b {
			t.Fatalf("%v: 01bfs=%d dijkstra=%d", grid, a, b)
		}
	}
}
