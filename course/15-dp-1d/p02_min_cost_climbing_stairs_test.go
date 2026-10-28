package dp1d

import (
	"math/rand"
	"testing"
)

var minCostImpls = map[string]func([]int) int{
	"recursive": MinCostClimbingStairsRecursive, "memo": MinCostClimbingStairsMemo,
	"tab": MinCostClimbingStairsTab, "optimized": MinCostClimbingStairs,
}

func TestMinCostClimbingStairs(t *testing.T) {
	tests := []struct {
		name string
		cost []int
		want int
	}{
		{"three stairs", []int{10, 15, 20}, 15},
		{"skip the expensive ones", []int{1, 100, 1, 1, 1, 100, 1, 1, 100, 1}, 6},
		{"two stairs", []int{5, 3}, 3},
		{"one stair: start on the top", []int{7}, 0},
		{"empty", nil, 0},
		{"all zero", []int{0, 0, 0, 0}, 0},
	}
	for _, tt := range tests {
		for name, fn := range minCostImpls {
			if got := fn(tt.cost); got != tt.want {
				t.Errorf("%s/%s = %d, want %d", tt.name, name, got, tt.want)
			}
		}
	}
}

func TestMinCostClimbingStairsRandom(t *testing.T) {
	r := rand.New(rand.NewSource(71))
	for iter := 0; iter < 300; iter++ {
		cost := make([]int, r.Intn(16))
		for i := range cost {
			cost[i] = r.Intn(50)
		}
		want := MinCostClimbingStairsRecursive(cost)
		for name, fn := range minCostImpls {
			if got := fn(cost); got != want {
				t.Fatalf("%v: %s=%d brute=%d", cost, name, got, want)
			}
		}
	}
}
