package slidingwindow

import (
	"math/rand"
	"testing"
)

func TestMaxProfit(t *testing.T) {
	tests := []struct {
		name   string
		prices []int
		want   int
	}{
		{"empty", nil, 0},
		{"single", []int{5}, 0},
		{"classic", []int{7, 1, 5, 3, 6, 4}, 5},
		{"falling", []int{7, 6, 4, 3, 1}, 0},
		{"rising", []int{1, 2, 3, 4, 5}, 4},
		{"flat", []int{3, 3, 3}, 0},
		{"min after max", []int{2, 4, 1}, 2},
		{"later lower buy", []int{3, 8, 1, 9}, 8},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := MaxProfitBrute(tc.prices); got != tc.want {
				t.Errorf("brute = %d, want %d", got, tc.want)
			}
			if got := MaxProfit(tc.prices); got != tc.want {
				t.Errorf("window = %d, want %d", got, tc.want)
			}
		})
	}
}

func TestMaxProfitRandomized(t *testing.T) {
	rng := rand.New(rand.NewSource(2))
	for range 1000 {
		prices := make([]int, rng.Intn(15))
		for i := range prices {
			prices[i] = rng.Intn(20)
		}
		if a, b := MaxProfitBrute(prices), MaxProfit(prices); a != b {
			t.Fatalf("%v: brute=%d window=%d", prices, a, b)
		}
	}
}
