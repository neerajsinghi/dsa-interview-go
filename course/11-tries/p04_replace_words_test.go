package tries

import (
	"math/rand"
	"strings"
	"testing"
)

func TestReplaceWords(t *testing.T) {
	tests := []struct {
		name     string
		roots    []string
		sentence string
		want     string
	}{
		{"classic", []string{"cat", "bat", "rat"}, "the cattle was rattled by the battery", "the cat was rat by the bat"},
		{"shortest wins", []string{"a", "aa"}, "aadsfasf absbs bbab", "a a bbab"},
		{"no roots", nil, "keep these words", "keep these words"},
		{"root equals word", []string{"go"}, "go gopher", "go go"},
		{"empty sentence", []string{"a"}, "", ""},
		{"root longer than word", []string{"catalog"}, "cat", "cat"},
	}
	for _, tt := range tests {
		for name, fn := range map[string]func([]string, string) string{"trie": ReplaceWords, "brute": ReplaceWordsBrute} {
			if got := fn(tt.roots, tt.sentence); got != tt.want {
				t.Errorf("%s/%s = %q, want %q", tt.name, name, got, tt.want)
			}
		}
	}
}

func TestReplaceWordsRandom(t *testing.T) {
	r := rand.New(rand.NewSource(14))
	word := func(maxLen int) string {
		b := make([]byte, 1+r.Intn(maxLen))
		for i := range b {
			b[i] = byte('a' + r.Intn(3))
		}
		return string(b)
	}
	for iter := 0; iter < 300; iter++ {
		roots := make([]string, r.Intn(5))
		for i := range roots {
			roots[i] = word(3)
		}
		words := make([]string, r.Intn(6))
		for i := range words {
			words[i] = word(5)
		}
		s := strings.Join(words, " ")
		if got, want := ReplaceWords(roots, s), ReplaceWordsBrute(roots, s); got != want {
			t.Fatalf("roots %q sentence %q: trie=%q brute=%q", roots, s, got, want)
		}
	}
}
