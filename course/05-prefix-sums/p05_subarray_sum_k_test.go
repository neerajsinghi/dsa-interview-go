package prefixsums

import (
	"math/rand"
	"testing"
)

func TestSubarraySum(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		k    int
		want int
	}{
		{"empty", nil, 0, 0}, // no non-empty subarrays
		{"single hit", []int{3}, 3, 1},
		{"single miss", []int{3}, 2, 0},
		{"classic", []int{1, 1, 1}, 2, 2},
		{"three", []int{1, 2, 3}, 3, 2},
		{"negatives and zero", []int{1, -1, 0}, 0, 3},
		{"all zeros", []int{0, 0, 0}, 0, 6}, // every subarray
		{"negative k", []int{-1, -1, 1}, -1, 3},
		{"large", []int{1 << 50, -(1 << 50), 1 << 50}, 1 << 50, 3},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := SubarraySumBrute(tc.nums, tc.k); got != tc.want {
				t.Errorf("brute = %d, want %d", got, tc.want)
			}
			if got := SubarraySum(tc.nums, tc.k); got != tc.want {
				t.Errorf("prefix+map = %d, want %d", got, tc.want)
			}
		})
	}
}

func TestSubarraySumRandomized(t *testing.T) {
	rng := rand.New(rand.NewSource(5))
	for range 1000 {
		nums := make([]int, rng.Intn(15))
		for i := range nums {
			nums[i] = rng.Intn(7) - 3
		}
		k := rng.Intn(9) - 4
		if a, b := SubarraySumBrute(nums, k), SubarraySum(nums, k); a != b {
			t.Fatalf("%v k=%d: brute=%d prefix=%d", nums, k, a, b)
		}
	}
}
