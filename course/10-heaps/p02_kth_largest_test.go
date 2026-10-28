package heaps

import (
	"math/rand"
	"testing"
)

func TestFindKthLargest(t *testing.T) {
	tests := []struct {
		nums []int
		k    int
		want int
	}{
		{[]int{3, 2, 1, 5, 6, 4}, 2, 5},
		{[]int{3, 2, 3, 1, 2, 4, 5, 5, 6}, 4, 4},
		{[]int{1}, 1, 1},
		{[]int{3, 3, 1}, 2, 3},
		{[]int{7, 7, 7, 7}, 3, 7},
		{[]int{-1, -5, 2}, 3, -5},
	}
	impls := map[string]func([]int, int) int{
		"sort": FindKthLargestSort, "heap": FindKthLargestHeap, "quickselect": FindKthLargestQuickselect,
	}
	for _, tt := range tests {
		for name, fn := range impls {
			if got := fn(tt.nums, tt.k); got != tt.want {
				t.Errorf("%s(%v, %d) = %d, want %d", name, tt.nums, tt.k, got, tt.want)
			}
		}
	}
}

func TestFindKthLargestRandom(t *testing.T) {
	r := rand.New(rand.NewSource(2))
	for iter := 0; iter < 500; iter++ {
		nums := randomInts(r, 1+r.Intn(40), 15)
		k := 1 + r.Intn(len(nums))
		want := FindKthLargestSort(nums, k)
		if got := FindKthLargestHeap(nums, k); got != want {
			t.Fatalf("heap(%v, %d) = %d, want %d", nums, k, got, want)
		}
		if got := FindKthLargestQuickselect(nums, k); got != want {
			t.Fatalf("quickselect(%v, %d) = %d, want %d", nums, k, got, want)
		}
	}
}
