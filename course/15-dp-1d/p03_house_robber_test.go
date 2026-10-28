package dp1d

import (
	"math/rand"
	"testing"
)

var robImpls = map[string]func([]int) int{
	"recursive": RobRecursive, "memo": RobMemo, "tab": RobTab, "optimized": Rob,
}

func TestRob(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want int
	}{
		{"classic", []int{1, 2, 3, 1}, 4},
		{"every other", []int{2, 7, 9, 3, 1}, 12},
		{"skip two in a row", []int{2, 1, 1, 2}, 4},
		{"single", []int{5}, 5},
		{"empty", nil, 0},
		{"zeros", []int{0, 0, 0}, 0},
	}
	for _, tt := range tests {
		for name, fn := range robImpls {
			if got := fn(tt.nums); got != tt.want {
				t.Errorf("%s/%s = %d, want %d", tt.name, name, got, tt.want)
			}
		}
	}
}

func TestRobCircular(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want int
	}{
		{"three in a circle", []int{2, 3, 2}, 3},
		{"four", []int{1, 2, 3, 1}, 4},
		{"three increasing", []int{1, 2, 3}, 3},
		{"single", []int{7}, 7},
		{"pair", []int{4, 9}, 9},
		{"empty", nil, 0},
	}
	for _, tt := range tests {
		for name, fn := range map[string]func([]int) int{"dp": RobCircular, "brute": RobCircularBrute} {
			if got := fn(tt.nums); got != tt.want {
				t.Errorf("%s/%s = %d, want %d", tt.name, name, got, tt.want)
			}
		}
	}
}

func TestRobRandom(t *testing.T) {
	r := rand.New(rand.NewSource(72))
	for iter := 0; iter < 300; iter++ {
		nums := make([]int, r.Intn(13))
		for i := range nums {
			nums[i] = r.Intn(20)
		}
		want := RobRecursive(nums)
		for name, fn := range robImpls {
			if got := fn(nums); got != want {
				t.Fatalf("%v: %s=%d brute=%d", nums, name, got, want)
			}
		}
		if a, b := RobCircular(nums), RobCircularBrute(nums); a != b {
			t.Fatalf("%v: circular dp=%d brute=%d", nums, a, b)
		}
	}
}
