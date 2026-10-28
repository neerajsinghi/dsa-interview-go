package binarysearch

import (
	"math/rand"
	"slices"
	"testing"
)

func TestSearchInsert(t *testing.T) {
	tests := []struct {
		nums   []int
		target int
		want   int
	}{
		{[]int{1, 3, 5, 6}, 5, 2},
		{[]int{1, 3, 5, 6}, 2, 1},
		{[]int{1, 3, 5, 6}, 7, 4},
		{[]int{1, 3, 5, 6}, 0, 0},
		{[]int{1, 3, 5, 6}, 1, 0},
		{[]int{1}, 1, 0},
		{nil, 3, 0},
	}
	for _, tt := range tests {
		if got := SearchInsert(tt.nums, tt.target); got != tt.want {
			t.Errorf("SearchInsert(%v, %d) = %d, want %d", tt.nums, tt.target, got, tt.want)
		}
	}
}

func TestSearchInsertRandom(t *testing.T) {
	r := rand.New(rand.NewSource(2))
	for iter := 0; iter < 500; iter++ {
		nums := r.Perm(30)[:r.Intn(30)]
		slices.Sort(nums)
		target := r.Intn(32) - 1
		// brute force: count elements strictly smaller than target
		want := 0
		for _, v := range nums {
			if v < target {
				want++
			}
		}
		if got := SearchInsert(nums, target); got != want {
			t.Fatalf("SearchInsert(%v, %d) = %d, want %d", nums, target, got, want)
		}
	}
}
