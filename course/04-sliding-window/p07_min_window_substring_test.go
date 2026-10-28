package slidingwindow

import (
	"math/rand"
	"testing"
)

func TestMinWindow(t *testing.T) {
	tests := []struct {
		s, t, want string
	}{
		{"", "", ""},
		{"abc", "", ""},
		{"", "a", ""},
		{"a", "a", "a"},
		{"a", "aa", ""},
		{"ADOBECODEBANC", "ABC", "BANC"},
		{"aa", "aa", "aa"},
		{"ab", "b", "b"},
		{"bba", "ab", "ba"},
		{"abcabc", "abc", "abc"}, // leftmost of equal-length windows
		{"aAbB", "AB", "AbB"},    // case-sensitive
	}
	for _, tc := range tests {
		if got := MinWindowBrute(tc.s, tc.t); got != tc.want {
			t.Errorf("brute(%q, %q) = %q, want %q", tc.s, tc.t, got, tc.want)
		}
		if got := MinWindow(tc.s, tc.t); got != tc.want {
			t.Errorf("window(%q, %q) = %q, want %q", tc.s, tc.t, got, tc.want)
		}
	}
}

func TestMinWindowRandomized(t *testing.T) {
	rng := rand.New(rand.NewSource(7))
	word := func(n int) string {
		b := make([]byte, n)
		for i := range b {
			b[i] = byte('a' + rng.Intn(3))
		}
		return string(b)
	}
	for range 1000 {
		s, tt := word(rng.Intn(12)), word(rng.Intn(4))
		if x, y := MinWindowBrute(s, tt), MinWindow(s, tt); x != y {
			t.Fatalf("s=%q t=%q: brute=%q window=%q", s, tt, x, y)
		}
	}
}
