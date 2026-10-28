package slidingwindow

import (
	"math/rand"
	"testing"
)

func TestLengthOfLongestSubstring(t *testing.T) {
	tests := []struct {
		s    string
		want int
	}{
		{"", 0},
		{"a", 1},
		{"bbbbb", 1},
		{"abcabcbb", 3},
		{"pwwkew", 3},
		{"abba", 2}, // left must never move backwards
		{"dvdf", 3},
		{" ", 1},
		{"a b!a", 4},
		{"abcdefghijklmnopqrstuvwxyz", 26},
	}
	for _, tc := range tests {
		if got := LengthOfLongestSubstringBrute(tc.s); got != tc.want {
			t.Errorf("brute(%q) = %d, want %d", tc.s, got, tc.want)
		}
		if got := LengthOfLongestSubstring(tc.s); got != tc.want {
			t.Errorf("window(%q) = %d, want %d", tc.s, got, tc.want)
		}
	}
}

func TestLengthOfLongestSubstringRandomized(t *testing.T) {
	rng := rand.New(rand.NewSource(4))
	for range 2000 {
		b := make([]byte, rng.Intn(15))
		for i := range b {
			b[i] = byte('a' + rng.Intn(4))
		}
		s := string(b)
		if x, y := LengthOfLongestSubstringBrute(s), LengthOfLongestSubstring(s); x != y {
			t.Fatalf("%q: brute=%d window=%d", s, x, y)
		}
	}
}
