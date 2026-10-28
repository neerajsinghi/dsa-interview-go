package binarysearch

import (
	"math/rand"
	"testing"
)

func TestMinEatingSpeed(t *testing.T) {
	tests := []struct {
		piles []int
		h     int
		want  int
	}{
		{[]int{3, 6, 7, 11}, 8, 4},
		{[]int{30, 11, 23, 4, 20}, 5, 30},
		{[]int{30, 11, 23, 4, 20}, 6, 23},
		{[]int{1}, 1, 1},
		{[]int{1000000000}, 2, 500000000},
		{[]int{5, 5, 5}, 100, 1},
	}
	for _, tt := range tests {
		if got := MinEatingSpeed(tt.piles, tt.h); got != tt.want {
			t.Errorf("MinEatingSpeed(%v, %d) = %d, want %d", tt.piles, tt.h, got, tt.want)
		}
	}
}

func TestMinEatingSpeedRandom(t *testing.T) {
	r := rand.New(rand.NewSource(7))
	for iter := 0; iter < 500; iter++ {
		piles := make([]int, 1+r.Intn(8))
		for i := range piles {
			piles[i] = 1 + r.Intn(50)
		}
		h := len(piles) + r.Intn(30)
		if got, want := MinEatingSpeed(piles, h), MinEatingSpeedBrute(piles, h); got != want {
			t.Fatalf("MinEatingSpeed(%v, %d) = %d, want %d", piles, h, got, want)
		}
	}
}
