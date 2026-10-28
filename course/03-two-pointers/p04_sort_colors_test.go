package twopointers

import (
	"math/rand"
	"slices"
	"testing"
)

func TestSortColors(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want []int
	}{
		{"empty", []int{}, []int{}},
		{"single", []int{1}, []int{1}},
		{"classic", []int{2, 0, 2, 1, 1, 0}, []int{0, 0, 1, 1, 2, 2}},
		{"three", []int{2, 0, 1}, []int{0, 1, 2}},
		{"all twos", []int{2, 2, 2}, []int{2, 2, 2}},
		{"reverse", []int{2, 2, 1, 1, 0, 0}, []int{0, 0, 1, 1, 2, 2}},
		{"no ones", []int{2, 0, 2, 0}, []int{0, 0, 2, 2}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			for name, fn := range map[string]func([]int){
				"Counting": SortColorsCounting, "DutchFlag": SortColors,
			} {
				got := slices.Clone(tc.nums)
				fn(got)
				if !slices.Equal(got, tc.want) {
					t.Errorf("%s(%v) = %v, want %v", name, tc.nums, got, tc.want)
				}
			}
		})
	}
}

func TestSortColorsRandomized(t *testing.T) {
	rng := rand.New(rand.NewSource(4))
	for range 1000 {
		nums := make([]int, rng.Intn(20))
		for i := range nums {
			nums[i] = rng.Intn(3)
		}
		want := slices.Clone(nums)
		slices.Sort(want)
		a, b := slices.Clone(nums), slices.Clone(nums)
		SortColorsCounting(a)
		SortColors(b)
		if !slices.Equal(a, want) || !slices.Equal(b, want) {
			t.Fatalf("%v: counting=%v flag=%v", nums, a, b)
		}
	}
}
