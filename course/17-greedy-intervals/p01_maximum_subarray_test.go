package greedy

import (
	"math/rand"
	"testing"
)

func TestMaxSubarray(t *testing.T) {
	tests := []struct {
		nums []int
		want int
	}{
		{[]int{-2, 1, -3, 4, -1, 2, 1, -5, 4}, 6},
		{[]int{1}, 1},
		{[]int{5, 4, -1, 7, 8}, 23},
		{[]int{-3, -1, -2}, -1},
		{[]int{0, 0}, 0},
	}
	for _, tt := range tests {
		if got := MaxSubarray(tt.nums); got != tt.want {
			t.Errorf("MaxSubarray(%v) = %d, want %d", tt.nums, got, tt.want)
		}
		if got := MaxSubarrayBrute(tt.nums); got != tt.want {
			t.Errorf("MaxSubarrayBrute(%v) = %d, want %d", tt.nums, got, tt.want)
		}
	}
}

func TestMaxSubarrayRandom(t *testing.T) {
	r := rand.New(rand.NewSource(17))
	for iter := 0; iter < 500; iter++ {
		nums := make([]int, 1+r.Intn(15))
		for i := range nums {
			nums[i] = r.Intn(21) - 10
		}
		want := MaxSubarrayBrute(nums)
		if got := MaxSubarray(nums); got != want {
			t.Fatalf("%v: kadane=%d brute=%d", nums, got, want)
		}
		sum, lo, hi := MaxSubarrayRange(nums)
		check := 0
		for _, x := range nums[lo : hi+1] {
			check += x
		}
		if sum != want || check != want {
			t.Fatalf("%v: range [%d,%d] sum=%d check=%d want=%d", nums, lo, hi, sum, check, want)
		}
	}
}
