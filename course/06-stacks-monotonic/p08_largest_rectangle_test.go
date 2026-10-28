package stacks

import (
	"math/rand"
	"testing"
)

func TestLargestRectangle(t *testing.T) {
	tests := []struct {
		in   []int
		want int
	}{
		{[]int{2, 1, 5, 6, 2, 3}, 10},
		{[]int{2, 4}, 4},
		{[]int{}, 0},
		{[]int{0}, 0},
		{[]int{7}, 7},
		{[]int{3, 3, 3, 3}, 12},
		{[]int{1, 2, 3, 4, 5}, 9},
		{[]int{5, 4, 3, 2, 1}, 9},
		{[]int{2, 0, 2}, 2},
		{[]int{6, 2, 5, 4, 5, 1, 6}, 12},
	}
	for _, tt := range tests {
		for name, fn := range map[string]func([]int) int{
			"brute": LargestRectangleBrute, "stack": LargestRectangle,
		} {
			if got := fn(tt.in); got != tt.want {
				t.Errorf("%s(%v) = %d, want %d", name, tt.in, got, tt.want)
			}
		}
	}
}

func TestLargestRectangleRandom(t *testing.T) {
	r := rand.New(rand.NewSource(8))
	for iter := 0; iter < 500; iter++ {
		h := make([]int, r.Intn(30))
		for i := range h {
			h[i] = r.Intn(10)
		}
		if got, want := LargestRectangle(h), LargestRectangleBrute(h); got != want {
			t.Fatalf("LargestRectangle(%v) = %d, want %d", h, got, want)
		}
	}
}
