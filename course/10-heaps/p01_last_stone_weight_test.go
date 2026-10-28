package heaps

import (
	"math/rand"
	"testing"
)

func TestLastStoneWeight(t *testing.T) {
	tests := []struct {
		stones []int
		want   int
	}{
		{[]int{2, 7, 4, 1, 8, 1}, 1},
		{[]int{3, 3}, 0},
		{[]int{1}, 1},
		{nil, 0},
		{[]int{10, 4}, 6},
		{[]int{5, 5, 5}, 5},
	}
	for _, tt := range tests {
		for name, fn := range map[string]func([]int) int{"brute": LastStoneWeightBrute, "heap": LastStoneWeight} {
			if got := fn(tt.stones); got != tt.want {
				t.Errorf("%s(%v) = %d, want %d", name, tt.stones, got, tt.want)
			}
		}
	}
}

func TestLastStoneWeightRandom(t *testing.T) {
	r := rand.New(rand.NewSource(1))
	for iter := 0; iter < 300; iter++ {
		stones := randomInts(r, r.Intn(20), 30)
		for i := range stones {
			stones[i]++ // weights are positive
		}
		if got, want := LastStoneWeight(stones), LastStoneWeightBrute(stones); got != want {
			t.Fatalf("LastStoneWeight(%v) = %d, want %d", stones, got, want)
		}
	}
}
