package shortestpaths

import (
	"math/rand"
	"testing"
)

func TestFindTheCity(t *testing.T) {
	tests := []struct {
		name      string
		n         int
		roads     [][]int
		threshold int
		want      int
	}{
		{"classic", 4, [][]int{{0, 1, 3}, {1, 2, 1}, {1, 3, 4}, {2, 3, 1}}, 4, 3},
		{"five cities", 5, [][]int{{0, 1, 2}, {0, 4, 8}, {1, 2, 3}, {1, 4, 2}, {2, 3, 1}, {3, 4, 1}}, 2, 0},
		{"no roads: tie, largest wins", 3, nil, 10, 2},
		{"threshold zero", 2, [][]int{{0, 1, 1}}, 0, 1},
		{"path is shorter than direct road", 3, [][]int{{0, 1, 1}, {1, 2, 1}, {0, 2, 10}}, 2, 2},
	}
	for _, tt := range tests {
		for name, fn := range map[string]func(int, [][]int, int) int{"floyd": FindTheCity, "dijkstra": FindTheCityDijkstra} {
			if got := fn(tt.n, tt.roads, tt.threshold); got != tt.want {
				t.Errorf("%s/%s = %d, want %d", tt.name, name, got, tt.want)
			}
		}
	}
}

func TestFindTheCityRandom(t *testing.T) {
	r := rand.New(rand.NewSource(45))
	for iter := 0; iter < 200; iter++ {
		n := 2 + r.Intn(7)
		var roads [][]int
		for k := r.Intn(14); k > 0; k-- {
			if u, v := r.Intn(n), r.Intn(n); u != v {
				roads = append(roads, []int{u, v, 1 + r.Intn(6)})
			}
		}
		th := r.Intn(12)
		if a, b := FindTheCity(n, roads, th), FindTheCityDijkstra(n, roads, th); a != b {
			t.Fatalf("n=%d th=%d %v: floyd=%d dijkstra=%d", n, th, roads, a, b)
		}
	}
}
