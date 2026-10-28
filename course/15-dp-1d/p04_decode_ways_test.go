package dp1d

import (
	"math/rand"
	"testing"
)

var decodeImpls = map[string]func(string) int{
	"recursive": NumDecodingsRecursive, "memo": NumDecodingsMemo,
	"tab": NumDecodingsTab, "optimized": NumDecodings,
}

func TestNumDecodings(t *testing.T) {
	tests := []struct {
		s    string
		want int
	}{
		{"12", 2}, {"226", 3}, {"06", 0}, {"10", 1}, {"0", 0},
		{"27", 1}, {"100", 0}, {"101", 1}, {"11106", 2}, {"1", 1},
		{"", 1}, // the empty string decodes one way: to the empty string
		{"111111", 13},
	}
	for _, tt := range tests {
		for name, fn := range decodeImpls {
			if got := fn(tt.s); got != tt.want {
				t.Errorf("%s(%q) = %d, want %d", name, tt.s, got, tt.want)
			}
		}
	}
}

func TestNumDecodingsRandom(t *testing.T) {
	r := rand.New(rand.NewSource(73))
	for iter := 0; iter < 300; iter++ {
		b := make([]byte, r.Intn(15))
		for i := range b {
			b[i] = "0122236"[r.Intn(7)] // bias toward digits that form 10..26
		}
		s := string(b)
		want := NumDecodingsRecursive(s)
		for name, fn := range decodeImpls {
			if got := fn(s); got != want {
				t.Fatalf("%q: %s=%d brute=%d", s, name, got, want)
			}
		}
	}
}
