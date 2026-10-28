package arrayshashing

import (
	"math/rand"
	"testing"
)

func TestLongestConsecutive(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want int
	}{
		{"empty", nil, 0},
		{"single", []int{7}, 1},
		{"classic", []int{100, 4, 200, 1, 3, 2}, 4},
		{"with duplicates", []int{0, 3, 7, 2, 5, 8, 4, 6, 0, 1}, 9},
		{"all same", []int{5, 5, 5}, 1},
		{"negatives", []int{-1, -3, -2, 10}, 3},
		{"two runs", []int{1, 2, 3, 10, 11, 12, 13}, 4},
		{"large values", []int{1 << 50, 1<<50 + 1, 1<<50 - 1}, 3},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := LongestConsecutiveSort(tc.nums); got != tc.want {
				t.Errorf("sort = %d, want %d", got, tc.want)
			}
			if got := LongestConsecutive(tc.nums); got != tc.want {
				t.Errorf("set = %d, want %d", got, tc.want)
			}
		})
	}
}

func TestLongestConsecutiveRandomized(t *testing.T) {
	rng := rand.New(rand.NewSource(8))
	for range 1000 {
		nums := make([]int, rng.Intn(25))
		for i := range nums {
			nums[i] = rng.Intn(40) - 20
		}
		if a, b := LongestConsecutiveSort(nums), LongestConsecutive(nums); a != b {
			t.Fatalf("%v: sort=%d set=%d", nums, a, b)
		}
	}
}
