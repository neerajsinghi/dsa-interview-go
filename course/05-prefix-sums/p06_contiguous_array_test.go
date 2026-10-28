package prefixsums

import (
	"math/rand"
	"testing"
)

func TestFindMaxLength(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want int
	}{
		{"empty", nil, 0},
		{"single", []int{0}, 0},
		{"pair", []int{0, 1}, 2},
		{"odd", []int{0, 1, 0}, 2},
		{"whole", []int{0, 1, 1, 1, 0, 0}, 6},
		{"all zeros", []int{0, 0, 0}, 0},
		{"inner", []int{1, 1, 0, 1, 1}, 2},
		{"long tail", []int{0, 0, 1, 0, 0, 0, 1, 1}, 6},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := FindMaxLengthBrute(tc.nums); got != tc.want {
				t.Errorf("brute = %d, want %d", got, tc.want)
			}
			if got := FindMaxLength(tc.nums); got != tc.want {
				t.Errorf("prefix+map = %d, want %d", got, tc.want)
			}
		})
	}
}

func TestFindMaxLengthRandomized(t *testing.T) {
	rng := rand.New(rand.NewSource(6))
	for range 1000 {
		nums := make([]int, rng.Intn(20))
		for i := range nums {
			nums[i] = rng.Intn(2)
		}
		if a, b := FindMaxLengthBrute(nums), FindMaxLength(nums); a != b {
			t.Fatalf("%v: brute=%d prefix=%d", nums, a, b)
		}
	}
}
