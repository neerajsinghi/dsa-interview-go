package twopointers

import (
	"math/rand"
	"testing"
)

func TestTrap(t *testing.T) {
	tests := []struct {
		name    string
		heights []int
		want    int
	}{
		{"empty", nil, 0},
		{"single", []int{4}, 0},
		{"two", []int{4, 2}, 0},
		{"classic", []int{0, 1, 0, 2, 1, 0, 1, 3, 2, 1, 2, 1}, 6},
		{"bowl", []int{4, 2, 0, 3, 2, 5}, 9},
		{"increasing", []int{1, 2, 3, 4}, 0},
		{"decreasing", []int{4, 3, 2, 1}, 0},
		{"flat", []int{2, 2, 2}, 0},
		{"tall walls", []int{1000000, 0, 1000000}, 1000000},
	}
	impls := map[string]func([]int) int{
		"Brute": TrapBrute, "PrefixMax": TrapPrefixMax, "TwoPointers": Trap,
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			for name, fn := range impls {
				if got := fn(tc.heights); got != tc.want {
					t.Errorf("%s = %d, want %d", name, got, tc.want)
				}
			}
		})
	}
}

func TestTrapRandomized(t *testing.T) {
	rng := rand.New(rand.NewSource(7))
	for range 1000 {
		h := make([]int, rng.Intn(20))
		for i := range h {
			h[i] = rng.Intn(8)
		}
		want := TrapBrute(h)
		if a, b := TrapPrefixMax(h), Trap(h); a != want || b != want {
			t.Fatalf("%v: brute=%d prefix=%d two pointers=%d", h, want, a, b)
		}
	}
}
