package stacks

import (
	"math/rand"
	"slices"
	"testing"
)

func TestNextGreaterCircular(t *testing.T) {
	tests := []struct {
		in, want []int
	}{
		{[]int{1, 2, 1}, []int{2, -1, 2}},
		{[]int{1, 2, 3, 4, 3}, []int{2, 3, 4, -1, 4}},
		{[]int{5, 4, 3, 2, 1}, []int{-1, 5, 5, 5, 5}},
		{[]int{2, 2, 2}, []int{-1, -1, -1}},
		{[]int{7}, []int{-1}},
		{[]int{}, []int{}},
	}
	for _, tt := range tests {
		for name, fn := range map[string]func([]int) []int{
			"brute": NextGreaterCircularBrute, "stack": NextGreaterCircular,
		} {
			if got := fn(tt.in); !slices.Equal(got, tt.want) {
				t.Errorf("%s(%v) = %v, want %v", name, tt.in, got, tt.want)
			}
		}
	}
}

func TestNextGreaterCircularRandom(t *testing.T) {
	r := rand.New(rand.NewSource(5))
	for iter := 0; iter < 500; iter++ {
		nums := make([]int, r.Intn(25))
		for i := range nums {
			nums[i] = r.Intn(8) - 3
		}
		if got, want := NextGreaterCircular(nums), NextGreaterCircularBrute(nums); !slices.Equal(got, want) {
			t.Fatalf("NextGreaterCircular(%v) = %v, want %v", nums, got, want)
		}
	}
}
