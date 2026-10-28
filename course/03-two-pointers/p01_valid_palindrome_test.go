package twopointers

import (
	"math/rand"
	"testing"
)

func TestIsPalindrome(t *testing.T) {
	tests := []struct {
		s    string
		want bool
	}{
		{"", true},
		{" ", true},
		{".,!", true},
		{"a", true},
		{"ab", false},
		{"A man, a plan, a canal: Panama", true},
		{"race a car", false},
		{"0P", false}, // digits count and are not letters
		{"No 'x' in Nixon", true},
		{"12321", true},
	}
	for _, tc := range tests {
		if got := IsPalindromeBrute(tc.s); got != tc.want {
			t.Errorf("IsPalindromeBrute(%q) = %v, want %v", tc.s, got, tc.want)
		}
		if got := IsPalindrome(tc.s); got != tc.want {
			t.Errorf("IsPalindrome(%q) = %v, want %v", tc.s, got, tc.want)
		}
	}
}

func TestIsPalindromeRandomized(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	alphabet := []byte("aAb ,1")
	for range 2000 {
		b := make([]byte, rng.Intn(8))
		for i := range b {
			b[i] = alphabet[rng.Intn(len(alphabet))]
		}
		s := string(b)
		if x, y := IsPalindromeBrute(s), IsPalindrome(s); x != y {
			t.Fatalf("%q: brute=%v optimal=%v", s, x, y)
		}
	}
}
