package dp2d

import (
	"math/rand"
	"testing"
)

func TestBurstBalloons(t *testing.T) {
	tests := []struct {
		nums []int
		want int
	}{
		{[]int{3, 1, 5, 8}, 167},
		{[]int{1, 5}, 10},
		{[]int{7}, 7},
		{nil, 0},
		{[]int{0, 0}, 0},
		{[]int{2, 3}, 9},
	}
	impls := map[string]func([]int) int{"brute": BurstBalloonsBrute, "memo": BurstBalloonsMemo, "tab": BurstBalloonsTab}
	for _, tt := range tests {
		for name, fn := range impls {
			if got := fn(tt.nums); got != tt.want {
				t.Errorf("%s(%v) = %d, want %d", name, tt.nums, got, tt.want)
			}
		}
	}
}

func TestBurstBalloonsRandom(t *testing.T) {
	r := rand.New(rand.NewSource(9))
	for iter := 0; iter < 100; iter++ {
		nums := make([]int, r.Intn(7))
		for i := range nums {
			nums[i] = r.Intn(10)
		}
		want := BurstBalloonsBrute(nums)
		if a, b := BurstBalloonsMemo(nums), BurstBalloonsTab(nums); a != want || b != want {
			t.Fatalf("%v: brute=%d memo=%d tab=%d", nums, want, a, b)
		}
	}
}
