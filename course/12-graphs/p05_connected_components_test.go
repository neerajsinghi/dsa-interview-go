package graphs

import (
	"math/rand"
	"testing"
)

func TestCountComponents(t *testing.T) {
	tests := []struct {
		name  string
		n     int
		edges [][]int
		want  int
	}{
		{"two groups", 5, [][]int{{0, 1}, {1, 2}, {3, 4}}, 2},
		{"one path", 5, [][]int{{0, 1}, {1, 2}, {2, 3}, {3, 4}}, 1},
		{"no edges", 3, nil, 3},
		{"no nodes", 0, nil, 0},
		{"cycle plus singleton", 4, [][]int{{0, 1}, {1, 2}, {2, 0}}, 2},
		{"duplicate edges", 2, [][]int{{0, 1}, {1, 0}}, 1},
	}
	for _, tt := range tests {
		for name, fn := range map[string]func(int, [][]int) int{"unionfind": CountComponents, "dfs": CountComponentsDFS} {
			if got := fn(tt.n, tt.edges); got != tt.want {
				t.Errorf("%s/%s = %d, want %d", tt.name, name, got, tt.want)
			}
		}
	}
}

func TestCountComponentsRandom(t *testing.T) {
	r := rand.New(rand.NewSource(27))
	for iter := 0; iter < 300; iter++ {
		n := 1 + r.Intn(10)
		var edges [][]int
		for k := r.Intn(10); k > 0; k-- {
			edges = append(edges, []int{r.Intn(n), r.Intn(n)})
		}
		if a, b := CountComponents(n, edges), CountComponentsDFS(n, edges); a != b {
			t.Fatalf("n=%d %v: unionfind=%d dfs=%d", n, edges, a, b)
		}
	}
}
