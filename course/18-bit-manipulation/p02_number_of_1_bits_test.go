package bits

import (
	"math/rand"
	"testing"
)

func TestHammingWeight(t *testing.T) {
	tests := []struct {
		n    uint32
		want int
	}{
		{0, 0}, {11, 3}, {128, 1}, {0xFFFFFFFF, 32}, {0x80000000, 1}, {0xFFFFFFFD, 31},
	}
	for _, tt := range tests {
		for name, fn := range map[string]func(uint32) int{"shift": HammingWeightShift, "kernighan": HammingWeight, "stdlib": HammingWeightStdlib} {
			if got := fn(tt.n); got != tt.want {
				t.Errorf("%s(%#x) = %d, want %d", name, tt.n, got, tt.want)
			}
		}
	}
}

func TestHammingWeightRandom(t *testing.T) {
	r := rand.New(rand.NewSource(2))
	for iter := 0; iter < 1000; iter++ {
		n := r.Uint32()
		want := HammingWeightStdlib(n)
		if a, b := HammingWeight(n), HammingWeightShift(n); a != want || b != want {
			t.Fatalf("%#x: want %d, kernighan=%d shift=%d", n, want, a, b)
		}
	}
}
