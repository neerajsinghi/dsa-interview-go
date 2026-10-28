package binarysearch

import (
	"math/rand"
	"slices"
	"testing"
)

func TestSearchRotated(t *testing.T) {
	tests := []struct {
		nums   []int
		target int
		want   int
	}{
		{[]int{4, 5, 6, 7, 0, 1, 2}, 0, 4},
		{[]int{4, 5, 6, 7, 0, 1, 2}, 3, -1},
		{[]int{4, 5, 6, 7, 0, 1, 2}, 4, 0},
		{[]int{4, 5, 6, 7, 0, 1, 2}, 2, 6},
		{[]int{4, 5, 6, 7, 0, 1, 2}, 7, 3},
		{[]int{1}, 0, -1},
		{[]int{1}, 1, 0},
		{[]int{3, 1}, 1, 1},
		{[]int{1, 3}, 3, 1},
		{nil, 5, -1},
	}
	for _, tt := range tests {
		for name, fn := range map[string]func([]int, int) int{
			"brute": SearchRotatedBrute, "bs": SearchRotated,
		} {
			if got := fn(tt.nums, tt.target); got != tt.want {
				t.Errorf("%s(%v, %d) = %d, want %d", name, tt.nums, tt.target, got, tt.want)
			}
		}
	}
}

func TestSearchRotatedRandom(t *testing.T) {
	r := rand.New(rand.NewSource(6))
	for iter := 0; iter < 2000; iter++ {
		sorted := r.Perm(40)[:r.Intn(25)]
		slices.Sort(sorted)
		nums := rotate(sorted, r.Intn(len(sorted)+1))
		target := r.Intn(42) - 1
		if got, want := SearchRotated(nums, target), SearchRotatedBrute(nums, target); got != want {
			t.Fatalf("SearchRotated(%v, %d) = %d, want %d", nums, target, got, want)
		}
	}
}
