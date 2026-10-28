package arrayshashing

import (
	"math/rand"
	"slices"
	"testing"
)

func TestTwoSum(t *testing.T) {
	tests := []struct {
		name   string
		nums   []int
		target int
		want   []int
	}{
		{"empty", nil, 5, nil},
		{"single", []int{5}, 10, nil}, // can't use the same element twice
		{"classic", []int{2, 7, 11, 15}, 9, []int{0, 1}},
		{"not at start", []int{3, 2, 4}, 6, []int{1, 2}},
		{"duplicates", []int{3, 3}, 6, []int{0, 1}},
		{"negatives", []int{-3, 4, 3, 90}, 0, []int{0, 2}},
		{"no answer", []int{1, 2, 3}, 100, nil},
		{"large", []int{1 << 40, 5, 1 << 40}, 1 << 41, []int{0, 2}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := TwoSumBrute(tc.nums, tc.target); !slices.Equal(got, tc.want) {
				t.Errorf("brute = %v, want %v", got, tc.want)
			}
			if got := TwoSum(tc.nums, tc.target); !slices.Equal(got, tc.want) {
				t.Errorf("optimal = %v, want %v", got, tc.want)
			}
		})
	}
}

// With several valid pairs the two versions may pick different ones, so we
// check that both agree on existence and that every answer is valid.
func TestTwoSumRandomized(t *testing.T) {
	rng := rand.New(rand.NewSource(3))
	valid := func(nums []int, target int, p []int) bool {
		return len(p) == 2 && p[0] < p[1] && nums[p[0]]+nums[p[1]] == target
	}
	for range 1000 {
		nums := make([]int, rng.Intn(15))
		for i := range nums {
			nums[i] = rng.Intn(21) - 10
		}
		target := rng.Intn(21) - 10
		b, o := TwoSumBrute(nums, target), TwoSum(nums, target)
		if (b == nil) != (o == nil) {
			t.Fatalf("%v target %d: brute=%v optimal=%v", nums, target, b, o)
		}
		if o != nil && !valid(nums, target, o) {
			t.Fatalf("%v target %d: invalid pair %v", nums, target, o)
		}
	}
}
