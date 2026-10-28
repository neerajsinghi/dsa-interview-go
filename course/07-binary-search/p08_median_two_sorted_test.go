package binarysearch

import (
	"math/rand"
	"slices"
	"testing"
)

func TestFindMedianSortedArrays(t *testing.T) {
	tests := []struct {
		a, b []int
		want float64
	}{
		{[]int{1, 3}, []int{2}, 2},
		{[]int{1, 2}, []int{3, 4}, 2.5},
		{[]int{}, []int{1}, 1},
		{[]int{2}, []int{}, 2},
		{nil, nil, 0},
		{[]int{1, 1, 1}, []int{1, 1}, 1},
		{[]int{1, 2, 3, 4}, []int{10, 20, 30}, 4},
		{[]int{-5, -3}, []int{-4}, -4},
		{[]int{100}, []int{1, 2, 3, 4, 5}, 3.5},
	}
	for _, tt := range tests {
		for name, fn := range map[string]func([]int, []int) float64{
			"brute": FindMedianSortedArraysBrute, "partition": FindMedianSortedArrays,
		} {
			if got := fn(tt.a, tt.b); got != tt.want {
				t.Errorf("%s(%v, %v) = %v, want %v", name, tt.a, tt.b, got, tt.want)
			}
		}
	}
}

func TestFindMedianSortedArraysRandom(t *testing.T) {
	r := rand.New(rand.NewSource(8))
	gen := func() []int {
		xs := make([]int, r.Intn(12))
		for i := range xs {
			xs[i] = r.Intn(20) - 10
		}
		slices.Sort(xs)
		return xs
	}
	for iter := 0; iter < 2000; iter++ {
		a, b := gen(), gen()
		if got, want := FindMedianSortedArrays(a, b), FindMedianSortedArraysBrute(a, b); got != want {
			t.Fatalf("FindMedianSortedArrays(%v, %v) = %v, want %v", a, b, got, want)
		}
	}
}
