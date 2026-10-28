package arrayshashing

import (
	"math/rand"
	"slices"
	"testing"
)

func TestTopKFrequent(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		k    int
		want []int
	}{
		{"empty", nil, 0, []int{}},
		{"single", []int{4}, 1, []int{4}},
		{"classic", []int{1, 1, 1, 2, 2, 3}, 2, []int{1, 2}},
		{"k zero", []int{1, 2}, 0, []int{}},
		{"all", []int{3, 1, 2, 1}, 3, []int{1, 2, 3}}, // ties → smaller value first
		{"negatives", []int{-1, -1, -2, -2, -2, 5}, 2, []int{-2, -1}},
		{"tie at boundary", []int{7, 8, 9}, 2, []int{7, 8}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := TopKFrequentSort(tc.nums, tc.k); !slices.Equal(got, tc.want) {
				t.Errorf("sort = %v, want %v", got, tc.want)
			}
			if got := TopKFrequent(tc.nums, tc.k); !slices.Equal(got, tc.want) {
				t.Errorf("bucket = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestTopKFrequentRandomized(t *testing.T) {
	rng := rand.New(rand.NewSource(5))
	for range 500 {
		nums := make([]int, 1+rng.Intn(30))
		for i := range nums {
			nums[i] = rng.Intn(10) - 5
		}
		k := rng.Intn(len(Counter(nums)) + 1)
		a, b := TopKFrequentSort(nums, k), TopKFrequent(nums, k)
		if !slices.Equal(a, b) {
			t.Fatalf("%v k=%d: sort=%v bucket=%v", nums, k, a, b)
		}
	}
}
