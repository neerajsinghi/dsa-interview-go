package greedy

import (
	"math/rand"
	"testing"
)

func TestJump(t *testing.T) {
	tests := []struct {
		nums []int
		want int
	}{
		{[]int{2, 3, 1, 1, 4}, 2},
		{[]int{2, 3, 0, 1, 4}, 2},
		{[]int{0}, 0},
		{[]int{1, 1, 1, 1}, 3},
		{[]int{5, 1, 1, 1}, 1},
	}
	for _, tt := range tests {
		if got := Jump(tt.nums); got != tt.want {
			t.Errorf("Jump(%v) = %d, want %d", tt.nums, got, tt.want)
		}
		if got := JumpDP(tt.nums); got != tt.want {
			t.Errorf("JumpDP(%v) = %d, want %d", tt.nums, got, tt.want)
		}
	}
}

func TestJumpRandom(t *testing.T) {
	r := rand.New(rand.NewSource(3))
	for iter := 0; iter < 500; iter++ {
		nums := make([]int, 1+r.Intn(15))
		for i := range nums {
			nums[i] = 1 + r.Intn(4) // ≥1 guarantees the end is reachable
		}
		if a, b := Jump(nums), JumpDP(nums); a != b {
			t.Fatalf("%v: greedy=%d dp=%d", nums, a, b)
		}
	}
}
