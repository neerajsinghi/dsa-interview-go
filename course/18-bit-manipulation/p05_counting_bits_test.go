package bits

import (
	"slices"
	"testing"
)

func TestCountBits(t *testing.T) {
	tests := []struct {
		n    int
		want []int
	}{
		{0, []int{0}},
		{2, []int{0, 1, 1}},
		{5, []int{0, 1, 1, 2, 1, 2}},
		{8, []int{0, 1, 1, 2, 1, 2, 2, 3, 1}},
	}
	impls := map[string]func(int) []int{"naive": CountBitsNaive, "shift": CountBitsShift, "lowest": CountBits}
	for _, tt := range tests {
		for name, fn := range impls {
			if got := fn(tt.n); !slices.Equal(got, tt.want) {
				t.Errorf("%s(%d) = %v, want %v", name, tt.n, got, tt.want)
			}
		}
	}
	want := CountBitsNaive(5000)
	if !slices.Equal(CountBits(5000), want) || !slices.Equal(CountBitsShift(5000), want) {
		t.Error("disagreement on n=5000")
	}
}
