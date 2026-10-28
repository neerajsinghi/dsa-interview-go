package stringalgos

import (
	"math/rand"
	"slices"
	"testing"
)

func TestFindRepeatedDNA(t *testing.T) {
	tests := []struct {
		s    string
		want []string
	}{
		{"AAAAACCCCCAAAAACCCCCCAAAAAGGGTTT", []string{"AAAAACCCCC", "CCCCCAAAAA"}},
		{"AAAAAAAAAAAAA", []string{"AAAAAAAAAA"}},
		{"AAAAAAAAAA", nil},
		{"ACGT", nil},
		{"", nil},
	}
	for _, tt := range tests {
		if got := FindRepeatedDNA(tt.s); !slices.Equal(got, tt.want) {
			t.Errorf("FindRepeatedDNA(%q) = %v, want %v", tt.s, got, tt.want)
		}
		if got := FindRepeatedDNASubstrings(tt.s); !slices.Equal(got, tt.want) {
			t.Errorf("FindRepeatedDNASubstrings(%q) = %v, want %v", tt.s, got, tt.want)
		}
	}
}

func TestFindRepeatedDNARandom(t *testing.T) {
	r := rand.New(rand.NewSource(7))
	for iter := 0; iter < 200; iter++ {
		s := randText(r, r.Intn(200), "ACGT"[:1+r.Intn(4)])
		if a, b := FindRepeatedDNA(s), FindRepeatedDNASubstrings(s); !slices.Equal(a, b) {
			t.Fatalf("%q: rolling=%v map=%v", s, a, b)
		}
	}
}
