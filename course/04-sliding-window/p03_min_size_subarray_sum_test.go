package slidingwindow

import (
	"math/rand"
	"testing"
)

func TestMinSubArrayLen(t *testing.T) {
	tests := []struct {
		name   string
		target int
		nums   []int
		want   int
	}{
		{"empty", 1, nil, 0},
		{"single hit", 3, []int{5}, 1},
		{"single miss", 6, []int{5}, 0},
		{"classic", 7, []int{2, 3, 1, 2, 4, 3}, 2},
		{"exact element", 4, []int{1, 4, 4}, 1},
		{"none", 11, []int{1, 1, 1, 1, 1, 1, 1, 1}, 0},
		{"whole array", 15, []int{1, 2, 3, 4, 5}, 5},
		{"large", 1 << 41, []int{1 << 40, 1, 1 << 40}, 3},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := MinSubArrayLenBrute(tc.target, tc.nums); got != tc.want {
				t.Errorf("brute = %d, want %d", got, tc.want)
			}
			if got := MinSubArrayLen(tc.target, tc.nums); got != tc.want {
				t.Errorf("window = %d, want %d", got, tc.want)
			}
		})
	}
}

func TestMinSubArrayLenRandomized(t *testing.T) {
	rng := rand.New(rand.NewSource(3))
	for range 1000 {
		nums := make([]int, rng.Intn(15))
		for i := range nums {
			nums[i] = 1 + rng.Intn(9) // positive, as the problem requires
		}
		target := 1 + rng.Intn(40)
		if a, b := MinSubArrayLenBrute(target, nums), MinSubArrayLen(target, nums); a != b {
			t.Fatalf("target=%d %v: brute=%d window=%d", target, nums, a, b)
		}
	}
}
