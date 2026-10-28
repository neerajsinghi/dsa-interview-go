package bits

import (
	"math/rand"
	"testing"
)

func TestReverseBits(t *testing.T) {
	tests := []struct{ in, want uint32 }{
		{0b00000010100101000001111010011100, 0b00111001011110000010100101000000},
		{0b11111111111111111111111111111101, 0b10111111111111111111111111111111},
		{0, 0},
		{1, 0x80000000},
		{0xFFFFFFFF, 0xFFFFFFFF},
	}
	for _, tt := range tests {
		for name, fn := range map[string]func(uint32) uint32{"loop": ReverseBitsLoop, "masks": ReverseBits, "stdlib": ReverseBitsStdlib} {
			if got := fn(tt.in); got != tt.want {
				t.Errorf("%s(%032b) = %032b, want %032b", name, tt.in, got, tt.want)
			}
		}
	}
}

func TestReverseBitsRandom(t *testing.T) {
	r := rand.New(rand.NewSource(6))
	for iter := 0; iter < 1000; iter++ {
		n := r.Uint32()
		want := ReverseBitsStdlib(n)
		if a, b := ReverseBitsLoop(n), ReverseBits(n); a != want || b != want {
			t.Fatalf("%032b: want %032b loop=%032b masks=%032b", n, want, a, b)
		}
		if ReverseBits(ReverseBits(n)) != n {
			t.Fatalf("reverse is not an involution for %032b", n)
		}
	}
}
