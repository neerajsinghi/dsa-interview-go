package slidingwindow

import (
	"math/rand"
	"testing"
)

func TestCheckInclusion(t *testing.T) {
	tests := []struct {
		s1, s2 string
		want   bool
	}{
		{"", "", true}, // the empty string is a permutation of itself
		{"", "abc", true},
		{"a", "", false},
		{"ab", "eidbaooo", true},
		{"ab", "eidboaoo", false},
		{"abc", "ab", false},
		{"adc", "dcda", true}, // match at the very end
		{"aab", "baa", true},
		{"aab", "abab", true},
		{"aa", "ab", false},
	}
	for _, tc := range tests {
		if got := CheckInclusionBrute(tc.s1, tc.s2); got != tc.want {
			t.Errorf("brute(%q, %q) = %v, want %v", tc.s1, tc.s2, got, tc.want)
		}
		if got := CheckInclusion(tc.s1, tc.s2); got != tc.want {
			t.Errorf("window(%q, %q) = %v, want %v", tc.s1, tc.s2, got, tc.want)
		}
	}
}

func TestCheckInclusionRandomized(t *testing.T) {
	rng := rand.New(rand.NewSource(6))
	word := func(n int) string {
		b := make([]byte, n)
		for i := range b {
			b[i] = byte('a' + rng.Intn(3))
		}
		return string(b)
	}
	for range 2000 {
		s1, s2 := word(rng.Intn(4)), word(rng.Intn(10))
		if x, y := CheckInclusionBrute(s1, s2), CheckInclusion(s1, s2); x != y {
			t.Fatalf("%q in %q: brute=%v window=%v", s1, s2, x, y)
		}
	}
}
