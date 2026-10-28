package binarysearch

import (
	"math/rand"
	"slices"
	"testing"
)

func TestSearch(t *testing.T) {
	tests := []struct {
		nums   []int
		target int
		want   int
	}{
		{[]int{-1, 0, 3, 5, 9, 12}, 9, 4},
		{[]int{-1, 0, 3, 5, 9, 12}, 2, -1},
		{[]int{-1, 0, 3, 5, 9, 12}, -1, 0},
		{[]int{-1, 0, 3, 5, 9, 12}, 12, 5},
		{[]int{-1, 0, 3, 5, 9, 12}, 13, -1},
		{[]int{-1, 0, 3, 5, 9, 12}, -5, -1},
		{[]int{5}, 5, 0},
		{[]int{5}, 4, -1},
		{nil, 1, -1},
	}
	for _, tt := range tests {
		for name, fn := range map[string]func([]int, int) int{
			"closed": Search, "lowerBound": SearchLowerBound,
		} {
			if got := fn(tt.nums, tt.target); got != tt.want {
				t.Errorf("%s(%v, %d) = %d, want %d", name, tt.nums, tt.target, got, tt.want)
			}
		}
	}
}

func TestSearchRandom(t *testing.T) {
	r := rand.New(rand.NewSource(1))
	for iter := 0; iter < 500; iter++ {
		// distinct sorted values
		nums := r.Perm(40)[:r.Intn(40)]
		slices.Sort(nums)
		target := r.Intn(42) - 1
		want := slices.Index(nums, target) // linear brute force
		if got := Search(nums, target); got != want {
			t.Fatalf("Search(%v, %d) = %d, want %d", nums, target, got, want)
		}
		if got := SearchLowerBound(nums, target); got != want {
			t.Fatalf("SearchLowerBound(%v, %d) = %d, want %d", nums, target, got, want)
		}
	}
}
