package twopointers

import (
	"math/rand"
	"slices"
	"testing"
)

func TestTwoSumSorted(t *testing.T) {
	tests := []struct {
		name   string
		nums   []int
		target int
		want   []int
	}{
		{"empty", nil, 1, nil},
		{"single", []int{1}, 2, nil},
		{"classic", []int{2, 7, 11, 15}, 9, []int{0, 1}},
		{"ends", []int{2, 3, 4}, 6, []int{0, 2}},
		{"negatives", []int{-1, 0}, -1, []int{0, 1}},
		{"duplicates", []int{1, 3, 3, 8}, 6, []int{1, 2}},
		{"none", []int{1, 2, 3}, 10, nil},
		{"large", []int{-(1 << 50), 0, 1 << 50}, 0, []int{0, 2}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := TwoSumSortedBrute(tc.nums, tc.target); !slices.Equal(got, tc.want) {
				t.Errorf("brute = %v, want %v", got, tc.want)
			}
			if got := TwoSumSorted(tc.nums, tc.target); !slices.Equal(got, tc.want) {
				t.Errorf("two pointers = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestTwoSumSortedRandomized(t *testing.T) {
	rng := rand.New(rand.NewSource(3))
	for range 1000 {
		nums := make([]int, rng.Intn(12))
		for i := range nums {
			nums[i] = rng.Intn(21) - 10
		}
		slices.Sort(nums)
		target := rng.Intn(21) - 10
		b, o := TwoSumSortedBrute(nums, target), TwoSumSorted(nums, target)
		if (b == nil) != (o == nil) {
			t.Fatalf("%v target %d: brute=%v two pointers=%v", nums, target, b, o)
		}
		if o != nil && (o[0] >= o[1] || nums[o[0]]+nums[o[1]] != target) {
			t.Fatalf("%v target %d: invalid answer %v", nums, target, o)
		}
	}
}
