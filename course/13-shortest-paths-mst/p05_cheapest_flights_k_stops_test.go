package shortestpaths

import (
	"math/rand"
	"testing"
)

func TestFindCheapestPrice(t *testing.T) {
	tests := []struct {
		name        string
		n           int
		flights     [][]int
		src, dst, k int
		want        int
	}{
		{"stop limit bites", 4, [][]int{{0, 1, 100}, {1, 2, 100}, {2, 0, 100}, {1, 3, 600}, {2, 3, 200}}, 0, 3, 1, 700},
		{"one stop allowed", 3, [][]int{{0, 1, 100}, {1, 2, 100}, {0, 2, 500}}, 0, 2, 1, 200},
		{"no stops allowed", 3, [][]int{{0, 1, 100}, {1, 2, 100}, {0, 2, 500}}, 0, 2, 0, 500},
		{"unreachable", 3, [][]int{{0, 1, 100}}, 0, 2, 5, -1},
		{"src is dst", 2, [][]int{{0, 1, 5}}, 0, 0, 0, 0},
		{"in-place relaxation trap", 3, [][]int{{0, 1, 1}, {1, 2, 1}, {0, 2, 5}}, 0, 2, 0, 5},
	}
	for _, tt := range tests {
		for name, fn := range map[string]func(int, [][]int, int, int, int) int{"bellmanford": FindCheapestPrice, "brute": FindCheapestPriceBrute} {
			if got := fn(tt.n, tt.flights, tt.src, tt.dst, tt.k); got != tt.want {
				t.Errorf("%s/%s = %d, want %d", tt.name, name, got, tt.want)
			}
		}
	}
}

func TestFindCheapestPriceRandom(t *testing.T) {
	r := rand.New(rand.NewSource(48))
	for iter := 0; iter < 300; iter++ {
		n := 2 + r.Intn(5)
		var flights [][]int
		for i := r.Intn(12); i > 0; i-- {
			if u, v := r.Intn(n), r.Intn(n); u != v {
				flights = append(flights, []int{u, v, 1 + r.Intn(20)})
			}
		}
		src, dst, k := r.Intn(n), r.Intn(n), r.Intn(4)
		if a, b := FindCheapestPrice(n, flights, src, dst, k), FindCheapestPriceBrute(n, flights, src, dst, k); a != b {
			t.Fatalf("n=%d %v %d->%d k=%d: bellmanford=%d brute=%d", n, flights, src, dst, k, a, b)
		}
	}
}
