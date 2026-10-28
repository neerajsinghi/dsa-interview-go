package binarysearch

import (
	"math/rand"
	"slices"
	"testing"
)

func TestSearchRange(t *testing.T) {
	tests := []struct {
		nums   []int
		target int
		want   [2]int
	}{
		{[]int{5, 7, 7, 8, 8, 10}, 8, [2]int{3, 4}},
		{[]int{5, 7, 7, 8, 8, 10}, 6, [2]int{-1, -1}},
		{[]int{5, 7, 7, 8, 8, 10}, 5, [2]int{0, 0}},
		{[]int{5, 7, 7, 8, 8, 10}, 10, [2]int{5, 5}},
		{[]int{2, 2, 2, 2}, 2, [2]int{0, 3}},
		{[]int{1}, 1, [2]int{0, 0}},
		{nil, 0, [2]int{-1, -1}},
	}
	for _, tt := range tests {
		for name, fn := range map[string]func([]int, int) [2]int{
			"brute": SearchRangeBrute, "bounds": SearchRange,
		} {
			if got := fn(tt.nums, tt.target); got != tt.want {
				t.Errorf("%s(%v, %d) = %v, want %v", name, tt.nums, tt.target, got, tt.want)
			}
		}
	}
}

func TestSearchRangeRandom(t *testing.T) {
	r := rand.New(rand.NewSource(3))
	for iter := 0; iter < 1000; iter++ {
		nums := make([]int, r.Intn(25))
		for i := range nums {
			nums[i] = r.Intn(6)
		}
		slices.Sort(nums)
		target := r.Intn(8) - 1
		if got, want := SearchRange(nums, target), SearchRangeBrute(nums, target); got != want {
			t.Fatalf("SearchRange(%v, %d) = %v, want %v", nums, target, got, want)
		}
	}
}
