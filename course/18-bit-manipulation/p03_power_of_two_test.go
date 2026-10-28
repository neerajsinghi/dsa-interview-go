package bits

import (
	"math"
	"testing"
)

func TestIsPowerOfTwo(t *testing.T) {
	tests := []struct {
		n    int
		want bool
	}{
		{1, true}, {2, true}, {16, true}, {1 << 40, true}, {1 << 62, true},
		{0, false}, {6, false}, {-8, false}, {math.MinInt, false}, {math.MaxInt, false},
	}
	for _, tt := range tests {
		if got := IsPowerOfTwo(tt.n); got != tt.want {
			t.Errorf("IsPowerOfTwo(%d) = %v, want %v", tt.n, got, tt.want)
		}
		if got := IsPowerOfTwoLoop(tt.n); got != tt.want {
			t.Errorf("IsPowerOfTwoLoop(%d) = %v, want %v", tt.n, got, tt.want)
		}
	}
	for n := -50; n <= 5000; n++ {
		if IsPowerOfTwo(n) != IsPowerOfTwoLoop(n) {
			t.Fatalf("disagree at %d", n)
		}
	}
}

func TestIsPowerOfFour(t *testing.T) {
	want := map[int]bool{}
	for p := 1; p <= 1<<20; p *= 4 {
		want[p] = true
	}
	for n := -5; n <= 1<<20; n++ {
		if got := IsPowerOfFour(n); got != want[n] {
			t.Fatalf("IsPowerOfFour(%d) = %v", n, got)
		}
	}
}
