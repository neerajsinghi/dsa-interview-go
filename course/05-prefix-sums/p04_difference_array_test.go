package prefixsums

import (
	"math/rand"
	"slices"
	"testing"
)

func TestApplyUpdates(t *testing.T) {
	tests := []struct {
		name    string
		n       int
		updates [][3]int
		want    []int
	}{
		{"no positions", 0, nil, []int{}},
		{"no updates", 3, nil, []int{0, 0, 0}},
		{"classic", 5, [][3]int{{1, 3, 2}, {2, 4, 3}, {0, 2, -2}}, []int{-2, 0, 3, 5, 3}},
		{"single cell", 3, [][3]int{{1, 1, 7}}, []int{0, 7, 0}},
		{"whole range", 4, [][3]int{{0, 3, 1}}, []int{1, 1, 1, 1}},
		{"cancelling", 3, [][3]int{{0, 2, 5}, {0, 2, -5}}, []int{0, 0, 0}},
		{"overlapping", 4, [][3]int{{0, 1, 1}, {1, 2, 1}, {2, 3, 1}}, []int{1, 2, 2, 1}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := ApplyUpdatesBrute(tc.n, tc.updates); !slices.Equal(got, tc.want) {
				t.Errorf("brute = %v, want %v", got, tc.want)
			}
			if got := ApplyUpdates(tc.n, tc.updates); !slices.Equal(got, tc.want) {
				t.Errorf("difference array = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestApplyUpdatesRandomized(t *testing.T) {
	rng := rand.New(rand.NewSource(4))
	for range 500 {
		n := 1 + rng.Intn(15)
		updates := make([][3]int, rng.Intn(10))
		for i := range updates {
			s := rng.Intn(n)
			e := s + rng.Intn(n-s)
			updates[i] = [3]int{s, e, rng.Intn(21) - 10}
		}
		if a, b := ApplyUpdatesBrute(n, updates), ApplyUpdates(n, updates); !slices.Equal(a, b) {
			t.Fatalf("n=%d %v: brute=%v diff=%v", n, updates, a, b)
		}
	}
}
