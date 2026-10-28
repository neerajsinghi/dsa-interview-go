package twopointers

import (
	"math/rand"
	"reflect"
	"slices"
	"testing"
)

func TestThreeSum(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want [][]int
	}{
		{"empty", nil, [][]int{}},
		{"too short", []int{0, 0}, [][]int{}},
		{"classic", []int{-1, 0, 1, 2, -1, -4}, [][]int{{-1, -1, 2}, {-1, 0, 1}}},
		{"no answer", []int{0, 1, 1}, [][]int{}},
		{"zeros", []int{0, 0, 0, 0}, [][]int{{0, 0, 0}}},
		{"many dups", []int{-2, 0, 0, 2, 2}, [][]int{{-2, 0, 2}}},
		{"several", []int{-4, -2, -2, -2, 0, 1, 2, 2, 2, 3, 3, 4, 4, 6, 6},
			[][]int{{-4, -2, 6}, {-4, 0, 4}, {-4, 1, 3}, {-4, 2, 2}, {-2, -2, 4}, {-2, 0, 2}}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := ThreeSumBrute(tc.nums); !reflect.DeepEqual(got, tc.want) {
				t.Errorf("brute = %v, want %v", got, tc.want)
			}
			if got := ThreeSum(tc.nums); !reflect.DeepEqual(got, tc.want) {
				t.Errorf("two pointers = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestThreeSumDoesNotMutate(t *testing.T) {
	nums := []int{3, -1, -2}
	ThreeSum(nums)
	if !slices.Equal(nums, []int{3, -1, -2}) {
		t.Fatalf("input mutated: %v", nums)
	}
}

func TestThreeSumRandomized(t *testing.T) {
	rng := rand.New(rand.NewSource(5))
	for range 500 {
		nums := make([]int, rng.Intn(14))
		for i := range nums {
			nums[i] = rng.Intn(11) - 5
		}
		if a, b := ThreeSumBrute(nums), ThreeSum(nums); !reflect.DeepEqual(a, b) {
			t.Fatalf("%v: brute=%v two pointers=%v", nums, a, b)
		}
	}
}
