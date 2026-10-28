package greedy

import (
	"math/rand"
	"testing"
)

func TestCanJump(t *testing.T) {
	tests := []struct {
		nums []int
		want bool
	}{
		{[]int{2, 3, 1, 1, 4}, true},
		{[]int{3, 2, 1, 0, 4}, false},
		{[]int{0}, true},
		{[]int{0, 1}, false},
		{[]int{1, 0}, true},
		{[]int{2, 0, 0}, true},
	}
	for _, tt := range tests {
		if got := CanJump(tt.nums); got != tt.want {
			t.Errorf("CanJump(%v) = %v, want %v", tt.nums, got, tt.want)
		}
		if got := CanJumpDP(tt.nums); got != tt.want {
			t.Errorf("CanJumpDP(%v) = %v, want %v", tt.nums, got, tt.want)
		}
	}
}

func TestCanJumpRandom(t *testing.T) {
	r := rand.New(rand.NewSource(2))
	for iter := 0; iter < 500; iter++ {
		nums := make([]int, 1+r.Intn(12))
		for i := range nums {
			nums[i] = r.Intn(4)
		}
		if a, b := CanJump(nums), CanJumpDP(nums); a != b {
			t.Fatalf("%v: greedy=%v dp=%v", nums, a, b)
		}
	}
}
