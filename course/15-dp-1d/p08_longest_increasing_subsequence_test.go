package dp1d

import (
	"math/rand"
	"testing"
)

var lisImpls = map[string]func([]int) int{
	"recursive": LengthOfLISRecursive, "memo": LengthOfLISMemo,
	"quadratic": LengthOfLISQuadratic, "patience": LengthOfLIS,
}

func TestLengthOfLIS(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want int
	}{
		{"classic", []int{10, 9, 2, 5, 3, 7, 101, 18}, 4},
		{"with repeats", []int{0, 1, 0, 3, 2, 3}, 4},
		{"all equal", []int{7, 7, 7, 7}, 1},
		{"empty", nil, 0},
		{"single", []int{5}, 1},
		{"sorted", []int{1, 2, 3, 4, 5}, 5},
		{"reversed", []int{5, 4, 3, 2, 1}, 1},
		{"negatives", []int{-2, -1, -3, 0}, 3},
	}
	for _, tt := range tests {
		for name, fn := range lisImpls {
			if got := fn(tt.nums); got != tt.want {
				t.Errorf("%s/%s = %d, want %d", tt.name, name, got, tt.want)
			}
		}
	}
}

func TestLengthOfLISRandom(t *testing.T) {
	r := rand.New(rand.NewSource(77))
	for iter := 0; iter < 300; iter++ {
		nums := make([]int, r.Intn(13))
		for i := range nums {
			nums[i] = r.Intn(8) // small range => plenty of duplicates
		}
		want := LengthOfLISRecursive(nums)
		for name, fn := range lisImpls {
			if got := fn(nums); got != want {
				t.Fatalf("%v: %s=%d brute=%d", nums, name, got, want)
			}
		}
	}
	// Larger inputs: only the polynomial versions, checked against each other.
	for iter := 0; iter < 50; iter++ {
		nums := make([]int, 200+r.Intn(200))
		for i := range nums {
			nums[i] = r.Intn(1000)
		}
		if a, b := LengthOfLISQuadratic(nums), LengthOfLIS(nums); a != b {
			t.Fatalf("n=%d: quadratic=%d patience=%d", len(nums), a, b)
		}
	}
}
