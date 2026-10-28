package arrayshashing

import (
	"math"
	"math/rand"
	"testing"
)

func TestContainsDuplicate(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want bool
	}{
		{"empty", nil, false},
		{"single", []int{1}, false},
		{"distinct", []int{1, 2, 3, 4}, false},
		{"dup at ends", []int{1, 2, 3, 1}, true},
		{"many dups", []int{1, 1, 1, 3, 3, 4, 3, 2, 4, 2}, true},
		{"negatives", []int{-1, 0, 1, -1}, true},
		{"large values", []int{math.MaxInt, math.MinInt, math.MaxInt}, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := ContainsDuplicateBrute(tc.nums); got != tc.want {
				t.Errorf("brute = %v, want %v", got, tc.want)
			}
			if got := ContainsDuplicate(tc.nums); got != tc.want {
				t.Errorf("optimal = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestContainsDuplicateRandomized(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	for range 1000 {
		nums := make([]int, rng.Intn(20))
		for i := range nums {
			nums[i] = rng.Intn(40) - 20
		}
		if a, b := ContainsDuplicateBrute(nums), ContainsDuplicate(nums); a != b {
			t.Fatalf("%v: brute=%v optimal=%v", nums, a, b)
		}
	}
}
