package dp1d

import (
	"math/rand"
	"testing"
)

func TestMaxProduct(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want int
	}{
		{"classic", []int{2, 3, -2, 4}, 6},
		{"zero wins", []int{-2, 0, -1}, 0},
		{"negatives cancel", []int{-2, 3, -4}, 24},
		{"single negative", []int{-3}, -3},
		{"zeros split", []int{2, -5, -2, -4, 3}, 24},
		{"all negative odd", []int{-1, -2, -3}, 6},
		{"zero in middle", []int{3, 0, 4}, 4},
	}
	for _, tt := range tests {
		for name, fn := range map[string]func([]int) int{"brute": MaxProductBrute, "dp": MaxProduct} {
			if got := fn(tt.nums); got != tt.want {
				t.Errorf("%s/%s = %d, want %d", tt.name, name, got, tt.want)
			}
		}
	}
}

func TestMaxProductRandom(t *testing.T) {
	r := rand.New(rand.NewSource(76))
	for iter := 0; iter < 500; iter++ {
		nums := make([]int, 1+r.Intn(12))
		for i := range nums {
			nums[i] = r.Intn(9) - 4 // -4..4, including zero
		}
		if a, b := MaxProduct(nums), MaxProductBrute(nums); a != b {
			t.Fatalf("%v: dp=%d brute=%d", nums, a, b)
		}
	}
}
