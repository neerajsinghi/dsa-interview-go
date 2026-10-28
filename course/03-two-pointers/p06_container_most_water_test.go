package twopointers

import (
	"math/rand"
	"testing"
)

func TestMaxArea(t *testing.T) {
	tests := []struct {
		name    string
		heights []int
		want    int
	}{
		{"empty", nil, 0},
		{"single", []int{5}, 0},
		{"two", []int{1, 1}, 1},
		{"classic", []int{1, 8, 6, 2, 5, 4, 8, 3, 7}, 49},
		{"increasing", []int{1, 2, 3, 4, 5}, 6},
		{"equal heights", []int{3, 3, 3, 3}, 9},
		{"zeros", []int{0, 0, 0}, 0},
		{"tall middle", []int{1, 100, 100, 1}, 100},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := ContainerBrute(tc.heights); got != tc.want {
				t.Errorf("brute = %d, want %d", got, tc.want)
			}
			if got := MaxArea(tc.heights); got != tc.want {
				t.Errorf("two pointers = %d, want %d", got, tc.want)
			}
		})
	}
}

func TestMaxAreaRandomized(t *testing.T) {
	rng := rand.New(rand.NewSource(6))
	for range 1000 {
		h := make([]int, rng.Intn(15))
		for i := range h {
			h[i] = rng.Intn(10)
		}
		if a, b := ContainerBrute(h), MaxArea(h); a != b {
			t.Fatalf("%v: brute=%d two pointers=%d", h, a, b)
		}
	}
}
