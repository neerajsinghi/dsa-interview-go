package bits

import (
	"math"
	"math/rand"
	"testing"
)

func TestSingleNumberII(t *testing.T) {
	tests := []struct {
		nums []int
		want int
	}{
		{[]int{2, 2, 3, 2}, 3},
		{[]int{0, 1, 0, 1, 0, 1, 99}, 99},
		{[]int{-2, -2, 1, 1, 4, 1, 4, 4, -4, -2}, -4},
		{[]int{5}, 5},
		{[]int{math.MinInt, 7, 7, 7}, math.MinInt},
	}
	for _, tt := range tests {
		if got := SingleNumberII(tt.nums); got != tt.want {
			t.Errorf("SingleNumberII(%v) = %d, want %d", tt.nums, got, tt.want)
		}
		if got := SingleNumberIICount(tt.nums); got != tt.want {
			t.Errorf("SingleNumberIICount(%v) = %d, want %d", tt.nums, got, tt.want)
		}
	}
}

func TestSingleNumberIIRandom(t *testing.T) {
	r := rand.New(rand.NewSource(8))
	for iter := 0; iter < 200; iter++ {
		want := r.Int() - math.MaxInt/2
		nums := []int{want}
		for k := r.Intn(8); k > 0; k-- {
			v := r.Int() - math.MaxInt/2
			if v != want {
				nums = append(nums, v, v, v)
			}
		}
		r.Shuffle(len(nums), func(i, j int) { nums[i], nums[j] = nums[j], nums[i] })
		if a, b := SingleNumberII(nums), SingleNumberIICount(nums); a != want || b != want {
			t.Fatalf("want %d, machine=%d count=%d", want, a, b)
		}
	}
}
