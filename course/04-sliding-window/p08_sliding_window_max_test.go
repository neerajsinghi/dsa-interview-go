package slidingwindow

import (
	"math/rand"
	"slices"
	"testing"
)

func TestMaxSlidingWindow(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		k    int
		want []int
	}{
		{"empty", nil, 1, nil},
		{"k zero", []int{1}, 0, nil},
		{"k too big", []int{1, 2}, 3, nil},
		{"single", []int{4}, 1, []int{4}},
		{"classic", []int{1, 3, -1, -3, 5, 3, 6, 7}, 3, []int{3, 3, 5, 5, 6, 7}},
		{"k one", []int{9, 8, 7}, 1, []int{9, 8, 7}},
		{"k equals n", []int{2, 9, 4}, 3, []int{9}},
		{"duplicates", []int{5, 5, 5, 5}, 2, []int{5, 5, 5}},
		{"decreasing", []int{5, 4, 3, 2, 1}, 2, []int{5, 4, 3, 2}},
		{"negatives", []int{-7, -8, -1, -9}, 2, []int{-7, -1, -1}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := MaxSlidingWindowBrute(tc.nums, tc.k); !slices.Equal(got, tc.want) {
				t.Errorf("brute = %v, want %v", got, tc.want)
			}
			if got := MaxSlidingWindow(tc.nums, tc.k); !slices.Equal(got, tc.want) {
				t.Errorf("deque = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestMaxSlidingWindowRandomized(t *testing.T) {
	rng := rand.New(rand.NewSource(8))
	for range 1000 {
		nums := make([]int, rng.Intn(20))
		for i := range nums {
			nums[i] = rng.Intn(11) - 5
		}
		k := rng.Intn(len(nums)+2) - 1
		if a, b := MaxSlidingWindowBrute(nums, k), MaxSlidingWindow(nums, k); !slices.Equal(a, b) {
			t.Fatalf("%v k=%d: brute=%v deque=%v", nums, k, a, b)
		}
	}
}
