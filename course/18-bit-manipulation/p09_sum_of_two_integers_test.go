package bits

import (
	"math"
	"math/rand"
	"testing"
)

func TestGetSum(t *testing.T) {
	tests := []struct{ a, b, want int32 }{
		{1, 2, 3},
		{2, 3, 5},
		{-2, 3, 1},
		{-5, -7, -12},
		{0, 0, 0},
		{math.MaxInt32, 1, math.MinInt32},
		{math.MinInt32, -1, math.MaxInt32},
	}
	for _, tt := range tests {
		if got := GetSum(tt.a, tt.b); got != tt.want {
			t.Errorf("GetSum(%d,%d) = %d, want %d", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestGetSumRandom(t *testing.T) {
	r := rand.New(rand.NewSource(9))
	for iter := 0; iter < 2000; iter++ {
		a, b := int32(r.Uint32()), int32(r.Uint32())
		if got := GetSum(a, b); got != a+b {
			t.Fatalf("GetSum(%d,%d) = %d, want %d", a, b, got, a+b)
		}
		if got := GetDifference(a, b); got != a-b {
			t.Fatalf("GetDifference(%d,%d) = %d, want %d", a, b, got, a-b)
		}
	}
}
