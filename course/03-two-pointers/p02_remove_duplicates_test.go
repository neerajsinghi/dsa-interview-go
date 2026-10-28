package twopointers

import (
	"math/rand"
	"slices"
	"testing"
)

func TestRemoveDuplicates(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want []int
	}{
		{"empty", []int{}, []int{}},
		{"single", []int{1}, []int{1}},
		{"all same", []int{2, 2, 2}, []int{2}},
		{"classic short", []int{1, 1, 2}, []int{1, 2}},
		{"classic long", []int{0, 0, 1, 1, 1, 2, 2, 3, 3, 4}, []int{0, 1, 2, 3, 4}},
		{"negatives", []int{-3, -3, -1, 0, 0}, []int{-3, -1, 0}},
		{"already distinct", []int{1, 2, 3}, []int{1, 2, 3}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			for name, fn := range map[string]func([]int) int{
				"Brute": RemoveDuplicatesBrute, "InPlace": RemoveDuplicates,
			} {
				nums := slices.Clone(tc.nums)
				k := fn(nums)
				if !slices.Equal(nums[:k], tc.want) {
					t.Errorf("%s: nums[:%d] = %v, want %v", name, k, nums[:k], tc.want)
				}
			}
		})
	}
}

func TestRemoveDuplicatesRandomized(t *testing.T) {
	rng := rand.New(rand.NewSource(2))
	for range 1000 {
		nums := make([]int, rng.Intn(15))
		for i := range nums {
			nums[i] = rng.Intn(6)
		}
		slices.Sort(nums)
		a, b := slices.Clone(nums), slices.Clone(nums)
		ka, kb := RemoveDuplicatesBrute(a), RemoveDuplicates(b)
		if !slices.Equal(a[:ka], b[:kb]) || !slices.Equal(b[:kb], slices.Compact(slices.Clone(nums))) {
			t.Fatalf("%v: brute=%v inplace=%v", nums, a[:ka], b[:kb])
		}
	}
}
