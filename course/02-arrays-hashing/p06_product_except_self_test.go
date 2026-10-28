package arrayshashing

import (
	"math/rand"
	"slices"
	"testing"
)

func TestProductExceptSelf(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want []int
	}{
		{"empty", []int{}, []int{}},
		{"single", []int{5}, []int{1}}, // empty product is 1
		{"two", []int{3, 4}, []int{4, 3}},
		{"classic", []int{1, 2, 3, 4}, []int{24, 12, 8, 6}},
		{"one zero", []int{-1, 1, 0, -3, 3}, []int{0, 0, 9, 0, 0}},
		{"two zeros", []int{0, 2, 0}, []int{0, 0, 0}},
		{"negatives", []int{-2, -3, 4}, []int{-12, -8, 6}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := ProductExceptSelfBrute(tc.nums); !slices.Equal(got, tc.want) {
				t.Errorf("brute = %v, want %v", got, tc.want)
			}
			if got := ProductExceptSelf(tc.nums); !slices.Equal(got, tc.want) {
				t.Errorf("optimal = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestProductExceptSelfRandomized(t *testing.T) {
	rng := rand.New(rand.NewSource(6))
	for range 500 {
		nums := make([]int, rng.Intn(10))
		for i := range nums {
			nums[i] = rng.Intn(7) - 3 // small values: no overflow, frequent zeros
		}
		if a, b := ProductExceptSelfBrute(nums), ProductExceptSelf(nums); !slices.Equal(a, b) {
			t.Fatalf("%v: brute=%v optimal=%v", nums, a, b)
		}
	}
}
