package bits

import (
	"math/rand"
	"testing"
)

func TestSingleNumber(t *testing.T) {
	tests := []struct {
		nums []int
		want int
	}{
		{[]int{2, 2, 1}, 1},
		{[]int{4, 1, 2, 1, 2}, 4},
		{[]int{7}, 7},
		{[]int{-3, 5, 5}, -3},
		{[]int{0, 9, 9}, 0},
	}
	for _, tt := range tests {
		if got := SingleNumber(tt.nums); got != tt.want {
			t.Errorf("SingleNumber(%v) = %d, want %d", tt.nums, got, tt.want)
		}
		if got := SingleNumberMap(tt.nums); got != tt.want {
			t.Errorf("SingleNumberMap(%v) = %d, want %d", tt.nums, got, tt.want)
		}
	}
}

func TestSingleNumberRandom(t *testing.T) {
	r := rand.New(rand.NewSource(18))
	for iter := 0; iter < 200; iter++ {
		want := r.Intn(2001) - 1000
		nums := []int{want}
		for k := r.Intn(10); k > 0; k-- {
			v := r.Intn(2001) - 1000
			if v == want {
				continue
			}
			nums = append(nums, v, v)
		}
		r.Shuffle(len(nums), func(i, j int) { nums[i], nums[j] = nums[j], nums[i] })
		if got := SingleNumber(nums); got != want {
			t.Fatalf("%v: got %d, want %d", nums, got, want)
		}
	}
}
