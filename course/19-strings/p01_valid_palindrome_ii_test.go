package stringalgos

import (
	"math/rand"
	"testing"
)

func TestValidPalindromeII(t *testing.T) {
	tests := []struct {
		s    string
		want bool
	}{
		{"aba", true}, {"abca", true}, {"abc", false}, {"", true}, {"a", true},
		{"ab", true}, {"deeee", true}, {"eeccccbebaeeabebccceea", false},
		{"cupuufxoohdfpgjdmysgvhmvffcnqxjjxqncffvmhvgsymdjgpfdhooxfuupucu", true},
	}
	for _, tt := range tests {
		if got := ValidPalindromeII(tt.s); got != tt.want {
			t.Errorf("ValidPalindromeII(%q) = %v, want %v", tt.s, got, tt.want)
		}
		if got := ValidPalindromeIIBrute(tt.s); got != tt.want {
			t.Errorf("Brute(%q) = %v, want %v", tt.s, got, tt.want)
		}
	}
}

func TestValidPalindromeIIRandom(t *testing.T) {
	r := rand.New(rand.NewSource(1))
	for iter := 0; iter < 1000; iter++ {
		s := randText(r, r.Intn(9), "ab")
		if a, b := ValidPalindromeII(s), ValidPalindromeIIBrute(s); a != b {
			t.Fatalf("%q: fast=%v brute=%v", s, a, b)
		}
	}
}
