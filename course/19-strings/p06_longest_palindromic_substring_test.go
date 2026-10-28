package stringalgos

import (
	"math/rand"
	"testing"
)

func TestLongestPalindrome(t *testing.T) {
	tests := []struct{ s, want string }{
		{"babad", "bab"},
		{"cbbd", "bb"},
		{"a", "a"},
		{"", ""},
		{"ac", "a"},
		{"forgeeksskeegfor", "geeksskeeg"},
		{"aaaa", "aaaa"},
		{"abacdfgdcaba", "aba"},
	}
	impls := map[string]func(string) string{"brute": LongestPalindromeBrute, "expand": LongestPalindromeExpand, "manacher": LongestPalindromeManacher}
	for _, tt := range tests {
		for name, fn := range impls {
			if got := fn(tt.s); got != tt.want {
				t.Errorf("%s(%q) = %q, want %q", name, tt.s, got, tt.want)
			}
		}
	}
}

func TestLongestPalindromeRandom(t *testing.T) {
	r := rand.New(rand.NewSource(6))
	for iter := 0; iter < 1000; iter++ {
		s := randText(r, r.Intn(16), "abc")
		want := LongestPalindromeBrute(s)
		if a, b := LongestPalindromeExpand(s), LongestPalindromeManacher(s); a != want || b != want {
			t.Fatalf("%q: brute=%q expand=%q manacher=%q", s, want, a, b)
		}
	}
}
