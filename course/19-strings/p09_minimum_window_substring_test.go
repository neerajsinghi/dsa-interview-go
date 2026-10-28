package stringalgos

import (
	"math/rand"
	"testing"
)

func TestMinWindow(t *testing.T) {
	tests := []struct{ s, t, want string }{
		{"ADOBECODEBANC", "ABC", "BANC"},
		{"a", "a", "a"},
		{"a", "aa", ""},
		{"", "a", ""},
		{"abc", "", ""},
		{"aaflslflsldkalskaaa", "aaa", "aaa"},
		{"bba", "ab", "ba"},
		{"abcabc", "cba", "abc"}, // leftmost of equal-length windows
	}
	for _, tt := range tests {
		if got := MinWindow(tt.s, tt.t); got != tt.want {
			t.Errorf("MinWindow(%q,%q) = %q, want %q", tt.s, tt.t, got, tt.want)
		}
		if got := MinWindowBrute(tt.s, tt.t); got != tt.want {
			t.Errorf("MinWindowBrute(%q,%q) = %q, want %q", tt.s, tt.t, got, tt.want)
		}
	}
}

func TestMinWindowRandom(t *testing.T) {
	r := rand.New(rand.NewSource(9))
	for iter := 0; iter < 1000; iter++ {
		s, tt := randText(r, r.Intn(20), "abc"), randText(r, r.Intn(4), "abc")
		if a, b := MinWindow(s, tt), MinWindowBrute(s, tt); a != b {
			t.Fatalf("(%q,%q): window=%q brute=%q", s, tt, a, b)
		}
	}
}
