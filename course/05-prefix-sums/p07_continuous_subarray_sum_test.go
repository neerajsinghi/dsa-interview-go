package prefixsums

import (
	"math/rand"
	"testing"
)

func TestCheckSubarraySum(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		k    int
		want bool
	}{
		{"empty", nil, 3, false},
		{"single multiple", []int{6}, 6, false}, // length must be ≥ 2
		{"classic", []int{23, 2, 4, 6, 7}, 6, true},
		{"whole array", []int{23, 2, 6, 4, 7}, 6, true}, // 42
		{"none", []int{23, 2, 6, 4, 7}, 13, false},
		{"two zeros", []int{0, 0}, 5, true}, // 0 is a multiple of k
		{"zero then value", []int{5, 0}, 3, false},
		{"k one", []int{1, 2}, 1, true},
		{"negatives", []int{-7, 1}, 6, true}, // -6
		{"negative prefix trap", []int{-1, 2, 9}, 6, false},
		{"mixed signs", []int{4, -3, 5}, 6, true}, // 4 + -3 + 5 = 6
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := CheckSubarraySumBrute(tc.nums, tc.k); got != tc.want {
				t.Errorf("brute = %v, want %v", got, tc.want)
			}
			if got := CheckSubarraySum(tc.nums, tc.k); got != tc.want {
				t.Errorf("prefix mod = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestCheckSubarraySumRandomized(t *testing.T) {
	rng := rand.New(rand.NewSource(7))
	for range 2000 {
		nums := make([]int, rng.Intn(10))
		for i := range nums {
			nums[i] = rng.Intn(21) - 10
		}
		k := 1 + rng.Intn(12)
		if a, b := CheckSubarraySumBrute(nums, k), CheckSubarraySum(nums, k); a != b {
			t.Fatalf("%v k=%d: brute=%v prefix=%v", nums, k, a, b)
		}
	}
}
