package stringalgos

import (
	"math/rand"
	"testing"
)

func TestLongestPrefix(t *testing.T) {
	tests := []struct{ s, want string }{
		{"level", "l"},
		{"ababab", "abab"},
		{"abc", ""},
		{"a", ""},
		{"", ""},
		{"aaaa", "aaa"},
		{"abcabcab", "abcab"},
	}
	impls := map[string]func(string) string{"brute": LongestPrefixBrute, "kmp": LongestPrefixKMP, "hash": LongestPrefixHash}
	for _, tt := range tests {
		for name, fn := range impls {
			if got := fn(tt.s); got != tt.want {
				t.Errorf("%s(%q) = %q, want %q", name, tt.s, got, tt.want)
			}
		}
	}
}

func TestLongestPrefixRandom(t *testing.T) {
	r := rand.New(rand.NewSource(8))
	for iter := 0; iter < 1000; iter++ {
		s := randText(r, r.Intn(20), "ab")
		want := LongestPrefixBrute(s)
		if a, b := LongestPrefixKMP(s), LongestPrefixHash(s); a != want || b != want {
			t.Fatalf("%q: brute=%q kmp=%q hash=%q", s, want, a, b)
		}
	}
}
