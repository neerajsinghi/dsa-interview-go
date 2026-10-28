package dp2d

import (
	"math/rand"
	"testing"
)

func TestKnapsack(t *testing.T) {
	tests := []struct {
		name     string
		w, v     []int
		capacity int
		want     int
	}{
		{"classic", []int{1, 3, 4, 5}, []int{1, 4, 5, 7}, 7, 9},
		{"too heavy", []int{5}, []int{10}, 4, 0},
		{"exact fit", []int{5}, []int{10}, 5, 10},
		{"no items", nil, nil, 10, 0},
		{"zero capacity", []int{1, 2}, []int{3, 4}, 0, 0},
		// Each item at most once: rolling forward would wrongly give 30.
		{"no reuse", []int{2}, []int{10}, 6, 10},
		{"zero-weight item", []int{0, 3}, []int{5, 4}, 2, 5},
	}
	impls := map[string]func([]int, []int, int) int{"brute": KnapsackBrute, "memo": KnapsackMemo, "tab": KnapsackTab, "rolling": KnapsackRolling}
	for _, tt := range tests {
		for name, fn := range impls {
			if got := fn(tt.w, tt.v, tt.capacity); got != tt.want {
				t.Errorf("%s/%s = %d, want %d", tt.name, name, got, tt.want)
			}
		}
	}
}

func TestKnapsackRandom(t *testing.T) {
	r := rand.New(rand.NewSource(5))
	for iter := 0; iter < 300; iter++ {
		n := r.Intn(10)
		w, v := make([]int, n), make([]int, n)
		for i := range w {
			w[i], v[i] = 1+r.Intn(8), r.Intn(20)
		}
		c := r.Intn(25)
		want := KnapsackBrute(w, v, c)
		if a, b, d := KnapsackMemo(w, v, c), KnapsackTab(w, v, c), KnapsackRolling(w, v, c); a != want || b != want || d != want {
			t.Fatalf("w=%v v=%v c=%d: brute=%d memo=%d tab=%d rolling=%d", w, v, c, want, a, b, d)
		}
	}
}
