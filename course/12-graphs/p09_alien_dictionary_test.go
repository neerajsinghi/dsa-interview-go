package graphs

import (
	"math/rand"
	"sort"
	"testing"
)

// sortedUnder reports whether order is a valid answer for words: it uses each
// letter of the words exactly once, and words are non-decreasing under it.
func sortedUnder(words []string, order string) bool {
	rank := map[byte]int{}
	for i := 0; i < len(order); i++ {
		if _, dup := rank[order[i]]; dup {
			return false
		}
		rank[order[i]] = i
	}
	letters := map[byte]bool{}
	for _, w := range words {
		for i := 0; i < len(w); i++ {
			if _, ok := rank[w[i]]; !ok {
				return false
			}
			letters[w[i]] = true
		}
	}
	if len(letters) != len(order) {
		return false
	}
	less := func(a, b string) bool {
		for i := 0; i < len(a) && i < len(b); i++ {
			if a[i] != b[i] {
				return rank[a[i]] < rank[b[i]]
			}
		}
		return len(a) < len(b)
	}
	for i := 0; i+1 < len(words); i++ {
		if less(words[i+1], words[i]) {
			return false
		}
	}
	return true
}

func TestAlienOrder(t *testing.T) {
	tests := []struct {
		name  string
		words []string
		valid bool // whether some ordering exists
	}{
		{"classic", []string{"wrt", "wrf", "er", "ett", "rftt"}, true},
		{"two letters", []string{"z", "x"}, true},
		{"cycle", []string{"z", "x", "z"}, false},
		{"prefix after longer word", []string{"abc", "ab"}, false},
		{"single word", []string{"abc"}, true},
		{"equal words", []string{"ab", "ab"}, true},
		{"prefix first is fine", []string{"ab", "abc"}, true},
	}
	for _, tt := range tests {
		got := AlienOrder(tt.words)
		if tt.valid && !sortedUnder(tt.words, got) {
			t.Errorf("%s: AlienOrder = %q is not a valid alphabet", tt.name, got)
		}
		if !tt.valid && got != "" {
			t.Errorf("%s: AlienOrder = %q, want \"\"", tt.name, got)
		}
	}
	if got := AlienOrder([]string{"wrt", "wrf", "er", "ett", "rftt"}); got != "wertf" {
		t.Errorf("classic has a unique answer: got %q, want \"wertf\"", got)
	}
}

// TestAlienOrderRandom invents a secret alphabet, sorts random words by it,
// and checks that AlienOrder recovers SOME alphabet that sorts them too.
func TestAlienOrderRandom(t *testing.T) {
	r := rand.New(rand.NewSource(31))
	for iter := 0; iter < 300; iter++ {
		secret := r.Perm(5) // letters a..e in a random order
		rank := map[byte]int{}
		for i, c := range secret {
			rank[byte('a'+c)] = i
		}
		words := make([]string, 1+r.Intn(8))
		for i := range words {
			b := make([]byte, 1+r.Intn(4))
			for j := range b {
				b[j] = byte('a' + r.Intn(5))
			}
			words[i] = string(b)
		}
		sort.Slice(words, func(i, j int) bool {
			a, b := words[i], words[j]
			for k := 0; k < len(a) && k < len(b); k++ {
				if a[k] != b[k] {
					return rank[a[k]] < rank[b[k]]
				}
			}
			return len(a) < len(b)
		})
		if got := AlienOrder(words); !sortedUnder(words, got) {
			t.Fatalf("words %q (secret %v): got %q", words, secret, got)
		}
	}
}
