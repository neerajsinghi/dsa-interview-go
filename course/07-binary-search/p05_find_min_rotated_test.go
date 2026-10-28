package binarysearch

import (
	"math/rand"
	"slices"
	"testing"
)

// rotate returns a copy of xs rotated left by k.
func rotate(xs []int, k int) []int {
	if len(xs) == 0 {
		return nil
	}
	k %= len(xs)
	return append(slices.Clone(xs[k:]), xs[:k]...)
}

func TestFindMin(t *testing.T) {
	tests := []struct {
		nums []int
		want int
	}{
		{[]int{3, 4, 5, 1, 2}, 1},
		{[]int{4, 5, 6, 7, 0, 1, 2}, 0},
		{[]int{11, 13, 15, 17}, 11},
		{[]int{2, 1}, 1},
		{[]int{1, 2}, 1},
		{[]int{9}, 9},
		{[]int{5, 1, 2, 3, 4}, 1},
	}
	for _, tt := range tests {
		for name, fn := range map[string]func([]int) int{"brute": FindMinBrute, "bs": FindMin} {
			if got := fn(tt.nums); got != tt.want {
				t.Errorf("%s(%v) = %d, want %d", name, tt.nums, got, tt.want)
			}
		}
	}
}

func TestFindMinRandom(t *testing.T) {
	r := rand.New(rand.NewSource(5))
	for iter := 0; iter < 1000; iter++ {
		sorted := r.Perm(50)[:1+r.Intn(30)]
		slices.Sort(sorted)
		nums := rotate(sorted, r.Intn(len(sorted)+1))
		if got, want := FindMin(nums), FindMinBrute(nums); got != want {
			t.Fatalf("FindMin(%v) = %d, want %d", nums, got, want)
		}
	}
}
