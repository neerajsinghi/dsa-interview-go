package bits

import (
	"math/rand"
	"testing"
)

func TestMissingNumber(t *testing.T) {
	tests := []struct {
		nums []int
		want int
	}{
		{[]int{3, 0, 1}, 2},
		{[]int{0, 1}, 2},
		{[]int{9, 6, 4, 2, 3, 5, 7, 0, 1}, 8},
		{[]int{0}, 1},
		{[]int{1}, 0},
		{[]int{}, 0},
	}
	impls := map[string]func([]int) int{"seen": MissingNumberSeen, "sum": MissingNumberSum, "xor": MissingNumber}
	for _, tt := range tests {
		for name, fn := range impls {
			if got := fn(tt.nums); got != tt.want {
				t.Errorf("%s(%v) = %d, want %d", name, tt.nums, got, tt.want)
			}
		}
	}
}

func TestMissingNumberRandom(t *testing.T) {
	r := rand.New(rand.NewSource(4))
	for iter := 0; iter < 200; iter++ {
		n := r.Intn(50)
		want := r.Intn(n + 1)
		var nums []int
		for v := 0; v <= n; v++ {
			if v != want {
				nums = append(nums, v)
			}
		}
		r.Shuffle(len(nums), func(i, j int) { nums[i], nums[j] = nums[j], nums[i] })
		if a, b, c := MissingNumber(nums), MissingNumberSum(nums), MissingNumberSeen(nums); a != want || b != want || c != want {
			t.Fatalf("%v: want %d, xor=%d sum=%d seen=%d", nums, want, a, b, c)
		}
	}
}
