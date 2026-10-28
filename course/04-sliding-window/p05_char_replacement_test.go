package slidingwindow

import (
	"math/rand"
	"testing"
)

func TestCharacterReplacement(t *testing.T) {
	tests := []struct {
		s    string
		k    int
		want int
	}{
		{"", 2, 0},
		{"A", 0, 1},
		{"ABAB", 2, 4},
		{"AABABBA", 1, 4},
		{"AAAA", 0, 4},
		{"ABCDE", 0, 1},
		{"ABCDE", 10, 5}, // k larger than the string
		{"ABBB", 0, 3},
	}
	for _, tc := range tests {
		if got := CharacterReplacementBrute(tc.s, tc.k); got != tc.want {
			t.Errorf("brute(%q, %d) = %d, want %d", tc.s, tc.k, got, tc.want)
		}
		if got := CharacterReplacement(tc.s, tc.k); got != tc.want {
			t.Errorf("window(%q, %d) = %d, want %d", tc.s, tc.k, got, tc.want)
		}
	}
}

func TestCharacterReplacementRandomized(t *testing.T) {
	rng := rand.New(rand.NewSource(5))
	for range 2000 {
		b := make([]byte, rng.Intn(15))
		for i := range b {
			b[i] = byte('A' + rng.Intn(3))
		}
		s, k := string(b), rng.Intn(4)
		if x, y := CharacterReplacementBrute(s, k), CharacterReplacement(s, k); x != y {
			t.Fatalf("%q k=%d: brute=%d window=%d", s, k, x, y)
		}
	}
}
