package shortestpaths

import (
	"math/rand"
	"testing"
)

func TestNetworkDelayTime(t *testing.T) {
	tests := []struct {
		name  string
		times [][]int
		n, k  int
		want  int
	}{
		{"classic", [][]int{{2, 1, 1}, {2, 3, 1}, {3, 4, 1}}, 4, 2, 2},
		{"single link", [][]int{{1, 2, 1}}, 2, 1, 1},
		{"unreachable", [][]int{{1, 2, 1}}, 2, 2, -1},
		{"single node", nil, 1, 1, 0},
		{"detour is faster", [][]int{{1, 2, 10}, {1, 3, 1}, {3, 2, 2}}, 3, 1, 3},
		{"zero weights", [][]int{{1, 2, 0}, {2, 3, 0}}, 3, 1, 0},
	}
	for _, tt := range tests {
		for name, fn := range map[string]func([][]int, int, int) int{"dijkstra": NetworkDelayTime, "bellmanford": NetworkDelayTimeBellmanFord} {
			if got := fn(tt.times, tt.n, tt.k); got != tt.want {
				t.Errorf("%s/%s = %d, want %d", tt.name, name, got, tt.want)
			}
		}
	}
}

func TestNetworkDelayTimeRandom(t *testing.T) {
	r := rand.New(rand.NewSource(44))
	for iter := 0; iter < 300; iter++ {
		n := 1 + r.Intn(7)
		times := make([][]int, r.Intn(20))
		for i := range times {
			times[i] = []int{1 + r.Intn(n), 1 + r.Intn(n), r.Intn(10)}
		}
		k := 1 + r.Intn(n)
		if a, b := NetworkDelayTime(times, n, k), NetworkDelayTimeBellmanFord(times, n, k); a != b {
			t.Fatalf("n=%d k=%d %v: dijkstra=%d bellmanford=%d", n, k, times, a, b)
		}
	}
}
