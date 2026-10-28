package stringalgos

import (
	"math/rand"
	"slices"
	"testing"
)

func TestPrefixFunction(t *testing.T) {
	tests := []struct {
		s    string
		want []int
	}{
		{"", []int{}},
		{"a", []int{0}},
		{"aaaa", []int{0, 1, 2, 3}},
		{"ababcababa", []int{0, 0, 1, 2, 0, 1, 2, 3, 4, 3}},
		{"abcd", []int{0, 0, 0, 0}},
	}
	for _, tt := range tests {
		if got := PrefixFunction(tt.s); !slices.Equal(got, tt.want) {
			t.Errorf("PrefixFunction(%q) = %v, want %v", tt.s, got, tt.want)
		}
	}
}

func TestZFunction(t *testing.T) {
	tests := []struct {
		s    string
		want []int
	}{
		{"", []int{}},
		{"aabxaab", []int{7, 1, 0, 0, 3, 1, 0}},
		{"aaaaa", []int{5, 4, 3, 2, 1}},
	}
	for _, tt := range tests {
		if got := ZFunction(tt.s); !slices.Equal(got, tt.want) {
			t.Errorf("ZFunction(%q) = %v, want %v", tt.s, got, tt.want)
		}
	}
}

func randText(r *rand.Rand, n int, alphabet string) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = alphabet[r.Intn(len(alphabet))]
	}
	return string(b)
}

// Brute-force definitions to check the linear algorithms against.
func TestPrefixAndZRandom(t *testing.T) {
	r := rand.New(rand.NewSource(19))
	for iter := 0; iter < 300; iter++ {
		s := randText(r, r.Intn(20), "ab")
		pi, z := PrefixFunction(s), ZFunction(s)
		for i := range s {
			want := 0
			for k := i; k > 0; k-- { // proper border of s[:i+1]
				if s[:k] == s[i+1-k:i+1] {
					want = k
					break
				}
			}
			if pi[i] != want {
				t.Fatalf("%q: pi[%d]=%d want %d", s, i, pi[i], want)
			}
			lcp := 0
			for i+lcp < len(s) && s[lcp] == s[i+lcp] {
				lcp++
			}
			if z[i] != lcp {
				t.Fatalf("%q: z[%d]=%d want %d", s, i, z[i], lcp)
			}
		}
	}
}

func TestRollingHash(t *testing.T) {
	s := "abcabcxabc"
	h := NewRollingHash(s)
	if h.Hash(0, 3) != h.Hash(3, 6) || h.Hash(0, 3) != h.Hash(7, 10) {
		t.Error("equal substrings must hash equally")
	}
	if h.Hash(0, 3) == h.Hash(1, 4) {
		t.Error("unexpected collision on tiny input")
	}
	if h.Hash(4, 4) != 0 {
		t.Error("empty substring should hash to 0")
	}
}
