package prefixsums

import (
	"math/rand"
	"slices"
	"testing"
)

func TestLeftRightDifference(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want []int
	}{
		{"empty", []int{}, []int{}},
		{"single", []int{1}, []int{0}},
		{"classic", []int{10, 4, 8, 3}, []int{15, 1, 11, 22}},
		{"symmetric", []int{2, 5, 2}, []int{7, 0, 7}},
		{"negatives", []int{-3, 1}, []int{1, 3}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := LeftRightDifferenceBrute(tc.nums); !slices.Equal(got, tc.want) {
				t.Errorf("brute = %v, want %v", got, tc.want)
			}
			if got := LeftRightDifference(tc.nums); !slices.Equal(got, tc.want) {
				t.Errorf("prefix = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestLeftRightDifferenceRandomized(t *testing.T) {
	rng := rand.New(rand.NewSource(3))
	for range 500 {
		nums := make([]int, rng.Intn(12))
		for i := range nums {
			nums[i] = rng.Intn(41) - 20
		}
		if a, b := LeftRightDifferenceBrute(nums), LeftRightDifference(nums); !slices.Equal(a, b) {
			t.Fatalf("%v: brute=%v prefix=%v", nums, a, b)
		}
	}
}
