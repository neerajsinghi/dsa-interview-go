package prefixsums

import (
	"math/rand"
	"testing"
)

func TestPivotIndex(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want int
	}{
		{"empty", nil, -1},
		{"single", []int{5}, 0}, // both sides empty
		{"classic", []int{1, 7, 3, 6, 5, 6}, 3},
		{"none", []int{1, 2, 3}, -1},
		{"left edge", []int{2, 1, -1}, 0},
		{"right edge", []int{-1, 1, 2}, 2},
		{"leftmost of several", []int{0, 0, 0}, 0},
		{"negatives", []int{-1, -1, -1, -1, -1, 0}, 2},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := PivotIndexBrute(tc.nums); got != tc.want {
				t.Errorf("brute = %d, want %d", got, tc.want)
			}
			if got := PivotIndex(tc.nums); got != tc.want {
				t.Errorf("prefix = %d, want %d", got, tc.want)
			}
		})
	}
}

func TestPivotIndexRandomized(t *testing.T) {
	rng := rand.New(rand.NewSource(2))
	for range 1000 {
		nums := make([]int, rng.Intn(10))
		for i := range nums {
			nums[i] = rng.Intn(7) - 3
		}
		if a, b := PivotIndexBrute(nums), PivotIndex(nums); a != b {
			t.Fatalf("%v: brute=%d prefix=%d", nums, a, b)
		}
	}
}
