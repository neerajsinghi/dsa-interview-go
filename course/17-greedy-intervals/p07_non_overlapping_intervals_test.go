package greedy

import (
	"math/rand"
	"testing"
)

func TestEraseOverlapIntervals(t *testing.T) {
	tests := []struct {
		in   []Interval
		want int
	}{
		{iv(1, 2, 2, 3, 3, 4, 1, 3), 1},
		{iv(1, 2, 1, 2, 1, 2), 2},
		{iv(1, 2, 2, 3), 0},
		{nil, 0},
		{iv(1, 100, 11, 22, 1, 11, 2, 12), 2},
		{iv(-5, -1, -3, 0), 1},
	}
	for _, tt := range tests {
		if got := EraseOverlapIntervals(tt.in); got != tt.want {
			t.Errorf("EraseOverlapIntervals(%v) = %d, want %d", tt.in, got, tt.want)
		}
		if got := EraseOverlapIntervalsBrute(tt.in); got != tt.want {
			t.Errorf("Brute(%v) = %d, want %d", tt.in, got, tt.want)
		}
	}
}

func TestEraseOverlapIntervalsRandom(t *testing.T) {
	r := rand.New(rand.NewSource(7))
	for iter := 0; iter < 200; iter++ {
		in := randIntervals(r, r.Intn(10), 20)
		if a, b := EraseOverlapIntervals(in), EraseOverlapIntervalsBrute(in); a != b {
			t.Fatalf("%v: greedy=%d brute=%d", in, a, b)
		}
	}
}
